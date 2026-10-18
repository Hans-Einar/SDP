# SDP development client using Codex app-server

| Field | Value |
| --- | --- |
| id | KB-SDP-038 |
| project | SDP |
| type | Proposal |
| CardState | backlog |
| created | 2026-09-28T15:10:13.034937+00:00 |
| source | Owner conversation 2026-09-28: register app-server-based system development and assess SDPTool MCP |
| next_review | Start the bounded PGI1-M1 core slice under planned PLAN-SDP-0018; design delivered in PGD1 |
| tags | codex, app-server, client, MCP, supervision, blueprints, observability |
| Systems | SDPTOOL |
| PlanId | PLAN-SDP-0018 |

## Current study entry point

[PLAN-SDP-0007 — shared StudyPlan](../../02--Requirements/RoutineGovernance/StudyPlan.md)
now coordinates governance, SDPTool MCP and Codex app-server research. It links
the existing RGS1 results and owns the result register for subsequent studies.
The selected studies are complete; this capability card is back in backlog. Read it before selecting further
study work or a successor DesignPlan. Earlier RGS1 work remains delivered under
PLAN-SDP-0006; this link does not reopen that completed plan.

## Owner intent

Explore using Codex app-server to develop systems with Codex through an
SDP-aware client. The owner needs coherent project dialogue, bounded assignments
and visible work state, including agent roles, skills, MCP activity and the
position in the work tree. A right-side panel is one presentation option, not
the whole capability. Preserve the ability to use ChatGPT Pro authentication.

The owner normally talks to Steering/Project Manager. A fresh Master coordinates
each bounded task, with Worker, independent Reviewer and appropriate Verifier
children. Carry forward the role clarification in the
[RGS1 study](../../02--Requirements/RoutineGovernance/Study.md); do not make the
project-facing conversation responsible for remembering every execution detail.

This request authorizes backlog registration, not a Codex fork, client
implementation, host configuration or automatic agent delegation.

## Proposed responsibility split

Recommendation: keep both app-server and an SDPTool MCP adapter, serving different
purposes. These are proposed boundaries, not delivered interfaces.

| Component | Responsibility |
| --- | --- |
| Codex app-server | Agent conversations/threads, model execution, authentication, host approvals and activity events |
| SDPTool core | Project discovery, current process/model context, selected work and validated routine/evidence operations as those capabilities are implemented |
| SDPTool MCP adapter | Let Codex agents query that context and invoke bounded SDPTool operations; reuse the same core as CLI |
| SDP-aware client | Owner dialogue, project/task navigation, role/thread association, observation and decision presentation |

App-server is the interface our client uses to operate Codex. SDPTool MCP is the
interface Codex agents use to access SDP. Neither replaces the other. A client
may also use an appropriate SDPTool snapshot/event API directly; it need not
pretend to be an agent or route every UI read through a model/MCP call.

Do not duplicate business rules in the client or MCP adapter. Keep durable SDP
routine/task identities separate from Codex thread/session identities. Link them
explicitly so replacement sessions and retries preserve the same work history.
MCP availability alone does not enforce procedure compliance or intercept arbitrary
shell/filesystem changes; retain RGS1's measured host-integration boundaries.

## Initial investigation and candidate scope

- Compare a separate Kitty observer, a small Codex TUI extension and a custom
  app-server client. Select presentation after a bounded prototype; do not assume
  a permanent Codex fork is necessary.
- Demonstrate project-facing supervision handing a bounded task to a fresh Master,
  with explicit context/blueprint references, child work and returned evidence.
  Respect actual host nesting/concurrency limits; a role tree need not map to
  unlimited recursively spawned agents.
- Correlate observed tool/agent events with SDP-owned plan, milestone, routine and
  gate state. Codex conversational plan completion is not SDP acceptance.
- Distinguish available skills, explicit invocation/loading and declared task use;
  unknown activation stays unknown. Separate available MCP tools from active calls
  and their results. Do not label instruction compliance as directly observed.
- Establish connection/session ownership and recovery. Do not assume passive
  attachment to an arbitrary running TUI, safe multi-client control or reliable
  event replay without testing the selected app-server version.
- Reuse Codex-managed ChatGPT sign-in for subscription access. Verify the active
  auth mode and applicable usage limits. Do not silently fall back to separately
  billed API-key use. Client/server transport authentication is a separate concern.
