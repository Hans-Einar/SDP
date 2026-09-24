package viewpoint

func (v *Views) goalGap(vp, name, code, message, mode, feature string) {
	var m, f any
	if mode != "" {
		m = mode
	}
	if feature != "" {
		f = feature
	}
	v.Gaps = append(v.Gaps, Fact{"viewpoint": vp, "model_id": name, "code": code, "message": message, "mode": m, "feature": f})
}
func (v *Views) goals() {
	r := v.Relations
	if len(v.names("usecase")) == 0 {
		v.goalGap("VP01", "—", "NO_USECASES", "No UseCase is declared.", "", "")
	}
	for _, g := range v.names("usecase") {
		actors, features := eq(r["pursues"], "object", g), eq(r["supports"], "object", g)
		targets := append(values(features, "subject"), g)
		contributions := filter(r["contributes-to"], func(f Fact) bool { return has(targets, f.S("object")) })
		v.diagram("VP01-"+g, "Use case: "+g, join(actors, features, eq(contributions, "object", g)), g)
		if len(actors) == 0 {
			v.goalGap("VP01", g, "NO_ACTOR", "No Actor pursues this use case.", "", "")
		}
		if len(contributions) == 0 {
			v.goalGap("VP01", g, "NO_CONTRIBUTION", "No direct or Feature-mediated Functionality contributions.", "", "")
		}
	}
	for _, a := range v.names("actor") {
		if len(eq(r["pursues"], "subject", a)) == 0 {
			v.goalGap("VP01", a, "NO_GOAL", "Actor has no pursues relation.", "", "")
		}
	}
	if len(v.names("feature")) == 0 {
		v.goalGap("VP07", "—", "NO_FEATURES", "No Feature is declared.", "", "")
	}
	for _, f := range v.names("feature") {
		c := eq(r["contributes-to"], "object", f)
		functions := values(c, "subject")
		owners := filter(r["owns"], func(r Fact) bool { return has(functions, r.S("object")) })
		alloc := filter(r["allocated-to"], func(r Fact) bool { return has(functions, r.S("subject")) })
		v.diagram("VP01-feature-"+f, "Functionality contributions to Feature: "+f, c, f)
		if len(eq(r["supports"], "subject", f)) == 0 {
			v.goalGap("VP01", f, "NO_GOAL", "Feature supports no declared UseCase.", "", "")
		}
		if len(c) == 0 {
			v.goalGap("VP07", f, "NO_CONTRIBUTION", "Feature has no Functionality contributions.", "", "")
		}
		modes := values(alloc, "mode")
		if len(modes) == 0 {
			modes = []string{""}
		}
		for _, m := range modes {
			placed := eq(alloc, "mode", m)
			suffix, title := m, "mode "+m
			if m == "" {
				suffix = "unallocated"
				title = "unallocated"
			}
			v.diagram("VP07-"+f+"-"+suffix, "Feature: "+f+" — "+title, join(c, owners, placed), f)
			for _, fn := range functions {
				if !has(values(placed, "subject"), fn) {
					v.goalGap("VP07", fn, "UNSPECIFIED_ALLOCATION", "Container allocation is unspecified for contributions to "+f+".", m, f)
				}
			}
		}
	}
}
func (v *Views) architecture() {
	r := v.Relations
	children := values(r["contains"], "object")
	roots := []string{}
	for _, n := range append(v.names("unit"), v.names("container")...) {
		if !has(children, n) {
			roots = append(roots, n)
		}
	}
	v.diagram("VP02-roots", "Architecture roots — no connection/allocation is inferred", nil, roots...)
	for _, p := range values(r["contains"], "subject") {
		v.diagram("VP02-"+p, "Logical decomposition: "+p, eq(r["contains"], "subject", p))
	}
	for _, c := range v.names("capability") {
		contrib := eq(r["realizes"], "object", c)
		functions := values(contrib, "subject")
		owners := filter(r["owns"], func(f Fact) bool { return has(functions, f.S("object")) })
		v.diagram("VP03-"+c, "Contributions to capability: "+c, join(contrib, owners), c)
		offers := eq(r["provides"], "object", c)
		if len(offers) > 0 {
			v.diagram("VP03-"+c+"-offers", "Capability providers: "+c, offers, c)
		}
	}
	for _, m := range v.names("mode") {
		v.diagram("VP05-"+m, "Required ports in mode: "+m, filter(v.Facts, func(f Fact) bool { return f.S("verb") == "requires" && f.S("mode") == m }))
	}
}
func (v *Views) plans() {
	r := v.Relations
	children := values(r["refines"], "subject")
	roots := []string{}
	for _, a := range v.names("activity") {
		if !has(children, a) {
			roots = append(roots, a)
		}
	}
	v.diagram("VP06-roots", "Activity roots", nil, roots...)
	for _, p := range values(r["refines"], "object") {
		c := eq(r["refines"], "object", p)
		as := append(values(c, "subject"), p)
		d := filter(r["delivers"], func(f Fact) bool { return has(as, f.S("subject")) })
		v.diagram("VP06-detail-"+p, "Activity decomposition: "+p, join(c, d))
	}
	for _, p := range values(r["addresses"], "subject") {
		a := eq(r["addresses"], "subject", p)
		fn := values(a, "object")
		o := filter(r["owns"], func(f Fact) bool { return has(fn, f.S("object")) })
		v.diagram("VP06-work-"+p, "Planned responsibilities: "+p, join(a, o))
	}
	v.diagram("VP06-dependencies", "Explicit activity dependencies", r["depends-on"])
}
