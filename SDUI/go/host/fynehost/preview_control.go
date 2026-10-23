package fynehost

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/markdown"
)

// A preview paints only admitted immutable content. Markdown itself remains in
// the shared background; its transparent native wrapper supplies accessibility.
type previewControl struct {
	widget.BaseWidget
	kind, description, status, accessible string
	image, text                           *canvas.Image
	frame                                 previewFrame
	disabled                              bool
	context                               func(fyne.Position)
}

type previewFrame struct {
	box, image, caption, status layout.Rect
	resource, text              fyne.Resource
}

func newPreviewControl(p *markdown.Previews, path string) *previewControl {
	o, _ := p.Outcome(path)
	c := &previewControl{kind: o.Kind, description: o.Description, status: o.Status, accessible: previewAccessible(o)}
	c.image = &canvas.Image{FillMode: canvas.ImageFillStretch}
	c.text = &canvas.Image{FillMode: canvas.ImageFillStretch}
	c.image.Hide()
	c.text.Hide()
	c.ExtendBaseWidget(c)
	return c
}

func (c *previewControl) AccessibilityLabel() string             { return c.accessible }
func (c *previewControl) AccessibilityRole() fyne.AccessibleRole { return fyne.AccessibleRoleText }
func (c *previewControl) Disabled() bool                         { return c.disabled }
func (c *previewControl) Disable()                               { c.disabled = true }
func (c *previewControl) Enable()                                { c.disabled = false }
func (c *previewControl) TappedSecondary(e *fyne.PointEvent) {
	if !c.disabled && c.context != nil {
		c.context(e.AbsolutePosition)
	}
}
func (c *previewControl) apply(f previewFrame) {
	c.frame = f
	c.image.Resource = f.resource
	c.text.Resource = f.text
	if f.resource == nil {
		c.image.Image = nil
	}
	if f.text == nil {
		c.text.Image = nil
	}
	c.Refresh()
}
func (c *previewControl) release() {
	c.context = nil
	c.frame = previewFrame{}
	c.image.Resource, c.image.Image = nil, nil
	c.text.Resource, c.text.Image = nil, nil
}
func (c *previewControl) CreateRenderer() fyne.WidgetRenderer { return &previewRenderer{control: c} }

type previewRenderer struct{ control *previewControl }

func (r *previewRenderer) MinSize() fyne.Size { return fyne.NewSize(0, 0) }
func (r *previewRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.control.image, r.control.text}
}
func (r *previewRenderer) Destroy() {}
func (r *previewRenderer) Refresh() { r.Layout(r.control.Size()); canvas.Refresh(r.control) }
func (r *previewRenderer) Layout(_ fyne.Size) {
	c := r.control
	place := func(image *canvas.Image, rect layout.Rect, resource fyne.Resource) {
		if resource == nil || rect.W <= 0 || rect.H <= 0 {
			image.Hide()
			return
		}
		image.Move(fyne.NewPos(float32(rect.X-c.frame.box.X), float32(rect.Y-c.frame.box.Y)))
		image.Resize(fyne.NewSize(float32(rect.W), float32(rect.H)))
		image.Show()
		image.Refresh()
	}
	place(c.image, c.frame.image, c.frame.resource)
	place(c.text, c.frame.box, c.frame.text)
}

var _ fyne.Accessible = (*previewControl)(nil)
