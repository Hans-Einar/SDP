package fynehost

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

// Fyne Entry.Refresh reveals the caret, including on unrelated theme/container
// refreshes. Retain its actual offset across identical-text publication; the
// public native ScrollToOffset still clamps against any accepted resize/reflow.
// These observations exist only for this synchronous paint, never in runtime.
func (b *Bundle) retainTextScroll(snapshot ui.Snapshot) func() {
	type saved struct {
		entry  *textEntry
		scroll *container.Scroll
		text   string
		offset fyne.Position
	}
	var offsets []saved
	for path, obj := range b.view.Controls {
		c, ok := obj.(*textControl)
		if !ok {
			continue
		}
		f, ok := snapshot.Fields[path]
		if !ok || f.Proposed.Text != c.entry.Text {
			continue
		}
		if s := c.entry.nativeScroll(); s != nil {
			offsets = append(offsets, saved{c.entry, s, c.entry.Text, s.Offset})
		}
	}
	return func() {
		for _, o := range offsets {
			if o.entry.Text == o.text {
				o.scroll.ScrollToOffset(o.offset)
			}
		}
	}
}

func (b *Bundle) inspectTextFields(fields map[string]map[string]any) {
	for path, obj := range b.view.Controls {
		c, ok := obj.(*textControl)
		if !ok {
			continue
		}
		cv := b.canvasFor(path)
		canvas, title := b.inspectionCanvas(cv)
		holder := b.view.controls[path]
		clip := inspectionRect(holder.clip).Intersect(inspectionBounds(cv))
		entry := c.entry
		parts := map[string]any{
			"kind": "input", "canvas": canvas, "title": title,
			"label": inspectionRect(c.label), "labelText": c.label.Text,
			"control": inspectionRect(entry), "entry": inspectionRect(entry), "feedback": inspectionRect(c.feedback),
			"clip": clip, "visible": holder.clip.Visible() && clip.W > 0 && clip.H > 0,
			"text": entry.Text, "placeholder": entry.PlaceHolder, "readOnly": entry.readOnly,
			"focused": cv != nil && cv.Focused() == entry, "cursorRow": entry.CursorRow, "cursorColumn": entry.CursorColumn,
			"selectedText": entry.SelectedText(), "multiline": entry.MultiLine,
		}
		if scroll := entry.nativeScroll(); scroll != nil {
			parts["scrollRect"] = inspectionRect(scroll)
			parts["scrollOffset"] = layout.Size{W: float64(scroll.Offset.X), H: float64(scroll.Offset.Y)}
		}
		fields[path] = parts
	}
}
func (b *Bundle) syncTextMenus() {
	for _, obj := range b.view.Controls {
		c, ok := obj.(*textControl)
		if !ok || c.editMenu == nil || c.editMenu.scoped {
			continue
		}
		w, ok := b.Session.Widget(c.path)
		if !ok || w.Handle != c.state.Target.Handle || !w.Visible || !w.Enabled {
			c.editMenu.close()
		}
	}
}
func (b *Bundle) closeTextMenus() {
	if b.view == nil {
		return
	}
	for _, obj := range b.view.Controls {
		if c, ok := obj.(*textControl); ok && c.editMenu != nil {
			c.editMenu.close()
		}
	}
}
