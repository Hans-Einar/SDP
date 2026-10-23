package fynehost

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"image/color"
)

// An empty dialog itself is its one native focus stop. It is only mounted while
// runtime has no eligible descendant, and creates no synthetic model widget.
type surfaceFocus struct {
	widget.BaseWidget
	bundle  *Bundle
	surface *nativeSurface
	shift   bool
}

func newSurfaceFocus(b *Bundle, s *nativeSurface) *surfaceFocus {
	f := &surfaceFocus{bundle: b, surface: s}
	f.ExtendBaseWidget(f)
	return f
}
func (f *surfaceFocus) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(canvas.NewRectangle(color.Transparent))
}
func (f *surfaceFocus) FocusGained() {
	if f.bundle.livePane() {
		f.bundle.owner.Mutate(func(s *ui.Session) error { return s.FocusSurface(&f.surface.target) })
	}
}
func (f *surfaceFocus) FocusLost()     { f.shift = false }
func (f *surfaceFocus) TypedRune(rune) {}
func (f *surfaceFocus) TypedKey(e *fyne.KeyEvent) {
	if commandKeyEvent(e, f.shift, false, f.TypedShortcut) {
		return
	}
	if e.Name == fyne.KeyEscape {
		f.bundle.escapeSurface(f.surface.path + "/")
	}
}
func (f *surfaceFocus) TypedShortcut(s fyne.Shortcut) { f.bundle.routeShortcut(f.surface.canvas, s) }

func (f *surfaceFocus) KeyDown(e *fyne.KeyEvent) {
	if shiftKey(e) {
		f.shift = true
	}
}
func (f *surfaceFocus) KeyUp(e *fyne.KeyEvent) {
	if shiftKey(e) {
		f.shift = false
	}
}
