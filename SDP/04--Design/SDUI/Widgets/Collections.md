# WCI1 collections and viewports — stage contract

| Field | Value |
| --- | --- |
| Assignment | SDP Architect; WCI1-M1 under PLAN-SDP-0022 |
| Authority | KB-SDUI-003; PLAN-SDP-0021; owner assignment of 2026-10-08 |
| Status | Independently reviewed for bounded WCI1 implementation; implementation and acceptance pending |
| Parent contracts | [Design](Design.md), [Acceptance](Acceptance.md), [Preparation](Preparation.md) |
| Plan | [Staged implementation](../../../05--Implementation/SDUI/Widgets/Plan.md) |
| Requirements | SDUI-R05, R12, R15–R18, R23, R26–R28; GAP-XFMD-SDUI-001/003/008 |
| Session | [SESSION-SDP-0010](../../../Sessions/session-%230010--SDUI_widgets.md), S2 |

## 1. Outcome, evidence and bounded decisions

Deliver a source-defined tree beside a grouped list and a visible text preview.
Supplied data contains stable folder/file identities, nonselectable headings and
separators, and lazy branches. Users can navigate, expand, select, activate, retry,
scroll and reload using keyboard or pointer. Only activation invokes a real SDL
action, with the selected item ID as text. Failed preparation retains the entire
previous UI. This is a complete WCI1 slice, not full XFMD navigation integration.

This elaborates the reviewed design rather than changing its full-inventory
promise. Tabs/splits and commands/transients remain WCI2; typed controls and
multiline input remain WCI3; all-family integration and matching consumer packages
remain WCI4. No new source sets, general collection query language, transaction
framework, virtualized/infinite list, drag-and-drop, multiselection, editable rows,
filesystem provider or decimal SDL transport is selected.

Observed original-workspace implementation:

| Source | Observation and consequence |
| --- | --- |
| `SDUI/go/parser/{parser,ast,validate,normalize}.go` | Exact 0.2 header; generic widget calls; only button/input/svg validated. Tagged JSON is emitted by `parser.Data`; normalization preserves declaration/use spans. Add a profile-specific schema, not another parser. |
| `SDUI/go/cmd/sdui/main.go` | Hardcoded `sdui-ast/0.2` envelope. Profile and format must be selected together. |
| `SDUI/go/runtime/{types,session,reload,properties}.go` | Widget handle, model revision and sequence exist; only button Activate and input Commit dispatch. Reload retains compatible named handles; root-only layout checking cannot validate collection state. |
| `SDUI/go/layout/{engine,contents,tracks,types,font}.go` | Rect/Clip already govern hit testing. Scroll rejects. Content and fixed/flexible tracks already have shared measurement. Extend that geometry rather than duplicating it in Fyne. |
| `SDUI/go/host/fynehost/{runtime,view,input}.go` | Mount creates controls; `Adopt` changes the session before mounting. Per-control ScrollNone wrappers clip controls. They are not collection viewport state or an atomic connected publication boundary. |
| `SDL/go/bridge/{bind,values}.go` | Rebind validates before installing handlers. Event sources require input widgets; output requires a text field and explicit input setHandle. Add one typed event-field selector. |
| `SDL/go/runtime/engine.go` | Engine construction checks Go registration signatures; Execute consumes a monotonic engine-wide sequence. Preserve command sequence across UI replacement. |
| `SDL/go/examples/application/reload.go` | Detached binding preview exists, followed by separate live model/rebind operations. WCI1 must not infer an atomic native swap from this preview. |
| `SDUI/go/preparation`, `host/fynehost/admission` | WCI0 detached admission is the prerequisite. Its root-only layout callback and hardcoded 0.2 admission need a state-aware 0.3 path. |
| `SDUI/go/presentation`, `go/svg`, `go/codegen` | Presentation has SVG fallback defaults; SVG may skip unknown widgets; codegen emits typed model constructors. Eliminate unknown-kind fallthrough for WCI1 consumers. |

Inspection began at `eae0149` with existing working changes. WCI0 is now delivered
in the original workspace at `b774907`; prerequisite generated-model provenance
was repaired at `b61a3de` (owner checkpoint and observed Git history). Use that
integrated prerequisite, not the earlier isolated worker snapshot. This Architect
inspection adds no implementation test claims; WCI1 acceptance remains pending.

Selected alternatives:

- Extend existing parser/runtime/layout/bridge packages. Encoding rows in Markdown,
  hidden inputs or JSON text loses typed identity and is rejected.
- Prefer mature Fyne tree/list controls with a bounded adapter. Prove stable IDs,
  nonselectable groups/separators, event semantics and nested scrolling against
  the contract. Shared layout owns outer geometry; the adapter may use native
  row geometry with one authoritative runtime offset model. Custom row painting
  is justified only by a demonstrated native-control limitation, not assumed in
  advance. Toolkit-internal row virtualization does not imply an infinite-data API.
- Permit one outstanding provider request per collection, independently for each
  reused instance. This bounds cancellation and ordering without a scheduler.
  A later request explicitly cancels the earlier one. Concurrent loading inside
  one collection is outside this stage, not a claim about future inventory.
- Collection SVG export may return an explicit unsupported diagnostic with source
  span and instance path. No collection SVG snapshot renderer is required.
  Text/composition remain explicitly static structural descriptions.

Section 5 defines observable interaction and event rules; section 6 defines the
permitted native adapter and geometry boundary.

## 2. Exact source profile and AST contract

### 2.1 Grammar and local validation

Keep `grammar/sdui-0.2.ebnf` and its fixtures unchanged. Add
`grammar/sdui-0.3.ebnf`: copy the 0.2 productions, replacing only the document
header terminal with `"0.3"`. Lexical rules, limits, definitions, reuse, formatting,
member-ref and setHandle productions are unchanged. The one parser recognizes
exact source spellings `0.2` and `0.3`; reject `0.30`, `3e-1`, `0.4` and omitted
headers. `Document.Profile` is respectively `sdui/0.2` or `sdui/0.3`.

The existing generic `widget-call` syntax remains. The following schema is the
complete new locally validated call surface, not a new expression grammar:

```text
tree(label: string, callback?: member-ref)
list(label: string, callback?: member-ref)
```

Here `?:` describes optional schema fields, not literal SDUI syntax. Examples:

```text
nodes=tree("Navigation", callback=nav.Activate.@invoke)
rows=list(label="Entries", callback=nav.Activate.@invoke)
```

The label is required, either one positional first argument or `label=...`.
Subsequent arguments must be named. Duplicate, extra, wrong-type, positional-after-
named and missing-label cases use existing argument diagnostics. Require a
nonempty, non-whitespace accessible label for tree/list. A callback requires a
named widget and a declared module alias, as today. The parser retains any
syntactically valid member reference; connected preflight requires `@invoke`.
Only Activate uses `callback`. Select, Expand, Collapse and Retry have runtime
semantics, with no additional source callback names in WCI1.

