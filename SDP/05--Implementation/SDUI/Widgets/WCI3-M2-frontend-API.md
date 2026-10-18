# WCI3-M2 frontend seam preparation — no implementation selected

Canonical selection note: M2 is now selected after WCI3-M1 delivery. The reviewed
Values-and-text native editing refinement supersedes earlier read-only/pending
status and broad history promises in this preparation record: identical bytes
preserve history; changed programmatic text resets it; an actual rejected native
edit may reset that Entry history/caret/selection/scroll while restoring current
authoritative draft muted. Failed Commit/reload/probe and admitted invalid drafts
retain history. Exact captured CR/LF clipboard content rejects before native paste.

Historical preparation status: read-only reconciliation plus this memo only, 2026-10-08. M1 frontend
remains frozen. SDP Worker and shared document workflow reused; main owns Session,
stage selection, dependency/IME changes, integration and native evidence.
Inspected shared clone branch `sdui/widgets-wci3`, HEAD
`a28b3ccc72b1ecc44300d4aa7c3cd679f4bc9625`, with concurrent uncommitted M1 work.
Authority is ORIGINAL `SDP/04--Design/SDUI/Widgets/Values-and-text.md`, SHA-256
`719cb7eab8375f20bbedb994fc0dc2eb40928b43903d5270238ce07629220309`.
Read with WCI3-frontend-API.md and WCI3-runtime-API.md, checking actual exports.
Proposals below are for selection/review, not new implementation authority.

## 1. Closed source schema and compatibility switch

Extend only .3 input, using its existing widget-call production:

```text
input(text, value="", multiline=false, readOnly=false,
      placeholder?, required=false, callback?)
```

text/value/placeholder are strings; multiline/readOnly/required are booleans;
callback remains the symbolic member-reference for Commit. Only text may be
positional, first; unknown/duplicate/wrong-type arguments reject. Retain named
callback ownership and connected @invoke validation. No body, textarea, wrap,
onChange, onCommit, validator expression or parser I/O. Formatting enabled/visible
keeps its current location. Do not impose a new nonblank-label rule on old input.

Proposed opt-in: presence of ANY new input argument selects extended policy,
including explicit false and empty placeholder. No such arguments means legacy
input behavior under both .2 and .3. James agrees with this handoff boundary;
coordinator should record it before M2 implementation. Profile alone, callback
presence and registration of a handler must not silently promote old sources.
All four new arguments remain invalid in .2. No source migration is necessary.

Proposed sole frontend facts API, not yet exported:

```go
type InputPolicy struct {
    Extended, Multiline, ReadOnly, Required bool
    Placeholder string
    PlaceholderSet bool
}
func InputOptions(n *Instance) (InputPolicy, error)
```

InputOptions accepts an input Instance, validates its effective profile and closed
argument types, and returns effective defaults without mutation. Placeholder is
the explicit string when supplied (even empty), otherwise text. PlaceholderSet
records presence; Extended is the new-argument presence test above. Invalid input
returns a parser.Diagnostic with the original offending span. Consumers do not
reimplement a lexical parser or infer these facts from rendered text.

AST fields and literal kinds remain unchanged: boolean/string/reference arguments
already suffice. Do not add normalized default arguments or derived input keys.
Keep Instance.Profile's empty .2 sentinel, sdui-ast/0.3, sdui-go-model/2, spans,
reuse chains and generated constructor metadata. NumericArguments remains numeric
only. ResolveInteractions remains the strict selected-root resolver and already
provides input dialog ownership; ResolveDialogField stays text-input-only.

## 2. Actual M1 runtime and proposed bounded extension

M1 Field(input) and Snapshot.Fields project Widget.Value/Draft as String values,
with Widget value/draft revisions and copied RawDraft. initFields registers only
the four scalar kinds. EditField, CaptureCommit, ObserveChanges and
ValidateFieldWith currently reject input; the projection is not typed-edit support.
Legacy Draft checks byte length and input eligibility; legacy Commit has no Control.

Reuse these existing signatures for extended inputs after selection:

```go
EditField(Handle, uint64, Value) (FieldChange, error) // model revision == Session.Revision
Field(Handle) (FieldState, bool)
CaptureCommit(FieldTarget) (Event, error)
ObserveChanges(Handle, ChangeHandler) error
ValidateFieldWith(Handle, FieldValidator) error
RevertField(Handle) error
```

