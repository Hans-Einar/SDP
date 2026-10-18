package viewpoint

import (
	"fmt"
	"strings"
)

// useCaseDetail changes presentation only: realization remains in companion text,
// and an ordered sequence is included only through an explicit exercises fact.
func (v *Views) useCaseDetails(out *Views, q Query) {
	if q.Viewpoint != "VP01" {
		return
	}
	focused := ""
	if v.Kinds[q.Focus] == "usecase" {
		focused = q.Focus
	}
	if strings.HasPrefix(q.Diagram, "VP01-") && v.Kinds[strings.TrimPrefix(q.Diagram, "VP01-")] == "usecase" {
		focused = strings.TrimPrefix(q.Diagram, "VP01-")
	}

	if focused != "" {
		diagrams := []Diagram{}
		for _, d := range out.Diagrams {
			if d.ID == "VP01-"+focused {
				diagrams = append(diagrams, d)
			}
		}
		out.Diagrams = diagrams
		features := values(eq(v.Relations["supports"], "object", focused), "subject")
		out.Gaps = filter(out.Gaps, func(g Fact) bool {
			return g.S("model_id") == focused || has(features, g.S("model_id"))
		})
		out.Facts = filter(out.Facts, func(f Fact) bool {
			return (f.S("verb") == "pursues" || f.S("verb") == "supports" || f.S("verb") == "contributes-to") && (f.S("object") == focused || (f.S("verb") == "contributes-to" && has(features, f.S("object"))))
		})
		out.Relations = map[string][]Fact{}
		for _, f := range out.Facts {
			out.Relations[f.S("verb")] = append(out.Relations[f.S("verb")], f)
		}
	}
	for i, old := range out.Diagrams {
		name := strings.TrimPrefix(old.ID, "VP01-")
		if v.Kinds[name] != "usecase" {
			continue
		}
		// Preserve explicit relation filtering. Realization facts stay in out.Facts.
		pursuits := eq(out.Relations["pursues"], "object", name)
		d := Diagram{ID: old.ID, Title: old.Title, Kind: "usecase", Nodes: map[string]any{}, Edges: []Edge{}, SourceFacts: values(pursuits, "id"), Elements: []Fact{}}
		names := unique(append(values(pursuits, "subject"), name))
		sourceNodes := v.nodes(names)
		var syntax strings.Builder
		syntax.WriteString("usecase-beta\ndirection LR\n")
		for _, n := range names {
			id := "n_" + n
			d.Nodes[id] = sourceNodes[id]
			if v.Kinds[n] == "actor" {
				syntax.WriteString("actor ")
			}
			fmt.Fprintf(&syntax, "%s(\"%s\")\n", id, n)
		}
		for _, f := range pursuits {
			a, b := "n_"+f.S("subject"), "n_"+name
			fmt.Fprintf(&syntax, "%s -- %s\n", a, b)
			d.Edges = append(d.Edges, Edge{a, "pursues", b, f.S("id")})
		}
		d.Syntax = syntax.String()
		out.Diagrams[i] = d
	}
	if focused == "" || len(q.Relations) > 0 {
		return
	}
	used := map[string]bool{}
	for _, f := range out.Facts {
		used[f.S("id")] = true
	}
	for _, link := range eq(v.Relations["exercises"], "object", focused) {
		for _, d := range v.Diagrams {
			if d.ID != "VP08-"+link.S("subject") {
				continue
			}
			modes := eq(v.Relations["runs-in"], "subject", link.S("subject"))
			if q.Mode != "" && !has(values(modes, "object"), q.Mode) {
				continue
			}
			d.SourceFacts = unique(append(append([]string{}, d.SourceFacts...), link.S("id")))
			out.Diagrams = append(out.Diagrams, d)
			for _, id := range d.SourceFacts {
				if used[id] {
					continue
				}
				for _, f := range v.Facts {
					if f.S("id") == id {
						out.Facts = append(out.Facts, f)
						used[id] = true
					}
				}
			}
		}
	}
}

func (v *Views) UseCaseContext() string {
	var b strings.Builder
	b.WriteString("## Use-case realization and conditions\n\nActors and use cases use UML shapes; associations mean modeled pursues. Features and Functionality are realization context, not use cases.\n\n| Source | Relation | Target | Fact |\n| --- | --- | --- | --- |\n")
	for _, f := range v.Facts {
		if has([]string{"supports", "contributes-to", "exercises"}, f.S("verb")) {
			fmt.Fprintf(&b, "| %s | %s | %s | %s |\n", f.S("subject"), f.S("verb"), f.S("object"), v.FactReference(f.S("id")))
		}
	}
	found := false
	for _, d := range v.Diagrams {
		if d.Kind == "sequence" {
			found = true
		}
	}
	if !found {
		b.WriteString("\nNo explicitly linked scenario sequence is selected. Ordering is not inferred from static relations.\n")
	}
	b.WriteString("\nPreconditions: not specified by the supported semantic profile.\n\nPostconditions: not specified by the supported semantic profile. Generic properties are not interpreted as conditions.\n\n")
	return b.String()
}
