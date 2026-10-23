package fynehost

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"image/color"
)

// paneHeader reuses AppTabs' native header drawing/pointer handling, adding one
// keyboard focus stop. Bodies are laid out separately by shared pane geometry.
// A request never becomes authority: sync always restores the accepted selection.
type paneHeader struct {
	widget.BaseWidget
	tabs                     *nativeTabs
	pointer                  *paneHeaderPointer
	pages                    []headerPage
	selected                 string
	muted, disabled, focused bool
	request                  func(string)
	focus                    func()
	enter                    func()
	backward                 func()
	escape                   func()
	shift                    bool
	shortcut                 func(fyne.Shortcut)
}
type headerPage struct {
	id, label string
	enabled   bool
	icon      fyne.Resource
}

func newPaneHeader() *paneHeader {
	h := &paneHeader{tabs: newNativeTabs()}
	h.ExtendBaseWidget(h)
	h.pointer = &paneHeaderPointer{header: h}
	h.pointer.ExtendBaseWidget(h.pointer)
	h.tabs.OnSelected = func(item *container.TabItem) {
		if h.muted {
			return
		}
		if h.disabled {
			h.restore()
			return
		}
		for i, it := range h.tabs.Items {
			if it == item && i < len(h.pages) {
				page := h.pages[i]
				// Undo AppTabs' eager selection before calling the application transaction.
				h.restore()
				if page.enabled && page.id != h.selected && h.request != nil {
					h.request(page.id)
				}
				return
			}
		}
	}
	return h
}
func (h *paneHeader) sync(pages []headerPage, selected string) {
	h.muted = true
	defer func() { h.muted = false }()
	h.selected = selected
	changed := len(h.pages) != len(pages)
	if !changed {
		for i, p := range pages {
			if h.pages[i] != p {
				changed = true
				break
			}
		}
	}
	// Replacing items recreates native header buttons before their next layout.
	// Drafts, provider progress and selection-only changes must keep those hit
	// targets (and their already laid-out inspection geometry) stable.
	if changed {
		h.pages = append([]headerPage(nil), pages...)
		items := make([]*container.TabItem, len(pages))
		for i, p := range pages {
			items[i] = container.NewTabItemWithIcon(p.label, p.icon, container.NewWithoutLayout())
		}
		h.tabs.SetItems(items)
		for i, p := range pages {
			if !p.enabled {
				h.tabs.DisableIndex(i)
			}
		}
	}
	h.restore()
	h.Refresh()
}
func (h *paneHeader) restore() {
	was := h.muted
	h.muted = true
	defer func() { h.muted = was }()
	for i, p := range h.pages {
		if p.id == h.selected {
			h.tabs.SelectIndex(i)
			return
		}
	}
}
func (h *paneHeader) CreateRenderer() fyne.WidgetRenderer {
	border := canvas.NewRectangle(color.Transparent)
	border.StrokeWidth = 2
	return &headerRenderer{h: h, border: border}
}
func (h *paneHeader) FocusGained() {
	h.focused = true
	if !h.disabled && h.focus != nil {
		h.focus()
	}
	h.Refresh()
}
func (h *paneHeader) FocusLost()     { h.shift = false; h.focused = false; h.Refresh() }
func (h *paneHeader) TypedRune(rune) {}
func (h *paneHeader) TypedKey(e *fyne.KeyEvent) {
	if h.disabled || h.muted {
		return
	}
	if commandKeyEvent(e, h.shift, false, h.shortcut) {
		return
	}
	if e.Name == fyne.KeyEscape && h.escape != nil {
		h.escape()
		return
	}
	if e.Name == fyne.KeyTab {
		if h.shift {
			if h.backward != nil {
				h.backward()
			}
		} else if h.enter != nil {
			h.enter()
		}
		return
	}
	eligible := []string{}
	at := -1
	for _, p := range h.pages {
		if p.enabled {
			if p.id == h.selected {
				at = len(eligible)
			}
			eligible = append(eligible, p.id)
		}
	}
	if len(eligible) == 0 {
		return
	}
	next := at
	switch e.Name {
	case fyne.KeyLeft, fyne.KeyUp:
		next--
	case fyne.KeyRight, fyne.KeyDown:
		next++
	case fyne.KeyHome:
		next = 0
	case fyne.KeyEnd:
		next = len(eligible) - 1
	default:
		return
	}
	if next < 0 {
		next = 0
	}
	if next >= len(eligible) {
		next = len(eligible) - 1
	}
	if eligible[next] != h.selected && h.request != nil {
		h.request(eligible[next])
	}
}
func (h *paneHeader) Disabled() bool { return h.disabled }
func (h *paneHeader) Disable()       { h.disabled = true; h.Refresh() }
func (h *paneHeader) Enable()        { h.disabled = false; h.Refresh() }

