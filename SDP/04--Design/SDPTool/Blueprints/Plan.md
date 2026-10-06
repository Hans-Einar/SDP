# BP2 — design a minimal SDL assignment-bundle contract

| Field | Value |
| --- | --- |
| id | PLAN-SDP-0001 |
| project | SDP |
| state | active |
| PlanType | DesignPlan |
| BranchPolicy | current |
| CommitPolicy | milestone |
| Systems | SDL, SDPTOOL |
| source | KB-SDP-050 / KB-SDP-031 / BP1 study recommendations |

## Outcome and authority

Turn the [BP1 findings](Study.md) into a reviewable contract for a small generated
assignment bundle. [Worked example](Pilot.md) provides the initial real model and
behavior boundary. This successor to the Study was selected on 2026-10-06.
Contract recommendations remain proposed where specific owner disposition is
not yet recorded; execution authority is distinct from approving every design choice.

The output is a detailed design and a proportionate ImplementationPlan, not a
production command, SDL grammar change or mandatory new project structure. Keep
one authored task authority and existing management/Traceability ownership.

## Git and scope

2026-10-06 reconciliation: use the current ModelGovernance implementation branch
for this bounded design milestone because its APIs are not yet on main. This revises
the unexecuted main-based branch proposal; preserve the BP1 study history and
concurrent ProjectGovernance edits. No shared branch switch is made. Use milestone commits and the existing phase-push/PR authorization.
No merge or release is implied. No XFMD application change, Issue #7 adoption,
source-set implementation, scheduling agent or automatic Go call-graph proof.
Do not create a parallel ledger or mandatory per-Issue assignment hierarchy.

## Phases and milestones

| Phase | Milestone | Observable acceptance | State |
| --- | --- | --- | --- |
| BP2-A — contract | BP2-A-M1 | Specify authored intent versus generated/observed facts; freeze identity, required inputs, gaps, NOW/TARGET and pinned before/after test baselines; define box/edge constraint marks; choose one owner-reviewed pilot with explicit invariants | Partial: specimen ready; explicit pilot disposition pending |
| BP2-B — selection and evidence | BP2-B-M1 | Define complete modeled impact context, explicit unknown frontier, typed inclusion/exclusion and atomic contracts; demonstrate current/target union, removed-edge neighbor, cycle, missing provider, size overflow and stale inputs; specify SDL code-tag mappings and evidence references with KB-SDP-004 without claiming its full delivery | Delivered BP2-B-M1 |
| BP2-C — executable handoff design | BP2-C-M1 | Specify Go library/SDPTool boundary, before/after verification, failure/publication behavior and annotated output; walk one authorized worker/reviewer through the bundle and a scope-violating change, recording observed results; write implementation slices with measurable acceptance | Planned |

Use the study probe as evidence input, not the production selector. Define any new
fixture or API as proposed until selected. Do not require a live XFMD viewer to
validate basic usefulness; generated Markdown and source-linked facts suffice.
An actual worker trial requires explicit permitted delegation and a bounded task;
if unavailable, report that evidence gap rather than simulating independent review.

## Proposed implementation roadmap to refine in BP2-C

These are candidate vertical increments, not activated implementation phases.

1. **One reproducible bundle:** one source-discovered structural model, one authored task,
   existing contract documents, current/target identity, deterministic Markdown and
   manifest, validation failures and provenance. A headless consumer can inspect it.
2. **Preserved context and change checks:** typed boundary rules, removed/added fact
   comparison, authored invariants, explicit unknowns, stale regeneration and
   contract-preserving publication. Exercise adversarial cases from the study.
3. **Evidence and workflow integration:** reviewed model-to-code mappings and
   scoped Traceability links; actual worker/reviewer pilot and acceptance report.
   Native XFMD navigation is a separately coordinated consumer enhancement only
   if the headless workflow establishes value.

Do not design every future source profile before delivering the first increment.
No production command syntax is promised by this plan; SDPTool remains the facade,
SDL owns model semantics, and existing document publishers/renderers are reused.

## Dependencies, verification and completion

KB-SDP-004 owns the general design/code/evidence contract. BP2 can represent unknown
or directly cited evidence without inventing a replacement ledger. Source composition is now implemented in design-core/0.6; this pilot exercises
a composed model and ModelGovernance snapshots without registering source files. KB-SDP-032/033
supply future navigation feedback/adoption experience rather than hard prerequisites.

Completion requires a concrete reviewed design with explicit remaining decisions,
a reproducible pilot and negative cases, evidence distinguishing machine checks
from human judgment, and a proposed implementation plan. No GUI acceptance,
implementation-conformance proof, production feature or owner agreement may be
inferred from parser success or the study's completion.

## Owner clarification incorporated — BP1-M3

The [owner clarification](Study.md#owner-clarification-after-bp1--2026-09-26)
requires the impacted environment, explicit box/connection marks, complete NOW
and TARGET SDL snapshots, reproducible before/after code checks and SDL code tags.
Study recommendations have not received a detailed owner review. TARGET is a
working term; final naming and annotation syntax remain open. This plan stays
planned and the producer contract/model remain unchanged.

Executable Channel tests and companion-language choices belong to
[KB-SDL-006](../../../KanBan/backlog/%23006--SDL--Study--Executable-channel-tests-and-unit-bindings.md).
Coordinate its eventual test/Unit binding contract; do not silently expand this
DesignPlan into runtime implementation. The first blueprint can cite named tests
and report missing executable mappings without claiming general SDL execution.

## Execution — 2026-10-06

Owner selects continuing KB050 in Session0008. BP2-A contract candidate and a
reproducible reduced MVP1 pilot are delivered in Contract.md and
experiments/blueprint_mvp1. Pilot-MVP1-evidence.json records real parser checks,
model-only release -> dirty WORK -> candidate and ordinary generated viewpoints.
The baseline/target distinction does not claim implemented Ponsse conformance.
The specific extraction is agent-selected and still awaits owner feedback; no
owner-reviewed pilot or production blueprint generator is claimed. BP2-A acceptance
is therefore partial (contract/pilot ready, owner pilot disposition pending).
Next executable design work: BP2-B selection/evidence rules and negative specimens;
BP2-C and production implementation remain open.

## BP2-B-M1 — 2026-10-07

Delivered Selection-and-Evidence.md with versioned relation closure, NOW/TARGET
union, atomic contract groups, exclusions/unknown frontier, PRESERVE checks and
code/evidence references coordinated with KB004. Eight experiment tests pass over
actual toolkit ASTs plus a labeled synthetic cycle fixture. The real parser accepts
a missing channel receiver; selection reports that gap explicitly. Parser success
also does not prevent the tested ownership drift. See Selection-evidence.json and
Selection-tests.txt. No production producer or code-conformance claim.
Next BP2-C API/publisher design and implementation slices; BP2-A pilot feedback
remains pending rather than being inferred from the owner's continue instruction.
