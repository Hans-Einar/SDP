package markdown

import (
	"fmt"
	"math"
)

// finiteArc checks the endpoint-to-center conversion and the actual swept
// extrema, in local and transformed coordinates. Radii are vectors, not points.
func finiteArc(m matrix, x, y, ex, ey, rx, ry, degrees float64, large, sweep bool) error {
	if e := transformed(m, ex, ey); e != nil {
		return e
	}
	if rx == 0 || ry == 0 || x == ex && y == ey {
		return nil
	} // SVG's line/empty-arc cases.
	phi := math.Mod(degrees, 360) * math.Pi / 180
	c, s := math.Cos(phi), math.Sin(phi)
	dx, dy := x/2-ex/2, y/2-ey/2
	px, py := c*dx+s*dy, -s*dx+c*dy
	if !finite(px) || !finite(py) {
		return fmt.Errorf("nonfinite arc endpoint conversion")
	}
	u, v := px/rx, py/ry
	scale := math.Hypot(u, v)
	if !finite(scale) {
		return fmt.Errorf("nonfinite arc radius correction")
	}
	if scale > 1 {
		rx *= scale
		ry *= scale
		u, v = px/rx, py/ry
	}
	norm := u*u + v*v
	if !finite(rx) || !finite(ry) || !finite(norm) || norm <= 0 {
		return fmt.Errorf("nonfinite arc center conversion")
	}
	factor := math.Sqrt(math.Max(0, (1-norm)/norm))
	if large == sweep {
		factor = -factor
	}
	cxp, cyp := (factor*v)*rx, (-factor*u)*ry
	cx, cy := c*cxp-s*cyp+x/2+ex/2, s*cxp+c*cyp+y/2+ey/2
	if !finite(cxp) || !finite(cyp) || !finite(cx) || !finite(cy) {
		return fmt.Errorf("nonfinite arc center")
	}
	// p(theta) = center + a*cos(theta) + b*sin(theta).
	ax, ay, bx, by := rx*c, rx*s, -ry*s, ry*c
	tax, tay := m[0]*ax+m[2]*ay, m[1]*ax+m[3]*ay
	tbx, tby := m[0]*bx+m[2]*by, m[1]*bx+m[3]*by
	for _, n := range []float64{ax, ay, bx, by, tax, tay, tbx, tby} {
		if !finite(n) {
			return fmt.Errorf("nonfinite arc axis")
		}
	}
	start := math.Atan2((py-cyp)/ry, (px-cxp)/rx)
	end := math.Atan2((-py-cyp)/ry, (-px-cxp)/rx)
	delta := end - start
	if sweep && delta < 0 {
		delta += 2 * math.Pi
	}
	if !sweep && delta > 0 {
		delta -= 2 * math.Pi
	}
	if !finite(start) || !finite(delta) {
		return fmt.Errorf("nonfinite arc sweep")
	}
	onSweep := func(a float64) bool {
		distance := a - start
		if !sweep {
			distance = -distance
		}
		distance = math.Mod(distance, 2*math.Pi)
		if distance < 0 {
			distance += 2 * math.Pi
		}
		return distance <= math.Abs(delta)+1e-12
	}
	// Check only extrema traversed by this arc; the unused side of its ellipse
	// must not create a false overflow rejection for a finite selected arc.
	for _, pair := range [][2]float64{{ax, bx}, {ay, by}, {tax, tbx}, {tay, tby}} {
		angle := math.Atan2(pair[1], pair[0])
		for _, a := range []float64{angle, angle + math.Pi} {
			if onSweep(a) {
				px, py := cx+ax*math.Cos(a)+bx*math.Sin(a), cy+ay*math.Cos(a)+by*math.Sin(a)
				if e := transformed(m, px, py); e != nil {
					return e
				}
			}
		}
	}
	return nil
}
