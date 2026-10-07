# WCI1 integration worker — in progress

Candidate: isolated `/tmp/sdp-sdui-widgets`, branch `sdui/widgets-wci1`, HEAD `d1c5c88d1989b173d610dad916f6f493f838ab7b` plus shared uncommitted changes. No commit or original-workspace edits. This is worker evidence, not independent review or full slice acceptance.

## Ownership

Host/preparation/bridge integration worker retains `SDUI/go/host/fynehost` except `document.go` and James's dedicated `document_lifecycle_test.go`, `SDUI/go/preparation`, and `SDL/go/bridge`. `document.go` is explicitly frozen and transferred to James for reentrant error reconciliation, failed-resize lifecycle recovery and effective-snapshot resource validation. `document_sync.go` now uses `b.size` for prepared image/container dimensions in support of that handoff. Existing `document_test.go` remains integration-owned. Entire `SDL/go/examples/collections` subtree (including native command/tests) is frozen and transferred to Gibbs. Prototype is Dalton-owned. No other lane's changes are claimed here.

## Implementation evidence

Detached document preparation validates actual Markdown/provider resources, prepares native controls and SVG omission only for actual native controls, and gates runtime state with snapshot layout. Bundle callbacks and asynchronous provider completions are scoped to live bundle/request identity. Cancellation revokes runtime acceptance before canceling provider contexts. Actual SDL bridge EventField resolves collection item IDs and checks the captured receiver's identity and revisions after execution.

The native collection adapter has one focus stop per collection, runtime-owned selection/expansion/offsets, explicit retry/cancel, and ASCII markers rendered with the bundled GoRegular font. Geometry uses the same row text and font metrics as painting. Page scrolling uses an effective clipped point; reveal uses measured row bounds. Thumb drag uses travel excluding thumb length and retains the initial grab offset.

Bounded custom adapter rationale: pinned Fyne collection controls keep inner scrolling/key behavior private; their independently maintained scroll position cannot be the runtime's authoritative offset. Existing native Input and Button remain reused. A forwarding clipping wrapper and bounded viewport surfaces route background/control wheel input and frame/group thumbs through shared layout/runtime offsets. This does not introduce a general scene framework.

## Current verification

Commands run with `GOCACHE=/tmp/sdp-wci0-gocache`:

- SDUI/go: `go test ./host/fynehost ./preparation -count=1 -timeout 90s` — PASS.
- SDL/go: `go test ./bridge ./examples/collections -count=1 -timeout 90s` — PASS.
- SDUI/go: targeted `TestDocumentNativeTabOrderAndDisabledSkip` and `TestRuntimeReloadNativeDraftFocusAndNoReplay` — PASS after fixing native wrapper base identity.
- SDUI/go: `go test -race ./host/fynehost ./preparation -count=1 -timeout 120s` — running at this checkpoint; subsequent event guard edits require another final run.

New native focus regression traverses the mounted Fyne object tree with empty providers: nav, entries, preview; disabled nav is skipped. The legacy reload regression initially failed (runtime focus and draft preserved, canvas focus nil). The forwarding wrapper had used an already self-extended Scroll, so extending the wrapper could not replace its base identity. Constructing its Scroll before extending to the wrapper fixes this; no test expectation was weakened.

Partially clipped text now has explicit MinSize paint bounds. Regression checks negative-X text has a nonempty visible extent; it does not prove OS rendering. The initial real SDL fixture canonical Actions order was corrected before handoff and its compile/run test passes.

## Remaining verification and limitations

Main's old-binary pilot has ten reported passing actual-input checks, but reported incorrect first-Tab order and potential negative-X text clipping. Those are preliminary external observations, not exact-final-candidate evidence. Main must rerun native empty-root/disabled focus, nested viewport/background input, Alt+arrow, thumb drag, resize/reload/lazy barriers and negative-X partial text painting on a settled rebuilt binary. Canvas.Capture black output is a reported harness limitation; use OS screenshots for visual acceptance.

James owns unresolved lifecycle fixes described above. Final integration race/desktop build and changed-file hashes remain pending. This report does not claim full WCI1 publication/acceptance.

## Native adapter checkpoint

