# SDUI ↔ SDL — implemented runtime boundary

Updated 2026-10-08 for the WCI2-M2 implementation candidate. SDUI 0.2 and 0.3 share
one Go runtime; SDL action-core 0.1 remains the bounded execution profile. WCI1
adds typed tree/list identity and collection/provider/viewport state without
changing 0.2 input values, button activation or the text result port. Native
integration acceptance is separate from the runtime package's unit/race evidence.
M1 adds tabs/page/split state and typed page activation; M2 adds shared commands,
toggles, menus and modal/nonmodal dialog state. WCI3 typed-value families remain
outside this implementation.

| Responsibility | Contract / implementation | Evidence boundary |
| --- | --- | --- |
| UI session, drafts, panes, collections, requests, snapshots and successors | [Runtime API](../go/runtime/README.md) | Runtime tests; historical [G3](../go/evidence/G3.md) covers 0.2 |
| Source profiles and normalized models | [0.3 profile](profile-0.3.md), [language](language.md) | Parser/generator tests; 0.2 encoding retained |
| Collection interaction and atomic publication obligations | [WCI1 stage contract](../../SDP/04--Design/SDUI/Widgets/Collections.md) | Stage acceptance, including native input, is separate |
| Pane interaction and measured publication obligations | [WCI2 stage contract](../../SDP/04--Design/SDUI/Widgets/Panes-and-commands.md) | M1 runtime/bridge tests; native pane acceptance remains separate |
| SDL actions, records and Go registration | [action-core 0.1](../../SDL/docs/profiles/SDL-Executable-Action-Profile.md) | [G4](../../SDL/go/evidence/G4.md) |
| Field bindings and result receiver | [Bridge](../../SDL/go/bridge/bind.go), [value extraction](../../SDL/go/bridge/values.go) | Actual loaded engine/registration and bridge tests |
| Generated models using the same runtime | [Go generation](go-generation.md) | Generator equivalence tests; historical [G5](../../SDL/go/evidence/G5.md) |

## Connection and explicit field sources

The SDUI parser stores `ref`, callbacks and `setHandle` as data. The application
supplies module aliases mapped to already loaded SDL engines, actual Go function
registrations and a `bridge.Plan`. Neither parser nor bridge opens reference paths
or infers provider bindings. Action callbacks use `module.Action.@invoke`;
`module.Action.setHandle(page.preview)` identifies an explicit result input.
Names and symbolic callbacks alone do not establish connected readiness.

Each bridge input field must select exactly one source:

| `bridge.Source` member | Meaning |
| --- | --- |
| `Widget` | Current draft from an explicitly named input |
| `Event: true` | Existing input-commit event value; requires an input callback |
| `Literal` | Typed SDL scalar value |
| `Context` | Named application-supplied typed context value |
| `EventField: bridge.CollectionItemID` | Stable item ID from a validated tree/list Activate event, mapped to SDL text |
| `EventField: bridge.TabPageID` | New stable direct page ID from ActivatePage, mapped to SDL text |
| `EventField: bridge.TabPreviousPageID` | Previous stable direct page ID from ActivatePage, mapped to SDL text |

`CollectionItemID` has wire/debug spelling `collection.item-id`. It is the only
collection event-field selector in this boundary. It rejects button/input sources,
integer/boolean destination fields, unknown selectors and combinations with another
source. `Event: true` retains its 0.2 input-value meaning; it does not become a
generic event payload or accept a collection. IDs are not encoded in hidden inputs,
labels, row indexes or JSON text.

A bounded navigation binding can use:

```go
map[string]bridge.Plan{
    "nav.Activate": {
        Inputs: map[string]bridge.Source{
            "ItemId": {EventField: bridge.CollectionItemID},
        },
        OutputField: "Preview",
    },
}
```

Both collection callbacks may reference this same action with an explicit visible
input as the result receiver. The bridge still publishes a text output field to
that input through runtime Apply. It does not transport arbitrary collection data,
add collection results to action-core records or make the preview read-only.
Provider loading is a separate typed application responsibility.

