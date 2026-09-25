package sdptool

import (
	"encoding/json"
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
