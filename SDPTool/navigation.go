package sdptool

import (
	"encoding/json"
	"fmt"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/documents"
)

// Navigation keeps optional unavailable services visible instead of claiming
// that absence means a validated empty system.
func Navigation(p Project, id string) (Tree, error) {
	if p.Status == "incomplete" {
		return Tree{}, failure("incomplete", fmt.Errorf("installation has not published navigation; resume its recorded operation"))
	}

	t := Tree{Schema: Version, Operation: "tree", Project: p.Inventory.ProjectID, Roots: []string{"sdl", "kanban", "sdui", "sessions"}, Nodes: []Node{}, ExpansionDepthLimit: 8}
	if id != "" {
		mt, e := ModelTree(p, id)
		if e != nil {
			return t, e
		}
		t.Model = mt.Model
		t.Revision = mt.Revision
		t.Nodes = mt.Nodes
	} else {
		tab := Node{ID: "sdl", Kind: "tab", Label: "SDL", State: "absent"}
		for _, m := range p.Inventory.Models {
			key := "sdl/" + m.ID
			n := Node{ID: key, Kind: "source", Label: m.Source, State: "available", Target: &Target{Operation: "tree", Project: p.Inventory.ProjectID, Model: m.ID}}
			tab.Children = append(tab.Children, key)
			tab.State = "available"
			var info *SourceInfo
			for i := range p.Sources {
				if p.Sources[i].ID == m.ID {
					info = &p.Sources[i]
					break
				}
			}
			if info != nil && info.State != "validated" {
				n.State = info.State
				n.Diagnostic = info.Diagnostic
				n.Target = nil
			} else {
				mt, e := ModelTree(p, m.ID)
				if e != nil {
					n.State = "invalid"
					n.Diagnostic = e.Error()
					n.Target = nil
				} else {
					n.State = "validated"
					n.Target.Revision = mt.Revision
					for _, mn := range mt.Nodes {
						if mn.ID == "sdl" {
							n.Children = mn.Children
						} else {
							t.Nodes = append(t.Nodes, mn)
						}
					}
				}
			}
			t.Nodes = append(t.Nodes, n)
		}
		t.Nodes = append(t.Nodes, tab)
	}
	if len(p.Files) > 0 {
		t.Roots = append([]string{"files"}, t.Roots...)
		t.Nodes = append(t.Nodes, p.Files...)
	}
	t.Nodes = append(t.Nodes, sessionsNode(p))
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
