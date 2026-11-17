# Publish capacity correction

| Field | Value |
| --- | --- |
| Session reference | SESSION-SDP-0024 |
| Status | active |
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
