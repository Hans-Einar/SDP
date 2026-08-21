# Example scenarios

Status: **pilot examples; not migration instructions**

All examples are validated by `validate_pilot.py`. Their UUIDs, URLs, commits,
Issue numbers, and IDs are synthetic except where the HSX handoff explicitly
labels observed evidence.

## Simple single-domain repository

[`examples/simple-single-domain.json`](examples/simple-single-domain.json)
declares exactly one default domain with `newRecordIdStyle: unscoped`. New work
therefore uses `FEAT-001`, `SLC-001`, and an illustrative `FIX-001`. The
assignment aggregates its authorized accepted Slice while the assignment and
Slice retain independently qualified candidates. `activeSlices` is empty after
acceptance, while the durable Feature remains `active`, demonstrating that one accepted assignment does not force a
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

The Debugger Feature authors its `depends_on` relation to the HSX Study once in
its embedded `relations` list. The global graph does not repeat it in the
top-level surface. The target is the absolute pair `{domainUid, id}`; the
optional `keyHint` is checked against the resolved stable key. A bare
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
truthful no-Issue case: it retains historical `authorityIssue: null`, an
explicit reason,
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

Additional positive controls close delivered Feature, Refactor, and sliced-Fix
owners over accepted Slices and assignments while deliberately allowing the
aggregate owner evidence candidate to differ from its historical Slice
candidate. Another control proves two distinct supported relation types for the
same source/target are not a duplicate, and a path-boundary control keeps
authorized and prohibited surfaces disjoint.

Positive lifecycle controls exercise the formerly contradictory
cardinalities: one delivered Feature has two Issues/two accepted Slices at
distinct candidates; one Issue aggregates two sequential accepted Slices at
distinct candidates; two accepted Issues retain distinct reservation sets and
bases; two and three terminal epochs reuse one stable reservation-set ID with
distinct immutable digests; a current later Issue coexists with accepted
history and legitimately reuses an earlier path; a complete unreferenced
preparatory reservation validates before binding; and a self-contained
exact-byte revision-2 chain has
its actual revision-1 assignment already containing accepted `SLC-001` before
the later revision appends `SLC-002`. The embedded-history driver exists only
inside the fixture harness and cannot masquerade as repository-local history.
Further controls
exercise an acyclic cross-scope diamond/mixed embedded-top-level chain, a
finite numeric digest extension, and documented pathless Study-only and
standalone-Fix assignments.

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
- missing schema/kind/experimental markers, including combined downgrade
  attempts, for every current v0 object kind and history snapshots;
- explicit `legacy-preserved` compatibility inventory attempting to satisfy a
  current reservation gate;
- terminal delivered/released Feature, Refactor, and sliced-Fix closure over
  accepted Slices/assignments, evidence-qualified Study decisions, and gating
  `depends_on` targets;
- one global semantic-edge set, including cross-surface duplicates, self edges,
  unsupported/candidate/extra fields, and a legitimate distinct-edge control;
- two/three-hop cycles in the global `depends_on`, `supersedes`, `refines`, and
  `requires_revision` DAGs, including cross-scope embedded/top-level mixes;
- allocation/execution-authority separation, missing/rewritten authorized
  Slice history, non-subset/multiple active Slices, and wrong Slice allocation;
- terminal missing/active Issue authority, standalone/sliced Fix closure, and
  empty/proposed/unaccepted corrected-behavior targets;
- active host-repository prohibition and missing implementation write surfaces;
- owned/shared versus prohibited path overlap, a JSON-escaped lone surrogate,
  and recursively nested `NaN`/positive/negative infinity through bound
  reservation extension objects; and
- missing, fabricated, tampered, gapped, truncated, or repository-unresolved
  multi-refreeze history, including Issue/source/work replacement, accepted
  Slice deletion/candidate rewrite, and historical reservation absence,
  tampering, or digest rewrite;
- coordinated historical assignment/reservation digest and pointer rewrites
  that delete or rewrite an authorized Slice, replace work, or remove/change
  reserved IDs, paths, and edges;
- duplicate exact reservation epoch objects, wrong digests among same-ID
  epochs, two same-ID current groups, and unreferenced sets with an invalid
  base, missing/duplicate row, self/cyclic dependency, incomplete/incorrect
  order, or blank/nonterminal convergence;
- multiple current reservation epochs while allowing preserved accepted sets,
  reverse active-Slice projection, and active Fix targets with no accepted
  represented outcome;
- wrong source/target kinds for top-level `corrects`, `owned_by`,
  `independent_of`, and `informs`; and
- exact raw duplicate-member assignment/reservation JSON plus a JSON-escaped
  surrogate branch, all rejected with stable diagnostics.

Each fixture declares its exact expected diagnostic codes, so a negative case
cannot “pass” merely because the validator failed for an unrelated reason.

## Issue #7 dogfood refreeze

`Steering/Assignments/ISSUE-007.yaml` is revision 7. Its immediate
`previousRevision` points to
`Steering/Assignments/History/ISSUE-007-revision-006.json`, sourced from exact
candidate `0655002cbe5e14543b007525cdc3e3bad82b6816`. That snapshot links revision
5 at `dfdea4fe259a9e342d651254656bcfac5363d0b9`, revision 4 at
`15a476dd76bc80de5573ab2abb65af8c963a9c0e`, revision 3 at
`60978c306f7ce3c09033603abf3300839b2a251d`, revision 2 at
`0de8a8957b3212404007160d90da47445d7b4e7b`, and revision 1 at
`f79e3dfc18c7a1650f1f3ae66167dda9b69692b4`, so the complete 1..6 chain is
retained. Revision 7 reserves exact `REV/VER-SDP-007-008` evidence paths and
records why the seventh evidence gate required another refreeze.
