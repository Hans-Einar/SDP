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

The scoped example contains Issue assignments 201 and 202 frozen against the
same integration commit and reservation-set digest. They reserve disjoint
domain IDs and paths. Both name the exact shared touchpoint
`SDP/Shared/contracts/interface.json`, its explicit owner, the merge order, and
the final convergence command. Changing path case/slashes, composing Unicode
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
- the new repository declares the same UID, key, issued IDs, and canonical
  `DBG-FEAT-001`, with new roots and predecessor metadata.

The external relation remains `{same UID, "DBG-FEAT-001"}`. A negative fixture
proves that rewriting it to a new key fails validation.

## Negative collision examples

`fixtures/negative/` covers:

- domain-key collision after casefold;
- within-domain record-ID collision;
- case/slash, Unicode-composition, and parent/child path collisions;
- an unqualified cross-domain relation;
- unresolved local work/relation references and nonreciprocal Slice ownership;
- invalid scoped-domain ID reservation;
- stale concurrent bases, duplicate reservations, asymmetric conflicts, and
  mismatched merge/convergence contracts; and
- record/key rewrite during a repository move.

Each fixture declares its exact expected diagnostic codes, so a negative case
cannot “pass” merely because the validator failed for an unrelated reason.
