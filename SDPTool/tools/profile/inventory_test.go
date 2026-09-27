package main

import (
	"encoding/json"
	"os"
	"path/filepath"
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
