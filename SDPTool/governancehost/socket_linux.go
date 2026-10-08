package governancehost

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
	"github.com/coder/websocket"
)

// OpenSocket connects to an explicitly selected private local server. It never
// starts, stops or restarts the host. Its identity binds socket and peer lifetime.
func OpenSocket(ctx context.Context, path string) (*Client, error) {
	if !filepath.IsAbs(path) || filepath.Clean(path) != path {
		return nil, fmt.Errorf("absolute clean host socket path required")
	}
	canonical, err := filepath.EvalSymlinks(path)
	if err != nil {
		return nil, fmt.Errorf("host socket unavailable: %w", err)
	}
	requestedDir, err := filepath.EvalSymlinks(filepath.Dir(path))
	if err != nil || requestedDir != filepath.Dir(path) {
		return nil, fmt.Errorf("host socket parent must not traverse symlinks")
	}
	alias, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	if canonical != path {
		// Codex 0.161.0 creates an absolute alias to a private per-user socket.
		target, err := os.Readlink(path)
		if err != nil || target != canonical {
			return nil, fmt.Errorf("host socket alias must point directly to canonical socket")
		}
		stat, ok := alias.Sys().(*syscall.Stat_t)
		if !ok || stat.Uid != uint32(os.Getuid()) {
			return nil, fmt.Errorf("host socket alias owner mismatch")
		}
	}
	privateDir := func(path string) error {
		info, err := os.Stat(path)
		if err != nil {
			return err
		}
		stat, ok := info.Sys().(*syscall.Stat_t)
		if !ok || !info.IsDir() || info.Mode().Perm()&0077 != 0 || stat.Uid != uint32(os.Getuid()) {
			return fmt.Errorf("host socket parents must be private and owned by current user")
		}
		return nil
	}
	if err = privateDir(requestedDir); err != nil {
		return nil, err
	}
	if err = privateDir(filepath.Dir(canonical)); err != nil {
		return nil, err
	}
	dir, err := os.Stat(filepath.Dir(canonical))
	if err != nil {
		return nil, err
	}
	info, err := os.Lstat(canonical)
	if err != nil {
		return nil, err
	}
	ds, dok := dir.Sys().(*syscall.Stat_t)
	ss, sok := info.Sys().(*syscall.Stat_t)
	if !dok || !sok || !dir.IsDir() || dir.Mode().Perm()&0077 != 0 || info.Mode()&os.ModeSocket == 0 || info.Mode().Perm()&0077 != 0 || ds.Uid != uint32(os.Getuid()) || ss.Uid != uint32(os.Getuid()) {
		return nil, fmt.Errorf("host socket and parent must be private and owned by current user")
	}
	identity := ""
	transport := &http.Transport{DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
		conn, err := (&net.Dialer{}).DialContext(ctx, "unix", canonical)
		if err != nil {
			return nil, err
		}
		fail := func(err error) (net.Conn, error) { conn.Close(); return nil, err }
		raw, err := conn.(*net.UnixConn).SyscallConn()
		if err != nil {
			return fail(err)
		}
		var peer *syscall.Ucred
		var peerErr error
		err = raw.Control(func(fd uintptr) {
			peer, peerErr = syscall.GetsockoptUcred(int(fd), syscall.SOL_SOCKET, syscall.SO_PEERCRED)
		})
		if err != nil {
			return fail(err)
		}
		if peerErr != nil {
			return fail(peerErr)
		}
		if peer.Uid != uint32(os.Getuid()) {
			return fail(fmt.Errorf("host peer user mismatch"))
		}
		stat, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", peer.Pid))
		if err != nil {
			return fail(err)
		}
		fields := strings.Fields(string(stat)[strings.LastIndex(string(stat), ")")+1:])
		if len(fields) < 20 {
			return fail(fmt.Errorf("host process identity unavailable"))
		}
		boot, err := os.ReadFile("/proc/sys/kernel/random/boot_id")
		if err != nil {
			return fail(err)
		}
		after, err := os.Lstat(canonical)
		if err != nil || !os.SameFile(info, after) {
			return fail(fmt.Errorf("host socket changed during connect"))
		}
		resolved, resolveErr := filepath.EvalSymlinks(path)
		if resolveErr != nil || resolved != canonical {
			return fail(fmt.Errorf("host socket alias changed during connect"))
		}
		sum := sha256.Sum256([]byte(fmt.Sprintf("%s/%d/%d/%d/%d/%d/%d/%s/%s", path, ss.Dev, ss.Ino, ss.Ctim.Sec, ss.Ctim.Nsec, peer.Uid, peer.Pid, fields[19], boot)))
		next := fmt.Sprintf("%x", sum)
		if identity != "" && identity != next {
			return fail(fmt.Errorf("host peer changed during connect"))
		}
		identity = next
		return conn, nil
	}}
	defer transport.CloseIdleConnections()
	ws, _, err := websocket.Dial(ctx, "ws://localhost/", &websocket.DialOptions{HTTPClient: &http.Client{Transport: transport, CheckRedirect: func(*http.Request, []*http.Request) error { return fmt.Errorf("host redirect refused") }}})
	if err != nil {
		return nil, err
	}
	ws.SetReadLimit(governance.MaxRequest)
	wire := &socketWire{ctx: ctx, conn: ws}
	c := newClient(wire, wire)
	c.HostIdentity = identity
	return c.initialize(ctx)
}

// Preserve one JSON-RPC frame per WebSocket message and adapt only at the existing
// bounded JSONL reader boundary. Pretty-printed JSON is compacted, never split.
type socketWire struct {
	ctx     context.Context
	conn    *websocket.Conn
	pending []byte
}

func (w *socketWire) Read(p []byte) (int, error) {
	if len(w.pending) == 0 {
		kind, b, err := w.conn.Read(w.ctx)
		if err != nil {
			return 0, err
		}
		if kind != websocket.MessageText {
			return 0, fmt.Errorf("host binary frame refused")
		}
		if err = governance.CheckWire(b); err != nil {
			return 0, err
		}
		var compact bytes.Buffer
		if err = json.Compact(&compact, b); err != nil {
			return 0, err
		}
		w.pending = append(compact.Bytes(), '\n')
	}
	n := copy(p, w.pending)
	w.pending = w.pending[n:]
	return n, nil
}
func (w *socketWire) Write(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	if err := w.conn.Write(w.ctx, websocket.MessageText, bytes.TrimSuffix(p, []byte{'\n'})); err != nil {
		return 0, err
	}
	return len(p), nil
}
func (w *socketWire) Close() error { return w.conn.CloseNow() }

var _ io.ReadWriteCloser = (*socketWire)(nil)
