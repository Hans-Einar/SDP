package documents

import (
	"context"
	"encoding/json"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/viewpoint"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRealFrontendCompositionProjection(t *testing.T) {
	v, input, e := viewpoint.Load("../examples/source-composition/System.design")
	if e != nil {
		t.Fatal(e)
	}
	old, e := os.ReadFile("../../../SDP/SDL/SDL/Frontend/System.design")
	if e != nil {
		t.Fatal(e)
	}
	baseline, e := viewpoint.New(string(old))
	if e != nil {
		t.Fatal(e)
	}
	if len(v.Model.Declarations) != len(baseline.Model.Declarations)+1 || len(v.Model.Statements) != len(baseline.Model.Statements)+1 {
		t.Fatal("pilot lost baseline facts")
	}
	facts := map[string]bool{}
	for _, s := range v.Model.Statements {
		facts[s.Sentence()] = true
	}
	for _, s := range baseline.Model.Statements {
		if !facts[s.Sentence()] {
			t.Fatal("missing baseline fact", s.Sentence())
		}
	}
	if v.System != "Frontend" || v.Revision != input.Snapshot.Revision() {
		t.Fatal("consumer identity disagrees")
	}
	b, e := Build(context.Background(), v, Options{Navigator: true, Project: "frontend"})
	if e != nil {
		t.Fatal(e)
	}
	var provenance struct {
		Revision string `json:"revision"`
		Sources  []struct {
			Path string `json:"path"`
		}
	}
	if e = json.Unmarshal(b.Files["sources.json"], &provenance); e != nil || provenance.Revision != v.Revision || len(provenance.Sources) != 3 {
		t.Fatal(provenance, e)
	}
	out := filepath.Join(t.TempDir(), "bundle")
	if e = b.Publish(out); e != nil {
		t.Fatal(e)
	}
	// Relocate a fact file: semantic identities survive, source revision/provenance changes.
	root := t.TempDir()
	for name, text := range input.Snapshot.Format() {
		p := filepath.Join(root, name)
		os.MkdirAll(filepath.Dir(p), 0700)
		os.WriteFile(p, []byte(text), 0600)
	}
	text, _ := os.ReadFile(filepath.Join(root, "System.design"))
	// An unlisted source never becomes part of the checked model.
	os.WriteFile(filepath.Join(root, "Unlisted.design"), []byte("invalid"), 0600)
	os.Rename(filepath.Join(root, "Features.design"), filepath.Join(root, "Moved.design"))
	os.WriteFile(filepath.Join(root, "System.design"), []byte(strings.ReplaceAll(string(text), "Features.design", "Moved.design")), 0600)
	relocated, _, e := viewpoint.Load(filepath.Join(root, "System.design"))
	if e != nil {
		t.Fatal(e)
	}
	if relocated.Revision == v.Revision {
		t.Fatal("source move did not change revision")
	}
	ids := map[string]bool{}
	for _, f := range v.Facts {
		ids[f.S("id")] = true
	}
	for _, f := range relocated.Facts {
		if !ids[f.S("id")] {
			t.Fatal("fact identity depends on source file", f)
		}
	}
}
