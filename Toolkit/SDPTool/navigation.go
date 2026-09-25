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

	t := Tree{Schema: Version, Operation: "tree", Project: p.Registration.ProjectID, Roots: []string{"sdl", "kanban", "sdui"}, Nodes: []Node{}, ExpansionDepthLimit: 8}
	if len(p.Registration.Models) > 0 {
		modelTree, e := ModelTree(p, id)
		if e != nil {
			return t, e
		}
		t.Model = modelTree.Model
		t.Revision = modelTree.Revision
		t.Nodes = modelTree.Nodes
	} else {
		t.Nodes = append(t.Nodes, Node{ID: "sdl", Kind: "tab", Label: "SDL", State: "absent"})
	}
	versions := []string{t.Revision}
	if p.Registration.KanBan != "" {
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
		Registration Registration
		Versions     []string
		Nodes        []Node
	}{p.Registration, versions, t.Nodes})
	if len(t.Nodes) > 20000 || len(b) > 32<<20 {
		return Tree{}, failure("limit", fmt.Errorf("combined inventory exceeds 20000 nodes or 32 MiB"))
	}
	t.InventoryRevision = documents.Hash(b)
	return t, nil
}
