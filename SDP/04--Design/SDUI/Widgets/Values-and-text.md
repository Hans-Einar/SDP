# WCI3 values and text — reviewed stage contract

| Field | Value |
| --- | --- |
| Assignment | SDP Architect; PLAN-SDP-0022 WCI3-M1/M2 |
| Status | Independent design review approved; WCI3-M1 delivered and reviewed; WCI3-M2 selected after scalar acceptance |
| Authority | KB-SDUI-003 full inventory; bounded owner design assignment, 2026-10-08 |
| Parents | [Design](Design.md), [Acceptance](Acceptance.md), [Collections](Collections.md), [Panes/commands draft](Panes-and-commands.md), [Plan](../../../05--Implementation/SDUI/Widgets/Plan.md) |
| Obligations | SDUI-R05/R12/R16–R18/R23/R26–R28; GAP-XFMD-SDUI-005/006/008; Session0010 S4 |

## 1. Boundary and observed foundation

Deliver typed checkbox, slider, stable-option select, numeric input and ordinary single/multiline text.
Prefer Fyne Check/Slider/Select/Entry with bounded adapters; do not add a form engine, editor language,
expression validator, hidden collection widget or parser/provider file I/O. Persistence remains application-owned.
Original runtime has string Value/Draft fields, draft revisions, Apply and Commit; the WCI1 candidate adds
CollectionTarget, Snapshot, StateGate, StateRevision and Successor. The bridge still principally receives
text into input; typed scalar results below are a proposed extension, not existing support evidence.
Reconcile these seams with the reviewed WCI1/WCI2 candidate before code. No product tests were run here.

## 2. Exact proposed source and representation

Extend the same **unreleased `sdui 0.3;`** and generic widget-call production; no new version per family.
Keep `sdui-ast/0.3`, `sdui-go-model/2`, Instance.Profile propagation and original declaration/use spans.
0.2 schemas, default output, string drafts, public method signatures and legacy event semantics stay intact.

| Call schema (not literal source syntax) | Defaults / validation |
| --- | --- |
| `checkbox(label, value=false, readOnly=false, callback?)` | String nonblank accessible label; value/readOnly boolean; no tri-state. |
| `slider(label, min, max, step, value, readOnly=false, callback?)` | Label and all four finite numeric arguments required; min < max, step > 0; initial value in range/on step. |
| `number(label, min, max, step, value, readOnly=false, placeholder="", callback?)` | Same numeric contract plus an editable numeric text draft and increment/decrement. |
| `select(label, value="", required=false, readOnly=false, callback?)` | value is an option ID string, empty means no selection; required/readOnly boolean; options supplied through typed Go API. |
| `input(text, value="", multiline=false, readOnly=false, placeholder?, required=false, callback?)` | Keep text as existing accessible label and value as string; omitted placeholder retains text-as-placeholder behavior. |

Only the first label/text argument may be positional; subsequent arguments named, duplicates/unknowns reject.
`callback` remains a declared symbolic member-ref for **Commit**; connected use requires `@invoke`.
Callbacks require named widgets. `enabled`/`visible` remain existing formatting properties; do not duplicate
them as call arguments. No `textarea`, `onCommit`, source onChange hook or source validator expression.
Change is a typed local notification through Go (§3); it cannot implicitly call the Commit SDL binding.
New fields/widgets are invalid under 0.2. Existing input("...") sources retain their previous behavior.
For input, presence of any new argument (multiline/readOnly/placeholder/required)
selects extended typed behavior, including explicit false or empty placeholder.
Absent new arguments preserves legacy input under both 0.2 and 0.3; profile, callback
or Go handler presence alone cannot opt in. Effective defaults are read-only facts,
never injected source/AST arguments. All new arguments remain invalid under 0.2.

```text
sdui 0.3;
ref: settings "settings.sdl";
page=[
  active=checkbox("Active", value=false, callback=settings.AcceptFlag.@invoke);
  level=slider("Level", min=0, max=100, step=1, value=50);
  count=number("Count", min=0, max=100, step=1, value=1, callback=settings.AcceptCount.@invoke);
  mode=select("Mode", value="compact", callback=settings.AcceptMode.@invoke);
  notes=input("Notes", multiline=true, placeholder="Write notes", callback=settings.AcceptText.@invoke)
];
settings.AcceptFlag.setHandle(page.active);
settings.AcceptCount.setHandle(page.count);
settings.AcceptMode.setHandle(page.mode);
settings.AcceptText.setHandle(page.notes);
```

