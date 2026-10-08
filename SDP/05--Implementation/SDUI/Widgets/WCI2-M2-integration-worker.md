# WCI2-M2 host/preparation implementation handoff

Status: product code refrozen after the bounded parent-surface and WM-owner-lifecycle corrections. The refreshed artifact below is ready for independent review and coordinator-owned exact-candidate native acceptance; earlier artifacts are historical. No independent approval or complete slice acceptance is claimed here.

## Scope and identity

Worked only in `/tmp/sdp-sdui-widgets`, branch `sdui/widgets-wci2`, current shared HEAD `a4f2c22435d91fe07935b8b9d5fcf46fafc6e36b`. The tested candidate includes uncommitted shared-lane changes; it is not an unchanged-HEAD result. M1 acceptance baseline was `403c5402c89532e7bf3130776b27f2300e7612fb`.

Owned host/fynehost, preparation, related documentation/tests and this report. Explicit delegated contributions from Gibbs: inspection.go, inspection_m2_test.go and command_admission_test.go. Other lanes own runtime, parser/layout/export, SDL bridge/fixture and module files. No original workspace edits, commits, branch changes, PM/Session writes, or WCI3 changes were made by this worker. Reused SDP Worker 2.0.0/document-workflow, applicable AGENTS, canonical original Panes-and-commands M2 contract, and runtime/layout/frontend/fixture API memos.

Current final desktop artifact: `/tmp/wci2-m2-noether-wm-final`

SHA-256: `2c1db6d173c67ff873da4a7fdadb2a32722d3153705523622328e0d509be7806`

Previous parent-correction artifact `/tmp/wci2-m2-noether-parent-final` SHA-256 `1f9b81af5470e9360d977137abcf97da907d839dd262bebe45b8bd2476042b21` remains unchanged as predecessor evidence.

Previous `/tmp/wci2-m2-noether-final` SHA-256 `2253938916e1a31dfa4ec6d6095daf9f82a7522abcd444833cb64977e69b1c33` remains unchanged as predecessor evidence.

Lorentz's included A06 plus nested-parent fixture source SHA-256: `d23b373149de45b082603673a9e975e72ed01a0d846fa4c4b055266b18883fb2`. Earlier pilot binaries remain untouched. Main owns the exact pinned `fancyfs v0.0.1` indirect module/checksum addition needed by Fyne dialog; this worker made no dependency changes.

## Implementation

- Detached readiness checks hidden/closed command and icon declarations, native key aliases/reserved forms, Markdown using the actual nil-renderer baseline, and connected command bindings. Prototype unbound counting includes ordinary commands without effects/toggle/handlers. Icon resources own their bytes while retaining the public ThemedResource color contract; flattening that contract had baked dark-theme white into light native controls and is corrected.
- DocumentHost advertises concrete M2 capabilities. Native Fyne Button adapters project accepted command labels/icons/check state, preserve plain/decorated legacy Activate dispatch, and publish actual focus through guarded runtime Focus/EnsureVisible. Tooltips use native container minimum size before available-canvas truncation.
- A single command capture/DispatchInteraction route serves toolbar, menu and key. Native Entry editing and collection/divider navigation have priority. Public Shortcutable forwarding handles commands swallowed by Fyne's focused-control dispatch; public Keyable handles Shift-only/function keys and Menu. Named Copy/Cut/etc. shortcut values normalize outside Entry editing. Extra-modifier custom shortcuts remain commands. Shift+F10 is explicitly reserved for context. Lifecycle-owned no-focus canvas callbacks restore the composition owner's callbacks on disposal; no duplicate focused dispatch.
- PopUpMenu keeps native private items untouched, with a cloned menu model and a same-tree public pointer/keyboard forwarding wrapper. Its synchronous per-opening selection scope retains capture only across dismiss-before-action. Hide/removal happens before Action; a claim executes at most once, scope exit revokes unclaimed capture, and Escape/outside revoke immediately. Late old dismissal cannot remove a replacement overlay. No timers, fyne.Do cleanup heuristic, reflection, unsafe, private fields, dependency patches or general routing framework.
- Context targets use captured runtime handles/items without changing selection. Native controls include bounded secondary-click targeting for the existing SVG placeholder; this adds no SVG execution or renderer framework.
- Accepted presentation tickets prepare each surface background with subtree-scoped native inventory and the full finalized InteractionRoot. Modal CustomWithoutButtons content geometry is distinct from measured native chrome, including a prospective parent-size check. Nonmodal app.NewWindow uses measured outer chrome with no implicit window padding; accepted body size remains shared geometry authority. Empty dialogs have one native dialog focus stop, without fabricating a model widget.
- Native Show precedes ConfirmSurfacePublication. Visible fixed status chrome preserves false/unknown/post-domain-conflict outcomes and AcceptBlocked; full bounded text stays in runtime. Explicit child gestures may request OS focus; background/parent disposal does not. Parent hidden/closed hooks, successful reload, Close, and direct native Window.Close revoke exact published generations and drain once. OnClosed uses the runtime's exact RevokeSurface fallback without a fallible gate. Observer panic is isolated; reentrant observer replacement cannot trigger an obsolete OnChange.
- Inspect reports actual native controls/menu rows/surfaces and canvas-local geometry/title, preserving M1 nativeSelected. Main canvas ID is the literal `"main"`; nonmodal surfaces use their path; a modal reports its actual parent canvas ID. It is diagnostic geometry, not paint-ready acknowledgement.

