package parser

import "strings"

type reader struct {
	tokens       []token
	source       string
	index, nodes int
}

func (p *reader) t() token { return p.tokens[p.index] }
func (p *reader) ahead(k string) bool {
	return p.index+1 < len(p.tokens) && p.tokens[p.index+1].kind == k
}
func (p *reader) take(k string, value ...string) token {
	t := p.t()
	if t.kind != k || len(value) > 0 && t.value != value[0] {
		fail("syntax", "Expected "+k, t.span)
	}
	p.index++
	return t
}
func (p *reader) accept(k string) bool {
	if p.t().kind == k {
		p.take(k)
		return true
	}
	return false
}
func (p *reader) ident() token {
	t := p.take("ID")
	s := t.value.(string)
	if strings.Contains(s, "-") || s == "sdui" || s == "ref" || s == "true" || s == "false" || s == "null" || s == "setHandle" {
		fail("syntax", "Reserved or invalid identifier", t.span)
	}
	return t
}
func (p *reader) span(s Span) Span { s.End = p.tokens[p.index-1].span.End; return s }

// Parse constructs a source-positioned AST, without resolving names or opening ref sources.
func Parse(source string) (doc *Document, err error) {
	defer recoverDiagnostic(&err)
	p := reader{tokens: lex(source), source: source}
	return p.document(), nil
}
func (p *reader) document() *Document {
	start := p.take("ID", "sdui").span
	v := p.take("NUMBER")
	if p.source[v.span.Start:v.span.End] != "0.2" {
		fail("version", "Only exact sdui 0.2 is supported", v.span)
	}
	p.take(";")
	d := &Document{Profile: "sdui/0.2"}
	for p.t().value == "ref" {
		s := p.take("ID").span
		p.take(":")
		a := p.ident().value.(string)
		path := p.take("STRING").value.(string)
		p.take(";")
		d.References = append(d.References, ModuleRef{a, path, p.span(s)})
	}
	for p.t().kind == "ID" && p.ahead("=") {
		s := p.t().span
		n := p.ident().value.(string)
		p.take("=")
		root := p.node(1, "")
		p.take(";")
		d.Definitions = append(d.Definitions, Definition{n, root, p.span(s)})
	}
	if len(d.Definitions) == 0 {
		fail("syntax", "At least one definition is required", p.t().span)
	}
	for p.t().kind != "EOF" {
		s := p.t().span
		m := p.ident().value.(string)
		p.take(".")
		o := p.ident().value.(string)
		p.take(".")
		p.take("ID", "setHandle")
		p.take("(")
		def := p.ident().value.(string)
		path := []string{}
		for p.accept(".") {
			path = append(path, p.ident().value.(string))
		}
		if len(path) == 0 {
			fail("syntax", "Expected component path", p.t().span)
		}
		p.take(")")
		p.take(";")
		d.Connections = append(d.Connections, Connection{m, o, def, path, p.span(s)})
	}
	end := p.take("EOF")
	start.End = end.span.End
	d.Span = start
	return d
}
func (p *reader) node(depth int, parent string) *Node {
	if depth > 64 {
		fail("depth-limit", "Component depth exceeds 64", p.t().span)
	}
	p.nodes++
	if p.nodes > 2048 {
		fail("node-limit", "More than 2048 components", p.t().span)
	}
	n := &Node{Span: p.t().span}
	if p.t().kind == "ID" && p.ahead("=") {
		name := p.ident().value.(string)
		p.take("=")
		if parent == "frame" && (name == "header" || name == "body" || name == "footer") {
			n.Role = str(name)
		} else {
			n.Name = str(name)
		}
	}
	switch p.t().kind {
	case "[", "<":
		open := p.t().kind
		p.take(open)
		close := ">"
		n.Kind = "group"
		if open == "[" {
			close = "]"
			n.Kind = "frame"
		}
		n.Rows = p.rows(close, depth, n.Kind)
		p.take(close)
	case "STRING":
		n.Kind = "markdown"
		l := p.literal()
		n.Text = &l
	case "ID":
		name := p.ident().value.(string)
		if p.accept("(") {
			n.Kind = "widget"
			n.Widget = str(name)
			n.Arguments = p.arguments()
			p.take(")")
		} else {
			n.Kind = "use"
			n.Target = str(name)
		}
	default:
		fail("syntax", "Expected UI component", p.t().span)
	}
	if p.accept("*") {
		n.Variant = str(p.ident().value.(string))
	}
	if p.t().kind == "{" {
		n.Layout = p.layout()
	}
	n.Span = p.span(n.Span)
	return n
}
func (p *reader) rows(close string, depth int, parent string) []Row {
	rows := []Row{}
	if p.t().kind == close {
		return rows
	}
	for {
		s := p.t().span
		items := []*Node{p.node(depth+1, parent)}
		for p.accept(",") {
			items = append(items, p.node(depth+1, parent))
		}
		rows = append(rows, Row{items, p.span(s)})
		if !p.accept(";") || p.t().kind == close {
			break
		}
	}
	return rows
}
func (p *reader) arguments() []Argument {
	args := []Argument{}
	if p.t().kind == ")" {
		return args
	}
	for {
		a := Argument{Span: p.t().span}
		if p.t().kind == "ID" && p.ahead("=") {
			a.Name = str(p.ident().value.(string))
			p.take("=")
		}
		a.Value = p.value()
		a.Span = p.span(a.Span)
		args = append(args, a)
		if len(args) > 32 {
			fail("argument-limit", "More than 32 arguments", p.t().span)
		}
		if !p.accept(",") {
			return args
		}
	}
}
func (p *reader) value() any {
	if p.t().kind == "ID" && p.t().value != "true" && p.t().value != "false" && p.t().value != "null" {
		s := p.t().span
		m := p.ident().value.(string)
		p.take(".")
		o := p.ident().value.(string)
		p.take(".")
		p.take("@")
		member := p.ident().value.(string)
		return Reference{m, o, member, p.span(s)}
	}
	return p.literal()
}
func (p *reader) literal() Literal {
	t := p.t()
	switch t.kind {
	case "STRING":
		p.take("STRING")
		return Literal{"string", t.value, t.span}
	case "NUMBER":
		p.take("NUMBER")
		return Literal{"number", t.value, t.span}
	case "ID":
		if t.value == "true" || t.value == "false" {
			p.take("ID")
			return Literal{"boolean", t.value == "true", t.span}
		}
		if t.value == "null" {
			p.take("ID")
			return Literal{"null", nil, t.span}
		}
	}
	fail("syntax", "Expected literal", t.span)
	return Literal{}
}
func (p *reader) layoutValue() Literal {
	s := p.t().span
	if p.accept("(") {
		values := []float64{p.take("NUMBER").value.(float64)}
		for p.accept(",") {
			values = append(values, p.take("NUMBER").value.(float64))
			if len(values) > 4 {
				fail("syntax", "Padding has at most four values", p.t().span)
			}
		}
		p.take(")")
		return Literal{"tuple", values, p.span(s)}
	}
	if p.t().kind == "ID" && p.t().value != "true" && p.t().value != "false" && p.t().value != "null" {
		return Literal{"enum", p.take("ID").value, s}
	}
	v := p.literal()
	if v.Kind == "number" && p.t().value == "fr" {
		p.take("ID")
		v.Kind = "fr"
		v.Span = p.span(s)
	}
	return v
}
func (p *reader) layout() []LayoutRule {
	p.take("{")
	out := []LayoutRule{}
	if p.t().kind != "}" {
		for {
			s := p.t().span
			name := ""
			var v Literal
			if p.t().kind == "ID" && p.ahead("=") {
				name = p.take("ID").value.(string)
				p.take("=")
				v = p.layoutValue()
			} else if p.t().kind == "NUMBER" && p.ahead(":") {
				x := p.take("NUMBER").value.(float64)
				p.take(":")
				y := p.take("NUMBER").value.(float64)
				name = "ratio"
				v = Literal{"ratio", []float64{x, y}, p.span(s)}
			} else {
				if p.t().kind == "," || p.t().kind == "}" || p.t().kind == "EOF" {
					fail("syntax", "Expected formatting rule", p.t().span)
				}
				for p.t().kind != "," && p.t().kind != "}" && p.t().kind != "EOF" {
					p.index++
				}
				name = "arrow"
				v = Literal{"arrow", strings.Join(strings.Fields(p.source[s.Start:p.span(s).End]), ""), p.span(s)}
			}
			out = append(out, LayoutRule{name, v, p.source[s.Start:p.span(s).End], p.span(s)})
			if len(out) > 32 {
				fail("layout-limit", "More than 32 rules", p.t().span)
			}
			if !p.accept(",") {
				break
			}
		}
	}
	p.take("}")
	return out
}