Additional passing targeted tests: `TestEmptyReloadPreservesNativeFocusAndExplicitLoadKey`, `TestCollectionPartiallyClippedTextHasPaintBounds`, `TestCollectionThumbGrabAndAltScroll`, `TestFrameBackgroundAndWrappedControlWheel`. The reload test delivers R through the actual canvas-focused control after a loading-root successor pauses, and requires a new request. The thumb test verifies a stationary grab does not jump and displacement follows track travel. Frame/background and wrapped-button wheel tests verify the runtime offset changes, rather than a separate Fyne scroll position.

Collection input now rejects a rendered collection generation/provider epoch that no longer matches runtime; viewport map publication checks the captured model and every viewport handle before atomic SetViewports. Hidden/zero-extent native controls are hidden without hide/show cycles for active controls. Input focus failure stops reveal rather than continuing after a rejected gate.

Latest whole-lane race tests and desktop build are running; this checkpoint supersedes neither James's lifecycle ownership nor main's required native acceptance. Earlier host/preparation race passed before these final adapter changes.

## Verification checkpoint after adapter fixes

- `GOCACHE=/tmp/sdp-wci0-gocache go test -race ./host/fynehost ./preparation -count=1 -timeout 120s` (SDUI/go): PASS, host 13.389s, preparation 1.079s.
- `GOCACHE=/tmp/sdp-wci0-gocache go test -race ./bridge ./examples/collections -count=1 -timeout 120s` (SDL/go): PASS, bridge 1.129s, collections 22.417s.
- `GOCACHE=/tmp/sdp-wci0-gocache go build -tags desktop -o /tmp/sdui-wci1-native-noether-check ./examples/collections/cmd/native` (SDL/go): PASS. This is a checkpoint binary, not final settled lifecycle candidate.

Corrected main observation: the old-binary nested pilot passed inner-to-outer wheel and reverse. Actual drag at x993,y100→170 moved OUTER body Y to 112 while inner Y stayed 0, consistent with outer-thumb travel. Input was delivered; the intended inner thumb was occluded. Separate nested-thumb accessibility remains FAILED/pending a geometry correction and rerun; unit drag math is not a native input claim. Logged geometry places collection gutter at x988–998 and outer frame overlay thumb at x988–1000, so x993 overlaps both. Frame gutter allocation is layout-owned and requires coordination; no layout edits made here.

### Integration-owned file inventory and checkpoint SHA256

These are uncommitted checkpoint bytes, not immutable final evidence. Shared document.go and the fixture/prototype lanes are intentionally excluded from this ownership hash list.

- `SDL/go/bridge/bind.go`: `0ba92ee0eb00df3f435851623b734d2cee3eec526ff1733fc57febc6426b2a91`
- `SDL/go/bridge/values.go`: `37742045dd6546f73aee4c8b12d3cc5589409abb884360ca9a1d1bc202827f4b`
- `SDL/go/bridge/collection_test.go`: `df6c28ac774b9c4ee9fa71db35a38da4bfee46116c28753009c6145aea112a8a`
- `SDUI/go/preparation/capabilities.go`: `f44c0361f89fdcb6d36e7065e73afb66a478cdd0adaa6cb824ca52365f320ab4`
- `SDUI/go/preparation/prepare.go`: `a1b346d1e806433c72b93fece23835053dfb470a5aed07f1bcd7a3f17b86998b`
- `SDUI/go/host/fynehost/admission/admission.go`: `638808148a3505d1fb9b3b7e9d4d530d0ca2fd1930d4364a8c61401a0965a82a`
- `SDUI/go/host/fynehost/view.go`: `04355dc2527332a157ae13787da2a87d71f49cd6a7afb6526bf647a4c7bc5900`
- `SDUI/go/host/fynehost/runtime_test.go`: `e9603f70083b738eaaef134213eba338208aed78192da0fb5995acb9b4ed40df`
- `SDUI/go/host/fynehost/collection_control.go`: `a73b85b350b225d6d73af43cd7c808d8dc5e7dffb0b90b301b8decebddf6a5e3`
- `SDUI/go/host/fynehost/collection_metrics.go`: `5480f1e2315289ce8082477a0effa57891dce43925a74fe945f5529232ae5827`
- `SDUI/go/host/fynehost/document_events.go`: `1474c6ffebe28d5845b2222c7b297037cbca13fc6a39871c9565e8e2bafb41c2`
- `SDUI/go/host/fynehost/document_sync.go`: `9b69da19ed9fad61b314f02c449d85333649c445bab033da95b73eaaa809031d`
- `SDUI/go/host/fynehost/document_test.go`: `c216d12203ebaa2fc3c8c7550035ad1cd2f88f8a39365e667956673c91fb83a2`
- `SDUI/go/host/fynehost/inspection.go`: `0fecdc67cce9c48bed4a0a0ff9a6df84f4a6fad58d3a262a6df94320fa50eb3f`
- `SDUI/go/host/fynehost/viewport_control.go`: `7a8bc665b481e222db55e96561d6afe038a6879e719d07e79d324f1aeeb2d782`
- `SDUI/go/host/fynehost/viewport_control_test.go`: `2906c656815b8f4ced1404266358d0946589078b69215935809c357d07bb460c`

