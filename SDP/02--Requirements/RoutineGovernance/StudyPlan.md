# RGS2 StudyPlan — request governance, SDPTool MCP and Codex app-server

| Field | Value |
| --- | --- |
| id | PLAN-SDP-0007 |
| project | SDP |
| state | completed |
| PlanType | RequirementPlan |
| BranchPolicy | current |
| CommitPolicy | milestone |
| Systems | SDPTOOL, SDL, SDUI |
| source | Owner request 2026-09-28; KB-SDP-036, KB-SDP-037, KB-SDP-038 |

## Purpose and authority

This is the shared **StudyPlan** and result index for request/routine governance,
an MCP connection to SDPTool, and a Codex app-server-based development client.
The owner requests one plan linked from the relevant cards, with separate study
results discoverable here. The owner subsequently authorized execution of these studies on 2026-09-28.
The authorized research is now complete; production implementation remains outside this assignment.

StudyPlan describes the purpose of this document. Its machine-readable PlanType
is RequirementPlan under the existing planning contract: the studies establish
needs, constraints, evidence and recommendations before detailed design. No new
plan enum, ledger schema or distributed process capability is introduced here.

The completed [RGS1 plan](Plan.md), PLAN-SDP-0006, remains closed. Its delivered
studies are reused as evidence; they are not copied, reclassified as unfinished,
or claimed as new work. This plan owns the remaining coordinated investigation
and synthesis. It is the current study entry point for all three cards.

## Linked work and ownership

| Card | Study responsibility |
| --- | --- |
| [KB-SDP-036](../../KanBan/backlog/%23036--Proposal--Request-routines-and-process-control.md) | Request classification, routine selection, missing-procedure handling and entry/skill responsibilities |
| [KB-SDP-037](../../KanBan/backlog/%23037--Proposal--Observable-routine-state-machine.md) | Durable routine state, evidence/gates, SDPTool MCP operations and read-only observation |
| [KB-SDP-038](../../KanBan/backlog/%23038--Proposal--Codex-app-server-development-client.md) | Codex client, supervision/task handoff, app-server integration and presentation alternatives |

Steering/Project Manager retains the project-facing dialogue and broader context.
A fresh task Master coordinates each bounded assignment and permitted Worker,
Reviewer and Verifier children. The studies must preserve this responsibility
split, actual host limits and independent-review meaning. Role names do not grant
owner decision or publication authority.

The [BP2 blueprint plan](../../04--Design/SDPTool/Blueprints/Plan.md) remains a
separate planned dependency for assignment context. This StudyPlan examines its
interface and readiness needs without activating or replacing that plan.

## Study results — start here

This table is the maintained result register. Add a result link when the report
exists, with its actual status, evidence and unresolved questions. Do not create
empty report files or label preliminary observations completed studies.

