package viewpoint

import (
	"fmt"
	"sort"
	"strings"
)

// interfaceOverview includes every Interface, even without consumption. Only
// consumes facts connect consumers; no provider or runtime relationship is guessed.
func (v *Views) interfaceOverview() *Views {
	out := *v
	out.Diagrams = []Diagram{}
	out.MessageSets = []Fact{}
	out.Facts = append([]Fact{}, v.Relations["consumes"]...)
	out.Relations = map[string][]Fact{"consumes": out.Facts}
	out.Gaps = filter(v.Gaps, func(f Fact) bool { return f.S("viewpoint") == "VP04" })
	out.diagram("VP04-overview", "All declared interfaces and explicit consumers", out.Facts, v.names("interface")...)
	if len(out.Diagrams) > 0 {
		out.Diagrams[0].Syntax = interfaceSyntax(out.Diagrams[0])
	}
	return &out
}

func interfaceSyntax(d Diagram) string {
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
		if data["kind"] == "interface" {
			fmt.Fprintf(&b, "    %s(\"%s\")\n", id, label)
		} else {
			fmt.Fprintf(&b, "    %s[\"%s\"]\n", id, label)
		}
	}
	for _, e := range d.Edges {
		fmt.Fprintf(&b, "    %s -.->|consumes| %s\n", e.Source, e.Target)
	}
	return b.String()
}

func (v *Views) InterfaceOverviewContext() string {
	var b strings.Builder
	b.WriteString("## Scope and notation\n\nEvery declared Interface is included, even without a consumer. Rectangles are explicit Unit/Container consumers; rounded nodes are interfaces. Dashed consumes arrows point from consumer to interface and mean designed use, not an observed call or mandatory availability. No provider, Channel, transport, call order or mode necessity is inferred. Units without consumes facts are outside this projection; VP02 shows the complete logical structure.\n\n")
	if len(v.names("interface")) == 0 {
		b.WriteString("No Interfaces are declared.\n\n")
		return b.String()
	}
	if len(v.Relations["consumes"]) == 0 {
		b.WriteString("No interface consumption is declared.\n\n")
	}
	references := func(names []string) string {
		labels := []string{}
		for _, name := range unique(names) {
			label := name
			for _, d := range v.Model.Declarations {
				if d.Name.Name == name {
					label = v.SpanReference(name, d.Span)
					break
				}
			}
			labels = append(labels, label)
		}
		return strings.Join(labels, ", ")
	}
	b.WriteString("## Interface register\n\n| Interface | Explicit consumers | Source facts |\n| --- | --- | --- |\n")
	for _, name := range v.names("interface") {
		uses := eq(v.Relations["consumes"], "object", name)
		consumers := references(values(uses, "subject"))
		if consumers == "" {
			consumers = "No declared consumer"
		}
		proof := v.FactReferences(values(uses, "id"))
		if proof == "" {
			proof = "—"
		}
		fmt.Fprintf(&b, "| %s | %s | %s |\n", references([]string{name}), consumers, proof)
	}
	return b.String() + "\nNo declared consumer means no consumes fact is present; it does not prove an unused or unnecessary interface. Provider information is not expressed by consumes.\n\n"
}
