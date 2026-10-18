# WCI3-M2 runtime API preparation

Canonical selection note: M2 is now selected after WCI3-M1 delivery. The reviewed
Values-and-text native editing refinement supersedes earlier read-only/pending
status and broad history promises in this preparation record: identical bytes
preserve history; changed programmatic text resets it; an actual rejected native
edit may reset that Entry history/caret/selection/scroll while restoring current
authoritative draft muted. Failed Commit/reload/probe and admitted invalid drafts
retain history. Exact captured CR/LF clipboard content rejects before native paste.

Historical preparation status: read-only handoff, 2026-10-08; M2 product implementation remains unselected.
Only this memo is written. M1 runtime/numeric stay frozen. Reused SDP 1.1.1,
Worker 2.0.0 and the shared document workflow. Main owns canonical records,
Session, selection, IME dependencies, integration and native proof.

Inspected clone `/tmp/sdp-sdui-widgets`, branch `sdui/widgets-wci3`, base HEAD
`a28b3ccc72b1ecc44300d4aa7c3cd679f4bc9625` plus concurrent M1 changes.
Authority: ORIGINAL Values-and-text §§2, 3, 5, 6 and final M2 reconciliation.
Read with WCI3-M2-frontend-API.md and WCI3-M2-layout-API.md. Dalton confirmed the InputOptions seam; subsequent canonical opt-in, self-TextResult,
required-empty and CR/LF decisions supersede earlier pending memo wording.

## 1. Compatibility and concrete seam

Canonical: presence of multiline/readOnly/placeholder/required selects extended
input, even explicit false/empty. Absent all four retains legacy behavior under
both .2 and .3. Profile, handler registration and event shape cannot select policy.
All four arguments reject under .2. Parser defaults are facts, not AST mutations.

Frontend proposes the single source-facts accessor:

```go
// parser; proposal, not a current export
type InputPolicy struct {
    Extended, Multiline, ReadOnly, Required bool
    Placeholder string
    PlaceholderSet bool
}
func InputOptions(n *Instance) (InputPolicy, error)
```

Runtime consumes it rather than inferring policy from strings or duplicating source
argument validation. Placeholder is explicit supplied text, including empty, or the
accessible text label when omitted. It is never an accepted value or draft.

Proposed only public runtime data addition:

```go
type InputState struct {
    Multiline bool
    Placeholder string // effective placeholder, not editable text
}
// Add to existing FieldState:
Input *InputState // nil for legacy input and non-input; nonnil means extended input
```

Use existing FieldState.ReadOnly, Required and Validation. No extra Extended bool
needed downstream; Input != nil is the runtime policy fact. Copy Input in every
Field/Snapshot/observer/candidate copy. PlaceholderSet is needed for frontend source
fidelity, not for runtime mutation. Keep kind input, Value String and numeric/option
payloads absent. Source metadata/explicit presence survives unchanged in the root.
This is a proposed additive API for lane agreement, not already implemented.

Reuse exact signatures:

```go
func (s *Session) Field(h Handle) (FieldState, bool)
func (s *Session) EditField(h Handle, modelRevision uint64, v Value) (FieldChange, error)
func (s *Session) Draft(h Handle, text string) error
func (s *Session) Revert(h Handle) error
func (s *Session) RevertField(h Handle) error
func (s *Session) ObserveChanges(h Handle, fn ChangeHandler) error
func (s *Session) ValidateFieldWith(h Handle, fn FieldValidator) error
func (s *Session) CaptureCommit(t FieldTarget) (Event, error)
```

No new Commit, transaction, text-scroll or reconciliation API. EditField's uint64
remains ModelRevision. Draft is the existing synchronous current-session helper;
its signature cannot claim a caller-captured model revision. Native callbacks must
use EditField with their captured model revision, then the returned Field.Target.

## 2. One authoritative text store

Actual M1 Field(input) is already a projection of Widget.Value/Draft and counters.
The private scalar field currently embeds FieldState and owns scalar values. M2
must not turn that embedding into a second accepted/proposed text store.

Recommended bounded change: register input policy/feedback in the existing private
field map, but keep its Accepted/Proposed/RawDraft unused for input. Field(input)
always reconstructs Text(w.Value), Text(w.Draft), a copied RawDraft and Dirty from
the Widget, then adds input metadata. Validation reads that projection. Existing
scalar storage remains unchanged. A small input-specific validation/reset branch
is sufficient; no generic storage abstraction or parallel input registry needed.

