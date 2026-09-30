package main

import (
	"context"
	"encoding/json"
	"github.com/Hans-Einar/SDP/SDPTool/install"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Exercise the real authoring inventory: a guide added to the canonical template
// must not silently disappear from the generated installation payload.
func TestCurrentTemplateInventoryComplete(t *testing.T) {
	repo := "../../.."
	b, err := os.ReadFile(filepath.Join(repo, "SDPTool/profiles/payload.json"))
	if err != nil {
		t.Fatal(err)
	}
	var inventory struct {
		Files []struct{ Source, Destination, Ownership string }
	}
	if err := json.Unmarshal(b, &inventory); err != nil {
		t.Fatal(err)
	}
	listed := map[string]bool{}
	for _, f := range inventory.Files {
		listed[f.Source] = true
	}
	err = filepath.WalkDir(filepath.Join(repo, "Template/sdp-root"), func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		relative, err := filepath.Rel(repo, path)
		if err != nil {
			return err
		}
		if !listed[filepath.ToSlash(relative)] {
			t.Errorf("template missing from current inventory: %s", relative)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestCurrentDescriptorReproducible(t *testing.T) {
	dir := t.TempDir()
	first, second := filepath.Join(dir, "first.json"), filepath.Join(dir, "second.json")
	for _, out := range []string{first, second} {
		if err := build("../../..", out, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "development-reproducible", ""); err != nil {
			t.Fatal(err)
		}
	}
	a, err := os.ReadFile(first)
	if err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(second)
	if err != nil {
		t.Fatal(err)
	}
	if string(a) != string(b) {
		t.Fatal("descriptor build is nondeterministic")
	}
}

func TestCurrentSessionDistributionAndPreservation(t *testing.T) {
	t.Setenv("SDP_CACHE_DIR", t.TempDir())
	artifact := filepath.Join(t.TempDir(), "release.json")
	if e := build("../../..", artifact, "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "development-sessions", ""); e != nil {
		t.Fatal(e)
	}
	b, _ := os.ReadFile(artifact)
	var d install.Descriptor
	if e := json.Unmarshal(b, &d); e != nil {
		t.Fatal(e)
	}
	sessions := map[string]bool{}
	for _, f := range d.Files {
		if strings.HasPrefix(f.Path, "SDP/Sessions/") {
			if f.Ownership != "initialize-if-missing" {
				t.Fatal("Session overwrite policy", f.Path)
			}
			sessions[f.Path] = true
		}
	}
	if len(sessions) != 2 || !sessions["SDP/Sessions/README.md"] || !sessions["SDP/Sessions/Session-template.md"] {
		t.Fatal("project conversation leak or missing templates", sessions)
	}
	root := t.TempDir()
	p, e := install.Preview(install.Options{Root: root, Operation: "install", Artifact: artifact, AllowUnreleased: true})
	if e != nil || !p.CanApply {
		t.Fatal(p.Conflicts, e)
	}
	if _, e = (install.Executor{}).Apply(context.Background(), root, p); e != nil {
		t.Fatal(e)
	}
	guide := filepath.Join(root, "SDP/Sessions/README.md")
	os.WriteFile(guide, []byte("project-owned guide"), 0600)
	doc := filepath.Join(root, "SDP/Sessions/session-#0001--Local.md")
	os.WriteFile(doc, []byte("owner work"), 0600)
	p, e = install.Preview(install.Options{Root: root, Operation: "upgrade", Artifact: artifact, AllowUnreleased: true})
	if e != nil || !p.NoChange {
		t.Fatal(p.Conflicts, e)
	}
	if _, e = (install.Executor{}).Apply(context.Background(), root, p); e != nil {
		t.Fatal(e)
	}
	b, _ = os.ReadFile(guide)
	if string(b) != "project-owned guide" {
		t.Fatal("guide overwritten")
	}
	b, _ = os.ReadFile(doc)
	if string(b) != "owner work" {
		t.Fatal("work overwritten")
	}
}
