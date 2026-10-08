package markdown

import "github.com/Hans-Einar/SDP/SDUI/go/parser"

// PreparedSVG binds an application's immutable resource identity to an exact
// normalized instance path. PreparePreviews copies the bytes before callbacks.
type PreparedSVG struct {
	Source             parser.Reference
	ProviderID, SHA256 string
	Resource           Resource
}

// MarkdownRenderer is an explicit preparation-only application registration.
type MarkdownRenderer struct {
	ProviderID, Revision string
	Renderer             Renderer
}

// PreviewBackend checks the actual composed destination representation. Nil
// means unavailable. Checks receive copies and are never retained by Previews.
type PreviewBackend struct {
	SVG      func(Resource) error
	Markdown func() error
	Mermaid  func(Resource) error
}

type DiagramOutcome struct {
	ID, SHA256 string
	Resource   *Resource
	Diagnostic string
}

type PreviewOutcome struct {
	Path, Kind, Description, Fallback string
	Span                              parser.Span
	Uses                              []parser.UseSite
	Source                            parser.Reference
	ProviderID, Revision, SHA256      string
	Status                            string // rendered, partial, label
	Diagnostic                        string
	Resource                          *Resource
	Document                          *Document
	Diagrams                          []DiagramOutcome
}

// Previews owns frozen per-instance outcomes. All methods after preparation are
// pure; no renderer, backend callback, caller map or caller byte slice is retained.
type Previews struct {
	records      map[string]PreviewOutcome
	fingerprints map[string]string
	projections  map[string]*Provider
	legacy       *Provider
}

// Outcome returns a detached copy, including resource bytes and provenance.
func (p *Previews) Outcome(path string) (PreviewOutcome, bool) {
	if p == nil {
		return PreviewOutcome{}, false
	}
	o, ok := p.records[path]
	if !ok {
		return PreviewOutcome{}, false
	}
	o.Uses = append([]parser.UseSite(nil), o.Uses...)
	if o.Resource != nil {
		r := copyResource(*o.Resource)
		o.Resource = &r
	}
	o.Document = copyDocument(o.Document)
	o.Diagrams = append([]DiagramOutcome(nil), o.Diagrams...)
	for i := range o.Diagrams {
		if o.Diagrams[i].Resource != nil {
			r := copyResource(*o.Diagrams[i].Resource)
			o.Diagrams[i].Resource = &r
		}
	}
	return o, true
}
func copyResource(r Resource) Resource { r.SVG = append([]byte(nil), r.SVG...); return r }
func copyDocument(d *Document) *Document {
	if d == nil {
		return nil
	}
	v := &Document{Blocks: append([]Block(nil), d.Blocks...), Diagrams: append([]Diagram(nil), d.Diagrams...)}
	for i := range v.Blocks {
		v.Blocks[i].Table = append([][]string(nil), v.Blocks[i].Table...)
		for j := range v.Blocks[i].Table {
			v.Blocks[i].Table[j] = append([]string(nil), v.Blocks[i].Table[j]...)
		}
	}
	return v
}
