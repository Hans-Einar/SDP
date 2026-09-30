package presentation

import (
	"encoding/json"
	"fmt"
	"io"
)

type node struct {
	ID, Label, State, Reference, ExternalReference, Diagnostic string
	Children                                                   []string
}

func tree(w io.Writer, b json.RawMessage) error {
	var t struct {
		Project             string
		Roots               []string
		Nodes               []node
		ExpansionDepthLimit int
	}
	if err := json.Unmarshal(b, &t); err != nil {
		return err
	}
	nodes := map[string]node{}
	for _, n := range t.Nodes {
		if _, exists := nodes[n.ID]; exists {
			return fmt.Errorf("duplicate navigation node %q", n.ID)
		}
		nodes[n.ID] = n
	}
	limit := t.ExpansionDepthLimit
	if limit <= 0 || limit > 64 {
		limit = 64
	}
	fmt.Fprintln(w, safe(t.Project))
	path := map[string]bool{}
	emitted := 0
	var walk func([]string, string, int) error
	walk = func(ids []string, prefix string, depth int) error {
		for i, id := range ids {
			emitted++
			if emitted > 20000 {
				return fmt.Errorf("navigation presentation exceeds 20000 entries")
			}
			branch, continuation := "├── ", "│   "
			if i == len(ids)-1 {
				branch, continuation = "└── ", "    "
			}
			n, exists := nodes[id]
			if !exists {
				return fmt.Errorf("missing navigation node %q", id)
			}
			label := safe(n.Label)
			if n.State != "" {
				label += " [" + safe(n.State) + "]"
			}
			if n.Diagnostic != "" {
				label += " — " + safe(n.Diagnostic)
			}
			if n.Reference != "" {
				label += " → " + safe(n.Reference)
			}
			if n.ExternalReference != "" {
				label += " → external: " + safe(n.ExternalReference)
			}
			if path[id] {
				label += " [cycle]"
			}
			fmt.Fprintln(w, prefix+branch+label)
			if path[id] || n.Reference != "" || n.ExternalReference != "" || len(n.Children) == 0 {
				continue
			}
			if depth >= limit {
				fmt.Fprintln(w, prefix+continuation+"└── … [depth limit]")
				continue
			}
			path[id] = true
			if err := walk(n.Children, prefix+continuation, depth+1); err != nil {
				return err
			}
			delete(path, id)
		}
		return nil
	}
	// A reference is shown as an edge, never recursively expanded.
	return walk(t.Roots, "", 1)
}
