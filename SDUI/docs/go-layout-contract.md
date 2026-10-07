# Go measurement contract

The shared `go/layout` package owns measured outer geometry. Its original G2
behavior remains the SDUI 0.2 path; WCI1 adds the bounded SDUI 0.3 snapshot and
viewport APIs described here. The [0.3 source profile](profile-0.3.md) and
[WCI1 collection contract](../../SDP/04--Design/SDUI/Widgets/Collections.md) distinguish
source acceptance, runtime state, native integration and their acceptance evidence.

## Common geometry and units

SVG and native hosts receive `Box` values containing the instance path, rectangle,
clip, font, enabled state and children. `Box.Rect` and `Box.Clip` use logical screen
coordinates. Clips constrain both painting and half-open hit tests. `Box.Hit`
returns the topmost enabled widget under the point, excluding clipped content.
Measurement executes no application callbacks and source dimensions contain no
pixel values.

`font` uses absolute logical display units (DIP). `font=10` means ten units in both
SVG viewBox and host. Physical unit size follows host DPI, not window size. The
reference measurer embeds Go Regular through OpenType at 72 DPI; Markdown headings
may use relative font sizes. SVG must use the same font or disclose substitution.
Native collection adapters supply their own actual rendering metrics as described
below; reference text measurement is not evidence of native metric equivalence.

Rows do not change the ancestor for relative dimensions. Fr tracks distribute
remaining space using `clamp(lambda * weight,min,max)`; minima are not added on top
of weights. Scaled/content-sized siblings do not shrink to conceal overflow.
Header/footer allocation precedes the body; regions use available width unless
another x policy is selected. `items=stretch` lets a single row component fill
width and equalizes component heights in multicolumn rows.

A content-sized ancestor with relatively sized children on the same axis rejects
with `layout-dependency`. A width-driven ratio derives height without requiring a
height-defined ancestor. Min/max never distort ratios. Padding/gap refer to the
node's explicit ancestor. Scrolling preserves these finite reference areas; it
never substitutes an infinite content plane or synthetic row reference.

The root viewport must be finite and positive, at most 32768 logical units per
axis. Computed sizes/content extents remain bounded by `1e7`; measurement retains
the 200000-operation budget. Parser source/depth/expansion limits also apply.

## Profiles and overflow

Both layout entry points call `parser.EffectiveProfile(root)`. Empty normalized
`Instance.Profile` means legacy 0.2; every 0.3 instance carries `sdui/0.3`. Unknown,
mixed or malformed instance trees reject. Normalized instances do not inherit a
profile merely because a caller selected a newer host.

The 0.2 geometry and `Engine.Layout(root, size)` API remain available. Existing
scroll declarations retain the `unsupported-scroll` behavior and diagnostic;
WCI1 does not reinterpret 0.2 scroll as supported clipping or scrolling. Ordinary
supported 0.2 static SVG behavior remains covered by regression tests.

For 0.3, supported scroll owners are frames, groups, trees and lists. Direct scroll
on button/input/svg/Markdown rejects. Every scroll axis needs an explicit fill/fr,
scale or resolved frame-ratio dimension; a content-dependent viewport rejects with
`scroll-layout`. A scroll owner requires `justify=start`. Owner/axis/justification
checks include hidden branches. Unknown 0.3 instance/widget kinds also reject.

| Overflow policy | Geometry behavior |
| --- | --- |
| `error` | Reject content outside the available inner rectangle. |
| `clip` | Preserve content geometry, apply clipping and keep offset zero. |
| `scroll` in 0.3 | Measure finite content extent, clip to the viewport and apply a clamped offset. |

Allocate fixed/fill/fr tracks against the finite parent area, then measure direct
child extents. An inner viewport contributes its outer rectangle to its parent's
extent, not all of its hidden content. On a scroll axis, overflowing center/end
alignment starts nonnegative. A scrolling container moves all content, including
header/footer. Regions retain viewport-relative allocation and are not sticky;
use a separate scrolling body child when page header/footer should remain fixed.

## Snapshot entry point and effective offsets

The original `Engine.Layout(root, Size) (*Box, error)` measures without requested
runtime offsets. For prospective runtime state use:

```go
func (*Engine) LayoutSnapshot(runtime.Snapshot, Size) (*SnapshotLayout, error)

type SnapshotLayout struct {
    Root      *Box
    Viewports map[string]Viewport
}

type Viewport struct {
    Path, Parent                     string
    Rect, Clip                       Rect
    HorizontalGutter, VerticalGutter Rect
    Content                          Size
    ScrollX, ScrollY                 bool
    Offset, Maximum                  runtime.ViewportState
}

func (*SnapshotLayout) EffectiveOffsets() map[string]runtime.ViewportState
```

