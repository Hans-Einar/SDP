# WCI2-M2 layout / native canvas seam

Scope: selected WCI2-M2 after reviewed Panes-and-commands.md; M1 and 0.2 remain.
Layout and host coordinated these additive interfaces before dependent code. This
memo specifies the interface; implementation evidence belongs to the worker report.

```go
type MenuMeasurer interface {
    MeasureMenu(*parser.Instance, float64) (Size, error)
}
type CanvasLayout struct {
    Main *SnapshotLayout
    Surfaces map[string]*SnapshotLayout
}
func (*Engine) SurfaceSize(runtime.Snapshot, string, Size) (Size, error)
func (*Engine) LayoutCanvases(runtime.Snapshot, Size, map[string]Size) (*CanvasLayout, error)
func (*CanvasLayout) PresentationState() runtime.PresentationState
```

MeasureMenu returns the native bar control minimum at the effective inherited
logical font size. It applies only to a menu with mode=bar; popup rows, nested
menus and groups remain native-adapter geometry. Missing, nonfinite, nonpositive or
out-of-range measurements reject. Source fill/scale may enlarge a bar control;
assigned dimensions below its native minimum reject. Bar menus are leaf boxes in
source position, never rows of ordinary layout controls.

Commands, menu groups/items/separators, context/submenu menus and dialog declarations
never participate in their enclosing flow, gaps, relative dependency checks or
minima. Open dialogs are still excluded there: their ordinary content is measured
in separate local canvases. Closed dialog content is not measured. Capability,
binding and resource preflight still inspects all declarations in other lanes.

SurfaceSize takes the exact normalized dialog path and finite reference area. It
evaluates that dialog's source-relative initial size without changing state. Host
supplies parent content area for a modal opening and parent host area for a new
nonmodal opening; it does not reimplement source scale/fill/min/max arithmetic.
The helper operates on a prospective open surface snapshot and returns finite
nonnegative logical content-canvas dimensions bounded by 32768. An empty natural
dimension may be zero; host combines it with the concrete content adapter's measured
minimum before supplying a positive accepted canvas. Native title/window decoration
chrome is outside this content canvas and remains host-owned.

LayoutCanvases receives main canvas size and accepted sizes keyed by exact dialog
path. Every open surface requires an explicit size. Unknown entries reject;
retained entries for closed declarations are ignored. Each open dialog root is
assigned (0,0,width,height), so root opening scale/fill/ratio/min/max policies are
not applied again on native resize. SurfaceSize resolves those policies once.
Root padding references this accepted local canvas, descendants their
ordinary finite allocated bodies. Ancestor fonts are inherited from source, while
coordinates, clipping and viewport ancestry restart at each canvas boundary.
Nested dialogs do not contribute to their parent's extent or child hit routing.

The result preserves separate Main and Surfaces geometry and merges their active
viewport offsets and expanded split bounds by globally unique normalized path.
PresentationState returns fresh runtime-owned types, never native pointers or a
second state authority. Runtime retains inactive/closed offsets and validates the
single aggregate gate result; host final preparation uses one finalized snapshot
and ticket for all canvases. A failed canvas returns no aggregate result. Each
accepted nonmodal size remains independent of parent resize. Native modal input
blocking/focus/lifetime are runtime/host responsibilities and do not hide parent
painting geometry.

LayoutSnapshot remains the main-canvas entry point and excludes auxiliary surfaces;
M2 hosts with open surfaces use LayoutCanvases for the complete gate. Surface hit
routing/EnsureVisible use that surface's SnapshotLayout and local coordinates;
there is no cross-window scroll ancestry. M2 tests do not establish SDL execution,
OS/window focus, menu selection ordering or native acceptance.