There is no `items=`, `provider=`, `onSelect=`, `onExpand=`, `onContext=`, inline
item array, callback lambda, new reserved keyword or collection literal. Providers
are supplied by the composition root. Anonymous collections are allowed for static
descriptions; a runtime provider may target their exact internal instance path,
but anonymous instances have no reload retention guarantee.

Validation dispatches by `Document.Profile`, including programmatically constructed
documents; unknown/empty profiles reject. `sdui/0.2` still rejects tree/list at
local validation. Its syntax-only mode may retain unknown calls exactly as before.
Profile 0.3 includes the existing three widget schemas unchanged and the two above.
Normalization must not consult a process-global schema without the document profile.

`overflow-x=scroll` and `overflow-y=scroll` reuse existing formatting. For the
0.3 execution profile their supported owners are frames, groups, trees and lists.
Direct scroll on button/input/svg/Markdown rejects during capability/layout
preflight; wrap such content in a supported container. Each scroll axis must have
a definite assigned size from fill/fr, scale or a resolved frame ratio, rather
than an auto/content-dependent viewport. Scroll with `justify=center/end/between`
on the same owner rejects with `scroll-layout`; WCI1 uses start-aligned scrolling
content. Other existing formatting and bounded layout failures remain applicable.
No pixel-valued source dimension or implicit scroll property is introduced.

### 2.2 Serialization and generated models

No source-AST struct or tag changes are necessary for these calls. Keep Document, Node,
Argument, Literal, Reference, Connection, Span, Instance, Region and UseSite field
names and representations, except the explicit normalized Instance addition below.
Collection nodes are `kind="widget"`,
`widget="tree"` or `"list"`. Labels remain tagged Literal strings; callbacks remain
tagged References. Normalized Arguments contain the same typed Go values, not
untyped decoded maps. Collection data/state is absent from source AST and Instance.

The canonical AST envelope has exactly the existing keys:

| Document profile | `astFormat` | `validation` | `document` |
| --- | --- | --- | --- |
| `sdui/0.2` | `sdui-ast/0.2` | Existing `syntax-only` or `local-profile` values | Existing `parser.Data(document)` |
| `sdui/0.3` | `sdui-ast/0.3` | Same validation vocabulary | Same tagged shape, with `profile: "sdui/0.3"` and expanded widget schema |

The format bump denotes the new valid profile/widget vocabulary, even though the
structural shape is unchanged. Keep 0.2 JSON output stable, including empty arrays,
null pointers, type tags and spans. Centralize envelope-format selection from the
validated profile; no independent caller-supplied format string. Any consumer of
an envelope must reject profile/format mismatch and unknown versions. There is no
new AST decoder or persisted runtime-state JSON schema in this milestone.

Add exactly `Profile string` with tag `json:"profile,omitempty"` to normalized
`parser.Instance`. Populate **every** instance of a 0.3 expansion with
`"sdui/0.3"`, including reused instances. For 0.2 leave it empty: empty means the
legacy 0.2 normalized representation, not an unknown document profile. A single
parser helper resolves effective Instance profile (`""` to `sdui/0.2`,
`"sdui/0.3"` to itself); other values reject. Admission/layout verify that every
descendant has the root's representation and reject mixed-profile trees. Prepared
Document.Profile must match that effective profile.

This deliberate legacy sentinel preserves existing normalized JSON and Go
constructors. Standard JSON omits empty Profile. `parser.Data` currently ignores
omitempty: add a narrow omission for empty Instance.Profile, rather than changing
its handling of all existing fields. The codegen literal writer similarly omits
only that newly added field for legacy instances. All existing declaration/use,
empty-array/null and type-tag behavior remains unchanged. Nonempty 0.3 Instance
Profile is emitted by both serializers and generated Root constructors; clone,
reuse, state snapshots and reload preserve it. Document remains explicitly
`sdui/0.2` or `sdui/0.3`; never serialize a 0.3 envelope/tag for 0.2 input.

Keep `codegen.Version = "sdui-go-model/1"` and output unchanged for 0.2 sources;
select generator identity `sdui-go-model/2` for 0.3. The 0.3 output still exports
`UISourceSHA256`, `Document()` and `Root()` with the same signatures. It additionally
exports `UIProfile = "sdui/0.3"` and `UIASTFormat = "sdui-ast/0.3"` constants.
Generated constructors contain data only; no providers, event closures or fetched
items. Normalize the generated Document and compare it with generated Root in
tests. Activation goes through document-based preparation, which already checks
root/document correspondence; calling Root alone never proves host readiness.

Bare-root admission/layout/host APIs use the effective Instance profile, replacing
the hardcoded `sdui/0.2` in admission.Check. They must not infer profile from widget
spelling. Legacy empty-profile trees retain 0.2 behavior, including rejected scroll.
A 0.3 root without required collection state/providers fails explicitly; use the
state-aware preparation/layout path for it. State-aware layout verifies profile
against Instance, rather than accepting unrelated caller assertions. Connected
activation still needs Document-based preparation for symbols and connections.

## 3. Capabilities and provider binding

Capabilities remain exact `(Dimension, ID, Major)` tuples. Slash notation below
shows ID and major separately where necessary; source-profile text itself contains
a slash. Extend the WCI0 checker without changing the 0.2 advertised set.

| Required tuple | Trigger / proof |
| --- | --- |
| Frontend, `sdui/0.3`, 1 | Every 0.3 document; all definitions locally validated |
| Layout, `relative`, 1 | Existing relative composition rules retained |
| Layout, `collections`, 1 | Shared collection row/content measurement |
| Widget, `tree` or `list`, 1 | Corresponding normalized collection instance |
| Provider, `collection-data`, 1 | Every live collection has an actual typed binding and validated initial snapshot |
| Provider, `collection-load`, 1 | A live binding declares lazy loading; its Load function is present |
| Host, `tree` or `list`, 1 | Native target implements the interaction contract below |
| Viewport, `scroll-x` / `scroll-y`, 1 | Each requested scroll axis, plus supported geometry and runtime offset handling |
| Host, `viewport`, 1 | Native target with any requested scroll axis |
| Host, `atomic-publication`, 1 | WCI1 native activation/reload path |

The existing button/input/Markdown/SVG-placeholder tuples remain separately
required when present. There is no blanket `all-widgets` capability. Check every
node of the selected expanded entry, including hidden descendants and reused
instances. Unselected definitions are syntax/local-validation inputs, not mounted
provider instances. Missing support identifies dimension, exact ID/major, instance
path, declaration span and use-site chain. Keep runtime data errors associated
with that widget path and offending item ID, not invented item source spans.

