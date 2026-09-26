//go:build linux

package broker

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
)

type Call struct {
	Operation string
	Request   Request
	Lease     string
}
type Reply struct {
	Result *Result `json:"result,omitempty"`
	Error  string  `json:"error,omitempty"`
}

func RuntimeDirectory() (string, error) {
	root := os.Getenv("XDG_RUNTIME_DIR")
	if root == "" {
		root = filepath.Join(os.TempDir(), fmt.Sprintf("sdl-runtime-%d", os.Getuid()))
	}
	if e := private(root); e != nil {
		return "", e
	}
	root = filepath.Join(root, "sdl")
	return root, private(root)
}
func private(p string) error {
	if !filepath.IsAbs(p) {
		return fmt.Errorf("absolute runtime path required")
	}
	if e := os.MkdirAll(p, 0700); e != nil {
		return e
	}
	i, e := os.Lstat(p)
	if e != nil || !i.IsDir() || i.Mode().Perm() != 0700 {
		return fmt.Errorf("private runtime directory required")
	}
	s, ok := i.Sys().(*syscall.Stat_t)
	if !ok || s.Uid != uint32(os.Getuid()) {
		return fmt.Errorf("runtime directory owner mismatch")
	}
	return nil
}
func (b *Broker) Serve(ctx context.Context, path string) error {
	if e := private(filepath.Dir(path)); e != nil {
		return e
	}
	if i, err := os.Lstat(path); err == nil {
		if i.Mode()&os.ModeSocket == 0 {
			return fmt.Errorf("endpoint exists and is not a socket")
		}
		c, err := net.DialTimeout("unixpacket", path, 200*time.Millisecond)
		if err == nil {
			c.Close()
			return fmt.Errorf("broker already running")
		}
		if !errors.Is(err, syscall.ECONNREFUSED) {
			return err
		}
		if err = os.Remove(path); err != nil {
			return err
		}
	}
	l, e := net.ListenUnix("unixpacket", &net.UnixAddr{Name: path, Net: "unixpacket"})
	if e != nil {
		return e
	}
	defer l.Close()
	if e = os.Chmod(path, 0600); e != nil {
		return e
	}
	b.Endpoint = path
	done := make(chan struct{})
	defer close(done)
	go func() {
		select {
		case <-ctx.Done():
			l.Close()
		case <-done:
		}
	}()
	slots := make(chan struct{}, 16)
	var jobs sync.WaitGroup
	defer jobs.Wait()
	for {
		c, e := l.AcceptUnix()
		if e != nil {
			if ctx.Err() != nil {
				return nil
			}
			return e
		}
		select {
		case slots <- struct{}{}:
			jobs.Go(func() { defer func() { <-slots; c.Close() }(); b.connection(ctx, c) })
		default:
			c.Close()
		}
	}
}
func (b *Broker) connection(ctx context.Context, c *net.UnixConn) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()
	_ = c.SetDeadline(time.Now().Add(3 * time.Minute))
	raw, e := c.SyscallConn()
	if e != nil {
		return
	}
	allowed := false
	_ = raw.Control(func(fd uintptr) {
		cred, err := syscall.GetsockoptUcred(int(fd), syscall.SOL_SOCKET, syscall.SO_PEERCRED)
		allowed = err == nil && cred.Uid == uint32(os.Getuid())
	})
	if !allowed {
		return
	}
	buf := make([]byte, 32769)
	n, e := c.Read(buf)
	if e != nil {
		return
	}
	reply := Reply{}
	if n > 32768 {
		reply.Error = "IPC packet limit"
	} else {
		packet := buf[:n]
		var call Call
		if strings.HasPrefix(string(packet), "SDLVIEW1\tRELEASE\t") {
			call.Operation = "release"
			call.Lease = strings.TrimSuffix(strings.TrimPrefix(string(packet), "SDLVIEW1\tRELEASE\t"), "\n")
		} else {
			d := json.NewDecoder(bytes.NewReader(packet))
			d.DisallowUnknownFields()
			e = d.Decode(&call)
			if e == nil {
				var extra any
				if d.Decode(&extra) != io.EOF {
					e = fmt.Errorf("trailing JSON")
				}
			}
		}
		if e != nil {
			reply.Error = e.Error()
		} else {
			switch call.Operation {
			case "select":
				r, err := b.Select(ctx, call.Request)
				if err != nil {
					if r.Lease != "" {
						reply.Result = &r
					}
					reply.Error = err.Error()
				} else {
					reply.Result = &r
				}
			case "release":
				e = b.Store.Release(call.Lease)
			case "sweep":
				e = b.Store.Sweep()
			default:
				e = fmt.Errorf("unknown IPC operation")
			}
			if e != nil {
				reply.Error = e.Error()
			}
		}
	}
	data, _ := json.Marshal(reply)
	_, _ = c.Write(data)
}
func Send(ctx context.Context, path string, call Call) (Reply, error) {
	var r Reply
	d := net.Dialer{}
	conn, e := d.DialContext(ctx, "unixpacket", path)
	if e != nil {
		return r, e
	}
	defer conn.Close()
	deadline := time.Now().Add(3 * time.Minute)
	if t, ok := ctx.Deadline(); ok && t.Before(deadline) {
		deadline = t
	}
	conn.SetDeadline(deadline)
	packet, _ := json.Marshal(call)
	if len(packet) > 32768 {
		return r, fmt.Errorf("IPC request limit")
	}
	if _, e = conn.Write(packet); e != nil {
		return r, e
	}
	buf := make([]byte, 32769)
	n, e := conn.Read(buf)
	if e != nil {
		return r, e
	}
	if n > 32768 {
		return r, fmt.Errorf("IPC response limit")
	}
	e = json.Unmarshal(buf[:n], &r)
	if e == nil && r.Error != "" {
		e = fmt.Errorf("%s", r.Error)
	}
	return r, e
}