## Verification

All commands used `GOCACHE=/tmp/sdp-wci0-gocache` and `-mod=readonly`; no network/sandbox bypass occurred.

From SDUI/go:

1. `go test -mod=readonly -race ./host/fynehost ./preparation -count=1 -timeout 120s` — PASS host 65.171s, preparation 1.045s. This broad integration run preceded the final narrow prospective-chrome/keyboard/tooltip corrections; it is not claimed as the final whole-suite hash result.
2. `go test -mod=readonly -race ./host/fynehost -run 'TestModalChrome|TestEmptyDialog|TestDeclaredShortcut|TestPreparedThemeIcon|TestDialog|TestCommandAdmission' -count=1 -timeout 90s` — PASS 19.019s after prospective chrome validation.
3. `go test -mod=readonly -overlay /tmp/wci2-key-review-mtgtj646/overlay.json ./host/fynehost -run TestReview -count=1 -timeout 60s` — all seven independent reviewer key repros PASS 0.200s: header modifier reset, collection Shift+K, native semantic Copy, extra modifiers in Entry, empty-dialog function key, reserved Shift+F10 and no-focused-control main-canvas F2.
4. Final owned regressions: `go test -mod=readonly -race ./host/fynehost -run 'TestDeclaredKeys|TestTooltip|TestNativeCommandKey|TestEmptyDialog|TestDeclaredShortcut|TestCommandAdmission|TestModalChrome' -count=1 -timeout 90s` — PASS 12.782s on final product bytes. Persistent tests cover these key/lifecycle cases, native tooltip width, actual modal/nonmodal body agreement, native close after failing resources, once-only result callbacks, missing hidden icons/capability/binding, immutable icon bytes, late admission rejection and retained predecessor native shortcut.
5. Native menu tests (included in the broad run) cover nested pointer/Enter/Space, Escape/outside, longer-than-root nested Home/End, caller-model reuse, disabled keyboard claim refusal and callback opening replacement surviving duplicate old dismissal.

From SDL/go:

`go build -mod=readonly -tags desktop -o /tmp/wci2-m2-noether-final ./examples/commands/cmd/native` — PASS for the predecessor artifact, including Lorentz's then-frozen three-level-Escape fixture.

Final `go vet -mod=readonly ./host/fynehost ./preparation` — PASS.
Final SDL `go test -mod=readonly ./examples/commands ./bridge -count=1` — PASS commands 3.068s, bridge 0.379s.
Independent reviewer Mendel separately reports all seven key overlay repros PASS under race (3.959s), plus persistent exact-editing/canvas-lifecycle/key-boundary regressions PASS (1.854s); this closes that bounded delta, not the full host/native gate. Coordinator owns the full SDUI/SDL integration rerun and independent approval.

