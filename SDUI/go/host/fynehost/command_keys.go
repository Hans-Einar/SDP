package fynehost

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/desktop"
)

// GLFW sends unmodified/function and Shift-only keys through TypedKey rather
// than canvas shortcuts. Adapt those exact forms; native editing/navigation
// remains with each control. Modifier state comes only from public Keyable.
func commandKeyEvent(e *fyne.KeyEvent, shift bool, editing bool, route func(fyne.Shortcut)) bool {
	key := string(e.Name)
	function := false
	for _, f := range []fyne.KeyName{fyne.KeyF1, fyne.KeyF2, fyne.KeyF3, fyne.KeyF4, fyne.KeyF5, fyne.KeyF6, fyne.KeyF7, fyne.KeyF8, fyne.KeyF9, fyne.KeyF10, fyne.KeyF11, fyne.KeyF12} {
		if e.Name == f {
			function = true
			break
		}
	}
	modifier := fyne.KeyModifier(0)
	if shift {
		modifier = fyne.KeyModifierShift
	}
	if e.Name == desktop.KeyMenu {
		function = true
		e = &fyne.KeyEvent{Name: fyne.KeyF10}
		modifier = fyne.KeyModifierShift
	}
	if !function && !(shift && !editing && len(key) == 1 && (key[0] >= 'A' && key[0] <= 'Z' || key[0] >= '0' && key[0] <= '9')) {
		return false
	}
	if route != nil {
		route(&desktop.CustomShortcut{KeyName: e.Name, Modifier: modifier})
	}
	return true
}
func shiftKey(e *fyne.KeyEvent) bool {
	return e.Name == desktop.KeyShiftLeft || e.Name == desktop.KeyShiftRight
}

// GLFW emits named editing shortcut values instead of CustomShortcut. Outside
// native Entry editing, normalize those public values back to declared keys.
func commandShortcut(s fyne.Shortcut) fyne.Shortcut {
	name := ""
	switch value := s.(type) {
	case *fyne.ShortcutCopy:
		name = "Primary+C"
		if value.Secondary {
			return &desktop.CustomShortcut{KeyName: fyne.KeyInsert, Modifier: fyne.KeyModifierControl}
		}
	case *fyne.ShortcutCut:
		name = "Primary+X"
		if value.Secondary {
			return &desktop.CustomShortcut{KeyName: fyne.KeyDelete, Modifier: fyne.KeyModifierShift}
		}
	case *fyne.ShortcutPaste:
		name = "Primary+V"
		if value.Secondary {
			return &desktop.CustomShortcut{KeyName: fyne.KeyInsert, Modifier: fyne.KeyModifierShift}
		}
	case *fyne.ShortcutSelectAll:
		name = "Primary+A"
	case *fyne.ShortcutUndo:
		name = "Primary+Z"
	case *fyne.ShortcutRedo:
		name = "Primary+Y"
	}
	if name != "" {
		if shortcut, err := nativeShortcut(name); err == nil {
			return shortcut
		}
	}
	return s
}

// A canvas with no focusable content still accepts declared function keys.
// These public callbacks are scoped to this bundle and restored on disposal.
type commandCanvasInput struct {
	typed, down, up func(*fyne.KeyEvent)
	shift           bool
}

func (b *Bundle) installCanvasKeys(c fyne.Canvas) {
	if b.canvasInputs == nil {
		b.canvasInputs = map[fyne.Canvas]*commandCanvasInput{}
	}
	state := &commandCanvasInput{typed: c.OnTypedKey()}
	b.canvasInputs[c] = state
	c.SetOnTypedKey(func(e *fyne.KeyEvent) {
		handled := false
		commandKeyEvent(e, state.shift, false, func(s fyne.Shortcut) {
			for _, key := range b.keys {
				if key.shortcut.ShortcutName() == s.ShortcutName() {
					handled = true
					b.routeShortcut(c, s)
					return
				}
			}
		})
		if !handled && state.typed != nil {
			state.typed(e)
		}
	})
	if d, ok := c.(desktop.Canvas); ok {
		state.down, state.up = d.OnKeyDown(), d.OnKeyUp()
		d.SetOnKeyDown(func(e *fyne.KeyEvent) {
			if shiftKey(e) {
				state.shift = true
			}
			if state.down != nil {
				state.down(e)
			}
		})
		d.SetOnKeyUp(func(e *fyne.KeyEvent) {
			if shiftKey(e) {
				state.shift = false
			}
			if state.up != nil {
				state.up(e)
			}
		})
	}
}
func (b *Bundle) removeCanvasKeys(c fyne.Canvas) {
	state := b.canvasInputs[c]
	if state == nil {
		return
	}
	delete(b.canvasInputs, c)
	c.SetOnTypedKey(state.typed)
	if d, ok := c.(desktop.Canvas); ok {
		d.SetOnKeyDown(state.down)
		d.SetOnKeyUp(state.up)
	}
}
