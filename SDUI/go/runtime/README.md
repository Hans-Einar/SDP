# SDUI runtime in Go

`runtime` owns UI state on one owner goroutine, normally the host UI goroutine.
It imports no layout or GUI package, retains no native pointers, performs no I/O,
and never invokes a collection provider. The host marshals background completions
back to that goroutine. SDL execution is connected through an explicit handler
installed by the [bridge](../../../SDL/go/bridge); runtime itself does not load SDL.

The WCI2 implementation adds source-profile 0.3 tabs/pages/splits, shared commands,
button toggles, menus and dialog state to WCI1 collections and viewports, retaining
the 0.2 button/input APIs and normalized empty-profile encoding.
`New` checks the bounded, uniform normalized profile, clones the supplied frame
and strictly resolves the selected root with `parser.ResolveInteractions`.
Collection/provider and native readiness require the additional steps below;
constructing a Session alone is not connected admission. See the shared
[runtime boundary](../../docs/runtime-contract.md) and the governing
[collection stage contract](../../../SDP/04--Design/SDUI/Widgets/Collections.md).

## Existing widgets and session state

`Bind` installs an explicit Go handler. `Dispatch` validates handle/model revision,
nonzero increasing sequence, payload and enabled/visible state. Button `Activate`
has no payload; input `Commit` carries the current `Value` and `DraftRevision`.
`Draft`, `Revert`, `Focus` and `Apply` update UI state without simulating user events.
`Apply` validates the complete property batch before publishing. External accepted
values conflict with dirty drafts; explicit draft acceptance must match that draft.
`Close` revokes the Session and any pending collection requests.

Widget `Handle` contains session ID, public path, generation and kind. Named
same-kind widgets can retain that identity across compatible reload; anonymous,
removed or incompatible widgets cannot. `Widget.InstancePath` retains the exact
normalized path, including anonymous ancestors. These are different namespaces:
widget lookup uses the public path; provider, snapshot and viewport maps use the
exact normalized path.

`Revision` identifies the published model/binding epoch. `BatchRevision` orders
property batches. `StateRevision` advances for accepted UI-state changes, including
drafts, focus, collection state, offsets and consumed event sequences. Failed
prospective validation does not advance it. An unbound activation still consumes
its sequence, preventing later replay. `Sequence()` exposes the consumed watermark;
the owning application allocates the next event as `Sequence()+1` across successors.

## Collection data and provider admission

Use `BindProviders(map[string]CollectionProvider)` for the exact set of normalized
collection instance paths. Missing, unused and wrong-kind entries reject atomically.
A provider has a stable `ID`, an application-managed `Epoch`, copied `Initial` data,
`RootLoaded`, and optional `Load(context.Context, LoadRequest) (CollectionData, error)`.
Changing the provider source/function requires a changed identity or epoch.
Rebinding the same ID/epoch retains live state; a changed binding validates its
seed before replacing state and revoking the old request. Admission never calls Load.

`CollectionData.Items` is an ordered slice of `CollectionItem` records with `ID`,
`Parent`, `Kind`, `Label`, `HasChildren` and `ChildrenLoaded`. `ItemID` is opaque;
labels and row indexes never identify items. `Row` can be selected/activated and
may branch; `Group` may branch but is not selectable/activatable; `Separator` is
an empty nonfocusable leaf. Lists are flat; trees use explicit parent IDs.

Validation bounds each collection to 4096 items, depth 64, IDs to 1024 UTF-8 bytes,
and single-line labels to 32768 UTF-8 bytes. IDs must be nonempty and unique;
control characters, missing parents, cycles, invalid kinds/flags and children of
leaves or unloaded branches reject. Leaves have `ChildrenLoaded=true`; unloaded
branches and unloaded empty roots require a loader. Duplicate labels are allowed.

`Target(handle, itemID)` returns a `CollectionTarget` carrying the widget handle,
model revision, collection generation and item ID. Empty item ID means the root.
`ValidateCollectionTarget` checks this identity and item existence; interaction
operations additionally check active widget state, visible ancestry and item kind.
`Collection(handle)` and `Provider(handle)` return copied state/data, not writable
views of Session internals.

