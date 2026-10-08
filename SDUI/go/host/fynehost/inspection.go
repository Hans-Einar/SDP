package fynehost

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

// Inspect returns owned diagnostic state on the Fyne owner goroutine. Rectangles
// use logical client coordinates of the named canvas/window, not OS coordinates.
// It does not invoke providers or user actions.
func (b *Bundle) Inspect() map[string]any {
	snapshot := b.Session.Snapshot()
	rows := map[string][]map[string]any{}
	for path, obj := range b.view.Controls {
		c, ok := obj.(*CollectionControl)
		if !ok {
			continue
		}
		origin := b.inspectionOrigin(path)
		canvas, title := b.inspectionCanvas(b.canvasFor(path))
		for _, r := range c.rows {
			rect := layout.Rect{X: c.viewport.Rect.X, Y: c.viewport.Rect.Y + r.Y - c.viewport.Offset.Y, W: c.viewport.Rect.W, H: r.H}
			rows[path] = append(rows[path], map[string]any{"item": r.Item.ID, "parent": r.Parent, "kind": r.Item.Kind, "rect": inspectionTranslate(rect, origin), "clip": inspectionTranslate(rect.Intersect(c.viewport.Clip), origin), "canvas": canvas, "title": title, "recovery": r.Recovery, "text": r.Text})
		}
	}
	tabs := map[string]map[string]any{}
	tabGeometry := map[string]layout.TabsLayout{}
	for path, g := range b.Geometry().Tabs {
		tabGeometry[path] = g
	}
	if b.presentation != nil {
		for _, frame := range b.presentation.canvases {
			for path, g := range frame.geometry.Tabs {
				tabGeometry[path] = g
			}
		}
	}
	for path, g := range tabGeometry {
		origin := b.inspectionOrigin(path)
		g.Header = inspectionTranslate(g.Header, origin)
		g.HeaderClip = inspectionTranslate(g.HeaderClip, origin)
		g.Body = inspectionTranslate(g.Body, origin)
		canvas, title := b.inspectionCanvas(b.canvasFor(path))
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
		tabs[path] = map[string]any{"header": g.Header, "clip": g.HeaderClip, "body": g.Body, "selected": g.Selected, "nativeSelected": nativeSelected, "pages": pages, "canvas": canvas, "title": title}
	}
	controls := map[string]map[string]any{}
	for path, c := range b.view.controls {
		pos, size := fyne.CurrentApp().Driver().AbsolutePositionForObject(c.clip), c.clip.Size()
		canvas, title := b.inspectionCanvas(b.canvasFor(path))
		surface := surfacePath(snapshot, path)
		live := true
		if surface != "" {
			native := b.surfaces[surface]
			live = native != nil && native.shown && !native.retiring
		}
		offset := c.clip.Offset
		clip := layout.Rect{X: float64(pos.X), Y: float64(pos.Y), W: float64(size.Width), H: float64(size.Height)}.Intersect(inspectionBounds(b.canvasFor(path)))
		controls[path] = map[string]any{
			"rect":   layout.Rect{X: float64(pos.X - offset.X), Y: float64(pos.Y - offset.Y), W: float64(c.fixed.size.Width), H: float64(c.fixed.size.Height)},
			"clip":   clip,
			"canvas": canvas, "title": title, "surface": surface,
			"visible": live && c.clip.Visible() && c.widget.Visible() && clip.W > 0 && clip.H > 0,
		}
	}
	fields, choices := b.inspectFields()
	widgets := b.Session.Widgets()
	return map[string]any{"fields": fields, "choices": choices, "controls": controls, "tabs": tabs, "splits": b.Geometry().Splits, "source": b.SourceRevision, "snapshot": snapshot, "surfaces": b.inspectSurfaces(), "menus": b.inspectMenus(), "viewports": b.Geometry().Viewports, "rows": rows, "widgets": widgets, "focused": b.Session.Focused()}
}

