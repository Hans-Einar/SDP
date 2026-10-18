# PGL1 — Inspect registered project work

2026-10-10, PLAN-SDP-0018 / KB-SDP-052 / Session0007 T014.
Owner selects Project Leader/SAD orchestration after trying the one-Master lab.
[Design](../../../04--Design/SDPTool/ProjectGovernance/Orchestration.md) and
[handoff](../../../04--Design/SDPTool/ProjectGovernance/Project-Leader-Handoff.md)
retain the full workflow and six acceptance scenarios. PGL1 delivers only the first
inspection part of scenario 5, not the hierarchy or scheduler.

## Candidate and behavior

[Candidate manifest](candidate-PGL1.json) identifies baseline 0b0e82c plus exact
changed production/test/schema hashes and the built binary. Go 1.27.1, Linux amd64.
`sdp-controller overview --project ROOT --state DIR` displays every registered run;
`--json` returns one strict `sdp-governance-overview/1` object. This command starts
no host, creates no state and uses one operational generation. Source/coverage reads
are per run, not an atomic filesystem snapshot. The command does not enumerate all
cards/plans, other state directories or unregistered work.

Task, selected Session/step/card/plan, routine state, coverage and last host observation
remain separate. Role evidence and explicit skill input acceptance do not imply
compliance. Tool counts are retained latest-attempt completion observations, including
failed/declined operations. Native children, current liveness and owner acceptance
remain unobserved. Input identity is retained separately from attempt identity.
Human output neutralizes terminal controls; JSON excludes transcript/owner-input
content and capability credentials. Existing worker MCP privileges are unchanged.

## Checks and observed outcome

- `go test -race ./governance ./cmd/sdp-controller` passes on the final code.
  Cases cover absent state without creation, two distinct runs, stable ordering,
  stale/deleted governing sources, corrupt store rejection, historical epochs/gaps,
  completed host versus verifying routine, private-content exclusion, control-code
  escaping, explicit selection and output errors.
- `go test -race ./governancecontroller ./governancemcp` passes on the candidate
  before the final additive overview skill-input identity fields. Those packages
  do not consume Overview; existing provenance behavior is also retested by the
  final governance suite. No mocked host is described as live execution.
- `go vet ./governance ./governancecontroller ./governancemcp ./cmd/sdp-controller`
  passes; final changed overview/core/command paths were checked again.
- Compiled human/JSON commands inspect the actual owner lab and an absent state
  directory. Both JSON objects pass the new schema. No empty-store directory is
  created; all 176 project/runtime file or symlink identities/modes are unchanged.
- Actual lab run 06a3e736-bbc5-44bf-a10d-5e57c4235130 is revision 24: routine running,
  last host status unknown, observation gap true. Two explicit skills have observed
  host input acceptance. This is historical inspection, not proof that a worker is
  currently running, finished or accepted. Do not repair/reset the owner trial as
  part of an inspection feature.

An initial synthetic provenance fixture incorrectly encoded nil lists; strict
stored-provenance validation rejected it. The fixture was corrected to the existing
non-null wire contract; production validation was retained. Management draft
validation also rejected an unsupported cardState payload field; it was removed
from the newly drafted event. Original historical/imported bytes remain intact.

The separate lab helper `overview` invokes a frozen preview binary; it leaves the
existing `lab`, candidate manifest, model task and runtime untouched. Its manifest
pins source files and executable. Preview usage does not need installation/release.
Toolkit/management, ledger-prefix/import-byte and maintained-link checks are recorded
at closeout. No new model task, independent review, online Steering installation,
merge, push or release is part of this result. PGL2 is next; full KB052 remains open.

Closeout checks: Toolkit validation passes; management replay passes with 62 cards,
39 management records and 529 events. All 192 checked local Markdown links resolve.
Both historical ledger prefixes and imported upstream event bytes are preserved;
final source hashes and the frozen preview binary match the candidate manifest.
Same-context boundary review found no added mutation/launch path; this is not an
independent review. PGL1 is delivered only at the stated inspection boundary.

Final preflight refresh observed concurrent Actions planning at a9edf28. The product
write boundary remains disjoint. Its newly committed event IDs superseded our draft
allocations; only our uncommitted events were reassigned before commit (PM177, KB335/336).
Original histories remain unchanged. See preflight-PGL1.json and Session T014.
