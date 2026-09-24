package viewpoint

import (
	"fmt"
	"sort"
	"strings"
)

func (v *Views) data() {
	r := v.Relations
	for _, d := range v.names("dataset") {
		sources := eq(r["from"], "object", d)
		families := append(values(sources, "subject"), d)
		holders := eq(r["holds"], "object", d)
		names := values(holders, "subject")
		owners := filter(r["owns"], func(f Fact) bool { return has(names, f.S("object")) })
		contracts := filter(r["upholds"], func(f Fact) bool { return has(families, f.S("subject")) })
		v.diagram("VP09-data-"+d, "Data origin and holder: "+d, join(sources, holders, owners, contracts), d)
	}
	for _, c := range v.names("contract") {
		variants := eq(r["defines"], "subject", c)
		names := append(values(variants, "object"), c)
		fields := filter(r["has-field"], func(f Fact) bool { return has(names, f.S("subject")) })
		v.diagram("VP09-contract-"+c, "Contract structure: "+c, join(variants, fields, eq(r["permits"], "subject", c)), c)
	}
	for _, e := range v.names("encoding") {
		variant := eq(r["encodes"], "subject", e)[0].S("object")
		contract := eq(r["defines"], "object", variant)[0].S("subject")
		entries := filter(v.Facts, func(f Fact) bool { return f.S("verb") == "places" && f.S("subject") == e })
		sort.Slice(entries, func(i, j int) bool { return entries[i].N("offset") < entries[j].N("offset") })
		fields := values(entries, "field")
		context := append(append([]string{}, fields...), e, variant, contract)
		proof := values(filter(v.Facts, func(f Fact) bool { return has(context, f.S("subject")) }), "id")
		lines := []string{"packet"}
		elements := []Fact{}
		for _, f := range entries {
			first, last := f.N("offset"), f.N("offset")+f.N("width")-1
			extent := fmt.Sprint(first)
			if first != last {
				extent = fmt.Sprintf("%d-%d", first, last)
			}
			lines = append(lines, fmt.Sprintf("    %s: \"%s\"", extent, f.S("field")))
			elements = append(elements, Fact{"field": f.S("field"), "first_bit": first, "last_bit": last, "fact": f.S("id")})
		}
		title := fmt.Sprintf("Packet: %s / %s — %s, %s", e, variant, v.prop(e, "byte-order"), v.prop(e, "bit-order"))
		v.Diagrams = append(v.Diagrams, Diagram{"VP10-" + e, title, "packet", v.nodes(fields), []Edge{}, proof, elements, strings.Join(lines, "\n") + "\n"})
	}
	if len(v.names("dataset")) == 0 {
		v.Gaps = append(v.Gaps, Fact{"viewpoint": "VP09", "model_id": "—", "code": "NO_DATASETS", "message": "No Dataset declared."})
	}
	encoded := values(r["encodes"], "object")
	for _, n := range v.names("variant") {
		if !has(encoded, n) {
			v.Gaps = append(v.Gaps, Fact{"viewpoint": "VP10", "model_id": n, "code": "NO_ENCODING", "message": "No fixed wire layout; a logical variant does not define a packet."})
		}
	}
	if len(v.names("encoding")) == 0 {
		v.Gaps = append(v.Gaps, Fact{"viewpoint": "VP10", "model_id": "—", "code": "NO_ENCODINGS", "message": "No Encoding declared."})
	}
}
func (v *Views) prop(s, p string) string {
	for _, f := range v.Facts {
		if f.S("subject") == s && f.S("property") == p {
			return f.S("value")
		}
	}
	return ""
}
