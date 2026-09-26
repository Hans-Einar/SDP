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


The canonical SDL model is [SDPTool.design](../03--Architecture/SDPTool.design).
Its actors/use cases/feature facts support A0/A1 views of these requirements.
Requirement prose and acceptance IDs are not SDL keywords. design-core 0.5 lacks
imports, so do not duplicate declarations in independent phase files or invent
#include to join them. Split the model only when a supported composition contract
is available. The phase folder classifies the work, not every model object.
