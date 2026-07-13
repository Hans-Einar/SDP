# Sprint-001 handoff

Active Sprint: `Sprint-001`
Active Iteration: `SPI-001`
Active Slice: `SPS-001`
Status: rework required by independent review

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

Assign a fresh Worker the complete remediation contract in
`CodeReview/REV-SPS-001-001.md`. The Worker must address all five high and ten
medium findings, add the required adversarial and exact `gh-sdp` regressions,
and rerun the full evidence matrix. Then record a new exact verification
candidate and assign a separate fresh Reviewer. Do not open the draft PR or close
the Slice with unresolved blocking, high or medium findings. Linux GitHub Actions
evidence remains required before final closure.