Bridge preflight validates aliases, callback member, action/signature, exact input
fields, selected source types and text result receiver before installing handlers.
Actual SDL engine construction/preview validates the registered Go signature.
Preparation executes no domain action and no provider Load. Prototype local
interaction with supplied providers can work while symbolic SDL activation remains
unbound; it must not be reported as connected-ready.

## UI identity, local events and result delivery

Widget handles contain Session ID, public path, generation and kind. Provider,
collection-snapshot and viewport maps instead use exact normalized instance paths,
including anonymous ancestors. Reused definitions have distinct runtime instances.
SDUI owns drafts/accepted values, focus, collection state and revisions; registered
Go functions own domain state. Runtime stores no native pointers.

CollectionTarget adds ModelRevision, CollectionGeneration and stable ItemID to the
widget handle. Data is supplied as ordered CollectionData/CollectionItem records
with parent IDs and row/group/separator kinds. Runtime validates complete bounded
data before publication; programmatic whole/subtree replacements advance generation
and invalidate old targets even when IDs or data are identical. A deleted/recreated
item cannot regain an old token. Labels and list indexes never substitute for IDs.

Collection Event carries a CollectionTarget and zero ordinary Value/DraftRevision.
Local Select/Expand/Collapse/Retry events consume the increasing event sequence
without calling SDL. Only Activate on a current visible selectable row invokes the
collection handler. Groups/separators cannot activate. The original button Activate
and input Commit payloads retain their meanings; programmatic property changes do
not synthesize events.

The bridge captures the receiver's handle and value/draft revisions before executing
SDL. It checks that exact receiver and the bound SDL revision before returning an
update, avoiding retargeting a new input that happens to reuse the same path. Runtime
rechecks model/session and the full collection target after the handler returns,
then applies the property batch atomically. External values cannot silently overwrite
a dirty preview draft. Accepted domain execution is not rolled back by a later UI
conflict, and that failure does not trigger automatic action replay.

The WCI1 action fixture executes synchronously on the owner goroutine. Lazy provider
loading alone runs asynchronously under host control. This contract is not a general
asynchronous action executor or a responsiveness guarantee for arbitrary blocking
Go handlers.

## M1 pane interaction and publication

Runtime registers tabs/page/split composition owners separately from legacy
Widgets(). CallbackOwners() and BindInteraction expose callback-owning tabs to the
bridge; page headers do not acquire separate domain bindings. TabPageID and
TabPreviousPageID require a valid ActivatePage owner/payload and SDL text input.
The existing explicit text OutputField/setHandle receiver remains unchanged.
M2 selectors/result mode extend this boundary below; no implicit persistence is introduced.

DispatchInteraction validates exact model/state/handle/page identity and sequence.
ActivatePage carries previous/new stable IDs and the new page handle. Same-page
activation, programmatic SelectPage, removal fallback and reload execute no action.
Prospective geometry is checked before the synchronous handler. Selection/focus
and returned updates publish together or remain unchanged. Domain annotations
survive post-Execute adapter conflicts: ui-conflict/succeeded does not mean rollback.
A called handler with no outcome annotation is unknown. No automatic action replay.
Captured receiver value/draft revisions and the post-consumption runtime baseline
are checked before publication. Accepted legacy reentrant work remains accepted;
the stale outer result cannot overwrite it. Nested interaction dispatch rejects.

Snapshot adds detached Tabs/Splits maps and Focused. Runtime retains page intent
separately from content activity, preserves hidden drafts/data/offsets, cancels loads
only with accepted hiding, and never auto-restarts a canceled root on reveal.
Split state includes explicit collapse and saved expanded ratio. Layout supplies
measured minima/divider geometry and effective proportions. Strict programmatic
ratios reject out-of-bounds values; user/resize/restore paths clamp. Failed restore
or candidate preparation preserves prior state. Disabled visible panes still paint.

