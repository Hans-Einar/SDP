# WCI4 host and preparation API preparation

Status: read-only preparation, 2026-10-08. This memo is the only file written for
this assignment. WCI4 product implementation and package work are not selected.
WCI3-M2 product files remain frozen; an actual M2 native finding takes priority.
Main owns stage selection, canonical records, integration, dependencies/build roots
and OS evidence. Reused SDP 1.1.1, Worker 2.0.0 and the shared document workflow.

Inspected clone: `/tmp/sdp-sdui-widgets`, branch `sdui/widgets-wci3`, HEAD
`d6742742f0968661d15698518df688ffb1d28946` plus the existing multi-owner changes.
Authority: ORIGINAL `SDP/04--Design/SDUI/Widgets/Providers-and-packaging.md`, including
the reviewed Session0010 T003 backend/accessibility refinement, SHA-256
`9d9e53bcae0f54f8f2a3c1bb1fd65633ef6a6c5403597b08b997d2324c4a706d`.
Related handoffs: root `WCI4-provider-API.md`, `WCI4-frontend-API.md` and
`WCI4-layout-API.md`. Their earlier generic/native-Mermaid suggestions are superseded
by the canonical concrete matrix below, not by an inferred richer implementation.

## Concrete native scope

| Source/outcome | Actual bounded representation |
| --- | --- |
| Opt-in supplied SVG, validated closed subset | A direct Fyne canvas.Image from frozen SVG bytes, at the shared fitted rectangle. |
| SVG caption or labelled failure | Separate text-only SVG overlay using existing Go-Regular glyph outlines. No embedded image elements. |
| Explicit bounded Markdown prose | Existing background SVG text/path composition, from per-instance prepared content. |
| Explicit Markdown diagram | Current composed native route is unsupported. Freeze the declared per-diagram label/reject outcome before gates; preserve ordinary prose. No rendered Mermaid capability. |
| Bare Markdown and legacy SVG | Preserve current APIs, placeholder behavior and capability boundary. |

The independent reproduction at `/tmp/wci4-review-backend-9k83akv8/main.go`
shows the relevant distinction: direct SVG rectangle has one oksvg path; the same
SVG embedded as a Markdown data-URI image has zero paths with no default-mode error.
Strict mode rejects the image element. This worker inspected that reproducer and
the pinned Fyne/provider source; the run result is attributed to the reviewer.
Validating only the inner diagram resource cannot establish native readiness.
No per-diagram image adapter, generic scene, new renderer framework or richer export
is needed for the selected labelled-preview boundary.

## Request, preparation and publication

Use James's canonical `markdown.PreparedSVG`, `MarkdownRenderer` and immutable
`Previews` types; do not introduce parallel host resource records or runtime state.
The additive application request fields are exactly:

```go
SVGResources      map[string]markdown.PreparedSVG
MarkdownRenderers map[string]markdown.MarkdownRenderer
```

They are keyed by exact normalized instance path. The proposed detached sequence is:

1. Normalize/validate the selected document before constructing runtime gates.
   Validate all resource/renderer bindings, including hidden pages and closed
   dialogs: path/kind/use, source reference, identity, digest and budgets. A known
   unsupported backend does not bypass identity validation.
2. Call `markdown.PreparePreviews` once outside the pure preparation package.
   Copy application bytes and renderer output, validate content, classify backend
   representation and resolve each path's declared policy. Equal text may share
   parsed storage, never policy/description/outcome identity. Retain no external
   renderer as a measurement/paint dependency.
3. Supply the resulting data as `preparation.Request.Previews`. The proposed
   `CheckWithPreviews` checks its exact inventory/fingerprint against the root
   actually normalized inside preparation, before session/gate construction.
   Double normalization is acceptable; substituting an unchecked external root is
   not. Old request/Check behavior remains unchanged for legacy content.
4. Compose the existing host measurer with this frozen provider. Replace the
   current gate-time `markdown.Prepare(root,nil)` on this document-host route with
   pure prepared-content checks and measurement. Prepare legacy Markdown once in
   the composite provider too; do not change its public legacy API.
