# Codex 0.164.0 and compatible later host support

| Field | Value |
| --- | --- |
| id | KB-SDP-057 |
| project | SDP |
| type | Change |
| CardState | in-progress |
| created | 2026-10-23T22:10:45.708406+00:00 |
| source | Owner Session0007 T029, 2026-10-24 |
| PlanId | PLAN-SDP-0018 |
| Systems | SDPTOOL, ProjectGovernance |

## Authorized change

Owner authorizes supporting Codex 0.164.0 and later compatible versions, rather
than replacing one release pin with another. Preserve actual host version and
all authority/recovery guards. [PGC1](../../05--Implementation/SDPTool/ProjectGovernance/Codex-Compatibility.md)
selects the bounded contract and verification; [Session0007](../../Sessions/session-%230007--Project_governance.md)
owns continuity. No complete leader workflow, Ponsse trial, merge or release selected.

## Worklog

| Time | Actor / event | Result and next action |
| --- | --- | --- |
| 2026-10-23T22:10:45.708406+00:00 | Codex; EVT-KB-SDP-000389 | Register selected PGC1; contract/preflight delivered. Implement minimum/profile admission and actual host identity with recovery preservation. |

2026-10-23T22:26:13.624760+00:00: PGC1-B delivers stable >=0.164.0 eligibility with observed-version identity,
wire-schema alignment and legacy/recovery preservation. Scoped race/negative/future
fixture checks pass; final PGC1-C verification next. See
[evidence](../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGC1.md).
CardState stays in-progress. EVT-KB-SDP-000395.
