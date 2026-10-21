# SDPTool — feature requirements

Feature: **SdpTool**. Design record: **DES-SDPTOOL-001**. Status: initial design.
[Mandate](../01--Mandate/Mandate.md) · [Architecture](../03--Architecture/SDPTool.md) ·
[Implementation plan](../05--Implementation/SDPTool.md).

## Requirements and acceptance

These requirement IDs are project records, not new SDL keywords. All describe
target behavior; model validity does not establish their implementation.

| ID | Requirement / use case | Acceptance to implement |
| --- | --- | --- |
| REQ-SDPTOOL-001 | Preview a supported saved design file through one facade; PreviewDesign | A standalone valid .design produces a source-revision-tagged Markdown/Mermaid/SVG bundle using existing Go SDL services and Rust layout. Relevant default views are bounded; malformed/profile-invalid source gives positioned diagnostics. No domain execution or source overwrite. |
| REQ-SDPTOOL-002 | Resolve an explicit project or current directory; BrowseProject | Recognize the selected valid SDP area or its SDP child; no implicit parent discovery. Specify validity against reconciled versioned metadata. Distinguish absent, invalid and unsupported contexts; behave equally in mono/separate repos. |
| REQ-SDPTOOL-003 | Expose general native/document navigation; BrowseProject | Stable node IDs, kinds, labels, references, revision and selectable targets for all supported groupable concepts. Non-UseCase collections, variable depth and graph cycles are covered; details generated on demand. KanBan and SDUI data retain their owning contracts. |
| REQ-SDPTOOL-004 | Preserve source and session identity; PreviewDesign | Saved-file producer first; later unsaved snapshot/revision contract rejects stale output. Errors preserve previous valid preview; typed bundle/resource ownership is explicit. XFMD implements its source editor/preview pairing separately. |
| REQ-SDPTOOL-005 | Propose and validate implementation plans; ReviewImplementationPlan | Explicit goals/prerequisites/acceptance, preserve approved plans and show proposed changes. Separate proposed/planned/implemented/verified using KB-SDP-004 evidence semantics. Later scope, not claimed by generated model status. |
| REQ-SDPTOOL-006 | Inspect card history without another version engine; InspectCardHistory | Use Git plus stable KanBan IDs/lineage for revision/draft diffs; report unavailable history. Later scope using existing shared KanBan contracts. |
| REQ-SDPTOOL-007 | Install and upgrade project SDP through the common Go command entry point; owner decision 2026-09-26 | SDPTool owns planning/apply, preservation, recovery and installed identity; no required PowerShell or competing Toolkit Go executable. Thin gh-sdp delegates to SDPTool. Support explicit adoption of a verified manual baseline with a selected release manifest; test on a disposable XFMD copy before live upgrade. Implemented for explicit development/signed-test inputs and verified on Linux, including a disposable XFMD copy; see PLAN-SDP-0003. Production trust/release selection, live adoption and other native platforms remain unselected. |


The canonical SDL model is [SDPTool.design](../03--Architecture/SDPTool.design).
Its actors/use cases/feature facts support A0/A1 views of these requirements.
Requirement prose and acceptance IDs are not SDL keywords. design-core 0.5 lacks
imports, so do not duplicate declarations in independent phase files or invent
#include to join them. Split the model only when a supported composition contract
is available. The phase folder classifies the work, not every model object.

## External KanBan references — REQ-SDPTOOL-003 clarification

Owner decision 2026-09-27 (KB-SDP-034): permit references to cards outside the
board's declared local namespaces without requiring their existence or a local
checkout. Report them explicitly as external and unverified; keep local cards
openable. Do not locate external cards or introduce a cross-project configuration.
Missing local references and invalid history remain errors. The
[consumer contract](../../SDPTool/Contract.md) defines the additive representation.

## Sessions browsing — REQ-SDPTOOL-003 clarification

Owner request 2026-10-01 (KB-SDP-047): discovery must expose SDP/Sessions for
browsing. Provide a directory-derived capability and dedicated navigation root,
with nested documents and typed open targets. Preserve ordinary Files navigation
and refresh on additions/edits/moves/removal without registration. Missing, empty
and unavailable directories are distinguishable; symlinks are not followed.
This is document navigation, not validation or execution of the Session process.
PLAN-SDP-0015 and the producer contract define the bounded implementation.

## Complete viewpoint overviews — REQ-SDPTOOL-003/004 clarification

KB-SDP-032 and XFMD Session0005 select optional, typed overview targets for VP01
and VP07. Selecting the viewpoint title generates its aggregate on demand;
expansion does not generate views. Preserve original detail targets and old-client
compatibility. VP01 retains all actors/use cases, including isolated declarations,
with separate realization context. VP07 retains the complete declared structural
base and proven Feature annotations, using logical or explicit named-mode scope
without mixing allocations. Reject unsupported filters and stale source snapshots
before publishing; composed sources bind the full source graph. No relationships,
scenario order or conditions may be invented. See the
[contract](../../SDPTool/Contract.md) and
[integration work unit](../05--Implementation/SDL/ViewpointNavigation/Plan.md).

