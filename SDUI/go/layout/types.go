// Package layout measures SDUI instances in logical screen units. It never executes callbacks.
package layout

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/Hans-Einar/SDP/SDUI/go/parser"
)

type Size struct{ W, H float64 }
type Rect struct{ X, Y, W, H float64 }

func (r Rect) Contains(x, y float64) bool { return x >= r.X && y >= r.Y && x < r.X+r.W && y < r.Y+r.H }
func (r Rect) Intersect(s Rect) Rect {
	x, y := math.Max(r.X, s.X), math.Max(r.Y, s.Y)
	return Rect{x, y, math.Max(0, math.Min(r.X+r.W, s.X+s.W)-x), math.Max(0, math.Min(r.Y+r.H, s.Y+s.H)-y)}
}

// Box geometry is shared by SVG and native hosts. Clip applies to both paint and hit tests.
type Box struct {
	Instance *parser.Instance `json:"-"`
	Path     string           `json:"path"`
	Rect     Rect             `json:"rect"`
	Clip     Rect             `json:"clip"`
	Font     float64          `json:"font"`
	Enabled  bool             `json:"enabled"`
	Children []*Box           `json:"children,omitempty"`
}

func (b *Box) Walk(fn func(*Box)) {
	fn(b)
	for _, c := range b.Children {
		c.Walk(fn)
	}
}
func (b *Box) Hit(x, y float64) *Box {
	if !b.Clip.Contains(x, y) {
		return nil
	}
	for i := len(b.Children) - 1; i >= 0; i-- {
		if v := b.Children[i].Hit(x, y); v != nil {
			return v
		}
	}
	if b.Enabled && b.Instance.Kind == "widget" && b.Rect.Contains(x, y) {
		return b
	}
	return nil
}

type Measurer interface {
	Measure(*parser.Instance, float64, float64) (Size, error)
}
type Engine struct {
	Measure    Measurer
	operations int
}

func number(n *parser.Instance, key string, fallback float64) float64 {
	if v, ok := n.Layout[key].(float64); ok {
		return v
	}
	return fallback
}
func choice(n *parser.Instance, key, fallback string) string {
	if v, ok := n.Layout[key].(string); ok {
		return v
	}
	return fallback
}
func visible(n *parser.Instance) bool { return n.Layout["visible"] != false }
func weight(n *parser.Instance, axis string) float64 {
	s := choice(n, axis, "content")
	if s == "fill" {
		return 1
	}
	if strings.HasSuffix(s, "fr") {
		f, _ := strconv.ParseFloat(strings.TrimSuffix(s, "fr"), 64)
		return f
	}
	return 0
}
func scale(n *parser.Instance, axis string) (float64, bool) {
	if v, ok := n.Layout["scale-"+axis].(float64); ok {
		return v, true
	}
	v, ok := n.Layout["scale"].(float64)
	return v, ok
}
func explicit(n *parser.Instance, axis string) bool {
	_, s := scale(n, axis)
	return s || weight(n, axis) > 0
}
func diag(n *parser.Instance, code, message string) error {
	return &parser.Diagnostic{Code: code, Message: fmt.Sprintf("%s: %s", n.Path, message), Span: n.Span}
}
func (e *Engine) step(n *parser.Instance) error {
	e.operations++
	if e.operations > 200000 {
		return diag(n, "layout-budget", "Measurement budget exceeded")
	}
	return nil
}
func padding(n *parser.Instance, ref Size) [4]float64 {
	v := 0.0
	if n.Variant == "box" {
		v = .01
	}
	p := [4]float64{v, v, v, v}
	switch a := n.Layout["padding"].(type) {
	case float64:
		p = [4]float64{a, a, a, a}
	case []float64:
		copy(p[:], a)
	}
	p[0] *= ref.H
	p[2] *= ref.H
	p[1] *= ref.W
	p[3] *= ref.W
	return p
}
func gaps(n *parser.Instance, ref Size) (float64, float64) {
	g := number(n, "gap", .01)
	return number(n, "gap-x", g) * ref.W, number(n, "gap-y", g) * ref.H
}
func bounds(n *parser.Instance, axis string, ref float64) (float64, float64) {
	return number(n, "min-"+axis, 0) * ref, number(n, "max-"+axis, math.Inf(1)) * ref
}
func clamp(v, lo, hi float64) float64 { return math.Max(lo, math.Min(hi, v)) }
func alignment(available, used float64, kind string) float64 {
	switch kind {
	case "center":
		return (available - used) / 2
	case "end":
		return available - used
	}
	return 0
}