Maps are keyed by normalized instance path. `Parent` is the nearest scroll-owner
ancestor, or empty at the window boundary. `Rect` is the screen-space viewport;
`Clip` intersects it with every ancestor clip and the window. `Content` is a local
size. `Offset` and `Maximum` use runtime's existing `{X,Y}` logical-unit type;
layout introduces no competing runtime state model or collection row scene.

For each scroll axis, `Maximum = max(0, Content - viewport size)` and
`Offset = clamp(requested, 0, Maximum)`. Negative/NaN/infinite requests reject;
requests beyond a new maximum clamp. A non-scroll axis has effective offset zero.
The result includes visible scroll owners and omits removed/hidden ones.
`EffectiveOffsets` returns a fresh complete map suitable for `runtime.StateGate`.

The calculation is pure with respect to Session: it consumes detached prospective
state and returns geometry plus matching offsets. A collection measurer must
capture that same snapshot. A host can stage these results in `CheckStateWith`,
then publish them only when the runtime operation succeeds. Failure leaves the
live state and last valid presentation intact; candidate construction must not
publish, start a provider load or invoke a domain callback.

Runtime owns accepted offsets. The host submits changes through its checked
viewport operations or an identity-guarded complete-map update. A routed operation
that changes several ancestors must publish the map atomically, not update each
ancestor separately. The host also owns bundle/lifecycle guards and native
publication; the geometry result alone does not establish those guarantees.

## Collection measurement and coordinate ownership

Collections require a `Measurer` that also implements:

```go
type CollectionMeasurer interface {
    MeasureCollection(*parser.Instance, float64, Size) (CollectionMetrics, error)
}

type CollectionMetrics struct {
    Minimum, Content Size
    Viewport         Rect
}
```

Arguments are instance, logical font size and assigned outer size. `Minimum` is
the actual native outer control minimum; it is separate from full loaded-row
extent. `Content` covers visible rows, indentation, rendered markers, separators
and status/recovery affordances, excluding the fixed title. `Viewport` is local
to the outer control and excludes the title strip and native scrollbar gutters.
Account for each gutter once and use the same metrics for measurement, painting
and hit testing. Missing, nonfinite, negative or out-of-bounds adapter metrics
reject with `collection-measurement`; assigned sizes below the native minimum
reject with `native-minimum`.

On a scroll axis, all rows need not fit the assigned native minimum. On a
content-sized non-scroll axis, layout adds native chrome to content size once;
explicit `error`/`clip` policies still apply. No default collection row renderer,
provider loading or guessed native metrics are supplied by layout.

| Boundary | Coordinate responsibility |
| --- | --- |
| Outer collection `Box` | Shared layout applies every ancestor scroll translation once. The collection's own offset does not move its outer control or title. |
| Collection `Viewport` | Shared layout returns screen-space bounds, effective clip and the runtime offset. |
| Native rows | The adapter adds its viewport origin to content-local row bounds and subtracts the collection's own offset once. |
| Pointer and focus targets | The adapter uses the same bounds/offsets, intersects with the effective clip, and reports actual measured target rectangles in screen coordinates. |

Native controls mirror accepted runtime offsets under a synchronization guard;
toolkit scroll state is not a second authority. Collection item identity,
selection, expansion and load status belong to runtime. The adapter may own row
measurement/rendering, but must keep stable item identities separate from row
indexes and status affordances. Public collection SVG export remains explicitly
unsupported; a measured collection does not imply a static snapshot renderer.

## Routing and focus revelation

```go
func (*SnapshotLayout) RouteScroll(x, y, dx, dy float64) (
    map[string]runtime.ViewportState, runtime.ViewportState, error)
func (*SnapshotLayout) EnsureVisible(path string, target Rect) (
    map[string]runtime.ViewportState, error)
```

`RouteScroll` takes a screen-space pointer and logical deltas; positive deltas
move toward content end. It finds the topmost visible branch, then consumes each
axis at eligible viewports from inner to outer. Eligibility includes each owner's
effective content clip and gutter strips. Only the unconsumed remainder
reaches a parent; siblings and obscured branches do not scroll. The second return
value is the unused delta. Zero ranges consume nothing. Nonfinite input rejects.

`EnsureVisible` takes an ordinary widget or collection path plus its actual
screen-space target rectangle. It adjusts minimally from inner to outer and
translates the target before considering each parent. A target larger than the
viewport aligns its leading edge. Missing/hidden/disabled targets reject.
Collection adapters supply measured row bounds; substituting the entire viewport
width can cause unnecessary horizontal movement.

