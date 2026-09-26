# K6 — clickable terminal card paths

Owner request, 2026-09-24. One milestone, K6-M1, on
sdp/phase-k6-terminal-links, stacked on R3 84089ae.

Deliver OSC 8 links in kanban status/state when stdout is a terminal and TERM
is not dumb. Use absolute, byte-wise percent-encoded file URLs; retain safe
labels and plain redirected output. Reuse existing terminal file-handler settings.
No viewer launch or Kitty configuration changes are part of the command.

## K6-M1 evidence

Eight CLI tests pass, including an actual pseudo-terminal assertion of OSC 8
opening/closing sequences and exact file URIs containing spaces, #, UTF-8,
percent, question mark, brackets, tab and escape bytes. Filename controls do
not leak into terminal output. Captured pipe output has no escapes/file URLs.
Existing scope, errors, metadata, installer backup and idempotence checks pass.
bash -n and git diff --check pass. The installer refreshed /home/warloc/bin/kanban,
preserving its previous version as a backup.

Read-only inspection confirms Kitty's existing file/Markdown rule selects the
executable /home/warloc/.local/bin/xfmd. No actual mouse click or native viewer
launch is claimed as tested here. Card/ledger replay and document links are checked
before the milestone commit. K6-M1 completes this bounded request.
