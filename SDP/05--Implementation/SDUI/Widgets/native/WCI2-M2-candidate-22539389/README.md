# WCI2-M2 pre-review candidate 22539389

This candidate is **not final acceptance**. Binary
2253938916e1a31dfa4ec6d6095daf9f82a7522abcd444833cb64977e69b1c33
passed 110 command/surface native checks and 58 pane regression checks. Exact
raw logs and screenshots are retained. A terminal-result audit found exactly
one result for each of 22 observed published openings.

Independent review then reproduced a missing case: a source-root sibling dialog
opened from a nonmodal surface was ordered by source path length instead of its
actual ParentSurface. It attached to the wrong canvas and used main-window size.
The existing checks remain valid for their cases but do not close this blocker.
The corrected final candidate must add this workflow and rerun relevant proof.

One keyboard run logged Fyne Preferences load EOF; all six behavior checks and
teardown passed. The original warning is retained in commands/keyboard/stderr.txt.
A single same-binary repeat with fresh isolated configuration passed with empty
stderr. Pinned Fyne creates/truncates its preferences file before setting the
savedRecently watcher guard, which is a plausible cause, not a proven causal
trace. No SDUI/domain failure or shared owner configuration was observed. This
upstream observation does not authorize a dependency patch in WCI2.

Native conditions: Go 1.27.1 linux/amd64, pinned Fyne2.8.1, Xvfb :190
1600x1200 (commands) and :189 1000x700 (panes), XTest input and OS captures.
No window manager, decoration-click or Wayland behavior is claimed. Window close
uses the actual WM_DELETE_WINDOW protocol; nonmodal focus testing explicitly
arranges windows with XMoveWindow to make both reachable.
