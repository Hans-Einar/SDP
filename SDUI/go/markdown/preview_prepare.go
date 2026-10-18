package markdown

import (
	"fmt"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"sort"
)

// PreparePreviews validates the entire binding inventory before any content
// fallback or callback. It performs registered renderer I/O only during this call.
func PreparePreviews(root *parser.Instance, resources map[string]PreparedSVG, renderers map[string]MarkdownRenderer, backend PreviewBackend) (*Previews, error) {
	if _, e := parser.ResolveInteractions(root); e != nil {
		return nil, e
	}
	inputs, e := previewInventory(root)
	if e != nil {
		return nil, e
	}
	byPath := map[string]previewInput{}
	for _, in := range inputs {
		byPath[in.outcome.Path] = in
	}
	resourceInputs := map[string]PreparedSVG{}
	rendererInputs := map[string]MarkdownRenderer{}
	// Sorting makes fatal diagnostics deterministic without depending on map order.
	keys := make([]string, 0, len(resources))
	for path := range resources {
		keys = append(keys, path)
	}
	sort.Strings(keys)
	for _, path := range keys {
		r := resources[path]
		in, ok := byPath[path]
		if !ok || in.outcome.Kind != "svg" {
			return nil, fmt.Errorf("preview-identity: extra/wrong SVG binding %s", path)
		}
		if !sameReference(r.Source, in.outcome.Source) || !identity(r.ProviderID) || !validDigest(r.SHA256) || digest(r.Resource.SVG) != r.SHA256 {
			return nil, previewError(in.outcome, "preview-identity", "invalid source/provider/digest binding")
		}
		resourceInputs[path] = r
	}
	keys = keys[:0]
	for path := range renderers {
		keys = append(keys, path)
	}
	sort.Strings(keys)
	for _, path := range keys {
		r := renderers[path]
		in, ok := byPath[path]
		if !ok || in.outcome.Kind != "markdown" || !hasDiagram(in.text) {
			return nil, fmt.Errorf("preview-identity: extra/wrong/unused Markdown binding %s", path)
		}
		if !identity(r.ProviderID) || !identity(r.Revision) || nilRenderer(r.Renderer) {
			return nil, previewError(in.outcome, "preview-identity", "invalid renderer identity or nil renderer")
		}
		rendererInputs[path] = r
	}
	// Validate and copy all supplied bytes before callbacks, one bounded resource
	// at a time. Do not stage an unbounded second copy of the application map.
	budget := previewBudget{resources: map[string]Resource{}}
	resourceErrors := map[string]error{}
	for _, in := range inputs {
		path := in.outcome.Path
		r, ok := resourceInputs[path]
		if !ok {
			continue
		}
		v, err := ValidateSVGResource(r.Resource)
		if err == nil {
			v, err = budget.keep(v)
			if err != nil {
				return nil, previewError(in.outcome, "preview-budget", err.Error())
			}
		} else {
			resourceErrors[path] = err
		}
		r.Resource = v
		resourceInputs[path] = r
	}

	p := &Previews{records: map[string]PreviewOutcome{}, fingerprints: map[string]string{}, projections: map[string]*Provider{}, legacy: &Provider{Documents: map[string]*Document{}, Resources: map[string]Resource{}}}
	var legacyErr error
	root.Walk(func(n *parser.Instance) {
		if legacyErr != nil || n.Kind != "markdown" || len(n.Arguments) != 0 {
			return
		}
		if _, ok := p.legacy.Documents[n.Text]; ok {
			return
		}
		d, e := Parse(n.Text)
		if e != nil {
			legacyErr = e
			return
		}
		p.legacy.Documents[n.Text] = d
	})
	if legacyErr != nil {
		return nil, legacyErr
	}
	for _, in := range inputs {
		o := in.outcome
		o.Status = "rendered"
		if o.Kind == "svg" {
			r, ok := resourceInputs[o.Path]
			if !ok {
				e = fmt.Errorf("resource provider unavailable")
			} else {
				o.ProviderID, o.SHA256 = r.ProviderID, r.SHA256
				v := r.Resource
				e = resourceErrors[o.Path]
				if e == nil {
					if backend.SVG == nil {
						e = fmt.Errorf("SVG backend unavailable")
					} else {
						e = backend.SVG(copyResource(v))
					}
					if e == nil {
						o.Resource = &v
					}
				}
			}
			if e != nil {
				if e = labelOutcome(&o, e); e != nil {
					return nil, e
				}
			}
		} else {
			if r, ok := rendererInputs[o.Path]; ok {
				o.ProviderID, o.Revision = r.ProviderID, r.Revision
			}
			d, err := Parse(in.text)
			if err == nil && len(d.Blocks) > 256 {
				err = fmt.Errorf("markdown-limit: more than 256 blocks")
			}
			if err != nil {
				if e = labelOutcome(&o, err); e != nil {
					return nil, e
				}
			} else {
				o.Document = d
				for _, diagram := range d.Diagrams {
					result := DiagramOutcome{ID: diagram.ID}
					r, ok := rendererInputs[o.Path]
					var v Resource
					if !ok {
						err = fmt.Errorf("diagram renderer unavailable")
					} else if err = diagramSource(diagram.Source); err == nil {
						v, err = r.Renderer.Render(diagram.Source)
						if err == nil {
							v, err = ValidateMermaidResource(v)
						}
					}
					if err == nil {
						v, err = budget.keep(v)
						if err != nil {
							return nil, previewError(o, "preview-budget", err.Error())
						}
						result.SHA256 = digest(v.SVG)
						if backend.Mermaid == nil {
							err = fmt.Errorf("diagram backend unavailable")
						} else {
							err = backend.Mermaid(copyResource(v))
						}
						if err == nil {
							result.Resource = &v
						}
					}
					if err != nil {
						if o.Fallback == "reject" {
							return nil, previewError(o, "preview-content", err.Error())
						}
						result.Diagnostic = reason(err)
						o.Status = "partial"
					}
					o.Diagrams = append(o.Diagrams, result)
				}
			}
			// Even an unavailable whole-Markdown backend cannot waive validation
			// and the admission tally for successful registered diagram output.
			if o.Document != nil {
				if backend.Markdown == nil {
					err = fmt.Errorf("Markdown backend unavailable")
				} else {
					err = backend.Markdown()
				}
				if err != nil {
					if e = labelOutcome(&o, err); e != nil {
						return nil, e
					}
					o.Document = nil
					for i := range o.Diagrams {
						o.Diagrams[i].Resource = nil
					}
				}
			}
			p.projections[o.Path] = projectMarkdown(o, in.text)
		}
		p.records[o.Path] = o
		p.fingerprints[o.Path] = in.fingerprint
	}
	return p, nil
}

type previewBudget struct {
	resources map[string]Resource
	bytes     int
}

func (b *previewBudget) keep(r Resource) (Resource, error) {
	key := digest(r.SVG)
	if old, ok := b.resources[key]; ok {
		return old, nil
	}
	if len(r.SVG) > candidateLimit-b.bytes {
		return Resource{}, fmt.Errorf("prepared distinct resources exceed 32 MiB")
	}
	b.bytes += len(r.SVG)
	b.resources[key] = r
	return r, nil
}
func labelOutcome(o *PreviewOutcome, e error) error {
	if o.Fallback == "reject" {
		return previewError(*o, "preview-content", e.Error())
	}
	o.Status = "label"
	o.Diagnostic = reason(e)
	return nil
}
