package markdown

import (
	"fmt"
	"strings"

	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
)

// Check verifies the entire exact explicit inventory against the selected root.
// Mutable caption, geometry and visibility do not form resource identity.
func (p *Previews) Check(root *parser.Instance) error {
	if _, e := parser.ResolveInteractions(root); e != nil {
		return e
	}
	inputs, e := previewInventory(root)
	if e != nil {
		return e
	}
	if p == nil {
		if len(inputs) == 0 {
			return nil
		}
		return fmt.Errorf("preview-unprepared: missing prepared previews")
	}
	if len(inputs) != len(p.records) {
		return fmt.Errorf("preview-stale: explicit inventory changed")
	}
	for _, in := range inputs {
		if p.fingerprints[in.outcome.Path] != in.fingerprint {
			return previewError(in.outcome, "preview-stale", "prepared source fingerprint mismatch")
		}
	}
	return nil
}
func (p *Previews) CheckPreview(n *parser.Instance) error {
	policy, e := parser.PreviewOptions(n)
	if e != nil {
		return e
	}
	if !policy.Explicit {
		return fmt.Errorf("preview-unprepared: expected explicit preview")
	}
	inputs, e := previewInventory(n)
	if e != nil {
		return e
	}
	if len(inputs) != 1 || p == nil {
		return fmt.Errorf("preview-unprepared: missing prepared preview")
	}
	in := inputs[0]
	if p.fingerprints[in.outcome.Path] != in.fingerprint {
		return previewError(in.outcome, "preview-stale", "prepared source fingerprint mismatch")
	}
	return nil
}
func unavailable(description, why string) string {
	return "Preview unavailable: " + description + " (" + why + ")"
}
func projectMarkdown(o PreviewOutcome, source string) *Provider {
	d := copyDocument(o.Document)
	resources := map[string]Resource{}
	if d == nil {
		d = &Document{Blocks: []Block{{Text: unavailable(o.Description, o.Diagnostic), Scale: 1}}}
	} else {
		diagramIndex := 0
		for i := range d.Blocks {
			block := &d.Blocks[i]
			if block.Diagram == "" {
				continue
			}
			result := o.Diagrams[diagramIndex]
			diagramIndex++
			if result.Resource == nil {
				*block = Block{Text: unavailable(o.Description, result.Diagnostic), Scale: 1}
			} else {
				// Repeated identical fences may produce different per-occurrence
				// results. Never merge their policy results through the legacy ID map.
				key := fmt.Sprintf("%s-%d", result.ID, diagramIndex)
				block.Diagram = key
				resources[key] = *result.Resource
			}
		}
	}
	return &Provider{Documents: map[string]*Document{source: d}, Resources: resources}
}
func (p *Previews) Measure(n *parser.Instance, font, width float64) (layout.Size, error) {
	if e := previewMetrics(font, width); e != nil {
		return layout.Size{}, e
	}
	if n == nil {
		return layout.Size{}, fmt.Errorf("preview-unprepared: nil node")
	}
	if n.Kind == "markdown" || n.Widget == "svg" {
		policy, e := parser.PreviewOptions(n)
		if e != nil {
			return layout.Size{}, e
		}
		if policy.Explicit {
			if e = p.CheckPreview(n); e != nil {
				return layout.Size{}, e
			}
			if n.Widget == "svg" {
				natural, _, e := p.MeasurePreview(n, font, width)
				return natural, e
			}
			return p.projections[n.Path].Measure(n, font, width)
		}
	}
	if p == nil || p.legacy == nil {
		return layout.Size{}, fmt.Errorf("preview-unprepared: missing provider")
	}
	return p.legacy.Measure(n, font, width)
}
func (p *Previews) Render(out *strings.Builder, box *layout.Box) error {
	if out == nil || box == nil || box.Instance == nil {
		return fmt.Errorf("preview-unprepared: nil render input")
	}
	if e := previewMetrics(box.Font, box.Rect.W); e != nil {
		return e
	}
	policy, e := parser.PreviewOptions(box.Instance)
	if e != nil {
		return e
	}
	if !policy.Explicit {
		if p == nil || p.legacy == nil {
			return fmt.Errorf("preview-unprepared: missing provider")
		}
		return p.legacy.Render(out, box)
	}
	if e = p.CheckPreview(box.Instance); e != nil {
		return e
	}
	if box.Instance.Widget == "svg" {
		return fmt.Errorf("preview-render: resource SVG requires direct image and RenderSVGText")
	}
	if p.records[box.Instance.Path].Status != "rendered" {
		// Markdown source may assign a smaller box than its desired size. A
		// labelled/partial outcome must not publish a clipped-away explanation.
		_, needed, err := p.projections[box.Instance.Path].items(box.Instance.Text, box.Font, box.Rect.W)
		if err != nil {
			return err
		}
		if !finite(box.Rect.H) || box.Rect.H < 0 || box.Rect.W+1e-7 < needed.W || box.Rect.H+1e-7 < needed.H {
			return fmt.Errorf("preview-geometry: Markdown fallback allocation is too small")
		}
	}
	// Build privately so a failed serializer never leaves a half-written fragment.
	var fragment strings.Builder
	if e = p.projections[box.Instance.Path].Render(&fragment, box); e != nil {
		return e
	}
	out.WriteString(fragment.String())
	return nil
}
func previewMetrics(font, width float64) error {
	if !finite(font) || font <= 0 || font > 512 || !finite(width) || width < 0 || width > 1e7 {
		return fmt.Errorf("preview-geometry: invalid font or available width")
	}
	return nil
}
