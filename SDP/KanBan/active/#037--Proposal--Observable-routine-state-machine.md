# Observable routine execution and agent work-package progress

| Field | Value |
| --- | --- |
| id | KB-SDP-037 |
| project | SDP |
| type | Proposal |
| CardState | in-progress |
| created | 2026-09-28T00:04:17.612149+00:00 |
| source | Owner conversation 2026-09-28: MCP-connected state machine and live side-by-side process visualization |
| next_review | At selection/start of the shared RGS2 StudyPlan; review its result register at study milestones |
| tags | workflow, state-machine, MCP, observability, terminal, blueprints |
| Systems | SDPTOOL, SDL |
| PlanId | PLAN-SDP-0007 |

## Current study entry point

[PLAN-SDP-0007 — shared StudyPlan](../../02--Requirements/RoutineGovernance/StudyPlan.md)
now coordinates governance, SDPTool MCP and Codex app-server research. It links
the existing RGS1 results and owns the result register for subsequent studies.
The owner has authorized the studies; this card is active/in-progress. Read it before selecting further
study work or a successor DesignPlan. Earlier RGS1 work remains delivered under
PLAN-SDP-0006; this link does not reopen that completed plan.

## Owner intent

The owner wants to observe the process while an agent works, without depending
on chat context or remembering process documents. A separate terminal in the same
Kitty window should show a flowchart of request categories and routines. Nodes
and paths change appearance as a request progresses. The observer can enter the
selected routine to follow its internal steps, prerequisites, reminders and
required approvals, then return to the enclosing flow.

The same view should follow agents executing work packages with blueprints and
show their routine checklists. A concrete pilot is changing the SDP process and
carrying that change through a new release. The process should hold durable
memory and make deviations visible, rather than merely display an agent's prose.

This is a requested capability and proposed design scope. No state-machine
service, MCP adapter or live terminal renderer is implemented by this card.

## Relationship to current work

[KB-SDP-036](../active/%23036--Proposal--Request-routines-and-process-control.md) owns request
classification, existing-routine selection and visible card/plan coverage. This
card owns executing/observing the selected routine across requests and agents.
Routine identity, versions and completion rules must be agreed across both.

[KB-SDP-031](../completed/%23031--Study--SDL-assignment-bundles-and-blueprints.md)
and its proposed [blueprint DesignPlan](../../04--Design/SDPTool/Blueprints/Plan.md)
own assignment context, NOW/TARGET models, boundaries and evidence. Reference
those artifacts; do not invent a competing blueprint schema or silently activate
that plan. The completed study remains completed.

## Candidate responsibility boundaries

These are recommendations for a DesignPlan, not an adopted architecture.

- SDPTool owns routine execution state and validates requested transitions.
  Request classification may still need agent judgment; retain its stated basis
  and unresolved ambiguity rather than presenting it as an objective fact.
- An MCP adapter exposes inspect, start/resume and transition operations over
  that core. CLI and other clients should use the same rules. MCP is an interface,
  not an automatic observer of every chat message or shell command.
- A read-only monitor obtains a current snapshot plus subsequent events. Losing
  or closing the monitor must not cancel work or change the underlying state.
  Reopening it restores actual state without reconstructing it from chat.
- A routine definition is distinct from an execution instance. Pin its version
  for each run; changing the definition must not silently rewrite a running or
  historical instance. Progress belongs to the instance, not to the diagram.

## State and evidence contract to design

Model hierarchical workflows with parallel child assignments rather than one
flat current-state value for the whole project. Keep stable request, routine-run,
step, assignment and agent-session identities and explicit parent relationships.
Changing an agent session must not lose the assignment's process state.

Candidate step states include not-started, ready, running, waiting-for-input,
waiting-for-review, blocked, failed, completed and canceled. Define allowed
transitions, prerequisites and recovery semantics before adopting names. Skipping
or accepting an exception requires its recorded reason and permitted authority;
it must remain distinguishable from successfully completing a step.

Separate an agent's progress report, independently observed tool/test results,
review disposition and owner approval. A generic complete call or a green box
must not manufacture approval or establish code conformance. Evidence should
identify the relevant commit/artifact and become stale when relevant inputs
change. Reuse existing valid owner authorization rather than prompting again at
every step. Unavailable/stale observations should appear explicitly as unknown.

Define idempotent transition requests, concurrency/revision checks and durable
recovery after disconnect or restart. Parallel agents may finish in a different
order; events must not overwrite each other's progress. Disconnection alone is
not evidence of task failure or completion.

Use existing ProjectManagement and Traceability authorities for their respective
lifecycle/evidence facts. Decide how operational execution events relate to those
records without duplicating their ownership or adding a second competing ledger.
A projection for live monitoring is not a replacement source of project truth.

## Live observer requirements

- Show overview -> selected routine -> step/child assignment, with an explicit
  path back to the enclosing request and linked KanBan card/plan/milestone.
- Distinguish available routes from the route actually taken. Show parallel work,
  dependencies, waiting reasons, last observed activity and stale/disconnected
  state. Make color optional by including labels or symbols as well.
- At a step, show its checklist, required inputs, governing instruction, evidence
  links, responsible role and any pending approval. A checklist item marked done
  must expose the basis for that status.
- Support a separate Kitty terminal as the owner's initial viewing scenario.
  Assess a portable text/TUI view and diagram rendering using existing components;
  do not select a new GUI framework or require XFMD application changes here.
