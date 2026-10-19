package viewpoint

import (
	"fmt"
	"strings"
)

func (v *Views) ActivityOverviewContext() string {
	ref := func(name string) string {
		for _, d := range v.Model.Declarations {
			if d.Name.Name == name {
				return v.SpanReference(name, d.Span)
			}
		}
		return name
	}
	refs := func(names []string) string {
		labels := []string{}
		for _, n := range unique(names) {
			labels = append(labels, ref(n))
		}
		return strings.Join(labels, ", ")
	}
	var b strings.Builder
	b.WriteString("## Scope and notation\n\nAll declared Activities are retained, including isolated activities. Refines points from detailed to abstract activity; depends-on points from dependent activity to its explicit prerequisite. Decomposition and prerequisites are separate diagrams. Neither relation inherits status or dependencies. No dates, durations, execution order or calendar schedule are inferred.\n\nImplementation status is a source claim, not test evidence or owner acceptance. Addresses models Functionality coverage; delivers models a Feature contribution, not complete delivery; illustrates links an example Scenario, not an executed test.\n\n")
	activities := v.names("activity")
	if len(activities) == 0 {
		b.WriteString("No Activities are declared. No delivery plan is inferred from Features, files or KanBan.\n\n")
	} else {
		b.WriteString("## Activity register\n\n| Activity | Source-claimed status | Refines | Source facts |\n| --- | --- | --- | --- |\n")
		for _, a := range activities {
			status := filter(v.Facts, func(f Fact) bool { return f.S("subject") == a && f.S("property") == "implementation-status" })
			state := "Unspecified"
			if len(status) > 0 {
				state = status[0].S("value")
			}
			parents := eq(v.Relations["refines"], "subject", a)
			parent := refs(values(parents, "object"))
			if parent == "" {
				parent = "No declared parent"
			}
			proof := v.FactReferences(append(values(status, "id"), values(parents, "id")...))
			if proof == "" {
				proof = "—"
			}
			fmt.Fprintf(&b, "| %s | %s | %s | %s |\n", ref(a), state, parent, proof)
		}
		b.WriteString("\n")
	}
	if len(v.Relations["depends-on"]) == 0 {
		b.WriteString("No activity prerequisites are declared. This does not prove independence.\n\n")
	}
	if len(v.Relations["addresses"]) == 0 {
		b.WriteString("No activity responsibility coverage is declared.\n\n")
	} else {
		b.WriteString("## Addressed responsibilities\n\n| Activity | Functionality | Explicit owner | Source facts |\n| --- | --- | --- | --- |\n")
		for _, f := range v.Relations["addresses"] {
			owners := eq(v.Relations["owns"], "object", f.S("object"))
			proof := append([]string{f.S("id")}, values(owners, "id")...)
			fmt.Fprintf(&b, "| %s | %s | %s | %s |\n", ref(f.S("subject")), ref(f.S("object")), refs(values(owners, "subject")), v.FactReferences(proof))
		}
		b.WriteString("\n")
	}
	if len(v.Relations["delivers"]) == 0 {
		b.WriteString("No activity Feature contributions are declared.\n\n")
	} else {
		b.WriteString("## Planned Feature contributions\n\n| Activity | Feature | Source fact |\n| --- | --- | --- |\n")
		for _, f := range v.Relations["delivers"] {
			fmt.Fprintf(&b, "| %s | %s | %s |\n", ref(f.S("subject")), ref(f.S("object")), v.FactReference(f.S("id")))
		}
		b.WriteString("\n")
	}
	if len(v.Relations["illustrates"]) == 0 {
		b.WriteString("No illustrative Scenarios are linked to activities.\n\n")
	} else {
		b.WriteString("## Illustrative Scenarios\n\n| Activity | Scenario | Source fact |\n| --- | --- | --- |\n")
		for _, f := range v.Relations["illustrates"] {
			fmt.Fprintf(&b, "| %s | %s | %s |\n", ref(f.S("object")), ref(f.S("subject")), v.FactReference(f.S("id")))
		}
		b.WriteString("\n")
	}
	b.WriteString("## Functionality without activity coverage\n\nThis lists missing addresses relations in the model, not unimplemented native code.\n\n")
	addressed := values(v.Relations["addresses"], "object")
	count := 0
	for _, f := range v.names("functionality") {
		if !has(addressed, f) {
			fmt.Fprintf(&b, "- %s\n", ref(f))
			count++
		}
	}
	if count == 0 {
		b.WriteString("No unaddressed Functionality declarations in this snapshot. This is not proof of completed implementation.\n")
	}
	return b.String() + "\n"
}
