package fynehost

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"image/color"
)

// paneDivider projects shared layout's usable length/bounds. It never asks a
// Fyne Split to reclamp minima or mutate its own independent proportion.
type paneDivider struct {
	widget.BaseWidget
	shift                            bool
	axis                             string
	proportion, lower, upper, usable float64
	collapsed                        bool
	disabled, focused                bool
	request                          func(string, float64)
	focus                            func()
	focusCanvas                      func()
	shortcut                         func(fyne.Shortcut)
	escape                           func()
}

func newPaneDivider() *paneDivider { d := &paneDivider{}; d.ExtendBaseWidget(d); return d }
func (d *paneDivider) CreateRenderer() fyne.WidgetRenderer {
	return &dividerRenderer{d: d, rect: canvas.NewRectangle(color.NRGBA{90, 110, 135, 255})}
}
func (d *paneDivider) Disabled() bool { return d.disabled }
func (d *paneDivider) Disable()       { d.disabled = true; d.Refresh() }
func (d *paneDivider) Enable()        { d.disabled = false; d.Refresh() }
func (d *paneDivider) FocusGained() {
	d.focused = true
	if !d.disabled && d.focus != nil {
		d.focus()
	}
	d.Refresh()
}
func (d *paneDivider) FocusLost()     { d.shift = false; d.focused = false; d.Refresh() }
func (d *paneDivider) TypedRune(rune) {}
func (d *paneDivider) Tapped(*fyne.PointEvent) {
	if !d.disabled && d.focusCanvas != nil {
		d.focusCanvas()
	}
}
func (d *paneDivider) TypedKey(e *fyne.KeyEvent) {
	if d.disabled || d.request == nil {
		return
	}
	if commandKeyEvent(e, d.shift, false, d.shortcut) {
		return
	}
	if e.Name == fyne.KeyEscape && d.escape != nil {
		d.escape()
		return
	}
	switch e.Name {
	case fyne.KeySpace:
		d.request("restore", 0)
	case fyne.KeyHome:
		d.request("ratio", d.lower)
	case fyne.KeyEnd:
		d.request("ratio", d.upper)
	case fyne.KeyLeft:
		if d.axis == "horizontal" {
			d.request("ratio", d.proportion-.05)
		}
	case fyne.KeyRight:
		if d.axis == "horizontal" {
			d.request("ratio", d.proportion+.05)
		}
	case fyne.KeyUp:
		if d.axis == "vertical" {
			d.request("ratio", d.proportion-.05)
		}
	case fyne.KeyDown:
		if d.axis == "vertical" {
			d.request("ratio", d.proportion+.05)
		}
	}
}
func (d *paneDivider) TypedShortcut(s fyne.Shortcut) {
	if d.disabled {
		return
	}
	if key, ok := s.(*desktop.CustomShortcut); ok && key.Modifier == fyne.KeyModifierControl && d.request != nil {
		if key.KeyName == fyne.KeyHome {
			d.request("collapse-first", 0)
			return
		}
		if key.KeyName == fyne.KeyEnd {
			d.request("collapse-second", 0)
			return
		}
	}
	if d.shortcut != nil {
		d.shortcut(s)
	}
}
func (d *paneDivider) Dragged(e *fyne.DragEvent) {
	if d.disabled || d.collapsed || d.usable <= 0 || d.request == nil {
		return
	}
	if d.focusCanvas != nil {
		d.focusCanvas()
	}
	delta := float64(e.Dragged.DX)
	if d.axis == "vertical" {
		delta = float64(e.Dragged.DY)
	}
	d.request("ratio", d.proportion+delta/d.usable)
}
func (d *paneDivider) DragEnd() {}
func (d *paneDivider) Cursor() desktop.Cursor {
	if d.axis == "vertical" {
		return desktop.VResizeCursor
	}
	return desktop.HResizeCursor
}

type dividerRenderer struct {
	d    *paneDivider
	rect *canvas.Rectangle
}

func (r *dividerRenderer) Layout(s fyne.Size)           { r.rect.Resize(s) }
func (r *dividerRenderer) MinSize() fyne.Size           { return fyne.NewSize(1, 1) }
func (r *dividerRenderer) Objects() []fyne.CanvasObject { return []fyne.CanvasObject{r.rect} }
func (r *dividerRenderer) Destroy()                     {}
func (r *dividerRenderer) Refresh() {
	r.rect.FillColor = color.NRGBA{90, 110, 135, 255}
	if r.d.focused {
		r.rect.FillColor = color.NRGBA{30, 100, 220, 255}
	}
	r.rect.Refresh()
}

func (d *paneDivider) KeyDown(e *fyne.KeyEvent) {
	if shiftKey(e) {
		d.shift = true
	}
}
func (d *paneDivider) KeyUp(e *fyne.KeyEvent) {
	if shiftKey(e) {
		d.shift = false
	}
}
