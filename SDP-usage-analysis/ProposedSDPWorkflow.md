# Proposed next canonical SDP workflow

Status: recommendation for Steering Group review; not implemented

## Decision summary

- **OWNER DIRECTION:** Product growth belongs to Feature or Refactor; a small
  Fix remains available for genuine correction. GitHub Issues are the
  operational assignment boundary, a bounded Codex session acts as Master for
  one Issue, and ChatGPT Steering Group supervises without replacing repository
  authority.
- **RECOMMENDATION:** Keep the initial Mandate/Study/Requirements/Architecture/
  Design progression as a living horizontal project foundation, not as a frozen
  implementation plan for the whole future product.
- **RECOMMENDATION:** Make `Feature | Refactor | Fix` the semantic work owners.
  Execute them through one or more small vertical Slices. Sprint becomes an
  optional grouping/timebox; Iteration is not a mandatory hierarchy level.
- **RECOMMENDATION:** Store authored decisions once and derive Git/GitHub-native
  facts. Maintain three distinguishable status projections: declared repository
  state, observed Git/GitHub state and accepted evidence state.

## Canonical hierarchy

```text
Project
├── Foundation (living horizontal system skeleton)
│   ├── Mandate
│   ├── Study or bounded independent Studies
│   ├── Requirements
│   ├── Architecture
│   ├── Design Analysis
│   └── Initial Design
├── Features/FEAT-NNN
│   ├── optional Feature-local Study / requirement / design deltas
│   ├── one or more GitHub Issue assignments over time
│   └── vertical SLC-NNN work
├── Refactors/REF-NNN
│   ├── optional local Study / architecture / migration contract
│   ├── one or more GitHub Issue assignments
│   └── behavior-preserving or explicitly compatible SLC-NNN work
├── Fixes/FIX-NNN
│   └── bounded correction and evidence
└── Releases/REL-X.Y.Z
    └── includes accepted Feature / Refactor / Fix outcomes
```

```text
Human owner / ChatGPT Steering Group
    -> GitHub Issue (one bounded assignment)
        -> Codex Issue Master
            -> optional Architect / Study Workers
            -> bounded implementation Worker(s)
            -> Verifier
            -> fresh independent Reviewer
            -> exact-head Master report / stop
        -> Steering disposition
            -> accept / changes required / blocked / cancel
```

Sprint and Iteration may annotate this graph when they add real coordination;
they do not sit between the work owner and its Slices by default.

## Stable identities

| Entity | Recommended identity | Rule |
|---|---|---|
| Feature | `FEAT-NNN` | Repository-local, stable, never reused; release is a relation, not encoded in the ID. |
| Refactor | `REF-NNN` | Peer of Feature; stable across several assignments/Slices. |
| Fix | `FIX-NNN` | Repository-local correction identity; do not encode a target release that may change. Preserve legacy version-bound Fix IDs as history. |
| Study | `STU-NNN` | Create only for a material evidence/decision boundary; local IDs may relate to a Feature/Refactor. |
| Slice | `SLC-NNN` | Repository-local, stable and never reused/reopened for material new behavior. |
| Release | `REL-X.Y.Z` | Public/internal release identity under existing SemVer rules. |
| GitHub assignment | full repository + Issue number | The Issue is already globally addressable; an `ASGN-<issue>` alias may be generated for tooling but should not become another hand-maintained authority. |

- **RECOMMENDATION:** Do not encode Sprint, Issue, release or calendar state in
  Feature/Refactor/Slice identity. Relations carry those changing associations.
- **RECOMMENDATION:** A bounded correction to an unaccepted Slice candidate may
  use a numeric `revision`. Materially new behavior after acceptance gets a new
  Slice, even when the same Feature and Issue continue.

## Feature contract

A Feature is a durable product capability or externally meaningful operational
capability. It is not an implementation task, branch or timebox.

Minimum record:

- stable ID, name and lifecycle state;
- problem, intended outcome, users/value;
- scope, non-goals, constraints and product invariants;
- related Use Cases and foundation Requirements/Architecture/Design;
- Feature-local Study/requirement/design deltas and explicit supersession;
- related GitHub Issues, Slices, verification/reviews and Releases;
- completion evidence and residual limitations.

Lifecycle recommendation:

```text
proposed -> studying -> ready -> active -> delivered -> released
                    \-> blocked / rejected / superseded
```

`delivered` and `released` are distinct. A Feature is delivered when all
required Slices are accepted at exact candidates and its acceptance criteria are
satisfied. It is released only through an actual Release inclusion/publication
relation.