5. Build native/background presentation tickets from exactly the finalized
   snapshot and its already measured geometry. The existing accepted-ticket
   handshake remains authoritative. Failed/consumed/stale tickets must not install
   speculative pending resources or mutate the published controls.
6. Run the existing application Guard before publication. It rechecks source/model,
   resource identity/digest and renderer Revision. Different resource bytes or
   identity need a new candidate even if source text is unchanged. Non-yielding
   publication only installs admitted content/geometry; it does not run providers.

Current `DocumentRequest.PrepareResources(snapshot)` runs after measurement and is
per-presentation; it is not the hook for source renderer execution. Preserve it for
existing application resources. A successful reentrant legacy mutation still uses
the established accepted-presentation synchronization rules even if its enclosing
interaction returns an error.

`Previews.Check`/`CheckPreview` fingerprints source resource identity, explicit text,
description and policy. Caption, current visibility and layout are mutable
presentation facts; pure placement reads them from the current snapshot without
repreparing resource bytes. Accessor copies should be acquired once for the host's
private resource projection, not repeatedly copied for every geometry probe.

## Backend validation and exact capabilities

Use the fixed preparation-only backend checks in James's memo. Direct SVG must
pass the provider's closed grammar/security/dimension checks and the concrete
pinned native decoder's support check. Existing `prepareIcons` already uses public
oksvg strict parsing, but icon presence and icon limits are not preview admission.
Do not import Fyne internal packages, patch dependencies, or treat a void
canvas.Image.Refresh / nonnil Resource as a successful backend check.

Native Markdown readiness covers the actual background text representation.
Its Mermaid backend check reports unsupported for the current composite route,
even for a valid inner Resource. Provider preparation still enforces identity,
copying, custom-output validation and limits in the agreed order. No renderer
invocation, resource resolution or fallback reconsideration occurs in Measure,
StateGate, resize or native publication. Fyne rasterization of already admitted
immutable bytes is distinct from invoking the application's source Renderer.

Outcome-aware requirements remain separate from concrete host implementation
capabilities. Actual rendered SVG requires widget svg/1, provider svg-resource/1,
layout preview-resource/1 and host svg-resource/1. Labelled SVG requires existing
placeholder facts and reports fallback. Explicit Markdown requires provider
markdown/1 plus native host markdown/1; this native route supplies no rendered
mermaid-flowchart fact. Keep all WCI1–WCI3 facts and exact major matching intact.
No broad boolean, resource dimension or generic registry is introduced.

Frontend's `svg.PreparedPreviewRenderer` check on Options.Content validates the
frozen explicit-Markdown path, including hidden/closed declarations. Native
SkipControls omits only the exact opt-in SVG path with an actually prepared native
preview adapter. A legacy SVG context proxy is not such an adapter. Bare SVG stays
background-painted. Public export may retain the proposed conservative
unsupported-resource-export rejection; no native omission map licenses public
export to silently hide unsupported content.

## One geometry and text computation

Agree with Gibbs's proposed `layout.PreviewMeasurer.MeasurePreview` returning
natural and minimum sizes for opt-in SVG, plus pure `layout.FitPreview`. Both the
normal desired-size path and pane minima need that distinction. The image minimum
is zero, subject to source/container minima; do not make intrinsic image size a
native chrome minimum. Required caption/status bands have their actual measured
minimum. Zero-area image allocation has no painted image; do not rasterize it or
invent a second pixel-size layout floor.

Agree with James/Gibbs's fixed geometry accessor:

```go
func (p *Previews) SVGRects(box *layout.Box) (image, caption, status layout.Rect, err error)
```

Use the same pure band/line plan for MeasurePreview, SVGRects and the following
bounded text serializer proposed to the provider owner:

```go
func (p *Previews) RenderSVGText(out *strings.Builder, box *layout.Box) error
```

It emits ONLY caption/status glyph paths using existing `layout.TextWidth/Lines`
and `svg.Text`, which use embedded Go Regular without font hinting. It does not
embed the resource image, invoke an external renderer or return a general scene.
This is preferred over the earlier native canvas.Text/metrics-callback proposal:
Fyne text shaping is not assumed identical merely because the font name matches.
James, Gibbs and host agree on this smaller approach; the independent reviewer
also gave scoped design agreement. This is API preparation, not code selection.

