# PGL2-C1 — Durable assignment preparation

2026-10-10, Session0007 T015 / PLAN-SDP-0018 / KB-SDP-052. Owner continuation
selects the [bounded preparation contract](../../../04--Design/SDPTool/ProjectGovernance/Orchestration.md#pgl2-c1-preparation-contract--2026-10-10).
[Candidate](candidate-PGL2-C1.json) identifies baseline 02e8c9d plus exact code/test/
schema/driver hashes and the compiled controller; Go 1.27.1 on Linux amd64.
This is C1 delivery; full PGL2, execution and owner acceptance remain unfinished.

## Delivered boundary

Shared core coverage accepts SAD DesignPlan/ArchitecturePlan or Master
ImplementationPlan through a separate preparation path. The original Maintenance
route retains its plan restrictions. Only a trusted local controller capability can
prepare work; a role label or worker capability cannot grant that right. Inputs pin
canonical coverage, role-specific artifact files, writable/excluded paths and return
criteria. These file references are not provider validation or owner design approval.

Immutable prepared IDs and operation receipts share the existing lock/generation
store. Generation CAS, retry-content comparison, two-way legacy receipt collision
checks and conservative scope conflict checks preserve identity on retry/restart.
New preparation refuses overlap with prepared/nonterminal legacy work in this same
store; legacy starts also honor prepared scope. Nonliteral legacy scope cannot be
used to evade those checks. Other stores, including another directory for the same
checkout, remain outside this local exclusion. No distributed/global claim exists.

Inspection reports stale canonical/artifact inputs without rewriting them. Results
always say dispatch unavailable and design acceptance not assessed. Old binaries
refuse the new operational 0.2 store; the original 0.1 owner-lab store remains intact.
No actual approval, task launch, retirement, phase continuation or model API adapter
is delivered by storing a design-disposition reference.

## Verification

Affected core/controller/MCP/command race tests and vet are recorded in the final
closeout. Tests cover unauthorized/no-write preparation, wrong binding, plan-family
separation, symlink/glob refusal, exact retries with changed inputs, concurrent stale
generation, before/after-HEAD crash recovery, overlapping/disjoint scope, legacy run
coexistence, shared operation identity and stale authority/TARGET.

The maintained `preparation_trial.py` runs the compiled controller against a separate
synthetic project through fresh subprocesses. [Actual requests/results](trial-PGL2-C1.json)
show two disjoint prepared assignments, stable retry identity, third-scope refusal,
TARGET change staling only Master, a fresh stale-pin refusal and refusal by the old
C6 binary after store upgrade without changing HEAD. Sixteen captured requests and
results validate against the new schema. These are actual local CLI/service tests
with fixture artifacts, not model execution, ModelGovernance/Blueprint integration,
independent review or owner acceptance evidence.

The preview is `/home/warloc/git/pg_lab/previews/pgl2-c1`; `pg_lab/assignments` displays
its result. The script refuses an existing destination. The original owner's project
and runtime retain all 176 file/symlink entries after new-binary inspection. Shared
untracked PG files still match prior PGI2 hashes; no other worktree was changed.
The new preview retains a copied controller, driver, evidence and candidate manifest.

## Remaining work and handoff

C2 must implement actual controller-bound design disposition, provider-aware revision
resolution, assignment revision/retirement and dispatch-boundary checks. Role-skill
adoption is sequenced with the integrator-owned distribution baseline before launch.
PGL3–5 and unfinished PGI review/evidence gates remain. Shared history/identity
reconciliation is separate; colliding MAINT-SDP-0015 is qualified as PLR1 from 5141da7,
not release RSP2. See [coordination handoff](Handoff-PGL2-C1.md).

Closeout: final `go test -race ./governance ./governancecontroller ./governancemcp
./cmd/sdp-controller`, corresponding `go vet`, and compiled legacy CLI regression
`go test . -run Governance -count=1` pass. Sixteen new preparation wire captures
and two existing provenance captures pass their schemas. Toolkit validation and
management replay pass (62 cards, 39 records, 531 events). Historical ledger
prefixes, exact source/binary evidence identities and maintained local links are
checked before the milestone commit. Same-context diff review preserves the stated
C1 boundary; no independent review is claimed.
