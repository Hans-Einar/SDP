# WCI3-M2 layout seam preparation

Canonical selection note: M2 is now selected after WCI3-M1 delivery. The reviewed
Values-and-text native editing refinement supersedes earlier read-only/pending
status and broad history promises in this preparation record: identical bytes
preserve history; changed programmatic text resets it; an actual rejected native
edit may reset that Entry history/caret/selection/scroll while restoring current
authoritative draft muted. Failed Commit/reload/probe and admitted invalid drafts
retain history. Exact captured CR/LF clipboard content rejects before native paste.

Historical preparation status: read-only preparation plus this memo only, 2026-10-08. M1 layout remains
frozen; M2 implementation and dependencies are not selected. Main owns selection,
Session records, integration and actual native evidence. SDP 1.1.1, Worker 2.0.0
and the shared document workflow were reused; Session0010 S4 recovered read-only.

Inspected clone `/tmp/sdp-sdui-widgets`, branch `sdui/widgets-wci3`, HEAD
`a28b3ccc72b1ecc44300d4aa7c3cd679f4bc9625`, with concurrent uncommitted M1 work.
This is an inspection of those working bytes, not evidence for unchanged HEAD.
Authority: ORIGINAL `SDP/04--Design/SDUI/Widgets/Values-and-text.md`, §§2, 3, 6,
7 and final M2 reconciliation, SHA-256
`cb52f4174f705d6d60c5ee2f68e76cd5b6efdb24d463945cc18b7c034a7a681e`
(refreshed after the approved exact CR/LF single-line refinement).
Read alongside `WCI3-runtime-API.md`, `WCI3-frontend-API.md`,
`WCI3-M2-frontend-API.md`, the subsequent `WCI3-M2-runtime-API.md`
and actual implementations below.

## 1. Compatibility and actual foundation

Canonical review now approves explicit presence of any new input argument
(multiline/readOnly/placeholder/required) as the extended-policy switch, including
false or an empty placeholder. No new arguments means legacy behavior in both
0.2 and 0.3. All new arguments reject in 0.2. The frontend memo's earlier pending
opt-in/self-echo wording is superseded by the final canonical reconciliation;
Dalton confirmed this during preparation. No new approval is requested here.

Current M1 `layout.scalarField` selects only checkbox/slider/select/number.
Input continues through generic `Measurer.Measure`, with legacy input sizing;
it does not acquire FieldMeasurer semantics because Snapshot.Fields contains it.
Runtime Field(input) currently projects Widget.Value/Draft into copied String
Accepted/Proposed and RawDraft. That projection is not extended editing support.
SnapshotRoot projects the current input draft into its value argument; current
host sync reads that argument, not the accepted value alone.

The legacy host creates NewInput, uses text as placeholder and wires ordinary
Draft/Commit callbacks. Bundle.apply retains controls and calls SetText only if
text differs; inactive controls are hidden, not deleted. The enclosing routedClip
currently forwards wheel events to outer scroll routing. These are the paths an
extended adapter must deliberately distinguish, not globally replace.

## 2. Smallest proposed shared seam

Reuse the existing exported types without adding a text viewport or second value:

```go
type FieldMeasurer interface {
    MeasureField(*parser.Instance, runtime.FieldState, float64, Size) (FieldMetrics, error)
}
type FieldMetrics struct {
    Minimum Size
    Label, Control, Feedback, Decrement, Increment Rect
}
// Existing FieldLayout projects the five rectangles and their clips into
// the owning canvas, with Enabled and ReadOnly. SnapshotLayout.Fields keys
// remain normalized InstancePath; runtime handles are not path substitutes.
```

After M2 selection, extend the internal field-geometry predicate to scalar kinds
OR opted-in input, using Dalton's proposed `parser.InputOptions(n).Extended`.
Propagate helper diagnostics; do not duplicate argument parsing, inject defaults,
or infer opt-in from profile, callback or Fields presence. InputOptions is a
proposed frontend helper, not an export already available in inspected M1 code.
James's subsequent runtime proposal adds `FieldState.Input *InputState`, with
`InputState{Multiline bool; Placeholder string}` and nil for legacy/non-input.
The frontend helper remains the source-policy authority; candidate snapshot
metadata must agree with that policy. Extended geometry requires the matching
nonnil Input projection rather than manufacturing it. Native metrics consume its
copied Multiline/effective Placeholder, and existing ReadOnly/Required/Validation.
Copy the new Input pointer defensively in layout's adapter input as runtime does
in snapshot copies. Missing/mismatched candidate metadata rejects; an incidental
legacy Field projection must never select extended geometry. This additive runtime
proposal leaves the FieldMeasurer signature and geometry types unchanged.
Layout neither validates text nor manufactures field state.

