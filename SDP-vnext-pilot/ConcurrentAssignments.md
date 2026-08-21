# Concurrent Issue assignment contract

Status: **pilot v0; not canonical**

Concurrent Issue Masters are safe only after Steering freezes one shared
reservation set. Merely creating separate branches from approximately the same
default branch is insufficient.

## Per-Issue record

The plural structure is:

```text
Steering/Assignments/ISSUE-123.yaml
Steering/Assignments/ISSUE-124.yaml
```

Each record declares:

- its exact v0 `schemaVersion`, `kind`, `experimental: true`, its own portable
  repository-relative `source`, canonical Issue/comment amendment URLs, and
  positive integer `revision`;
- primary `workRef` as `{domainUid, id}`, the complete preserved
  `authorizedSlices` set, and its `0..1` currently executing `activeSlices`
  subset;
- `baseline.branch`, exact 40-character `baseline.commit`, and the identical
  `coordination.integrationBase` used by the reservation set;
- intended delivery branch and PR URL, not their mutable live state;
- `ownedPaths`, `sharedTouchpoints`, and reserved qualified IDs;
- dependency/conflict Issue URLs, convergence gate, and merge order;
- stale-base policy, required evidence, stop condition, and boundaries.

`CurrentAssignments` is a generated observation/read model and is not committed
by default. Closing, clearing, or advancing one assignment never edits another
assignment record.

## Reservation/refreeze protocol

1. Steering selects an exact common integration branch and commit.
2. Against that commit, enumerate every active/tombstoned domain declaration,
   issued/reserved ID, active assignment, and path mapping.
3. Normalize keys, IDs, and paths using [WorkDomains.md](WorkDomains.md).
4. Produce one reservation set that contains every concurrent assignment,
   qualified ID, private path, symmetric shared touchpoint, dependency,
   conflict, convergence gate, and merge order. Hash its canonical JSON bytes.
   Pilot v0 canonical bytes are UTF-8 without BOM from JSON with object keys
   sorted lexicographically, original array order retained, no insignificant
   whitespace, separators `,` and `:`, Unicode characters unescaped, and only
   finite JSON numbers. Parsing rejects `NaN` and positive/negative infinity;
   encoding uses the equivalent of `allow_nan=false`. Recursive invalid scalar
   checks include outer extension objects before any hash is returned.
5. Commit/publish the coordination record so every work branch can start from
   the exact commit containing the same reservation set. When default cannot
   yet receive it, use an explicitly named integration/coordination branch and
   draft PR; do not pretend it is default truth.
6. Each Master records that integration base and reservation-set digest before
   product work. Workers receive only their reserved IDs/private paths and the
   named shared protocol.
7. Immediately before candidate sign-off and merge/convergence, compare the
   current integration head and active reservation set. Any unaccounted
   integration-base change, domain registry change, overlapping reservation,
   or dependency-head change makes the assignment stale.
8. On stale state, stop writes, rebase/merge only under the Issue policy,
   recompute all normalization/collisions, update the common reservation set,
   increment affected assignment revisions, and obtain Steering refreeze.

A changed reservation digest/content cannot reuse an assignment revision.
Revision 1 omits `previousRevision` and `refreeze`; every revision greater than
1 carries both. The current assignment links its immediate predecessor by
portable snapshot path, revision, exact source candidate, and prior reservation
digest, while `refreeze` repeats that candidate and gives a nonblank reason.
Every retained snapshot has `previousSnapshot: null` at revision 1 or an exact
immediate-prior link thereafter. Revision N retains exactly revisions
`1..N-1`, with no gaps, truncation, or unreferenced same-assignment snapshots.

The repository driver resolves every local assignment, reservation, record,
and history source, verifies each source candidate exists as a Git commit,
loads the historical assignment from `assignmentSource`, and compares canonical
Issue/source/revision/reservation identity plus every recorded context pointer.
If repository resolution is unavailable for a declared revision chain, pilot
v0 fails closed with `REPOSITORY_DRIVER_REQUIRED`; field-consistent fabricated
history is never accepted.

Issue #7 demonstrates an important distinction: closed, Steering-accepted
Issue #5 is evidence/prerequisite context for this pilot, not a current active
reservation-DAG node. Revision 1 recorded it as a dependency; revision 2
refroze the current execution graph without it and recorded why. Revision 3
retains both earlier snapshots while reserving the fourth exact evidence gate
and adopting the third-pass validation corrections. Revision 4 retains the
exact revision-3 snapshot from the fourth-review candidate, reserves the fifth
independent evidence gate, and records the terminal-work, marker, global-graph,
surrogate, and path-boundary rework. Revision 5 retains the exact revision-4 snapshot from
`15a476dd76bc80de5573ab2abb65af8c963a9c0e`, reserves exact
`REV/VER-SDP-007-006` paths, and refreezes allocation/execution authority,
preserved authorized Slices, terminal aggregate closure, semantic DAG, strict
JSON, and active-work satisfiability rules. A completed
historical prerequisite remains linked as evidence without pretending another
Issue Master is concurrently active.

