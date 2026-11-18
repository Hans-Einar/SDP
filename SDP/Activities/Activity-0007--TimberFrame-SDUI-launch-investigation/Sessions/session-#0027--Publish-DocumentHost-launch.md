# Publish DocumentHost launch

| Field | Value |
| --- | --- |
| Session reference | SESSION-SDP-0027 |
| Status | active |
| Activity | ACTIVITY-SDP-0007 |
| Phase | 06-Verification |
| Phase contributions | 06-Verification |
| Plan | MAINT-SDP-0023 |

## Goal

Publish signed SDPTool and gh-sdp2.4.0 and verify public upgrade in isolation

## Linked records

- plan: [MAINT-SDP-0023](../../../Maintenance/R240/Plan.md)

## T001 — Release authorization and preflight

Owner requests publication of the completed DocumentHost launcher. MAINT0023 selects
2.4.0 as an additive minor release. Plan preflight records clean historical worktrees,
shared write ownership and preserved ledger lock. Loaded SDP, Release, Versioning
and Planning; reused shared document and lifecycle routines. Product candidate735c865
and model0.1.0 are already verified; only release metadata/default/predecessor changes
are selected. One extension worker and one consolidated review conserve quota.

| State | Step | Outcome |
| --- | --- | --- |
| on-going | R240-M1 | Build/sign and qualify exact release candidate |
| planned | R240-M2 | Publish, verify public upgrade, reconcile records |

Next: exact signed package, tests/CI, review, authorized publication and closeout.
