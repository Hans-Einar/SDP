package viewpoint

import (
	"fmt"
	"strings"
)

// Dependencies are checked Dependency facts, not ordinary binary Relations.
func (v *Views) modeRequirements(mode string) []Fact {
	return filter(v.Facts, func(f Fact) bool { return f.S("verb") == "requires" && f.S("mode") == mode })
}
func (v *Views) modeOverview(mode string) *Views {
	out := *v
	out.Diagrams = []Diagram{}
	out.MessageSets = []Fact{}
	out.Facts = []Fact{}
	out.Relations = map[string][]Fact{}
	out.Gaps = filter(v.Gaps, func(f Fact) bool { return f.S("viewpoint") == "VP05" && (mode == "" || f.S("mode") == mode) })
	// A missing mode means an index only, never a union or an implicit default.
	if mode == "" {
		return &out
	}
	out.Facts = v.modeRequirements(mode)
	out.Relations["requires"] = out.Facts
	out.diagram("VP05-overview-"+mode, "Required interfaces — mode "+mode, out.Facts, v.names("capability")...)
	if len(out.Diagrams) > 0 {
		d := &out.Diagrams[0]
		d.Syntax = strings.ReplaceAll(d.Mermaid(), " -->|requires| ", " -.->|requires| ")
	}
	return &out
}

// ModeOverviewContext reads the original checked inventory to count modes, but
// emits dependency facts/register rows only after an explicit mode selection.
func (v *Views) ModeOverviewContext(mode string) string {
	ref := func(name string) string {
		for _, d := range v.Model.Declarations {
			if d.Name.Name == name {
				return v.SpanReference(name, d.Span)
			}
		}
		return name
	}
	var b strings.Builder
	b.WriteString("## Scope and notation\n\nRequires states that an Interface is necessary for a Capability in the selected Mode. It does not establish sufficiency, a provider, consumption, transport or runtime execution. Modes have no implicit inheritance. Dashed requires arrows point from Capability to Interface.\n\n")
	if mode == "" {
		b.WriteString("## Choose a mode\n\nSelect a named mode under VP05 in the navigation tree. No mode is selected and no dependency graph is generated; dependencies from different modes are never combined.\n\n")
		if len(v.names("mode")) == 0 {
			return b.String() + "No Modes are declared.\n\n"
		}
		b.WriteString("| Mode | Declared requirements |\n| --- | --- |\n")
		for _, m := range v.names("mode") {
			fmt.Fprintf(&b, "| %s | %d |\n", ref(m), len(v.modeRequirements(m)))
		}
		return b.String() + "\nA zero count means no requires fact is declared for that mode, not proof of independence.\n\n"
	}
	fmt.Fprintf(&b, "Selected mode: **%s**. Mode source: %s.\n\n", mode, ref(mode))
	fs := v.modeRequirements(mode)
	if len(fs) == 0 {
		b.WriteString("No required interfaces are declared for this mode. This does not prove independence or sufficiency.\n\n")
	}
	if len(v.names("capability")) == 0 {
		return b.String() + "No Capabilities are declared. No dependency is inferred from Features or consumers.\n\n"
	}
	b.WriteString("## Capability requirements\n\n| Capability | Required interface in selected mode | Source facts |\n| --- | --- | --- |\n")
	for _, c := range v.names("capability") {
		reqs := eq(fs, "subject", c)
		if len(reqs) == 0 {
			fmt.Fprintf(&b, "| %s | No declared requirement | — |\n", ref(c))
		}
		for _, f := range reqs {
			fmt.Fprintf(&b, "| %s | %s | %s |\n", ref(c), ref(f.S("object")), v.FactReference(f.S("id")))
		}
	}
	return b.String() + "\nCapabilities without a declared requirement remain visible; no claim of completeness is inferred. Interfaces not required in this mode are outside this projection.\n\n"
}
