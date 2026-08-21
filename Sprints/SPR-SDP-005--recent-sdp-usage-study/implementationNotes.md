# SPR-SDP-005 Implementation Notes

## 2026-08-20 — Study initialization

- Read Issue #5 and confirmed it has no comments at study start.
- Fast-forwarded local `main` to authoritative `origin/main` at
  `e398ebaf3a4ace6a5d92fd9ce22736a7427a9e15`.
- Confirmed pull request #4 is a separate open draft based on
  `codex/sdp-install-contract-v1`; it is not this study's base.
- Created branch `codex/issue-5-sdp-usage-study`.
- Fixed the activity window to the four calendar months ending at Issue #5's
  creation timestamp: `2026-04-20T21:08:52Z` through
  `2026-08-20T21:08:52Z`, inclusive.
- Activated `ITR-SDP-005-001` / `SLC-SDP-005-001` for the exact inventory.

No repository analysis conclusion has been accepted yet.

## 2026-08-20 — Inventory finalized

- Enumerated all 35 repositories owned by `Hans-Einar` and visible to the
  authenticated account: 27 private, 8 public, 3 owned forks and no archived
  repositories.
- Used fresh filtered mirrors of all advertised refs and the inclusive UTC
  committer-time window `2026-04-20T21:08:52Z` through
  `2026-08-20T21:08:52Z`.
- Classified 17 repositories in scope and 18 excluded. `HSX` is the additional
  in-scope repository absent from Issue #5's preliminary discovery list.
- Master validation confirmed 35 inventory rows, 17 unique report paths and all
  17 recorded default-branch heads against GitHub.
- A separate fresh-mirror Master reproduction completed with 35 repositories,
  17 in scope, 18 excluded, zero clone failures and no set differences. Its
  verified temporary mirror directory was removed after the check.
- Corrected the method labels so the exact timestamp/committer/all-ref rule is
  recorded as study methodology rather than owner direction.
- Completed `ITR-SDP-005-001` / `SLC-SDP-005-001` and activated
  `ITR-SDP-005-002` / `SLC-SDP-005-002` for the 17 reports.

## 2026-08-20 — Ledger contract correction

- GitHub Actions run `32423995193`, contracts job `96601936007`, rejected the
  study's four initially appended Sprint/Slice events because the current root
  `Ledger.ndjson` schema permits release events only and requires `releaseId`.
- Removed the invalid branch-local lines in a corrective commit. Git history
  retains the failed attempt; the current Ledger remains valid and release-only.
- Did not relabel study transitions as release events and did not change the
  Toolkit event schema because both would violate evidence truth or Issue #5's
  no-implementation boundary.
- Current study state remains durable in the Sprint records,
  `CurrentIndex.yaml`, `Relations.yaml`, draft PR and Issue milestone comments.
  General work-event support remains a migration/design finding for synthesis.

## 2026-08-20 — Repository corpus completed

- Created and Master-reviewed one report for each of the 17 repositories in the
  inventory manifest. Reports distinguish default-tree truth from non-default,
  open-PR and proposal evidence.
- Used dedicated deep Workers for ActionCam, TerrainAnalyzer, SDP, gh-sdp,
  SDP-Analyzer, weight_app_flutter, Lyndata and HSX, plus bounded batches for
  the remaining repositories. The Master read every report, corrected evidence
  classification/identity/coverage issues and committed progressively.
- Added the standard-library corpus validator
  `SDP-usage-analysis/validate_analysis.py`.
- Full validation passed at `dc270f788a0df9d3269080ba05e9783ee0d55e40`:
  35 considered, 17 in scope, 18 excluded, 17 reports, digest
  `a38e7b74ad9c6c58fe93763b3aa82e28f972b903388e22bee90d9f99451b9916`.
- The validator also proved representative negative cases fail: duplicate
  headings, wrong study commit, broken local links, tabs and trailing
  whitespace.
- Completed `ITR-SDP-005-002` / `SLC-SDP-005-002` and activated
  `ITR-SDP-005-003` / `SLC-SDP-005-003` for Master-owned synthesis.

## 2026-08-20 — Cross-repository synthesis completed

- Added the 17-row observational `EvidenceMatrix.md` and evidence/counterexample
  synthesis in `CrossRepositoryPatterns.md`.
- The Master authored `ProposedSDPWorkflow.md`, making Feature/Refactor/Fix the
  semantic work owners, GitHub Issue the execution assignment, Slice the
  vertical outcome and Release the publication owner. Sprint is optional and
  Iteration is no longer mandatory.
- Added `LegacyAndDeprecation.md`, `MigrationImpact.md`, `SkillsAndRoles.md` and
  dependency-ordered `FollowUpIssues.md`.
