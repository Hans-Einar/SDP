// Package admission describes the concrete Fyne 0.2 adapter without importing
// Fyne, so standalone prototype preflight uses the same provider/layout facts.
package admission

import (
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/markdown"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/preparation"
)

// Capabilities describes actual controls and the bounded Markdown provider with
// symbolic SVG placeholders. It advertises neither vector execution nor scrolling.
func Capabilities() preparation.Capabilities {
	return preparation.Capabilities{
		{Dimension: preparation.Frontend, ID: "sdui/0.2", Major: 1},
		{Dimension: preparation.Layout, ID: "relative", Major: 1},
		{Dimension: preparation.Widget, ID: "button", Major: 1},
		{Dimension: preparation.Widget, ID: "input", Major: 1},
		{Dimension: preparation.Widget, ID: "svg-placeholder", Major: 1},
		{Dimension: preparation.Provider, ID: "svg-placeholder", Major: 1},
		{Dimension: preparation.Provider, ID: "markdown", Major: 1},
		{Dimension: preparation.Host, ID: "button", Major: 1},
		{Dimension: preparation.Host, ID: "input", Major: 1},
	}
}
func Layout(root *parser.Instance, size layout.Size) error {
	provider, err := markdown.Prepare(root, nil)
	if err != nil {
		return err
	}
	_, err = (&layout.Engine{Measure: provider}).Layout(root, size)
	return err
}
func Check(root *parser.Instance) error {
	profile, err := parser.EffectiveProfile(root)
	if err != nil {
		return err
	}
	return preparation.Check(profile, root, Capabilities())
}

// CollectionCapabilities is advertised only by the document-based bundle host,
// which supplies actual provider bindings, state-aware geometry and publication.
func CollectionCapabilities() preparation.Capabilities {
	return append(Capabilities(), preparation.Capabilities{
		{Dimension: preparation.Frontend, ID: "sdui/0.3", Major: 1},
		{Dimension: preparation.Layout, ID: "collections", Major: 1},
		{Dimension: preparation.Widget, ID: "tree", Major: 1}, {Dimension: preparation.Widget, ID: "list", Major: 1},
		{Dimension: preparation.Provider, ID: "collection-data", Major: 1}, {Dimension: preparation.Provider, ID: "collection-load", Major: 1},
		{Dimension: preparation.Host, ID: "tree", Major: 1}, {Dimension: preparation.Host, ID: "list", Major: 1},
		{Dimension: preparation.Host, ID: "viewport", Major: 1}, {Dimension: preparation.Host, ID: "atomic-publication", Major: 1},
		{Dimension: preparation.Viewport, ID: "scroll-x", Major: 1}, {Dimension: preparation.Viewport, ID: "scroll-y", Major: 1},
	}...)
}
