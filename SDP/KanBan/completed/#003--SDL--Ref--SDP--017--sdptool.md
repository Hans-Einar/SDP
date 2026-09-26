# SDP tools as a consumer of SDL

| Field | Value |
| --- | --- |
| id | KB-SDL-003 |
| CardState | completed |
| ScrumId | SCRUM-SDP-0001 |
| project | SDL |
| type | Ref |
| created | 2026-09-23T18:25:13Z |
| source | owner-conversation-2026-09-23 |
| next_review | 2026-09-30 |
| primary | KB-SDP-017 |

Registered from the owner conversation on 2026-09-23. The timestamp records registration, not a reconstructed discussion time. The directory and ledger record lifecycle status.

Primary card: [KB-SDP-017 — sdptool and project navigation](../active/%23017--Proposal--sdptool-and-project-navigation.md).

## Local impact

SDP process commands should consume SDL as a library for parsing/validation and model facts. Resolve language-package versus process/viewpoint responsibilities before repository extraction. Do not create a competing parser or move Go modules now.

## Handling

Map existing APIs/CLIs and the export boundaries needed for plan proposals, validation and viewing. Record concrete deliveries/dependencies in the approved plan.

## Navigation and consolidation review — 2026-09-25

SDL tooling owns the ViewPoint catalog, typed model collections and selected projections. The native XFMD SDL subtab needs a source-derived tree, such as VP01 → UseCase → actual use cases. Inventory existing exports and define the smallest adapter needed; XFMD and sdptool must not duplicate parsing or infer model identity from rendered Markdown.

The previous primary KB-SDP-002 and discovery card KB-SDP-016 were fully merged
into KB-SDP-017. This Ref keeps its identity and backlog state; only its primary
and filename change. Recorded 2026-09-25T01:41:26Z, Codex, EVT-KB-SDL-000012.

## Scrum-0001 routing outcome

Local consequence is recorded in [KB-SDP-017](../active/%23017--Proposal--sdptool-and-project-navigation.md). The former separate-board
routing is no longer needed. This closes only this Ref, not language, tool or
model implementation. Historical identity and original primary remain recorded.
