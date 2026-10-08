package fynehost

import (
	"fmt"
	"fyne.io/fyne/v2"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

// InteractionError preserves the domain outcome when a native pane interaction
// fails to publish its UI state. The wrapped cause remains available to Is/As.
type InteractionError struct {
	Result ui.InteractionResult
	Err    error
}

func (e *InteractionError) Error() string {
	return fmt.Sprintf("interaction status=%s domain=%s sequence=%d: %v", e.Result.Status, e.Result.Domain, e.Result.Sequence, e.Err)
}
func (e *InteractionError) Unwrap() error { return e.Err }

func interactionError(result ui.InteractionResult, err error) error {
	if err == nil {
		return nil
	}
	return &InteractionError{Result: result, Err: err}
}

func (b *Bundle) livePane() bool {
	return b.owner.current == b && !b.closed && !b.owner.closed && !b.muted && b.presentation != nil && b.Session.Revision == b.presentation.snapshot.ModelRevision
}
func (b *Bundle) connectPanes() {
	for path, obj := range b.view.Controls {
		path := path
		switch c := obj.(type) {
		case *paneHeader:
			c.request = func(id string) {
				if !b.livePane() {
					return
				}
				rendered, ok := b.presentation.snapshot.Tabs[path]
				if !ok {
					return
				}
				current, ok := b.Session.Tabs(rendered.Handle)
				if !ok || current.Selected != rendered.Selected {
					return
				}
				for _, page := range rendered.Pages {
					if page.ID == id {
						b.owner.Mutate(func(s *ui.Session) error {
							result, err := s.DispatchInteraction(ui.Event{Handle: rendered.Handle, ModelRevision: b.presentation.snapshot.ModelRevision, StateRevision: s.StateRevision, Sequence: s.Sequence() + 1, Kind: ui.ActivatePage, Page: &ui.PageActivation{PreviousID: rendered.Selected, PageID: id, Page: page.Handle}})
							return interactionError(result, err)
						})
						b.owner.restoreFocus(b)
						return
					}
				}
			}
			c.focus = func() {
				if !b.livePane() {
					return
				}
				t := b.presentation.snapshot.Tabs[path]
				b.owner.status(b.owner.Mutate(func(s *ui.Session) error { return s.Focus(t.Handle) }))
			}
			c.enter = func() {
				if !b.livePane() {
					return
				}
				t := b.presentation.snapshot.Tabs[path]
				b.owner.status(b.owner.Mutate(func(s *ui.Session) error { return s.EnterPage(t.Handle) }))
				b.owner.restoreFocus(b)
			}
			c.backward = func() {
				if b.livePane() && b.owner.canvas != nil {
					b.canvasFor(path).FocusPrevious()
				}
			}
		case *paneDivider:
			c.request = func(operation string, proportion float64) {
				if !b.livePane() {
					return
				}
				rendered, ok := b.presentation.snapshot.Splits[path]
				if !ok {
					return
				}
				current, ok := b.Session.Split(rendered.Handle)
				if !ok || current.Collapsed != rendered.Collapsed || current.Proportion != rendered.Proportion {
					return
				}
				b.owner.Mutate(func(s *ui.Session) error {
					result, err := s.DispatchInteraction(ui.Event{Handle: rendered.Handle, ModelRevision: b.presentation.snapshot.ModelRevision, StateRevision: s.StateRevision, Sequence: s.Sequence() + 1, Kind: ui.AdjustSplit, Split: &ui.SplitChange{Operation: operation, Proportion: proportion}})
					return interactionError(result, err)
				})
				b.owner.restoreFocus(b)
			}
			c.focus = func() {
				if !b.livePane() {
					return
				}
				v := b.presentation.snapshot.Splits[path]
				b.owner.status(b.owner.Mutate(func(s *ui.Session) error { return s.Focus(v.Handle) }))
			}
			c.focusCanvas = func() {
				if b.livePane() && b.owner.canvas != nil {
					b.canvasFor(path).Focus(c)
				}
			}
		}
	}
}

// Keep a checked native pane operation scoped to the displayed bundle.
func (b *Bundle) paneHandle(path string) (ui.Handle, error) {
	if !b.livePane() {
		return ui.Handle{}, fmt.Errorf("stale pane bundle")
	}
	if t, ok := b.presentation.snapshot.Tabs[path]; ok {
		return t.Handle, nil
	}
	if t, ok := b.presentation.snapshot.Splits[path]; ok {
		return t.Handle, nil
	}
	return ui.Handle{}, fmt.Errorf("unknown pane %s", path)
}

var _ fyne.Focusable = (*paneHeader)(nil)
var _ fyne.Focusable = (*paneDivider)(nil)
