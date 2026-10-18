# WCI3-M2 integration worker handoff

Status: host/preparation product files frozen for independent final review and
main's exact-candidate integration/native acceptance. Component implementation is
complete in the assigned lane; no whole-stage or OS IME acceptance is claimed.
Date: 2026-10-08. Role: Noether, SDP Worker 2.0.0; SDP 1.1.1 and the shared document
workflow reused. Main owns Session/management/traceability, canonical system docs,
module/dependency changes, external native input and final disposition.

## Authority and scope

Isolated clone `/tmp/sdp-sdui-widgets`, branch `sdui/widgets-wci3`, base HEAD
`ea49991f2065679c93e39fe02a3d458fe72803df` (independently approved M1).
At freeze, main's evidence/selection commit advanced observed HEAD to
`d6742742f0968661d15698518df688ffb1d28946`; this worker made no commit.
Read ORIGINAL Values-and-text and its reviewed native editing refinement, the
root M2 host/frontend/runtime/layout memos, Worker instructions and SDUI/AGENTS.
Latest inspected canonical SHA-256:
`8199c275fecff3f6d03835540a3e586dd418ff70470f641590d2d75ace8f51be`.
Explicit coordinator M2 authorization permits this host/preparation lane only.
No commits, branch changes, module/dependency edits, original-workspace edits,
management writes, other-owner changes or prior evidence rewrites were made.
All 17 changed lane paths are listed in the exact manifest below; this root report
is the only additional handoff file written for implementation.

## Delivered behavior

* Only explicit new-argument presence opts into extended input. Parser.InputOptions
  selects it and runtime FieldState.Input supplies copied policy. Legacy .2/basic
  .3 still use their previous native Input/events and geometry. Hidden/closed
  multiline/readonly requirements are checked before native publication.
* The native text composition retains one Entry and measured optional label/fixed
  validation feedback. Three-row multiline word wrapping and its private editing
  state belong to Entry; runtime Widget.Value/Draft remain the sole model store.
  Whole Entry receives FieldLayout.Control. No reflection, parallel editor/history,
  copied runtime scroll authority, generic event framework or Fyne fork was added.
* Checked EditField is the only native draft publication. Single-line Enter and
  exact Primary+Enter perform typed CaptureCommit/Dispatch; multiline Enter and
  Shift+Enter insert newline. Tab/blur never commit. Native undo/redo only edit
  drafts. Exact Primary+Shift+Z is translated from GLFW CustomShortcut to native
  named Redo, alongside the native Primary+Y form; extra modifiers remain distinct.
* Readonly retains focus/selection/copy. Guards cover typing, delete/word-delete,
  cut/paste, Undo/Redo, Commit and late editing-menu actions. A bounded native
  editing menu routes through the same guards; declared context menus keep WCI2
  semantics. Keyboard and context Paste capture clipboard content once, reject
  CR/LF before single-line insertion, and delegate those exact checked bytes.
  Pre-delegation refusal leaves native history unchanged.
* Identical displayed bytes preserve history, including Apply and self-echo.
  Different authoritative text uses one muted SetText. Actual rejected native
  EditField restores the CURRENT authoritative draft on the same Entry and clears
  native undo/redo, under the reviewed caret/selection/scroll reset exception.
  Accepted-invalid drafts, failed Commit/reload/probe, page hiding and ordinary
  synchronization do not use that exception. Reentrant accepted work survives.
* An actual mounted-component test exposed Fyne Refresh revealing the caret during
  unrelated synchronization. Same-text publication now retains the observed native
  offset using public ScrollToOffset, with native resize clamping. Native wheel
  consumes within Entry even at its limit; user typing/arrows still reveal the
  caret. Focus loss/muted restoration preserve observed internal offset. These
  observations are scoped to synchronous painting, never persisted in runtime.
* Detached measurement uses an empty, identically configured native Entry. Pinned
  Fyne's scrolling minimum uses themed character metrics and visible row count.
  Real short/long Unicode/unbroken/wrapped-CRLF controls at fonts 14/20/28 prove
  equality with empty minimums. This removes repeated full-draft SetText/shaping
  from probes; live controls still receive exact text. No measured performance or
  native latency improvement is claimed solely from that structural change.
* Inspector fields expose actual Entry text/placeholder/readonly/focus/cursor and
  selection, plus actual control/entry/label/feedback rectangles in named canvases.
  Actual public renderer Scroll adds scrollRect/scrollOffset (Size W=X,H=Y), with
  no fabricated geometry or second state API. README documents usage and limits.

