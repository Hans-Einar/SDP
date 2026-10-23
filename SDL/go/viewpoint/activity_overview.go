package viewpoint

import "strings"

// activityOverview separates decomposition from prerequisites; neither is a
// schedule. The remaining checked plan facts support source-linked registers.
func (v *Views) activityOverview() *Views {
	out := *v
	out.Diagrams = []Diagram{}
	out.MessageSets = []Fact{}
	functions := values(v.Relations["addresses"], "object")
	out.Facts = filter(v.Facts, func(f Fact) bool {
		return has([]string{"refines", "depends-on", "addresses", "delivers", "illustrates"}, f.S("verb")) ||
			f.S("verb") == "has" && f.S("property") == "implementation-status" && v.Kinds[f.S("subject")] == "activity" ||
			f.S("verb") == "owns" && has(functions, f.S("object"))
	})
	out.Relations = map[string][]Fact{}
	for _, f := range out.Facts {
		out.Relations[f.S("verb")] = append(out.Relations[f.S("verb")], f)
	}
	out.Gaps = filter(v.Gaps, func(f Fact) bool { return f.S("viewpoint") == "VP06" })
	out.diagram("VP06-overview", "All activities and explicit refinement", out.Relations["refines"], v.names("activity")...)
	if len(out.Relations["depends-on"]) > 0 {
		out.diagram("VP06-overview-dependencies", "Explicit activity prerequisites", out.Relations["depends-on"], v.names("activity")...)
		d := &out.Diagrams[len(out.Diagrams)-1]
		d.Syntax = strings.ReplaceAll(d.Mermaid(), " -->|depends-on| ", " -.->|depends-on| ")
	}
	return &out
}
