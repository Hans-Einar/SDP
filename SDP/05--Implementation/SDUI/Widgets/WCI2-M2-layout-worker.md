# WCI2-M2 layout worker handoff

## Authorized outcome and boundaries

Coordinator selected M2 after independent M1 approval (403c540, 58 native checks)
and the final reviewed M2 contract substantive hash 9a9c8710. Read canonical
original SDP/04--Design/SDUI/Widgets/Panes-and-commands.md including the latest M2
transition, the clone WCI2-M2-runtime-API.md and WCI2-M2-frontend-API.md. Reused
SDP 1.1.1, Worker 2.0.0 and document-workflow; recovered original Session0010 S3
read-only. Main retains Session/plan/integration/native ownership.

Worked only in the isolated clone, branch sdui/widgets-wci2: layout code, related
English docs/tests and this report. No host/runtime/parser/SVG/product fixture,
management or Session edits. No commits or branch changes. No WCI3 implementation.
Concurrent exclusive-lane files remain owned by their workers.

## Frozen seams and implementation

Noether agreed the bounded API before dependent code; Dalton published and
implemented parser.IsAuxiliary before dependent integration. James froze existing
PresentationState/gate/ticket types and the concrete Snapshot.Surfaces map. The
exact layout API is SDUI/docs/wci2-m2-layout-api.md:

- MenuMeasurer.MeasureMenu(instance,font) returns actual native bar minimum.
  Native bar controls are source-positioned leaf boxes, with no duplicate popup
  row tree. Missing, nonfinite, nonpositive or excessive metrics reject, and final
  assigned size must meet the native minimum even after source max bounds.
- Frontend IsAuxiliary excludes nonvisual commands/items/separators/groups,
  context/submenu menus and dialog declarations from enclosing rows, gaps, minima
  and relative dependencies. This includes open dialogs; source declarations stay
  intact for other lanes' full preflight. Direct auxiliary split children reject
  through frontend validation. Pane dispatch remains limited to tabs/page/split.
- SurfaceSize(snapshot,path,reference) resolves initial dialog sizing against the
  supplied parent content/host reference. Empty natural dimensions may be zero;
  Noether confirmed host combines them with its concrete content-adapter minimum.
  Native title/window chrome remains host-owned and no constants are guessed.
- LayoutCanvases(snapshot,mainSize,acceptedSurfaceSizes) returns CanvasLayout with
  independent Main/Surfaces geometry. Open surfaces require positive finite sizes
  bounded by 32768; missing/unknown paths reject, closed cache entries are ignored.
  Each root fills its exact accepted local canvas at (0,0). Opening root scale,
  fill, ratio and min/max are applied once, never reapplied on actual native resize.
  Local padding, ordinary child layout, clipping and minima remain enforced.
- Nested dialog content is excluded from parent canvas extent and hit routing.
  Source fonts inherit, but coordinate clips and viewport ancestry restart at each
  canvas. Existing nonmodal sizes are explicit and unaffected by parent resize.
  The whole measurement keeps the common operation budget and rejects duplicate
  normalized paths or geometry crossing canvas boundaries.
- CanvasLayout.PresentationState returns fresh aggregate active viewport/split
  maps using runtime types. Any failed canvas returns no aggregate result.
  Runtime owns retained closed offsets, modality, publication and final ticket;
  layout adds no native objects, handlers, window state or second offset authority.

Initial root max bounds also must not be reapplied: max-x=.8 resolves 320 against
parent width 400, then the accepted local root stays 320 rather than shrinking to
256. This clarification was sent to Noether/main and documented. Zero natural
content size handling was explicitly agreed with Noether; accepted canvas remains
positive. No unreviewed workflow departure was selected.

Native background integration uses Dalton's svg.Options.InteractionRoot. Each
surface Box.Instance is the exact dialog pointer inside the finalized snapshot,
not a synthetic frame or clone. Passing that same snapshot root preserves command
references outside the dialog subtree. Scoped prepared inventory still covers
closed nested adapters; native omission does not authorize public M2 SVG export.

## Verification

Go go1.27.1 linux/amd64, commands from SDUI/go unless stated otherwise:

- go test ./layout -count=1: PASS, including M1/WCI1 and new M2 cases.
- Final go test -race ./layout ./svg ./prototype -count=1: PASS; layout 1.310s,
  svg 1.050s, prototype 1.054s.
- git diff --check on assigned layout/docs: PASS; gofmt -l layout: empty.

Concrete oracles:

- Auxiliary-only rows create no gaps/tracks, closed poisoned content is not
  measured, menu bars use native leaf metrics and popup descendants have no boxes.