Flags alone do not satisfy provider admission. The composition root supplies an
exact instance-path map whose values contain a stable binding ID, binding epoch,
initial data and optional loader. Validate missing, duplicate, unused and
wrong-kind bindings. An eagerly populated collection may have no loader, but
cannot contain an unloaded branch. Hidden lazy collections are checked but do not
load until visible and requested. A supplied binding's ID/epoch controls reload
compatibility; replacing its function/data source requires a new epoch.

Native prototype mode may use explicitly supplied fixture providers and run local
selection/expansion/loading. It reports unbound SDL callbacks and refuses domain
activation with `unbound`; it never upgrades itself to connected readiness.
There is no automatic filesystem or demo-data provider. Static description export
does not require a provider and makes no native-readiness claim. Connected mode
requires the real bridge adapter even when no callbacks occur, as in WCI0.

## 4. Typed collection data, mutation and loading

### 4.1 Runtime data

Use concrete exported Go types in `SDUI/go/runtime`; labels never encode identity.
The following fields and meanings are the stage contract (ordinary Go naming and
file splitting may follow repository conventions):

```go
type ItemID string
type ItemKind string // "row", "group", "separator"
type CollectionItem struct {
    ID ItemID
    Parent ItemID       // "" is the root, never an item's identity
    Kind ItemKind
    Label string
    HasChildren bool
    ChildrenLoaded bool
}
type CollectionData struct { Items []CollectionItem }
type CollectionTarget struct {
    Handle Handle
    ModelRevision uint64
    CollectionGeneration uint64
    ItemID ItemID
}
type LoadRequest struct {
    Target CollectionTarget // ItemID is parent; empty means root request
    RequestID uint64
    ProviderEpoch uint64
}
type CollectionProvider struct {
    ID string
    Epoch uint64
    Initial CollectionData
    RootLoaded bool
    Load func(context.Context, LoadRequest) (CollectionData, error)
}
```

These are proposed interfaces, not existing declarations. Runtime owns copied item
records, generation, selected ID, focused row ID, expanded IDs and request status.
The application supplies ordered data and I/O. Parser, layout and generic controls
never open paths or call a domain action to populate a collection. The host owns
loader goroutines/contexts and marshals completions to the owner UI goroutine.
Provider functions receive copied requests and return owned/copied data; they
cannot access or mutate a live Session.

Add a Session `StateRevision` counter for preparation's state guard. Advance it
on every accepted UI-state mutation (including draft/focus, collection data/local
state, loading state, offsets and consumed event sequence), independently of the
source ModelRevision and collection data generation. Failed validation does not
advance it. The detached candidate records the old counter; final admission must
still observe it unchanged. This closes the gap left by checking only source and
batch revisions while a user edits a draft.

Rules for a validated snapshot:

- At most 4096 items per collection after a mutation, maximum parent-chain depth
  64 (root rows are depth 1), maximum 32768 UTF-8 bytes per label. Retain existing
  source/node and layout-operation limits independently.
- IDs are opaque nonempty valid UTF-8 strings, at most 1024 bytes, with no control
  characters. Compare bytes without case folding, path interpretation or Unicode
  normalization. Duplicate IDs reject even under different parents. Duplicate
  labels are valid. Escape IDs/labels for text, SVG and diagnostics.
- Every nonempty parent exists; cycles reject. Slice order determines sibling
  order by filtering on Parent; flatten the tree by ordered depth-first traversal
  of expanded branches. No runtime sorting, paging protocol or index-based IDs.
- `row` is selectable/activatable; it may also be a branch. `group` is never
  selectable/activatable but may be expanded. `separator` has empty label, no
  children and is never focusable/selectable/activatable. Rows/groups require a
  nonempty label. All labels are single-line plain text; reject control characters
  rather than interpreting markup or line breaks.
- A leaf has `HasChildren=false, ChildrenLoaded=true` and no child records.
  An unloaded branch has `HasChildren=true, ChildrenLoaded=false` and no supplied
  descendants. A loaded branch may contain zero children; no automatic retry loop
  follows a successful empty result. Separator flags must describe a leaf.
- Lists are flat: every Parent is empty, every item is a leaf. Groups/separators
  delimit sections by their position, without inventing another hierarchy.
- An unloaded root has no items. It needs a loader; a loaded empty root is a real
  empty state. Initial data must pass all rules before native publication.

### 4.2 Bounded mutation API

Provide two owner-goroutine operations, both requiring a current collection
handle, model revision and expected collection generation:

1. `ReplaceCollection(target, data)` replaces the entire snapshot, with loaded
   root. This is explicit application refresh.
2. `ReplaceChildren(target, data)` replaces all descendants of target.ItemID
   (empty ID means root). Returned top-level records name that parent; other
   records must be descendants of those records. The parent must be a current
   tree branch. Root replacement works for both tree and list.

This subtree replacement is the bounded incremental API: unrelated branches are
unchanged. There is no patch algebra or separate insert/move/delete protocol.
The caller supplies the intended sibling order. The candidate whole snapshot
must satisfy all limits and layout before any data/selection/geometry is published.
Duplicate IDs against retained branches reject. Invalid input leaves data,
generation, offsets, selection and outstanding requests untouched.

Every accepted data replacement advances CollectionGeneration exactly once,
including identical data. It invalidates every old collection target/request and
cancels the collection's outstanding load. Retain selection/row focus/expanded
IDs only when the same ID and compatible kind survive; row-to-group clears
selection. Deletion clears selection, never selecting whatever inherited its
index. IDs removed and later recreated cannot revive an old token. Pure selection,
focus, expansion or load-status changes do not advance the data generation.
Programmatic replacements and selection setters emit no user callback.

### 4.3 Requests, cancellation and errors

Runtime exposes begin/complete/cancel operations; the adapter executes Load only
after Begin has accepted and published loading state. There is at most one active
request per collection. RequestID increases monotonically for its lifetime and
does not reset on data replacement. Provider epoch is independent of UI revision.

Retain explicit per-parent/root error or canceled status independently of whether
the branch is expanded. The root also has a one-shot `AutoLoadPending` flag: true
only for a fresh binding's unloaded root, cleared on the first request, and never
reset by cancellation, visibility changes or compatible reload. Thus an initially
visible unloaded root loads once; an unloaded root after cancellation or a
compatible reload of an outstanding root request stays paused.
Show a visible `Load — R` control in that root state and `Retry — R` for an error.
An errored or canceled-unloaded branch shows the corresponding recovery control
with its branch status when expanded. These are existing runtime Retry operations,
not new source syntax, item IDs or SDL callbacks.

