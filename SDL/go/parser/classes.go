package parser

import (
	"fmt"
	"sort"
	"strings"
)

type ClassEnd struct {
	Class, Role Identifier
	Minimum     Integer
	Maximum     *Integer
	Span        Span
}
type ClassStatement struct {
	Subject         Identifier
	Verb            string
	Name            Identifier
	Type, Ownership string
	First, Second   ClassEnd
	Span            Span
}
type ClassModel struct {
	Header       Header
	Declarations []Declaration
	Statements   []ClassStatement
	Span         Span
}

func ParseClasses(text string) (model *ClassModel, err error) {
	defer func() {
		if v := recover(); v != nil {
			if d, ok := v.(Diagnostic); ok {
				model = nil
				err = d
			} else {
				panic(v)
			}
		}
	}()
	p, err := newReader(text)
	if err != nil {
		return nil, err
	}
	start := p.expect("language").span.Start
	p.expect("class-core")
	p.expect("version")
	if p.current().text != "0.1" {
		p.fail("UNSUPPORTED_VERSION", "Expected class-core 0.1")
	}
	p.take()
	m := &ClassModel{Header: Header{"class-core", "0.1", p.finish(start)}, Declarations: []Declaration{}, Statements: []ClassStatement{}}
	for has([]string{"class", "association"}, p.current().text) {
		t := p.take()
		n := p.identifier()
		m.Declarations = append(m.Declarations, Declaration{t.text, n, p.finish(t.span.Start)})
	}
	end := func() ClassEnd {
		e := ClassEnd{Class: p.identifier()}
		p.expect("role")
		e.Role = p.identifier()
		p.expect("multiplicity")
		e.Minimum = p.integer()
		p.expect("to")
		if p.current().text == "many" {
			p.take()
		} else {
			x := p.integer()
			e.Maximum = &x
		}
		e.Span = p.source.span(e.Class.Span.Start, p.tokens[p.index-1].span.End)
		return e
	}
	for p.current().text != "" {
		s := ClassStatement{Subject: p.identifier()}
		s.Verb = p.take().text
		switch s.Verb {
		case "attribute", "operation":
			s.Name = p.identifier()
			if s.Verb == "attribute" {
				p.expect("as")
			} else {
				p.expect("returns")
			}
			s.Type = p.take().text
		case "links":
			s.First = end()
			p.expect("with")
			s.Second = end()
		case "has":
			p.expect("ownership")
			p.expect("=")
			s.Ownership = p.take().text
		default:
			p.fail("UNSUPPORTED_SYNTAX", "Unknown class statement")
		}
		s.Span = p.finish(s.Subject.Span.Start)
		m.Statements = append(m.Statements, s)
	}
	m.Span = p.source.span(start, p.tokens[p.index-1].span.End)
	return m, nil
}
func (e ClassEnd) Multiplicity() string {
	max := "*"
	if e.Maximum != nil {
		max = fmt.Sprint(e.Maximum.Value)
	}
	min := fmt.Sprint(e.Minimum.Value)
	if min == max {
		return min
	}
	return min + ".." + max
}
func (e ClassEnd) sentence() string {
	max := "many"
	if e.Maximum != nil {
		max = fmt.Sprint(e.Maximum.Value)
	}
	return fmt.Sprintf("%s role %s multiplicity %d to %s", e.Class.Name, e.Role.Name, e.Minimum.Value, max)
}
func (s ClassStatement) Sentence() string {
	base := s.Subject.Name + " " + s.Verb + " "
	switch s.Verb {
	case "attribute":
		return base + s.Name.Name + " as " + s.Type + "."
	case "operation":
		return base + s.Name.Name + " returns " + s.Type + "."
	case "links":
		return base + s.First.sentence() + " with " + s.Second.sentence() + "."
	default:
		return base + "ownership = " + s.Ownership + "."
	}
}
func CanonicalClasses(m *ClassModel) string {
	decl := append([]Declaration{}, m.Declarations...)
	sort.Slice(decl, func(i, j int) bool { return decl[i].Name.Name < decl[j].Name.Name })
	lines := []string{"language class-core version 0.1."}
	for _, d := range decl {
		lines = append(lines, d.Kind+" "+d.Name.Name+".")
	}
	facts := []string{}
	for _, s := range m.Statements {
		facts = append(facts, s.Sentence())
	}
	sort.Strings(facts)
	return strings.Join(append(lines, facts...), "\n") + "\n"
}
func CheckClasses(text string) (*ClassModel, error) {
	m, e := ParseClasses(text)
	if e != nil {
		return nil, e
	}
	if e = ValidateClasses(m); e != nil {
		return nil, e
	}
	if CanonicalClasses(m) != text {
		return nil, Diagnostic{"NONCANONICAL_FORM", "Use canonical class ordering and spacing", m.Span}
	}
	return m, nil
}
