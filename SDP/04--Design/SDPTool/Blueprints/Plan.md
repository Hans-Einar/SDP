# BP2 — design a minimal SDL assignment-bundle contract

| Field | Value |
| --- | --- |
| id | PLAN-SDP-0001 |
| project | SDP |
| state | planned |
| PlanType | DesignPlan |
| BranchPolicy | current |
| CommitPolicy | milestone |
| Systems | SDL, SDPTOOL |
| source | KB-SDP-031 / BP1 study recommendations |

## Outcome and authority

Turn the [BP1 findings](Study.md) into a reviewable contract for a small generated
assignment bundle. [Worked example](Pilot.md) provides the initial real model and
behavior boundary. This is the explicit successor to the Study; it is planned,
not selected for execution by the request to perform BP1. Recommendations remain
proposed until the owner selects the intended delivery and dispositions below.

The output is a detailed design and a proportionate ImplementationPlan, not a
production command, SDL grammar change or mandatory new project structure. Keep
one authored task authority and existing management/Traceability ownership.

## Git and scope

If selected, start one working branch from then-current main; preserve the BP1
study branch. Use milestone commits and the existing phase-push/PR authorization.
No merge or release is implied. No XFMD application change, Issue #7 adoption,
source-set implementation, scheduling agent or automatic Go call-graph proof.
Do not create a parallel ledger or mandatory per-Issue assignment hierarchy.

## Phases and milestones

| Phase | Milestone | Observable acceptance | State |
| --- | --- | --- | --- |
| BP2-A — contract | BP2-A-M1 | Specify authored intent versus generated/observed facts; freeze identity, required inputs, gaps, current/target and code-baseline handling; choose one owner-reviewed pilot with explicit invariants | Planned |
| BP2-B — selection and evidence | BP2-B-M1 | Define typed context inclusion/exclusion and atomic contracts; demonstrate current/target union, removed-edge neighbor, cycle, missing provider, size overflow and stale inputs; align evidence references with KB-SDP-004 without claiming its full delivery | Planned |
| BP2-C — executable handoff design | BP2-C-M1 | Specify Go library/SDPTool boundary, failure/publication behavior and readable output; walk one authorized worker/reviewer through the bundle and a scope-violating change, recording observed results; write implementation slices with measurable acceptance | Planned |

Use the study probe as evidence input, not the production selector. Define any new
fixture or API as proposed until selected. Do not require a live XFMD viewer to
validate basic usefulness; generated Markdown and source-linked facts suffice.
An actual worker trial requires explicit permitted delegation and a bounded task;
if unavailable, report that evidence gap rather than simulating independent review.

## Proposed implementation roadmap to refine in BP2-C

These are candidate vertical increments, not activated implementation phases.

1. **One reproducible bundle:** one registered structural model, one authored task,
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
or directly cited evidence without inventing a replacement ledger. KB-SDL-005 owns
multi-file source identity; the first pilot can remain single-file. KB-SDP-032/033
supply future navigation feedback/adoption experience rather than hard prerequisites.

Completion requires a concrete reviewed design with explicit remaining decisions,
a reproducible pilot and negative cases, evidence distinguishing machine checks
from human judgment, and a proposed implementation plan. No GUI acceptance,
implementation-conformance proof, production feature or owner agreement may be
inferred from parser success or the study's completion.