- **RECOMMENDATION:** One Issue normally authorizes one bounded execution
  episode for one primary Feature. A Feature may have multiple Issues and
  sequential Slices over time. An Issue may relate to several Features when a
  vertical integration genuinely crosses them, but it must identify one primary
  work owner or an explicit integration assignment.
- **RECOMMENDATION:** Extend a delivered Feature with a new Issue and new Slice;
  preserve the earlier delivered evidence. Use `supersedes` or `refines`
  relations rather than rewriting prior acceptance.

## Refactor contract

A Refactor is a bounded structural objective whose primary intent is to improve
ownership, dependency direction, maintainability, safety or architecture while
preserving declared behavior or explicitly controlling compatibility change.

Minimum record adds to the common work contract:

- verified behavioral baseline;
- target ownership/dependency structure;
- compatibility and migration range;
- preserved/changed behavior list;
- architecture/design deltas;
- ordered vertical migration Slices;
- adapter/removal conditions and exit evidence.

- **RECOMMENDATION:** Refactor uses the same Issue/Master/branch/draft-PR/
  verification/review/Steering mechanics as Feature.
- **RECOMMENDATION:** When structural preparation has no owner-visible
  capability, do not force it into Feature (for example TerrainAnalyzer
  FEAT-007). Use Refactor, or a preparation assignment related to the future
  Feature when no product implementation occurs.

## Fix contract

A Fix is a proportionate correction to accepted behavior where no new
capability, architecture/public-contract change or broad migration is required.

- **RECOMMENDATION:** Use `FIX-NNN` with goal, defect evidence, scope,
  invariants, affected release/work owner, verification, review and stop.
- **RECOMMENDATION:** A tiny low-risk Fix may use one Issue/PR without a Slice
  ID when the Fix itself is the smallest reviewed work contract. A risky or
  cross-layer Fix may use a Slice.
- **RECOMMENDATION:** Review rework on an unaccepted candidate remains rework;
  it is not a new Fix. A production correction after acceptance is a Fix.

## Sprint and Iteration decision

### Sprint

- **RECOMMENDATION:** Sprint is optional execution grouping or timebox metadata
  inside one or more Feature/Refactor assignments. It is not the persistent
  owner of product intent and is not required for every Issue.
- Create a Sprint only when it groups several Slices toward one coherent
  near-term objective, coordinates several agents/repositories, or represents a
  real planning/timebox decision.
- A Sprint ID, when used, must not be required to locate Feature, Slice or
  evidence records. Small repositories may omit Sprints entirely.

### Iteration

- **RECOMMENDATION:** Remove Iteration from the mandatory canonical hierarchy.
  Use it optionally for a real replanning/learning cycle spanning several
  Slices or a convergence checkpoint.
- Do not create an Iteration that is always 1:1 with a Sprint or Slice. Review
  rework uses candidate revisions; material new scope uses a new Slice.

Existing Sprint/Iteration IDs remain valid historical evidence and are not
retroactively rewritten.

## Vertical Slice definition

A Slice is the smallest independently implementable, runnable, verifiable and
reviewable end-to-end outcome through the horizontal boundaries it needs.

### Entry criteria

- one Feature/Refactor (or related Fix/integration) owns the work;
- Issue authority, exact baseline, branch/PR policy and stop boundary exist;
- required Study/requirements/architecture/design decisions are accepted;
- expected areas, invariants, non-goals and ownership are explicit;
- verification and independent-review contracts are executable/proportionate;
- dependencies and concurrent ownership are clear.

### Required Slice contract

- goal and why now;
- owner work ID and Issue;
- observable outcome;
- linked requirement/design decisions and allowed local refinement;
- expected modules/areas, owners and shared touchpoints;
- invariants, non-goals and prohibited opportunistic work;
- verification matrix and fresh-review gate;
- discoveries/escalation rule;
- completion signal and hard stop.

### Completion evidence

- immutable candidate commit/tree and PR identity;
- deterministic checks with environment/artifact/run identities;
- required real-world/manual evidence or explicit pending limitation;
- independent review at the current candidate with zero unresolved Blocking,
  High or Medium finding;
- accepted traceability and Handoff/next-decision context;
- Steering disposition when the Issue requires it.

- **RECOMMENDATION:** Activation, baseline capture, Study, publication,
  governance closeout and release preparation are bounded tasks/gates, not
  automatically vertical Slices.
- **RECOMMENDATION:** An Issue may contain several sequential Slices. Give a
  Slice its own Issue when it has independent authorization/risk/ownership,
  needs a separate branch/PR or can be accepted/rejected independently.