Audit every map-presence assumption: initFields must not clear input strings;
Field must prioritize the input projection over embedded scalar state;
edit/Apply/reset/successor must modify Widget strings; capture/Dispatch should
recognize opted-in input rather than only scalar(kind). copyFields copies only
input metadata/feedback, and snapshots copy the projected text. Reload's compatible
closure retention can use the existing metadata entry once these branches exist.

## 3. Editing, programmatic changes and notifications

| Path | Extended input behavior | Legacy behavior |
| --- | --- | --- |
| EditField | String only; check model/handle/activity/readOnly; write Widget.Draft and increment draft revision; validate prospective projection and publish once; one copied Change | Continue rejecting as a typed-edit API |
| Draft | Delegate once to the same checked edit path at current model revision; no extra notification or unchecked direct assignment | Preserve current signature, limits and behavior |
| Revert / RevertField | Restore Widget.Draft from Widget.Value, clear application feedback, recompute built-ins, publish once; never Commit | Preserve current behavior |
| Apply AcceptedValue | Checked String write to Widget.Value/Draft, exact baseline value/draft revisions; readOnly is allowed; no Change/Commit | Preserve existing legacy update requirements |
| Apply ReadOnly/ValidationState | Existing typed-property guards and batch baseline; validation metadata cannot bypass built-ins | Do not silently add typed policy |
| Observer / validator | Permit only opted-in input plus M1 scalars; bounded copies, pure validation, reentrance barriers | Continue rejecting registration as typed field |
| Initialize / reload / dialog reset | Silent; reconstruct projection from Widget strings | Preserve old contracts |

Proposed Revert details: extended public Revert is an edit cancellation and checks
activity/readOnly; lifecycle Cancel/Close uses the existing internal reset path,
which can discard drafts even when a control is hidden or has become readOnly.
Keep Revert silent like M1 RevertField; native undo/redo is an ordinary EditField
change and does notify. This distinguishes cancellation from programmatic
replacement without adding an event kind. Host must not synthesize Change again.

UTF-8 and 32768 bytes are admission bounds for extended text: reject malformed or
oversized edits atomically before storing them, without truncation. Required means
non-whitespace (Go strings.TrimSpace for the check only); preserve every supplied
byte. Whitespace-only bounded drafts remain visible with validation and blocked
Commit. Pure Go validators run after built-ins, return the existing 4096-byte
code/message contract, and cannot overwrite work accepted by reentrant callbacks.

Programmatic acceptance validates the final candidate text before publication.
AcceptDraft requires exact current proposal, matching captured value/draft revisions
and valid feedback; dirty unrelated Load delivery rejects rather than replacing it.
A readOnly preview receiver may receive a checked Load update. All distinct
properties compare with the original batch baseline, preserving the M1 P2 fix.
Required/Multiline/Placeholder remain source policy; do not add mutable properties
for them. ReadOnly uses the already selected runtime property.

## 4. Exact Commit, bridge and dialog shape

Extended CaptureCommit produces the existing Event with:

```go
Kind: Commit
Value: Text(exactWidgetDraft)
ModelRevision, StateRevision, DraftRevision: captured current revisions
Control: &ControlCommit{ValueRevision: capturedValueRevision}
// Control.RawDraft == nil; Control.Option == nil
```

String already preserves the exact text; number-only provenance rules do not apply.
FieldState.RawDraft remains a copied text projection. Capture checks current
FieldTarget, validation/activity/readOnly without consuming sequence. Dispatch
requires the typed envelope for opted-in input; nil Control cannot bypass policy.
Legacy input keeps its existing nil-Control/zero-StateRevision envelope and rejects
a typed envelope. No event-supplied flag can promote/demote the current source.

Keep M1 callback order: pure probe, consume sequence before handler, exact source
and target reentrance guards, checked acceptance update, final prospective gate,
accepted ticket publication. A successful callback must explicitly accept its
captured source with exact String and both revision guards. Missing acceptance,
non-echo, stale receiver, validation or geometry failure publishes no accepted text;
after Execute it never causes automatic replay. Observers remain separate from SDL.
Unbound extended Commit outside a dialog may accept locally; inside an open dialog
it keeps the valid proposal pending for owner Accept. Declared unbound callback
still errors. Explicit child commits remain intentional and survive Cancel.

