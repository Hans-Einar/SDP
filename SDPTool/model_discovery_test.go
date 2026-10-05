package sdptool

import (
	"github.com/Hans-Einar/SDP/SDPTool/model"
	"os"
	"path/filepath"
	"testing"
)

func TestArtifactDiscoveryPrunesOnlyOwnedHistory(t *testing.T) {
	root := t.TempDir()
	area := filepath.Join(root, "SDP", "SDL")
	if e := os.MkdirAll(area, 0700); e != nil {
		t.Fatal(e)
	}
	r, e := model.CreateWork(area, "Demo", "", true)
	if e != nil {
		t.Fatal(e)
	}
	put := func(path string) {
		t.Helper()
		if e := os.MkdirAll(filepath.Dir(path), 0700); e != nil {
			t.Fatal(e)
		}
		if e := os.WriteFile(path, []byte("language design-core version 0.5.\nunit Demo.\n"), 0600); e != nil {
			t.Fatal(e)
		}
	}
	put(filepath.Join(r.Path, "Main.design"))
	_, e = model.Commit(area, "work:Demo", "save", false)
	if e != nil {
		t.Fatal(e)
	}
	put(filepath.Join(area, "ordinary", ".commits", "Also.design"))
	p, e := Discover(root)
	if e != nil {
		t.Fatal(e)
	}
	if len(p.Sources) != 2 {
		t.Fatalf("history leaked or ordinary source hidden: %+v", p.Sources)
	}
	owned := 0
	for _, s := range p.Sources {
		if s.ArtifactKind == "work" {
			owned++
			if !s.Preliminary || s.ArtifactID == "" {
				t.Fatal(s)
			}
		}
	}
	if owned != 1 {
		t.Fatal(p.Sources)
	}
	view, e := model.Snapshot(area, "work:Demo")
	if e != nil {
		t.Fatal(e)
	}
	if !view.Preliminary || len(view.Files) != 1 {
		t.Fatal(view)
	}
	before := view.Artifact.Head
	put(filepath.Join(r.Path, "Other.design"))
	after, e := model.Snapshot(area, "work:Demo")
	if e != nil {
		t.Fatal(e)
	}
	if before != after.Artifact.Head || after.LiveDigest == view.LiveDigest || len(view.Files) != 1 {
		t.Fatal("snapshot mutated or committed")
	}
}
