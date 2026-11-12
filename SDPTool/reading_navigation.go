package sdptool

import (
	"fmt"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/sourcegraph"
	"net/url"
	"path/filepath"
	"sort"
	"strings"
)

const readingPresentation = "producer-reading/1"

func readingTarget(p Project, m Model, revision, vp, collection, focus, register string) *Target {
	q := url.Values{"projection": {"reading"}, "target": {"main"}, "consumer": {"xfmd"}}
	for k, v := range map[string]string{"collection": collection, "focus": focus, "register": register} {
		if v != "" {
			q.Set(k, v)
		}
	}
	return &Target{Operation: "select", Project: p.Inventory.ProjectID, Model: m.ID, Revision: revision, URI: "sdl-view://" + p.Inventory.ProjectID + "/" + vp + "?" + q.Encode()}
}
func sourceTarget(p Project, m Model, revision, source string) *Target {
	q := url.Values{"projection": {"source"}, "target": {"main"}, "consumer": {"xfmd"}}
	if source != "" {
		q.Set("source", source)
	}
	return &Target{Operation: "select", Project: p.Inventory.ProjectID, Model: m.ID, Revision: revision, URI: "sdl-view://" + p.Inventory.ProjectID + "/VP11?" + q.Encode()}
}

// Membership follows the explicit loader namespace even when the owning root is
// broken. It never compiles unrelated files into a synthetic system.
func (p *Project) discoverMembership() {
	indices := map[string]int{}
	for i, s := range p.Sources {
		indices[filepath.Clean(filepath.Join(p.Root, filepath.FromSlash(s.Source)))] = i
	}
	for _, root := range p.Sources {
		if root.Kind != "sdl" || !root.DeclaresSystem {
			continue
		}
		entry := filepath.Join(p.Root, filepath.FromSlash(root.Source))
		base := filepath.Dir(entry)
		seen := map[string]bool{}
		var visit func(string)
		visit = func(full string) {
			if seen[full] || len(seen) >= sourcegraph.MaxFiles {
				return
			}
			seen[full] = true
			if i, ok := indices[full]; ok && full != entry {
				p.Sources[i].IncludedBy = append(p.Sources[i].IncludedBy, root.ID)
			}
			b, e := readSource(full)
			if e != nil {
				return
			}
			m, e := parser.Parse(string(b))
			if e != nil {
				return
			}
			for _, inc := range m.Includes {
				if !sourcegraph.ValidPath(inc.Path) {
					continue
				}
				next := filepath.Join(base, filepath.FromSlash(inc.Path))
				if _, e = resolvePath(base, inc.Path); e == nil {
					visit(next)
				}
			}
		}
		visit(entry)
	}
	for i := range p.Sources {
		sort.Strings(p.Sources[i].IncludedBy)
	}
}
func readingSourceNodes(p Project, deferred bool) ([]Node, error) {
	groups := []Node{{ID: "sdl/systems", Kind: "group", Label: "Systems", State: "available"}, {ID: "sdl/standalone", Kind: "group", Label: "Standalone designs", State: "available"}, {ID: "sdl/fragments", Kind: "group", Label: "Unassigned fragments", State: "available"}}
	nodes := []Node{}
	tab := Node{ID: "sdl", Kind: "tab", Label: "SDL", State: "absent"}
	for _, m := range p.Inventory.Models {
		var info *SourceInfo
		for i := range p.Sources {
			if p.Sources[i].ID == m.ID {
				info = &p.Sources[i]
				break
			}
		}
		group := 1
		if info != nil && info.DeclaresSystem {
			group = 0
		} else if info != nil && info.State == "context-required" {
			if len(info.IncludedBy) > 0 {
				continue
			}
			group = 2
		}
		n := Node{ID: "sdl/" + m.ID, Kind: "source", Label: m.System + " — " + filepath.Base(m.Source), State: "available"}
		if group == 2 {
			n.Label = m.Source
		}
		if info != nil && info.State != "validated" {
			n.State = info.State
			n.Diagnostic = info.Diagnostic
			if info.State == "context-required" {
				n.State = "available"
				n.Diagnostic = "Local syntax only; no System context is available"
				n.Target = sourceTarget(p, m, info.Revision, "")
			}
			// Included files of broken systems remain reachable with explicit diagnostics.
			for _, s := range p.Sources {
				for _, owner := range s.IncludedBy {
					if owner == m.ID {
						id := n.ID + "/source/" + s.ID
						n.Children = append(n.Children, id)
						nodes = append(nodes, Node{ID: id, Kind: "file", Label: s.Source, State: "available", Diagnostic: "Owning System is invalid; open source to repair", Target: &Target{Operation: "open", Project: p.Inventory.ProjectID, Path: filepath.Join(p.Root, s.Source)}})
					}
				}
			}
		} else if deferred && info != nil {
			n.State = "validated"
			n.Target = readingTarget(p, m, info.Revision, "VP02", "", "", "")
			n.ChildrenDeferred = true
			n.ExpansionTarget = &Target{Operation: "tree", Project: p.Inventory.ProjectID, Model: m.ID, Revision: info.Revision}
		} else {
			mt, e := ModelTree(p, m.ID)
			if e != nil {
				n.State = "invalid"
				n.Diagnostic = e.Error()
			} else {
				n.State = "validated"
				n.Target = readingTarget(p, m, mt.Revision, "VP02", "", "", "")
				for _, child := range mt.Nodes {
					if child.ID == "sdl" {
						n.Children = child.Children
					} else {
						nodes = append(nodes, child)
					}
				}
			}
		}
		groups[group].Children = append(groups[group].Children, n.ID)
		nodes = append(nodes, n)
	}
	for _, g := range groups {
		if len(g.Children) > 0 {
			tab.Children = append(tab.Children, g.ID)
			tab.State = "available"
			nodes = append(nodes, g)
		}
	}
	nodes = append(nodes, tab)
	if len(nodes) > 20000 {
		return nil, failure("limit", fmt.Errorf("source navigation exceeds 20000 nodes"))
	}
	return nodes, nil
}

func isSourceSelection(raw string) bool {
	u, e := url.Parse(raw)
	return e == nil && u.Scheme == "sdl-view" && strings.EqualFold(u.Query().Get("projection"), "source")
}
