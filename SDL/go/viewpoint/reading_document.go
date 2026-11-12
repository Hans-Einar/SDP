package viewpoint

import (
	"fmt"
	"sort"
	"strings"
)

var readingQuestions = map[string]string{
	"VP01": "Who wants what from the system?", "VP02": "What parts exist and how are they organized?", "VP03": "Who is responsible for which behavior?", "VP04": "Which interfaces are used, and by whom?", "VP05": "What is required in each operating mode?", "VP06": "What work delivers the design, and what depends on it?", "VP07": "Where and how are features realized?", "VP08": "Who communicates, what is exchanged, and in what order?", "VP09": "What data exists, who holds it, and how is it represented?", "VP10": "How are datagrams encoded?", "VP11": "What is specified, traceable, or missing?",
}

// ReadingMarkdown shares object chapters between collection and focused pages.
// The caller renders selected.Diagrams under diagrams/ before publication.
func (v *Views) ReadingMarkdown(q Query, selected *Views, rendered bool, supplements ...string) string {
	spec, _ := SpecFor(q.Viewpoint)
	var b strings.Builder
	fmt.Fprintf(&b, "# %s — %s\n\n%s\n\n%s\n\n", q.Viewpoint, spec.Title, readingQuestions[q.Viewpoint], spec.Note)
	if q.Mode != "" {
		fmt.Fprintf(&b, "Selected mode: **%s**. Global definitions remain visible; mode applicability is stated where declared. Modes have no implicit inheritance.\n\n", q.Mode)
	}
	modes := v.names("mode")
	if len(modes) > 0 && v.ReadingProject != "" {
		b.WriteString("Operating mode: ")
		all := q
		all.Mode = ""
		all.Page = 0
		b.WriteString(v.readingLink("All declared modes", all))
		for _, name := range modes {
			selectedMode := q
			selectedMode.Mode = name
			selectedMode.Page = 0
			b.WriteString(" · " + v.readingLink(name, selectedMode))
		}
		b.WriteString("\n\n")
	}
	landing := q
	landing.Collection = ""
	landing.Register = ""
	landing.Focus = ""
	landing.Page = 0
	if q.Collection != "" || q.Focus != "" || q.Register != "" {
		fmt.Fprintf(&b, "%s\n\n", v.readingLink("Viewpoint overview", landing))
	}
	if q.Register != "" {
		v.readingRegister(&b, q)
		return b.String()
	}
	supplement := strings.Join(supplements, "\n")
	names := v.readingNames(q)
	if q.Focus == "" && q.Collection == "" {
		fmt.Fprintf(&b, "This page describes the checked model: %d declarations and %d source statements. Model claims do not establish implemented behavior or execution evidence.\n\n", len(v.Model.Declarations), len(v.Facts))
		b.WriteString(supplement)
		for _, d := range selected.Diagrams {
			b.WriteString(v.DiagramMarkdown(d, "diagrams/", rendered))
		}
		if len(selected.Diagrams) == 0 && supplement == "" {
			b.WriteString("The reading scope remains available without a combined diagram. A missing diagram does not establish absence of declarations or relationships.\n\n")
		}
		b.WriteString("## Reading chapters\n\n")
		nonempty := 0
		for _, c := range v.ReadingCollections(q.Viewpoint) {
			if len(c.Names) == 0 {
				continue
			}
			nonempty++
			cq := landing
			cq.Collection = c.ID
			fmt.Fprintf(&b, "- %s — %d objects.\n", v.readingLink(c.Label, cq), len(c.Names))
			limit := len(c.Names)
			if limit > 3 {
				limit = 3
			}
			for _, n := range c.Names[:limit] {
				fmt.Fprintf(&b, "  %s.\n", v.readingObjectLink(n, q))
			}
		}
		if nonempty == 0 {
			b.WriteString("No relevant declarations are present for this viewpoint. This concerns the model, not runtime impossibility.\n")
		}
		b.WriteString("\n")
		if q.Viewpoint == "VP11" {
			fmt.Fprintf(&b, "## Coverage and dedicated registers\n\n%d explicit property statements; %d reported model gaps. Gap rules report modeling omissions, not a complete requirements audit.\n\n", len(eq(v.Facts, "verb", "has")), len(v.Gaps))
			for _, reg := range []struct{ id, label string }{{"facts", "Complete fact registry"}, {"declarations", "Declarations"}} {
				rq := landing
				rq.Register = reg.id
				fmt.Fprintf(&b, "- %s\n", v.readingLink(reg.label, rq))
			}
			b.WriteString("\n")
		}
		v.readingGaps(&b, selected.Gaps)
		b.WriteString(readingLimits(q.Viewpoint))
		return b.String()
	}
	if q.Collection != "" {
		label := q.Collection
		for _, c := range v.ReadingCollections(q.Viewpoint) {
			if c.ID == q.Collection {
				label = c.Label
			}
		}
		fmt.Fprintf(&b, "## %s\n\n", label)
		v.readingPagination(&b, q)
		b.WriteString("### Contents\n\n")
		for _, n := range names {
			fmt.Fprintf(&b, "- %s (%s)\n", v.readingObjectLink(n, q), v.Kinds[n])
		}
		b.WriteString("\nEach chapter uses the same model content as its object detail. Definitions without relationships remain visible. Large diagram sets are deferred to object details; all chapter statements and object links remain present.\n\n")
		if len(names) == 0 {
			b.WriteString("No declarations belong to this collection.\n\n")
		}
	}
	if q.Focus != "" {
		b.WriteString(supplement)
	}
	for _, n := range names {
		v.readingChapter(&b, q, n, selected, rendered)
	}
	if q.Collection != "" {
		v.readingPagination(&b, q)
	}
	return b.String()
}

