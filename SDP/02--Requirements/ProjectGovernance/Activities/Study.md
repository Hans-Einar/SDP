# Need-driven Activities — workflow proposal

2026-11-11. [KB-SDP-061](../../../KanBan/active/%23061--Proposal--Need-driven-activities.md),
[Session0015](../../../Sessions/session-%230015--SDL_browsing_planning.md).
Proposed management extension, not an adopted schema, installed process migration,
implemented Activity board or runtime enforcement claim.

Owner follow-up, 2026-11-11: Activity integration is selected ahead of SDL browsing.
The [complete design candidate](../../../04--Design/SDPTool/ProjectGovernance/Activities/Design.md)
and Session0016 now refine this study. Earlier suggestions below that browsing
could proceed independently are superseded by that priority decision.

Owner refinement, 2026-11-12: the design now fixes standard README guidance,
Mandate-owned need/outcome, root-level Sessions and Activity-scoped implementation
while allowing unrelated discussion/capture. Commit and ledger event markers
supplement the state-duration timeline. These refinements supersede earlier
layout/selection alternatives in this initial study, without claiming runtime adoption.

## Need and recommendation

Sessions are discoverable but do not give the owner a clear overview of which
customer need they serve across studies, design changes and implementation plans.
The owner proposes Activity as that durable connection, with successive Sessions
and one active Activity at a time. Recommend this concept, with explicit ownership
so it does not become a second KanBan/plan/Session state store.

Current Session guidance already owns a long-lived goal across plans. The proposed
change moves the need-level roadmap above Sessions; it is a real responsibility
change requiring compatibility/adoption, not just another folder. Preserve old
Sessions and link them retrospectively with provenance rather than rewriting history.

## Concepts and sources of truth

| Concept | Owns | Does not own |
| --- | --- | --- |
| Need / need brief | Customer/owner context, changed situation, desired outcome, constraints and success evidence | Detailed solution or testable system requirements by itself |
| Activity | Pursuit/disposition of that need, scope, accountable owner, linked Sessions/plans/cards, current route and next step | Copies of plan milestones, every card's status, host-attempt liveness |
| Requirement | Agreed testable obligation derived from need and analysis | The entire rationale or development journey |
| KanBan card | Discoverable technical wish, defect, question or bounded obligation | Mandatory one-to-one Session or Activity wrapper |
| Typed plan | Deliverables, milestones, dependencies, acceptance, Git policy and evidence | Conversation/turn history |
| Session | Bounded work objective, selected plan milestones, decisions/turns and successor handoff | Another implementation milestone authority |
| Blueprint | Reviewed design/change context, authored scope, preserved obligations, source identities and unknowns | A complete schedule, execution authorization or proof of code conformance |
| PG run/attempt | Actual bounded agent execution, receipts and operational recovery | Customer need acceptance or canonical lifecycle truth |

Cards may be captured before an Activity exists; an Activity may draw on several
cards, and a card may contribute to several Activities over time. Record selected
scope/contribution so one delivery does not falsely complete a broader card.
Each Session has one primary Activity when that feature is adopted; incidental
findings may create cards for other future Activities without changing its purpose.

An Activity can span SDP/SDL and collaborating XFMD work: choose one coordinating
home and reference other repositories' plans/cards. Avoid one independent Activity
copy per repository for the same need. MINI and MICRO could be distinct HEOS
Activities if they represent distinct needs; do not relabel existing history now.

## State, stage and focus

Recommend lifecycle states planned, active, on-hold, completed and canceled.
Do not make active and in-progress two overlapping states. Study, requirements,
architecture/design, planning, implementation and verification are **stages**,
not additional lifecycle states. Waiting/blocked reasons are attention conditions;
they do not falsely complete the Activity.

Recommend at most one active Activity per explicitly selected project/workstream
authority. This matches focused owner work without globally locking unrelated
projects. The exact scope of that authority remains a design choice. A spanning
Activity uses its coordinating home. Default Sessions are successive; deliberate
parallel assignments inside the same Activity must remain possible if authorized.
Switching active Activity requires an explicit disposition of the previous one;
it must not silently cancel, pause or terminate existing agents.

This is currently a proposed human policy. Runtime enforcement later needs a
versioned, atomic selection/transition rule and recovery across clients; a Markdown
Active field alone cannot enforce a concurrency limit. Existing active Sessions
and work are not automatically paused by registering this proposal.

