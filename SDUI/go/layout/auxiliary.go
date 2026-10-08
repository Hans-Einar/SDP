package layout

import "github.com/Hans-Einar/SDP/SDUI/go/parser"

// Auxiliary declarations remain in the normalized tree for preflight, but never
// reserve parent tracks or gaps. An open dialog is measured in its own canvas.
func auxiliary(n *parser.Instance) bool { return parser.IsAuxiliary(n) }
func inFlow(n *parser.Instance) bool    { return visible(n) && !auxiliary(n) }
func menuMode(n *parser.Instance) string {
	mode := n.Argument("mode")
	if mode == "" {
		return "bar"
	}
	return mode
}
func menuBar(n *parser.Instance) bool {
	return n.Profile == "sdui/0.3" && n.Kind == "composition" && n.Widget == "menu" && menuMode(n) == "bar"
}
func dialog(n *parser.Instance) bool {
	return n.Profile == "sdui/0.3" && n.Kind == "composition" && n.Widget == "dialog"
}

// MenuMeasurer returns actual native bar-button minimum dimensions. Popup item,
// separator/group and submenu geometry belongs to the native menu adapter.
type MenuMeasurer interface {
	MeasureMenu(*parser.Instance, float64) (Size, error)
}

func (e *Engine) menuMinimum(n *parser.Instance, font float64) (Size, error) {
	m, ok := e.Measure.(MenuMeasurer)
	if !ok {
		return Size{}, diag(n, "menu-measurement", "Menu bars require native adapter metrics")
	}
	size, err := m.MeasureMenu(n, font)
	if err != nil {
		return Size{}, err
	}
	if !finiteExtent(size) || size.W <= 0 || size.H <= 0 {
		return Size{}, diag(n, "menu-measurement", "Menu minimum must be finite, positive and bounded")
	}
	return size, nil
}
