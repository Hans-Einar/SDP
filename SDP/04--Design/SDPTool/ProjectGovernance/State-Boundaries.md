# Current model, Blueprint and ProjectGovernance state responsibilities

Session0013 / WR2-M2, 23 October 2026. This explains existing implemented boundaries;
it does not merge their state machines or add new authority. The
[orchestration design](Orchestration.md), [PG design](Design.md),
[model design](../ModelGovernance/Design.md) and
[Blueprint lifecycle](../Blueprints/Assignment-Lifecycle.md) retain behavioral authority.

## One owner for each meaning

| Record / identity | Authoritative owner | Meaning and permitted consequence |
| --- | --- | --- |
| Model artifact ID, kind, head and source digest | SDPTool/model; model artifact metadata/checkpoints | WORK is mutable; frozen proposal/candidate/release preserves exact bytes. Model promotion is not assignment approval or code verification. |
| Blueprint task, analysis revision and retainedRevision | SDPTool/blueprints and immutable retained bundle | Reproducible NOW/TARGET/task context and analysis; retention/integrity is not readiness, acceptance or execution. |
| Blueprint assignment ID, expectedEvent and proof | SDPTool/blueprintstate plus blueprints/lifecycle.go; canonical project history | Revision-bound work/evidence lifecycle. A completed host turn or reported PG run cannot write an acceptance proof by implication. |
| Prepared assignment ID/revision, artifact pins and return criteria | SDPTool/governance preparation/lifecycle services; private checkout store | SAD/Master package independent of a run. Preparation alone leaves design acceptance unrecorded and dispatch unavailable. |
| Review revision and attributed design decision | governance/design_check.go and design_decision.go | Check canonical Blueprint assignment and recaptured live model/task/preparation. Trusted controller records an explicit decision against that exact basis; source drift invalidates applicability. |
| Activated execution run ID/revision | governance/activation.go and shared run core | Controller activates one approved Master package. Legacy Maintenance can start independently without models; activation uses the prepared ImplementationPlan path. |
| Attempt/thread/turn and observation/agent epoch | governancehost/controller with durable run observations | Observed execution, delivery and recovery identities. Epoch changes and uncertain sends must be reconciled; a role string is not a capability. |
| Worker report / report receipt | restricted governance MCP run_transition; report inbox projection | Report moves work to verifying. Inbox consumption currently reports not-tracked; a report is not reviewed delivery or owner acceptance. |
| Owner questions, replies and continuation references | governance attention/continuation services; trusted reply adapter | Durable attributed input, applicability and guarded fresh attempt. Answer acknowledgement is not approval of implementation or another phase. |
| Card/typed-plan lifecycle and Session roadmap | canonical ProjectManagement ledger, current records and authored Session | Canonical project work disposition. Private PG operational state does not automatically publish a completed card/Session. |

The private store and canonical project history deliberately hold different facts.
They are not duplicate global authorities. No automatic agreement between their
states is claimed. Current code in [design_check.go](../../../../SDPTool/governance/design_check.go)
reads both Blueprint assignment state and ModelGovernance snapshots;
[activation.go](../../../../SDPTool/governance/activation.go) checks the current
review decision before creating a run.

## Transitions to preserve

1. Resolve current card/plan/Session and actual execution authority. Inspection
   must not create a run or infer a decision from owner-role text.
2. For model-bound implementation, retain exact model/task/Blueprint context,
   prepare the Master package, recheck providers and record the actual design
   disposition. A stale source blocks activation. Ordinary Maintenance retains its
   separate model-free route.
3. Controller owns activation, start, host/attempt binding and dispatch. Worker
   MCP receives only the selected run/attempt/epoch capability; an unbound static
   MCP connection is inspection-only.
4. Work may report, fail, request scope change or ask for input within that
   capability. Changed authority/coverage and blocking attention prevent affected
   advancement; uncertain transport does not justify blind resubmission.
5. Checking and independent review follow reporting. Canonical publication and
   final acceptance remain unfinished PG gates. Do not complete Blueprint/card/
   Session records merely to mirror host or run state.

## Current versus remaining verification

WR2 checks source/provider integration, real MCP transport and existing negative/
recovery tests on the consolidated candidate. Its fixtures are not native agent
execution or independent review. The [execution evidence](../../../Maintenance/WorkstationReadiness/Execution-Evidence.md)
and [MCP adoption study](../../../Maintenance/WorkstationReadiness/MCP-Adoption-Study.md)
record the exact result. Host compatibility, leader report consumption, SAD dispatch,
capacity/recovery and canonical evidence/publication gates remain with PLAN-SDP-0018.
No broad code refactor is selected without a demonstrated conflicting rule.