Bridge owner extends ControlText to policy-checked typed Commit while preserving
legacy ControlText and Source.Event/Widget mappings. TextResult stays text/input;
ScalarResult must not absorb input. Canonical now requires an explicit SELF receiver
for extended Commit using TextResult regardless of request selector; reject nonself
at detached preflight. Exact returned echo is checked only after Execute. Legacy
nonself Commit and Load-to-input remain valid. Add ExpectedDraftRevision to extended
receiver updates and preserve engine/model/value/draft guards before/after Execute.

DialogControls includes extended input through its projection; DialogFields and
DialogField remain text-only. CaptureDialog validates projected input state even on
inactive pages. Extend DraftField optional metadata to opted-in input if needed,
without changing legacy DraftField envelope bytes or text selector semantics.
Both text-only SDL Accept and mixed Go Accept use the current exact capture,
combined 256 writes, Domain outcomes, AcceptBlocked and published receipt rules.
Closing/resetting a field restores Widget strings first, then recomputes feedback;
there is no copied scalar proposal to reset independently.

## 5. Initial state, reload and newline precedence

The following concrete cases separate admission, an invalid editable proposal and
retained accepted data. The required-empty and newline-conversion precedence
below now have explicit main/reviewer approval in the canonical contract. Main
also recorded CR or LF as the bounded single-line line-break predicate.

| Case | Proposed result |
| --- | --- |
| New extended required input with empty/whitespace source value | Admit as editable invalid state; Value and Draft retain exact text, Dirty=false, required feedback; Commit/Accept blocked. Do not make required-empty a parser grammar error. |
| New extended text with malformed UTF-8 or >32768 bytes | Reject runtime/preparation before activation with source diagnostic; no truncation. |
| Required changed false -> true while retained accepted Value is blank, even if Draft is now nonblank | Reject successor: unaccepted draft cannot rescue an accepted value violating the tightened constraint. Old bundle/form stays live. |
| Nonblank accepted Value but blank dirty Draft with required=true | Keep draft visible and invalid through compatible reload; accepted state remains valid. |
| Unchanged required=true with accepted exact empty String absence | Approved: retain this same editable invalid baseline on unchanged-policy reload, including closed forms. Keep invalid feedback, accepted revision and no Commit; no automatic choice or general invalid-value retention. No accepted/initialized flag or new store is proposed. |
| Multiline -> single-line, retained accepted text contains a line break | Reject successor; do not strip, replace, normalize or split content. |
| Multiline -> single-line, accepted text valid but surviving main/page draft contains a line break | Reject successor: explicit single-line conversion rule takes precedence over generic invalid-draft retention. Include inactive pages. |
| Multiline -> single-line in a dialog closed by successful successor, only unaccepted dialog draft has line breaks | Reset draft to current accepted value first, then check surviving text. Permit if accepted text is valid; prior explicit child commits still constrain acceptance. Failed candidate does not reset the live form. |
| readOnly/placeholder changes alone | Preserve content; readOnly changes eligibility, placeholder changes display only. |
| Legacy input becomes extended | Same named input strings/counters may survive, but validate accepted text and the single-line conversion rule under the newly selected policy. A legacy event cannot operate the new policy after revision changes. |
| Extended input becomes legacy | Retain compatible Widget text under legacy policy; remove extended metadata/observers/validators, and reject old typed captures through model/policy checks. |

The approved absence exception is exact empty String; it is not permission to
retain arbitrary invalid accepted text. Whitespace-only required initial text may
be editable invalid, but is not silently folded into the exact-empty reload
exception. Newly-required blank text includes all whitespace-only content.

The approved line-break predicate is CR or LF, treating CRLF as the original
two bytes, never normalizing. Extended single-line source/programmatic accepted
text containing either rejects. A bounded user/paste draft containing either can
remain visibly invalid with `text-single-line` feedback and no Commit; native
adapter should normally prevent single-line insertion. The runtime never silently
flattens it. Do not broaden this check to other Unicode separator code points,
normalize text, or tighten .2/basic input behavior. Frontend, runtime, bridge and
host use the same approved bounded rule.