## Native evidence boundary and limitations

Main reported provisional actual XTest passes on predecessor binaries for shared commands, context and focused-tree key forwarding, modal/nonmodal input and focus, failure classification, lifecycle, nested menu keys and visible outcome/status. The fc1dfe pilot also showed the corrected icon. These are coordinator observations, not this worker's final-artifact acceptance claim. Final artifact reruns, OS screenshots and the mandatory menu → dirty draft revert → clean Cancel workflow remain main-owned.

Fyne app.NewWindow has no native transient-parent ownership and RequestFocus may be ignored on Wayland. Applications must call the documented parent lifecycle hooks. Actual OS paint can lag Inspect state; Canvas.Capture has previously returned black, so native visual evidence must use OS capture after paint. Native test-driver repeated Window.Close is not idempotent; the regression closes the real native test window once and repeats the captured fallback method to verify duplicate callback safety. No portability claim beyond the tested native environment is made.

The original handoff was reopened for one independently reproduced parent-surface blocker. The bounded correction and evidence are below. No other product scope was reopened. Main remains responsible for slice acceptance and integration records.

## Bounded correction: actual surface parent topology

Independent review reproduced a root-sibling modal `page/x` opened from nonmodal `page/longparent`: sorting source path lengths prepared the child first, attaching it to the main canvas and deriving 240x150 from main instead of 144x90 from its 480x300 parent. Predecessor artifact `225393...` and its native results remain historical evidence; they do not close this case.

Only product file `document_surfaces.go` and new `surface_parent_test.go` changed in this correction. `orderedSurfacePaths` now follows exact open runtime SurfaceTarget/ParentSurface edges with deterministic parent-before-child traversal. Missing/stale parents and cycles reject. Size calculation requires the prepared parent size; modal creation requires the prepared exact parent canvas. Publication uses the same order. No source-path ownership inference or missing-parent main-canvas fallback remains.

- Original independent overlay repro: `GOCACHE=/tmp/sdp-wci0-gocache go test -mod=readonly -race -overlay /tmp/wci2-final-review-jxqaw7a_/overlay.json ./host/fynehost -run TestReviewDynamic -count=1 -timeout 60s` — PASS 2.217s. Mendel separately reported race PASS 2.403s.
- Persistent and surrounding surface regressions: `go test -mod=readonly -race ./host/fynehost -run 'TestDynamicSurface|TestModalChrome|TestEmptyDialog|TestDialog|TestInspectionM2Surface' -count=1 -timeout 90s` — PASS 15.220s. The new test checks actual overlay canvas, parent-relative content size, stale-parent rejection and cycle rejection.
- Product SHA-256 `document_surfaces.go`: `74344547a6ea54140fa0795f31124c9dca20196f30af8e2c4612a2416a7c9351`.
- Test SHA-256 `surface_parent_test.go`: `95f20ec1175efd23762e8b86598e2bb47229014922e8827707cd79c90af098a9`.

Refreshed build: from SDL/go, `GOCACHE=/tmp/sdp-wci0-gocache go build -mod=readonly -tags desktop -o /tmp/wci2-m2-noether-parent-final ./examples/commands/cmd/native` — PASS. Artifact SHA-256 `1f9b81af5470e9360d977137abcf97da907d839dd262bebe45b8bd2476042b21`, with Lorentz's final nested fixture SHA above. `go vet -mod=readonly ./host/fynehost` also PASS after this correction. Lorentz separately reports actual fixture Request/Adopt modal/nonmodal race PASS 3.757s on those source bytes.

The refreshed artifact includes Lorentz's actual native root-sibling child fixture. Main owns child-on-parent OS overlay, parent-relative body, focus/Escape and parent-close-descendant acceptance. No native passing claim is inferred from the headless test.

## Bounded correction: native WM owner close

On parent-correction artifact `1f9b81...`, main's actual nested native test passed six owner/size/focus/Escape/pointer-Close checks, then WM_DELETE_WINDOW on the nonmodal parent was incorrectly blocked by the child's modal CaptureDialog guard. That is a host lifecycle defect; the harness obligation was retained.

