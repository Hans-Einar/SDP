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
