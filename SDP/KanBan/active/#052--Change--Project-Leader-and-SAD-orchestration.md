# Project Leader and SAD orchestration — ProjectGovernance handoff

| Field | Value |
| --- | --- |
| id | KB-SDP-052 |
| project | SDP |
| type | Change |
| CardState | in-progress |
| created | 2026-10-09T12:08:50.581137+00:00 |
| source | Owner Session0011 T004 role definition and T005 KanBan handoff request |
| Systems | SDPTOOL |
| next_review | PGL3 responsive coordination and durable reporting; retain PGL2 closure/SAD restrictions |
| PlanId | PLAN-SDP-0018 |
| tags | project-leader, steering, SAD, master, MCP, app-server, coordination |

## Outcome and owning workstream

Make Project Leader the normal owner-facing coordination agent, able to supervise
multiple bounded SAD/Master assignments while remaining available to the owner.
The owner can also talk directly with SAD/Master; consequential feedback must reach
shared records and the leader without requiring manual chat-to-chat copying.

This is a distinct runtime/UI increment related to
[KB-SDP-038](%23038--Proposal--Codex-app-server-development-client.md).
Its owning workstream is Session0007 / ProjectGovernance. The registration branch's KB038 copy was backlog, while the owning worktree has
it active. This imported card now belongs to that active workstream; its original
handoff does not reset KB038 or duplicate the controller implementation.

## Authoritative handoff and delivered foundation

Read [Project-Leader-Handoff.md](../../04--Design/SDPTool/ProjectGovernance/Project-Leader-Handoff.md)
for the full owner direction, source observations, required increments and acceptance
scenarios. Keep runtime design details there and in the owning design/plan, rather
than duplicating them in this card.

[MAINT-SDP-0016](../../Maintenance/PLR1/Plan.md) delivered local skills and the shared
[role contract](https://github.com/Hans-Einar/SDP/blob/edfac185bb11b06ae9c6dd39a1809d166ad9541f/Skills/sdp/references/roles.md): new Project Leader, strategic
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

2026-10-10: Session0007 T014 owner explicitly selects KB052 after trying the lab.
Imported from blueprint commit edfac185bb11b06ae9c6dd39a1809d166ad9541f; retain original registration history.
[Orchestration design](../../04--Design/SDPTool/ProjectGovernance/Orchestration.md)
and [PLAN-SDP-0018](../../05--Implementation/SDPTool/ProjectGovernance/Plan.md)
now own the staged delivery. PGL1 is the first selected runtime increment.

2026-10-10: PGL1 implementation starts within the scoped preflight boundary.

2026-10-10: PGL1 delivers registered-work overview in readable and JSON forms,
verified on existing owner lab state without mutation. See [evidence](../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGL1.md).
PGL2 assignment preparation/revision guards are next. The card stays in-progress:
Project Leader dialogue, multi-assignment execution, durable inbox, SAD integration,
phase continuation, full dashboard and online Steering evidence remain outstanding.

2026-10-10 / T015: Owner continues. PGL2-C1 delivers immutable SAD/Master preparation,
exact source/generation checks, local scope refusal, retry recovery and stale
inspection. [Evidence](../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGL2-C1.md)
is CLI/service-level with fixture artifacts. PGL2 remains in-progress for actual
disposition/provider/revision/retirement gates; PGL3–5 and owner acceptance remain.

2026-10-18 / T018: PGL2-C2a delivers stable-ID package revision and terminal
retirement with preserved receipts and local scope release.
[Evidence](../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGL2-C2a.md)
records scoped tests and the separate compiled CLI trial. CardState remains
in-progress; actual disposition/provider checks and PGL3–5 remain unfinished.

2026-10-18 / T020: PGL2-C2b delivers read-only checks against canonical Blueprint
and live ModelGovernance captures, exact task/scope and preparation revision.
[Evidence](../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGL2-C2b.md)
records source-drift refusals, compiled CLI/schema checks and the isolated lab preview.
No owner acceptance is inferred. CardState stays in-progress for trusted disposition,
dispatch and PGL3–5. EVT-KB-SDP-000372 links this review to the active plan delivery.

2026-10-18 / T022: Owner authorizes continuing while reading the walkthrough.
PGL2-C2c delivers the local controller decision adapter with immutable revision-bound
records, supersession and current/stale applicability. [Evidence](../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGL2-C2c.md)
uses scripted fixture decisions, not actual owner approval. Final conversation/UI and
same-user identity authentication are not claimed. CardState remains in-progress;
dispatch prerequisites and PGL3–5 remain. EVT-KB-SDP-000373 records this update.

2026-10-19 / T024: Owner requests continued work toward Session completion.
[C2d evidence](../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGL2-C2d.md)
delivers exact approved Master activation, live run authority and a real isolated Go
implementation trial. External behavioral tests pass and protected files remain intact.
Native report submission was refused under never; no complete reporting, acceptance
or independent review is claimed. Scope closure/replanning, SAD dispatch and PGL3–5
remain open. CardState stays in-progress. EVT-KB-SDP-000374 records this boundary.