| Trigger | State transition |
| --- | --- |
| Explicit expansion of unloaded branch without an error, or initial visible unloaded root with AutoLoadPending | Mark expanded where applicable; allocate request token, clear root AutoLoadPending if applicable, show loading status, then invoke provider off UI goroutine |
| Expand a loaded branch | Show existing children; no request and no action |
| New request for another branch/root | Revoke old token and cancel its context; old branch becomes canceled/unloaded, not an error; show the new request's loading status |
| Collapse a loading branch or its ancestor | Revoke/cancel its request; preserve existing loaded data and record canceled status; next explicit expansion may load once with a fresh token |
| Escape while collection has a load | Cancel request; retain loaded data and selection, record canceled status and show Load recovery; do not immediately auto-reload or invoke a domain action |
| Matching successful completion | Validate prospective subtree and geometry, replace atomically, mark loaded, consume request; empty success is empty content |
| Matching provider/data/layout failure | Consume request; preserve prior data, selection and offsets; show escaped error and an explicit Retry affordance at that parent/root |
| Retry via visible recovery control or R while collection focused | Target errored or canceled-unloaded parent/root, or an unloaded root with RootLoaded=false and AutoLoadPending=false after compatible reload; require loader, expand branch if needed, allocate fresh RequestID, clear error/canceled status, show loading and run Load once. A second Retry while that target is loading does nothing and starts no request. |
| Collapse/re-expand an errored branch | Preserve its error and show Retry again; no provider request or SDL action until explicit Retry |
| Stale/canceled completion | Discard without data, status, selection, geometry or domain changes |

The acceptance predicate for a completion is the conjunction of live published
bundle, open session, exact handle, model revision, collection generation,
provider epoch, active RequestID, surviving parent and still-requested expansion
(or root request). Check again on UI-goroutine delivery, even if cancellation was
ignored by the provider. Cancellation revokes acceptance before calling cancel.
Hide/disable, provider replacement, collection removal, reload or disposal also
revokes outstanding requests; restoring visibility does not resurrect tokens.

Copy and bound completion data before retaining it; reject oversized results.
Capture loader panic as a load failure in the adapter. Error display is bounded
to 4096 UTF-8 bytes without cutting a code point; detailed diagnostics may be
reported separately. A late failure cannot overwrite a newer success. Cleanup
never waits for an uncooperative provider on the UI goroutine; canceled goroutines
are application resource responsibilities, while their publication rights are
already revoked. No claim of forced cancellation of arbitrary Go code is made.

Keyboard R (either case) targets the focused branch when it has recoverable
error/canceled status; otherwise, if no row is focusable, it targets the root's
recovery state. Pointer recovery always targets the displayed branch/root control.
An empty-root collection remains a focusable Tab stop, so Load/Retry is reachable
without rows. A selectable branch's Enter retains Activate semantics, while R
performs only recovery. Neither retry nor restarting a canceled load changes
selection or invokes SDL. Visible recovery labels document the keyboard shortcut.

## 5. Typed events, focus and native interaction

Extend the existing Event with `Collection *CollectionTarget`. Keep Value and
DraftRevision semantics unchanged for 0.2. Collection events have zero Value and
zero DraftRevision; the target's Handle/ModelRevision must equal the envelope.
Add event kinds `select`, `expand`, `collapse`, `retry`; reuse `activate` with a
discriminated collection target. A button Activate still requires no payload.

Before consuming an event or running callbacks, validate handle/kind, revision,
nonzero increasing Session sequence, enabled/visible ancestry, collection
generation, item existence, kind and current expanded-row visibility. Retry alone
may target empty ItemID for root error, canceled-unloaded recovery, or a paused
unloaded root with RootLoaded=false and AutoLoadPending=false. Select/Activate
require a row; Expand/Collapse require a branch; Retry requires the matching
error/canceled-unloaded state or that root-only paused state, and a loader. Repeated Retry for a currently loading
target is a no-op after identity/sequence checks, never a second request.
Reject extra/mismatched payloads. A focused group is not a selected item. A
collection has at most one selected row; activation does not depend on a display
label and never silently substitutes the current index.

Successful local Select/Expand/Collapse/Retry events consume the sequence without
requiring an SDL callback. Only Activate reaches the collection callback handler.
Activation in a source with no callback reports `unbound`; local navigation remains
usable. Repeated events and old native closures cannot execute twice. Revalidate
the full collection target after a handler returns, before Apply, in addition to
the existing model-revision check. A synchronous handler that caused a data
replacement invalidates its result too. Preserve existing atomic update/draft
conflict checks on the preview receiver.

Host interaction is fixed for this milestone:

| Input | Behavior |
| --- | --- |
| Tab / Shift-Tab | Enter/leave a collection as one focus stop in normalized reading order: header, body/rows left-to-right, footer; skip disabled/hidden controls |
| Enter collection | Restore surviving visible focused row, otherwise first navigable row; do not auto-select merely on focus |
| Up / Down, Home / End | Move focused row among visible rows/groups, skipping separators; select when the destination is a selectable row; focusing a group leaves existing selection unchanged |
| Left in tree | Collapse focused expanded branch; otherwise move to visible parent; selection follows only a selectable destination |
| Right in tree | Expand focused collapsed branch; otherwise move to first navigable loaded child; no action invocation |
| Space | Select a focused selectable row; no activation; group does nothing |
| Enter | Activate focused selectable row once; on focused group error Retry, otherwise toggle that group's expansion; empty-root error retries |
| Primary single click on row body | Focus and select a row; focus group without selection; never activate |
| Primary double click on row body | Activate a selectable row once after selection; toggle a group; never confuse a disclosure/retry hit with activation |
| Disclosure click | Focus branch and expand/collapse; does not change selection or activate |
| Recovery control click / R while collection focused | Load/Retry the focused errored or canceled-unloaded branch, or the recoverable root when there are no focusable rows; one provider request, no selection or activation. The visible control names the R shortcut. |
| Wheel / touchpad | Scroll by shared logical delta at nearest eligible viewport under pointer; selection unchanged |
| PageUp / PageDown | Scroll nearest vertical viewport by 90% of visible height; preserve selection and row focus |
| Alt + arrow keys | Scroll focused widget's nearest eligible viewport on the indicated axis by one row height (or 40 logical units outside a collection) |
| Escape | Cancel active load as above; otherwise no action |

Collapsing/re-expanding an errored branch retains the error and reveals Retry;
it does not retry automatically. Explicit re-expansion of a canceled-unloaded
branch starts one fresh request. A canceled root remains paused until Load click
or R; focus changes, repaint and compatible reload never restart it. R is handled
only with collection focus and does not intercept text input in another widget.

List Left/Right have no item operation. Horizontal wheel gestures and Alt-Left/
Right provide horizontal scroll access. Text-entry key handling is retained when
an input has focus; host scrolling must not steal editing keys. Focus movement to
an offscreen row/control runs EnsureVisible. Collapsing an ancestor of focused
row moves row focus to that ancestor; a selected descendant remains selected but
cannot activate while hidden. Deleting a focused item clears its row focus and
selection if selected; next navigation chooses a surviving row without synthesizing
an activation. Programmatic state sync never simulates native user events.

