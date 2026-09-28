# Request classification and mandatory routine selection

| Field | Value |
| --- | --- |
| id | KB-SDP-036 |
| project | SDP |
| type | Proposal |
| CardState | in-progress |
| created | 2026-09-27T23:51:37.811184+00:00 |
| source | Owner conversation 2026-09-28: mandatory routines for every request and scope change |
| next_review | At selection/start of the shared RGS2 StudyPlan; review its result register at study milestones |
| tags | workflow, governance, skills, request-routing, scope-control |
| Systems | SDP, SDPTOOL |
| PlanId | PLAN-SDP-0007 |

## Current study entry point

[PLAN-SDP-0007 — shared StudyPlan](../../02--Requirements/RoutineGovernance/StudyPlan.md)
now coordinates governance, SDPTool MCP and Codex app-server research. It links
the existing RGS1 results and owns the result register for subsequent studies.
The owner has authorized the studies; this card is active/in-progress. Read it before selecting further
study work or a successor DesignPlan. Earlier RGS1 work remains delivered under
PLAN-SDP-0006; this link does not reopen that completed plan.

## Need and owner direction

Established routines become ineffective when agents and owners must remember to
invoke them. Enthusiasm, interrupted attention, context changes and scope drift
can bypass existing KanBan, planning and review work or cause another competing
routine to be invented. The process must carry this memory rather than depend on
the owner recalling the correct procedure in each conversation.

The owner requests a mandatory entry routine for every request: classify intent,
look for an applicable established routine and identify the required path before
acting. Apply this again when a new message changes scope during ongoing work.
Changes to SDL/SDUI models and product code must be covered by KanBan and a plan.
Process/template/installer cleanup must have proportionate Maintenance planning.
The owner asks the agent to be more consistent about these constraints, including
briefly explaining the selected category and routine when action is involved.

This card captures the direction and proposed acceptance criteria. It does not
claim a new skill, enforced hook or released workflow already exists.

## Existing authority and the concrete recent example

[The SDP entrypoint](../../../Skills/sdp/SKILL.md) already routes role selection;
[KanBan](../README.md), [Plans](../../ProjectManagement/Plans.md) and the role
skills already define much of the work. Prefer strengthening this shared entry
and referring to existing procedures over adding a competing workflow hierarchy.

The recent template consolidation was recorded under
[MAINT-SDP-0009 / TS1](../../Maintenance/TS1/Plan.md). Go-only removal and release
were recorded in [MAINT-SDP-0010 / RP3](../../Maintenance/RP3/Plan.md), including
the owner's scope and version corrections. The gap is not an absence of these
plans: the request-to-routine/plan connection needs to be visible, repeatable and
checked when scope changes. Do not rewrite this history as undocumented work.

## Proposed first category set

These are candidates for study and trial, not newly adopted process rules.
Categories may combine; the resulting obligations must compose without creating
one new plan or card for every message.

| Request category | Candidate existing route |
| --- | --- |
| Explanation or status | Read current facts and answer; no mutation or unnecessary plan |
| Idea or new requirement | Find related cards, capture owner intent and reconcile overlap |
| Study or unresolved decision | Existing study/architecture route; bound unknowns and deliver a decision basis |
| Planning or reprioritization | Planning skill; select/revise the linked plan, optional Scrum/Sprint |
| Implementation or defect correction | Applicable card and plan before model/code mutation; analysis where the cause is unknown |
| Maintenance, migration or cleanup | MaintenancePlan covering affected files, preservation and verification |
| Review or verification | Appropriate independent review/evidence route; findings return to tracked work |
| Integration, release, installation or upgrade | Existing release/install routine, exact candidate and authorization boundaries |
| Continuation, correction, pause or cancellation | Reuse the active assignment; check scope, evidence and state changes |

## Proposed routine contract

1. Classify the request using current project authority and active work. Several
   categories or a non-project/out-of-scope result are valid. Do not force every
   question into an implementation workflow.
2. Locate the authoritative matching routine and current version/status. Reuse a
   matching active card/plan when it covers the work; avoid duplicate proposals.
3. State the route briefly for substantive work: category, routine, card/plan,
   next permitted step and any material gap. Do not require the owner to remember
   or supply these identifiers. Routine selection is not another approval prompt.
4. Before mutation, check prerequisites and record coverage. For code/model work,
   select the card and plan first. For a scope change, assess whether the current
   plan can be revised or separate work is needed; do not silently continue under
   stale scope or abandon the previous objective.
5. If no routine fits, identify the gap and register it. Do not invent a durable
   replacement routine implicitly. Clarify only the material unresolved decision;
   continue independent already-authorized work where possible.
6. Verify the applicable completion conditions and update existing records before
   saying done. Reclassify at milestones and new substantive input, not only at
   session startup. Reuse loaded unchanged instructions where practical.

Keep a routine record small: identity, purpose, matching criteria, required
inputs, current authority, allowed steps, completion evidence, exception handling
and supersession. Link existing skills/plans instead of duplicating them. Decide
where this registry belongs during planning; no new directory is assumed here.

