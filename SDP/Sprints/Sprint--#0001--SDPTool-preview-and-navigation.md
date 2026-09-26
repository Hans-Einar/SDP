# Sprint-0001 — SDPTool preview and project navigation

| Field | Value |
| --- | --- |
| id | SPR-SDP-0001 |
| project | SDP |
| state | completed |
| ScrumId | SCRUM-SDP-0002 |
| Members | KB-SDP-021, KB-SDP-022, KB-SDP-023, KB-SDP-024, KB-SDP-025, KB-SDP-026 |
| source | KB-SDP-017; owner-conversation-2026-09-25 |

## Goal and boundary

Deliver a usable Go SDPTool facade for saved-design preview, explicit project
recognition, configured plan viewing and current-source navigation data consumable
by XFMD. Keep SDL/SDUI semantics in their owning tools. This is a delivery group,
not an invented calendar commitment; no timebox was requested.

The [single feature plan](../05--Implementation/SDPTool.md) remains authoritative
for phases/milestones. This record owns membership, order and sprint acceptance.
No implementation starts during this planning Scrum. On actual start, move all
six cards from backlog to active/ready together, record their transitions and the
Sprint start; mark only the currently worked card in-progress.

## Members and order

| Card | Phase | Depends on | Outcome |
| --- | --- | --- | --- |
| [KB-SDP-021](../KanBan/completed/%23021--Change--SDPTool-saved-design-preview.md) | P0-M1 | None; existing SDL pipeline | Saved-file design preview |
| [KB-SDP-022](../KanBan/completed/%23022--Change--SDPTool-discovery-contract.md) | T1-M1/M2 | Existing manifests and P0 contract alignment | Project recognition and delegation contract |
| [KB-SDP-023](../KanBan/completed/%23023--Change--SDPTool-project-viewer-bridge.md) | T2-M1/M2 | 021, 022 | Project resolver and configured viewer bridge |
| [KB-SDP-024](../KanBan/completed/%23024--Change--SDPTool-model-navigation.md) | T3-M1/M2 | 021–023 | Model-derived navigation and selected generation |
| [KB-SDP-025](../KanBan/completed/%23025--Change--SDPTool-KanBan-and-SDUI-inventory.md) | T3-M3 | 022, 023; coordinate 024 | KanBan and SDUI navigation services |
| [KB-SDP-026](../KanBan/completed/%23026--CodeReview--SDPTool-consumer-contract-review.md) | T4-M1 | 021–025 | Consumer contract and end-to-end delivery review |

Each phase uses a branch stacked on the previous phase; commit each completed
milestone with plan and evidence updates, then push the phase. Register actual
system changes/evidence in Traceability with system-prefixed IDs and Sprint/card/
management-event references. Card completion is not implementation evidence.

## Completion and dependencies

Demonstrate saved-source preview and regeneration, deterministic discovery,
configured plan viewing, all supported grouped SDL concepts, KanBan/SDUI inventory,
versioned consumer fixtures and stale/error/resource handling. Use real tool
output and preserved-source checks. Publish exact supported profiles/capabilities
and concrete limitations. KB-SDP-026 reviews the combined result; close the Sprint
only after every member has an explicit disposition. Removed/unfinished work
requires recorded membership changes, never silently completed cards.

Native GUI implementation/verification belongs to XFMD KB-XFMD-014/015. The
consumer harness proves producer behavior, not a delivered XFMD sidebar. No
native consumer release is required to close this producer Sprint; real GUI
interoperability remains T4-M2 in KB-SDP-017. P0-M2 unsaved buffers and T5 slice
proposal/validation, evidence-aware plans/roadmaps and Git history/diff also remain
there. Existing saved single-file models avoid a dependency on KB-SDL-005 and
KB-SDP-020; new language/source-set claims still require those owning deliveries.

KB-SDP-010 owner review is non-blocking. KB-SDP-011's existing Toolkit ID failures
must be reported honestly if affected; do not silently repair or suppress them.
KB-SDP-014 compatibility work is coordinated through T1, not assumed delivered.
Skills and installer Scrums (#027/#028), the broad Toolkit audit and the history
graph are outside membership. No daemon, startup compilation or second parser is
required. Current project metadata gaps are addressed in T1, not hidden by a
folder-exists check.

## Planning evidence

PM2-M1 records Scrum-0002, the partial split KBO-SDP-000004 and this planned
membership. Validation results are recorded in the Scrum document; no product
implementation or native UI evidence is claimed here.

## Execution

2026-09-25T13:08:54Z — EVT-PM-SDP-000010: Sprint started on owner direction. #027/#028 remain backlog; review them after Sprint outcomes are available.

## Completion

EVT-PM-SDP-000022: all six members are completed. P0-M1, T1-M1/M2, T2-M1/M2,
T3-M1/M2/M3 and T4-M1 have separate milestone commits/evidence. The producer is
implemented; native XFMD T4-M2, unsaved P0-M2 and advanced T5 remain in #017.
[Review](../Verification/REVIEW-SDP-0001.md) records tested scope and limitations.
KB-SDP-027/028 received concrete implementation findings but remain backlog;
neither Maintenance Scrum has started. No unrelated backlog scope is closed.