## Nested thumb correction and lane boundary

Main corrected the initial no-event interpretation: x993 hit the outer overlay thumb, and body Y112 is the observed result. Required repair is separately reachable inner and ancestor thumbs using shared measured gutter geometry, not a relaxed harness assertion. Main is coordinating Gibbs's minimal frame/group gutter proposal. Integration will consume that shared geometry after the API/ownership handoff; no parser pixel syntax, generic scene, or layout-lane edit is introduced here. Native nested thumbs remain an open acceptance issue.

## Composition documentation

Added `SDUI/go/host/fynehost/README.md` describing actual Prepare/Commit/Mutate/Close use, detached candidate disposal, owner-goroutine Post delivery, typed provider and external identity responsibilities, resource hook ownership, and legacy 0.2 RuntimeView differences. Relative links to the real SDL fixture, runtime, preparation API/design and collection contract were checked and exist. Documentation explicitly retains pending native acceptance and nested-thumb geometry issue. This update changes only the README and this worker report; James retains lifecycle-file ownership. No code tests repeated for the documentation-only change.

## Approved gutter integration in progress

Independent stage review and coordinator authorization received. Added native `viewport_metrics.go`: pure declared-axis frame/group insets (right for Y, bottom for X), with legacy/collection behavior unchanged. Updated viewport surfaces to consume layout-owned HorizontalGutter/VerticalGutter strips, including background wheel targets when a declared axis has zero range; thumb origin uses the uncut content viewport while clipping uses the supplied effective strip. No padding is reconstructed in host code. Added regressions for declared-axis insets and separately reachable nested tracks with independent runtime offset changes. Gibbs's layout API/types are not yet present at this checkpoint; build/test verification awaits that shared lane. James's document.go remains untouched.

## Coordinator native empty-root evidence

Main reports `/tmp/sdui-wci1-harness-empty-pilot4b` passed all seven actual-native checks on the source build before gutter refinements: Escape rejects late completion; R creates a fresh request; reload pauses loading; native R after reload creates another fresh request; pre-reload completion is rejected; empty success becomes loaded; teardown closes the session and leaves no pending barriers. Initial pilot4 failure was an ended Xvfb server, restarted with `-noreset`, not a product defect. This is coordinator-reported exact-pilot evidence, not a claim for subsequent gutter bytes. Final native rerun remains required after gutter integration.

## Final integration lane freeze checkpoint

Gibbs's approved shared gutter API is integrated. Native adjunct reserves right only for declared Y and bottom only for declared X; viewport surfaces consume effective layout strips without reconstructing padding. Native collection metrics remain unchanged. Independent nested-track geometry and per-owner offset regressions pass.

PageUp/Down now moves 90% of `Viewport.Clip.H`, not the full offscreen `Rect.H`. New `TestCollectionPageKeysUseAncestorClippedHeight` exercises a collection twice its ancestor height, verifies the visible-height forward step, reverse step, and no premature ancestor movement.

Scoped verification after this final code edit: `GOCACHE=/tmp/sdp-wci0-gocache go test -race ./host/fynehost -run 'TestCollectionPageKeysUseAncestorClippedHeight|TestNestedNativeThumbsUseDisjointSharedGutters|TestNativeViewportInsetsDeclaredAxesOnly' -count=1 -timeout 60s` — PASS (1.999s). Main reports the immediately preceding full SDUI race suite passed (host 16.571s), plus SDL bridge/runtime/codegen/fixture passed. Full suites were not unnecessarily repeated for the two-line paging fix. Desktop checkpoint build is being completed below.

