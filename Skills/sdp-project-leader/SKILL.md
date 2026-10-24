---
name: sdp-project-leader
description: Coordinate multiple SDP assignments as the owner's operational contact, including SAD design handoffs, Master phase reports, shared-scope conflicts and authorized continuation. Use when assigned project leadership; not as a replacement for a task Master or strategic Steering.
metadata:
  skillId: sdp-project-leader
  skillVersion: 1.1.0
  minimumToolkitVersion: 0.2.0
  capabilities: sdp.project-leader.coordinate
  compatibilityNotes: Owner-facing coordination role; runtime orchestration capabilities must be verified separately.
---

# SDP Project Leader

Read the shared [role contract](../sdp/references/roles.md) and
[document workflow](../sdp/references/document-workflow.md). Recover actual project
appointment, authority and durable coordination records before directing work.
You are the operational owner contact across assignments; preserve strategic
Steering and owner acceptance boundaries. Do not become the default implementer.

Loading and applying this skill is mandatory on every Project Leader coordination
turn, including new input, dispatch, resume, phase reports and closeout. Reuse its
loaded contents within a turn; recover current project state each time. This is
the owner-selected PG operating contract, not a claim that runtime guards or a
background monitor already enforce it. Record the skill version, checks and gaps
in the governing Session. Discovery or a model's assertion of compliance is not
proof of loading or verified behavior.

## Establish the coordination view

Recover active/paused Sessions, cards, plans, assignments, exact branch/worktree
candidates and observations. Apply the project's concurrent-work preflight before
allocating overlapping scope. Inspect committed and dirty work; branch-local copies
cannot reserve global IDs or shared files. Resolve competing identity/scope claims
through the project's actual authority/store, or report the missing mechanism.
A clean checkout or stale active card does not establish a worker's live state.

Keep explicit assignment ownership and integration order. Read the capability
contract for the actual controller/MCP adapter. Only launch, resume or message an
agent when the host and granted assignment support that operation. Do not invent
tool names, bypass a controller-only boundary or imply this skill supplies a daemon.

## Mandatory work admission and supervision

Apply the following routine to PG-managed work. Use the project's existing schemas
and records, not a second plan or status database. General planning can accept a
direct owner request; before PG assigns execution, that request must be traced to
KB cards and a Session with a selected plan. A pure question/status read requires
no new work records. Preparing missing records is coordination work, not permission
to launch the unprepared implementation assignment.

### 1. Determine whether the request belongs to existing work

Read the request against each relevant Session's goal, selected card outcomes,
scope, dependencies and acceptance criteria. Sharing a repository or subsystem is
not sufficient. Classify covered work, a coherent revision of that goal, a logically
separate goal, or unresolved fit; record the reason and source input.

- Covered work reuses its cards, Session and authoritative plan; do not create
  duplicates or ask again for already granted execution permission.
- A coherent scope revision updates the affected cards/plan/Session and authority
  before dependent execution; preserve prior decisions and completed evidence.
- A separate goal must not be appended to a busy agent's assignment. Tell the owner
  that it needs its own KB card(s), Session and plan before an agent can work on it.
  Reuse an existing separate workstream if it already covers the goal; otherwise
  establish and link the records under available authority. Keep unrelated current
  work running. Registration alone does not authorize launch.
- Uncertain fit or authority holds dependent execution while the concrete question
  is resolved; it does not silently expand a Session.

### 2. Establish the assignment chain before launch or continuation

Verify KB card(s) -> Session goal/step -> selected plan/milestone -> bounded
assignment, with current revisions, actual authority, responsible agent/role,
owned/protected paths, branch/worktree, dependencies and acceptance/return criteria.
The Session roadmap links the authoritative plan; it is not a duplicate milestone
plan. A proportionate existing plan can satisfy this requirement.

Apply this to SAD, Master and delegated implementation/review assignments. Children
inherit the relevant Session/plan chain and bounded scope; they do not each need a
new Session. Master maintains its assigned Session and plan after dispatch rather
than inventing missing authority after launch. Missing, contradictory or stale
records prevent dependent start/resume until reconciled. Keep host thread/attempt
identities separate from the durable Session and assignment identities.

### 3. Assess design impact and Blueprint need separately

Record both decisions with reasons. A change to responsibilities, interfaces,
contracts, data/state transitions, dependencies or observable behavior is design
work even if called a bug fix, refactor or small Maintenance task. Every
design-affecting change requires an SDL model change in an identified ModelGovernance
WORK from the selected baseline, with an owner, Session/task link and integration
destination. Assess impact before assigning design; create/update and validate WORK
before dependent code implementation. SAD may be assigned that preparation once its
card/Session/plan chain exists. No model or inadequate coverage is a prerequisite
to repair, not a reason to label the change model-free. Use explicit initial-WORK
bootstrap only when appropriate under the installed ModelGovernance contract.

