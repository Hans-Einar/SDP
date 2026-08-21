# Provisional workflow contract

Status: **pilot v0; not canonical**

## 1. Operating model

The project foundation is a living horizontal skeleton:

```text
Mandate -> Study or Studies -> Requirements -> Architecture
        -> Design Analysis -> Initial Design
```

Ongoing work grows through bounded intent and vertical outcomes:

```text
Feature | Refactor | genuine Fix
    -> one or more GitHub Issue assignments over time
    -> zero or more bounded Studies and one or more vertical Slices
    -> exact-candidate verification + independent review
    -> Steering disposition
    -> optional Release inclusion/publication
```

An Issue can be Study-only and therefore authorize no implementation Slice.
Sprint and Iteration are optional coordination metadata; neither owns product
intent nor becomes a mandatory path segment.

Repository and GitHub evidence are authoritative. Conversation memory and agent
summaries are not evidence.

## 2. Entities

### Feature

A Feature owns a durable product or externally meaningful operational
capability. It states the problem, outcome, acceptance criteria, scope,
non-goals, constraints, foundation links, local requirement/design refinements,
and residual limitations. It is not a branch, Issue, Sprint, or Release.

Suggested declared lifecycle:

```text
proposed -> studying -> ready -> active -> delivered -> released
                 \-> blocked | rejected | superseded
```

`delivered` requires accepted implementation evidence. `released` additionally
requires an explicit inclusion in an actual Release; merge alone is neither.

### Refactor

A Refactor owns a bounded structural outcome whose primary purpose is improved
ownership, dependency direction, maintainability, safety, or architecture. It
MUST name the behavior baseline, behavior/compatibility to preserve or change,
target structure, migration range, temporary adapters, and exit evidence. It
uses the same Issue, Slice, verification, review, and Steering gates as a
Feature.

### Fix

A Fix owns a proportionate correction to accepted behavior. It MUST NOT conceal
new capability, a public/architectural contract change, or a broad migration.
A tiny low-risk Fix MAY itself be the smallest verified/reviewed unit and omit
a Slice only when `risk: low`, `standaloneReviewedUnit: true`, exactly one
Issue assignment/authority names that Fix, and the Fix itself carries the same
qualified accepted evidence otherwise required for a Slice. A medium-, high-,
or safety-critical-risk Fix, or any non-standalone Fix, uses one or more
Slices. Review rework on an unaccepted candidate is not a Fix; a correction
after acceptance is. Every delivered/released Fix, sliced or standalone, has at
least one qualified, locally resolved `affectedWork` reference to an
evidence-qualified accepted behavior surface. A durable work owner may remain
`active` after an accepted assignment, but owner-level aggregate evidence alone
is insufficient: the target MUST retain at least one accepted assignment and
accepted authorized Slice at its immutable candidate (or, for a proportional
standalone Fix target, one accepted direct zero-Slice outcome). A proposed,
unresolved, evidence-empty, or active owner with no accepted represented
outcome does not prove that the Fix corrects accepted behavior.

### Study

A Study owns a material question and bounded evidence/decision boundary. It may
be foundation-wide or local to later Feature, Refactor, Fix, architecture, or
integration work. Parallel Studies MUST declare independent questions and a
convergence/refreeze gate. Study results inform work; only an accepted decision
or relation makes them implementation authority.

### Slice

A Slice is the smallest independently implementable, runnable, verifiable, and
reviewable end-to-end outcome through the horizontal boundaries it needs.

A Slice contract MUST name:

- one primary Feature, Refactor, or Fix owner and one authorizing assignment;
- observable outcome and why it is the smallest coherent vertical increment;
- linked requirements/decisions and permitted local refinement;
- expected areas, ownership, shared touchpoints, invariants, and non-goals;
- executable verification and independent-review criteria;
- discovery/escalation rule, completion signal, and hard stop.

Activation, baseline capture, a pure Study, publication, verification-only
work, and governance closeout are bounded tasks or gates, not automatically
Slices. A domain-foundation Slice is an exceptional pilot case: it MUST have a
frozen interface, local runnable evidence, and a named later integration Slice.

### Sprint and Iteration

A Sprint MAY group one or more assignments/Slices into a near-term objective or
timebox. It owns no capability and is not required to locate records.

