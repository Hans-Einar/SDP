# WCI2 nonmodal surface probe

Temporary Architect experiment, 2026-10-08. **Compiled only; native run pending.**
No product imports, SDL calls, persistence, release work or repository source edits.
Fyne v2.8.1 dependency graph/checksums copied from SDUI/go; standalone module name.
`manifest.json` identifies the inspected sources, probe and build. No headless test
result is presented as native evidence.

## Decision and source evidence

Choose `app.NewWindow` with application-owned parent lifetime for persistent
nonmodal dialogs. It allows independent canvases and keeps the parent available.
It is not an OS parent/transient relationship: Fyne.Window supplies no such API.
Children must not be SetMaster. All operations run on the Fyne UI goroutine.

Inspected in the pinned module:

- app/app.go NewWindow delegates to driver.CreateWindow; window.go exposes
  SetCloseIntercept, SetOnClosed, Canvas, Hide, Close and RequestFocus.
- internal/driver/glfw/window.go Close invokes OnClosed before marking closing;
  processClosed invokes the close interceptor, whereas direct Close bypasses it.
  Hide does not invoke OnClosed. Revoke the owner token before calling Close;
  make the fallback OnClosed notification idempotent and never recursively Close.
- widget/popup.go Show uses an OverlayContainer with Hide as outside dismissal.
  internal/widget/overlay_container.go Tapped/Secondary invoke dismissal;
  internal/driver/util.go FindObjectAtPositionMatching searches only the overlay
  when one exists. Thus ordinary NewPopUp is not persistent parent pass-through:
  a parent click dismisses it instead of operating the parent underneath.
- internal/driver/glfw/window_desktop.go RequestFocus returns without action on
  Wayland; even on X11 it is a request. Canvas.Focused is not an OS focus query.
  Do not raise the parent on parent-driven/background child disposal.

Reject the ordinary popup alternative for this workflow; do not create a custom
overlay/event-routing framework. Native failure would require a bounded design
correction, not make the selected adapter unspecified again.

## Build and native launch

Executed successfully from this directory with Go 1.27.1:

```sh
go build -mod=readonly -o surface-probe .
```

Coordinator may run on the allocated isolated Xvfb display :189, with isolated
XDG directories and synthetic text. A window manager is optional for the bounded
probe; decorations and WM focus policy require one. This guide starts no server
and does not claim that :189 is currently free of other coordinator work.

```sh
mkdir -p /tmp/wci2-surface-probe/xdg/config /tmp/wci2-surface-probe/xdg/data /tmp/wci2-surface-probe/xdg/cache
env DISPLAY=:189 XDG_CONFIG_HOME=/tmp/wci2-surface-probe/xdg/config XDG_DATA_HOME=/tmp/wci2-surface-probe/xdg/data XDG_CACHE_HOME=/tmp/wci2-surface-probe/xdg/cache ./surface-probe >native.ndjson 2>native.stderr
```

Use XTest/manual OS input and screenshots, not direct invocation of Go callbacks.
Preserve display/WM/environment details and exact binary hash beside the log.
Entry logs hash/length rather than text. Probe field Commit is local simulation,
explicitly labelled; it is not SDL persistence or acceptance evidence.

Without a WM, limit claims to actual OS-delivered control input, independent
drafts, concurrent parent/child interaction, button-driven lifecycle, one result
per opening and no orphan child. The harness may position/focus X windows to
reach controls; record that assistance. Skip decoration-close steps and leave
WM focus/restoration/no-focus-stealing proof pending. `existing-opening-focused`
and `opener-focus-requested` log requests, not WM success. No-WM findings do not
close the full WCI2-A06/A08 acceptance rows.

## Native protocol — all pending

1. Edit the parent field, open child, edit its field without Enter. Alternate OS
   clicks/typing between both windows. Parent action must log childOpen=true;
   child stays visible and both draft fingerprints survive focus changes.
2. Press Enter in child: exactly one local-field-commit. Edit again, Escape
   reverts only the draft; second Escape produces one cancel surface-result.
   Explicit Cancel button requests opener focus. Check actual focus with OS
   input: the request log/Canvas.Focused alone does not prove WM delivery.
3. Reopen, edit child, switch to parent and choose Close child from parent.
   Exactly one close result; no opener-focus-requested and parent interaction
   continues. Reopen creates a new opening ID; old callbacks cannot commit.
4. Repeat with parent Hide for one second and Simulate reload. Each closes the
   child once with its distinct reason and discards only its unaccepted draft.
   Parent reappears without reopening the child; reload advances generation.
   This simulates ownership invalidation, not the SDUI atomic reload mechanism.
5. Repeat native child decoration close, parent decoration close and the Close
   parent button, including dirty child text. One result per opening, no duplicate
   callback from OnClosed fallback, no orphan child and process exits when both
   windows close. Do not call Hide directly outside the owner route.
6. Also inspect resizing, keyboard Tab order, parent action with child open and
   explicit-child cancellation focus. No IME/multiline/typed-control proof here.

Full SDUI WCI2-A06/A08 must later repeat these semantics with actual surface tokens,
dialog results, drafts and publication. This standalone probe has no SDUI state,
failed-publication injection, generalized reconciliation or implementation authority.