Preserve EditField's M1 model-revision parameter; FieldTarget separately guards
draft/value/state revisions for capture. Do not reinterpret that uint64 parameter.

Keep Widget.Value/Draft and their revisions as the sole authoritative text store.
Extended field metadata/validators may be registered, but Accepted/Proposed and
RawDraft in returned FieldState are copied projections, not a second mutable text
store. No new Value kind, Session, field framework or input-specific event type.
Draft/Revert public signatures remain; extended input calls must use the same
guards/publication as the new path, not bypass readOnly or emit duplicate Change.
Legacy calls retain their present behavior. Checked Apply remains the distinct
programmatic path, including readOnly receiver updates, without user notifications.

James prefers CaptureCommit for opted-in input with Event.Value=Text(exact draft),
Control.ValueRevision plus existing model/state/draft/sequence guards. Proposed
Control.RawDraft=nil and Option=nil: String already preserves all text bytes;
numeric provenance rules stay number-only. FieldState/DraftField RawDraft may
remain the existing copied text projection. Reject stray payloads and stale captures.
Keep the old envelope valid only for legacy input; no envelope-based policy bypass.

Extended edits publish copied String Change after accepted prospective-state
validation; no implicit SDL Commit. required tests non-whitespace without trimming.
UTF-8 and 32768-byte bounds remain explicit. Invalid drafts remain visible and
cannot Commit; Go validators use the existing pure bounded FieldValidation contract.
ReadOnly blocks user edits/Commit, not focus/selection/copy or checked Apply.
AcceptDraft must match the captured proposal AND value/draft revisions. A bound
handler without source acceptance cannot claim accepted Commit. Existing unbound
local acceptance versus open-dialog proposal-only behavior is reused.

Reuse StateGate, PresentationGate/tickets, atomic Apply and Successor. Mixed Go
Accept validates all owned fields, including hidden pages and excluding nested
dialogs; DialogFieldValue remains text-only. No new form transaction mechanism.
Invalid retained main/page drafts remain visible; invalid retained accepted values
reject reload. Single-line conversion with newlines rejects. Closed-successor
dialogs reset unaccepted drafts to current accepted values, retaining child commits.

## 3. ControlText/TextResult bridge reconciliation

Actual ControlText accepts ONLY basic input Commit and explicitly rejects nonnil
Event.Control. This is a real extension seam, not existing M2 support. Extend it
with two policy-checked paths: legacy validation unchanged; extended input requires
the exact typed capture above, validated again after execution. Both map exact
String to SDL text, preserving Unicode/newlines. No label/placeholder/value guessing.
Keep Source.Event/Widget conversions and .2 callers unchanged.

TextResult remains the zero/default text/input result mode for Commit and Load.
ScalarResult remains restricted to checkbox/slider/number/select; never switch a
text receiver into ScalarResult by widget guessing. Existing TextResult already
checks receiver handle/value/draft and engine revision around Execute. Extend its
returned Update with the captured ExpectedDraftRevision for extended receivers and
the required final guarded publication; do not fabricate a second acceptance step.

Material decision for selection: existing TextResult permits a Commit result to
target another input and returns AcceptDraft only when result equals its draft.
The extended Commit contract requires acceptance of the SOURCE's exact capture.
Proposed minimal bound: extended input Commit using TextResult requires explicit
self setHandle and an exact echo; preflight rejects a different receiver. Explicit
Go handlers may return checked source acceptance plus other updates within the
existing batch limit. Preserve legacy nonself bindings and Load-button -> input
receivers, including readOnly. Do not discover this incompatibility after persisting
domain work or silently accept the source on behalf of the application.

Load delivery checks the captured receiver/model/engine and both revisions;
newer dirty draft or replaced receiver rejects delivery. Validation/shape/preflight
failures dispatch zero actions. Delivery failure after domain success is reported,
never automatically replayed. Text-only SDL dialog Accept and mixed Go Accept keep
their existing distinct contracts and atomic publication.

## 4. Native, capabilities, exports and snapshots

