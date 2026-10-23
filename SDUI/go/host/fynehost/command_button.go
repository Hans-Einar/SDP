package fynehost

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

// commandButton retains native Button input/drawing; checked state is projected
// only after runtime publication, never toggled optimistically by the widget.
type commandButton struct {
	widget.Button
	tooltip          string
	shift            bool
	shortcut         func(fyne.Shortcut)
	tipChanged       func(bool)
	focused, hovered bool
	escape           func()
	focus            func()
	context          func(fyne.Position)
}

func newCommandButton(label string, tapped func()) *commandButton {
	b := &commandButton{}
	b.ExtendBaseWidget(b)
	b.Text, b.OnTapped = label, tapped
	return b
}
func (b *commandButton) FocusGained() {
	b.Button.FocusGained()
	b.focused = true
	if b.focus != nil {
		b.focus()
	}
	b.showTip()
}
func (b *commandButton) FocusLost() {
	b.shift = false
	b.Button.FocusLost()
	b.focused = false
	if !b.hovered {
		b.hideTip()
	}
}
func (b *commandButton) MouseIn(e *desktop.MouseEvent) {
	b.Button.MouseIn(e)
	b.hovered = true
	b.showTip()
}
func (b *commandButton) MouseOut() {
	b.Button.MouseOut()
	b.hovered = false
	if !b.focused {
		b.hideTip()
	}
}
func (b *commandButton) TypedKey(e *fyne.KeyEvent) {
	if commandKeyEvent(e, b.shift, false, b.shortcut) {
		return
	}
	if e.Name == fyne.KeyEscape && b.escape != nil {
		b.escape()
		return
	}
	b.Button.TypedKey(e)
}
func (b *commandButton) TappedSecondary(e *fyne.PointEvent) {
	if !b.Disabled() && b.context != nil {
		b.context(e.AbsolutePosition)
	}
}
func (b *commandButton) showTip() {
	if b.tooltip == "" || b.Disabled() {
		return
	}
	if b.tipChanged != nil {
		b.tipChanged(true)
	}
}
func (b *commandButton) hideTip() {
	if b.tipChanged != nil {
		b.tipChanged(false)
	}
}
func commandLabel(label string, c ui.CommandState) string {
	if c.Toggle {
		if c.Checked {
			return "[x] " + label
		}
		return "[ ] " + label
	}
	return label
}
func (m *collectionMeasure) MeasureMenu(n *parser.Instance, font float64) (layout.Size, error) {
	b := newCommandButton(n.Argument("label"), nil)
	th := container.NewThemeOverride(b, componentTheme{float32(font)})
	s := th.MinSize()
	return layout.Size{W: float64(s.Width), H: float64(s.Height)}, nil
}

func (m *collectionMeasure) measureButton(n *parser.Instance, font float64) layout.Size {
	label, icon := n.Argument("label"), n.Argument("icon")
	if p, ok := m.snapshot.Presentations[n.Path]; ok {
		label, icon = p.Label, p.Icon
		for _, c := range m.snapshot.Commands {
			if c.Handle == p.Command {
				label = commandLabel(label, c)
				break
			}
		}
	}
	b := newCommandButton(label, nil)
	b.SetIcon(m.icons[icon])
	s := container.NewThemeOverride(b, componentTheme{float32(font)}).MinSize()
	return layout.Size{W: float64(s.Width), H: float64(s.Height)}
}

func (b *commandButton) KeyDown(e *fyne.KeyEvent) {
	if shiftKey(e) {
		b.shift = true
	}
}
func (b *commandButton) KeyUp(e *fyne.KeyEvent) {
	if shiftKey(e) {
		b.shift = false
	}
}

func (b *commandButton) Tapped(e *fyne.PointEvent) {
	if b.Disabled() {
		return
	}
	if c := fyne.CurrentApp().Driver().CanvasForObject(b); c != nil {
		c.Focus(b)
	}
	b.Button.Tapped(e)
}
