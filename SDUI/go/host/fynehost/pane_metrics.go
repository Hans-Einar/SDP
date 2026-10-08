package fynehost

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"

	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

func headerPages(state ui.TabsState, label string) []headerPage {
	out := []headerPage{}
	for _, p := range state.Pages {
		if p.Visible {
			out = append(out, headerPage{id: p.ID, label: p.Label, enabled: p.Enabled})
		}
	}
	if len(out) == 0 {
		out = append(out, headerPage{label: label})
	}
	return out
}
func (m *collectionMeasure) MeasureTabs(n *parser.Instance, font float64, _ layout.Size) (layout.TabsMetrics, error) {
	state, ok := m.snapshot.Tabs[n.Path]
	if !ok {
		return layout.TabsMetrics{}, fmt.Errorf("pane-state: missing tabs %s", n.Path)
	}
	h := newPaneHeader()
	th := container.NewThemeOverride(h, componentTheme{float32(font)})
	h.sync(headerPages(state, n.Argument("label")), state.Selected)
	minimum := th.MinSize()
	// Measure all native headers at once so a popup overflow menu is not required
	// for the M1 header keyboard and stable-ID hit rectangle contract.
	th.Resize(fyne.NewSize(1e6, minimum.Height))
	width := float32(0)
	if r := h.tabs.renderer; r != nil {
		objects := r.Objects()
		if len(objects) > 0 {
			width = objects[0].MinSize().Width
		}
	}
	return layout.TabsMetrics{Header: layout.Size{W: float64(width), H: float64(minimum.Height)}}, nil
}
func (m *collectionMeasure) MeasureSplit(_ *parser.Instance, font float64) (float64, error) {
	native := container.NewHSplit(container.NewWithoutLayout(), container.NewWithoutLayout())
	wrapped := container.NewThemeOverride(native, componentTheme{float32(font)})
	return float64(wrapped.MinSize().Width), nil
}