The existing private accepted-presentation ticket remains the only publication
handshake. No new runtime/frontend/layout stubs or signatures were introduced.
Source, provider readiness and bridge behavior remain owned by their respective
lanes. Main's separately selected pinned GLFW X11 filter is consumed through the
main-owned Go replacements; this lane made no driver or dependency changes.

## Commands and results

All Go commands below use `GOCACHE=/tmp/sdp-wci0-gocache`. All reported successful
commands completed with exit 0. No network/sandbox bypass or permission change.

From `SDUI/go`:

| Command | Result and candidate scope |
| --- | --- |
| `go test ./host/fynehost ./preparation -count=1` | PASS host 15.004s, preparation 0.045s; before the final Redo/measurement deltas. |
| `go test -race ./host/fynehost ./preparation -count=1` | PASS host 222.640s, preparation 1.055s; before those final deltas. The process remained CPU-active; it was not hung. |
| `go test -race ./host/fynehost -run 'TestExtendedTextSubmissionAndHistory\|TestExtendedTextClipboardAndReadOnly\|TestExtendedTextScrollPreservationAndUserReveal\|TestExtendedTextDetachedProbesPreserveEditing' -count=1` | PASS 7.187s on Redo-fixed candidate; before empty measurement. |
| `go test ./host/fynehost -run 'TestExtendedTextEmptyMeasurementMatchesNativeContent\|TestExtendedTextClipboardAndReadOnly' -count=1` | PASS 0.868s after empty measurement and deterministic readonly Redo setup. |
| `go test -race ./host/fynehost ./preparation -run 'TestExtended\|TestHiddenExtended' -count=1` | PASS host 25.100s, preparation 1.031s on final lane bytes below. |

From `SDL/go`:

| Command | Result |
| --- | --- |
| `go build -tags desktop -o /tmp/wci3-m2-noether-redo ./examples/text/cmd/native` | PASS; Redo-fixed predecessor artifact preserved. |
| `go build -tags desktop -o /tmp/wci3-m2-noether-text-fast ./examples/text/cmd/native` | PASS; current artifact with Redo and empty measurement. |

`git diff --check` PASS. Meaningful regressions cover original-target/current-draft
reentrance, same Entry/focus after rejection, zero Commit/no publication on failed
edit, declared history reset, invalid required draft undo, same/different Apply,
failed Commit/reload retention, inactive pages/collapse, exact clipboard one-read
and pre-delegation refusal, readonly late popup actions, source capability checks,
actual native scroll hit targeting/limits and pure detached preparation.

Mendel's two independent scroll/caret-reveal and detached-probe regressions were
persisted as `text_probe_test.go`, with attribution. Reviewer separately reports
host race 41.110s on the core/Redo source plus those cases, and empty-measurement
race 3.576s with scoped approval. This is independent reported evidence, not a
worker-run independent review or final whole-host/native verdict.

During development, tests exposed unmounted test-canvas focus assumptions and a
readonly fixture relying on merged undo grouping; the fixtures were corrected.
Readonly now uses a second explicit native Paste action, Undo, and an immediate
baseline assertion before testing blocked Redo. No product guard was weakened.
The observed same-text scroll reset and actual OS Primary+Shift+Z miss were real
product findings and received the bounded fixes described above.

## Artifacts, reported native observations and remaining acceptance

Current desktop artifact:
`/tmp/wci3-m2-noether-text-fast`
SHA-256 `0078f72eec7fe61bd41a1b090bb7ca2f26ba48c3a1e54c4eff6bd240c010f001`.

Preserved predecessor `/tmp/wci3-m2-noether-redo`, SHA-256
`4bb0f6a07694712eb37bafb6c3fb27d35760f16fa4e677fc5c39157500645af2`.
Main reports actual multiline eight-check PASS on that predecessor including
Primary+Shift+Z Redo, plus native inner wheel/no outer route, unrelated-sync offset
retention and page-hide retention. Those reports are historical candidate evidence;
main owns final native/IME/SDL receipts on its selected integrated binary.
Main subsequently reports readonly five-check PASS on `0078`, 13.21 seconds
including cold launch in that run. This is a coordinator-observed pilot duration,
not a general latency guarantee. Main's history harness was corrected to use a
second clipboard Paste as a distinct undo action: native grouping may legitimately
merge a paste followed by typed letters. No product defect or source change was
inferred from that harness assumption; forced-rejection assertions remain strict.

