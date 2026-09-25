# KanBan and SDUI navigation services

| Field | Value |
| --- | --- |
| id | KB-SDP-025 |
| project | SDP |
| type | Change |
| CardState | completed |
| ScrumId | SCRUM-SDP-0002 |
| SprintId | SPR-SDP-0001 |
| Systems | SDPTOOL |
| created | 2026-09-25T12:38:14Z |
| source | KB-SDP-017; owner-conversation-2026-09-25 |
| next_review | Sprint start, then each owning phase milestone |

## Delivery and acceptance — T3-M3

Expose KanBan lifecycle folders, CardState, stable card identity/current paths,
Refs and optional SprintId/ScrumId using the supported board contract. Read the
shared board descriptor and management ledger without counting frozen history
copies twice. Define refresh after moves/edits, unreadable boards and unsupported
versions. Keep lifecycle, work state, Git revisions and implementation evidence
separate. Do not infer feature completion from a completed card.

Expose the SDUI source inventory and only explicitly verified existing Go preview
or documentation operations. Publish supported, unavailable and unsupported
capabilities honestly; no new SDUI runtime, widgets or XFMD rendering is included.

Depends on KB-SDP-022/023; align IDs/revisions with KB-SDP-024. Acceptance: fixtures
cover moved cards, Ref targets, sprint grouping, optional/missing boards and SDUI,
malformed metadata and revision changes. Inventory and dispatch use one host
registration contract. KB-SDP-003 retains the history timeline; KB-SDP-014 retains
reusable contract/distribution work.

## Source and scope transfer

| Source | Transferred scope | Retained scope |
| --- | --- | --- |
| [KB-SDP-017](../active/%23017--Proposal--sdptool-and-project-navigation.md) | T3-M3: KanBan and SDUI navigation services | P0-M2, T4-M2 and T5 remain in the source; sibling cards own the other selected milestones |

Typed partial split: KBO-SDP-000004. The [single feature plan](../../05--Implementation/SDPTool.md)
owns phase/milestone acceptance; the [Sprint](../../Sprints/Sprint--%230001--SDPTool-preview-and-navigation.md) owns
membership and completion. No implementation is delivered by this registration.

## Worklog

2026-09-25T12:38:14Z — EVT-KB-SDP-000128: registered with acceptance and dependency boundaries. Await Sprint start.

2026-09-25T13:08:53Z — EVT-KB-SDP-000138: Owner selects Sprint execution before the skills and installer Scrums; backlog to active/ready.

2026-09-25T13:30:23Z — EVT-KB-SDP-000148: T3-M3: add KanBan status/card and explicit SDUI preview inventory to the common navigation contract.

2026-09-25T13:35:18Z — EVT-KB-SDP-000149: T3-M3 delivered: board/Ref/work-state inventory and explicit SDUI structural Markdown preview with source/entry targets; see SDPTOOL-VER-T3-M3.
