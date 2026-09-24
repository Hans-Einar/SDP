# CardState and scoped shell commands

| Field | Value |
| --- | --- |
| id | KB-SDP-015 |
| project | SDP |
| type | Change |
| created | 2026-09-24T17:49:15Z |
| source | owner-conversation-2026-09-24 |
| CardState | in-progress |
| owner | Codex |
| next_review | At K5-M2 completion |

## Scope and acceptance

Use one CardState in visible metadata to distinguish ready, in-progress,
gate-review and queued without a separate gate file or queue directory. Keep
folders as lifecycle stages. Queue explanations identify the next bounded work.
Implement Toolkit/scripts/cli/kanban.sh and a CLI installer copying shell scripts
to ~/bin without their .sh suffix. status/state group by state; board root lists
its statuses, a status directory only itself. No matches must fail with exactly
`no KanBan cards found`. Ignore body examples, preserve paths with spaces and
existing installed commands. No new sdptool or graph implementation.

## Milestones and worklog

[Plan](../../../Maintenance/K5/Plan.md): M1 metadata/workflow, M2 CLI/installation.
2026-09-24T17:49:15Z: EVT-KB-SDP-000043 registered the authorized scope. Validation and
installation evidence will be recorded at their actual completion.
