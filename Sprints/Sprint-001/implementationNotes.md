# Sprint-001 implementation notes

## SPS-001

Status: active (`SPS-001` revision 1)

### Master preparation

- Confirmed `main` matched `origin/main` at
  `bc110bb5fd60009ba67015cf640ad6ddbfe1b04b` and created branch
  `codex/sdp-install-contract-v1`.
- Confirmed no Git tags or GitHub Releases exist for `Hans-Einar/SDP` and no
  other SDP pull request is open.
- Inspected `Hans-Einar/gh-sdp` draft PR #1 at exact head
  `ed205c1ef193ab8a6e5cd1c50e558c3049ce6def`. Its Phase 1 evidence reproduces
  the `REL-0.2.0` repeat-initialization proposal and records the validator and
  schema gaps addressed by this Slice.
- Recorded the alternatives and selected authoritative explicit JSON manifest,
  separated neutral project templates and runtime JSON plan design in
  `InstallationContractStudy.md`.
- Registered active `Sprint-001 / SPI-001 / SPS-001` in CurrentIndex and
  Relations before product implementation.

### Baseline evidence

Environment: Windows `10.0.26200`, Windows PowerShell `5.1.26100.8655`, Python
`3.11.5`, Git `2.43.0.windows.1`.

- `python Toolkit/scripts/validate_sdp.py` — passed.
- `python -m unittest discover -s Toolkit/tests -p "test_*.py" -v` — 9 tests
  passed.
- `./Toolkit/tests/Install-SDP.Tests.ps1` — installer fixture tests passed.

The first environment-probe command attempted unavailable `pwsh` and stopped
before tests; the successful baseline used the active Windows PowerShell host.

### Implementation

- Added one explicit 40-entry installation manifest with two generators and 22
  exclusions, plus install-manifest and deterministic-plan schemas.
- Moved neutral, project-owned seeds into a 22-file
  `Toolkit/project-templates/` boundary; live root records are excluded.
- Refactored the PowerShell compatibility installer to derive entry selection,
  policy, content and actions from the manifest and expose read-only `-PlanJson`.
- Added archive-safe source-commit handling, pre-mutation policy/destination
  gates and the corrected neutral initialization behavior.
- Added Toolkit/project validator modes, CurrentIndex, Relations and generic
  Ledger schemas, specialized release events and namespaced extensions.
- Added 49 Python tests, expanded Windows installer fixtures, public contract and
  validation documentation, examples, release notes and synchronized capability
  metadata.

Worker implementation is committed at
`ce9278b8b78f9c320a65799aefd13101582d1eb8`.

#### Revision 1 downstream-conformance implementation

- Added the required
  `orderingPolicy: migration-first-manifest-order-v1` authority to the
  installation manifest and plan. PowerShell and Python now enforce migration
  prefix, exact manifest-array entry order, adjacent identity-equal
  backup/mutation pairs, final contiguous sequence assignment and the canonical
  single-action blocked shape.
- Completed AGENTS conflict preservation with exact-byte SHA-256,
  `targetSourceSha256`, `destinationPrecondition: absent`, identical-target
  idempotency, closed fatal failure classes and immediate pre-apply source and
  destination race checks. Migration uses exclusive creation before ordinary
  entry mutations and never overwrites project-owned content.
- Added `Toolkit/conformance/install-v1/`, a schema-validated language-neutral
  package with 17 portable scenarios and committed authoritative plan/failure
  outcomes. Its Python harness validates without PowerShell, compares the
  PowerShell reference under Windows PowerShell 5.1 or PowerShell 7, checks
  planning immutability and apply assertions, and writes new candidates only
  through an explicit maintainer flag that CI never uses.
- Expanded PowerShell and Python regression coverage, integrated Linux
  conformance validation and Windows reference comparison in CI, and updated
  the public installation, migration, validation, Toolkit, example, Analyzer
  and Unreleased release-note surfaces.
- Worker validation passed Toolkit mode, 80 Python tests, full Windows
  PowerShell 5.1 and PowerShell 7.6.3 installer suites, and all 17 conformance
  scenarios against both PowerShell hosts. The host could not create file
  symlinks, so the optional symlink cases reported their documented skip while
  directory/reparse fail-closed coverage passed.

### Verification and review

Master verification is recorded in `Verification/VER-SPS-001.md` and passed on
exact clean candidate `ce9278b8b78f9c320a65799aefd13101582d1eb8`: Toolkit
validation, 49 Python tests, expanded PowerShell fixtures, archive/project
validation, diff check and publication-state checks passed.

