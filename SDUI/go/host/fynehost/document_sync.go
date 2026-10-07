package fynehost

import (
	"context"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"sort"
)

func (b *Bundle) connect() {
	for _, w := range b.Session.Widgets() {
		w := w
		live := func() bool { return b.owner.current == b && !b.closed && !b.owner.closed && !b.muted }
		switch control := b.view.Controls[w.InstancePath].(type) {
		case *widget.Button:
			b.view.Actions[w.InstancePath] = func(_, _ string) error {
				if !live() {
					return nil
				}
				return b.owner.Mutate(func(s *ui.Session) error {
					return s.Dispatch(ui.Event{Handle: w.Handle, ModelRevision: b.Session.Revision, Sequence: s.Sequence() + 1, Kind: ui.Activate})
				})
			}
		case *Input:
			b.view.Actions[w.InstancePath] = func(_, value string) error {
				if !live() {
					return nil
				}
				return b.owner.Mutate(func(s *ui.Session) error { return s.Draft(w.Handle, value) })
			}
			control.OnFocus = func() {
				if live() {
					if err := b.Session.Focus(w.Handle); err != nil {
						b.owner.status(err)
						return
					}
					b.owner.ensureWidget(w.InstancePath)
				}
			}
			control.OnRevert = func() {
				if live() {
					b.owner.status(b.owner.Mutate(func(s *ui.Session) error { return s.Revert(w.Handle) }))
				}
			}
			control.OnSubmitted = func(string) {
				if !live() {
					return
				}
				b.owner.status(b.owner.Mutate(func(s *ui.Session) error {
					current, ok := s.Widget(w.Handle.Path)
					if !ok {
						return fmt.Errorf("stale input")
					}
					return s.Dispatch(ui.Event{Handle: w.Handle, ModelRevision: s.Revision, Sequence: s.Sequence() + 1, Kind: ui.Commit, Value: ui.Text(current.Draft), DraftRevision: current.DraftRevision})
				}))
			}
		}
	}
	b.view.OnStatus = b.owner.status
	for _, control := range b.view.controls {
		control := control
		control.clip.onScroll = func(e *fyne.ScrollEvent) {
			if b.owner.current != b || b.closed {
				return
			}
			b.owner.scroll(float64(control.clip.Position().X+e.Position.X), float64(control.clip.Position().Y+e.Position.Y), -float64(e.Scrolled.DX), -float64(e.Scrolled.DY))
		}
	}
}
func (b *Bundle) apply() {
	p := b.pending
	if p == nil {
		return
	}
	b.presentation = p
	b.pending = nil
	b.muted = true
	defer func() { b.muted = false }()
	v := b.view
	v.Container.Objects = p.objects
	v.Root = p.snapshot.Root
	v.Geometry = p.geometry.Root
	v.image.Resource = p.background
	v.image.Resize(fyne.NewSize(float32(b.size.W), float32(b.size.H)))
	v.image.Refresh()
	active := map[string]bool{}
	p.geometry.Root.Walk(func(box *layout.Box) {
		if box.Rect.W > 0 && box.Rect.H > 0 && box.Clip.W > 0 && box.Clip.H > 0 {
			active[box.Path] = true
		}
	})
	for path, c := range v.controls {
		if !active[path] {
			c.clip.Hide()
		}
	}
	p.geometry.Root.Walk(func(box *layout.Box) {
		c := v.controls[box.Path]
		if c == nil {
			return
		}
		rect, clip := box.Rect, box.Clip
		if rect.W <= 0 || rect.H <= 0 || clip.W <= 0 || clip.H <= 0 {
			return
		}
		switch obj := c.widget.(type) {
		case *CollectionControl:
			obj.state = p.snapshot.Collections[box.Path]
			obj.rows = collectionRows(obj.state, box.Font)
			obj.font = box.Font
			obj.title = box.Instance.Argument("label")
			if viewport, ok := p.geometry.Viewports[box.Path]; ok {
				obj.viewport = viewport
			} else {
				height := rowHeight(box.Font)
				obj.viewport = layout.Viewport{Path: box.Path, Rect: layout.Rect{X: rect.X, Y: rect.Y + height, W: rect.W, H: rect.H - height}, Clip: clip}
			}
			obj.Refresh()
		case *Input:
			obj.SetPlaceHolder(box.Instance.Argument("text"))
			value := box.Instance.Argument("value")
			if obj.Text != value {
				obj.SetText(value)
			}
		case *widget.Button:
			obj.SetText(box.Instance.Argument("label"))
		}
		c.theme.Theme = componentTheme{float32(box.Font)}
		c.theme.Refresh()
		c.fixed.size = fyne.NewSize(float32(rect.W), float32(rect.H))
		c.clip.Content.Resize(c.fixed.size)
		c.clip.Resize(fyne.NewSize(float32(clip.W), float32(clip.H)))
		c.clip.Move(fyne.NewPos(float32(clip.X), float32(clip.Y)))
		c.clip.Offset = fyne.NewPos(float32(clip.X-rect.X), float32(clip.Y-rect.Y))
		c.clip.Refresh()
		c.clip.Show()
		if d, ok := c.widget.(fyne.Disableable); ok {
			if box.Enabled {
				d.Enable()
			} else {
				d.Disable()
			}
		}
	})
	v.Container.Resize(fyne.NewSize(float32(b.size.W), float32(b.size.H)))
}
func (h *DocumentHost) after() {
	b := h.current
	if b == nil || b.closed || h.closed {
		return
	}
	snapshot := b.Session.Snapshot()
	paths := make([]string, 0, len(snapshot.Collections))
	for p := range snapshot.Collections {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	for _, path := range paths {
		state := snapshot.Collections[path]
		w, ok := b.Session.Widget(state.Handle.Path)
		if !ok || !w.Enabled || !w.Visible {
			continue
		}
		if state.Request == nil && state.AutoLoadPending {
			target, e := b.Session.Target(state.Handle, "")
			if e == nil {
				_, e = b.Session.BeginLoad(target)
			}
			h.status(e)
		}
	}

	b.apply()
	h.reconcileLoads(b)
	if h.OnChange != nil {
		h.OnChange(b)
	}
}
func (h *DocumentHost) reconcileLoads(b *Bundle) {
	snapshot := b.Session.Snapshot()
	for path, flight := range b.flights {
		state, ok := snapshot.Collections[path]
		if !ok || state.Request == nil || *state.Request != flight.request {
			flight.cancel()
			delete(b.flights, path)
		}
	}
	for path, state := range snapshot.Collections {
		if state.Request == nil {
			continue
		}
		request := *state.Request
		if _, ok := b.flights[path]; ok {
			continue
		}
		provider, ok := b.Session.Provider(state.Handle)
		if !ok || provider.Load == nil {
			continue
		}
		ctx, cancel := context.WithCancel(context.Background())
		b.flights[path] = providerFlight{request, cancel}
		post := h.Post
		go func(path string, request ui.LoadRequest) {
			var data ui.CollectionData
			var err error
			func() {
				defer func() {
					if v := recover(); v != nil {
						err = fmt.Errorf("provider panic: %v", v)
					}
				}()
				data, err = provider.Load(ctx, request)
			}()
			if len(data.Items) > 4096 {
				data = ui.CollectionData{}
				err = fmt.Errorf("collection-limit: more than 4096 items")
			} else {
				data.Items = append([]ui.CollectionItem(nil), data.Items...)
			}
			post(func() {
				if h.closed || b.closed || h.current != b {
					return
				}
				flight, ok := b.flights[path]
				if !ok || flight.request != request {
					return
				}
				if e := b.Session.CompleteLoad(request, data, err); e != nil {
					h.status(e)
					return
				}
				h.after()
			})
		}(path, request)
	}
}
