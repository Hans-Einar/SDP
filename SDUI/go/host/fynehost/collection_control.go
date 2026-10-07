package fynehost

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"image/color"
	"math"
)

// CollectionControl is one native focus stop. Its rows have stable runtime IDs;
// provider status affordances never pretend to be items. No local selection or
// scroll state is authoritative: all input goes to the published bundle owner.
type CollectionControl struct {
	disabled         bool
	scrollX, scrollY bool
	dragAxis         string
	dragGrab         float64
	widget.BaseWidget
	owner    *DocumentHost
	bundle   *Bundle
	path     string
	state    ui.CollectionState
	rows     []collectionRow
	viewport layout.Viewport
	font     float64
	title    string
	focused  bool
}

func newCollectionControl(h *DocumentHost, b *Bundle, path string) *CollectionControl {
	c := &CollectionControl{owner: h, bundle: b, path: path, font: 14}
	b.view.Root.Walk(func(n *parser.Instance) {
		if n.Path == path {
			c.scrollX = n.Layout["overflow-x"] == "scroll"
			c.scrollY = n.Layout["overflow-y"] == "scroll"
		}
	})
	c.ExtendBaseWidget(c)
	return c
}
func (c *CollectionControl) live() bool {
	if c.owner.current != c.bundle || c.bundle.closed || c.owner.closed || c.bundle.muted || c.disabled {
		return false
	}
	current, ok := c.bundle.Session.Collection(c.state.Handle)
	return ok && current.Generation == c.state.Generation && current.ProviderEpoch == c.state.ProviderEpoch
}
func (c *CollectionControl) CreateRenderer() fyne.WidgetRenderer {
	return &collectionRenderer{c: c, box: container.NewWithoutLayout()}
}
func (c *CollectionControl) FocusGained() {
	c.focused = true
	if !c.live() {
		return
	}
	c.focused = true
	c.owner.focusCollection(c)
	c.Refresh()
}
func (c *CollectionControl) FocusLost() { c.focused = false; c.Refresh() }
func (c *CollectionControl) TypedRune(r rune) {
	if (r == 'r' || r == 'R') && c.live() {
		c.owner.recoverCollection(c)
	}
}
func (c *CollectionControl) TypedKey(e *fyne.KeyEvent) {
	if c.live() {
		c.owner.collectionKey(c, e.Name)
	}
}
func (c *CollectionControl) Scrolled(e *fyne.ScrollEvent) {
	if c.live() {
		c.owner.scroll(c.viewport.Rect.X+float64(e.Position.X), c.viewport.Rect.Y-rowHeight(c.font)+float64(e.Position.Y), -float64(e.Scrolled.DX), -float64(e.Scrolled.DY))
	}
}
func (c *CollectionControl) rowAt(p fyne.Position) (collectionRow, bool) {
	x, y := float64(p.X), float64(p.Y)
	title := rowHeight(c.font)
	if x < 0 || y < title || x >= c.viewport.Rect.W || y >= title+c.viewport.Rect.H {
		return collectionRow{}, false
	}
	// Outer clipping can hide part of the collection even though its own rect fits.
	sx := c.viewport.Rect.X + x
	sy := c.viewport.Rect.Y + y - title
	if !c.viewport.Clip.Contains(sx, sy) {
		return collectionRow{}, false
	}
	y = y - title + c.viewport.Offset.Y
	for _, r := range c.rows {
		if y >= r.Y && y < r.Y+r.H {
			return r, true
		}
	}
	return collectionRow{}, false
}
func (c *CollectionControl) Tapped(e *fyne.PointEvent)       { c.tap(e, false) }
func (c *CollectionControl) DoubleTapped(e *fyne.PointEvent) { c.tap(e, true) }
func (c *CollectionControl) tap(e *fyne.PointEvent, double bool) {
	if !c.live() {
		return
	}
	r, ok := c.rowAt(e.Position)
	if !ok {
		return
	}
	c.owner.canvas.Focus(c)
	if r.Status {
		if r.Recovery {
			c.owner.collectionEvent(c, r.Parent, ui.Retry)
		}
		return
	}
	if r.Item.Kind == ui.Separator {
		return
	}
	c.owner.focusItem(c, r.Item.ID, false)
	if r.Item.HasChildren && float64(e.Position.X)+c.viewport.Offset.X < float64(r.Depth)*18+24 {
		c.owner.toggle(c, r.Item.ID)
		return
	}
	if r.Item.Kind == ui.Row {
		c.owner.collectionEvent(c, r.Item.ID, ui.Select)
		if double {
			c.owner.collectionEvent(c, r.Item.ID, ui.Activate)
		}
	} else if double && r.Item.HasChildren {
		c.owner.toggle(c, r.Item.ID)
	}
}

