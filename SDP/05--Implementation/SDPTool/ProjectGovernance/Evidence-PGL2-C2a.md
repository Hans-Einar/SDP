# PGL2-C2a — Revision and retirement of prepared assignments

2026-10-18, Session0007 T018 / PLAN-SDP-0018 / KB-SDP-052. Owner continuation selects
the bounded lifecycle contract in [Orchestration](../../../04--Design/SDPTool/ProjectGovernance/Orchestration.md).
[Candidate](candidate-PGL2-C2a.json) identifies the integrated baseline plus exact
product/test/schema/driver bytes and compiled controller. This is local lifecycle
delivery, not completion of PGL2 or acceptance of KB052.

## Behavior and compatibility

The trusted local controller can revise a prepared package without changing its
ID, role or selection. Exact generation/revision checks and current coverage/pins
guard replacement. Scope replacement and the immutable receipt commit together;
failed replacement retains the old claim. Retirement increments revision, records
a reason and releases only the unlaunched preparation claim. It works with stale
or missing sources, retains the package/history, and is terminal for that identity.
It neither cancels a host run nor declares work complete. New work needs a new ID.

Exact retries preserve historical results across later changes and restarts;
changed content, cross-operation reuse, worker capabilities and wrong bindings
refuse. Existing nonterminal run/preparation claims still exclude overlapping
replacement. Retired preparations no longer block legacy starts. Local exclusion
does not cover other state stores or worktrees. Records and receipts retain their
existing bounded capacities; retirement is not deletion or garbage collection.

Successful lifecycle mutations explicitly upgrade the selected store to 0.3.
New readers retain 0.1/0.2 support; later preparation and legacy writes preserve 0.3.
C1 readers refuse the new state. Lifecycle requests/results and 0.3 inspection use
assignment wire profile 2. Original preparation/context requests and retry receipts
retain profile 1. Both paths still report dispatch unavailable and design acceptance
not-assessed. No provider or owner approval is inferred from pinned file contents.

## Verification and limitations

Core/service tests exercise source drift and refresh, identity preservation,
atomic scope replacement, overlap refusal, retirement with missing sources,
terminal refusal, immutable historical retries after restart, worker/no-write and
binding guards, strict wire decoding, shared operation IDs, concurrent revise/retire
CAS, before/after-HEAD faults and legacy scope/profile preservation.

The first compiled trial found that embedded request fields were incompatible with
the existing strict decoder. Request types were made explicit without changing
shared JSON validation; wire decoding regression tests cover both new requests.
Final checks and their exact scope are recorded in the candidate manifest.

Passing commands: `go test -race ./governance ./governancecontroller ./governancemcp
./cmd/sdp-controller`, `go vet` on those same packages, and `go test . -run Governance
-count=1` from SDPTool. ProjectManagement validation passes with 66 cards, 45 records,
4 lineage operations and 602 events; Toolkit validation passes. These checks do not
claim an independent review or repeat unrelated SDL/SDUI suites.

[Actual CLI requests/results](trial-PGL2-C2a.json) come from the maintained
`SDPTool/labs/project-governance/lifecycle_trial.py` driver and fresh subprocesses.
The separate fixture prepares SAD/Master, changes TARGET, refreshes Master pins and
scope, rejects conflicting replacement, retires with missing TARGET, proves
terminal behavior and immutable retries, then gives the released scope to a new ID.
An actual C1 controller refuses 0.3 without rewriting HEAD. Thirty-eight successful
request/result objects validate against the profile-specific JSON Schema definitions.
These are synthetic artifact/workflow checks, not a model task or provider integration.

The preview is `/home/warloc/git/pg_lab/previews/pgl2-c2a`; the read-only
`/home/warloc/git/pg_lab/assignments-lifecycle` helper displays readiness, revision,
identity, scope and change reason. The original owner project/runtime retains all
176 file/symlink entries. Existing C1/C6 previews and attempts are untouched.

## Remaining gates and integration handoff

Next: controller-bound design disposition and provider-aware ModelGovernance/
Blueprint revision checks before dispatch eligibility. PGL3–5 and unfinished PGI
evidence/review gates remain. No live agent, owner approval, independent review,
publication or main merge is part of this delivery.

Integrate on ACT0 ancestry (main 8aa16c2 or its successor), retaining canonical
history reconciliation. Writes are PG core/controller, its fixture driver and
owning records; Actions' generic API/CLI, provider services and dependencies are
unchanged. Refresh main and concurrent dirty/committed work before the next PR.
New management/traceability events append to current ledgers; imported historical
bytes and the original PGL2-C1 evidence are preserved.