An Iteration MAY describe a real learning/replanning cycle containing one or
more Slices. It has exactly one primary work owner, MAY belong to one Sprint,
and MUST record the learning or decision that justifies it. A 1:1 wrapper added
only for hierarchy is invalid pilot practice.

### GitHub Issue assignment and Issue Master

One GitHub Issue is one bounded operational assignment. The complete Issue body
and owner/Steering comments are the human-readable authority. One Codex session
acts as Master for that Issue, reads repository evidence first, creates or uses
the declared branch and early draft PR, delegates implementation and fresh
review, inspects exact evidence, records results, and stops at the Issue
boundary. It does not become a permanent project Master.

The assignment record is a compact recovery and coordination binding. It links
to the Issue rather than copying its full prose or mutable GitHub facts.

### Steering disposition

Steering supervises and makes product/architecture/acceptance decisions. It may
accept, require changes, block, reject, cancel, split, or supersede an
assignment. ChatGPT may assist Steering, but repository/GitHub records—not the
chat—are durable evidence. An acceptance disposition MUST identify the exact
candidate and qualifying verification/review.

### Release inclusion

A Release explicitly includes accepted Feature, Refactor, or Fix outcomes.
Work identity never encodes a target release. A delivered work owner may be
included in zero or more Releases; inclusion/publication does not rewrite its
identity or earlier evidence.

## 3. Explicit cardinalities

| Relation | Pilot cardinality and rule |
|---|---|
| Work owner to Issue | A proposed work owner has `0..N` Issues; delivery requires `1..N` over its lifetime. |
| Issue to primary work | Exactly `1`. It is Feature, Refactor, Fix, or Study. Related work is `0..N` and cannot silently share primary ownership. |
| Issue to authorized Slice | `0..N` preserved over the Issue lifetime. Study-only assignments may have zero; implementation assignments normally have `1..N`; a tiny standalone Fix may have zero. |
| Issue to active Slice | Pilot v0 has `0..1`, always a subset of `authorizedSlices`. It is the currently executing Slice only, never the historical Slice set. An accepted assignment has zero active Slices. |
| Slice to owner | Exactly `1` Feature, Refactor, or Fix, plus `0..N` related/dependency references. |
| Slice to Issue assignment | Exactly `1`. Give it a separate Issue when authority, risk, ownership, branch/PR, or independent acceptance is separate. |
| Study to owner | `0..1` local owner plus `0..N` informed work/decision relations. A foundation Study may have no local owner. |
| Sprint membership | Work/Issue/Slice has `0..1` active Sprint membership; a Sprint groups `1..N` assignments or Slices. |
| Iteration membership | Slice has `0..1`; Iteration has exactly one primary owner, `1..N` Slices, and `0..1` Sprint. |
| Verification/review | Each accepted Slice or standalone Fix has `1..N` qualifying verification results and exactly one current independent review disposition for the exact candidate. Earlier attempts remain history. |
| Steering disposition | Assignment has `0..N` dispositions; only the latest disposition for the exact current candidate governs. Completion requires one terminal accepted disposition when the Issue contract calls for Steering acceptance. |
| Release inclusion | Accepted Feature/Refactor/Fix outcome has `0..N` Release relations; a Release has `1..N` accepted inclusions. |

An integration assignment that crosses several domains still identifies one
primary work owner and makes other participants/dependencies explicit. If no
honest primary exists, Steering creates a dedicated integration work owner
rather than weakening the cardinality.

For every represented Slice, validation walks both directions: the owner lists
the Slice exactly once; the Slice names that owner and one Issue; exactly one
assignment for that Issue preserves the Slice in `authorizedSlices`; the
assignment's primary `workRef` is that owner; and the owner's
`issueAuthorities` contains that Issue exactly once. Each authorized entry
stores `acceptedCandidate: null` until acceptance and thereafter stores the
immutable qualified Slice candidate. `activeSlices` is a `0..1` subset of that
set. For an active implementation assignment it is also an exact reverse
projection: represented authorized Slices for that Issue/owner whose
`declaredState` is `active` equal `activeSlices`. Accepted earlier Slices remain
authorized without appearing in the active projection. A missing
owner-referenced or authorized Slice is an error rather than an
unresolved future placeholder. Slice `decisionRefs` are qualified and locally
resolvable, private paths are portable and contained by the assignment's owned
reservations, and shared paths exactly match its declared touchpoints.

