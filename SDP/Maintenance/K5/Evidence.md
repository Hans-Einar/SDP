# K5 verification

## K5-M1

Added CardState to the 17 previously committed project cards and the new
KB-SDP-015, with ready for KB-SDP-001/010, in-progress for K5, and states matching
other lifecycle directories. The concurrently registered, uncommitted KB-SDP-014
is preserved outside this milestone. Updated the shared template, workflow and
root agent instructions; queue/review explanations live in cards, not another store.

Three-board schema/replay/physical placement and historical prefixes checked.
Current document file/fragment links and git diff --check pass. State-only changes
append reviewed events; historical event bytes and language profiles are unchanged.
K5-M2 supplies the executable scripts and installation tests.

## K5-M2

Seven CLI behavior tests pass: scoped/aliased/ordered listings and quoted paths;
missing cards/body examples; invalid/duplicate/empty states; body-table isolation;
arguments; full installation, backup, idempotence and installed execution; symlink
destination preservation. bash -n passes. Both shell scripts are installed in
/home/warloc/bin and kanban resolves through the existing PATH. Actual board-root
and active-directory listings match scope. No shell startup file was changed.

The owner subsequently authorized including concurrent KB-SDP-014. Its uncommitted
event was assigned the next free ID and appended after the committed prefix,
with provenance recorded in the card/event. No committed ledger bytes changed.
It gains CardState=backlog; its future version/distribution proposal is not implemented.

All 19 cards on the three boards now have one CardState consistent with their
folder. K5 is completed; KB-SDP-001/010 are ready for R3. Current file/fragment links,
board replay and historical evidence preservation pass.

## K5-M3 — explicit correction of an unintegrated event collision

The M2 board check actually failed: the temporary event writer used line count,
which reused 000057 after the concurrent registration filled a reserved ID gap.
The shell sequence did not stop on that failure and committed/pushed c510fb7.
The earlier M2 blanket board-pass statement is corrected here. CLI tests and
file/fragment checks did pass independently.

Only the new K5 closing event ID changes from 000057 to 000058; its prior chain,
time and payload remain. The original two colliding records are retained in
duplicate-event-original.ndjson and Git c510fb7; no Git history is rewritten.
A reviewed event records the correction. Allocation now uses maximum ID + 1,
and check/commit command groups stop on failure. This narrow unintegrated-phase
repair is explicit; pre-K5 historical records remain byte-identical.

Final K5-M3 checks pass: 19 cards/70 events across three boards, schema/replay/
placement and R2 ledger prefixes; 105 frozen records, 574 generated outputs,
2003 local file links/94 fragments; all seven CLI tests; git diff --check.