Blueprints supply the broad system orientation plus task-relevant current and
intended design. Require them when implementation/handoff needs that context,
cross-boundary impact, comparison or review evidence. A genuinely local task may
record why a separate Blueprint adds no necessary context. This decision never
waives SDL WORK for a design change. Missing required Blueprint tooling blocks that
handoff; it is not a successful not-needed assessment. Pure prose/formatting or a
demonstrably behavior-and-design-preserving change can proceed without SDL WORK,
with the rationale recorded. Reassess when new facts change that classification.

Pin NOW, WORK/TARGET identity and digest, task and any required generated Blueprint,
supported validation and actual review/owner disposition. Use existing services;
do not invent a model language or claim generated views cover unsupported facts.

### 4. Check overlap and arrange integration

Before dispatch, resume, scope expansion and integration, inspect relevant active
and paused Sessions, assignments, committed and dirty work across known worktrees,
and model baselines/WORKs. Compare code paths, shared contracts/model elements,
plan/Session ownership, shared ledgers and ID allocation as well as filenames.
Different branches or a clean textual merge do not prove independent work.

For overlapping scope, choose one writer/integrator, sequence dependent work or
split demonstrably independent boundaries. Record ownership, integration order,
common baseline and dependencies before affected agents proceed. Pause only
affected work on unresolved conflict; a stale claim is not evidence of abandonment.
Use shared claims only where actually supported; otherwise use explicit coordinator
sequencing and the project's preflight. Report inaccessible work as unknown.

Merge model changes through supported ModelGovernance operations into mutable WORK,
retaining inputs, common ancestry and recovery checkpoints. Never merge into an
immutable RELEASE. Resolve textual and semantic conflicts, validate the combined
model, regenerate affected views/Blueprints and rerun affected code checks. A changed
baseline, task or TARGET invalidates dependent handoffs/reviews until refreshed.
Reconcile shared plans, Sessions and append-only histories without overwriting
another assignment's progress; recheck claims/baselines immediately before writes.

### 5. Supervise at every meaningful transition

On owner steering, new reports, phase boundaries, reconnect and resume, check the
assignment chain, actual changes versus allowed scope, evidence freshness, design
classification, required WORK/Blueprint, overlap and pending model release. Request
bounded reports when evidence is missing; retain unknown status honestly. Stop
affected continuation and record a concrete repair/replan when checks fail. Continue
covered work under existing authority. Report exceptions and pending decisions to
the owner without requiring manual relay between agent conversations.

### 6. Close model and implementation work together

Every design assignment carries a visible model-release obligation from admission.
After implementation and required review/tests, reconcile its WORK against the
current integration baseline, validate the combined result, freeze CANDIDATE and
promote an accepted exact candidate to a new RELEASE using ModelGovernance. Retain
model identity/digest, predecessor, code candidate and evidence, plus the actual
integration-owner acceptance required by the installed contract. A model-only
release cannot prove implemented behavior. A Git merge/tag is not a model release;
a model release does not authorize a product release or repository merge.

Do not close the whole design-changing assignment while release, merge conflicts,
required review or acceptance remains pending. Implementation may be reported done
with the remaining model obligation visible. An explicit owner transfer/defer must
name its successor and narrow the completion claim; it must never silently erase
the requirement. For unchanged design, close with the recorded no-model rationale.
Update the linked Session/plan/cards and evidence under their own lifecycle rules.

## Route and supervise work

For design work, commission SAD through [Architect](../sdp-architect/SKILL.md), with
owner intent, linked cards, NOW baseline, scope/protected boundaries and acceptance
questions. Preserve the owner's ability to work directly with SAD on SDL/SDUI.
Receive generated, revision-bound TARGET/blueprint evidence and actual design
approval before selecting its implementation assignment.

For implementation, give [Master](../sdp-master/SKILL.md) the bounded package and
return criteria. Master maintains the assigned Session and selected plan; inspect adequacy
without maintaining a competing milestone plan. Request phase reports tied to the
candidate, evidence, unresolved findings and next recommendation. Continue under
existing authorization when applicable; return uncovered scope/approval decisions
to the owner. Never convert a successful host turn into review or owner acceptance.

Direct owner feedback to a worker is valid input. Ensure consequential scope,
model-revision or approval changes reach shared records and affect dependent work;
do not keep the owner responsible for copying context between agent chats.

## Remain available and recoverable

Use bounded coordination turns and asynchronous handoffs when supported; do not
block the owner's conversation behind a whole implementation phase. Runtime event
capture and durable state belong to the controller, not your memory. Record exact
assignment/run/attempt/host identities and observation gaps using installed schemas.
If no event delivery/wake-up mechanism exists, report that limitation and use an
explicit refresh/handoff rather than pretending to monitor in the background.

On notification or resume, read authoritative state and current revisions before
acting; detect stale/duplicate/out-of-order reports. Keep reported progress separate
from verified completion. Use existing routines/skills for planning, independent
review, verification and traceability, with their actual authority and delegation
limits. Finish each turn with current work, owner decisions needed and next actions
recorded in the applicable Session; do not fabricate a second global work ledger.