// Dragging the visible gutter mirrors native thumb semantics through runtime.
func (c *CollectionControl) Dragged(e *fyne.DragEvent) {
	if !c.live() {
		return
	}
	x, y := float64(e.Position.X), float64(e.Position.Y)-rowHeight(c.font)
	if c.dragAxis == "" {
		startX, startY := x-float64(e.Dragged.DX), y-float64(e.Dragged.DY)
		if startX >= c.viewport.Rect.W && c.viewport.Maximum.Y > 0 {
			length := thumbLength(c.viewport.Rect.H, c.viewport.Content.H)
			start := (c.viewport.Rect.H - length) * c.viewport.Offset.Y / c.viewport.Maximum.Y
			if startY < start || startY > start+length {
				return
			}
			c.dragAxis = "y"
			c.dragGrab = startY - start
		} else if startY >= c.viewport.Rect.H && c.viewport.Maximum.X > 0 {
			length := thumbLength(c.viewport.Rect.W, c.viewport.Content.W)
			start := (c.viewport.Rect.W - length) * c.viewport.Offset.X / c.viewport.Maximum.X
			if startX < start || startX > start+length {
				return
			}
			c.dragAxis = "x"
			c.dragGrab = startX - start
		} else {
			return
		}
	}
	next := c.viewport.Offset
	if c.dragAxis == "y" {
		travel := c.viewport.Rect.H - thumbLength(c.viewport.Rect.H, c.viewport.Content.H)
		next.Y = math.Max(0, y-c.dragGrab) / math.Max(1, travel) * c.viewport.Maximum.Y
	} else {
		travel := c.viewport.Rect.W - thumbLength(c.viewport.Rect.W, c.viewport.Content.W)
		next.X = math.Max(0, x-c.dragGrab) / math.Max(1, travel) * c.viewport.Maximum.X
	}
	c.owner.setViewport(c.path, next)
}
func (c *CollectionControl) DragEnd()       { c.dragAxis = "" }
func (c *CollectionControl) Disabled() bool { return c.disabled }
func (c *CollectionControl) Disable()       { c.disabled = true }
func (c *CollectionControl) Enable()        { c.disabled = false }
func (c *CollectionControl) TypedShortcut(s fyne.Shortcut) {
	key, ok := s.(*desktop.CustomShortcut)
	if !ok || key.Modifier != fyne.KeyModifierAlt || !c.live() {
		return
	}
	dx, dy := 0.0, 0.0
	step := rowHeight(c.font)
	switch key.KeyName {
	case fyne.KeyLeft:
		dx = -step
	case fyne.KeyRight:
		dx = step
	case fyne.KeyUp:
		dy = -step
	case fyne.KeyDown:
		dy = step
	default:
		return
	}
	c.owner.scroll(c.viewport.Clip.X+c.viewport.Clip.W/2, c.viewport.Clip.Y+c.viewport.Clip.H/2, dx, dy)
}
func thumbLength(view, content float64) float64 {
	return math.Min(view, math.Max(15, view*view/math.Max(1, content)))
}

type collectionRenderer struct {
	c   *CollectionControl
	box *fyne.Container
}

