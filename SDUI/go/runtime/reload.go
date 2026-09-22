package runtime

import (
	"strings"

	"github.com/Hans-Einar/SDP/SDUI/go/parser"
)

// Reload publishes a fully validated replacement. Named widgets of the same
// kind retain value/draft and logical handle. Defaults initialize new instances
// only. Source labels, visibility and enabled rules belong to the new model.
func (s *Session) Reload(root *parser.Instance) error {
	if s.closed {
		return fault("closed", "Session is closed")
	}
	next, err := New(s.ID, root)
	if err != nil {
		return err
	}
	generation := s.generation
	for _, item := range next.Widgets() {
		path := item.Handle.Path
		w := next.widgets[path]
		old := s.widgets[path]
		if old != nil && old.Handle.Kind == w.Handle.Kind && !strings.HasPrefix(path, "@") {
			w.Handle = old.Handle
			w.Value = old.Value
			w.Draft = old.Draft
			w.Dirty = old.Dirty
			w.ValueRevision = old.ValueRevision
			w.DraftRevision = old.DraftRevision
			if old.Binding == w.Binding {
				if handler := s.handlers[path]; handler != nil {
					next.handlers[path] = handler
				}
			}
		} else {
			generation++
			w.Handle.Generation = generation
		}
	}
	focused := s.focused
	if w := next.widgets[focused]; w == nil || !w.Enabled || !w.Visible || s.widgets[focused].Handle != w.Handle {
		focused = ""
	}
	candidate := *s
	candidate.root = next.root
	candidate.widgets = next.widgets
	candidate.handlers = next.handlers
	candidate.generation = generation
	candidate.focused = focused
	candidate.Revision++
	if s.check != nil {
		if err := s.check(candidate.SnapshotRoot()); err != nil {
			return err
		}
	}
	*s = candidate
	return nil
}
