package viewpoint

import "strings"

// dataOverview combines dataset neighborhoods with full inventories. Contract
// structure remains a separately selectable detail, never an all-fields union.
func (v *Views) dataOverview() *Views {
	out := *v
	out.Diagrams = []Diagram{}
	out.MessageSets = []Fact{}
	projectors := values(eq(v.Facts, "verb", "projects"), "subject")
	out.Facts = filter(v.Facts, func(f Fact) bool {
		kind := v.Kinds[f.S("subject")]
		return has([]string{"from", "holds", "defines", "has-field", "permits", "projects"}, f.S("verb")) ||
			f.S("verb") == "upholds" && has([]string{"dataset", "datagram"}, kind) ||
			f.S("verb") == "owns" && (v.Kinds[f.S("object")] == "database" || has(projectors, f.S("object"))) ||
			f.S("verb") == "has" && has([]string{"contract", "field"}, kind)
	})
	out.Relations = map[string][]Fact{}
	for _, f := range out.Facts {
		out.Relations[f.S("verb")] = append(out.Relations[f.S("verb")], f)
	}
	out.Gaps = filter(v.Gaps, func(f Fact) bool { return f.S("viewpoint") == "VP09" })
	for _, d := range v.Diagrams {
		if strings.HasPrefix(d.ID, "VP09-data-") {
			d.ID = strings.Replace(d.ID, "VP09-data-", "VP09-overview-data-", 1)
			out.Diagrams = append(out.Diagrams, d)
		}
	}
	return &out
}
