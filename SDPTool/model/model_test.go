package model

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func must(t *testing.T, e error) {
	t.Helper()
	if e != nil {
		t.Fatal(e)
	}
}
func initial(t *testing.T, area, name string) string {
	t.Helper()
	r, e := CreateWork(area, name, "", true)
	must(t, e)
	return r.Path
}
func TestCreateStatusNoGit(t *testing.T) {
	area := t.TempDir()
	p := initial(t, area, "Example")
	r, e := Status(area, "work:Example")
	must(t, e)
	if r.Status != "clean" || r.Artifact.Kind != "work" {
		t.Fatal(r)
	}
	must(t, os.WriteFile(filepath.Join(p, "A.design"), []byte("broken intermediate"), 0600))
	r, e = Status(area, "work:Example")
	must(t, e)
	if r.Status != "dirty" {
		t.Fatal(r)
	}
	if _, e = CreateWork(area, "example", "", true); e == nil {
		t.Fatal("duplicate accepted")
	}
	must(t, os.RemoveAll(filepath.Join(p, ".merge")))
	_, e = Status(area, "work:Example")
	must(t, e)
}
func TestStrictAndPaths(t *testing.T) {
	for _, s := range []string{"schema: x\nschema: y\n", "id: &a foo\nname: *a\n", "schema: x\n---\nname: y\n", "unknown: true\n", "id: !custom value\n"} {
		var a Artifact
		if strict([]byte(s), &a) == nil {
			t.Fatalf("accepted %s", s)
		}
	}
	for _, s := range []string{"../x", "a/../../b", "/etc/passwd", "a\\b", "CON", "a:b"} {
		if safePath(s) {
			t.Fatal(s)
		}
	}
	if !safePath(".commits/#00000/files/Main.design") {
		t.Fatal("history path rejected")
	}
}
func TestSymlinkRejected(t *testing.T) {
	area := t.TempDir()
	p := initial(t, area, "A")
	must(t, os.Symlink("/etc/passwd", filepath.Join(p, "evil")))
	if _, e := Status(area, "work:A"); e == nil {
		t.Fatal("symlink accepted")
	}
}
func TestArtifactTamper(t *testing.T) {
	area := t.TempDir()
	p := initial(t, area, "A")
	meta := filepath.Join(p, metadataName(p))
	b, e := os.ReadFile(meta)
	must(t, e)
	must(t, os.WriteFile(meta, []byte(strings.Replace(string(b), "name: A", "name: B", 1)), 0600))
	if _, e = Status(area, "work:B"); e == nil {
		t.Fatal("metadata tamper accepted")
	}
}
