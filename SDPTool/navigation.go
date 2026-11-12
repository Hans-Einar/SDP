package sdptool

import (
	"encoding/json"
	"fmt"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/documents"
)

// Navigation keeps optional unavailable services visible instead of claiming
// that absence means a validated empty system.
func Navigation(p Project, id string) (Tree, error) {
	return navigation(p, id, false)
}

func navigation(p Project, id string, deferredSDL bool) (Tree, error) {
	if p.Status == "incomplete" {
		return Tree{}, failure("incomplete", fmt.Errorf("installation has not published navigation; resume its recorded operation"))
	}

	t := Tree{Presentation: readingPresentation, Schema: Version, Operation: "tree", Project: p.Inventory.ProjectID, Roots: []string{"sdl", "kanban", "sdui", "sessions", "blueprints"}, Nodes: []Node{}, ExpansionDepthLimit: 8}
	if id != "" {
		mt, e := ModelTree(p, id)
		if e != nil {
			return t, e
		}
		t.Model = mt.Model
		t.Revision = mt.Revision
		t.Nodes = mt.Nodes
	} else {
		nodes, e := readingSourceNodes(p, deferredSDL)
		if e != nil {
			return t, e
		}
		t.Nodes = append(t.Nodes, nodes...)
	}

	if len(p.Files) > 0 {
		t.Roots = append([]string{"files"}, t.Roots...)
		t.Nodes = append(t.Nodes, p.Files...)
	}
	if activityPresent(p) {
		for i, root := range t.Roots {
			if root == "sessions" {
				t.Roots[i] = "activities"
			}
		}
		nodes, e := initialActivityTree(p)
		if e != nil {
			nodes = []Node{{ID: "activities", Kind: "tab", Label: "Activities", State: "invalid", Diagnostic: e.Error()}}
		}
		t.Nodes = append(t.Nodes, nodes...)
	} else {
		t.Nodes = append(t.Nodes, sessionsNode(p))
	}
	t.Nodes = append(t.Nodes, blueprintNodes(p)...)
	versions := []string{t.Revision}
	if p.Inventory.KanBan != "" {
		nodes, hash, e := BoardNodes(p)
		if e != nil {
			t.Nodes = append(t.Nodes, Node{ID: "kanban", Kind: "tab", Label: "KanBan", State: "unavailable", Diagnostic: e.Error()})
		} else {
			t.Nodes = append(t.Nodes, nodes...)
			versions = append(versions, hash)
		}
	} else {
		t.Nodes = append(t.Nodes, Node{ID: "kanban", Kind: "tab", Label: "KanBan", State: "absent"})
	}
	nodes, hash, e := UINodes(p)
	if e != nil {
		t.Nodes = append(t.Nodes, Node{ID: "sdui", Kind: "tab", Label: "SDUI", State: "unavailable", Diagnostic: e.Error()})
	} else {
		t.Nodes = append(t.Nodes, nodes...)
		versions = append(versions, hash)
	}
	b, _ := json.Marshal(struct {
		Inventory Inventory
		Versions  []string
		Nodes     []Node
	}{p.Inventory, versions, t.Nodes})
	if len(t.Nodes) > 20000 || len(b) > 32<<20 {
		return Tree{}, failure("limit", fmt.Errorf("combined inventory exceeds 20000 nodes or 32 MiB"))
	}
	t.InventoryRevision = documents.Hash(b)
	return t, nil
}

// Sessions is a second entry point into the canonical file tree, not another
// registry or a claim of validation of the manual Session document format.
func sessionsNode(p Project) Node {
	tab := Node{ID: "sessions", Kind: "tab", Label: "Sessions", State: "absent"}
	for _, n := range p.Files {
		if n.ID != "files/Sessions" {
			continue
		}
		if n.Kind != "directory" {
			tab.State = "unavailable"
			tab.Diagnostic = "SDP/Sessions must be a real directory; symlinks are not followed"
			return tab
		}
		tab.State, tab.Diagnostic = n.State, n.Diagnostic
		tab.Children = n.Children
		if tab.State == "available" && len(tab.Children) == 0 {
			tab.State = "empty"
		}
		return tab
	}
	return tab
}