CheckPresentationWith installs a pure PresentationGate returning effective viewport
offsets and SplitGeometry bounds/ratios. Runtime validates the result, retains
inactive pane offsets and constructs the finalized state. PreparePresentationWith
then prepares a private PresentationTicket from that finalized detached snapshot.
Only the runtime's successful state swap permits ticket.Publish to promote native
pending; failure discards only that ticket. Publish is non-failing/non-reentrant.
Speculative geometry never prepares a ticket or changes pending; sequence/outcome
revision changes are not publication. Host synchronization may therefore retain an
independently accepted reentrant presentation after an outer callback error.

Compatible successors preserve named pane identities/selection/focus and relative
split state; axis/kind changes reset splits, and collapsible=false requires a valid
expanded candidate. Successors inherit no interaction handlers or presentation
hooks. Runtime imports no layout or GUI packages; all native resource/focus work
belongs to the host. See [runtime API](../go/runtime/README.md) for concrete calls.

## Provider requests and presentation state

The application supplies an exact path-to-CollectionProvider map. Each binding has
stable ID/epoch, initial data, RootLoaded and optional Load function. An unloaded
root or branch requires a loader. Runtime admission validates hidden collections
too, without loading them. Provider ID/epoch determines reload compatibility;
changing a function/source requires a changed epoch or identity.

BeginLoad publishes a token and loading state; only afterwards may the host invoke
Load. Runtime allows one request per collection, independently per reused instance.
The host owns goroutines/contexts and marshals CompleteLoad back to the UI goroutine.
Completion checks the live bundle in the host and exact Session/target/generation/
epoch/request/expansion in runtime. Supersession, collapse, hide/disable, refresh,
provider replacement, reload and disposal revoke publication rights before host
context cancellation. A provider that ignores cancellation cannot publish a late
success or failure.

Matching data and geometry are validated atomically. Failure retains prior data and
selection and becomes a bounded visible error with explicit Retry when the fallback
presentation can be prepared. A gate failure even for that fallback is returned to
the host's out-of-band status surface while the last valid presentation is retained.
Successful empty data is a loaded empty state, not an automatic retry loop.

A fresh unloaded root has one-shot AutoLoadPending. Cancellation, focus, visibility
and compatible reload do not reset it. An errored branch keeps its error through
collapse/re-expansion until explicit Retry. A canceled unloaded branch can load on
explicit re-expansion. A paused unloaded root with AutoLoadPending=false, including
one whose previous load was removed during reload, restarts only through explicit
Load/Retry. Repeated Retry while loading does not start another request or invoke SDL.

Snapshot includes detached Root, collection states, requested/effective viewport
offsets, viewport identities and model/state/sequence revisions. CheckStateWith
installs a pure prospective gate whose return value contains all clamped offsets.
Layout owns extents, ancestor transforms and clipping; runtime owns offsets. Native
scroll uses SetViewport with exact owner handle and model revision; trusted bulk
SetViewports supports coordinated offset changes. Old viewport handles cannot move
a replacement owner. The adapter mirrors accepted values under muted programmatic
sync, avoiding independent native and runtime scroll state.

## Reload, publication and compatibility

StateRevision covers every accepted UI-state mutation, including draft/focus edits
and consumed event sequence. It is independent of ModelRevision, collection generation
and property-batch order. A host preparing a replacement captures the predecessor's
StateRevision and must recheck it before publication; a source revision alone cannot
protect a user edit made during preparation.

Successor creates a detached Session with the same logical ID, next model revision,
monotonic generation/request watermarks and the consumed event sequence. Named
compatible main/page widgets retain values/drafts/focus. Dialog inputs retain
accepted values but discard unaccepted drafts; every successor surface starts closed. Collections with the same compatible
handle/kind/provider ID/epoch retain data, selection, focus, expansion and stable
status; collection generation advances and no request token survives. Loading
becomes unloaded without resetting the root auto flag. Other bindings use validated
initial data. Compatible named viewport offsets survive, removed axes become zero,
and the new geometry gate clamps offsets. Old handlers and gates are not copied.

