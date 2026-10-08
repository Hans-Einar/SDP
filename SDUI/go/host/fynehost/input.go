package fynehost

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// Input adds the SDUI draft-revert hook to Fyne's native text editing widget.
type Input struct {
	widget.Entry
	OnRevert     func()
	OnFocus      func()
	OnEscape     func() bool
	OnContext    func(fyne.Position)
	OnShortcut   func(fyne.Shortcut)
	commandShift bool
}

func NewInput() *Input { i := &Input{}; i.ExtendBaseWidget(i); return i }
func (i *Input) TypedKey(e *fyne.KeyEvent) {
	if commandKeyEvent(e, i.commandShift, true, i.OnShortcut) {
		return
	}
	if e.Name == fyne.KeyEscape && i.OnEscape != nil && i.OnEscape() {
		return
	}
	if e.Name == fyne.KeyEscape && i.OnRevert != nil {
		i.OnRevert()
		return
	}
	i.Entry.TypedKey(e)
}
func (i *Input) TappedSecondary(e *fyne.PointEvent) {
	if i.OnContext != nil {
		i.OnContext(e.AbsolutePosition)
		return
	}
	i.Entry.TappedSecondary(e)
}
func (i *Input) FocusGained() {
	i.Entry.FocusGained()
	if i.OnFocus != nil {
		i.OnFocus()
	}
}

func nativeEditingShortcut(s fyne.Shortcut) bool {
	switch s.(type) {
	case *fyne.ShortcutCopy, *fyne.ShortcutCut, *fyne.ShortcutPaste, *fyne.ShortcutSelectAll, *fyne.ShortcutUndo, *fyne.ShortcutRedo:
		return true
	}
	return false
}
func (i *Input) TypedShortcut(s fyne.Shortcut) {
	if nativeEditingShortcut(s) {
		i.Entry.TypedShortcut(s)
		return
	}
	i.Entry.TypedShortcut(s)
	if i.OnShortcut != nil {
		i.OnShortcut(s)
	}
}

func (i *Input) KeyDown(e *fyne.KeyEvent) {
	i.Entry.KeyDown(e)
	if shiftKey(e) {
		i.commandShift = true
	}
}
func (i *Input) KeyUp(e *fyne.KeyEvent) {
	i.Entry.KeyUp(e)
	if shiftKey(e) {
		i.commandShift = false
	}
}
func (i *Input) FocusLost() { i.commandShift = false; i.Entry.FocusLost() }