Prefer Fyne collection controls with a bounded adapter. Use custom BaseWidget
behavior only where a demonstrated native limitation prevents the required
semantics. Existing button/Input controls remain native Fyne controls. The adapter
keeps rendering and hit testing consistent with shared outer geometry and the
authoritative runtime offsets; clipped or disabled rows never dispatch. Build
native closures with bundle identity, handle, model revision and current collection
generation, refreshing them after data changes. A double click must produce at
most one Activate regardless of native tap notifications.

The visible title supplies the accessible label; focused/selected, expanded,
loading, empty and error states must be distinguishable without color alone.
Keyboard/pointer parity is required. Do not claim an OS screen-reader integration
or IME proof from BaseWidget interfaces or headless tests; record actual native
coverage and limitations. Full text-edit/IME acceptance remains WCI3.

## 6. Shared geometry and nested scrolling

### 6.1 Ownership and calculation

Shared layout owns outer widget/container rectangles, viewport extents, ancestor
transforms and clips. Runtime owns offsets and viewport identities. A bounded
Fyne adapter may use native tree/list row measurement, painting and hit testing;
it must map native row IDs/indexes to the current stable ItemID and report the
content extent and target bounds needed for clamping and EnsureVisible. Native
paint and hit testing must agree with the effective offsets and ancestor clips.
No exported full-scene API or shared row renderer is required.

Extend existing layout/preparation inputs only as needed for validated profile,
prospective collection measurements and requested offsets, keyed by normalized
instance path. Keep the existing Box tree for outer geometry. The layout/adapter
boundary needs the following facts, which may remain internal implementation
records rather than a new public scene model:

```text
owner instance path; nearest viewport ancestor; scroll axes;
viewport rect in screen coordinates; content size in local coordinates;
requested offset; clamped effective offset; maximum offset; effective clip
```

For each axis, `maxOffset = max(0, contentExtent - viewportExtent)` and
`effectiveOffset = clamp(requestedOffset, 0, maxOffset)`; a non-scroll axis has
offset zero. Reject NaN/infinity/negative requested offsets. An offset beyond the
new maximum is clamped, not rejected. Publish effective offsets with the geometry
that established their bounds. Runtime imports no layout or Fyne package.

Native scroll callbacks submit requested changes to that runtime authority; the
adapter mirrors accepted offsets into the control under a sync guard. Native
scrolling, selection and expansion must not maintain competing authoritative
state or emit duplicate user events during synchronization. Demonstrate this
boundary with the actual pinned Fyne controls. If a required behavior cannot be
adapted, document the specific limitation and use a bounded custom implementation
only for that gap; no generalized collection rendering framework is selected.

Preserve current relative-size references: the explicit parent's finite available
body, not an infinite scroll plane or a synthetic row, determines percentages,
padding and gaps. Allocate fixed/fill/fr tracks against that finite size once;
sum/union actual child extents without shrinking overflowing fixed content. On a
scroll axis, start positions stay nonnegative; child center/end alignment uses
`max(0, available-used)` rather than generating unreachable negative overflow.
For `error`, compare content extent with viewport and fail; `clip` clips with zero
offset; `scroll` permits positive overflow and reports extent. Extents remain
finite and within the existing `1e7` computed-size and operation budgets.

Scrolling a container moves its entire content, including its regions. Header/
footer retain their existing viewport-relative allocation, but are not sticky
inside a scrolling owner. To make fixed page header/footer with scrolling body,
put scroll on an explicit body child. For nested viewports, parent extent includes
the nested viewport's outer box, not its hidden content extent. Do not propagate
an inner collection's 4096-row content height into the outer page's extent.

The adapter measures title, rows, separators, indentation and status affordances
using the same native metrics used for rendering. No custom row-height constants
or renderer are mandated. Labels remain single-line; content extent includes
their measured width and indentation. Title text clips within its strip and does
not scroll. Empty/loading/error/retry affordances are UI targets, never provider
ItemIDs. Native index-based controls require a current index-to-ID mapping; row
reuse or reorder must not retain a handler for the previous item.

On a scroll axis, distinguish the native viewport minimum from full content size;
do not apply `native-minimum` to all loaded rows. Without scroll, ordinary content
sizing and explicit error/clip policies apply. A collapsed branch contributes no
visible descendant extent. Collection row clipping intersects its own viewport,
every ancestor clip and the root window clip. A native adapter is sufficient only
if these metrics, clipped hits and stable-target events are demonstrably coherent.

### 6.2 Transforms, routing and lifecycle

For a content-local point, screen position is its owner's content origin plus the
point minus that owner's effective offset, with each ancestor content transform
applied once. The owner's viewport itself moves only with ancestors, not its own
offset. Compose clips by intersection in screen space. Paint and hit tests use
the same transformed bounds. Native local hit testing is permitted inside that
clip; outer Box hit testing retains its half-open bounds. Assign each translation
to either layout or the native adapter and apply it once, not in both.

Wheel routing hit-tests the topmost visible branch and walks its viewport ancestry
from inner to outer. Per axis, the nearest scroll-enabled viewport consumes the
delta it can apply; only the unconsumed remainder reaches the parent. No scrolling
of a sibling or obscured viewport, no double application at a boundary. A clipped
child does not receive pointer events through its offscreen area. Reuse native
scrollbar thumbs where possible; their drags use the same runtime scroll operation.
Account for any native gutter in the measured viewport once. Thumbs never own a
competing offset. Provide visible thumbs whenever an axis has positive range.

EnsureVisible starts at the innermost scroll ancestor, adjusts minimally to reveal
the focused row/control, then recomputes its screen rect before adjusting the
next ancestor. If a target is larger than the viewport, align its leading edge.
No scroll update emits selection/activation. Resize, collapse, deletion, refresh
and offset changes validate prospective geometry/native state and publish the
corresponding clamped offsets. Failed geometry retains the last valid data and
presentation and reports a diagnostic. If the physical window becomes too small,
clip that retained presentation to the real
window and show the failure; do not publish a partly measured layout.

Viewport runtime handles are separate from widget handles because frames/groups
also scroll: session ID, exact normalized path and monotonic generation. Retain
offsets on compatible named owner paths/kinds during reload, zero removed axes,
clamp against new extent, and discard anonymous/removed/kind-changed viewport state.
Nested offsets are independent. A stale viewport handle cannot move a replacement.

## 7. Real SDL Activate fixture and bridge extension

Use the existing executable SDL action profile, unchanged:

```text
language action-core version 0.1.
action Activate.
record ActivateInput.
record ActivateOutput.
Activate invokes GoActivate.
Activate takes ActivateInput.
Activate returns ActivateOutput.
ActivateInput field ItemId as text.
ActivateOutput field Preview as text.
```