## VP02 complete architecture — REQ-SDPTOOL-003/004

Owner continuation after the VP01/VP07 delivery selects VP02 (XFMD Session0007).
Overview activation must retain every declared System, Container and Unit and
explicit contains edge, including isolated roots. Keep original parent/root
slices separate and preserve exact source provenance. Logical structure must not
be represented as deployment, Functionality ownership or runtime communication.
Empty structural inventories are explicit; unsupported filters cannot turn a
partial graph into a claimed complete overview. See ViewpointNavigation/N5.md.

### VP03 overview acceptance (N6)

The optional VP03 overview must retain all structural declarations and all declared
responsibilities/capabilities. Distinguish immediate ownership, realization
contribution and explicit capability provision; do not infer one from another or
substitute Features. Show isolated declarations and missing providers/contributors.
Keep source evidence, old detail targets and complete-selection freshness checks.
Architecture annotations and linked registers may summarize explicit fact paths
without adding language relations. Verify semantic, composed-source and native
consumer behavior; large-format publication remains a separate obligation.

## VP04 interface overview — N7

VP04 advertises the optional overviewTarget with projection=overview. It must
retain every declared interface, every explicit consumes fact and its consumer
(including Containers). Unconsumed interfaces remain visible. Dashed consumes
arrows describe design-time use, not observed calls or required availability.
No provider, Channel, transport, call ordering or mode necessity is inferred.
Source-linked inventory/register accompanies the unchanged Interface use table.
Empty inventories are explicit. Other structural units belong to VP02. Existing
sourceLinks/composed-source freshness rules apply; no schema change or new
count-only cap. Complete overview rejects mode/focus/relation/diagram/traversal
filters and preserves ordinary VP04 table selection.

## VP05 mode dependency overview — N8

VP05 overviewTarget without mode returns a declared-mode index, counts and source
links, never a union dependency graph. Advertise an existing-style overview-mode
target for every declared mode. With mode, select all Capabilities and the exact
mode's requires facts/Interfaces. Empty capabilities/requirements are explicit.
No inheritance, sufficiency, provider, consumption or runtime-flow inference.
Preserve existing detail IDs, VP07 targets and composed source freshness. Reject
unknown modes and other partial-overview filters; add no count-only truncation.
Consumers must retain viewpoint as well as mode when identifying mode rows.

## VP06 delivery-plan overview — N9

Advertise VP06 overviewTarget for all Activities with refines, a separate explicit
depends-on graph, source-linked status and addresses/owns/delivers/illustrates
registers. Retain isolated/statusless activities and unaddressed Functionality
inventory. Implementation status is a source claim, never completion evidence;
no inherited status/dependencies, scheduling, dates or scenario execution is
inferred. Empty cases are explicit. Keep existing detail and implementation.md
behavior; reject partial overview filters, preserve composed-source freshness
and add no count-only truncation or native/schema changes.

## VP08 communication overview — N10

Advertise an on-demand VP08 overview index and explicit named-mode overview targets.
The index lists all channels/modes/scenarios without a merged graph. Per-mode output
retains every channel and only explicit participants of that mode, grouping topology
edges by channel/participant/role while preserving all supporting source facts and
individual message bindings. Scenario and MessageSet scope is mode-isolated. Existing
scenario sequence targets remain unchanged. Empty scope is explicit; no timeline,
mode inheritance, transport or runtime success is inferred. Composed-source links,
revision freshness and generic consumer contracts apply unchanged.

The selected mode page partitions topology into one diagram per Channel. Every
channel neighborhood and supporting uses fact is retained; participants repeated
across diagrams keep the same declaration identity. This is a semantic partition,
not a size threshold or truncated graph. It avoids an all-channel union that cannot
currently be routed within the native renderer budget. Sequence targets and the
complete source-linked registers remain unchanged.

## VP09 data/storage overview — N11

Advertise a complete on-demand VP09 overview. Partition diagrams by Dataset with
explicit origin, holders, Database owners and data contracts. Include all data,
Database, Contract, Variant and Field definitions in source-linked registers,
including idle/unused definitions; retain explicit shape membership and ternary
projection/owner proofs. Detail contract diagrams remain independently selectable.
Do not infer SQL/ER semantics, persistence for Unit-held data, transformations or
wire encodings. Missing inventories/properties are explicit. Existing source and
freshness guarantees apply; no silent content truncation or new consumer schema.

## Deferred model-tree navigation — N15a refinement

REQ-SDPTOOL-002/003/004 and KB-XFMD-025 require source visibility, on-demand work
and revision-bound selection. A project must be able to expose all its discovered
SDL models without eagerly combining every model graph into one navigation reply.
Provide opt-in compact discovery and explicit revision-bound single-model expansion
through both CLI and the JSON action catalogue. Preserve invalid/context-dependent
sources and diagnostics, original node/detail/overview identities, non-SDL tabs and
legacy eager behavior. Changes to composed inputs invalidate prior expansion.
No policy silently excludes ordinary evidence models. Native expansion/refresh and
retirement of the remaining filesystem/source/per-model/GUI caps are subsequent
work; this bounded producer service does not claim unlimited project scale.