The connected host binds and prepares that candidate, including actual geometry and
native resources, before publishing Session/handlers/controls as one owner-goroutine
bundle. Final guards include source/bundle/state/size/provider/SDL identity. Failure
before publication disposes only the candidate; the previous Session, requests and
resources remain live. After successful publication, the host revokes the old bundle,
closes its Session and cancels old loads. Retained application-owned SDL engines and
providers are not disposed merely because a view changes.

Event sequences continue through Successor rather than resetting after New. A
retained SDL Engine therefore sees monotonically increasing commands across both
collection callbacks and UI reloads. This fixture has one UI owner per engine;
sharing an engine between independent streams needs an application-owned allocator.
The legacy synchronous Reload and 0.2 CheckWith/SnapshotRoot APIs remain available.
Connected atomic native publication is a host responsibility, not a claim established
by calling live Reload followed by rebinding.

SDL reload validates registrations while preserving Go-owned domain state. Changed
Go code requires build/restart and does not automatically preserve memory state.
Generated constructors still use the same normalized model/runtime rather than a
second execution implementation.

## Bounds and earlier proposals

WCI1 covers bounded tree/list data, navigation, lazy recovery and viewports;
WCI2 adds tabs/pages/splits, typed page activation, commands/toggles, menus and
dialog acceptance/lifetime. This boundary makes no claim for typed numeric controls, multiline input,
a general collection query language, collection results in SDL or an installed
external consumer upgrade. Existing text/composition exports remain structural;
collection SVG export explicitly rejects unsupported interactive collection output.
The SVG widget remains a placeholder.

The bridge executes action-core, not design-core facts. EditAptCell remains an explicit
simulation, not a machine/Ponsse runtime, distributed transport or exactly-once
guarantee. There is no C ABI, FOX pointer or automatic SDL loader. The historical
SDUI-RUNTIME-001 proposal remains in Git before R2-M2; these current package/profile
contracts supersede its proposed payloads and methods.

## M2 command and dialog bridge

Basic legacy buttons keep Activate/Handler. Explicit behavioral M2 fields opt in
to a canonical command and InvokeCommand/InteractionHandler; icon/tooltip alone
keep legacy behavior. Shared button/menu/key presentations use one canonical
binding. Widgets retains real controls; CallbackOwners adds commands and dialogs,
with promoted buttons' legacy Binding cleared to avoid double installation.
The strict selected-root frontend resolver supplies command/target/dialog/scope
identities. Runtime DialogField/DialogFields permit closed-dialog preflight.

| Additional bridge selector | Captured source | SDL destination |
| --- | --- | --- |
| CommandContextItemID | CommandInvocation.Context.Item.ItemID | text |
| CommandChecked | Proposed CommandInvocation.Checked | boolean |
| DialogFieldValue plus FieldPath | Exact field in DialogRequest.Fields | text |

FieldPath is a named relative input owned by the callback dialog, including hidden
page fields and excluding nested dialogs. Accept extraction never reads live drafts
after capture. Existing CollectionItemID, TabPageID, TabPreviousPageID and legacy
Event sources retain their behavior. Unknown selectors and mixed sources reject.

Plan.ResultMode defaults to TextResult, preserving OutputField/setHandle delivery.
DialogAcceptResult instead requires a boolean AcceptField and text MessageField;
it forbids OutputField/setHandle/RevisionField/RevisionContext. It returns the typed
AcceptDecision and zero widget Updates. Admission validates the actual registered
Go signature and field paths without execution. Runtime checks returned decisions,
UTF-8/32768-byte message limits, capture identity and the combined 256-field/update
budget. A known excessive capture rejects before calling; unknowable oversized Go
reply Updates reject afterwards without erasing the executed domain outcome.

