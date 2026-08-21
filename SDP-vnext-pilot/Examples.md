# Example scenarios

Status: **pilot examples; not migration instructions**

All examples are validated by `validate_pilot.py`. Their UUIDs, URLs, commits,
Issue numbers, and IDs are synthetic except where the HSX handoff explicitly
labels observed evidence.

## Simple single-domain repository

[`examples/simple-single-domain.json`](examples/simple-single-domain.json)
declares exactly one default domain with `newRecordIdStyle: unscoped`. New work
therefore uses `FEAT-001`, `SLC-001`, and an illustrative `FIX-001`. The
assignment and Slice are accepted at one candidate while the durable Feature
remains `active`, demonstrating that one accepted assignment does not force a
multi-assignment capability to `delivered`. Even though the visible ID is unscoped,
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
the absolute pair `{domainUid, id}`; the optional `keyHint` is checked against
the resolved stable key. A bare
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
Issues 201 and 202 also demonstrate a valid symmetric conflict: 201 is active
and 202 is explicitly blocked. Making both active fails validation.

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

## Accepted Study-only assignment

[`examples/accepted-study-only.json`](examples/accepted-study-only.json)
contains an accepted assignment with no Slice. Its primary Study is accepted
at the same candidate with qualified verification, review, and Steering
evidence. Changing only the assignment or Study state/candidate is invalid.

`fixtures/positive/accepted-standalone-fix.json` derives a second valid
assignment from the simple example. Its low-risk zero-Slice Fix is delivered at
the assignment candidate and points to the evidence-qualified delivered
Feature it corrects; the accepted prerequisite assignment precedes it.

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
- assignment/Slice/Study candidate-state mismatch, unsatisfied prerequisites,
  and simultaneously active conflict endpoints;
- bare, unresolved, wrong-kind, duplicate, and key-hint-conflicting embedded
  Feature/Refactor/Study/Fix semantic edges;
- complete conservative portable-path failures including `|`, `"`, `<`, C1 and
  bidi controls, prohibited-path escape, and blank shared mutation;
- exact/case/slash/Unicode prospective-source collisions and missing
  inventory/record/materialized bindings;
- unknown v0 schema markers and duplicate same-row qualified reservations; and
- missing, fabricated, tampered, gapped, truncated, or repository-unresolved
  multi-refreeze history.

Each fixture declares its exact expected diagnostic codes, so a negative case
cannot “pass” merely because the validator failed for an unrelated reason.

## Issue #7 dogfood refreeze

`Steering/Assignments/ISSUE-007.yaml` is revision 3. Its immediate
`previousRevision` points to
`Steering/Assignments/History/ISSUE-007-revision-002.json`, sourced from exact
candidate `0de8a8957b3212404007160d90da47445d7b4e7b`. That snapshot links revision
1 at `f79e3dfc18c7a1650f1f3ae66167dda9b69692b4`, so the complete 1..2 chain is
retained. Revision 2 removed accepted Issue #5 from the live execution DAG;
revision 3 reserves exact `REV/VER-SDP-007-004` evidence paths and records why
the third review gate required another refreeze.
