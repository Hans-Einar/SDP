# Portable human-readable SDPTool output

| Field | Value |
| --- | --- |
| id | KB-SDP-045 |
| project | SDP |
| type | Change |
| CardState | in-progress |
| PlanId | PLAN-SDP-0013 |
| created | 2026-09-30 |
| source | Owner CLI output proposal and Bash prototype clarification |

## Intent

Run gh sdp tree without a Bash filter. Default output is human-readable; --json
selects unchanged machine data. Go renderers dispatch by result type/operation,
with raw JSON fallback for unknown results. json-tree.sh is an example, not a
portable implementation or required runtime dependency.

## Execution

[Plan and design](../../05--Implementation/SDPTool/Output/Plan.md) owns scope,
consumer compatibility, milestones and evidence. No live XFMD application edits,
installation, merge or release. The public default-output change requires an
updated release proposal and explicit machine-client migration.

## Worklog

2026-09-30: owner-authorized execution activated after inspecting Go CLI,
installation output, bootstrap probe and compiled consumer tests.

## Session tracking

[Session 0002](../../Sessions/session-%230002--SDPTool_output.md) records the goal, roadmap, turns and handoff.
It was registered after implementation began, following the owner correction;
retrospective entries are explicitly labeled.
