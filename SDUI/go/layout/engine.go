package layout

import (
	"math"

	"github.com/Hans-Einar/SDP/SDUI/go/parser"
)

type assigned struct{ x, y bool }
type placement struct {
	node  *parser.Instance
	rect  Rect
	ref   Size
	force assigned
}

func (e *Engine) Layout(root *parser.Instance, viewport Size) (*Box, error) {
	// Each measurement has its own state; a prior snapshot's offsets never leak
	// into this legacy entry point.
	run := &Engine{Measure: e.Measure}
	return run.layout(root, viewport)
}

func (e *Engine) layout(root *parser.Instance, viewport Size) (*Box, error) {
	if root == nil {
		return nil, &parser.Diagnostic{Code: "layout-root", Message: "Missing root"}
	}
	if viewport.W <= 0 || viewport.H <= 0 || viewport.W > 32768 || viewport.H > 32768 || math.IsNaN(viewport.W) || math.IsNaN(viewport.H) {
		return nil, diag(root, "viewport", "Viewport must be finite and in (0,32768]")
	}
	profile, err := parser.EffectiveProfile(root)
	if err != nil {
		return nil, err
	}
	e.profile = profile
	if profile == "sdui/0.3" {
		if err := validateScrollOwners(root); err != nil {
			return nil, err
		}
	}
	if e.Measure == nil {
		e.Measure = TextMetrics{}
	}
	e.operations = 0
	size, err := e.desired(root, viewport, viewport, assigned{}, 14)
	if err != nil {
		return nil, err
	}
	rect := Rect{alignment(viewport.W, size.W, choice(root, "align-x", "start")), alignment(viewport.H, size.H, choice(root, "align-y", "start")), size.W, size.H}
	clip := Rect{0, 0, viewport.W, viewport.H}
	if err = overflow(root, rect, clip); err != nil {
		return nil, err
	}
	return e.arrange(root, rect, viewport, clip, 14, true, "")
}
func (e *Engine) desired(n *parser.Instance, ref, slot Size, force assigned, inherited float64) (Size, error) {
	if err := e.step(n); err != nil {
		return Size{}, err
	}
	if !visible(n) || auxiliary(n) && n != e.surfaceRoot {
		return Size{}, nil
	}
	font := number(n, "font", inherited)
	w, h := 0., 0.
	knownW, knownH := force.x, force.y
	if knownW {
		w = slot.W
	}
	if knownH {
		h = slot.H
	}
	if v, ok := scale(n, "x"); ok && !force.x {
		w = v * ref.W
		knownW = true
	}
	if v, ok := scale(n, "y"); ok && !force.y {
		h = v * ref.H
		knownH = true
	}
	if weight(n, "x") > 0 && !force.x {
		w = slot.W
		knownW = true
	}
	if weight(n, "y") > 0 && !force.y {
		h = slot.H
		knownH = true
	}
	if ratio, ok := n.Layout["ratio"].([]float64); ok {
		r := ratio[0] / ratio[1]
		switch {
		case explicit(n, "x") || force.x:
			h = w / r
		case explicit(n, "y") || force.y:
			w = h * r
		default:
			w = math.Min(slot.W, slot.H*r)
			h = w / r
		}
		knownW, knownH = true, true
	}
	if e.profile == "sdui/0.3" {
		if scrolls(n, "x") && !knownW || scrolls(n, "y") && !knownH {
			return Size{}, diag(n, "scroll-layout", "Scroll axes require a definite assigned size")
		}
	}
	if menuBar(n) {
		minimum, err := e.menuMinimum(n, font)
		if err != nil {
			return Size{}, err
		}
		if knownW && w < minimum.W-.01 || knownH && h < minimum.H-.01 {
			return Size{}, diag(n, "native-minimum", "Assigned menu size is below native minimum")
		}
		if !knownW {
			w = minimum.W
		}
		if !knownH {
			h = minimum.H
		}
	} else if pane(n) && n.Widget != "page" {
		minimum, err := e.minimum(n, ref, font)
		if err != nil {
			return Size{}, err
		}
		if knownW && w < minimum.W-.01 || knownH && h < minimum.H-.01 {
			return Size{}, diag(n, "pane-minimum", "Assigned pane size is below measured minimum")
		}
		if !knownW {
			w = minimum.W
		}
		if !knownH {
			h = minimum.H
		}
	} else if n.Kind == "markdown" || n.Kind == "widget" {
		limit := slot.W
		if knownW {
			limit = w
		}
		s, minimum, err := e.measureWidget(n, font, limit, Size{w, h}, slot, assigned{knownW, knownH})
		if err != nil {
			return Size{}, err
		}
		if n.Kind == "widget" && (knownW && w < minimum.W-.01 || knownH && h < minimum.H-.01) {
			return Size{}, diag(n, "native-minimum", "Assigned size is below widget minimum")
		}
		if !knownW {
			w = s.W
		}
		if !knownH {
			h = s.H
		}
	} else {
		insets, err := e.containerInsets(n, font)
		if err != nil {
			return Size{}, err
		}
		p := padding(n, ref)
		inner := Size{math.Max(0, slot.W-p[1]-p[3]), math.Max(0, slot.H-p[0]-p[2])}
		if knownW {
			inner.W = math.Max(0, w-p[1]-p[3])
		}
		if knownH {
			inner.H = math.Max(0, h-p[0]-p[2])
		}
		inner, err = insetSize(n, inner, insets)
		if err != nil {
			return Size{}, err
		}
		for _, row := range allRows(n) {
			for _, c := range row {
				if !inFlow(c) {
					continue
				}
				if !knownW && relative(c, "x") {
					return Size{}, diag(c, "layout-dependency", "Relative x requires a definite ancestor width")
				}
				if !knownH && relative(c, "y") {
					return Size{}, diag(c, "layout-dependency", "Relative y requires a definite ancestor height")
				}
			}
		}
		_, extent, err := e.contents(n, inner, ref, font, knownH)
		if err != nil {
			return Size{}, err
		}
		if !knownW {
			w = extent.W + p[1] + p[3] + insets.Right
		}
		if !knownH {
			h = extent.H + p[0] + p[2] + insets.Bottom
		}
	}
	loW, hiW := bounds(n, "x", ref.W)
	loH, hiH := bounds(n, "y", ref.H)
	if _, ratio := n.Layout["ratio"]; ratio {
		if w < loW || w > hiW || h < loH || h > hiH {
			return Size{}, diag(n, "ratio-bounds", "Aspect ratio conflicts with min/max")
		}
	} else {
		w = clamp(w, loW, hiW)
		h = clamp(h, loH, hiH)
	}
	if math.IsNaN(w) || math.IsNaN(h) || math.IsInf(w, 0) || math.IsInf(h, 0) || w > 1e7 || h > 1e7 {
		return Size{}, diag(n, "layout-range", "Computed size exceeds layout bounds")
	}
	return Size{w, h}, nil
}
func relative(n *parser.Instance, axis string) bool {
	_, s := scale(n, axis)
	if axis == "y" {
		if _, r := n.Layout["ratio"]; r && !explicit(n, "y") {
			return false
		}
	}
	return s || weight(n, axis) > 0
}
func allRows(n *parser.Instance) [][]*parser.Instance {
	r := append([][]*parser.Instance{}, n.Rows...)
	for _, region := range n.Regions {
		r = append(r, []*parser.Instance{region.Node})
	}
	return r
}

