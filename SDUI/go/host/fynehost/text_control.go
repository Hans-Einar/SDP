package fynehost

import (
	"fmt"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

// textControl retains one native editor. Runtime owns its accepted/draft text;
// Entry alone owns selection, caret, undo and its internal editing viewport.
type textControl struct {
	widget.BaseWidget
	bundle          *Bundle
	path            string
	state           ui.FieldState
	label, feedback *widget.Label
	entry           *textEntry
	parts           layout.FieldMetrics
	editMenu        *nativeMenu
}

func extendedInput(n *parser.Instance) bool {
	if n.Widget != "input" {
		return false
	}
	p, err := parser.InputOptions(n)
	return err == nil && p.Extended
}

func newTextControl(b *Bundle, n *parser.Instance) *textControl {
	return newTextControlValue(b, n, n.Argument("value"))
}
func newTextControlValue(b *Bundle, n *parser.Instance, value string) *textControl {
	c := &textControl{bundle: b, path: n.Path, label: widget.NewLabel(n.Argument("text")), feedback: widget.NewLabel("")}
	c.ExtendBaseWidget(c)
	c.label.Truncation = fyne.TextTruncateEllipsis
	c.feedback.Truncation = fyne.TextTruncateEllipsis
	c.entry = &textEntry{owner: c}
	c.entry.ExtendBaseWidget(c.entry)
	policy, _ := parser.InputOptions(n) // detached preparation has validated policy
	c.entry.MultiLine = policy.Multiline
	if policy.Multiline {
		c.entry.Wrapping = fyne.TextWrapWord
		c.entry.SetMinRowsVisible(3)
	}
	c.entry.SetPlaceHolder(policy.Placeholder)
	c.entry.SetText(value)
	c.entry.OnChanged = c.changed
	return c
}
func (c *textControl) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.New(textPartsLayout{c}, c.label, c.entry, c.feedback))
}
func (c *textControl) MinSize() fyne.Size {
	return fyne.NewSize(float32(c.parts.Minimum.W), float32(c.parts.Minimum.H))
}
func (c *textControl) Disable()       { c.entry.Disable() }
func (c *textControl) Enable()        { c.entry.Enable() }
func (c *textControl) Disabled() bool { return c.entry.Disabled() }

type textPartsLayout struct{ c *textControl }

