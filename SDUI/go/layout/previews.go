package layout

import "github.com/Hans-Einar/SDP/SDUI/go/parser"

// PreviewMeasurer is the optional geometry adjunct for explicitly opted-in SVG.
// It reads immutable prepared outcomes; it must not resolve resources, invoke a
// renderer or change state. Arguments match Measurer. A scalable image can have
// zero minimum size; caption/status text and source constraints still apply.
type PreviewMeasurer interface {
	MeasurePreview(*parser.Instance, float64, float64) (natural Size, minimum Size, err error)
}

func (e *Engine) explicitPreview(n *parser.Instance) (bool, error) {
	if e.profile != "sdui/0.3" || n.Kind != "widget" || n.Widget != "svg" {
		return false, nil
	}
	p, err := parser.PreviewOptions(n)
	return p.Explicit, err
}

func (e *Engine) previewMetrics(n *parser.Instance, font, width float64) (Size, Size, error) {
	m, ok := e.Measure.(PreviewMeasurer)
	if !ok {
		return Size{}, Size{}, diag(n, "preview-measurement", "Explicit SVG layout requires prepared preview metrics")
	}
	if !finite(font) || font <= 0 || !finiteExtent(Size{width, 0}) {
		return Size{}, Size{}, diag(n, "preview-measurement", "Preview measurement requires finite font and available width")
	}
	natural, minimum, err := m.MeasurePreview(n, font, width)
	if err != nil {
		return Size{}, Size{}, err
	}
	if !finiteExtent(natural) || !finiteExtent(minimum) {
		return Size{}, Size{}, diag(n, "preview-measurement", "Preview sizes must be finite, nonnegative and within layout bounds")
	}
	return natural, minimum, nil
}

// FitPreview centers intrinsic content in the full allocated rectangle while
// preserving its aspect ratio, including upscaling. Clip the result afterward;
// fitting a clipped fragment would change the image as its ancestor scrolls.
// Zero-area allocation returns Rect{}. Negative/nonfinite sizes, nonpositive
// intrinsic dimensions and sizes exceeding the layout bound (1e7) reject.
// Resource-byte validation and preparation limits belong to the provider.
func FitPreview(content Rect, intrinsic Size) (Rect, error) {
	invalid := func() (Rect, error) {
		return Rect{}, &parser.Diagnostic{Code: "preview-geometry", Message: "Preview fit requires finite bounded geometry and positive intrinsic dimensions"}
	}
	if !finite(content.X) || !finite(content.Y) || !finite(content.X+content.W) || !finite(content.Y+content.H) || !finiteExtent(Size{content.W, content.H}) || !finiteExtent(intrinsic) || intrinsic.W <= 0 || intrinsic.H <= 0 {
		return invalid()
	}
	if content.W == 0 || content.H == 0 {
		return Rect{}, nil
	}
	// Normalize along the larger intrinsic axis. This avoids overflowing a
	// scale factor for very small (but finite positive) intrinsic dimensions.
	w, h := content.W, content.H
	if intrinsic.W >= intrinsic.H {
		ratio := intrinsic.H / intrinsic.W
		h = w * ratio
		if h > content.H {
			h = content.H
			w = h / ratio
		}
	} else {
		ratio := intrinsic.W / intrinsic.H
		w = h * ratio
		if w > content.W {
			w = content.W
			h = w / ratio
		}
	}
	r := Rect{content.X + (content.W-w)/2, content.Y + (content.H-h)/2, w, h}
	if !finite(r.X) || !finite(r.Y) || !finiteExtent(Size{w, h}) || w <= 0 || h <= 0 {
		return invalid()
	}
	return r, nil
}
