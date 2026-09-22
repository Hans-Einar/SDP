package runtime

import (
	"sort"

	"github.com/Hans-Einar/SDP/SDUI/go/parser"
)

type Session struct {
	check                   func(*parser.Instance) error
	ID                      string
	Revision, BatchRevision uint64
	root                    *parser.Instance
	widgets                 map[string]*Widget
	handlers                map[string]Handler
	focused                 string
	generation, sequence    uint64
	closed                  bool
}

func New(id string, root *parser.Instance) (*Session, error) {
	if id == "" || root == nil || root.Kind != "frame" {
		return nil, fault("session", "Session ID and frame root are required")
	}
	s := &Session{ID: id, Revision: 1, root: clone(root), widgets: map[string]*Widget{}, handlers: map[string]Handler{}}
	if err := s.register(s.root, true, true); err != nil {
		return nil, err
	}
	return s, nil
}
func (s *Session) register(n *parser.Instance, enabled, visible bool) error {
	parentEnabled, parentVisible := enabled, visible
	enabled = enabled && n.Layout["enabled"] != false
	visible = visible && n.Layout["visible"] != false
	if n.Kind == "widget" {
		path := public(n.Path)
		if s.widgets[path] != nil {
			return fault("ambiguous-instance", path)
		}
		s.generation++
		w := &Widget{ancestorEnabled: parentEnabled, ancestorVisible: parentVisible, Handle: Handle{s.ID, path, s.generation, n.Widget}, InstancePath: n.Path, Label: n.Argument("label"), Value: n.Argument("value"), Draft: n.Argument("value"), ValueRevision: 1, DraftRevision: 1, Enabled: enabled, Visible: visible}
		if n.Widget == "input" {
			w.Label = n.Argument("text")
		}
		if ref, ok := n.Arguments["callback"].(parser.Reference); ok {
			w.Binding = ref
		}
		s.widgets[path] = w
	}
	for _, r := range n.Regions {
		if err := s.register(r.Node, enabled, visible); err != nil {
			return err
		}
	}
	for _, row := range n.Rows {
		for _, c := range row {
			if err := s.register(c, enabled, visible); err != nil {
				return err
			}
		}
	}
	return nil
}
func (s *Session) Widgets() []Widget {
	keys := []string{}
	for k := range s.widgets {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := []Widget{}
	for _, k := range keys {
		out = append(out, *s.widgets[k])
	}
	return out
}
func (s *Session) Widget(path string) (Widget, bool) {
	w, ok := s.widgets[path]
	if !ok {
		return Widget{}, false
	}
	return *w, true
}
func (s *Session) Bind(handle Handle, handler Handler) error {
	if _, e := s.lookup(handle); e != nil {
		return e
	}
	if handler == nil {
		return fault("binding", "Nil callback")
	}
	s.handlers[handle.Path] = handler
	return nil
}
func (s *Session) Draft(handle Handle, value string) error {
	w, e := s.lookup(handle)
	if e != nil {
		return e
	}
	if w.Handle.Kind != "input" || !w.Enabled || !w.Visible {
		return fault("draft", "Input is not editable")
	}
	if len(value) > 32768 {
		return fault("value-limit", "Draft exceeds 32768 bytes")
	}
	w.Draft = value
	w.DraftRevision++
	w.Dirty = w.Draft != w.Value
	return nil
}
func (s *Session) Revert(handle Handle) error {
	w, e := s.lookup(handle)
	if e != nil {
		return e
	}
	if w.Handle.Kind != "input" {
		return fault("widget-type", "Revert requires input")
	}
	w.Draft = w.Value
	w.Dirty = false
	w.DraftRevision++
	return nil
}
func (s *Session) Focus(handle Handle) error {
	w, e := s.lookup(handle)
	if e != nil {
		return e
	}
	if !w.Enabled || !w.Visible {
		return fault("focus", "Widget is not focusable")
	}
	s.focused = handle.Path
	return nil
}
func (s *Session) Focused() string { return s.focused }
func (s *Session) Dispatch(event Event) error {
	w, err := s.lookup(event.Handle)
	if err != nil {
		return err
	}
	if event.ModelRevision != s.Revision {
		return fault("stale-event", "Event belongs to another model revision")
	}
	if event.Sequence == 0 || event.Sequence <= s.sequence {
		return fault("duplicate-event", "Event sequence already consumed")
	}
	if !w.Enabled || !w.Visible {
		return fault("inactive-widget", w.Handle.Path)
	}
	switch event.Kind {
	case Activate:
		if w.Handle.Kind != "button" || event.Value != (Value{}) {
			return fault("event-type", "Activate requires a button and no payload")
		}
	case Commit:
		if w.Handle.Kind != "input" || event.DraftRevision != w.DraftRevision || !validValue(event.Value, String) || event.Value.Text != w.Draft {
			return fault("event-type", "Commit requires the current input draft")
		}
	default:
		return fault("event-type", string(event.Kind))
	}
	s.sequence = event.Sequence
	handler := s.handlers[w.Handle.Path]
	if handler == nil {
		return fault("unbound", w.Handle.Path)
	}
	revision := s.Revision
	updates, err := handler(event)
	if err != nil {
		return err
	}
	if s.closed || s.Revision != revision {
		return fault("stale-result", "Session changed while callback executed")
	}
	return s.Apply(revision, s.BatchRevision+1, updates)
}
func (s *Session) Close() {
	s.closed = true
	s.handlers = map[string]Handler{}
	s.focused = ""
	s.Revision++
}

// CheckWith installs a pure presentation gate. Failed geometry never publishes
// partial state; callbacks are never invoked by the gate.
func (s *Session) CheckWith(check func(*parser.Instance) error) error {
	if check != nil {
		if err := check(s.SnapshotRoot()); err != nil {
			return err
		}
	}
	s.check = check
	return nil
}
