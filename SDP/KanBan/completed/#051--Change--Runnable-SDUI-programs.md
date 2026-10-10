# Discover and run connected SDUI programs

| Field | Value |
| --- | --- |
| id | KB-SDP-051 |
| project | SDP |
| type | Change |
| CardState | completed |
| created | 2026-10-09T10:08:18.980176+00:00 |
| source | Owner Session0010 T010 |
| Systems | SDPTOOL, SDUI |
| PlanId | PLAN-SDP-0023 |

## Selected outcome

SDPTool owns runnable application discovery and launch. Keep existing SDUI source
inventory and add explicit connected program declarations; widget lab is the first
consumer. XFMD can later consume the same contract. The owner's correction moves
ownership from a proposed XFMD-specific launcher into SDPTool.

[Plan and boundary](../../05--Implementation/SDPTool/Programs/Plan.md).
Predecessor KB-SDUI-006 remains completed. KB-SDUI-005 owns the separate generic
preview workflow; do not mix its concurrent changes into this implementation.

## Worklog

- 2026-10-09T10:08:18.980176+00:00: EVT-KB-SDP-000318; created active/in-progress under explicit owner request. RSP1 implementation and verification pending; no release or XFMD change selected.

## Outcome

Implemented on sdp/runnable-programs at 99631595eb2ccf5d1e8cc52be5ce30a02b492341; full race suite passed.
Consumer 75b6dca declares widget-lab; actual gh-sdp discovery and native Run passed,
including a visible SDL/Go return and normal exit 0. [Evidence](../../05--Implementation/SDPTool/Programs/Evidence.md).
The local development descriptor enables testing; published defaults and XFMD UI
are separate. No selected RSP1 implementation remains.

- 2026-10-09T10:18:54.830417+00:00: EVT-KB-SDP-000319; completed both RSP1 milestones and linked verification.
