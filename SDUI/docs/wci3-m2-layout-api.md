# WCI3-M2 extended input layout API

The selected M2 implementation reuses the M1 native field geometry API. This
contract describes layout behavior; native editing and OS IME acceptance require
separate host/integration evidence. The governing source is
[Values and text](../../SDP/04--Design/SDUI/Widgets/Values-and-text.md), including
the reviewed native editing refinement. No text value or viewport authority moves
into layout.

## Selection and snapshot requirements

`parser.InputOptions(instance)` is the sole source-policy accessor. Presence of
any new input argument (multiline, readOnly, placeholder, required), including
false/empty, opts a 0.3 input into field measurement. Basic input under 0.2 or 0.3
uses its existing generic measurement and output. Having an input entry in
Snapshot.Fields does not opt in; nor does a callback or profile alone.

Active extended input requires its matching Snapshot.Fields entry by normalized
InstancePath. Public runtime Handle.Path may differ for anonymous/reused nodes;
layout checks identity/kind without inventing another path resolver. Input metadata
must be nonnil, with Multiline/effective Placeholder/Required matching source
policy. ReadOnly is mutable snapshot state, not compared to its source default.
Runtime recomputes omitted-placeholder fallback when the input label changes.
Accepted/Proposed/RawDraft remain copied projections of Widget.Value/Draft;
layout does not validate text, parse numbers or change acceptance.

The existing source-only Layout method cannot supply an active extended field's
snapshot. Use LayoutSnapshot or LayoutCanvases through the established
PresentationGate/preparation-ticket path. Hidden/closed input declarations still
undergo source-policy validation, but need no native measurement until active.

## Existing native measurement signature

```go
type FieldMeasurer interface {
    MeasureField(*parser.Instance, runtime.FieldState, float64, Size) (FieldMetrics, error)
}
type FieldMetrics struct {
    Minimum Size
    Label, Control, Feedback, Decrement, Increment Rect
}
```

The source node, detached field projection, inherited logical font and outer
allocation are passed to the existing Engine.Measure adjunct. Input pointer
metadata and RawDraft are defensively copied. The adapter must remain pure: use
detached native measurement objects, never mutate live Entry state, publish
resources or invoke application callbacks. No default text chrome is guessed.

For extended input:

* Minimum is the finite positive native outer minimum, including label, Entry and
  reserved feedback exactly once. The native policy uses three visible rows and
  TextWrapWord for multiline input. Draft length is not an outer content extent.
* Control is the entire Entry including padding, border and internal scroll chrome.
  Shared geometry does not expose a private text viewport or its offset.
* Label is the native text-label region; an empty text label may omit it as Rect{}.
  Nonempty source text is not trimmed into absence. Label must not overlap Control.
* Feedback is the separate validation region. Active validation requires a positive
  rectangle. The native adapter reserves a fixed row so diagnostics do not resize
  the editor. Feedback must not overlap Label or Control.
* Decrement and Increment are absent zero rectangles for input.

Intrinsic/minimum probes may pass an undersized outer allocation: only Minimum is
consumed there. Final allocation checks all positive/absent parts against the
actual outer bounds and checks native minimum after source maxima. Malformed
metrics return field-measurement; insufficient space returns native-minimum;
missing/mismatched snapshot policy returns field-snapshot. Parser policy errors
retain their own source diagnostics. Layout never silently shrinks native chrome.

## Returned geometry and publication

Existing SnapshotLayout.Fields maps each active input path to FieldLayout with
Label/Control/Feedback rectangles and corresponding clips; numeric button parts
stay zero. Rectangles are translated once into their owning canvas. Clips include
window/ancestor content clips and exclude ancestor gutter strips. Enabled includes
disabled ancestors; ReadOnly remains readable/focusable and does not disable hits.
Runtime and host enforce editing eligibility.

Existing relative layout and recursive pane minima apply. Inactive page/collapsed
pane/hidden field/closed surface geometry is omitted. Open dialogs start at their
own canvas origin; an existing nonmodal canvas keeps its explicitly supplied size
when its parent resizes. Failed any-canvas measurement returns no partial result.
PresentationState still contains only outer viewport offsets and split geometry.
Fields publish with the same accepted ticket as the rest of the canvas.

No Entry text offset enters shared Viewports. Native editing internally reveals
the caret and consumes text-area wheel events even at its limits. The host must
first verify that the existing Entry subtree receives exactly one event before
adding an adapter; it must not forward that same event to outer RouteScroll.
Ancestor gutters/exposed background keep existing routing. Shared EnsureVisible
reveals the Control rectangle in its ancestors, not a private caret position.

## Lifecycle and evidence boundary

Layout does not own native object retention. Ordinary sync, page hide/show,
collapse/restore and failed reload/probe retain the Entry through the host.
Identical displayed bytes retain history, including self-echo and programmatic
Apply. Different-byte programmatic replacement resets native history once muted.
Successful compatible reload may recreate Entry/reset native edit state under the
runtime's main/page draft/focus and closed-dialog reset rules.

A pure geometry gate cannot undo a native edit already performed before OnChanged.
The reviewed host exception restores the latest authoritative draft on that same
Entry when runtime rejects such an edit; it may reset native caret/selection/scroll
and clears undo/redo. Accepted-but-invalid drafts, failed Commit/reload/probe and
unrelated publication failures do not acquire that exception. These are host
responsibilities; there is no layout edit-provenance or native-history API.

Synthetic metrics tests prove shared geometry and admission: explicit opt-in,
legacy geometry, typed projection isolation, label/feedback boundaries, finite
relative minima, nested clips/gutters, outer reveal, and hidden/canvas omission.
Actual runtime gate tests prove failed layout candidates do not publish drafts,
observers, tickets or partial canvas geometry. Native font minima, three-row word
wrapping, caret/history retention, clipboard, wheel targeting and actual SDL/OS
IME are separate obligations, not inferred from these component tests.
