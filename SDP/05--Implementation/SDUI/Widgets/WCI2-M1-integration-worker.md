# WCI2-M1 integration worker — lane frozen for integration review

Scope: only SDUI/go/host/fynehost and preparation, plus this report. Bridge and real panes fixture belong to Lorentz; runtime to James; layout to Gibbs; frontend/export to Dalton. No commits, branch changes or management writes. WCI1 baseline is delivered; WCI2-M2 is not selected.

Loaded/reused sdp-worker 2.0.0 and document-workflow, reviewed Panes-and-commands contract, runtime seams, WCI2-runtime-API.md and WCI2-frontend-API.md. Current handshake accepted by host: pure typed PresentationGate first; PresentationPrepare only on finalized candidate; runtime invokes non-failing ticket Publish after successful state replacement, otherwise Discard. Only Publish may promote Bundle.pending. Consumed sequence/error alone is not presentation authority; separately accepted reentrant legacy work must retain its own ticket.

Native composition strategy: reuse pinned Fyne AppTabs through one bounded keyboard/focus wrapper; bodies remain under shared geometry. AppTabs selects eagerly, so wrapper restores accepted native header selection before requesting the atomic application transaction. Renderer-returned actual CanvasObjects support header hit-coordinate inspection, not guessed labels. Fyne Split privately owns drag behavior and minimum clamping, so use a bounded native divider projecting layout-owned ratio/bounds and thickness, not an independent Split state. This is no general scene/input framework.

Agreed layout adjunct: MeasureTabs(instance,font,outer) -> TabsMetrics{Header Size}, MeasureSplit(instance,font)->divider thickness. Whole header minimum is measured from actual native AppTabs; layout returns Header/HeaderClip/Body and Divider/DividerClip/First/Second plus typed effective split bounds. Runtime owns selection/collapse/focus. Inspect must expose real header/page-ID and divider screen hit rectangles for Lorentz/main's XTest fixture.

Current native source is preliminary; no runnable pane/native acceptance or final test claim yet. Host is coordinating dependent signatures through main because direct thread messaging has repeatedly stalled. Earlier untracked WCI1 reports/binaries preserved.

## Recovery/implementation checkpoint

Owner resumed after interruption; preserved all lane changes. Runtime/layout API files reread, no restart. Host uses the agreed PaneMeasurer signatures. Preparation now accepts the typed presentation gate/final ticket hooks and discovers new callback owners without altering legacy Widgets. Pure stage never writes pending; finalized resources produce a ticket whose Publish alone promotes pending. Native input routes tabs/split via DispatchInteraction; Inspect adds `tabs[path]` with header/clip/body/selected/pages (per-page id/label/enabled/actual native rect/clip), and `splits[path]` with the shared divider/child geometry. Runtime Snapshot supplies all pane state.

Initial compile attempts are blocked by in-progress runtime undefined InteractionHandler/PageActivation/SplitChange, not a completed-lane test failure. No tests/native behavior claimed yet.

## Integrated M1 checkpoint

Implemented host pure presentation gate + final accepted preparation ticket, M1 capability admission, callback-owner/actual interaction binding postcheck, native AppTabs header wrapper, layout-owned divider adapter, shared native pane metrics, pane focus/keyboard/pointer interaction, and logical screen diagnostics. `controls` inspection includes actual mounted wrapper rect/clip/visible; `tabs` and `splits` expose native header/geometry for OS input. M2 remains rejected. Missing page icons explicitly reject provider icon/1, including hidden pages.

Meaningful host tests cover user header keys and real native TabItem pointer invocation, focus/remembered draft, Inspect targets, split collapse/restore, zero eligible/hidden pages, disabled skip, retired native callbacks, callback error with consumed sequence leaving accepted presentation unchanged, successful reentrant draft retained after stale outer reply, final resource error rejecting page publication, and pending provider preservation on rejected hide / cancellation on accepted hide / no revival on reveal. Preparation tests check no-op Binder rejection, zero callback calls during preparation, hidden icon requirements and no M2 widening.

Initial divider metric read the default theme's zero SplitThickness; fixed by measuring an actual empty native Fyne Split (its built-in fallback), rather than inventing native minima. Real SDL draft-conflict test exposed an extra application Guard call in final state preparation causing SDL Engine.Revision to deadlock inside Execute. Removed that call: Guard stays at Prepare/Commit, local source/bundle/size checks stay with the ticket, bridge checks the executed result. Real SDL panes + bridge tests subsequently passed (7.034s / 0.079s). A duplicate child name in the first host source was adjusted to distinct names; frontend/fixture own their naming semantics.

`go test -race ./host/fynehost ./preparation -count=1 -timeout 120s` passed at the preceding checkpoint (28.555s / 1.043s); final run including provider/inspection/empty-header changes is in progress. Lorentz independently reports bridge/panes/collections race passing. Desktop native panes build succeeded. No actual OS input acceptance claimed by this worker; main owns it.