Before refinement, M1 `validAcceptedField` rejected every required-empty retained
select, even with unchanged constraints. Main and the independent reviewer now
explicitly approved the narrow already-required accepted-empty exception for M1
select and future M2 input, preserving new-constraint and removed/disabled-ID
rejection. This refines the old contract rather than relabeling its implementation
a defect. Main authorized only field_successor.go plus related tests/report for
the M1 adjustment; that separate bounded delta is recorded in
WCI3-M1-runtime-worker.md. M2 remains unselected. Required-empty initialization,
unchanged baseline, tightened policy and newline precedence need distinct tests.

## 6. Native object and publication ownership

Runtime stores no caret, selection, undo stack, preedit or Entry scroll offset.
InputState and Snapshot.Fields describe content/policy/feedback; layout publishes
bounded rectangles via the existing FieldMeasurer and accepted presentation ticket.
Probe must never create pending publication or mutate a live Entry. Failed
preparation/reload preserves the live native object, focus/history/scroll and all
runtime state; provider cancellation remains a consequence of accepted publication.

Host retains Entry through ordinary sync, validation, page hide/show, split
collapse/restore and resize. Do not SetText when the draft already matches; an
echoing Commit increments accepted revisions but is not a new text replacement.
Explicit programmatic replacement resets that field's history once, muted. A
successful compatible reload may recreate Entry/reset native editing history and
scroll while preserving main/page draft/focus; closed dialogs discard unaccepted
drafts as above. No new runtime counter for native history is proposed: host knows
local edit/callback context and the actual replacement value.

Native editing changes can occur before runtime accepts an OnChanged callback;
host must verify rejection rollback/history behavior separately rather than claim
that a pure runtime gate alone preserved native edit state. Candidate measurement
must never SetText/Resize the live Entry. The existing publication ticket is the
only state/geometry publication handshake.

Native host owns Enter versus Primary+Enter, copy/selection under readOnly,
undo/redo as draft edits, Escape composition/draft/surface ordering and wheel
consumption inside Entry. Preedit never calls EditField as committed text and never
Commit; Enter consumed by IME must not reach the SDL binding. Main owns dependency
patch selection and actual OS IME evidence. No runtime IME event or parallel editor.

## 7. Implementation ordering and proof after selection

1. Frontend freezes InputOptions with explicit opt-in and source diagnostics.
   Runtime freezes the small InputState projection addition with host/layout/bridge.
2. Runtime implements input metadata/projection, edit/reset/Apply validation and
   capture/Dispatch; preserve legacy tests byte-for-byte where applicable. Bridge
   can independently implement typed ControlText and self-TextResult against the
   agreed fields. No competing stubs or numeric package change needed.
3. Host/layout consume InputOptions and FieldState.Input, retain Entry and prepare
   pure measured candidates. Use Field.Target for exact typed captures. Main owns
   dependency/native integration. Presentation gate and schemas remain unchanged.
4. Test required-empty initialization, blank tightening, unchanged-policy exception,
   newline accepted/draft/dialog precedence, exact Unicode/CRLF, size limits,
   readOnly bypass through Draft/Revert, one Change/no implicit Commit, all stale
   captures, callback/validator reentrance, atomic mixed Apply and failed gates.
5. Test legacy nil-Control/nonself bindings, typed self/echo preflight/post-domain
   failures with zero replay, Load-to-readOnly, 256 mixed Accept and child Commit
   retention; .2/basic .3 exports/events/storage remain unchanged.
6. Host proves object identity, no redundant SetText, undo/scroll/focus retention,
   failed candidate purity, exact submission gestures and readOnly copy. OS proof
   separately covers clipboard, native scrolling and real IME; runtime tests do
   not establish those outcomes.

No product tests were run for this read-only handoff and no implementation claim
is made. CR/LF, the exact-empty select/input baseline exception, closed-dialog
newline precedence, opt-in and self-TextResult are approved canonical refinements.
The later native-history seam below requires coordinator disposition before
dependent M2 implementation.
Main should record the handoff in Session0010; this worker
makes no management/documentation edits beyond this specifically authorized memo.

## Inspection identity

Canonical original at initial inspection SHA-256: `424c67b4a88827c7c38a41dfed49f94cfe983c3d1e5e4b49f8f3a4dd3d4dbbba`.