Prospective inventory records a stable ID's one-time immutable
`allocationIssue`; it does not claim that every future assignment executes
under that Issue. The allocating Issue must match each newly created/reserved
Slice, Study, and Fix identity. A later assignment may reference an already
issued Feature, Refactor, or Fix only when its canonical Issue occurs exactly
once in that owner's `issueAuthorities`; this is execution authority, not ID
reclaim. Legacy-preserved inventory separately retains historical
`authorityIssue` (or a truthful null plus reason) as provenance.

Embedded and top-level semantic edges are two authored surfaces of one global
set. Feature and Refactor `relations` contain one supported relation `type`
plus one qualified `targetRef`; Study `ownerRef`, `independentStudyRefs`, and
`informs` normalize respectively to `owned_by`, `independent_of`, and
`informs`; Fix `affectedWork` normalizes to `corrects`. Top-level edges contain
exactly `type`, `sourceRef`, and `targetRef`. Edge objects and qualified
references are strict: candidate or extension fields are not supported in
pilot v0, even though outer record objects remain extension-tolerant. A
normalized `(type, sourceRef, targetRef)` may be authored once across all
surfaces. Bare, unresolved, wrong-kind, unsupported, self-referential, and
cross-surface duplicate edges are errors.

The complete pilot v0 relation-type vocabulary is `depends_on`, `informs`,
`refines`, `supersedes`, `preserves`, `requires_revision`, `owned_by`,
`independent_of`, and `corrects`. Embedded and top-level forms use this same
set; neither surface accepts a candidate-only relation type.

Both surfaces enforce this one endpoint-kind matrix before duplicate or DAG
processing:

| Relation | Allowed source kinds | Allowed target kinds |
|---|---|---|
| `corrects` | Fix | Feature, Refactor, Fix |
| `owned_by` | Study | Feature, Refactor, Fix |
| `independent_of` | Study | Study |
| `informs` | Study | Study, Feature, Refactor, Fix |
| `depends_on`, `refines`, `supersedes`, `requires_revision`, `preserves` | Study, Feature, Refactor, Fix | Study, Feature, Refactor, Fix |

`depends_on` is a gating edge. When its source is represented active/accepted/
delivered/released implementation work, its target is evidence-qualified
`accepted`, `delivered`, or `released`; a proposed/active Study or work owner
cannot authorize implementation. Pilot v0 Slice `decisionRefs` are deliberately
narrower: each is a qualified, locally resolved, evidence-qualified accepted
Study. Other decision-record kinds may be considered after pilot evidence.

Four directional relation types form separate globally normalized DAGs across
embedded and top-level surfaces. Direction is always authored source to target:

- `depends_on`: the source depends on the target;
- `supersedes`: the source replaces the target;
- `refines`: the source elaborates the target; and
- `requires_revision`: the source requires revision of the target.

Two- or longer-hop cycles are invalid across scopes and repositories.
`preserves`, `informs`, `owned_by`, `independent_of`, and `corrects` are not
DAG-gated; they retain the same self-edge, duplicate, reference, and kind rules
described above.

## 4. Identity and revision rules

The record types use `FEAT`, `REF`, `FIX`, `STU`, and `SLC` prospectively.
Existing forms such as `RF`, `ST`, Sprint-bound Slice IDs, or version-bound Fix
IDs remain historical evidence and MUST NOT be rewritten.

- A record ID is stable, never reused, and never encodes branch, Issue, Sprint,
  calendar, or release state.
- Before acceptance, corrections to the same outcome increment `revision` and
  invalidate evidence tied to the previous candidate.
- Material new behavior, risk, or acceptance scope creates a new Slice.
- Rejected or superseded records remain addressable; their IDs are not recycled.

Every stateful record has a normalized, non-recursive, portable repository-
relative `source` path. For prospective work this value is byte-for-byte equal
to the source in its issued-ID inventory member. A repository validator also
resolves that path to the represented JSON record when the record exists as a
local file. Within one canonical repository, NFKC/casefold/slash-normalized
prospective sources are one-to-one: two identities cannot claim one source and
each active prospective inventory member binds exactly one represented record.
This prevents an inventory entry from naming one artifact while the record
claims another.