Pinned Fyne on Linux does not register Ctrl+End for Entry. This adapter deliberately
matches actual native shortcuts and adds only the selected submission/Redo routes.
Ctrl+End is not a new navigation promise; real repeated Down/navigation proves
caret reveal at the end. Main was informed when its pilot observed the unsupported
shortcut. Main confirmed and corrected the harness to use plain End for single
line and actual selection collapse/arrows for multiline, retaining the earlier
failed harness record. Public native PageUp/PageDown/Home/End/arrows retain their
Fyne behavior; no additional keybinding was requested or added.

No remaining confirmed host defect at this handoff. Still required: independent
final review of exact lane manifest, main's full integration suites and final actual
native/SDL/IME acceptance. Main is also investigating startup/transient native
minimum reports and observed 15-second pilot waits under concurrent race/text load;
this lane does not claim a cause or close those observations from the optimization.
No OS IME composition is inferred from component tests or injected Unicode. Main
owns the pinned driver and configured IBus proof. Source files are held unless a
new concrete finding is assigned; no WCI4 product implementation was started.

## Exact lane freeze

The following manifest covers all changed product/docs/tests in this lane, relative
to the clone root. The worker report itself is excluded from its own digest.
Manifest SHA-256: `17904e385d05e28a1640877986c2e8f520100151d095f8239859a269546744b6`.

Subsequent independent disposition: Mendel approved this exact 17-path host/prep
manifest after recomputing it and verifying unchanged bytes before/after independent
`TestExtended|TestHiddenExtended` race tests (host 24.652s, preparation 1.020s,
PASS). No remaining source/component blocker was reported. This closes the scoped
lane review requirement above, not main's whole-candidate/native/IME acceptance.

```text
736abf8586beb0745344017ee3c87900e4a01b8cfbabf64ea6f47ebba173692c  SDUI/go/host/fynehost/README.md
f12cd318811279141b20892c699a2fbaa71e4635ac0dd5a59f9cdad483fefad9  SDUI/go/host/fynehost/admission/admission.go
fe7c7db6ea0f13aa74461099c84152fb3f5b8b16b28fee49bf1cf738f141ad6b  SDUI/go/host/fynehost/document.go
b903c4d342ae530920224f72aa70ee4dc506b039587608844d4e9ef6636b86f8  SDUI/go/host/fynehost/document_commands.go
b717864642da3544ad7bd43cf176559dec062e92db9f88efc1014cfd4f47cc53  SDUI/go/host/fynehost/document_sync.go
9dc49ef927b616125f53a86d48b4cebb33c6cb3a0ee0d5e4ec1115d9b459fe11  SDUI/go/host/fynehost/scalar_choices.go
aa594ed252838261624a9134803c6e7e0d197e91af09005e958af3dacb2a1f18  SDUI/go/host/fynehost/scalar_control.go
fc4bfe754738452c1465aba39c1c1bcfa8041aed9ceb795fd6517475c52b70ca  SDUI/go/host/fynehost/text_control.go
f0437a0045cb8b6d4980bd2f49b9ae17be792db0b3d4008558eee1692eb7dcf6  SDUI/go/host/fynehost/text_control_test.go
db69c3648969340fa424aefab0ce7597db3ec9a82134a26ec88464bc1b749a9f  SDUI/go/host/fynehost/text_entry.go
4178f811747bebdff332aa0e36234a51d61508be87d5adea5f242bc9f5488414  SDUI/go/host/fynehost/text_inspection.go
ba031af050c7e48510bdc11a864325be947b3d546f31b0c04838090c3432f498  SDUI/go/host/fynehost/text_lifecycle_test.go
71c8573cbf0e90a7bcd141c3c602dae4b997133f87f990851e8f10f569e15f95  SDUI/go/host/fynehost/text_metrics_test.go
fec8284842f17dc798ec7089355599c9e4a9e28cfaf9c8eb6f50c6ab2fd05535  SDUI/go/host/fynehost/text_probe_test.go
7331c2e0b201f1057fcd696b781805197630ca712c9aeba2f420cea06d5237c8  SDUI/go/host/fynehost/view.go
85823f21f160b841c7b4e7c0c4fa995ebcbf2850248729c70d2b068e5ab391b8  SDUI/go/preparation/capabilities.go
f2b4a1f88fa2246c80b62b1d34603c41131899978669ffa780d6bac3a5f48240  SDUI/go/preparation/text_test.go
```
