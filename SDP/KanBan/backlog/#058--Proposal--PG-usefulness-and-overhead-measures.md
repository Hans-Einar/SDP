# Measure PG usefulness and coordination overhead

| Field | Value |
| --- | --- |
| id | KB-SDP-058 |
| project | SDP |
| type | Proposal |
| CardState | backlog |
| created | 2026-10-23T22:10:45.708406+00:00 |
| source | Agent recommendation responding to owner Session0007 T029 design/value question |
| next_review | First complete owner-to-reviewed-result PG workflow |
| Systems | ProjectGovernance |

## Proposal, not owner adoption

Evaluate PG by reduced owner coordination effort and reliable recovery, as well as
protocol correctness. Possible observations: owner interventions per delivered
assignment, time to find the current state, and time/effort to resume after a break.
Separate measured elapsed time from unknown active effort; no conversation content
telemetry or new automatic recorder is implied. Use existing observations where
possible. Scope and collection/privacy choices need a concrete design if selected.

## Acceptance and next action

At the first complete workflow, propose a small baseline and compare the process
with manual SDP use. Decide which measurements are useful before implementing any
instrumentation. This card does not activate work or add a PG readiness gate.

## Worklog

2026-10-23T22:10:45.708406+00:00: EVT-KB-SDP-000391 registers this agent recommendation; owner has asked for
ideas but has not adopted or authorized instrumentation. Compatibility remains focus.
