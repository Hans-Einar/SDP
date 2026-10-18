# Session0011 T002 — concurrent-work preflight

Verdict: **RED for SDPTool action catalogue/JSON implementation**. The owner permits
implementation only after a green result. No product implementation was started.
The new [local procedure](../../../ProjectManagement/Concurrent-Work-Preflight.md)
was exercised; [inventory.json](inventory.json) records exact heads/status and
committed merge forecasts. These are observations, not identification of live agents.

## Intended write boundary

KB-SDP-051 (Discoverable SDPTool actions, blueprint branch) would modify CLI input,
command registry/dispatch, presentation/error handling, consumer contract/docs and
possibly discovery. New packages alone would not remove the shared API boundary.
The shared board, canonical ledger and ID namespace also need coordination.

## Findings

| Workstream | Observed work | Overlap / consequence |
| --- | --- | --- |
| ProjectGovernance, Session0007 / KB038; SDP-project-governance at 0b0e82c | Clean checkout, unintegrated governance CLI/MCP/controller changes; Session active | cli.go command routing, output_arguments.go, presentation, go.mod/go.sum; source authority and invocation design must be reconciled |
| Runnable programs, Session0010 / new KB051; sdp-runnable-programs at 9963159 | Clean committed branch plus related dirty primary checkout | cli.go, output_arguments.go, discovery/project/SDUI and consumer contract/presentation changes; catalogue must account for actual run operations |
| Primary SDP-vNow at 74e87aa | 328 status entries at capture, including shared SDPTool files and process records | Cannot assume branch-only forecasts include the current candidate; preserve uncommitted work |
| Card identity collision | Primary active/#051--Change--Runnable-SDUI-programs.md declares KB-SDP-051; our backlog/#051--Proposal--Discoverable-SDPTool-actions.md also declares KB-SDP-051 | Distinct subjects created in separate workstreams. Original local registration predates the other card; neither identity may be silently rewritten during integration |

The primary Session0008 active status is a stale branch projection; current
blueprint-branch Session0008 is paused. Historical release/research worktrees are
not considered active workers solely because their old records say active.
The old rp3 candidate has extensive deletions; preserve it and do not interpret
those as authorized cleanup or as an active new SDPTool implementation.

## Merge forecasts

All three committed forecasts returned conflict status 1:

- blueprint vs governance: management and traceability ledgers, README,
  output_arguments.go, presentation/presentation.go.
- blueprint vs runnable programs: both ledgers, Contract.md, README,
  output_arguments.go.
- governance vs runnable programs: both ledgers, README, cli.go,
  output_arguments.go.

No branch was merged and no worktree files were changed by these forecasts.
Concurrent dirty edits, live-owner activity and remote-only work are not verified.

## Conditions for green

1. Establish which exact governance/program revisions are selected for integration
   and which owner handles shared dispatch, output and discovery contracts.
2. Reconcile the duplicate card identity and any event-ID collisions explicitly,
   preserving historical records and updating current references across workstreams.
3. Select a combined or clearly sequenced implementation baseline. Resolve/test
   known conflicts in a separately authorized integration effort; do not merge to
   main implicitly. Agree any remaining shared-file ownership.
4. Refresh local statuses and semantic boundaries, then rerun this preflight.

Safe work now: this procedure, evidence, Session queue and coordination handoff.
No automatic cross-agent notification, reservation service or global lock exists.
No native XFMD work, code edits, merge or publication is included in this check.
