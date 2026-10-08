# WCI4 layout/resource API preparation

Status: read-only preparation plus this root memo only, 2026-10-08. WCI3-M2
layout is approved/frozen; WCI4 implementation is not selected. No product code,
layout docs, management, module, dependency or prior evidence files were changed.
Main owns selection, Session/PM/traceability, integration and actual native proof.
SDP 1.1.1, Worker 2.0.0 and the shared document workflow were reused; Session0010
S4 recovered read-only.

Inspected `/tmp/sdp-sdui-widgets`, branch `sdui/widgets-wci3`, HEAD
`d6742742f0968661d15698518df688ffb1d28946` plus concurrent WCI3-M2 working changes.
Authority: ORIGINAL `SDP/04--Design/SDUI/Widgets/Providers-and-packaging.md`, SHA-256
`83d54b218c7c2dfbebff17b641c2fdc67e4ca7d1f5fb73fe65851de5bee8ad46`, especially
§§2–4 and WCI4-A01–A05/A08. This is inspection/proposal evidence, not tested
WCI4 behavior or a change to that reviewed contract.

## 1. Actual paths and concrete gaps

* Layout.Measurer takes an Instance, font and available width. Generic widget
  measurement returns the same Size as both preferred size and minimum. SVG
  currently uses a 160x90 placeholder from TextMetrics. Source-assigned widget
  dimensions below that returned minimum reject in Engine.desired.
* markdown.Provider implements both Measure and svg.ContentRenderer.Render. Its
  Documents map is keyed by source text and Resources by diagram hash. Prepare
  optionally invokes a Renderer, stores its returned Resource without a defensive
  copy or independent custom-output validation, and shares equal text/diagram IDs.
  These public legacy caches cannot represent different per-instance policies.
* Provider.items places rendered diagrams at the left edge with width
  min(availableWidth, intrinsicWidth), preserving ratio through derived height.
  Missing legacy diagrams use the existing placeholder. Measure and Render repeat
  this pure item calculation; they do not call the renderer or open files.
* DocumentHost.measureCanvases calls markdown.Prepare(root,nil) inside every pure
  geometry probe. This reparses legacy Markdown but does not invoke a renderer.
  Merely replacing nil with a WCI4 renderer would execute it from StateGate and
  repeated resize/ticket work, directly violating the reviewed contract.
* svg.Render delegates only Markdown to Options.Content. SVG widgets always paint
  the old placeholder; SkipControls explicitly excludes svg in both check and
  render. Native resource SVG therefore needs a narrowly opted-in inventory path.
* Layout.Box already supplies one outer Rect, inherited Font and effective Clip
  for each active node. The clip includes ancestor scroll translations/content
  gutters and the owning canvas. No resource-specific hit tree is necessary.
* PresentationState contains outer offsets and split bounds only. Current host
  preparation tickets consume finalized snapshots and Commit invokes the existing
  application Guard. Resources belong to that bundle/ticket lifecycle, not runtime.

Two important consequences: current natural-as-minimum SVG sizing blocks fitting a
wide image into a shorter assigned box, and text-keyed legacy caches cannot decide
explicit fallback or readiness for two equal Markdown texts at different paths.
Neither should be fixed by calling providers during measurement.

## 2. Coordinated ownership and source facts

Dalton proposes, not yet implements:

```go
// parser: source facts only
func PreviewOptions(*Instance) (PreviewPolicy, error)
type PreviewPolicy struct {
    Explicit bool
    Description, Fallback string
}
```

It applies only to normalized SVG/Markdown. SVG description/fallback presence
requires both; bare Markdown has empty Arguments and Explicit=false. Explicit
Markdown remains Kind markdown, empty Widget, Text matching its retained text
argument, with description/fallback arguments. Source references stay symbolic;
Normalize checks declared modules. Layout must not infer opt-in from descriptions,
resource bindings, profile alone or provider presence.

