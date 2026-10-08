package layout

import (
	"math"

	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

type minimumKey struct {
	node *parser.Instance
	ref  Size
	font float64
}

// minimum solves the intrinsic active subtree size using that subtree's own
// finite body as the reference for descendants. Parent-relative padding/bounds
// still use ref. Caching is per layout run; no native state is retained.
func (e *Engine) minimum(n *parser.Instance, ref Size, font float64) (Size, error) {
	if err := e.step(n); err != nil {
		return Size{}, err
	}
	if !visible(n) || auxiliary(n) && n != e.surfaceRoot {
		return Size{}, nil
	}
	font = number(n, "font", font)
	key := minimumKey{n, ref, font}
	if v, ok := e.minima[key]; ok {
		return v, nil
	}
	loW, _ := bounds(n, "x", ref.W)
	loH, _ := bounds(n, "y", ref.H)
	available := Size{}
	previousDelta := Size{}
	for attempt := 0; attempt < 128; attempt++ {
		result, err := e.minimumAt(n, ref, available, font)
		if err != nil {
			return Size{}, err
		}
		result.W = math.Max(result.W, loW)
		result.H = math.Max(result.H, loH)
		if !finiteExtent(result) {
			return Size{}, diag(n, "pane-minimum", "Measured minimum exceeds finite layout bounds")
		}
		if n.Kind == "widget" || n.Kind == "markdown" || math.Abs(result.W-available.W) < 1e-7 && math.Abs(result.H-available.H) < 1e-7 {
			if e.minima == nil {
				e.minima = map[minimumKey]Size{}
			}
			e.minima[key] = result
			return result, nil
		}
		delta := Size{result.W - available.W, result.H - available.H}
		available = Size{acceleratedMinimum(result.W, delta.W, previousDelta.W, loW), acceleratedMinimum(result.H, delta.H, previousDelta.H, loH)}
		previousDelta = delta
		if err = e.step(n); err != nil {
			return Size{}, err
		}
	}
	return Size{}, diag(n, "pane-minimum", "Relative minimum measurement did not converge within its bounded budget")
}

// Relative padding/minima are piecewise affine. Extrapolate a contracting linear
// segment, then measure the candidate again before accepting it. Starting from
// zero finds the least fixed point even for a child with min-x=1.
func acceleratedMinimum(value, delta, previous, floor float64) float64 {
	if previous != 0 {
		ratio := delta / previous
		if ratio > 0 && ratio < 1 {
			candidate := value + delta*ratio/(1-ratio)
			if finite(candidate) && candidate <= 1e7 {
				return math.Max(floor, candidate)
			}
		}
	}
	return value
}

func (e *Engine) minimumAt(n *parser.Instance, ref, available Size, font float64) (Size, error) {
	if scalarField(n) {
		m, err := e.fieldMetrics(n, font, ref)
		return m.Minimum, err
	}
	if menuBar(n) {
		return e.menuMinimum(n, font)
	}
	if n.Kind == "widget" || n.Kind == "markdown" {
		if collection(n) {
			adapter, ok := e.Measure.(CollectionMeasurer)
			if !ok {
				return Size{}, diag(n, "collection-measurement", "Collection layout requires adapter metrics")
			}
			// Intrinsic probes may be smaller than native chrome. Only Minimum
			// participates here; full viewport metrics are checked at allocation.
			m, err := adapter.MeasureCollection(n, font, ref)
			if err == nil && !finiteExtent(m.Minimum) {
				err = diag(n, "collection-measurement", "Collection minimum must be finite, nonnegative and bounded")
			}
			return m.Minimum, err
		}
		return e.Measure.Measure(n, font, ref.W)
	}
	p := padding(n, ref)
	insets, err := e.containerInsets(n, font)
	if err != nil {
		return Size{}, err
	}
	inner := Size{math.Max(0, available.W-p[1]-p[3]-insets.Right), math.Max(0, available.H-p[0]-p[2]-insets.Bottom)}
	var result Size
	switch {
	case pane(n) && n.Widget == "tabs":
		m, e2 := e.tabsMetrics(n, font, inner)
		if e2 != nil {
			return Size{}, e2
		}
		child, _, e2 := e.selectedPage(n)
		if e2 != nil {
			return Size{}, e2
		}
		if child != nil {
			result, err = e.minimum(child, Size{inner.W, math.Max(0, inner.H-m.Header.H)}, font)
		}
		result.W = math.Max(result.W, m.Header.W)
		result.H += m.Header.H
	case pane(n) && n.Widget == "split":
		result, err = e.splitMinimum(n, inner, font)
	default:
		result, err = e.rowsMinimum(n, inner, ref, font)
	}
	if scrolls(n, "x") {
		result.W = 0
	}
	if scrolls(n, "y") {
		result.H = 0
	}
	result.W += p[1] + p[3] + insets.Right
	result.H += p[0] + p[2] + insets.Bottom
	return result, err
}

func (e *Engine) splitMinimum(n *parser.Instance, inner Size, font float64) (Size, error) {
	children, err := parser.PaneChildren(n)
	if err != nil {
		return Size{}, err
	}
	state, err := e.splitStateFor(n)
	if err != nil {
		return Size{}, err
	}
	divider, err := e.dividerMetric(n, font)
	if err != nil {
		return Size{}, err
	}
	horizontal := state.Axis == "horizontal"
	ref := orientedSize(math.Max(0, axisSize(inner, horizontal)-divider), crossSize(inner, horizontal), horizontal)
	var m [2]Size
	for i, c := range children {
		if i == 0 && state.Collapsed == runtime.SplitFirst || i == 1 && state.Collapsed == runtime.SplitSecond {
			continue
		}
		m[i], err = e.minimum(c.Node, ref, font)
		if err != nil {
			return Size{}, err
		}
	}
	a, b := axisSize(m[0], horizontal), axisSize(m[1], horizontal)
	u := a + b
	if state.Collapsed == runtime.SplitNone {
		// max(minFirst*U,a)+max(minSecond*U,b)<=U has this exact
		// minimum solution; both relative minima vanish for explicit collapse.
		u = math.Max(u, math.Max(a/(1-state.MinSecond), b/(1-state.MinFirst)))
	}
	return orientedSize(u+divider, math.Max(crossSize(m[0], horizontal), crossSize(m[1], horizontal)), horizontal), nil
}

func (e *Engine) rowsMinimum(n *parser.Instance, inner, ancestor Size, font float64) (Size, error) {
	gx, gy := gaps(n, ancestor)
	result := Size{}
	rows := n.Rows
	if body := n.Region("body"); body != nil {
		rows = [][]*parser.Instance{{body}}
	}
	count := 0
	vertical := []track{}
	verticalFixed, verticalLambda := 0., 0.
	for _, row := range rows {
		width, height, fixed, lambda := 0., 0., 0., 0.
		mins := make([]Size, 0, len(row))
		nodes := make([]*parser.Instance, 0, len(row))
		for _, c := range row {
			if !inFlow(c) {
				continue
			}
			m, err := e.minimum(c, inner, font)
			if err != nil {
				return Size{}, err
			}
			mins = append(mins, m)
			nodes = append(nodes, c)
			width += m.W
			height = math.Max(height, m.H)
			if w := weight(c, "x"); w > 0 {
				lo, _ := bounds(c, "x", inner.W)
				if m.W > lo {
					lambda = math.Max(lambda, m.W/w)
				}
			} else {
				fixed += m.W
			}
		}
		if len(nodes) == 0 {
			continue
		}
		gaps := gx * float64(len(nodes)-1)
		width += gaps
		// Invert the existing clamped fr rule: source minima already meeting the
		// native minimum do not force a larger common lambda on sibling tracks.
		if lambda > 0 {
			free := 0.
			for _, c := range nodes {
				if w := weight(c, "x"); w > 0 {
					lo, hi := bounds(c, "x", inner.W)
					free += clamp(lambda*w, lo, hi)
				}
			}
			width = math.Max(width, fixed+gaps+free)
		}
		for i, c := range nodes {
			if s, ok := scale(c, "x"); ok && s > 0 {
				width = math.Max(width, mins[i].W/s)
			}
			if s, ok := scale(c, "y"); ok && s > 0 {
				height = math.Max(height, mins[i].H/s)
			}
			if v, ok := c.Layout["max-x"].(float64); ok && v > 0 {
				width = math.Max(width, mins[i].W/v)
			}
			if v, ok := c.Layout["max-y"].(float64); ok && v > 0 {
				height = math.Max(height, mins[i].H/v)
			}
		}
		result.W = math.Max(result.W, width)
		result.H += height
		if len(nodes) == 1 && weight(nodes[0], "y") > 0 {
			c := nodes[0]
			lo, hi := bounds(c, "y", inner.H)
			w := weight(c, "y")
			vertical = append(vertical, track{weight: w, min: lo, max: hi})
			if height > lo {
				verticalLambda = math.Max(verticalLambda, height/w)
			}
		} else {
			verticalFixed += height
		}
		count++
	}
	if len(vertical) > 0 {
		height := verticalFixed
		for _, v := range vertical {
			height += clamp(verticalLambda*v.weight, v.min, v.max)
		}
		result.H = math.Max(result.H, height)
	}
	if count > 1 {
		result.H += gy * float64(count-1)
	}
	for _, role := range []string{"header", "footer"} {
		c := n.Region(role)
		if c == nil || !inFlow(c) {
			continue
		}
		m, err := e.minimum(c, inner, font)
		if err != nil {
			return Size{}, err
		}
		result.W = math.Max(result.W, m.W)
		result.H += m.H + gy
	}
	return result, nil
}
