# PGF1-C2c4 — Reviewed quiescent Session rebind

T050 implements explicit reviewed revision of a remaining Session binding while
retaining the same Session, plan, workflow, original adoption and completed evidence.
This is the quiescent subset of replan: it does not abandon an unfinished selected
step, retire a running/uncertain host, remove existing work or migrate the workflow.

## Behavior and authority

CLI/JSON and explicitly scoped Session MCP share `session-rebind-context` and
`session-rebind` core operations. Preview identifies the current instance revision,
old binding digest and proposed full context. Preparation requires an external,
hash-pinned passing review naming those exact values. It is controller-attributed
review, not native independent review or human acceptance. Caller identity cannot
be supplied through request prose or JSON.

The revision retains every existing step/card reference. Completed step bindings,
required review/owner policy and protected paths cannot weaken. Canonical metadata,
board/history consistency, unchanged instructions, completed cards, old evidence and
existing PG journal tails are checked. The prior management ledger remains an exact
prefix. Pending publication, unfinished selected work and conflicting shared/local
claims refuse admission. At least one remaining step is required; otherwise use
normal completion or explicit successor transfer.

Preparation saves exact input preimages and deterministic postimages. Resume uses
the shared writer reservation and checked canonical publisher to append a rebind
artifact, plan/Session/live-card journals and management events. No card status is
changed by rebind. The operational commit archives the prior binding epoch and clears
current selection/assessment/assignment/owner-gate state. Fresh selection and assessment
are mandatory; completed predecessors still satisfy dependencies. Original receipts
and old runs remain addressable. Repeated binding revisions are flat, not recursive.

## Verification and candidate

[Candidate](evidence-T050/rebind-candidate.json) records baseline, changed-file and
compiled-binary hashes, toolchain, actual check outputs and the durable lab commit.
[Scoped task](evidence-T050/rebind-blueprint-task.json) passes diagnostic Blueprint
constraints against unchanged unreleased SessionWorkflow WORK :00002. The target
model remains unaccepted/unreleased; this is not a new NOW or release claim.

Service tests cover revision before selection, repeated epochs, revision after
committed disposition and after next-step advancement, shared contributions,
completed whole cards, fresh assessment/execution and overall closure. The completed
card is unchanged by rebind; later overall closure may append its own audit journal.
Negative tests cover review pins/verdict/context, policy/identity changes, missing
scope, modified instructions/artifacts/journals, shared claims, prepared external
edits, conflicting retry identities and forged retained snapshots/profile downgrade.
Fault injection covers preparation HEAD boundaries, writer reservation, every
canonical file effect, and final HEAD boundaries with fresh-service retry.

The compiled test compares CLI, JSON invocation and MCP preview, applies through CLI,
replays and commits through MCP, and repeats through CLI without changing HEAD.
Wrong-instance and forged-authority calls refuse. Observed objects pass the shared
wire schema. Focused race checks and vet pass. The full SDPTool regression passes at the frozen product candidate; no mixed-source
run is used as evidence. Complete outputs and durations are retained in the candidate.

## Durable pg_lab evidence

Lab commit a2dbe0e retains `session-pgf1-c2c4-rebind`, its explicit synthetic review,
wire captures and final result. External private state is
`/home/warloc/.cache/sdp-workstation/T050/lab-rebind-state`.
[Result](evidence-T050/rebind-lab-result.json) and
[validation](evidence-T050/rebind-lab-validation.json) show committed effects/current
coverage, one card/one plan/four canonical events and all 615 earlier tracked paths
unchanged. The actual previous binary refuses the new `rebindings` field. The new
binary reads the earlier transfer store without changing its HEAD.

No native agent was launched in this fixture. Its review does not qualify as independent
PG implementation review or actual owner acceptance. Earlier native evidence retains
its earlier candidate; it is not relabeled as final-source verification.

## Remaining obligations

Operational profile 0.22 is additive and preserves older supported reads. Limits
remain eight binding epochs, existing receipt/store bounds and bounded effect sizes;
there is no rollover. In-flight replan/retirement, workflow migration, removed-scope
disposition, controller-owned leader scheduling, SDL WORK/Blueprint/merge/release
integration and combined native recovery/acceptance remain in the existing plan.
Independent PG review, actual owner trial and final scope disposition remain open.
PLAN-SDP-0018 and KB038/052 stay active/in-progress. Nothing here authorizes merge,
software publication or model acceptance/release.