Integration lane is frozen after this checkpoint, except genuine review findings or explicit coordinator follow-up. James retains document.go/lifecycle-test ownership. Main owns final exact-candidate native input/screenshots and full acceptance. No native nested-thumb pass is claimed by this worker.

Final integration-owned checkpoint hashes (supersede earlier corresponding hashes):

- `SDUI/go/host/fynehost/README.md`: `931bc6d009e6196b5140760544bef114c4c221b89890fddda2e2a3304076619b`
- `SDUI/go/host/fynehost/admission/admission.go`: `638808148a3505d1fb9b3b7e9d4d530d0ca2fd1930d4364a8c61401a0965a82a`
- `SDUI/go/host/fynehost/view.go`: `04355dc2527332a157ae13787da2a87d71f49cd6a7afb6526bf647a4c7bc5900`
- `SDUI/go/host/fynehost/runtime_test.go`: `e9603f70083b738eaaef134213eba338208aed78192da0fb5995acb9b4ed40df`
- `SDUI/go/host/fynehost/collection_control.go`: `a73b85b350b225d6d73af43cd7c808d8dc5e7dffb0b90b301b8decebddf6a5e3`
- `SDUI/go/host/fynehost/collection_metrics.go`: `5480f1e2315289ce8082477a0effa57891dce43925a74fe945f5529232ae5827`
- `SDUI/go/host/fynehost/document_events.go`: `5a3d0895cce9642250c0643ea6e863236f56b06691637b26da79622675fd528c`
- `SDUI/go/host/fynehost/document_sync.go`: `9b69da19ed9fad61b314f02c449d85333649c445bab033da95b73eaaa809031d`
- `SDUI/go/host/fynehost/document_test.go`: `c216d12203ebaa2fc3c8c7550035ad1cd2f88f8a39365e667956673c91fb83a2`
- `SDUI/go/host/fynehost/inspection.go`: `0fecdc67cce9c48bed4a0a0ff9a6df84f4a6fad58d3a262a6df94320fa50eb3f`
- `SDUI/go/host/fynehost/viewport_control.go`: `652a758bcc3975316285f43fcb22e0423fed0c8cbcd19df25ebd20ab457e5fdb`
- `SDUI/go/host/fynehost/viewport_control_test.go`: `add0ff04dc14315df9d0df39b81fd0fc8ab874dfcecdc633393aa0523576ccb6`
- `SDUI/go/host/fynehost/viewport_metrics.go`: `238b0fc76a2a3a1062947fe97da43f987a43d85d836decb8c91f6ca53082db28`
- `SDUI/go/preparation/capabilities.go`: `f44c0361f89fdcb6d36e7065e73afb66a478cdd0adaa6cb824ca52365f320ab4`
- `SDUI/go/preparation/prepare.go`: `a1b346d1e806433c72b93fece23835053dfb470a5aed07f1bcd7a3f17b86998b`
- `SDL/go/bridge/bind.go`: `0ba92ee0eb00df3f435851623b734d2cee3eec526ff1733fc57febc6426b2a91`
- `SDL/go/bridge/values.go`: `37742045dd6546f73aee4c8b12d3cc5589409abb884360ca9a1d1bc202827f4b`
- `SDL/go/bridge/collection_test.go`: `df6c28ac774b9c4ee9fa71db35a38da4bfee46116c28753009c6145aea112a8a`

Desktop build PASS: `GOCACHE=/tmp/sdp-wci0-gocache go build -tags desktop -o /tmp/sdui-wci1-native-noether-check ./examples/collections/cmd/native` from SDL/go. Binary SHA256: `fc42424fcfc537f2a82a224ce2b41390d73e9d5376ecae34a3d0be16d3a9a4da`. This build includes the visible-height paging fix and shared gutter integration. No tests or builds remain running from this checkpoint.

## Independent review recovery blocker repair

