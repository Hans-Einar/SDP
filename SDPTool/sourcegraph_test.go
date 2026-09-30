package sdptool

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestComposedNavigationAndInventory(t *testing.T) {
	root, r := projectFixture(t)
	r.Models = []Model{{"model", "Demo", "System.design", "design-core/0.6"}}
	r.ImplementationPlan = "SDP/plan.md"
	saveInventory(t, root, r)
	const header = "language design-core version 0.6.\n"
	write := func(name, text string) {
		t.Helper()
		if e := os.WriteFile(filepath.Join(root, name), []byte(text), 0600); e != nil {
			t.Fatal(e)
		}
	}
	write("System.design", header+"system Demo.\nincludes \"Child.design\".\nDemo contains Child.\n")
	write("Child.design", header+"unit Child.\n")
	write(r.ImplementationPlan, "plan")
	p, e := Discover(root)
	p.Inventory = r
	if e != nil {
		t.Fatal(e)
	}
	tree, e := ModelTree(p, "model")
	if e != nil {
		t.Fatal(e)
	}
	out := filepath.Join(t.TempDir(), "preview")
	o := PreviewOptions{Output: out, URI: "sdl-view://" + r.ProjectID + "/VP02", Revision: tree.Revision}
	result, e := SelectProject(context.Background(), p, "model", o)
	if e != nil || result.Profile != "design-core/0.6" {
		t.Fatal(result, e)
	}
	provenance, e := os.ReadFile(filepath.Join(out, "sources.json"))
	if e != nil || !strings.Contains(string(provenance), "Child.design") {
		t.Fatal(string(provenance), e)
	}
	manifest, _ := os.ReadFile(filepath.Join(out, "manifest.json"))
	write("Child.design", header+"unit Child.\nunit Extra.\n")
	if _, e = SelectProject(context.Background(), p, "model", o); e == nil || !strings.Contains(e.Error(), "stale") {
		t.Fatal("old tree accepted", e)
	}
	after, _ := os.ReadFile(filepath.Join(out, "manifest.json"))
	if string(manifest) != string(after) {
		t.Fatal("failed refresh replaced output")
	}
	updated, e := ModelTree(p, "model")
	if e != nil || updated.Revision == tree.Revision {
		t.Fatal(updated, e)
	}
	o.Revision = updated.Revision
	viewer := executable(t, "if [ \"$1\" = '--help' ]; then echo '--navigator --sdl-tool'; exit 0; fi\nexit 0\n")
	tool := executable(t, "exit 0\n")
	if e = ViewPlan(context.Background(), p, "model", Host{Viewer: viewer, SDLTool: tool}); e != nil {
		t.Fatal(e)
	}
	for _, bad := range []Model{{"model", "Wrong", "System.design", "design-core/0.6"}, {"model", "Demo", "System.design", "design-core/0.5"}} {
		p.Inventory.Models = []Model{bad}
		if _, e = ModelTree(p, "model"); e == nil {
			t.Fatal("tree accepted mismatched inventory")
		}
		if _, e = SelectProject(context.Background(), p, "model", o); e == nil {
			t.Fatal("selection accepted mismatched inventory")
		}
		if e = ViewPlan(context.Background(), p, "model", Host{Viewer: viewer, SDLTool: tool}); e == nil {
			t.Fatal("viewer accepted mismatched inventory")
		}
	}
	// A dependency directory is protected even though the entry is outside it.
	os.Mkdir(filepath.Join(root, "Children"), 0700)
	os.Rename(filepath.Join(root, "Child.design"), filepath.Join(root, "Children", "Child.design"))
	write("System.design", header+"system Demo.\nDemo contains Children/Child.\n")
	if _, e = Preview(context.Background(), PreviewOptions{Source: filepath.Join(root, "System.design"), Output: filepath.Join(root, "Children")}); e == nil {
		t.Fatal("publication could replace dependency")
	}
	if _, e = os.Stat(filepath.Join(root, "Children", "Child.design")); e != nil {
		t.Fatal("source removed", e)
	}
}
