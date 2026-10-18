package fynehost

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image/color"
)

// nativeMenu keeps Fyne's private menu items intact. Only synchronous selection
// forwarding can authorize an Action after Fyne dismisses its popup first.
// Every instance belongs to one root opening and is never reused.
type nativeMenu struct {
	model                              *fyne.Menu
	popup                              *widget.PopUpMenu
	canvas                             fyne.Canvas
	overlay                            fyne.CanvasObject
	input                              *menuInput
	dismiss                            func()
	nativeTransition                   func(func())
	scoped, claimed, dismissed, closed bool
}

func newNativeMenu(model *fyne.Menu, c fyne.Canvas, dismiss func()) *nativeMenu {
	m := &nativeMenu{canvas: c, dismiss: dismiss}
	model = copyNativeMenu(model)
	m.model = model
	var bind func(*fyne.Menu)
	bind = func(menu *fyne.Menu) {
		for _, item := range menu.Items {
			if item.ChildMenu != nil {
				bind(item.ChildMenu)
			}
			if action := item.Action; action != nil {
				item.Action = func() {
					if item.Disabled || item.ChildMenu != nil || m.closed || !m.scoped || !m.dismissed || m.claimed {
						return
					}
					m.claimed = true // claim before domain execution, including errors/reentry
					action()
				}
			}
		}
	}
	bind(model)
	m.popup = widget.NewPopUpMenu(model, c)
	m.popup.OnDismiss = m.onDismiss
	m.input = &menuInput{menu: m}
	m.input.ExtendBaseWidget(m.input)
	return m
}

func copyNativeMenu(menu *fyne.Menu) *fyne.Menu {
	copy := &fyne.Menu{Label: menu.Label, Items: make([]*fyne.MenuItem, len(menu.Items))}
	for i, item := range menu.Items {
		value := *item
		if item.ChildMenu != nil {
			value.ChildMenu = copyNativeMenu(item.ChildMenu)
		}
		copy.Items[i] = &value
	}
	return copy
}
func (m *nativeMenu) show(pos fyne.Position) {
	m.popup.ShowAtPosition(pos)
	// Keep the original native overlay and its full menu/scroll subtree in the
	// same hit-test tree. A separate top overlay would mask hover and scrolling.
	native := m.canvas.Overlays().Top()
	m.canvas.Overlays().Remove(native)
	m.overlay = container.NewStack(native, m.input)
	m.canvas.Overlays().Add(m.overlay)
	m.canvas.Focus(m.input)
}
func (m *nativeMenu) hide() {
	if m.nativeTransition != nil {
		m.nativeTransition(m.hideNative)
		return
	}
	m.hideNative()
}
func (m *nativeMenu) hideNative() {
	// Remove only this opening. Remove on a non-top overlay would also remove
	// a replacement above it; never allow obsolete callbacks to do that.
	if m.overlay != nil && m.canvas.Overlays().Top() == m.overlay {
		m.canvas.Overlays().Remove(m.overlay)
	}
	m.popup.Hide()
}
func (m *nativeMenu) onDismiss() {
	if m.closed {
		return
	}
	m.dismissed = true
	m.hide()
	if !m.scoped {
		m.finish()
	}
}
func (m *nativeMenu) finish() {
	if m.closed {
		return
	}
	m.closed = true
	if m.dismiss != nil {
		m.dismiss()
	}
}
func (m *nativeMenu) selection(fn func()) {
	if m.closed || m.dismissed || m.scoped {
		return
	}
	m.scoped = true
	defer func() {
		m.scoped = false
		if m.dismissed {
			m.finish()
		}
	}()
	fn()
}
func (m *nativeMenu) close() {
	if m.closed {
		return
	}
	m.dismissed = true
	m.hide()
	m.finish()
}

// nativeMenuItemChildren accesses an exported method through a public interface;
// no private concrete type assertion, field access, renderer surgery or unsafe.
type nativeMenuItemChildren interface{ Child() *widget.Menu }

func walkNativeMenu(menu *widget.Menu, visit func(*widget.Menu, fyne.CanvasObject)) {
	if !menu.Visible() {
		return
	}
	for _, obj := range menu.Items {
		visit(menu, obj)
		if nested, ok := obj.(nativeMenuItemChildren); ok {
			if child := nested.Child(); child != nil {
				walkNativeMenu(child, visit)
			}
		}
	}
}
func (m *nativeMenu) hit(point fyne.Position) fyne.CanvasObject {
	var hit fyne.CanvasObject
	driver := fyne.CurrentApp().Driver()
	walkNativeMenu(m.popup.Menu, func(menu *widget.Menu, obj fyne.CanvasObject) {
		mp, ms := driver.AbsolutePositionForObject(menu), menu.Size()
		// Root PopUpMenu is the mounted widget, its embedded Menu is not a tree node.
		if menu == m.popup.Menu {
			mp = driver.AbsolutePositionForObject(m.popup)
		}
		p, s := driver.AbsolutePositionForObject(obj), obj.Size()
		inside := func(pos fyne.Position, size fyne.Size) bool {
			return point.X >= pos.X && point.Y >= pos.Y && point.X < pos.X+size.Width && point.Y < pos.Y+size.Height
		}
		if obj.Visible() && inside(mp, ms) && inside(p, s) {
			hit = obj
		}
	})
	return hit
}

type menuInput struct {
	widget.BaseWidget
	menu *nativeMenu
}

func (p *menuInput) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(canvas.NewRectangle(color.Transparent))
}
func (p *menuInput) FocusGained()   {}
func (p *menuInput) FocusLost()     {}
func (p *menuInput) TypedRune(rune) {}
func (p *menuInput) TypedKey(e *fyne.KeyEvent) {
	m := p.menu
	if m.closed || m.dismissed {
		return
	}
	switch e.Name {
	case fyne.KeyEnter, fyne.KeyReturn, fyne.KeySpace:
		m.selection(func() { m.popup.TypedKey(e) })
	case fyne.KeyHome:
		for i := 0; i < m.visibleItems(); i++ {
			m.popup.ActivatePrevious()
		}
	case fyne.KeyEnd:
		for i := 0; i < m.visibleItems(); i++ {
			m.popup.ActivateNext()
		}
	default:
		m.popup.TypedKey(e) // Escape dismisses outside selection scope.
	}
}

func (m *nativeMenu) visibleItems() int {
	count := 0
	walkNativeMenu(m.popup.Menu, func(_ *widget.Menu, _ fyne.CanvasObject) { count++ })
	return count
}
func (p *menuInput) Tapped(e *fyne.PointEvent) {
	m := p.menu
	if m.closed || m.dismissed {
		return
	}
	point := fyne.CurrentApp().Driver().AbsolutePositionForObject(p).Add(e.Position)
	target := m.hit(point)
	if target == nil {
		m.popup.Dismiss() // outside, padding and separators never authorize an Action
		return
	}
	if tappable, ok := target.(fyne.Tappable); ok {
		pos := fyne.CurrentApp().Driver().AbsolutePositionForObject(target)
		m.selection(func() { tappable.Tapped(&fyne.PointEvent{Position: point.Subtract(pos), AbsolutePosition: point}) })
	}
}
