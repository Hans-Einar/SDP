# RGS2 synthesis — a governed Codex workflow with shared SDPTool state

This report belongs to [PLAN-SDP-0007](StudyPlan.md). It integrates the governance,
MCP and app-server studies into a recommendation for subsequent design. It is not
an implementation plan, a selected product architecture or permission to configure
an account, publish a release or migrate a project.

## Decision to prepare

Recommend one local pilot that connects project supervision, a bounded task Master,
SDP context/routine operations and an observable return disposition. Preserve the
existing Go SDPTool as the process and model facade. MCP exposes its operations to
agents; app-server exposes Codex to a client. The two interfaces are complementary.

Keep the project-facing agent responsible for the wider system and selected work.
Start a fresh Master for a bounded assignment with enough surrounding context to
understand its impact. The Master coordinates permitted work/review children and
returns evidence. Do not make a long-lived chat the sole owner of project memory,
or expect an isolated Worker to infer the entire product from a file list.

## Findings that affect the next step

The [MCP study](MCP-Study.md) verifies useful existing SDPTool read/navigation and
installation boundaries, but finds no generic routine service or adapter in the
inspected Go sources. Preview/viewer operations can write files or launch programs;
they must not be exposed as if they were passive knowledge reads. The proposed
core and adapter therefore need explicit side-effect classifications.

The MCP evidence distinguishes the earlier 2025-11-25 baseline from the pinned
2026-07-28 specification and Go SDK v1.8.0. Session and notification/replay semantics
differ between protocol eras. The common protocol supported by the installed
Codex candidate and a selected SDK has not been exercised. Design a compatibility
probe first and keep SDP run/event identity independent of the chosen transport.
The study does not recommend “latest” as a compatibility policy.

The [app-server study](AppServer-Study.md) measures schema generation from installed
Codex 0.158.0, with 314 normal and 440 experimental files. It identifies concrete
differences from the running documentation: collaboration item naming, detached
review deprecation and an internal-only external-token variant. Use the generated
contract plus measured behavior for a pinned client; exported schema alone is not
permission to use internal features or proof that a method works.

No account, model turn, app-server session or MCP integration was exercised. Managed
ChatGPT login is the documented candidate for subscription use. Passive attachment
to an existing TUI, multi-controller operation and notification replay remain
unverified. The next slice should have one controller, explicit ownership and a
read-only observer projection. A second network client is not needed merely to
show a second terminal pane.

Native child delegation remains the owner's desired role mapping. Client-managed
fresh threads are a technical alternative to evaluate if needed for context or
review isolation. The Master still decides and coordinates authorized task work;
a UI client does not become an autonomous project manager by launching threads.

## Boundaries to retain in the next design

| Responsibility | Owner | Must not be inferred from |
| --- | --- | --- |
| Product intent, priorities and owner-reserved decisions | Owner, assisted by Steering/PM | Agent completion text or transport success |
| Current project authority and selected work | Existing SDP records, resolved through SDPTool | Whichever worktree/default branch was first indexed |
| Routine applicability, run identity and validated transitions | Proposed SDPTool core capability | An MCP session ID or Codex thread ID |
| Agent/thread execution, host approvals and activity | Codex/app-server | SDP CardState alone |
| Assignment coordination and integrated candidate | Bounded Master role | A single child Worker summary |
| Review and verification claims | Evidence records and assigned roles | A plan checkbox or a role name supplied by the caller |
| Owner view of project/work state | Read-only projections in the client/observer | A second independently maintained workflow database |

The shared core must distinguish recommendations, authorizations, observations and
accepted transitions. Its callers need revision checks and explicit binding to
project/checkout/assignment. A model-supplied `role=owner` is not a credential.
An external card reference may remain unverified; required peer evidence may not.
The existing ProjectManagement and Traceability records retain their different
purposes. Operational activity needs recoverable correlation, not duplicate
competing management histories.

## Four initial routine paths

