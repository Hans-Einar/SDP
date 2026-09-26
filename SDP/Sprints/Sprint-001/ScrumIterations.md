# Sprint-001 — Portable installation and project-validation contract

Status: complete (including `SPS-001` revision 1)
Release target: `REL-0.2.0` (`0.2.0`, unreleased)
Started: `2026-07-13T20:54:45Z`
Original completion: `2026-07-14T01:42:15Z`
Reopened: `2026-07-17T12:42:02Z`
Revision 1 completed: `2026-07-17T15:24:56Z`

## Sprint goal

Make the SDP Toolkit installation boundary canonical, portable and safe for both
the supported PowerShell installer and independent clients such as `gh-sdp`,
without publishing Toolkit `0.2.0`.

## SPI-001 — Installation contract hardening

Status: complete

### SPS-001 — Canonical installation contract v1

Status: complete

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

H1-R3 Worker commit `1c650578109a361cf2f187da665eabdc4ac85e81`
rejects normalization-sensitive extended paths before prefix removal while
preserving canonical extended drive, UNC and available 8.3 paths. Master
verification `VER-SPS-001-004` passed 73 Python tests, the full PowerShell suite
and all validator modes. The Slice returned to review for a new independent
approval.

Fresh independent review `REV-SPS-001-004` approved exact integrated candidate
`66759afb45813f9c890e67d4938e5b796b9ed05d` with no unresolved blocking, high,
medium or low findings. The Slice remains in review only for draft-PR Linux and
Windows GitHub Actions evidence and final SDP integration.

Draft PR `Hans-Einar/SDP#4` exposed two hosted-runtime gaps before completion.
Run `29295427882` showed that the Windows job did not install the declared
Python schema dependency. After commit
`64a89b0f4e9378f97c727cf3f1b5a7fdfa39f543`, run `29295843790` reached the
installer and exposed PowerShell 7 quoted-timestamp coercion. Fresh bounded
Workers and separate Reviewers closed that defect, a PowerShell 7 UNC
share-root traversal difference, and a trailing ProjectRoot separator false
rejection in commits `64d874b6c6a19a61d8e06a1f39d525b53cc50d90`,
`70e38d663846c67a984363f80defe4d44e798fa9` and
`7c78132f31788d7494eb93c6cba5310c5503e416`.

Final independent review `REV-SPS-001-005` approved candidate `7c78132` with no
blocking, high, medium or low findings. GitHub Actions run `29298949100` then
passed the Linux contracts job and the full Windows installer job on that exact
head. `VER-SPS-001-005` is the authoritative completion evidence.

#### Completion evidence

The canonical contract is consumed by PowerShell, neutral templates and project
validation work, deterministic tests and archive checks pass on Windows and
Linux, independent review has no unresolved finding, SDP records are current,
and draft PR `#4` is open against `main`. Toolkit `0.2.0` remains unreleased with
no tag or GitHub Release. `SPS-001`, `SPI-001` and `Sprint-001` are complete;
this does not complete or publish `REL-0.2.0`.

#### Revision 1 — downstream contract-consumability remediation

Status: complete
Opened: `2026-07-17T12:42:02Z`
Completed: `2026-07-17T15:24:56Z`
Trigger review: `CodeReview/REV-SPS-001-006.md`
Assessed upstream head: `bf20832bed618ab240cf87c17517fc31ea721311`
Draft pull request: `Hans-Einar/SDP#4` (must remain draft)

The `gh-sdp` Steering Group assessed the completed candidate and returned
`UPSTREAM_REWORK_REQUIRED` through the human repository owner. The earlier
verification and approvals remain truthful for their exact candidates, but are
superseded as completion evidence by three unresolved Medium downstream
contract findings.

Revision 1 is a bounded correction inside `SPS-001`. It must:

- define and enforce the canonical v1 plan ordering: blocked plans contain one
  sequence-1 block; deterministic target-to-target migrations precede ordinary
  entries; ordinary entries follow exact installation-manifest array order;
  backup/mutation pairs are adjacent and identity-consistent; sequence is
  assigned last, contiguous from 1; clients must not privately sort actions;
- define the exact-byte SHA-256 `AGENTS.md` preservation destination and
  fail-closed absent, identical-file, differing-file, unsupported-object and
  plan/apply race semantics without overwriting project-owned content;
- add a versioned, language-neutral installation-v1 conformance package with a
  machine-readable scenario index, portable before-state declarations,
  authoritative expected normalized plans or closed-vocabulary fatal failure
  outcomes, preservation assertions and a Python reference harness;
- cover at least empty default/initialize, repeat default/initialize, legacy
  migration, new/already-preserved/different/invalid AGENTS hash targets,
  managed upgrade, same-version force, project-owned preservation, archive
  without Git, unsupported project and installed schemas, downgrade blocking
  and malformed manifest;
- update schemas, PowerShell, Python validation, tests, public documentation,
  examples and Unreleased notes only as needed for those three findings; and
- retain Windows PowerShell 5.1 and current PowerShell 7 compatibility, exact
  pinned `gh-sdp` project validation, source-archive support and all prior safety
  invariants.

Revision 1 does not authorize `gh-sdp` implementation, unrelated installer
redesign, Toolkit `0.2.0` publication, version change, tagging, a GitHub Release,
merging PR #4 or marking it ready. Completion requires `VER-SPS-001-006`, fresh
independent review `REV-SPS-001-007` with no unresolved Blocking/High/Medium
finding, Linux and Windows GitHub Actions on the exact new head, and a request
for a new bounded downstream acceptance pass.

##### Hosted Windows CI remediation

GitHub Actions run `29585982002` passed Linux contracts/conformance on merge head
`6b326543b3712562338e26dfe978b04440411e84` but failed the hosted Windows
installer job. The hosted runner could create the optional preservation-target
file symlink that the local host could not. Planning rejected that link before
`Get-PathObjectState` classified it, so the failure did not expose the required
stable `agents-migration-destination-unsupported-object` class.

The bounded remediation is to preserve every existing link/reparse containment
guard while mapping the deterministic AGENTS preservation destination's
unsupported link/reparse state to the public failure class before mutation. It
must add or refine the hosted-capable regression, preserve zero target mutation
and project-owned content, rerun both local PowerShell hosts and the full
conformance matrix, and obtain a new exact-head GitHub Actions pass. No other
installer behavior or contract vocabulary is in scope.

##### Revision 1 completion evidence

Exact contract candidate
`a37a8fa298fa3ace2ec5826bc3a30f26f90a64ac` implements the required
`migration-first-manifest-order-v1` policy, exact-byte AGENTS preservation
semantics and the 17-scenario language-neutral
`Toolkit/conformance/install-v1/` package.

Master verification `VER-SPS-001-006` passed the complete local matrix and
GitHub Actions run `29590324256` on that exact head. Linux `contracts` job
`87917499632` passed 80 tests and portable conformance. Hosted Windows
`installer` job `87917499648` passed the full fixture suite, exercised the
file-symlink regression and passed all 17 reference scenarios.

Fresh independent review `REV-SPS-001-007` completed all seven required
adversarial attempts and approved the exact candidate with zero Blocking, High,
Medium or Low findings. All three findings in `REV-SPS-001-006` are closed.

`SPS-001` revision 1, `SPI-001` and `Sprint-001` are complete. This returns the
repository to the downstream acceptance gate; it does not merge PR #4 or
complete/publish `REL-0.2.0`. PR #4 remains draft pending a new bounded
`gh-sdp` Steering Group disposition.
