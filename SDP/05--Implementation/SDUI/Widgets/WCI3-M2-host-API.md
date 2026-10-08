# WCI3-M2 host API preparation

Canonical selection note: M2 is now selected after WCI3-M1 delivery. The reviewed
Values-and-text native editing refinement supersedes earlier read-only/pending
status and broad history promises in this preparation record: identical bytes
preserve history; changed programmatic text resets it; an actual rejected native
edit may reset that Entry history/caret/selection/scroll while restoring current
authoritative draft muted. Failed Commit/reload/probe and admitted invalid drafts
retain history. Exact captured CR/LF clipboard content rejects before native paste.

Historical preparation status: read-only handoff, 2026-10-08. M2 implementation is not selected.
Only this root memo is written. M1 product stays frozen; no product, module,
management, Session or original-workspace changes. SDP Worker 2.0.0 and the
already loaded SDP/document workflow are reused. Main owns stage disposition,
integration, pinned GLFW changes/build roots and actual OS IME proof.

Inspected clone `/tmp/sdp-sdui-widgets`, branch `sdui/widgets-wci3`, HEAD
`a28b3ccc72b1ecc44300d4aa7c3cd679f4bc9625` plus the frozen M1 candidate.
Authority: ORIGINAL `SDP/04--Design/SDUI/Widgets/Values-and-text.md`, SHA-256
`cb52f4174f705d6d60c5ee2f68e76cd5b6efdb24d463945cc18b7c034a7a681e`, especially
§6 and the final opt-in/CR-LF/required-empty/closed-dialog refinements. Read with
the three root WCI3-M2 frontend, runtime and layout API memos. Earlier pending
wording in those memos does not override the reviewed canonical refinements.

## 1. Bounded seams and compatibility

Consume the proposed sole frontend `parser.InputOptions(*Instance)` helper.
Presence of any new argument selects extended input, including false/empty;
legacy .2 and basic .3 stay on their current Input/events/geometry path. Do not
infer opt-in from profile, callback presence, a Field projection or defaults.

James proposes `FieldState.Input *InputState`, where InputState contains
Multiline and effective Placeholder; nil means legacy/non-input. ReadOnly,
Required, Validation and Field.Target already exist. Widget.Value/Draft remain
the sole runtime text store. Host keeps no second accepted/proposed text model.
The Entry's native text is its editing buffer, reconciled with that projection.

Use `EditField(handle, modelRevision, Text(raw))`, `CaptureCommit(FieldTarget)`,
`RevertField` and existing Dispatch/publication tickets. Newline and required
feedback comes from runtime; source parsing, validation and formatting are not
duplicated in the widget. Field.Target returned by an edit remains the target
for an automatic action associated with that edit; never recapture after an
observer to erase a conflict. Enter is a separate explicit commit gesture.
No text provider, generic event system or new transaction API is proposed.

Preparation inventories every extended input, including hidden/closed contents.
Retain widget+host input capability checks; require host input-multiline/1 when
true and read-only/1 when true, with actual adapters behind the claims. Bindings
still establish SDL readiness, including extended TextResult's explicit self
receiver. Neither native control presence nor parser success proves binding.

## 2. Public Fyne facts and the rejected-edit boundary

Inspected pinned Fyne v2.8.1 `widget/entry.go`, not a later upstream version.
Public Entry exposes Text, MultiLine, Wrapping, Scroll, placeholder, cursor
row/column, CursorTextOffset, SelectedText, OnChanged, OnSubmitted, Undo/Redo,
SetMinRowsVisible and native keyboard/clipboard methods. It exposes no complete
selection/history snapshot, edit veto, undo transaction or preedit query.

Concrete source order:

* TypedRune/TypedKey update text, cursor/selection and the private undo stack
  before invoking OnChanged. Validator supplies feedback; it does not veto edits.
* Adjacent characters can merge in entryUndoStack.MergeOrAdd. Typing an accepted
  `a`, then a rejected `b`, can leave one native `ab` action. Calling Undo to roll
  back `b` can remove `a` too and changes redo/cursor state.
* SetText restores text and then clears the entire undo stack. Calling it under
  callback mute prevents another runtime edit, but cannot preserve native history.
* Undo/Redo themselves change the stack before OnChanged. Direct Text assignment
  plus Refresh does not restore the private stack/selection consistently.

