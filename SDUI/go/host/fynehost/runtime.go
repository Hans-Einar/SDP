package fynehost

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/markdown"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/reload"
	uiruntime "github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

// RuntimeView is the owner-goroutine adapter. Watcher candidates must arrive
// through fyne.Do; native OnChanged edits drafts, OnSubmitted commits them.
type RuntimeView struct {
	Session       *uiruntime.Session
	View          *View
	Container     *fyne.Container
	Canvas        fyne.Canvas
	Controller    reload.Controller
	OnStatus      func(error)
	sequence      uint64
	muted, closed bool
}

func NewRuntime(session *uiruntime.Session, canvas fyne.Canvas) (*RuntimeView, error) {
	r := &RuntimeView{Session: session, Canvas: canvas, Controller: reload.Controller{Session: session}}
	r.Container = container.NewStack()
	if err := r.mount(); err != nil {
		return nil, err
	}
	if err := session.CheckWith(r.check); err != nil {
		return nil, err
	}
	return r, nil
}
func (r *RuntimeView) check(root *parser.Instance) error {
	provider, err := markdown.Prepare(root, nil)
	if err != nil {
		return err
	}
	size := r.Container.Size()
	if size.Width < 2 || size.Height < 2 {
		return nil
	}
	_, err = (&layout.Engine{Measure: provider}).Layout(root, layout.Size{W: float64(size.Width), H: float64(size.Height)})
	return err
}
func (r *RuntimeView) status(err error) {
	if r.OnStatus != nil {
		r.OnStatus(err)
	}
}
func (r *RuntimeView) mount() error {
	root := r.Session.SnapshotRoot()
	provider, err := markdown.Prepare(root, nil)
	if err != nil {
		return err
	}
	old := r.View
	next := New(root)
	next.Measure = provider
	next.Content = provider
	next.OnStatus = r.status
	revision := r.Session.Revision
	for _, state := range r.Session.Widgets() {
		w := state
		switch obj := next.Controls[w.InstancePath].(type) {
		case *widget.Button:
			next.Actions[w.InstancePath] = func(_, _ string) error { return r.dispatch(w.Handle, revision, uiruntime.Activate) }
		case *Input:
			next.Actions[w.InstancePath] = func(_, value string) error {
				if r.muted || r.closed {
					return nil
				}
				return r.Session.Draft(w.Handle, value)
			}
			obj.OnFocus = func() {
				if !r.closed && r.View == next {
					_ = r.Session.Focus(w.Handle)
				}
			}
			obj.OnSubmitted = func(string) { r.status(r.dispatch(w.Handle, revision, uiruntime.Commit)) }
			obj.OnRevert = func() {
				if r.closed || r.View != next {
					return
				}
				err := r.Session.Revert(w.Handle)
				if err == nil {
					err = r.Sync()
				}
				r.status(err)
			}
		}
	}
	r.View = next
	r.Container.Objects = []fyne.CanvasObject{next.Container}
	r.Container.Refresh()
	if old != nil {
		old.Close()
	}
	if focused, ok := r.Session.Widget(r.Session.Focused()); ok {
		if obj, ok := next.Controls[focused.InstancePath].(fyne.Focusable); ok {
			r.Canvas.Focus(obj)
		}
	}
	return nil
}
func (r *RuntimeView) dispatch(handle uiruntime.Handle, revision uint64, kind uiruntime.EventKind) error {
	if r.closed {
		return fmt.Errorf("closed: runtime view")
	}
	if revision != r.Session.Revision {
		return fmt.Errorf("stale-event: native callback belongs to another model")
	}
	r.sequence++
	w, ok := r.Session.Widget(handle.Path)
	if !ok {
		return fmt.Errorf("stale-handle: %s", handle.Path)
	}
	event := uiruntime.Event{Handle: handle, ModelRevision: r.Session.Revision, Sequence: r.sequence, Kind: kind}
	if kind == uiruntime.Commit {
		event.Value = uiruntime.Text(w.Draft)
		event.DraftRevision = w.DraftRevision
	}
	err := r.Session.Dispatch(event)
	if err != nil {
		return err
	}
	return r.Sync()
}
func (r *RuntimeView) Sync() error {
	if r.closed {
		return fmt.Errorf("closed: runtime view")
	}
	r.muted = true
	defer func() { r.muted = false }()
	root := r.Session.SnapshotRoot()
	provider, err := markdown.Prepare(root, nil)
	if err != nil {
		return err
	}
	r.View.Root = root
	r.View.Measure = provider
	r.View.Content = provider
	for _, w := range r.Session.Widgets() {
		switch obj := r.View.Controls[w.InstancePath].(type) {
		case *widget.Button:
			obj.SetText(w.Label)
		case *Input:
			obj.SetPlaceHolder(w.Label)
			if obj.Text != w.Draft {
				obj.SetText(w.Draft)
			}
		}
	}
	r.View.Container.Refresh()
	return nil
}
func (r *RuntimeView) Adopt(candidate reload.Candidate) error {
	if r.closed {
		return fmt.Errorf("closed: runtime view")
	}
	for _, w := range r.Session.Widgets() {
		if obj := r.View.Controls[w.InstancePath]; obj != nil && any(obj) == any(r.Canvas.Focused()) {
			_ = r.Session.Focus(w.Handle)
		}
	}
	if err := r.Controller.Adopt(candidate); err != nil {
		r.status(err)
		return err
	}
	err := r.mount()
	r.status(err)
	return err
}
func (r *RuntimeView) Close() { r.closed = true; r.View.Close(); r.Session.Close() }
