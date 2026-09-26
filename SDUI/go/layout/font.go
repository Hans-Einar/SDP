package layout

import (
	"math"
	"strings"
	"sync"

	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
)

// TextMetrics uses embedded Go Regular at 72 DPI: font=10 means ten logical
// screen units, independent of window size. Device scaling belongs to the host.
// Faces are per call because font.Face implementations need not be concurrent-safe.
type TextMetrics struct{}

var once sync.Once
var typeface *opentype.Font

func Face(size float64) (font.Face, error) {
	once.Do(func() { typeface, _ = opentype.Parse(goregular.TTF) })
	return opentype.NewFace(typeface, &opentype.FaceOptions{Size: size, DPI: 72, Hinting: font.HintingNone})
}
func TextWidth(text string, size float64) float64 {
	f, err := Face(size)
	if err != nil {
		return 0
	}
	defer f.Close()
	return float64(font.MeasureString(f, text)) / 64
}
func Lines(text string, width, size float64) []string {
	f, err := Face(size)
	if err != nil {
		return []string{text}
	}
	defer f.Close()
	lines := []string{}
	for _, line := range strings.Split(strings.ReplaceAll(text, "\t", "    "), "\n") {
		if width <= 0 {
			lines = append(lines, line)
			continue
		}
		current := ""
		for _, r := range line {
			candidate := current + string(r)
			if current != "" && float64(font.MeasureString(f, candidate))/64 > width {
				lines = append(lines, current)
				current = string(r)
			} else {
				current = candidate
			}
		}
		lines = append(lines, current)
	}
	return lines
}
func (TextMetrics) Measure(n *parser.Instance, size, width float64) (Size, error) {
	if size <= 0 || size > 512 {
		return Size{}, diag(n, "font-range", "Font must be in (0,512] logical units")
	}
	text := n.Text
	if n.Kind == "widget" {
		text = n.Argument("label")
		if n.Widget == "input" {
			text = n.Argument("value")
			if text == "" {
				text = n.Argument("text")
			}
		}
	}
	if n.Kind == "widget" {
		switch n.Widget {
		case "button":
			return Size{math.Max(48, TextWidth(text, size)+24), math.Max(28, size*1.5+12)}, nil
		case "input":
			return Size{math.Max(120, TextWidth(text, size)+20), math.Max(28, size*1.5+12)}, nil
		case "svg":
			return Size{160, 90}, nil
		}
	}
	lines := Lines(text, width, size)
	w := 0.
	for _, line := range lines {
		w = math.Max(w, TextWidth(line, size))
	}
	return Size{w, float64(len(lines)) * size * 1.4}, nil
}
