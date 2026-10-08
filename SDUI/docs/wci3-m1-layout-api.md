# WCI3-M1 scalar layout API

Status: native measurement seam agreed with the host lane and implemented;
scoped verification described below. Native integration evidence remains separate. WCI3-M1
is selected under the reviewed Values-and-text contract. This note describes only
checkbox, slider, select and number; extended input/text/IME is not selected here.
Runtime owns typed values, validation and options; layout owns measured geometry.

## Additive seam

```go
type FieldMeasurer interface {
    MeasureField(*parser.Instance, runtime.FieldState, float64, Size) (FieldMetrics, error)
}
type FieldMetrics struct {
    Minimum Size
    Label, Control, Feedback, Decrement, Increment Rect
}
type FieldLayout struct {
    Label, Control, Feedback, Decrement, Increment Rect
    LabelClip, ControlClip, FeedbackClip, DecrementClip, IncrementClip Rect
    Enabled, ReadOnly bool
}
// SnapshotLayout adds Fields map[string]FieldLayout, omitted from JSON when empty.
```

Arguments are the exact normalized source node, a detached copy of that field's
prospective runtime snapshot, inherited logical font size and outer size. Native
metrics must use the same objects/font/layout policy as the prepared adapter.
They perform no application callbacks, numeric parsing, state changes or native
publication. No default scalar renderer or guessed chrome is supplied.

All metric rectangles are local to the assigned outer control. Label identifies
native label placement. Control is the actual value widget; for number it is the
editable entry, excluding the decrement/increment buttons. The checkbox's native
label may overlap its combined Check control rectangle. Number button rectangles
are distinct and do not overlap the entry or each other. Other kinds have no
number-button rectangles. Feedback is the native validation message region,
associated with this field's label, and must be present when validation has a
nonempty code or message. An absent optional region is exactly the zero rectangle.

The native host reserves a fixed feedback row, so feedback updates do not grow its
minimum. The native outer minimum includes label, value control, button chrome and any
reserved feedback space once. Intrinsic probes can pass zero or undersized outer
axes; only the returned positive finite Minimum participates in those probes.
Final allocation validates every rectangle against the actual assigned outer
size, including source max bounds. Missing adapter/state, nonfinite/negative
metrics, exhausted native space and invalid subcontrol allocation reject the
candidate. Relative layout still references finite enclosing content, not an
unbounded plane or the value's numeric range.

## Geometry and ownership

LayoutSnapshot and LayoutCanvases consume Snapshot.Fields by exact normalized
path/InstancePath, independent of public runtime handle spelling for anonymous or
reused nodes. Layout(root,size) has no field snapshot and cannot lay out these controls.
Source Arguments remain lexical source facts; layout never projects typed values
or raw numeric drafts into Arguments and never parses or rounds numeric values.

Each active scalar field contributes one ordinary outer Box and one FieldLayout
in its own SnapshotLayout.Fields. Rectangles gain ancestor scroll translation
once and clips intersect the outer Box and ancestor/window content clips. Open
dialogs retain their independent canvas origins; closed surfaces and inactive or
hidden branches contribute no field geometry. No field-specific viewport/state
authority is introduced. Existing EnsureVisible accepts the measured subcontrol
rectangle, and existing routing/gutters remain authoritative.

Enabled includes inherited disabled state. ReadOnly comes from the prospective
FieldState and does not make a focusable/readable control disabled. Native
adapters enforce the distinct mutation policy. Geometry does not authorize a
Commit, a choice target or any numeric value.

The existing PresentationState still contains only runtime viewport offsets and
split constraints. Scalar rectangles are prepared with the same final ticket as
the rest of each canvas; invalid geometry returns no candidate and cannot publish
partial typed edits. Empty .2 profiles and existing .2 layout/scroll behavior do
not acquire this adjunct.

## Scoped evidence

Synthetic native metric tests establish exact label/control/feedback/button
translation and clipping, finite relative references, nested viewport revelation,
read-only versus disabled hits, hidden/page/surface omission, split minimums and
whole-candidate rejection. Invalid metrics and injected prospective minimum failures
leave the input snapshot/source unchanged. Real runtime gate tests check failed
field edits before draft/observer/preparation/publication and invalid numeric raw
draft retention across rejected surface resize.
Exact .2 regressions show no scalar adjunct invocation. Native metrics,
paint, pointer/keyboard gestures, SDL actions and accessibility acceptance remain
with the host/integration lanes; synthetic rectangles do not prove those outcomes.
