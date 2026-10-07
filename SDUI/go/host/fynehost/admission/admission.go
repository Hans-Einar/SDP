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
func Check(root *parser.Instance) error { return preparation.Check("sdui/0.2", root, Capabilities()) }