| Owner-goroutine operation | Effect |
| --- | --- |
| `ReplaceCollection(target, data)` | Replace the whole collection with a loaded root |
| `ReplaceChildren(target, data)` | Replace all descendants of a tree branch; root target replaces the entire collection |
| `SelectItem(target)` | Select a visible row without invoking a handler |
| `FocusItem(target)` | Focus a visible row/group or the collection root without selecting it |
| `ExpandItem(target)` | Expand a branch; begin one load for an unloaded non-error branch |
| `CollapseItem(target)` | Collapse a branch, cancel affected requests and move hidden row focus to the collapsed ancestor |
| `RetryItem(target)` | Explicit recovery of an error, canceled-unloaded target, or paused unloaded root |

Replacements validate the complete prospective collection and presentation before
publishing. Every accepted data replacement advances collection generation once,
even if data is identical, and revokes all prior targets and the active request.
Surviving compatible IDs retain selection/focus/expansion; deleted IDs and
row-to-group changes do not retain row selection. Unrelated subtree order is kept.
Local selection/focus/expansion/loading status does not advance data generation.

## Requests and typed events

`BeginLoad(target)` publishes loading state and returns a `LoadRequest` with target,
monotonic `RequestID` and `ProviderEpoch`. Only then may the host run Provider.Load.
There is one active request per collection; a different request supersedes it.
`CancelLoad(handle)` first revokes publication rights and records canceled status;
the host then cancels the corresponding context. Runtime neither runs goroutines
nor waits for uncooperative loaders.

Host adapters compare copied `CollectionState.Request` before/after accepted
operations to cancel removed flights and start newly allocated requests exactly
once. `CompleteLoad(request, data, err)` runs on the owner goroutine. It rejects
stale handles, revisions, generations, epochs, requests, hidden/disabled collections
and no-longer-expanded parents. Matching successful data is validated and published
atomically. Matching provider/data/geometry failure preserves data/selection and
records bounded error status; a nil return means either success or this recoverable
error state was accepted. Inspect `Status[parent]`. If even the error presentation
fails the pure gate, the operation returns that failure and preserves the last
valid state; the host reports the diagnostic outside the collection.

`Status` maps parent IDs (empty for root) to `Unloaded`, `Loading`, `Loaded`,
`LoadError` or `Canceled`, plus bounded UTF-8 error text. `AutoLoadPending` is a
one-shot flag for a fresh unloaded root. First load clears it; cancellation, focus,
visibility changes and compatible reload do not reset it. Error branches require
explicit Retry; collapsing/re-expanding does not retry. Re-expanding a canceled
unloaded branch may load once. A root left unloaded by compatible reload with
`AutoLoadPending=false` supports explicit Load/Retry, never an automatic restart.
Repeated Retry while already loading allocates no second request.

Collection events carry `Event.Collection *CollectionTarget`, zero Value and zero
DraftRevision; the target handle/model revision must match the envelope.
`Select`, `Expand`, `Collapse` and `Retry` are local events and consume sequence
without requiring a callback. Only `Activate` reaches a handler, and only for a
visible selectable row. After the handler returns, runtime revalidates the full
collection target before applying updates. Reentrant replacement, collapse,
hide/disable, reload or Close therefore cannot publish a stale action result.

## Snapshots and prospective presentation gates

`SnapshotRoot()` returns the detached presentation model, including current
widget properties. `Snapshot()` also contains model/state revisions, consumed
sequence, collection states, viewport offsets/handles, Tabs/Splits, Commands,
Presentations, Menus and Surfaces maps, ActiveSurface and
Focused public handle path. All maps,
item slices, request pointers and model provenance (`Uses`) are copied.

`CheckWith(func(*parser.Instance) error)` remains the legacy pure root gate.
`CheckStateWith(StateGate)` adds the complete state gate:

```go
type ViewportState struct { X, Y float64 }
type StateGate func(Snapshot) (map[string]ViewportState, error)
```