func (v *Views) readingPagination(b *strings.Builder, q Query) {
	a, z, total, _ := v.readingPage(q)
	page := q.Page
	if page == 0 {
		page = 1
	}
	first := a + 1
	if total == 0 {
		first = 0
	}
	fmt.Fprintf(b, "Showing %d–%d of %d; page %d. The complete scope is retained across pages.\n\n", first, z, total, page)
	size := ReadingChapterPageSize
	if q.Register != "" {
		size = ReadingRegisterPageSize
	}
	pages := (total + size - 1) / size
	links := []string{}
	if page > 1 {
		prev := q
		prev.Page = page - 1
		links = append(links, v.readingLink("Previous page", prev))
	}
	if page < pages {
		next := q
		next.Page = page + 1
		links = append(links, v.readingLink("Next page", next))
	}
	if len(links) > 0 {
		fmt.Fprintf(b, "%s\n\n", strings.Join(links, " · "))
	}
}

func (v *Views) readingRegister(b *strings.Builder, q Query) {
	title := "Declarations"
	if q.Register == "facts" {
		title = "Complete fact registry"
	}
	fmt.Fprintf(b, "## %s\n\nThe complete checked model is the register scope. Selected mode is retained for navigation; global statements and declarations are not filtered. Source links retain exact original spans.\n\n", title)
	v.readingPagination(b, q)
	a, z, _, _ := v.readingPage(q)
	if q.Register == "facts" {
		b.WriteString("| Statement and source location |\n| --- |\n")
		for _, f := range v.Facts[a:z] {
			fmt.Fprintf(b, "| %s |\n", v.FactReference(f.S("id")))
		}
	} else {
		b.WriteString("| Name | Kind | Source location |\n| --- | --- | --- |\n")
		for _, d := range v.Model.Declarations[a:z] {
			fmt.Fprintf(b, "| %s | %s | %s |\n", v.readingObjectLink(d.Name.Name, q), d.Kind, v.SpanReference(traceabilityPosition(v.sourceLabel(d.Span.Source), d.Span.Line), d.Span))
		}
	}
	b.WriteString("\n")
	v.readingPagination(b, q)
}

func (v *Views) readingGaps(b *strings.Builder, gaps []Fact) {
	if len(gaps) == 0 {
		return
	}
	b.WriteString("## Reported model gaps\n\nThese observations concern the declared model. They do not establish missing implementation or runtime behavior.\n\n")
	for _, g := range gaps {
		context := g.S("model_id")
		if v.Kinds[context] == "" {
			context = "Model scope"
		}
		if g.S("mode") != "" {
			context += "; mode " + g.S("mode")
		}
		fmt.Fprintf(b, "- %s: %s\n", context, g.S("message"))
	}
	b.WriteString("\n")
}

