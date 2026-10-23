package fynehost

import (
	"fmt"
	"runtime"
	"strings"
	"unicode/utf8"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type textEntry struct {
	widget.Entry
	owner           *textControl
	shift, readOnly bool
	renderer        fyne.WidgetRenderer
}

func (e *textEntry) CreateRenderer() fyne.WidgetRenderer {
	e.renderer = e.Entry.CreateRenderer()
	return e.renderer
}
func (e *textEntry) nativeScroll() *container.Scroll {
	if e.renderer != nil {
		for _, o := range e.renderer.Objects() {
			if s, ok := o.(*container.Scroll); ok {
				return s
			}
		}
	}
	return nil
}
func (e *textEntry) AcceptsTab() bool { return false }
func (e *textEntry) FocusGained() {
	if e.owner.bundle != nil && e.owner.bundle.muted {
		defer e.retainScroll()()
	}
	e.Entry.FocusGained()
	e.owner.focus()
}
func (e *textEntry) FocusLost() {
	defer e.retainScroll()()
	e.shift = false
	e.Entry.KeyUp(&fyne.KeyEvent{Name: desktop.KeyShiftLeft})
	e.Entry.FocusLost()
}
func (e *textEntry) retainScroll() func() {
	if scroll := e.nativeScroll(); scroll != nil {
		offset := scroll.Offset
		return func() { scroll.ScrollToOffset(offset) }
	}
	return func() {}
}
func (e *textEntry) KeyDown(k *fyne.KeyEvent) {
	e.Entry.KeyDown(k)
	if shiftKey(k) {
		e.shift = true
	}
}
func (e *textEntry) KeyUp(k *fyne.KeyEvent) {
	e.Entry.KeyUp(k)
	if shiftKey(k) {
		e.shift = false
	}
}
func (e *textEntry) TypedRune(r rune) {
	if !e.owner.editable() {
		return
	}
	if !e.MultiLine && (r == '\r' || r == '\n') {
		e.refuse("text-single-line")
		return
	}
	e.Entry.TypedRune(r)
}
func (e *textEntry) TypedKey(k *fyne.KeyEvent) {
	if e.Disabled() || !e.owner.live() {
		return
	}
	if commandKeyEvent(k, e.shift, true, e.TypedShortcut) {
		return
	}
	switch k.Name {
	case fyne.KeyEscape:
		e.owner.escape()
		return
	case fyne.KeyTab:
		return // native canvas/source traversal owns this key
	case fyne.KeyReturn, fyne.KeyEnter:
		if !e.owner.editable() {
			return
		}
		if !e.MultiLine {
			e.owner.commit()
			return
		}
	case fyne.KeyBackspace, fyne.KeyDelete:
		if !e.owner.editable() {
			return
		}
	}
	// OnSubmitted stays nil: native multiline Shift+Enter is also a newline.
	e.Entry.TypedKey(k)
}
func (e *textEntry) TypedShortcut(s fyne.Shortcut) {
	if e.Disabled() || !e.owner.live() {
		return
	}
	if key, ok := s.(*desktop.CustomShortcut); ok {
		// Pinned GLFW emits Primary+Shift+Z as CustomShortcut (Primary+Y
		// becomes named Redo). Entry registers only the named shortcut.
		if key.KeyName == fyne.KeyZ && key.Modifier == fyne.KeyModifierShortcutDefault|fyne.KeyModifierShift {
			if e.owner.editable() {
				e.Entry.TypedShortcut(&fyne.ShortcutRedo{})
			}
			return
		}
		if key.KeyName == fyne.KeyF10 && key.Modifier == fyne.KeyModifierShift && !e.owner.bundle.hasContext(e.owner.path) {
			e.owner.openEditMenu(fyne.CurrentApp().Driver().AbsolutePositionForObject(e))
			return
		}
		if (key.KeyName == fyne.KeyReturn || key.KeyName == fyne.KeyEnter) && key.Modifier == fyne.KeyModifierShortcutDefault {
			e.owner.commit()
			return
		}
		if native, editing := textNavigationShortcut(key); native {
			if editing && !e.owner.editable() {
				return
			}
			e.Entry.TypedShortcut(s)
			return
		}
	}
	switch sh := s.(type) {
	case *fyne.ShortcutCopy, *fyne.ShortcutSelectAll:
		e.Entry.TypedShortcut(s)
		return
	case *fyne.ShortcutPaste:
		if !e.owner.editable() {
			return
		}
		// Read the OS clipboard exactly once. Fyne must receive these same bytes,
		// not a second read that might evade newline/encoding checks.
		raw := sh.Clipboard.Content()
		if !utf8.ValidString(raw) || len(raw) > 32768 {
			e.refuse("text-limit")
			return
		}
		if !e.MultiLine && strings.ContainsAny(raw, "\r\n") {
			e.refuse("text-single-line")
			return
		}
		e.Entry.TypedShortcut(&fyne.ShortcutPaste{Clipboard: &capturedTextClipboard{text: raw}, Secondary: sh.Secondary})
		return
	case *fyne.ShortcutCut, *fyne.ShortcutUndo, *fyne.ShortcutRedo:
		if !e.owner.editable() {
			return
		}
		e.Entry.TypedShortcut(s)
		return
	}
	e.owner.shortcut(s)
}

// Match only native Entry's registered public shortcut forms. Extra modifiers
// remain available to declared commands instead of disappearing into Entry.
func textNavigationShortcut(s *desktop.CustomShortcut) (native, editing bool) {
	word := fyne.KeyModifierShortcutDefault
	if runtime.GOOS == "darwin" {
		word = fyne.KeyModifierAlt
	}
	switch s.KeyName {
	case fyne.KeyLeft, fyne.KeyRight:
		return s.Modifier == word || s.Modifier == word|fyne.KeyModifierShift || runtime.GOOS == "darwin" && s.Modifier == fyne.KeyModifierSuper, false
	case fyne.KeyBackspace, fyne.KeyDelete:
		return s.Modifier == word, true
	}
	return false, false
}
func (e *textEntry) refuse(code string) {
	if e.owner.live() {
		e.owner.bundle.owner.status(fmt.Errorf("%s: %s", code, e.owner.path))
	}
}

type capturedTextClipboard struct{ text string }

func (c *capturedTextClipboard) Content() string        { return c.text }
func (c *capturedTextClipboard) SetContent(text string) { c.text = text }

func (e *textEntry) TappedSecondary(p *fyne.PointEvent) {
	c := e.owner
	if !c.live() || e.Disabled() {
		return
	}
	b := c.bundle
	b.canvasFor(c.path).Focus(e)
	if b.hasContext(c.path) {
		b.openContext(c.path, "", p.AbsolutePosition)
		return
	}
	c.openEditMenu(p.AbsolutePosition)
}
func (c *textControl) openEditMenu(pos fyne.Position) {
	if !c.live() || c.Disabled() {
		return
	}
	if c.editMenu != nil {
		c.editMenu.close()
	}
	b := c.bundle
	clipboard := fyne.CurrentApp().Clipboard()
	model := fyne.NewMenu("")
	add := func(label string, s fyne.Shortcut) {
		model.Items = append(model.Items, fyne.NewMenuItem(label, func() {
			if c.live() {
				c.entry.TypedShortcut(s)
			}
		}))
	}
	if c.editable() {
		add("Undo", &fyne.ShortcutUndo{})
		add("Redo", &fyne.ShortcutRedo{})
		add("Cut", &fyne.ShortcutCut{Clipboard: clipboard})
	}
	add("Copy", &fyne.ShortcutCopy{Clipboard: clipboard})
	if c.editable() {
		add("Paste", &fyne.ShortcutPaste{Clipboard: clipboard})
	}
	add("Select all", &fyne.ShortcutSelectAll{})
	var menu *nativeMenu
	menu = newNativeMenu(model, b.canvasFor(c.path), func() {
		if c.editMenu != menu {
			return
		}
		c.editMenu = nil
		if c.live() {
			b.owner.restoreFocus(b)
		}
	})
	menu.nativeTransition = func(fn func()) { muted := b.muted; b.muted = true; defer func() { b.muted = muted }(); fn() }
	c.editMenu = menu
	muted := b.muted
	b.muted = true
	menu.show(pos)
	b.muted = muted
}