For extended input, Control denotes the entire native Entry allocation, including
its border, padding and internal scroll chrome. It does not denote a guessed text
baseline or the private inner editing viewport. Decrement/Increment are absent
zero rectangles. No EntryViewport/EntryOffset field, new PresentationState map,
text extent in runtime.Viewports, GUI import or custom text renderer is needed.
Native inspection can expose additional measured observations later if evidence
requires them; those must not become authoritative shared edit state.

Label is measured separately from placeholder and value. Nonempty text requires
its measured label region. Empty text remains legal: permit Label=Rect{} only for
an extended input with empty text, preserving scalar label requirements. Do not
trim or reject source text to achieve this. This is a bounded proposed validation
refinement; there is no reason to add a new metrics type. Dalton agrees; Noether subsequently proposed this same empty-label rule
for the extended adapter, pending M2 selection.

Feedback uses the existing measured region and snapshot validation. Reuse the
native fixed feedback row policy so an error need not resize the editor or jump
its caret. Active validation must have positive, nonoverlapping feedback geometry;
when no feedback is allocated, absence is exactly Rect{}. ReadOnly is projected
independently of Enabled: it must not make the field disappear or disable focus.
Required/invalid state changes feedback, not text or authoritative acceptance.

Intrinsic probes consume only a finite positive Minimum. Final allocation must
fit that minimum and validate every local part, after source max constraints.
Apply the existing bounds, cost accounting, defensive snapshot copies and failure
diagnostics. Translate parts once and intersect with the owning Box.Clip. Existing
relative sizing, split minimum/reveal, scroll extents and canvas-local coordinates
remain the shared geometry; no text-content-sized outer scroll scene is created.

## 3. Measurable bounded multiline control

The inspected Fyne v2.8.1 Entry exposes MultiLine, Wrapping and
SetMinRowsVisible. `TextWrapWord` selects native vertical scrolling in its Entry
renderer. Its scrolling MinSize uses native font/padding and a bounded visible
row count (default three), rather than the entire draft's line count. SetText
explicitly resets undo history. These are source observations, not native proof.

Proposed host adapter: measure a detached Entry configured with the same effective
font/theme, multiline/wrap/placeholder and chrome as the live control. Use native
MinSize for the Entry region, plus measured label/feedback once. Noether now proposes exactly `SetMinRowsVisible(3)` and `TextWrapWord`,
with the same fixed feedback row and exact-font measurement. This freezes a
concrete proposed adapter policy for selection, not a new source argument or pixel
value; no product implementation is selected by this memo.

Assigned size stays finite. A long draft, many newlines or long unbroken token
must not turn outer intrinsic height into the document height. At final layout,
Entry receives exactly Control's allocated size; word-wrap and native vertical
scroll handle additional content there. Narrow/exhausted allocations reject before
publication instead of silently overlapping label/feedback. Resize/reflow clamps
native scrolling after accepted publication without text replacement or Commit.
Single-line extended input retains native horizontal caret reveal. No new source
wrap/axis flags, mandatory full-scene API or source overflow ownership for input.

## 4. Scroll, event and publication ownership

| Region/action | Owner and required behavior |
| --- | --- |
| Inside Entry text area, including its scroll limit | Entry consumes wheel; no forwarded remainder and no second outer route. Native scroll chrome belongs to Entry. |
| Outer frame/group gutter or exposed background | Existing WCI1 routing and checked runtime offsets; native text scrolling must not obscure the ancestor gutter. |
| Label/feedback outside Entry | Existing outer geometry/routing; do not blanket-capture the entire field wrapper as an editing surface. |
| Caret navigation/typing | Entry reveals caret internally; runtime stores only the text draft. |
| Field focus | Existing EnsureVisible reveals the control through ancestor content clips; it does not export or set an Entry offset. |
| Clipped child/canvas | Native hits honor ControlClip and ancestor content clips, including exclusion of ancestor gutter strips. Open surfaces retain their own canvas coordinates and size. |