func readingLimits(vp string) string {
	text := map[string]string{
		"VP01": "Preconditions and postconditions are unspecified by this profile. Only explicitly exercises-linked scenarios establish scenario context; static relations do not establish an execution order.",
		"VP02": "Contains describes logical structure. Deployment, scheduling, runtime calls, and unnamed internal units are not inferred. Channel endpoints refer to actual declared participants.",
		"VP03": "Owns assigns responsibility; realizes contributes to a capability; provides declares an offer. Missing relationships remain unspecified rather than being inferred from neighboring objects.",
		"VP04": "Consumes declares use of an interface. It does not identify a provider or imply a channel, transport, or runtime call.",
		"VP05": "Requires is scoped to its explicit mode. Requirements in one mode do not imply inheritance or defaults in another mode.",
		"VP06": "Delivers targets are a derived presentation group. Dependency order does not establish execution, completion, or test evidence. Implementation status is a source assertion.",
		"VP07": "Feature contributions, responsibility, containment, and mode allocation are separate relations. A missing allocation remains unspecified.",
		"VP08": "Participation topology is static. Only numbered steps within one scenario establish an order. Contracts declare allowed payloads; participation does not prove delivery. No transport, precondition, or runtime scheduling is inferred.",
		"VP09": "Dataset and Datagram definitions describe logical data. Database denotes a persistent source, not a storage engine or SQL schema. Optional fields do not imply default values, and projections do not establish executable conversion.",
		"VP10": "Logical shapes do not establish wire layout. Only explicit encoding placements define bit ranges; byte order and bit order remain source claims. Missing placements are unspecified rather than inferred padding.",
		"VP11": "Properties and status are source assertions. Traceability is not conformance evidence, and absence of reported gaps does not prove completeness.",
	}
	return "## Interpretation and unspecified behavior\n\n" + text[vp] + "\n\n"
}

func diagramForName(d Diagram, vp, name string) bool {
	p := strings.TrimPrefix(d.ID, vp+"-")
	return strings.HasPrefix(d.ID, vp+"-") && (p == name || p == "contract-"+name || p == "data-"+name || p == "feature-"+name || p == "detail-"+name || p == "work-"+name || strings.HasPrefix(p, name+"-"))
}

func (v *Views) readingChapter(b *strings.Builder, q Query, name string, selected *Views, rendered bool) {
	kind := v.Kinds[name]
	fmt.Fprintf(b, "## %s\n\n%s declaration: %s.\n\n", name, kind, v.readingDeclaration(name))
	fs := v.readingFacts(q, []string{name})
	switch kind {
	case "contract":
		v.readingContract(b, q, name)
	case "channel":
		v.readingChannel(b, q, name)
	case "scenario":
		v.readingScenario(b, q, name)
	case "encoding":
		v.readingEncoding(b, q, name)
	}
	if has([]string{"dataset", "datagram", "message"}, kind) {
		fmt.Fprintf(b, "Declared holder: %s. Source Dataset: %s. Payload/shape contract: %s.\n\n", v.readingNamesText(values(eq(v.Relations["holds"], "object", name), "subject"), q), v.readingNamesText(values(eq(v.Relations["from"], "subject", name), "object"), q), v.readingNamesText(values(eq(v.Relations["upholds"], "subject", name), "object"), q))
		if kind == "message" {
			fmt.Fprintf(b, "Message kind: %s. Replies to: %s.\n\n", or(v.prop(name, "message-kind"), "unspecified"), v.readingNamesText(values(eq(v.Relations["replies-to"], "subject", name), "object"), q))
		}
	}
	properties := filter(fs, func(f Fact) bool { return f.S("subject") == name && f.S("verb") == "has" })
	if len(properties) > 0 {
		b.WriteString("### Explicit properties\n\n| Property | Value | Source statement |\n| --- | --- | --- |\n")
		for _, f := range properties {
			fmt.Fprintf(b, "| %s | %s | %s |\n", readingEscape(f.S("property")), readingEscape(fmt.Sprint(f["value"])), v.FactReference(f.S("id")))
		}
		b.WriteString("\n")
	}
	for _, d := range selected.Diagrams {
		relevant := diagramForName(d, q.Viewpoint, name)
		if q.Viewpoint == "VP01" && d.Kind == "sequence" {
			for _, f := range eq(v.Relations["exercises"], "object", name) {
				if d.ID == "VP08-"+f.S("subject") {
					relevant = true
				}
			}
		}
		if relevant {
			b.WriteString(v.DiagramMarkdown(d, "diagrams/", rendered))
		}
	}
	relationships := filter(fs, func(f Fact) bool { return !has([]string{"has", "step", "places"}, f.S("verb")) })
	if len(relationships) == 0 {
		b.WriteString("No relationships are declared in this viewpoint and mode scope. The declaration remains available; no missing relationship is invented.\n\n")
	} else {
		b.WriteString("### Declared relationships\n\n")
		for _, f := range sortedFacts(relationships) {
			fmt.Fprintf(b, "- %s\n", v.FactReference(f.S("id")))
		}
		b.WriteString("\n")
	}
	related := []string{}
	for _, f := range fs {
		for _, key := range []string{"subject", "object", "dataset", "datagram", "channel", "message", "sender", "receiver", "field", "variant", "mode"} {
			n := f.S(key)
			if n != "" && n != name && v.Kinds[n] != "" {
				related = append(related, n)
			}
		}
	}
	if len(related) > 0 {
		fmt.Fprintf(b, "### Related objects in this viewpoint\n\n%s\n\n", v.readingNamesText(unique(related), q))
	}
	v.readingGaps(b, filter(v.Gaps, func(g Fact) bool {
		return g.S("model_id") == name && (q.Viewpoint == "VP11" || g.S("viewpoint") == q.Viewpoint) && (q.Mode == "" || g.S("mode") == "" || g.S("mode") == q.Mode)
	}))
	b.WriteString(readingLimits(q.Viewpoint))
}

