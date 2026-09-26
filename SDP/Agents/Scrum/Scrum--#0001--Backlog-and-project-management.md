# Scrum-0001 — shared backlog and project-management history

| Field | Value |
| --- | --- |
| id | SCRUM-SDP-0001 |
| project | SDP |
| state | completed |
| source | owner-conversation-2026-09-25 |
| participants | Project Owner via conversation; Codex review/recording |
| outcome | MAINT-SDP-0001 |

## Decisions and selected outcome

The owner requests one SDP board for all language/tool systems, optional
Sprint/Scrum grouping and a separate project-management ledger. KanBan moves
to SDP/KanBan because its scope is project-wide; Scrum remains under Agents as
requested. No separate queue/sprint card copies are introduced. Keep gate-review
as the existing canonical review state. This review selects
[MAINT-SDP-0001](../../Maintenance/PM1/Plan-and-Evidence.md); no Sprint starts.

Reviewed all 25 existing cards, including active review/handoff cards and closed
history. Six Ref cards become unnecessary separate-board routing; their distinct
local consequences transfer to the named primary before closing the Ref. This is
not a semantic merge of independent primary requirements, so no fake merge lineage
is created. KB-SDL-004 is actually split into KB-SDL-005 and KB-SDP-020 with typed
lineage KBO-SDP-000003. Its responsibilities remain fully assigned.

## Review and disposition

| Card | Resulting state | Disposition |
| --- | --- | --- |
| [KB-SDP-010](../../KanBan/active/%23010--Proposal--Document-consolidation.md) | gate-review | Owner review remains pending; non-blocking |
| [KB-SDP-017](../../KanBan/active/%23017--Proposal--sdptool-and-project-navigation.md) | ready | Keep bounded scope; reviewed dependencies and acceptance |
| [KB-SDL-001](../../KanBan/backlog/%23001--SDL--Proposal--Requirements-narrative.md) | backlog | Keep bounded scope; reviewed dependencies and acceptance |
| [KB-SDL-002](../../KanBan/backlog/%23002--SDL--Proposal--Links-through.md) | backlog | Keep bounded scope; reviewed dependencies and acceptance |
| [KB-SDP-003](../../KanBan/backlog/%23003--Idea--KanBan-graph.md) | backlog | Keep bounded scope; reviewed dependencies and acceptance |
| [KB-SDP-004](../../KanBan/backlog/%23004--Proposal--Design-traceability.md) | backlog | Keep bounded scope; reviewed dependencies and acceptance |
| [KB-SDL-005](../../KanBan/backlog/%23005--SDL--Change--System-and-source-sets.md) | queued | Next proposed implementation after maintenance |
| [KB-SDP-011](../../KanBan/completed/%23011--Bug--Traceability-id-conformance.md) | backlog | Keep bounded scope; reviewed dependencies and acceptance |
| [KB-SDP-014](../../KanBan/backlog/%23014--Proposal--KanBan-version-contract-and-distribution.md) | backlog | Keep bounded scope; reviewed dependencies and acceptance |
| [KB-SDP-018](../../KanBan/backlog/%23018--Study--Toolkit-audit-and-organization.md) | backlog | Keep bounded scope; reviewed dependencies and acceptance |
| [KB-SDP-020](../../KanBan/backlog/%23020--Change--Shared-design-source-organization.md) | backlog | Keep bounded scope; reviewed dependencies and acceptance |
| [KB-SDP-001](../../KanBan/completed/%23001--Proposal--Project-structure.md) | completed | Preserve completed or superseded history |
| [KB-SDUI-001](../../KanBan/completed/%23001--SDUI--Ref--SDP--017--sdptool.md) | completed | Routing consequence transferred to its primary; no product completion claimed |
| [KB-SDUI-002](../../KanBan/completed/%23002--SDUI--Ref--SDL--004--Language-source-organization.md) | completed | Routing consequence transferred to its primary; no product completion claimed |
| [KB-SDL-003](../../KanBan/completed/%23003--SDL--Ref--SDP--017--sdptool.md) | completed | Routing consequence transferred to its primary; no product completion claimed |
| [KB-SDP-005](../../KanBan/completed/%23005--Ref--SDL--002--Links-through.md) | completed | Routing consequence transferred to its primary; no product completion claimed |
| [KB-SDP-006](../../KanBan/completed/%23006--Ref--SDL--001--Requirements-narrative.md) | completed | Routing consequence transferred to its primary; no product completion claimed |
| [KB-SDP-007](../../KanBan/completed/%23007--Change--KanBan-foundation.md) | completed | Preserve completed or superseded history |
| [KB-SDP-008](../../KanBan/completed/%23008--Change--Visible-card-metadata.md) | completed | Preserve completed or superseded history |
| [KB-SDP-009](../../KanBan/completed/%23009--Change--Card-merge-and-split.md) | completed | Preserve completed or superseded history |
| [KB-SDP-012](../../KanBan/completed/%23012--Change--Card-history.md) | completed | Preserve completed or superseded history |
| [KB-SDP-013](../../KanBan/completed/%23013--Change--English-documentation.md) | completed | Preserve completed or superseded history |
| [KB-SDP-015](../../KanBan/completed/%23015--Change--CardState-and-shell-cli.md) | completed | Preserve completed or superseded history |
| [KB-SDP-019](../../KanBan/completed/%23019--Ref--SDL--004--Language-source-organization.md) | completed | Routing consequence transferred to its primary; no product completion claimed |
| [KB-SDP-002](../../KanBan/superseded/%23002--Proposal--sdptool.md) | superseded | Preserve completed or superseded history |
| [KB-SDL-004](../../KanBan/superseded/%23004--SDL--Change--Language-source-organization.md) | superseded | Split into System/source contract and model migration |
| [KB-SDP-016](../../KanBan/superseded/%23016--Proposal--SDP-discovery-and-viewer-capabilities.md) | superseded | Preserve completed or superseded history |

## Deliberately separate scopes

Narrative language and weak-path semantics have different acceptance criteria.
Evidence aggregation, historical-ID CI repair and KanBan distribution are also
separate deliveries. Toolkit audit is broader than KanBan distribution. They may
share a future Sprint without losing their separate acceptance. SDPTool keeps its
phase plan and uses these capabilities; it does not own their implementations.

## Next selection

Propose KB-SDL-005 queued after this maintenance; KB-SDP-020 depends on its
supported contract. KB-SDP-017 remains active/ready, and KB-SDP-010 remains on the
owner's explicitly deferred gate-review. No held items were found/selected. A new
Scrum can change this order with a recorded reason; queued is not automatic work.

## Record limits

This is the actual conversational review, not a reconstructed calendar meeting.
It records a maintenance outcome and future work selection, not completed system
features. The common ledger links this Scrum to its source cards and maintenance.