AST uses existing Node/Argument/Literal/Reference shapes. For **new slider/number numeric arguments only**,
retain the exact numeric token as `Literal{Kind:"number-lexeme", Value:<string>, Span:...}`; quoted numeric
strings still reject. Normalized Arguments retain that literal; runtime parses a finite binary64 value.
This narrowly preserves precision provenance for connected integer checks without changing 0.2 literals
or adding fields to its structs. Existing lexer/source numeric limits still apply. Generated constructors
retain these tokens; no reparse of original source or extra source-file dependency is needed at activation.

## 3. State, events and validation

Keep Widget.Value/Widget.Draft as strings with their current input meaning. Extend the existing typed Value
with Number (finite float64) and OptionID (ItemID, empty permitted), alongside String/Boolean; unused payload
fields must be zero. Add one Fields map to Snapshot for new control state: accepted/proposed typed value,
draft/value revisions, Dirty, ReadOnly and Validation{code,message}. Numeric editable draft remains a separate
string in that field state; never put formatted numeric/boolean/option values into the legacy input fields.
Input keeps its existing draft storage; its readOnly/placeholder/multiline/validation metadata uses the same
field-state extension. No second Session or generic property/event registry is introduced.

Reuse Handle/model/sequence/draft revision checks, StateRevision and WCI2 PresentationGate for every prospective
state change. Extend Apply with typed AcceptedValue updates and validation metadata; batches check every
field/value revision, option generation and receiver before publishing any result. ReadOnly rejects user
mutation/Commit but permits checked programmatic updates, focus, selection and copying. Hidden/disabled
controls reject user events. Programmatic updates, initialization and reload never emit Change/Commit.

User editing updates the draft/proposed value and revision, then emits Change to an optional typed Go
observer registered separately from the existing Commit handler. Text Change carries String; checkbox
Boolean; slider Number; select OptionID plus its generation target. Number typing carries its raw String
draft on Change, with validation status; only a valid Commit carries Number. Invalid drafts remain visible
and never become accepted values or SDL requests. Observers receive copied state, not native pointers.

Commit validates the current draft and captured revisions before invoking its handler. A callback's typed
Apply update accepts the value; no successful callback without such an update is reported as accepted.
For an absent callback, the local prototype adapter may explicitly accept the valid draft through Apply;
a declared but unbound callback reports unbound and retains the draft. Preserve existing 0.2 behavior.
AcceptDraft requires the returned value equal the validated proposed value and captured draft revision;
normalization to a different value requires an explicit application decision, not silently overwritten edits.

Built-in validation checks finite/range/step, selected-option existence/enabled state, required text/choice,
and current UTF-8/32768-byte text limits. Required text means non-whitespace content without trimming it.
Optional pure Go validators return a bounded code/message (4096 UTF-8 bytes), receive the typed proposal/raw
numeric draft, and perform no I/O. No asynchronous validation subsystem. Validation rejection retains draft
and accepted value, shows feedback associated with the label, and dispatches no domain action.

## 4. Numbers and choices

Number draft grammar is the existing finite source-number syntax, with no locale/group separators; empty,
"-" and incomplete decimal/exponent text may exist as invalid drafts. Do not round on text Commit or accept
NaN/infinity/overflow. Bounds/step/value are checked before activation and on programmatic updates.
Use one exact decimal grid from the retained min/max/step lexemes. Source initial
values and numeric text Commit must satisfy exact rational range and integer
`k=(value-min)/step` membership; never round an off-grid source or draft. This also
avoids cancellation in `(value-min)` at a large origin. Before arbitrary-precision
conversion, scan the decimal lexeme with overflow-safe bounded arithmetic. Permit
at most 32768 mantissa digits (all digits on both sides of the decimal point,
before trimming). For a nonzero coefficient, remove leading zeros and strip its
t trailing zeros; define effective decimal exponent as explicit exponent (zero
when omitted) minus fractional digit count plus t. Require its absolute value
<=4096. Scan explicit exponent with a saturating bound derived from 4096 plus
fractional digit count plus t; never parse an unbounded integer exponent or
construct its power before checking. An all-zero coefficient becomes exact zero
with effective exponent zero, without constructing any exponent power. Existing
whole-text/source limits still apply, including exponent spelling length.

