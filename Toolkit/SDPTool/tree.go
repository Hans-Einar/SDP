package sdptool

import (
	"encoding/json"
	"fmt"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/documents"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/viewpoint"
	"net/url"
	"sort"
	"strings"
)

type Target struct {
	Entry     string `json:"entry,omitempty"`
	Operation string `json:"operation"`
	Project   string `json:"project"`
	Model     string `json:"model,omitempty"`
	URI       string `json:"uri,omitempty"`
	Path      string `json:"path,omitempty"`
	Revision  string `json:"revision,omitempty"`
}
type Node struct {
	Diagnostic string   `json:"diagnostic,omitempty"`
	ID         string   `json:"id"`
	Kind       string   `json:"kind"`
	Label      string   `json:"label"`
	State      string   `json:"state"`
	Children   []string `json:"children,omitempty"`
	Reference  string   `json:"reference,omitempty"`
	Target     *Target  `json:"target,omitempty"`
	WorkState  string   `json:"cardState,omitempty"`
	Sprint     string   `json:"sprintId,omitempty"`
	Scrum      string   `json:"scrumId,omitempty"`
}
type Tree struct {
	InventoryRevision   string   `json:"inventoryRevision,omitempty"`
	Schema              string   `json:"schema"`
	Operation           string   `json:"operation"`
	Project             string   `json:"project"`
	Model               string   `json:"model,omitempty"`
	Revision            string   `json:"revision,omitempty"`
	Roots               []string `json:"roots"`
	Nodes               []Node   `json:"nodes"`
	ExpansionDepthLimit int      `json:"expansionDepthLimit"`
}

