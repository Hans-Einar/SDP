package layout

import (
	"math"

	"github.com/Hans-Einar/SDP/SDUI/go/parser"
)

// CollectionMetrics separates the native control minimum from the full visible
// content extent. Viewport is local to the assigned outer rectangle and excludes
// the title and any native scrollbar gutters. Content excludes the fixed title.
type CollectionMetrics struct {
	Minimum, Content Size
	Viewport         Rect
}

// CollectionMeasurer is an optional extension to Measurer. Implementations use
// their actual row/title/status metrics and prospective snapshot data. Neither
// method may load data, publish UI state or invoke application callbacks.
type CollectionMeasurer interface {
	MeasureCollection(*parser.Instance, float64, Size) (CollectionMetrics, error)
}

func (e *Engine) collectionMetrics(n *parser.Instance, font float64, outer Size) (CollectionMetrics, error) {
	measure, ok := e.Measure.(CollectionMeasurer)
	if !ok {
		return CollectionMetrics{}, diag(n, "collection-measurement", "Collection layout requires adapter metrics")
	}
	m, err := measure.MeasureCollection(n, font, outer)
	if err != nil {
		return CollectionMetrics{}, err
	}
	if !finiteExtent(m.Minimum) || !finiteExtent(m.Content) || !finite(m.Viewport.X) || !finite(m.Viewport.Y) || !finiteExtent(Size{m.Viewport.W, m.Viewport.H}) || m.Viewport.X < 0 || m.Viewport.Y < 0 || m.Viewport.X+m.Viewport.W > outer.W+.01 || m.Viewport.Y+m.Viewport.H > outer.H+.01 {
		return CollectionMetrics{}, diag(n, "collection-measurement", "Collection metrics must be finite and viewport must fit its outer rectangle")
	}
	return m, nil
}

func (e *Engine) measureWidget(n *parser.Instance, font, limit float64, size, slot Size, known assigned) (Size, Size, error) {
	if !collection(n) || e.profile != "sdui/0.3" {
		s, err := e.Measure.Measure(n, font, limit)
		if err == nil && e.profile == "sdui/0.3" && !finiteExtent(s) {
			return Size{}, Size{}, diag(n, "layout-range", "Measured size exceeds layout bounds")
		}
		return s, s, err
	}
	outer := slot
	if known.x {
		outer.W = size.W
	}
	if known.y {
		outer.H = size.H
	}
	m, err := e.collectionMetrics(n, font, outer)
	if err != nil {
		return Size{}, Size{}, err
	}
	// Content-sized axes include native chrome once. On assigned scroll axes,
	// only the viewport minimum constrains size; all rows need not fit at once.
	natural := Size{math.Max(m.Minimum.W, m.Content.W+outer.W-m.Viewport.W), math.Max(m.Minimum.H, m.Content.H+outer.H-m.Viewport.H)}
	return natural, m.Minimum, nil
}

func (e *Engine) arrangeCollection(n *parser.Instance, r, clip Rect, font float64, parent string) error {
	m, err := e.collectionMetrics(n, font, Size{r.W, r.H})
	if err != nil {
		return err
	}
	if r.W < m.Minimum.W-.01 || r.H < m.Minimum.H-.01 {
		return diag(n, "native-minimum", "Assigned size is below widget minimum")
	}
	view := m.Viewport
	view.X += r.X
	view.Y += r.Y
	_, err = e.viewport(n, view, clip.Intersect(view), m.Content, parent)
	if err == nil {
		e.gutters(n, r, view, clip)
	}
	return err
}
