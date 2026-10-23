package fynehost

import (
	"fmt"
	"fyne.io/fyne/v2"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

func (h *DocumentHost) collectionEvent(c *CollectionControl, id ui.ItemID, kind ui.EventKind) {
	if !c.live() {
		return
	}
	h.status(h.Mutate(func(s *ui.Session) error {
		t, e := s.Target(c.state.Handle, id)
		if e != nil {
			return e
		}
		return s.Dispatch(ui.Event{Handle: t.Handle, ModelRevision: t.ModelRevision, Sequence: s.Sequence() + 1, Kind: kind, Collection: &t})
	}))
}
func (h *DocumentHost) focusItem(c *CollectionControl, id ui.ItemID, selectRow bool) {
	if !c.live() {
		return
	}
	if e := h.Mutate(func(s *ui.Session) error {
		t, e := s.Target(c.state.Handle, id)
		if e != nil {
			return e
		}
		if e = s.Focus(c.state.Handle); e != nil {
			return e
		}
		return s.FocusItem(t)
	}); e != nil {
		return
	}
	if selectRow {
		h.collectionEvent(c, id, ui.Select)
	}
	h.ensureItem(c, id)
}
func (h *DocumentHost) focusCollection(c *CollectionControl) {
	if !c.live() {
		return
	}
	if err := c.bundle.Session.Focus(c.state.Handle); err != nil {
		h.status(err)
		return
	}
	id := c.state.Focused
	found := false
	for _, r := range c.rows {
		if r.Item.ID == id && id != "" {
			found = true
		}
	}
	if !found {
		for _, r := range c.rows {
			if !r.Status && r.Item.Kind != ui.Separator {
				id = r.Item.ID
				break
			}
		}
	}
	if id != "" {
		h.focusItem(c, id, false)
	} else {
		h.after()
	}
}
func (h *DocumentHost) toggle(c *CollectionControl, id ui.ItemID) {
	kind := ui.Expand
	if c.state.Expanded[id] {
		kind = ui.Collapse
	}
	h.collectionEvent(c, id, kind)
}
func (h *DocumentHost) recoverCollection(c *CollectionControl) {
	id := c.state.Focused
	st := c.state.Status[id]
	if st.Phase == ui.LoadError || st.Phase == ui.Canceled {
		h.collectionEvent(c, id, ui.Retry)
		return
	}
	navigable := false
	for _, r := range c.rows {
		if !r.Status && r.Item.Kind != ui.Separator {
			navigable = true
		}
	}
	if !navigable {
		st = c.state.Status[""]
		if st.Phase == ui.LoadError || st.Phase == ui.Canceled || st.Phase == ui.Unloaded && !c.state.RootLoaded && !c.state.AutoLoadPending {
			h.collectionEvent(c, "", ui.Retry)
		}
	}
}
func (h *DocumentHost) collectionKey(c *CollectionControl, key fyne.KeyName) {
	rows := []collectionRow{}
	idx := -1
	for _, r := range c.rows {
		if !r.Status && r.Item.Kind != ui.Separator {
			if r.Item.ID == c.state.Focused {
				idx = len(rows)
			}
			rows = append(rows, r)
		}
	}
	move := func(i int) {
		if i < 0 || i >= len(rows) {
			return
		}
		r := rows[i]
		h.focusItem(c, r.Item.ID, r.Item.Kind == ui.Row)
	}
	switch key {
	case fyne.KeyEscape:
		if c.state.Request == nil && c.bundle.escapeSurface(c.path) {
			return
		}
		h.status(h.Mutate(func(s *ui.Session) error { return s.CancelLoad(c.state.Handle) }))
	case fyne.KeyDown:
		if idx < len(rows)-1 {
			move(idx + 1)
		}
	case fyne.KeyUp:
		if idx > 0 {
			move(idx - 1)
		}
	case fyne.KeyHome:
		move(0)
	case fyne.KeyEnd:
		move(len(rows) - 1)
	case fyne.KeyPageDown:
		h.scrollIn(c.path, c.viewport.Clip.X+c.viewport.Clip.W/2, c.viewport.Clip.Y+c.viewport.Clip.H/2, 0, c.viewport.Clip.H*.9)
	case fyne.KeyPageUp:
		h.scrollIn(c.path, c.viewport.Clip.X+c.viewport.Clip.W/2, c.viewport.Clip.Y+c.viewport.Clip.H/2, 0, -c.viewport.Clip.H*.9)
	case fyne.KeySpace:
		if idx >= 0 && rows[idx].Item.Kind == ui.Row {
			h.collectionEvent(c, rows[idx].Item.ID, ui.Select)
		}
	case fyne.KeyReturn, fyne.KeyEnter:
		if idx < 0 {
			h.recoverCollection(c)
			return
		}
		r := rows[idx]
		if r.Item.Kind == ui.Row {
			h.collectionEvent(c, r.Item.ID, ui.Activate)
		} else if c.state.Status[r.Item.ID].Phase == ui.LoadError {
			h.collectionEvent(c, r.Item.ID, ui.Retry)
		} else if r.Item.HasChildren {
			h.toggle(c, r.Item.ID)
		}
	case fyne.KeyLeft:
		if c.state.Handle.Kind != "tree" || idx < 0 {
			return
		}
		r := rows[idx]
		if r.Item.HasChildren && c.state.Expanded[r.Item.ID] {
			h.collectionEvent(c, r.Item.ID, ui.Collapse)
		} else {
			for i, row := range rows {
				if row.Item.ID == r.Item.Parent {
					move(i)
					break
				}
			}
		}
	case fyne.KeyRight:
		if c.state.Handle.Kind != "tree" || idx < 0 {
			return
		}
		r := rows[idx]
		if r.Item.HasChildren && !c.state.Expanded[r.Item.ID] {
			h.collectionEvent(c, r.Item.ID, ui.Expand)
		} else if idx+1 < len(rows) && rows[idx+1].Item.Parent == r.Item.ID {
			move(idx + 1)
		}
	}
}
func (h *DocumentHost) setViewport(path string, offset ui.ViewportState) {
	b := h.current
	if b == nil {
		return
	}
	handle, ok := b.Session.Viewport(path)
	if !ok {
		return
	}
	revision := b.Session.Revision
	h.status(h.Mutate(func(s *ui.Session) error { return s.SetViewport(handle, revision, offset) }))
}
func (h *DocumentHost) scroll(x, y, dx, dy float64) { h.scrollIn("", x, y, dx, dy) }
func (h *DocumentHost) scrollIn(path string, x, y, dx, dy float64) {
	b := h.current
	if b == nil || b.Geometry() == nil {
		return
	}
	offsets, _, e := b.geometryFor(path).RouteScroll(x, y, dx, dy)
	if e != nil {
		h.status(e)
		return
	}
	h.status(h.Mutate(func(s *ui.Session) error {
		return h.commitOffsets(b, offsets)
	}))
}
func (h *DocumentHost) ensureItem(c *CollectionControl, id ui.ItemID) {
	for _, r := range c.rows {
		if r.Item.ID == id {
			h.ensure(c.path, layout.Rect{X: c.viewport.Rect.X + float64(r.Depth)*18 - c.viewport.Offset.X, Y: c.viewport.Rect.Y + r.Y - c.viewport.Offset.Y, W: float64(nativeText(rowText(r, c.state), c.font).Width) + 8, H: r.H})
			return
		}
	}
}
func (h *DocumentHost) ensureWidget(path string) {
	b := h.current
	if b == nil || b.Geometry() == nil {
		return
	}
	b.geometryFor(path).Root.Walk(func(box *layout.Box) {
		if box.Path == path {
			h.ensure(path, box.Rect)
		}
	})
}
func (h *DocumentHost) ensure(path string, rect layout.Rect) {
	b := h.current
	if b == nil {
		return
	}
	offsets, e := b.geometryFor(path).EnsureVisible(path, rect)
	if e != nil {
		h.status(e)
		return
	}
	h.status(h.Mutate(func(s *ui.Session) error { return h.commitOffsets(b, offsets) }))
}

// Validate every captured viewport before publishing the offset map as one state.
func (h *DocumentHost) commitOffsets(b *Bundle, offsets map[string]ui.ViewportState) error {
	if h.current != b || b.closed || b.presentation == nil {
		return fmt.Errorf("stale viewport bundle")
	}
	captured := b.presentation.snapshot
	if b.Session.Revision != captured.ModelRevision {
		return fmt.Errorf("stale viewport model")
	}
	for path := range offsets {
		handle, ok := b.Session.Viewport(path)
		if !ok || handle != captured.ViewportHandles[path] {
			return fmt.Errorf("stale viewport: %s", path)
		}
	}
	return b.Session.SetViewports(offsets)
}
