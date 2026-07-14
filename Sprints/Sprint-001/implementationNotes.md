# Sprint-001 implementation notes

## SPS-001

Status: review — locally approved; draft-PR CI required

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

### Residual limitations

- Remediation Worker commit
  `f81a75b96fdcc47bff4a11e9381bb62ff459a494` closes all five high and ten
  medium findings with expanded installer, validator, schema, exact `gh-sdp`,
  documentation and adversarial-test coverage.
- Master verification `VER-SPS-001-002` passed on that exact clean candidate:
  72 Python tests, the full Windows PowerShell fixture suite, both Toolkit
  invocations, exact pinned project validation and diff/publication checks.
- File-symlink creation was unavailable on the Windows host; junction/reparse
  coverage passed. A fresh Reviewer must decide the remediated gate on the
  integrated candidate.
- Linux evidence remains pending from draft-PR CI.
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