The current routedClip.Scrolled forwards when invoked. This alone does not mean
an Entry wheel reaches that wrapper: Noether's subsequent read-only inspection
reports the pinned Fyne wheel finder descends into the native renderer and finds
the inner Scroll, which consumes even at its limit. Preserve that native subtree
and test exactly one consumer and zero outer routes before adding a routing adapter.
A wrapper-level coordinate check alone is not proof of native delivery. Distinguish
Entry from label/feedback/background through actual event targeting. Do not alter
legacy input wheel behavior. If evidence needs inner scroll observations, Noether
reports the public container.Scroll alias permits bounded renderer-object inspection;
that is host-side observation, never a new shared offset or geometry authority.
The existing gutter mechanism still reserves each outer owner's own strip once;
never reserve a second shared gutter for Entry's private scroll bar.

Measurement/gate work must use detached metrics and cannot resize, hide, focus or
SetText on the live Entry. Existing accepted presentation tickets publish geometry
and field state together across main and open surface canvases. Failed geometry,
resource preparation, field publication or reload leaves the live object and edit
state untouched by candidate preparation. A native OnChanged callback can occur
after Entry has already edited itself: a pure runtime gate alone does not prove
native rollback/history preservation on rejection. Host must separately verify
its rejection handling rather than claim transaction purity from layout tests.
Hidden/inactive/closed surfaces produce no active geometry, while
source admission and runtime validation still cover the selected tree.

ReadOnly must allow focus, selection and copy while blocking editing/cut/paste;
Disable is not an implementation of this contract. Fyne's multiline OnSubmitted
default is Shift+Enter, whereas the reviewed contract requires Primary+Enter.
Current Input delegates native editing shortcuts and command handling, so an
extended host path must establish the reviewed ordering without changing legacy
input: Enter newline, Primary+Enter one explicit Commit, blur/Tab no Commit,
Escape composition first, then dirty draft, then clean surface Cancel. Layout
cannot certify these behaviors or IME-consumed Enter through geometry.

## 5. Native Entry retention policy

| Transition | Required native and runtime treatment |
| --- | --- |
| Ordinary sync, validation, command or unrelated state change | Keep the same Entry and caret/selection/scroll/history. Do not SetText for identical draft. Feedback/enable/readOnly sync must not reconstruct the editor. |
| Page hide/show or split collapse/restore | Keep the existing Entry and its edit state even when shared geometry is absent. Do not build a replacement when it becomes visible. |
| Resize/wrap reflow | Keep Entry; resize only after accepted geometry and clamp native offset. Text/draft/Commit counts unchanged. |
| Failed edit publication or failed reload | No live Entry or accepted ticket replacement, no text/history/focus/offset change from candidate measurement. |
| Explicit accepted programmatic text replacement | Replace text once and reset that field's native history; do not turn it into user Change/Commit. |
| Successful compatible main/page reload | May recreate Entry and reset caret/selection/internal scroll/history. Retain compatible accepted/draft text and field focus, including inactive pages; do not promise exact cursor position. |
| Successful reload with closed dialog successors | Reset unaccepted dialog drafts to current accepted values, preserving explicit child commits. Failed reload instead retains the old open form and edit state. |
| Removed/type-changed field or disposal | Retire native object/callback authority under existing host/runtime identity rules. |

Do not serialize native editing state into layout or Snapshot.Viewports to simulate
retention. Compatible text revalidation and single-line conversion rejection are
runtime responsibilities. Native object lifetime and programmatic-versus-user sync
are host responsibilities; layout publishes only pure candidate rectangles.

## 6. Evidence plan after selection

