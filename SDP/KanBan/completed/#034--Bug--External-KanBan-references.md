# KB-SDP-034 — Keep valid boards available with external Ref cards

| Field | Value |
| --- | --- |
| id | KB-SDP-034 |
| project | SDP |
| type | Bug |
| CardState | completed |
| PlanId | PLAN-SDP-0004 |
| Systems | SDPTOOL |
| created | 2026-09-27T08:22:06.815649Z |
| source | RP1 live XFMD upgrade verification, 2026-09-27 |
| next_review | Before the next SDPTool/XFMD navigation delivery |

## Observed behavior and impact

Released SDPTool 0.2.0 `gh sdp tree` reports the entire XFMD KanBan tab unavailable:
`unresolved Ref KB-SDP-014`. XFMD's KB-XFMD-012 legitimately points to that card in
SDP's separate board. Project discovery and installation succeed; the board itself
passes schema/history/CardState validation. The present BoardNodes implementation
resolves primary only among local card IDs and fails the whole board if absent.
This is a navigation limitation, not an installation failure. The current contract
promises local Ref resolution; extend its behavior explicitly for cross-project Refs.

## Selected bounded outcome

Keep valid local cards selectable while representing an external/unavailable primary
honestly. Preserve strict rejection of genuinely broken local references and invalid
history. Do not erase primary metadata, import another project's cards, follow arbitrary
external paths automatically, or treat the foreign ID as a local alias. ER1 selects the
consumer representation and tests, including real XFMD input.

## Evidence and scope

[RP1](../../Maintenance/RP1/Plan.md) records the installed signed release and
[node-state output](../../Maintenance/RP1/evidence/live-tree.json).
The initial registration selected no implementation. The subsequent owner decision
below authorizes ER1; patch release and XFMD application changes remain separate.

## Worklog

2026-09-27: EVT-KB-SDP-000199. Registered the live navigation finding; no implementation selected.

## Owner decision and execution

2026-09-27: External references are permitted without existence guarantees. Do not
resolve them through local checkouts, remote lookup or an external configuration.
[ER1](../../05--Implementation/SDPTool/ExternalReferences/Plan.md) implements explicit
external metadata while keeping local references strict and local cards openable.

2026-09-27: EVT-KB-SDP-000200. Selected ER1; backlog → active/in-progress.

## Delivery

2026-09-27: EVT-KB-SDP-000201. ER1-M1 complete; active/in-progress → completed.
The optional externalReference field labels foreign primary IDs without lookup
or existence guarantees. Local references and history remain strict. Eleven
regression cases, the full race suite and vet pass; the built CLI returns all 18
XFMD cards with a validated KanBan tab. See the linked ER1 plan for evidence.
The installed 0.2.0 engine still requires a subsequent release to receive this fix.