## Activation and enforcement questions

- Determine whether to extend the mandatory SDP entrypoint or introduce an
  sdp-routine skill beneath it. A discoverable SKILL.md alone is not proof that
  an agent actually loaded or followed it.
- Separate host/session instruction activation, agent instruction compliance and
  programmatic checks. Establish what each supported host can actually enforce;
  do not promise that prose can intercept every request or prevent every write.
- Consider a later SDPTool/MCP preflight that returns the selected routine,
  linked work and missing prerequisites. A machine-checkable transition can
  strengthen controls but does not guarantee that arbitrary shell/file writes
  elsewhere are intercepted. Treat this as a separately bounded implementation.
- Cover compound requests, ambiguous categories, existing authorization, scope
  changes, missing routines, resumption after compaction and competing versions.
- Preserve proportionate work: direct factual questions should stay easy; a
  matching active plan should not be recreated, and explicit owner authorization
  should not trigger repeated permission requests.

## Owner extension — observable execution

The owner also requests an MCP-connected state machine and a live side-by-side
terminal flowchart, with drill-down into routines, approvals and agent assignment
checklists. [KB-SDP-037](%23037--Proposal--Observable-routine-state-machine.md)
now owns that runtime/observer capability. This card retains the entry/routing
contract and initial skill-process Maintenance scope. Align the routine identity,
version and evidence contracts; do not grow the initial Maintenance task into a
new workflow engine without a separate selected plan.

## Next action and completion criteria

The completed [RGS1 study](../../02--Requirements/RoutineGovernance/Study.md)
now recommends selecting a bounded shared DesignPlan with KB-SDP-037 for the
request/routine contract. Entry-skill Maintenance can then implement its agreed
portion without silently adopting the whole engine. Trial actual routing before
changing distributed skills/templates. No implementation has been selected.

Acceptance should demonstrate representative requests from the table, including
"also remove PowerShell" during release preparation, template restructuring,
"continue", a pure status question and an unrelated factual request. Evidence
must show when routing was performed, which authority/card/plan was reused or
revised, what gaps blocked dependent mutations and how scope/closeout was kept
current. A declarative instruction or file-presence check alone is insufficient.

Related: [KB-SDP-029: typed plans](../completed/%23029--Proposal--Typed-plans-and-planning-skill.md),
[KB-SDP-027: skill activation](../completed/%23027--Study--Skills-review-and-project-activation.md),
[KB-SDP-031: assignment context](../completed/%23031--Study--SDL-assignment-bundles-and-blueprints.md).
No reopening of completed cards or MCP implementation is implied by this link.

## Worklog and revisions

| Time (RFC3339) | Actor / event | Work, finding or decision | Evidence / remaining work |
| --- | --- | --- | --- |
| 2026-09-27T23:51:37.811184+00:00 | codex; EVT-KB-SDP-000207 | Captured owner direction and proposed routine/category/activation contract during RP3 closeout | Registration only; select the MaintenancePlan before implementation |
| 2026-09-28T00:04:17.612149+00:00 | codex; EVT-KB-SDP-000209 | Linked owner extension for durable execution state and live visualization to KB-SDP-037 | Entry/routing scope retained; runtime/observer requires separate design and implementation planning |
| 2026-09-28T00:43:32.046442+00:00 | codex; EVT-KB-SDP-000212 | Linked completed RGS1 research and aligned the next-plan recommendation with KB-SDP-037 | Remains backlog; 16 routine families are proposals, not adopted policy |

## Study findings to retain in the selected scope

RGS1 identifies stale inventory references in the current skill source map and
SDPTool contract, plus verified stale external PR status in sampled projects.
See its Evidence.md for exact sources and limits. The selected future work should
cover authority reconciliation and discovery of missing/incompatible procedures,
rather than papering over these cases with another generic instruction. These
findings are captured here; this study has not repaired those sources or authorized
changes in external projects.

## Owner follow-up: responsibility for routing

The owner clarifies the intended agent hierarchy: the project-facing Steering/
Project Manager agent has primary responsibility for request classification,
procedure coverage and project scope; each bounded task starts with a fresh Master
that coordinates Worker/review roles and enforces the assignment boundary.
Workers still report missing context, procedure gaps and scope conflicts. The
[RGS1 study](../../02--Requirements/RoutineGovernance/Study.md) records the detailed
responsibility split and handoff implications. This does not adopt a new skill,
automatic spawning rule or grant publication/owner-decision authority to agents.

2026-09-28T15:00:58.823677+00:00: EVT-KB-SDP-000214 — recorded owner clarification and proposed design implications; CardState remains backlog.

2026-09-28T15:20:29.600210+00:00: EVT-KB-SDP-000217 — linked the shared StudyPlan and result index; remaining research planned, CardState unchanged.

2026-09-28T21:31:58.382764+00:00: EVT-KB-SDP-000220 — owner starts RGS2 studies; backlog -> active/in-progress, implementation remains unselected.
