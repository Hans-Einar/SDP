# Publish capacity correction

| Field | Value |
| --- | --- |
| Session reference | SESSION-SDP-0024 |
| Status | completed |
| Activity | ACTIVITY-SDP-0006 |
| Phase | 06-Verification |
| Phase contributions | 06-Verification |
| Plan | MAINT-SDP-0022 |

## Goal

Publish signed SDPTool and gh-sdp2.3.2 and verify public extension upgrade in isolation

## Linked records

- plan: [MAINT-SDP-0022](../../../Maintenance/R232/Plan.md)

## T001 — Release authorization and preflight

Owner requests public distribution through `gh extension upgrade sdp`. Reopened
Activity0006 for the same need and selected MAINT0022; completed BC1 remains intact.
Preflight/branch/write ownership are recorded in that plan. Loaded SDP release and
versioning, reusing Worker/Planning and shared routines. Patch2.3.2 preserves
wire/schema/source semantics. One Sol/medium worker prepares the extension; one
combined independent release review follows exact signed candidate preparation.
No live extension or consuming-project Framework upgrade is authorized here.
Next: prepare clean candidate, package/sign, verify, publish and reconcile.

## T002 — Signed candidate qualification

R232-M1 has clean canonical f7a6c9f and extension c1d86bd packages. Nine signed
predecessor upgrades, extracted-source checks, targeted tests and the actual
135-file regression pass. Native helper toolchain differs from 2.3.1; recorded an
isolated startup smoke rather than claiming identical bytes. Independent review
continues in the same context. Final CI and public consumer checks remain pending.

| State | Step | Outcome / dependency |
| --- | --- | --- |
| completed | R232-M1 | Signed candidate, CI and independent review pass |
| completed | R232-M2 | Both public releases and isolated upgrade verified |

## T003 — Publication and closeout

Owner-authorized 2.3.2 is published in both repositories. Final CI, signature/hash
checks, independent combined review, all nine predecessor rehearsals, actual
isolated extension upgrade/fresh installation and public-default Framework checks
pass. Source and production capacity evidence show 512 files / 128 MiB available.
No live extension, real consuming project or main branch was modified. Loaded
Release/Versioning and reused SDP, Worker, Planning and shared document routines.
Both roadmap steps and Activity0006 are complete. The service recorded all four
whole-scope dispositions in operation 4ef36c57-ded1-4858-b862-f649feffc75a.
Final records are validated before pushing the closeout commit.
This is a work summary; the owner can now run `gh extension upgrade sdp`.

Final validation caught an invalid cross-domain link in the newly drafted
Maintenance completion event. Its links now reference the existing card and
predecessor Maintenance record; release identity remains in Traceability. All
committed ledger bytes remain an exact prefix. Management and Toolkit validation
pass, and the Activity service confirms completed state.
