# Provisional work-domain and scoped identity contract

Status: **pilot v0; not canonical**

A work domain is a stable bounded area of product/process ownership. It may be
hosted beside other domains in one repository or alone in another. Domain
identity, record identity, repository location, and physical path are separate
concepts.

## Decisions for the twelve Issue #7 questions

1. **Several domains in one repository:** yes. Pilot v0 uses exactly one
   registry object per canonical GitHub repository; that object declares one
   or more active/moved domain entries. Partitioned registry objects are
   rejected provisionally because they create unsafe collision-check
   boundaries.
2. **Stable short namespace:** yes. Each domain has an immutable uppercase
   ASCII `key` of 2–12 characters (`[A-Z][A-Z0-9]{1,11}`), selected and
   collision-checked before use.
3. **Scoped IDs:** yes. New multi-domain work uses
   `<KEY>-FEAT|REF|FIX|STU|SLC-NNN`, for example `CORE-FEAT-001`.
4. **Unscoped IDs:** yes, but only for new work in a repository with exactly
   one active domain, a declared `defaultDomainUid`, and
   `newRecordIdStyle: unscoped`. Historic unscoped IDs remain valid after a
   repository later gains domains when marked `identityStatus:
   legacy-preserved` with source evidence; new work then uses scoped IDs.
5. **Prefix meaning:** for a scoped record, the key is part of its stable
   canonical visible ID, not a mutable display alias. Its globally durable
   identity is the pair `(domainUid, recordId)`.
6. **Machine declaration:** a repository registry declares `domainUid`, key,
   name, lifecycle/hosting state, owners, ID style, roots, issued-ID inventory,
   and move metadata. See `templates/work-domain-registry.template.json`.
7. **Cross-scope relations:** a relation to another domain MUST use an absolute
   object `{domainUid, id}`. A `keyHint` is not identity, but when supplied it
   MUST equal the stable key resolved from `domainUid`; a contradictory hint is
   invalid. Bare cross-domain IDs are invalid.
8. **Shared/common work:** shared work is an ordinary explicitly owned domain
   with its own UID, key, owners, roots, records, and reservations. `shared` is
   not a magical bucket or a permission to use ownerless paths.
9. **Concurrent reservations:** before work starts, every Issue assignment
   freezes a common integration base, domain UID, reserved IDs, private owned
   paths, symmetric shared touchpoints, dependency/conflict edges, and
   convergence/merge order. See [ConcurrentAssignments.md](ConcurrentAssignments.md).
10. **Future repository split:** the moved domain keeps its `domainUid`, key,
    and every record ID. The old registry retains a tombstone and successor;
    the new registry imports the issued/reserved-ID inventory and declares the
    active host. References continue to store the same `(domainUid, recordId)`.
11. **Namespace versus folder:** independent. Roots are declared repository-
    relative mappings and may change without identity change. Tools do not
    infer keys from folder names.
12. **Discovery:** a future Analyzer/gh-sdp reads the versioned registry from a
    manifest capability pointer, then follows its explicit roots/tombstones.
    Pilot repositories without a manifest place the registry at a stated
    Issue/assignment path. Recursive folder guessing is legacy discovery and
    must be reported as observed/ambiguous, never promoted automatically.

These are pilot decisions. The open choices that still require real pilot
evidence are isolated in [OpenQuestions.md](OpenQuestions.md).

## Identity layers

### Immutable domain UID

`domainUid` is a lowercase RFC 4122 UUID URN, for example
`urn:uuid:5f4d6d0e-8fab-4bc7-9794-0f7d9b0b16a1`. It is minted once, is never
derived from an organization, repository, folder, display name, or key, and is
never reassigned. This makes it portable across hosting moves and renames.

### Stable domain key

`key` is the human-facing namespace used in new scoped IDs. It is immutable
once the first scoped ID is issued. Keys must be unique after portable
normalization among all active and tombstoned domains in one repository.
No global key uniqueness is claimed; `domainUid` provides that boundary.

If a destination repository already retains a colliding key, the move MUST
stop for Steering. Neither domain may be silently renamed. A future mounting or
alias mechanism may be piloted, but it cannot replace either canonical key or
absolute reference.

An optional alias is a redirect for an old repository URL, root, or
human-navigation label. It MUST resolve to one canonical `domainUid` and key,
MUST NOT be used to mint record IDs, and MUST NOT replace the UID/key stored in
a relation. Tombstones retain aliases needed to follow historical locations;
aliases are never an escape hatch for a destination key collision.