- The recommendation distinguishes owner direction from observed success: only
  two default trees exercise first-class Feature delivery, two exercise
  first-class Refactor delivery and three show accepted Issue-authoritative
  execution.
- Extended `validate_analysis.py` to cover every required synthesis document,
  matrix identity, required workflow/migration/skill/deprecation topics and
  follow-up dependency ordering.
- Repeated validation passed at
  `ae4b7bc143b569d9ea01b445d47df5a98476ff38`: 35 considered, 17 in scope,
  18 excluded, 17 reports, digest
  `2d470c36f861d3f5bbb86d2fe123b960613588c8c5b504c5fff8bea6aff20bd0`.
- Completed `ITR-SDP-005-003` / `SLC-SDP-005-003` and activated
  `ITR-SDP-005-004` / `SLC-SDP-005-004` for fresh verification and independent
  adversarial review.

## 2026-08-20 — Exact-candidate verification and review rework

- Fresh `VER-SDP-005-001` passed exact candidate
  `721ae53a0a845d0a1cf5286ef8f4d873d926fd29`, including complete live inventory
  reproduction, corpus/toolkit validation, 9 Python tests, PowerShell installer
  fixtures, exact PR checks and negative mutation attempts.
- Fresh `REV-SDP-005-001` returned changes required at the same candidate:
  Blocking 0, High 0, Medium 2, Low 2.
- Resolved M-001 by replacing singleton `CurrentAssignment.yaml` with per-Issue
  assignment records, a generated current-assignment set and explicit
  default/branch/concurrency/ownership/conflict/merge/reconciliation semantics.
- Resolved M-002 by adding an early non-destructive workflow pilot before schema
  freeze and downstream write-capable tooling, while retaining a later migration
  apply pilot.
- Resolved L-001 by relabelling separate skills and the exact Analyzer graph as
  recommendations rather than owner direction.
- Resolved L-002 by correcting the mutable-state/contract drift aggregate from
  eleven to thirteen reports and including HEOS and canonical SDP.
- Reworked corpus validation passes with 35 considered, 17 in scope, 18
  excluded, 17 reports and digest
  `6b5e4ef00fa598c115ae641ff28ad07f731fb5e16a76e9752998e774e1975abc`.
- The changed candidate requires fresh verification and independent re-review;
  the earlier dispositions remain exact-candidate historical evidence.

## 2026-08-20 — Reworked candidate approved

- Fresh `VER-SDP-005-002` passed exact reworked candidate
  `2b51b924ff538e73986eb8a9480343f7b8f89a04`, including deterministic digest
  `6b5e4ef00fa598c115ae641ff28ad07f731fb5e16a76e9752998e774e1975abc`, all
  35 live default heads and fixed-window counts, tests, installer fixtures,
  exact Actions and negative mutations.
- Fresh `REV-SDP-005-002` approved the same candidate: Blocking 0, High 0,
  Medium 0, Low 1, Note 4. It confirmed all prior findings resolved.
- Resolved remaining Low L-003 by removing changing Iteration/Slice and active
  phase claims from the analysis README; it now points to Sprint/CurrentIndex
  and review/verification records for live state.
- Post-capture advertised ref/total-commit counts grew in HSX, SDP and
  TerrainAnalyzer during verification, but all 35 default heads, every frozen
  fixed-window count and the 35/17/18 classification still matched.
- The README remediation and committed evidence records create a new candidate;
  final closure verification/review remains required before completion.

## 2026-08-20 — Study closure accepted

- Fresh closure `VER-SDP-005-003` passed exact candidate
  `834283bb38b5e00e4f24f2b5a7e5b8454d0c4aba`: deterministic corpus digest
  `e346fc230b26c3045305c1ffee22ff43edfc61aa99597005134a314b210b8114`,
  Toolkit validation, 9 unit tests plus pytest subtests, installer fixtures,
  structured formats/links/diff/relations/Ledger, exact Actions run
  `32430079507` and live 35-repository/default-head continuity all passed.
- Fresh closure `REV-SDP-005-003` approved the same candidate with Blocking 0,
  High 0, Medium 0, Low 0 and Note 4. All earlier findings M-001, M-002, L-001,
  L-002 and L-003 are resolved.
- Remaining Notes are evidence limitations, not unresolved findings: mutable
  refs/GitHub state, unrerun private product/physical claims, single-account
  actor provenance and final Master publication/closure integration.
- Completed `ITR-SDP-005-004` / `SLC-SDP-005-004`; cleared active development
  coordinates while keeping Toolkit Release `REL-0.2.0` unreleased.
- The release-only Ledger remains unchanged because its current schema cannot
  truthfully represent study/review transitions without a canonical method
  implementation change.
- Final disposition: `STUDY_READY_FOR_STEERING_REVIEW`. Issue #5 remains open;
  no merge, migration, Toolkit implementation, tag or GitHub Release is
  authorized.
