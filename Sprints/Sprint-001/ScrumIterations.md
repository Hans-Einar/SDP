# Sprint-001 — Portable installation and project-validation contract

Status: active
Release target: `REL-0.2.0` (`0.2.0`, unreleased)
Started: `2026-07-13T20:54:45Z`

## Sprint goal

Make the SDP Toolkit installation boundary canonical, portable and safe for both
the supported PowerShell installer and independent clients such as `gh-sdp`,
without publishing Toolkit `0.2.0`.

## SPI-001 — Installation contract hardening

Status: active

### SPS-001 — Canonical installation contract v1

Status: rework

#### Goal

Deliver one schema-validated installation contract that authoritatively defines
copied and generated files, destinations, ownership and update policies; make the
PowerShell installer consume it; separate neutral project templates from this
repository's live records; and add an explicit consuming-project validator.

#### Why now

`Hans-Einar/gh-sdp` draft PR #1 proved that the current initializer can propose
the Toolkit repository's `REL-0.2.0` record in a consuming project and that an
external installer otherwise has to reconstruct PowerShell-only behavior. The
same bootstrap also exposed missing reusable schemas and no project validation
mode.

#### Authoritative study and decisions

- Contract study: `Sprints/Sprint-001/InstallationContractStudy.md`
- Toolkit release state: `SDP.manifest.yaml`
- Release record: `Releases/REL-0.2.0.yaml`
- Upstream consumer evidence: `Hans-Einar/gh-sdp` draft PR #1 at exact reviewed
  closure head `ed205c1ef193ab8a6e5cd1c50e558c3049ce6def`

#### Scope

- Add canonical installation-contract v1 plus JSON Schema.
- Make the PowerShell installer derive all installable entries from the contract.
- Add deterministic portable preview-plan JSON suitable for conformance tests.
- Create a physically separate neutral project-template payload.
- Prevent repository-instance Sprint, release, review, verification and ledger
  state from entering consuming projects.
- Add Toolkit and consuming-project validator modes.
- Add the minimum reusable CurrentIndex, Relations and generic Ledger schemas,
  retaining a specialized release-event contract.
- Validate release/Fix records, release-note structure, installed skills and
  deterministic path references where present.
- Document external-client, archive, ownership, migration and Analyzer contracts.
- Add regression, archive and cross-platform-oriented deterministic tests.
- Update Toolkit `0.2.0` Unreleased notes and SDP records.

#### Expected files

- `Toolkit/SDP-install.manifest.json`
- `Toolkit/schemas/SDP-install-manifest.schema.json`
- `Toolkit/schemas/current-index.schema.json`
- `Toolkit/schemas/relations.schema.json`
- `Toolkit/schemas/ledger-event.schema.json`
- `Toolkit/schemas/release-event.schema.json`
- `Toolkit/project-templates/**`
- `Toolkit/scripts/Install-SDP.ps1`
- `Toolkit/scripts/validate_sdp.py`
- `Toolkit/tests/Install-SDP.Tests.ps1`
- `Toolkit/tests/test_validate_sdp.py`
- root and Toolkit READMEs, installation/manifest/project/validation/Analyzer and
  release/archive documentation, and examples
- `RELEASE-NOTES.md`, `SDP.manifest.yaml`, `Releases/REL-0.2.0.yaml`
- this Sprint's implementation, verification, review and handoff records
- `Traceability/CurrentIndex.yaml`, `Traceability/Relations.yaml` and
  `Traceability/Ledger.ndjson`

The exact file set may add narrowly named fixtures or documentation when needed
to satisfy the contract, but may not broaden into a `gh-sdp` implementation or
release publication.

#### Invariants

- `0.2.0` remains `unreleased`; `gitTag`, release commit and GitHub Release fields
  remain null and no tag or Release is created.
- The installation manifest is authoritative. PowerShell contains no second
  installable-file inventory that can silently drift.
- Every installable source and destination is normalized, relative and free of
  traversal; entry IDs are stable and unique.
- Normal install, forced managed refresh and one-time structure initialization
  preserve existing compatibility, backup and project-ownership guarantees.
- Preview and JSON-plan modes make zero target-project mutations.
- Existing project-owned files are never replaced, including under force.
- A fresh project never receives live Toolkit repository records, including
  `Releases/REL-0.2.0.yaml`, populated ledger history, active work, review or
  verification evidence.