func (v *Views) readingNamesText(names []string, q Query) string {
	if len(names) == 0 {
		return "unspecified"
	}
	out := []string{}
	for _, n := range unique(names) {
		out = append(out, v.readingObjectLink(n, q))
	}
	return strings.Join(out, ", ")
}

func (v *Views) readingContract(b *strings.Builder, q Query, name string) {
	fmt.Fprintf(b, "Completeness: **%s**. This is a declared shape constraint, not observed validation or delivery.\n\nUsed by: %s.\n\n", or(v.prop(name, "completeness"), "unspecified"), v.readingNamesText(values(eq(v.Relations["upholds"], "object", name), "subject"), q))
	allowed := values(eq(v.Relations["permits"], "subject", name), "object")
	if len(allowed) > 0 {
		fmt.Fprintf(b, "Allowed payloads: %s. Request/result correlation requires explicit replies-to relationships and scenario steps.\n\n", v.readingNamesText(allowed, q))
	}
	shapes := append([]string{name}, values(eq(v.Relations["defines"], "subject", name), "object")...)
	count := 0
	for _, shape := range shapes {
		members := eq(v.Relations["has-field"], "subject", shape)
		if len(members) == 0 {
			continue
		}
		count += len(members)
		fmt.Fprintf(b, "### %s fields\n\n| Field | Value type | Presence | Source statement |\n| --- | --- | --- | --- |\n", shape)
		for _, f := range members {
			n := f.S("object")
			fmt.Fprintf(b, "| %s | %s | %s | %s |\n", v.readingObjectLink(n, q), or(v.prop(n, "value-type"), "unspecified"), or(v.prop(n, "presence"), "unspecified"), v.FactReference(f.S("id")))
		}
		b.WriteString("\n")
	}
	if count == 0 {
		b.WriteString("No record fields are declared for this contract or its variants. A channel contract may describe permitted messages without defining record fields.\n\n")
	}
	if len(shapes) > 1 {
		fmt.Fprintf(b, "Declared variants: %s.\n\n", v.readingNamesText(shapes[1:], q))
	}
}