func (h *paneHeader) AcceptsTab() bool { return true }
func (h *paneHeader) KeyDown(e *fyne.KeyEvent) {
	if e.Name == desktop.KeyShiftLeft || e.Name == desktop.KeyShiftRight {
		h.shift = true
	}
}
func (h *paneHeader) KeyUp(e *fyne.KeyEvent) {
	if e.Name == desktop.KeyShiftLeft || e.Name == desktop.KeyShiftRight {
		h.shift = false
	}
}

// Keep the renderer returned for this actual AppTabs instance, so diagnostics
// read its real public CanvasObject positions rather than estimating hit areas.
type nativeTabs struct {
	container.AppTabs
	renderer fyne.WidgetRenderer
}

func newNativeTabs() *nativeTabs { n := &nativeTabs{}; n.ExtendBaseWidget(n); return n }
func (n *nativeTabs) CreateRenderer() fyne.WidgetRenderer {
	n.renderer = n.AppTabs.CreateRenderer()
	return n.renderer
}
func (n *nativeTabs) buttons() []fyne.CanvasObject {
	if n.renderer == nil {
		return nil
	}
	objects := n.renderer.Objects()
	if len(objects) == 0 {
		return nil
	}
	bar, ok := objects[0].(*fyne.Container)
	if !ok || len(bar.Objects) == 0 {
		return nil
	}
	group, ok := bar.Objects[0].(*fyne.Container)
	if !ok {
		return nil
	}
	return group.Objects
}

type headerRenderer struct {
	h      *paneHeader
	border *canvas.Rectangle
}

func (r *headerRenderer) Layout(s fyne.Size) {
	r.h.tabs.Resize(s)
	r.border.Resize(s)
	r.h.pointer.Resize(s)
}
func (r *headerRenderer) MinSize() fyne.Size { return r.h.tabs.MinSize() }
func (r *headerRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.h.tabs, r.border, r.h.pointer}
}
func (r *headerRenderer) Destroy() {}
func (r *headerRenderer) Refresh() {
	r.border.StrokeColor = color.Transparent
	if r.h.focused {
		r.border.StrokeColor = color.NRGBA{30, 100, 220, 255}
	}
	r.h.tabs.Refresh()
	if native := r.h.tabs.renderer; native != nil {
		objects := native.Objects()
		if len(objects) > 2 {
			if r.h.selected == "" {
				objects[2].Hide()
			} else {
				objects[2].Show()
			}
		}
	}
	r.border.Refresh()
}

// AppTabs suppresses OnSelected for a click on its selected item. A bounded
// pointer surface uses the actual native buttons' hit rectangles to focus that
// header without an activation. Changed-page clicks still use AppTabs' callback
// path; focus is published only if that transaction succeeds. This surface is
// deliberately not Focusable, so it cannot eagerly focus a rejected target.
type paneHeaderPointer struct {
	widget.BaseWidget
	header *paneHeader
}

func (p *paneHeaderPointer) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(canvas.NewRectangle(color.Transparent))
}
func (p *paneHeaderPointer) Tapped(e *fyne.PointEvent) {
	h := p.header
	if h.muted || h.disabled {
		return
	}
	driver := fyne.CurrentApp().Driver()
	origin := driver.AbsolutePositionForObject(p)
	point := origin.Add(e.Position)
	for i, button := range h.tabs.buttons() {
		if i >= len(h.pages) {
			break
		}
		pos, size := driver.AbsolutePositionForObject(button), button.Size()
		if point.X < pos.X || point.Y < pos.Y || point.X >= pos.X+size.Width || point.Y >= pos.Y+size.Height {
			continue
		}
		page := h.pages[i]
		if !page.enabled || page.id == "" {
			return
		}
		if page.id == h.selected {
			if h.focus != nil {
				h.focus()
			}
		} else if target, ok := button.(fyne.Tappable); ok {
			target.Tapped(&fyne.PointEvent{Position: point.Subtract(pos), AbsolutePosition: point})
		}
		return
	}
}