// inspectionCanvas identifies actual windows, including a modal's parent window.
func (b *Bundle) inspectionCanvas(c fyne.Canvas) (string, string) {
	id := "main"
	for path, surface := range b.surfaces {
		if surface.window != nil && surface.canvas == c {
			id = path
			break
		}
	}
	if w := windowForCanvas(c); w != nil {
		return id, w.Title()
	}
	return id, ""
}

func inspectionRect(obj fyne.CanvasObject) layout.Rect {
	p, s := fyne.CurrentApp().Driver().AbsolutePositionForObject(obj), obj.Size()
	return layout.Rect{X: float64(p.X), Y: float64(p.Y), W: float64(s.Width), H: float64(s.Height)}
}
func inspectionBounds(c fyne.Canvas) layout.Rect {
	if c == nil {
		return layout.Rect{}
	}
	s := c.Size()
	return layout.Rect{W: float64(s.Width), H: float64(s.Height)}
}
func inspectionTranslate(r layout.Rect, p fyne.Position) layout.Rect {
	r.X += float64(p.X)
	r.Y += float64(p.Y)
	return r
}
func (b *Bundle) inspectionOrigin(path string) fyne.Position {
	if b.presentation != nil {
		if owner := surfacePath(b.presentation.snapshot, path); owner != "" {
			if surface := b.surfaces[owner]; surface != nil {
				return fyne.CurrentApp().Driver().AbsolutePositionForObject(surface.content)
			}
		}
	}
	return fyne.CurrentApp().Driver().AbsolutePositionForObject(b.view.Container)
}

func (b *Bundle) inspectSurfaces() map[string]map[string]any {
	out := map[string]map[string]any{}
	for path, surface := range b.surfaces {
		canvas, title := b.inspectionCanvas(surface.canvas)
		r := inspectionRect(surface.content)
		clip := r.Intersect(inspectionBounds(surface.canvas))
		out[path] = map[string]any{
			"target": surface.target, "modal": surface.modal != nil,
			"canvas": canvas, "title": title, "rect": r, "clip": clip,
			"visible":    surface.shown && !surface.retiring && surface.content.Visible() && clip.W > 0 && clip.H > 0,
			"chromeRect": inspectionRect(surface.chrome), "size": surface.size,
		}
	}
	return out
}

func (b *Bundle) inspectMenus() map[string]map[string]any {
	out := map[string]map[string]any{}
	for path, menu := range b.menus {
		native := menu.native
		if native == nil || native.closed || native.dismissed || !native.popup.Visible() {
			continue
		}
		canvas, title := b.inspectionCanvas(native.canvas)
		bounds := inspectionBounds(native.canvas)
		items := []map[string]any{}
		var walk func(*widget.Menu, *fyne.Menu, string)
		walk = func(view *widget.Menu, model *fyne.Menu, prefix string) {
			if !view.Visible() {
				return
			}
			menuRect := inspectionRect(view)
			if view == native.popup.Menu {
				// Embedded Menu is not itself mounted; the PopUpMenu is.
				menuRect = inspectionRect(native.popup)
			}
			for i, obj := range view.Items {
				if i >= len(model.Items) || !obj.Visible() {
					continue
				}
				index := fmt.Sprintf("%s/%d", prefix, i)
				item := model.Items[i]
				r := inspectionRect(obj)
				clip := r.Intersect(menuRect).Intersect(bounds)
				items = append(items, map[string]any{
					"path": menu.rows[index], "index": index, "label": item.Label,
					"rect": r, "clip": clip, "submenu": item.ChildMenu != nil,
					"enabled": !item.Disabled && !item.IsSeparator, "checked": item.Checked,
					"separator": item.IsSeparator, "canvas": canvas, "title": title,
				})
				if child, ok := obj.(nativeMenuItemChildren); ok && item.ChildMenu != nil && child.Child() != nil {
					walk(child.Child(), item.ChildMenu, index)
				}
			}
		}
		walk(native.popup.Menu, native.model, "")
		r := inspectionRect(native.popup)
		out[path] = map[string]any{"title": title, "canvas": canvas, "rect": r, "clip": r.Intersect(bounds), "items": items, "scope": menu.scope}
	}
	return out
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
