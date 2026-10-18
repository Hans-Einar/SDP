# WCI3 IME boundary investigation — in progress

Not product acceptance. Actual configured IBus 1.5.32 Simple engine through XIM,
X11/Xvfb :189 and pinned Fyne 2.8.1/GLFW, isolated D-Bus and XDG directories.
Official AlmaLinux RPMs were extracted only under /tmp/sdui-ime-probe; no system
installation or desktop configuration change. IBUS_ENABLE_CTRL_SHIFT_U=1 enables
Simple-engine Unicode composition. Initially disabled hex mode and a launcher
readiness bug accepting `(null)` were corrected before the recorded reproduction.

XTest Ctrl+Shift+U, 4 e 2 d produces no Entry.OnChanged during composition.
Return then produces Entry.TypedKey(Return), OnSubmitted("") and only afterward
OnChanged("中"). Escape during a second composition reaches Entry.TypedKey(Escape).
Thus the current adapter must not claim consumed IME keys cannot submit/revert.
No SDUI/SDL application action was attached; this is an upstream boundary probe.

Pinned GLFW x11_window.c calls XFilterEvent, but delivers the key callback before
checking its filtered result; it only guards character lookup. A bounded isolated
trial now guards the key callback/timestamp update too. No repository dependency
or product source has been changed; trial and independent review remain pending.

Primary source context: https://github.com/ibus/ibus/blob/main/src/ibusenginesimple.c
and https://github.com/ibus/ibus/blob/main/bus/ibus-daemon.1.in . The checked-out
pinned GLFW source is the actual dependency used, not an inferred latest version.
No WM, Wayland, external candidate-panel or CJK engine coverage is claimed.

## Bounded filter trial result

The one-condition trial patch guards BOTH key delivery and its duplicate-event
timestamp update with !filtered. The same real IBus input now produces one
OnChanged("中") and no Submit on composition Return; composition Escape causes
neither TypedKey(Escape) nor Submit. A subsequent ordinary Return submits "中"
once. Ordinary a, Left and Backspace still arrive once and edit correctly.
Logs, patch and binary/source hashes are retained here. This is a successful
mechanism experiment, not SDUI acceptance or a general GLFW compatibility claim.

Candidate architectural disposition: retain the exact pinned GLFW source with
this tiny X11-only patch as a documented local dependency, preserving licenses,
upstream identity and a reproducible patch. Apply explicit module replacements
to each maintained build root; WCI4 must verify private helper packaging includes
it. This avoids timers/heuristics in Entry and makes the actual consumed-key
boundary authoritative. Independent review and WCI3-stage selection still precede
product dependency edits. Other platforms retain their pinned upstream source.
Generated-code consumers outside maintained build roots must explicitly select
this dependency policy; import alone cannot impose a transitive Go replacement.
