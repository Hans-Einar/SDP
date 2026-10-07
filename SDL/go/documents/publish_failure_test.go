package documents

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPublicationRollback(t *testing.T) {
	for _, failAt := range []int{1, 2} {
		t.Run(fmt.Sprint(failAt), func(t *testing.T) {
			out := filepath.Join(t.TempDir(), "out")
			old := &Bundle{Files: map[string][]byte{"index.md": []byte("old")}}
			old.Seal()
			if e := old.Publish(out); e != nil {
				t.Fatal(e)
			}
			next := &Bundle{Files: map[string][]byte{"index.md": []byte("new")}}
			next.Seal()
			calls := 0
			e := next.publish(out, func(a, b string) error {
				calls++
				if calls == failAt {
					return fmt.Errorf("injected rename failure")
				}
				return os.Rename(a, b)
			})
			if e == nil {
				t.Fatal("failure ignored")
			}
			data, e := os.ReadFile(filepath.Join(out, "index.md"))
			if e != nil || string(data) != "old" {
				t.Fatal("previous output lost", e)
			}
		})
	}
}
func TestPublicationWriteFailure(t *testing.T) {
	out := filepath.Join(t.TempDir(), "out")
	old := &Bundle{Files: map[string][]byte{"index.md": []byte("old")}}
	old.Seal()
	if e := old.Publish(out); e != nil {
		t.Fatal(e)
	}
	bad := &Bundle{Files: map[string][]byte{"../escape": []byte("invalid")}}
	bad.Seal()
	if e := bad.Publish(out); e == nil {
		t.Fatal("unsafe write accepted")
	}
	data, e := os.ReadFile(filepath.Join(out, "index.md"))
	if e != nil || string(data) != "old" {
		t.Fatal("prior output lost")
	}
}
func TestPublicationFailedRecoveryReportsBackup(t *testing.T) {
	out := filepath.Join(t.TempDir(), "out")
	old := &Bundle{Files: map[string][]byte{"index.md": []byte("old")}}
	old.Seal()
	if e := old.Publish(out); e != nil {
		t.Fatal(e)
	}
	calls := 0
	e := old.publish(out, func(a, b string) error {
		calls++
		if calls >= 2 {
			return fmt.Errorf("injected")
		}
		return os.Rename(a, b)
	})
	if e == nil || !strings.Contains(e.Error(), "prior output retained at") {
		t.Fatal(e)
	}
	backups, _ := filepath.Glob(filepath.Join(filepath.Dir(out), ".sdl-publish-*-previous", "index.md"))
	if len(backups) != 1 {
		t.Fatal(backups)
	}
	data, e := os.ReadFile(backups[0])
	if e != nil || string(data) != "old" {
		t.Fatal("backup lost", e)
	}
}
