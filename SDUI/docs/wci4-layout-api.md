# WCI4 layout preview API

Selected implementation under Providers-and-packaging.md and the final four WCI4
API memos. This document describes the layout seam; native acceptance belongs to
the integrated candidate.

```go
type PreviewMeasurer interface {
    MeasurePreview(*parser.Instance, float64, float64) (natural Size, minimum Size, err error)
}
func FitPreview(content Rect, intrinsic Size) (Rect, error)
```

The arguments are the exact normalized SVG instance, inherited font and available
width, matching `Measurer`. The composed measurer must implement both interfaces.
Only explicit SVG uses this adjunct, as determined by `parser.PreviewOptions`.
Legacy SVG and Markdown use `Measure`; explicit Markdown uses the same existing
method on its frozen per-path provider. No runtime preview state is introduced.

Measurement consumes immutable prepared outcomes and performs no resource lookup,
renderer invocation, fallback reconsideration, native mutation or I/O. The layout
engine checks finite nonnegative bounded natural/minimum sizes independently;
it imposes no `natural >= minimum` relationship. Zero image minimum is legal.
Caption/status minimums belong to the provider's shared text-band plan. Desired
allocation, pane minimum solving and final arrangement all consume the adjunct.
Source minima/maxima, relative tracks and split constraints retain their existing
meaning. Failed geometry returns an error before state/ticket publication.

FitPreview returns centered aspect-fit geometry, including enlargement. Its input
is the full caption-excluded image allocation in owning canvas coordinates, not a
visible clipped fragment. Nonpositive intrinsic dimensions, negative/nonfinite or
unbounded sizes, nonfinite positions and unrepresentable fits return the diagnostic
`preview-geometry`. Sizes use the existing layout bound of 1e7; resource limits and
byte validation remain the provider's responsibility. Valid zero-area allocation
returns exactly `Rect{}`. Apply `Box.Clip` after fitting, with no second translation.

Missing active preview metrics and invalid measurement sizes return
`preview-measurement`; final allocation below required minimum returns
`native-minimum`. Provider errors propagate. Inactive geometry is not measured,
but invalid preview source policy cannot hide in a closed surface or inactive page.
Resource identity/preparation for those declarations is a separate provider check.

The provider's `SVGRects` owns its fixed image/caption/status rectangles and calls
FitPreview once for the image. `RenderSVGText` owns matching text glyph placement.
Layout imports neither Markdown nor GUI packages and adds no font callback or
preview offset. Native adapters must consume the returned rectangles and existing
ancestor clips; they must not fit the clipped fragment again.