James proposes immutable `markdown.Previews` in the existing markdown package,
which already owns Resource and the parser/layout/SVG dependency direction.
PreparedSVG and MarkdownRenderer remain the canonical application input shapes.
Previews owns a private path map; copied outcome access and pure Measure/Render
consume frozen SVG or explicit-Markdown outcomes. Each outcome retains identity,
description/policy/status, resource digest/dimensions, and ordered parsed blocks
with per-diagram resource or label result. No runtime preview state is proposed.
James additionally proposes `preparation.Request.Previews *markdown.Previews` as
an already-prepared input, validated against the separately normalized selected
root before the first gate. A source/policy/text fingerprint must match; mutable
caption/layout are not resource identity. This does not authorize preparation to
invoke a renderer. The exact outcome/accessor names remain James's API to freeze
at selection; layout
must not create a parallel resource record or import markdown (which imports layout).

All bindings, including hidden pages and closed dialogs, are resolved and validated
before the first presentation gate. Provider preparation copies bytes/documents,
validates identity before content, resolves label/reject once and checks aggregate
budgets. Shared immutable bytes may be deduplicated; policy/outcome records may not
be shared by equal text alone. A missing prepared path is an admission error, not
permission for Measure to render, resolve a resource or choose fallback again.

## 3. Smallest proposed geometry seam

Reuse Engine.Measure and the existing provider composition. Add only a bounded
optional geometry adjunct for explicitly opted-in SVG, because the current
Measurer does not distinguish preferred size from minimum:

```go
// layout: proposal, not an existing export
// Sizes contain no bytes, provider objects, policy or runtime state.
type PreviewMeasurer interface {
    MeasurePreview(*parser.Instance, float64, float64) (natural Size, minimum Size, err error)
}
func FitPreview(content Rect, intrinsic Size) (Rect, error)
```

Arguments mirror existing Measure: exact node, inherited font, finite available
width. The provider looks up the already frozen outcome by normalized path and
returns pure dimensions. Host's existing metrics wrapper delegates this adjunct
alongside its field/collection/pane metrics; absence for an active explicit SVG
must reject rather than accidentally report placeholder readiness.

For a rendered image, natural size comes from validated prepared dimensions,
constrained by available width while preserving aspect. Its image minimum is zero:
it is scalable content, not fixed native widget chrome. Source minima and any
measured caption still apply. For a labelled outcome, MeasurePreview measures the
actual English status text and optional caption, not rejected SVG dimensions or
an arbitrary 160x90 placeholder. Minimum should preserve the measured text region
required by the selected adapter; if that cannot fit, reject before publication
rather than suppressing status behind a caption. Exact text minimum policy needs
host/provider agreement at selection, using their existing text metrics.

Use this natural/minimum split in desired widget size and pane-minimum calculation;
otherwise a change in only one path leaves relative split layouts inconsistent.
Zero image minimum deliberately invents no pixel floor: a wholly unassigned pane
containing only scalable imagery has no intrinsic readability minimum. Confirm
source allocation/minimum expectations in the selection review; do not silently
reuse full intrinsic image dimensions as native minima and defeat shrink-to-fit.
An alternative of globally treating every SVG as zero-minimum is rejected because
it changes legacy placeholders and cannot preserve a prepared fallback's text needs.

Explicit Markdown keeps existing Measure/Content.Render signatures; select its
frozen record by path, not Text. Its ordered pure block-layout routine uses the
same font, width, status strings and prepared dimensions for measurement and
painting. A diagram failure replaces that block only; document-level unsupported
content replaces the whole explicit preview. The provider owns that distinction.
No block may reevaluate fallback while laying out, rendering a background or resizing.
The external `markdown.Renderer.Render(diagramSource)` prepares bytes once; the
existing `svg.ContentRenderer.Render(builder, box)` only serializes already prepared
content. They are distinct interfaces despite sharing the method name. Previews
must not retain external renderer objects as paint-time dependencies.

## 4. Aspect fit, caption, clipping and hit semantics

FitPreview operates only on validated dimensions and a final content allocation.
For intrinsic size iw,ih and available content x,y,w,h:

```text
s = min(w/iw, h/ih)
paint = (x + (w-iw*s)/2, y + (h-ih*s)/2, iw*s, ih*s)
paintClip = paint intersect Box.Clip
```

For content (10,20,100,20) and intrinsic (400,200), paint is (40,20,40,20).
Unused space is centered; the image is not stretched to 100x20. Explicit larger
allocations may upscale by the same ratio. Zero-area content has no paint area;
never divide by zero. Reject nonfinite/negative geometry and nonpositive intrinsic
dimensions; respect existing finite layout bounds. Resource dimension validation
and its 32768 limit remain preparation responsibilities, independently of this
pure arithmetic guard.

