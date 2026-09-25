# K3-M1: traceable card merging and splitting

| Field | Value |
| --- | --- |
| id | KB-SDP-009 |
| CardState | completed |
| ScrumId | SCRUM-SDP-0001 |
| project | SDP |
| type | Change |
| created | 2026-09-23T23:57:37Z |
| source | owner-conversation-2026-09-24; XFMD conversation board |
| owner | Codex |

## Need and basis

The owner wants backlog cards to merge before activation or split into multiple work cards, with ledger traceability. New knowledge must clarify/replace older proposals without losing information or remaining work.

XFMD's local KanBan rules propose a new consolidated card, source table, bidirectional links, superseded for fully replaced sources and reviewed for sources retaining work. They also mention splitting. Reuse this workflow without requiring XFMD to adopt SDP. XFMD files were uncommitted when inspected; none are modified here.

Inspected source: `/home/warloc/git/xfmd-sdl-navigation/Agents/KanBan/README.md`. Registration: 2026-09-23T23:57:37Z. SHA-256: `eb7488d631dc8c9c934a71da15fe3e3f52239989998596648e64154a8ff2c7c8`. This records local provenance, not a published XFMD contract.

## Implementation plan K3

One phase branch `sdp/phase-k3-card-lineage` from K2 (`321e193`); one milestone K3-M1:

1. Describe full/partial merge and split, source/target roles, preservation of questions and activation with bounded, authorized scope.
2. Extend ledger payload with versioned, typed lineage without rewriting old events; reuse created/moved/reviewed and stable card IDs.
3. Update card template and graph idea; create isolated, explicitly fictional examples and check positive/negative cases, replay and existing boards.

This delivers workflow/data contracts, not automatic file movement, a graph or sdptool implementation. Do not merge existing product proposals without substantive review. XFMD sources are read only.

## Acceptance and evidence

**K3-M1 delivered on 2026-09-24.** Lineage contract, payload schema 0.2, updated card template and graph clarification are delivered. Payload schema 0.1 and historical events remain unchanged. Merge/split creates new cards; partial transfer retains explicit remaining work. No production cards were merged or split.

`python3 SDP/Agents/KanBan/examples/verify_lineage.py` passed: fictional full/partial merge/split histories, historical payload 0.1 and 15 negative cases. Tests reject missing participants/links, inconsistent definitions, reused target IDs, overlapping sets and invalid remaining-work handling.

A separate one-off production-board check passed: three boards, 21 status directories, 13 cards, four Refs and 112 local Markdown links. Schema/replay and locations also passed after closure; ledgers contain 17 events. Old ledger lines were verified as byte-preserved prefixes against the parent commit; SDL/SDUI ledgers and payload 0.1 are unchanged. `git diff --check` passed.

This is a manual process and tested data contract. No general move command, KanBan graph or cross-repository transaction engine was delivered. XFMD's proposal was used as local, uncommitted evidence, not modified or published.
