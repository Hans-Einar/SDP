package viewpoint

import (
	"fmt"
	"strings"
)

// HasOverview advertises implemented projections only. It never generates output.
func HasOverview(vp string) bool { return vp == "VP01" || vp == "VP07" }

// overview selects the complete validated inventory, not a union of detail views.
func (v *Views) overview(q Query) (*Views, error) {
	if q.Viewpoint == "VP07" {
		return v.featureOverview(q)
	}
	out := *v
	out.Diagrams = []Diagram{}
	out.MessageSets = []Fact{}
	out.Facts = filter(v.Facts, func(f Fact) bool { return has([]string{"pursues", "supports", "contributes-to"}, f.S("verb")) })
	out.Relations = map[string][]Fact{}
	for _, f := range out.Facts {
		out.Relations[f.S("verb")] = append(out.Relations[f.S("verb")], f)
	}
	out.Gaps = filter(v.Gaps, func(f Fact) bool { return f.S("viewpoint") == q.Viewpoint })
	names := append(v.names("actor"), v.names("usecase")...)
	if len(names) == 0 {
		return &out, nil
	}
	// Match the documented native profile's total record/identity boundaries.
	// Reject the whole projection rather than silently omitting declarations.
	if len(names) > 128 || len(names)+len(v.Relations["pursues"])+1 > 512 {
		return nil, fmt.Errorf("overview-limit: complete VP01 exceeds native profile; no partial overview generated")
	}
	d := Diagram{ID: "VP01-overview", Title: "All declared actors and use cases", Kind: "usecase", Nodes: map[string]any{}, Edges: []Edge{}, SourceFacts: []string{}, Elements: []Fact{}}
	ids := map[string]string{}
	sourceNodes := v.nodes(names)
	var syntax strings.Builder
	syntax.WriteString("usecase-beta\ndirection LR\n")
	for i, name := range names {
		id := fmt.Sprintf("u%d", i)
		ids[name] = id
		d.Nodes[id] = sourceNodes["n_"+name]
		if v.Kinds[name] == "actor" {
			syntax.WriteString("actor ")
		}
		fmt.Fprintf(&syntax, "%s(\"%s\")\n", id, name)
	}
	for _, f := range v.Relations["pursues"] {
		a, b := ids[f.S("subject")], ids[f.S("object")]
		if a == "" || b == "" {
			return nil, fmt.Errorf("invalid checked pursues endpoints")
		}
		fmt.Fprintf(&syntax, "%s -- %s\n", a, b)
		d.Edges = append(d.Edges, Edge{a, "pursues", b, f.S("id")})
		d.SourceFacts = append(d.SourceFacts, f.S("id"))
	}
	d.Syntax = syntax.String()
	out.Diagrams = append(out.Diagrams, d)
	return &out, nil
}

// OverviewContext keeps realization separate from use-case notation. Fact IDs
// and the bundle's full declaration/source registry preserve traceability.
func (v *Views) OverviewContext() string {
	var b strings.Builder
	b.WriteString("## Scope and notation\n\nComplete declared actor/use-case inventory in this source snapshot, including isolated declarations. An association means SDL pursues; no System boundary or UML include/extend is inferred.\n\n")
	b.WriteString("## Functionality realization\n\nFeatures and Functionality are not use cases. Only modeled supports and contributes-to relations are listed.\n\n| Source | Relation | Target | Fact |\n| --- | --- | --- | --- |\n")
	count := 0
	for _, f := range v.Facts {
		if has([]string{"supports", "contributes-to"}, f.S("verb")) {
			fmt.Fprintf(&b, "| %s | %s | %s | %s |\n", f.S("subject"), f.S("verb"), f.S("object"), v.FactReference(f.S("id")))
			count++
		}
	}
	if count == 0 {
		b.WriteString("\nNo functionality realization relations are specified.\n")
	}
	b.WriteString("\nScenario sequences and pre/postconditions are not part of this overview; they are not inferred from static relations.\n\n")
	return b.String()
}