Reply.Domain is read even on error. Succeeded remains succeeded through later
engine/receiver/revision/geometry/resource conflicts; empty/malformed called
outcomes are unknown. False Accept is rejected/rejected, local true acceptance is
committed/not-called, called true acceptance requires succeeded. Succeeded UI
conflict and unknown execution block Accept for the same opening. Editing/Revert
and token-based Cancel/Close remain available; no automatic retry or ResolveAccept.
Cancel/Close never call the acceptance adapter or undo a child's prior Commit.

Snapshot adds detached Commands, Presentations, Menus, Surfaces and ActiveSurface.
MenuScope identifies an exact root opening; CloseMenu compares that stored opening,
not current StateRevision. InvokeCommand separately requires current state and the
original item target. The host retains capture only within its synchronous native
selection scope, removes the native wrapper before Action, and immediately revokes
Escape/outside dismissal. Success closes the root with command publication;
failure dismisses without replay/restamping and preserves Domain.

OpenSurfaceFrom records the actual opener/parent. SurfaceTarget includes handle,
model and monotonic opening generation. New surfaces pass the existing pure gate
and final ticket; sequence/outcome changes never promote pending. Host acknowledges
actual mounting with ConfirmSurfacePublication and drains DialogResult receipts
after synchronization/teardown. Unpublished candidates emit none. Reentrant closure
holds its receipt until synchronous Accept outcome finalization, and reopening waits
until the queued result drains. Only Accept results contain committed fields.

Prospective parent hide/collapse/page/reload stages child closure and provider
revocation until accepted publication. Actual native hide/close/dispose uses
RevokeSurfaces without a fallible gate. RevokeSurface targets one exact lost native
opening, including itself and children, with stale replacement protection; it is
the native OnClosed-bypass fallback. Successful replacement closes predecessor
openings with reload reason; successor surfaces start closed, retain compatible
checked/accepted values and discard unaccepted dialog drafts. Main/page drafts
retain existing behavior. Failed preparation preserves live openings and requests.
Runtime remains independent of layout/GUI; multiple canvases, actual modal/nonmodal
windows, focus and completion observers belong to the host and native evidence.

## WCI3-M1 typed scalar fields — development 0.3

Checkbox, slider, select and number use detached `Snapshot.Fields` keyed by exact
instance path. Number and OptionID are typed Value payloads; legacy input's string
Value/Draft remains its single authority. Invalid numeric raw drafts remain visible
with bounded validation feedback and cannot become accepted Number or SDL calls.
ReadOnly rejects user mutation/Commit but permits focus, copying and checked Apply.

EditField/EditTick validate the current ModelRevision; ChooseOption carries an exact
OptionTarget. Accepted edits publish before one copied Change notification. Automatic
Commit must use the returned FieldTarget, including its original state/value/draft
revisions and option generation. An observer's newer accepted edit invalidates that
capture; hosts may not recapture it to make the old gesture current. CaptureCommit
and existing Dispatch/Handler perform the checked typed Commit. A successful callback
must explicitly accept that exact source proposal; an unrelated update alone is not
an accepted Commit. No automatic domain replay follows rejected UI delivery.

AcceptedValue, ReadOnly and ValidationState compare required revision guards against
the same original batch baseline. Distinct-property order is irrelevant; invalid
mixed batches publish nothing. Validator output is bounded, detached and pure;
accepted reentrant work invalidates the outer result. Programmatic changes, initial
binding, reload and RevertField are silent. An unbound automatic Commit outside a
dialog may accept locally; inside an open dialog it keeps the proposal for owner
Accept. An explicitly bound child Commit persists through later Cancel.

