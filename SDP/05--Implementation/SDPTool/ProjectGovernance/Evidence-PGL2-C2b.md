# PGL2-C2b — Provider-backed design review basis

2026-10-18, Session0007 T020 / PLAN-SDP-0018 / KB-SDP-052. This delivers the read-only
design check selected in [Orchestration](../../../04--Design/SDPTool/ProjectGovernance/Orchestration.md).
[Candidate](candidate-PGL2-C2b.json) pins baseline, exact source bytes, binaries and
check outputs. Actual owner disposition and dispatch remain unfinished C2 work.

## Behavior and boundary

The standalone controller checks an unretired Master preparation against the
canonical Blueprint assignment/event, governing plan, retained bundle, live
ModelGovernance NOW/TARGET captures, exact task bytes and allowed code paths.
Two provider reads and preparation/generation rechecks reject observed drift.
Current results include a reproducible review digest; blocked results have explicit
conditions and no digest. All results say designAcceptance not-recorded and dispatch
unavailable. No operational writes, profile upgrade or new MCP privileges occur.

This reuses exported provider APIs without changing provider code. The full live
model capture catches an uncommitted TARGET change even when its metadata-file pin
is unchanged. Coverage, unknowns, preliminary WORK status and Blueprint evidence
readiness remain visible. They do not establish human acceptance or code conformance.

The provider checks at most two distinct live bindings in a view. A selected binding
beyond that budget reports unknown and the PG check blocks. The test preserves this
restriction; focused provider verification is not invented here. Writable paths must
match literal task permissions exactly. No wildcard expansion, global exclusion,
source locking or exhaustive impact analysis is claimed. Future disposition/dispatch
must repeat the check against current facts.

## Verification

Passing checks: race tests and vet for governance, governancecontroller, governancemcp
and cmd/sdp-controller; the complete SDPTool `go test ./...` suite; management and
Toolkit validators. These are same-context implementation checks, not independent
review or native client acceptance.

Service tests use actual model creation, retained Blueprint generation and canonical
assignment creation. They cover changed sources with unchanged metadata, wrong/copy
pins, task drift, broadened scope, stale event, retirement, provider hold, retained-file
tampering, concurrent observation drift, changed preparation review digest, generation
guards and provider budget exhaustion. A malformed lowercase event ID in the extra
fixture was corrected to the existing history contract; production validation was
not weakened.

The maintained `design_trial.py` runs the built binaries in fresh subprocesses.
[Actual requests/results](trial-PGL2-C2b.json) prove current/repeated checks, unchanged
project/state during reads, TARGET drift and restoration of the identical review
digest. Eight design-check request/result objects validate against the JSON Schema.
No model task was launched and no owner approval was recorded.

The isolated preview is `/home/warloc/git/pg_lab/previews/pgl2-c2b`; the read-only
`/home/warloc/git/pg_lab/design-check` helper displays its blocked result, with TARGET
deliberately changed. `--json` exposes the complete result. All 176 original owner
project/runtime file/symlink entries retain their bytes/targets. Other previews and
installed tools are unchanged. This is a reproducible synthetic workflow using real
provider services, not evidence of a completed owner trial in a production project.

## Integration and next step

The baseline includes Actions main e8760d3. C2b writes PG core/controller, its lab
drivers and owning documents only; provider implementations and generic Actions
dispatch are untouched. The refreshed worktree preflight sees XFMD viewpoint work
at ecd9dcd with no scoped PG overlap. Preserve canonical ledger bytes and reconcile
future concurrent appends during integration. No main merge or release is included.

Next: trusted owner-input disposition bound to this review digest, with invalidation
on changed design and explicit dispatch prerequisites. The Project Leader dialogue,
multi-assignment execution, durable inbox and full overview remain later milestones.
