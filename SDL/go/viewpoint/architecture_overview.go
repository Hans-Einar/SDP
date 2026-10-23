package viewpoint

import (
	"fmt"
	"strings"
)

// architectureOverview is the complete logical inventory, never the roots slice
// or a union of allocation/ownership relations from other viewpoints.
func (v *Views) architectureOverview() *Views {
	out := *v
	out.Diagrams = []Diagram{}
	out.MessageSets = []Fact{}
	out.Facts = append([]Fact{}, v.Relations["contains"]...)
	out.Relations = map[string][]Fact{"contains": out.Facts}
	out.Gaps = filter(v.Gaps, func(f Fact) bool { return f.S("viewpoint") == "VP02" })
	names := append(append(v.names("system"), v.names("container")...), v.names("unit")...)
	out.diagram("VP02-overview", "Complete logical architecture", out.Facts, names...)
	return &out
}

func (v *Views) ArchitectureOverviewContext() string {
	var b strings.Builder
	b.WriteString("## Scope and notation\n\nAll declared Systems, Containers and Units in the checked model are shown, including isolated roots. Arrows mean explicit contains. This is logical decomposition, not deployment, Functionality ownership or runtime communication. No parent or connection is inferred.\n\n")
	names := append(append(v.names("system"), v.names("container")...), v.names("unit")...)
	if len(names) == 0 {
		b.WriteString("No System, Container or Unit is declared in this model.\n\n")
		return b.String()
	}
	b.WriteString("## Structural register\n\n| Component | Kind | Explicit parent | Contains facts |\n| --- | --- | --- | --- |\n")
	for _, name := range names {
		label := name
		for _, d := range v.Model.Declarations {
			if d.Name.Name == name {
				label = v.SpanReference(name, d.Span)
				break
			}
		}
		parents := eq(v.Relations["contains"], "object", name)
		parent := strings.Join(values(parents, "subject"), ", ")
		if parent == "" {
			parent = "Root (no declared parent)"
		}
		proof := v.FactReferences(values(parents, "id"))
		if proof == "" {
			proof = "—"
		}
		fmt.Fprintf(&b, "| %s | %s | %s | %s |\n", label, v.Kinds[name], parent, proof)
	}
	return b.String() + "\n"
}
