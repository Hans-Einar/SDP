# Project Leader and SAD — handoff to ProjectGovernance

Owner refinement, Session0011 T004, 2026-10-09. Role guidance is delivered under
[MAINT-SDP-0015](../../../Maintenance/PLR1/Plan.md). Runtime changes below belong to
the existing ProjectGovernance workstream, Session0007 / KB-SDP-038. Its current
Design.md and ImplementationPlan remain untouched in this branch; reconcile this
handoff into them before implementation. No second controller implementation is
selected here. [Shared role contract](../../../../Skills/sdp/references/roles.md).

Tracked by [KB-SDP-052](../../../KanBan/backlog/%23052--Change--Project-Leader-and-SAD-orchestration.md),
linked to the owning KB038 workstream. Registration does not activate runtime work.

## Owner-selected operating model

The owner primarily talks with **Project Leader**. Steering is strategic support,
possibly in ChatGPT online with GitHub access; it supplies direction/constraints
and actual owner dispositions, not automatically executable orders. Project Leader
coordinates multiple SAD/Master assignments and remains responsive to the owner.

SAD (Software Architecture and Design) uses the existing Architect skill to develop
SDL/SDUI in a ModelGovernance WORK selected as TARGET. It produces runnable previews
only where supported, validated model evidence and generated preliminary blueprints.
The owner can talk directly with SAD and inspect TARGET in XFMD. Review binds exact
NOW/TARGET/task/blueprint revisions. Project Leader then supplies the accepted package
to Master, which creates/maintains its Session and ImplementationPlan with phases and
milestones. Master reports per phase; existing granted continuation may cover later
phases. New scope or reserved approval goes back to the owner. Completed phase work
can leave the Master idle and resumable; idle is not acceptance or assignment closure.

Project Leader should normally commission design rather than author all model/code
changes itself. Direct owner contact with Master remains available. Consequential
feedback is captured once in shared records, delivered to the leader, and revalidated
against dependent assignments. Do not require the owner to relay three chat histories.

## Current implementation observed

Inspected /home/warloc/git/SDP-project-governance, HEAD 0b0e82c, clean when read.
Session0007 T013 identifies PGI3-C6 evidence and an owner lab with feedback pending.

- governancecontroller: bounded existing authorized Maintenance run, one fresh Master
  turn; steering/interruption, approvals, correlated recovery/reattachment and explicit
  skill input. Its README states no native child spawning is implemented.
- governancemcp: context_read, route_evaluate, run_read and attempt-bound events_read /
  run_transition. Transitions limited to report/fail/scope-change. No create/resume/
  launch/complete/owner-decision operations are exposed to this worker adapter.
- governancehost: thread/start and turn/start/steer/interrupt, observed host identities.
  A transport capability is not a proven multi-Master scheduler or parent wake-up.
- The broader Steering/PM conversation, native child attribution and automatic
  routine compliance are not established by these candidate tests.

These facts are read-only source evidence, not a claim that the owner has accepted
PGI3 or that future commits retain the same limitations. Reinspect before selection.

## Required runtime increments for the proposed workflow

### Privileged project-leader operations

Keep worker MCP scope restricted. Add a separately authorized Project Leader adapter
or role-bound capabilities for assignment preparation, dispatch, inspection, targeted
messages and authorized continuation. Reuse the controller and core state machine;
MCP cannot grant controller rights merely because a payload says role=ProjectLeader.
Use stable operation identity, expected revisions, allowed scopes, concurrency limits
and explicit launch provenance. A request must not launch a second Master after an
uncertain retry. Reconcile the original attempt before replay.

Track assignment/run/attempt/thread separately from human role labels. A Master can
retain its logical Session across idle turns, reconnection or a replacement host
thread, with explicit mappings and preserved history. Do not force one perpetual
model context or one new agent per phase. SAD and implementation get separate bounded
assignments; their dependencies pin accepted design revisions.

### Event delivery and a responsive owner conversation

Proposed path: worker/app-server event -> trusted controller -> durable correlated
inbox/state -> notification to UI and a scheduled short Project Leader turn.
MCP may expose reports or event reads, but does not itself guarantee wake-up of the
leader. Choose a supported host mechanism and prove it. A controller event loop must
continue even when the leader is idle or the owner is talking to another agent.

Do not inject simultaneous conflicting turns into one thread. Queue/coalesce status
updates with stable event IDs/cursors; prioritize owner input, serialize decision
turns, and preserve pending reports across disconnect/restart. Store observation
freshness/gaps. Reconcile duplicate/out-of-order terminal messages and status changes
before dispatching subsequent work. Test a delayed report against newer owner steering.
The model performs decisions; it is not the always-running background event loop.

### UI requirement, not framework selection

The owner's desired TUI or supported host extension should keep the Project Leader
conversation usable alongside a panel of Masters/SAD assignments: current role,
Session, phase/milestone, assigned task, running/idle/waiting state, last observation,
known tool activity and mapped child agents. Separate declared skill use, explicit
skill-input acknowledgment and verified behavior. A stale heartbeat cannot be shown
as live work. Selecting a worker should allow direct dialogue without losing the
leader overview, with permission/context preserved.

No suitable stock-TUI plugin surface has been verified in this handoff. Compare
existing app-server client/lab reuse with a small dedicated UI before proposing a
Codex TUI fork. This does not authorize XFMD changes or implementation in this branch.

### Shared coordination and model revision gates

The preflight found a real duplicate KB051 and conflicting dispatch/output edits.
A project leader personality cannot make branch-local registries globally unique.
Specify a shared allocation/claim authority across worktrees and recovery of abandoned
claims, retaining project history. Do not use machine-local checkout state as proof
of cross-checkout exclusion. Start with explicit human/coordinator sequencing until
the shared capability exists; do not advertise atomic reservations prematurely.

Owner-approved design, tool validation and implementation evidence are separate.
Any changed TARGET, task or protected boundary stales the affected approval/bundle;
pause dependent implementation for the relevant rework/disposition. Use existing
ModelGovernance, Blueprints and Traceability services, not a second status database.

## Acceptance scenarios for the owning plan

1. Two disjoint Masters run while Project Leader answers the owner; a shared-scope
   third assignment is refused or queued with a reason.
2. Owner speaks directly with SAD; accepted TARGET changes and the old implementation
   package becomes stale before a Master can continue under it.
3. A Master completes a phase and becomes idle. Leader continues the same Session only
   under current authority and candidate evidence; a gated next phase awaits disposition.
4. Completion arrives during owner dialogue or controller reconnect: it is retained,
   visible and processed once without duplicate launch or lost report.
5. Dashboard truthfully shows tools, children and skill evidence including unknown/stale
   states; process completion, assignment completion and owner acceptance stay distinct.
6. Steering sees a committed branch snapshot but not local dirty work; it reports that
   limit and supplies a revision-bound brief rather than claiming a live global view.

## Skills and online distribution

Canonical skill changes: new sdp-project-leader; clarify sdp-steering; extend
sdp-architect for SAD and sdp-master for phase reporting; shared roles.md and router.
Local discovery and development install inventory are updated together. This does
not install a plugin in an online account or prove runtime role enforcement.

Official [skills documentation](https://developers.openai.com/plugins/concepts/skills)
and [plugin packaging](https://developers.openai.com/plugins/build/plugins) describe
sharing skills with ChatGPT/Codex. Package canonical Steering guidance with its shared
references for the selected host, verify GitHub connector visibility and observed
skill loading there. Do not create a separately maintained online role definition.
No online plugin was installed or submitted by this maintenance.