Fit against the full allocated content rectangle, never its already clipped visible
fragment. Otherwise scrolling an image would resize/recenter it. Intersect the
result with the existing Box.Clip afterward, without translating scroll offsets a
second time. This applies identically in main, nested clipped panes and independent
open-surface canvases. Caption, when present, consumes one measured band before
fitting. Label fallback uses measured text rather than pretending to be an image.

The provider's pure content placement routine must determine caption/status/block
bands consistently for Measure, background and native preparation. A native ticket
stores the resulting destination rectangles and paints them directly; no separate
host fitting or hit geometry. James's subsequent provider memo/direct handoff proposes the bounded accessor
needed to carry that single placement into native preparation:

```go
// package markdown; proposal on the immutable prepared provider
func (p *Previews) SVGRects(box *layout.Box) (image, caption, status layout.Rect, err error)
```

These three fixed slots are in the owning Box canvas coordinates. Optional absent
parts are exactly Rect{}; labelled fallback has no image. Image is already fitted
against the caption-excluded allocation once. Native preparation stores each
rectangle and its intersection with Box.Clip; it does not refit the clipped
fragment. MeasurePreview, this accessor and content serialization must use the same
pure text-band calculation and frozen path outcome. The accessor cannot resolve
resources, invoke external renderers or reconsider fallback. It adds no hit tree
or general scene and requires no runtime state. I agree with this bounded addition;
its exact implementation/native-font agreement still requires selection.

Noether's subsequent inspection identifies a concrete font mismatch between native
canvas.Text and provider layout.TextWidth, even with Go Regular. Host and layout
therefore agree on the smaller proposed text-paint accessor instead of introducing
a new PreviewTextMetrics callback or native line-layout protocol:

```go
// package markdown; pure serialization, not the external diagram Renderer
func (p *Previews) RenderSVGText(out *strings.Builder, box *layout.Box) error
```

It emits only caption/status glyph outlines using existing svg.Text. Inspection
of svg.Text and layout.TextWidth/Lines shows the same embedded Go Regular and
HintingNone; MeasurePreview, SVGRects and RenderSVGText must share one pure band/
line routine. The host mounts a direct text-only SVG overlay, clipped to Box.Clip,
separate from the directly rendered resource image. It contains no nested image or
data URI, so it does not reintroduce the unsupported composed-image route.

A caption may use one measured ellipsized row; status wraps the complete bounded
fallback text into its measured band and cannot be suppressed by the caption.
Reject exhausted text allocation before publication; no guessed font metrics or
second image fit. The native accessibility wrapper still exposes description/status
independently: SVG title/outlined glyphs alone do not prove native accessibility.
Actual native paint/wrap/clip and supported glyph evidence remain required; shared
font source is not a universal text-fidelity claim. This is a coordinated proposal,
not a selected WCI4 implementation or a change to frozen product files.

No extra public scene, SnapshotLayout resource map or runtime viewport is needed.
Native controls/background can consume these existing Rect/Clip values once.

Outer Box hit/context behavior remains authoritative; fitted artwork and unused
letterbox space do not create additional actions. No new source scrolling owner
is introduced for SVG or Markdown. Existing frame/group scroll extents and offsets
remain based on shared outer layout, never a second image camera or provider offset.

## 5. Preparation, background and publication ordering

1. Normalize and validate the full selected source. Validate binding paths/kinds,
   source triples, IDs/revisions/digests and budgets before content fallback.
2. In detached application/host preparation, prepare immutable per-instance outcomes
   for the entire selected tree. Renderer calls occur here only, using each fenced
   diagram's source. Never infer a global renderer for an explicit Markdown path.
3. Compute actual outcome capabilities, then register/run pure geometry gates using
   the frozen provider. Current host measureCanvases needs this lifecycle change;
   current PrepareResources(snapshot) runs too late to supply renderer-dependent
   dimensions to its preceding geometry gate and must not become a back door.
4. Prepare each finalized snapshot's native/background ticket from the same frozen
   outcomes and measured geometry. Repeated gates, tickets and resize may calculate
   text/fit geometry but must make zero renderer calls and no resource resolution.