func (r *collectionRenderer) MinSize() fyne.Size {
	m := collectionMinimum(r.c.font, r.c.scrollX, r.c.scrollY)
	return fyne.NewSize(float32(m.W), float32(m.H))
}
func (r *collectionRenderer) Layout(size fyne.Size)        { r.box.Resize(size) }
func (r *collectionRenderer) Objects() []fyne.CanvasObject { return []fyne.CanvasObject{r.box} }
func (r *collectionRenderer) Destroy()                     {}
func (r *collectionRenderer) Refresh() {
	c := r.c
	objects := []fyne.CanvasObject{}
	height := rowHeight(c.font)
	addText := func(text string, x, y float64, col color.Color) {
		t := canvas.NewText(text, col)
		t.TextSize = float32(c.font)
		t.FontSource = regular
		t.Resize(t.MinSize())
		t.Move(fyne.NewPos(float32(x), float32(y)))
		objects = append(objects, t)
	}
	addText(c.title, 4, 2, color.NRGBA{R: 25, G: 40, B: 60, A: 255})
	// Draw only visible rows; the enclosing ScrollNone viewport clips horizontal
	// labels and partial rows. Title remains outside the translated row content.
	content := container.NewWithoutLayout()
	for _, row := range c.rows {
		y := row.Y - c.viewport.Offset.Y
		if y+row.H <= 0 || y >= c.viewport.Rect.H {
			continue
		}
		if row.Item.Kind == ui.Separator {
			line := canvas.NewRectangle(color.NRGBA{R: 140, G: 150, B: 165, A: 255})
			line.Move(fyne.NewPos(0, float32(y+row.H/2)))
			line.Resize(fyne.NewSize(float32(c.viewport.Rect.W), 1))
			content.Add(line)
			continue
		}
		if row.Item.ID != "" && (row.Item.ID == c.state.Selected || c.focused && row.Item.ID == c.state.Focused) {
			bg := canvas.NewRectangle(color.NRGBA{R: 210, G: 225, B: 245, A: 255})
			bg.Move(fyne.NewPos(0, float32(y)))
			bg.Resize(fyne.NewSize(float32(c.viewport.Rect.W), float32(row.H)))
			content.Add(bg)
		}
		label, indent := fittedRow(row, c.state, c.font, c.viewport.Rect.W)
		offsetX := c.viewport.Offset.X
		if row.Status {
			offsetX = 0
		}
		text := canvas.NewText(label, color.NRGBA{R: 30, G: 40, B: 50, A: 255})
		text.FontSource = regular
		text.TextSize = float32(c.font)
		text.Resize(text.MinSize())
		text.Move(fyne.NewPos(float32(indent+4-offsetX), float32(y+3)))
		content.Add(text)
	}
	clip := container.NewClip(content)
	clip.Move(fyne.NewPos(0, float32(height)))
	clip.Resize(fyne.NewSize(float32(c.viewport.Rect.W), float32(c.viewport.Rect.H)))
	objects = append(objects, clip)
	if c.viewport.Maximum.Y > 0 {
		thumb := canvas.NewRectangle(color.NRGBA{R: 90, G: 110, B: 135, A: 255})
		h := thumbLength(c.viewport.Rect.H, c.viewport.Content.H)
		y := height + (c.viewport.Rect.H-h)*c.viewport.Offset.Y/c.viewport.Maximum.Y
		thumb.Move(fyne.NewPos(float32(c.viewport.Rect.W), float32(y)))
		thumb.Resize(fyne.NewSize(10, float32(h)))
		objects = append(objects, thumb)
	}
	if c.viewport.Maximum.X > 0 {
		thumb := canvas.NewRectangle(color.NRGBA{R: 90, G: 110, B: 135, A: 255})
		w := thumbLength(c.viewport.Rect.W, c.viewport.Content.W)
		x := (c.viewport.Rect.W - w) * c.viewport.Offset.X / c.viewport.Maximum.X
		thumb.Move(fyne.NewPos(float32(x), float32(height+c.viewport.Rect.H)))
		thumb.Resize(fyne.NewSize(float32(w), 10))
		objects = append(objects, thumb)
	}
	r.box.Objects = objects
	r.box.Refresh()
}
