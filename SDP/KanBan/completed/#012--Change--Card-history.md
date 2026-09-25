# K4-M1: living cards, revision logs and Git diff

| Field | Value |
| --- | --- |
| id | KB-SDP-012 |
| CardState | completed |
| ScrumId | SCRUM-SDP-0001 |
| project | SDP |
| type | Change |
| created | 2026-09-24T14:52:35Z |
| source | owner-conversation-2026-09-24 |
| owner | Codex |

## Request and scope

The owner wants cards to document work while active, with revisions/diffs throughout their lifecycle. Reuse existing Git. Keep Maintenance for phase plans, inventories and evidence linked back to active cards. This introduces no separate Git engine, automatic snapshots or viewer.

## Plan K4

Branch `sdp/phase-k4-card-history` from R1 `f722dc2`. One milestone K4-M1: workflow/template, visible active-card log, reproducible history/diff and acceptance criteria for future tools/graphs. KB-SDP-010 remains the next documentation assignment; registration does not implement other backlog proposals.

## Results and evidence

[History](../History.md) defines the manual workflow. The [card template](../Card-template.md) includes a worklog; KB-SDP-001 summarizes delivered/remaining work. KB-SDP-002/003 retain future history/diff needs. SDL/SDUI documents use the same shared workflow. No payload schema or old ledger event changed.

Git examples were run against KB-SDP-001: `--follow` shows K1, K2 and R1; exact blob diff from `bb3728c` backlog to `f722dc2` active shows 32 additions / 11 deletions. Looking up eventId 000015 identifies `d269bc7`. K3 lineage tests pass for full/partial merge/split, old payload and 15 negative cases. Schema/replay and locations pass for 16 cards / 26 events across three boards. Ledger prefixes are byte-preserved against R1 `f722dc2`; R1 checks pass for 1900 local Markdown file targets, 574 generated artifacts and migration boundaries. `git diff --check` passes. This verifies the documentation workflow, not integrated history viewing.

## Worklog and revisions

| Time | Actor / event | Work and outcome | Remaining work |
| --- | --- | --- | --- |
| 2026-09-24T14:52:35Z | Codex; EVT-KB-SDP-000018 | Registers and bounds K4 from the owner's request. | Manual history and active worklog. |
| 2026-09-24T14:52:35Z | Codex; EVT-KB-SDP-000022 | K4-M1 delivered as a documented, verified workflow. | KB-SDP-002/003 still own tools/graph work. |

The card has one committed delivery revision. Created/moved within one milestone are two handling events, not separate Git snapshots of its content.
