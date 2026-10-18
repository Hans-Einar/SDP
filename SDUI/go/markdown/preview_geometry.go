package markdown

import (
	"fmt"
	"math"
	"strings"

	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/svg"
)

type previewBands struct {
	caption                                          string
	status                                           []string
	captionHeight, statusHeight, minWidth, textWidth float64
}

func textMinimum(s string, font float64) float64 {
	w := 0.
	for _, r := range strings.ReplaceAll(s, "\t", "    ") {
		if r != '\n' && r != '\r' {
			w = math.Max(w, layout.TextWidth(string(r), font))
		}
	}
	return w
}
func ellipsis(s string, font, width float64) string {
	s = strings.NewReplacer("\n", " ", "\r", " ", "\t", "    ").Replace(s)
	if layout.TextWidth(s, font) <= width {
		return s
	}
	rs := []rune(s)
	lo, hi := 0, len(rs)
	for lo < hi {
		mid := (lo + hi + 1) / 2
		if layout.TextWidth(string(rs[:mid])+"…", font) <= width {
			lo = mid
		} else {
			hi = mid - 1
		}
	}
	return string(rs[:lo]) + "…"
}
func bands(n *parser.Instance, o PreviewOutcome, font, width float64) (previewBands, error) {
	b := previewBands{}
	if e := previewMetrics(font, width); e != nil {
		return b, e
	}
	if caption := n.Argument("label"); caption != "" {
		b.caption = ellipsis(caption, font, width)
		b.captionHeight = font * 1.4
		b.minWidth = math.Min(layout.TextWidth(strings.NewReplacer("\n", " ", "\r", " ", "\t", "    ").Replace(caption), font), layout.TextWidth("…", font))
		b.textWidth = layout.TextWidth(b.caption, font)
	}
	if o.Status == "label" {
		text := unavailable(o.Description, o.Diagnostic)
		b.status = layout.Lines(text, width, font)
		b.statusHeight = float64(len(b.status)) * font * 1.4
		b.minWidth = math.Max(b.minWidth, textMinimum(text, font))
		for _, line := range b.status {
			b.textWidth = math.Max(b.textWidth, layout.TextWidth(line, font))
		}
	}
	return b, nil
}
func (p *Previews) MeasurePreview(n *parser.Instance, font, width float64) (natural, minimum layout.Size, err error) {
	if err = p.CheckPreview(n); err != nil {
		return
	}
	o := p.records[n.Path]
	if o.Kind != "svg" {
		err = fmt.Errorf("preview-geometry: MeasurePreview requires SVG")
		return
	}
	b, e := bands(n, o, font, width)
	if e != nil {
		err = e
		return
	}
	minimum = layout.Size{W: b.minWidth, H: b.captionHeight + b.statusHeight}
	natural = layout.Size{W: b.textWidth, H: minimum.H}
	if o.Resource != nil {
		w := math.Min(width, o.Resource.Width)
		natural.W = math.Max(natural.W, w)
		natural.H += (w / o.Resource.Width) * o.Resource.Height
	}
	if !finite(natural.H) || natural.H > 1e7 || natural.W > 1e7 {
		err = fmt.Errorf("preview-geometry: measured content exceeds layout bounds")
	}
	return
}
func (p *Previews) svgBands(box *layout.Box) (b previewBands, o PreviewOutcome, err error) {
	if box == nil || box.Instance == nil {
		err = fmt.Errorf("preview-geometry: nil box")
		return
	}
	if err = p.CheckPreview(box.Instance); err != nil {
		return
	}
	o = p.records[box.Instance.Path]
	if o.Kind != "svg" || box.Path != box.Instance.Path {
		err = fmt.Errorf("preview-geometry: incorrect SVG box")
		return
	}
	r := box.Rect
	if !finite(r.X) || !finite(r.Y) || !finite(r.W) || !finite(r.H) || r.W < 0 || r.H < 0 || r.W > 1e7 || r.H > 1e7 || !finite(r.X+r.W) || !finite(r.Y+r.H) {
		err = fmt.Errorf("preview-geometry: invalid rectangle")
		return
	}
	b, err = bands(box.Instance, o, box.Font, r.W)
	if err != nil {
		return
	}
	if r.W+1e-7 < b.minWidth || r.H+1e-7 < b.captionHeight+b.statusHeight {
		err = fmt.Errorf("preview-geometry: caption/status allocation is too small")
	}
	return
}

// SVGRects returns already fitted absolute canvas geometry. Clip is applied by
// the owner afterward; it never changes fitting or text wrapping.
func (p *Previews) SVGRects(box *layout.Box) (image, caption, status layout.Rect, err error) {
	b, o, e := p.svgBands(box)
	if e != nil {
		err = e
		return
	}
	r := box.Rect
	if b.captionHeight > 0 {
		caption = layout.Rect{X: r.X, Y: r.Y + r.H - b.captionHeight, W: r.W, H: b.captionHeight}
	}
	if b.statusHeight > 0 {
		status = layout.Rect{X: r.X, Y: r.Y, W: r.W, H: b.statusHeight}
	}
	if o.Resource != nil {
		image, err = layout.FitPreview(layout.Rect{X: r.X, Y: r.Y, W: r.W, H: math.Max(0, r.H-b.captionHeight)}, layout.Size{W: o.Resource.Width, H: o.Resource.Height})
	}
	return
}

// RenderSVGText emits only glyph paths for the caption and unavailable status.
// The host supplies the direct SVG root/viewBox and separate resource image.
func (p *Previews) RenderSVGText(out *strings.Builder, box *layout.Box) error {
	if out == nil {
		return fmt.Errorf("preview-render: nil output")
	}
	b, _, e := p.svgBands(box)
	if e != nil {
		return e
	}
	_, caption, status, e := p.SVGRects(box)
	if e != nil {
		return e
	}
	var fragment strings.Builder
	if b.captionHeight > 0 {
		if e = svg.Text(&fragment, b.caption, caption.X, caption.Y+box.Font, box.Font, "#334155"); e != nil {
			return e
		}
	}
	for i, line := range b.status {
		if e = svg.Text(&fragment, line, status.X, status.Y+box.Font+float64(i)*box.Font*1.4, box.Font, "#334155"); e != nil {
			return e
		}
	}
	out.WriteString(fragment.String())
	return nil
}
