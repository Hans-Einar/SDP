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