### Pilot v0 value and state floor

The eventual canonical vocabulary remains open, but pilot v0 supports only:

| Kind | Supported declared states |
|---|---|
| Feature / Refactor | `proposed`, `studying`, `ready`, `active`, `blocked`, `rejected`, `superseded`, `delivered`, `released` |
| Fix | `proposed`, `ready`, `active`, `blocked`, `rejected`, `superseded`, `delivered`, `released` |
| Study | `proposed`, `active`, `blocked`, `rejected`, `superseded`, `accepted` |
| Slice | `proposed`, `active`, `blocked`, `rejected`, `superseded`, `accepted` |
| Issue assignment | `proposed`, `active`, `blocked`, `accepted`, `rejected`, `cancelled`, `superseded` |

`revision` is an integer of at least 1. Required titles, intent/question/
outcome/correction strings, convergence/stop strings, and evidence identifiers
are nonblank. Required collections retain their advertised array/object types;
acceptance, verification, review, domain ownership, Slice invariants, and
similar contract-bearing arrays cannot be empty. Domain owners are unique
nonblank strings. Branches use the conservative portable Git subset documented
in [IssueContract.md](IssueContract.md). `independentReview` is Boolean and
`maximumUnresolvedSeverity` is one of `none`, `low`, `medium`, `high`, or
`blocking`. These are a minimum usable pilot contract, not a canonical schema
freeze.

Every prospective current pilot v0 record, assignment, reservation set,
active/moved registry, and assignment-history snapshot carries its exact
documented `schemaVersion`, exact `kind`, and `experimental: true`. Missing as
well as unknown markers are errors; removing markers never downgrades a current
object into a loose mode. Compatibility handling is limited to structured
`issuedIds` members explicitly marked `status: legacy-preserved` with their
provenance, and prior Git bytes reconstructed through a fully typed current v0
history snapshot. Loose historical bytes remain evidence only; neither they
nor legacy inventory can by itself satisfy a current reservation, decision,
implementation, evidence, dependency, or acceptance gate. Historical
relations may remain inspectable, but they do not become current authority.

### Reopen and extension

- An Issue reopened before its assignment is accepted may increment the
  assignment `revision`, refreeze its baseline/reservations, and continue. Past
  candidate evidence remains immutable.
- Reopening an accepted Issue does not reopen accepted Slices in place. Prefer
  a new Issue; if the same Issue is intentionally reused, create a new
  assignment revision and new Slice IDs. Preserve every earlier Slice in
  `authorizedSlices` with its immutable accepted candidate; set only the new
  Slice active. Accepted revision snapshots and the bound reservation digest
  make an attempted earlier-evidence rewrite fail closed.
- Across every assignment revision, canonical Issue, assignment `source`, and
  primary qualified `workRef` are immutable. `authorizedSlices` is append-only:
  every earlier identity remains, every non-null accepted candidate remains
  byte-for-byte unchanged, and a later revision may only qualify a formerly
  null candidate or append a new Slice. `activeSlices` and the refrozen
  reservation/base may change under the Issue authority.
- A delivered Feature or Refactor is extended by a new Issue and new Slice(s).
  The work owner may move back to declared `active`, but the earlier delivery
  and release evidence remains intact.
- A defect in accepted behavior becomes a Fix. It does not revise the accepted
  Slice as if the defect never existed.

## 5. Declared, observed, and accepted state

These projections MUST be kept separate:

- **declared** — project-owned intent and lifecycle in work, Slice, domain, and
  assignment records;
- **observed** — timestamped Git/GitHub facts such as branch head, Issue/PR
  state, checks, reviews, merge, tag, or GitHub Release;
- **accepted** — evidence-qualified state supported by an exact candidate,
  applicable verification, current independent review, and required Steering
  disposition.

Every Feature, Refactor, Fix, Study, Slice, and Issue-assignment shape carries
the same explicit projection:

```json
{
  "acceptedEvidence": {
    "candidate": null,
    "verificationRefs": [],
    "currentReviewRef": null,
    "steeringDisposition": null,
    "releaseRefs": []
  }
}
```

