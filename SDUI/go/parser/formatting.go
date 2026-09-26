package parser

import (
	"fmt"
	"math"
	"strings"
)

var arrows = map[string]map[string]any{
	"^<": {"align-x": "start", "align-y": "start"}, "-^": {"align-x": "center", "align-y": "start"}, ">^": {"align-x": "end", "align-y": "start"},
	"<-": {"align-x": "start", "align-y": "center"}, "--": {"align-x": "center", "align-y": "center"}, "->": {"align-x": "end", "align-y": "center"},
	"v<": {"align-x": "start", "align-y": "end"}, "-v": {"align-x": "center", "align-y": "end"}, ">v": {"align-x": "end", "align-y": "end"},
	"<->": {"x": "fill"}, "^|v": {"y": "fill"}, ">-<": {"x": "content"}, ">|<": {"y": "content"},
}
var aliases = map[string]string{"<^": "^<", "^-": "-^", "^>": ">^", "-<": "<-", ">-": "->", "<v": "v<", "v-": "-v", "v>": ">v"}
var enums = map[string]string{"align-x": "start center end", "align-y": "start center end", "justify": "start center end between", "items": "start center end stretch", "overflow-x": "error clip scroll", "overflow-y": "error clip scroll", "wrap": "none wrap"}

func member(s, list string) bool {
	for _, v := range strings.Fields(list) {
		if v == s {
			return true
		}
	}
	return false
}
func formatting(n *Node, kind string, rowCount int) map[string]any {
	out := map[string]any{}
	for _, r := range n.Layout {
		key, v := r.Name, r.Value.Value
		add := map[string]any{}
		if key == "arrow" {
			shape, _ := v.(string)
			if canonical, ok := aliases[shape]; ok {
				shape = canonical
			}
			var ok bool
			add, ok = arrows[shape]
			if !ok {
				fail("layout-arrow", "Unknown arrow", r.Span)
			}
		} else {
			valid := false
			if options, ok := enums[key]; ok {
				s, _ := v.(string)
				valid = r.Value.Kind == "enum" && member(s, options)
			} else if member(key, "scale scale-x scale-y min-x min-y max-x max-y gap gap-x gap-y padding font") {
				num, ok := v.(float64)
				valid = ok && r.Value.Kind == "number" && num >= 0
				if member(key, "scale scale-x scale-y font") {
					valid = valid && num > 0
				}
				if key == "padding" && r.Value.Kind == "tuple" {
					xs, ok := v.([]float64)
					valid = ok && len(xs) == 4
					for _, x := range xs {
						valid = valid && x >= 0
					}
				}
			} else if key == "x" || key == "y" {
				if r.Value.Kind == "enum" {
					s, _ := v.(string)
					valid = member(s, "content fill")
				}
				if r.Value.Kind == "fr" {
					num, ok := v.(float64)
					valid = ok && num > 0
					v = fmt.Sprintf("%gfr", num)
				}
			} else if key == "ratio" {
				xs, ok := v.([]float64)
				valid = r.Value.Kind == "ratio" && ok && len(xs) == 2
				for _, x := range xs {
					valid = valid && x > 0
				}
			} else if key == "enabled" || key == "visible" {
				_, ok := v.(bool)
				valid = ok && r.Value.Kind == "boolean"
			} else {
				fail("layout-property", "Unknown formatting property "+key, r.Span)
			}
			if !valid {
				fail("layout-type", "Invalid value for "+key, r.Span)
			}
			if xs, ok := v.([]float64); ok {
				v = append([]float64(nil), xs...)
			}
			add[key] = v
		}
		for k, v := range add {
			if _, ok := out[k]; ok {
				fail("layout-conflict", "Multiple rules for "+k, r.Span)
			}
			out[k] = v
		}
	}
	checkCombination(out, kind, rowCount, n.Span)
	return out
}
func number(props map[string]any, key string, def float64) float64 {
	if v, ok := props[key].(float64); ok {
		return v
	}
	return def
}
func checkCombination(p map[string]any, kind string, rows int, s Span) {
	axes := 0
	for _, axis := range []string{"x", "y"} {
		count := 0
		for _, key := range []string{"scale", "scale-" + axis, axis} {
			if _, ok := p[key]; ok {
				count++
			}
		}
		if count > 1 {
			fail("layout-conflict", "Competing "+axis+" sizes", s)
		}
		if number(p, "min-"+axis, 0) > number(p, "max-"+axis, math.Inf(1)) {
			fail("size-range", "Minimum exceeds maximum", s)
		}
		_, scale := p["scale-"+axis]
		_, size := p[axis]
		if scale || size {
			axes++
		}
	}
	if _, ok := p["ratio"]; ok {
		if kind != "frame" {
			fail("ratio-scope", "Ratio applies to frames", s)
		}
		_, scale := p["scale"]
		if scale || axes > 1 {
			fail("ratio-axis", "Ratio has only one controlling axis", s)
		}
		for _, axis := range []string{"x", "y"} {
			if v, ok := p[axis]; ok && v != "fill" {
				fail("ratio-axis", "Ratio requires a scale or fill axis", s)
			}
		}
	}
	if _, ok := p["wrap"]; ok && kind != "group" {
		fail("wrap-scope", "Wrap applies to groups", s)
	}
	if p["wrap"] == "wrap" && (rows != 1 || p["x"] == "fill" || strings.HasSuffix(fmt.Sprint(p["x"]), "fr")) {
		fail("wrap-layout", "Wrap needs one row without horizontal fill/fr", s)
	}
}

// Formatting validates and canonicalizes the rules on an unexpanded component.
func Formatting(n *Node) (props map[string]any, err error) {
	defer recoverDiagnostic(&err)
	return formatting(n, n.Kind, len(n.Rows)), nil
}
