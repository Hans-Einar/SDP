# SDUI runtime in Go

`runtime` owns UI state on one owner goroutine, normally the host UI goroutine.
It imports no layout or GUI package, retains no native pointers, performs no I/O,
and never invokes a collection provider. The host marshals background completions
back to that goroutine. SDL execution is connected through an explicit handler
installed by the [bridge](../../../SDL/go/bridge); runtime itself does not load SDL.

The current WCI1 implementation adds source-profile 0.3 collections and viewports
while retaining the 0.2 button/input APIs and normalized empty-profile encoding.
`New` checks the bounded, uniform normalized profile and clones the supplied frame.
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
sequence, collection states, viewport offsets and viewport handles. All maps,
item slices, request pointers and model provenance (`Uses`) are copied.

`CheckWith(func(*parser.Instance) error)` remains the legacy pure root gate.
`CheckStateWith(StateGate)` adds the complete state gate:

```go
type ViewportState struct { X, Y float64 }
type StateGate func(Snapshot) (map[string]ViewportState, error)
```

A gate receives prospective detached state before publication and returns **all**
clamped effective offsets keyed by normalized owner path. It must not change live
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
watermarks. Compatible named widgets retain handle, accepted value, draft and
focus. Matching collection kind/handle/provider ID/epoch retains data and local
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

Original 0.2 regression evidence remains in [G3](../evidence/G3.md). WCI1 runtime
unit/race tests cover typed data/events, cancellation, copied state, viewport guards
and successors. They do not establish native input, application-wide publication
or installed-consumer acceptance; those remain separate stage evidence.
