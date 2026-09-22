package viewpoint

import (
	"fmt"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	"strings"
)

func Classes(m *parser.ClassModel) (Diagram, []Fact, error) {
	if e := parser.ValidateClasses(m); e != nil {
		return Diagram{}, nil, e
	}
	d := Diagram{ID: "CL01-classes", Title: "Eksplisitt klassestruktur", Kind: "class", Nodes: map[string]any{}, Edges: []Edge{}, SourceFacts: []string{}, Elements: []Fact{}}
	lines := []string{"classDiagram"}
	facts := []Fact{}
	ownership := map[string]string{}
	ownFacts := map[string]string{}
	for i, s := range m.Statements {
		id := fmt.Sprintf("f%04d", i)
		facts = append(facts, Fact{"id": id, "text": s.Sentence(), "span": parser.Data(s.Span), "line": s.Span.Line, "subject": s.Subject.Name, "verb": s.Verb})
		d.SourceFacts = append(d.SourceFacts, id)
		if s.Verb == "has" {
			ownership[s.Subject.Name] = s.Ownership
			ownFacts[s.Subject.Name] = id
		}
	}
	for _, c := range m.Declarations {
		if c.Kind != "class" {
			continue
		}
		n := c.Name.Name
		d.Nodes[n] = map[string]any{"model_id": n, "kind": "class", "span": parser.Data(c.Span)}
		lines = append(lines, "    class "+n+" {")
		for _, s := range m.Statements {
			if s.Subject.Name != n {
				continue
			}
			switch s.Verb {
			case "attribute":
				lines = append(lines, fmt.Sprintf("        +%s %s", s.Type, s.Name.Name))
			case "operation":
				lines = append(lines, fmt.Sprintf("        +%s() %s", s.Name.Name, s.Type))
			}
		}
		lines = append(lines, "    }")
	}
	for i, s := range m.Statements {
		if s.Verb != "links" {
			continue
		}
		kind := ownership[s.Subject.Name]
		arrow := "--"
		if kind == "aggregation" {
			arrow = "o--"
		}
		if kind == "composition" {
			arrow = "*--"
		}
		a, b := s.First, s.Second
		lines = append(lines, fmt.Sprintf("    %s \"%s %s\" %s \"%s %s\" %s : %s", a.Class.Name, a.Multiplicity(), a.Role.Name, arrow, b.Multiplicity(), b.Role.Name, b.Class.Name, s.Subject.Name))
		f := fmt.Sprintf("f%04d", i)
		d.Edges = append(d.Edges, Edge{a.Class.Name, kind, b.Class.Name, f})
		d.Elements = append(d.Elements, Fact{"association": s.Subject.Name, "ownership": kind, "first": a, "second": b, "fact": f, "ownership_fact": ownFacts[s.Subject.Name]})
	}
	d.Syntax = strings.Join(lines, "\n") + "\n"
	return d, facts, nil
}
