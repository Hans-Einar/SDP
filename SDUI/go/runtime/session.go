package runtime

import (
	"fmt"
	"sort"

	"github.com/Hans-Einar/SDP/SDUI/go/parser"
)

type Session struct {
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
func (s *Session) Apply(revision, batch uint64, updates []Update) error {
	if s.closed {
		return fault("closed", "Session is closed")
	}
	if revision != s.Revision || batch != s.BatchRevision+1 {
		return fault("stale-batch", "Expected current model and next batch revision")
	}
	if len(updates) > 256 {
		return fault("update-limit", "Batch exceeds 256 updates")
	}
	next := map[string]*Widget{}
	seen := map[string]bool{}
	for _, u := range updates {
		original, e := s.lookup(u.Handle)
		if e != nil {
			return e
		}
		key := u.Handle.Path + "/" + string(u.Property)
		if seen[key] {
			return fault("duplicate-property", key)
		}
		seen[key] = true
		w := next[u.Handle.Path]
		if w == nil {
			copy := *original
			w = &copy
			next[u.Handle.Path] = w
		}
		switch u.Property {
		case Label:
			if !validValue(u.Value, String) {
				return fault("property-type", "Label requires string")
			}
			w.Label = u.Value.Text
		case AcceptedValue:
			if w.Handle.Kind != "input" || !validValue(u.Value, String) {
				return fault("property-type", "Value requires string input")
			}
			if u.ExpectedValueRevision != w.ValueRevision {
				return fault("value-conflict", w.Handle.Path)
			}
			if w.Dirty && !u.AcceptDraft {
				return fault("draft-conflict", w.Handle.Path)
			}
			if u.AcceptDraft && u.Value.Text != w.Draft {
				return fault("draft-conflict", "Accepted value differs from current draft")
			}
			w.Value = u.Value.Text
			w.Draft = w.Value
			w.Dirty = false
			w.ValueRevision++
			w.DraftRevision++
		case Enabled, Visible:
			if !validValue(u.Value, Boolean) {
				return fault("property-type", "Boolean property required")
			}
			if u.Value.Bool && (u.Property == Enabled && !w.ancestorEnabled || u.Property == Visible && !w.ancestorVisible) {
				return fault("inactive-ancestor", w.Handle.Path)
			}
			if u.Property == Enabled {
				w.Enabled = u.Value.Bool
			} else {
				w.Visible = u.Value.Bool
			}
		default:
			return fault("property", fmt.Sprint(u.Property))
		}
	}
	for path, w := range next {
		s.widgets[path] = w
		if s.focused == path && (!w.Enabled || !w.Visible) {
			s.focused = ""
		}
	}
	s.BatchRevision = batch
	return nil
}
func (s *Session) SnapshotRoot() *parser.Instance {
	root := clone(s.root)
	root.Walk(func(n *parser.Instance) {
		if n.Kind != "widget" {
			return
		}
		w := s.widgets[public(n.Path)]
		if w == nil {
			return
		}
		n.Layout["enabled"] = w.Enabled
		n.Layout["visible"] = w.Visible
		key := "label"
		if w.Handle.Kind == "input" {
			key = "text"
			n.Arguments["value"] = parser.Literal{Kind: "string", Value: w.Draft, Span: n.Span}
		}
		n.Arguments[key] = parser.Literal{Kind: "string", Value: w.Label, Span: n.Span}
	})
	return root
}
func (s *Session) Close() {
	s.closed = true
	s.handlers = map[string]Handler{}
	s.focused = ""
	s.Revision++
}
