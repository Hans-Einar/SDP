package fynehost

import (
	"fyne.io/fyne/v2"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

// Inspect returns owned diagnostic state and logical screen coordinates for the
// native acceptance harness. It does not invoke providers or user actions.
func (b *Bundle) Inspect() map[string]any {
	rows := map[string][]map[string]any{}
	for path, obj := range b.view.Controls {
		c, ok := obj.(*CollectionControl)
		if !ok {
			continue
		}
		for _, r := range c.rows {
			rect := layout.Rect{X: c.viewport.Rect.X, Y: c.viewport.Rect.Y + r.Y - c.viewport.Offset.Y, W: c.viewport.Rect.W, H: r.H}
			rows[path] = append(rows[path], map[string]any{"item": r.Item.ID, "parent": r.Parent, "kind": r.Item.Kind, "rect": rect, "clip": rect.Intersect(c.viewport.Clip), "recovery": r.Recovery, "text": r.Text})
		}
	}
	tabs := map[string]map[string]any{}
	for path, g := range b.Geometry().Tabs {
		pages := []map[string]any{}
		nativeSelected := ""
		if header, ok := b.view.Controls[path].(*paneHeader); ok {
			if index := header.tabs.SelectedIndex(); index >= 0 && index < len(header.pages) {
				nativeSelected = header.pages[index].id
			}
			for i, obj := range header.tabs.buttons() {
				if i >= len(header.pages) {
					break
				}
				if header.pages[i].id == "" {
					continue
				}
				pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(obj)
				size := obj.Size()
				rect := layout.Rect{X: float64(pos.X), Y: float64(pos.Y), W: float64(size.Width), H: float64(size.Height)}
				pages = append(pages, map[string]any{"id": header.pages[i].id, "label": header.pages[i].label, "enabled": header.pages[i].enabled, "rect": rect, "clip": rect.Intersect(g.HeaderClip)})
			}
		}
		tabs[path] = map[string]any{"header": g.Header, "clip": g.HeaderClip, "body": g.Body, "selected": g.Selected, "nativeSelected": nativeSelected, "pages": pages}
	}
	controls := map[string]map[string]any{}
	for path, c := range b.view.controls {
		pos, size := c.clip.Position(), c.clip.Size()
		offset := c.clip.Offset
		controls[path] = map[string]any{
			"rect":    layout.Rect{X: float64(pos.X - offset.X), Y: float64(pos.Y - offset.Y), W: float64(c.fixed.size.Width), H: float64(c.fixed.size.Height)},
			"clip":    layout.Rect{X: float64(pos.X), Y: float64(pos.Y), W: float64(size.Width), H: float64(size.Height)},
			"visible": c.clip.Visible() && c.widget.Visible() && size.Width > 0 && size.Height > 0,
		}
	}
	widgets := b.Session.Widgets()
	return map[string]any{"controls": controls, "tabs": tabs, "splits": b.Geometry().Splits, "source": b.SourceRevision, "snapshot": b.Session.Snapshot(), "viewports": b.Geometry().Viewports, "rows": rows, "widgets": widgets, "focused": b.Session.Focused()}
}

// CancelCollection is an application operation used by the interactive fixture.
func (h *DocumentHost) CancelCollection(path string) error {
	return h.Mutate(func(s *ui.Session) error {
		for _, w := range s.Widgets() {
			if w.InstancePath == path {
				return s.CancelLoad(w.Handle)
			}
		}
		return nil
	})
}
