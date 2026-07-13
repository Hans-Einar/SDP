# Sprint-001 implementation notes

## SPS-001

Status: active

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

Pending delegated Worker results.

### Verification and review

Pending exact candidate evidence and fresh independent review.
