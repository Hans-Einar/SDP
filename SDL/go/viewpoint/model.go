// Package viewpoint projects only validated structural SDL facts. Rendering,
// publication and document hosts are separate consumers of this model.
package viewpoint

import (
	"crypto/sha256"
	"fmt"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	"sort"
	"strings"
)

const Version = "sdl-viewpoints-go-1.2"

type Fact map[string]any

func (f Fact) S(k string) string { s, _ := f[k].(string); return s }
func (f Fact) N(k string) int    { n, _ := f[k].(int); return n }

type Edge struct {
	Source   string `json:"source"`
	Relation string `json:"relation"`
	Target   string `json:"target"`
	Fact     string `json:"fact"`
}
type Diagram struct {
	ID          string         `json:"id"`
	Title       string         `json:"title"`
	Kind        string         `json:"kind"`
	Nodes       map[string]any `json:"nodes"`
	Edges       []Edge         `json:"edges"`
	SourceFacts []string       `json:"source_facts"`
	Elements    []Fact         `json:"elements"`
	Syntax      string         `json:"-"`
}
type Views struct {
	Model       *parser.Model
	Revision    string
	Kinds       map[string]string
	Facts       []Fact
	Relations   map[string][]Fact
	Diagrams    []Diagram
	Gaps        []Fact
	MessageSets []Fact
}

func New(source string) (*Views, error) {
	m, ds := parser.Check(source)
	if len(ds) > 0 {
		return nil, ds[0]
	}
	v := &Views{Model: m, Revision: fmt.Sprintf("%x", sha256.Sum256([]byte(source))), Kinds: map[string]string{}, Relations: map[string][]Fact{}, Facts: []Fact{}, Diagrams: []Diagram{}, Gaps: []Fact{}, MessageSets: []Fact{}}
	for _, d := range m.Declarations {
		v.Kinds[d.Name.Name] = d.Kind
	}
	for i, s := range m.Statements {
		f := Fact{"id": fmt.Sprintf("f%04d", i), "node": s.Kind, "line": s.Span.Line, "span": parser.Data(s.Span), "text": s.Sentence(), "subject": s.Subject.Name}
		switch s.Kind {
		case "Relation":
			f["verb"] = s.Verb
			f["object"] = s.Object.Name
		case "Dependency":
			f["verb"] = "requires"
			f["object"] = s.Interface.Name
			f["mode"] = s.Mode.Name
		case "Allocation":
			f["verb"] = "allocated-to"
			f["object"] = s.Container.Name
			f["mode"] = s.Mode.Name
		case "Projection":
			f["verb"] = "projects"
			f["dataset"] = s.Dataset.Name
			f["datagram"] = s.Datagram.Name
		case "Placement":
			f["verb"] = "places"
			f["field"] = s.Field.Name
			f["offset"] = s.Offset.Value
			f["width"] = s.Width.Value
		case "Participation":
			f["verb"] = "uses"
			f["channel"] = s.Channel.Name
			f["role"] = s.Role
			f["message"] = s.Message.Name
			f["mode"] = s.Mode.Name
		case "Step":
			f["verb"] = "step"
			f["ordinal"] = s.Ordinal.Value
			f["message"] = s.Message.Name
			f["variant"] = nil
			f["reply_to"] = nil
			if s.Variant != nil {
				f["variant"] = s.Variant.Name
			}
			if s.ReplyTo != nil {
				f["reply_to"] = s.ReplyTo.Value
			}
			f["sender"] = s.Sender.Name
			f["receiver"] = s.Receiver.Name
			f["channel"] = s.Channel.Name
		default:
			f["verb"] = "has"
			f["property"] = s.Property
			f["value"] = s.Value
		}
		v.Facts = append(v.Facts, f)
		if s.Kind == "Relation" || s.Kind == "Allocation" {
			v.Relations[f.S("verb")] = append(v.Relations[f.S("verb")], f)
		}
	}
	v.goals()
	v.data()
	v.channels()
	v.plans()
	v.architecture()
	return v, nil
}
func filter(fs []Fact, p func(Fact) bool) []Fact {
	o := []Fact{}
	for _, f := range fs {
		if p(f) {
			o = append(o, f)
		}
	}
	return o
}
func eq(fs []Fact, k, s string) []Fact { return filter(fs, func(f Fact) bool { return f.S(k) == s }) }
func values(fs []Fact, k string) []string {
	o := []string{}
	for _, f := range fs {
		if s := f.S(k); s != "" {
			o = append(o, s)
		}
	}
	return unique(o)
}
func unique(xs []string) []string {
	set := map[string]bool{}
	for _, s := range xs {
		set[s] = true
	}
	o := []string{}
	for s := range set {
		o = append(o, s)
	}
	sort.Strings(o)
	return o
}
func has(xs []string, s string) bool {
	for _, x := range xs {
		if x == s {
			return true
		}
	}
	return false
}
func join(groups ...[]Fact) []Fact {
	o := []Fact{}
	for _, g := range groups {
		o = append(o, g...)
	}
	return o
}
func (v *Views) names(kind string) []string {
	xs := []string{}
	for n, k := range v.Kinds {
		if k == kind {
			xs = append(xs, n)
		}
	}
	sort.Strings(xs)
	return xs
}
func (v *Views) nodes(names []string) map[string]any {
	o := map[string]any{}
	for _, d := range v.Model.Declarations {
		if has(names, d.Name.Name) {
			o["n_"+d.Name.Name] = map[string]any{"model_id": d.Name.Name, "kind": d.Kind, "span": parser.Data(d.Span)}
		}
	}
	return o
}
func (v *Views) diagram(id, title string, fs []Fact, names ...string) {
	edges := []Edge{}
	proof := []string{}
	for _, f := range fs {
		names = append(names, f.S("subject"), f.S("object"))
		edges = append(edges, Edge{"n_" + f.S("subject"), f.S("verb"), "n_" + f.S("object"), f.S("id")})
		proof = append(proof, f.S("id"))
	}
	if len(names) == 0 {
		return
	}
	sort.Slice(edges, func(i, j int) bool {
		a, b := edges[i], edges[j]
		return a.Source+"\x00"+a.Relation+"\x00"+a.Target+"\x00"+a.Fact < b.Source+"\x00"+b.Relation+"\x00"+b.Target+"\x00"+b.Fact
	})
	v.Diagrams = append(v.Diagrams, Diagram{id, title, "flowchart", v.nodes(names), edges, unique(proof), []Fact{}, ""})
}
func (d Diagram) Mermaid() string {
	if d.Syntax != "" {
		return d.Syntax
	}
	lines := []string{"flowchart LR"}
	names := []string{}
	for n := range d.Nodes {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		p := d.Nodes[n].(map[string]any)
		lines = append(lines, fmt.Sprintf("    %s[\"%s (%s)\"]", n, p["model_id"], p["kind"]))
	}
	for _, e := range d.Edges {
		lines = append(lines, fmt.Sprintf("    %s -->|%s| %s", e.Source, e.Relation, e.Target))
	}
	return strings.Join(lines, "\n") + "\n"
}