Register GoActivate with exactly `{ItemId: text} -> {Preview: text}` in a real
`sdl.Engine`. Its fixture implementation looks up an opaque ID in a supplied map
and returns `Preview: <fixture text>`; unknown ID fails without changing UI. Keep
a call log so tests distinguish provider requests, local navigation and domain
execution. No real file open or XFMD write is implied.

Use this complete SDUI source, with injected providers for `page/body/nav` and
`page/body/entries` (the latter is also the exact normalized path here):

```text
sdui 0.3;
ref: nav "collections.sdl";
page = [
  header="Collection navigation fixture";
  body=<
    nav=tree("Navigation", callback=nav.Activate.@invoke)
      {x=1fr, y=fill, overflow-x=scroll, overflow-y=scroll},
    entries=list("Entries", callback=nav.Activate.@invoke)
      {x=1fr, y=fill, overflow-x=scroll, overflow-y=scroll}
  > {x=fill, y=fill};
  footer=<preview=input("Preview", value="No activation") {x=fill}> {x=fill}
] {x=fill, y=fill};
nav.Activate.setHandle(page.footer.preview);
```

The group inside footer preserves existing region and public-path rules; the
preview is a named input at `page/footer/preview`.

Extend `bridge.Source` with `EventField`, a closed enum whose only new nonzero
value is `CollectionItemID` (wire/debug spelling `collection.item-id`). Exactly
one of Widget/Event/Literal/Context/EventField must be selected. `Event=true`
retains its existing input-value behavior and does not accept collections.
Preflight permits CollectionItemID only for tree/list Activate callbacks and an
SDL text input field; it rejects integer/boolean destinations, unknown selectors,
button/input sources, wrong member, missing action/module and invalid result
targets. Extraction uses the already validated `event.Collection.ItemID`.

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

Both collection callbacks may bind the same action and visible preview connection;
retain the existing one-setHandle-per-object/target rule. No hidden receiver or
generic serialization is added. The preview remains an ordinary visible input;
WCI1 does not invent the WCI3 read-only property. A user draft in that receiver
continues to cause a draft conflict instead of being overwritten. Tests must show
that accepted domain execution is not automatically retried if UI result delivery
then fails; cancellation is not domain rollback.

Validate actual registration via `sdl.New`/Preview, and signature/connection via
real bridge preparation on the detached session. Preparation executes zero
GoActivate calls and zero provider Load calls. Capture result target handle and
expected value/draft revisions before executing the action; do not retarget a
newly created input found later by the same path. Check current collection target,
bound SDL revision and destination before publishing its update.

The WCI1 fixture keeps synchronous SDL Execute on the owner goroutine and uses a
bounded in-memory handler. Lazy provider loading alone is asynchronous. No generic
async action executor is introduced. Long-running domain handlers remain an
application integration concern; this fixture does not establish responsive
native behavior for arbitrary blocking actions.

Session event sequences must be allocated by the owning application across bundle
swaps, never reset to 1 when a replacement Session starts. Carry the consumed
Session watermark to its detached successor and continue the host counter. A
retained SDL Engine sees that same increasing command sequence, including across
both collection callbacks. The fixture has one UI owner per engine; sharing an
engine among independent event streams requires a separate application-owned
sequence allocator and is not silently supported here.

## 8. Atomic native preparation, publication and reload

WCI0 remains the pure capability/binding foundation. Add one composition-root
publication owner, with a stable host shell pointing to one published bundle:
Document/source identity, Session, bindings/module revisions, provider bindings,
native view/resources, validated geometry and lifecycle token. Native callbacks always
check that exact bundle is current before dispatch. Do not call live Session.Reload
followed by mount/rebind for WCI1 connected reload.

The transaction is:

1. Read/parse/normalize bounded source off-thread where useful; keep the candidate
   sequence/hash, entry and module/provider epochs. No Fyne work on that thread.
2. On the UI goroutine, snapshot the old bundle's UI state. Prepare a detached
   successor Session with the same logical session ID, revision old+1 and
   non-reused generation counters. Retain named same-kind widget handles,
   accepted/draft values and compatible focus. Never copy old handler closures.
3. For a same-path/kind collection with unchanged provider ID/epoch, copy loaded
   data, selection, row focus and expansion, advance collection generation and
   remove all request tokens. Loading becomes unloaded; preserve stable error,
   canceled and loaded state, and the root's AutoLoadPending flag. Otherwise
   validate/initialize the new provider's initial data. Carry compatible viewport
   offsets, then clamp using new geometry. An expanded branch that was loading
   before successful reload may start a fresh request after publication; error or
   canceled recovery remains explicit and cannot be restarted by the reload. A
   root that was loading becomes paused unloaded with AutoLoadPending=false;
   expose Load/R and accept root Retry to allocate one fresh request.
4. Check the entire selected entry's profile/capabilities/providers and actual
   loaded SDL registration/signatures/connections. Bind only the detached Session.
   Its model revision and binding epoch must be fixed before closures are built.
5. Prepare actual Markdown/content resources, collection extents/clipping, native
   controls and their adapters, handlers and focus target for the current window size.
   Exercise any fallible measure/render/resource construction now, including hidden
   supported controls that may later appear. Do not mount, start Load or execute
   actions. A zero/uninitialized native size is not proof of layout readiness:
   require a real intended finite size, then revalidate the actual mount size.
6. Recheck latest source candidate sequence/hash, old published-bundle identity,
   old UI-state revision, window size, provider epochs and SDL revisions. Any
   mismatch abandons or rebuilds the candidate; never copy stale drafts back over
   edits made during preparation. Single synchronous owner-goroutine preparation
   is the default; any deferred work must retain these guards.
7. Commit by switching the prepared bundle/host child as one non-yielding operation
   on the UI goroutine, with native sync muted. The final operation performs no
   provider/domain callback, fallible reflow, binding validation, decoding or
   application resource allocation. Native layout/refresh uses the prepared controls
   and validated bounds; ordinary toolkit repaint does not require a custom scene
   representation or reopen source/binding validation. Restore focus without emitting
   a domain event. No event can observe a new Session with old controls/handlers.
8. Revoke/dispose the old bundle and cancel its loads after successful publication;
   release its native/provider-owned resources once. Retained application-owned SDL
   engines/providers are not closed merely because a view changed. Start new lazy
   requests from the published bundle after commit, respecting the one-request rule.

All failures before step 7 dispose only the candidate and leave the old model,
handlers, provider requests, native resources, focus and offsets untouched. Failure
reporting belongs to an out-of-band status surface, not a half-published widget.
For this bounded host, select a prepared non-failing swap; do not implement a
partly fallible live mutation and call candidate disposal rollback. If native
integration cannot establish that swap, it must restore the complete previous
bundle before returning an error and requires renewed review of that boundary.

Do not publish a candidate twice or after disposal. On application Close, revoke
its lifecycle token first, cancel requests, release resources and close Session;
any queued completion then fails its acceptance predicate. Failed reload must not
close the existing Session. Raw 0.2 consumers retain existing APIs; a 0.2 document
may also use this new document-based host path, with its existing semantics.

