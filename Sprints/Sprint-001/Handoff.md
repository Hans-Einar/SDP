# Sprint-001 handoff

Active Sprint: `Sprint-001`
Active Iteration: `SPI-001`
Active Slice: `SPS-001`
Status: locally approved; draft PR and CI pending

## Entry points

- Slice contract: `Sprints/Sprint-001/ScrumIterations.md`
- Contract study: `Sprints/Sprint-001/InstallationContractStudy.md`
- Implementation notes: `Sprints/Sprint-001/implementationNotes.md`
- Current state: `Traceability/CurrentIndex.yaml`
- Relations: `Traceability/Relations.yaml`
- Release target: `Releases/REL-0.2.0.yaml`

## Current boundary

Implement and verify the canonical installation contract v1 only. Do not modify
`Hans-Einar/gh-sdp`, publish Toolkit `0.2.0`, create a tag or create a GitHub
Release. Product implementation must be delegated to Workers and reviewed by a
fresh independent Reviewer.

## Next action

Push `codex/sdp-install-contract-v1`, open a draft PR against `main`, and require
the Linux contract and Windows installer GitHub Actions checks on the pushed
head. Local approval is `REV-SPS-001-004`. Keep `0.2.0` unreleased and do not
create a tag or GitHub Release. Record real CI evidence before final Slice
closure.