## Initial project design and later evolution

### Initial horizontal foundation

The initial pass establishes enough stable structure to prevent accidental
architecture:

- system/context boundaries and responsibilities;
- domain and state ownership;
- dependency direction and major layers;
- public/internal interfaces, schemas and error contracts;
- persistence, integration and security boundaries;
- migration/compatibility policy;
- cross-cutting quality/verification strategy.

It does not attempt to design every future Feature or decompose the lifetime of
the product into Sprints.

### Feature/Refactor-local refinement

- Local Study and requirements/design deltas belong with the work owner.
- An accepted delta links back with `refines`, `preserves`, `supersedes` or
  `requires_revision` rather than silently rewriting history.
- Reconcile the global foundation only when the current system-wide contract
  actually changes.

### Architecture-revision threshold

Require an explicit architecture revision and Architect/fresh design review when
work changes any of:

- component/domain/state ownership;
- dependency direction or trust/security boundary;
- public/cross-component contract or compatibility range;
- persistence/integration topology or migration policy;
- cross-cutting performance, reliability or deployment invariant.

An interface/local state/implementation choice within accepted ownership can
remain Feature/Slice-local design.

## GitHub Issue assignment contract

A high-quality SDP Issue contains:

1. work type and primary `FEAT`/`REF`/`FIX`/Study identity;
2. goal, why now and observable completion;
3. authoritative default/integration baseline and prerequisite state;
4. scope, expected areas and non-goals;
5. requirements/architecture/design references and allowed refinement;
6. invariants, safety/security/physical boundaries and forbidden actions;
7. role/delegation expectations and concurrent ownership;
8. branch and early draft-PR policy;
9. verification, real-world evidence and independent-review requirements;
10. Steering/merge/closure authority;
11. stop/escalation/split conditions;
12. milestone and final reporting contract.

The Issue holds the complete human-readable assignment. Repository state mirrors
only fields needed for local machine recovery and validation.

## Codex Master-per-Issue

- Read repository instructions, current declared/observed/accepted state and the
  complete Issue/comments before changing work state.
- Reconcile stale repository/GitHub facts and record the exact baseline.
- Create the branch, activation/binding record and early draft PR.
- Refine only the bounded work contract; return architecture/product choices to
  Steering when not authorized.
- Delegate product code to bounded Workers; use Architect/Study Workers when the
  horizontal contract may change; use a Verifier and separate fresh Reviewer.
- Inspect actual exact-candidate evidence, not summaries.
- Resolve every Blocking/High/Medium finding; request re-review after head
  changes.
- Update authored records and append truthful material transitions.
- Report, stop at the Issue boundary and never auto-start the next Issue.

The Issue Master differs from the older permanent project Master: it owns one
assignment episode, while the repository and Steering Group preserve continuity
across Issue Masters.

## Steering Group model

Steering is supervisory decision authority, not repository truth and not an
implementation role.

Responsibilities:

- product/architecture direction and Feature/Refactor disposition;
- author/refine bounded Issues and resolve scope/architecture escalations;
- inspect actual PR, CI, review, verification and repository evidence;
- accept, require changes, block, cancel or split assignments;
- coordinate cross-repository dependencies and release inclusion;
- create follow-up Issues and stop at decision boundaries.

The durable evidence of Steering action is the Issue/comment plus repository
decision/assignment relation, not ChatGPT conversation memory. Human owner
authority and ChatGPT assistance must not be conflated with a GitHub actor name.

## CurrentAssignment recommendation

Replace the singleton experiment with one compact record per Issue:

```text
Steering/
└── Assignments/
    ├── ISSUE-123.yaml
    └── ISSUE-124.yaml
```

`Steering/CurrentAssignments.yaml`, when useful, is a generated local/read-model
view of all active assignments at an `observedAt` instant. It is not an authored
authority file and is not committed by default.

Recommended authored fields:

```yaml
schemaVersion: sdp-assignment-v1
authority:
  issue: https://github.com/OWNER/REPO/issues/123
  comment: null
work:
  type: feature
  id: FEAT-017
baseline:
  branch: main
  commit: <40-char-sha>
delivery:
  branch: codex/issue-123-short-name
  pullRequest: https://github.com/OWNER/REPO/pull/456
  activeSlice: SLC-042
coordination:
  integrationBase: <40-char-sha>
  dependsOnIssues: []
  conflictsWithIssues: []
  ownedPaths: []
  sharedTouchpoints: []
  mergeOrder: null
permissions: []
invariants: []
requiredEvidence:
  verification: []
  review: []
stopCondition: awaiting-steering-review
```

