# PLAN-SDP-0026 — Complete SDPTool architecture migration

| Field | Value |
| --- | --- |
| id | PLAN-SDP-0026 |
| project | SDP |
| state | completed |
| PlanType | ImplementationPlan |
| BranchPolicy | current |
| CommitPolicy | milestone |
| Systems | SDPTOOL; collaborating SDL, SDUI and XFMD |
| source | Owner request for a complete migration proposal, Session0016 T5; KB-SDP-062 |

## Authority, readiness and outcome

Owner decision, 2026-11-12, Session0016 T6: select restructuring as the first
manually registered Activity, ACTIVITY-SDP-0001. Activity implementation follows
as a separate second Activity; browsing is intended as the first tool-created
Activity, third in the sequence. No new SDPTool functionality belongs in this
restructuring plan. Owner continuation in Session0017 T002 authorizes execution
through this Activity. M0 starts now; later milestones retain their design and
evidence gates. Publication and owner acceptance remain separate.

Deliver one coherent architecture for existing SDPTool entrypoints and operation
families while preserving their public behavior and exposure. Existing conforming
packages remain with evidence. Do not add Activity commands, new public actions,
new browsing behavior or new Blueprint features. Verification tooling may establish
migration correctness; it does not expand the shipped product capability set.

This plan owns milestones and removal obligations. KB cards own needs/issues;
Sessions record execution and handoffs, never redefine completion. The manual
[Activity mandate](../../../Activities/Activity-0001--SDPTool-Restructuring/01-Mandate/Mandate.md)
owns scope. Its [Session0017](../../../Activities/Activity-0001--SDPTool-Restructuring/Sessions/session-%230017--SDPTool_restructuring.md)
starts coordination; original Session0016 remains at its historical path.

References: [architecture assessment](../../../04--Design/SDPTool/ProjectGovernance/Activities/ArchitectureReview.md),
[Activity design](../../../04--Design/SDPTool/ProjectGovernance/Activities/Design.md),
[command design](../../../04--Design/SDPTool/ProjectGovernance/Activities/Commands.md).
KB062 is selected for restructuring. KB061 is on hold for the next Activity;
KB060 browsing stays queued behind Activity support. Both remain outside execution
scope here. The earlier combined migration/Activity proposal is superseded by this
explicit owner split; historical proposal records remain unchanged.

## Finite inventory and architectural acceptance

M0 must enumerate each public operation and actual consumer across sdptool,
sdp-controller, sdp-mcp, library callers, gh-sdp and XFMD. Reconcile CLI dispatch,
action registry, MCP registrations, Go imports, docs and tests; no single list
proves completeness. Cover discovery/navigation/selection, models, Blueprints,
governance/Sessions, program/preview execution, installation/upgrades and release
utilities. Include default paths, errors, cancellation, retries and durable state.

For every row record current entrypoint/handler/state owner, target owner, existing
adapter exposure or justified exclusion, preserved contract, replacement/deletion,
milestone, candidate revision and evidence. Classify conforming/retained, migrate,
retire, or external boundary. No unclassified row at M0 exit. Record newly found
rows during execution; agents cannot silently narrow this denominator.

Target invariants to refine in M1:

- CLI, actions and permitted MCP adapters call typed application operations;
  internal domain calls never manufacture argv or parse presentation output.
- Application/domain ownership and dependency direction are explicit. Domain
  services do not import CLI, GUI or host adapters.
- Public operation metadata and adapter exposure derive from one operation
  definition per capability, composed centrally. Domain authorization remains
  authoritative; metadata is not permission and no universal write role is assumed.
- Each durable state has one owner. Transactions, idempotency, cancellation at
  commit boundaries and interrupted-operation recovery have explicit contracts.
- SDPTool owns consumer navigation/projections; SDL/SDUI keep language semantics.
  Existing external protocols remain stable unless separately selected to change.
- No migration-only adapter, obsolete dispatch path or duplicate authority list
  remains at final acceptance. Supported old CLI syntax may remain as a permanent
  thin adapter into the target service; it must not retain the old internal path.

## Phases and observable milestones

| Milestone | Delivery | Exit evidence |
| --- | --- | --- |
| M0 Baseline and scope | Complete operation/consumer inventory; code-matched NOW; known defects separated from preserved behavior | Inventory reconciliation, pinned baseline code/tool versions, representative successful/error/recovery workflows; model gaps explicit |
| M1 Target and contracts | Whole-system responsibility/dependency model, detailed contracts for changed paths, code mapping, removal matrix and revised milestone sizing | Validated SDL TARGET; every inventory row has a disposition; architecture rules and unsupported semantics have named verification methods |
| M2 Enforcement and Blueprint trial | Architecture checks and a real existing workflow migrated end-to-end: Session action/CLI through typed service | Adapter parity and authority tests; old argv/JSON bridge removed for the selected path; a fresh-context reviewer can assess the pinned work package without chat history |
| M3 Complete existing-operation migration | M3a discovery/navigation/selection; M3b models/Blueprints; M3c all remaining PG/controller/MCP paths; M3d programs/preview/install/release utilities | Each group fully classified/migrated or retained with evidence, consumers exercised, corresponding old paths removed; complete dependency sweep |
| M4 Transferred by owner; outside this plan | Activity and Session creation, lifecycle, tree/timeline and admission belong to the future second Activity (KB061) | Not implemented or counted as delivered here; identity retained to explain the scope change |
| M5 Whole-candidate qualification | All inventory rows closed against one candidate; documentation/examples aligned | Independent integrated review, architecture and regression checks, zero migration leftovers or missing required evidence, owner acceptance disposition |

