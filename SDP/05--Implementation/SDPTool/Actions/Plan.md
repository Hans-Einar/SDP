# ACT — discoverable actions and JSON invocation

| Field | Value |
| --- | --- |
| id | PLAN-SDP-0024 |
| project | SDP |
| state | planned |
| PlanType | ImplementationPlan |
| BranchPolicy | current |
| CommitPolicy | milestone |
| Systems | SDPTOOL |
| source | Session0011; discoverable-actions KB-SDP-051 |

## Outcome, authority and baseline

Deliver the [proposed action contract](../../../04--Design/SDPTool/Actions/Contract.md)
so XFMD can discover and invoke a bounded blueprint workflow through JSON without
shell construction or help scraping. Owner selected Session0011 and continued it
in T007, retaining the earlier green concurrent-work prerequisite. Planning is
permitted; product execution remains blocked by the refreshed RED finding.

This plan starts from e2caa22 on sdp/blueprint-implementation. It is not a release
or a selected combined code baseline. Preserve the separate runnable-program and
ProjectGovernance deliveries. KB051 here means the discoverable-actions card at
backlog/#051--Proposal--Discoverable-SDPTool-actions.md on this branch, not the
runnable-program card with the same ID. Resolve that collision before execution.
PG owns KB052 and multi-agent orchestration. Native XFMD work is external KB-XFMD-030.

## Git policy

Planning stays on the current branch in the isolated blueprint checkout. After an
explicit integration sequence and green preflight, create sdp/actions-json-api from
the selected baseline in an isolated checkout. That becomes this plan's current
working branch for all phases, with one commit per meaningful milestone and phase
pushes. This refines Session0011's provisional branch suggestion; no code branch is
created while the baseline is unresolved. No main merge, publication, installation
replacement or editing other worktrees is authorized by this plan.

## Phases and milestones

| Phase | Milestone | Complete delivery and acceptance | State |
| --- | --- | --- | --- |
| ACT0 Preparation | ACT0-M1 | Reconcile shared-file ownership, KB051 identity and event provenance; record combined baseline and fresh green preflight | waiting |
| ACT1 Contract | ACT1-M1 | Finalize machine schemas/typed operation boundary and update validated SDL design sources on selected baseline; agree bounded consumer vocabulary | planned; prose proposal ready |
| ACT2 Read workflow | ACT2-M1 | Registry-derived actions catalogue, strict invoke stdin/file transport and project.discover through shared CLI service | planned |
| ACT3 Blueprint workflow | ACT3-M1 | Generate/retain, assess and list typed actions; preserve preliminary sources, immutable bundles, blocked results and domain diagnostics | planned |
| ACT3 Blueprint workflow | ACT3-M2 | Embedded assignment request, separately supplied principal, shared human/JSON handler, exact revision checks and safe retry reconciliation | planned |
| ACT4 Consumer verification | ACT4-M1 | Compiled-process tests, headless consumer and actual gh-sdp transport with isolated signed development distribution; review and regressions | planned |
| ACT5 Handoff | ACT5-M1 | Usage/examples, API/CLI contract, evidence and external XFMD handoff; reconcile card/Session and paused blueprint resume pointer | planned |

### ACT0 — coordination before product edits

Use the [T007 evidence](../../../Sessions/evidence/0011-concurrent-work/refresh-T007.md).
Select integration ownership/order for CLI, output_arguments, presentation, Contract,
SDPTool SDL sources and management/traceability records. Preserve original ledger
bytes and document identity reconciliation; do not silently reuse or renumber KB051.
This plan does not authorize resolving another agent's working copy. Keep design
changes in Actions documents while waiting. Recheck before the first product write.

### ACT1 — contract and model

Produce executable request/catalogue/response schemas and concrete examples covering
all five actions. Model operation registry, CLI/JSON adapters and existing services
under the actual SdpToolHost, not as fictitious separately deployed containers.
Validate .design sources through the Go parser and generate changed views with tools.
Ensure the proposed trust/cancellation/path/error rules have realizable service APIs.
Record necessary deviations in the contract before implementation. No separate
RequirementPlan or ArchitecturePlan is needed for this bounded change.

### ACT2 — useful read slice

Add an explicit input-reader entrypoint while preserving existing Run callers.
Factor typed registry/handlers without circular imports or shell invocation.
`actions --json` and `invoke --request -|FILE` work in the compiled binary; readable
catalogue and existing CLI behavior remain. Discover must use the same service and
return the same domain data. Catalogue must be usable without a project, with
contextual availability and producer identity. No manual discovery index is written.

### ACT3 — complete blueprint slice

Reuse Generate/GenerateRetained, Assess, AssignmentViews and ApplyAssignment. Do not
fork capture, semantic analysis, domain lifecycle validation or history storage.
Human CLI and JSON invocation use the same handler paths for included operations.
Preserve existing CLI flags/results; JSON envelope belongs only to invoke. Scoped
cancellation improvements must not falsify committed mutation outcomes. Assignment
principal comes from local process attribution/host injection, never request JSON.
Unsupported platform refusal must survive catalogue and invocation honestly.

### ACT4 — evidence at actual boundaries

Required checks (record results only when run):

- Table-driven malformed, duplicate/case-variant/unknown fields, overflow/depth,
  null/types, trailing data, unknown action/version, invalid selectors and paths.
- Registry metadata/handler consistency; deterministic metadata; absent context,
  unavailable platform, unsupported form schema, fallback icon and persisted action ID.
- Real compiled stdin and file requests, exact single JSON stdout, errors/exit codes,
  cancellation and broken-output behavior. No success claim after write failure.
- Human/JSON parity for discovery, generated bundle hashes and assessment; fresh
  non-Git fixture; no-write snapshots for reads/refusals.
- Retained generation and assignment create/update/retry, stale expectedEvent,
  forged request authority, changed source/evidence and lost-response reconciliation.
- Current SDPTool suite/race/vet as appropriate after scoped tests; meaningful CLI
  regressions for parallel governance/program additions on the combined baseline.
- Build actual gh-sdp wrapper route with stdin preserved and a signed nonproduction
  descriptor in isolated cache/config; pin both producer and wrapper identities.
- Independent review if authorized/available; otherwise explicitly record the
  verification limitation, never claim same-context checking is independent.

Headless consumer evidence does not prove native XFMD integration. No release assets
or published extension are changed. Treat failure/timeouts as evidence, not as passed.

### ACT5 — closeout

Update SDPTool README/Contract and XFMD integration guide, including exact commands,
request examples, availability, local trust limits and nonzero JSON result handling.
Remove proposed labels only for delivered behavior. Review backlog/onHold for related
scope before finalizing. Close only the scoped primary card after evidence; retain
consumer/runtime/release work with named owners. Session0008 remains paused until
owner resumes its BP2-A pilot review; this delivery cannot approve the pilot.

## Worklog and evidence

2026-10-10, Session0011 T007: plan registered as planned with contract proposal.
No product code, SDL source, executable schema or test was changed/run for ACT delivery.
Concurrency preflight still RED for integration; isolated planning records only.
