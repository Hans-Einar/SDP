package markdown

import (
	"fmt"
	"math"
	"regexp"
	"strings"

	"golang.org/x/image/colornames"
)

var xmlDeclaration = regexp.MustCompile(`^version\s*=\s*(?:"1\.0"|'1\.0')(?:\s+encoding\s*=\s*(?:"UTF-8"|'UTF-8'|"utf-8"|'utf-8'))?(?:\s+standalone\s*=\s*(?:"yes"|'yes'|"no"|'no'))?\s*$`)
var hexPaint = regexp.MustCompile(`^#(?:[0-9a-fA-F]{3}|[0-9a-fA-F]{6})$`)
var localID = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_.-]*$`)
var localMarker = regexp.MustCompile(`^url\(#[A-Za-z_][A-Za-z0-9_.-]*\)$`)

var shapeAttrs = map[string]string{
	"svg": "viewBox xmlns width height", "g": "", "path": "d", "rect": "x y width height rx ry", "circle": "cx cy r", "ellipse": "cx cy rx ry", "line": "x1 y1 x2 y2", "polyline": "points", "polygon": "points", "title": "", "desc": "",
}
var mermaidAttrs = map[string]string{
	"defs": "", "marker": "id viewBox refX refY markerUnits markerWidth markerHeight orient", "text": "x y dx dy text-anchor font-size font-family font-weight dominant-baseline", "tspan": "x y dx dy text-anchor font-size font-family font-weight dominant-baseline",
}

func allowed(set, key string) bool { return strings.Contains(" "+set+" ", " "+key+" ") }
func validateElement(name string, a map[string]string, m matrix, mermaid bool) error {
	attrs, ok := shapeAttrs[name]
	if !ok && mermaid {
		attrs, ok = mermaidAttrs[name]
	}
	if !ok {
		return fmt.Errorf("unsupported element %s", name)
	}
	common := "transform fill stroke stroke-width opacity fill-opacity stroke-opacity"
	if name == "title" || name == "desc" {
		common = ""
	}
	if mermaid {
		common += " id stroke-linecap stroke-linejoin stroke-dasharray marker-start marker-end"
	}
	for k, v := range a {
		if !allowed(attrs+" "+common, k) {
			return fmt.Errorf("unsupported %s attribute %s", name, k)
		}
		switch k {
		case "xmlns":
			if v != svgNamespace {
				return fmt.Errorf("invalid SVG namespace")
			}
		case "viewBox":
			values, e := numbers(v)
			if e != nil || len(values) != 4 || values[2] <= 0 || values[3] <= 0 || values[2] > 32768 || values[3] > 32768 {
				return fmt.Errorf("invalid viewBox")
			}
		case "transform": // checked including ancestor composition
		case "d":
			if e := validatePath(v, m); e != nil {
				return e
			}
		case "points":
			values, e := numbers(v)
			if e != nil || len(values)%2 != 0 || len(values) < 4 {
				return fmt.Errorf("invalid points")
			}
			for i := 0; i < len(values); i += 2 {
				if e := transformed(m, values[i], values[i+1]); e != nil {
					return e
				}
			}
		case "fill", "stroke":
			if !solidPaint(v, mermaid) {
				return fmt.Errorf("unsupported solid paint")
			}
		case "opacity", "fill-opacity", "stroke-opacity":
			n, e := scalar(v)
			if e != nil || n < 0 || n > 1 {
				return fmt.Errorf("invalid opacity")
			}
		case "id":
			if !localID.MatchString(v) {
				return fmt.Errorf("invalid local identifier")
			}
		case "marker-start", "marker-end":
			if !localMarker.MatchString(v) {
				return fmt.Errorf("invalid local marker reference")
			}
		case "markerUnits":
			if v != "userSpaceOnUse" && v != "strokeWidth" {
				return fmt.Errorf("unsupported marker units")
			}
		case "orient":
			if v != "auto" && v != "auto-start-reverse" {
				if _, e := scalar(v); e != nil {
					return e
				}
			}
		case "text-anchor":
			if !allowed("start middle end", v) {
				return fmt.Errorf("invalid text anchor")
			}
		case "font-family":
			if strings.ContainsAny(v, "(){};\\") || len(v) > 256 {
				return fmt.Errorf("unsupported font family")
			}
		case "font-weight":
			if !allowed("normal bold 100 200 300 400 500 600 700 800 900", v) {
				return fmt.Errorf("unsupported font weight")
			}
		case "dominant-baseline":
			if !allowed("auto middle central hanging alphabetic", v) {
				return fmt.Errorf("unsupported baseline")
			}
		case "stroke-linecap":
			if !allowed("butt round square", v) {
				return fmt.Errorf("invalid line cap")
			}
		case "stroke-linejoin":
			if !allowed("miter round bevel", v) {
				return fmt.Errorf("invalid line join")
			}
		case "stroke-dasharray":
			if v != "none" {
				vs, e := numbers(v)
				if e != nil {
					return e
				}
				for _, n := range vs {
					if n < 0 {
						return fmt.Errorf("negative dash")
					}
				}
			}
		default:
			n, e := scalar(v)
			if e != nil {
				return e
			}
			if allowed("width height rx ry r stroke-width markerWidth markerHeight font-size", k) && n < 0 {
				return fmt.Errorf("negative dimension")
			}
			if k == "stroke-width" && (!finite(math.Hypot(m[0]*n, m[1]*n)) || !finite(math.Hypot(m[2]*n, m[3]*n))) {
				return fmt.Errorf("nonfinite transformed stroke width")
			}
		}
	}
	return validatePrimitive(name, a, m)
}
func solidPaint(v string, mermaid bool) bool {
	if v == "none" {
		return true
	}
	if _, ok := colornames.Map[v]; ok {
		return true
	}
	if hexPaint.MatchString(v) {
		return true
	}
	prefix := "rgb("
	count := 3
	if mermaid && strings.HasPrefix(v, "rgba(") {
		prefix = "rgba("
		count = 4
	}
	if !strings.HasPrefix(v, prefix) || !strings.HasSuffix(v, ")") {
		return false
	}
	ns, e := numbers(v[len(prefix) : len(v)-1])
	if e != nil || len(ns) != count {
		return false
	}
	for i, n := range ns {
		limit := 255.
		if i == 3 {
			limit = 1
		}
		if n < 0 || n > limit {
			return false
		}
	}
	return true
}
