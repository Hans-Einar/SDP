package layout

import (
	"math"

	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

// Viewport describes measured scrolling geometry, not UI state. Rect and Clip
// are screen coordinates; Content and offsets use local logical units. Parent
// is the nearest scroll owner, or empty at the window boundary.
type Viewport struct {
	Path, Parent                     string
	Rect, Clip                       Rect
	HorizontalGutter, VerticalGutter Rect
	Content                          Size
	ScrollX, ScrollY                 bool
	Offset, Maximum                  runtime.ViewportState
}

// SnapshotLayout joins existing outer Box geometry with measured viewports.
// Collection rows remain owned by the native adapter, not this geometry tree.
type SnapshotLayout struct {
	Root      *Box
	Viewports map[string]Viewport
	Tabs      map[string]TabsLayout
	Splits    map[string]SplitLayout
}

// LayoutSnapshot measures a detached prospective runtime state. A collection
// measurer must capture that same snapshot, rather than read a live Session.
// No state or native resource is published by this method.
func (e *Engine) LayoutSnapshot(snapshot runtime.Snapshot, size Size) (*SnapshotLayout, error) {
	run := &Engine{Measure: e.Measure, requested: snapshot.Viewports, viewports: map[string]Viewport{}, tabState: snapshot.Tabs, splitState: snapshot.Splits}
	for _, offset := range snapshot.Viewports {
		if !validOffset(offset) {
			return nil, &parser.Diagnostic{Code: "viewport-offset", Message: "Viewport offsets must be finite and nonnegative"}
		}
	}
	root, err := run.layout(snapshot.Root, size)
	if err != nil {
		return nil, err
	}
	return &SnapshotLayout{Root: root, Viewports: run.viewports, Tabs: run.tabs, Splits: run.splits}, nil
}

// PresentationState returns detached active geometry for the one runtime gate.
// Runtime retains inactive pane offsets and decides strict versus clamped ratio
// operations. Collapsed and inactive splits never contribute expanded bounds.
func (g *SnapshotLayout) PresentationState() runtime.PresentationState {
	out := runtime.PresentationState{Viewports: g.EffectiveOffsets(), Splits: map[string]runtime.SplitGeometry{}}
	for path, s := range g.Splits {
		if s.Collapsed == runtime.SplitNone {
			out.Splits[path] = s.Geometry
		}
	}
	return out
}

// EffectiveOffsets returns a fresh map suitable for runtime.StateGate. Removed
// or hidden viewports are omitted; non-scroll axes have effective offset zero.
func (g *SnapshotLayout) EffectiveOffsets() map[string]runtime.ViewportState {
	out := make(map[string]runtime.ViewportState, len(g.Viewports))
	for path, v := range g.Viewports {
		out[path] = v.Offset
	}
	return out
}

func scrolls(n *parser.Instance, axis string) bool {
	return choice(n, "overflow-"+axis, "error") == "scroll"
}
func collection(n *parser.Instance) bool {
	return n.Kind == "widget" && (n.Widget == "tree" || n.Widget == "list")
}
func validOffset(p runtime.ViewportState) bool {
	return finite(p.X) && finite(p.Y) && p.X >= 0 && p.Y >= 0
}
func finite(v float64) bool { return !math.IsNaN(v) && !math.IsInf(v, 0) }
func finiteExtent(s Size) bool {
	return finite(s.W) && finite(s.H) && s.W >= 0 && s.H >= 0 && s.W <= 1e7 && s.H <= 1e7
}

func validateScrollOwners(n *parser.Instance) error {
	switch n.Kind {
	case "frame", "group", "markdown":
	case "composition":
		if n.Widget != "tabs" && n.Widget != "page" && n.Widget != "split" {
			return diag(n, "unsupported-kind", "Unsupported pane composition")
		}
	case "widget":
		switch n.Widget {
		case "button", "input", "svg", "tree", "list":
		default:
			return diag(n, "unsupported-widget", "Unsupported widget kind "+n.Widget)
		}
	default:
		return diag(n, "unsupported-kind", "Unsupported instance kind "+n.Kind)
	}
	if scrolls(n, "x") || scrolls(n, "y") {
		if n.Kind != "frame" && n.Kind != "group" && !collection(n) {
			return diag(n, "unsupported-scroll", "Scroll owner must be a frame, group, tree or list")
		}
		// A scroll viewport cannot acquire its dimension from the content it
		// scrolls. Check declarations even on hidden branches, which are not
		// otherwise measured until they become visible.
		_, ratio := n.Layout["ratio"].([]float64)
		for _, axis := range []string{"x", "y"} {
			if scrolls(n, axis) && !explicit(n, axis) && !ratio {
				return diag(n, "scroll-layout", "Scroll axes require fill, fr, scale or a resolved frame ratio")
			}
		}
		if choice(n, "justify", "start") != "start" {
			return diag(n, "scroll-layout", "Scroll owners require start justification")
		}
	}
	for _, row := range allRows(n) {
		for _, c := range row {
			if err := validateScrollOwners(c); err != nil {
				return err
			}
		}
	}
	return nil
}

func contentAlignment(owner *parser.Instance, axis string, available, used float64, kind string) float64 {
	if owner.Profile == "sdui/0.3" && scrolls(owner, axis) {
		available = math.Max(available, used)
	}
	return alignment(available, used, kind)
}

func (e *Engine) viewport(n *parser.Instance, r, clip Rect, extent Size, parent string) (runtime.ViewportState, error) {
	if e.profile != "sdui/0.3" {
		return runtime.ViewportState{}, nil
	}
	if !finiteExtent(extent) {
		return runtime.ViewportState{}, diag(n, "layout-range", "Content extent exceeds layout bounds")
	}
	if err := overflow(n, Rect{r.X, r.Y, extent.W, extent.H}, r); err != nil {
		return runtime.ViewportState{}, err
	}
	sx, sy := scrolls(n, "x"), scrolls(n, "y")
	if !sx && !sy {
		return runtime.ViewportState{}, nil
	}
	if sx && r.W <= 0 || sy && r.H <= 0 {
		return runtime.ViewportState{}, diag(n, "scroll-layout", "Scroll viewport must have positive size")
	}
	requested := e.requested[n.Path]
	if !validOffset(requested) {
		return runtime.ViewportState{}, diag(n, "viewport-offset", "Viewport offsets must be finite and nonnegative")
	}
	v := Viewport{Path: n.Path, Parent: parent, Rect: r, Clip: clip, Content: extent, ScrollX: sx, ScrollY: sy}
	if sx {
		v.Maximum.X = math.Max(0, extent.W-r.W)
		v.Offset.X = clamp(requested.X, 0, v.Maximum.X)
	}
	if sy {
		v.Maximum.Y = math.Max(0, extent.H-r.H)
		v.Offset.Y = clamp(requested.Y, 0, v.Maximum.Y)
	}
	if e.viewports == nil {
		e.viewports = map[string]Viewport{}
	}
	e.viewports[n.Path] = v
	return v.Offset, nil
}