func (p textPartsLayout) MinSize([]fyne.CanvasObject) fyne.Size { return p.c.MinSize() }
func (p textPartsLayout) Layout(_ []fyne.CanvasObject, _ fyne.Size) {
	for _, part := range []struct {
		o fyne.CanvasObject
		r layout.Rect
	}{{p.c.label, p.c.parts.Label}, {p.c.entry, p.c.parts.Control}, {p.c.feedback, p.c.parts.Feedback}} {
		part.o.Move(fyne.NewPos(float32(part.r.X), float32(part.r.Y)))
		part.o.Resize(fyne.NewSize(float32(part.r.W), float32(part.r.H)))
	}
}
func (c *textControl) live() bool { return c.bundle != nil && c.bundle.livePane() }
func (c *textControl) field() (ui.FieldState, bool) {
	if !c.live() {
		return ui.FieldState{}, false
	}
	return c.bundle.Session.Field(c.state.Target.Handle)
}
func (c *textControl) editable() bool {
	f, ok := c.field()
	if !ok || f.ReadOnly || c.Disabled() {
		return false
	}
	w, ok := c.bundle.Session.Widget(c.path)
	return ok && w.Handle == f.Target.Handle && w.Visible && w.Enabled
}
func (c *textControl) changed(raw string) {
	if !c.live() {
		return
	}
	b := c.bundle
	err := b.owner.Mutate(func(s *ui.Session) error {
		_, err := s.EditField(c.state.Target.Handle, c.state.Target.ModelRevision, ui.Text(raw))
		return err
	})
	// OnChanged is post-mutation in Fyne. Only this rejection path deliberately
	// clears native history, even when reentrant work has left equal text.
	if err != nil && b.owner.current == b && !b.closed {
		if f, ok := b.Session.Field(c.state.Target.Handle); ok {
			muted := b.muted
			b.muted = true
			c.entry.SetText(f.Proposed.Text)
			b.muted = muted
		}
	}
}
func (c *textControl) commit() {
	if !c.editable() {
		return
	}
	c.bundle.owner.Mutate(func(s *ui.Session) error {
		f, ok := s.Field(c.state.Target.Handle)
		if !ok {
			return fmt.Errorf("stale-field")
		}
		event, err := s.CaptureCommit(f.Target)
		if err != nil {
			return err
		}
		return s.Dispatch(event)
	})
}
func (c *textControl) escape() {
	if !c.live() {
		return
	}
	if c.editMenu != nil {
		c.editMenu.close()
		return
	}
	if f, ok := c.field(); ok && f.Dirty {
		c.bundle.owner.Mutate(func(s *ui.Session) error { return s.RevertField(f.Target.Handle) })
		return
	}
	c.bundle.escapeSurface(c.path)
}
func (c *textControl) focus() {
	if !c.live() || c.Disabled() {
		return
	}
	if err := c.bundle.owner.Mutate(func(s *ui.Session) error { return s.Focus(c.state.Target.Handle) }); err == nil {
		c.bundle.owner.ensureWidget(c.path)
	}
}
func (c *textControl) shortcut(s fyne.Shortcut) {
	if c.live() {
		c.bundle.routeShortcut(c.bundle.canvasFor(c.path), s)
	}
}
func (c *textControl) sync(f ui.FieldState, parts layout.FieldLayout, box *layout.Box) {
	c.state = f
	c.entry.readOnly = f.ReadOnly
	local := func(r layout.Rect) layout.Rect {
		if r.W == 0 || r.H == 0 {
			return layout.Rect{}
		}
		r.X -= box.Rect.X
		r.Y -= box.Rect.Y
		return r
	}
	c.parts = layout.FieldMetrics{Label: local(parts.Label), Control: local(parts.Control), Feedback: local(parts.Feedback)}
	c.label.SetText(box.Instance.Argument("text"))
	if c.label.Text == "" {
		c.label.Hide()
	} else {
		c.label.Show()
	}
	message := f.Validation.Message
	if message == "" {
		message = f.Validation.Code
	}
	c.feedback.SetText(message)
	if f.Input != nil && c.entry.PlaceHolder != f.Input.Placeholder {
		c.entry.SetPlaceHolder(f.Input.Placeholder)
	}
	if c.entry.Text != f.Proposed.Text {
		c.entry.SetText(f.Proposed.Text)
	}
	c.Refresh()
}
func measureTextField(n *parser.Instance, f ui.FieldState, font float64, outer layout.Size) (layout.FieldMetrics, error) {
	p, err := parser.InputOptions(n)
	if err != nil {
		return layout.FieldMetrics{}, err
	}
	if !p.Extended || f.Input == nil || f.Input.Multiline != p.Multiline || f.Input.Placeholder != p.Placeholder {
		return layout.FieldMetrics{}, fmt.Errorf("native-input-policy: %s", n.Path)
	}
	// Both configured Entry modes scroll. Pinned Entry.MinSize uses the themed
	// character minimum and row count, not its text extent. Avoid shaping the
	// complete draft for every detached admission probe; live editors still get
	// their exact text. Font/short/long equivalence is tested against real Entry.
	c := newTextControlValue(nil, n, "")
	th := componentTheme{float32(font)}
	body := container.NewThemeOverride(c.entry, th).MinSize()
	label := fyne.Size{}
	if n.Argument("text") != "" {
		c.label.Truncation = fyne.TextTruncateOff
		label = container.NewThemeOverride(c.label, th).MinSize()
	}
	feedback := widget.NewLabel("Invalid")
	feedback.Truncation = fyne.TextTruncateEllipsis
	fh := float64(container.NewThemeOverride(feedback, th).MinSize().Height)
	w, h := math.Max(float64(body.Width), float64(label.Width)), float64(body.Height+label.Height)+fh
	width, height := math.Max(w, outer.W), math.Max(h, outer.H)
	r := layout.FieldMetrics{Minimum: layout.Size{W: w, H: h}, Control: layout.Rect{Y: float64(label.Height), W: width, H: height - float64(label.Height) - fh}, Feedback: layout.Rect{Y: height - fh, W: width, H: fh}}
	if label.Height > 0 {
		r.Label = layout.Rect{W: width, H: float64(label.Height)}
	}
	return r, nil
}
