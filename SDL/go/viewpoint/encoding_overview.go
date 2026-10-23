package viewpoint

import "strings"

// encodingOverview keeps independent fixed layouts separate. A logical Variant
// without an Encoding has no packet; declarations alone never imply bit positions.
func (v *Views) encodingOverview() *Views {
	out := *v
	out.Diagrams = []Diagram{}
	out.MessageSets = []Fact{}
	variants := v.names("variant")
	contracts := values(filter(v.Relations["defines"], func(f Fact) bool { return has(variants, f.S("object")) }), "subject")
	shapes := append(append([]string{}, variants...), contracts...)
	fields := values(filter(v.Relations["has-field"], func(f Fact) bool { return has(shapes, f.S("subject")) }), "object")
	context := append(append(shapes, fields...), v.names("encoding")...)
	out.Facts = filter(v.Facts, func(f Fact) bool {
		return has(context, f.S("subject")) || f.S("verb") == "upholds" && v.Kinds[f.S("subject")] == "datagram" && has(contracts, f.S("object"))
	})
	out.Relations = map[string][]Fact{}
	for _, f := range out.Facts {
		out.Relations[f.S("verb")] = append(out.Relations[f.S("verb")], f)
	}
	out.Gaps = filter(v.Gaps, func(f Fact) bool { return f.S("viewpoint") == "VP10" })
	for _, d := range v.Diagrams {
		if strings.HasPrefix(d.ID, "VP10-") {
			d.ID = strings.Replace(d.ID, "VP10-", "VP10-overview-", 1)
			out.Diagrams = append(out.Diagrams, d)
		}
	}
	return &out
}
