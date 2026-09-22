package svg

import (
	"fmt"
	"html"
	"strings"
	"sync"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

var once sync.Once
var regular *sfnt.Font

// Text emits outlines from the same embedded font used for layout. This avoids
// renderer font substitution. The literal text remains in accessible SVG titles.
func Text(out *strings.Builder, text string, x, y, size float64, color string) error {
	once.Do(func() { regular, _ = sfnt.Parse(goregular.TTF) })
	var buf sfnt.Buffer
	scale := fixed.Int26_6(size * 64)
	path := strings.Builder{}
	last := sfnt.GlyphIndex(0)
	for i, r := range text {
		g, err := regular.GlyphIndex(&buf, r)
		if err != nil {
			return err
		}
		if i > 0 {
			kern, _ := regular.Kern(&buf, last, g, scale, font.HintingNone)
			x += float64(kern) / 64
		}
		segments, err := regular.LoadGlyph(&buf, g, scale, nil)
		if err != nil {
			return err
		}
		for _, s := range segments {
			n := 1
			op := "M"
			switch s.Op {
			case sfnt.SegmentOpLineTo:
				op = "L"
			case sfnt.SegmentOpQuadTo:
				op = "Q"
				n = 2
			case sfnt.SegmentOpCubeTo:
				op = "C"
				n = 3
			}
			path.WriteString(op)
			for j := 0; j < n; j++ {
				fmt.Fprintf(&path, "%.3f %.3f ", x+float64(s.Args[j].X)/64, y+float64(s.Args[j].Y)/64)
			}
		}
		advance, err := regular.GlyphAdvance(&buf, g, scale, font.HintingNone)
		if err != nil {
			return err
		}
		x += float64(advance) / 64
		last = g
	}
	fmt.Fprintf(out, `<g role="img" aria-label="%s"><title>%s</title><path fill="%s" d="%s"/></g>`+"\n", html.EscapeString(text), html.EscapeString(text), color, path.String())
	return nil
}
