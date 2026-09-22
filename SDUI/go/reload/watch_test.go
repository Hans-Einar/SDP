package reload

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	uiruntime "github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

func receive(t *testing.T, ch <-chan Candidate) Candidate {
	t.Helper()
	select {
	case c, ok := <-ch:
		if !ok {
			t.Fatal("watch closed")
		}
		return c
	case <-time.After(3 * time.Second):
		t.Fatal("watch timeout")
	}
	return Candidate{}
}
func TestAtomicSaveErrorAndLastGood(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "page.sdui")
	first := `sdui 0.2; page=[go=button("One")];`
	if err := os.WriteFile(path, []byte(first), 0600); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := Watch(ctx, path, "page", 10*time.Millisecond, nil)
	c := receive(t, ch)
	if c.Err != nil {
		t.Fatal(c.Err)
	}
	s, err := uiruntime.New("test", c.Root)
	if err != nil {
		t.Fatal(err)
	}
	controller := Controller{Session: s, Sequence: c.Sequence, SourceHash: c.Hash}
	w, _ := s.Widget("page/go")
	os.WriteFile(path, []byte(`sdui 0.2; page=[`), 0600)
	bad := receive(t, ch)
	if controller.Adopt(bad) == nil {
		t.Fatal("invalid accepted")
	}
	after, _ := s.Widget("page/go")
	if after != w || s.Revision != 1 {
		t.Fatal("invalid save changed UI")
	}
	replacement := filepath.Join(dir, "replacement")
	os.WriteFile(replacement, []byte(`sdui 0.2; page=[go=button("Two")];`), 0600)
	os.Rename(replacement, path)
	good := receive(t, ch)
	if e := controller.Adopt(good); e != nil {
		t.Fatal(e)
	}
	after, _ = s.Widget("page/go")
	if after.Label != "Two" || s.Revision != 2 || controller.LastError != nil {
		t.Fatal(after, controller.LastError)
	}
	if controller.Adopt(bad) == nil {
		t.Fatal("late result accepted")
	}
	cancel()
	select {
	case _, ok := <-ch:
		if ok {
			t.Fatal("pending unexpected candidate")
		}
	case <-time.After(time.Second):
		t.Fatal("watch did not stop")
	}
}