- Release source archives work without `.git`, fixed absolute paths, Windows path
  semantics or PowerShell interpretation of the contract; unknown source commit
  is represented truthfully as null.
- Schemas validate stable machine-readable boundaries without claiming full
  semantic validation of arbitrary Markdown.

#### Non-goals

- Implementing or modifying `Hans-Einar/gh-sdp`.
- Publishing, tagging or creating a GitHub Release for Toolkit `0.2.0`.
- Changing the Toolkit target to `0.3.0`.
- Adding a dedicated release asset without evidence that source archives fail.
- Defining a comprehensive ontology for all project documents or arbitrary
  Markdown semantics.
- Prescribing private Go implementation structure beyond the public contract.

#### Verification

- Pre-change baseline: canonical validator passed; 9 Python tests passed; Windows
  installer fixtures passed on PowerShell `5.1.26100.8655`, Python `3.11.5` and
  Git `2.43.0.windows.1`.
- Validate Toolkit mode and every JSON Schema.
- Validate a real generated consuming-project fixture and malformed variants.
- Run the full Python unit suite.
- Run the full PowerShell installer fixture suite on Windows.
- Exercise empty, non-empty, legacy, upgrade, force, backup, initialization,
  repeated initialization, preview, downgrade, unsupported-schema and sibling
  path cases.
- Exercise a `gh-sdp`-like bootstrap, one-root-`.git` invariant and explicit
  absence of `SDP/Releases/REL-0.2.0.yaml` and live Toolkit history.
- Exercise extracted-source-archive behavior with no `.git`.
- Compare deterministic JSON plans and manifest inventory coverage.
- Run `git diff --check` and repository validation against `origin/main`.
- Require Linux contract tests and Windows installer tests from GitHub Actions on
  the pushed draft-PR head.
- Require a fresh independent Reviewer to inspect the actual diff and evidence.

#### Review gate REV-SPS-001-001

The first independent review of candidate
`877fa7693359e7ff74e8dde9284654a8a61ef341` returned changes required. It found
five high and ten medium issues covering physical containment, YAML and manifest
preflight, Windows path aliases, exact `gh-sdp` compatibility, plan/apply
equivalence, SemVer, traceability, publication semantics and governing metadata.
The authoritative finding and remediation contract is
`CodeReview/REV-SPS-001-001.md`. `SPS-001` remains active in rework until a fresh
Worker remedies the findings and a separate fresh Reviewer approves the new
exact candidate.

Remediation Worker commit
`f81a75b96fdcc47bff4a11e9381bb62ff459a494` addresses the complete review
contract. Master verification `VER-SPS-001-002` passed 72 Python tests, the full
PowerShell suite, both Toolkit validator invocations and exact offline `gh-sdp`
project validation. The Slice has returned to review; it is not complete until a
fresh independent Reviewer approves the integrated candidate and draft-PR CI
provides Linux and Windows evidence.

Second independent review `REV-SPS-001-002` returned changes required on exact
integrated candidate `8f22614ada03effe0d7f044315ee76414554d098`. It closed
twelve original findings but confirmed two residual highs—physical local/UNC
root identity and destination-topology partial mutation—and one residual medium
covering four missing Python schema/capability pairing checks. The Slice returned
to rework for this bounded three-finding contract.

Bounded Worker commit `25fdf5cfe5a119192f512bf5322a15776e26836f`
closes the three residual findings. Master verification `VER-SPS-001-003`
passed parser/schema sanity, 73 Python tests, the full Windows PowerShell suite,
both Toolkit validators and exact pinned `gh-sdp` project validation. The Slice
has returned to review and still requires a separate fresh approval before draft
PR publication.

Third independent review `REV-SPS-001-003` closed H3-M8-R2 and M10-R2, but
reproduced one high H1-R3 wrong-project mutation caused by normalization of a
distinct extended-only trailing-space/dot path. The Slice returned to rework for
this single bounded path-normalization finding.

#### Completion signal

The Slice is complete only when the canonical contract is consumed by
PowerShell, neutral templates and project validation work, all deterministic
tests and archive checks pass, independent review has no unresolved blocking,
high or medium findings, SDP records are current, a draft PR against `main` is
open, and Toolkit `0.2.0` is still unreleased with no tag or GitHub Release.