func (v *Views) readingChannel(b *strings.Builder, q Query, name string) {
	contracts := values(eq(v.Relations["upholds"], "subject", name), "object")
	fmt.Fprintf(b, "Channel contract: %s.\n\nAllowed payloads: %s.\n\n", v.readingNamesText(contracts, q), v.readingNamesText(values(filter(v.Relations["permits"], func(f Fact) bool { return has(contracts, f.S("subject")) }), "object"), q))
	bindings := filter(v.Facts, func(f Fact) bool {
		return f.S("verb") == "uses" && f.S("channel") == name && (q.Mode == "" || f.S("mode") == q.Mode)
	})
	if len(bindings) == 0 {
		b.WriteString("No participation is declared in the selected mode scope. Global contract definitions remain available.\n\n")
		return
	}
	b.WriteString("### Actual participants\n\n| Participant | Role | Payload | Mode | Source statement |\n| --- | --- | --- | --- | --- |\n")
	for _, f := range bindings {
		fmt.Fprintf(b, "| %s | %s | %s | %s | %s |\n", v.readingObjectLink(f.S("subject"), q), f.S("role"), v.readingObjectLink(f.S("message"), q), v.readingObjectLink(f.S("mode"), q), v.FactReference(f.S("id")))
	}
	b.WriteString("\nA participant may send and receive on the same channel. Roles do not introduce additional SDL entities.\n\n")
}

func (v *Views) readingScenario(b *strings.Builder, q Query, name string) {
	modes := values(eq(v.Relations["runs-in"], "subject", name), "object")
	fmt.Fprintf(b, "Declared mode: %s. Completeness: %s.\n\n", v.readingNamesText(modes, q), or(v.prop(name, "completeness"), "unspecified"))
	if q.Mode != "" && !has(modes, q.Mode) {
		b.WriteString("This scenario has no declared applicability to the selected mode. Its global definition is retained; no steps are presented as behavior in that mode.\n\n")
		return
	}
	steps := sortedFacts(eq(eq(v.Facts, "verb", "step"), "subject", name))
	if len(steps) == 0 {
		b.WriteString("No ordered steps are declared. No runtime order is inferred.\n\n")
		return
	}
	b.WriteString("### Ordered scenario steps\n\n| Step | Sender → receiver | Channel | Payload / variant | Reply to step | Source statement |\n| --- | --- | --- | --- | --- | --- |\n")
	for _, f := range steps {
		payload := v.readingObjectLink(f.S("message"), q)
		if f.S("variant") != "" {
			payload += " / " + v.readingObjectLink(f.S("variant"), q)
		}
		reply := "unspecified"
		if f.N("reply_to") != 0 {
			reply = fmt.Sprint(f.N("reply_to"))
		}
		fmt.Fprintf(b, "| %d | %s → %s | %s | %s | %s | %s |\n", f.N("ordinal"), v.readingObjectLink(f.S("sender"), q), v.readingObjectLink(f.S("receiver"), q), v.readingObjectLink(f.S("channel"), q), payload, reply, v.FactReference(f.S("id")))
	}
	b.WriteString("\nOnly these numbered steps establish ordering. Preconditions, postconditions, timing, and execution evidence are unspecified.\n\n")
}

func (v *Views) readingEncoding(b *strings.Builder, q Query, name string) {
	fmt.Fprintf(b, "Encoded variants: %s. Byte order: %s. Bit order: %s.\n\n", v.readingNamesText(values(eq(v.Relations["encodes"], "subject", name), "object"), q), or(v.prop(name, "byte-order"), "unspecified"), or(v.prop(name, "bit-order"), "unspecified"))
	entries := eq(eq(v.Facts, "verb", "places"), "subject", name)
	sort.Slice(entries, func(i, j int) bool { return entries[i].N("offset") < entries[j].N("offset") })
	if len(entries) == 0 {
		b.WriteString("No explicit bit placements are declared. A complete wire layout is unspecified.\n\n")
		return
	}
	b.WriteString("### Explicit packet layout\n\n| Field | Bit range | Width | Value type | Presence | Source statement |\n| --- | --- | --- | --- | --- | --- |\n")
	for _, f := range entries {
		n := f.S("field")
		fmt.Fprintf(b, "| %s | %d–%d | %d | %s | %s | %s |\n", v.readingObjectLink(n, q), f.N("offset"), f.N("offset")+f.N("width")-1, f.N("width"), or(v.prop(n, "value-type"), "unspecified"), or(v.prop(n, "presence"), "unspecified"), v.FactReference(f.S("id")))
	}
	b.WriteString("\n")
}