func (e *Engine) arrange(n *parser.Instance, r Rect, ref Size, clip Rect, font float64, enabled bool, ancestor string) (*Box, error) {
	font = number(n, "font", font)
	enabled = enabled && n.Layout["enabled"] != false
	b := &Box{Instance: n, Path: n.Path, Rect: r, Clip: clip.Intersect(r), Font: font, Enabled: enabled}
	if !visible(n) || auxiliary(n) && n != e.surfaceRoot {
		b.Rect.W = 0
		b.Rect.H = 0
		return b, nil
	}
	if menuBar(n) {
		minimum, err := e.menuMinimum(n, font)
		if err != nil {
			return nil, err
		}
		if r.W < minimum.W-.01 || r.H < minimum.H-.01 {
			return nil, diag(n, "native-minimum", "Assigned menu size is below native minimum")
		}
		return b, nil
	}
	if pane(n) && n.Widget != "page" {
		return e.arrangePane(n, r, ref, clip, font, enabled, ancestor)
	}
	if n.Kind == "widget" || n.Kind == "markdown" {
		preview, err := e.explicitPreview(n)
		if err != nil {
			return nil, err
		}
		if preview {
			_, minimum, err := e.previewMetrics(n, font, r.W)
			if err != nil {
				return nil, err
			}
			if r.W < minimum.W-.01 || r.H < minimum.H-.01 {
				return nil, diag(n, "native-minimum", "Assigned preview size is below measured minimum")
			}
		}
		field, err := measuredField(n)
		if err != nil {
			return nil, err
		}
		if field {
			if err := e.arrangeField(n, b); err != nil {
				return nil, err
			}
		}
		if e.profile == "sdui/0.3" && collection(n) {
			if err := e.arrangeCollection(n, r, clip, font, ancestor); err != nil {
				return nil, err
			}
		}
		return b, nil
	}
	insets, err := e.containerInsets(n, font)
	if err != nil {
		return nil, err
	}
	p := padding(n, ref)
	inner := Rect{r.X + p[3], r.Y + p[0], math.Max(0, r.W-p[1]-p[3]), math.Max(0, r.H-p[0]-p[2])}
	allocation := inner
	size, err := insetSize(n, Size{inner.W, inner.H}, insets)
	if err != nil {
		return nil, err
	}
	inner.W, inner.H = size.W, size.H
	places, extent, err := e.contents(n, Size{inner.W, inner.H}, ref, font, true)
	if err != nil {
		return nil, err
	}
	offset, err := e.viewport(n, inner, clip.Intersect(inner), extent, ancestor)
	if err != nil {
		return nil, err
	}
	if e.profile == "sdui/0.3" && (scrolls(n, "x") || scrolls(n, "y")) {
		e.gutters(n, allocation, inner, clip)
		ancestor = n.Path
	}
	for _, pos := range places {
		child := pos.rect
		child.X += inner.X
		child.Y += inner.Y
		if err := overflow(n, child, inner); err != nil {
			return nil, err
		}
		child.X -= offset.X
		child.Y -= offset.Y
		c, err := e.arrange(pos.node, child, pos.ref, clip.Intersect(inner), font, enabled, ancestor)
		if err != nil {
			return nil, err
		}
		b.Children = append(b.Children, c)
	}
	return b, nil
}
func overflow(owner *parser.Instance, r, area Rect) error {
	for _, a := range []struct {
		axis               string
		start, end, lo, hi float64
	}{{"x", r.X, r.X + r.W, area.X, area.X + area.W}, {"y", r.Y, r.Y + r.H, area.Y, area.Y + area.H}} {
		policy := choice(owner, "overflow-"+a.axis, "error")
		if policy == "scroll" && owner.Profile != "sdui/0.3" {
			return diag(owner, "unsupported-scroll", "Scrolling requires a viewport host and is not enabled in this layout profile")
		}
		if policy == "error" && (a.start < a.lo-.01 || a.end > a.hi+.01) {
			return diag(owner, "overflow-"+a.axis, "Child exceeds inner available rectangle")
		}
	}
	return nil
}