## Final host/preparation freeze — 2026-10-08

All owned source is held for coordinator integration/review. No running test/build jobs remain. No commits, branch switches, original-workspace edits, management writes or M2 implementation were made. Earlier checkpoints above are historical; this checkpoint supersedes their pending-test status.

The last native pilot observation was transient header coordinates after unrelated provider state publication: AppTabs items were recreated at every sync. The bounded correction calls SetItems only when ordered page ID/label/enabled metadata changes. Draft/provider/selection-only publication retains actual native targets. A regression verifies native TabItem and button identity, laid-out position/size across an unrelated draft and accepted selection, and metadata-change rebuilding including disabled state. The targeted header/provider tests passed in 0.887s.

Final commands, using cached dependencies and GOCACHE=/tmp/sdp-wci0-gocache:

- In SDUI/go: `go test -race ./host/fynehost ./preparation -count=1 -timeout 120s` — PASS, host 34.857s; preparation 1.048s, after the header stability correction.
- In SDL/go: `go build -tags desktop -o /tmp/sdui-wci2-native-noether-check ./examples/panes/cmd/native` — PASS, after the header correction.
- In clone root: `git diff --check -- SDUI/go/host/fynehost SDUI/go/preparation` — PASS.
- Earlier integrated real SDL validation: `go test ./examples/panes ./bridge -count=1 -timeout 60s` — PASS (7.034s / 0.079s). Lorentz separately reported affected bridge/panes/collections race PASS; this is delegated evidence, not a rerun by this worker after the final header change.
- Host README relative links validated. No network/sandbox block remained.

Coordinator reported native panes 25 PASS, expanded horizontal split 11 PASS, provider lifecycle 7 PASS and vertical split 9 PASS (52 checks), with clear screenshots. Those are coordinator pilot results on preceding candidate bytes, not an exact-final-binary acceptance claim by this worker. Recorded harness corrections: clear dirty receiver before independent resource failure; assert final drag state rather than intermediate drag; await paint-ready header coordinates. Main owns final full-suite/native rerun and independent acceptance.

Remaining scope/limits: no known host test failure or integration blocker. Native AppTabs introspection is bounded to the pinned Fyne renderer's actual CanvasObjects. Missing icon provider rejects admission, including hidden declarations. M2 command/menu/dialog code is absent. Whole WCI2 delivery/approval remains the coordinator's responsibility.

## Bounded outcome-reporting correction and final refreeze

Coordinator explicitly reopened only native outcome reporting after the preceding freeze. Native pane dispatch previously discarded InteractionResult when forwarding errors. Exported `fynehost.InteractionError{Result, Err}` now preserves it; Error includes status/domain/sequence, and Unwrap retains the original error for errors.Is/errors.As. Both header and divider dispatch wrap failures. The host's existing Mutate/OnStatus path reports once, avoiding a duplicate outer status call. No runtime, lifecycle, preparation, bridge or other lane edits were required.

Regression covers final resource failure after DomainSucceeded and DomainUnknown, unchanged accepted UI, exactly one native error notification, precise result/text/sequence, and original typed fault identity through both Is and As. A handler returning an error with unspecified outcome reports domain=unknown, status=ui-conflict and preserves its cause. This distinguishes a completed or uncertain domain call from UI publication failure; it does not claim domain rollback.

After this final correction:

- SDUI/go: `GOCACHE=/tmp/sdp-wci0-gocache go test -race ./host/fynehost -run 'TestPane(FinalResourceFailure|HandlerError|FailedPane|NativeKeyboard|Header|Retired)|TestFailedPaneCallback|TestReentrantAcceptedDraft' -count=1 -timeout 60s` — PASS (20.347s).
- SDL/go: `GOCACHE=/tmp/sdp-wci0-gocache go build -tags desktop -o /tmp/sdui-wci2-native-noether-check ./examples/panes/cmd/native` — PASS.
- `git diff --check -- SDUI/go/host/fynehost SDUI/go/preparation` — PASS.

Files are frozen again. No jobs remain running. The broader host/preparation race result above predates only this bounded correction; main explicitly owns the final full suite and native rerun on the following bytes. No further implementation scope was opened.

## Selected-header focus and disabled-header correction — final refreeze

Coordinator reopened only the native header adapter after an actual native same-page click left the input focused. Pinned Fyne AppTabs selectIndex returns before OnSelected on a same-index click. A bounded transparent, non-focusable pointer surface now tests the actual native button rectangles: an eligible selected-page click requests only accepted tabs focus; changed-page clicks delegate to the real native button/callback transaction. No reflection, unsafe access, guessed tab widths or eager changed-page focus were introduced. Native drawing and hover remain AppTabs-owned.

