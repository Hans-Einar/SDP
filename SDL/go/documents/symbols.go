package documents

import (
	"fmt"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/viewpoint"
	"html"
	"math"
	"strings"
)

const SymbolProfile = "sdl-symbols-1"

func semanticSVG(d viewpoint.Diagram, g graphGeometry) ([]byte, error) {
	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.3f" height="%.3f" viewBox="0 0 %.3f %.3f" data-profile="%s"><title>%s</title><rect width="100%%" height="100%%" fill="white"/><defs><marker id="sdl-arrow" viewBox="0 0 10 10" refX="10" refY="5" markerWidth="8" markerHeight="8" orient="auto-start-reverse"><path d="M0 0 L10 5 L0 10 Z" fill="#334155"/></marker><marker id="uml-dependency" viewBox="0 0 10 10" refX="10" refY="5" markerWidth="8" markerHeight="8" orient="auto-start-reverse"><path d="M0 0 L10 5 L0 10" fill="none" stroke="#334155"/></marker></defs>`, g.Width, g.Height, g.Width, g.Height, SymbolProfile, html.EscapeString(d.Title))
	nodes := map[string]nodeGeometry{}
	for _, n := range g.Nodes {
		nodes[n.ID] = n
	}
	for i, e := range g.Edges {
		e.Points = append([][2]float64{}, e.Points...)
		sourceKind := d.Nodes[e.From].(map[string]any)["kind"]
		targetKind := d.Nodes[e.To].(map[string]any)["kind"]
		if sourceKind == "actor" {
			n := nodes[e.From]
			e.Points = append([][2]float64{{n.X + n.Width/2, n.Y + 31}}, e.Points...)
		}
		if sourceKind == "usecase" {
			e.Points[0] = ellipseBoundary(nodes[e.From], e.Points[1])
		}
		if targetKind == "usecase" {
			i := len(e.Points) - 1
			e.Points[i] = ellipseBoundary(nodes[e.To], e.Points[i-1])
		}

		p := d.Edges[i]
		marker, dash := "sdl-arrow", ""
		if p.Relation == "consumes" {
			marker = "uml-dependency"
			dash = ` stroke-dasharray="6 4"`
		}
		fmt.Fprintf(&b, `<g data-fact="%s" data-relation="%s" data-source="%s" data-target="%s"><polyline fill="none" stroke="#334155" stroke-width="1.5" stroke-linejoin="round" marker-end="url(#%s)"%s points="`, p.Fact, p.Relation, p.Source, p.Target, marker, dash)
		for _, xy := range e.Points {
			fmt.Fprintf(&b, "%.3f,%.3f ", xy[0], xy[1])
		}
		b.WriteString(`"/>`)
		if e.Anchor != nil {
			fmt.Fprintf(&b, `<text x="%.3f" y="%.3f" text-anchor="middle" font-family="DejaVu Sans,sans-serif" font-size="13" fill="#111827" paint-order="stroke" stroke="white" stroke-width="3" stroke-linejoin="round">%s</text>`, e.Anchor[0], e.Anchor[1]+4, html.EscapeString(p.Relation))
		}
		b.WriteString(`</g>`)
	}
	for _, n := range g.Nodes {
		p := d.Nodes[n.ID].(map[string]any)
		kind := p["kind"].(string)
		fmt.Fprintf(&b, `<g id="%s" data-model-id="%s" data-kind="%s" transform="translate(%.3f %.3f)" stroke="#334155" stroke-width="1.5" fill="white">`, n.ID, p["model_id"], kind, n.X, n.Y)
		shape(&b, kind, n.Width, n.Height)
		font := 14.0
		line := 24.0
		scale := 1.0
		if kind == "usecase" {
			scale = math.Min(1, math.Min(n.Width*.78/math.Max(1, n.LabelWidth), n.Height*.70/math.Max(1, n.LabelHeight)))
			font *= scale
			line *= scale
		}
		y := n.Height/2 - (float64(len(n.Lines))-1)*line/2 + font*.35
		for _, text := range n.Lines {
			if text != "" {
				fmt.Fprintf(&b, `<text x="%.3f" y="%.3f" text-anchor="middle" stroke="none" fill="#111827" font-family="DejaVu Sans,sans-serif" font-size="%.3f">%s</text>`, n.Width/2, y, font, html.EscapeString(text))
			}
			y += line
		}
		b.WriteString(`</g>`)
	}
	b.WriteString(`</svg>`)
	data := []byte(b.String())
	return data, ValidateSVG(data)
}
func shape(b *strings.Builder, k string, w, h float64) {
	rect := func(rx float64, extra string) {
		fmt.Fprintf(b, `<rect width="%.3f" height="%.3f" rx="%.3f" %s/>`, w, h, rx, extra)
	}
	switch k {
	case "actor":
		cx := w / 2
		fmt.Fprintf(b, `<circle data-symbol="actor-head" cx="%.3f" cy="13" r="8"/><path data-symbol="actor-body" fill="none" d="M%.3f 21 V49 M%.3f 31 H%.3f M%.3f 49 L%.3f 65 M%.3f 49 L%.3f 65"/>`, cx, cx, cx-18, cx+18, cx, cx-15, cx, cx+15)
	case "usecase":
		fmt.Fprintf(b, `<ellipse data-symbol="usecase" cx="%.3f" cy="%.3f" rx="%.3f" ry="%.3f"/>`, w/2, h/2, w/2, h/2)
	case "feature":
		fmt.Fprintf(b, `<path data-symbol="feature-tab" d="M0 8 V0 H%.3f V8 H%.3f V%.3f H0 Z"/>`, math.Min(w/3, 70), w, h)
	case "functionality":
		rect(10, `data-symbol="functionality"`)
		fmt.Fprintf(b, `<path d="M8 8 V%.3f" stroke-width="3"/>`, h-8)
	case "capability":
		fmt.Fprintf(b, `<path data-symbol="capability" d="M16 0 H%.3f L%.3f %.3f L%.3f %.3f H16 L0 %.3f Z"/>`, w-16, w, h/2, w-16, h, h/2)
	case "activity":
		rect(20, `data-symbol="activity"`)
	case "mode":
		rect(3, `data-symbol="mode" stroke-dasharray="6 3"`)
	case "container":
		rect(0, `data-symbol="container" stroke-width="2"`)
		fmt.Fprintf(b, `<rect x="4" y="4" width="%.3f" height="%.3f"/>`, w-8, h-8)
	case "interface":
		rect(0, `data-symbol="interface"`)
		fmt.Fprintf(b, `<path d="M0 8 H8 V%.3f H0"/>`, h-8)
	case "database":
		fmt.Fprintf(b, `<path data-symbol="persistent-source" d="M0 8 C0 -2 %.3f -2 %.3f 8 V%.3f C%.3f %.3f 0 %.3f 0 %.3f Z M0 8 C0 18 %.3f 18 %.3f 8"/>`, w, w, h-8, w, h+2, h+2, h-8, w, w)
	default:
		rect(3, `data-symbol="typed-card"`)
	}
}

func ellipseBoundary(n nodeGeometry, p [2]float64) [2]float64 {
	cx, cy := n.X+n.Width/2, n.Y+n.Height/2
	dx, dy := p[0]-cx, p[1]-cy
	div := math.Sqrt(dx*dx/(n.Width*n.Width/4) + dy*dy/(n.Height*n.Height/4))
	if div == 0 {
		return [2]float64{cx - n.Width/2, cy}
	}
	return [2]float64{cx + dx/div, cy + dy/div}
}
