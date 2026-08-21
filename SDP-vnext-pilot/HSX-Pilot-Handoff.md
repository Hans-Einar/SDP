# First real HSX pilot handoff

Status: **provisional pilot handoff; HSX remains unchanged**

## Evidence to anchor, without promoting it

The accepted Issue #5 report
[`SDP-usage-analysis/repositories/HSX.md`](../SDP-usage-analysis/repositories/HSX.md)
records two simultaneously true states:

- merged/default truth is
  `Hans-Einar/HSX main@e374da88f4dd470bad2d8ec6a2a14f1ce367e40e`
  (tree `fa16ce642f437e1fe8f9be0398a9424f6d389932`), which contains no
  `SDP/`; and
- evolved SDP practice is unmerged branch evidence on a lineage beginning at
  `Implementation/vscode@a1daa1c62605c44ac67e58e2b71320006f73cdd9`.

The frozen study checkpoints on that lineage are:

```text
sdp/debugger-gap-analysis@e5a50ab45acdcb515ccd3602ce99487bd668cdfd
codex/dbg-rf-001@977da1c19117e805a341a6055d881eefaff58ddf
codex/dbg-da-001@c0003d070c6840f0d55487d73b227880cfd27494
codex/dbg-st-006@bf92c9be6cf81a7cb704778dafe88e55fee2e235
codex/dbg-rf-002-003@69a54aeb3394d3cd4792bce620748e15bab69f1f
codex/dbg-rf-004@82154c614a31284723bf3e6a337c5bedfb8aba5d
```

None is merged to `main`. PR #49 targets `Implementation/vscode` and stops at
the earlier `sdp/debugger-gap-analysis` head; the later delivery branches had no
PR in the study. The study also records that `codex/dbg-rf-002-003@69a54a...`
has accepted first-wave publication evidence and that
`codex/dbg-rf-004@82154c...` is the latest observed operational checkpoint.

For the first pilot, Steering should name `codex/dbg-rf-004@82154c...` as the
frozen operational source evidence **if it is still the intended lineage**,
while continuing to name `main@e374da...` as merged truth. At Issue creation,
freshly resolve advertised refs and GitHub state. If the operational branch has
advanced or a newer accepted descendant exists, amend the Issue with the exact
selected commit and refreeze; never silently substitute “latest.”

## Proposed bounded HSX Issue

Title:

```text
Pilot: register HSX work domains and one provisional scoped assignment
```

Primary prospective work owner:

```text
SHARED-REF-001 — adopt provisional work-domain/assignment metadata
```

First prospective Slice:

```text
SHARED-SLC-001 — declare domains, preserve legacy ID inventory, and validate one assignment
```

These names are prospective reservations, not facts about HSX today. Before
activation, scan every advertised operational ref/track for normalized
collisions and reserve them in a common-base coordination record. If either is
already issued, select the next free number without renaming any history.

The work is a Refactor because it changes process/metadata structure while
preserving product behavior and historical meaning. It MUST NOT include product
code, canonical Toolkit installation, default-branch migration, or cleanup of
existing traceability contradictions.

## Branch and PR topology

1. Observe `main` and the intended operational branch separately; record exact
   commits and the evidence caveat in the new Issue.
2. From the exact selected operational commit, create a dedicated pilot branch,
   for example `codex/issue-<N>-vnext-scope-pilot`.
3. Open an early draft PR from that branch **to the named operational
   integration branch**, not to `main`, unless a separate owner Issue has first
   authorized reconciliation to default.
4. Add only isolated pilot records and validation under the operational SDP
   surface: the domain registry, plural per-Issue assignment, prospective
   `SHARED-REF-001`/`SHARED-SLC-001` contracts, and project-local validator
   fixture/config. Mark all records experimental.
5. Report the PR as branch-scoped in-flight evidence. Passing checks, review,
   or Steering acceptance on that PR does not make its contents merged/default
   truth. A later separately authorized convergence Issue decides how the
   operational chain reaches `main`.

This preserves the study's evidence precedence: default is merged truth; the
exact Issue branch/PR is bounded in-flight truth; Issue comments are durable
authorization; merge/reconciliation alone promotes state.

Use exactly one pilot registry object for canonical repository
`https://github.com/Hans-Einar/HSX`. The pilot Issue, amendment comments, and
planned PR use the canonical GitHub URL forms; the Issue and PR must belong to
that repository. Reject zero-padded/trailing-slash/query/fragment aliases
before reservation. Host/owner/repository comparison is case-insensitive, so
a spelling variation cannot create another Issue node or registry partition.

## Domain declarations and history preservation

Declare four ordinary domains, each with a newly minted immutable UUID UID,
stable key, owners, and current root mapping:

| Key | Current evidence root | Meaning in the pilot |
|---|---|---|
| `HSX` | `SDP/HSX/` | HSX core/portable-contract work domain |
| `DBG` | `SDP/Debugger/` | Debugger work domain |
| `AVR` | `SDP/AVR/` | AVR work domain |
| `SHARED` | `SDP/Shared/` | Explicitly owned common process/contract domain |