Use the [governance reconciliation](Study.md#10-rgs2-a-m2--operational-responsibility-and-procedure-reconciliation)
and [catalog](Routine-Catalog.md) to design a small executable subset:

1. Inspect/triage: answer directly, reuse current work, capture a proposal or
   identify the missing prerequisite.
2. Bounded Maintenance: selected plan, adequate task handoff, execution, candidate
   verification, applicable review and return disposition.
3. Resume/reconcile: restore durable work identity, refresh observations, recover
   interrupted record updates and reassess affected evidence.
4. Procedure gap: distinguish absence from unavailability/conflict/incompatibility,
   deduplicate the finding and recommend a disposition without inventing policy.

Unsupported routes must be visible. A release request can initially be shown as
an unselected or unsupported extension; no real publication is needed to test
that boundary. The sixteen-family inventory remains the longer-term coverage map.

## Decision register

| ID | Disposition | Reason and next responsibility |
| --- | --- | --- |
| D01 | Owner-selected: project supervision plus fresh bounded task Masters | Preserve this in handoffs and acceptance; measure native delegation/context behavior before promising isolation |
| D02 | Existing direction: shared Go SDPTool core | CLI, MCP and client projections should reuse its rules; new routine APIs are still proposed |
| D03 | Recommend: local STDIO MCP first | Smallest deployment boundary; prove the actual host/SDK protocol intersection before selecting versions |
| D04 | Recommend: one app-server controller and a read-only terminal projection | Avoid unverified competing controller/approval behavior; no Codex fork required for the first compatibility/pilot work |
| D05 | Design work: trusted caller and approval provenance | A self-declared role or tool success must not satisfy an owner/reviewer gate |
| D06 | Design work: operational persistence and recoverable record projection | Existing installer journals offer precedent; they are not a generic routine transaction implementation |
| D07 | Design work: compatible routine distribution and active-run retention | Reuse signed install mechanisms; preserve project refinements and exact historical definitions |
| D08 | Owner selection pending: bounded DesignPlan/pilot execution | Studies prepare the decision; they do not authorize configuration, model trials or product implementation |

## Candidate next DesignPlan brief

**Outcome:** a reviewable contract and executable test design for one local
supervision-to-task workflow, using current SDP records, compatible Codex/MCP
interfaces and one authoritative state core.

**Inputs:** the three studies, actual installed capability evidence, owner role
clarification, R01–R16 and the scenario extensions. Reuse the separate BP2 blueprint
work where applicable; do not require full experimental MVP1 language support.

**Required design deliverables:**

- Identity and authority contract: project, checkout, request, routine version/run,
  assignment, candidate, Codex thread/session and evidence relationships.
- A small operation/resource contract shared by CLI, MCP and client projections,
  with actual authorization, revision, error, idempotency and recovery semantics.
- A verified compatibility profile and explicit host/tool interception coverage.
- One selected client/observer approach and lifecycle ownership contract.
- One routine definition and exception/resume/gap paths, with representative
  request examples and executable acceptance fixtures.
- A bounded implementation sequence with milestones, preservation rules and
  evidence requirements. Existing Git policy must be selected in that plan.

**Non-goals:** a generic graphical workflow editor, complete portfolio onboarding,
a new agent runtime, automatic role approval, arbitrary shell containment, all
sixteen routines, complete SDL runtime or generated production code, live release
publication, and an immediate long-lived Codex fork.

This is a brief for selecting the next DesignPlan, not activation of that plan.
Technical choices inside the selected outcome can be resolved by its responsible
agents; only material unresolved product/authority choices belong with the owner.

## Pilot evidence, before broader rollout

The first pilot should use an isolated SDP project fixture and one real bounded
Maintenance change, preserving a deliberately unrelated file. Supply a selected
card/plan, current context and explicit non-goals. The task Master must perform
useful work rather than spawn agents merely to animate the display.

A successful pilot must demonstrate the following observable outcomes:

| Case | Required evidence |
| --- | --- |
| Explain versus act | A question creates no work item; an implementation request resolves its existing work coverage |
| Fresh task Master | A new bounded context receives assignment plus surrounding invariants and returns the identified candidate |
| Scope drift | A changed consumer/constraint produces an explicit scope delta and appropriate escalation |
| Review boundary | Self-check, independent review and owner-reserved disposition remain distinguishable |
| Missing procedure | Missing/conflicting/unavailable routes have distinct results and repeated gaps do not create duplicate cards |
| Stale candidate | Earlier verification cannot satisfy a gate after relevant inputs change |
| Retry/crash | An interrupted accepted operation reconciles once; no lost or duplicated management event |
| Observer disconnect | Work survives closing/reopening the view, with restored actual status rather than invented progress |
| Host coverage | Known unguarded paths are reported; tested invalid guarded calls are rejected |
| Account/transport | Selected authentication mode and compatibility are established without silent API-key fallback |

These are acceptance proposals. No integrated runtime or paid model trial was run
by writing this synthesis. Record both successful and failing trials; measure owner
reminders, false routing, duplicate records and recovery effort. Product behavior
in XFMD/Ponsse requires a later actual application-level trial, not extrapolation
from a fixture or a passing JSON schema check.

## Sequencing and adoption

First settle identity/authority and the protocol intersection. Then implement
one bounded workflow through the shared core and a minimal client, with observation
and recovery included. Expand blueprint context and a second project's refinement
only after that workflow is useful. Distribute tested definitions through the
existing Go installer and signed inventory when a release is actually authorized.

Recommendations remain pending selection. The owner has already selected the
role model and research scope; do not ask again whether supervision and bounded
Masters are desired. The unresolved decision is the next bounded delivery and
its real prototype/implementation authority, not permission to remember this study.

## Reading and evidence map

- [StudyPlan result register](StudyPlan.md): study state and delivery links.
- [Governance reconciliation](Study.md#10-rgs2-a-m2--operational-responsibility-and-procedure-reconciliation): responsibilities, handoffs and initial routine subset.
- [MCP study](MCP-Study.md) and [source pins](Evidence/RGS2-MCP-source-pins.json): actual core reuse, protocol/SDK evidence and candidate API.
- [App-server study](AppServer-Study.md), [probe evidence](Evidence/RGS2-AppServer-evidence.json) and [schema extract](Evidence/RGS2-AppServer-schema-extract.json): measured exported contract and untested runtime boundaries.
- [Analytical scenarios](Scenarios.md): RGS1 cases and RGS2 S17–S19 role/scope extensions; not executed engine tests.
