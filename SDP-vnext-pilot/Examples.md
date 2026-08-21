# Example scenarios

Status: **pilot examples; not migration instructions**

All examples are validated by `validate_pilot.py`. Their UUIDs, URLs, commits,
Issue numbers, and IDs are synthetic except where the HSX handoff explicitly
labels observed evidence.

## Simple single-domain repository

[`examples/simple-single-domain.json`](examples/simple-single-domain.json)
declares exactly one default domain with `newRecordIdStyle: unscoped`. New work
therefore uses `FEAT-001` and `SLC-001`. Even though the visible ID is unscoped,
references still carry its immutable `domainUid`, so adding another domain or
moving the original domain later does not make the old identity ambiguous.

## HSX-like scoped monorepo

[`examples/scoped-monorepo.json`](examples/scoped-monorepo.json) declares four
ordinary domains with synthetic UIDs and keys `HSX`, `DBG`, `AVR`, and
`SHARED`. The generic contract does not know those names; this example uses
them because Issue #7 asks for an HSX-shaped pilot.

New records use composite visible IDs such as `DBG-FEAT-001`,
`HSX-STU-001`, and `SHARED-REF-001`. The shared domain has named owners and a
root like every other domain.

The Debugger Feature has a `depends_on` relation to the HSX Study. The target is
the absolute pair `{domainUid, id}`; `keyHint` is advisory. A bare
`HSX-STU-001` would fail in this multi-domain repository.

## Concurrent Issue Masters

The scoped example contains Issue assignments 201, 202, and 203 frozen against
the same integration commit and a real embedded canonical reservation object.
Issue 203 owns `SHARED-SLC-001` and depends on the Debugger and AVR assignments;
the merge order lists all three Issues once before the explicit terminal
convergence. They reserve disjoint domain IDs and paths. All participants name the exact shared touchpoint
`SDP/Shared/contracts/interface.json`, its explicit owner, the merge order, and
the final convergence command. The validator recomputes the reservation digest
and compares every assignment projection to it. Changing path case/slashes, composing Unicode
differently, or reserving a parent directory does not evade collision checks.

## Cross-scope dependency

The `DBG-FEAT-001 -> HSX-STU-001` relation demonstrates that record paths and
repository URLs are not identity. If the HSX domain moves, the same relation is
still valid because the UID and original record ID are unchanged.

## Domain moves to a new repository

[`examples/domain-move.json`](examples/domain-move.json) contains two registry
views for one domain:

- the old monorepo retains a `moved` tombstone with no active roots and a
  successor registry; and
- the new repository declares the same UID, key, structured identity inventory,
  and proposed `DBG-FEAT-001`, with new roots and predecessor metadata.

The inventory also preserves historical `DBG-RF-001` with its exact source
repository, commit, path, and authority Issue. `DBG-ST-000` is a positive
truthful no-Issue case: it retains `authorityIssue: null`, an explicit reason,
and exact provenance rather than inventing authority. The move example
deliberately models no delivered work, Issue assignment, or Slice; it
demonstrates identity portability without inventing missing delivery evidence.

The external relation remains `{same UID, "DBG-FEAT-001"}`. A negative fixture
proves that rewriting it to a new key fails validation.

## Negative collision examples

`fixtures/negative/` covers:

- domain-key collision after casefold;
- within-domain record-ID collision;
- case/slash, Unicode-composition, and parent/child path collisions;
- an unqualified cross-domain relation;
- unresolved local work/relation/Slice-decision references, bidirectional
  owner/Issue/Slice cardinality, and Slice path authorization;
- invalid scoped-domain ID reservation;
- stale concurrent bases, duplicate Issue/ID reservations, asymmetric
  conflicts, digest/content mismatch, dependency cycles, incomplete
  merge/convergence contracts, self/duplicate/overlapping graph edges, and
  order violations;
- noncanonical GitHub repository/Issue/comment/PR identities, aliases,
  repository mismatches, and duplicate registry partitions;
- unsupported accepted/delivered/released claims, invalid minimum state/value/
  collection policy, and an unsafe zero-Slice Fix;
- malformed, whitespace, unresolved-format, or wrong-candidate accepted
  evidence and contradictory Study authority;
- structured-inventory duplication, malformed prospective IDs, missing legacy
  provenance, unsafe source paths/mismatches, truthful-null-authority rules,
  exact `-NNN` grammar, and cross-Issue reclaim;
- recursive domain-root collisions plus NFKC-introduced colon/glob/separator
  semantics;
- record/key/inventory rewrite during a repository move; and
- assignment revision/digest reuse against a durable prior snapshot.

Each fixture declares its exact expected diagnostic codes, so a negative case
cannot “pass” merely because the validator failed for an unrelated reason.

## Issue #7 dogfood refreeze

`Steering/Assignments/ISSUE-007.yaml` is revision 2. Its
`previousRevision` points to
`Steering/Assignments/History/ISSUE-007-revision-001.json`, which binds
revision 1 at exact source candidate `f79e3dfc18c7a1650f1f3ae66167dda9b69692b4`
and reservation digest
`sha256:50d6768eecce05e540c7674b34a78ba6b4432cabae6e079b88fb51e669bebd4d`.
The current reservation digest differs and cannot reuse revision 1. Closed,
accepted Issue #5 remains evidence/prerequisite context; it is not a live node
in the revision-2 reservation DAG.
