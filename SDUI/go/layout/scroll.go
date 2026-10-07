package layout

import (
	"math"

	"github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

// RouteScroll consumes each delta at the innermost viewport under the pointer,
// then sends only the remainder to its ancestors. Positive deltas move toward
// the content end. It returns a fresh complete offset map and unused deltas.
// Publish the map with Session.SetViewports; this method never mutates geometry.
func (g *SnapshotLayout) RouteScroll(x, y, dx, dy float64) (map[string]runtime.ViewportState, runtime.ViewportState, error) {
	out := g.EffectiveOffsets()
	remainder := runtime.ViewportState{X: dx, Y: dy}
	if !finite(x) || !finite(y) || !finite(dx) || !finite(dy) {
		return nil, remainder, &runtime.Fault{Code: "viewport-delta", Message: "Scroll coordinates and deltas must be finite"}
	}
	branch := hitBranch(g.Root, x, y)
	if len(branch) == 0 || !branch[len(branch)-1].Enabled {
		return out, remainder, nil
	}
	for i := len(branch) - 1; i >= 0; i-- {
		v, ok := g.Viewports[branch[i].Path]
		if !ok || !(v.Clip.Contains(x, y) || v.HorizontalGutter.Contains(x, y) || v.VerticalGutter.Contains(x, y)) {
			continue
		}
		next := v.Offset
		if v.ScrollX {
			next.X = clamp(v.Offset.X+remainder.X, 0, v.Maximum.X)
			remainder.X -= next.X - v.Offset.X
		}
		if v.ScrollY {
			next.Y = clamp(v.Offset.Y+remainder.Y, 0, v.Maximum.Y)
			remainder.Y -= next.Y - v.Offset.Y
		}
		out[v.Path] = next
	}
	return out, remainder, nil
}

// EnsureVisible minimally reveals a screen-space target, walking inner-to-outer.
// Use a collection's path with its adapter-measured row rectangle, or an ordinary
// widget's path/Box.Rect. For an oversized target its leading edge is aligned.
// Returned offsets are prospective only; remeasure after publishing them.
func (g *SnapshotLayout) EnsureVisible(path string, target Rect) (map[string]runtime.ViewportState, error) {
	if !finite(target.X) || !finite(target.Y) || !finiteExtent(Size{target.W, target.H}) {
		return nil, &runtime.Fault{Code: "viewport-target", Message: "Reveal target must be a finite nonnegative rectangle"}
	}
	branch := pathBranch(g.Root, path)
	if len(branch) == 0 || !branch[len(branch)-1].Enabled || !visible(branch[len(branch)-1].Instance) {
		return nil, &runtime.Fault{Code: "viewport-target", Message: "Reveal target is missing, hidden or disabled"}
	}
	out := g.EffectiveOffsets()
	for i := len(branch) - 1; i >= 0; i-- {
		v, ok := g.Viewports[branch[i].Path]
		if !ok {
			continue
		}
		next := v.Offset
		if v.ScrollX {
			next.X = clamp(v.Offset.X+revealDelta(target.X, target.W, v.Rect.X, v.Rect.W), 0, v.Maximum.X)
		}
		if v.ScrollY {
			next.Y = clamp(v.Offset.Y+revealDelta(target.Y, target.H, v.Rect.Y, v.Rect.H), 0, v.Maximum.Y)
		}
		target.X -= next.X - v.Offset.X
		target.Y -= next.Y - v.Offset.Y
		out[v.Path] = next
	}
	return out, nil
}

func revealDelta(start, size, viewStart, viewSize float64) float64 {
	if size > viewSize || start < viewStart {
		return start - viewStart
	}
	return math.Max(0, start+size-viewStart-viewSize)
}

// hitBranch includes container backgrounds, unlike Box.Hit which returns only
// controls. Reverse paint order prevents routing into a sibling behind the hit.
func hitBranch(b *Box, x, y float64) []*Box {
	if b == nil || !b.Clip.Contains(x, y) || !b.Rect.Contains(x, y) {
		return nil
	}
	for i := len(b.Children) - 1; i >= 0; i-- {
		if child := hitBranch(b.Children[i], x, y); len(child) > 0 {
			return append([]*Box{b}, child...)
		}
	}
	return []*Box{b}
}
func pathBranch(b *Box, path string) []*Box {
	if b == nil || !visible(b.Instance) {
		return nil
	}
	if b.Path == path {
		return []*Box{b}
	}
	for _, c := range b.Children {
		if child := pathBranch(c, path); len(child) > 0 {
			return append([]*Box{b}, child...)
		}
	}
	return nil
}
