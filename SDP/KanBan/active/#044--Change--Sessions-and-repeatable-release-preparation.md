# Sessions distribution and repeatable release preparation

| Field | Value |
| --- | --- |
| id | KB-SDP-044 |
| project | SDP |
| type | Change |
| CardState | gate-review |
| PlanId | MAINT-SDP-0011 |
| created | 2026-09-30 |
| source | Owner request after Session 0001 closeout |

## Selected outcome

Generate one release log from canonical notes per release, include the manual
Sessions format in installation/upgrade, and require ReleaseChecklist evidence
for payload inventory, predecessor digests, signed descriptor and consumer gates.
Owner will perform the XFMD upgrade manually after publication is ready.
[Plan](../../Maintenance/RL1/Plan.md) owns execution and evidence.
[KB042](../backlog/%23042--Proposal--Goal-oriented-sessions-and-roadmaps.md) owns
broader Session identity/capture integration, which this delivery does not implement.

## Worklog

Registered and activated for the explicitly requested Maintenance. SDL 0.6 and
SDP release numbering are separate; latest product release is 1.0.0, next proposed
additive release is 1.1.0. No merge, publication or XFMD mutation yet.

## Concrete review handoff

MAINT-SDP-0011 preparation is completed and independently reviewed. See its
[checklist](../../Maintenance/RL1/ReleaseChecklist.md) and
[evidence](../../Maintenance/RL1/Evidence.md). Requested decision: select proposed
SDP 1.1.0 integration/publication (including the preceding SDL composition stack).
The latest published product is still 1.0.0. This gate-review awaits that concrete
owner disposition; it does not imply production gates or publication already passed.
XFMD remains unchanged and its upgrade remains manual.
