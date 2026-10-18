package fynehost

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/markdown"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/preparation"
	"github.com/Hans-Einar/SDP/SDUI/go/svg"
	"github.com/fyne-io/oksvg"
)

// Resource.Content does not expose the bundle's authoritative backing slice.
type previewResource struct {
	name string
	data []byte
}

func (r *previewResource) Name() string    { return r.name }
func (r *previewResource) Content() []byte { return append([]byte(nil), r.data...) }
func newPreviewResource(data []byte) fyne.Resource {
	return &previewResource{name: fmt.Sprintf("preview-%x.svg", sha256.Sum256(data)), data: append([]byte(nil), data...)}
}

func checkPreviewSVG(data []byte) error {
	_, err := oksvg.ReadIconStream(bytes.NewReader(data), oksvg.StrictErrorMode)
	return err
}
func previewTextBackend() error {
	var out strings.Builder
	out.WriteString(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 320 40">`)
	if err := svg.Text(&out, "Preview unavailable", 0, 20, 14, "#000000"); err != nil {
		return err
	}
	out.WriteString(`</svg>`)
	return checkPreviewSVG([]byte(out.String()))
}
func (b *Bundle) preparePreviews() error {
	if b.Document == nil {
		return fmt.Errorf("preparation: document required")
	}
	roots, err := parser.Normalize(b.Document)
	if err != nil {
		return err
	}
	root := roots[b.request.Entry]
	if root == nil || root.Kind != "frame" {
		return fmt.Errorf("preparation: entry %q must be a frame", b.request.Entry)
	}
	explicit := false
	root.Walk(func(n *parser.Instance) {
		policy, _ := parser.PreviewOptions(n) // Normalize already validated source.
		explicit = explicit || policy.Explicit
	})
	if explicit {
		// Every native outcome, including missing-resource fallback, needs this
		// actual glyph representation. An unavailable status renderer cannot label.
		if err := previewTextBackend(); err != nil {
			return err
		}
	}
	// These checks certify the actual direct image / background representation.
	backend := markdown.PreviewBackend{
		SVG:      func(r markdown.Resource) error { return checkPreviewSVG(r.SVG) },
		Markdown: previewTextBackend,
		Mermaid:  func(markdown.Resource) error { return fmt.Errorf("native Markdown diagram images are unsupported") },
	}
	b.previews, err = markdown.PreparePreviews(root, b.request.SVGResources, b.request.MarkdownRenderers, backend)
	b.request.SVGResources, b.request.MarkdownRenderers = nil, nil
	if err != nil {
		return err
	}
	b.previewResources = map[string]fyne.Resource{}
	b.previewOutcomes = map[string]markdown.PreviewOutcome{}
	resources := map[string]fyne.Resource{}
	root.Walk(func(n *parser.Instance) {
		o, ok := b.previews.Outcome(n.Path)
		if !ok {
			return
		}
		if o.Resource != nil {
			r := resources[o.SHA256]
			if r == nil {
				r = newPreviewResource(o.Resource.SVG)
				resources[o.SHA256] = r
			}
			b.previewResources[n.Path] = r
		}
		// Diagnostics need identity/status, not another retained document/byte tree.
		o.Resource, o.Document = nil, nil
		for i := range o.Diagrams {
			o.Diagrams[i].Resource = nil
		}
		b.previewOutcomes[n.Path] = o
	})
	return nil
}

func previewCapabilities(outcomes map[string]markdown.PreviewOutcome) preparation.Capabilities {
	var caps preparation.Capabilities
	svgReady, markdownReady := false, false
	for _, o := range outcomes {
		svgReady = svgReady || o.Kind == "svg" && o.Status == "rendered"
		markdownReady = markdownReady || o.Kind == "markdown"
	}
	if svgReady {
		caps = append(caps, preparation.Capabilities{
			{Dimension: preparation.Widget, ID: "svg", Major: 1},
			{Dimension: preparation.Layout, ID: "preview-resource", Major: 1},
			{Dimension: preparation.Provider, ID: "svg-resource", Major: 1},
			{Dimension: preparation.Host, ID: "svg-resource", Major: 1},
		}...)
	}
	if markdownReady {
		caps = append(caps, preparation.Capability{Dimension: preparation.Host, ID: "markdown", Major: 1})
	}
	return caps
}
func (b *Bundle) hasPreview(path string) bool { _, ok := b.previewOutcomes[path]; return ok }

func previewAccessible(o markdown.PreviewOutcome) string {
	if o.Diagnostic != "" {
		return "Preview unavailable: " + o.Description + " (" + o.Diagnostic + ")"
	}
	label := o.Description + " (" + o.Status + ")"
	for _, d := range o.Diagrams {
		if d.Diagnostic != "" {
			label += "\nPreview unavailable: " + o.Description + " (" + d.Diagnostic + ")"
		}
	}
	return label
}

func (b *Bundle) preparePreviewFrames(p *nativePresentation) error {
	p.previews = map[string]previewFrame{}
	var failure error
	walkPresentation(p, func(box *layout.Box, _ *layout.SnapshotLayout) {
		if failure != nil || !b.hasPreview(box.Path) {
			return
		}
		f := previewFrame{box: box.Rect}
		if box.Instance.Widget == "svg" {
			f.image, f.caption, f.status, failure = b.previews.SVGRects(box)
			if failure != nil {
				return
			}
			f.resource = b.previewResources[box.Path]
			if box.Rect.W > 0 && box.Rect.H > 0 && (f.caption.H > 0 || f.status.H > 0) {
				var out strings.Builder
				fmt.Fprintf(&out, `<svg xmlns="http://www.w3.org/2000/svg" width="%g" height="%g" viewBox="0 0 %g %g"><g transform="translate(%g %g)">`, box.Rect.W, box.Rect.H, box.Rect.W, box.Rect.H, -box.Rect.X, -box.Rect.Y)
				failure = b.previews.RenderSVGText(&out, box)
				if failure != nil {
					return
				}
				out.WriteString(`</g></svg>`)
				data := []byte(out.String())
				if failure = checkPreviewSVG(data); failure != nil {
					return
				}
				f.text = newPreviewResource(data)
			}
		}
		p.previews[box.Path] = f
	})
	return failure
}
func (b *Bundle) closePreviews() {
	if b.view != nil {
		for _, control := range b.view.Controls {
			if c, ok := control.(*previewControl); ok {
				c.release()
			}
		}
	}
	if b.presentation != nil {
		b.presentation.previews = nil
	}
	if b.pending != nil {
		b.pending.previews = nil
	}
	for _, p := range b.retiredPresentations {
		p.previews = nil
	}
	b.previewOutcomes, b.previewResources, b.previews = nil, nil, nil
	b.request.SVGResources, b.request.MarkdownRenderers = nil, nil
}