A state gate receives prospective detached state before publication and returns
clamped effective offsets keyed by normalized owner path. Inactive pane descendants
retain their stored offsets until they are measured again; ordinary legacy hidden
viewport omission retains its previous meaning. It must not change live
state, execute callbacks or invoke providers. Geometry and native resource owners
may prepare detached results, then publish them only after runtime accepts the
operation. Runtime validates returned finite/nonnegative offsets and scroll-owner
paths; layout supplies extents and clamping. `Snapshot` exposes no custom scene.

`Viewport(path)` returns a separate owner `Handle` with a `viewport:` kind prefix.
Frames/groups can own viewports without being widgets. Native callbacks use
`SetViewport(handle, modelRevision, offset)`, which rejects replaced owners and old
model revisions. `SetViewports(map[string]ViewportState)` is trusted bulk input;
it merges requested offsets and runs the same gate. Nonempty offset changes require
a state gate. Runtime is the single offset authority; the host mirrors accepted
offsets under a programmatic-sync guard. Geometry, clipping and scroll chaining
belong to [layout](../layout), not runtime.

## Detached successors and compatibility

`Successor(root, providers)` creates an unmounted candidate with the same logical
Session ID, next model revision and preserved consumed sequence/generation
watermarks. Compatible named main/page widgets retain handle, accepted value, draft and
focus. Dialog inputs retain accepted values but discard unaccepted drafts. Matching collection kind/handle/provider ID/epoch retains data and local
state, advances collection generation and removes all request tokens. Loading
becomes unloaded; stable errors/canceled/loaded states and AutoLoadPending persist.
Changed providers start from their validated seed.

Compatible named viewport owners retain offsets/identity, removed axes become
zero, and anonymous/removed/kind-changed owners discard viewport state. The new
state gate clamps retained offsets against the new geometry. The successor copies
neither handler closures nor old gates. Install new bindings and validate native
resources before publishing it. Capture the predecessor StateRevision and recheck
it at final publication; disposing a failed candidate must not close its predecessor.

`Reload(root)` remains the synchronous compatibility API, including compatible
legacy handler reuse. Connected WCI1 hosts prepare a detached Successor and publish
Session, bindings, resources and native controls together. Only after success do
they close the predecessor and cancel its provider contexts. Retained SDL engines
receive the continuing event sequence; creating a fresh New session instead would
reset it and is not the connected reload protocol.

## Tabs, pages and split state

`Pane(publicPath)` returns the exact owner/page handle; `Tabs(handle)` and
`Split(handle)` return copied state. Snapshot maps use normalized InstancePath.
Tabs contains ordered PageState records with stable direct IDs, page handles,
labels/icons and remembered focus. Page Enabled/Visible expresses direct declaration
intent; Tabs Enabled/Visible includes ancestor and containing-page activity. Thus an
unselected page can be eligible while its content is inactive. Snapshot.Root and
leaf Widgets expose derived effective activity for layout/input/provider admission.
`Widgets()` still enumerates legacy leaf controls; `CallbackOwners()` enumerates
tabs, canonical commands and dialogs with symbolic callbacks for discovery.

`SelectPage(handle,id)` changes selection silently. `Apply` supports pane
Enabled/Visible and tabs/page Label (nonempty UTF-8), with the existing 256-update
atomic bound. Disabling/hiding/removing the selected page picks the first eligible
page or an empty body. Hidden drafts, data and scroll offsets survive; request
tokens are revoked only when hiding successfully publishes. Reveal never resets
AutoLoadPending or revives a canceled request. Focus accepts tabs header and split
divider handles; `EnterPage(tabsHandle)` enters remembered valid or first content
focus. Invalid content focus returns to a surviving pane affordance. Collection
focus participates in page memory without changing selection semantics.

Split state stores axis, Proportion, declared minima, explicit Collapsed side
(`SplitNone`, `SplitFirst`, `SplitSecond`), Collapsible and SavedProportion.
`SetSplitProportion(handle,p)` rejects invalid or out-of-measured-range programmatic
values. `CollapseSplit(handle,side)` and `RestoreSplit(handle)` preserve the last
expanded proportion. User movement/resize/restore clamps against layout's measured
bounds. Collapsed hidden children preserve their state and offsets; both expanded
minima are suspended while the visible child's minimum still must fit. Failed
restore leaves collapsed state/focus unchanged. Disabled visible splits still need
geometry; disabled affects input rather than painting.