Coordinator/architect explicitly selected and clarified the canonical contract: WM/decoration close is native owner loss. The native close interceptor revokes the exact opening and descendants via existing `RevokeSurface(target, "parent-closed")` before native teardown. Receipts carry reason `parent-closed` and sequence zero, retaining AcceptSequence/Domain. It never synthesizes a source interaction or calls Accept. Source Close/Cancel remain user interactions carrying their triggering sequence and retain modal guards. No runtime API/source changes were made for this correction.

`document_surfaces.go` now installs `nativeCloseRequested` directly as the nonmodal Close interceptor. It revokes before destruction, restores logical focus, reconciles native menus/providers and drains once. OS focus restoration is restricted to closing the currently active surface or its actual ancestor, not an unrelated nonmodal. Old callbacks cannot revoke a reopened replacement. The existing direct Window.Close fallback remains idempotent native-owner teardown. README now documents the receipt distinction and actual Inspect canvas IDs.

- Final targeted command from SDUI/go: `GOCACHE=/tmp/sdp-wci0-gocache go test -mod=readonly -race ./host/fynehost -run 'TestNativeOwnerClose|TestDynamicSurface|TestEmptyDialog|TestDialogPublication' -count=1 -timeout 90s` — PASS 11.202s, exit 0. Persistent regression includes modal child, failing resource gate (zero calls), no extra Accept, prior unknown domain preservation, both drafts discarded, two exact lifecycle receipts, duplicate callbacks, reopened replacement safety and unchanged unrelated nonmodal focus.
- `go vet -mod=readonly ./host/fynehost` — PASS, exit 0.
- Mendel independently reported bounded source approval and race PASS 5.535s on the WM product bytes. The last test-only addition extends this to background-owner close without stealing another root nonmodal's focus.
- Product SHA-256 `document_surfaces.go`: `9a96857577e8b58848d7eef31c19735b9f142b80141bd08c68c063ad896c20a8`.
- Persistent test SHA-256 `surface_parent_test.go`: `33c36f03af4139432ba4cc02999923f28b3d995f565a6865587be9847b4c21b9`.
- Build from SDL/go: `GOCACHE=/tmp/sdp-wci0-gocache go build -mod=readonly -tags desktop -o /tmp/wci2-m2-noether-wm-final ./examples/commands/cmd/native` — PASS, exit 0. Artifact SHA-256 `2c1db6d173c67ff873da4a7fdadb2a32722d3153705523622328e0d509be7806`. Fixture source remains `d23b373149de45b082603673a9e975e72ed01a0d846fa4c4b055266b18883fb2`.

Product code is refrozen. Main owns the final WM cascade and all native/integration reruns; no broad worker suite repetition was performed. Main reports early rejected tentative native-minimum sizes before the valid initial resize in fixture logs; these are recorded initialization transitions, not a claim of error-free logs or a newly introduced failure.

## Changed-file SHA-256 inventory

The following manifest includes all modified/new host and preparation files (including the explicitly delegated Gibbs evidence). Manifest SHA-256, over these exact `hash  path\n` lines: `6843c3b58ee960b3e2412dcb32b48e2ffb8aafa18966a456d57e8872dcfca7ff`.