| Study area | Current result status | Available result / evidence | Result / follow-up |
| --- | --- | --- | --- |
| Request and routine governance | RGS1 delivered; RGS2-A-M2 role/routing reconciliation delivered | [Main study](Study.md), [routine catalog](Routine-Catalog.md), [scenario challenge](Scenarios.md), [evidence](Evidence.md), [source pins](Source-pins.json) | [RGS2-A-M2 findings](Study.md#10-rgs2-a-m2--operational-responsibility-and-procedure-reconciliation); same authoritative reports maintained |
| SDPTool MCP connection | RGS2-B-M1 delivered; independently reviewed | [Existing study](Study.md), especially sections 3–6; [KB038 boundary](../../KanBan/backlog/%23038--Proposal--Codex-app-server-development-client.md) | [MCP study](MCP-Study.md) and [source pins](Evidence/RGS2-MCP-source-pins.json) |
| Codex app-server/client | RGS2-C-M1 delivered; installed schemas inspected; runtime behavior untested | [Evidence E11](Evidence.md#e11--codex-panel-follow-up), [client proposal](../../KanBan/backlog/%23038--Proposal--Codex-app-server-development-client.md) | [App-server study](AppServer-Study.md), [probe evidence](Evidence/RGS2-AppServer-evidence.json) and [schema extract](Evidence/RGS2-AppServer-schema-extract.json) |
| Joint recommendation and next delivery | RGS2-D-M1 delivered; independently reviewed | Existing staged recommendation in [RGS1](Study.md), section 8 | [Synthesis](Synthesis.md) and [verification/closeout](RGS2-Verification.md) |

RGS1 contains analytical scenarios, not executed routine-engine tests. Its role
clarifications and subsequent official-documentation findings retain their dates.
App-server authentication support is a documented capability; no local account
configuration or integrated client behavior has been verified by these reports.

## Scope and research questions

### A — request/routine governance

Reuse the 16 routine families, requirements and challenge cases. Focus further
work on the owner's clarified agent hierarchy and unresolved contract choices:

- Which obligations belong to project supervision, the task Master and children?
  What evidence must cross the boundary before assignment and on return?
- How do a new request, clarification, scope change and fresh session map to an
  existing routine/task rather than create duplicate plans or cards?
- How are missing, unavailable, incompatible and conflicting procedures handled?
  Who proposes, reviews, adopts, versions and retires a routine?
- Which state belongs to management records, operational execution, observations
  and system evidence? What can be checked deterministically, and what requires
  semantic judgment or an owner decision?
- What is the smallest useful initial routine set and how do we measure overhead,
  false blocking, missed scope changes and repeated owner reminders?

### B — SDPTool MCP connection

Study an adapter over the existing Go SDPTool core, shared with CLI/client use:

- Separate read-only knowledge/resources from operations that change work state.
  Identify existing capabilities versus missing routine-engine functionality.
- Define candidate tools/resources for context, routes, selected work, assignment
  references, evidence, transitions and observations. Specify inputs, provenance,
  authorization, revision checks, idempotency and failure results at study depth.
- Compare STDIO and network transports for the actual hosts. Examine project and
  worktree identity, simultaneous clients, reconnect, version negotiation and
  installed-profile compatibility without selecting a deployment prematurely.
- Explain how MCP, skills and host integration cooperate. Document what an agent
  can bypass; never equate an MCP connection with mandatory process compliance.
- Determine what the observer/client reads directly and what an agent invokes
  through MCP. Keep business rules and durable identity in the shared core.

### C — Codex app-server and development client

Study the owner-facing development workflow rather than only a visual sidebar:

- Inspect the selected app-server version's thread, turn, collaboration, approval,
  plan and tool-event interfaces. Identify documented versus verified behavior.
- Map project supervision, fresh task Masters and child agents to supported host
  lifecycles and limits. Determine which context is supplied and how escalation,
  returned evidence and replacement sessions preserve work identity.
- Assess ChatGPT Pro sign-in and usage reporting, with explicit auth mode and no
  silent API-billing fallback. Keep transport credentials separate from account
  authentication; do not expose secrets in monitoring.
- Compare a separate Kitty pane, a small upstreamable Codex TUI change and an
  app-server client. Include lifecycle ownership, maintenance cost, portability
  and whether any claimed observer connection is actually supported.
- Identify reliable telemetry for skills, MCP calls and the work tree. Preserve
  the distinction between available, invoked, observed and declared status.

### D — synthesis and bounded next step

Reconcile the studies into one responsibility map, dependency order and decision
register. Recommend the smallest end-to-end pilot and its acceptance criteria.
Specify which remaining choices need owner disposition and which are routine
engineering decisions. Produce a concrete brief for a subsequent DesignPlan;
do not execute or silently activate that successor.

## Phases and milestone acceptance

| Phase | Milestone | Required outcome | Current state |
| --- | --- | --- | --- |
| A — governance | RGS2-A-M1 | Resolve existing RGS1 report links and preserve completed predecessor evidence | Available predecessor result; indexed during plan creation |
| A — governance | RGS2-A-M2 | Reconcile role ownership, routing/gap lifecycle and remaining decisions; update the existing reports with dated findings | Delivered |
| B — MCP | RGS2-B-M1 | Deliver MCP-Study.md with capability inventory, candidate interface/transport alternatives, authority and failure boundaries | Delivered |
| C — app-server | RGS2-C-M1 | Deliver AppServer-Study.md with versioned evidence, lifecycle/telemetry/auth limits and client-option comparison | Delivered |
| D — synthesis | RGS2-D-M1 | Deliver Synthesis.md with reconciled boundaries, unresolved decisions, one bounded pilot and linked acceptance cases | Delivered |

A's existing results are reusable now; A-M2 need not re-audit the whole repository
portfolio. B and C can investigate independent questions after the shared role
and identity baseline is clear, then reconcile in D. Separate agents may be used
only within explicit assignment and host authorization. No mandatory Scrum or
Sprint is added solely to coordinate these studies.

## Verification and research boundaries

Reports must distinguish owner intent, current project authority, documented
external capabilities, observations, design inferences and untested proposals.
Pin source/version/date for technical findings and state coverage limits.

Use read-only inspection and analytical walkthroughs by default. Any executable
probe must be bounded in the selected research assignment, isolated from product
work, and documented with its actual candidate and result. The plan does not
authorize account/config changes, billable model trials, Codex forks, installations,
product migrations, production engine changes or external publication. Do not
start real workers or claim independent review just to illustrate a diagram.

Carry forward the existing scenarios and add concrete acceptance obligations for
supervisor-to-Master handoff, MCP failures, stale evidence, host bypasses, interrupted
updates and observer recovery. Proposed acceptance tests remain proposals until
run against an implementation. Keep private source content outside public reports.

For each milestone, check report/card/plan links, update the result register,
append the actual management event and validate record consistency. Pure study
management belongs in ProjectManagement, not system implementation Traceability.
At closeout, every planned result must have a delivered or explicitly deferred
scope disposition; do not mark the whole plan completed with required study work
still open. Capability cards remain open unless their own outcome is delivered.

## Git and lifecycle policy

BranchPolicy current: use the appropriate existing working branch, currently
sdp/request-routine. Recheck status and branch before execution; never develop on
main. CommitPolicy milestone: commit each meaningful study delivery together with
its report index, card references and management history. Preserve unrelated work.
Existing phase-push/PR authorization applies; merging and release remain separate.

Study execution authorized on 2026-09-28 is complete. KB036–038 were active/
in-progress during execution and now return to backlog for their unimplemented
capabilities. Their PlanId retains this completed study; select the successor
plan before activating implementation. Retain PLAN-SDP-0006 as the completed predecessor;
do not reopen it to track this work. Study completion does not automatically
implement routines, adopt a process schema or activate BP2/client development.

## Planning record

RGS2-P0: the owner requested this shared StudyPlan on 2026-09-28. Plan creation
links existing results and assigns destinations for later studies. The new
research milestones remain planned; no study result or runtime evidence is
manufactured by creating their entries above.

## Execution log

- RGS2 start: owner authorizes the studies; KB036–038 move to active/in-progress.
  Research reports and bounded read-only probes are in scope, not product changes.
- RGS2-A-M2: role/authority matrix, handoff contents, scope-delta routing, gap
  lifecycle and four-path pilot subset delivered in Study.md; S17–S19 add
  analytical role/escalation cases. No engine behavior or policy adoption claimed.
- RGS2-B-M1: dedicated MCP report and 19 local / 10 external source pins
  delivered. Proposed interfaces, protocol intersection gap and durable recovery
  boundaries are explicit. No MCP integration or product tests were executed.
- RGS2-C-M1: app-server report, isolated 0.158.0 schema-generation evidence
  and selected schema extract delivered. No account/model/runtime probe performed.
  Fresh independent review approves report content; final lifecycle reconciliation
  and study closeout remain in phase D.
- RGS2-D-M1: Synthesis.md reconciles responsibilities, version/transport findings,
  decisions and the next DesignPlan brief. Independent study review approved the
  report candidate with no material findings. RGS2-Verification.md records the
  measured checks and runtime gaps. All selected study outcomes are delivered;
  capability implementation remains unselected.

## Study closeout

All study milestones are delivered, with report links in the result register.
There is no deferred required research milestone hidden behind completion. Open
runtime/compatibility questions are explicit inputs and acceptance probes for the
next DesignPlan, not claims that a client or routine engine already works.
KB036–038 return to backlog rather than falsely closing the proposed capabilities
or retaining in-progress status after their selected study work has ended.
No owner review, integration or publication approval is manufactured by this closeout.
