// Package svg exports a measured SDUI tree. Export is static and has no domain callbacks.
package svg

import (
	"fmt"
	"html"
	"strings"

	"github.com/Hans-Einar/SDP/SDUI/go/layout"
)

// ContentRenderer can replace raw Markdown rendering without changing geometry.
type ContentRenderer interface {
	Render(*strings.Builder, *layout.Box) error
}
type Options struct {
	Width, Height float64
	Content       ContentRenderer
	SkipControls  bool
}

func Render(root *layout.Box, options Options) (string, error) {
	var out strings.Builder
	fmt.Fprintf(&out, `<svg xmlns="http://www.w3.org/2000/svg" width="%g" height="%g" viewBox="0 0 %g %g" role="img"><title>SDUI layout</title>`+"\n", options.Width, options.Height, options.Width, options.Height)
	out.WriteString(`<rect width="100%" height="100%" fill="#f1f5f9"/>` + "\n")
	var renderErr error
	index := 0
	root.Walk(func(b *layout.Box) {
		if renderErr != nil || b.Rect.W <= 0 || b.Rect.H <= 0 {
			return
		}
		index++
		n := b.Instance
		r := b.Rect
		c := b.Clip
		fmt.Fprintf(&out, `<defs><clipPath id="c%d"><rect x="%g" y="%g" width="%g" height="%g"/></clipPath></defs><g data-path="%s" data-kind="%s" clip-path="url(#c%d)">`+"\n", index, c.X, c.Y, c.W, c.H, html.EscapeString(b.Path), html.EscapeString(n.Kind), index)
		switch {
		case n.Variant == "box":
			rect(&out, r, "#ffffff", "#94a3b8", 4)
		case n.Kind == "widget":
			if options.SkipControls && n.Widget != "svg" {
				break
			}
			fill, stroke, color := "#e2e8f0", "#94a3b8", "#0f172a"
			if !b.Enabled {
				color = "#64748b"
				fill = "#f1f5f9"
			}
			label := n.Argument("label")
			switch n.Widget {
			case "button":
				rect(&out, r, fill, stroke, 5)
				x := r.X + (r.W-layout.TextWidth(label, b.Font))/2
				renderErr = Text(&out, label, x, r.Y+(r.H+b.Font*.7)/2, b.Font, color)
			case "input":
				rect(&out, r, "#fff", stroke, 3)
				label = n.Argument("value")
				if label == "" {
					label = n.Argument("text")
				}
				renderErr = Text(&out, label, r.X+8, r.Y+(r.H+b.Font*.7)/2, b.Font, color)
			case "svg":
				rect(&out, r, "#f8fafc", stroke, 0)
				renderErr = Text(&out, "SVG: "+label, r.X+4, r.Y+b.Font*1.3, b.Font, color)
			}
		case n.Kind == "markdown":
			if options.Content != nil {
				renderErr = options.Content.Render(&out, b)
			} else {
				for i, line := range layout.Lines(n.Text, r.W, b.Font) {
					if err := Text(&out, line, r.X, r.Y+b.Font+float64(i)*b.Font*1.4, b.Font, "#0f172a"); err != nil {
						renderErr = err
						break
					}
				}
			}
		}
		out.WriteString("</g>\n")
		if out.Len() > 16<<20 {
			renderErr = fmt.Errorf("SVG export exceeds 16 MiB")
		}
	})
	out.WriteString("</svg>\n")
	if renderErr != nil {
		return "", renderErr
	}
	return out.String(), nil
}
func rect(out *strings.Builder, r layout.Rect, fill, stroke string, radius float64) {
	fmt.Fprintf(out, `<rect x="%g" y="%g" width="%g" height="%g" rx="%g" fill="%s" stroke="%s"/>`+"\n", r.X, r.Y, r.W, r.H, radius, fill, stroke)
}