BindChoices requires the exact select inventory and at most 4096 unique nonempty
stable IDs per field. Duplicate labels are valid. Unsupplied generation zero differs
from supplied-empty. ReplaceChoices advances generation, revokes old popup tokens
and retains removed/disabled accepted IDs for diagnosis while clearing the proposal.
There is no index-based substitution, hidden collection or implicit first choice.
SuccessorWithChoices validates supplied sets before retaining compatible values;
main/inactive-page drafts survive. A compatible already-required accepted-empty
select retains its invalid editable initialization baseline without revision advance
or invented selection; newly-required blank and nonempty ineligible accepted IDs
still reject. Closed-dialog unaccepted proposals reset to current
accepted values, and earlier child commits survive. Detached successors copy no
closures; compatibility Reload retains compatible handlers/validators/observers and
checks reentrance. Failed replacement preserves the live bundle.

DialogControls captures mixed typed fields under the existing exact opening, 256
combined-write, validation and domain-outcome protocol. DialogFields/DialogField
and SDL DialogFieldValue remain text-only. Explicit Go InteractionHandler owns mixed
form persistence; SDUI does not invent heterogeneous SDL record mappings.

SDL ScalarResult is closed typed Commit delivery with one explicit self setHandle,
OutputField and matching actual signature/result. ControlBoolean maps Boolean,
ChoiceOptionID maps the stable ID to text, and ControlNumber maps only checked safe53
integer grids/values. Original numeric source/raw lexemes are checked exactly before
conversion; fractional Go values never travel as SDL text. ControlText/TextResult
continues existing basic input behavior; the opt-in extended policy below is WCI3-M2.

The shared stdlib-only numeric package performs bounded exact decimal admission,
raw range/grid checks and typed binary64 reconstructed-point round trips. It exposes
immutable Grid operations, including exact Text(tick); only native gestures may snap.
See [numeric bounds](../go/numeric/README.md). This implementation description is
not native/whole-stage acceptance; final evidence is recorded by PLAN-SDP-0022.

## WCI3-M2 extended input — development 0.3

parser.InputOptions selects policy only through explicit new argument presence.
FieldState.Input is nonnil only for that source policy; Multiline and effective
Placeholder are defensively copied. Accepted/Proposed/RawDraft for input project
Widget.Value/Draft without a second text store. Required/ReadOnly/Validation reuse
the typed field contract. Extended text is valid UTF-8, at most 32768 bytes; required
checks TrimSpace without changing supplied bytes. Single-line line breaks mean
exactly CR or LF. Invalid bounded user drafts remain visible with blocked Commit,
while invalid programmatic acceptance rejects atomically.

EditField keeps its ModelRevision argument and returns the original FieldTarget.
Draft delegates once for extended input; legacy input remains unchanged. Revert
restores accepted text silently. Observers/validators and Apply use the established
copy/reentrance/original-batch guards. CaptureCommit carries String plus
Control.ValueRevision, with nil Control.RawDraft/Option: String already preserves
the exact text. Legacy nil-Control events cannot bypass the new source policy.

Every extended Commit with SDL TextResult requires an explicit self receiver at
preflight, regardless of the request selector, and exact returned echo after
Execute. A malformed/non-echo/newer-draft reply never accepts stale source text or
replays domain work. Legacy nonself text bindings and Load into read-only receivers
remain valid. DialogFields stays text-only; mixed DialogControls projects extended
input under the same 256-write/capture/outcome/replay rules.

Compatible main/inactive-page text drafts survive reload. Exact accepted empty
text under an unchanged required policy retains its invalid editable baseline;
newly required blank accepted values reject, and whitespace-only values are not
covered by the exact-empty exception. CR/LF in accepted or surviving draft text
rejects multiline-to-single-line conversion. Successful closed-dialog replacement
first discards unaccepted drafts to current accepted values, retaining prior child
commits. Failed replacement preserves the old form and all its proposals.

Native editing history is outside runtime. Identical displayed bytes preserve it;
changed programmatic text resets it. Actual post-mutation native edit rejection
restores current authoritative Draft muted on the same Entry, with the declared
history/caret/selection/scroll reset exception. Invalid-but-admitted drafts and
failed Commit/reload/probe retain history. Entry editing/clipboard/scroll/IME proof
is separate from these state guarantees and remains pending final M2 acceptance.
