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

func TestDomainSchemaRejectsMalformedLineage(t *testing.T) {
	area := t.TempDir()
	p := initial(t, area, "A")
	original, e := readArtifact(p)
	must(t, e)
	cases := map[string]func(*Artifact){
		"counter":          func(a *Artifact) { a.Sequence = 1000000000 },
		"kind":             func(a *Artifact) { a.Ledger[0].Kind = "anything" },
		"id":               func(a *Artifact) { a.Ledger[0].ID = "not-an-id" },
		"cycle":            func(a *Artifact) { a.Ledger[0].Parents = []string{a.Head} },
		"missing parent":   func(a *Artifact) { a.Ledger[0].Parents = []string{uuid() + ":00000"} },
		"release ancestry": func(a *Artifact) { a.BaseRelease = uuid() },
		"conflict state":   func(a *Artifact) { a.PendingParents = []string{a.Head} },
		"author":           func(a *Artifact) { a.Ledger[0].Author = "" },
		"path":             func(a *Artifact) { a.Ledger[0].Inventory = map[string]string{"../escape": hash(nil)} },
	}
	for name, modify := range cases {
		t.Run(name, func(t *testing.T) {
			a := original
			a.Ledger = append([]Record(nil), original.Ledger...)
			modify(&a)
			a.MetadataDigest = metadataHash(a)
			if e := validate(a); e == nil {
				t.Fatal("malformed schema accepted")
			}
		})
	}
}