The canvas-hit regression starts with a raw input draft, clicks the selected header, and verifies native/runtime tabs focus with zero action calls, no consumed sequence and preserved draft. It then clicks another page whose handler rejects and verifies the original input focus, page and draft remain. Targeted non-race focus tests passed (0.856s).

Independent disabled-header repro was read from `/tmp/wci2-review-disabled-5scitliz/review_test.go`: disable the tabs owner then call the second actual tabButton.Tapped; native selection previously diverged from runtime. Disabled OnSelected now restores the accepted selection instead of returning after AppTabs' eager change. The pointer surface also suppresses disabled interaction. Regression covers actual canvas click and direct native tabButton invocation with zero callbacks/state mutation, owner re-enable retaining individually disabled page eligibility, and subsequently enabling that page restoring operability. Inspect tabs now reports `nativeSelected`, mapped from actual AppTabs.SelectedIndex to stable ID, separately from accepted geometry `selected`.

Final commands after these changes:

- SDUI/go: `GOCACHE=/tmp/sdp-wci0-gocache go test -race -overlay /tmp/wci2-review-disabled-5scitliz/overlay.json ./host/fynehost -run 'TestReviewDisabledNativeHeader|TestPane|TestFailedPaneCallback|TestReentrantAcceptedDraft' -count=1 -timeout 90s` — PASS (32.373s), including the independent exact repro, all pane regressions and domain outcome reporting.
- SDL/go: `GOCACHE=/tmp/sdp-wci0-gocache go build -tags desktop -o /tmp/sdui-wci2-native-noether-check ./examples/panes/cmd/native` — PASS.
- `git diff --check -- SDUI/go/host/fynehost SDUI/go/preparation` — PASS.

Host/preparation lane is frozen again; no tests/builds remain running. Main owns the pending actual OS-input regression rerun and final full-suite/reviewer acceptance. Previous pilot success counts do not certify these final bytes. Domain reporting correction remains included and passing. No new scope, M2 work, commits or other-lane edits.

### Candidate identity and owned changed files

HEAD plus uncommitted shared-lane work; this is not evidence for unchanged HEAD alone. Other lanes remain their owners' responsibility. This inventory supersedes previous freeze hashes and excludes the report's own bytes.

HEAD: `53d031e9a9a6e10bdcf56a2ac7b6e863adc81a9d`

Native binary: `/tmp/sdui-wci2-native-noether-check`

SHA-256: `f1de87cefd039f84321c145d02f4b79bbd219916a30b76dfaa3a2007600239cc`

| Owned changed file | SHA-256 |
| --- | --- |
| `SDUI/go/host/fynehost/README.md` | `d3a1e4c7bc107c68237a65185aed0dbfc0b2a6d7c337467ebe45ed7ec293e330` |
| `SDUI/go/host/fynehost/admission/admission.go` | `19618ec9552d82b20db1f96532d614ed1bba056be5e900f43bb40ea8b5477ca0` |
| `SDUI/go/host/fynehost/document.go` | `afd0d058735583aade9ab63184d6009575c4c5827f5e09c44a69331a12ab09e9` |
| `SDUI/go/host/fynehost/document_sync.go` | `3639a1163692c4dce7e374cc06a892332d2da3d0c18eecf40444a091ac3088ad` |
| `SDUI/go/host/fynehost/inspection.go` | `e962e66f0de19f0a4a423a1d7d080ff66e3165fff09de2322618a7f081c3cf06` |
| `SDUI/go/host/fynehost/pane_divider.go` | `9c1f68c2b8754e75ca1be264ff40779d4211f3d1694439091b07b572552b4187` |
| `SDUI/go/host/fynehost/pane_events.go` | `87dcd76dfce028018675c755a944fa47084bd12442415542c8acf66e03f5aec1` |
| `SDUI/go/host/fynehost/pane_header.go` | `af270cf7113e306c416997ace3e532b16e6426fa5eef730f402b1e6822e5ac80` |
| `SDUI/go/host/fynehost/pane_metrics.go` | `93baaffc4115ea4b8135c1943cc9148c9091cc7d742b6a3490eb9b91808564bf` |
| `SDUI/go/host/fynehost/pane_test.go` | `ecfba480753d8bda48153c898d5f6b6f4819caa9588708874f7a651b6ce96c9d` |
| `SDUI/go/host/fynehost/view.go` | `80086dd29aec41a16516d0c3942a76c14311753ea9626d01e8415004a00148ce` |
| `SDUI/go/preparation/capabilities.go` | `824242c742b9e66d25bdbf39e4571f677e3dfce7ef455c1efab60e27cb1691f8` |
| `SDUI/go/preparation/panes_test.go` | `7853ae676a6d02fcb1fd004d1a116c79bdbb7b491dfeefbe4dc1a3ab4b9c4f30` |
| `SDUI/go/preparation/prepare.go` | `f0413e0804bd7e8d5ba448e4f64cda6ed8893c1375d82354653818b336071746` |