- Nested surfaces have separate origins, inherited font, local scrolling ancestry
  and correct aggregate split/offset maps; nested content is absent from parent.
- Parent resize leaves accepted nonmodal geometry unchanged; actual child resize
  uses the new local size without repeated source scale or resizing another canvas.
- Missing/unknown/nonfinite/negative/zero/excessive or below-minimum accepted
  surface sizes return nil geometry and leave the supplied snapshot unchanged.
- Empty natural content is allowed only before host allocation; zero accepted
  canvas rejects. Opening max bounds apply once, closed size caches are ignored.
- Missing/bad bar metrics and post-bound native-minimum violations reject.
  Exact-source 0.2 geometry equals the legacy entry point with no menu calls.
- A real runtime PresentationGate and PresentationPrepare ticket test rejects an
  invalid open-surface resize before preparation/publication, preserving the
  complete live snapshot and accepted geometry. Successful parent resize retains
  child size and finalized runtime split proportions equal prepared geometry.
- External-package layout/SVG background test accepts outside-subtree command
  references with exact selected snapshot, renders active body only, and rejects
  absent/foreign snapshot context or missing closed nested adapter with no artifact.

Temporary runtime compile gaps while James added actual M2 methods were not
worked around with stubs. Once implementations arrived, all above checks passed.
These are pure/synthetic metric and runtime ticket/background checks. They do not
prove Fyne menu dismissal order, native title metrics, pointer/keyboard actions,
modal blocking, window/OS focus, SDL execution or terminal receipts. No native
binary/display was built or exercised in this lane. Main/Noether/bridge own the
remaining actual native/SDL proof and independent final candidate review.

## Exact scoped candidate

Tested uncommitted clone at base HEAD a4f2c22435d91fe07935b8b9d5fcf46fafc6e36b,
branch sdui/widgets-wci2. This is not an unchanged-commit test result. Other lanes
have concurrent uncommitted dependencies; main captures the final aggregate.

| File | SHA-256 |
| --- | --- |
| SDUI/go/layout/auxiliary.go | 3839c558816db40928bbc3fef75c69e170e21af118a4d131f97447abe9030ed0 |
| SDUI/go/layout/canvases.go | 47a6ce48122125ff6492308689f7405553199d83027ee03b1873a99178ee7db9 |
| SDUI/go/layout/canvases_test.go | 7d2c3ac356675c2892445e1b6a993b56cea2fe1eaab08a5028d099521b333f96 |
| SDUI/go/layout/canvas_gate_test.go | ee53c6eb821426a84ece91db78ad6cced196ed44860b36ccde98763655e3b80d |
| SDUI/go/layout/canvas_background_test.go | 9e7aec18bc69f51e55230243e916b7c598992923e5979d1cb4b6df160d9af5e8 |
| SDUI/go/layout/contents.go | 1b0a1be2434e0c8574ee1191e916a2ba8f2b1f1c76b5b7415cf5808bc106a8ef |
| SDUI/go/layout/engine.go | 5522bcb38fa08ca53257739d5d1cfb0ea89411bd0762ba8f216752beb4720384 |
| SDUI/go/layout/pane_minimum.go | 44a0723a58d34b0599947d8766997be5f1433e55904f30c0bee41ac361f1c212 |
| SDUI/go/layout/panes.go | d233d9501aa0266de349265f4ffaf02ad48a55f4ca93a4a753bfa6163efe48c5 |
| SDUI/go/layout/types.go | 4f03529f75d341a29e00d9a97e31b39f3d380a0a8d0b1dcc39e0b9a8c11123b7 |
| SDUI/go/layout/viewport.go | dac32c82826b94a1e2c19d8f8d7398556e21c3d09a73613b200bb79d3dda7192 |
| SDUI/docs/go-layout-contract.md | 4f0039719bd7755611574c96472b9416ebfcf04061ff0a7da54247db1f99dccf |
| SDUI/docs/wci2-m2-layout-api.md | 15d8b348c6369543dce36f64a471e133302d16469a5e1d9c0aeeaed5d738411a |

## Coordinator Session0010 S3 handoff

Owner input selected M2 only, confirmed the reviewed CanvasLayout seam and asked
for adversarial background/layout rejection before publication plus existing
nonmodal reference tests. Worker recovered authority/memos, froze host/frontend
interfaces, delivered scoped geometry/docs, coordinated empty minima and SVG
snapshot provenance, and produced the passing evidence above. Loaded routines are
SDP 1.1.1, Worker 2.0.0 and document-workflow. Next: native host integration and
actual M2 SDL/window/menu acceptance, then independent exact-candidate review.
Main maintains Session/plan/records; this report does not close M2 or the full card.
