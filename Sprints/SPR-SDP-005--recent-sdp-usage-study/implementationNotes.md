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
