package markdown

import (
	"encoding/base64"
	"fmt"
	"html"
	"math"
	"strings"

	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/svg"
)

type Resource struct {
	SVG           []byte
	Width, Height float64
}
type Renderer interface {
	Render(source string) (Resource, error)
}
type Provider struct {
	Documents map[string]*Document
	Resources map[string]Resource
}

// Prepare performs optional registered renderer I/O before measurement. No
// rendering or filesystem access occurs during Measure or Render.
func Prepare(root *parser.Instance, renderer Renderer) (*Provider, error) {
	p := &Provider{Documents: map[string]*Document{}, Resources: map[string]Resource{}}
	var failure error
	root.Walk(func(n *parser.Instance) {
		if failure != nil || n.Kind != "markdown" {
			return
		}
		if _, ok := p.Documents[n.Text]; ok {
			return
		}
		d, err := Parse(n.Text)
		if err != nil {
			failure = &parser.Diagnostic{Code: "markdown-profile", Message: err.Error(), Span: n.Span}
			return
		}
		p.Documents[n.Text] = d
		for _, diagram := range d.Diagrams {
			if renderer == nil {
				continue
			}
			if _, ok := p.Resources[diagram.ID]; ok {
				continue
			}
			r, err := renderer.Render(diagram.Source)
			if err != nil {
				failure = err
				return
			}
			p.Resources[diagram.ID] = r
		}
	})
	return p, failure
}

type item struct {
	text       string
	x, y, font float64
	diagram    string
	w, h       float64
	strong     bool
	cell       bool
}

func (p *Provider) items(source string, font, width float64) ([]item, layout.Size, error) {
	d := p.Documents[source]
	if d == nil {
		return nil, layout.Size{}, fmt.Errorf("markdown-unprepared: source was not prepared")
	}
	items := []item{}
	y, maxW := 0., 0.
	for _, block := range d.Blocks {
		if block.Diagram != "" {
			r, ok := p.Resources[block.Diagram]
			if !ok {
				label := "[Mermaid: " + block.Diagram + " · provider ikke konfigurert]"
				for _, line := range layout.Lines(label, width, font) {
					items = append(items, item{text: line, y: y, font: font})
					y += font * 1.4
					maxW = math.Max(maxW, layout.TextWidth(line, font))
				}
			} else {
				w := math.Min(width, r.Width)
				h := w * r.Height / r.Width
				items = append(items, item{diagram: block.Diagram, y: y, w: w, h: h})
				y += h
				maxW = math.Max(maxW, w)
			}
		} else if len(block.Table) > 0 {
			cols := len(block.Table[0])
			if cols == 0 {
				continue
			}
			cw := width / float64(cols)
			for ri, row := range block.Table {
				height := font*1.4 + 8
				for _, cell := range row {
					height = math.Max(height, float64(len(layout.Lines(cell, math.Max(1, cw-8), font)))*font*1.4+8)
				}
				for ci, cell := range row {
					items = append(items, item{x: float64(ci) * cw, y: y, w: cw, h: height, cell: true})
					for li, line := range layout.Lines(cell, math.Max(1, cw-8), font) {
						items = append(items, item{text: line, x: float64(ci)*cw + 4, y: y + 4 + float64(li)*font*1.4, font: font, strong: ri == 0})
					}
				}
				y += height
			}
			maxW = math.Max(maxW, width)
		} else {
			f := font * block.Scale
			for _, line := range layout.Lines(block.Text, width, f) {
				items = append(items, item{text: line, y: y, font: f, strong: block.Strong})
				y += f * 1.4
				maxW = math.Max(maxW, layout.TextWidth(line, f))
			}
		}
		y += font * .35
	}
	if len(d.Blocks) > 0 {
		y -= font * .35
	}
	return items, layout.Size{W: maxW, H: y}, nil
}
func (p *Provider) Measure(n *parser.Instance, font, width float64) (layout.Size, error) {
	if font <= 0 || font > 512 {
		return layout.Size{}, fmt.Errorf("font-range: font must be in (0,512]")
	}
	if n.Kind != "markdown" {
		return (layout.TextMetrics{}).Measure(n, font, width)
	}
	_, size, err := p.items(n.Text, font, width)
	return size, err
}
func (p *Provider) Render(out *strings.Builder, b *layout.Box) error {
	items, _, err := p.items(b.Instance.Text, b.Font, b.Rect.W)
	if err != nil {
		return err
	}
	for _, i := range items {
		x, y := b.Rect.X+i.x, b.Rect.Y+i.y
		if i.diagram != "" {
			r := p.Resources[i.diagram]
			fmt.Fprintf(out, `<image data-resource="%s.svg" x="%g" y="%g" width="%g" height="%g" href="data:image/svg+xml;base64,%s"/>`, html.EscapeString(i.diagram), x, y, i.w, i.h, base64.StdEncoding.EncodeToString(r.SVG))
		} else if i.cell {
			fmt.Fprintf(out, `<rect x="%g" y="%g" width="%g" height="%g" fill="#f8fafc" stroke="#cbd5e1"/>`, x, y, i.w, i.h)
		} else {
			color := "#334155"
			if i.strong {
				color = "#0f172a"
			}
			if err := svg.Text(out, i.text, x, y+i.font, i.font, color); err != nil {
				return err
			}
		}
	}
	return nil
}
