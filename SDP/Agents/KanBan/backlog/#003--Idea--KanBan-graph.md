# KanBan graph with a time axis and progressive detail

| Field | Value |
| --- | --- |
| id | KB-SDP-003 |
| project | SDP |
| type | Idea |
| created | 2026-09-23T18:25:13Z |
| source | owner-conversation-2026-09-23 |
| next_review | 2026-09-30 |

Registered from the owner conversation on 2026-09-23. The timestamp records registration, not a reconstructed discussion time. The directory and ledger record lifecycle status.

## Owner's idea

Display each status directory as a vertical column. Time runs downwards, initially one row per week. Each document receives a point when created or moved; a continuous path through columns shows its lifecycle. Clicking a week expands its days; further clicks focus on hours with activity. Points should open the document and event explanation.

## Basis and scope

Read the append-only KanBan ledger, not file mtime or assumed Git commit times. Preserve real timestamps and deterministic ordering for simultaneous events. Distinguish events from current location and show reopenings/backward moves. Display Ref cards as references to primary cards, not additional implemented deliveries. Project, type/tag and time filters can improve readability.

Interactive viewing, time zoom and host selection remain open. Markdown/SVG can provide static export, but static Markdown alone cannot provide clickable time zoom. Do not build a new renderer or select Fyne/XFMD/SDUI before agreeing on a bounded trial.

## Possible first trial

Generate a static weekly view from a small validated ledger with moves and reopenings. Check status, paths, timestamps and links. Then decide the interactive host, timezone/week boundaries, dense-event handling and hidden inactive weeks. The KanBan foundation implements no graph.

## Clarification, 2026-09-24: convergence and branching

Show merge (multiple sources → one target) and split (one source → multiple targets) using typed ledger lineage, not guesses from prose or ordinary links. Partial transfer shows remaining work in the source. OperationId connects participant events; mark incomplete operations or inaccessible boards as unknown/incomplete. Preserve source history and timestamps. [Contract](../Lineage.md).

## Clarification, 2026-09-24: content revisions

Also show significant reviewed events while a card remains active. A point should open the card content as it was at that event and a diff against a selected revision. Label working drafts separately from committed history. Git owns text versions; the ledger owns lifecycle events. [K4 workflow](../History.md) explains their connection and limits. History/diff requires available Git objects.

## Worklog and revisions

| Time | Actor / event | Handling | Remaining work |
| --- | --- | --- | --- |
| 2026-09-24T14:52:35Z | Codex; EVT-KB-SDP-000020 | Adds active content revisions and historical diff to the graph idea; graph remains unimplemented. | Agree on a vertical tool delivery before implementation. |