Both methods return fresh complete offset maps without changing geometry,
Session, selection or activation. Publish through the runtime boundary, then use
the matching newly measured geometry. Keyboard routing must select the focused
viewport ancestry or a point inside its effective clip; an unclipped viewport
origin can be outside the visible window after ancestor scrolling.

## Evidence boundary

Existing layout tests retain fr/scale, minima/maxima, fixed fonts, ratios, regions,
clipped hits and the Concept1 model geometry. WCI1 tests add exact nested
translation/clip/extent oracles, finite track references, error/clip/scroll,
profile rejection, remainder routing, focus revelation and collection minimum/
chrome/extent separation. A real runtime state-gate integration test checks failed
candidate preservation and matching offset clamping after content shrink. Gutter
oracles cover padding, independent nested strips, corner exclusion, finite fr/scale
references, content sizing, resize, zero ranges, collection chrome, rejected native
metrics and exhausted-space state-gate preservation. The 0.2 regression checks
ensure the native inset adjunct is never invoked on that profile.

These tests use synthetic collection adapter metrics. They do not prove native
Fyne painting, row metrics, thumb dragging, external keyboard/pointer input or
atomic native resource publication. Those WCI1 acceptance obligations remain
pending integrated native verification. This contract adds no later widget-family
or future-stage support claim.

## Native frame/group gutters

The [native pilot refinement](../../SDP/04--Design/SDUI/Widgets/Collections.md#native-pilot-refinement--independently-review-before-implementation)
selects a bounded correction for overlapping parent/child thumbs. Following
independent design approval, shared layout implements these APIs and strip rules.
Native adapter integration and actual separate thumb reachability require their
own evidence; layout tests do not establish native acceptance.

```go
type ViewportInsets struct { Right, Bottom float64 }
type ViewportMeasurer interface {
    MeasureViewport(*parser.Instance, float64) (ViewportInsets, error)
}
```

An optional `ViewportMeasurer` adjunct on `Engine.Measure` returns pure native
metrics for 0.3 scrolling frames/groups; the float argument is the effective font
size. `Right` is permitted only for declared `overflow-y=scroll`, and `Bottom`
only for declared `overflow-x=scroll`. Values must be finite, nonnegative and at
most `1e7`. Invalid metrics or a gutter consuming all remaining width or height
reject with `viewport-measurement`; measurer errors propagate. Native gutters are
reserved for declared axes, including zero ranges. The existing CollectionMeasurer
continues to own collection title/gutter measurements. Default measurers return
zero implicit insets; 0.2 measurement does not consult this extension.

Let `A` be the screen-space assigned rectangle after source padding, before native
gutters, and `C = (A.X, A.Y, A.W-Right, A.H-Bottom)` the content viewport. Subtract
insets before contents in both desired measurement and arrangement; add them back
once when deriving content-sized non-scroll outer dimensions. Keep each run's
pure measurements consistent across both operations: the adjunct is called at
most once per visible scrolling frame/group in a layout run and cached only for
that run. Children use the remaining
finite content/body area for relative dimensions and tracks.

The `Viewport` gutter fields have these exact meanings:

| Field | Unclipped strip | Returned geometry |
| --- | --- | --- |
| `VerticalGutter` | `(C.X+C.W, C.Y, Right, C.H)` | Intersection with the owner allocation, ancestors' content clips and window. |
| `HorizontalGutter` | `(C.X, C.Y+C.H, C.W, Bottom)` | Intersection with the owner allocation, ancestors' content clips and window. |

Neither strip is intersected with `C`, which excludes them. The bottom-right
`Right × Bottom` corner belongs to neither strip. Bounds are half-open, so a
shared edge does not create overlapping hit regions. An absent inset/axis or
fully clipped strip produces an empty hit region. The owner's own content offset
never translates its gutters; ancestor offsets translate them once.

`Viewport.Rect`/`Clip` continue to describe content, excluding gutters. Children
paint/hit only within that content clip. Thumb placement and travel use the full
unclipped content viewport length; returned gutter strips only constrain visible
painting/hits. A partially clipped strip must not shorten the scroll range or
rescale drag movement. Zero range can retain reserved gutter space without a
visible/draggable thumb. Wheel routing treats a hit in an owner's effective strip
as a hit on that owner, then chains only its unconsumed remainder outward.

For collections, corresponding right/bottom strips are derived from their
existing local CollectionMetrics.Viewport and assigned outer rectangle, excluding
the title. This does not add another native inset call or reserve collection
chrome twice. The collection measurer defines any chrome already inside that
allocation; only declared scroll axes produce corresponding gutter strips.
Runtime offset/state types, source syntax and row rendering remain unchanged.
