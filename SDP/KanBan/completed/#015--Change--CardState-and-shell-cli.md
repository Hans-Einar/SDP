# CardState and scoped shell commands

| Field | Value |
| --- | --- |
| id | KB-SDP-015 |
| project | SDP |
| type | Change |
| created | 2026-09-24T17:49:15Z |
| source | owner-conversation-2026-09-24 |
| CardState | completed |
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

[Plan](../../Maintenance/K5/Plan.md): M1 metadata/workflow, M2 CLI/installation.
2026-09-24T17:49:15Z: EVT-KB-SDP-000043 registered the authorized scope. Validation and
installation evidence will be recorded at their actual completion.

## Outcome

2026-09-24T17:51:13Z: K5-M1/M2 complete, EVT-KB-SDP-000058. Installed read-only kanban
and install-cli in /home/warloc/bin. Seven behavior tests passed; no compiler
or daemon is involved. [Evidence](../../Maintenance/K5/Evidence.md).

## Next selection

R3 resumes KB-SDP-001 and KB-SDP-010, already ready in active, because English
documentation and work-state tracking are now delivered. Finish their bounded
profile/consolidation results before selecting new tooling implementation.

## K5-M3 correction

2026-09-24T17:52:09Z: EVT-KB-SDP-000059. Initial K5-M2 board replay found a duplicate
ID but the shell sequence continued to commit/push. The closing event is now
000058; original bytes remain in Git c510fb7 and the K5 evidence capture.
This is an explicit correction to an unintegrated phase record, not rewritten Git
history. New allocation uses max ID + 1. See the final K5 checks.

## K6 — terminal links

2026-09-24T18:20:34Z: EVT-KB-SDP-000066. Owner requests clickable OSC 8 card paths.
One milestone: percent-encoded absolute file URLs on terminals, plain redirected
output, real PTY verification and installation. Opening follows the terminal’s
configured file handler; kanban itself remains read-only.

2026-09-24T18:22:08Z: EVT-KB-SDP-000067. K6-M1 complete; [scope and evidence](../../Maintenance/K6/Plan.md).