- Keep private project context scoped to the selected project and avoid exposing
  credentials, complete transcripts or unnecessary source material in the monitor.

## Existing findings and boundaries

[Official app-server documentation](https://learn.chatgpt.com/docs/app-server)
describes ChatGPT-managed authentication as well as API-key authentication, and
agent/tool event interfaces. ChatGPT Pro can be used; app-server does not require
switching to API-key billing. Applicable subscription limits still apply; see
[pricing](https://learn.chatgpt.com/docs/pricing). This records documentation checked
on 2026-09-28, not a tested local app-server session or account configuration.

[RGS1 Evidence E11](../../02--Requirements/RoutineGovernance/Evidence.md#e11--codex-panel-follow-up)
records existing Codex inspection commands and the limits of the panel inquiry.
The SDL assignment compiler, generic SDP routine engine and client proposed here
must not be described as existing functionality.

## Relationships and next action

- [KB-SDP-036](%23036--Proposal--Request-routines-and-process-control.md) owns request
  routing, routine selection and entry guidance.
- [KB-SDP-037](%23037--Proposal--Observable-routine-state-machine.md) owns durable
  routine execution, MCP-facing process operations and observation contracts.
  This card owns the Codex client/integration experience and consumes those
  contracts; it does not create a second workflow engine or duplicate MCP project.
- [KB-SDP-031](../completed/%23031--Study--SDL-assignment-bundles-and-blueprints.md)
  and planned [BP2](../../04--Design/SDPTool/Blueprints/Plan.md) own assignment context
  and blueprint design. Do not activate BP2 merely by linking it.

Next action: select a proportionate DesignPlan or bounded investigation together
with KB036/037, resolving presentation and lifecycle ownership. No new software
System or source directory is adopted by this registration.

A selected pilot should prove one small end-to-end assignment: context resolution,
supervision-to-Master handoff, permitted child execution/review, visible evidence
and a correct return disposition. Include scope escalation, stale candidate,
disconnection/reconnect, restart and missing-routine cases. Demonstrate actual
ChatGPT sign-in and telemetry without manufacturing owner approvals. A UI mock or
successful MCP connection alone does not complete the capability.

## Worklog and revisions

| Time (RFC3339) | Actor / event | Work, finding or decision | Evidence / remaining work |
| --- | --- | --- | --- |
| 2026-09-28T15:10:13.034937+00:00 | codex; EVT-KB-SDP-000215 | Registered the requested app-server client investigation and complementary SDPTool MCP boundary | Backlog only; plan, prototype and implementation unselected |

2026-09-28T15:20:29.600210+00:00: EVT-KB-SDP-000219 — linked the shared StudyPlan and result index; remaining research planned, CardState unchanged.

2026-09-28T21:31:58.382764+00:00: EVT-KB-SDP-000222 — owner starts RGS2 studies; backlog -> active/in-progress, implementation remains unselected.

2026-09-28T21:40:44.055216+00:00: EVT-KB-SDP-000225 — RGS2-C-M1 [app-server study](../../02--Requirements/RoutineGovernance/AppServer-Study.md) delivered with generated-schema evidence; report independently reviewed, synthesis/closeout continues.

## RGS2 study outcome

[PLAN-SDP-0007](../../02--Requirements/RoutineGovernance/StudyPlan.md) is complete.
Its result register links the governance, MCP, app-server and synthesis reports.
[Independent review and verification](../../02--Requirements/RoutineGovernance/RGS2-Verification.md)
support study delivery, not implemented runtime/client behavior. The recommended
next action is selection of the bounded DesignPlan described in
[Synthesis.md](../../02--Requirements/RoutineGovernance/Synthesis.md).

This card was active/in-progress during the authorized studies and now returns
to backlog because its proposed capability is not implemented. The completed
study plan remains linked as evidence; no implementation plan is activated here.

2026-09-28T21:42:41.640062+00:00: EVT-KB-SDP-000228 — selected RGS2 research complete; active/in-progress -> backlog for unimplemented capability, with concrete study results and next-plan brief.

## Ecosystem modeling outcome — 2026-09-29

[PLAN-SDP-0008](../../03--Architecture/Ecosystems/Plan.md) adds a [bounded CodexClient architecture](../../SDL/ProjectGovernance/CodexClient/README.md) and registered SDL model. It applies RGS2 findings as proposed responsibilities and illustrative contracts. The implementation remains unselected; this backlog card remains open. Select a bounded DesignPlan against the shared catalog rather than creating a rival process engine.

## Session proposal — 2026-09-30

[KB042](../backlog/%23042--Proposal--Goal-oriented-sessions-and-roadmaps.md) and
[proposed SDL Session](../../Sessions/session-%230001--SDL_expansion.md).

Client should capture submitted/steered prompts and completed visible final responses with thread/turn/item identity; preserve reconnect deduplication and distinguish observed activity from routine compliance.

EVT-KB-SDP-000245: Existing CardState and execution selection unchanged.

## ProjectGovernance workstream — 2026-10-03

The owner prioritizes starting this card and asks which capabilities should
precede it. Owner scope: ProjectGovernance covers method, methodology and tooling
for collaboration between the owner and Codex agents, including KB036/037/042/043.
Another agent is progressing ModelGovernance; it retains model artifact, revision,
checkpoint, merge and promotion ownership. This work references that contract
without taking over its implementation.

[Session 0007](../../Sessions/session-%230007--Project_governance.md) records the
assessment, recommended route and per-turn continuity for this goal. Reuse RGS2
instead of repeating the broad study. Recommendation: first design the minimum
KB036 routing contract with KB037 durable run/transition state and KB042 Session
identity; then exercise it through one bounded KB038 client workflow. The client
supplies the controlled request/capture surface, while SDPTool owns shared process
rules. Full completion of KB036 is not a prerequisite to starting client work.
Capture reliable activity identity/timing in the pilot so KB043 can later derive
its timeline; the complete timeline renderer is not an entry dependency.

The initial pilot can use explicit assignment/context/evidence references without
waiting for ModelGovernance, KB050 semantic blueprints or KB040's full KanBan TUI.
The shared DesignPlan is the next proposed bounded step. This turn delivers
dependency analysis, not a client implementation or activation of all five cards.
CardState remains backlog; PLAN-SDP-0007 remains completed study evidence until a
successor plan is actually registered.

2026-10-03T09:19:45.269241+00:00: EVT-KB-SDP-000286 — Record owner ProjectGovernance scope, KB038 priority and Session0007 dependency assessment; backlog state retained during preparation of the shared design.

2026-10-03T21:05:00.519753+00:00: EVT-KB-SDP-000287 — Owner continues Session0007 S2; activate bounded PGD1 design under [PLAN-SDP-0017](../../04--Design/SDPTool/ProjectGovernance/Plan.md). Related capability cards remain backlog.

## PGD1 shared design delivery — 2026-10-03

[Session 0007](../../Sessions/session-%230007--Project_governance.md) T002 records
the owner-selected design continuation. [PGD1 contract](../../04--Design/SDPTool/ProjectGovernance/Design.md)
and its acceptance cases now specify the minimum routing/run/Session/client
boundary. [PLAN-SDP-0018](../../05--Implementation/SDPTool/ProjectGovernance/Plan.md)
is the planned implementation handoff. Codex 0.160.0 schema export was inspected;
no live client, routine enforcement or integrated model execution is delivered.
The shared pilot covers a bounded contribution; broader card acceptance remains
open. ModelGovernance continues independently.

2026-10-03T21:14:27.153979+00:00: EVT-KB-SDP-000288 — PGD1 design delivered; return KB038 from in-progress to backlog for planned PGI implementation, preserving the completed design and unimplemented capability.

## Session0011 T004 — Project Leader / SAD handoff

The owner refines normal dialogue ownership to Project Leader, with separate
strategic Steering, SAD model design and implementation Masters. Read the
[role and runtime handoff](../../04--Design/SDPTool/ProjectGovernance/Project-Leader-Handoff.md).
It records current controller/MCP limits, multiple idle/resumable assignments,
direct owner feedback, revision-bound design approval and durable leader notification.
MAINT-SDP-0015 delivers local skills; runtime work remains owned by Session0007.
This branch's backlog state is a historical projection: the governance worktree has
KB038 active. Reconcile this additive handoff there without resetting its lifecycle.

Session0011 T005: [KB-SDP-052](%23052--Change--Project-Leader-and-SAD-orchestration.md)
now tracks this runtime/UI handoff explicitly. Review it during the owning
ProjectGovernance planning reconciliation; no current phase or card state is reset.
