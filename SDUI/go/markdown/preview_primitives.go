package markdown

import (
	"fmt"
	"math"
)

// validatePrimitive runs after attribute grammar validation. Coordinates must be
// paired with their actual partner, not an unrelated scalar of the same value.
func validatePrimitive(name string, attributes map[string]string, m matrix) error {
	value := func(key string) float64 {
		if s, ok := attributes[key]; ok {
			v, _ := scalar(s)
			return v
		}
		return 0
	}
	switch name {
	case "line":
		if e := transformed(m, value("x1"), value("y1")); e != nil {
			return e
		}
		return transformed(m, value("x2"), value("y2"))
	case "rect":
		x, y, w, h := value("x"), value("y"), value("width"), value("height")
		for _, p := range [][2]float64{{x, y}, {x + w, y}, {x, y + h}, {x + w, y + h}} {
			if e := transformed(m, p[0], p[1]); e != nil {
				return e
			}
		}
	case "circle", "ellipse":
		rx, ry := value("rx"), value("ry")
		if name == "circle" {
			rx, ry = value("r"), value("r")
		}
		return finiteEllipse(m, value("cx"), value("cy"), rx, ry)
	case "svg", "marker":
		if raw, ok := attributes["viewBox"]; ok {
			v, _ := numbers(raw)
			for _, p := range [][2]float64{{v[0], v[1]}, {v[0] + v[2], v[1]}, {v[0], v[1] + v[3]}, {v[0] + v[2], v[1] + v[3]}} {
				if e := transformed(m, p[0], p[1]); e != nil {
					return e
				}
			}
		}
	case "text", "tspan":
		return transformed(m, value("x")+value("dx"), value("y")+value("dy"))
	}
	return nil
}
func finiteEllipse(m matrix, x, y, rx, ry float64) error {
	// Local extents are real geometry too, even if a later transform shrinks it.
	for _, v := range []float64{x - rx, x + rx, y - ry, y + ry} {
		if !finite(v) {
			return fmt.Errorf("nonfinite ellipse extent")
		}
	}
	if e := transformed(m, x, y); e != nil {
		return e
	}
	cx, cy := m[0]*x+m[2]*y+m[4], m[1]*x+m[3]*y+m[5]
	// Exact axis extents of the transformed ellipse; artificial bounding-box
	// corners would reject finite ellipses under some rotations/shears.
	ex, ey := math.Hypot(m[0]*rx, m[2]*ry), math.Hypot(m[1]*rx, m[3]*ry)
	for _, v := range []float64{ex, ey, cx - ex, cx + ex, cy - ey, cy + ey} {
		if !finite(v) {
			return fmt.Errorf("nonfinite transformed ellipse extent")
		}
	}
	return nil
}
