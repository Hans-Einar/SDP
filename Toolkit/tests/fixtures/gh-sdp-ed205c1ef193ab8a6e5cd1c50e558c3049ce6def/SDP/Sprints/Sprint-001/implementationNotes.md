# Sprint-001 Implementation Notes

## SPS-001

Status: complete — Phase 1 ready for Steering Group assessment

### Worker Changes

- Replaced the installer seed in `AGENTS-project.md` with bounded `gh-sdp`
  instructions that preserve the managed role hierarchy and prohibit production
  implementation in this Slice.
- Added the project-owned `SDP/Instructions/StableIdentifiers.md` convention
  because the canonical Toolkit defines no Mandate/outcome/boundary ID grammar.
- Replaced the Mandate template with drafting `MAN-001`, including product
  problem, stakeholders, future command intent, version/repository boundaries,
  preservation, stable-release, cross-platform, security/reliability, success,
  assumptions, constraints, open questions, and the Phase 1 stop boundary.
- Expanded `SDP/Traceability/Relations.yaml` so the Mandate, its major outcomes,
  boundaries and success criteria, active work coordinates, working-tree
  verification, and proposed release target resolve to real records. The release
  remains unreleased, publication fields remain null, and no new ledger event was
  appended.
- Added `SDP/Verification/VER-SPS-001.md` and updated the handoff for the initial
  Master integration and fresh Reviewer pass.
- Remediated the four medium findings in `REV-SPS-001-001`: exact-candidate
  evidence now names commit `381b38ed6fc59936e7e2cb30b877e156ff57914c`
  and draft PR #1, `git diff --check` is recorded, the scope check includes
  hidden tracked/untracked paths, and all four assumption plus seven question
  IDs resolve through `Relations.yaml` and validation.

### Verification

Deterministic checks and their exact commands/results are recorded in
`SDP/Verification/VER-SPS-001.md`. The reviewed commit
`381b38ed6fc59936e7e2cb30b877e156ff57914c` was rerun directly from Git objects,
including `git diff --check` and all 46 tracked paths, so hidden `.codex` content
and the `.github` exclusion were evaluated without relying on the remediation
working tree. Committed remediation candidate
`bb16bee815d40742f29e15eeac4254bd1310376e` was then rerun directly from Git
objects. It passes the expanded relation check for outcomes, boundaries, success
criteria, assumptions, and questions, along with canonical schemas, paths,
managed hashes, untouched later-lifecycle templates, the one-Git-root invariant,
the complete 47-file tree, full-PR `git diff --check`, and the no-product boundary.

Repeating the one-time `-InitializeProjectStructure` preview proposed one change:
re-creating the Toolkit-specific `SDP/Releases/REL-0.2.0.yaml`. It was not applied.
The project warning and agent instruction record this canonical initializer gap.

### Decisions

- `MAN-001` subordinate IDs use an explicitly project-local Phase 1 convention;
  no upstream format is claimed.
- The `gh-sdp` `0.1.0` client target and SDP Toolkit `0.2.0` compatibility target
  are separate proposals and both remain unreleased.
- Go is a preferred direction only; Study and Architecture retain the decision.
- `VER-SPS-001` passed as Slice evidence confirmed by independent follow-up
  review. It is not release-gate verification; the proposed release remains
  unreleased and its verification/review arrays remain empty.

### Independent Review Pass 1

Fresh Reviewer `/root/phase1_reviewer` inspected committed candidate
`381b38ed6fc59936e7e2cb30b877e156ff57914c` and draft PR #1. The disposition was
changes required: zero blocking, zero high, four medium, and one low finding.
`SDP/CodeReview/REV-SPS-001-001.md` is the authoritative review record. The four
medium findings were addressed in committed candidate
`bb16bee815d40742f29e15eeac4254bd1310376e`; the low
empty-repository-description finding is unchanged for Steering Group disposition.

### Independent Review Pass 2

Fresh Reviewer `/root/phase1_followup_reviewer` inspected exact pushed draft PR
head `f9d97c5a19d7b7f035f8dbfa2dd1f9e464b4378f`, the complete diff, live PR
metadata, and pinned canonical source. `REV-SPS-001-002` approved `SPS-001`
closure at `2026-07-13T10:16:59Z` with zero blocking, high, or medium findings
and one unchanged low finding. The approval is limited to Slice closure and does
not approve publication, Toolkit support, Steering acceptance, or Phase 2.

### Residual Risks, Gaps, and Discoveries

- The canonical project-manifest schema has no fields for repository URL,
  lifecycle phase, Toolkit compatibility target, or document/stable IDs. These
  facts remain in supported manifest fields plus authoritative project records;
  no unsupported manifest fields were invented.
- The canonical `validate_sdp.py` validates the Toolkit repository's own layout
  and cannot be represented as a consuming-project validator. Consumer checks
  therefore compose the canonical JSON schemas with project-local YAML/NDJSON,
  path/ID, hash, repository-layout, and content-boundary checks.
- The initial reviewed commit, draft PR, committed remediation, and exact pushed
  follow-up candidate exist. Exact validation and fresh review passed with no
  blocking, high, or medium findings.
- Normal same-version installation is idempotent, but the optional initializer is
  not idempotent after correcting its Toolkit-specific release seed. The current
  mitigation is preview-first, one-time use; correcting the reusable installer
  requires a separate authorized change in `Hans-Einar/SDP`.
- Detailed artifact/provenance, compatibility, rollback, alternate-source,
  authentication, offline, and platform filesystem policy remains input to
  Study. No later-lifecycle template was populated.
- `Ledger.ndjson` intentionally remains unchanged after the original
  `release-planned` event. The canonical ledger schema supports release events,
  not generic Sprint, Slice, or review completion; inventing such an event would
  misstate schema authority.
