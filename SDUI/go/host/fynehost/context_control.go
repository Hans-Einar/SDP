package fynehost

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"image/color"
)

// The existing SVG placeholder remains background-painted. This transparent,
// nonfocusable native target adds only its declared secondary-click context.
type contextControl struct {
	widget.BaseWidget
	disabled bool
	context  func(fyne.Position)
}

func newContextControl() *contextControl { c := &contextControl{}; c.ExtendBaseWidget(c); return c }
func (c *contextControl) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(canvas.NewRectangle(color.Transparent))
}
func (c *contextControl) TappedSecondary(e *fyne.PointEvent) {
	if !c.disabled && c.context != nil {
		c.context(e.AbsolutePosition)
	}
}
func (c *contextControl) Disabled() bool { return c.disabled }
func (c *contextControl) Disable()       { c.disabled = true }
func (c *contextControl) Enable()        { c.disabled = false }
