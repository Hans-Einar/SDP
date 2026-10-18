# Widget capabilities, activation and interaction

Design candidate under PLAN-SDP-0021, KB-SDUI-003, Session0010 T001.
This document selects implementation contracts; it does not claim parser or
host support. The implemented source remains SDUI 0.2 until evidence says otherwise.

## Evidence and alternatives

The current parser accepts only button/input/svg. Runtime Value carries string
and boolean, but editable state and Commit are input/string-specific. The Fyne
view mounts input/button and skips unrecognized controls. Layout rejects scroll.
SDL bridge validates its complete plan before binding, but only supports a text
result into an input; its existence is not a generic collection binding contract.
KB-SDUI-005 checks a local prototype and reports unbound callbacks; retain that
truthful distinction. These observations come from parser/validate.go,
runtime/session.go, host/fynehost/view.go, layout and SDL/go/bridge/bind.go.

Extend these packages rather than introduce another UI model, renderer or binding
engine. A generic opaque native-object escape hatch loses typed identity and
export diagnostics and is rejected. Static Markdown rows cannot satisfy collection
interaction. Fyne supplies controls, while SDUI owns state and shared geometry;
Fyne alone cannot define source semantics or prove keyboard/accessibility support.

## Profile and capability decision

Preserve exact SDUI 0.2 semantics and frozen fixtures. Implement added syntax under
exact `sdui 0.3;`, with a corresponding EBNF/profile and explicitly versioned AST
representation. Do not make 0.2 silently accept new controls. A parser supporting
both versions uses one implementation with profile-specific validation, not a
legacy fallback parser. Profile 0.3 is selected for development, not published.
Existing valid 0.2 sources remain useful consumer compatibility fixtures.

Separate source profile from execution capabilities. Each required feature has
an identifier, version, source span and instance path. A concrete target declares
frontend, layout, provider/export and native-host support separately. Unknown or
missing capabilities reject the candidate before controls/handlers become active.
Hidden or initially inactive pages are checked too; visibility is not an escape
from preflight. Providers advertise real rendering separately from placeholders.

Static text/composition exports may describe every parsed family with an explicit
static label. SVG must implement a declared snapshot or reject; a fallback must
be explicitly selected by the caller and visibly labelled. Code generation must
construct equivalent versioned AST/instances or reject. No exporter may silently
fall through an unknown kind. A syntax-only parse never establishes support.

## Preparation and publication

1. Parse/profile-check and normalize a bounded source revision, preserving
   declaration/use spans and reusable-instance paths.
2. Select the entry and inspect all required capabilities, including provider
   resources, scroll, transients and event payloads. Validate measured layout.
3. For connected operation, load modules outside the frontend and prepare the
   existing SDL bridge's complete typed binding plan. Reject missing module,
   action/signature, handle or payload before publication. Local prototype mode
   identifies unbound declarations and cannot claim connected readiness.
4. Prepare controls/resources and a complete handler set against the same source,
   UI revision and provider/binding epoch. Preparation invokes no domain action.
5. Publish the candidate on the owner UI goroutine, then dispose replaced resources.
   Any preparation failure leaves the previous model, handlers and focus intact.
   The application composition root owns publication. Select a non-failing final
   swap of a complete prepared model/handler/resource bundle, after all fallible
   operations and a final revision check. No domain callback or fallible allocation
   belongs in that swap. A host that cannot provide this contract must retain the
   previous complete bundle and restore it before returning a publication error;
   candidate disposal alone is not rollback of a partly installed bundle.

Extend bridge preparation/installation at its existing boundary when necessary;
do not accept a caller-supplied boolean as evidence of a validated binding plan.
An asynchronously completed preparation must recheck revision before publication.
Cancellation is not domain rollback. An already accepted application transaction
cannot be undone by a later UI failure; report that outcome without automatic retry.

## Typed state and event ownership

Keep Session/Path/Generation/Kind handles and the existing model revision, sequence
and draft revision checks. Add explicit typed state rather than encode booleans,
numbers, item IDs or structured events into a label/value string. Numeric values
must be finite; reject invalid bounds/step. Invalid intermediate numeric text is
a draft, not an accepted number. Programmatic state updates never emit user events.

Collection target = UI handle + model revision + collection generation + stable
item ID. A data refresh advances collection generation independently of UI source
reload. A request token additionally contains parent item ID and monotonically
increasing request generation. Stale, canceled, deleted, replaced or disposed
targets cannot publish a result. Reusing an item ID does not revive an old token.
Page, option and command IDs are stable identities, never display labels or indexes.

Typed event variants distinguish activate, select, expand/collapse, retry,
change, commit, context, accept/cancel/close and split adjustment. Each variant has
an explicit payload and applicable kinds. Validation precedes user callbacks.
Reject mismatched kinds, invisible/disabled/read-only targets where applicable,
stale revision and duplicate sequence. A draft update has its own revision and
cannot overwrite a newer accepted value. Atomic batches validate every update
before publishing any. Application adapters own persistence and form transactions.

## Collection and viewport slice

Canonical widgets are `tree` and `list`; initial argument is an accessible label.
Provider data belongs in a typed runtime API, not embedded filesystem calls or a
second source language. Fixture providers exercise unbound local mode. Callback
references remain symbolic. New event references require explicit typed bridge
mapping; do not overload today's text-only callback source.

