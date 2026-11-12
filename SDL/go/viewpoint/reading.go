package viewpoint

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

// ReadingCollection is the semantic authority shared by navigation and documents.
// IDs are stable presentation selectors, not new SDL declaration kinds.
type ReadingCollection struct {
	ID    string   `json:"id"`
	Label string   `json:"label"`
	Names []string `json:"names"`
}

func (v *Views) ReadingCollections(vp string) []ReadingCollection {
	out := []ReadingCollection{}
	add := func(id, label string, kinds ...string) {
		ns := []string{}
		for _, k := range kinds {
			ns = append(ns, v.names(k)...)
		}
		out = append(out, ReadingCollection{id, label, unique(ns)})
	}
	derived := func(id, label string, ns []string) { out = append(out, ReadingCollection{id, label, unique(ns)}) }
	switch vp {
	case "VP01":
		add("usecases", "Use cases", "usecase")
		add("actors", "Actors", "actor")
		add("features", "Features", "feature")
	case "VP02":
		add("containers", "Containers", "container")
		add("units", "Units", "unit")
	case "VP03":
		add("units", "Units", "unit")
		add("capabilities", "Capabilities", "capability")
		add("functionalities", "Functionalities", "functionality")
	case "VP04":
		add("interfaces", "Interfaces", "interface")
		add("consumers", "Consumers", "unit", "container")
	case "VP05":
		add("modes", "Modes", "mode")
		add("capabilities", "Capabilities", "capability")
		add("interfaces", "Interfaces", "interface")
	case "VP06":
		add("activities", "Activities", "activity")
		derived("deliverables", "Deliverables", values(v.Relations["delivers"], "object"))
	case "VP07":
		add("features", "Features", "feature")
		add("containers", "Containers", "container")
		add("units", "Units", "unit")
		add("modes", "Modes", "mode")
	case "VP08":
		add("channels", "Channels", "channel")
		add("scenarios", "Scenarios", "scenario")
		add("messages", "Messages and datagrams", "message", "datagram")
		add("contracts", "Contracts", "contract")
		add("modes", "Modes", "mode")
		add("participants", "Participants", "unit", "container")
	case "VP09":
		add("datasets", "Datasets", "dataset")
		add("datagrams", "Datagrams", "datagram")
		add("databases", "Databases", "database")
		add("contracts", "Contracts", "contract")
	case "VP10":
		add("encodings", "Encodings", "encoding")
		add("datagrams", "Datagrams", "datagram")
	case "VP11":
		properties := values(eq(v.Facts, "verb", "has"), "subject")
		gaps := []string{}
		for _, g := range v.Gaps {
			n := g.S("model_id")
			if v.Kinds[n] != "" && !has(properties, n) {
				gaps = append(gaps, n)
			}
		}
		rest := []string{}
		for n := range v.Kinds {
			if !has(properties, n) && !has(gaps, n) {
				rest = append(rest, n)
			}
		}
		derived("properties", "Properties", properties)
		derived("traceability", "Traceability", rest)
		derived("model-gaps", "Model gaps", gaps)
	}
	return out
}

const ReadingChapterPageSize = 24
const ReadingRegisterPageSize = 100