Reject inadmissible source constraints during preparation; an inadmissible raw
draft stays visible with validation and cannot Commit. Source/text length alone
does not bound compact exponent expansion. Derive intermediate rational bounds
from these admitted coefficients/exponents, finite binary64 operands and bounded
legal ticks; operate directly on rationals rather than reparsing expanded decimal
strings. Reconstructed grid points do not acquire the lexical digit/exponent
limits again. No decimal expression language or unbounded arbitrary-precision input.

For a typed Go Number (which has no decimal lexeme), find the nearest integer tick
using the exact rational representation of its binary64 value and the exact grid.
Accept only if that Number equals the correctly rounded binary64 representation
of the reconstructed exact point `min+k*step`, with k in the legal range. Slider
movement and increment keys construct this same point then round once. Thus source
0.3 on a 0.1 grid is valid, large-origin legal increments remain valid, and a typed
off-grid value cannot pass through an expanding quotient tolerance. Raw text still
requires exact grid membership, even if an off-grid decimal would round to a legal
binary64 Number. No membership tolerance and no silent text normalization.

For all Go grids outside the exact safe53 integer exception below (including
fractional grids and integral constraints beyond safe53), define N=floor((max-min)/step) and legal ticks
0 <= k <= N <= 2^26. Require step strictly greater than the largest binary64
adjacent spacing at either rounded endpoint, comparing exact rationals (check
both directions, finite differences). Equality is insufficient: midpoint ties
can make adjacent legal ticks round to the same even significand. Reject an endpoint/step whose neighboring
representations or bounded range cannot be established. This prevents distinct
legal ticks collapsing to the same representable value; reject constraints at
preparation with a source diagnostic, never discover ambiguous increments during
editing. The no-cap exception applies only when min/max/step are exact integers within
±(2^53-1). These grids retain exact integer arithmetic and uint64 tick indexes;
all other Go grids use the bounded N/spacing policy above. Connected SDL grids
always require the explicit safe53 integer constraints below.
Pointer/step-key movement deliberately chooses a legal tick and clamps to the
last legal tick <= max. Native raw-Commit validation does not substitute that tick.


Supply `ChoiceOption{ID ItemID, Label string, Enabled bool}` slices by normalized instance path; at most
4096 options, unique nonempty IDs and label/ID limits reused from WCI1. Duplicate labels are allowed.
Initial explicit selection must exist and be enabled. Empty optional choice is valid; required empty choice
is an invalid editable state with Commit blocked.
A compatible already-required field whose accepted value is the empty absence may
retain that same invalid editable baseline across successful reload, including a
closed form. Keep validation invalid, retain its accepted revision and choose nothing
automatically. This narrow exception does not admit newly-required blank accepted
values or removed/disabled nonempty accepted IDs; those reject replacement. No parser I/O or automatic first-option selection.
Runtime owns an option-set generation: WCI3 introduces `OptionTarget{Handle, ModelRevision, OptionGeneration,
OptionID}` for selection events. Reuse WCI1 identity principles, but validate through a dedicated option-set
validator; `ValidateCollectionTarget` requires a collection and cannot validate a select. No hidden list or
loader. Replacement validates the whole set, increments generation
and revokes old menu/selection targets. Surviving enabled IDs retain value; removed/disabled accepted IDs
become explicitly invalid (retain ID for diagnosis), proposed selection clears, and user must choose again.
Never substitute an option at the old index or revive a deleted/recreated ID's old event token.

## 5. Real SDL bindings and transactions

Keep SDL action-core 0.1. Add closed bridge EventField selectors ControlBoolean, ControlNumber, ControlText
and ChoiceOptionID for validated Commit only. Legacy Source.Event/Widget text conversions remain unchanged
for existing 0.2 callers; new selectors never parse booleans/numbers from arbitrary text or stringify them.

| New selector / result receiver | SDL field type and exact rule |
| --- | --- |
| ControlBoolean / checkbox | boolean; direct typed value, no "true"/"false" string conversion. |
| ControlText / single or multiline input | text; preserve supplied Unicode/newlines; no label or placeholder as value. |
| ChoiceOptionID / select | text from/to explicit stable ID; capture/recheck option generation and eligibility; no label/index mapping. |
| ControlNumber / slider or number | integer only; min/max/step/value integral, absolute value <= 9007199254740991, checked before conversion both directions. |