The reusable pilot `reservation-set` shape is in
[`templates/reservation-set.template.json`](templates/reservation-set.template.json).
Self-contained examples embed it in `reservationSets`; a repository record may
instead bind an external path, as Issue #7 dogfood does. Both routes use the
same generic validator. It resolves the set by ID, recomputes the canonical
digest, and compares this complete projection for every assignment:

```text
Issue + primary workRef + authorizedSlices + activeSlices
+ qualified reservedIds + ownedPaths + sharedTouchpoints
+ dependsOnIssues + conflictsWithIssues
+ common integrationBase + mergeOrder + terminal convergence
```

A missing reference, unresolved set, incomplete row, digest mismatch, or any
assignment/reservation content difference blocks activation.

Starting two branches first and attempting to reconcile colliding assignments
later is invalid. Reservation must precede concurrent implementation.

## Ownership and sharing

- `ownedPaths` are private write reservations. They may not overlap another
  assignment's owned or shared paths.
- `sharedTouchpoints` are exact or recursive shared write reservations. Every
  participant declares the same normalized path, explicit domain/file owner,
  allowed mutation, convergence verification, and merge order. A shared path
  cannot also appear in any participant's private paths.
- Within one assignment, normalized `ownedPaths` and `sharedTouchpoints` are
  disjoint from `boundaries.prohibitedPaths`. A prohibition that contains or
  is contained by an authorized write surface makes the assignment internally
  unsatisfiable and blocks activation.
- An active Feature/Refactor/non-standalone-Fix assignment cannot prohibit its
  own hosting repository and exposes at least one owned/shared write surface.
  It has exactly one active Slice in pilot v0, and that Slice also exposes an
  owned/shared write surface. Study-only and explicitly standalone zero-Slice
  non-implementation assignments may be pathless only when a nonblank
  `pathlessReason` mirrors the boundary documented by their Issue.
- Read-only dependencies need no path reservation but MUST be named with an
  exact candidate or contract reference when drift would matter.
- Discovery outside reserved areas stops the Worker. The Master either returns
  it to the owning assignment, requests a shared-path refreeze, or splits work.

## Dependencies, conflicts, and merge order

Every current assignment has a unique Issue authority in the set.
`dependsOnIssues` is a directed acyclic prerequisite edge to another Issue in
that set. A dependent candidate
cannot receive accepted state until the named prerequisite candidate/gate is
satisfied, represented in-set by an accepted prerequisite assignment.
`conflictsWithIssues` is symmetric and blocks concurrent activation
until Steering changes ownership/order or refreezes a safe activation plan;
it is a blocked relationship declaration, not permission to write
concurrently. Every conflict endpoint also resolves inside the set. An Issue
cannot conflict with itself. Dependency and conflict lists are mathematical
edge sets: each peer appears once, and the same peer cannot be both a
dependency and conflict because those semantics give contradictory activation
instructions.
Two symmetric conflict endpoints therefore cannot both be `active` or
`accepted`; one active endpoint plus an explicitly `blocked` peer is valid.
`mergeOrder` is a non-empty list containing every current Issue
exactly once and no other member; every prerequisite precedes its dependent.
The separate `convergence` object has a non-empty owner and executable command
plus `terminal: true`. This is the explicit final convergence gate, not
narrative “coordinate later.”

After each predecessor merges, successors compare their recorded base to the
new integration head. A changed base is expected, but it is still stale until
the successor refreezes. The convergence candidate receives fresh integrated
verification and independent review; passing isolated branches is not enough.

## Authored versus derived

| Authored durable contract | Derived timestamped observation |
|---|---|
| Issue/amendment URLs | Issue open/closed state, labels, assignee |
| common integration branch/base SHA | current integration and branch heads |
| planned delivery branch and PR URL | PR draft/open/merged state and merge SHA |
| domain/work/Slice references | check runs and current review objects |
| owned/shared paths and reserved IDs | changed-file diff and observed collisions |
| dependency/conflict/convergence/merge order | whether a dependency/merge gate is currently satisfied |
| required evidence and stop condition | latest CI result or generated current-assignment list |

Observed differences create diagnostics. They do not overwrite authored terms.

## Required stale-base outcomes

- Same exact integration head and reservation digest: continue.
- Integration advanced only by a declared predecessor and no collision: rebase
  or merge as authorized, rerun reservation validation, record refreeze.
- Integration/domain/path/ID state advanced unexpectedly: block until Steering
  refreeze.
- Source base is no longer remotely resolvable: block; a local object alone is
  not a recoverable common base.
- Base change alters requirements, interface, safety, or ownership: return to
  Study/architecture/design convergence before implementation resumes.