func (q Query) validateReading(v *Views) error {
	if _, ok := SpecFor(q.Viewpoint); !ok {
		return fmt.Errorf("unknown-viewpoint: %s", q.Viewpoint)
	}
	selectors := 0
	for _, s := range []string{q.Collection, q.Register, q.Focus} {
		if s != "" {
			selectors++
		}
	}
	if selectors > 1 || q.Diagram != "" || q.Level != "" || len(q.Relations) > 0 || (q.Direction != "" && q.Direction != "both") || (q.Depth != 0 && q.Depth != 2) {
		return fmt.Errorf("unsupported-reading-selection")
	}
	if q.Focus != "" && v.Kinds[q.Focus] == "" {
		return fmt.Errorf("unknown-focus: %s", q.Focus)
	}
	if q.Mode != "" && v.Kinds[q.Mode] != "mode" {
		return fmt.Errorf("unknown-mode: %s", q.Mode)
	}
	if q.Page < 0 || q.Page > 0 && q.Collection == "" && q.Register == "" {
		return fmt.Errorf("invalid-reading-page")
	}
	if q.Register != "" && (q.Viewpoint != "VP11" || !has([]string{"facts", "declarations"}, q.Register)) {
		return fmt.Errorf("unknown-reading-register")
	}
	if q.Collection != "" {
		found := false
		for _, c := range v.ReadingCollections(q.Viewpoint) {
			if c.ID == q.Collection {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("unknown-reading-collection: %s", q.Collection)
		}
	}
	_, _, _, err := v.readingPage(q)
	return err
}

func (v *Views) readingPage(q Query) (start, end, total int, err error) {
	size := ReadingChapterPageSize
	if q.Register != "" {
		size = ReadingRegisterPageSize
		if q.Register == "facts" {
			total = len(v.Facts)
		} else {
			total = len(v.Model.Declarations)
		}
	}
	if q.Collection != "" {
		for _, c := range v.ReadingCollections(q.Viewpoint) {
			if c.ID == q.Collection {
				total = len(c.Names)
			}
		}
	}
	page := q.Page
	if page == 0 {
		page = 1
	}
	pages := (total + size - 1) / size
	if pages == 0 {
		pages = 1
	}
	if page > pages {
		return 0, 0, total, fmt.Errorf("reading-page-out-of-range")
	}
	start = (page - 1) * size
	end = start + size
	if end > total {
		end = total
	}
	return
}

func (v *Views) readingNames(q Query) []string {
	if q.Focus != "" {
		return []string{q.Focus}
	}
	if q.Collection != "" {
		for _, c := range v.ReadingCollections(q.Viewpoint) {
			if c.ID == q.Collection {
				a, b, _, _ := v.readingPage(q)
				return c.Names[a:b]
			}
		}
	}
	return nil
}

var readingVerbs = map[string][]string{
	"VP01": {"pursues", "supports", "contributes-to", "exercises"},
	"VP02": {"contains", "owns", "allocated-to", "uses"},
	"VP03": {"owns", "provides", "realizes"},
	"VP04": {"consumes"},
	"VP05": {"requires"},
	"VP06": {"refines", "addresses", "delivers", "depends-on", "owns"},
	"VP07": {"supports", "contributes-to", "owns", "contains", "allocated-to", "realizes", "provides"},
	"VP08": {"uses", "upholds", "permits", "replies-to", "runs-in", "step", "exercises", "illustrates", "defines", "has-field", "from"},
	"VP09": {"holds", "from", "owns", "upholds", "defines", "has-field", "projects", "permits"},
	"VP10": {"encodes", "places", "defines", "has-field", "upholds", "from"},
}

func factTouches(f Fact, ns []string) bool {
	for _, k := range []string{"subject", "object", "dataset", "datagram", "channel", "message", "sender", "receiver", "field", "variant", "mode"} {
		if has(ns, f.S(k)) {
			return true
		}
	}
	return false
}

// readingFacts follows only semantic companion relations, then retains statements
// adjacent to that scope. It never walks an arbitrary transitive graph.
func (v *Views) readingFacts(q Query, names []string) []Fact {
	scope := append([]string{}, names...)
	// Include supporting realization and nested data shapes, without pulling all
	// unrelated users of a shared interface/channel into an object chapter.
	for pass := 0; pass < 3; pass++ {
		next := append([]string{}, scope...)
		for _, f := range v.Facts {
			if q.Mode != "" && f.S("mode") != "" && f.S("mode") != q.Mode {
				continue
			}
			verb := f.S("verb")
			if has(scope, f.S("subject")) && has([]string{"upholds", "defines", "has-field", "encodes", "owns"}, verb) {
				next = append(next, f.S("object"))
			}
			if has(scope, f.S("object")) && has([]string{"supports", "contributes-to", "realizes"}, verb) {
				next = append(next, f.S("subject"))
			}
		}
		scope = unique(next)
	}
	return filter(v.Facts, func(f Fact) bool {
		if q.Mode != "" && f.S("mode") != "" && f.S("mode") != q.Mode {
			return false
		}
		if q.Viewpoint != "VP11" && f.S("verb") != "has" && !has(readingVerbs[q.Viewpoint], f.S("verb")) {
			return false
		}
		return factTouches(f, scope)
	})
}

func (v *Views) reading(q Query) (*Views, error) {
	out := *v
	out.Diagrams = []Diagram{}
	out.MessageSets = []Fact{}
	out.Facts = []Fact{}
	out.Gaps = []Fact{}
	out.Relations = map[string][]Fact{}
	if q.Register != "" {
		if q.Register == "facts" {
			a, b, _, _ := v.readingPage(q)
			out.Facts = v.Facts[a:b]
		}
		return &out, nil
	}
	names := v.readingNames(q)
	if len(names) == 0 && q.Collection == "" && q.Focus == "" {
		// Reuse the existing language-specific overview, never reject a readable
		// landing solely because the optional diagram exceeds its geometry limits.
		overview, err := v.overview(q)
		if err == nil {
			out.Diagrams = overview.Diagrams
		}
		out.Gaps = filter(v.Gaps, func(g Fact) bool {
			return (q.Viewpoint == "VP11" || g.S("viewpoint") == q.Viewpoint) && (q.Mode == "" || g.S("mode") == "" || g.S("mode") == q.Mode)
		})
		if len(out.Diagrams) > 40 {
			out.Diagrams = nil
		}
		v.readingProof(&out)
		return &out, nil
	}
	out.Facts = v.readingFacts(q, names)
	for _, f := range out.Facts {
		out.Relations[f.S("verb")] = append(out.Relations[f.S("verb")], f)
	}
	out.Gaps = filter(v.Gaps, func(g Fact) bool {
		return has(names, g.S("model_id")) && (q.Viewpoint == "VP11" || g.S("viewpoint") == q.Viewpoint) && (q.Mode == "" || g.S("mode") == "" || g.S("mode") == q.Mode)
	})
	seen := map[string]bool{}
	for _, name := range names {
		for _, d := range v.Diagrams {
			primary := strings.TrimPrefix(d.ID, q.Viewpoint+"-")
			if q.Mode != "" && d.Kind == "sequence" && !has(values(eq(v.Relations["runs-in"], "subject", primary), "object"), q.Mode) {
				continue
			}
			if !strings.HasPrefix(d.ID, q.Viewpoint+"-") || !(primary == name || primary == "contract-"+name || primary == "data-"+name || primary == "feature-"+name || primary == "detail-"+name || primary == "work-"+name || strings.HasPrefix(primary, name+"-")) {
				continue
			}
			if q.Mode != "" {
				var keep bool
				d, keep = v.readingModeDiagram(d, q)
				if !keep {
					continue
				}
			}

			if !seen[d.ID] && len(d.Nodes) <= 128 && len(d.Edges) <= 256 {
				out.Diagrams = append(out.Diagrams, d)
				seen[d.ID] = true
			}
		}
	}
	// Preserve the explicit exercises-linked scenario companions of VP01 details.
	for _, name := range names {
		if v.Kinds[name] == "usecase" && q.Viewpoint == "VP01" {
			detail := *v
			detail.Facts = out.Facts
			detail.Relations = out.Relations
			detail.Diagrams = out.Diagrams
			v.useCaseDetails(&detail, Query{Viewpoint: "VP01", Focus: name, Mode: q.Mode})
			for _, d := range detail.Diagrams {
				if !seen[d.ID] {
					out.Diagrams = append(out.Diagrams, d)
					seen[d.ID] = true
				} else {
					for i := range out.Diagrams {
						if out.Diagrams[i].ID == d.ID {
							out.Diagrams[i] = d
						}
					}
				}
			}
		}
	}
	if len(out.Diagrams) > 40 {
		out.Diagrams = nil
	}
	v.readingProof(&out)
	return &out, nil
}

// Mode-specific proof must be filtered with geometry, never restored from the
// unfiltered diagram after text selection. Atomic diagrams cannot be pruned.
func (v *Views) readingModeDiagram(d Diagram, q Query) (Diagram, bool) {
	facts := map[string]Fact{}
	for _, f := range v.Facts {
		facts[f.S("id")] = f
	}
	eligible := func(id string) bool {
		f := facts[id]
		if f.S("mode") != "" && f.S("mode") != q.Mode {
			return false
		}
		if f.S("verb") == "runs-in" && f.S("object") != q.Mode {
			return false
		}
		if f.S("verb") == "step" {
			modes := values(eq(v.Relations["runs-in"], "subject", f.S("subject")), "object")
			if len(modes) > 0 && !has(modes, q.Mode) {
				return false
			}
		}
		return true
	}
	proof := []string{}
	for _, id := range d.SourceFacts {
		if eligible(id) {
			proof = append(proof, id)
		} else if d.Kind != "flowchart" {
			return d, false
		}
	}
	if d.Kind != "flowchart" {
		return d, true
	}
	// Filtered geometry must not render a cached unfiltered syntax string.
	d.Syntax = ""
	edges := []Edge{}
	active := map[string]bool{}
	for _, e := range d.Edges {
		if eligible(e.Fact) {
			edges = append(edges, e)
			active[e.Source] = true
			active[e.Target] = true
		}
	}
	nodes := map[string]any{}
	for id, p := range d.Nodes {
		node, _ := p.(map[string]any)
		if active[id] || node["model_id"] == q.Focus || len(d.Edges) == 0 {
			nodes[id] = p
		}
	}
	elements := []Fact{}
	for _, element := range d.Elements {
		keep := element.S("mode") == "" || element.S("mode") == q.Mode
		if id := element.S("fact"); id != "" && !eligible(id) {
			keep = false
		}
		for _, key := range []string{"proof", "source_facts"} {
			if ids, ok := element[key].([]string); ok {
				for _, id := range ids {
					if !eligible(id) {
						keep = false
					}
				}
			}
		}
		if keep {
			elements = append(elements, element)
		}
	}
	d.Elements = elements
	d.Nodes = nodes
	d.Edges = edges
	d.SourceFacts = proof
	return d, len(nodes) > 0
}

func (v *Views) readingProof(out *Views) {
	used := map[string]bool{}
	present := map[string]bool{}
	for _, d := range out.Diagrams {
		for _, id := range d.SourceFacts {
			used[id] = true
		}
	}
	for _, f := range out.Facts {
		present[f.S("id")] = true
	}
	for _, f := range v.Facts {
		if used[f.S("id")] && !present[f.S("id")] {
			out.Facts = append(out.Facts, f)
			out.Relations[f.S("verb")] = append(out.Relations[f.S("verb")], f)
		}
	}
}

// ReadingTarget keeps context selection in producer-generated actions. The host
// binds model/revision using its retained document context; those identities are
// never copied into visible labels.
func (v *Views) ReadingTarget(q Query) string {
	if v.ReadingProject == "" {
		return ""
	}
	values := url.Values{"projection": {"reading"}, "target": {"main"}, "consumer": {"xfmd"}}
	for k, s := range map[string]string{"focus": q.Focus, "collection": q.Collection, "register": q.Register, "mode": q.Mode} {
		if s != "" {
			values.Set(k, s)
		}
	}
	if q.Page > 1 {
		values.Set("page", strconv.Itoa(q.Page))
	}
	return "sdl-view://" + v.ReadingProject + "/" + q.Viewpoint + "?" + values.Encode()
}

func (v *Views) readingLink(label string, q Query) string {
	if target := v.ReadingTarget(q); target != "" {
		return "[" + label + "](" + target + ")"
	}
	return label
}
func (v *Views) readingObjectLink(n string, q Query) string {
	q.Focus = n
	q.Collection = ""
	q.Register = ""
	q.Page = 0
	return v.readingLink(n, q)
}
func readingEscape(s string) string {
	return strings.NewReplacer("|", "\\|", "[", "\\[", "]", "\\]", "<", "&lt;", ">", "&gt;", "\n", " ", "`", "\\`", "*", "\\*", "_", "\\_").Replace(s)
}
func (v *Views) readingDeclaration(name string) string {
	for _, d := range v.Model.Declarations {
		if d.Name.Name == name {
			return v.SpanReference(readingEscape(name)+" — "+traceabilityPosition(v.sourceLabel(d.Span.Source), d.Span.Line), d.Span)
		}
	}
	return readingEscape(name)
}

func sortedFacts(fs []Fact) []Fact {
	out := append([]Fact{}, fs...)
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].S("verb") == "step" && out[j].S("verb") == "step" {
			return out[i].N("ordinal") < out[j].N("ordinal")
		}
		return out[i].S("text") < out[j].S("text")
	})
	return out
}