5. Existing application Guard rechecks resource identity and renderer Revision
   before publication. Changed bytes/identity require a new guarded candidate even
   if source is unchanged. Failed/stale candidates discard their objects only;
   successful replacement/disposal releases references under existing bundle rules.

Dalton's read memo `WCI4-frontend-API.md` proposes optional
`svg.PreparedPreviewRenderer`, embedding ContentRenderer and
`CheckPreview(*parser.Instance) error`, on Options.Content
for pure per-path preflight. Missing/foreign outcome must fail even with label
policy; this check cannot invoke a renderer or reconsider content fallback. Validate
hidden/closed declarations before painting. Native SkipControls may omit only an
explicit SVG with an exact actually prepared native SVG inventory entry. Legacy
SVG stays drawn as before. A context hit proxy alone is not a prepared resource
renderer and cannot authorize omission. Existing descendant/unknown-kind checks
continue; collection/pane/value export restrictions do not disappear.

Dalton's conservative public-export proposal rejects explicit preview with
unsupported-resource-export; the contract permits that bounded choice. If a richer
or explicit-label public route is subsequently selected, prove its actual output
and failure artifact retention. Native Markdown background uses its prepared
Content provider, not an unverified raw-text fallback. Capabilities follow actual
rendered versus labelled outcome; no rich readiness follows from parser acceptance.

## Backend finding: composed Markdown images are not proven

After the initial handoff, James relayed the independent reviewer's concrete
pinned-backend result: direct SVG produced one oksvg path, whereas the current
Markdown-style nested `<image href="data:image/svg+xml;base64,...">` produced zero
paths without a default-mode error; strict mode reported unsupported content.
Read-only inspected reproducer: `/tmp/wci4-review-backend-9k83akv8/main.go`.
This worker did not rerun it or claim an OS screenshot; the reported result is
reviewer evidence and the source visibly exercises those two different routes.

A Resource-only Mermaid backend callback therefore cannot certify the composed
native background route. Valid bytes and successful direct-image decoding do not
prove that an embedded image survives the outer SVG renderer. The proposal above
must not freeze a rendered outcome on that evidence alone.

The independent reviewer subsequently corrected the scope after rereading the
original card preview inventory, Acceptance010 and Providers §4, as recorded in
James's updated provider memo. The approved bounded matrix is:

| Native route | Required outcome |
| --- | --- |
| Supplied closed-subset SVG | Positive direct resource rendering with shared aspect fit and clip. |
| Supported Markdown prose | Positive bounded prose rendering. |
| Mermaid on the unsupported embedded-image route | Per-diagram label/reject resolved before gates; preserve surrounding prose and advertise no rendered-mermaid capability. |

The previous suggestion that a positive native Mermaid route might be required to
close this inventory is withdrawn. No per-diagram native image helper, general
scene or extra runtime API is required by this bounded matrix. SVGRects remains
sufficient for the standalone resource SVG route; a richer diagram-image route
would need separate selection and its own verified representation.

Known backend unavailability still cannot bypass full binding identity validation,
resource/copy/revision checks or zero-rerender evidence. Unsupported diagrams use
source policy before outcome/capability freeze; they cannot silently disappear or
acquire fallback during Measure/gate. Label status remains measured and visible
at that diagram block while prose survives; reject retains the old bundle.

Native evidence must inspect positive direct SVG/prose output and actual labelled
unsupported diagrams, not only a nil backend error. This reviewed matrix correction
is not WCI4 code selection and does not authorize changes to frozen M2 product.

## 6. Proposed verification and remaining coordination

Synthetic layout tests after selection: asymmetric aspect-fit/centering, finite and
zero allocations, negative/nonfinite dimensions, source-relative minima/maxima,
small assigned height versus large intrinsic SVG, nested translation and clipping,
gutter exclusion, inactive panes, and independent nonmodal resize. Fit must remain
unchanged when only the visible clip/outer scroll offset changes, apart from the
single shared translation. Label fallback tests use the exact measured status and
ensure caption cannot suppress it. Legacy 0.2/basic-0.3 placeholders and bare
Markdown geometry/exports remain byte-compatible.