Bound each supplied snapshot to 4096 items, hierarchy depth 64 and 32768 UTF-8
bytes per label. Reject duplicate IDs, missing parents and cycles atomically.
These are selected initial resource limits, not virtualization/performance claims.
Application owns ordering and pagination; runtime supports bounded replacement and
incremental updates. Item kinds include selectable row, group and separator.
Group/separator rows cannot be selected or activated. Branch expansion does not
activate an item. Loading/error/retry are visible per branch/request. Distinguish
selection from activation (pointer single/double click; keyboard arrows/Enter).
Removing a selected item clears selection; never select the row now at its index.

Use existing overflow-x/y scroll syntax for negotiated viewport support. Layout
owns viewport and content extent; runtime owns offset and clamping; Fyne routes
wheel/keyboard and paints/hit-tests the same translated/clipped geometry. Nested
scroll reaches the nearest eligible viewport and does not dispatch through clipped
controls. Focus movement scrolls its target into view. Resize/content deletion
clamps offsets. Reload retains compatible offsets and identities, invalidates
outstanding requests and clears removed selection. Export declares captured offset
and clipping, or requires an explicit whole-content/static fallback.

## Remaining widget decisions

The first connected collection fixture maps an explicit event item-ID field to
SDL action-core 0.1 text and a result field to the existing text preview input.
Collection loading remains a typed Go provider responsibility; do not insert
hidden inputs or serialize a generic collection into text. Existing action-core
supports text/integer/boolean only. General numeric controls may use finite Go
numbers locally. For this inventory select checked integer binding to action-core
0.1: connected bounds, step and value must be integral and losslessly representable
within both binary64 and signed 64-bit domains (absolute value at most 2^53-1).
Reject fractional/out-of-range bindings before activation; never round or encode
numeric payloads into text. Fractional controls remain usable with typed Go adapters;
decimal SDL transport is explicitly unsupported, with a capability diagnostic.
A later versioned action profile extension is optional separate work.

| WCI1 connected fixture field | Mapping and ownership |
| --- | --- |
| User action | Collection Activate only; Expand, Select and Retry do not invoke this action |
| Input `ItemId` | Validated stable item ID extracted as SDL text by an explicit bridge event-field source |
| Result `Preview` | SDL text applied through the existing explicit setHandle target, a visible input/preview receiver |
| UI/collection identity | Handle, source revision, collection generation, item existence, request and sequence checked locally before dispatch; rechecked before result publication |
| Domain generation | Supplied application context when required; never substitute it for UI or collection revision |
| Loading/result lists | Typed Go provider API; not generic SDL text serialization |

WCI1 owns the bridge event-field extension and its failure tests. WCI0 verifies
the existing input/button bridge on a detached candidate, without adding collection
or numeric language support. The receiver is visible fixture content, not a hidden
input used to smuggle generic state.

| Family | Canonical representation selected for elaboration | State / lifecycle |
| --- | --- | --- |
| Button/toggle | Extend button with typed checked state, semantic command and icon identity; no separate toggle keyword | Shared command identity across toolbar/menu/keyboard; exclusive group clears peers atomically |
| Single-line / multiline | input plus multiline property (not a separate editor language) | Text drafts, read-only, validation, explicit change/commit; Fyne undo/redo, Unicode/IME and clipboard verified natively |
| Tabs | tabs composition of named page content | Selected page ID; preserve hidden page drafts/focus; removal revokes identity and chooses a documented enabled fallback |
| Split pane | split composition of two content children, axis and relative proportion | Minimum relative extents, drag/keyboard, collapse/restore; no source pixel dimensions |
| Menu/context | menu composition backed by semantic command IDs | Stable context target; grouping/submenus; Escape dismisses without action |
| Dialog | dialog surface containing ordinary composed content | Explicit modality and result; parent generation/lifetime, focus capture/restoration; cancel never accepts drafts |
| Checkbox | checkbox with boolean value | Change/commit, read-only/enabled; no implicit tri-state |
| Slider | slider with finite numeric bounds/step/value | Distinguish changing from committed value; clamp user movement, validate programmatic updates |
| Choice | select with stable option IDs | Empty/disabled options and explicit invalid selection; option removal cannot select by old index |
| Numeric | number input with editable numeric draft and increment/decrement | Preserve invalid draft, validate acceptance and ranges; use same numeric contract as slider |
| SVG/Markdown | Existing symbolic content with versioned provider/resource capability | No automatic fetching; accessible description and explicit labelled fallback; glyph fidelity owned by KB-SDUI-004 |

Container forms need EBNF/AST elaboration before their respective implementation
milestones. The table chooses semantic representations, not an already accepted
call syntax; no rejected XM-M2 spelling is adopted by implication. Shared semantic
properties include visible, enabled, accessible label, read-only where meaningful,
focus order, tooltip, validation and command/icon identity. Define one common
property/event schema before adding per-family cases. Domain Load/Save actions
remain explicit application/SDL operations.

## Consumer and release boundary

Keep XM-M2 sources, expected results and hashes as dated evidence. New profile
tests get a new baseline; normal discovery continues to show historical invalid
sources and diagnostics. SDPTool must recognize the new profile/capabilities and
delegate to matching SDUI-owned helpers. XFMD privately bundles its helpers;
updating gh-sdp alone is insufficient. Prepare matching distributions and an
external consumer handoff; actual XFMD writes, installation or upstream publication
need their own applicable authorization. No new FOX host, Rust renderer, component
source sets, full editor or full XFMD parity is required by this design.
