# PGF1-C1 — Durable Session binding and registration

Session0007 T042 / PLAN-SDP-0018 / KB038 and KB052. Owner authorizes autonomous
continuation of T041's selected design. BranchPolicy current / CommitPolicy
milestone; baseline 0457a2742ad3c56fe5b330d6c1884ab4e6512039. The Session records
preflight and boundaries. [Candidate hashes](evidence-PGF1-C1/candidate.json)
identify the tested source delta; documentation closeout is separate.

## Delivered behavior

An explicitly covered manual Session can be adopted as one version-pinned local
instance. Binding validation resolves the actual Session, typed plan, milestones,
source cards and authority document; rejects cycles, unsupported versions/digests,
invalid scope and incomplete coverage. Cards retain their canonical lifecycle.

Adoption persists revision 1 plus a prepared registration effect. Resume safely
appends that exact journal entry and commits revision 2. Repeating the same operation
returns its original receipt; restarting the process preserves identity. Interrupted
writes recover from the exact preimage/postimage, while external edits stop recovery.
Read/list expose stages, holds and current coverage without initiating work.

Six CLI commands and six discoverable JSON invocation actions use the same core.
Controller authority is injected by the host/local CLI, never granted by payload.
The additive state profile 0.13 preserves legacy run/review records and is refused
by the prior executable. Stored receipts are checked for identity, complete revisions
and exact registration intent. Reads/refused authority create no state.

## Verification

| Obligation | Evidence and result |
| --- | --- |
| Binding and authority admission | Focused tests reject unsupported/cyclic/duplicate/missing records, invalid paths/policies, unknown/duplicate/null JSON and unauthorized calls |
| Durable recovery | Before/after-head and post-journal injected failures recover with one journal entry; changed requests and stale revisions refuse |
| Preserve source/history | Plan, Session, binding, authority and ledger drift refuse without state advancement or overwriting external edits; legacy records and unrelated file preserved |
| Actual adapter journey | Fresh compiled CLI processes adopt/read/resume/list and retry; JSON invoke reads the same instance; 13 wire objects pass strict schema validation |
| Disposable project | [Lab checks](evidence-PGF1-C1/lab-checks.json): pg_lab/session-c1, same actual-process journey, deliberate stale-authority refusal/restoration; 110 prior tracked files/links unchanged |
| Compatibility | [Old reader](evidence-PGF1-C1/old-reader.json) refuses the new profile without changing HEAD; empty stores and existing profile coexistence tested |
| Regression and static checks | [Check record](evidence-PGF1-C1/checks.json) records final commands, outcomes and raw logs |
| Model scope | WORK :00002 still validates as canonical design-core/0.6; scoped diagnostic Blueprint constraints pass; source/generated hashes retained |

`session_cli_test.go` runs the compiled journey in a temporary project by default.
For a fresh lab fixture, set PG_SESSION_LAB_ROOT to an absent directory,
PG_SESSION_LAB_STATE outside the project and PG_SESSION_BINARY to the built engine.
The test refuses to overwrite an existing lab root. Session wire validation uses
Python jsonschema, as existing governance schema tests do. Final source/environment
and command results are pinned in the linked machine records.

## Blueprint and model disposition

The explicit C1 task limits code permissions and states that full target context
does not authorize later workflow implementation. Final diagnostic revision is
66990e264a4f2dfbcc23441e51d5cf69465c72e3d6e1f89b4b1ada1e0f373157.
[Task](evidence-PGF1-C1/blueprint-task.json),
[result](evidence-PGF1-C1/blueprint-result.json) and
[file hashes](evidence-PGF1-C1/blueprint-files.json) retain reproducibility.
The initial scoped generation occurred during implementation, not as a previously
completed independent handoff. Its code permission list was extended for actual
adapter coverage and regenerated. The observed comparison baseline and preliminary
status remain explicit. Same-context constraint inspection is not independent review.

WORK 29973936-d34d-4751-b15a-66c2d0f64a2b:00002 stays TARGET. C1 implements a subset
of its binding/persistence responsibility; the aggregate SessionWorkflowDelivery
correctly remains planned. No model source semantics changed in this increment.
Complete integration, exact-candidate review and accepted model release remain
required before closing the wider design-changing assignment.

## Deliberate limits and next increment

Every selected step remains at admission with an activation hold. Resume currently
recovers registration and checkpoints coverage; it does not resume a host attempt,
start agents or advance workflow. No automatic KB move, model gate, privileged
leader MCP, acceptance or final Session disposition is claimed.

C1 pins the whole management ledger, conservatively blocking even unrelated changes.
PGF1-C2 must implement checked lifecycle effects and scoped own/unrelated-effect
reconciliation before normal ongoing work can use this engine. One instance per
Session is enforced in the selected store, not globally across clones/stores. Shared
writer locking and exact postimages provide bounded recovery, not atomicity against
arbitrary same-account editors. Capacity remains explicitly bounded.

The fixture's inherited run/review is synthetic protocol setup. This turn makes no
native model call, independent product review, owner acceptance, merge, software
release or consuming-project installation claim. KB038/052 and Session S3/S6 remain
active; broader KB037, Blueprint implementation and the deferred Ponsse pilot remain
outside this delivery.
