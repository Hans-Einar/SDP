package runtime

import (
	"fmt"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
)

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

	if s.check != nil {
		candidate := *s
		candidate.widgets = map[string]*Widget{}
		for path, w := range s.widgets {
			candidate.widgets[path] = w
		}
		for path, w := range next {
			candidate.widgets[path] = w
		}
		if err := s.check(candidate.SnapshotRoot()); err != nil {
			return err
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
