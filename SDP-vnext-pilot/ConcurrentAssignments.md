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
   Duplicate object member names are rejected recursively before normalization
   or hashing; no first-wins/last-wins interpretation can become authority.
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

The repository driver resolves the exact validated Git `HEAD`, then resolves
every local assignment, reservation, record, and history source. Every
historical `sourceCandidate` MUST resolve as exactly that 40-hex commit, MUST be
a strict ancestor of the validated `HEAD`, and MUST be a strict ancestor of the
next revision candidate. Disconnected roots, sibling/non-ancestor commits,
descendants/future commits, the current commit itself, duplicates, reversals,
and out-of-order chains fail before cited bytes can claim history. The driver
strictly parses the historical assignment from `assignmentSource` and the
historical reservation path declared by those exact bytes, recomputes that
reservation digest, and compares canonical Issue/source/revision/reservation
identity. It also compares the actual revision sequence: Issue, assignment
source, and primary `workRef` never change; authorized Slice identities are
append-only; and every non-null historical `acceptedCandidate` survives
unchanged. For every snapshot it resolves exactly one reservation row for the
canonical Issue and requires its complete `issue`, `workRef`, authorized/active
Slices, reserved IDs, owned/shared paths, and dependency/conflict projection to
equal the canonical projection of the exact historical assignment. It also
requires exact historical base, merge order, and convergence equality. Updating
both historical files, their digest, and every snapshot pointer cannot conceal
a row divergence. `recordedContext` adds audit pointers but cannot weaken those
mandatory comparisons. Duplicate JSON members, missing/tampered historical
reservation bytes, digest disagreement, or a missing/duplicate/mutated row fail
closed.

The only compatibility normalization recognizes the exact early Issue #7 pilot
shapes and receives the exact historical revision. Revision 1 alone may use
its exact typed set-level edge/string-ID form; revisions 2–4 alone may omit
only `authorizedSlices`; revision 5 and later require the current row shape.
A late legacy form, mixed modern/legacy rows, a partial row, or a hybrid
top-level shape fails closed even if bytes, digests, context, and pointers were
coordinated. The normalizer maps only fields present in the allowed historical
bytes, requires them to equal the exact historical assignment, and never
infers authority from current state.
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
JSON, and active-work satisfiability rules.
Revision 6 retains the exact revision-5 assignment and reservation at
`dfdea4fe259a9e342d651254656bcfac5363d0b9`, reserves exact
`REV/VER-SDP-007-007` paths, and closes reservation epochs, immutable revision
semantics, reverse active projection, accepted Fix targets, relation kinds,
duplicate JSON members, and branch-Unicode robustness. A completed
historical prerequisite remains linked as evidence without pretending another
Issue Master is concurrently active.
Revision 7 retains the exact revision-6 assignment and reservation at
`0655002cbe5e14543b007525cdc3e3bad82b6816`, reserves exact
`REV/VER-SDP-007-008` paths, and requires exact historical row projection,
tuple-keyed epochs, and complete validation of every typed preparatory or bound
reservation set.
Revision 8 retains the exact revision-7 assignment and reservation at
`133cfaee9cb194b6181ac1e8fa9e1b474f79c00f`, reserves exact
`REV/VER-SDP-007-009` paths, anchors the chain to exact Git ancestry/order,
closes repository-wide preparatory collisions, and makes early compatibility
strictly revision-scoped.
Revision 9 retains the exact revision-8 assignment and reservation at
`a6a23d50e1bb807c9db350758aad1737c01f7846`, reserves exact
`REV/VER-SDP-007-010` paths, requires every bound/preparatory ID claim to exist
in durable inventory, introduces the pre-materialization `reserved` state,
forbids cross-epoch identity reuse, and keeps every preparatory authorized
Slice candidate null.

The reusable pilot `reservation-set` shape is in
[`templates/reservation-set.template.json`](templates/reservation-set.template.json).
Self-contained examples embed it in `reservationSets`; a repository record may
instead bind an external path, as Issue #7 dogfood does. Both routes use the
same generic validator. It computes the epoch key as `(reservation-set ID,
canonical digest)`, rejects duplicate objects for the same pair, resolves each
assignment by its exact pair, and compares this complete projection:

```text
Issue + primary workRef + authorizedSlices + activeSlices
+ qualified reservedIds + ownedPaths + sharedTouchpoints
+ dependsOnIssues + conflictsWithIssues
+ common integrationBase + mergeOrder + terminal convergence
```

A missing reference, unresolved set, incomplete row, digest mismatch, or any
assignment/reservation content difference blocks activation.
Every typed reservation object is fully validated before binding, including a
preparatory set referenced by no assignment yet. Exact base, nonempty unique
rows, qualified IDs/references, portable paths/shared values, local edge
endpoints, self/duplicate/cycle/overlap rules, complete dependency-consistent
merge order, and nonblank terminal convergence are mandatory. Being unbound is
not a semantic-validation bypass. A preparatory row's Issue belongs to the
active repository host; its primary domain and work record exist; that record
names the Issue exactly once; reserved domains, ID grammar/style, inventory
allocation authority, and shared owners resolve; and each authorized new Slice
remains in `reservedIds` with `acceptedCandidate: null`. Every reserved ID
resolves to exactly one `reserved` or materialized `prospective` inventory
member in its active owning domain and the row Issue is that member's immutable
allocator. Missing or legacy-preserved members fail. A reserved member carries
its prospective grammar/style and planned portable source, but cannot serve as
current work or acceptance evidence until promoted to a materialized
prospective record. Its ID, private-path, and shared-path claims are
compared with every nonterminal bound epoch and every other preparatory set.
Terminal bound epochs remain non-current writers, but their identities remain
permanently allocated and cannot be reintroduced by a preparatory/current
epoch. Several disjoint preparatory sets are valid; colliding or repeated
alternatives must be combined/refrozen or rejected.

## Current reservation epoch versus preserved history

Assignments remain durable after acceptance; they are not all current writers.
Pilot v0 groups assignments by their bound `(reservation-set ID, digest)`.
Within each group it validates one integration base, exact rows, dependency
DAG, symmetric conflicts, complete merge order, terminal convergence, private/
shared path collisions, and shared-touchpoint symmetry. A group is current when
at least one member is `proposed`, `active`, or `blocked`; a group whose members
are all `accepted`, `rejected`, `cancelled`, or `superseded` is historical.
There may be at most one current group in one repository validation document.
The ID is a stable coordination-record identity and MAY be reused across
terminal immutable epochs with different digests, objects, and bases. Those
objects remain distinct by the pair; two objects with the same pair are
ambiguous and invalid. An assignment pair must resolve exactly one object.

Accepted historical groups retain their exact base, reservation, paths, IDs,
and evidence as terminal graph facts. Their paths no longer collide with a
later epoch and they do not satisfy its current dependency, conflict, order,
sharing, or convergence gate. Their IDs remain permanently allocated in the
inventory and cannot be reserved by another epoch. Therefore a later Issue may
legitimately reuse a path after the earlier epoch is terminal without rewriting
the earlier assignment, but cannot reuse its identity. Two assignments in the
same current epoch still collide and stale-base/refreeze checks remain
mandatory.

Starting two branches first and attempting to reconcile colliding assignments
later is invalid. Reservation must precede concurrent implementation.

## Ownership and sharing

- `ownedPaths` are private write reservations. Within their reservation epoch
  they may not overlap another assignment's owned or shared paths.
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

Every current assignment has a unique Issue authority in its set.
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
