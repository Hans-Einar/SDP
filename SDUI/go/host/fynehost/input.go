package fynehost

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// Input adds the SDUI draft-revert hook to Fyne's native text editing widget.
type Input struct {
	widget.Entry
	OnRevert func()
	OnFocus  func()
}

func NewInput() *Input { i := &Input{}; i.ExtendBaseWidget(i); return i }
func (i *Input) TypedKey(e *fyne.KeyEvent) {
	if e.Name == fyne.KeyEscape && i.OnRevert != nil {
		i.OnRevert()
		return
	}
	i.Entry.TypedKey(e)
}
func (i *Input) FocusGained() {
	i.Entry.FocusGained()
	if i.OnFocus != nil {
		i.OnFocus()
	}
}