Therefore a pure runtime gate does not make a post-OnChanged native edit atomic.
An exact native-history rollback is not implementable with these public methods.
Prechecking readonly, UTF-8, size and known paste policy avoids some failures,
but cannot eliminate arbitrary gate/reentrant rejection. Rebuilding a parallel
text editor/history or accessing private fields is outside the bounded adapter.

**Disposition needed before dependent M2 code:** recommend a narrow recovery
exception for a rejected native edit: retain the Entry object, reread the current
authoritative runtime draft after the complete callback/reentrant operation,
restore it with muted SetText, reset native undo/redo, and report the rejection.
Do not restore a stale pre-callback draft over accepted reentrant work. Exact
selection/caret/internal-scroll retention is not promised on this recovery path;
use native clamping/reveal. This exception is a proposal, not current authority.
It must not apply to accepted-but-invalid drafts: those stay visible/editable with
feedback and keep history. Rejection before delegation changes no native state.

Failed reload, failed detached preparation, failed Commit and unrelated failed
publication still preserve the live editor/history; none requires native text
rollback. A rejected OnChanged candidate has already crossed a different boundary.
If exact rollback is mandatory there too, return a dependency/API decision for
a bounded Fyne edit transaction. The selected GLFW filtering work does not supply
that Entry API. No such dependency or exception is implemented by this memo.

## 3. Retention and programmatic replacement

Keep the same Entry through ordinary sync, validation, page hide/show, split
collapse, resize and failed reload. Candidate measurement uses detached controls;
it never SetText/Resize/hides/focuses the live editor. Publish only through the
existing accepted preparation ticket. Preserve accepted reentrant work even when
an outer callback returns an error, as the current host lifecycle already does.

On accepted sync, avoid SetText when the native text equals the runtime draft.
An exact self-echo Commit changes accepted revisions without replacing displayed
text; resetting history on ValueRevision changes would be wrong. Different-byte
programmatic replacement uses one muted SetText and intentionally resets history.
Successful compatible reload may recreate Entry and reset native editing state
while preserving the runtime's declared main/page draft/focus policy. Closed
dialog successor reset and single-line conversion precedence remain runtime-owned.

**Second disposition needed:** identical-byte explicit Apply and self-echo
acceptance cannot be distinguished by current value/draft revisions or Dirty.
James confirmed that reentrant Apply defeats an inference from local callback
scope as well. Recommend clarifying that identical displayed bytes retain native
history; actual changed-text replacement resets it. If identical-byte Apply must
reset history, accepted-publication provenance is needed and must be selected
explicitly with runtime. Do not invent an origin flag or infer it from counters.

## 4. Keyboard, readonly and clipboard adapter

Implement an extended-only Entry adapter/composite; preserve the existing legacy
Input behavior. Use the same native renderer, selection, clipboard and undo engine.

| Input | Proposed bounded route |
| --- | --- |
| Single-line Enter | Intercept Return/Enter before native selection handling; capture current typed field and dispatch once. No text mutation. |
| Multiline Enter or Shift+Enter | Delegate native newline insertion with OnSubmitted unset, so Fyne's default Shift+Enter submit branch cannot run. One ordinary draft edit, no Commit. |
| Primary+Enter | Match exactly the platform primary modifier through CustomShortcut, before native Entry handling; one typed Commit without erasing selected text. Extra modifiers are not silently equivalent. |
| Tab/Shift+Tab | Override AcceptsTab to false; retain WCI2 source/surface traversal. Blur does not Commit. |
| Escape | Driver first filters IME-consumed key; then existing open-menu dismissal, dirty-field Revert, otherwise clean-surface Cancel. A dirty-field Escape must not also close the surface. |
| Undo/redo | Delegate named native shortcuts after eligibility checks; resulting OnChanged edits the draft once, never Commit. Preserve native redo alternatives and exact modifier matching. |
| Readonly | Remain enabled/focusable for pointer selection, navigation, SelectAll and Copy. Consume typing, newline/Commit, delete/backspace, cut/paste, Undo/Redo and word-delete shortcuts without mutation. |
| Other command shortcuts | Forward only unhandled declared shortcuts to the existing canonical command route, once. Native navigation/editing wins only for its exact supported modifier forms. |