Connected numeric preflight checks retained decimal lexemes exactly (e.g. rational arithmetic), then their
binary64 round trip; `9007199254740991.1` must not become an accepted integer through rounding. Repeat exact
checking for typed numeric drafts before conversion; programmatic Number values must be finite/integral
and bounded. Fractional values work only with typed Go adapters, never SDL text encoding or silent coercion.
An incompatible numeric binding rejects the detached candidate with a source-linked capability diagnostic.

Add closed Plan.ResultMode `ScalarResult` beside WCI2 TextResult/DialogAcceptResult.
ScalarResult is valid only for a WCI3 typed Commit callback with one explicit
setHandle field receiver; OutputField is required and its declared/registered/actual
result type must match the receiver (boolean, checked integer, or text option ID).
It forbids AcceptField/MessageField and retains existing revision guards. Legacy
TextResult still means text/input; it is never reinterpreted by widget guessing.
Text Commit/Load continue using TextResult. Unknown modes reject preflight.
An extended input Commit using TextResult must have one explicit self receiver,
regardless of which input selector populates the SDL request. Reject a different
receiver during detached preflight; neither legacy Source.Event nor Widget selectors
may bypass this rule. The actual returned text must exactly echo the captured draft
to accept it. Non-echo output is knowable only after Execute: reject delivery without
accepting the source or automatically replaying domain work. Preserve legacy input
nonself bindings and explicit Load-to-input receivers, including readOnly. Explicit
Go handlers may accept the captured source and return other checked updates.
No generic multi-result mapper. Fixture actions AcceptFlag/AcceptCount/AcceptMode/AcceptText each take and return a
record with one Value field, respectively boolean/integer/text/text, registered with real matching Go
signatures and a call log. Each uses the selector above and OutputField="Value" with the self receiver in
§2. Errors leave accepted state unchanged; echoes accept the proposed draft through one checked update.
Also exercise explicit Load button -> text receiver and Save text Commit -> text receiver with supplied
in-memory data. Input widgets/parser never read or write files. ReadOnly preview receivers accept programmatic
text results; a newer dirty draft or replaced handle rejects delivery instead of being overwritten.

Preflight validates all callback/result fields, modules, signatures and hidden controls with zero actions.
Before and after Execute, check receiver handle/value/draft revision, model/engine revisions and option
generation where applicable; retain global monotonic event/command sequence across Successor. Never replay
accepted domain work after UI delivery failure. Domain persistence, rollback and form transactions belong
to application adapters. For WCI2 dialogs, validate captured typed fields together and persist only on Accept
when atomic forms are intended. A field with no explicitly bound Commit handler
inside an open dialog retains its validated proposal as dirty on its automatic
checkbox/choice/slider/number-step Commit; it does not advance accepted value/value revision.
The same unbound control outside a dialog may explicitly accept locally. A source
callback or explicitly supplied Go Commit handler is an intentional child commit:
its accepted updates persist and later Cancel does not undo them. Thus applications
choose atomic forms by binding persistence only to Accept, without a new form mode
or hidden transaction layer. Rejected validation remains visible in either case; Cancel restores unaccepted drafts without Save. Earlier explicit child
domain commits cannot be undone by Cancel. Mixed valid/invalid Apply batches publish nothing.

Extend WCI2's captured `DraftField` with typed proposed Value, optional RawDraft (text/number),
FieldValidation and OptionTarget for select; invalid raw drafts never become accepted typed values.
WCI2's `DialogFieldValue` remains text-only. For mixed typed forms in WCI3,
select the explicit Go `InteractionHandler` Accept adapter: it receives validated
captured typed DraftFields and returns WCI2 AcceptDecision/Domain under the same
combined 256-write limit. It may call its own application persistence/SDL service;
SDUI does not invent typed record/schema mappings. A source SDL Accept declaration
whose plans require nontext DialogFieldValue rejects preparation, even if a Go
adapter might have been possible. Native acceptance must include a real mixed
Go Accept form and a text-only SDL Accept form, with Cancel and post-domain conflict.
Future generic typed dialog selectors are not required for this inventory.
Check captured handle/model/state/value/draft revisions and option eligibility immediately before atomic
Apply, without yielding. Current Update guards only ExpectedValueRevision; AcceptDraft/equal text alone
cannot prove that the captured draft is current. Keep these checks in the owning dispatch.

