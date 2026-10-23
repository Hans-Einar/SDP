# WCI2-M1 layout integration API

This lane implements only tabs/page/split under the reviewed
[Panes contract](../../SDP/04--Design/SDUI/Widgets/Panes-and-commands.md).
Runtime and frontend types follow the WCI2 API memos at the clone root. The API
below is the concrete layout/host seam; implementation verification is recorded
in the layout worker report. It is not native acceptance evidence.

```go
type TabsMetrics struct { Header Size }
type PaneMeasurer interface {
    MeasureTabs(*parser.Instance, float64, Size) (TabsMetrics, error)
    MeasureSplit(*parser.Instance, float64) (float64, error)
}
type TabsLayout struct {
    Header, HeaderClip, Body Rect
    Selected string
}
type SplitLayout struct {
    Divider, DividerClip, First, Second Rect
    Geometry runtime.SplitGeometry
    Collapsed runtime.SplitSide
}
```

Engine.Measure may implement PaneMeasurer. Tabs/split require it; absent native
metrics reject. The font argument is the effective inherited logical font size.
MeasureTabs receives the finite available inner size after source padding. Header.W
is the native minimum header width, Header.H the fixed native header height; both
are finite/nonnegative/bounded, with positive height. MeasureSplit returns the
actual positive finite divider thickness. Adapters do not return child content
sizes, selection or independent split offsets. Header and divider chrome are
measured from the native adapter, never guessed by shared layout.
Intrinsic probes may provide a zero or smaller-than-minimum size; measurements
still return the native minimum. Actual assigned rectangles are checked separately.

SnapshotLayout adds Tabs and Splits maps keyed by exact normalized path. Header,
Body, First, Second and Divider are full screen rectangles; HeaderClip and
DividerClip intersect chrome with the owner allocation, ancestor content clips
and window. Child Box clips exclude header/divider areas. Header/divider input and
native painting use these rectangles. Runtime supplies tabs header eligibility
separately from inactive page activity. Only the selected page has a body Box;
collapsed children have zero extent in SplitLayout and no Box or measurements.
The visible child receives the usable extent and keeps its measured minimum.

Padding precedes chrome allocation. A tabs page fills the body below its header;
a split's children fill the cross-axis and divide usable axis U after subtracting
the divider. Shared recursive minima include active descendants, native control
minima, source bounds, padding/gap, headers and nested dividers. Collection row
extent is not its native minimum. For expanded splits, bounds are
max(minFirst*U, measuredFirst)/U and 1-max(minSecond*U, measuredSecond)/U. Effective
is exactly the clamped requested proportion; impossible minima reject. Explicit
collapse suspends both relative minima and all hidden-child measurement. Saved
expanded proportion belongs to runtime and is never changed by layout.
Relative intrinsic minima are solved against each subtree's own finite body,
with a per-run cache and at most 128 refinement steps inside the existing layout
operation budget. Nonconvergent or out-of-range measurements reject. Actual child
allocation is validated after the split ratio is resolved; no invalid measurement
publishes partially arranged geometry.

SnapshotLayout.PresentationState() returns detached runtime.PresentationState:
measured active viewport offsets and SplitGeometry for every visible expanded
split, including disabled-but-visible splits. Collapsed/inactive splits are absent
from the gate result. Runtime overlays retained inactive pane offsets and clamps
on reveal; layout does not copy hidden offsets into its active geometry map.
Runtime enforces strict programmatic ratios versus user/resize/restore clamps.
Preparation consumes the same finalized snapshot, through the runtime ticket
protocol. A geometry result or consumed event sequence is not native publication.

The WCI1 viewport/gutter APIs, 0.2 source/geometry and unsupported-scroll behavior
remain unchanged. M2 command/menu/dialog support is not part of this interface.