### Record ID and absolute reference

For new records:

```text
scoped:   <KEY>-(FEAT|REF|FIX|STU|SLC)-NNN
unscoped:       (FEAT|REF|FIX|STU|SLC)-NNN
```

Optional grouping IDs use the same rule with `SPR` and `ITR`. A record's
`domainUid` MUST match its owning declaration. For scoped records, its prefix
MUST equal that declaration's key exactly. The absolute reference is:

```json
{
  "domainUid": "urn:uuid:5f4d6d0e-8fab-4bc7-9794-0f7d9b0b16a1",
  "id": "CORE-FEAT-001"
}
```

Within one domain, the normalized record ID may be issued or reserved only
once, across all current and historical repositories. A move imports the full
issued/reserved inventory before new allocation. Numeric gaps are not evidence
that an ID is available.

Every new form ends with exactly one three-digit allocation suffix (`-NNN`).
An undocumented form such as `FEAT-001-002` is not prospective pilot v0
identity. Empty IDs are invalid. Historical spellings are preserved exactly,
including spelling/case that would not satisfy the prospective grammar.

`issuedIds` is a structured inventory, never a list of bare strings. A
prospective member carries `id`, `status: prospective`, the allocating
`allocationIssue`, and its repository-relative `source`. It must satisfy the
current domain key/style/type grammar. `source` is a normalized, portable,
non-recursive repository-relative record path and equals the represented
record's own `source`; when local, it resolves to that exact JSON record. No
two prospective identities hosted by one canonical repository may claim the
same NFKC/casefold/slash-normalized source, and inventory-to-record binding is
one-to-one. A
preserved historical spelling such as `DBG-RF-001` carries `status:
legacy-preserved`, its original source, and exact provenance `{repository,
commit, path}`; it is not normalized into a new `REF` identity. Its
historical `authorityIssue` is canonical when an Issue existed. If none existed, it is
`null` and a nonblank `authorityMissingReason` states that fact—history is not
invented. Provenance repository and Issue values use the canonical GitHub
forms below, while provenance/source paths remain portable. Inventory IDs are
unique after NFKC/casefold normalization. A reservation for a newly created ID
must use the prospective member's immutable `allocationIssue`; another Issue
cannot reclaim it. That allocation fact is distinct from execution authority:
a later Issue may use an already-issued Feature/Refactor/Fix as `workRef`
without reserving it again when the owner lists that canonical Issue exactly
once in `issueAuthorities`. Each new Slice/Study/Fix reserved by that later
Issue still records that Issue as its own allocator.

## Canonical GitHub identities

Pilot v0 accepts only these authored forms:

```text
repository:     https://github.com/<owner>/<repo>
Issue:          https://github.com/<owner>/<repo>/issues/<positive-decimal>
pull request:   https://github.com/<owner>/<repo>/pull/<positive-decimal>
Issue comment:  https://github.com/<owner>/<repo>/issues/<positive-decimal>#issuecomment-<positive-decimal>
```

Numbers have no leading zero. Values are nonblank; the scheme/host and path
shape are exact. Arbitrary hosts/schemes, `.git` transport forms, query
strings, extra fragments, trailing slashes, zero, and path/number aliases are
rejected. Host, owner, and repository comparisons are case-insensitive, so a
case spelling cannot create another graph node or another registry partition.
One registry object owns each canonical repository in pilot v0; duplicate and
case-alias registry objects are rejected before default-domain, key, or root
checks. Every assignment Issue and planned PR belongs to the registry
repository that actively hosts its primary work domain. Canonical Issue
identity is also used for work/Study authorities, structured inventory,
reservation rows and edges/order, amendments, and Steering evidence.

## Declaration and discovery

A registry has its exact schema marker, exact `kind`, and `experimental: true`,
identifies the hosting repository in canonical GitHub form, and lists
declarations. Omission is an error, never a legacy/untyped downgrade. Pilot v0
allows exactly one registry object for that canonical repository. Each
declaration includes:

- immutable `domainUid`, stable `key`, name, explicit owners, and
  `newRecordIdStyle`;
