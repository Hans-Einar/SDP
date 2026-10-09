# Project Leader and SAD orchestration — ProjectGovernance handoff

| Field | Value |
| --- | --- |
| id | KB-SDP-052 |
| project | SDP |
| type | Change |
| CardState | backlog |
| created | 2026-10-09T12:08:50.581137+00:00 |
| source | Owner Session0011 T004 role definition and T005 KanBan handoff request |
| Systems | SDPTOOL |
| next_review | ProjectGovernance / Session0007 planning reconciliation, before extending the owner lab |
| tags | project-leader, steering, SAD, master, MCP, app-server, coordination |

## Outcome and owning workstream

Make Project Leader the normal owner-facing coordination agent, able to supervise
multiple bounded SAD/Master assignments while remaining available to the owner.
The owner can also talk directly with SAD/Master; consequential feedback must reach
shared records and the leader without requiring manual chat-to-chat copying.

This is a distinct runtime/UI increment related to
[KB-SDP-038](%23038--Proposal--Codex-app-server-development-client.md).
Its owning workstream is Session0007 / ProjectGovernance. This branch's KB038 copy
is backlog, while the inspected owning worktree has it active: reconcile the handoff
without resetting that workstream or duplicating its controller implementation.

## Authoritative handoff and delivered foundation

Read [Project-Leader-Handoff.md](../../04--Design/SDPTool/ProjectGovernance/Project-Leader-Handoff.md)
for the full owner direction, source observations, required increments and acceptance
scenarios. Keep runtime design details there and in the owning design/plan, rather
than duplicating them in this card.

[MAINT-SDP-0015](../../Maintenance/PLR1/Plan.md) delivered local skills and the shared
[role contract](../../../Skills/sdp/references/roles.md): new Project Leader, strategic
Steering, Architect as SAD, and Master as the owner of a bounded implementation
Session/ImplementationPlan. Metadata/distribution checks pass; this does not deliver
runtime orchestration, online installation or automatic notification.

## Required scope

- Project Leader prepares bounded SAD assignments, then implementation assignments
  pinned to reviewed NOW/TARGET/task/blueprint revisions and linked cards.
- SAD creates model WORK/TARGET with ModelGovernance and generates previews/blueprints;
  owner design feedback and approval bind exact revisions. Changed design invalidates
  the affected handoff rather than silently widening implementation scope.
- Masters own Sessions and phased plans, report phase evidence, and can remain idle
  for authorized continuation. Preserve direct owner access and leader visibility.
- Extend privileged controller/MCP capabilities deliberately; current worker MCP
  does not expose launch/resume rights. Reuse actual runtime services and preserve
  per-assignment authority, optimistic revision checks and retry safety.
- Durable event delivery and a responsive leader conversation: correlate reports,
  retain them while the leader is idle/busy, and process them without duplicate
  launch or losing newer owner steering. MCP alone is not a wake-up scheduler.
- A TUI/client overview exposes agent hierarchy, role, Session/phase, observed state,
  tools and child activity with freshness/unknowns. Select a supported UI approach;
  a Codex TUI plugin/fork is not assumed available or selected.
- Reconcile shared-work claims and ID allocation across worktrees. Use the existing
  preflight until stronger coordination is implemented; do not claim global exclusion
  from branch-local records. Package/test Steering in the selected online host as a
  separately verified integration, reusing canonical skill sources.

## Acceptance and next action

Before implementation, the ProjectGovernance agent should reconcile the handoff with
its actual current design and PLAN-SDP-0018, then select bounded milestones or a linked
successor plan. Exercise all six scenarios in the handoff: parallel disjoint work,
shared-scope refusal, changed design, phase/idle continuation, durable notification,
truthful dashboard and remote evidence boundaries. Host completion must not imply
assignment acceptance. No main merge, release or native XFMD changes are authorized
by this registration.

## Worklog

2026-10-09: Owner requests a concrete KanBan entry for the handoff. Registered as
KB-SDP-052 after inspecting IDs in all nine known local worktrees. Session0011
records this transfer; its KB051 product implementation remains behind the red
preflight. The existing duplicate KB051 is recorded separately and not renumbered.
