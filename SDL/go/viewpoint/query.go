package viewpoint

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

type Query struct {
	Viewpoint string   `json:"viewpoint"`
	Focus     string   `json:"focus,omitempty"`
	Relations []string `json:"relations,omitempty"`
	Direction string   `json:"direction,omitempty"`
	Depth     int      `json:"depth,omitempty"`
	Level     string   `json:"level,omitempty"`
	Mode      string   `json:"mode,omitempty"`
	Diagram   string   `json:"diagram,omitempty"`
}
type Selection struct {
	Project          string
	Query            Query
	Target, Consumer string
}

func ParseURI(raw string) (Selection, error) {
	u, e := url.Parse(raw)
	if e != nil {
		return Selection{}, e
	}
	if u.Scheme != "sdl-view" || u.Host == "" || u.User != nil || u.Port() != "" || u.Fragment != "" {
		return Selection{}, fmt.Errorf("invalid sdl-view URI")
	}
	s := Selection{Project: u.Host, Query: Query{Viewpoint: strings.TrimPrefix(u.Path, "/"), Direction: "both", Depth: 2}, Target: "main", Consumer: "xfmd"}
	q, e := url.ParseQuery(u.RawQuery)
	if e != nil {
		return s, e
	}
	for k, vs := range q {
		if len(vs) != 1 {
			return s, fmt.Errorf("duplicate query parameter %s", k)
		}
		value := vs[0]
		switch k {
		case "focus":
			s.Query.Focus = value
		case "relations":
			if value != "" {
				s.Query.Relations = strings.Split(value, ",")
			}
		case "direction":
			s.Query.Direction = value
		case "depth":
			s.Query.Depth, e = strconv.Atoi(value)
			if e != nil {
				return s, e
			}
		case "level":
			s.Query.Level = value
		case "mode":
			s.Query.Mode = value
		case "diagram":
			s.Query.Diagram = value
		case "target":
			s.Target = value
		case "consumer":
			s.Consumer = value
		default:
			return s, fmt.Errorf("unknown query parameter %s", k)
		}
	}
	if (s.Target != "main" && s.Target != "navigation") || s.Consumer != "xfmd" {
		return s, fmt.Errorf("unregistered target or consumer")
	}
	return s, nil
}
func (q Query) Validate(v *Views) error {
	s, ok := SpecFor(q.Viewpoint)
	if !ok {
		return fmt.Errorf("unknown-viewpoint: %s", q.Viewpoint)
	}
	if q.Focus != "" && v.Kinds[q.Focus] == "" {
		return fmt.Errorf("unknown-focus: %s", q.Focus)
	}
	if q.Direction != "" && !has([]string{"in", "out", "both"}, q.Direction) {
		return fmt.Errorf("invalid-direction")
	}
	if q.Depth < 0 || q.Depth > 8 {
		return fmt.Errorf("depth-limit: 0..8")
	}
	if q.Level != "" && !has(s.Levels, q.Level) {
		return fmt.Errorf("unsupported-level: %s for %s", q.Level, q.Viewpoint)
	}
	if q.Mode != "" && v.Kinds[q.Mode] != "mode" {
		return fmt.Errorf("unknown-mode: %s", q.Mode)
	}
	supported := strings.Fields("contains owns realizes provides consumes requires allocated-to refines pursues supports contributes-to addresses delivers depends-on illustrates upholds from holds defines has-field encodes permits replies-to runs-in exercises projects places uses step has")
	for _, r := range q.Relations {
		if !has(supported, r) {
			return fmt.Errorf("unsupported-relation: %s", r)
		}
	}
	if q.Diagram != "" {
		found := false
		for _, d := range v.Diagrams {
			if d.ID == q.Diagram && strings.HasPrefix(d.ID, q.Viewpoint+"-") {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("unknown-diagram")
		}
	}
	return nil
}

// Select never reparses a truncated source. Validation belongs to the complete
// snapshot; this operation selects already validated facts and projections.
func (v *Views) Select(q Query) (*Views, error) {
	if e := q.Validate(v); e != nil {
		return nil, e
	}
	out := *v
	out.Diagrams = []Diagram{}
	out.Gaps = []Fact{}
	out.Facts = []Fact{}
	out.Relations = map[string][]Fact{}
	out.MessageSets = []Fact{}
	facts := map[string]Fact{}
	for _, f := range v.Facts {
		facts[f.S("id")] = f
	}
	eligible := func(f Fact) bool {
		if len(q.Relations) > 0 && !has(q.Relations, f.S("verb")) {
			return false
		}
		return q.Mode == "" || f.S("mode") == "" || f.S("mode") == q.Mode
	}
	reachable := map[string]bool{}
	if q.Focus != "" {
		reachable[q.Focus] = true
		for i := 0; i < q.Depth; i++ {
			next := map[string]bool{}
			for n := range reachable {
				next[n] = true
			}
			for _, f := range v.Facts {
				if !eligible(f) {
					continue
				}
				a, b := f.S("subject"), f.S("object")
				if b == "" {
					continue
				}
				if q.Direction != "in" && reachable[a] {
					next[b] = true
				}
				if q.Direction != "out" && reachable[b] {
					next[a] = true
				}
			}
			reachable = next
		}
	}
	used := map[string]bool{}
	for _, original := range v.Diagrams {
		if !strings.HasPrefix(original.ID, q.Viewpoint+"-") || (q.Diagram != "" && original.ID != q.Diagram) {
			continue
		}
		d := original
		// Sequence/packet contracts remain atomic; pruning steps would corrupt their proof.
		if d.Kind != "flowchart" {
			if len(q.Relations) > 0 {
				return nil, fmt.Errorf("unsupported-selection: relation pruning of %s", d.Kind)
			}
			hit := q.Focus == ""
			modeOK := q.Mode == ""
			for _, id := range d.SourceFacts {
				f := facts[id]
				if f.S("mode") == q.Mode || f.S("object") == q.Mode {
					modeOK = true
				}
				for _, k := range []string{"subject", "object", "sender", "receiver", "channel", "message", "field"} {
					if reachable[f.S(k)] {
						hit = true
					}
				}
			}
			if !hit || !modeOK {
				continue
			}
		} else {
			nodes := map[string]any{}
			edges := []Edge{}
			proof := []string{}
			filtered := q.Focus != "" || len(q.Relations) > 0 || q.Mode != ""
			for n, p := range d.Nodes {
				if q.Focus == "" || reachable[strings.TrimPrefix(n, "n_")] {
					nodes[n] = p
				}
			}
			for _, e := range d.Edges {
				f := facts[e.Fact]
				if eligible(f) && nodes[e.Source] != nil && nodes[e.Target] != nil {
					edges = append(edges, e)
					proof = append(proof, e.Fact)
				}
			}
			if filtered {
				active := map[string]bool{}
				for _, e := range edges {
					active[e.Source] = true
					active[e.Target] = true
				}
				if q.Focus != "" {
					active["n_"+q.Focus] = true
				}
				for n := range nodes {
					if !active[n] {
						delete(nodes, n)
					}
				}
			}
			if len(nodes) == 0 {
				continue
			}
			d.Nodes = nodes
			d.Edges = edges
			d.SourceFacts = unique(proof)
		}
		if len(d.Nodes) > 160 || len(d.Edges) > 320 || len(out.Diagrams) >= 48 {
			return nil, fmt.Errorf("selection-limit: narrow focus, relation or diagram")
		}
		out.Diagrams = append(out.Diagrams, d)
		for _, f := range d.SourceFacts {
			used[f] = true
		}
	}
	for _, f := range v.Facts {
		if eligible(f) && (q.Focus == "" || reachable[f.S("subject")] || used[f.S("id")]) {
			out.Facts = append(out.Facts, f)
			if f.S("object") != "" {
				out.Relations[f.S("verb")] = append(out.Relations[f.S("verb")], f)
			}
		}
	}
	for _, g := range v.Gaps {
		if g.S("viewpoint") == q.Viewpoint && (q.Focus == "" || reachable[g.S("model_id")]) && (q.Mode == "" || g.S("mode") == "" || g.S("mode") == q.Mode) {
			out.Gaps = append(out.Gaps, g)
		}
	}
	for _, m := range v.MessageSets {
		if (q.Mode == "" || m.S("mode") == q.Mode) && (q.Focus == "" || reachable[m.S("channel")] || reachable[m.S("message")]) {
			out.MessageSets = append(out.MessageSets, m)
		}
	}
	return &out, nil
}
func (q Query) Canonical() Query {
	q.Relations = append([]string{}, q.Relations...)
	sort.Strings(q.Relations)
	if q.Direction == "" {
		q.Direction = "both"
	}
	return q
}