The optional caption occupies one measured row and may ellipsize using the same
metrics. Full bounded fallback text wraps in its own measured band, independently
of caption presence. Its required text region must fit admission; never suppress
status or substitute rejected image dimensions as its minimum. Full description
and diagnostic remain in the prepared outcome/accessibility label. Exact paint,
wrap and minimum-size tests are required before claiming this implementation.

The native adapter uses two ordinary image objects: direct admitted resource and
the provider's text-only SVG overlay. The host creates the overlay root/viewBox
from the same owning Box allocation, translating the serialized coordinates once.
Image destination comes directly from SVGRects after caption allocation and
FitPreview. Native ImageFillStretch paints into that already aspect-correct
rectangle; it must not independently refit the clipped fragment. The text overlay
contains paths/title/group only, never a nested data-URI image. Shared Box.Clip is
applied afterward to both through existing clipping wrappers. No new hit tree,
scroll authority or preview camera is introduced. Context/wheel routing continues
through the existing bounded outer control adapter.

Validate the actual complete text-overlay SVG representation while preparing its
ticket, before publication; glyph/serialization failure explicitly rejects that
ticket rather than publishing blank status or choosing another fallback. This
pure serialization of current caption/geometry does not
reopen frozen resource policy or invoke the application's renderer. Backend
availability for that representation is established before outcome freeze.

For intrinsic 400x200 and allocated image content (10,20,100,20), the stored fitted
rect is (40,20,40,20). Scrolling clips that fixed rect; it must not resize it to the
visible fragment. Main, modal-parent and nonmodal canvases use the same calculation.
Explicit Markdown retains the prepared provider's existing ordered block geometry;
the labelled diagram replaces only its block, not the surrounding prose.

## Native accessibility and inspection

Pinned Fyne 2.8.1 publicly exposes only AccessibilityLabel() string and
AccessibilityRole() fyne.AccessibleRole. There is no separate public image
description method or image role. Implement those methods on the actual mounted
SVG and explicit-Markdown preview wrapper, using exact prepared description and
outcome status, independent of caption. Use its text/container role truthfully;
do not add a fake action or keyboard focus stop just for metadata.

The mounted clipped wrapper has the shared visible rectangle, not the unclipped
intrinsic artwork bounds. Hide it when its page/surface is hidden. Keep the full
Box and fitted artwork rect separately for inspection. No broad change to Fyne's
accessibility traversal or existing controls is selected. Fyne's Linux GLFW
accessibility backend is a no-op even with the build tag. Darwin/Windows have
additional build/traversal constraints and no evidence here. Claim the actual
native adapter interface and visible fallback only, not OS screen-reader delivery.

Proposed additive diagnostic `Inspect()["previews"][path]`:

```text
kind, status, description, diagnostic, providerID, revision, sha256
rect, clip, image, caption, statusRect, visible, canvas, title
accessibleLabel, accessibleRole
```

Description/outcome identity is frozen-provider evidence; accessibleLabel/Role must
be queried from the mounted fyne.Accessible object. Rect/clip/visible/canvas/title
come from actual mounted objects and the existing per-canvas inspector transforms,
not reconstructed hypothetical controls. Main canvas is literal "main"; modal
uses its actual parent's canvas, nonmodal its surface path. Closed declarations may
appear as prepared but unmounted/hidden; never pretend they have a native window.
No resource bytes, fake provider readiness or reconstructed accessibility proof
is exposed. Native screenshot evidence still comes from main's OS harness.

## Resource ownership and disposal

The bundle owns the copied validated bytes and frozen per-path outcomes. Public
input maps and renderer references should not remain retained in its stored request
after preparation; retain the prepared identities and application Guard instead.
The caller still owns its own maps, renderer process/configuration and cancellation.
No global application resource cache or provider service is added.

Immutable bytes may share by full digest, but native CanvasObjects belong to one
canvas/adapter at a time. Separate main/nonmodal/surface instances have independent
geometry and lifetime; never share one canvas.Image object between canvases.
Resource names should include the full digest to prevent accidental native cache
identity collisions. Do not expose writable authoritative byte slices through
diagnostics or prepared accessors. A private fyne.Resource wrapper can return owned
copies from Content if public access would otherwise expose the frozen backing
slice. Count distinct authoritative bytes by the provider's budget, not temporary
accessor/decoder copies.

