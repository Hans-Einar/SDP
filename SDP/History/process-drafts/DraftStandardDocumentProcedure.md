# Draft Standard Document Procedure

Status: working draft  
Audience: agents first, humans second  
Method-Target-Version: 0.2.0 (unreleased)

## Purpose

The Standard Document Procedure (SDP) is a repository-local documentation,
execution and release system. It reduces drift between discussion, requirements,
design intent, implementation, verification, review, release and handoff.
Repository documents are the source of truth; chat memory is not enough.

## Baseline Folder Structure

Recommended project-local areas:

- `01--Mandate` through `07--Implementation`
- `Sprints`, `Refactors` and bounded `Fixes`
- `Releases`, `CodeReview`, `Verification`, `Traceability`, `Instructions`
- `SDP-project.manifest.yaml` and `RELEASE-NOTES.md`

The numbered documents progress from intent to implementation readiness.
Architecture/design work is horizontal; implementation is delivered in coherent
vertical Slices.

## Traceability Layer

Traceability normally contains:

- `CurrentIndex.yaml` — actual current release/work coordinates
- `Relations.yaml` — stable links among requirements, design, work, review,
  verification, migration and releases
- `Ledger.ndjson` — append-only state transitions

Ledger corrections are new events. A tag or GitHub Release event is prohibited
until that object really exists.

## Release And Development Identity

Released software and SDP Toolkits use Semantic Versioning:
`MAJOR.MINOR.PATCH`. Compatibility impact selects the increment.

The following are separate development coordinates and never overload the SemVer
core:

- active Sprint or Refactor
- Iteration
- Slice or bounded Fix
- optional revision
- Git commit and generated build state

Unreleased UI/build identities must be visibly marked as development builds.
Static installed facts, editable project release state and generated Git/build
facts are separate machine-readable records.

## Sprints, Slices And Fixes

There should normally be one active Sprint or Refactor, one Iteration and one
Slice. A Slice contract states goal, why now, expected files, invariants,
non-goals, traceability, verification and completion signal.

A Fix record is allowed only when a normal Iteration/Slice is clearly
proportionally excessive. It still requires bounded scope, release relationship,
verification, review and ledger entries. Capability, architecture or public
contract changes require normal planned work.

## Roles

- Master coordinates contracts, versions, traceability, delegation, evidence and
  stopping at the work boundary.
- Worker implements one assigned Slice or Fix.
- Reviewer uses fresh context and actual diff/evidence.
- Architect owns long-horizon structure and compatibility.
- Release/versioning/auditor/verifier task skills handle their concrete workflows.

No role may invent evidence or publication state.

## Design And Implementation Rule

> Design horizontally. Implement vertically.

Horizontal documents define boundaries, ownership and contracts. Vertical work
cuts through the necessary layers to produce one coherent verified capability.

## Verification And Review

A Slice/Fix or release candidate is not complete until records identify what
changed, linked IDs, exact evidence, independent review disposition, append-only
events and residual limitations.

Documentation work still validates links, schemas, YAML/JSON/NDJSON and
cross-file agreement. Product work uses real build/test/manual evidence. Release
work additionally validates version ordering, notes, manifests, immutable
history, migration, clean state and publication preconditions.

## Release Transaction

Release preparation and publication are two phases:

1. prepare, verify, review and approve a clean release-preparation commit without
   claiming a future tag or GitHub Release
2. after explicit human authorization, create the annotated tag and GitHub
   Release, then commit a truthful reconciliation record containing their real
   identities

Agents stop before publication unless explicitly authorized.

## Living Document Rule

When a repeated failure mode appears, add a concrete reusable rule here, in the
Framework, or in a versioned skill. Project-specific rules belong in project
records and `AGENTS-project.md`.
