package fynehost

import (
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
)

// MeasureViewport reserves stable native tracks on declared frame/group axes.
// Layout owns padding, finite child references and ancestor clipping; collection
// controls continue to report their own title/gutters through MeasureCollection.
func (m *collectionMeasure) MeasureViewport(n *parser.Instance, _ float64) (layout.ViewportInsets, error) {
	insets := layout.ViewportInsets{}
	if n.Profile != "sdui/0.3" || (n.Kind != "frame" && n.Kind != "group") {
		return insets, nil
	}
	if n.Layout["overflow-y"] == "scroll" {
		insets.Right = collectionGutter
	}
	if n.Layout["overflow-x"] == "scroll" {
		insets.Bottom = collectionGutter
	}
	return insets, nil
}