## Typed pane interactions and presentation tickets

`DispatchInteraction(Event)` handles ActivatePage and AdjustSplit plus the M2
command/dialog interactions described below.
Event carries exact Handle, ModelRevision, expected StateRevision, nonzero increasing
Sequence and exactly one kind-specific payload: panes use `Page *PageActivation` or
`Split *SplitChange`. Legacy
Value/Collection/DraftRevision are empty; legacy Dispatch rejects pane payloads.
PageActivation carries PreviousID, PageID and exact direct Page handle. SplitChange
Operation is ratio/collapse-first/collapse-second/restore; only ratio carries
Proportion. Same-page activation is a no-op. Native page activation keeps focus on
the tabs header; silent selection, fallback and reload never invoke the callback.

For pane interactions, `BindInteraction(handle, InteractionHandler)` installs one tabs callback;
`HasInteractionBinding` checks it. Omitted callback permits local navigation; a
declared unbound callback rejects activation. The handler returns
`InteractionReply{Updates, Domain}`; InteractionResult reports Sequence, Status
(committed/rejected/ui-conflict) and Domain (not-called/succeeded/rejected/unknown).
Domain annotations survive adapter errors; an unspecified/invalid called-handler
outcome is unknown. Consumed sequence survives callback failure; no automatic replay.
Pane events reject command/dialog payloads and non-nil reply Accept; M2 events
use their separate kind-specific path below.

Dispatch validates prospective geometry without preparing or publishing native
resources, then consumes one sequence before calling the synchronous handler once.
Widget identities and value/draft revisions plus the internal post-consumption
StateRevision are captured. Reentrant interaction dispatch rejects. Accepted legacy
mutations, including Reload, remain accepted and invalidate the outer reply rather
than being rolled back. Returned updates cannot replace reserved pane state.
Final local state and accepted updates publish atomically; the speculative gate's
clamps are discarded and final geometry is computed against the final candidate.

```go
type SplitGeometry struct { Lower, Upper, Effective float64 }
type PresentationState struct {
    Viewports map[string]ViewportState
    Splits map[string]SplitGeometry
}
type PresentationGate func(Snapshot) (PresentationState, error)
type PresentationTicket struct { Publish func(); Discard func() }
type PresentationPrepare func(Snapshot) (PresentationTicket, error)
```

`CheckPresentationWith` installs the pure typed geometry authority. Split-bearing
models require this gate; legacy CheckStateWith keeps its signature for other
models and replaces the typed gate, never adds a competing authority. Layout must
return exactly one finite legal bounds/effective-ratio entry for every visible
expanded split, including disabled ones; no collapsed/inactive/unknown entries.
Runtime validates the clamp, enforces strict programmatic bounds, retains inactive
pane offsets and accepts the effective state. No layout/GUI package is imported.

`PreparePresentationWith` installs a separate final resource preparer. It receives
only the finalized detached snapshot after geometry and runtime validation. Failed
preparation or changed live model/state calls that ticket's Discard. On success the
runtime swaps its state and synchronously calls the non-failing, non-reentrant
Publish, which only promotes prepared resources to accepted host pending. No live
mutation/domain callback belongs in Prepare; Publish performs no fallible work or
allocation. The host owns source/bundle/size checks and native synchronization.
Installing either hook publishes a validated current candidate, so failure retains
the old hooks/state. A speculative probe never calls Prepare or writes pending;
consuming Sequence/StateRevision alone never grants publication authority.
A separately accepted reentrant legacy change retains its own accepted ticket.

Successor retains compatible named pane identity, selection, hidden drafts/offsets,
remembered focus and split state. Removed/ineligible pages fall back silently;
axis/kind changes reset split state. Changing collapsible to false stages expansion
and must pass preparation. New pane generations continue the global watermark;
removed/recreated pages do not regain old handles. No old interaction handler,
presentation gate or presentation preparer is copied. Reload keeps compatible
handlers/hooks for its existing synchronous use. The governing
[pane stage contract](../../../SDP/04--Design/SDUI/Widgets/Panes-and-commands.md)
and host integration cover native acceptance separately.

