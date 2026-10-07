# Native nonmodal adapter selection probe

Session0010 T002, WCI2 draft architecture selection. This is a standalone Fyne
experiment with synthetic local field commits, not SDUI/SDL implementation or
WCI2 acceptance. The selected adapter is app.NewWindow with explicit owner-managed
parent lifetime. The frozen source/README.md and source/manifest.json record the
preceding compile-only experiment; this run record adds actual native observations.

## Environment and candidate

Go 1.27.1, pinned Fyne v2.8.1, Xvfb :189 (1000x700x24, no TCP, no reset), XTest
input through the WCI1 x11_input.py tool. Separate XDG directories under
/tmp/wci2-surface-probe/xdg. No window manager was available in this display.
Binary SHA-256: a5694a15611793b7ed7af819bb1c176543f11e9b32935ca70ea1ac1214026fc0.
Sources and dependency hashes are retained in source/manifest.json. To rebuild,
run go build -mod=readonly -o surface-probe . from source/. No binary is committed.

## Actual input and observations

All field edits, button invocations and keyboard operations below used XTest.
The probe did not receive a direct callback or control-channel invocation.

1. Parent: click field (100,60), End/a; click Open (330,100). Child 1: End/b.
   Click exposed parent action (640,140). Log confirms child still open and dirty,
   parent draft length 16 and unchanged child draft fingerprint.
2. Child 1: Enter/b/Escape/Escape. Exactly one local field commit, then revert
   to its accepted fingerprint, then one Cancel result. Parent Space reopened
   child 2 through the opener button, establishing actual keyboard focus restoration.
3. Child 2: End/b; parent Close child (640,180). One Close result with
   parent-command reason and no additional opener-focus request.
4. Reopen child 3; End/c; parent Hide (640,260). Exactly one parent-hidden Close,
   followed by parent shown without child reopening. Reopen child 4; End/d;
   parent Simulate reload (640,220). One reload Close and parent generation advance.
5. Reopen child 5; End/e; parent Close parent (330,300). One parent-closed Close,
   no orphan child, and process exited successfully.

The log assertions passed: five opening IDs each have exactly one terminal result;
there is one local field commit and one opener-focus request; the parent action
observes its live dirty child; opening 5 belongs to parent generation 2. OS captures
were inspected: parent controls, child field and buttons render, and the parent
remains available behind the nonmodal child. native.stderr is empty.

## Limits

No window-manager decoration click or WM focus policy was tested. Keyboard focus
restoration is demonstrated only in this X11 environment. Parent Hide/reload are
explicit ownership simulations, not SDUI atomic reload evidence. There are no SDL
calls, model tokens, failed-publication checks, modal behavior, IME or persistence.
WCI2 implementation must repeat its actual contract cases with SDUI runtime state,
real bridge bindings and independent integrated review. This probe supports the
adapter selection; it does not deliver or close a widget family.