Unfroze only collection metrics/control and tests for the confirmed long-provider-error readiness failure. Shared `fittedRow` now fits diagnostic/status text and caps its indentation to the admitted content viewport. Measurement and painting use the same fitted label/indent; runtime retains the complete bounded diagnostic. Retry/R stays visible even at minimum viewport width. Diagnostic rows remain horizontally anchored while data rows keep their full labels, depth and original overflow policy. No runtime/document.go changes.

New `collection_recovery_test.go` drives actual asynchronous completion through the host Post queue with default overflow-x=error, at root and a fourteen-level expanded branch. It verifies LoadError clears the request, full diagnostic retention, visible fitted recovery label, repeated R producing exactly one fresh load, accepted subsequent completion, and rejection of a wide data-row replacement. Minimum-width/deep-indent label test also passes. The first full host race run failed only the new assertion that depth12 must cap indentation: the shorter Retry/R label legitimately fit at that depth. The fixture now uses depth14 so the test actually exercises capping; the scoped recovery race tests pass (2.635s). Full host race rerun is running at this note. Desktop rebuild passed with the product fix.

Main reported native gutter PageDown passing; background gap-wheel observation remains under harness investigation and is not treated as a confirmed product finding.

Recovery repair final checkpoint: `GOCACHE=/tmp/sdp-wci0-gocache go test -race ./host/fynehost -count=1 -timeout 120s` — PASS (15.197s). Desktop checkpoint binary `/tmp/sdui-wci1-native-noether-check` SHA256 `2e71d97c44c7332f48c039c1e58553adff179584a149893f5ed5eb7708af3647`. No remaining worker test/build processes. Metrics/control/test lane refrozen; actual-native recovery rerun and independent review remain main-owned.

Recovery repair file hashes:
- `SDUI/go/host/fynehost/collection_metrics.go`: `217c277427d6adb641951ae2b3e946edd8e5191e4b798b36f366088144f1d3fb`
- `SDUI/go/host/fynehost/collection_control.go`: `d3aca0675f7a104f1506d33bc4a729846888f7d584d9b15be931599c8b849229`
- `SDUI/go/host/fynehost/collection_recovery_test.go`: `569e9abd521915fb7f010ae60570ebe877281a0305213263dce75d598341dae9`

## Residual recovery-label review correction

The recovery affordance is now structural: native fitting derives `Retry/R` or `Load/R` from runtime status, reserves that complete prefix first, and truncates only the separate diagnostic string. Neither ellipsis nor deep indentation can consume the R hint. The same fitted result is used for measurement and drawing; data rows and runtime diagnostics remain unchanged. Minimum-width regression now covers fonts 14/20/28, depth0/depth100, and error/canceled phases, including an old verbose row.Text to prove the fixed prefix does not depend on parsing its display string. Full host race/build results follow below.

Main reports `/tmp/sdui-wci1-native-gutter-paced` passed all 20 native nested checks on the `fc42424fcfc537f2a82a224ce2b41390d73e9d5376ecae34a3d0be16d3a9a4da` predecessor. Motion followed immediately by wheel used Fyne's cached pointer; the external helper now explicitly moves and settles 0.1s before wheel. This is a corrected input harness, not a product gap-wheel defect. The later recovery-label candidate still requires its affected native rerun.

Final residual-fix verification: `GOCACHE=/tmp/sdp-wci0-gocache go test -race ./host/fynehost -count=1 -timeout 120s` PASS (15.721s). Desktop build command unchanged, PASS. Final binary `/tmp/sdui-wci1-native-noether-check` SHA256 `d97eac538dbb80816da863e07e3ec4aa4bb6ef0416b62c2e37211591a8badbd8`. Metrics/control/tests held; no tests/builds running. Runtime/document.go untouched. Native affected rerun and independent acceptance remain main-owned.

- `SDUI/go/host/fynehost/collection_metrics.go`: `86f03ce4c5ad5b37fd1857ed9b597d54ee85ffd3a6fcbf59d60d9b5da7005483`
- `SDUI/go/host/fynehost/collection_control.go`: `d3aca0675f7a104f1506d33bc4a729846888f7d584d9b15be931599c8b849229`
- `SDUI/go/host/fynehost/collection_recovery_test.go`: `6834572f0525e8dba7607fbbc28056b519e45e6cee84b7023ab581999fcdcbca`