```text
a52ac229bf1d6e3ea509b322cf9763162add5a33ac0c6990c61d772aba84a4ad  SDUI/go/host/fynehost/README.md
33f856cde1a6e2ff844305f40e3561da7061eb752ba206820abba171e78f46b4  SDUI/go/host/fynehost/admission/admission.go
1718b2ae8b2990d22ea0f0c285b5ff433d9f2b4f7c964bbed3c4ec5773334601  SDUI/go/host/fynehost/collection_control.go
d686f3fe8cb6b0d411cf9b5fcd79bf38639eb646a8ee054774faf5b2f76f5fb3  SDUI/go/host/fynehost/collection_metrics.go
7bba53d95efbba2d7744075153c76d63a1c76b62c989449a8e4bfe9e67cbb797  SDUI/go/host/fynehost/command_admission_test.go
f34a9f2253bd77a1fdf892e4b8e43471385549f399c0acd6ffb54b560ea36aa9  SDUI/go/host/fynehost/command_button.go
eed66cac075232dc8bb8af906eb323d7948acdae76ace2cf4195770a14ff1053  SDUI/go/host/fynehost/command_keys.go
f3a517fcf6c5732b4097bbd1e522f458772e9ab1fb6f3611819ee3f4e9d5a639  SDUI/go/host/fynehost/command_resources.go
84ecb5d72644a2da8086a40fb795aa1dd43546e2d707ac31e256e7bf368eea3e  SDUI/go/host/fynehost/commands_test.go
f2130d8457f1162fcbde05f5672fbb22bd7bf1315bb26de08b213ad4197e9f58  SDUI/go/host/fynehost/context_control.go
0c8512ea78580bb0870ecdaaded5627c2423c3c7aef13be12a92c46bf2a52a89  SDUI/go/host/fynehost/document.go
06cabe3e1c926621f4d5ea03282581e957585ef3cb6274445bfa9a5391cd9a5a  SDUI/go/host/fynehost/document_commands.go
5e76dc828b4e06d5606c484961fc54ababd8e6df0042f73cc56821039c507472  SDUI/go/host/fynehost/document_events.go
9a96857577e8b58848d7eef31c19735b9f142b80141bd08c68c063ad896c20a8  SDUI/go/host/fynehost/document_surfaces.go
7ec9418051fdbcf1024af30d6a0dac52f82d3f5c612ea84aaa7c9c508bdeee15  SDUI/go/host/fynehost/document_sync.go
c337660906f214d3ee5b36093ce27a639c5e919f2b719d64581df55acfa55623  SDUI/go/host/fynehost/input.go
8eaaecbea8876c893c64db9cfa08c011abcd0f8f4588f7bb71d1f9554ef8881b  SDUI/go/host/fynehost/inspection.go
981c958994820e2c8a05117136c77cdd0b79e017d3c39a22b1b656753a416547  SDUI/go/host/fynehost/inspection_m2_test.go
4189927bc6d80eef4bddca749f7ffb408b7d02e20d6547ec7609b7ade72223dd  SDUI/go/host/fynehost/native_menu.go
8b1f1673aaaf44feb1c0530806d772a1c87ad4b9b6c2998918490f459f10bd17  SDUI/go/host/fynehost/native_menu_test.go
74de6201bcbac84ae79c45d15f551c0d9454236dbc829f9480e88dbd0fa08276  SDUI/go/host/fynehost/pane_divider.go
256724db93264f3c57fae76f92c448585600a2ebab95348aab8529097678a644  SDUI/go/host/fynehost/pane_events.go
7f81085d4cb9ae572bc592d3e366e7c84a12c501aa9ef18d3b7f0ab2d4b83417  SDUI/go/host/fynehost/pane_header.go
2e631c1854776730f5b3a8e9441904682ec7b5b86e85618ea7d4a0b41cf9e0ab  SDUI/go/host/fynehost/pane_metrics.go
3b8f2cf391a69d307358513f2649e7536ec1f07c39364c39fc5265fc4b352ec8  SDUI/go/host/fynehost/surface_focus.go
33c36f03af4139432ba4cc02999923f28b3d995f565a6865587be9847b4c21b9  SDUI/go/host/fynehost/surface_parent_test.go
28c0b218bf363b63c5cd5dfcf2c9551660787100c32e19403477caaef194642c  SDUI/go/host/fynehost/tooltip.go
105d863e7856823d4982182dbbfb0dba37077db9228a993098d544fce9a5eba2  SDUI/go/host/fynehost/view.go
103162ec62036104ef4ca95ff815933876cb225a6f81164d8052dedf7b7d5888  SDUI/go/host/fynehost/viewport_control.go
5da51d9c77502627582a8e521c1ddce459c94fdfe755f7bcc4867bdb9baaf867  SDUI/go/preparation/capabilities.go
78f3a4719512bf9b12070df62409e21dd145d4e30416483c9f1d22e55f05d03f  SDUI/go/preparation/prepare.go
```
