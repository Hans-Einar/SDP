# Markdown and prepared previews

The legacy `Prepare(root, renderer)`, `Provider`, `Parse`, `Mmdr` and
`WriteResources` APIs retain their existing behavior. Bare Markdown and legacy SVG
calls do not opt into the explicit preview contract. Source parsing never opens
resource references or selects executables.

`PreparePreviews(root, resources, renderers, backend)` prepares an entire selected
normalized root once. `PreparedSVG` and `MarkdownRenderer` maps use exact instance
paths, including reused, hidden and closed-surface declarations. Inputs carry
source/provider/digest or renderer/provider/revision identities. Missing bindings
follow the declaration's `label`/`reject` policy; extra, wrong-kind, unused, nil,
source-mismatched or digest-mismatched bindings are fatal before callbacks/fallback.
No renderer or backend callback is retained by the returned `Previews`.

All supplied SVG inputs are validated/copied one bounded resource at a time before
application callbacks. Registered diagrams receive only their fenced source.
Custom output and backend-check arguments are copied. `Outcome(path)` returns deep
copies, including resource bytes, parsed blocks/table cells and provenance. Equal
Markdown text at distinct paths never merges descriptions, policies or outcomes.
Repeated fences retain their individual prepared outcomes. Successful validated
bytes can share by their full SHA-256 digest.

The 4 MiB per-resource bound precedes content validation. The 32 MiB candidate
admission tally counts distinct copied and validated payloads before backend
classification, including subsequent labelled failures. Identical bytes count once
across SVG and diagram paths. Invalid content does not enter the aggregate tally.
Backend failure does not refund it. This is an admission limit, not a claim about
live retained memory: labelled outcomes release the failed resource bytes and keep
identity/digest plus bounded diagnostics. Aggregate overflow is always fatal.

`ValidateSVGResource` admits one complete, finite SVG document with matching
positive supplied/viewBox dimensions at most 32768. The supplied subset is
svg/g/path/rect/circle/ellipse/line/polyline/polygon/title/desc, supported finite
geometry/path commands/transforms, solid paint, stroke width and opacity. Unknown
attributes/elements, CSS, external references, nested SVG documents, active content,
DTD/directives and extra roots/trailing data reject. An initial UTF-8 XML declaration
is allowed. No URL, image, font or entity acquisition occurs.

`ValidateMermaidResource` is a separate registered-output profile: it additionally
admits bounded text/tspan and local defs/markers with explicitly checked attributes.
Both validators return copied bytes. Existing legacy `Mmdr` behavior is unchanged;
its process has the existing five-second/4 MiB limits. Arbitrary Go renderers remain
application-owned calls; this API adds no timeout/cancellation service around them.
Source is limited to flowchart/graph, 12000 bytes and no configuration directive.
Explicit Markdown retains 32768 bytes, 256 blocks and eight diagrams.

`PreviewBackend` has three fixed preparation-only checks: SVG, Markdown and Mermaid.
Nil means unavailable; a check must certify the actual composed representation.
The selected native route supports direct closed-subset SVG and bounded Markdown
prose. Its composite Markdown diagram route is unavailable: freeze each diagram's
label/reject outcome while preserving ordinary prose. Do not advertise rendered
Mermaid merely because the inner SVG validated. Unsupported document content or a
whole-Markdown backend failure labels/rejects the whole explicit preview.

After preparation, `Check(root)` verifies exact explicit inventory and source
fingerprints. `CheckPreview(node)` verifies one explicit node. Caption, layout and
visibility can change without new content preparation; text/source/description/
policy changes require a new candidate. The host's application Guard separately
rechecks provider/digest/renderer revision before publication. Resource outcome
state belongs to the prepared bundle, never runtime state.

`Measure`/`Render` project frozen per-path Markdown; legacy content uses the old
provider projection. `MeasurePreview` returns natural/minimum sizes for explicit
SVG. Explicit Markdown label/partial rendering also checks the complete frozen
projection at its final width and rejects insufficient box allocation before writing
the output fragment. Fully rendered prose and legacy paths keep their behavior.
The scalable image minimum is zero; required caption/status text is measured.
`SVGRects` returns fixed image/caption/status rectangles in owning-box canvas
coordinates. It fits the full image allocation once; the caller applies Box.Clip
afterward. The caption occupies one ellipsized row below the image. Full bounded
fallback status wraps in its own region; insufficient allocation rejects.

`RenderSVGText` uses the same band/line plan and embedded Go Regular as measurement.
It emits only caption/status glyph paths, never an embedded resource image. The
host uses a separate direct image plus text-only SVG overlay and validates the final
representation before its accepted ticket publishes. Measurement/paint/resize/gates
never invoke application renderers or reconsider fallback. Existing glyph coverage
is unchanged; this is not an arbitrary-Unicode fidelity promise.

Native owners expose exact description/status through the actual mounted Fyne
Accessible wrapper independently of caption. This package provides data, not a GUI
or OS accessibility bridge. Pinned Linux OS screen-reader delivery is unsupported/
unverified. Public SVG export currently rejects explicit previews; native admission
and public export are separate contracts.

Run component evidence from SDUI/go with
`GOWORK=off go test -mod=readonly -race ./markdown`.
Provider tests do not prove native pixel output, OS accessibility, publication or
package installation; those remain integration/owner verification.
