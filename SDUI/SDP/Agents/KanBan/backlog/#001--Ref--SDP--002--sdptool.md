# SDUI as a subproject and library for SDP tools

| Field | Value |
| --- | --- |
| id | KB-SDUI-001 |
| project | SDUI |
| type | Ref |
| created | 2026-09-23T18:25:13Z |
| source | owner-conversation-2026-09-23 |
| next_review | 2026-09-30 |
| primary | KB-SDP-002 |

Registered from the owner conversation on 2026-09-23. The timestamp records registration, not a reconstructed discussion time. The directory and ledger record lifecycle status.

Primary card: [KB-SDP-002 — sdptool](../../../../../SDP/Agents/KanBan/backlog/%23002--Proposal--sdptool.md).

## Local impact

SDUI should have its own SDP area and KanBan, and work as a library/subproject with the same structure in a monorepo or later separate repository. `sdptool` must support explicit selection of the SDUI project. [KB-SDP-001](../../../../../SDP/Agents/KanBan/active/%23001--Proposal--Project-structure.md) owns the structural migration.

Preserve the current Go/Fyne/SVG direction. This does not request a FOX host, another parser, a KanBan GUI or full Markdown/Mermaid support. The [KB-SDP-003 graph idea](../../../../../SDP/Agents/KanBan/backlog/%23003--Idea--KanBan-graph.md) does not select SDUI as its host.

## Handling

Resolve local project configuration, model/plan sources and library boundaries before migration. Keep the code plan's G phases distinct from this ideas backlog.

[SDUI implementation plan](../../../../docs/implementation-plan.md)
