package governancehost

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"
)

func TestPrivateSocketReconnectIdentityAndFrameBoundary(t *testing.T) {
	dir, err := os.MkdirTemp("", "pgi4-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	path := filepath.Join(dir, "host.sock")
	listen := func(version string) (net.Listener, *http.Server) {
		l, err := net.Listen("unix", path)
		if err != nil {
			t.Fatal(err)
		}
		if err = os.Chmod(path, 0600); err != nil {
			t.Fatal(err)
		}
		server := &http.Server{Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ws, err := websocket.Accept(w, r, nil)
			if err != nil {
				return
			}
			defer ws.CloseNow()
			for {
				_, b, err := ws.Read(r.Context())
				if err != nil {
					return
				}
				var f frame
				if json.Unmarshal(b, &f) != nil {
					return
				}
				if f.Method == "initialize" {
					// Pretty printing stays one JSON message on the adapter's JSONL boundary.
					b, _ = json.MarshalIndent(map[string]any{"id": f.ID, "result": map[string]string{"userAgent": "sdp_governance/" + version + " fixture", "platformOs": "linux"}}, "", "  ")
					if ws.Write(r.Context(), websocket.MessageText, b) != nil {
						return
					}
				}
			}
		})}
		go server.Serve(l)
		return l, server
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	l, server := listen(HostVersion)
	defer server.Close()
	c, err := OpenSocket(ctx, path)
	if err != nil {
		t.Fatal(err)
	}
	first := c.HostIdentity
	c.Close()
	c, err = OpenSocket(ctx, path)
	if err != nil {
		t.Fatal(err)
	}
	if first == "" || c.HostIdentity != first {
		t.Fatal("identity changed for same server")
	}
	c.Close()
	if err = os.Chmod(dir, 0755); err != nil {
		t.Fatal(err)
	}
	if c, err = OpenSocket(ctx, path); err == nil {
		c.Close()
		t.Fatal("public parent accepted")
	}
	os.Chmod(dir, 0700)
	link := filepath.Join(dir, "alias.sock")
	if err = os.Symlink(path, link); err != nil {
		t.Fatal(err)
	}
	if c, err = OpenSocket(ctx, link); err != nil {
		t.Fatal(err)
	}
	c.Close()
	nested := filepath.Join(dir, "nested.sock")
	if err = os.Symlink(link, nested); err != nil {
		t.Fatal(err)
	}
	if c, err = OpenSocket(ctx, nested); err == nil {
		c.Close()
		t.Fatal("indirect alias accepted")
	}
	l.Close()
	server.Close()
	_, server2 := listen(HostVersion)
	defer server2.Close()
	c, err = OpenSocket(ctx, path)
	if err != nil {
		t.Fatal(err)
	}
	if c.HostIdentity == first {
		t.Fatal("replacement socket retained identity")
	}
	c.Close()
	server2.Close()
	_, server3 := listen("0.0.0")
	defer server3.Close()
	if c, err = OpenSocket(ctx, path); err == nil {
		c.Close()
		t.Fatal("wrong host version accepted")
	} else if !strings.Contains(err.Error(), "version mismatch") {
		t.Fatal(err)
	}
}