The Issue must name real owners for `SHARED`; the folder name alone does not
grant common ownership. Root paths are mappings and may differ from keys.

Inventory, preserve, and qualify every existing ID exactly as written,
including `DBG-RF-001..009`, `DBG-DA-001`, `DBG-ST-002..006`,
`HSX-ST-001..008`, `AVR-ST-001`, and all additional IDs discovered at the
selected exact head. Existing `RF` and `ST` spellings remain historical. Do not
convert them to `REF`/`STU`, renumber them, or change path case. New pilot work
uses the prospective v0 forms `SHARED-REF-001` and `SHARED-SLC-001`.

Store each historical spelling as a structured `legacy-preserved` inventory
member with its exact repository, commit, path/source, and authority Issue when
one exists. If no Issue existed, store `authorityIssue: null` plus a nonblank
`authorityMissingReason`; do not invent one. Store each new pilot ID as
`prospective` with the allocating Issue and a normalized portable
repository-relative source, and put the identical `source` in the represented
record. New IDs end exactly in `-NNN`. Check normalized uniqueness before
allocation and reject any reservation by an Issue other than the recorded
authority. The first move exercise must preserve the complete structured
member data, not only the ID strings.

Create absolute references by pairing each old ID with the new UID of its
owning domain. This adds qualification without rewriting the original ID. A
later repository split keeps both values.

## Assignment and concurrency exercise

The pilot assignment must contain:

- exact pilot v0 `schemaVersion`, exact `kind`, and `experimental: true` on
  every prospective record, registry, assignment, reservation, and history
  snapshot; preserve loose historical material only as structured
  `legacy-preserved` inventory that cannot satisfy current gates;
- selected operational integration branch and exact base commit;
- the `SHARED` UID and `SHARED-REF-001`/`SHARED-SLC-001` reservations;
- owned pilot metadata paths, with the four existing domain trees otherwise
  read-only;
- explicit shared touchpoints, if any, and their domain owner;
- Issues #36–#48 and PR #49 as observed/related evidence, not copied current
  status;
- dependency/conflict/merge order with any live HSX Issue Masters discovered
  at pilot start;
- one canonical reservation-set object containing every active Issue row,
  qualified ID, private/shared path, dependency/conflict edge, complete
  dependency-consistent merge order, and explicit terminal convergence; every
  assignment binds its recomputed digest and identical projection;
- collision validation across NFKC/casefold/slash normalization (important
  because Issue #36 already exposed the `agents.md`/`AGENTS.md` Windows case
  collision);
- owned/shared paths disjoint from prohibited paths and stable diagnostics for
  invalid Unicode scalar content before any reservation digest is computed;
- stale-base/refreeze before any change if the selected operational head moves;
  and
- fresh exact-candidate verification and separate review.

Accepted pilot state uses qualified evidence objects: each verification and
current-review reference repeats the exact candidate, the current review is
`approved`, and the Steering acceptance cites a canonical Issue-comment URL
and the same candidate. Opaque strings, whitespace IDs, unresolved
placeholders, or evidence for a different candidate do not qualify.

The existing contradictory current-state evidence—open child Issues after
branch-local completion, differing HSX verified heads, and future-dated Ledger
events—must be reported, not normalized as part of this Slice. They are useful
pilot inputs for declared/observed/accepted projections and future Analyzer
work.

## What stays project-local

Keep HSX-specific debugger/ABI/address/epoch contracts, anti-monolith domain
decomposition, device/platform constraints, test matrices, legacy oracle,
controller/gateway generation semantics, and product acceptance thresholds in
HSX. The generic pilot supplies identity, authority, reservation, relation,
evidence, and move semantics only.

## Pilot success and failure evidence

Success means a fresh Master can discover all four domains, resolve old and new
qualified references, start two synthetic or real disjoint assignments from
one frozen base without ID/path collision, detect stale-base drift, validate
the branch from a fresh checkout, and explain why accepted branch evidence is
not default truth. It must also demonstrate a no-rewrite export/tombstone plan
for one domain without actually moving it.

Failure evidence is equally valuable before schema freeze: ambiguous domain
ownership, inability to inventory legacy IDs, excessive manual duplication,
false collision positives/negatives across Windows/Linux paths, inconvenient
qualified references, unsafe concurrent sharing, or a move that seems to
require identity rewrite. Record these against [OpenQuestions.md](OpenQuestions.md)
and return them to the SDP Steering Group; do not patch around them with an
HSX-only schema.

## Stop condition for the downstream Issue

Stop when the isolated branch/draft PR has a validated domain registry, one
pilot assignment and Slice, preserved legacy identities, fresh review with no
unresolved Blocking/High/Medium findings, and an explicit success/failure
report. Do not merge to `main`, start product work, move a domain, or install a
canonical Toolkit unless separately authorized.
