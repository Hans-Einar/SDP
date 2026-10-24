# PGI4-C1 — Candidate evidence and guarded rework

Session0007 T035 / PLAN-SDP-0018 / KB038 and KB052. This implements the evidence
prerequisite in the authorized PG completion route. It does not close PGI4/PGI5,
the orchestration plan or the Session.

## Candidate and design

[Preflight](preflight-PGI4-C1.json) records the clean 090fb7f SDP baseline and clean
73ccee5 pg_lab baseline, observed worktree ownership and bounded scope. The
[candidate manifest](evidence-PGI4-C1/candidate.json) identifies the baseline plus
exact changed source hashes, toolchain and compiled controller. The
[contract](../../../04--Design/SDPTool/ProjectGovernance/Design.md#pgi4-c1--candidate-bound-evidence-and-rework-t035)
refines PG-A06/07/08/16/17/18 without adding model or code-derived context scope.

Implemented core and CLI operations capture selected source bytes and exact
executable bits, retain attributed external evidence, inspect present applicability
and return failed verification to ready. Rework retains Session/assignment identity,
preserves original results and revokes the prior agent capability. No new MCP
mutation authority is granted. Profile 0.10 extends the existing atomic state store;
current writers preserve it and the actual C8 reader refuses it.

## Verification

[Commands and results](evidence-PGI4-C1/checks.json) record the final full SDPTool
race suite and vet. Focused tests cover source/dependency creation, deletion,
executable-bit changes, exclusions, path/type/size limits, forged or stale requests,
authority/artifact drift, restart/exact retry, before/after-head crash, shared UUID
collisions, profile preservation, review-claim limits and guarded rework. Old agent
capability reuse is refused after rework. These are same-context verification tests,
not an independent review.

[Schema checks](evidence-PGI4-C1/schema-checks.json) validate actual core, compiled
CLI and lab requests/results against the application schema. The preliminary full
suite and first lab trial passed before the final executable-bit representation and
read-generation checks; they are retained under preliminary and do not substitute
for final-candidate evidence.

## Actual pg_lab workflow

[Trial runner](evidence-PGI4-C1/lab_trial.py), [summary](evidence-PGI4-C1/lab-summary.json)
and [wire capture](evidence-PGI4-C1/lab-wire.json) exercise the C8 native worker's
completed result in a separate copy of its operational store. No model is launched.

1. Capture FormatLabel and its explicit read-only dependencies. Run the actual Go
   tests and attach the passing log; work stays verifying.
2. Change an unselected unrelated file temporarily: evidence remains applicable.
3. Introduce a deliberate separator fault: the original candidate becomes stale.
   Exact attachment retry still returns its original historical receipt.
4. Capture the failing source, run the actual failing tests, attach that result and
   request rework. The run becomes ready exactly once, in the same Session.
5. Restore original sources and rerun the Go tests successfully. Both receipts
   remain historical because the original report has been returned for rework.

All 43 baseline regular files were preserved during each trial. The original C8
state HEAD is unchanged. The deliberate fault is a verification fixture, not a
regression in the original native implementation. Lab `./pg evidence` reads this
trial; existing status/leader/reports continue reading the original native trial.
The lab result commit is recorded in [lab delivery](evidence-PGI4-C1/lab-delivery.json).

## Limits and next delivery

Evidence is controller-attributed external documentation. PG does not infer that a
command ran or a review was independent from attached prose. The external harness
records the actual test processes separately. Explicit dependency selection is
bounded and does not prove complete behavioral coverage. Filesystem capture assumes
quiescent external writers; it does not claim a filesystem-wide transaction.

Completion still refuses without applicable independent review and checked canonical
publication. Native review provenance, publication/recovery and full PG verification
remain; PGL2/4 SAD/owner interaction and PGL3 parallel/active-host/capacity gates are
also retained. Existing SDL-only scope, Blueprint pause, Ponsse deferral and no
merge/release boundary are unchanged. Backlog/onHold was reviewed without selecting
additional work. No card or Session is closed by this increment.