| Inspected file | SHA-256 |
| --- | --- |
| `WCI3-M2-frontend-API.md` | `8599b76129a26f36163a0f1d1a94602edd1d7fbd6ea82a44479fb1dfd7a6da1f` |
| `WCI3-M2-layout-API.md` | `e4843b783944a1fdd45514b55a3f8c7e68dca9341395d30bf201daf4f562f2bd` |
| `SDUI/go/runtime/session.go` | `f19c8014dbb31b270956dedb930d25ac7e5295a96b6d6b5f49e385369d82d1c9` |
| `SDUI/go/runtime/fields.go` | `417b88501efcc3e19d07d7239fec4c64473ee7ba2834845c4e6f3f7762ace525` |
| `SDUI/go/runtime/field_types.go` | `88a91fa823e2272f17d5cb34bb2dff9aabb5ab34ab35b790a93545e92bb670c6` |
| `SDUI/go/runtime/field_commit.go` | `d68ccd9bfc05a26e14b1f6b88720a9d70c4763d875ad05c623f52d33dd492f63` |
| `SDUI/go/runtime/properties.go` | `193d0f5366706a9b2d5a7815c1a943c5d324dee296c1d5093bf76cfb3b95f46d` |
| `SDUI/go/runtime/field_successor.go` | `8ff617e8e74319fe7db1be2276215e805096068888e6fb49c517ddae64ec9670` |
| `SDUI/go/runtime/successor.go` | `b980001006c9756d89d2380959c13d0365de8201378c74c689c16f659e067465` |
| `SDUI/go/runtime/command_capture.go` | `3f9e3376c6e28bf5905c9b49b48960e25aef56ccfa32a2d3fcaac3214b9cab7a` |
| `SDUI/go/runtime/reload.go` | `ec670327fdb27800e97cc579a4649d035853ec8d0449e9023683cf033d7018eb` |
| `SDUI/go/host/fynehost/input.go` | `c337660906f214d3ee5b36093ce27a639c5e919f2b719d64581df55acfa55623` |
| `SDL/go/bridge/values.go` | `3de306696660d107d66fabb5b3b43319df4c43cfabb8c7a22f6509b26bf4eef5` |
| `SDL/go/bridge/scalars.go` | `ac890c9376e3acd34ec61164df0a70b00f1e189b19683c9b32e3ae4176481d80` |

Approved-refinement canonical SHA-256: `461c70a63f5a3111c2ebff837b3a12c21e59cd8c28c3ede4645129321da2415f`.

Main subsequently reported native pilot6034 `--required-empty`: four checks PASS,
including real beta SDL repair. This is coordinator-reported M1 evidence, not a
worker-run native test or M2 implementation claim. M1 source remains frozen.
Current CR/LF canonical SHA-256: `cb52f4174f705d6d60c5ee2f68e76cd5b6efdb24d463945cc18b7c034a7a681e`.

## Follow-up: identical-byte acceptance and native history

Noether identified a concrete limit to the no-marker proposal in §6. Exact source
self-echo Commit and an explicit Apply AcceptedValue of the same displayed bytes
can both advance value/draft revisions and clear Dirty. Current snapshots do not
identify operation origin. Callback scope cannot establish it reliably when an
explicit Apply reenters during a callback. Do not infer reset intent from counters,
Dirty transitions or a broad callback flag.

Recommended bounded clarification, not yet selected: preserve native history for
identical displayed text regardless of acceptance origin; reset it for actual
programmatic text replacement. This preserves self-echo history and needs no new
runtime field. If canonical instead requires a reset for identical-byte explicit
replacement, a minimal accepted-publication provenance signal is required; it must
publish with the same accepted ticket, distinguish source acceptance from explicit
replacement, and stay unchanged on failed candidates. No such API is implemented
or frozen by this memo. Main has been notified to resolve the observable policy.

Noether also confirmed that native Entry mutates private undo/selection before
OnChanged. SetText rollback clears history and Undo may combine accepted/rejected
characters. Pure runtime admission cannot establish exact native-history rollback.
Host is reporting this public-API boundary separately; do not claim that a proposed
origin signal would solve rejection rollback. M2 remains unselected, and product
source remains frozen.
