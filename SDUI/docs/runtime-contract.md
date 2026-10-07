# SDUI ↔ SDL — implemented runtime boundary

Updated 2026-10-08 for the WCI1 implementation candidate. SDUI 0.2 and 0.3 share
one Go runtime; SDL action-core 0.1 remains the bounded execution profile. WCI1
adds typed tree/list identity and collection/provider/viewport state without
changing 0.2 input values, button activation or the text result port. Native
integration acceptance is separate from the runtime package's unit/race evidence.

| Responsibility | Contract / implementation | Evidence boundary |
| --- | --- | --- |
| UI session, drafts, collections, requests, snapshots and successors | [Runtime API](../go/runtime/README.md) | Runtime tests; historical [G3](../go/evidence/G3.md) covers 0.2 |
| Source profiles and normalized models | [0.3 profile](profile-0.3.md), [language](language.md) | Parser/generator tests; 0.2 encoding retained |
| Collection interaction and atomic publication obligations | [WCI1 stage contract](../../SDP/04--Design/SDUI/Widgets/Collections.md) | Stage acceptance, including native input, is separate |
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
compatible widgets retain values/drafts/focus. Collections with the same compatible
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

WCI1 covers bounded tree/list data, navigation, lazy recovery and viewports. It makes
no claim for tabs/splits, command/transient families, typed numeric controls, multiline
input, a general collection query language, collection results in SDL or an installed
external consumer upgrade. Existing text/composition exports remain structural;
collection SVG export explicitly rejects unsupported interactive collection output.
The SVG widget remains a placeholder.

The bridge executes action-core, not design-core facts. EditAptCell remains an explicit
simulation, not a machine/Ponsse runtime, distributed transport or exactly-once
guarantee. There is no C ABI, FOX pointer or automatic SDL loader. The historical
SDUI-RUNTIME-001 proposal remains in Git before R2-M2; these current package/profile
contracts supersede its proposed payloads and methods.