Reuse native Entry with bounded adapter: single-line Enter Commit; multiline Enter
newline and Primary+Enter Commit; blur/Tab no Commit. ReadOnly selection/copy must
remain available. Undo/redo edits draft only. Escape first cancels preedit, then
dirty draft, then permits clean-field surface Cancel. IME-consumed Enter must not
reach SDL. Main owns the reviewed pinned GLFW patch/build roots and actual IME proof.

Native Entry owns caret/selection/undo/internal scrolling; runtime Viewports remain
outer geometry only. Retain Entry across sync/page hiding/collapse/failed publication.
Successful compatible reload may recreate/reset native edit history while retaining
declared text/draft/focus; explicit programmatic replacement resets field history.
Multiline word-wrap/vertical scroll stays inside the measured control rectangle.
Entry consumes its text-area wheel even at a limit; outer gutter remains reachable.
Shared layout/EnsureVisible still owns outer rectangles/clips/ancestor reveal.

Preparation walks the full selected tree, including hidden pages/dialogs, using
InputOptions plus strict ResolveInteractions. Keep current frontend/layout/widget/
viewport/provider/host dimensions: existing widget+host input/1, additional host
input-multiline/1 when true and read-only/1 when true. required/placeholder need
actual adapter support but no invented provider or bridge capability dimension.
Binders establish SDL readiness; parser success/native inventory alone cannot.
Standalone native preview remains explicitly unsupported for .3 document-host
adapters, with source-linked diagnostics and no connected/provider readiness claim.
Main/reviewer confirm no mandatory legacy-launcher migration in WCI4; the actual
connected route must be explicitly named and staged, not inferred from that launcher.

Proposed public SVG policy: any extended input rejects unsupported-text-export with
normalized path and original span, even false/empty arguments, hidden nodes or
supplied state/geometry. No snapshot renderer required. Native SkipControls omission
still requires the exact actual prepared input control inventory; no broad skip or
unknown-kind fallthrough. Kind remains input, not a new multiline widget.
Text/composition truthfully describe explicit properties and source initial text;
callback references are not proof of binding. Legacy default exports remain byte-
identical, including .2 AST/JSON/Go/text/Markdown/SVG output.

SnapshotRoot's existing input value projection must remain schema-valid and preserve
source metadata; add no validation/caret/draft/provider arguments. Native text and
validation read Snapshot.Fields plus Widget strings. Generated constructors retain
exact supplied new arguments, including explicit false/empty, without live state.

## 5. Remaining decisions and proposed verification

Two material seams need coordinator/bridge disposition before dependent code:
the explicit-argument opt-in boundary (§1, agreed by runtime) and extended TextResult
self-acceptance constraint (§3). Neither is implemented or an M1 blocker. Also make
extended initial/single-line text validation explicit in implementation tests:
required-empty must remain an editable invalid draft, not a new grammar error;
invalid UTF-8/oversize and newline-bearing single-line accepted text must not bypass
runtime/preparation checks. Preserve previously admitted legacy input behavior.

* A01: all new args/types/duplicates; .2 rejection; absent versus false/empty;
  reuse/selected-root ownership/spans; compile and run generated constructors;
  exact legacy AST/Go/export bytes; no default injection or new AST fields.
* A02/A06: legacy Commit unchanged; typed extended capture accepted; wrong envelope,
  invalid/required/readOnly/stale draft/value/state rejected; one Change and one
  explicit Commit; missing/extraneous acceptance and failed gates publish nothing.
* A06: actual ControlText/TextResult echo; Load -> readOnly input; hidden-control
  preflight zero actions; bad result/receiver, dirty receiver during Execute, engine
  replacement and post-domain delivery failure; legacy nonself result regression.
* A07/A08: native Unicode/newlines/clipboard/undo/redo/readOnly, Enter distinctions,
  text/outer scrolling, IME composition/commit/cancel and zero preedit actions.
  XTest Unicode insertion or headless tests cannot substitute for actual OS IME.
* A09: ordinary sync preserves Entry; failed reload preserves editing state;
  successful main/page and closed-dialog policies, tightened single-line constraints,
  removed handles, truthful descriptions and explicit SVG/native inventory boundary.

This is inspection evidence only; no M2 product/dependency edits or test-pass claim.
Only WCI3-M2-frontend-API.md is written. Canonical/Session updates and selection
remain with main; runtime/bridge/host/layout extensions remain with their owners.
