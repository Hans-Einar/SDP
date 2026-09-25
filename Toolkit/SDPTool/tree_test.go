package sdptool

import (
	"context"
	"encoding/json"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/viewpoint"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTreeCatalogCollectionsAndStableIdentity(t *testing.T) {
	root, r := projectFixture(t)
	r.Models = []Model{{"model", "System", "model.design", "design-core/0.5"}}
	saveRegistration(t, root, r)
	path := filepath.Join(root, "model.design")
	os.WriteFile(path, []byte(sample), 0600)
	p, _ := Discover(root)
	tree, e := ModelTree(p, "model")
	if e != nil {
		t.Fatal(e)
	}
	vp, count := 0, 0
	ids := map[string]bool{}
	for _, n := range tree.Nodes {
		ids[n.ID] = true
		if n.Kind == "viewpoint" {
			vp++
		}
		if n.Kind == "collection" && n.Label == "container" {
			count++
		}
	}
	if vp != 11 || count == 0 {
		t.Fatalf("vp=%d container=%d", vp, count)
	}
	for _, n := range tree.Nodes {
		for _, c := range n.Children {
			if !ids[c] {
				t.Fatal("unresolved child", c)
			}
		}
		if n.Reference != "" && !ids[n.Reference] {
			t.Fatal("unresolved reference", n.Reference)
		}
	}
	again, _ := ModelTree(p, "model")
	a, _ := json.Marshal(tree)
	b, _ := json.Marshal(again)
	if string(a) != string(b) {
		t.Fatal("unstable tree")
	}
	os.WriteFile(path, []byte(strings.Replace(sample, "version 0.5.\n", "version 0.5.\nactivity Added.\n", 1)), 0600)
	changed, e := ModelTree(p, "model")
	if e != nil {
		t.Fatal(e)
	}
	if tree.Revision == changed.Revision {
		t.Fatal("stale revision")
	}
	for _, n := range changed.Nodes {
		if n.Kind == "container" && !ids[n.ID] {
			t.Fatal("identity changed with unrelated property")
		}
	}
}

func TestRelationshipReferencesAndSelectionRevision(t *testing.T) {
	root, r := projectFixture(t)
	r.Models = []Model{{"model", "System", "model.design", "design-core/0.5"}}
	saveRegistration(t, root, r)
	path := filepath.Join(root, "model.design")
	// Two opposite actor/use-case relationships are not valid SDL. Use a valid
	// collaboration model, then inspect graph representation with cyclic facts
	// separately: renderer grouping must never recurse into references.
	os.WriteFile(path, []byte(sample), 0600)
	p, _ := Discover(root)
	tree, e := ModelTree(p, "model")
	if e != nil {
		t.Fatal(e)
	}
	relationships := 0
	for _, n := range tree.Nodes {
		if n.Kind == "relationship" {
			relationships++
			if n.Reference == "" || len(n.Children) != 0 {
				t.Fatal("relationship cloned a subtree")
			}
		}
	}
	if relationships == 0 {
		t.Fatal("no relations")
	}
	o := PreviewOptions{URI: selectionURI("trial", "VP02", "", "VP02-roots"), Revision: tree.Revision, Output: filepath.Join(t.TempDir(), "selected")}
	if _, e = SelectProject(context.Background(), p, "model", o); e != nil {
		t.Fatal(e)
	}
	o.URI = selectionURI("another", "VP02", "", "")
	if _, e = SelectProject(context.Background(), p, "model", o); e == nil {
		t.Fatal("foreign project accepted")
	}
	o.URI = selectionURI("trial", "VP02", "", "")
	os.WriteFile(path, []byte(strings.Replace(sample, "Child", "Edited", -1)), 0600)
	if _, e = SelectProject(context.Background(), p, "model", o); e == nil || !strings.Contains(e.Error(), "stale") {
		t.Fatalf("%v", e)
	}
}

func TestCyclicReferenceRepresentationIsFinite(t *testing.T) {
	v, e := viewpoint.New(sample)
	if e != nil {
		t.Fatal(e)
	}
	// Projection-boundary fixture: navigation must safely represent a graph even
	// when a future legal relation profile admits cycles; no syntax adoption claim.
	v.Facts = append(v.Facts, viewpoint.Fact{"subject": "Child", "object": "Main", "verb": "contains"})
	tree, e := modelTree(Project{Registration: Registration{ProjectID: "trial"}}, Model{ID: "model"}, v)
	if e != nil {
		t.Fatal(e)
	}
	if len(tree.Nodes) > 200 || tree.ExpansionDepthLimit != 8 {
		t.Fatal("unbounded expansion")
	}
	refs := 0
	for _, n := range tree.Nodes {
		if n.Kind == "relationship" {
			refs++
		}
	}
	if refs != 2 {
		t.Fatal(refs)
	}
}
