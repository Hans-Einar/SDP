# SDUI as a subproject and library for SDP tools

| Field | Value |
| --- | --- |
| id | KB-SDUI-001 |
| CardState | completed |
| ScrumId | SCRUM-SDP-0001 |
| project | SDUI |
| type | Ref |
| created | 2026-09-23T18:25:13Z |
| source | owner-conversation-2026-09-23 |
| next_review | 2026-09-30 |
| primary | KB-SDP-017 |

Registered from the owner conversation on 2026-09-23. The timestamp records registration, not a reconstructed discussion time. The directory and ledger record lifecycle status.

Primary card: [KB-SDP-017 — sdptool and project navigation](../active/%23017--Proposal--sdptool-and-project-navigation.md).

## Local impact

SDUI should have its own SDP area and KanBan, and work as a library/subproject with the same structure in a monorepo or later separate repository. `sdptool` must support explicit selection of the SDUI project. [KB-SDP-001](../completed/%23001--Proposal--Project-structure.md) owns the structural migration.

Preserve the current Go/Fyne/SVG direction. This does not request a FOX host, another parser, a KanBan GUI or full Markdown/Mermaid support. The [KB-SDP-003 graph idea](../backlog/%23003--Idea--KanBan-graph.md) does not select SDUI as its host.

## Handling

Resolve local project configuration, model/plan sources and library boundaries before migration. Keep the code plan's G phases distinct from this ideas backlog.

[SDUI implementation plan](../../../SDUI/docs/implementation-plan.md)

## Navigation and consolidation review — 2026-09-25

The owner requests an SDUI subtab inside native XFMD’s SDP tab. Define its initial source/documentation/preview actions explicitly with the host; preserve SDUI-owned parsing, layout and export. This entry point does not authorize a FOX widget runtime or change the Go/Fyne/SVG architecture.

The previous primary KB-SDP-002 and discovery card KB-SDP-016 were fully merged
into KB-SDP-017. This Ref keeps its identity and backlog state; only its primary
and filename change. Recorded 2026-09-25T01:41:26Z, Codex, EVT-KB-SDUI-000004.

## Scrum-0001 routing outcome

Local consequence is recorded in [KB-SDP-017](../active/%23017--Proposal--sdptool-and-project-navigation.md). The former separate-board
routing is no longer needed. This closes only this Ref, not language, tool or
model implementation. Historical identity and original primary remain recorded.
