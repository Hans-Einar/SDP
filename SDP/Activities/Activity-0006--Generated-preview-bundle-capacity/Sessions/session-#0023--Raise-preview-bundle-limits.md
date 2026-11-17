# Raise preview bundle limits

| Field | Value |
| --- | --- |
| Session reference | SESSION-SDP-0023 |
| Status | completed |
| Activity | ACTIVITY-SDP-0006 |
| Phase | 05-Implementation |
| Phase contributions | 05-Implementation |
| Plan | MAINT-SDP-0021 |

## Goal

Implement and verify fourfold generated output capacity without changing SDL source semantics

## Linked records

- plan: [MAINT-SDP-0021](../../../Maintenance/BundleCapacity/Plan.md)

## T001 — Owner-reported bundle capacity

Owner requests the reason for the 128-file/32-MiB preview failure and at least a
fourfold increase. Recovered completed Activities 1–5 and release 2.3.1 before selecting
this independent correction; Session0007 PG implementation is not resumed here.
Preflight is green on 23c5193, sdp/bundle-capacity, with two clean detached release
snapshots and only the preserved management lock untracked. Scope: shared output
capacity check, preview/source-reading callers, focused tests and contract records.
Existing semantic source limits and publication freshness/ownership remain intact.
Loaded SDP, Change Analysis, Worker and Planning plus their shared document/plan
routines. No delegation, native agent test or separate review is needed.
Next: implement BC1, verify boundary and composed-source behavior, record results.

## T001 result

BC1 implemented and targeted tests pass. The same sealed-output helper serves
preview and source reading, with actual usage in limit errors. Contract and
KB064/MAINT0021 projections are updated. No structural SDL/model change, separate
review or additional agent was required. This Session completes the bounded source
correction; any public release/installed rollout is a separate authorized task.
Verification: [BC1](../../../Maintenance/BundleCapacity/Verification.md).