func selectionURI(project, vp, focus, diagram string) string {
	q := url.Values{"target": {"main"}, "consumer": {"xfmd"}}
	if focus != "" {
		q.Set("focus", focus)
	}
	if diagram != "" {
		q.Set("diagram", diagram)
	}
	return "sdl-view://" + project + "/" + vp + "?" + q.Encode()
}
func ModelTree(p Project, id string) (Tree, error) {
	m, path, e := p.model(id, false)
	if e != nil {
		return Tree{}, e
	}
	v, _, e := loadModel(path)
	if e != nil {
		return Tree{}, e
	}
	return modelTree(p, m, v)
}
func modelTree(p Project, m Model, v *viewpoint.Views) (Tree, error) {
	t := Tree{Schema: Version, Operation: "tree", Project: p.Registration.ProjectID, Model: m.ID, Revision: v.Revision, Roots: []string{"sdl"}, Nodes: []Node{}, ExpansionDepthLimit: 8}
	if len(v.Kinds) > 2000 || len(v.Facts) > 10000 {
		return Tree{}, failure("limit", fmt.Errorf("navigation model exceeds 2000 declarations or 10000 facts; select a smaller source"))
	}
	factsByID := map[string]viewpoint.Fact{}
	for _, f := range v.Facts {
		factsByID[f.S("id")] = f
	}
	prefix := "sdl/" + m.ID + "/"
	nodes := map[string]*Node{}
	add := func(n Node) { copy := n; nodes[n.ID] = &copy }
	target := func(vp, focus, diagram string) *Target {
		return &Target{Operation: "select", Project: t.Project, Model: m.ID, URI: selectionURI(t.Project, vp, focus, diagram), Revision: v.Revision}
	}
	objectID := func(name string) string { return prefix + "object/" + v.Kinds[name] + "/" + url.PathEscape(name) }
	add(Node{ID: "sdl", Kind: "tab", Label: "SDL", State: "validated", Children: []string{prefix + "requirements", prefix + "architecture", prefix + "design", prefix + "implementation"}})
	for _, phase := range []struct {
		id, label string
		levels    []string
	}{{"requirements", "Requirements", []string{"A0", "A1"}}, {"architecture", "Architecture", []string{"A2", "A3"}}, {"design", "Design", []string{"A4"}}, {"implementation", "Implementation", []string{"A5", "delivery"}}} {
		n := Node{ID: prefix + phase.id, Kind: "phase", Label: phase.label, State: "available"}
		for _, level := range phase.levels {
			key := prefix + level
			n.Children = append(n.Children, key)
			add(Node{ID: key, Kind: "abstraction", Label: level, State: "available"})
		}
		add(n)
	}
	names := make([]string, 0, len(v.Kinds))
	for name := range v.Kinds {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		add(Node{ID: objectID(name), Kind: v.Kinds[name], Label: name, State: "validated", Target: target("VP11", name, "")})
	}
	for _, spec := range viewpoint.Catalog {
		vpID := prefix + spec.ID
		n := Node{ID: vpID, Kind: "viewpoint", Label: spec.ID + " — " + spec.Title, State: "empty", Target: target(spec.ID, "", "")}
		levels := spec.Levels
		if len(levels) == 0 {
			levels = []string{"delivery"}
		}
		for _, level := range levels {
			ref := prefix + level + "/" + spec.ID
			nodes[prefix+level].Children = append(nodes[prefix+level].Children, ref)
			add(Node{ID: ref, Kind: "reference", Label: spec.ID, State: "available", Reference: vpID})
		}
		used := map[string]bool{}
		for _, d := range v.Diagrams {
			if !strings.HasPrefix(d.ID, spec.ID+"-") {
				continue
			}
			n.State = "available"
			key := vpID + "/diagram/" + d.ID
			n.Children = append(n.Children, key)
			add(Node{ID: key, Kind: "diagram", Label: d.Title, State: "available", Target: target(spec.ID, "", d.ID)})
			for node, data := range d.Nodes {
				if facts, ok := data.(map[string]any); ok {
					if name, ok := facts["model_id"].(string); ok && v.Kinds[name] != "" {
						used[name] = true
					}
				}
				name := strings.TrimPrefix(node, "n_")
				if v.Kinds[name] != "" {
					used[name] = true
				}
			}
			// Sequence/packet membership is carried by source facts, not necessarily nodes.
			for _, fid := range d.SourceFacts {
				for _, f := range []viewpoint.Fact{factsByID[fid]} {
					for _, field := range []string{"subject", "object", "sender", "receiver", "channel", "message", "field", "dataset", "datagram", "mode"} {
						name := f.S(field)
						if v.Kinds[name] != "" {
							used[name] = true
						}
					}
				}
			}
		}
		if spec.ID == "VP11" {
			for _, name := range names {
				used[name] = true
			}
			if len(names) > 0 {
				n.State = "available"
			}
		}
		kinds := map[string][]string{}
		for _, name := range names {
			if used[name] {
				kind := v.Kinds[name]
				kinds[kind] = append(kinds[kind], name)
			}
		}
		ks := []string{}
		for kind := range kinds {
			ks = append(ks, kind)
		}
		sort.Strings(ks)
		for _, kind := range ks {
			key := vpID + "/kind/" + kind
			collection := Node{ID: key, Kind: "collection", Label: kind, State: "available"}
			for _, name := range kinds[kind] {
				ref := key + "/" + url.PathEscape(name)
				collection.Children = append(collection.Children, ref)
				add(Node{ID: ref, Kind: "reference", Label: name, State: "available", Reference: objectID(name), Target: target(spec.ID, name, "")})
			}
			add(collection)
			n.Children = append(n.Children, key)
		}
		add(n)
	}
	// Relationship targets are explicit references, never recursively copied subtrees.
	// Hash semantic endpoints rather than source-order fact IDs for stable edges.
	for _, f := range v.Facts {
		subject := f.S("subject")
		if v.Kinds[subject] == "" {
			continue
		}
		for _, field := range []string{"object", "sender", "receiver", "channel", "message", "field", "dataset", "datagram", "mode"} {
			name := f.S(field)
			if v.Kinds[name] == "" {
				continue
			}
			identity, _ := json.Marshal([]string{subject, f.S("verb"), field, name, f.S("mode")})
			key := prefix + "relation/" + documents.Hash(identity)
			if nodes[key] != nil {
				continue
			}
			add(Node{ID: key, Kind: "relationship", Label: f.S("verb") + " (" + field + "): " + name, State: "validated", Reference: objectID(name), Target: target("VP11", subject, "")})
			parent := nodes[objectID(subject)]
			parent.Children = append(parent.Children, key)
		}
	}

	if len(nodes) > 20000 {
		return Tree{}, failure("limit", fmt.Errorf("tree exceeds 20000 nodes"))
	}
	keys := []string{}
	for k := range nodes {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		t.Nodes = append(t.Nodes, *nodes[k])
	}
	return t, nil
}
