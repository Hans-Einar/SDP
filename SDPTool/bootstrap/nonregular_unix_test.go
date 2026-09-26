//go:build linux || darwin || freebsd

package bootstrap

import (
	"context"
	"path/filepath"
	"syscall"
	"testing"
	"time"
)

func TestFIFORejectedBeforeOpen(t *testing.T) {
	p := filepath.Join(t.TempDir(), "fifo")
	if e := syscall.Mkfifo(p, 0600); e != nil {
		t.Fatal(e)
	}
	done := make(chan error, 1)
	go func() {
		_, e := Resolve(context.Background(), Config{Descriptor: p, CacheDir: filepath.Join(t.TempDir(), "cache")})
		done <- e
	}()
	select {
	case e := <-done:
		if e == nil {
			t.Fatal("FIFO accepted")
		}
	case <-time.After(time.Second):
		t.Fatal("FIFO read blocked")
	}
}
