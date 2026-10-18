package viewpoint

import (
	"fmt"
	"sort"
	"strings"
)

// featureOverview keeps complete logical decomposition and only the requested
// mode's explicit allocation. The empty mode is logical, never a union of modes.
func (v *Views) featureOverview(q Query) (*Views, error) {
	out := *v
	out.Diagrams = []Diagram{}
	out.MessageSets = []Fact{}
	out.Facts = []Fact{}
	out.Relations = map[string][]Fact{}
	out.Gaps = filter(v.Gaps, func(f Fact) bool { return f.S("viewpoint") == "VP07" && f.S("code") != "UNSPECIFIED_ALLOCATION" })
	contributions := filter(v.Relations["contributes-to"], func(f Fact) bool { return v.Kinds[f.S("object")] == "feature" })
	functions := values(contributions, "subject")
	names := append(append(append(append(v.names("system"), v.names("container")...), v.names("unit")...), v.names("feature")...), functions...)
	for _, f := range v.Facts {
		a, b, verb := f.S("subject"), f.S("object"), f.S("verb")
		keep := verb == "contains" && has(names, a) && has(names, b)
		keep = keep || verb == "owns" && has(functions, b)
		keep = keep || verb == "contributes-to" && v.Kinds[b] == "feature"
		keep = keep || verb == "allocated-to" && q.Mode != "" && f.S("mode") == q.Mode && has(functions, a)
		if keep {
			out.Facts = append(out.Facts, f)
			out.Relations[verb] = append(out.Relations[verb], f)
		}
	}
	if q.Mode != "" {
		for _, f := range contributions {
			if len(eq(out.Relations["allocated-to"], "subject", f.S("subject"))) == 0 {
				out.goalGap("VP07", f.S("subject"), "UNSPECIFIED_ALLOCATION", "Container allocation is unspecified for contributions to "+f.S("object")+".", q.Mode, f.S("object"))
			}
		}
	}
	title := "Complete logical architecture and feature realization"
	if q.Mode != "" {
		title += " — mode " + q.Mode
	}
	// The complete structural base stays readable: feature keys annotate actual
	// owners/allocations; function paths and their proof remain in the companion table.
	structural := append(append(v.names("system"), v.names("container")...), v.names("unit")...)
	features := v.names("feature")
	elements := []Fact{}
	for i, feature := range features {
		placed := false
		for _, c := range eq(out.Relations["contributes-to"], "object", feature) {
			for _, role := range []string{"owns", "allocated-to"} {
				for _, f := range out.Relations[role] {
					endpoint := f.S("subject")
					matches := role == "owns" && f.S("object") == c.S("subject")
					if role == "allocated-to" {
						matches = f.S("subject") == c.S("subject")
						endpoint = f.S("object")
					}
					if matches {
						placed = true
						elements = append(elements, Fact{"feature": feature, "key": fmt.Sprintf("F%02d", i+1), "node": "n_" + endpoint, "role": role, "source_facts": []string{c.S("id"), f.S("id")}})
					}
				}
			}
		}
		if !placed {
			structural = append(structural, feature)
		}
	}
	out.diagram("VP07-overview", title, out.Relations["contains"], structural...)
	if len(out.Diagrams) == 0 {
		return &out, nil
	}
	d := &out.Diagrams[0]
	d.Elements = elements
	d.SourceFacts = values(out.Facts, "id")
	if len(d.Nodes) > 128 || len(d.Edges) > 320 {
		return nil, fmt.Errorf("overview-limit: complete VP07 exceeds flowchart profile; no partial overview generated")
	}
	d.Syntax = featureSyntax(*d)
	return &out, nil
}

func featureSyntax(d Diagram) string {
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
		for _, role := range []string{"owns", "allocated-to"} {
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
		switch data["kind"] {
		case "feature":
			fmt.Fprintf(&b, "    %s(\"%s\")\n", id, label)
		default:
			fmt.Fprintf(&b, "    %s[\"%s\"]\n", id, label)
		}
	}
	for _, e := range d.Edges {
		fmt.Fprintf(&b, "    %s -->|%s| %s\n", e.Source, e.Relation, e.Target)
	}
	return b.String()
}

func (v *Views) FeatureOverviewContext(mode string) string {
	var b strings.Builder
	b.WriteString("## Scope and notation\n\nAll declared Systems, containers and units form the logical architecture base, including isolated declarations. Rectangles show the complete structural base with contains edges. Feature keys annotate logical owners (owns) and selected-mode containers (allocated-to), derived only through explicit Functionality contribution paths. Rounded standalone Features have no represented placement. These annotations are not new direct SDL relations. This is decomposition and responsibility, not an inferred deployment or call graph.\n\n")
	if mode == "" {
		b.WriteString("Logical scope: allocation annotations are intentionally absent. Select a named mode under VP07 in the navigation tree to inspect its explicit placement. No allocations from different modes are combined.\n\n")
	} else {
		fmt.Fprintf(&b, "Selected mode: **%s**. Only this mode's explicit allocated-to facts inform placement annotations. Ownership does not imply allocation; missing placement is reported below.\n\n", mode)
	}
	b.WriteString("## Feature legend\n\n| Key | Feature |\n| --- | --- |\n")
	for i, f := range v.names("feature") {
		fmt.Fprintf(&b, "| F%02d | %s |\n", i+1, f)
	}
	b.WriteString("\n")
	b.WriteString("## Feature contributions\n\n| Feature | Functionality | Logical owner | Allocation in selected scope | Source facts |\n| --- | --- | --- | --- | --- |\n")
	for _, feature := range v.names("feature") {
		cs := eq(v.Relations["contributes-to"], "object", feature)
		if len(cs) == 0 {
			fmt.Fprintf(&b, "| %s | Not specified | — | — | — |\n", feature)
		}
		for _, c := range cs {
			fn := c.S("subject")
			owners := strings.Join(values(eq(v.Relations["owns"], "object", fn), "subject"), ", ")
			placed := "Not selected (logical scope)"
			if mode != "" {
				placed = strings.Join(values(eq(v.Relations["allocated-to"], "subject", fn), "object"), ", ")
				if placed == "" {
					placed = "Unspecified"
				}
			}
			proof := []string{c.S("id")}
			proof = append(proof, values(eq(v.Relations["owns"], "object", fn), "id")...)
			proof = append(proof, values(eq(v.Relations["allocated-to"], "subject", fn), "id")...)
			fmt.Fprintf(&b, "| %s | %s | %s | %s | %s |\n", feature, fn, owners, placed, v.FactReferences(proof))
		}
	}
	return b.String() + "\n"
}
