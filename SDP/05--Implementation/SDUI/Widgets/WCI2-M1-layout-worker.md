# WCI2-M1 layout worker handoff

## Authority and scope

Coordinator selected tabs/page/split only after WCI1 phase c39b330/evidence 53d031e
and independent review of the Panes-and-commands.md substantive contract cf93dea6.
Worked solely in /tmp/sdp-sdui-widgets on sdui/widgets-wci2, preserving concurrent
lanes. No commits, branch changes, host/parser/runtime/bridge/fixture edits or
management writes. WCI1 fixture ownership is frozen and Lorentz owns the panes
fixture. M2 command/menu/dialog implementation was not started.

Loaded/reused SDP 1.1.1, SDP Worker 2.0.0 and document-workflow. Recovered original
Session0010 S3 and the clone's current reviewed Panes contract, temporary runtime
seams analysis, WCI2-runtime-API.md and WCI2-frontend-API.md. Coordinator owns
Session, plan, native evidence and final integration. Accidental-interruption
recovery preserved all files; runtime's transient missing-type build dependency
cleared once James published its implementations. No competing runtime types or
stubs were added.

## Delivered behavior and coordination

The concrete API memo is SDUI/docs/wci2-layout-api.md. Sent API proposals to
Noether/James/main before dependent implementation and consumed frontend's
PaneChildren and runtime's TabsState/SplitState/PresentationState. Noether's current
pane_metrics.go uses the exact agreed methods; inspected it read-only. Raised the
visible-disabled split gate discrepancy; James corrected the expected expanded
set to visible regardless of enabled, consistent with retained painting geometry.

- Optional PaneMeasurer supplies actual TabsMetrics.Header minimum width/fixed
  height and split divider thickness. Missing or invalid metrics reject; no native
  chrome fallback or source pixels.
- SnapshotLayout.Tabs/Splits expose screen header/body/divider/child rectangles and
  effective chrome clips. Only active pages have body boxes; collapsed children
  have zero extent and leave measurement/layout/hit testing. Ancestor content
  clips/window constrain chrome; child clips exclude chrome.
- Expanded split geometry uses divider-excluded U, native recursive minima and
  minFirst/minSecond fractions. Layout returns runtime.SplitGeometry legal bounds
  and effective clamp. Both relative minima are suspended while explicitly
  collapsed; visible native minimum still applies. Runtime owns saved ratio.
- Recursive intrinsic minima use the subtree's own finite body. This prevents
  descendant min-x=.8/.99/1 from acquiring the full outer split reference.
  Nested padding and horizontal/vertical fr minima are included, with bounded
  refinement and per-run cache. Collection probes use native Minimum, not loaded
  extent; assigned collection viewport metrics remain fully validated.
- PresentationState() returns detached active viewport offsets and only visible
  expanded split geometry. Runtime overlays retained inactive pane offsets and
  enforces strict versus clamped operations. Layout never changes selection,
  proportions, saved ratio, focus or accepted runtime state. Final native resource
  publication belongs to the runtime/host preparation-ticket contract.
- Existing WCI1 viewport/gutter/collection behavior and exact-source 0.2 geometry
  remain covered. No M2 renderer, callbacks, command/menu/dialog code or general
  scene model was introduced.

## Verification and limitations

Environment: Go go1.27.1 linux/amd64. Commands from SDUI/go:

- go test ./layout -count=1: PASS after each concrete M1 geometry change.
- Final go test -race ./layout ./svg ./prototype -count=1: PASS; layout 1.270s,
  svg 1.030s, prototype 1.044s.
- git diff --check for assigned layout/doc paths: PASS.

Exact geometry oracles cover native header/body separation, active-only measurement
with a deliberately failing inactive measurer, empty selection, split bounds and
first/divider/second rectangles, explicit positive-min collapse, visible minimum
rejection, nested tabs+split relative minima, descendant relative minima, nested
padding, both-axis 1fr:3fr minima, collection native-minimum/extent separation,
invalid/missing metrics and exact-source 0.2 geometry equality. Clipped header
and nested EnsureVisible verify ancestor translation/reveal without mutation.

Real runtime PresentationGate integration tests prove inactive offsets retained
while absent from active geometry, clamping on reveal after resize, strict ratio
rejection without snapshot change, collapse, collapsed resize preserving saved
ratio, failed restore preserving collapsed state/focus, successful restore clamp,
and visible-disabled split geometry without input hits. Gate results are detached.

The measurements in layout tests are synthetic adapter oracles. They do not prove
native Fyne header paint, pointer/keyboard selection, divider drag, focused chrome,
actual minima or final resource-publication behavior. No native binary was rebuilt
or display exercised here. Main/Noether own these native checks and integration;
independent candidate review and WCI2-M1 acceptance remain pending. Passing this
lane does not authorize or establish WCI2-M2 support.

## Exact uncommitted candidate

Base HEAD 53d031e9a9a6e10bdcf56a2ac7b6e863adc81a9d, branch sdui/widgets-wci2.
These hashes identify the final assigned source/doc bytes, not an unchanged-commit
result. Other lanes have concurrent changes under their exclusive ownership.

| File | SHA-256 |
| --- | --- |
| SDUI/go/layout/engine.go | 19dd480608dcc20452c20892fc3f45750cd0d853fd95df3880b541693fb4fc61 |
| SDUI/go/layout/types.go | e2e0f3d1587665494056ec0334450cedcd5db7079307fb47b725e0e2ba2c545b |
| SDUI/go/layout/viewport.go | 1ac8be06ac37ec231f0fa888ad45e344a5a8761b072155093079f585acf0c5dc |
| SDUI/go/layout/panes.go | fbb728df6a146d20b12dd976cbf2c18c6db1f8f6e0497537591c8e1e7c0a27cf |
| SDUI/go/layout/pane_minimum.go | e4122b787ada52402b43d0441bffbc6bc58f49a32491e2c13fab1132c60d9603 |
| SDUI/go/layout/panes_test.go | 5c867ae28f0e3c6573123fae79a92b3c99f340e255ebcfdd19635e882ad4d6ad |
| SDUI/go/layout/pane_gate_test.go | 9d4e8f8011c5b53ca20c4005f438f3ce42b4167ea9eeb2e13b8b4d5df76d122b |
| SDUI/docs/go-layout-contract.md | 88c291588eecba46e601abcae9aa350251d69e7a298f67b4af6c155c93219adb |
| SDUI/docs/wci2-layout-api.md | 049259553ea332ac5ea2001dba35b889be2b10739d48acf18f8074fcb7e9ff82 |

## Coordinator Session0010 S3 handoff

Owner input: resume the same exclusive M1 lane after accidental interruption;
preserve WCI1/0.2, no M2 or fixture work. Worker recovered files/memos, coordinated
frontend/runtime/host seams, delivered the bounded geometry and English contract,
corrected nested relative minimum references, and ran the scoped evidence above.
The API memo/results are ready for Noether/James/main. Next step: integrate exact
candidate, inspect actual native metrics and run tab/divider/reveal/failed-resize
acceptance, followed by independent candidate review. Coordinator updates Session,
plan and acceptance records; this worker report is not a management transition.