Original 0.2 regression evidence remains in [G3](../evidence/G3.md). WCI1 runtime
unit/race tests cover typed data/events, cancellation, copied state, viewport guards
and successors. M1 tests add pane lifecycle, measured bounds, callback/reentrant
conflicts and ticket publication. They do not establish native input, application-wide publication
or installed-consumer acceptance; those remain separate stage evidence.

## Shared commands and captured menus

`Command(publicPath)` and `CommandState(handle)` expose the canonical owner.
Snapshot `Commands` holds its checked state, symbolic binding, context/effect,
resolved target and definition-instance exclusive scope. `Presentations` maps
actual button/item handles to one command and their effective label/icon/tooltip
and local restrictions. Frontend resolution supplies exact identities; runtime
never resolves source references by guessing lexical scope from public paths.

Basic 0.2 and existing 0.3 buttons remain `Activate`/`Bind` controls. Explicit
command/toggle/checked/exclusive/key/context/target/effect arguments opt in;
icon/tooltip alone do not. Promoted buttons remain in `Widgets()` with their
legacy Binding cleared. Canonical symbolic owners appear once in CallbackOwners.
Use `BindInteraction` for commands and dialogs; local effects reject handlers.

`CaptureCommand(origin, via, context)` returns a detached Event with the next
sequence, exact state/model identity and `CommandInvocation`. Via is button/menu/key;
a key uses the canonical owner. Capture consumes no sequence and publishes nothing.
Dispatch validates the complete capture, computes toggle/exclusive changes and
publishes them atomically with handler Updates. Selecting an already checked
exclusive command is a no-op. Apply supports Boolean `Checked` only for toggle
commands and validates the final exclusive group; callers include peer clearing
in their silent programmatic batch. Reply updates cannot overlap reserved toggles.

Button/key item context captures the current selected row. Menu item context
captures the clicked row without changing selection. Empty selection disables
item-context buttons. Context carries the widget/model identity and full WCI1 item
target; replacement, hidden rows or state changes cannot retarget an old action.
Modal routing blocks outside native events, drafts, focus and checked viewport
operations while preserving the underlying visible geometry/provider state.

`OpenMenu(menu, context)` publishes a root opening and returns `MenuScope`:
Handle, ModelRevision and captured StateRevision. `CloseMenu(scope)` compares the
stored opening, independently of the current global revision, so even stale menus
can be dismissed. Old cleanup cannot close a replacement. Dispatch separately
requires the current revision. Submenu navigation does not restamp runtime capture.
The host owns the synchronous native selection scope: hide/remove the old native
wrapper before Action, retain the runtime capture once, and close an unclaimed
opening on scope exit. Successful InvokeCommand stages root dismissal in its own
publication; any failure dismisses only that opening without action replay.
Escape/outside dismissal revokes immediately. Runtime stores no GUI scope object.

## Dialog opening, acceptance and lifetime

`Surface(publicPath)` and `SurfaceState(handle)` expose each declaration.
`OpenSurfaceFrom(dialog, opener, context)` publishes a new monotonically generated
`SurfaceTarget` and records the actual opener and parent surface. `OpenSurface`
uses current logical focus/root as a programmatic convenience. Reopening focuses
the same live token without clearing drafts. `FocusSurface(nil)` selects main;
a target selects its active canvas. Native OS focus remains a host responsibility.

`DialogFields(dialog)` returns owned input Widgets in source order, including
hidden pages and excluding nested dialogs. `DialogField(dialog, fieldPath)` uses
the frontend's strict named relative field resolver, including anonymous wrappers.
Both work for closed declarations so bridge admission executes no action.
`CaptureDialog(target, kind)` captures Accept/Cancel/Close. Accept includes exact
owned handles, raw String drafts and value/draft revisions, valid UTF-8 and at most
32768 bytes per field. External payloads must match exactly. Cancel/Close uses only
the opening token, never obsolete Accept field revisions.