Collection updates, loading/status changes and resize need the same preparation
ordering within a bundle: validate prospective runtime state, measured bounds and
native-adapter changes before committing. Extend the root-only CheckWith gate
only enough to cover that state; no generalized scene/publication framework is
required. Do not mutate data, then discover that its geometry cannot mount.
Provider invocation happens only after loading-state publication; provider failure
is recoverable visible state, not failed activation.

## 9. Presentation, export and consumer truthfulness

| Consumer | WCI1 contract |
| --- | --- |
| AST / normalized inspection | Profile and widget kinds exact; original spans/reuse paths retained; no provider data invented |
| Combined composition | Show tree/list kind, accessible label, symbolic Activate callback, viewport intent, declaration/use links and hidden instances; mark diagram as source composition |
| Dump / Markdown | Explicit Tree/List static description, provider-data-not-supplied label when source-only; no SVG default branch for unknown kinds |
| Collection SVG export | Return `unsupported-collection-export` with widget kind, exact instance path and original source span; supplying runtime state does not imply snapshot support. This explicit diagnostic satisfies WCI1's collection SVG boundary. |
| Native SVG background with SkipControls | May omit only kinds for which this prepared bundle actually has native controls; verify all such controls exist; public SVG export cannot use this internal option to hide unsupported widgets |
| Go codegen | Equivalent typed Document/Root for both profiles, with the explicit generator/version metadata above; user supplies providers/bindings at composition |
| Prototype readiness | State whether providers are supplied, and distinguish supported local interaction from unbound symbolic SDL activation; no connected-ready status from prototype Check |

WCI1 does not require a collection SVG renderer, runtime-state export API or new
snapshot file format. Explicit unsupported output is sufficient; never silently
omit a collection or replace it with an SVG placeholder. Unknown kinds/profiles,
unsupported scroll owners and missing capabilities also fail explicitly. Text
and composition are selected structural formats with source/instance identity,
not interactive fallbacks. No exporter invokes provider loading or callbacks.
Existing supported 0.2 SVG behavior remains covered. New diagnostics and maintained
descriptions are English.

Preserve ordinary 0.2 outputs and frozen XM-M2 sources/results/hashes; add new
versioned fixtures instead. Old binaries correctly reject 0.3; they must not receive
a relabelled 0.2 model. Update SDUI-owned helpers' advertised support as implemented.
Check current SDPTool delegation against those helpers; no second parser in SDPTool.
Matching external private helpers/package preparation remains WCI4, and this stage
does not claim an installed XFMD consumer now supports 0.3.

## 10. Acceptance IDs and required evidence

All IDs below are local WCI1 acceptance scenarios under the existing matrix, not
new global requirement/traceability IDs. Status for every row is **pending**.
Implementation evidence must name exact commit or source hashes, command, result,
artifact/native recording and independent review disposition.

| ID | Test and required result | Parent obligation |
| --- | --- | --- |
| WCI1-A01 | Exact 0.2/0.3 headers, positive tree/list, named/positional labels; negative versions, bad args/refs/profiles; 0.2 rejects new widgets locally and keeps syntax-only behavior | R27; Tree/List frontend |
| WCI1-A02 | Reused collection definitions get distinct paths/state/providers; Unicode spans/use chains intact; exact AST envelopes, stable 0.2 golden output and normalized serialization; 0.3 Instance.Profile survives clone/reuse/codegen, admission rejects mixed/unknown profiles and format mismatch | R05/R26; identity/consumer compatibility |
| WCI1-A03 | Compile and run generated 0.3 constructors; compare Document, normalized Root, spans and callbacks to parsed source; load provider state independently; 0.2 generation regression | R26; generated constructors |
| WCI1-A04 | Reject duplicate/empty/oversized IDs, bad UTF-8, missing parent/cycle, depth/item/label limits, wrong flags/kinds, list hierarchy, mixed valid/invalid batch; no partial state/generation changes | R16/R27; typed bounded data |
| WCI1-A05 | Replace one subtree and whole list; preserve unrelated branch/order and surviving selection; remove/recreate identical ID and row-to-group change; old targets reject and no index selection | R16/R27; incremental data/identity |
| WCI1-A06 | Controllable loader barriers: success, empty success, failure, invalid data/layout, Retry, superseding request, collapse/ancestor collapse, hide/disable, refresh, provider change, reload and Close; cancel empty-root load then explicitly restart, retry errored selectable branch, verify error collapse/re-expand stays paused and canceled branch re-expand starts once; each recovery gets a fresh token and exactly one request, zero SDL actions; late canceled success/error never publishes | R16/R27; loading/lifecycle |
| WCI1-A07 | Reject wrong payload/kind/sequence/revision/generation/hidden row; group/separator never activate; Select/Expand/Collapse/Retry never call SDL; reentrant data replacement invalidates action result | R16/R18; typed events |
| WCI1-A08 | Native keyboard and pointer exercise all interaction-table operations on long tree/list, nonselectable headings, mixed folder/file IDs and duplicate labels; double click activates exactly once; native row reuse/reorder retains correct stable-ID mapping; Escape pauses empty-root load, Tab/R and pointer Load restart it, focused selectable-branch R and pointer Retry recover without Activate; logs prove fresh token, one request, zero SDL actions and rejection of late canceled results | R27; native Tree/List |
| WCI1-A09 | Measured geometry oracle for vertical/horizontal scrolling, deep indent and long labels; error/clip/scroll, bounded viewport minimum, no infinite relative sizing; unsupported owner/indefinite axis/justify rejects | R15/R17; Scroll viewports |
| WCI1-A10 | Two nested viewports plus a sibling: transforms once, coherent native paint/hit clipping, one runtime offset authority with muted native sync, delta remainder chaining, focus EnsureVisible inner-to-outer, thumb drag, keyboard scroll, zero ranges; hidden/outside controls receive no events | R17/R27; nested viewports |
| WCI1-A11 | Resize/content deletion/collapse clamps offsets atomically; compatible reload retains state/focus/offsets, incompatible provider/kind/removed viewport revokes; too-small geometry keeps last valid presentation with diagnostic | R16/R17; reload and geometry |
| WCI1-A12 | Actual action-core parser/Engine/Go registration/Bridge: Activate extracts ItemId into Preview; both widgets share action; zero calls during preparation/navigation/loading; missing module, wrong registration/signature/source/target rejects | R12/R18/R28; real connected fixture |
| WCI1-A13 | Retained SDL engine across at least two UI reloads: increasing event/command sequence, old closures cannot invoke; dirty/deleted/recreated receiver rejects result; successful domain call is not replayed after UI failure | R16/R18; binding lifecycle |
| WCI1-A14 | Failure injection at profile, hidden capability, provider seed, bindings, Markdown/layout, native resource preparation and final stale check; old bundle/controls/focus/handlers/loads intact; candidate resources disposed once | R28; atomic native publication |
| WCI1-A15 | Successful publication and disposal with queued old load/events: no mixed-bundle observation, no old callback effect, no leaked native resources; current source/hash/window/provider/SDL revisions checked; two-phase prepare vs commit asserted | R28; connected activation |
| WCI1-A16 | Source composition/text explicitly identify collections/bindings, instance identity and no-data state; collection SVG rejects with unsupported diagnostic and correct path/source span, without a partial artifact; unknown kind rejects; no exporter callback/load execution; existing supported 0.2 SVG regression | R23/R27/R28; truthful export |
| WCI1-A17 | Full SDUI race suite, SDL bridge/runtime/codegen and affected SDPTool consumer checks; preserve 0.2 and frozen evidence; record unrelated baseline failures separately, never claim them passed | Plan verification; compatibility |

