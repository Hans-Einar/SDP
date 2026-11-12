package sdptool

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestActivityNavigationEmptyAndInvalid(t *testing.T) {
	area := t.TempDir()
	if e := os.MkdirAll(filepath.Join(area, "ProjectManagement"), 0755); e != nil {
		t.Fatal(e)
	}
	if e := os.WriteFile(filepath.Join(area, "ProjectManagement", "Ledger.ndjson"), nil, 0644); e != nil {
		t.Fatal(e)
	}
	p := Project{Area: area, Root: filepath.Dir(area), Inventory: Inventory{ProjectID: "test"}}
	nodes := activityNodes(p)
	for _, state := range []string{"planned", "active", "on-hold", "completed"} {
		found := false
		for _, n := range nodes {
			if n.ID == "activities/state/"+state {
				found = true
				if n.State != "empty" || n.Target == nil || n.Target.Model != "" {
					t.Fatal(n)
				}
			}
		}
		if !found {
			t.Fatalf("missing %s: %+v", state, nodes)
		}
	}
	for _, n := range nodes {
		if n.ID == "activities/state/canceled" {
			t.Fatal("empty canceled group")
		}
	}
	page, e := renderActivityView(p, "activity-view://test/state/planned")
	if e != nil || !strings.Contains(page.Markdown, "No Activities") {
		t.Fatal(page, e)
	}
	if _, e := renderActivityView(p, "activity-view://other/overview"); e == nil {
		t.Fatal("cross-project target accepted")
	}
	if e := os.WriteFile(filepath.Join(area, "ProjectManagement", "Ledger.ndjson"), []byte("invalid\n"), 0644); e != nil {
		t.Fatal(e)
	}
	nodes = activityNodes(p)
	if len(nodes) != 1 || nodes[0].State != "invalid" || nodes[0].Diagnostic == "" {
		t.Fatal(nodes)
	}
}