## 6. Native editing, focus and lifecycle

Prefer Fyne controls; adapt read-only selection/copy, option-ID lookup and step validation without parallel
state. Checkbox Space/click and accepted choice selection produce Change then one Commit; select Escape
cancels its popup proposal. Slider movement produces Change, release/key-up one Commit; Escape reverts the
uncommitted drag. Number typing edits its string draft; Enter commits, arrows/step buttons propose one legal
increment and Commit. Invalid numeric drafts block increments with validation; they do not silently reset.
Single-line Enter commits; multiline Enter inserts newline and Primary+Enter commits. Blur/Tab never commits.

Tab/Shift-Tab follow WCI2 surface/source order; focus movement uses WCI1 EnsureVisible and authoritative
offsets.

Multiline text defaults to word wrapping with a bounded vertical native Entry
scroll area inside the assigned shared-layout control rectangle. No additional
source wrap/axis syntax in WCI3. Single-line text retains native horizontal caret
reveal. The native Entry owns its internal editing viewport, caret/selection and
undo history; this is distinct from WCI1 outer frame/group/tree/list offsets.
Ordinary sync, page switches, collapse, failed Commit and failed reload retain that
native object and its internal edit state. Rejected native edit handling is subject
to the explicit bounded exception below. Entry scroll consumes wheel events inside its
text area (including at its limit); outer viewports remain reachable through their
own gutter/background. Do not route one gesture twice or promise remainder chaining
from this native editing surface. Typing/navigation reveals the caret internally;
focus navigation also uses WCI1 EnsureVisible for the control in its ancestors.
Resize/reflow clamps the native scroll without changing text or emitting Commit.
Successful compatible document reload may recreate Entry, reset caret/selection,
internal scroll and undo history while retaining compatible main/page accepted/draft
text and field focus, including inactive pages. Successfully closed dialog successors
reset unaccepted proposals to current accepted values, preserving explicit earlier
child commits, as required by WCI2;
failed reload retains all of them. No duplicate runtime text-scroll authority or
private-field reflection/unsafe access is required. Native evidence covers long
wrapped text, caret at both ends, wheel/scrollbar, outer-gutter access, page hiding,
resize, read-only copy and successful/failed reload policies.
 Read-only text remains focusable/selectable/copyable; paste/cut/typing are blocked. Use native
clipboard and Unicode editing; preserve content without normalization or splitting UTF-8 sequences.
Undo/redo changes drafts only; expose native Primary+Z and Primary+Shift+Z/redo bindings. Escape first cancels
active IME composition, then an editing draft; only an already clean field lets WCI2 surface Cancel handle
Escape. Preedit is not accepted text and never Commit; Enter consumed by IME must not dispatch SDL.

Successor retains compatible named same-kind main/page value/draft/selection/focus,
including inactive pages, revalidating new constraints. The WCI2 closed-dialog rule
overrides this general retention: successful reload resets unaccepted dialog proposals
to current accepted values, preserving explicit child commits; failed reload preserves
the original open form and all its proposals. For surviving main/page fields,
invalid retained drafts stay visible, not silently clamped. Retained accepted values violating new constraints
make reload fail, preserving the old bundle.
The already-required accepted-empty initialization exception above also applies to
extended input: unchanged required emptiness stays invalid/editable across reload.
Newly required blank accepted text rejects. For multiline-to-single-line conversion,
CR or LF in retained accepted text or a surviving main/page draft rejects; these
are the line-break characters checked by the bounded single-line policy. Discard a
closed dialog's unaccepted draft before checking that discarded text. This is not
a general exception for invalid accepted values or incompatible text constraints. Multiline/readOnly changes preserve valid text; single-line
conversion with newlines rejects. Keep native undo history across ordinary sync; explicit programmatic
replacement with different displayed bytes resets that field's history; identical
bytes preserve it regardless of explicit Apply or an echoing Commit. Successful reload may reset undo history, retaining content/draft;
failed reload preserves it. Removal/type change/Close revokes events and native callbacks. Extend existing
WCI2 PresentationGate/accepted-ticket publication for typed state and resources; do not invent another swap mechanism.