- Permit following the active run or inspecting a previous run. Live progress and
  historical replay must be clearly distinguished. Start with snapshot/polling
  if adequate; choose transport after testing the actual host capabilities.

Illustrative process-change routine: classify request -> resolve related card ->
select/revise MaintenancePlan -> inspect impact -> implement selected milestone ->
verify -> review -> publish when authorized -> reconcile release/installation.
This example must be reconciled with current release rules, including correction,
failure, rework and cancellation paths; it is not a newly enforced pipeline.

## Enforcement boundary

The core can reject invalid transitions through its own API. That does not stop
an agent using unrelated filesystem or shell tools. Explicitly distinguish
visibility from enforcement, and test how requests enter the workflow in each
supported host. Any stronger mutation control needs an actual integration and
authorization design. Do not claim complete enforcement from a skill, an MCP
connection or a colored diagram alone.

## Next action and acceptance

First align the routine contract with KB-SDP-036. Select a bounded DesignPlan for
execution state, evidence and observer interfaces, followed by a proportionate
implementation plan. This runtime/UI capability exceeds the initial Maintenance
work on skill entry and classification; do not absorb it into that work silently.

A small pilot should follow one SDP maintenance-to-release routine in a second
terminal, with a worker child assignment and blueprint reference. Demonstrate
visible progress, prerequisite rejection, review/rework, existing authorization,
missing evidence, context/session restart, observer reconnect and replay. A
second child assignment tests parallel progress without shared-state corruption.
Use recorded or isolated test releases; exercising the workflow must not publish
an unrequested release. Verify behavior against persisted facts, not animations.

## Worklog and revisions

| Time (RFC3339) | Actor / event | Work, finding or decision | Evidence / remaining work |
| --- | --- | --- | --- |
| 2026-09-28T00:04:17.612149+00:00 | codex; EVT-KB-SDP-000208 | Captured the owner's MCP-connected state machine, live drill-down visualization and work-package checklist request | Linked KB-SDP-036 and blueprint design; proposals only, plan selection and implementation remain pending |

## Selected extended study — RGS1

Owner request 2026-09-28 selects an extended study covering the entire project
workflow, the named repository portfolio, distributed procedures and detection
of missing routines. [PLAN-SDP-0006](../../02--Requirements/RoutineGovernance/Plan.md)
authorizes this research only. Production runtime/observer implementation remains
unselected. Return the card to backlog after delivering the study, with remaining
capability scope visible.

2026-09-28T00:26:19.580373+00:00: EVT-KB-SDP-000210 — backlog -> active/in-progress for the selected RGS1 study.

## Extended study delivered; implementation remains open

[RGS1 study](../../02--Requirements/RoutineGovernance/Study.md) completes the
selected research under PLAN-SDP-0006. It covers the named project portfolio,
work before assignment, 16 routine families and requirements, procedure-gap
discovery, versioned project distribution, evidence/authority boundaries and
16 analytical challenge cases. Pinned source evidence distinguishes observed
facts from project assertions; no product audit or runtime trial is claimed.

Recommended next work is a bounded shared DesignPlan with KB-SDP-036, then one
Go SDPTool/CLI/MCP maintenance workflow and read-only terminal observer pilot.
BP2 remains the separate assignment/blueprint contract. This card returns to
backlog because none of the proposed engine/observer capability has been delivered.
No mandatory owner review of the study is invented as a reason to keep it active.

2026-09-28T00:43:32.046442+00:00: EVT-KB-SDP-000211 — active/in-progress -> backlog; RGS1 study complete, capability implementation unselected.

## Owner follow-up: supervision, task Masters and Codex presentation

The owner clarifies that their normal dialogue is with a Steering/Project Manager
agent responsible for the broader project and procedure coverage. A fresh Master
coordinates each bounded task and delegates implementation/review as permitted.
Record project supervision, task coordination and implementing Worker as distinct
responsibilities, with durable handoff and escalation rather than relying on the
parent chat's memory. See the dated clarification in the
[RGS1 study](../../02--Requirements/RoutineGovernance/Study.md).

The owner also asks about a Codex TUI right-side panel for skills, MCP use and the
work tree. [Evidence E11](../../02--Requirements/RoutineGovernance/Evidence.md#e11--codex-panel-follow-up)
records available command/interface building blocks and uncertainty about a
combined built-in panel. Retain separate-pane, TUI-extension and custom app-server
client alternatives. Skill availability is not skill compliance. This inquiry
records an observer option; it does not select a Codex fork or implementation.

2026-09-28T15:00:58.823677+00:00: EVT-KB-SDP-000213 — recorded owner clarification and proposed design implications; CardState remains backlog.

## Related client proposal

[KB-SDP-038](%23038--Proposal--Codex-app-server-development-client.md) now owns the
requested app-server-based development client investigation. This card retains
the shared routine/MCP/state/observer contract. The client should consume it,
with explicit links between SDP work identities and Codex threads, rather than
implement another process engine. Neither card is activated by registration.

2026-09-28T15:10:13.034937+00:00: EVT-KB-SDP-000216 — linked the new client proposal and retained shared-core ownership; backlog unchanged.

2026-09-28T15:20:29.600210+00:00: EVT-KB-SDP-000218 — linked the shared StudyPlan and result index; remaining research planned, CardState unchanged.

2026-09-28T21:31:58.382764+00:00: EVT-KB-SDP-000221 — owner starts RGS2 studies; backlog -> active/in-progress, implementation remains unselected.
