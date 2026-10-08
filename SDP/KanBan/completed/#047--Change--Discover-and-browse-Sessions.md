# Discover and browse Sessions

| Field | Value |
| --- | --- |
| id | KB-SDP-047 |
| project | SDP |
| type | Change |
| CardState | completed |
| PlanId | PLAN-SDP-0015 |
| created | 2026-10-01 |
| source | Owner request to discover Sessions for browsing |

## Outcome

Expose SDP/Sessions as a discoverable capability and navigation group with typed
open targets, reusing source-owned directory discovery. Preserve the generic
Files tree, existing SDL/SDUI/KanBan and refresh/containment behavior.

[Plan](../../05--Implementation/SDPTool/Sessions/Plan.md) owns acceptance/evidence.
[Session 0004](../../Sessions/session-%230004--Session_navigation.md) records the
route. Completed KB046/Session0003 are predecessors, not reopened work.
KB042/043 retain capture, Session semantics and generated timelines.

## Worklog

2026-10-01: analyzed Files-only visibility; activated the authorized bounded
implementation. No release or native XFMD changes are included.

2026-10-01 SN1-M1: implemented and verified capability, directory-derived Sessions
tab, shared open targets and content revisions. Full Go suite, focused race tests,
compiled CLI and process validators pass; plan records scope and evidence.
Completed implementation is unreleased; no XFMD application changes were made.