Accept invokes only the dialog handler; accept/cancel/close command effects derive
that dialog event with the same sequence. A successful Accept requires
`Reply.Accept`; false requires rejected outcome and no Updates. Local unbound-free
acceptance uses committed/not-called. A declared but unbound callback rejects.
Captured fields plus Updates share the 256-write bound; captured-field overlap
rejects. Oversized known captures reject before execution. Oversized custom replies
reject after execution, retaining the truthful domain outcome without replay.

True acceptance commits drafts and Updates, closes the surface and queues one
result only after publication. False/error leaves it open. Succeeded with failed UI
publication, or unknown execution, records AcceptSequence/Domain/AcceptBlocked on
that opening without preparing a ticket. Further Accept rejects; edit/Revert and
token-based Cancel/Close remain available. No reconciliation or ResolveAccept API
exists. Reentrant legacy changes remain accepted and invalidate the outer reply.
A reentrant close retains the attempt on its queued result until the handler unwinds.

`CloseSurface(target, "cancel"|"close")` discards unaccepted drafts to current
accepted values. `ConfirmSurfacePublication(target)` is the host acknowledgment
following actual native publication; it changes neither geometry nor state revision.
`DrainDialogResults()` returns detached terminal results once, and holds unfinished
synchronous Accept receipts. A declaration cannot reopen before its queued result
is drained. Unmounted candidates never emit terminal results. Only Accept results
contain fields; automatic lifetime results have Sequence zero and retain the last
AcceptSequence/Domain. Reasons are user/programmatic/parent-closed/parent-hidden/
reload/dispose. Observers run after native synchronization, outside runtime mutation.

Parent hide/page/split transitions stage closure, draft discard, request revocation
and receipts on the candidate; rejected geometry/resources preserve the live state.
For an actual native owner loss, `RevokeSurfaces(parent, reason)` bypasses fallible
geometry. Zero parent means the entire Session. `RevokeSurface(target, reason)`
includes the exact lost opening itself (for native OnClosed bypass), using the same
lifecycle reasons and rejecting stale callbacks against a replacement opening.
Repeated exact-token cleanup emits no second result. Close performs dispose revocation.
For successful replacement, the host revokes the predecessor with reason reload,
closes it and drains receipts. Successor starts all surfaces/menus closed, carries
opening/handle/sequence watermarks and compatible checked state, resets unaccepted
dialog drafts, and retains ordinary main/page drafts. Failed successor preparation
changes none of the predecessor's surfaces, requests or focus.

M2 runtime tests cover these captures, outcomes, budgets, lifetime and ticket
contracts. Unit/race success does not establish native menu, window or OS keyboard
acceptance; those are separate host/integration evidence.

## WCI3-M1 typed scalar fields

M1 adds checkbox, slider, number and select within source profile 0.3. Existing
0.2 and basic 0.3 input Draft/Commit remain unchanged; extended input, multiline,
IME and undo are later work. Runtime imports the standard-library-only numeric
helper; it still imports no GUI/layout package or application file loader.

`Value` adds finite `Number` and stable `OptionID`, constructed with `Numeric` and
`Choice`. Only the selected payload may be nonzero. Widget.Value/Draft retain
legacy input strings; new controls never serialize their typed values there.
`Field(handle)` and `Snapshot.Fields` expose copied accepted/proposed values,
RawDraft for number/input, Dirty, ReadOnly, Required, bounded Validation, numeric
source constraints and options. The FieldTarget contains exact handle/model/state,
value/draft revisions and option generation. Inputs are a read projection only in
M1; they do not acquire the new editing policy.

`EditField(handle, modelRevision, value)` edits checkbox Boolean, slider Number or
number raw String. Invalid numeric intermediate text remains visible with validation
and no proposed Number. `EditTick` constructs one exact numeric grid point; invalid
number drafts cannot step. `ChooseOption` requires an exact OptionTarget. Accepted
edits publish through the existing gate/ticket, then call an optional ObserveChanges
observer once with copied FieldChange: Value is raw String for number, the typed
proposal for other controls, and Option carries the select token. Change never
calls the Commit binding. Initialization, Apply, revalidation and reload stay silent.

