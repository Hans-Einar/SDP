# WCI2-M2 native pilot — evolving candidate

Initial native binary SHA-256
`1b213abfb32a825c7d86a844111fb1e958e4fd2bb8d2ac941b7af454242077ed`
ran on separate Xvfb :190 at 1600x1200. This is development feedback, not final
candidate acceptance.

The commands variant passed 13 checks: zero preflight actions, one shared Run
through toolbar/key/menu, typed Flag state, atomic exclusivity, a nested pointer
selection, disabled toolbar/key/menu Enter refusal, re-enable recovery and clean
teardown. An OS capture showed the toolbar, tree and distinct visible receivers.

The context variant proved right-click Beta preserves selected Alpha and passes
Beta through real SDL. Replacing the collection closed its stale menu. A following
Ctrl+R from the focused tree was dropped: pinned GLFW dispatches shortcuts to a
focused Shortcutable without falling through to canvas handlers. Collection's
existing adapter only handles Alt arrows. Host correction must forward otherwise
unhandled declared commands while retaining native editing/navigation precedence;
the actual-input regression remains in the harness. No context retargeting failure
was observed. Complete dialog/native lifecycle evidence remains pending.

## Second evolving candidate

Binary `3b5efbcc0ca1b75b1c0a205dd1822db4ffeadf6f9c16662d53efcd6e4f7c286c`
corrects focused-control shortcut forwarding and scoped surface preparation.
Context passed four checks, modal/nonmodal dialogs seven each, acceptance errors
seventeen, nonmodal lifecycle twelve, and modal/nonmodal focus ten each.
Real SDL receipts preserve unknown/succeeded domain outcomes, prevent replay after
an ambiguous or post-domain failure, and distinguish Cancel/Close from Accept.

Keyboard Home/Enter, nested End/Space and Escape passed. Shift+F10 did not open
the focused collection context menu; pinned GLFW deliberately excludes Shift-only
CustomShortcut dispatch. This remains a host correction before final acceptance.

Initial modal capture at 0.2 seconds preceded completed desktop painting despite
accepted input. A later capture after editing and two seconds shows the actual
modal and draft. The harness now adds capture-only paint time, without changing
interaction assertions. A nonmodal parent click initially hit the overlapping
undecorated child window; explicit XMoveWindow test arrangement makes both windows
reachable. Parent action and restoration of actual keyboard focus then pass.
These are declared tooling corrections, not product changes or final evidence.

Native inspection identified a small nonmodal content clip discrepancy; host is
correcting initial chrome sizing. Run icon visibility and final-candidate native
proof remain pending. No WCI2-M2 delivery is claimed from these pilot results.
