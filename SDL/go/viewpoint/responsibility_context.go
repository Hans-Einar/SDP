package viewpoint

import (
	"fmt"
	"strings"
)

// responsibilityNames links original declarations, not generated diagram labels.
func (v *Views) responsibilityNames(names []string) string {
	labels := []string{}
	for _, name := range unique(names) {
		label := name
		if v.Model != nil {
			for _, d := range v.Model.Declarations {
				if d.Name.Name == name {
					label = v.SpanReference(name, d.Span)
					break
				}
			}
		}
		labels = append(labels, label)
	}
	if len(labels) == 0 {
		return "Not specified"
	}
	return strings.Join(labels, ", ")
}

func (v *Views) ResponsibilityOverviewContext() string {
	var b strings.Builder
	b.WriteString("## Scope and notation\n\nThe complete logical architecture uses contains edges, including isolated components. R keys identify responsibilities on their immediate owners (owns). C keys identify capabilities: contributes summarizes an owns + realizes path; provides records an explicit offered capability. These annotations are not new SDL relations. Realizes is a contribution, not proof of complete or verified realization. Ownership does not imply provision, and provision does not imply realization. Rounded standalone capabilities have neither a represented contributor nor provider. Keys are sorted within this snapshot. Features and runtime allocation belong to other viewpoints.\n\n")
	if len(v.names("functionality")) == 0 {
		b.WriteString("No Functionality responsibilities are declared.\n\n")
	} else {
		b.WriteString("## Responsibility register\n\n| Key | Functionality | Immediate owner | Capability contributions | Source facts |\n| --- | --- | --- | --- | --- |\n")
		for i, fn := range v.names("functionality") {
			owners := eq(v.Relations["owns"], "object", fn)
			cs := eq(v.Relations["realizes"], "subject", fn)
			proof := append(values(owners, "id"), values(cs, "id")...)
			fmt.Fprintf(&b, "| R%02d | %s | %s | %s | %s |\n", i+1, v.responsibilityNames([]string{fn}), v.responsibilityNames(values(owners, "subject")), v.responsibilityNames(values(cs, "object")), v.FactReferences(proof))
		}
		b.WriteString("\n")
	}
	if len(v.names("capability")) == 0 {
		b.WriteString("No Capabilities are declared. No capability realization or provision is inferred.\n\n")
	} else {
		b.WriteString("## Capability register\n\n| Key | Capability | Explicit providers | Contributing responsibilities | Source facts |\n| --- | --- | --- | --- | --- |\n")
		for i, cap := range v.names("capability") {
			ps := eq(v.Relations["provides"], "object", cap)
			cs := eq(v.Relations["realizes"], "object", cap)
			proof := v.FactReferences(append(values(ps, "id"), values(cs, "id")...))
			if proof == "" {
				proof = "—"
			}
			fmt.Fprintf(&b, "| C%02d | %s | %s | %s | %s |\n", i+1, v.responsibilityNames([]string{cap}), v.responsibilityNames(values(ps, "subject")), v.responsibilityNames(values(cs, "subject")), proof)
		}
		b.WriteString("\nNot specified means no such fact is declared; it does not assert impossibility.\n\n")
	}
	return b.String()
}
