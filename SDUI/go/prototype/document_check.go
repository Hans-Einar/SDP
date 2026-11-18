package prototype

import (
	"fmt"

	"github.com/Hans-Einar/SDP/SDUI/go/host/fynehost/admission"
	"github.com/Hans-Einar/SDP/SDUI/go/markdown"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/preparation"
)

// CheckDocument is headless semantic preflight. Native renderer preparation and
// geometry admission remain pending until DocumentHost adopts the snapshot.
func CheckDocument(source, entry, revision, contextPath string) (Report, error) {
	r := Report{Schema: "sdptool/0.2", Operation: "sdui-check", Source: source, Entry: entry, Host: "document", NativeLayout: "pending"}
	input, err := LoadDocument(source, entry, revision, contextPath)
	if err != nil {
		return r, err
	}
	r.Revision = input.Candidate.Hash
	r.Profile = input.Candidate.Document.Profile
	r.ContextRevision = input.ContextRevision
	previews, err := markdown.PreparePreviews(input.Candidate.Root, input.SVGResources, nil, markdown.PreviewBackend{
		SVG:      func(resource markdown.Resource) error { _, err := markdown.ValidateSVGResource(resource); return err },
		Markdown: func() error { return nil },
		Mermaid:  func(markdown.Resource) error { return fmt.Errorf("native Markdown diagram images are unsupported") },
	})
	if err != nil {
		return r, err
	}
	caps := preflightCapabilities(input, previews)
	// Input/choice state was validated by LoadDocument. This does not invoke a
	// no-op geometry validator or advertise a successful native presentation.
	if err = preparation.CheckWithPreviews(input.Candidate.Document.Profile, input.Candidate.Root, caps, previews); err != nil {
		return r, err
	}
	if err = input.Guard(); err != nil {
		return r, err
	}
	r.Status = "preflight"
	r.Diagnostic = "Headless DocumentHost preflight only; native resource and layout admission pending; no SDL execution"
	return r, nil
}
func preflightCapabilities(input *DocumentInput, previews *markdown.Previews) preparation.Capabilities {
	caps := preparation.Capabilities{}
	for _, cap := range admission.TextCapabilities() {
		// The prototype supplies no application collection providers.
		if cap.Dimension == preparation.Provider && (cap.ID == "collection-data" || cap.ID == "collection-load") {
			continue
		}
		caps = append(caps, cap)
	}
	if len(input.Choices) > 0 {
		caps = append(caps, preparation.Capability{Dimension: preparation.Provider, ID: "choice-options", Major: 1})
	}
	svgReady, markdownReady := false, false
	input.Candidate.Root.Walk(func(n *parser.Instance) {
		if o, ok := previews.Outcome(n.Path); ok {
			svgReady = svgReady || o.Kind == "svg" && o.Status == "rendered"
			markdownReady = markdownReady || o.Kind == "markdown"
		}
	})
	if svgReady {
		caps = append(caps, preparation.Capability{Dimension: preparation.Widget, ID: "svg", Major: 1}, preparation.Capability{Dimension: preparation.Layout, ID: "preview-resource", Major: 1}, preparation.Capability{Dimension: preparation.Provider, ID: "svg-resource", Major: 1}, preparation.Capability{Dimension: preparation.Host, ID: "svg-resource", Major: 1})
	}
	if markdownReady {
		caps = append(caps, preparation.Capability{Dimension: preparation.Host, ID: "markdown", Major: 1})
	}
	return caps
}