Fresh independent review `REV-SPS-001-001` inspected integrated candidate
`877fa7693359e7ff74e8dde9284654a8a61ef341` and returned changes required: five
high and ten medium findings. The original verification is retained as evidence
for its exact commit, but is superseded as a completion gate. Remediation must
cover every finding in `CodeReview/REV-SPS-001-001.md`, rerun the full evidence
matrix and receive a new independent review.

### Remediation history

- Remediation Worker commit
  `f81a75b96fdcc47bff4a11e9381bb62ff459a494` closes all five high and ten
  medium findings with expanded installer, validator, schema, exact `gh-sdp`,
  documentation and adversarial-test coverage.
- Master verification `VER-SPS-001-002` passed on that exact clean candidate:
  72 Python tests, the full Windows PowerShell fixture suite, both Toolkit
  invocations, exact pinned project validation and diff/publication checks.
- File-symlink creation was unavailable on the local Windows host;
  junction/reparse coverage passed. The later hosted Windows suite passed.
- Linux contract evidence passed in draft-PR run `29298949100`.
- This is Slice verification, not the final `0.2.0` release gate.

### Second review gate

Fresh review `REV-SPS-001-002` inspected candidate
`8f22614ada03effe0d7f044315ee76414554d098`. Twelve original findings were
closed. Two high findings remain for physical local/UNC root identity and
destination ancestor/prefix topology; one medium remains for four missing
Python governing-pair assertions. A new Worker must implement exactly this
bounded remediation and add the required zero-mutation/parameterized tests.

### Bounded remediation

Worker commit `25fdf5cfe5a119192f512bf5322a15776e26836f`
implements OS-backed symmetric Windows directory identity, complete destination
prefix/ancestor preflight and all seven Python governing-pair checks. It also
clarifies dirty-checkout `sourceCommit` semantics. Master verification
`VER-SPS-001-003` passed 73 Python tests, the full PowerShell suite, all three
validator invocations and clean publication-state checks. The local/UNC,
extended and integrated short-name cases all ran; only file-symlink creation was
unavailable on this host.

### Third review gate

Fresh review `REV-SPS-001-003` closed the destination-topology and seven-pair
findings, but reproduced H1-R3: a normalization-sensitive `\\?\` project path
could be redirected to a different ordinary directory and applied there. The
next Worker is limited to raw extended-path equivalence/rejection and trailing-
space/dot zero-mutation regressions.

### Extended-path remediation

Worker commit `1c650578109a361cf2f187da665eabdc4ac85e81`
validates raw extended drive/UNC input before any prefix removal. Master
verification `VER-SPS-001-004` passed the full matrix, including ProjectRoot,
BackupRoot and source trailing-space/dot zero-mutation cases for drive and UNC,
plus canonical extended and 8.3 positive cases.

### Local approval

Fresh review `REV-SPS-001-004` approved integrated candidate
`66759afb45813f9c890e67d4938e5b796b9ed05d` with no unresolved blocking, high,
medium or low findings. Draft-PR Linux contract and Windows installer checks are
the remaining Slice evidence gate.

### Draft-PR CI remediation

Draft PR `https://github.com/Hans-Einar/SDP/pull/4` was opened from
`codex/sdp-install-contract-v1` against `main`.

- Run `29295427882` passed Linux contracts and failed Windows because the job
  had not installed `jsonschema`. Fresh Worker/review commit
  `64a89b0f4e9378f97c727cf3f1b5a7fdfa39f543` added Python 3.12 setup and the
  existing declared test requirements.
- Run `29295843790` passed dependencies and Linux contracts, then reproduced a
  PowerShell 7 repeat-install failure. `ConvertFrom-Json` had coerced a quoted
  RFC 3339 scalar to DateTime before strict YAML validation. Fresh Worker/review
  commit `64d874b6c6a19a61d8e06a1f39d525b53cc50d90` preserved exact quoted
  strings and added a repeat PlanJson/apply regression.
- A complete portable PowerShell 7.6.3 run then exposed a UNC share-root parent
  traversal difference. Fresh Worker commit
  `70e38d663846c67a984363f80defe4d44e798fa9` added namespace-root-aware parent
  traversal. Its fresh review found a pre-existing medium trailing-separator
  ProjectRoot false rejection, which fresh Worker commit
  `7c78132f31788d7494eb93c6cba5310c5503e416` closed with local and loopback UNC
  PlanJson/apply regressions.

