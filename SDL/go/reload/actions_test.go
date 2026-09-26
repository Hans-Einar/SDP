package reload

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestSDLWatchValidInvalidAndAtomicReplacement(t *testing.T) {
	source, e := os.ReadFile("../examples/echo.sdl")
	if e != nil {
		t.Fatal(e)
	}
	path := filepath.Join(t.TempDir(), "actions.sdl")
	os.WriteFile(path, source, 0600)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := Watch(ctx, path, 10*time.Millisecond)
	next := func() Candidate {
		select {
		case c := <-ch:
			return c
		case <-time.After(3 * time.Second):
			t.Fatal("SDL watcher timed out")
		}
		return Candidate{}
	}
	first := next()
	if first.Err != nil || first.Program.Actions["Echo"].GoSymbol != "GoEcho" {
		t.Fatal(first)
	}
	os.WriteFile(path, []byte("invalid"), 0600)
	bad := next()
	if bad.Err == nil || bad.Sequence <= first.Sequence {
		t.Fatal("missing syntax diagnostic")
	}
	tmp := path + ".new"
	os.WriteFile(tmp, []byte(strings.Replace(string(source), "GoEcho", "GoOther", 1)), 0600)
	os.Rename(tmp, path)
	good := next()
	if good.Err != nil || good.Program.Actions["Echo"].GoSymbol != "GoOther" || good.Sequence <= bad.Sequence {
		t.Fatal(good)
	}
}