Run targeted tests in parser/runtime/layout/preparation/host/presentation/svg/
codegen during implementation; integration uses `go test -race ./...` from
`SDUI/go`, and `go test ./bridge ./runtime ./codegen` from `SDL/go` with the matching
SDUI dependency. Use barriers/channels for race scenarios instead of timing sleeps.
Unit/headless Fyne tests establish contract behavior, not native input evidence.

Use the available Xvfb with a separate display/user area. X11/Xtst and ImageMagick
`import` are present; xdotool and Python Xlib are not required. A small Python
ctypes harness may call XOpenDisplay/XTestFakeMotionEvent/XTestFakeButtonEvent/
XTestFakeKeyEvent and XFlush for actual OS pointer/keyboard delivery, with explicit
ctypes signatures and XKeysymToKeycode for the display keymap. Use Fyne canvas
Capture on the UI goroutine for screenshots, optionally paired with `import` for
the native window. Capture alone does not prove OS input delivery. No dependency
installation is required for this evidence route.

The native fixture must include two lazy
branches with controlled delayed/error replies and a list longer than the viewport,
and capture keyboard/pointer, nested scrolling, failed/successful reload and Close
with pending replies. Record source and binary identities, Fyne/environment,
actions and observed effects. Inspect at least empty, loaded, loading, error/retry,
selected and nonzero-offset native views. Inspect structural text/composition and
the source-linked collection SVG rejection separately. Screenshots alone do not establish
stale-response or callback-count acceptance; pair them with correlated fixture logs.

## 11. Review and coordinator handoff

Independent design review must challenge the grammar/envelope compatibility,
provider-generation rules, nested geometry, event/action split, retained engine
sequence and the actual non-failing native swap. Only after that review should
the implementation coordinator assign WCI1 code on its planned isolated phase
branch. Routine local implementation names may vary; changing these behavioral
contracts requires updating and reviewing this stage document.

This assignment authorizes writing only this new file in the original workspace.
No product code, Session, plan, card, ledger or other documentation is changed by
this Architect handoff. For coordinator Session upkeep, the work-summary payload
is: owner requested concrete WCI1 design; loaded sdp 1.1.1, sdp-architect 2.0.0 and
document-workflow; recovered Session0010 S1/S2; read the original workspace and
WCI0 report; selected the bounded contracts above; implementation/native acceptance
remain pending; next step is independent design review and reconciliation with the
actual WCI0 candidate. Host/turn IDs and exact prompt timestamp are unknown.

Before implementation closeout, coordinator updates language/EBNF/AST and consumer
docs, runtime/layout/architecture requirements, the acceptance matrix, plan,
Session and traceability with actual outcomes. The earlier Delivery-Proposal's
optional inventory wording is historical; the card and reviewed Design retain
the complete WCI2–WCI4 inventory. Do not close KB-SDUI-003 or any future family
from approval or delivery of this collection slice.

## Independent stage review disposition — Session0010 T002

Reviewer `01a11854-523d-7c11-adf4-f622b19b68c6` required accessible recovery
for canceled roots and failed selectable branches, and pruning mandatory collection
SVG snapshots/general scene/custom row rendering. Revised Load/Retry click and R
paths, A06/A08 and explicit export rejection satisfy those findings. The reviewer
approved bounded WCI1 implementation with no remaining blocking design findings.
This is design readiness, not native behavior or widget completion evidence.

Implementation clarification, same turn: a compatible reload of an outstanding
root request retains AutoLoadPending=false and changes Loading to Unloaded.
The root-only explicit Load/R path accepts this paused state; branch retry rules
and automatic loading are unchanged. A successor/recovery regression and native
recovery evidence are required; this is a contract clarification, not a pass claim.

## Native pilot refinement — independently review before implementation

The nested native pilot exposed overlapping outer frame and inner collection
scrollbar hit areas. Z-order changes merely exchange which thumb is unreachable;
host-only child shifts would violate the shared geometry contract. Reserve native
frame/group scroll gutters in shared measurement, with this bounded interface:

```go
type ViewportInsets struct { Right, Bottom float64 }
type ViewportMeasurer interface {
    MeasureViewport(*parser.Instance, float64) (ViewportInsets, error)
}
```

For 0.3 scrolling frames/groups only, an optional adjunct on Engine.Measure returns
pure native right/bottom insets. Native Fyne reserves its vertical/horizontal
gutter for each declared axis even when the current range is zero, avoiding a
range-dependent sizing loop. Other measurers default to zero; 0.2 behavior stays
unchanged. Values must be finite, nonnegative and bounded; a consumed viewport
rejects before publication. Collections retain their existing CollectionMeasurer.

Subtract insets after source padding and before child contents in both desired
measurement and arrangement. Relative/fr children reference the remaining finite
body. Content-sized non-scroll axes include the owning gutter once. Viewport.Rect
and Clip describe content; add HorizontalGutter/VerticalGutter screen rectangles
clipped by ancestors/window rather than by the content clip. Fyne paints and hits
thumbs in these declared strips without recomputing source padding. Track travel
uses the content viewport length, and runtime offsets remain authoritative.

This is a bounded native metric extension, not new source syntax or a scene model.
Add geometry tests for independent nested strips, relative allocation, shrink,
invalid metrics, padding and zero ranges; rerun actual native inner/outer drag and
wheel tests with both thumbs separately reachable. No acceptance is inferred from
this proposal. Coordinator-selected refinement stays inside WCI1 scope; independent
stage review approved this bounded implementation. Both track strips exclude
the corner intersection; ancestor clipping means the ancestors' content clips,
not their outer boxes. Reserve right only for declared vertical scrolling and
bottom only for declared horizontal scrolling. Wheel routing includes the owner's
effective gutter strips; child hit/paint clips exclude the parent gutters.
Reviewer 01a11854-523d-7c11-adf4-f622b19b68c6 reported no blocking design findings.
Native acceptance and full WCI1 completion remain pending.
