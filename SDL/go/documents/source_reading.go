package documents

import (
	"fmt"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	"strings"
)

// SourceReading describes syntax at its actual file span. A fragment is not
// checked in isolation and unresolved names never acquire invented kinds.
func SourceReading(m *parser.Model, name, binding string, checked bool) (string, []byte) {
	esc := func(s string) string {
		return strings.NewReplacer("\\", "\\\\", "|", "\\|", "[", "\\[", "]", "\\]", "\n", " ").Replace(s)
	}
	link := func(label string, sp parser.Span) string {
		if sp.End <= sp.Start {
			return esc(label)
		}
		return fmt.Sprintf("[%s](sdl-source://%s/%d/%d)", esc(label), binding, sp.Start, sp.End)
	}
	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n\n", esc(name))
	if checked {
		b.WriteString("Declarations and statements authored in this file, within the selected checked System context. Referenced declarations can live in its other source files.\n\n")
	} else {
		b.WriteString("Local syntax only. This fragment has no checked System context. Names declared elsewhere remain unresolved here; type, contract and relationship validity require a complete model.\n\n")
	}
	b.WriteString("![Local declarations](source.svg)\n\n## Declarations\n\n")
	for _, d := range m.Declarations {
		fmt.Fprintf(&b, "- %s — %s\n", link(d.Name.Name, d.Span), d.Kind)
	}
	if len(m.Declarations) == 0 {
		b.WriteString("No local declarations.\n")
	}
	b.WriteString("\n## Authored statements\n\n")
	for _, s := range m.Statements {
		fmt.Fprintf(&b, "- %s — line %d\n", link(s.Sentence(), s.Span), s.Span.Line)
	}
	if len(m.Statements) == 0 {
		b.WriteString("No local relationships or properties.\n")
	}
	if len(m.Includes) > 0 {
		b.WriteString("\n## Included sources\n\n")
		for _, s := range m.Includes {
			fmt.Fprintf(&b, "- %s\n", link(s.Path, s.Span))
		}
	}

	rows := (len(m.Declarations) + 2) / 3
	if rows < 1 {
		rows = 1
	}
	boxH := 65
	for _, d := range m.Declarations {
		boxH = max(boxH, 35+16*svgLineCount(d.Name.Name, 30))
	}
	relations := []parser.Statement{}
	for _, r := range m.Statements {
		if r.Kind == "Relation" {
			relations = append(relations, r)
		}
	}
	if len(relations) > 80 {
		fmt.Fprintf(&b, "\nThe source contains %d relation statements; use the complete source-linked statement list above for their detail. The diagram shows declarations.\n", len(relations))
		relations = nil
	}
	height := rows*(boxH+17) + 25 + len(relations)*80
	var svg strings.Builder
	fmt.Fprintf(&svg, `<svg xmlns="http://www.w3.org/2000/svg" width="900" height="%d" viewBox="0 0 900 %d"><defs><marker id="relation" markerWidth="8" markerHeight="8" refX="7" refY="4" orient="auto"><path d="M0,0 L8,4 L0,8 Z" fill="#48677d"/></marker></defs><rect width="900" height="100%%" fill="#fff"/>`, height, height)
	for i, d := range m.Declarations {
		x := 15 + (i%3)*295
		y := 15 + (i/3)*(boxH+17)
		fmt.Fprintf(&svg, `<rect x="%d" y="%d" width="280" height="%d" rx="6" fill="#edf4fa" stroke="#48677d"/>`, x, y, boxH)
		svgText(&svg, x+10, y+20, 35, 13, d.Kind, "#35424c")
		svgText(&svg, x+10, y+43, 30, 14, d.Name.Name, "#132d42")
	}
	if len(m.Declarations) == 0 {
		svgText(&svg, 20, 40, 80, 16, "No local declarations; see authored statements below.", "#253d52")
	}
	for i, r := range relations {
		y := rows*(boxH+17) + 45 + i*80
		svgText(&svg, 20, y, 30, 13, r.Subject.Name, "#253d52")
		svgText(&svg, 610, y, 30, 13, r.Object.Name, "#253d52")
		fmt.Fprintf(&svg, `<path d="M310,%d L590,%d" stroke="#48677d" stroke-width="1.5" marker-end="url(#relation)"/>`, y+8, y+8)
		svgText(&svg, 365, y-5, 28, 12, r.Verb, "#35424c")
	}
	svg.WriteString("</svg>")
	return b.String(), []byte(svg.String())
}
