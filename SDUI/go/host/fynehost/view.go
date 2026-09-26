// Package fynehost adapts shared SDUI geometry to native Fyne controls.
// All methods are called on Fyne's UI goroutine. Callbacks are registered by
// the embedding Go application; SDL references are never resolved here.
package fynehost

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/svg"
)

type Action func(path, value string) error
type control struct {
	widget fyne.CanvasObject
	theme  *container.ThemeOverride
	clip   *container.Scroll
	fixed  *fixedLayout
}
type View struct {
	Root      *parser.Instance
	Container *fyne.Container
	Geometry  *layout.Box
	Controls  map[string]fyne.CanvasObject
	Actions   map[string]Action
	OnStatus  func(error)
	Measure   layout.Measurer
	Content   svg.ContentRenderer
	controls  map[string]*control
	image     *canvas.Image
	reflowing bool
	closed    bool
}

func New(root *parser.Instance) *View {
	v := &View{Root: root, Controls: map[string]fyne.CanvasObject{}, Actions: map[string]Action{}, controls: map[string]*control{}}
	v.image = canvas.NewImageFromResource(fyne.NewStaticResource("empty.svg", []byte(`<svg xmlns="http://www.w3.org/2000/svg" width="1" height="1"/>`)))
	v.image.FillMode = canvas.ImageFillStretch
	v.Container = container.New(v, v.image)
	root.Walk(func(n *parser.Instance) {
		if n.Kind != "widget" || n.Widget == "svg" {
			return
		}
		var obj fyne.CanvasObject
		switch n.Widget {
		case "button":
			obj = widget.NewButton(n.Argument("label"), func() { v.invoke(n.Path, "") })
		case "input":
			entry := NewInput()
			entry.SetPlaceHolder(n.Argument("text"))
			entry.SetText(n.Argument("value"))
			entry.OnChanged = func(text string) { v.invoke(n.Path, text) }
			obj = entry
		}
		if obj == nil {
			return
		}
		th := container.NewThemeOverride(obj, componentTheme{14})
		fixed := &fixedLayout{}
		holder := container.New(fixed, th)
		clip := container.NewScroll(holder)
		clip.Direction = container.ScrollNone
		v.Controls[n.Path] = obj
		v.controls[n.Path] = &control{obj, th, clip, fixed}
		v.Container.Add(clip)
	})
	return v
}
func (v *View) invoke(path, value string) {
	if v.closed {
		return
	}
	fn := v.Actions[path]
	var err error
	if fn == nil {
		err = fmt.Errorf("unbound: %s", path)
	} else {
		err = fn(path, value)
	}
	if v.OnStatus != nil {
		v.OnStatus(err)
	}
}
func (v *View) Close()                                { v.closed = true; v.Actions = map[string]Action{} }
func (v *View) MinSize([]fyne.CanvasObject) fyne.Size { return fyne.NewSize(1, 1) }
func (v *View) Layout(_ []fyne.CanvasObject, size fyne.Size) {
	if v.reflowing || size.Width <= 0 || size.Height <= 0 {
		return
	}
	v.reflowing = true
	defer func() { v.reflowing = false }()
	if err := v.Reflow(layout.Size{W: float64(size.Width), H: float64(size.Height)}); v.OnStatus != nil {
		v.OnStatus(err)
	}
}
func (v *View) Reflow(size layout.Size) error {
	tree, err := (&layout.Engine{Measure: v.Measure}).Layout(v.Root, size)
	if err != nil {
		return err
	}
	background, err := svg.Render(tree, svg.Options{Width: size.W, Height: size.H, SkipControls: true, Content: v.Content})
	if err != nil {
		return err
	}
	v.Geometry = tree
	v.image.Resource = fyne.NewStaticResource("sdui.svg", []byte(background))
	v.image.Resize(fyne.NewSize(float32(size.W), float32(size.H)))
	v.image.Refresh()
	for _, c := range v.controls {
		c.clip.Hide()
	}
	tree.Walk(func(b *layout.Box) {
		c := v.controls[b.Path]
		if c == nil {
			return
		}
		r := b.Rect
		clip := b.Clip
		if clip.W <= 0 || clip.H <= 0 {
			return
		}
		c.theme.Theme = componentTheme{float32(b.Font)}
		c.theme.Refresh()
		c.fixed.size = fyne.NewSize(float32(r.W), float32(r.H))
		c.clip.Content.Resize(c.fixed.size)
		c.clip.Resize(fyne.NewSize(float32(clip.W), float32(clip.H)))
		c.clip.Move(fyne.NewPos(float32(clip.X), float32(clip.Y)))
		c.clip.Offset = fyne.NewPos(float32(clip.X-r.X), float32(clip.Y-r.Y))
		c.clip.Refresh()
		c.clip.Show()
		if dis, ok := c.widget.(fyne.Disableable); ok {
			if b.Enabled {
				dis.Enable()
			} else {
				dis.Disable()
			}
		}
	})
	return nil
}

type fixedLayout struct{ size fyne.Size }

func (f *fixedLayout) MinSize([]fyne.CanvasObject) fyne.Size { return f.size }
func (f *fixedLayout) Layout(objects []fyne.CanvasObject, _ fyne.Size) {
	for _, o := range objects {
		o.Move(fyne.NewPos(0, 0))
		o.Resize(f.size)
	}
}