Final integrated review `REV-SPS-001-005` approved `7c78132` with zero blocking,
high, medium or low findings. Independent local evidence passed full PowerShell
7.6.3 and Windows PowerShell 5.1 suites, 73 Python tests, every validator mode,
all schemas/parsers and the path-safety matrix.

### CI and closure

GitHub Actions run `29298949100` passed on exact head
`7c78132f31788d7494eb93c6cba5310c5503e416`:

- Linux `contracts`: passed in `9s`, including 73 Python tests.
- Windows `installer`: passed in `5m12s`, ending with
  `Installer fixture tests passed.`

The exact evidence is `VER-SPS-001-005`. `SPS-001`, `SPI-001` and `Sprint-001`
are complete. The draft PR remains open, and Toolkit `0.2.0` remains unreleased
with no tag or GitHub Release. Release verification/review records remain a
separate future gate.

### Downstream reopening — revision 1

The `gh-sdp` Steering Group supplied a downstream contract-consumability
assessment through the human repository owner for exact upstream head
`bf20832bed618ab240cf87c17517fc31ea721311`. The disposition is
`UPSTREAM_REWORK_REQUIRED`; the evidence is recorded as
`REV-SPS-001-006`. No GitHub review URL is claimed for this external assessment.

Three Medium findings reopen the Slice: normative plan ordering is incomplete,
AGENTS deterministic migration collision/idempotency/race behavior is
incompletely public, and no shared language-neutral conformance package exists.
CurrentIndex and Relations now identify
`Sprint-001 / SPI-001 / SPS-001 / revision 1` as active. Historical evidence is
preserved, but cannot close the revised candidate.

The bounded implementation contract is in
`Sprints/Sprint-001/ScrumIterations.md`. The next Master verification ID is
`VER-SPS-001-006`; the next fresh implementation-review ID is
`REV-SPS-001-007`. PR #4 remains the only PR, remains draft, and matched the
assessed head when reopened. Toolkit `0.2.0` remains unreleased; local tags,
remote tags and GitHub Releases were absent.

### Revision 1 hosted-CI remediation

`origin/main` advanced with feature-governance documentation while revision 1
was in progress. It was merged additively at
`6b326543b3712562338e26dfe978b04440411e84`; the only conflict retained both the
installation-contract and feature-governance README links. PR #4 returned to
`MERGEABLE/CLEAN` and remained draft.

GitHub Actions run `29585982002` passed Linux `contracts` including portable
conformance, then failed Windows job `87902965964`. The hosted runner exercised
the locally unavailable file-symlink preservation fixture and observed no
`sdpFailureClass`; the expected class was
`agents-migration-destination-unsupported-object`. The containment guard
rejected the reparse path before the deterministic AGENTS destination state was
classified. This is a bounded failure-class propagation defect, not evidence of
target mutation or project-owned overwrite. The remediation contract is
recorded in `ScrumIterations.md`.

The bounded Worker remediation classifies only the deterministic preservation
destination before its containment assertion, retains that assertion for every
state, and translates a containment rejection only when the destination was
already classified as an unsupported object. Absent and regular-file
destinations still pass through the existing physical containment checks;
ancestor-link failures retain their existing fail-closed behavior, and
apply-time destination races remain
`agents-migration-destination-changed`.

Regression coverage now proves the stable unsupported-object class and unchanged
project-owned content for a regular directory, a directory junction/reparse
destination and, when the host permits creation, a file symlink. Exact final
Worker evidence passed on Windows `10.0.26200`:

- the complete installer fixture suite under Windows PowerShell
  `5.1.26100.8875` and portable PowerShell `7.6.3`, both ending with
  `Installer fixture tests passed.`;
- Toolkit validation both bare and against `origin/main`, all 80 Python tests,
  portable conformance-package validation and exact pinned `gh-sdp` project
  validation;
- all 17 committed conformance scenarios against both PowerShell reference
  hosts; and
- both PowerShell parser checks, all 12 schema JSON parses, repository
  JSON/YAML/NDJSON parsing and working-tree/full-branch diff checks.

The local host still could not create the optional file symlink, so that fixture
reported its documented skip while the directory junction/reparse case ran.
Hosted Windows exact-head evidence and fresh review remain Master gates; no
verification or review record is created by this Worker pass.
