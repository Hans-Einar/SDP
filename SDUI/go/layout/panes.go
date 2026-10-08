package layout

import (
	"math"

	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

// TabsMetrics describes native header chrome, not page content. Header.W is its
// minimum width; Header.H is the fixed header height in logical units.
type TabsMetrics struct{ Header Size }

// PaneMeasurer is required for tabs/split layout. Measurements are pure and use
// the same prospective snapshot as LayoutSnapshot; no native defaults are guessed.
type PaneMeasurer interface {
	MeasureTabs(*parser.Instance, float64, Size) (TabsMetrics, error)
	MeasureSplit(*parser.Instance, float64) (float64, error)
}

type TabsLayout struct {
	Header, HeaderClip, Body Rect
	Selected                 string
}
type SplitLayout struct {
	Divider, DividerClip, First, Second Rect
	Geometry                            runtime.SplitGeometry
	Collapsed                           runtime.SplitSide
}

func pane(n *parser.Instance) bool { return n.Kind == "composition" }
func argumentNumber(n *parser.Instance, key string, fallback float64) float64 {
	if v, ok := n.Arguments[key].(parser.Literal); ok {
		if x, ok := v.Value.(float64); ok {
			return x
		}
	}
	return fallback
}
func (e *Engine) tabsMetrics(n *parser.Instance, font float64, outer Size) (TabsMetrics, error) {
	m, ok := e.Measure.(PaneMeasurer)
	if !ok {
		return TabsMetrics{}, diag(n, "pane-measurement", "Tabs require native header metrics")
	}
	v, err := m.MeasureTabs(n, font, outer)
	if err != nil {
		return v, err
	}
	if !finiteExtent(v.Header) || v.Header.H <= 0 {
		return v, diag(n, "pane-measurement", "Header metrics must be finite, bounded and have positive height")
	}
	return v, nil
}
func (e *Engine) dividerMetric(n *parser.Instance, font float64) (float64, error) {
	m, ok := e.Measure.(PaneMeasurer)
	if !ok {
		return 0, diag(n, "pane-measurement", "Split requires native divider metrics")
	}
	v, err := m.MeasureSplit(n, font)
	if err != nil {
		return 0, err
	}
	if !finite(v) || v <= 0 || v > 1e7 {
		return 0, diag(n, "pane-measurement", "Divider must be finite, bounded and positive")
	}
	return v, nil
}
func (e *Engine) selectedPage(n *parser.Instance) (*parser.Instance, string, error) {
	children, err := parser.PaneChildren(n)
	if err != nil {
		return nil, "", err
	}
	selected := n.Argument("selected")
	state, has := e.tabState[n.Path]
	if has {
		selected = state.Selected
	}
	if selected == "" && !has {
		for _, c := range children {
			if visible(c.Node) && c.Node.Layout["enabled"] != false {
				selected = c.ID
				break
			}
		}
	}
	if selected == "" {
		return nil, "", nil
	}
	for _, c := range children {
		if c.ID == selected {
			if !visible(c.Node) {
				return nil, "", diag(n, "pane-state", "Selected page is inactive")
			}
			return c.Node, selected, nil
		}
	}
	return nil, "", diag(n, "pane-state", "Selected page does not belong to tabs")
}
func (e *Engine) splitStateFor(n *parser.Instance) (runtime.SplitState, error) {
	s, ok := e.splitState[n.Path]
	if !ok {
		s = runtime.SplitState{Axis: n.Argument("axis"), Proportion: argumentNumber(n, "proportion", .5), MinFirst: argumentNumber(n, "minFirst", 0), MinSecond: argumentNumber(n, "minSecond", 0), Collapsed: runtime.SplitNone, Collapsible: true}
		if v, ok := n.Arguments["collapsible"].(parser.Literal); ok {
			s.Collapsible, _ = v.Value.(bool)
		}
	}
	if s.Axis != "horizontal" && s.Axis != "vertical" || !finite(s.Proportion) || s.Proportion < 0 || s.Proportion > 1 || !finite(s.MinFirst) || !finite(s.MinSecond) || s.MinFirst < 0 || s.MinSecond < 0 || s.MinFirst+s.MinSecond >= 1 || s.Collapsed != runtime.SplitNone && s.Collapsed != runtime.SplitFirst && s.Collapsed != runtime.SplitSecond {
		return s, diag(n, "pane-state", "Invalid split state")
	}
	if !s.Collapsible && s.Collapsed != runtime.SplitNone {
		return s, diag(n, "pane-state", "Noncollapsible split cannot be collapsed")
	}
	return s, nil
}
func axisSize(s Size, horizontal bool) float64 {
	if horizontal {
		return s.W
	}
	return s.H
}
func crossSize(s Size, horizontal bool) float64 {
	if horizontal {
		return s.H
	}
	return s.W
}
func orientedSize(axis, cross float64, horizontal bool) Size {
	if horizontal {
		return Size{axis, cross}
	}
	return Size{cross, axis}
}

// splitAllocation uses content fractions after divider subtraction. Collapsed
// panes contribute neither source minima nor measured minima; the visible side
// still has to fit. No caller-owned state is modified by clamping.
func (e *Engine) splitAllocation(n *parser.Instance, inner Size, font float64) (SplitLayout, error) {
	children, err := parser.PaneChildren(n)
	if err != nil {
		return SplitLayout{}, err
	}
	s, err := e.splitStateFor(n)
	if err != nil {
		return SplitLayout{}, err
	}
	divider, err := e.dividerMetric(n, font)
	if err != nil {
		return SplitLayout{}, err
	}
	horizontal := s.Axis == "horizontal"
	u := axisSize(inner, horizontal) - divider
	if u <= 0 {
		return SplitLayout{}, diag(n, "pane-minimum", "Divider exhausts split allocation")
	}
	ref := orientedSize(u, crossSize(inner, horizontal), horizontal)
	var mins [2]Size
	for i, c := range children {
		if i == 0 && s.Collapsed == runtime.SplitFirst || i == 1 && s.Collapsed == runtime.SplitSecond {
			continue
		}
		mins[i], err = e.minimum(c.Node, ref, font)
		if err != nil {
			return SplitLayout{}, err
		}
		if crossSize(mins[i], horizontal) > crossSize(inner, horizontal)+.01 {
			return SplitLayout{}, diag(n, "pane-minimum", "Split child cross-axis minimum does not fit")
		}
	}
	first := 0.
	g := runtime.SplitGeometry{}
	switch s.Collapsed {
	case runtime.SplitFirst:
		if axisSize(mins[1], horizontal) > u+.01 {
			return SplitLayout{}, diag(n, "pane-minimum", "Visible split child minimum does not fit")
		}
	case runtime.SplitSecond:
		first = u
		if axisSize(mins[0], horizontal) > u+.01 {
			return SplitLayout{}, diag(n, "pane-minimum", "Visible split child minimum does not fit")
		}
	default:
		a := math.Max(s.MinFirst*u, axisSize(mins[0], horizontal))
		b := math.Max(s.MinSecond*u, axisSize(mins[1], horizontal))
		if a+b > u {
			return SplitLayout{}, diag(n, "pane-minimum", "Expanded split minima do not fit")
		}
		g = runtime.SplitGeometry{Lower: a / u, Upper: 1 - b/u}
		g.Effective = clamp(s.Proportion, g.Lower, g.Upper)
		first = u * g.Effective
	}
	out := SplitLayout{Geometry: g, Collapsed: s.Collapsed}
	if horizontal {
		out.First = Rect{W: first, H: inner.H}
		out.Divider = Rect{X: first, W: divider, H: inner.H}
		out.Second = Rect{X: first + divider, W: u - first, H: inner.H}
	} else {
		out.First = Rect{W: inner.W, H: first}
		out.Divider = Rect{Y: first, W: inner.W, H: divider}
		out.Second = Rect{Y: first + divider, W: inner.W, H: u - first}
	}
	return out, nil
}

func shifted(r Rect, x, y float64) Rect { r.X += x; r.Y += y; return r }
func (e *Engine) arrangePane(n *parser.Instance, r Rect, ref Size, clip Rect, font float64, enabled bool, ancestor string) (*Box, error) {
	b := &Box{Instance: n, Path: n.Path, Rect: r, Clip: clip.Intersect(r), Font: font, Enabled: enabled}
	p := padding(n, ref)
	inner := Rect{r.X + p[3], r.Y + p[0], math.Max(0, r.W-p[1]-p[3]), math.Max(0, r.H-p[0]-p[2])}
	size := Size{inner.W, inner.H}
	contentClip := clip.Intersect(inner)
	switch n.Widget {
	case "tabs":
		m, err := e.tabsMetrics(n, font, size)
		if err != nil {
			return nil, err
		}
		if size.W < m.Header.W-.01 || size.H < m.Header.H {
			return nil, diag(n, "pane-minimum", "Tabs header does not fit")
		}
		child, id, err := e.selectedPage(n)
		if err != nil {
			return nil, err
		}
		header := Rect{inner.X, inner.Y, inner.W, m.Header.H}
		body := Rect{inner.X, inner.Y + m.Header.H, inner.W, inner.H - m.Header.H}
		if e.tabs == nil {
			e.tabs = map[string]TabsLayout{}
		}
		e.tabs[n.Path] = TabsLayout{Header: header, HeaderClip: contentClip.Intersect(header), Body: body, Selected: id}
		if child != nil {
			c, err := e.arrangeAssigned(child, body, Size{body.W, body.H}, contentClip.Intersect(body), font, enabled, ancestor)
			if err != nil {
				return nil, err
			}
			b.Children = append(b.Children, c)
		}
	case "split":
		g, err := e.splitAllocation(n, size, font)
		if err != nil {
			return nil, err
		}
		children, err := parser.PaneChildren(n)
		if err != nil {
			return nil, err
		}
		g.First = shifted(g.First, inner.X, inner.Y)
		g.Second = shifted(g.Second, inner.X, inner.Y)
		g.Divider = shifted(g.Divider, inner.X, inner.Y)
		g.DividerClip = contentClip.Intersect(g.Divider)
		if e.splits == nil {
			e.splits = map[string]SplitLayout{}
		}
		e.splits[n.Path] = g
		s, _ := e.splitStateFor(n)
		childRef := size
		if s.Axis == "horizontal" {
			childRef.W -= g.Divider.W
		} else {
			childRef.H -= g.Divider.H
		}
		for i, c := range children {
			if i == 0 && g.Collapsed == runtime.SplitFirst || i == 1 && g.Collapsed == runtime.SplitSecond {
				continue
			}
			rect := g.First
			if i == 1 {
				rect = g.Second
			}
			child, err := e.arrangeAssigned(c.Node, rect, childRef, contentClip.Intersect(rect), font, enabled, ancestor)
			if err != nil {
				return nil, err
			}
			b.Children = append(b.Children, child)
		}
	default:
		return nil, diag(n, "unsupported-kind", "Unsupported pane composition")
	}
	return b, nil
}
func (e *Engine) arrangeAssigned(n *parser.Instance, r Rect, ref Size, clip Rect, font float64, enabled bool, ancestor string) (*Box, error) {
	size, err := e.desired(n, ref, Size{r.W, r.H}, assigned{true, true}, font)
	if err != nil {
		return nil, err
	}
	if math.Abs(size.W-r.W) > .01 || math.Abs(size.H-r.H) > .01 {
		return nil, diag(n, "pane-minimum", "Pane allocation conflicts with child bounds")
	}
	return e.arrange(n, r, ref, clip, font, enabled, ancestor)
}
