package viewpoint

import (
	"fmt"
	"strings"
)

func (v *Views) DiagramMarkdown(d Diagram, assetPrefix string, rendered bool) string {
	s := "## " + d.Title + "\n\n"
	if rendered {
		s += "![" + d.Title + "](" + assetPrefix + d.ID + ".svg)\n\n"
	} else {
		s += "```mermaid\n" + d.Mermaid() + "```\n\n"
	}
	s += "Source facts: " + strings.Join(d.SourceFacts, ", ") + ".\n\n"
	if len(d.SourceFacts) == 0 {
		s += "Declarations only.\n\n"
	}
	return s
}
func (v *Views) Tables(vp string) string {
	var b strings.Builder
	if vp == "VP04" {
		b.WriteString("## Interface use\n\nNo Channel or provider is inferred.\n\n| Consumer | Interface | Fact | Line |\n| --- | --- | --- | --- |\n")
		for _, r := range v.Relations["consumes"] {
			fmt.Fprintf(&b, "| %s | %s | %s | %d |\n", r.S("subject"), r.S("object"), r.S("id"), r.N("line"))
		}
	}
	if vp == "VP09" {
		b.WriteString("## Fields and contract properties\n\nDatabase means a persistent data source, not necessarily SQL.\n\n| Fact | ID |\n| --- | --- |\n")
		for _, r := range v.Facts {
			if r.S("verb") == "has" && has([]string{"field", "contract", "encoding"}, v.Kinds[r.S("subject")]) {
				fmt.Fprintf(&b, "| %s | %s |\n", r.S("text"), r.S("id"))
			}
		}
		b.WriteString("\n## Projection responsibilities\n\n| Functionality | Dataset | Datagram | Fact |\n| --- | --- | --- | --- |\n")
		for _, r := range eq(v.Facts, "verb", "projects") {
			fmt.Fprintf(&b, "| %s | %s | %s | %s |\n", r.S("subject"), r.S("dataset"), r.S("datagram"), r.S("id"))
		}
	}
	if vp == "VP08" {
		b.WriteString("## Derived MessageSet\n\n| Channel | Mode | Datagram | Sender | Receiver | Source IDs |\n| --- | --- | --- | --- | --- | --- |\n")
		for _, r := range v.MessageSets {
			fmt.Fprintf(&b, "| %s | %s | %s | %s | %s | %s |\n", r.S("channel"), r.S("mode"), r.S("message"), strings.Join(r["senders"].([]string), ", "), strings.Join(r["receivers"].([]string), ", "), strings.Join(r["source_facts"].([]string), ", "))
		}
	}
	if vp == "VP11" {
		b.WriteString("## Complete fact registry\n\n| ID | Statement | Line |\n| --- | --- | --- |\n")
		for _, r := range v.Facts {
			fmt.Fprintf(&b, "| %s | %s | %d |\n", r.S("id"), r.S("text"), r.N("line"))
		}
		b.WriteString("\n## Declarations\n\n| ID | Type | Line |\n| --- | --- | --- |\n")
		for _, d := range v.Model.Declarations {
			fmt.Fprintf(&b, "| %s | %s | %d |\n", d.Name.Name, d.Kind, d.Span.Line)
		}
	}
	for _, g := range v.Gaps {
		if g.S("viewpoint") == vp {
			fmt.Fprintf(&b, "\nModel gap %s: %s (%s; mode=%s).\n", g.S("code"), g.S("message"), g.S("model_id"), g.S("mode"))
		}
	}
	return b.String() + "\n"
}
func (v *Views) Implementation() string {
	var b strings.Builder
	b.WriteString("# Generated implementation plan\n\nStatus is a model claim, not execution evidence.\n\n[Navigator](navigator.md)\n\n")
	addressed := values(v.Relations["addresses"], "object")
	for _, a := range v.names("activity") {
		fmt.Fprintf(&b, "## %s\n\nStatus: %s.\n\n", a, or(v.prop(a, "implementation-status"), "unspecified"))
		b.WriteString("| Fact | Source ID |\n| --- | --- |\n")
		for _, r := range v.Facts {
			if r.S("subject") == a || r.S("object") == a {
				fmt.Fprintf(&b, "| %s | %s |\n", r.S("text"), r.S("id"))
			}
		}
		for _, r := range eq(v.Relations["addresses"], "subject", a) {
			for _, owner := range eq(v.Relations["owns"], "object", r.S("object")) {
				fmt.Fprintf(&b, "| %s | %s |\n", owner.S("text"), owner.S("id"))
			}
		}
		b.WriteString("\n")
	}
	b.WriteString("## Unaddressed model responsibilities\n\n")
	for _, n := range v.names("functionality") {
		if !has(addressed, n) {
			b.WriteString("- " + n + "\n")
		}
	}
	return b.String()
}
func or(s, f string) string {
	if s == "" {
		return f
	}
	return s
}