- hosting state `active` or `moved`;
- active repository-relative `roots`, or an empty root list for a tombstone;
- `issuedIds` and optional reservation-set reference/digest;
- for a tombstone, canonical `successorRepository` and a portable local
  `successorRegistry` path resolved inside that successor repository;
- for the active successor, `predecessorRepositories`.

`issuedIds` on a tombstone is the frozen structured inventory at the move gate.
The active successor begins with byte-equivalent member content for that full
set, including status, authority, source, and provenance, and MAY append new
IDs after the move; it may never omit, rename, weaken provenance, or reuse an
ID from the tombstone.

At most one active declaration exists for a UID across a validated registry
set. Repeated declarations for the UID MUST preserve the key. A moved entry
cannot mint IDs or own paths. Tools resolve a relation by UID first, then ID;
URLs and roots are location hints.

The root mappings are allowed to cover different shapes, such as a complete
domain tree, selected record folders, or a nested task-local tree. Overlapping
active roots across domains are invalid unless ownership is narrowed into
non-overlapping mappings. Shared physical files belong to the explicit shared
domain or are assignment-level shared touchpoints; they do not acquire
ambiguous multi-domain ownership.

## Portable normalization and collision rules

All pre-work reservation checks use the same portable form:

### Keys and IDs

- apply Unicode NFKC, then Unicode casefold;
- keys/IDs for new work must also satisfy their uppercase ASCII grammar;
- compare the full normalized ID within one `domainUid`;
- compare normalized keys across active and tombstoned declarations in one
  repository.

### Paths

- apply Unicode NFKC before checking separators, glob syntax, colon, dot
  segments, reserved names, or any other path semantics, then casefold and
  compare `/` and `\` as equivalent separators;
- reject compatibility characters whose NFKC result introduces `:`, a glob
  token, `/`, `\`, or other invalid path semantics;
- require a repository-relative path; reject drive prefixes, leading `/`, empty
  segments, `.`/`..`, the conservative Windows-forbidden component set
  `< > : " | ? *`, reserved names, and every Unicode control/format/surrogate/
  private-use category (`Cc`, `Cf`, `Cs`, `Co`);
- reject invalid Unicode scalar content and non-finite numeric values
  recursively before canonical JSON encoding, including outer extension
  objects; a JSON-escaped lone surrogate or `NaN`/positive or negative
  infinity produces deterministic diagnostics and never an encoding exception
  or a usable reservation digest;
- trim trailing spaces/dots for collision comparison and reject a path whose
  normalized segment changes for that reason;
- reservations are either an exact file/directory path or one recursive
  directory expressed only as a final `/**`; other glob syntax is invalid;
- assignment-path collision exists when normalized paths are equal or one recursive reservation
contains the other. This catches case-only, slash-style, composed/decomposed
Unicode, and parent/child overlaps.

The same path routine governs domain roots, record/inventory/provenance
sources, Slice paths, assignment owned/shared/prohibited paths, and local
reservation/history/successor-registry paths. `allowedMutation` on every shared
touchpoint is a nonblank contract, and `prohibitedRepositories` contains only
canonical repository URLs.

Every active domain root is a recursive ownership tree even though registry
roots MUST NOT use assignment glob syntax such as `/**`. Therefore `SDP`,
`sdp/hsx`, `SDP\HSX`, and an
NFKC-equivalent spelling all overlap `SDP/HSX` when owned by another active
domain.

Two assignments may name the same shared touchpoint only when both declare the
same normalized path and the reservation set supplies an owner, dependency or
merge order, and a convergence check. A shared path cannot also be a private
owned path. Any other overlap blocks activation.

## Move/split protocol

1. Freeze new reservations for the domain at an exact integration base.
2. Export the domain declaration, issued/reserved IDs, semantic relations, and
   exact source candidate; validate their digest.
3. Create the destination declaration with the same UID and key, unchanged
   record IDs, predecessor repository, and new roots.
4. Replace the source's active entry with a tombstone naming the destination;
   do not delete its history or reuse its key.
5. Verify one active host, inventory equality, preserved absolute references,
   and no destination key/path collision.
6. Refreeze dependent assignments. External/moved relations retain UID plus
   original ID; only optional location hints may change.
7. Merge destination and tombstone changes in the declared convergence order.

Changing `CORE-SLC-007` to `NEWCORE-SLC-007` during a move is identity rewrite
and fails the pilot. Repository aliases, redirects, and folder changes are
location aids only.