Prepare new native objects detached. A failed ticket discards only objects/resources
owned by that candidate; existing published controls and immutable content remain.
Publishing replacement revokes old callbacks first, then drops old bundle/surface
references under the existing exact-target lifecycle. Closing a surface removes
its native image objects while the bundle retains declaration resources needed for
reopening. Bundle.Close drops preview maps/resources/overlay references after
detaching views and closing its surfaces, including abandoned pending tickets.
Current Close does not explicitly clear every presentation/image reference, so
bounded preview-specific cleanup must be part of the eventual implementation.

Fyne has no public per-image GPU/cache Dispose contract. Claim release of our
references and native object detachment, with driver/GC cache reclamation owned by
Fyne; do not promise immediate GPU-memory erasure. Disposal must not destroy shared
bytes still owned by the newly published bundle or another live canvas.

## Required implementation evidence, not results from this memo

- Identity-first failures under label policy, hidden/closed declarations, missing
  prepared paths, equal Markdown text with distinct path policies, input/output
  buffer mutation, digest limits and renderer Revision/Guard changes.
- Actual direct SVG decoder and OS image proof for the closed subset; supported
  prose with unsupported diagrams becomes visible labelled blocks or rejects.
  Never advertise rendered Mermaid from a resource-only validation success.
- Zero source-renderer calls during gates, resize, caption/visibility changes,
  accepted tickets, failed tickets and Commit. Failed reload keeps the live bundle.
- Shared fitting/negative-position clipping in main and nonmodal/modal-parent
  canvases. Caption and full status remain consistent with minimum measurement at
  narrow widths and several fonts; native overlay output has no image elements.
- Exact native omission inventory and pure explicit-Markdown fingerprint checks;
  legacy outputs and all prior family admission remain intact.
- Inspector reads actual mounted Accessible methods and geometry; visible fallback
  and hidden/closed behavior are checked independently. No OS accessibility claim.
- Failed candidate, successful replacement, surface close/reopen and bundle teardown
  release only their own references without affecting another live image/control.

Likely implementation files, not permission to change them now: host document
request/preparation, metrics composition, view/preview adapter, native inventory,
surface ticket wiring, inspector and scoped tests; preparation request/capability
checks; host README. Provider/parser/layout/export changes remain their owners'
exclusive lanes. No runtime, SDL transport, module or package changes are proposed
for this host seam. Main coordinates the later actual SDL fixture and package set.

## Inspection and coordination evidence

Read the original canonical contract and the three root WCI4 memos; inspected
document.go, document_surfaces.go, view/clip/inspector/admission/preparation code,
markdown/provider.go, layout/font.go, svg/text.go, and pinned Fyne image,
accessibility and SVG-decoder source. Used read-only git status/rev-parse, sed, rg
and SHA-256 checks. No WCI4 build, test, native execution or product change occurred.

Coordinator/reviewer confirmed the backend matrix and actual Accessible adapter
boundary. James, Gibbs and host agreed the fixed three rectangles plus
`RenderSVGText(out *strings.Builder, box *layout.Box) error`; Mendel gave scoped
design agreement. James owns its implementation. The earlier native metrics/line
callback alternative is superseded; no metrics callback or line protocol is needed.
No competing stub or callback framework was implemented. Main still owns canonical
projection, stage selection and later native proof.

Final read-only verification compared each SHA-256 entry in the frozen M2 report:
**17/17 product files match exactly**. The prior M2 report itself remains SHA-256
`42090c685ddd0f8315979b9872c1f56a006b218d5eca6446ed621447d498701f`.
No product test was repeated for this memo. Only `WCI4-host-API.md` was created;
the shared clone's other dirty files belong to the ongoing existing lanes.

Final coordinator acknowledgement retains a single measured caption row with
ellipsis, full bounded fallback status wrapping at its admitted minimum, and
explicit rejection on glyph-paint failure. WCI4 remains unselected until M2 close.
The subsequently requested full M2 SDUI race run is separate evidence in a `/tmp`
run directory, not WCI4 implementation or verification.