Fyne's selectingKeyHandler can erase a selection before typedKeyReturn, so
submission interception must precede Entry.TypedKey. Delegate KeyDown/KeyUp for
native Shift selection and clear adapter modifier bookkeeping on FocusLost.
Test focus loss with held modifiers; no key-state timers or OS polling fallback.

Readonly is not Entry.Disable: disabled has different focus/selection behavior.
Guard at the native entry points and again through runtime eligibility, including
late popup actions after readonly/activity changes. The native secondary menu
contains direct Entry.Undo/Redo actions, so blindly delegating that menu can bypass
wrapper guards. Use a bounded editing menu routed through the guarded adapter;
readonly offers Copy/SelectAll. Preserve the existing declared context-menu route.
Do not clone or reimplement the underlying selection/editing engine.

Pinned Entry.pasteFromClipboard replaces LF with spaces for single-line input.
That violates extended text preservation if delegated unchanged. Intercept paste:
read the clipboard once, reject the whole insertion if it contains CR or LF in
single-line mode, and leave selection/text/history untouched. Do not flatten it.
For legal single-line and multiline paste, pass the exact captured clipboard text
to native paste via a bounded clipboard value, preserving CRLF/Unicode and using
native selection replacement/undo. Enforce UTF-8/byte bounds without truncation;
post-edit runtime rejection still follows the separately dispositioned recovery.
Copy/selection retain exact bytes; readonly Cut must not alter either text or
clipboard. No extra normalization of Unicode separators or legacy behavior.

## 5. Geometry and one scroll owner

Agree with Gibbs's existing FieldMeasurer/FieldMetrics shape. Control is the entire
Entry including border and internal native scroll chrome; no new text viewport
state or +/- regions. Permit zero Label only for empty-label extended input.
Keep a fixed measured feedback row independent of validation message length.

Concrete proposed native policy: multiline Entry.MultiLine=true,
Wrapping=TextWrapWord, SetMinRowsVisible(3). Measure a detached identically themed
control with the actual font; add measured label/feedback exactly once. Long text
must not expand intrinsic height to document length. The same native configuration
receives the accepted finite Control rectangle. Single-line keeps native horizontal
caret reveal. Do not derive minima from placeholder/value length or guess pixels.

Pinned GLFW processMouseScrolled chooses one matching Scrollable through the
visible renderer tree; Entry's renderer contains its own Scroll. Public
container.Scroll aliases that native scroll type. Word wrapping selects vertical
scrolling, and its consumed event has no outer remainder callback. First preserve
that actual subtree and prove it wins over the ancestor routedClip; no synthetic
forwarding is needed merely because Entry itself lacks Scrolled.

Inside the Entry text region, native scrolling must consume at top/bottom too.
Label/feedback/background and the separately reserved ancestor gutter keep WCI1
outer routing. Never send the same event to Entry and outer RouteScroll. If native
hit testing in the composite requires an adapter, restrict it to this owned Entry
subtree and public renderer Objects/container.Scroll; do not access private fields
or introduce general routing. Public object inspection can report actual scroll
geometry/offset, but must not create runtime Viewports entries or a second authority.
Actual hit routing, limit consumption, clipping and thumb reachability need native
tests; source inspection is not proof of the final mounted composite.

## 6. IME boundary and required evidence

Main owns the pinned licensed GLFW source, X11 filtered-key patch, explicit
replacements in every maintained build root, packaging and actual IBus evidence.
The checked original ime-probe records the unpatched defect: composition Return
reached OnSubmitted before final character delivery; composition Escape also
escaped to Entry. Its isolated one-condition trial suppressed both consumed keys
and retained ordinary editing. That is mechanism evidence, not SDUI M2 acceptance.

Host relies on authoritative driver filtering: preedit produces no EditField or
SDL call; committed characters become ordinary draft edits; IME-consumed Return
cannot reach Commit and consumed Escape cannot Revert/Cancel. A later ordinary
Enter follows the selected submission policy. Entry has no public composition
state to reconstruct. No delay, text-change heuristic, key restamp or fake preedit
flag is acceptable. If actual configured IME violates the boundary, return the
concrete driver case to main; do not patch around it in this adapter.

After explicit M2 selection, bounded fixtures/tests should cover:

1. Legacy .2/basic .3 unchanged; false/empty opt-in, blank label, placeholders,
   readonly copy, required feedback and exact typed self-echo binding.
