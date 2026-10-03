# PGD1 — First governed Codex workflow

| Field | Value |
| --- | --- |
| id | PLAN-SDP-0017 |
| project | SDP |
| state | completed |
| PlanType | DesignPlan |
| BranchPolicy | current |
| CommitPolicy | phase |
| Systems | SDPTOOL |
| source | KB-SDP-038; owner continuation, 2026-10-03, Session0007 T002 |

## Outcome and authority

The owner accepted continuing the next step proposed in Session0007: design one
governed workflow joining KB036 routing, KB037 execution state, KB042 Session
continuity and KB038 client integration. Deliver an implementable pilot boundary
and a planned implementation handoff. The continuation selects this design work;
it does not select full implementation of all related cards, merge or release.

## Scope and baseline

Reuse [RGS2](../../../02--Requirements/RoutineGovernance/Synthesis.md) and existing
SDPTool discovery, board reading and Session browsing. The owner's ProjectGovernance
scope covers method/tooling for owner-agent collaboration. ModelGovernance remains
the other workstream's responsibility. KB043 is a downstream timeline consumer;
the first pilot captures its required identities without delivering its renderer.

Deliver [Design.md](Design.md), [Acceptance.md](Acceptance.md), compatibility
inspection evidence and a bounded ImplementationPlan. No product code, new
distributed process profile, account configuration or live model execution here.

## Git policy

Use the current `sdp/blueprint-model-history-proposal` branch for this documentation
phase, preserving the concurrent ModelGovernance workstream and existing branch
history. Do not switch the shared worktree underneath another agent. Commit PGD1
as one phase delivery with scoped paths after validation; exclude unrelated files.
The implementation handoff must reassess branch/worktree isolation before code
work. Existing phase-push permission is separate; no merge/release is authorized.

## Phase and milestones

| Phase | Milestone | Acceptance | State |
| --- | --- | --- | --- |
| PGD1 Design | PGD1-M1 | Ownership, identities, routine states, API and recovery contract cover the first workflow | completed |
| PGD1 Design | PGD1-M2 | Installed Codex version/schema inspected with provenance and explicit runtime limits | completed |
| PGD1 Design | PGD1-M3 | Acceptance cases and bounded implementation handoff preserve existing authority and non-goals | completed |

## Evidence and outcome

Evidence is recorded in [Evidence.md](Evidence.md). Design consistency and schema
inspection do not establish runtime enforcement, account access, MCP negotiation,
independent review or successful model execution. Completion state is updated only
after all three design milestones and record validation are complete.

PGD1-M1–M3 delivered: shared contract and 18 mapped acceptance cases; offline Codex
0.160.0 default-schema probe; [PLAN-SDP-0018](../../../05--Implementation/SDPTool/ProjectGovernance/Plan.md)
registered planned. This completes the design assignment, not implementation or
owner acceptance of a runtime. See Evidence.md for verification and remaining
integration conditions.