`CaptureCommit(fieldTarget)` creates a typed Commit envelope without consuming a
sequence. Automatic gestures use the Target returned in `change.Field`, so observer
reentrance cannot restamp an older edit. Explicit Enter/release captures a current
Field.Target. Event.Control includes captured ValueRevision, numeric RawDraft and
select OptionTarget. All identity, state and proposal fields must match exactly.
Existing Bind/Handler and Dispatch signatures remain; typed dispatch captures all
receivers before execution and rejects returned updates after reentrant mutation.
A successful typed callback must explicitly accept its source proposal with a
checked AcceptedValue update; otherwise the entire reply rejects. No domain replay.

A declared unbound callback rejects. With no source or Go Commit binding, valid
controls outside dialogs may accept locally. Automatic unbound Commit inside an
open dialog keeps the proposal dirty for owner Accept. Explicit child Commit effects
persist through later Cancel. ReadOnly blocks user edits/Commit while allowing
focus and checked programmatic changes. RevertField silently restores accepted
state; input delegates to the unchanged Revert operation.

Typed AcceptedValue, `ReadOnly` and `ValidationState` updates require the captured
ExpectedValueRevision and ExpectedDraftRevision; select additionally requires
ExpectedOptionGeneration. Every property compares those guards against the same
pre-batch baseline. Distinct property order cannot make already advanced candidate
revisions reject another property. Validation metadata stages last; mixed invalid
batches publish nothing. AcceptDraft must equal a valid captured proposal. Built-in
numeric/choice errors cannot be bypassed by clearing application feedback.
`ValidateFieldWith` installs an optional pure bounded validator; its code/message
are each at most 4096 UTF-8 bytes. Reentrant accepted work invalidates its outer
result instead of being overwritten.

`BindChoices(map[exactInstancePath][]ChoiceOption)` requires the exact select
inventory. Option generation zero means unsupplied; supplied-empty has a positive
generation. Sets contain at most 4096 unique nonempty WCI1-bounded IDs and plain
single-line labels; duplicate labels are valid. Initial explicit selection must
exist and be enabled. Required-empty remains invalid editable state. ReplaceChoices
checks owner/generation and the whole set, advances generation and invalidates old
popup results even if IDs reappear. Removed/disabled accepted IDs remain diagnostic;
the proposal clears and requires a new choice. There is no first-option/index
fallback, hidden collection or loader. Option/ValidateOptionTarget are dedicated
choice identity APIs, separate from collection targets.

`DialogFields`/`DialogField` stay text-input-only. `DialogControls` enumerates copied
input and typed fields for mixed Go Accept, including inactive pages and excluding
nested dialogs. CaptureDialog validates every proposal; DraftField now also carries
RawDraft, FieldValidation and OptionTarget. Mixed acceptance keeps the existing
256 combined writes, exact token/revisions, reserved targets, Domain/replay barriers,
published-only receipts and accepted-ticket protocol. Go InteractionHandler handles
mixed persistence; the SDL DialogFieldValue mapping stays text-only.

`SuccessorWithChoices(root, providers, choices)` supplies complete option data before
retained-value checks. The existing Successor signature can carry compatible supplied
sets. Neither copies observers/validators/handlers. The synchronous compatibility
Reload retains compatible Commit handlers, validators and Change observers, checks
validator reentrance and publishes them with the candidate. Compatible main/page drafts,
including invalid numeric intermediates, survive. Accepted values violating new
constraints reject the candidate. Successful closed-dialog successors reset dirty
proposals to current accepted values, retaining prior explicit child commits. Failed
reload preserves the live session. Option generations advance; old tokens cannot
move or accept replacement choices.

The numeric package preserves original decimal lexemes for exact source/text grid
membership and separately checks typed binary64 values by reconstructed-point round
trip. Gesture snapping is explicit; text and typed Apply never use a tolerance.
See [numeric API](../numeric/README.md) for bounded arithmetic and SDL conversion.
Runtime tests prove state and publication semantics; native gestures/SDL bindings
and whole-stage delivery require their separate integration evidence.
