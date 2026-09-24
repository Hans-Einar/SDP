package parser

import (
	"sort"
	"strings"
)

type ScalarType string

const (
	TextType    ScalarType = "text"
	IntegerType ScalarType = "integer"
	BooleanType ScalarType = "boolean"
)

type ActionStatement struct {
	Subject Identifier `json:"subject"`
	Verb    string     `json:"verb"`
	Object  Identifier `json:"object"`
	Type    ScalarType `json:"type,omitempty"`
	Span    Span       `json:"span"`
}
type ActionModel struct {
	Header       Header            `json:"header"`
	Declarations []Declaration     `json:"declarations"`
	Statements   []ActionStatement `json:"statements"`
	Span         Span              `json:"span"`
}
type RecordType map[string]ScalarType
type Action struct {
	Name, Input, Output, GoSymbol string
	Span                          Span
}
type Program struct {
	Records map[string]RecordType
	Actions map[string]Action
	Model   *ActionModel
}

func ParseActions(text string) (model *ActionModel, err error) {
	defer func() {
		if v := recover(); v != nil {
			if d, ok := v.(Diagnostic); ok {
				err = d
				model = nil
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
	p.expect("action-core")
	p.expect("version")
	if p.current().text != "0.1" {
		p.fail("UNSUPPORTED_VERSION", "Only action-core 0.1 is supported")
	}
	p.take()
	m := &ActionModel{Header: Header{"action-core", "0.1", p.finish(start)}, Declarations: []Declaration{}, Statements: []ActionStatement{}}
	for p.current().text == "action" || p.current().text == "record" {
		t := p.take()
		name := p.identifier()
		m.Declarations = append(m.Declarations, Declaration{t.text, name, p.finish(t.span.Start)})
	}
	for p.current().text != "" {
		s := ActionStatement{Subject: p.identifier()}
		s.Verb = p.current().text
		if !has([]string{"takes", "returns", "invokes", "field"}, s.Verb) {
			p.fail("UNSUPPORTED_SYNTAX", "Unknown executable action statement")
		}
		p.take()
		s.Object = p.identifier()
		if s.Verb == "field" {
			p.expect("as")
			s.Type = ScalarType(p.current().text)
			if !has([]string{"text", "integer", "boolean"}, string(s.Type)) {
				p.fail("UNSUPPORTED_SYNTAX", "Unknown executable scalar type")
			}
			p.take()
		}
		s.Span = p.finish(s.Subject.Span.Start)
		m.Statements = append(m.Statements, s)
	}
	m.Span = p.source.span(start, p.tokens[p.index-1].span.End)
	return m, nil
}
func (s ActionStatement) Sentence() string {
	line := s.Subject.Name + " " + s.Verb + " " + s.Object.Name
	if s.Verb == "field" {
		line += " as " + string(s.Type)
	}
	return line + "."
}
func CanonicalActions(m *ActionModel) string {
	decl := append([]Declaration{}, m.Declarations...)
	sort.Slice(decl, func(i, j int) bool { return decl[i].Name.Name < decl[j].Name.Name })
	lines := []string{"language action-core version 0.1."}
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
func CompileActions(text string) (*Program, error) {
	m, err := ParseActions(text)
	if err != nil {
		return nil, err
	}
	program, err := ValidateActions(m)
	if err != nil {
		return nil, err
	}
	canonical := CanonicalActions(m)
	if text != canonical {
		offset := 0
		for offset < len(text) && offset < len(canonical) && text[offset] == canonical[offset] {
			offset++
		}
		return nil, Diagnostic{"NONCANONICAL_FORM", "Use canonical ordering and spacing", newSource(text).span(offset, min(offset+1, len(text)))}
	}
	return program, nil
}
func ValidateActions(m *ActionModel) (*Program, error) {
	if m == nil {
		return nil, Diagnostic{Code: "ACTION_MODEL", Message: "Missing action model"}
	}
	if m.Header.Language != "action-core" || m.Header.Version != "0.1" {
		return nil, Diagnostic{"UNSUPPORTED_VERSION", "Expected action-core 0.1", m.Header.Span}
	}
	p := &Program{Records: map[string]RecordType{}, Actions: map[string]Action{}, Model: m}
	symbols := map[string]Declaration{}
	bad := func(code, msg string, span Span) (*Program, error) { return nil, Diagnostic{code, msg, span} }
	for _, d := range m.Declarations {
		name := d.Name.Name
		if !namePattern.MatchString(name) {
			return bad("UNSUPPORTED_SYNTAX", "Invalid declaration name", d.Name.Span)
		}
		if _, ok := symbols[name]; ok {
			return bad("DUPLICATE_DECLARATION", name, d.Name.Span)
		}
		symbols[name] = d
		switch d.Kind {
		case "record":
			p.Records[name] = RecordType{}
		case "action":
			p.Actions[name] = Action{Name: name, Span: d.Span}
		default:
			return bad("EXECUTABLE_TYPE", d.Kind, d.Span)
		}
	}
	if len(p.Records) > 128 || len(p.Actions) > 128 {
		return bad("ACTION_LIMIT", "At most 128 records and actions", m.Span)
	}
	seen := map[string]bool{}
	for _, s := range m.Statements {
		if !namePattern.MatchString(s.Subject.Name) || !namePattern.MatchString(s.Object.Name) {
			return bad("UNSUPPORTED_SYNTAX", "Invalid action identifier", s.Span)
		}
		subject, ok := symbols[s.Subject.Name]
		if !ok {
			return bad("UNDECLARED_NAME", s.Subject.Name, s.Subject.Span)
		}
		k := s.Subject.Name + "/" + s.Verb
		if s.Verb == "field" {
			k += "/" + s.Object.Name
		}
		if seen[k] {
			return bad("ACTION_CARDINALITY", "Repeated action property or record field", s.Span)
		}
		seen[k] = true
		if s.Verb == "field" {
			if subject.Kind != "record" {
				return bad("SUBJECT_TYPE_MISMATCH", "field requires record", s.Subject.Span)
			}
			if s.Type != TextType && s.Type != IntegerType && s.Type != BooleanType {
				return bad("EXECUTABLE_TYPE", "Unknown field type", s.Span)
			}
			p.Records[s.Subject.Name][s.Object.Name] = s.Type
			continue
		}
		if subject.Kind != "action" {
			return bad("SUBJECT_TYPE_MISMATCH", "Action relation requires action", s.Subject.Span)
		}
		a := p.Actions[s.Subject.Name]
		switch s.Verb {
		case "takes", "returns":
			if symbols[s.Object.Name].Kind != "record" {
				return bad("OBJECT_TYPE_MISMATCH", "Action signature requires declared record", s.Object.Span)
			}
			if s.Verb == "takes" {
				a.Input = s.Object.Name
			} else {
				a.Output = s.Object.Name
			}
		case "invokes":
			a.GoSymbol = s.Object.Name
		default:
			return bad("UNSUPPORTED_SYNTAX", s.Verb, s.Span)
		}
		p.Actions[a.Name] = a
	}
	for _, d := range m.Declarations {
		if d.Kind == "record" {
			n := len(p.Records[d.Name.Name])
			if n < 1 || n > 32 {
				return bad("RECORD_FIELDS", "Record requires 1–32 fields", d.Span)
			}
		} else {
			a := p.Actions[d.Name.Name]
			if a.Input == "" || a.Output == "" || a.GoSymbol == "" {
				return bad("ACTION_CARDINALITY", "Action needs takes, returns and invokes", d.Span)
			}
		}
	}
	return p, nil
}