2. One Change per native edit, no implicit Commit; Enter versus Primary+Enter
   with a selection; Shift+Enter newline, Tab/blur, precise shortcut modifiers,
   held-modifier focus loss, readonly keyboard and late editing-menu actions.
3. Native merged typing/undo/redo, accepted-invalid draft retention, rejected
   gate/size/stale/reentrant edit recovery under the approved policy; failed
   Commit/reload preserve history. Same-byte echo/Apply and changed-byte Apply
   need separate history assertions once the disposition is selected.
4. Long wrapped Unicode/CRLF, long tokens, actual clipboard round trip,
   single-line newline-paste refusal without flattening, empty selection and
   insertion at beginning/middle/end. Copy must remain available readonly.
5. Entry identity/caret/selection/history/scroll across sync, page hiding,
   collapse and failed reload; successful reload may reset native state but
   preserves runtime draft/focus. Resize both main and nonmodal surface editors.
6. Actual wheel at both text limits, text scrollbar, separately reachable outer
   gutter, clipped controls and caret reveal; assert only one scroll owner.
7. Real configured IBus composition/commit/cancel inside the actual SDL fixture,
   including dirty draft and dialog Escape hierarchy, with correlated SDL counts
   and OS screenshots. XTest Unicode injection/headless tests cannot prove A08.

Retain existing controls/fields inspection geometry in each actual canvas. Expose
actual native Entry.Text and public cursor/selection observations only as bounded
diagnostics if needed; no duplicate state API. Report native inner-scroll geometry
only when observed from actual mounted objects. Main owns the external harness.

## 7. Inspection evidence and handoff

Read-only commands: git status/branch, rg, sed/cat and sha256sum against the clone,
canonical original and pinned local Fyne source. No M2 build, product test, native
run, browser research or capability acceptance is claimed. Relevant Fyne functions
inspected: TypedKey/TypedRune, selectingKeyHandler/typedKeyReturn, TypedShortcut,
registerShortcut, TappedSecondary, pasteFromClipboard, SetText, Undo/Redo,
entryUndoStack.MergeOrAdd, CreateRenderer and updateScrollDirections; GLFW
processMouseScrolled/capturesTab/shortcut dispatch and visible object traversal.

| Inspected source | SHA-256 |
| --- | --- |
| Fyne v2.8.1 widget/entry.go | 67ed73c6c8825ed28df1997368c5ec0a9e11db8f8f26419bb92d67df44804b86 |
| host/fynehost/input.go | c337660906f214d3ee5b36093ce27a639c5e919f2b719d64581df55acfa55623 |
| host/fynehost/document_sync.go | 0daba7ec0a627586fc183a6b8f46ce4fcb9da2274adda42124a4627112754e0c |
| Frozen scalar_control.go | 19049aa6b53560eabd0cb695ab85f2515a0f1fbda61328a0150ae9278cc29244 |
| Frozen scalar_control_test.go | 6367c744db1165f32b46f16cf3cd972123b87b4ec49442196bb255b10178f57d |
| Frozen scalar_choices.go | 7fbe95923eef796490d8f6a87b43aa97b3b420b3deb68b643c4fb44dd5e7bb01 |

M1 native artifact `/tmp/wci3-m1-noether-value` remains
`3c1abba76adca3cdaf9d710302ebc6a16575150acfb1e6effcd85bb9b18b2ef6`.
Main reports an independent matching build and slider11 native PASS, including
held proposal/value label/revert/programmatic update. Main subsequently reports
whole M1 approval received and closeout underway. These are coordinator-reported
results, not additional tests performed for this read-only M2 handoff.

Direct coordination sent the geometry policy to Gibbs, rejection/history findings
to main/Mendel, and the identical-byte provenance question to James. James confirmed
the indistinguishability and recommends the same bounded clarification. Main's
architectural proposal to the reviewer now selects these recommendations for
review: identical displayed bytes preserve history regardless of Apply/self-echo;
different programmatic text resets; actual rejected native EditField restores the
current authoritative draft muted on the same Entry with an explicit history-reset
exception; failed reload/Commit/probe retain edit state. CR/LF paste is refused
before delegation. Both material history refinements remain pending reviewer
disposition and M2 selection, not implementation authority;
routine keyboard/geometry/clipboard choices above require no competing API stubs.
Next step: main records the two dispositions and selects M2 before dependent code.
Main owns Session/roadmap updates; this worker changes only this authorized memo.
