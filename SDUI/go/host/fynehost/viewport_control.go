package fynehost

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"image/color"
	"sort"
)

// routedClip preserves Fyne clipping while forwarding wheel gestures from a
// button/input wrapper instead of swallowing them in ScrollNone.
type routedClip struct {
	*container.Scroll
	onScroll func(*fyne.ScrollEvent)
}

func newRoutedClip(content fyne.CanvasObject) *routedClip {
	s := &routedClip{Scroll: &container.Scroll{Content: content}}
	s.Direction = container.ScrollNone
	s.ExtendBaseWidget(s)
	return s
}
func (s *routedClip) Scrolled(e *fyne.ScrollEvent) {
	if s.onScroll != nil {
		s.onScroll(e)
	} else {
		s.Scroll.Scrolled(e)
	}
}

type viewportSurface struct {
	widget.BaseWidget
	bundle   *Bundle
	viewport layout.Viewport
	handle   ui.Handle
	axis     string
}

func (v *viewportSurface) CreateRenderer() fyne.WidgetRenderer {
	col := color.NRGBA{}
	if v.axis != "" {
		col = color.NRGBA{R: 90, G: 110, B: 135, A: 255}
	}
	return widget.NewSimpleRenderer(canvas.NewRectangle(col))
}
func (v *viewportSurface) live() bool {
	return !v.bundle.closed && v.bundle.owner.current == v.bundle && !v.bundle.owner.closed
}
func (v *viewportSurface) Scrolled(e *fyne.ScrollEvent) {
	if v.live() {
		v.bundle.owner.scroll(float64(v.Position().X+e.Position.X), float64(v.Position().Y+e.Position.Y), -float64(e.Scrolled.DX), -float64(e.Scrolled.DY))
	}
}
func (v *viewportSurface) Dragged(e *fyne.DragEvent) {
	if !v.live() || v.axis == "" {
		return
	}
	s := v.bundle.Session
	current, ok := s.Viewport(v.viewport.Path)
	if !ok || current != v.handle {
		return
	}
	state := s.Snapshot().Viewports[v.viewport.Path]
	if v.axis == "y" {
		travel := v.viewport.Rect.H - thumbLength(v.viewport.Rect.H, v.viewport.Content.H)
		if travel <= 0 {
			return
		}
		state.Y += float64(e.Dragged.DY) * v.viewport.Maximum.Y / travel
	} else {
		travel := v.viewport.Rect.W - thumbLength(v.viewport.Rect.W, v.viewport.Content.W)
		if travel <= 0 {
			return
		}
		state.X += float64(e.Dragged.DX) * v.viewport.Maximum.X / travel
	}
	if state.X < 0 {
		state.X = 0
	}
	if state.Y < 0 {
		state.Y = 0
	}
	v.bundle.owner.status(v.bundle.owner.Mutate(func(s *ui.Session) error { return s.SetViewport(v.handle, v.bundle.Session.Revision, state) }))
}
func (v *viewportSurface) DragEnd() {}
func (b *Bundle) viewportObjects(snapshot ui.Snapshot, g *layout.SnapshotLayout) ([]fyne.CanvasObject, []fyne.CanvasObject) {
	backgrounds, thumbs := []fyne.CanvasObject{}, []fyne.CanvasObject{}
	paths := []string{}
	for path := range g.Viewports {
		if _, collection := snapshot.Collections[path]; !collection {
			paths = append(paths, path)
		}
	}
	sort.Strings(paths)
	for _, path := range paths {
		v := g.Viewports[path]
		makeSurface := func(axis string, rect layout.Rect) *viewportSurface {
			o := &viewportSurface{bundle: b, viewport: v, handle: snapshot.ViewportHandles[path], axis: axis}
			o.ExtendBaseWidget(o)
			o.Move(fyne.NewPos(float32(rect.X), float32(rect.Y)))
			o.Resize(fyne.NewSize(float32(rect.W), float32(rect.H)))
			return o
		}
		for _, strip := range []layout.Rect{v.Clip, v.VerticalGutter, v.HorizontalGutter} {
			if strip.W > 0 && strip.H > 0 {
				backgrounds = append(backgrounds, makeSurface("", strip))
			}
		}
		if v.Maximum.Y > 0 {
			height := thumbLength(v.Rect.H, v.Content.H)
			rect := layout.Rect{X: v.VerticalGutter.X, Y: v.Rect.Y + (v.Rect.H-height)*v.Offset.Y/v.Maximum.Y, W: v.VerticalGutter.W, H: height}.Intersect(v.VerticalGutter)
			if rect.W > 0 && rect.H > 0 {
				thumbs = append(thumbs, makeSurface("y", rect))
			}
		}
		if v.Maximum.X > 0 {
			width := thumbLength(v.Rect.W, v.Content.W)
			rect := layout.Rect{X: v.Rect.X + (v.Rect.W-width)*v.Offset.X/v.Maximum.X, Y: v.HorizontalGutter.Y, W: width, H: v.HorizontalGutter.H}.Intersect(v.HorizontalGutter)
			if rect.W > 0 && rect.H > 0 {
				thumbs = append(thumbs, makeSurface("x", rect))
			}
		}
	}
	return backgrounds, thumbs
}

// readingOrder follows header, body/ordinary rows, footer rather than putting
// every region before the ordinary body merely because Instance.Walk does so.
func readingOrder(n *parser.Instance, visit func(*parser.Instance)) {
	visit(n)
	if c := n.Region("header"); c != nil {
		readingOrder(c, visit)
	}
	if c := n.Region("body"); c != nil {
		readingOrder(c, visit)
	}
	for _, row := range n.Rows {
		for _, c := range row {
			readingOrder(c, visit)
		}
	}
	if c := n.Region("footer"); c != nil {
		readingOrder(c, visit)
	}
}
