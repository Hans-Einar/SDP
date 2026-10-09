package sdptool

import (
	"context"
	"encoding/json"
	"github.com/Hans-Einar/SDP/SDPTool/blueprints"
	"github.com/Hans-Einar/SDP/SDPTool/model"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/blueprint"
	"os"
	"path/filepath"
	"testing"
)

func retainedFixture(t *testing.T, root string) blueprints.Result {
	t.Helper()
	area := t.TempDir()
	for _, name := range []string{"Before", "After"} {
		w, e := model.CreateWork(area, name, "", true)
		if e != nil {
			t.Fatal(e)
		}
		if e = os.CopyFS(w.Path, os.DirFS("../experiments/blueprint_mvp1/NOW")); e != nil {
			t.Fatal(e)
		}
	}
	task := filepath.Join(t.TempDir(), "task.json")
	b, _ := json.Marshal(blueprint.Task{ID: "navigation", Intent: "Inspect model"})
	os.WriteFile(task, b, 0600)
	r, e := blueprints.GenerateRetained(context.Background(), blueprints.Options{Area: area, From: "work:Before", To: "work:After", Entry: "System.design", Task: task}, filepath.Join(root, "SDP", "Blueprints"))
	if e != nil {
		t.Fatal(e)
	}
	return r
}
func TestBlueprintDiscoveryProjection(t *testing.T) {
	root := t.TempDir()
	os.Mkdir(filepath.Join(root, "SDP"), 0700)
	p, e := Discover(root)
	if e != nil {
		t.Fatal(e)
	}
	if p.Capabilities["blueprints"] != "absent" {
		t.Fatal(p.Capabilities)
	}
	r := retainedFixture(t, root)
	p, e = Discover(root)
	if e != nil {
		t.Fatal(e)
	}
	if p.Inventory.Blueprints != "SDP/Blueprints" || p.Capabilities["blueprints"] != "discovered" || len(p.Inventory.Models) != 0 {
		t.Fatal("retained sources counted as active models", p.Inventory)
	}
	seen := map[string]bool{}
	tab := false
	revision := false
	targets := 0
	for _, n := range p.Navigation.Nodes {
		if seen[n.ID] {
			t.Fatal("duplicate ID", n.ID)
		}
		seen[n.ID] = true
		if n.ID == "blueprints" {
			tab = n.Kind == "tab" && n.State == "available"
		}
		if n.BlueprintRevision == r.Revision {
			revision = n.AssignmentState == "unknown" && n.Preliminary
		}
		if n.Kind == "document" && n.Target != nil {
			targets++
			if n.Target.Operation != "open" || len(n.Target.Revision) != 64 {
				t.Fatal(n)
			}
			if _, e = os.Stat(n.Target.Path); e != nil {
				t.Fatal(e)
			}
		}
	}
	if !tab || !revision || targets != 7 {
		t.Fatalf("tab=%t revision=%t targets=%d", tab, revision, targets)
	}
	for _, n := range p.Navigation.Nodes {
		for _, id := range n.Children {
			if !seen[id] {
				t.Fatal("dangling", id)
			}
		}
	}
	tree, e := Navigation(p, "")
	if e != nil {
		t.Fatal(e)
	}
	a, _ := json.Marshal(tree)
	b, _ := json.Marshal(p.Navigation)
	if string(a) != string(b) {
		t.Fatal("tree/discovery disagreement")
	}
	previous := p.Navigation.InventoryRevision
	path := filepath.Join(r.Path, "index.md")
	st, _ := os.Stat(path)
	os.WriteFile(path, []byte("bad"), 0600)
	os.Chtimes(path, st.ModTime(), st.ModTime())
	p, e = Discover(root)
	if e != nil || p.Navigation.InventoryRevision == previous {
		t.Fatal(e, "refresh missed edit")
	}
	bad := false
	for _, n := range p.Navigation.Nodes {
		if n.Kind == "blueprint-revision" && n.State == "invalid" {
			bad = true
			if n.Target != nil {
				t.Fatal("invalid open target")
			}
		}
	}
	if !bad {
		t.Fatal("invalid entry hidden")
	}
	os.RemoveAll(filepath.Join(root, "SDP", "Blueprints"))
	p, e = Discover(root)
	if e != nil || p.Capabilities["blueprints"] != "absent" {
		t.Fatal(e)
	}
}

func TestBlueprintAssignmentGroupsWithoutBundles(t *testing.T) {
	root := t.TempDir()
	area := filepath.Join(root, "SDP")
	os.MkdirAll(filepath.Join(area, "ProjectManagement"), 0700)
	raw, e := os.ReadFile("blueprintstate/testdata/lifecycle.ndjson")
	if e != nil {
		t.Fatal(e)
	}
	history := filepath.Join(area, "ProjectManagement/Ledger.ndjson")
	os.WriteFile(history, raw, 0600)
	p, e := Discover(root)
	if e != nil {
		t.Fatal(e)
	}
	found := false
	for _, n := range p.Navigation.Nodes {
		if n.ID == "blueprints" && n.State != "available" {
			t.Fatal("assignments hidden by absent bundle catalogue", n)
		}
		if n.AssignmentID == "BPA-one" {
			found = true
			if n.AssignmentState != "completed" || n.Target != nil || n.ReadinessStatus != "changed-or-unavailable" || n.EvidenceStatus != "changed-or-unavailable" {
				t.Fatal(n)
			}
		}
	}
	if !found {
		t.Fatal("missing completed assignment")
	}
	os.WriteFile(history, append(raw, []byte("{incomplete")...), 0600)
	p, e = Discover(root)
	if e != nil {
		t.Fatal(e)
	}
	found = false
	for _, n := range p.Navigation.Nodes {
		if n.ID == "blueprints" {
			found = true
			if n.State != "unavailable" || n.Diagnostic == "" {
				t.Fatal(n)
			}
		}
	}
	if !found {
		t.Fatal("missing unavailable tab")
	}
}