## Proportionate Activity checklist

Each stage records its deliverable/disposition or a reason it is not applicable.
Small maintenance should not require empty model releases or seven ceremonial plans.
Stages can revisit earlier decisions when evidence changes.

1. Capture the need brief, success conditions, source cards and constraints.
2. Study the existing system, feasibility, alternatives and impact. No change or
   a deferred/rejected solution is a valid explicit outcome.
3. Reconcile requirements and acceptance with the selected direction.
4. Study architecture/design; update applicable SDL WORK and prose contracts,
   distinguishing preliminary drafts from reviewed targets.
5. Prepare the Blueprint/assignment context and detailed ImplementationPlan.
6. Execute selected milestones through one or more Sessions; retain discoveries,
   evidence, reviews and scope decisions.
7. Verify the need-level outcome, complete applicable model/software release
   obligations separately, and record explicit acceptance, deferral or cancellation.

Need brief and workflow checklist are not replacements for existing mandate,
requirements or model governance. Existing release/design gates remain authoritative.

## Blueprint to implementation plan

Yes: use the Blueprint as a primary design input to implementation planning.
Supply NOW/TARGET identities, authored task and permitted scope, requirement/acceptance
links, preserved behavior, code/build/test context and explicit unknowns. SDL alone
does not provide code mappings, deployment risk or all behavioral requirements.

Preliminary Blueprint analysis can inform estimates before final design disposition;
execution must bind to the actually selected/reviewed target under existing PG
rules. Work-model review, model release and software release are different events.
Do not assume that generating a Blueprint freezes WORK or completes a model release.

A Worker explicitly assigned planning may draft the bounded ImplementationPlan and
propose execution Sessions. The coordinating role reconciles dependencies and scope;
that assignment does not grant the Worker authority to select unrelated work,
approve its own design or create independent execution agents. Preserve existing
authorization; avoid adding a redundant approval when execution was already selected.

Implementation milestones belong in the ImplementationPlan. Sessions reference
the milestones they will address, and record what actually happened. Plan future
Sessions as route slots with outcome/dependencies; instantiate them successively.
One milestone may span Sessions; one Session may deliver multiple small milestones.
A new chat thread or compaction is not by itself a new Session.

## Fit with existing SDP and ProjectGovernance

- [Session guide](../../../Sessions/README.md) and KB042 own current continuity;
  reconcile the proposed move of long-lived goal ownership explicitly.
- [Typed plans](../../../ProjectManagement/Plans.md) already own milestones. No
  new SessionPlan or RefactorPlan schema is required for this proposal.
- [PG state boundaries](../../../04--Design/SDPTool/ProjectGovernance/State-Boundaries.md)
  and [Session execution](../../../04--Design/SDPTool/ProjectGovernance/Session-Execution.md)
  distinguish canonical records from operational runs. Link an Activity to those
  identities; do not build a second agent execution engine.
- SDL already uses activity for modeled delivery work. Management Activity must
  be explicitly qualified and must not reinterpret every SDL activity declaration.
- Extend the existing canonical management history when adopted. Do not append
  invented Activity events to the current schema or introduce a parallel ledger.

## Proposed adoption sequence

First review this conceptual boundary and pilot a linked need/route overview using
ordinary documents. The owner now selects Activity integration before browsing;
qualify the adopted workflow before resuming PLAN-SDP-0025.

Then specify identifiers, lifecycle, membership, adoption of old records, canonical
history and the scope of the one-active rule. Deliver versioned read/write contracts
and tests before claiming enforcement or adding Activity to installed templates.
Expose Activity → Sessions → selected plan/milestones and cards through SDPTool's
producer-owned navigation, including current stage, next action and explicit holds.
Use the same existing PG runtime for execution and recovery.

Acceptance should cover several Sessions serving one need, shared/partial cards,
a cross-repository Activity, explicit focus switching, interrupted Sessions,
historical unassigned Sessions and completed plans with unmet need-level acceptance.
The UI must not infer completion from closed-card counts or successful agent turns.

Open design choices: exact active-selection authority, identifier/metadata contract,
adoption mapping and completion/acceptance authority. The proposal records the
owner's direction without claiming those choices or implementation already exist.
