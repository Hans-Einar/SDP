# Sprint-001 — Phase 1 SDP Bootstrap and Mandate

Status: complete
Iteration: SPI-001
Completed Slice: SPS-001
Mandate: MAN-001
Proposed gh-sdp target: REL-0.1.0 (`0.1.0`, unreleased)
Proposed first SDP Toolkit compatibility target: `0.2.0` (unreleased)

## Sprint Goal

Establish the repository-local SDP governance foundation and a precise Mandate
for `gh-sdp` without creating any product implementation or advancing into Study,
Requirements, Architecture, Design, or implementation planning.

## SPI-001 — Repository Bootstrap and Mandate

Status: complete

### SPS-001 — Authoritative Bootstrap, Mandate, and Initial Traceability

Status: complete

#### Goal

Bootstrap the current authoritative SDP project structure, replace Toolkit seed
facts with truthful `gh-sdp` project facts, define project-specific agent rules,
write Mandate `MAN-001`, establish initial traceability and verification evidence,
and prepare a reviewed draft pull request.

#### Why Now

The repository has only its initial empty commit. Governance, product boundaries,
version identities, verification expectations, and handoff records must exist
before later lifecycle or production work can begin.

#### Expected Files

- Toolkit-managed bootstrap files: `AGENTS.md`, `.codex/skills/`, and
  `SDP/Framework/`.
- Project-owned governance: `AGENTS-project.md`,
  `SDP/SDP-project.manifest.yaml`, `SDP/01--Mandate/mandate.md`, and the untouched
  lifecycle placeholders installed by the authoritative bootstrap.
- Phase 1 operating records under `SDP/Sprints/Sprint-001/`,
  `SDP/Traceability/`, `SDP/Verification/`, `SDP/CodeReview/`, and
  `SDP/Releases/REL-0.1.0.yaml`.
- No product source, build, packaging, workflow, or executable files.

#### Invariants

- `Hans-Einar/SDP` remains an independent, unmodified upstream repository.
- The project contains exactly one Git repository rooted at `gh-sdp`.
- Toolkit-managed files match source commit
  `bc110bb5fd60009ba67015cf640ad6ddbfe1b04b`.
- SDP Toolkit `0.2.0` is described as an unreleased source target, never as a
  published or supported release.
- gh-sdp client `0.1.0` and SDP Toolkit `0.2.0` remain separate identities.
- Existing project-owned content is preserved; Toolkit seed facts that describe
  the SDP repository itself are not retained as `gh-sdp` facts.
- Cross-platform, preservation, security, and reliability intent stays at Mandate
  level and does not become an irreversible design decision.

#### Non-Goals

- Product code, a Go module, CLI commands, handlers, installer/downloader/archive
  logic, migration logic, build or release workflows, tests for product behavior,
  or packaging.
- Study, Requirements, Architecture, Design Analysis, Design, Implementation
  planning, or Phase 2 work beyond installer-created placeholders.
- Changes, tags, releases, or pull requests in `Hans-Einar/SDP`.
- Publishing either gh-sdp `0.1.0` or SDP Toolkit `0.2.0`.

#### Traceability IDs

- Mandate: `MAN-001`.
- Sprint: `Sprint-001`.
- Iteration: `SPI-001`.
- Slice: `SPS-001`.
- Proposed client target: `REL-0.1.0`.

#### Verification Evidence

- Authoritative SDP source revision, release state, installer preview/apply, and
  installer idempotency are recorded exactly.
- Project and release manifests validate against the canonical JSON schemas.
- all YAML and NDJSON records parse; stable IDs and paths resolve; relations agree
  with document status.
- only the repository-root `.git` exists and no product implementation exists.
- a fresh independent Reviewer inspects the actual diff and evidence; every
  blocking, high, and medium finding is resolved and re-reviewed.

#### Completion Signal

- The authoritative bootstrap, `MAN-001`, manifest, traceability, verification,
  implementation notes, and handoff are internally consistent.
- Deterministic checks pass or any authoritative tooling gap is evidenced without
  fabrication.
- The focused branch is committed and pushed, a draft PR exists against `main`,
  and independent review has no unresolved blocking, high, or medium findings.
- The repository contains no production implementation and stops at the Phase 1
  boundary for Steering Group assessment.

#### Completion Evidence

- `VER-SPS-001` passed deterministic validation against committed remediation
  `bb16bee815d40742f29e15eeac4254bd1310376e` and exact pushed review head
  `f9d97c5a19d7b7f035f8dbfa2dd1f9e464b4378f`.
- Fresh follow-up review `REV-SPS-001-002` approved `SPS-001` closure with zero
  blocking, high, or medium findings and one non-blocking low finding.
- The draft PR remains open for Steering Group assessment. `REL-0.1.0`, client
  `0.1.0`, and Toolkit `0.2.0` remain proposed and unreleased; no publication or
  Phase 2 authorization is implied.
