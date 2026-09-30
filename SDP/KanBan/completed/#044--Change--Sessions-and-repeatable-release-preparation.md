# Sessions distribution and repeatable release preparation

| Field | Value |
| --- | --- |
| id | KB-SDP-044 |
| project | SDP |
| type | Change |
| CardState | completed |
| PlanId | PLAN-SDP-0014 |
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

## Historical preparation review handoff

MAINT-SDP-0011 preparation is completed and independently reviewed. See its
[checklist](../../Maintenance/RL1/ReleaseChecklist.md) and
[evidence](../../Maintenance/RL1/Evidence.md). Requested decision: review the completed preparation package. The additive
1.1.0 publication proposal is superseded by the output compatibility change in
[KB045](../completed/%23045--Change--Portable-SDPTool-presentation.md); its plan recommends
2.0.0 and requires explicit JSON machine clients before a new release candidate.
Integration/publication remains separately selected.
The latest published product is still 1.0.0. This gate-review awaits that concrete
owner disposition; it does not imply production gates or publication already passed.
XFMD remains unchanged and its upgrade remains manual.

2026-09-30 follow-up: Session 0002 records the owner-selected output change.
Before publication, update gh-sdp's bootstrap dependency and verify native XFMD
machine callers use --json. Previous preparation evidence remains historical;
it is not evidence that these new consumer gates have passed.

2026-10-01: owner selects implementation and publication in Session 0003.
PLAN-SDP-0014 executes the new combined candidate; MAINT-SDP-0011 remains completed.

## Delivery — Session 0003

Completed under PLAN-SDP-0014. SDP 2.0.0 and gh-sdp 0.2.0 are published,
independently reviewed and verified through actual public downloads/default.
Discovery derives inventory and navigation from source files/directories without
registration; invalid sources remain visible. Sessions guides/templates and
release logs/checklist are distributed. Signed predecessor upgrades preserve
project documents and inert navigation.json; new installs create no registry.

[Evidence](../../05--Implementation/SDPTool/Discovery/Evidence.md),
[publication identities](../../05--Implementation/SDPTool/Discovery/Publication.json),
[manual upgrade](../../05--Implementation/SDPTool/Discovery/Manual-upgrade.md) and
[Session 0003](../../Sessions/session-%230003--SDP_discovery_and_release.md).
Native XFMD schema adoption/watching remains its own workstream. Broader automatic
Session capture stays with KB042; neither is silently claimed implemented here.
Earlier present-tense entries above record the decision/preparation history.
