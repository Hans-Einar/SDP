package viewpoint

import (
	"fmt"
	"sort"
	"strings"
)

// responsibilityOverview summarizes only checked owns/realizes/provides paths.
// The graph stays structural; annotations retain the proofs behind each summary.
func (v *Views) responsibilityOverview() *Views {
	out := *v
	out.Diagrams = []Diagram{}
	out.MessageSets = []Fact{}
	out.Facts = filter(v.Facts, func(f Fact) bool { return has([]string{"contains", "owns", "realizes", "provides"}, f.S("verb")) })
	out.Relations = map[string][]Fact{}
	for _, f := range out.Facts {
		out.Relations[f.S("verb")] = append(out.Relations[f.S("verb")], f)
	}
	out.Gaps = filter(v.Gaps, func(f Fact) bool { return f.S("viewpoint") == "VP03" })
	names := append(append(v.names("system"), v.names("container")...), v.names("unit")...)
	elements := []Fact{}
	for i, fn := range v.names("functionality") {
		for _, owner := range eq(out.Relations["owns"], "object", fn) {
			elements = append(elements, Fact{"functionality": fn, "key": fmt.Sprintf("R%02d", i+1), "node": "n_" + owner.S("subject"), "role": "owns", "source_facts": []string{owner.S("id")}})
		}
	}
	for i, cap := range v.names("capability") {
		placed := false
		for _, realization := range eq(out.Relations["realizes"], "object", cap) {
			for _, owner := range eq(out.Relations["owns"], "object", realization.S("subject")) {
				elements = append(elements, Fact{"capability": cap, "key": fmt.Sprintf("C%02d", i+1), "node": "n_" + owner.S("subject"), "role": "contributes", "source_facts": []string{owner.S("id"), realization.S("id")}})
				placed = true
			}
		}
		for _, provider := range eq(out.Relations["provides"], "object", cap) {
			elements = append(elements, Fact{"capability": cap, "key": fmt.Sprintf("C%02d", i+1), "node": "n_" + provider.S("subject"), "role": "provides", "source_facts": []string{provider.S("id")}})
			placed = true
		}
		if !placed {
			names = append(names, cap)
		}
	}
	out.diagram("VP03-overview", "Responsibilities and capabilities across the architecture", out.Relations["contains"], names...)
	if len(out.Diagrams) > 0 {
		d := &out.Diagrams[0]
		d.Elements = elements
		d.SourceFacts = values(out.Facts, "id")
		d.Syntax = responsibilitySyntax(*d)
	}
	return &out
}

func responsibilitySyntax(d Diagram) string {
	var b strings.Builder
	b.WriteString("flowchart LR\n")
	ids := []string{}
	for id := range d.Nodes {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		data := d.Nodes[id].(map[string]any)
		label := fmt.Sprintf("%s (%s)", data["model_id"], data["kind"])
		for _, role := range []string{"owns", "contributes", "provides"} {
			keys := []string{}
			for _, e := range d.Elements {
				if e.S("node") == id && e.S("role") == role {
					keys = append(keys, e.S("key"))
				}
			}
			if len(keys) > 0 {
				label += " | " + role + ": " + strings.Join(unique(keys), " ")
			}
		}
		if data["kind"] == "capability" {
			fmt.Fprintf(&b, "    %s(\"%s\")\n", id, label)
		} else {
			fmt.Fprintf(&b, "    %s[\"%s\"]\n", id, label)
		}
	}
	for _, e := range d.Edges {
		fmt.Fprintf(&b, "    %s -->|%s| %s\n", e.Source, e.Relation, e.Target)
	}
	return b.String()
}