`proposed` and `active` records may retain this null/empty projection.
The non-null pilot v0 form is deliberately qualified rather than a list of
opaque strings:

```json
{
  "acceptedEvidence": {
    "candidate": "1111111111111111111111111111111111111111",
    "verificationRefs": [
      {"id": "VER-101", "candidate": "1111111111111111111111111111111111111111"}
    ],
    "currentReviewRef": {
      "id": "REV-101",
      "candidate": "1111111111111111111111111111111111111111",
      "disposition": "approved"
    },
    "steeringDisposition": {
      "authority": "https://github.com/example/app/issues/101#issuecomment-501",
      "candidate": "1111111111111111111111111111111111111111",
      "decision": "accepted"
    },
    "releaseRefs": [
      {"id": "REL-001", "candidate": "1111111111111111111111111111111111111111"}
    ]
  }
}
```

Evidence IDs are stable nonblank portable tokens, never whitespace or an
`unresolved:` placeholder. Every nested candidate equals the exact top-level
candidate. The current review disposition is `approved`; the Steering
authority is the canonical comment URL on an applicable authority Issue and
its decision is `accepted`. `accepted` or `delivered` requires the candidate,
one or more verification objects, current review, and Steering disposition.
`released` additionally requires one or more qualified Release objects. A
delivered Feature/Refactor also requires at least one Issue authority and one
resolvable Slice; the qualified low-risk standalone Fix remains the only
proportional zero-Slice delivery exception.

Acceptance is also cross-record state, not only a well-shaped evidence object.
An accepted implementation assignment has `activeSlices: []` and requires
every `authorizedSlices` member to be `accepted` with qualified evidence. Each
authorized entry repeats that Slice's own immutable candidate; different
Slices and the final aggregate assignment may have different candidates. Its
durable Feature/Refactor may remain `active` for later assignments. An
accepted Study-only assignment requires its primary Study to be `accepted` at
that candidate. An accepted standalone Fix assignment requires its primary Fix
to be delivered/released at the same candidate and requires exactly one
accepted, evidence-qualified, zero-authorized-Slice assignment. An in-set
dependency must already be accepted, and both endpoints of a symmetric
conflict cannot simultaneously be `active` or `accepted`; an active endpoint
with its peer explicitly `blocked` is the valid reservation state.

A Feature, Refactor, or non-standalone Fix declared `delivered` or `released`
resolves every `issueAuthorities` member to exactly one preserved assignment
for that owner. Every implementation assignment is evidence-qualified
`accepted`; all of its authorized Slices and every Slice in the owner's
complete `slices` list are evidence-qualified `accepted` at their own preserved
candidates. The aggregate work-owner evidence candidate does not have to equal
every historical assignment/Slice candidate: delivery aggregates accepted
outcomes over time. A delivered/released standalone Fix instead closes exactly
  one accepted zero-Slice assignment at the Fix candidate. Every terminal Fix
  also closes its qualified affected-work targets. An active affected owner
  retains at least one accepted assignment plus accepted authorized Slice or
  direct standalone outcome; its aggregate evidence object cannot substitute
  for that represented outcome. None of these terminal rules
changes the valid case where one accepted assignment and Slice leave a durable
multi-assignment Feature `active`.

An active Feature/Refactor/non-standalone-Fix assignment has exactly one active
Slice in pilot v0, and that pointer set exactly equals all represented
declared-active Slices for the Issue and primary owner. Both assignment and Slice expose at least one owned or
shared write surface, and the assignment cannot list its own hosting repository
under `prohibitedRepositories`. Study-only and explicitly standalone
zero-Slice non-implementation units may be pathless when the Issue documents
that boundary and the compact assignment mirrors a nonblank `pathlessReason`.

Authored records may contain stable URLs, the declared integration base, target
branch, policy, and intended PR URL. They MUST NOT claim mutable Issue/PR/check/
head/merge state as authored truth. Generated observations are timestamped and
replaceable. A contradiction is reported, never silently reconciled.

Default branch is merged repository truth. A named Issue branch/PR head is
bounded in-flight truth. Steering comments authorize a change but do not make
unmerged files default truth. Only a merge/reconciliation decision promotes
the branch state.