Provider/runtime-ticket tests: equal Markdown text at different paths with different
policies/descriptions/renderers; missing/extra/wrong bindings; identity failure cannot
fallback; safe copied renderer bytes; block versus whole-document label outcomes;
resource limits and aggregate-budget fatality. Count calls: once per prepared
required diagram outcome, zero from repeated Measure/gates/tickets/resize. Candidate
failure, stale Guard and renderer revision changes preserve the prior publication.
No runtime resource state or offset map should appear.

Actual native/export proof remains mandatory for accepted SVG shape/transform/solid
paint subset, clipping/aspect ratio/caption/status visibility and accessibility,
the approved per-diagram Mermaid fallback and its diagnostics. Pure geometry/XML acceptance does
not prove backend support or glyph fidelity. Public-export failure must preserve
the previous artifact. All-family package/helper/consumer and OS IME evidence stays
with main; this memo authorizes no build, install, launcher migration or publication.

James and Dalton agree on existing-package immutable outcomes/source helper and
no renderer in gates; James also agrees with the size-only adjunct and pure fit. Dalton accepts the natural/minimum adjunct and pure fitting
without AST geometry changes. Noether gives preliminary read-only agreement while
prioritizing M2, explicitly deferring full WCI4 adapter commitment. Remaining items
for selection: final provider accessors and text-overlay/native-accessibility proof, zero-image
minimum/source-allocation semantics and exact native/background resource route.
The reviewed bounded matrix above resolves the composed-Markdown route without
a new native block-image API. These are proposed implementation seams, not existing
exports or WCI4 selection.

## 7. Inspection identity and main handoff

| Inspected working path | SHA-256 |
| --- | --- |
| SDUI/go/layout/engine.go | 5f5405af50fc821c155d9289e85a306fa65c369a465841168e6d4a363efdd8c6 |
| SDUI/go/layout/collection.go | c9ae6e16c60997d90404956f63aa86abcf9ced0e211e63c9909d39d71ac1041c |
| SDUI/go/markdown/provider.go | 0166418cd0a1c6f67a93a530e08692562617bc361e788df1610d68ed3e1ac55e |
| SDUI/go/markdown/mermaid.go | 4c09f5152ced885b08b53fbc4dc2c4f98a34070e8d05a9eb80b933b718218efe |
| SDUI/go/host/fynehost/document.go | fe7c7db6ea0f13aa74461099c84152fb3f5b8b16b28fee49bf1cf738f141ad6b |
| SDUI/go/svg/render.go | 7fb8e908085bb6c28593a4eb6572b9532177a4041e14bd1e42847cfbd35f5e8a |
| SDUI/go/svg/check.go | 5a515532477bf90dc727078fef0f641ee30517558e2c00359f68d990b4405738 |

Read frontend memo SHA-256:
`532ef76335076cb80b76c2b39ad5602ccde73ca8cc4e26131fa24b963485311b`.
Read `WCI4-provider-API.md` after its publication; it agrees on immutable existing-
package ownership, preparation-only backend checks, root fingerprint validation,
size-only measurement and no runtime preview state. The initial read still spells
FitPreview without an error result; James confirms alignment to (Rect,error) and
proposes SVGRects as documented above. Verify final naming before implementation.
Provider memo SHA-256 at this read:
`e73bf764ce97126419927e61856c4d66ff5a654c4772200102283547fe117781`.
All ten hashes in WCI3-M2-layout-worker.md were checked unchanged after this task.

Inspection used git status/rev-parse, rg, sed/cat and SHA-256. No product tests or
renderer/native runs were executed for this read-only task. Only WCI4-layout-API.md
was written; existing M2 hashes/evidence remain untouched by this lane.

Session0010 S4 work summary for main: user authorized bounded WCI4 read-only layout
preparation while M2 stays frozen; identified natural/minimum and repeated host
Prepare boundaries, proposed small geometry-only seam, coordinated provider and
frontend ownership, recorded tests and remaining selection decisions. Main records
Session/roadmap updates and chooses/reviews the stage before dependent code.

Follow-up provider memo read after the scope correction, SHA-256:
`ae5e2d7813c2c3b2c0e6d712f20344cdf631668694db3dc93e4eb1e97cb88f35`.
Earlier memo hashes identify prior inspections, not the corrected current provider
memo. Main owns canonical recording and WCI4 selection; product remains frozen.
