package fynehost

import (
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
	widgets := b.Session.Widgets()
	return map[string]any{"source": b.SourceRevision, "snapshot": b.Session.Snapshot(), "viewports": b.Geometry().Viewports, "rows": rows, "widgets": widgets, "focused": b.Session.Focused()}
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