Synthetic layout tests can prove exact rectangles/clips, finite bounded minimums,
optional empty input label without scalar regression, positive feedback, separate
readOnly/disabled flags, invalid adapter rejection, and detached snapshot behavior.
Cover long/newline/Unicode drafts with fixed synthetic metrics, nested relative
scrolling, ancestor gutters, tabs, collapsed splits, hidden fields, dialog-local
canvases and independent nonmodal resize. Force a failing field/surface geometry
candidate and assert runtime revisions, offset maps and accepted ticket unchanged.
Assert no input entry in Viewports and no inner-scroll contribution to outer extent.
Retain exact 0.2/legacy-0.3 geometry and unsupported-scroll regression assertions;
explicit false/empty opt-in must use the new path only after parser support lands.
Synthetic metrics cannot prove actual wrapping, cursor, clipping or native wheel
routing, and fixed measured minima do not prove the native adapter is bounded.

Host tests with real Fyne controls can inspect native MinSize across font, long
text and width changes, object identity through sync/page/collapse, SetText/history
reset boundaries, emitted event counts and failed candidate purity. They supplement
rather than replace OS evidence. Main's native cases must exercise long wrapped
text/caret at both ends, Entry wheel at both limits, Entry scrollbar and separately
reachable outer gutter, resize/clipping, Unicode selection/copy/cut/paste,
readOnly copy, keyboard submission distinctions, undo/redo and all retention rows.
Use accepted surface-local geometry when targeting existing nonmodal windows.
Actual OS IME preedit/commit/cancel is separate A08 evidence; software canvas tests
or injected Unicode do not prove it. Existing black Canvas.Capture limitation
requires actual OS screenshots for visual acceptance, not fixture capture claims.

The canonical follow-up now approves retention of the same invalid editable
exactly empty accepted baseline for compatible already-required select/input fields,
including closed forms, with invalid feedback, unchanged accepted revision and no
automatic selection. Whitespace-only text is not covered by this narrow exception.
Newly-required blank accepted values and removed/disabled
nonempty option IDs still reject replacement. Closed-dialog successor draft reset
precedes single-line conversion validation; retained accepted text still constrains
conversion. These are approved runtime semantics, not pending layout decisions.
The canonical single-line predicate is now approved as exactly CR or LF, without
broader Unicode normalization or tighter legacy 0.2 behavior. Geometry presents
copied validation without deciding text validity; no layout dependency or frozen
M1 layout change is needed.

## 7. Inspection identity and handoff

SHA-256 of key inspected working files:

| File | SHA-256 |
| --- | --- |
| SDUI/go/layout/fields.go | 889f8f3b78ab28515b40ad51ff73b8dd185bd81f34321349cced413f8e827d61 |
| SDUI/go/host/fynehost/input.go | c337660906f214d3ee5b36093ce27a639c5e919f2b719d64581df55acfa55623 |
| SDUI/go/host/fynehost/document_sync.go | 0daba7ec0a627586fc183a6b8f46ce4fcb9da2274adda42124a4627112754e0c |
| SDUI/go/host/fynehost/scalar_control.go | a42aa7a49f9207b35feec989681b89006647b8d8b426581e44ac93064b39de4b |
| SDUI/go/host/fynehost/viewport_control.go | 103162ec62036104ef4ca95ff815933876cb225a6f81164d8052dedf7b7d5888 |
| SDUI/go/runtime/field_types.go | 88a91fa823e2272f17d5cb34bb2dff9aabb5ab34ab35b790a93545e92bb670c6 |
| Fyne v2.8.1 widget/entry.go (local module source) | 67ed73c6c8825ed28df1997368c5ec0a9e11db8f8f26419bb92d67df44804b86 |

Inspection used git status/rev-parse, rg, sed/cat and sha256sum; no product tests,
native run or M2 capability claim. Only this root memo was written. Direct read-only
coordination sent to Noether, Dalton and James. Noether agrees with whole-Entry
ownership and subsequently proposes three visible rows, word wrap, empty-label
omission and the existing fixed feedback row; Dalton agrees with the sole
opt-in helper and legal empty label. James confirms no projection mismatch and
that the proposed seam preserves current snapshot/gate boundaries. No request to
disturb M1 freeze was made.

Main Session0010 S4 handoff: owner requested bounded M2 preparation only; existing
FieldMeasurer shape suffices, with a proposed input-only empty-label refinement
and native adapter/routing/lifetime work. No material contract departure identified.
Next step is M1 closure and explicit M2 selection, then owners confirm the proposed
helper and exact native measurement policy before dependent implementation. Main
records this work summary in Session/roadmap; no management records were edited.
