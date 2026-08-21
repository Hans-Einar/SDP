# gh-sdp SDP usage study

## Evidence frame and limitations

- **OBSERVED:** This report studies `Hans-Einar/gh-sdp` at the exact current
  default-branch identity `main@32613734781bf39f2fce176db2acfb2284dfc92f`
  (tree `2cf00e181a6ce8c13cf3fb3e578fc06ae2b1fc6a`), committed
  `2026-07-15T06:23:24Z` as the normal merge of PR #2. Remote `main`, the local
  `origin/main` object, and GitHub repository metadata were checked on
  2026-08-21. The studied checkout was not changed.
- **OBSERVED:** GitHub state was inspected through the authenticated GitHub API
  and CLI on 2026-08-21. The repository has no GitHub Issues. It has three PRs:
  merged PRs #1 and #2 and open draft PR #3. All PR issue comments, native review
  submissions, review/check rollups, workflow runs, commit statuses, branches,
  releases, and tags were inspected.
- **OBSERVED:** Default-tree truth is kept separate from open/unmerged evidence.
  Draft PR [#3](https://github.com/Hans-Einar/gh-sdp/pull/3) has exact head
  `5eea5567eac8139815f45aee369160acca03b9f9` (tree
  `d99f9b6a41c6f367d3307ae98bc4e18857658647`), but that tree is not on `main`.
  Its Steering-acceptance and reconciliation records are therefore evidence of
  a proposed/latest working practice, not current default-tree state.
- **OBSERVED:** The installed Toolkit source pin is
  `Hans-Einar/SDP@bc110bb5fd60009ba67015cf640ad6ddbfe1b04b` from 2026-07-12.
  The current canonical SDP baseline for Issue #5 is later
  `main@e398ebaf3a4ace6a5d92fd9ce22736a7427a9e15` from 2026-07-17. Both call
  themselves unreleased Toolkit `0.2.0`; the source commit, not the version
  string alone, is needed to distinguish them.
- **OBSERVED:** Upstream SDP draft PR
  [#4](https://github.com/Hans-Einar/SDP/pull/4) is also unmerged evidence. The
  gh-sdp PR #3 reconciliation assesses older PR #4 commit
  `bf20832bed618ab240cf87c17517fc31ea721311` as
  `UPSTREAM_REWORK_REQUIRED`. Current PR #4 head is later
  `d611b8bf72aeb30d86c5ef28902469462b06a803`, for which GitHub reports two
  successful checks and no native reviews or PR comments. The older three-
  Medium assessment must not be silently projected onto that later head.
- **OBSERVED:** No product tests were run for this read-only study. In fact,
  default `main` contains no CLI product source, Go module, executable,
  workflow, or product test. Verification discussed below is the documentary
  and command evidence committed by the project, plus current GitHub state.
- **INFERENCE:** The recorded Codex reviewer identities such as
  `/root/study_reviewer_final` are credible process evidence, but GitHub cannot
  independently prove fresh context: all three PRs have zero native GitHub
  review submissions, and every durable PR comment is authored under GitHub
  identity `Hans-Einar`.

## Repository facts

- **OBSERVED:** Repository: `Hans-Einar/gh-sdp`; public, non-archived, not a
  fork; default branch `main`; created `2026-07-13T09:17:08Z`. Its description
  is “GitHub CLI extension for installing, inspecting, previewing, and updating
  the Standard Document Procedure (SDP).” GitHub reports no detected language,
  no license, no tags, and no releases.
- **OBSERVED:** All 14 commits reachable from default `main` fall inside Issue
  #5's fixed window (`2026-04-20T21:08:52Z` through
  `2026-08-20T21:08:52Z`). The first is
  `2dcc18423a10497678eb19ace8fd9668efda9551` at
  `2026-07-13T09:17:53Z`; the latest/default commit is `326137347...` at
  `2026-07-15T06:23:24Z`. GitHub `pushedAt` is later
  `2026-07-17T13:32:09Z` because open PR #3's branch was pushed after `main`.
- **OBSERVED:** SDP is extensively installed. Root authority is
  [`AGENTS.md`](https://github.com/Hans-Einar/gh-sdp/blob/32613734781bf39f2fce176db2acfb2284dfc92f/AGENTS.md),
  with repository-specific additions in
  [`AGENTS-project.md`](https://github.com/Hans-Einar/gh-sdp/blob/32613734781bf39f2fce176db2acfb2284dfc92f/AGENTS-project.md).
  Role skills are installed under `.codex/skills/sdp-*`. Project governance is
  rooted at `SDP/`.
- **OBSERVED:** The default tree has 55 regular files. Relevant locations are
  `SDP/01--Mandate/`, `02--Study/`, the still-template
  `03--Requirements/` through `07--Implementation/`, `Instructions/`,
  `Sprints/`, `Refactors/README.md`, `Fixes/README.md`, `Traceability/`,
  `CodeReview/`, `Verification/`, `Releases/`, and `Framework/`. There is no
  `SDP/Features/`, no `SDP/Steering/`, no populated Refactor, and no Fix record
  on `main`. PR #3 proposes the first Fix record.
- **OBSERVED:** Installed identity in
  `SDP/Framework/installed-toolkit.manifest.yaml` is Toolkit `0.2.0`, Framework
  `1.0.0`, AGENTS contract `1.0.0`, installer `0.2.0`, installed
  `2026-07-13T09:26:11Z` from exact canonical source `bc110bb5...`, with ten
  version-`1.0.0` SDP skills and five installed capabilities. The Toolkit was
  and remains unreleased; this is a source installation, not an official
  release installation.
- **OBSERVED:** Project-local product identity in
  `SDP/SDP-project.manifest.yaml` is current client `0.0.0`, next target
  `0.1.0`, state `unreleased`, with null latest tag/commit and all development
  coordinates clear. `REL-0.1.0` is only a proposed release record with empty
  release verification/review arrays and null publication identities.
- **OBSERVED:** Canonical/current identity is different from installed identity.
  Current `Hans-Einar/SDP@e398eba...` still declares unreleased Toolkit
  `0.2.0`, but it is three commits newer than the installed pin and adds
  `docs/Feature-Governance-And-SDP-2.0.md` plus a changed Architect skill. That
  document labels itself `Status: design proposal`; it is useful aspirational
  evidence, not installed gh-sdp behavior or a released contract.
- **OBSERVED:** Phase 1 used `Sprint-001` / `SPI-001` / `SPS-001` to bootstrap
  governance and approve Mandate `MAN-001`. Phase 2 used `Sprint-002` /
  `SPI-002` / `SPS-002` to produce the single broad Study `STU-001`. Both
  Sprints and their one Iteration/one Slice are complete on `main`; active
  coordinates are null.
- **OBSERVED:** `CodeReview/` contains five committed review records:
  `REV-SPS-001-001/002` and `REV-SPS-002-001/002/003`. `Verification/` contains
  `VER-SPS-001` and `VER-SPS-002`. Both handoff sets live inside their Sprint
  folders. `Relations.yaml` has 580 lines; `Ledger.ndjson` has exactly one
  event, the initial `release-planned` event.
- **OBSERVED:** PR
  [#1](https://github.com/Hans-Einar/gh-sdp/pull/1), “Phase 1: bootstrap SDP
  governance and define gh-sdp Mandate,” was created
  `2026-07-13T09:49:08Z`, has exact head `38a7ea46...`, two PR comments, no
  native reviews/checks, and merged normally as `3a3b9ece...` at
  `2026-07-13T21:24:07Z`.
- **OBSERVED:** PR
  [#2](https://github.com/Hans-Einar/gh-sdp/pull/2), “Phase 2 Study: gh-sdp
  feasibility and safety groundwork,” was created `2026-07-13T22:09:49Z`, has
  exact head `a9fa6980...`, one later reconciliation comment, no native
  reviews/checks, and merged normally as the studied commit `326137347...` at
  `2026-07-15T06:23:24Z`.
- **OBSERVED:** PR
  [#3](https://github.com/Hans-Einar/gh-sdp/pull/3), “Record Phase 2 Study
  acceptance and reconcile PR #2 merge,” was created
  `2026-07-17T13:12:33Z`; it remains open and draft at exact head `5eea5567...`.
  It has one terminal-review comment, no native reviews/checks, and proposes
  nine changed governance files, including `FIX-0.1.0-001`, one review record,
  and one verification record.
- **OBSERVED:** The repository has four remote branches—`main` and the three
  `codex/phase-*` PR branches—and none is protected. `gh run list` is empty;
  the four relevant PR/default commits have zero status contexts and no
  workflow/check runs.

## How SDP is used

### Work start and authority

- **OBSERVED:** Work did not start from a GitHub Issue; the repository has no
  Issues at all. Phase 1 began from an external assignment against an almost
  empty repository, then wrote the governing Mandate and Sprint/Slice contract.
  Phase 2 began only after two externally stated conditions were met: PR #1
  merged and a “separate Study assignment” was issued. `CurrentIndex.yaml`
  records that authority only as the string
  `separate-study-assignment-2026-07-13`; it has no durable Issue, comment,
  assignment ID, or Steering record behind it.
- **OBSERVED:** Once work began, the repository became the detailed work
  contract. Each Slice states goal, why now, scope/non-goals, expected files,
  invariants, traceability, verification, independent review, completion signal,
  and stop boundary. The Mandate records owner/Steering decisions and explicitly
  keeps later phases unauthorized.
- **OBSERVED:** External Steering decisions were reconciled into existing
  project-owned files because installed SDP had neither a Steering decision
  schema nor a suitable Ledger event. On `main`, Phase 1 acceptance appears in
  `mandate.md`, `CurrentIndex.yaml`, `Relations.yaml`, and Sprint-001 handoff.
  Phase 2 Steering acceptance appears only in unmerged PR #3.
- **INFERENCE:** gh-sdp demonstrates a strong repository contract after
  activation but weak durable assignment provenance before activation. A future
  session can reconstruct what was executed, but cannot inspect the actual
  Phase 1/Phase 2 owner assignment that authorized it.
- **OWNER DIRECTION:** Issue #5 makes the GitHub Issue the future operational
  assignment boundary. gh-sdp's no-Issue history is therefore migration
  evidence, not a counterexample that should preserve conversational authority.
- **RECOMMENDATION:** Future gh-sdp work should begin from one high-quality
  Issue or exact Issue comment, mirrored into a compact assignment record before
  broad work. The repository record should point to, not paraphrase, the full
  authority.

### Codex Master, Worker, Architect, Verifier, and Reviewer

- **OBSERVED:** Root `AGENTS.md` says every non-delegated agent is Master. The
  Master reads repository instructions/manifests/active work, maintains release
  versus development identities, delegates bounded work and fresh review,
  inspects real verification, updates traceability, and stops at the Slice/Fix
  boundary. Worker, Reviewer, Architect, Verifier, Traceability, Release,
  Auditor, Versioning, and Vertical-Refactor skills are installed separately.
- **OBSERVED:** Phase 2 implementation notes explicitly say the Master delegated
  three separate read-only research passes—GitHub CLI/Go, Toolkit installer and
  ownership, and security/platform/transaction—and personally inspected cited
  primary evidence. This is concrete Master-led parallel work, although worker
  identities and individual deliverables are not persisted.
- **OBSERVED:** Review used fresh named Codex contexts. The Phase 1 sequence went
  from 4 Medium findings to a follow-up pass with none; the Study sequence went
  from 2 Medium/2 Low, to 1 Medium, to 0 findings. Exact candidate commits are
  recorded for every pass. Reviewers did not edit the authored content.
- **OBSERVED:** Verification is first-class but role separation varies. The
  tracked verifier records contain exact commands and environment; final fresh
  Reviewers reran the matrix on later exact heads. PR #1 also has a terminal
  comment naming `/root/phase1_acceptance_verifier`. Phase 2 has no separate
  durable Verifier identity beyond the verification record and Reviewer rerun.
- **OBSERVED:** No separate Architect pass occurred. That is appropriate for the
  Study-only boundary: ten `STU-ADI-*` records are explicitly Architecture
  inputs, not decisions, and Architecture remained unauthorized/template-only.
- **INFERENCE:** The role model is conceptually mature but its actor evidence is
  mostly prose. It proves disciplined separation better than it proves who ran
  which context or what exact worker output the Master accepted.
- **RECOMMENDATION:** Preserve Master-per-assignment, bounded Workers, fresh
  Reviewer, and exact-head Verifier evidence. Require Architect only when the
  assignment may alter accepted architecture. Give each delegated pass a small
  durable identity, inputs, output path, and inspected commit rather than
  expanding the Ledger with every chat action.

### Branches and draft PRs

- **OBSERVED:** Every phase used a fresh `codex/` branch against exact `main`
  and a PR against `main`. PR #1 and PR #2 were normal merges, deliberately
  preserving referenced commit SHAs. PR evidence distinguishes PR head, merge
  commit, merge parents, and tree equality.
- **OBSERVED:** The Study activation commit `50bade647...` was made before the
  substantive Study candidate, which is good scope-first practice. However,
  draft PR #2 opened only after the first full Study candidate
  `f1f2aba...`; PR #1 similarly appeared after the first substantive governance
  commits. Thus this repository uses a draft PR as a review surface, but does
  not consistently open it immediately after activation.
- **OBSERVED:** Branch/PR state is repeated manually in implementation notes,
  handoff, verification, review records, and PR descriptions. Review caught
  stale PR/handoff text twice. After PR #2 merged, the merged tree necessarily
  continued to say “open draft” and “ready for Steering assessment,” requiring
  PR #3 for reconciliation.
- **RECOMMENDATION:** Open the draft PR immediately after the assignment/SDP
  activation commit. Store the stable binding once; derive live PR state and
  merge identity from GitHub. Keep reviewed-head and merge-commit distinctions,
  because gh-sdp uses them well.

### Feature, Refactor, Sprint, Iteration, Fix, and Slice

- **OBSERVED:** Features have no actual meaning in the installed/default
  repository because `SDP/Features/` does not exist. `MAN-001` describes the
  future client capability; Sprint-001 owns governance bootstrap and Sprint-002
  owns the entire Study phase.
- **OBSERVED:** Refactor is only a Toolkit README. No Refactor ID, contract, or
  execution exists. The default hierarchy is exactly the installed older model:
  `Sprint or Refactor -> Iteration -> Slice/Fix`.
- **OBSERVED:** Each of the two Sprints has exactly one Iteration and one Slice.
  `SPS-001` is a governance-document bootstrap, and `SPS-002` is a broad Study
  corpus. Neither is a product vertical slice through a horizontal software
  architecture because the repository has no product code.
- **OBSERVED:** Open PR #3 proposes `FIX-0.1.0-001` for a bounded governance
  reconciliation, explicitly arguing that reopening the completed Sprint or
  creating a new lifecycle Slice would misrepresent the work. This is a useful
  correction use case, but it also couples a process repair to the proposed
  product target `0.1.0`.
- **INFERENCE:** Sprint/Iteration/Slice are functioning here mostly as mandatory
  ceremony around lifecycle phases. A one-Iteration/one-Slice Sprint adds little
  decomposition value, while the broad Study is much larger than the future
  concept of one small end-to-end delivery Slice.
- **OWNER DIRECTION:** Issue #5 requires new implementation work to belong to a
  Feature or Refactor and retains Slice as a vertical capability. It explicitly
  rejects Sprint as capability owner merely for historical compatibility.
- **RECOMMENDATION:** Treat `gh-sdp`'s approved Mandate as the horizontal product
  foundation. Future capabilities such as `status`, `preview`, and verified
  install/update should be explicit Features or Refactors with optional
  execution timeboxes. A Slice should be a small testable vertical outcome, not
  a synonym for an entire Study or governance phase. Keep Fix for genuine
  bounded corrections, including governance reconciliation when clearly
  justified, but do not require a SemVer-derived Fix ID for every process repair.

### Initial lifecycle versus later requirements and design

- **OBSERVED:** `MAN-001` provides a strong product boundary, stakeholders,
  outcomes, command intents, safety objectives, and seven strategic questions.
  `STU-001` supplies evidence and candidate requirements/architecture inputs.
  The repository deliberately left Requirements, Architecture, Design Analysis,
  Design, and Implementation templates unchanged because the phase assignment
  did not authorize them.
- **OBSERVED:** That restraint is excellent scope discipline, but the documented
  lifecycle is sequential: Study must be Steering-accepted and separately
  authorize Requirements/Architecture. No installed mechanism permits a later
  Feature to add a local Study, requirements, or design after the initial
  project design.
- **OBSERVED:** Current canonical `SDP@e398eba...` contains a non-normative
  Feature-governance proposal with Feature, Integration Study, Release, Sprint,
  Slice, Steering, and CurrentAssignment concepts. It explicitly retains the
  old Sprint-first 1.x implementation hierarchy. Issue #5's later owner
  direction goes further by making Feature/Refactor the implementation owner and
  GitHub Issue the assignment boundary.
- **RECOMMENDATION:** Preserve gh-sdp's initial Mandate/Study rigor as the
  horizontal skeleton, then allow Feature-local Study/requirement/design
  refinement. Require an explicit architecture revision only when a Feature
  changes stable boundaries, ownership, protocols, persistence, or cross-cutting
  constraints; otherwise keep refinement local to the Feature/Slice.

### Study practice and convergence

- **OBSERVED:** gh-sdp used one large `STU-001`: 18 sections, 21 evidence IDs,
  14 findings, 8 recommendations, 17 requirement candidates, 10 Architecture
  inputs, and 7 cross-repository findings. The Master delegated independent
  research topics, but all results converged directly into one Study document
  and one Sprint/Slice.
- **OBSERVED:** The Study is exceptionally explicit about evidence
  classification: verified external fact, verified repository fact,
  observation, inference, recommendation, assumption, and unresolved question.
  Sources include exact upstream commits, official source/docs, actual commands,
  what each proves, and limitations.
- **OBSERVED:** There is no separate multi-Study manifest, no independently
  accepted `STU-*` children, and no formal convergence gate. The final review is
  a quality/closure gate, not a cross-Study design convergence decision.
- **INFERENCE:** The single document worked for a tightly related product-
  feasibility question, but the 77 child IDs and 580-line relations file show
  the cost of flattening every evidence atom into one manual graph.
- **RECOMMENDATION:** Keep one broad Study when the question is coherent. Permit
  multiple independent Studies and a convergence gate only when uncertainties
  can be investigated independently or have different safety/evidence owners.
  Do not require stable IDs for every paragraph-level evidence item unless a
  downstream decision or audit genuinely references it.

### Independent review and verification

- **OBSERVED:** Review is adversarial and iterative. Medium findings blocked
  closure, remediation occurred in new commits, and a fresh context reviewed
  the corrected exact head. Review found real evidence defects: stale candidate
  statements, incomplete hidden-file checking, missing relation edges, a dead
  URL, inaccurate source interpretation, and missing durable review evidence.
- **OBSERVED:** Verification records include OS/tool versions, exact source and
  candidate commits, commands, schema checks, ID/path resolution, full-tree and
  forbidden-artifact checks, diff hygiene, and limitations. The project
  truthfully says Toolkit `validate_sdp.py` validates the Toolkit repository, not
  a consuming project, and it composes additional project-local checks instead
  of claiming a nonexistent canonical validator.
- **OBSERVED:** The final tracked record cannot self-name its containing commit.
  The project handles this through an external exact-head PR comment. PR #1's
  two comments preserve terminal Phase 1 evidence; PR #2's later comment and PR
  #3's terminal comment preserve reconciliation evidence.
- **CONTRADICTORY EVIDENCE:** There are no native GitHub reviews, no review
  threads, no workflow runs, and no status checks. Therefore “independent review
  passed” is an SDP/Codex convention supported by committed reports and owner-
  authored comments, not a GitHub-enforced branch gate.
- **RECOMMENDATION:** Carry forward exact-head verification, changed-head
  re-review, severity gates, and candid limitations. Add CI once executable
  validators/product tests exist. Preserve terminal evidence as a signed/
  attributable acceptance object or generated post-merge snapshot rather than
  relying solely on a same-account PR comment.

### Handoff

- **OBSERVED:** Each Sprint has a focused `Handoff.md` naming authoritative
  entry points, exact branch/base/candidates, review/verification state,
  residual limitations, next decision, and stop boundary. This is a useful
  fresh-session recovery surface.
- **OBSERVED:** Handoff is manually duplicated state. Default
  `Sprint-002/Handoff.md` still says PR #2 is open/draft and Steering assessment
  is next even though the file was merged by PR #2. PR #3 exists principally to
  reconcile that transition.
- **RECOMMENDATION:** Retain a concise prose handoff for context, residual risk,
  and next decision. Generate coordinates, Issue/PR status, exact heads/checks,
  merge state, and assignment disposition from canonical machine records plus
  live GitHub evidence.

### CurrentIndex, Relations, and Ledger

- **OBSERVED:** `CurrentIndex.yaml` is concise and useful: it separates project,
  lifecycle, release, and active development coordinates. On default `main` it
  says Study is a reviewed candidate pending Steering and all execution
  coordinates are clear.
- **OBSERVED:** `Relations.yaml` is a comprehensive explicit graph: Mandate
  outcomes/boundaries/success criteria/assumptions/questions, every Study ID,
  Sprints/Iterations/Slices, review, verification, and proposed release. Exact
  candidate commits are first-class for review/verification.
- **OBSERVED:** The graph does not contain branch or GitHub Issue/PR binding.
  Those facts remain prose in handoff/reviews/verifications/PR bodies.
- **OBSERVED:** The 580-line Relations file repeats the path
  `SDP/02--Study/study.md` for almost every one of 77 Study IDs. This is
  mechanically checkable but expensive to maintain and was itself a repeated
  source of review findings.
- **OBSERVED:** The Ledger contains only one `release-planned` event because the
  installed schema accepts release lifecycle events only. Sprint, Slice, Study,
  review, verification, Steering decisions, PR merge, and reconciliation are
  absent. The project correctly refused to invent unsupported event types.
- **INFERENCE:** Traceability is simultaneously over-detailed in Relations and
  under-expressive in the Ledger. CurrentIndex/Relations also become stale at
  GitHub transitions because they manually mirror live state.
- **RECOMMENDATION:** Retain stable relations for Feature/Refactor, important
  Study/requirement/architecture decisions, Slice, review, verification, PR/
  commit, and release. Generate child-ID/path indexes and live GitHub fields.
  Extend the Ledger only with decision/recovery/audit-significant events; do not
  turn every file edit or agent message into an event.

### Steering and CurrentAssignment

- **OBSERVED:** Default `main` has no `SDP/Steering/` and no
  `CurrentAssignment`. Steering decisions are text embedded in Mandate, current
  state, Relations, and handoff. The separate Study authorization has no durable
  external coordinate.
- **OBSERVED:** Open PR #3 proposes richer Steering disposition fields in
  `CurrentIndex` and records `ACCEPTED_WITH_LOW_FINDINGS` dated 2026-07-17. It
  still adds no `Steering/CurrentAssignment`; its bounded work contract is a Fix
  record. The terminal result is attached to the PR conversation.
- **OBSERVED:** The current canonical SDP design proposal at `e398eba...`
  suggests `Steering/CurrentAssignment.yaml`, Master Reports, Decisions, and
  Blockers, but explicitly calls these conceptual contracts whose schemas need
  separate work. This is aspirational, not installed gh-sdp capability.
- **OWNER DIRECTION:** Issue #5 establishes ChatGPT Steering Group above the
  bounded Codex Master but keeps repository/GitHub evidence authoritative.
- **RECOMMENDATION:** Add a compact machine-readable assignment that binds the
  exact Issue/comment authority, Feature/Refactor, baseline, branch/draft PR,
  active Slice(s), permissions, invariants, stop boundary, and current
  disposition. Separate immutable assignment terms from generated live GitHub
  status so a merge cannot make the contract itself false.

### GitHub binding

- **OBSERVED:** Current binding is PR-centric, not Issue-centric. PR URLs and
  candidate commits appear in review/verification/handoff; branch and base
  appear in prose; the Mandate and Study have no GitHub Issue. Relations has no
  GitHub section. There is no stable assignment ID linking all these identities.
- **OBSERVED:** PR #3 illustrates identity repair: it verifies PR #2 head,
  merge commit, parents, tree equality, and corrected body, then records the
  evidence in a Fix/review/verification and PR comments. This is strong audit
  discipline but a costly manual reconciliation.
- **RECOMMENDATION:** Canonical binding should model
  `Issue/comment <-> assignment <-> Feature/Refactor <-> Slice <-> branch/PR
  <-> reviewed head/check/review <-> merge <-> release`. Store stable IDs and
  repository-owned decisions; query GitHub for current mutable state and flag
  contradictions.

### Automated versus conversational convention

- **OBSERVED — automated outside this repository:** The canonical PowerShell
  installer was previewed/applied and its upstream tests were run during Phase
  1. Custom Python/PowerShell matrices parsed schemas and Git objects. GitHub
  CLI/API commands verified repository/PR identities.
- **OBSERVED — not automated here:** There is no `.github/`, product script,
  reusable consuming-project validator, or `gh-sdp` executable on default
  `main`. GitHub has no workflow run or check context. All validation logic is
  embedded in long Markdown evidence records or run from the separate canonical
  SDP checkout.
- **OBSERVED — conversational/manual:** owner authorization, Steering
  disposition, branch/PR creation, role delegation, evidence integration,
  CurrentIndex/Relations/Handoff updates, review identity, PR comments, and
  merge/reconciliation are manual conventions.
- **INFERENCE:** gh-sdp currently specifies the automation it wants more than it
  supplies automation. It is a governance and feasibility repository for a
  future CLI, not yet a CLI implementation.

## CLI, installer, and consumer implications

- **OBSERVED:** `MAN-001` defines future command intents `install`, `status`,
  `preview`, `update`, `version`, and `help`, but no executable semantics are
  approved or implemented. The project clearly separates `gh extension upgrade
  sdp` (client update) from future `gh sdp update` (Toolkit content update).
- **OBSERVED:** `STU-001` already identifies the most important installation
  boundary: canonical SDP must own a versioned machine-readable plan, ownership
  and payload contract, while gh-sdp validates and applies it portably. It also
  identifies the false consumer seed risk (`REL-0.2.0.yaml`), missing
  deletion/tombstones, incomplete identity/capability projection, lack of
  preview/apply binding, archive/link hazards, and recovery/concurrency gaps.
- **OBSERVED:** The installed manifest demonstrates why `toolkitVersion` alone
  is insufficient. `bc110...` and `e398...` both identify as unreleased 0.2.0
  but differ in managed/proposed method content. Status and migration decisions
  must include exact installed source/release/plan/payload identity.
- **RECOMMENDATION — future status contract:** `gh sdp status` should present
  separately: client version/release; installed Toolkit version plus exact
  release/source, plan and payload digest; latest/selected eligible Toolkit;
  project schema/capabilities; declared active Feature/Refactor/assignment/
  Slice; linked Issue/branch/PR; observed PR/check/review/merge state; release
  inclusion; and explicit stale/contradictory/offline fields. It should never
  overwrite declared repository state merely because GitHub differs.
- **RECOMMENDATION — future Issue/Feature contract:** One implementation Issue
  should normally authorize one bounded Feature or Refactor execution episode.
  A long-lived Feature may have several Issues over time, and one Issue may
  contain several sequential Slices. Study, reconciliation, release, or
  migration Issues need not fabricate a Feature. The binding needs an explicit
  work type rather than assuming 1:1 from numbers.
- **RECOMMENDATION — future assignment contract:** Store `assignment_id`, Issue
  and authority-comment URLs, `work_type`, Feature/Refactor/Fix identity,
  authoritative base, branch and draft PR, required records, permissions/non-
  goals, verification/review gate, stop condition, and current disposition.
  Generate mutable PR/check/merge facts, or store them as timestamped snapshots,
  rather than hand-copying them into every document.
- **RECOMMENDATION — migration contract:** Installation/update must remain
  additive for project-owned records. Add missing `Features/` and `Steering/`
  templates/schemas without inventing Features for historical Sprint-only work.
  Preserve old Sprint/Study IDs and mark their generation. Migrate only with a
  previewed deterministic plan, explicit collision policy, conformance fixtures,
  backup/journal/recovery, and a machine-readable result. Never seed a Toolkit
  release record or live Toolkit traceability into a consumer.
- **RECOMMENDATION — Feature/Refactor migration:** Existing Sprints remain
  historical execution evidence. New work after migration must point to a new
  Feature/Refactor/Fix or an explicit non-implementation governance assignment.
  Do not retroactively relabel `Sprint-001`/`Sprint-002` as Features merely to
  satisfy a new schema.
- **RECOMMENDATION — schema/versioning:** Installed facts should distinguish
  Toolkit release version from contract schema, payload version/digest,
  installation-plan schema, installer/client compatibility, and source
  provenance. A repeated version number at different unreleased commits must be
  reported as different source builds, not “same installation.”
- **RECOMMENDATION — validation:** gh-sdp should consume canonical shared
  fixtures and schemas, validate project-local Feature/assignment/traceability
  semantics, and compare declared state with GitHub. It should not embed a
  second independent interpretation of SDP ownership or migration rules.
- **RECOMMENDATION — project-local safety:** Archive, path, link, locking,
  journaling, recovery, credentials, Windows/Linux/macOS behavior, and exact
  client/Toolkit compatibility remain real product requirements from STU-001.
  They should remain gh-sdp-specific engineering requirements, not be confused
  with the generic future SDP workflow contract.

## What worked well

- **OBSERVED:** Version truth is unusually disciplined. Client, installed
  Toolkit, candidate Toolkit, release target, Sprint/Slice, and publication
  identities are never collapsed into one number.
- **OBSERVED:** The project safely bootstrapped from an unreleased canonical
  source, removed Toolkit-specific consumer seed facts, preserved project-owned
  files, and documented the repeat-initializer false release-record proposal
  without applying it.
- **OBSERVED:** Mandate and Study scope boundaries are exceptionally clear.
  Requirements, Architecture, product code, packaging, release, and upstream
  modifications were repeatedly blocked rather than opportunistically begun.
- **OBSERVED:** Progressive commits separate activation, first candidate,
  remediation, review evidence, and closure. Exact SHA and tree identities make
  later audit/reconciliation possible.
- **OBSERVED:** Fresh iterative review materially improved evidence quality.
  Findings were specific, remediated, and re-reviewed rather than waived after
  a green self-check.
- **OBSERVED:** `STU-001` is a strong technical Study: primary sources are pinned,
  limitations are stated, facts are separated from inference/recommendation,
  alternatives are compared, and downstream owner boundaries are explicit.
- **OBSERVED:** Handoff and exact verification evidence support recovery after
  session loss. A fresh agent can identify what is complete, what was reviewed,
  what remains unauthorized, and which exact objects matter.
- **OBSERVED:** The project did not fabricate generic Ledger events or claim a
  consumer validator/release gate that the installed schemas do not support.
  Restraint preserved traceability truth.
- **OBSERVED:** PR #3's reconciliation work shows careful handling of reviewed
  head versus merge commit, tree equality, post-merge Steering state, and
  unavailable historical evidence—even though it remains unmerged.

## Pain points and accidental complexity

- **OBSERVED:** The repository has 55 governance files and no product code.
  `STU-001` is 676 lines, Relations is 580 lines, and review/verification records
  are extensive. The cost is defensible for safety research but disproportionate
  as a default model for ordinary Features.
- **OBSERVED:** The same live state is hand-maintained across Mandate/Study,
  CurrentIndex, Relations, manifest, Sprint notes, Handoff, review,
  verification, PR body, and PR comments. Review repeatedly found stale copies.
- **OBSERVED:** Merging PR #2 merged a tree that said the PR was still open/draft
  and awaiting Steering. The method then required a second governance-only PR
  to reconcile the inevitable transition. That PR is still unmerged, so
  default-tree state remains stale six weeks later.
- **OBSERVED:** There is no GitHub Issue authority and no durable Steering
  assignment. Strings such as “separate assignment” and role names cannot be
  audited back to their source.
- **OBSERVED:** Relations over-specifies paragraph-level Study IDs while the
  Ledger cannot represent any actual Study, Slice, review, Steering, or merge
  transition. The two traceability surfaces have opposite granularity problems.
- **OBSERVED:** Local identifier grammar was invented in
  `Instructions/StableIdentifiers.md` because the installed Toolkit did not
  define Mandate/Study IDs. Future canonical migration must preserve these IDs
  without pretending they were standard.
- **OBSERVED:** One-Sprint/one-Iteration/one-Slice wrappers make lifecycle phases
  look like vertical delivery while providing no actual product decomposition.
- **OBSERVED:** The same unreleased Toolkit version `0.2.0` names materially
  different source states. A consumer that compares only version strings may
  falsely consider itself current.
- **OBSERVED:** Review independence is not enforceable from GitHub. There are no
  native reviews/checks and only one GitHub user identity on comments/merges.
- **OBSERVED:** PR #3's assessment of upstream PR #4 is already commit-specific
  historical evidence; PR #4 advanced substantially afterward. Manual
  dependency status becomes stale unless the assessed commit is always shown.
- **OBSERVED:** The project has no reusable validator or actual `gh sdp status`
  command to detect the very contradictions its Study identifies.

## Carry forward

- **RECOMMENDATION:** Repository-authoritative work contracts with explicit
  scope, non-goals, invariants, verification, review, completion, and stop
  boundaries.
- **RECOMMENDATION:** Exact baseline/branch/PR/head/tree/merge identities and
  progressive candidate/remediation/review commits.
- **RECOMMENDATION:** Client, Toolkit, project, work, and release identities as
  separate concepts, including exact source/provenance for unreleased installs.
- **RECOMMENDATION:** Master coordination with bounded research/implementation
  Workers, fresh independent review, exact-head verification, and changed-head
  re-review for unresolved Blocking/High/Medium findings.
- **RECOMMENDATION:** Primary-evidence Study discipline with explicit fact,
  observation, inference, recommendation, assumption, and limitation labels.
- **RECOMMENDATION:** Concise recovery handoff and truthful refusal to invent
  unsupported schema events or publication claims.
- **RECOMMENDATION:** Canonical SDP ownership of the machine-readable install/
  migration plan and conformance fixtures; gh-sdp as a portable validating
  consumer, not a second rule owner.
- **RECOMMENDATION:** Feature/Refactor-owned future delivery, with optional
  Sprint/Iteration grouping and small vertical Slices, plus bounded Fixes for
  genuine corrections.
- **RECOMMENDATION:** Issue/comment authority and a compact assignment binding,
  with mutable GitHub state generated or captured as timestamped evidence.

## Legacy / do not carry forward

- **RECOMMENDATION:** Do not carry forward conversational “separate assignment”
  strings as sufficient authority. Replace them prospectively with exact GitHub
  Issue/comment plus assignment identity.
- **RECOMMENDATION:** Do not preserve Sprint as the conceptual owner of a
  capability or use one Iteration/one Slice merely to wrap every phase.
- **RECOMMENDATION:** Do not treat broad lifecycle Study closure as evidence of
  a vertical product Slice; use distinct Study and delivery semantics.
- **RECOMMENDATION:** Do not mirror live PR/check/merge state manually across
  many files. Generate it, detect contradiction, and reconcile additively.
- **RECOMMENDATION:** Do not carry forward the current paragraph-level ID and
  relation explosion as a universal minimum. Stable IDs should serve authority,
  decisions, delivery, evidence, release, or recovery.
- **RECOMMENDATION:** Do not use a release-only Ledger schema as a pretext for an
  empty operational history, and do not solve that gap by logging every trivial
  action. Add a small set of high-value assignment/decision/acceptance/recovery
  events.
- **RECOMMENDATION:** Do not infer installed equivalence from Toolkit version
  alone, especially for unreleased source installations. Exact release/source,
  plan, payload, and schema identities are required.
- **RECOMMENDATION:** Do not treat open PR #3, current canonical design proposal
  `SDP-PROP-001`, or upstream PR #4 as accepted/default-tree contract. They are
  useful commit-bound evidence until merged and released.
- **RECOMMENDATION:** Do not generalize gh-sdp's archive, provenance,
  cross-platform filesystem, authentication, or transaction rules into every
  SDP project. They are product-specific requirements; only installer ownership,
  migration truth, and consumer-status implications belong in generic SDP.
- **RECOMMENDATION:** Do not claim `gh-sdp` exists as a functioning CLI yet.
  Default `main` is a reviewed governance/Study foundation for future product
  work, with client version `0.0.0`, no product source, no CI, no tag, and no
  release.
