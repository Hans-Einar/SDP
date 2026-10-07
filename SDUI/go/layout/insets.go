package layout

import "github.com/Hans-Einar/SDP/SDUI/go/parser"

// ViewportInsets are native logical-unit gutters reserved after source padding.
// Right belongs only to scroll-y; Bottom belongs only to scroll-x. Insets depend
// on declared axes, not current ranges, so measurement has no visibility loop.
type ViewportInsets struct{ Right, Bottom float64 }

// ViewportMeasurer optionally extends Measurer for 0.3 scrolling frames/groups.
// It must be pure. CollectionMeasurer already reports collection chrome; its
// owners never invoke this adjunct. Other measurers use zero container insets.
type ViewportMeasurer interface {
	MeasureViewport(*parser.Instance, float64) (ViewportInsets, error)
}

func (e *Engine) containerInsets(n *parser.Instance, font float64) (ViewportInsets, error) {
	if e.profile != "sdui/0.3" || (n.Kind != "frame" && n.Kind != "group") || (!scrolls(n, "x") && !scrolls(n, "y")) {
		return ViewportInsets{}, nil
	}
	if v, ok := e.insets[n]; ok {
		return v, nil
	}
	m, ok := e.Measure.(ViewportMeasurer)
	if !ok {
		return ViewportInsets{}, nil
	}
	v, err := m.MeasureViewport(n, font)
	if err != nil {
		return ViewportInsets{}, err
	}
	if !finiteExtent(Size{v.Right, v.Bottom}) || v.Right > 0 && !scrolls(n, "y") || v.Bottom > 0 && !scrolls(n, "x") {
		return ViewportInsets{}, diag(n, "viewport-measurement", "Insets must be finite, nonnegative, bounded and belong to declared scroll axes")
	}
	if e.insets == nil {
		e.insets = map[*parser.Instance]ViewportInsets{}
	}
	e.insets[n] = v
	return v, nil
}

func insetSize(n *parser.Instance, size Size, insets ViewportInsets) (Size, error) {
	size.W -= insets.Right
	size.H -= insets.Bottom
	if (insets.Right > 0 || insets.Bottom > 0) && (size.W <= 0 || size.H <= 0) {
		return Size{}, diag(n, "viewport-measurement", "Native gutters exhaust the available viewport")
	}
	return size, nil
}

// gutters records only effective visible strips. Their full track lengths remain
// in Viewport.Rect, including when ancestors partially clip a strip. The corner
// is neither a horizontal nor vertical track, and own scroll does not move it.
func (e *Engine) gutters(n *parser.Instance, allocation, content, ancestorClip Rect) {
	v, ok := e.viewports[n.Path]
	if !ok {
		return
	}
	clip := ancestorClip.Intersect(allocation)
	if v.ScrollY {
		right := allocation.X + allocation.W - content.X - content.W
		if right > 0 {
			v.VerticalGutter = Rect{content.X + content.W, content.Y, right, content.H}.Intersect(clip)
		}
	}
	if v.ScrollX {
		bottom := allocation.Y + allocation.H - content.Y - content.H
		if bottom > 0 {
			v.HorizontalGutter = Rect{content.X, content.Y + content.H, content.W, bottom}.Intersect(clip)
		}
	}
	e.viewports[n.Path] = v
}
