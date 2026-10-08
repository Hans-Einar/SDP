package markdown

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

var svgNumber = regexp.MustCompile(`^[+-]?(?:[0-9]+(?:\.[0-9]*)?|\.[0-9]+)(?:[eE][+-]?[0-9]+)?`)

func finite(v float64) bool { return !math.IsNaN(v) && !math.IsInf(v, 0) }
func numbers(s string) ([]float64, error) {
	var out []float64
	s = strings.TrimSpace(s)
	for s != "" {
		token := svgNumber.FindString(s)
		if token == "" {
			return nil, fmt.Errorf("invalid SVG number")
		}
		v, e := strconv.ParseFloat(token, 64)
		if e != nil || !finite(v) {
			return nil, fmt.Errorf("nonfinite SVG number")
		}
		out = append(out, v)
		s = s[len(token):]
		before := len(s)
		s = strings.TrimLeft(s, " \t\r\n")
		if strings.HasPrefix(s, ",") {
			s = strings.TrimLeft(s[1:], " \t\r\n")
			if s == "" || s[0] == ',' {
				return nil, fmt.Errorf("invalid SVG separator")
			}
		} else if len(s) > 0 && before == len(s) && s[0] != '-' && s[0] != '+' && s[0] != '.' {
			return nil, fmt.Errorf("missing SVG separator")
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("missing SVG number")
	}
	return out, nil
}
func scalar(s string) (float64, error) {
	v, e := numbers(s)
	if e != nil || len(v) != 1 {
		return 0, fmt.Errorf("expected one finite SVG number")
	}
	return v[0], nil
}

type matrix [6]float64

var unitMatrix = matrix{1, 0, 0, 1, 0, 0}

func multiply(a, b matrix) (matrix, error) {
	c := matrix{a[0]*b[0] + a[2]*b[1], a[1]*b[0] + a[3]*b[1], a[0]*b[2] + a[2]*b[3], a[1]*b[2] + a[3]*b[3], a[0]*b[4] + a[2]*b[5] + a[4], a[1]*b[4] + a[3]*b[5] + a[5]}
	for _, v := range c {
		if !finite(v) {
			return c, fmt.Errorf("nonfinite composed SVG transform")
		}
	}
	return c, nil
}
func transformed(m matrix, x, y float64) error {
	if !finite(x) || !finite(y) || !finite(m[0]*x+m[2]*y+m[4]) || !finite(m[1]*x+m[3]*y+m[5]) {
		return fmt.Errorf("nonfinite transformed SVG geometry")
	}
	return nil
}
func transforms(s string) (matrix, error) {
	result := unitMatrix
	for strings.TrimSpace(s) != "" {
		s = strings.TrimSpace(s)
		end := strings.IndexByte(s, '(')
		if end < 1 {
			return result, fmt.Errorf("invalid SVG transform")
		}
		name := strings.TrimSpace(s[:end])
		s = s[end+1:]
		end = strings.IndexByte(s, ')')
		if end < 0 {
			return result, fmt.Errorf("unterminated SVG transform")
		}
		v, e := numbers(s[:end])
		if e != nil {
			return result, e
		}
		s = s[end+1:]
		m := unitMatrix
		switch name {
		case "matrix":
			if len(v) != 6 {
				return result, fmt.Errorf("matrix requires six operands")
			}
			copy(m[:], v)
		case "translate":
			if len(v) < 1 || len(v) > 2 {
				return result, fmt.Errorf("translate operands")
			}
			m[4] = v[0]
			if len(v) == 2 {
				m[5] = v[1]
			}
		case "scale":
			if len(v) < 1 || len(v) > 2 {
				return result, fmt.Errorf("scale operands")
			}
			m[0] = v[0]
			m[3] = v[0]
			if len(v) == 2 {
				m[3] = v[1]
			}
		case "rotate":
			if len(v) != 1 && len(v) != 3 {
				return result, fmt.Errorf("rotate operands")
			}
			angle := math.Mod(v[0], 360) * math.Pi / 180
			c, sn := math.Cos(angle), math.Sin(angle)
			m = matrix{c, sn, -sn, c, 0, 0}
			if len(v) == 3 {
				m[4] = v[1] - c*v[1] + sn*v[2]
				m[5] = v[2] - sn*v[1] - c*v[2]
			}
		case "skewX", "skewY":
			if len(v) != 1 || math.Abs(math.Mod(v[0], 180)) == 90 {
				return result, fmt.Errorf("invalid skew")
			}
			t := math.Tan(math.Mod(v[0], 180) * math.Pi / 180)
			if name == "skewX" {
				m[2] = t
			} else {
				m[1] = t
			}
		default:
			return result, fmt.Errorf("unsupported SVG transform %s", name)
		}
		result, e = multiply(result, m)
		if e != nil {
			return result, e
		}
		s = strings.TrimSpace(s)
		if strings.HasPrefix(s, ",") {
			s = strings.TrimSpace(s[1:])
			if s == "" {
				return result, fmt.Errorf("trailing transform separator")
			}
		}
	}
	return result, nil
}

// validatePath checks command arity, arc flags and finite accumulated/control
// coordinates. It does not ask a permissive backend parser to validate grammar.
func validatePath(s string, m matrix) error {
	s = strings.TrimSpace(s)
	if s == "" {
		return fmt.Errorf("empty SVG path")
	}
	x, y, sx, sy := 0., 0., 0., 0.
	first := true
	for s != "" {
		cmd := s[0]
		upper := cmd
		relative := cmd >= 'a' && cmd <= 'z'
		if relative {
			upper -= 32
		}
		s = s[1:]
		counts := map[byte]int{'M': 2, 'L': 2, 'H': 1, 'V': 1, 'C': 6, 'S': 4, 'Q': 4, 'T': 2, 'A': 7, 'Z': 0}
		count, ok := counts[upper]
		if !ok || first && upper != 'M' {
			return fmt.Errorf("invalid SVG path command")
		}
		first = false
		if count == 0 {
			x, y = sx, sy
			s = strings.TrimSpace(s)
			continue
		}
		end := 0
		for end < len(s) {
			c := s[end]
			if (c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z') && c != 'e' && c != 'E' {
				break
			}
			end++
		}
		if upper == 'A' {
			// Arc flags are the literal characters 0/1, not numbers that merely
			// round to those values (1.0, 1e0, +1 and 01 are not flags).
			raw := strings.TrimSpace(s[:end])
			for index := 0; raw != ""; index++ {
				token := svgNumber.FindString(raw)
				if token == "" {
					return fmt.Errorf("invalid SVG arc operand")
				}
				if (index%7 == 3 || index%7 == 4) && token != "0" && token != "1" {
					return fmt.Errorf("invalid SVG arc flag")
				}
				raw = strings.TrimLeft(raw[len(token):], " ,\t\r\n")
			}
		}
		v, e := numbers(s[:end])
		if e != nil || len(v)%count != 0 {
			return fmt.Errorf("invalid SVG path operands")
		}
		s = strings.TrimSpace(s[end:])
		for i := 0; i < len(v); i += count {
			a := v[i : i+count]
			ox, oy := 0., 0.
			if relative {
				ox, oy = x, y
			}
			switch upper {
			case 'H':
				x = a[0] + ox
			case 'V':
				y = a[0] + oy
			case 'A':
				if a[0] < 0 || a[1] < 0 || (a[3] != 0 && a[3] != 1) || (a[4] != 0 && a[4] != 1) {
					return fmt.Errorf("invalid SVG arc")
				}
				if e := transformed(m, a[0], a[1]); e != nil {
					return e
				}
				x, y = a[5]+ox, a[6]+oy
			default:
				for j := 0; j < count; j += 2 {
					if e := transformed(m, a[j]+ox, a[j+1]+oy); e != nil {
						return e
					}
				}
				x, y = a[count-2]+ox, a[count-1]+oy
			}
			if e := transformed(m, x, y); e != nil {
				return e
			}
			if upper == 'M' && i == 0 {
				sx, sy = x, y
			}
		}
	}
	return nil
}