## 7. Capabilities, exports and acceptance

Keep WCI1's frontend/layout/widget/viewport/provider/host dimensions. Require widget AND host checkbox,
slider, select, number major 1; host input-multiline/read-only major 1 when used; provider choice-options/1
for select. Actual Binder signature/typed-plan preflight establishes SDL readiness, not a new bridge
dimension or readiness flag. No advertising unsupported host behavior. Composition/text
show kind, label, initial typed value/constraints and unbound/unsupplied state with original paths/spans.
Codegen reconstructs exact 0.3 models, numeric lexemes and metadata; no live options, providers or closures.
Public SVG may reject new controls/properties with `unsupported-value-export` or `unsupported-text-export`
and path/span before writing output. No snapshot renderer required; legacy default 0.2 exports stay unchanged.
Native background omission still requires actual prepared-control inventory; unknown kinds never fall through.

| Pending acceptance ID | Required proof |
| --- | --- |
| WCI3-A01 | Positive/negative schemas, exact .3/.2 rejection, typed/lexeme AST, reuse/spans, generated-constructor round trip and byte-stable default .2 outputs. |
| WCI3-A02 | Typed draft/accepted separation, validation feedback, read-only/programmatic mute, stale/duplicate events, StateGate failure and mixed invalid batch atomicity. |
| WCI3-A03 | Numeric invalid intermediate text, finite/range/step boundaries, fractional Go adapters, exact ±(2^53-1) SDL endpoints, fractional-near-integer lexemes, large-origin legal increments, exact decimal 0.3/0.1 membership, off-grid lexemes rounding to an on-grid Number, large quotient/indistinguishable-step rejection, and out-of-range rejection without coercion. |
| WCI3-A04 | Choice empty/disabled/repeated labels, replace/remove/recreate IDs, stale popup result and generation, no index fallback; native pointer/keyboard parity. |
| WCI3-A05 | Native checkbox/slider/number gestures and one Commit, canceled slider draft, disabled/read-only controls, focus/EnsureVisible across nested viewports and WCI2 surfaces. |
| WCI3-A06 | Real SDL four typed echo actions plus explicit Load/Save, bad signatures/result kinds, zero preflight calls, stale receiver/engine/options, failed persistence/delivery without automatic replay. |
| WCI3-A07 | Native multiline/newlines, wrapped text scrolling/caret reveal, outer-gutter routing, resize/page/reload retention policy, Unicode selection/cut/copy/paste, undo/redo, explicit Commit versus blur, invalid/required text, cancel without Save, read-only selection/copy. |
| WCI3-A08 | Actual OS IME composition/commit/cancel in native Entry, no preedit SDL calls, Unicode/newline clipboard round trip; keyboard injection or headless tests alone are insufficient. |
| WCI3-A09 | Successful/failed reload retains declared drafts/focus/history policy, tightened constraints, newer draft/result conflict, removed/disposed events; inspected text/composition and source-linked SVG rejection. |

Use controlled Xvfb/XTest, native screenshots and correlated event/domain logs plus an actual configured IME
for A08. If unavailable/unsupported, record A08 pending with the precise gap; do not claim IME delivery from
Unicode insertion. Run meaningful targeted tests, SDUI race suite and affected SDL/consumer regressions.
All rows remain pending; this draft closes no inventory family. Only this file is authorized for this turn.
Coordinator handoff: Architect 2.0.0 reloaded, shared SDP context reused, Session0010 S2/S4 recovered; next
step is WCI1/WCI2 pilot reconciliation and independent WCI3 design review, then separate implementation.

## Draft review disposition — Session0010 T003

Independent review required explicit scalar result mode, atomic-form proposal
policy, multiline edit-scroll ownership and bounded numeric precision. The revised
sections above resolve those design omissions without a generic form or mapping
framework. WCI3 owns OptionTarget; WCI2 PresentationGate and accepted-ticket
publication are the current foundation. Independent reviewer approved substantive candidate `497031bb1132865b037721620ebaeb864cffd07d1deee8443a3fcbd962a4fbf8`, including the bounded numeric scan. WCI3-M1 is now selected after WCI2 delivery and the reconciliation below; WCI3-M2 remains subsequent.
[IME boundary experiment](ime-probe/README.md) records a concrete pinned-stack
Submit-before-composition result and successful bounded filter trial. The reviewed
architecture selects the documented exact pinned, licensed GLFW source plus its
one-condition X11 patch for WCI3 implementation. Each maintained native build
root needs an explicit replacement; external consumers cannot inherit it through
a Go import. Product dependency edits await stage selection, and actual SDUI/SDL
IME acceptance remains mandatory.

