# WCI4 provider/content API preparation

Status: bounded read-only handoff, 2026-10-08. WCI4 is not selected. Only this
root memo is written; M2 runtime source and its evidence remain frozen. Reused
SDP Worker and shared document workflow. Main owns canonical decisions, Session,
stage selection, package/dependency work and native/integration evidence.

Inspected `/tmp/sdp-sdui-widgets`, branch `sdui/widgets-wci3`, HEAD
`d6742742f0968661d15698518df688ffb1d28946` at final inspection, plus concurrent
M2 work. The coordinator advanced HEAD; this worker made no commits or branch changes. Authority is
ORIGINAL Providers-and-packaging.md, independently reviewed substantive contract
909851e3. All signatures below are proposals for lane agreement, not current exports.

## 1. Existing implementation and required bounded changes

| Inspected behavior | WCI4 consequence |
| --- | --- |
| markdown.Resource is SVG bytes plus Width/Height; Renderer.Render receives source | Reuse these values/signature; copy and validate every returned payload, including custom renderers. |
| markdown.Prepare caches Documents by source text and Resources by truncated diagram hash | Preserve that legacy API, but never use its caches as the authoritative explicit-preview policy/outcome. Equal text at distinct paths must stay independent. |
| Prepare accepts custom Renderer output without validation/copying | New preparation must enforce bounds, dimensions, SVG safety/profile and ownership for all outputs. |
| Mmdr enforces flowchart/graph, 12000-byte source, no %%{ configuration, five-second process deadline and capped stdout/stderr | Preserve the trusted explicit executable boundary. It is not a discovery mechanism or a deadline wrapper for arbitrary Go Renderer implementations. |
| validateResource scans XML but does not track a unique root, trailing content, namespaces, processing instructions or a closed attribute/element grammar | Reuse/refactor the security scanner; it is insufficient for the new closed supplied-SVG subset without stricter validation. ViewBox x/y and shape/transform operands also need validation. |
| markdown.Parse limits source, blocks and diagrams and rejects HTML/images | Keep bounded parsing; audit final block-count enforcement (the current check occurs before appending, so the last appended block needs an explicit final bound). This is an inspection concern, not a reproduced defect claim. |
| DocumentHost.measureCanvases calls markdown.Prepare(snapshot.Root,nil) on each gate | Prepare explicit content once before the first gate and reuse frozen records. Merely passing a nonnil renderer into this existing gate call would violate the contract. |
| preparation.Prepare normalizes internally, calls capability Check before session construction and installs a pure presentation gate before native preparation | Supply frozen content as data, validate it against that actual normalized root, and resolve outcome-aware capability facts before those gates. Keep preparation free of renderer/process I/O. |
| svg.ContentRenderer serves Markdown; native SVG placeholders are always background-painted | Explicit Markdown needs prepared-path evidence; actual resource SVG needs an exact native inventory entry/adapter. Legacy placeholder behavior remains. |
| Generic layout widget minimum equals natural size | A shrinkable resource image needs the layout owner's bounded preview measurement seam; fallback labels still require their measured text minimum. |

No runtime state, Handle kind, event, loader, Snapshot preview map or StateRevision
extension is needed. A preview is immutable content in one prepared document bundle.
Existing model/source identity and application Guard already govern replacement.

## 2. Source facts and ownership

Dalton proposes the sole parser accessor:

```go
type PreviewPolicy struct {
    Explicit bool
    Description, Fallback string
}
func PreviewOptions(n *Instance) (PreviewPolicy, error)
```

It accepts normalized svg or Markdown; bare Markdown has empty Arguments and
Explicit=false. SVG opts in by description/fallback presence, requiring both.
Explicit Markdown has Kind markdown, empty Widget, Text matching Arguments[text],
and retains all three typed arguments. Source reference remains the validated
Arguments[source].(Reference); module membership stays normalization's responsibility.
No provider package is imported into parser. Captions remain independent from status.

Smallest owner proposal: extend the existing markdown package with bounded preview
preparation alongside the unchanged legacy Provider. It already owns Resource,
Markdown parsing, Mmdr validation and pure layout/svg interfaces. A new generic
resource registry/package is unnecessary. Keep new files separated by preparation,
SVG validation and pure projection so the old public map-based Provider contract
is not silently tightened.

Dependencies remain acyclic: preparation -> markdown -> layout/svg -> runtime/parser.
Layout must not import markdown; it uses geometry-only interfaces. Runtime imports
neither the preview provider nor preparation/GUI. Host owns native images and final
backend admission. Frontend/export owner owns parser/schema/codegen/static checks.

## 3. Proposed concrete data and methods

```go
// package markdown; canonical application inputs
type PreparedSVG struct {
    Source parser.Reference
    ProviderID, SHA256 string
    Resource Resource
}
type MarkdownRenderer struct {
    ProviderID, Revision string
    Renderer Renderer
}

// Closed preparation-only backend checks, supplied by native/static owner.
// They never run during Measure, paint, gates or Commit. Nil means unavailable.
// Each check certifies its actual composed backend route, not bytes alone.
type PreviewBackend struct {
    SVG func(Resource) error
    Markdown func() error
    Mermaid func(Resource) error
}

type DiagramOutcome struct {
    ID, SHA256 string
    Resource *Resource // nil for labelled failure
    Diagnostic string // bounded English reason, empty on success
}
type PreviewOutcome struct {
    Path, Kind, Description, Fallback string // Kind: svg or markdown
    Span parser.Span
    Uses []parser.UseSite
    Source parser.Reference // symbolic SVG triple; zero for Markdown
    ProviderID, Revision, SHA256 string
    Status string // closed values: rendered, partial, label
    Diagnostic string // whole-node failure only
    Resource *Resource // SVG success only
    Document *Document // valid Markdown's ordered blocks/diagrams
    Diagrams []DiagramOutcome
}
type Previews struct { /* private records and private legacy projection */ }

func PreparePreviews(root *parser.Instance,
    resources map[string]PreparedSVG,
    renderers map[string]MarkdownRenderer,
    backend PreviewBackend) (*Previews, error)
func (p *Previews) Outcome(path string) (PreviewOutcome, bool)
func (p *Previews) Check(root *parser.Instance) error
func (p *Previews) CheckPreview(n *parser.Instance) error
func (p *Previews) Measure(n *parser.Instance, font, width float64) (layout.Size, error)
func (p *Previews) Render(out *strings.Builder, box *layout.Box) error
func (p *Previews) MeasurePreview(n *parser.Instance, font, width float64) (natural, minimum layout.Size, err error)
func (p *Previews) SVGRects(box *layout.Box) (image, caption, status layout.Rect, err error)
func (p *Previews) RenderSVGText(out *strings.Builder, box *layout.Box) error
```

Backend callbacks are three fixed route checks, not a new provider registry,
capability dimension or readiness flag. Their owner checks its existing exact
capabilities and renderer representation. Nil/unavailable routes are content failures
subject to source policy. They receive defensive resource copies and are discarded
with all Renderer references after preparation. If the selected adapter guarantees
a route by construction, its corresponding check is still explicit and deterministic.
For the now-reviewed native matrix, Mermaid is unavailable: its check rejects
(or is nil) and each diagram resolves to declared label/reject before freeze.
Validating standalone diagram bytes cannot certify the unsupported composite route.
No alternative native Mermaid representation is a prerequisite of this scope.
Main should confirm this small admission seam with host before implementation;
there is no need to publish a generic validation callback framework.

Outcome returns a deep copy: Resource bytes, Document.Blocks/Table rows, Diagrams,
UseSites and nested resources. Internal validated bytes/parsed documents can share,
but records and policy results never share by text. Public callers cannot mutate
records or force a rendered state. SHA256 is lowercase full digest of rendered bytes;
Markdown's record retains provider ID/revision while individual diagrams hold hashes.
A private source fingerprint covers path, kind, text/reference triple, description
and policy. It excludes mutable caption/layout/visibility so normal runtime Label,
page or viewport updates do not invalidate immutable content. Check(root) verifies
complete exact explicit inventory and each fingerprint; CheckPreview handles a node.

The composite can privately prepare the unchanged bare Markdown route once with
its existing nil/global renderer policy as appropriate to the caller. Existing
Prepare(root, renderer), exported Documents/Resources and legacy output remain
available and semantically unchanged. Never route explicit nodes through the legacy
text-keyed policy lookup, even if parsing/payload bytes are shared internally.

## 4. Identity-first preparation order

1. Normalize/validate the selected source and enumerate ALL explicit previews,
   including inactive pages/closed dialogs/reused instances. Freeze application map
   entries and scalar identity fields for the candidate. No renderer is called yet.
2. Validate the ENTIRE binding inventory before content fallback: exact paths/kinds,
   no extra entries, source triples, nonempty provider IDs/revisions, nonnil bound
   renderers, canonical full lowercase SHA256 and equality with supplied bytes.
   Missing entries are availability failures, not malformed bound entries. A bound
   renderer on prose-only/bare Markdown is an unused binding and fatal.
3. Determine diagram presence using bounded syntax inspection independent of content
   support. Existing Parse may return early for HTML/images; do not use that error
   to skip identity checks or guess whether a renderer binding is unused. A fenced
   diagram inside an otherwise unsupported explicit document can still establish
   binding relevance; the whole content then follows fallback. No renderer runs
   merely to discover identity or diagram presence.
4. Prepare each valid declaration under its own policy. Plain supported Markdown
   requires no renderer. Keep valid prose when a diagram fails; replace the whole
   preview only if the document profile itself fails. Render each needed diagram
   from its fence source only. Flowchart/graph and source/configuration restrictions
   must apply even when the registered Renderer is custom, not just Mmdr.
5. Copy, validate, digest and budget each successful resource before storing it.
   Check backend representation and exact capability availability here, then freeze
   rendered/partial/label outcomes. No later gate may select fallback anew. Fatal
   candidate errors discard the detached content and leave the live bundle untouched.

Identity errors anywhere remain fatal even if earlier nodes use fallback=label.
In particular, a digest mismatch cannot be downgraded to malformed SVG fallback.
Custom renderer errors are bounded (propose <=4096 UTF-8 bytes), sanitized for a
single English status reason; no arbitrary invalid bytes reach image decoding.
Fallback text is always `Preview unavailable: <description> (<reason>)` and cannot
be hidden by a nonempty SVG caption. Whole-node and per-diagram reasons remain in
the preparation report and accessible description/status.

## 5. Validation and budgets

Share a pure low-level validator in markdown with Mmdr, with two explicit entry
points rather than a caller-supplied grammar: proposed
`ValidateSVGResource(Resource) (Resource,error)` for supplied shapes and
`ValidateMermaidResource(Resource) (Resource,error)` for registered verified output.
Both return copied validated bytes/dimensions. The existing private Mmdr validator
can use the common document/security checks without making the new shape whitelist
an accidental ban on its separately verified Mermaid profile. Legacy API behavior
needs regression tests if common internals are tightened; do not claim rich-preview
support merely because the legacy XML scanner returned success.

Common checks: nonempty <=4 MiB, exactly one complete SVG document, finite viewBox
x/y/width/height with positive width/height <=32768, supplied dimensions finite and
equal to viewBox width/height. Track root start/end/depth and reject second roots,
trailing non-whitespace or malformed end state. Reject directives/DTD, active/external
content, unrecognized namespaces, duplicate/unknown attributes and processing
instructions other than a valid initial XML declaration. No automatic URL, image,
font, stylesheet or entity acquisition. Byte limits apply before copying/parsing;
stream the supplied digest before content classification where required by identity
precedence. Invalid oversized bytes are never retained/rendered.

Supplied-shape closed subset: svg/g/path/rect/circle/ellipse/line/polyline/polygon/
title/desc only. Declare each allowed attribute by element: root viewBox/namespace,
finite geometry/path data, supported finite transforms, solid fill/stroke,
nonnegative stroke width and bounded opacity. Optional root width/height must have
unambiguous finite units consistent with the validated dimensions. Reject CSS/style,
class/text/fonts, gradients, href/url references and any other unsupported element
or attribute. Parse actual path commands/operand counts/arc flags, points and
transform operands; a regex for script or a successful XML parse is not validation.
Reject nonfinite intermediate/composed transforms and backend-unsupported forms.
Do not add arbitrary pixel/element/depth limits silently beyond the reviewed bounds;
if concrete parser/backend constraints require them, report the bounded difference.

Mermaid resource validation is distinct from the supplied-shape whitelist; do not
blindly apply one grammar to the other or generalize either to arbitrary SVG.
Custom outputs still require common security/profile checks. The reviewed current
native explicit-Markdown diagram route is unsupported and resolves per-diagram
label/reject; resource validity never promotes it to rendered readiness. Legacy
registered renderer behavior and records remain unchanged.
Native image/backend preflight is preparation-only; later native image decoding can
consume only frozen admitted bytes, never invoke the source renderer.

Count distinct copied and validated resource bytes admitted during preparation (before backend classification; no fallback refund) by full SHA256, <=32 MiB for the
whole candidate across supplied SVG and rendered diagrams/all surfaces. Deduplicate
storage/budget by full digest after validation, not source text/truncated diagram ID,
provider ID, or fallback policy. Equal bytes can share storage, while per-instance
outcomes stay independent. Charge before allocating another retained copy; crossing
32 MiB is a fatal aggregate failure, never label fallback. Per-resource size,
dimensions or unsupported content uses node policy; a failed resource is not stored.
Defensive accessor/temporary backend copies do not become new authoritative budget
entries. Keep the source 32768-byte/256-block/eight-diagram limits and Mmdr's existing
process bounds. No new async loader, remote fetch or process launcher discovery.

## 6. Preparation, gate and publication wiring

Keep preparation free of provider I/O. Proposed additive data seam:

```go
// package preparation.Request
Previews *markdown.Previews
// Existing Check unchanged for legacy callers; explicit preview callers use:
func CheckWithPreviews(profile string, root *parser.Instance,
    supported Capabilities, previews *markdown.Previews) error
```

Host DocumentRequest receives canonical SVGResources and MarkdownRenderers maps;
Bundle privately owns the prepared Previews. Before calling preparation.Prepare,
host normalizes its frozen document for content preparation. preparation.Prepare
still normalizes internally, then Check(root) proves the supplied records match
THAT actual root before capability checks/session/gates. Do not add an injectable
normalized root or trust a provider prepared for an unrelated source. A second
pure normalization is acceptable; a new staged preparation framework is unnecessary.
The application Guard captures source/model and resource digest/provider/renderer
revision identities and rechecks them before final publication. Bytes alone do not
make stale bindings current. A changed renderer Revision or resource identity needs
a new candidate even with unchanged source text.

CheckWithPreviews retains every old family requirement and adds facts from actual
outcomes. Rendered resource SVG requires widget svg/provider svg-resource/layout
preview-resource (and host svg-resource on native route). Label SVG uses only the
placeholder facts. Explicit Markdown uses provider markdown/host markdown; successful
Mermaid blocks on a separately supported route would require mermaid-flowchart;
the reviewed current native route never emits that fact. A missing native route must
already have produced authorized fallback during content preparation or rejected;
Check does not mutate outcomes or invoke callbacks. Capabilities never license a
wrong inventory/source or suppress failures elsewhere in the tree.

Replace explicit-preview handling in measureCanvases with b.previews.Measure and
frozen outcomes. Repeated gates, resource tickets, resize, page switches and separate
dialog canvases reuse them. Neither markdown.Parse nor Renderer.Render nor fallback
selection is needed there. PreparePresentation still builds only native objects and
background from the finalized measured snapshot; only its accepted ticket sets
pending presentation. Commit only guards/swaps. Failed/abandoned candidates release
their own references/native objects; old content remains through rejection and is
released only after successful replacement or disposal.

Do not use DocumentRequest.PrepareResources for source renderer invocation: it is
called for each finalized runtime presentation, not once per source bundle. It can
continue validating application native resources against the final snapshot.

## 7. Layout, native and export seams agreed in preparation

Gibbs proposes this optional geometry-only interface on the existing measurer:

```go
// package layout; markdown.Previews implements it, layout imports no markdown
type PreviewMeasurer interface {
    MeasurePreview(n *parser.Instance, font, width float64) (natural, minimum Size, err error)
}
func FitPreview(content Rect, intrinsic Size) (Rect, error)
```

Use it only for opted-in resource SVG. Resource dimensions give natural size;
minimum may be zero for shrinkable image content, subject to existing source/split
allocation constraints. Label fallback minimum is measured visible status text,
not zero or raw SVG dimensions. FitPreview preserves aspect ratio and centers
unused space using the minimum axis scale; final native image also honors Box.Clip.
Fit against the full allocated content rectangle after removing the measured
caption band once, not against its visible clip; scrolling must not resize/recenter
the artwork. Proposed SVGRects returns only three fixed geometry slots computed by
the same pure text-band/fit helper: image, caption, status (zero for absent slots).
Gibbs agrees with this fixed-slot seam. Rectangles use the owning Box canvas
coordinates; absent slots are exactly Rect{}. The caller intersects each with
Box.Clip, without refitting the clipped fragment. It calls layout.FitPreview and
does no parsing or provider work. The accepted native
ticket stores these rectangles and paints them directly; host must not refit them.
This avoids duplicate caption sizing without a public scene or outcome map in layout.
Host and layout agree on RenderSVGText below to keep text measurement and paint
on the existing embedded-font path; no native metrics callback is needed.
No new runtime viewport, hit geometry, custom scene, pixel floor or content-sized
outer extent. Explicit Markdown stays on existing Measure/Render geometry, with
path-based outcomes and ordered diagram resources/statuses.

Dalton proposes an optional `svg.PreparedPreviewRenderer` interface embedding
ContentRenderer plus CheckPreview(*parser.Instance) error. Previews implements it.
Native explicit Markdown checks this pure evidence even in hidden descendants;
renderer presence alone is insufficient. Native resource SVG omission requires the
exact path/kind of an actually prepared SVG control, never a broad skip. Legacy
placeholder SVG remains background-painted. Public SVG conservatively rejects all
explicit previews with unsupported-resource-export before writing output; supplied
Content cannot bypass that policy. That is allowed by the reviewed contract and
does not promise a general embedded-image export implementation. Existing other
family export rejections remain.

## 8. Parallel ownership and acceptance after selection

- Frontend/export: PreviewOptions, explicit normalization/codegen/schema tests,
  source-only unprovided status, prepared-content check interface and conservative
  public SVG rejection. No provider imports into parser or resource resolution.
- Provider/content lane: existing markdown package's new inputs/outcomes/identity
  pass, validation, copying/budget/fallback, pure projection and adversarial tests.
  Preserve legacy APIs and exact old output where promised.
- Layout: PreviewMeasurer/FitPreview, natural/minimum constraints, aspect/clip tests,
  split/relative/hidden geometry; no byte/XML parsing or provider invocation.
- Host/preparation: request maps, bundle-owned Previews and root matching, backend
  checks/outcome capabilities, one-time preparation, native accessibility/status,
  image lifetimes and accepted-ticket integration. Explicitly assign preparation
  files to this lane to avoid overlap with provider ownership.
- Runtime: no product additions identified; existing snapshot/ticket/successor
  semantics suffice. Main owns package roots, staged helper payloads, dependencies,
  external-owner handoff, source identity guards and native/integration evidence.

Freeze inputs/outcome/accessor and layout geometry signatures first; then provider
and frontend/layout work can proceed in parallel. Host composes the actual helpers,
not duplicate caches/stubs. No SDL action/bridge protocol change is needed for
immutable previews. WCI4 package preparation remains a distinct selected unit, not
permission to install/publish helpers or infer standalone .3 readiness.

Required focused tests: global identity fatality despite earlier fallback; malformed
lowercase/full digest and source mismatch; missing/extra/unused/nil renderer including
unsupported documents; equal text across two paths with different renderers/policies;
mutation of supplied/custom bytes and returned outcome copies; exact one-document/
namespace/attribute/path/transform rejection; finite/matching dimensions; 4/32 MiB
boundaries and full-digest dedup; custom renderer safety; mixed prose/diagram fallback;
whole-document fallback for unsupported Markdown; labels despite captions; all hidden
nodes; zero extra renderer calls through gates/resize/Apply/tickets; stale guard,
failed preparation/reload and actual native geometry/clip/status. Inspect actual
direct-SVG and Markdown-prose native positive output, plus explicit diagram
unavailability and actual mounted accessibility interfaces. No tests or native proof were run for
this read-only memo.

No material contradiction with the reviewed architecture was found. Concrete
implementation gaps are listed in §1, especially identity/outcome ownership, custom
output validation and current gate-time Prepare. The proposed fixed backend checks,
double-normalization/root-match and fixed SVGRects caption geometry seams are
aligned with host/layout in the final follow-up below. Main still owns selection
and exact implementation review; these proposals do not expand runtime responsibility.

## Inspection identity

Canonical original SHA-256: `83d54b218c7c2dfbebff17b641c2fdc67e4ca7d1f5fb73fe65851de5bee8ad46`.

| Inspected file | SHA-256 |
| --- | --- |
| `SDUI/go/markdown/document.go` | `7359732245e33d010cde39b9392da794d57cf1dbc582d9201acd034b5443c73a` |
| `SDUI/go/markdown/provider.go` | `0166418cd0a1c6f67a93a530e08692562617bc361e788df1610d68ed3e1ac55e` |
| `SDUI/go/markdown/mermaid.go` | `4c09f5152ced885b08b53fbc4dc2c4f98a34070e8d05a9eb80b933b718218efe` |
| `SDUI/go/markdown/provider_test.go` | `f2ae3003b20e34cc82ec78e09d072fb282227b72edb5af2418a303cb166d1c71` |
| `SDUI/go/preparation/prepare.go` | `fbd200d3d38837021ca9d762b9769796f3321e465696a24d8faa76dc3d106771` |
| `SDUI/go/preparation/capabilities.go` | `85823f21f160b841c7b4e7c0c4fa995ebcbf2850248729c70d2b068e5ab391b8` |
| `SDUI/go/host/fynehost/document.go` | `fe7c7db6ea0f13aa74461099c84152fb3f5b8b16b28fee49bf1cf738f141ad6b` |
| `SDUI/go/host/fynehost/collection_metrics.go` | `d686f3fe8cb6b0d411cf9b5fcd79bf38639eb646a8ee054774faf5b2f76f5fb3` |
| `SDUI/go/host/fynehost/command_resources.go` | `f3a517fcf6c5732b4097bbd1e522f458772e9ab1fb6f3611819ee3f4e9d5a639` |
| `SDUI/go/svg/render.go` | `7fb8e908085bb6c28593a4eb6572b9532177a4041e14bd1e42847cfbd35f5e8a` |
| `SDUI/go/svg/check.go` | `5a515532477bf90dc727078fef0f641ee30517558e2c00359f68d990b4405738` |
| `SDUI/go/layout/engine.go` | `5f5405af50fc821c155d9289e85a306fa65c369a465841168e6d4a363efdd8c6` |

Read frontend/layout memos after their publication and aligned PreviewOptions,
PreparedPreviewRenderer and FitPreview signatures. All 12 entries of the frozen
WCI3-M2 runtime manifest were verified unchanged at completion. Only this memo was
written; no prior evidence manifest was modified.


## Independent follow-up: native Markdown image route is not supported as written

Independent reviewer supplied `/tmp/wci4-review-backend-9k83akv8/main.go`; this worker
read the reproducer but did not rerun it. Reviewer-reported pinned Fyne 2.8.1/oksvg
results: a direct SVG rect decodes to one path, while the same SVG wrapped in
`<image href="data:image/svg+xml;base64,...">` decodes to zero paths with no default
error; Strict mode rejects the image element. Existing markdown.Provider.Render
emits that exact embedded-image representation for diagrams. Thus the current
native Markdown background route can silently omit a valid diagram.

This is a concrete implementation seam, not a new resource security rule. A
resource-only Mermaid backend callback cannot justify Status=rendered or a
mermaid-flowchart fact. Before outcomes/capabilities freeze, host must either select
and prove an actual supported native representation (for example a bounded native
image at the provider's prepared diagram rectangle) or resolve the affected diagram
through its declared label/reject policy. Prose remains; the visible failure status
must not be hidden. Do not defer discovering unsupported embedding until paint.
The reviewer subsequently corrected the acceptance interpretation: the owner card
permits a labelled preview placeholder and richer preview is optional. The approved
bounded matrix requires positive direct closed-subset SVG and bounded Markdown
prose; native Mermaid is unsupported and may use explicit per-diagram label/reject.
No rendered-mermaid capability may be inferred. A positive native Mermaid image
workflow is not mandatory for this selected scope; the earlier contrary statement
is withdrawn.

Native block images are not required by the approved bounded matrix. Only if
separately selected later, expose the needed immutable per-diagram placement from the same pure block routine and store it in the existing native
presentation ticket. Avoid a general scene/hit graph, runtime preview state or a
second fit calculation. Exact placement/accessor names and native text metric
agreement remain with main/host/layout before selection; the proposed SVGRects
for standalone resource SVG does not alone place Markdown diagram blocks.

Main, Noether and Gibbs were notified. M2 product source remains frozen; this
follow-up changes only the authorized provider API memo.


## Reviewed bounded route matrix correction

Independent reviewer re-read the original card preview inventory (line 134),
Acceptance010 and Providers §4 and explicitly approved this bounded matrix:

| Native route | Required outcome |
| --- | --- |
| Supplied closed-subset SVG | Positive directly rendered resource proof with shared fit/clip geometry. |
| Supported Markdown prose | Positive rendered bounded prose proof. |
| Mermaid diagram on the current embedded-image route | Unsupported; freeze per-diagram label/reject after full identity validation, preserve surrounding prose, advertise no rendered-mermaid fact. |

Known unsupported backend does not bypass binding identity, copying, resource
validation/revision checks or zero-rerender evidence. Missing/reject policy retains
the old bundle; label policy records and displays the unavailable status before
gates. Legacy Markdown/global renderer APIs remain unchanged. There is no requirement
to add per-diagram native image placement, a generic scene or a runtime preview API
merely to close the inventory. The optional richer route discussion above is future
work only, not a dependency of this bounded handoff. Main has now recorded the matrix in the final canonical preimplementation
reconciliation section and still owns WCI4 selection; this memo is not code authority.


## Final canonical alignment and mounted accessibility boundary

Main's final ORIGINAL Providers-and-packaging preimplementation reconciliation is
now the governing disposition, superseding earlier tentative route alternatives.
Positive direct closed-subset SVG and bounded Markdown prose remain required.
The current native explicit-Markdown diagram route is unsupported: identity-first
checks, copied/budgeted resources, revision guards and declared per-diagram
label/reject complete before outcome/capability freeze. Preserve prose and never
claim a rendered-Mermaid fact. No new scene or native diagram-image API is needed.

Host must put Fyne's public Accessible interface on the actual mounted SVG and
explicit-Markdown wrappers, exposing the exact prepared description and outcome
status independently of caption. Inspection evidence must query that mounted
interface; reconstructing expected accessibility text from runtime/provider metadata
is insufficient. Per-diagram failures must remain represented in the wrapper's
outcome status while ordinary prose stays visible. This adds no runtime preview or
accessibility state: host projects the frozen record it already paints.

This is an adapter-interface and visible-label obligation, not a claim of OS
screen-reader delivery. The pinned Linux accessibility bridge is a no-op even with
its build tag; Linux OS delivery remains unsupported/unverified. No bridge expansion
or other-platform claim is selected. Product implementation and actual wrapper/native
proof remain future work under WCI4 selection; this memo records the required seam.

Coordination sent to Noether for the actual wrapper description/status projection.
Only this read-only memo changes; WCI4 remains unselected and M2 runtime frozen.

Final canonical observed SHA-256: `9d9e53bcae0f54f8f2a3c1bb1fd65633ef6a6c5403597b08b997d2324c4a706d`.


## Final host/layout text-overlay agreement

Noether and Gibbs agree on the single additional pure method:

```go
func (p *Previews) RenderSVGText(out *strings.Builder, box *layout.Box) error
```

It emits only caption/status glyph outlines through existing svg.Text. One private
band/line routine supplies MeasurePreview, SVGRects and this serializer. Existing
layout.TextWidth/Lines and svg.Text use the same embedded Go Regular with
HintingNone; native canvas.Text metrics are not substituted. The optional caption
uses one measured row with measured ellipsis where needed. The complete bounded
fallback status wraps in its own measured band and cannot be hidden by a caption.
Required text allocation is checked before publication. Exact full description
and diagnostic remain available through the frozen outcome and mounted Accessible
wrapper even when the optional caption is ellipsized.

Host keeps the direct resource image separate from a direct text-only SVG overlay.
SVGRects retains exactly three slots in the owning Box canvas coordinates, with
Rect{} for absent slots. The overlay root/viewBox uses that same allocation with
one coordinate translation; no resource image or data-URI image is embedded in it.
Image placement uses the already fitted rectangle, and both objects receive the
existing shared Box.Clip after placement. No second fit, native text-line protocol,
metrics callback, scene API or runtime preview state is introduced. Earlier direct
messages considering PreviewTextMetrics/WrapPreviewText/PreviewTextLine are
superseded; they are not proposed exports.

The host establishes support for this actual text-overlay representation before
outcome freeze and validates the complete serialized overlay before its ticket can
publish. Serialization errors reject that ticket without blanking live status or
reopening frozen resource fallback. This is pure painting of current caption and
geometry, not a call to the application's Renderer; repeated gates/resize still
make zero provider-renderer calls. Existing svg.Text propagates its font errors,
but does not explicitly reject glyph index zero: shared metrics alone is not proof
of arbitrary Unicode glyph fidelity. No expanded font-coverage claim or font work
is selected by this seam; actual supported output still needs inspected evidence.

The approved native matrix and accessibility boundary above remain unchanged.
Current explicit-Markdown diagrams resolve per-diagram label/reject before freeze;
no native Mermaid fact is claimed. Actual mounted wrappers supply Accessible;
Linux OS screen-reader delivery remains unsupported/unverified. Only this memo is
updated by this worker. No WCI4 source implementation, test run or selection follows
from this agreement, and frozen M2 runtime files remain unchanged.