M0/M1 may refine submilestones, not silently remove the completion obligations.
M2 is a real production-path trial inside this migration, not a throwaway example
or permission to stop after a pilot. Do not postpone all deletions to M5: each
M3 group removes its replaced paths before closing. M5 verifies the whole result.
External consumer work must have named ownership and verified integration; a
headless fixture alone cannot prove native XFMD acceptance.

## Execution disposition — Session0017

Product candidate: 527cc9008bdf8f567f8b201b665b2960865518f1.
[Verification and inventory closure](../../../Activities/Activity-0001--SDPTool-Restructuring/06-Verification/Verification.md)
owns exact results and provenance; the following is the authoritative milestone
status, not another independent plan.

| Milestone | Current disposition |
| --- | --- |
| M0 | Complete: pinned code/library/action inventory, controller/MCP/consumer reconciliation, observed SDL and passing baseline suite |
| M1 | Complete: validated target WORK, immutable Blueprint plus behavioral/removal contracts and authored source mapping |
| M2 | Implemented and independently reviewed: typed governance application and direct Session adapters; negative architecture specimens detect prohibited paths |
| M3a | Implemented/retained: one root CLI registration; existing navigation/selection services and native producer contracts preserved |
| M3b | Retained with model/Blueprint/assignment package evidence; caller and availability policy consolidated without catalogue change |
| M3c | Implemented/retained: common MCP response policy; existing capability, governance, controller and history boundaries preserved |
| M3d | Implemented/retained: typed program execution split; installation/bootstrap/release utilities preserved |
| M4 | Transferred by owner to the next Activity; no functionality implemented here |
| M5 | Technical qualification complete: exact-candidate full suite, native gates, independent review and inventory closure pass; owner accepted |

No temporary adapter or deferred migration cleanup remains. Technical delivery is
complete. Owner explicitly accepted the concrete result after branch publication and
authorized Activity0002. Plan and Activity are completed; merge/release remain separate.

## Blueprint's bounded role

Current `SDPTool/blueprints` can capture NOW/TARGET, task rules/protection, retained
bundles and revision-bound evidence/assignments. Assessment does not execute tests,
interpret Go symbols or prove whole-repository conformance. The new Activity WORK
is not a complete or code-verified NOW model of SDPTool.

Use a fixed whole-system target plus per-milestone retained bundles from verified
current state to the next state. Each bundle and any authored companion contract
must identify intent/rationale, preserved behavior, permitted changes, code mapping,
required deletions, acceptance checks and evidence locations. Use supported Task
fields and linked authored contracts; do not invent machine-enforced schema fields.
Plan/target/task changes invalidate dependent handoffs and require regeneration
and review. Commit source checkpoints through supported ModelGovernance operations.
Generate required Blueprint bundles only at the selected model/implementation
handoff. No generated viewpoints are maintained manually.

At M2 demonstrate that a worker can use this package, perform the change and produce
candidate-bound evidence, and that independent review detects a deliberately
nonconforming specimen. Use authored companion contracts where current Blueprint lacks support. New
Blueprint capabilities belong in a later Activity; if a missing capability makes
verification impossible, record the blocker rather than claiming completion.

## Make architectural understanding executable

Implement focused Go import/AST and registry checks for the chosen boundaries,
typed registration/decoding, adapter exposure and forbidden bridging. Verify them
against violating fixtures so a vacuous green check cannot satisfy the gate.
Use behavioral parity and integration tests for properties syntax cannot prove.
Trace each invariant to its check or explicit review procedure. Existing defects
are not automatically the expected behavior of characterization tests.

The tool's generic unknowns remain honest. Supply independent checks/review for
required code and runtime properties; a scoped Blueprint disposition cannot waive
the plan's completion obligations. A ready assessment alone is insufficient.

## Sessions, continuity and temporary code

Session0016 holds the proposal history; Session0017 owns current restructuring
coordination. Create new Session logs under this Activity root manually until
tooling exists; each points to this plan and selected milestone.
Do not precreate imaginary completed Sessions or duplicate milestone lists.

Each handoff records target/task/bundle revision, actual candidate commit and dirty
state, completed/open inventory rows, evidence and its freshness, temporary paths
and exact next step. Incoming workers explain intended ownership, constraints and
deletions from that package before changing code. Contradictions require a design
update; existing transitional code is never implicit architectural authority.

If an intermediate adapter is unavoidable, register its actual consumers, purpose,
named removal milestone within this plan and a final zero-consumer/deletion check
when introducing it. New consumers may not attach. Failure to remove it keeps the
plan incomplete; moving it to a backlog card does not satisfy acceptance.

## Completion and publication

Completion requires all inventory rows dispositioned with evidence on the final
candidate, all target invariants established, required integrated behavior passing,
zero migration-only paths, current canonical docs/models and independent review.
Parser validity, passing unit tests or an empty issue list alone are insufficient.

New unrelated wishes may go to backlog. Discoveries necessary for the selected
architecture or existing consumer compatibility stay inside this plan. Scope may
change only through an explicit owner decision, retaining the prior obligation
and reason; reduced scope must not be reported as the original complete delivery.
Time/context/budget exhaustion means incomplete work with a resumable handoff.

Plan creation used sdp/sdl-browsing-planning. Authorized execution now uses the
sdp/sdptool-architecture-migration integration branch
from a reconciled baseline, commits per milestone, and no direct main edits.
BranchPolicy current means continue that execution branch across phases, avoiding
separate long-lived phase branches. Preflight before branch creation/execution must
preserve existing uncommitted design work and other assignments.
No push, merge or release is authorized by this proposal. Verify the exact built
delivery candidate and its supported entrypoints; publication remains a distinct
disposition with release-checklist evidence if selected.
