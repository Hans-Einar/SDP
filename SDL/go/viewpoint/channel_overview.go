package viewpoint

import "sort"

func (v *Views) channelBindings(mode string) []Fact {
	return filter(v.Facts, func(f Fact) bool { return f.S("verb") == "uses" && f.S("mode") == mode })
}
func (v *Views) channelScenarios(mode string) []string {
	if mode == "" {
		return v.names("scenario")
	}
	return values(eq(v.Relations["runs-in"], "object", mode), "subject")
}

// channelOverview is a static, mode-isolated participation topology. It never
// concatenates scenario timelines or invents a direct participant connection.
func (v *Views) channelOverview(mode string) *Views {
	out := *v
	out.Diagrams = []Diagram{}
	out.MessageSets = []Fact{}
	out.Relations = map[string][]Fact{}
	channels := v.names("channel")
	scenarios := v.channelScenarios(mode)
	contracts := values(filter(v.Relations["upholds"], func(f Fact) bool { return has(channels, f.S("subject")) }), "object")
	proof := []string{}
	if mode != "" {
		out.MessageSets = eq(v.MessageSets, "mode", mode)
		for _, m := range out.MessageSets {
			proof = append(proof, m["source_facts"].([]string)...)
		}
	}
	bindings := v.channelBindings(mode)
	out.Facts = filter(v.Facts, func(f Fact) bool {
		return has(proof, f.S("id")) || has(scenarios, f.S("subject")) ||
			f.S("verb") == "upholds" && has(channels, f.S("subject")) ||
			f.S("verb") == "permits" && has(contracts, f.S("subject")) ||
			mode != "" && f.S("verb") == "uses" && f.S("mode") == mode
	})
	for _, f := range out.Facts {
		out.Relations[f.S("verb")] = append(out.Relations[f.S("verb")], f)
	}
	out.Gaps = filter(v.Gaps, func(f Fact) bool { return f.S("viewpoint") == "VP08" && (f.S("mode") == "" || f.S("mode") == mode) })
	if mode == "" || len(channels) == 0 {
		return &out
	}

	// Partition by actual Channel responsibility, not an arbitrary size threshold.
	// Preserve every neighborhood and proof; repeated participants remain the same
	// source identity. This avoids an unrouteable all-channel union graph.
	for _, ch := range channels {
		out.Diagrams = append(out.Diagrams, v.channelDiagram(mode, ch, eq(bindings, "channel", ch)))
	}
	return &out
}

func (v *Views) channelDiagram(mode, channel string, bindings []Fact) Diagram {
	d := Diagram{ID: "VP08-overview-" + mode + "-" + channel, Title: "Channel " + channel + " — mode " + mode, Kind: "flowchart", Edges: []Edge{}, SourceFacts: values(bindings, "id"), Elements: []Fact{}}
	groups := map[string][]Fact{}
	names := []string{channel}
	for _, f := range bindings {
		key := f.S("subject") + "\x00" + f.S("role")
		groups[key] = append(groups[key], f)
		names = append(names, f.S("subject"))
	}
	keys := []string{}
	for key := range groups {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		fs := groups[key]
		f := fs[0]
		a, b := "n_"+f.S("subject"), "n_"+channel
		if f.S("role") == "receiver" {
			a, b = b, a
		}
		d.Edges = append(d.Edges, Edge{a, f.S("role"), b, f.S("id")})
		d.Elements = append(d.Elements, Fact{"participant": f.S("subject"), "channel": channel, "role": f.S("role"), "messages": values(fs, "message"), "proof": values(fs, "id")})
	}
	d.Nodes = v.nodes(names)
	return d
}