- The file is created on the Issue branch before broad work. Its scope is that
  Issue/branch/PR; concurrent Issue Masters use different filenames and do not
  overwrite a singleton pointer.
- Default branch contains only assignment records that reached it through an
  accepted merge/reconciliation. An in-flight branch record remains declared
  assignment evidence, not default-tree current truth.
- Omit PR state, head, checks, merge status and Issue open/closed state; derive
  them live or capture them in a timestamped generated observation.
- Store full scope/non-goals in the Issue/Slice, link rather than copy. Issue
  comment amendments are appended/referenced, not silently folded into the
  original contract.
- On merge, preserve the per-Issue record as historical declared terms; generated
  current views stop listing it only when observed GitHub state and accepted
  evidence satisfy the terminal rule. A disagreement is a finding, not an
  overwrite.
- Before parallel assignments start, Steering establishes a common integration
  base, disjoint ownership/reserved IDs, shared touchpoints, dependency/conflict
  edges and merge/convergence order. Overlapping assignments require a
  preparation/refreeze gate.
- Clearing or advancing one assignment never changes another. Tools validate
  path/ID collisions, stale bases and merge-order violations across the active
  set.

## Declared, observed and accepted state

- `declared`: project-owned intent/status from Feature/Refactor/Slice and its
  per-Issue assignment record.
- `observed`: Git/GitHub facts such as Issue state, branch head, PR state,
  checks, review objects, merge commit, tag and Release at `observedAt`.
- `accepted`: evidence-qualified state from exact verification, independent
  review and Steering disposition.

Tools display all three and raise contradictions; they do not silently overwrite
declared decisions with GitHub or treat repository prose as proof of GitHub
state.

## Traceability simplification

### Keep stable authored graph nodes

- Feature, Refactor, Fix;
- material Study, Requirement and Architecture/Design decision;
- Slice;
- verification and review disposition/finding set;
- Release.

Use GitHub Issue/PR/commit/tag/Release as external stable keys, not duplicated
local IDs unless offline/tooling requirements justify an alias.

### CurrentIndex

Keep only actual current pointers, target Release and compatibility/profile
facts. Do not use it as a complete historical registry.

### Relations

Store non-derivable semantic edges such as `authorizes`, `implements`,
`refines`, `supersedes`, `depends_on`, `verified_by`, `reviewed_by` and
`included_in`. Generate containment, reverse links, paths and live GitHub edges.

### Ledger

Use a versioned general event envelope and append material transitions only:

- assignment/Study/Feature/Refactor/Fix authorized or canceled;
- Slice started, blocked, reopened, accepted;
- architecture/design decision/refreeze;
- review changes required/resolved;
- Steering disposition;
- correction/supersession;
- merge/release/publication events.

Do not hand-log every Worker message, file edit, branch push or check state.
Generate those observations from Git/GitHub. Reject unsupported, backdated or
future-dated events under deterministic timestamp/source rules.

## Verification, review, handoff and release

- Verification records identify exact subject/candidate, environment, commands,
  artifacts/runs, result and limitations.
- Review records identify independent reviewer attestation, exact candidate,
  findings by severity and disposition. Changed head requires re-review.
- Handoff is concise timestamped narrative: decisions, residual risk,
  unavailable evidence and next authorized decision. Generated status is linked,
  not copied.
- Release inclusion is an explicit relation from accepted Feature/Refactor/Fix
  outcomes to `REL-X.Y.Z`. Delivery, merge and release remain distinct states.
- Publication retains the existing truthful two-phase tag/GitHub Release model.

## Proportional profiles

- **Minimal:** Issue + Feature/Refactor/Fix record + Slice/Fix contract + draft
  PR + verification/review + compact assignment/handoff.
- **Standard:** adds foundation links, CurrentIndex/Relations, material Ledger
  and CI validation.
- **Safety/complex:** adds independent Studies, convergence/refreeze,
  claim/evidence levels, separate Architect/Verifier and physical/manual gates.

Projects select a declared profile/capability set. A legacy or small repository
is not invalid merely because it does not adopt every optional record.

## Completion gate for an Issue assignment

An Issue Master may report ready only when:

- scope and work-owner acceptance criteria are satisfied;
- exact branch/PR/head/merge state is observed and non-contradictory;
- required verification passes at the candidate;
- independent current-head review has no unresolved Blocking/High/Medium item;
- required Steering disposition exists;
- declared/observed/accepted state and durable relations agree;
- residual Low/Note findings and limitations are explicit;
- the Master stops without starting follow-up implementation.
