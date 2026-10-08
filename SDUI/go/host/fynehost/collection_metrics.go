package fynehost

import (
	"fmt"
	"fyne.io/fyne/v2"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/markdown"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"math"
	"strings"
)

// Pinned Fyne Tree/List keep scroll callbacks/offsets private and their keyboard
// selection policy cannot skip group/separator rows. This bounded adapter uses
// Fyne text metrics and controls, while runtime owns all selection and offsets.
const collectionGutter = 12.0

type collectionRow struct {
	Item     ui.CollectionItem
	Parent   ui.ItemID
	Text     string
	Depth    int
	Y, H     float64
	Recovery bool
	Status   bool
}
type collectionMeasure struct {
	snapshot ui.Snapshot
	markdown *markdown.Provider
	icons    map[string]fyne.Resource
}

func rowHeight(font float64) float64 {
	return float64(nativeText("Mg", font).Height) + 8
}
func collectionRows(state ui.CollectionState, font float64) []collectionRow {
	rows := []collectionRow{}
	y := 0.0
	height := rowHeight(font)
	status := func(parent ui.ItemID, depth int) {
		st := state.Status[parent]
		text := ""
		recovery := false
		switch st.Phase {
		case ui.Loading:
			text = "Loading…"
		case ui.LoadError:
			text = "Retry/R: " + st.Error
			recovery = true
		case ui.Canceled:
			text = "Load — R"
			recovery = true
		}
		if text != "" {
			rows = append(rows, collectionRow{Parent: parent, Text: text, Depth: depth, Y: y, H: height, Recovery: recovery, Status: true})
			y += height
		}
	}
	children := map[ui.ItemID][]ui.CollectionItem{}
	for _, item := range state.Data.Items {
		children[item.Parent] = append(children[item.Parent], item)
	}
	var walk func(ui.ItemID, int)
	walk = func(parent ui.ItemID, depth int) {
		for _, item := range children[parent] {
			h := height
			if item.Kind == ui.Separator {
				h = height / 3
			}
			rows = append(rows, collectionRow{Item: item, Text: item.Label, Depth: depth, Y: y, H: h})
			y += h
			if item.HasChildren && state.Expanded[item.ID] {
				walk(item.ID, depth+1)
				status(item.ID, depth+1)
			}
		}
	}
	walk("", 0)
	status("", 0)
	if len(rows) == 0 {
		text := "Empty"
		if !state.RootLoaded {
			text = "Load — R"
		}
		rows = append(rows, collectionRow{Text: text, H: height, Status: true, Recovery: !state.RootLoaded})
	}
	return rows
}
func (m *collectionMeasure) Measure(n *parser.Instance, font, width float64) (layout.Size, error) {
	if n.Widget == "button" && (parser.IsCommandButton(n) || n.Argument("icon") != "" || n.Argument("tooltip") != "") {
		return m.measureButton(n, font), nil
	}
	if n.Widget == "tree" || n.Widget == "list" {
		x, e := m.MeasureCollection(n, font, layout.Size{W: width, H: rowHeight(font) * 4})
		if e != nil {
			return layout.Size{}, e
		}
		w, h := x.Content.W, x.Content.H+rowHeight(font)
		if n.Layout["overflow-x"] == "scroll" {
			w = x.Minimum.W
		}
		if n.Layout["overflow-y"] == "scroll" {
			h = x.Minimum.H
		}
		return layout.Size{W: math.Max(w, x.Minimum.W), H: math.Max(h, x.Minimum.H)}, nil
	}
	return m.markdown.Measure(n, font, width)
}
func (m *collectionMeasure) MeasureCollection(n *parser.Instance, font float64, outer layout.Size) (layout.CollectionMetrics, error) {
	state, ok := m.snapshot.Collections[n.Path]
	if !ok {
		return layout.CollectionMetrics{}, fmt.Errorf("collection-provider: missing %s", n.Path)
	}
	rows := collectionRows(state, font)
	height := rowHeight(font)
	gutterX, gutterY := 0.0, 0.0
	if n.Layout["overflow-y"] == "scroll" {
		gutterX = collectionGutter
	}
	if n.Layout["overflow-x"] == "scroll" {
		gutterY = collectionGutter
	}
	w, h := 0.0, 0.0
	for _, r := range rows {
		text, indent := fittedRow(r, state, font, math.Max(0, outer.W-gutterX))
		label := float64(nativeText(text, font).Width) + indent + 8
		w = math.Max(w, label)
		h = r.Y + r.H
	}
	return layout.CollectionMetrics{Minimum: collectionMinimum(font, n.Layout["overflow-x"] == "scroll", n.Layout["overflow-y"] == "scroll"), Content: layout.Size{W: w, H: h}, Viewport: layout.Rect{Y: height, W: math.Max(0, outer.W-gutterX), H: math.Max(0, outer.H-height-gutterY)}}, nil
}

func nativeText(text string, font float64) fyne.Size {
	size, _ := fyne.CurrentApp().Driver().RenderedTextSize(text, float32(font), fyne.TextStyle{}, regular)
	return size
}
func rowText(r collectionRow, state ui.CollectionState) string {
	if r.Status {
		return r.Text
	}
	branch, selected, focused := "  ", "  ", "  "
	if r.Item.HasChildren {
		branch = "+ "
		if state.Expanded[r.Item.ID] {
			branch = "- "
		}
	}
	if r.Item.ID != "" && state.Selected == r.Item.ID {
		selected = "* "
	}
	if r.Item.ID != "" && state.Focused == r.Item.ID {
		focused = "> "
	}
	return branch + selected + focused + r.Text
}

func collectionMinimum(font float64, sx, sy bool) layout.Size {
	height := rowHeight(font)
	w, h := height*3, height*2
	if sx {
		h += collectionGutter
	}
	if sy {
		w += collectionGutter
	}
	return layout.Size{W: w, H: h}
}

// Diagnostic affordances must remain presentable inside an already admitted
// viewport, including under overflow-x=error. Data rows keep their full width.
// The runtime retains the complete diagnostic; only its native projection fits.
func fittedRow(r collectionRow, state ui.CollectionState, font, width float64) (string, float64) {
	text := rowText(r, state)
	indent := float64(r.Depth) * 18
	if !r.Status {
		return text, indent
	}
	prefix := ""
	if r.Recovery {
		prefix = "Load/R"
		text = ""
		if status := state.Status[r.Parent]; status.Phase == ui.LoadError {
			prefix = "Retry/R"
			text = status.Error
		}
	}
	text = strings.Join(strings.Fields(text), " ")
	available := math.Max(0, width-8)
	minimum := float64(nativeText(prefix, font).Width)
	indent = math.Min(indent, math.Max(0, available-minimum))
	available = math.Max(0, available-indent)
	join := func(detail string) string {
		if prefix == "" {
			return detail
		}
		if detail == "" {
			return prefix
		}
		return prefix + ": " + detail
	}
	if float64(nativeText(join(text), font).Width) <= available {
		return join(text), indent
	}
	suffix := "..."
	if float64(nativeText(join(suffix), font).Width) > available {
		return prefix, indent
	}
	runes := []rune(text)
	lo, hi := 0, len(runes)
	for lo < hi {
		mid := (lo + hi + 1) / 2
		if float64(nativeText(join(string(runes[:mid])+suffix), font).Width) <= available {
			lo = mid
		} else {
			hi = mid - 1
		}
	}
	return join(string(runes[:lo]) + suffix), indent
}
