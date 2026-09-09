# SDP vNext provisional pilot contract

Status: **provisional / experimental / pilot-only**

This directory is the bounded output of
[Hans-Einar/SDP Issue #7](https://github.com/Hans-Einar/SDP/issues/7). It turns
the Steering-accepted Issue #5 study into a contract that real repositories can
pilot. It does not change the canonical Toolkit, freeze a future schema, or
claim compatibility with a later vNext release.

The words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** in this directory state
pilot rules only.

## Reading order

1. [WorkflowContract.md](WorkflowContract.md) defines work ownership,
   cardinalities, evidence state, reopening, Steering, and release inclusion.
2. [WorkDomains.md](WorkDomains.md) defines scoped monorepo domains, stable
   identity, references, discovery, moves, and collision rules.
3. [IssueContract.md](IssueContract.md) is the reusable human-readable GitHub
   Issue contract.
4. [ConcurrentAssignments.md](ConcurrentAssignments.md) defines per-Issue
   assignments, reservation/refreeze, stale-base handling, and convergence.
5. [Profiles.md](Profiles.md) keeps ceremony proportional.
6. [Examples.md](Examples.md) walks through the machine-readable examples.
7. [HSX-Pilot-Handoff.md](HSX-Pilot-Handoff.md) is the concrete first real
   pilot handoff; it is evidence about HSX, not a change to HSX.
8. [OpenQuestions.md](OpenQuestions.md) records choices intentionally left for
   real pilots and later schema work.

`templates/` contains eight illustrative pilot-only record shapes, including a
reusable `reservation-set`. `examples/`
contains positive scenarios. `fixtures/positive/` contains lifecycle,
proportionality, DAG, digest, and terminal-closure controls, and `fixtures/negative/`
contains deliberately invalid cases and their required diagnostic codes. JSON is used as the machine-readable
pilot encoding because Python's standard library can parse it deterministically
in strict mode (`NaN` and infinities are rejected and canonical encoding uses
`allow_nan=False`);
this does not decide the eventual canonical serialization. The existing
`Steering/Assignments/ISSUE-007.yaml` is JSON-compatible YAML for the same
reason.

Each self-contained assignment example embeds the exact canonical reservation
object it hashes, so validation recomputes the digest and compares every Issue,
work/authorized/active-Slice reference, ID, private/shared path, dependency/conflict edge, merge
position, and convergence field. Negative mutation fixtures are applied to
deep copies of positive examples and assert the exact diagnostic set.

Issue #7 dogfoods the model through `Studies/STU-007.json`, the compact
assignment, the single default-domain declaration, and
`Steering/Reservations/RSV-ISSUE-007-001.json`. The assignment binds the real
reservation file by canonical JSON SHA-256; the validator resolves the Study,
domain, paths, IDs, base, dependency, and convergence contract together.
Assignment revision 10 retains the complete revision-1 through revision-9
snapshot chain under `Steering/Assignments/History/`. The generic repository
driver reconstructs all nine exact Git candidates plus each candidate's exact
historical reservation, rehashes them, and requires the one canonical Issue
row to equal the exact historical assignment projection. It rejects missing,
fabricated, tampered, gapped, truncated, identity-replacing,
evidence-erasing, or coordinated assignment/reservation divergence. Every
source candidate is also an exact strict ancestor of the validated Git `HEAD`,
and revision candidates form one strictly chronological ancestor chain.

## Validate

From the repository root, run:

```powershell
python SDP-vnext-pilot/validate_pilot.py
```

The validator uses only the Python standard library. It validates templates,
positive examples and mutation controls, exact negative fixtures, bidirectional owner/Issue/Slice
cardinality, canonical GitHub identities/repository coherence, qualified
accepted evidence, identity inventory/source/provenance/no-reuse, reservation
digests and tuple-keyed reservation-epoch concurrency graphs, complete
standalone validation of bound and preparatory typed reservation sets,
repository/domain/inventory/current-epoch coherence for preparatory rows,
durable reserved-before-materialization allocation, cross-epoch ID no-reuse,
and null-only preparatory accepted candidates,
portable roots/paths and branch names,
cross-record acceptance/prerequisite/conflict state, qualified embedded
and top-level semantic edges as one global set with relation-specific DAGs,
allocation versus execution authority, complete authorized-Slice history,
terminal aggregate authority/affected-work closure, strict finite JSON,
recursive duplicate-member rejection before canonical hashing,
evidence-qualified Slice decisions and gating dependencies, unique materialized
sources, mandatory exact v0 schema/kind/experimental markers, invalid Unicode
scalar robustness, active implementation satisfiability,
authorized/prohibited path disjointness, assignment revision
history with real-Git ancestry/order controls and exact revision-scoped early
compatibility, dogfood external binding, local Markdown links, required status
markers, and trailing whitespace.

## Boundary

This corpus does not modify or authorize changes to HSX, gh-sdp,
SDP-Analyzer, PR #4, canonical Toolkit contracts, Toolkit release state, tags,
or GitHub Releases. Tool behavior described here is a future consumer contract,
not an implementation in this Issue.
