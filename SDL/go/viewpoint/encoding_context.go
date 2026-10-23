package viewpoint

import (
	"fmt"
	"sort"
	"strings"
)

func (v *Views) EncodingOverviewContext() string {
	ref := func(name string) string {
		for _, d := range v.Model.Declarations {
			if d.Name.Name == name {
				return v.SpanReference(name, d.Span)
			}
		}
		return name
	}
	var b strings.Builder
	b.WriteString("## Scope and notation\n\nEach packet is one validated fixed Encoding of a Variant in a closed Contract. Bit offsets are logical contract positions; byte-order and bit-order are explicit source declarations. The view does not implement serialization or infer padding, variable lengths, optional-field layout or an executable codec. Header fields belong to the Contract; variant fields belong to the Variant. Independent Encodings are never concatenated.\n\n")
	b.WriteString("## Encoding register\n\n")
	if len(v.names("encoding")) == 0 {
		b.WriteString("No Encodings are declared; no packet layout is inferred.\n\n")
	}
	for _, name := range v.names("encoding") {
		encodes := eq(v.Relations["encodes"], "subject", name)
		variant := encodes[0].S("object")
		defines := eq(v.Relations["defines"], "object", variant)
		contract := defines[0].S("subject")
		props := filter(v.Relations["has"], func(f Fact) bool {
			return f.S("subject") == name || f.S("subject") == contract && f.S("property") == "completeness"
		})
		entries := eq(v.Relations["places"], "subject", name)
		sort.Slice(entries, func(i, j int) bool { return entries[i].N("offset") < entries[j].N("offset") })
		bits := 0
		for _, f := range entries {
			bits += f.N("width")
		}
		fmt.Fprintf(&b, "### Encoding %s\n\nVariant %s; Contract %s (%s); %s, %s; %d bits. Source facts: %s.\n\n", ref(name), ref(variant), ref(contract), v.prop(contract, "completeness"), v.prop(name, "byte-order"), v.prop(name, "bit-order"), bits, v.FactReferences(values(join(encodes, defines, props), "id")))
		families := eq(v.Relations["upholds"], "object", contract)
		labels := []string{}
		for _, f := range families {
			labels = append(labels, ref(f.S("subject")))
		}
		fmt.Fprintf(&b, "Datagram family: %s. Source facts: %s.\n\n", strings.Join(labels, ", "), v.FactReferences(values(families, "id")))
		b.WriteString("| Field | Shape owner | Bit range / width | Value type / presence | Source facts |\n| --- | --- | --- | --- | --- |\n")
		for _, f := range entries {
			field := f.S("field")
			members := filter(v.Relations["has-field"], func(m Fact) bool { return m.S("object") == field && has([]string{contract, variant}, m.S("subject")) })
			owner := members[0].S("subject")
			role := "Variant"
			if owner == contract {
				role = "Header"
			}
			properties := eq(v.Relations["has"], "subject", field)
			fmt.Fprintf(&b, "| %s | %s: %s | %d–%d / %d | %s / %s | %s |\n", ref(field), role, ref(owner), f.N("offset"), f.N("offset")+f.N("width")-1, f.N("width"), v.prop(field, "value-type"), v.prop(field, "presence"), v.FactReferences(values(join([]Fact{f}, members, properties), "id")))
		}
		b.WriteString("\n")
	}
	b.WriteString("## Variant coverage\n\nA logical Variant does not define a wire layout. Missing Encoding is a modeling gap, not evidence about native codec implementation.\n\n")
	if len(v.names("variant")) == 0 {
		b.WriteString("No Variants are declared.\n\n")
	} else {
		b.WriteString("| Variant | Contract | Encodings | Source facts |\n| --- | --- | --- | --- |\n")
		for _, name := range v.names("variant") {
			defines := eq(v.Relations["defines"], "object", name)
			encodes := eq(v.Relations["encodes"], "object", name)
			labels := []string{}
			for _, e := range encodes {
				labels = append(labels, ref(e.S("subject")))
			}
			label := strings.Join(labels, ", ")
			if label == "" {
				label = "No Encoding declared; no fixed wire layout"
			}
			fmt.Fprintf(&b, "| %s | %s | %s | %s |\n", ref(name), ref(defines[0].S("subject")), label, v.FactReferences(values(join(defines, encodes), "id")))
		}
		b.WriteString("\n")
	}
	return b.String()
}
