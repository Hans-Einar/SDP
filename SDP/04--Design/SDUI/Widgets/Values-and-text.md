# WCI3 values and text — provisional stage draft

| Field | Value |
| --- | --- |
| Assignment | SDP Architect, DRAFT ONLY; PLAN-SDP-0022 WCI3-M1/M2 |
| Status | Provisional pending WCI1 + WCI2 pilot and independent stage review; no implementation authorization |
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

Reuse Handle/model/sequence/draft revision checks, StateRevision and pure StateGate for every prospective
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
For fractional Go-only controls, step membership uses `q=(value-min)/step`, finite q and
`abs(q-round(q)) <= 8*2^-52*max(1,abs(q))`; validation never rewrites the supplied value. Pointer/step-key
movement deliberately chooses a legal step from min and clamps to the last legal step <= max.
Connected integer controls instead check exact integer arithmetic for bounds, range and step membership.

Supply `ChoiceOption{ID ItemID, Label string, Enabled bool}` slices by normalized instance path; at most
4096 options, unique nonempty IDs and label/ID limits reused from WCI1. Duplicate labels are allowed.
Initial explicit selection must exist and be enabled. Empty optional choice is valid; required empty choice
is an invalid editable state with Commit blocked. No parser I/O or automatic first-option selection.
Runtime owns an option-set generation: use WCI2's `OptionTarget{Handle, ModelRevision, OptionGeneration,
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

Extend the existing Plan.OutputField + explicit setHandle result path to these receiver kinds; no generic
multi-result mapper. Fixture actions AcceptFlag/AcceptCount/AcceptMode/AcceptText each take and return a
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
when atomic forms are intended; Cancel restores unaccepted drafts without Save. Earlier explicit child
domain commits cannot be undone by Cancel. Mixed valid/invalid Apply batches publish nothing.

Extend WCI2's captured `DraftField` with typed proposed Value, optional RawDraft (text/number),
FieldValidation and OptionTarget for select; invalid raw drafts never become accepted typed values.
WCI2's `DialogFieldValue` selector remains text-only until an explicit typed mapping is implemented.
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
offsets. Read-only text remains focusable/selectable/copyable; paste/cut/typing are blocked. Use native
clipboard and Unicode editing; preserve content without normalization or splitting UTF-8 sequences.
Undo/redo changes drafts only; expose native Primary+Z and Primary+Shift+Z/redo bindings. Escape first cancels
active IME composition, then an editing draft; only an already clean field lets WCI2 surface Cancel handle
Escape. Preedit is not accepted text and never Commit; Enter consumed by IME must not dispatch SDL.

Successor retains compatible named same-kind value/draft/selection/focus, revalidating new constraints;
invalid retained drafts stay visible, not silently clamped. Retained accepted values violating new constraints
make reload fail, preserving the old bundle. Multiline/readOnly changes preserve valid text; single-line
conversion with newlines rejects. Keep native undo history across ordinary sync; explicit programmatic
replacement resets that field's history. Successful reload may reset undo history, retaining content/draft;
failed reload preserves it. Removal/type change/Close revokes events and native callbacks. Extend existing
StateGate/atomic prepared publication for typed state and resources; do not invent another swap mechanism.

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
| WCI3-A03 | Numeric invalid intermediate text, finite/range/step boundaries, fractional Go adapters, exact ±(2^53-1) SDL endpoints, fractional-near-integer lexemes and out-of-range rejection without coercion. |
| WCI3-A04 | Choice empty/disabled/repeated labels, replace/remove/recreate IDs, stale popup result and generation, no index fallback; native pointer/keyboard parity. |
| WCI3-A05 | Native checkbox/slider/number gestures and one Commit, canceled slider draft, disabled/read-only controls, focus/EnsureVisible across nested viewports and WCI2 surfaces. |
| WCI3-A06 | Real SDL four typed echo actions plus explicit Load/Save, bad signatures/result kinds, zero preflight calls, stale receiver/engine/options, failed persistence/delivery without automatic replay. |
| WCI3-A07 | Native multiline/newlines, Unicode selection/cut/copy/paste, undo/redo, explicit Commit versus blur, invalid/required text, cancel without Save, read-only selection/copy. |
| WCI3-A08 | Actual OS IME composition/commit/cancel in native Entry, no preedit SDL calls, Unicode/newline clipboard round trip; keyboard injection or headless tests alone are insufficient. |
| WCI3-A09 | Successful/failed reload retains declared drafts/focus/history policy, tightened constraints, newer draft/result conflict, removed/disposed events; inspected text/composition and source-linked SVG rejection. |

Use controlled Xvfb/XTest, native screenshots and correlated event/domain logs plus an actual configured IME
for A08. If unavailable/unsupported, record A08 pending with the precise gap; do not claim IME delivery from
Unicode insertion. Run meaningful targeted tests, SDUI race suite and affected SDL/consumer regressions.
All rows remain pending; this draft closes no inventory family. Only this file is authorized for this turn.
Coordinator handoff: Architect 2.0.0 reloaded, shared SDP context reused, Session0010 S2/S4 recovered; next
step is WCI1/WCI2 pilot reconciliation and independent WCI3 design review, then separate implementation.