## WCI3 preimplementation reconciliation — Session0010 T003

Coordinator and independent reviewer reconciled the concrete frontend/runtime API
memos with delivered WCI2 semantics. Both approved the closed-dialog successor
exception (including retained explicit child commits and inactive-page drafts) and
the bounded Go-only integer-grid regime outside exact safe53 constraints. These
are design clarifications within the existing inventory, not owner quotations or
new product evidence. M1 scalar fields and M2 extended native text remain separate.

## WCI3-M2 preimplementation reconciliation — Session0010 T003

Independent reviewer approved explicit-new-argument opt-in for extended input and
the TextResult self-receiver/exact-echo constraint above. These preserve old inputs
and prevent a new typed Commit from claiming source acceptance without it. Exact
output echo remains a post-execution check, not a false preflight guarantee.
EditField keeps its existing ModelRevision argument; draft/value/state guards belong
to captured FieldTarget. This records future design readiness only; M1 delivery
still precedes M2 product/dependency selection and actual text/IME evidence.

Coordinator and independent reviewer additionally approve the narrow already-required
accepted-empty reload exception. This refines the previous broad invalid-accepted
rule rather than retroactively labeling its literal implementation a defect. M1
select and M2 extended input share the rule; newly tightened constraints remain guarded.

## Native text editing refinement — reviewed before WCI3-M2 code

The coordinator inspected pinned Fyne 2.8.1 Entry: TypedRune merges private undo
state before OnChanged; SetText resets that history, and public Undo can include
earlier merged accepted keystrokes. Runtime snapshots also cannot distinguish an
identical-byte explicit Apply from exact Commit acceptance. These are observed
public-API constraints, not evidence of delivered text behavior.

Coordinator selection, independently reviewed: retain Entry and
history when displayed bytes are identical, including explicit programmatic Apply.
For different programmatic text, replace once muted and reset history. For an actual
native edit rejected by EditField or its publication gate, restore the latest
authoritative runtime Draft muted on the same focused Entry; this exceptional rollback may
reset native caret/selection/internal scroll and clears undo/redo history. Preserve
accepted text, runtime draft/revisions, existing publication and zero Commit. Use
the latest authoritative draft so callback reentrance cannot restore stale text.
This is a deliberate narrowing of the draft's earlier broad failed-publication
history promise, not an owner quotation or logically unavoidable product behavior.

Rejected alternatives are guessing callback provenance, invoking Undo on merged
actions, private reflection, a second editor/history store, or extending/forking
Fyne with a new editing-transaction API. They add ambiguity or maintenance beyond
the card's basic native editing scope. Normal valid edits, invalid-but-retained
validation drafts, self-echo Commit, failed Commit/reload/probe, page hiding and
ordinary synchronization retain their stated history behavior. Pre-delegation
refusal preserves history because Entry has not edited yet; do not apply the
exception to unrelated failures. Reject CR/LF single-line paste before Fyne's
normal LF-to-space conversion, preserving the no-normalization rule. Apply this to
keyboard and context-menu Paste, delegating exactly the validated clipboard snapshot
without a second content read.

Required evidence: real typed/pasted Unicode, ordinary undo/redo and retention,
identical and changed programmatic text, one forced rejected native edit with
authoritative text restoration and declared history reset, no extra Change/Commit,
and failed reload/Commit history retention. This resolves native adapter policy
without a runtime provenance API or an additional dependency patch. Independent reviewer approved this bounded refinement. Main retains Session/card
traceability; implementation still requires the recorded M2 stage selection.

WCI3-M2 is selected after independently reviewed WCI3-M1 delivery. The bounded
explicit opt-in/native editing contract and pinned X11 filter dependency are now
implementation authority. Earlier preparation-only statements remain historical;
no extended-text delivery or all-platform IME acceptance is inferred.
