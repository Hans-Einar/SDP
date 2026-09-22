package parser

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

const MaxBytes = 2 << 20

var tokenPattern = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9-]*|^[0-9]+(?:\.[0-9]+)*|^[.=]`)
var namePattern = regexp.MustCompile(`^[A-Z][A-Za-z0-9]*$`)
var integerPattern = regexp.MustCompile(`^(0|[1-9][0-9]{0,5})$`)

type source struct {
	text  string
	lines []int
}

func newSource(text string) *source {
	s := &source{text: text, lines: []int{0}}
	for i, c := range text {
		if c == '\n' {
			s.lines = append(s.lines, i+1)
		}
	}
	return s
}
func (s *source) span(start, end int) Span {
	a := sort.Search(len(s.lines), func(i int) bool { return s.lines[i] > start }) - 1
	b := sort.Search(len(s.lines), func(i int) bool { return s.lines[i] > end }) - 1
	return Span{start, end, a + 1, start - s.lines[a] + 1, b + 1, end - s.lines[b] + 1}
}

type token struct {
	text string
	span Span
}
type reader struct {
	source *source
	tokens []token
	index  int
}

func (p *reader) current() token            { return p.tokens[p.index] }
func (p *reader) fail(code, message string) { panic(Diagnostic{code, message, p.current().span}) }
func (p *reader) take() token {
	t := p.current()
	if t.text == "" {
		p.fail("UNSUPPORTED_SYNTAX", "Unexpected end of input")
	}
	p.index++
	return t
}
func (p *reader) expect(text string) token {
	if p.current().text != text {
		p.fail("UNSUPPORTED_SYNTAX", fmt.Sprintf("Expected %q, received %q", text, p.current().text))
	}
	return p.take()
}
func (p *reader) identifier() Identifier {
	if !namePattern.MatchString(p.current().text) {
		p.fail("UNSUPPORTED_SYNTAX", "Expected an ASCII identifier matching [A-Z][A-Za-z0-9]*")
	}
	t := p.take()
	return Identifier{t.text, t.span}
}
func (p *reader) integer() Integer {
	if !integerPattern.MatchString(p.current().text) {
		p.fail("UNSUPPORTED_SYNTAX", "Expected a canonical integer between 0 and 65536")
	}
	t := p.take()
	v, _ := strconv.Atoi(t.text)
	if v > 65536 {
		p.fail("LAYOUT_LIMIT", "Integer exceeds fixed-layout limit")
	}
	return Integer{v, t.span}
}
func (p *reader) finish(start int) Span { return p.source.span(start, p.expect(".").span.End) }
func Parse(text string) (model *Model, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			if d, ok := recovered.(Diagnostic); ok {
				err = d
				model = nil
			} else {
				panic(recovered)
			}
		}
	}()
	s := newSource(text)
	if len(text) > MaxBytes {
		return nil, Diagnostic{"SOURCE_LIMIT", "Source exceeds 2 MiB", s.span(0, 0)}
	}
	if !utf8.ValidString(text) {
		return nil, Diagnostic{"INPUT_ERROR", "Invalid UTF-8", s.span(0, 0)}
	}
	p := &reader{source: s}
	for offset := 0; offset < len(text); {
		if strings.ContainsRune(" \t\r\n", rune(text[offset])) {
			offset++
			continue
		}
		match := tokenPattern.FindString(text[offset:])
		if match == "" {
			return nil, Diagnostic{"UNSUPPORTED_SYNTAX", "Unsupported character", s.span(offset, offset+1)}
		}
		p.tokens = append(p.tokens, token{match, s.span(offset, offset+len(match))})
		offset += len(match)
		if len(p.tokens) > 250000 {
			return nil, Diagnostic{"TOKEN_LIMIT", "More than 250000 tokens", s.span(offset, offset)}
		}
	}
	p.tokens = append(p.tokens, token{"", s.span(len(text), len(text))})
	start := p.expect("language").span.Start
	p.expect("design-core")
	p.expect("version")
	if p.current().text != "0.5" {
		p.fail("UNSUPPORTED_VERSION", "Only design-core version 0.5 is supported")
	}
	p.take()
	m := &Model{Header: Header{"design-core", "0.5", p.finish(start)}, Declarations: []Declaration{}, Statements: []Statement{}}
	for has(kinds, p.current().text) {
		kind := p.take()
		name := p.identifier()
		m.Declarations = append(m.Declarations, Declaration{kind.text, name, p.finish(kind.span.Start)})
	}
	for p.current().text != "" {
		m.Statements = append(m.Statements, p.statement())
	}
	m.Span = s.span(start, p.tokens[p.index-1].span.End)
	return m, nil
}
func (p *reader) statement() Statement {
	s := Statement{Subject: p.identifier()}
	verb := p.current().text
	if _, ok := signatures[verb]; ok {
		p.take()
		s.Kind = "Relation"
		s.Verb = verb
		s.Object = p.identifier()
	} else {
		switch verb {
		case "has":
			p.take()
			s.Kind = "PropertyAssignment"
			s.Property = p.current().text
			if _, ok := properties[s.Property]; !ok {
				p.fail("UNSUPPORTED_SYNTAX", "Unknown property")
			}
			p.take()
			p.expect("=")
			s.Value = p.current().text
			known := false
			for _, values := range properties {
				known = known || has(values, s.Value)
			}
			if !known {
				p.fail("UNSUPPORTED_SYNTAX", "Unknown property value")
			}
			p.take()
		case "requires", "allocated-to":
			p.take()
			target := p.identifier()
			p.expect("in")
			p.expect("mode")
			s.Mode = p.identifier()
			if verb == "requires" {
				s.Kind = "Dependency"
				s.Interface = target
			} else {
				s.Kind = "Allocation"
				s.Container = target
			}
		case "projects":
			p.take()
			s.Kind = "Projection"
			s.Dataset = p.identifier()
			p.expect("into")
			s.Datagram = p.identifier()
		case "places":
			p.take()
			s.Kind = "Placement"
			s.Field = p.identifier()
			p.expect("at")
			s.Offset = p.integer()
			p.expect("bits")
			s.Width = p.integer()
		case "uses":
			p.take()
			s.Kind = "Participation"
			s.Channel = p.identifier()
			p.expect("as")
			s.Role = p.current().text
			if s.Role != "sender" && s.Role != "receiver" {
				p.fail("UNSUPPORTED_SYNTAX", "Expected sender or receiver")
			}
			p.take()
			p.expect("of")
			s.Message = p.identifier()
			p.expect("in")
			p.expect("mode")
			s.Mode = p.identifier()
		case "step":
			p.take()
			s.Kind = "Step"
			s.Ordinal = p.integer()
			p.expect("sends")
			s.Message = p.identifier()
			if p.current().text == "variant" {
				p.take()
				v := p.identifier()
				s.Variant = &v
			}
			p.expect("from")
			s.Sender = p.identifier()
			p.expect("to")
			s.Receiver = p.identifier()
			p.expect("via")
			s.Channel = p.identifier()
			if p.current().text == "reply-to" {
				p.take()
				v := p.integer()
				s.ReplyTo = &v
			}
		default:
			p.fail("UNSUPPORTED_SYNTAX", "Unknown relation")
		}
	}
	s.Span = p.finish(s.Subject.Span.Start)
	return s
}
func Check(text string) (*Model, []Diagnostic) {
	m, err := Parse(text)
	if err != nil {
		return nil, []Diagnostic{err.(Diagnostic)}
	}
	diagnostics := Validate(m)
	if len(diagnostics) == 0 {
		canonical, _ := Canonical(m)
		if text != canonical {
			offset := 0
			for offset < len(text) && offset < len(canonical) && text[offset] == canonical[offset] {
				offset++
			}
			diagnostics = append(diagnostics, Diagnostic{"NONCANONICAL_FORM", "Use canonical ordering, spacing and LF lines with one final newline", newSource(text).span(offset, min(offset+1, len(text)))})
		}
	}
	return m, diagnostics
}
