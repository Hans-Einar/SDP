package runtime

import (
	"fmt"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"strings"
	"unicode/utf8"
)

func (s *Session) Apply(revision, batch uint64, updates []Update) error {
	n, err := s.applyCandidate(revision, batch, updates)
	if err != nil {
		return err
	}
	return s.publish(n)
}
func (s *Session) applyCandidate(revision, batch uint64, updates []Update) (*Session, error) {
	if s.closed {
		return nil, fault("closed", "Session is closed")
	}
	if revision != s.Revision || batch != s.BatchRevision+1 {
		return nil, fault("stale-batch", "Expected current model and next batch revision")
	}
	if len(updates) > 256 {
		return nil, fault("update-limit", "Batch exceeds 256 updates")
	}
	next := map[string]*Widget{}
	seen := map[string]bool{}
	for _, u := range updates {
		original, e := s.lookupControl(u.Handle)
		if e != nil {
			return nil, e
		}
		key := u.Handle.Path + "/" + string(u.Property)
		if seen[key] {
			return nil, fault("duplicate-property", key)
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
				return nil, fault("property-type", "Label requires string")
			}
			if s.panes[u.Handle.Path] != nil {
				if u.Handle.Kind == "split" {
					return nil, fault("property", "Split has no label property")
				}
				if !utf8.ValidString(u.Value.Text) || strings.TrimSpace(u.Value.Text) == "" {
					return nil, fault("property-type", "Pane label must be nonempty UTF-8")
				}
			}
			w.Label = u.Value.Text
		case AcceptedValue:
			if w.Handle.Kind != "input" || !validValue(u.Value, String) {
				return nil, fault("property-type", "Value requires string input")
			}
			if u.ExpectedValueRevision != w.ValueRevision {
				return nil, fault("value-conflict", w.Handle.Path)
			}
			if w.Dirty && !u.AcceptDraft {
				return nil, fault("draft-conflict", w.Handle.Path)
			}
			if u.AcceptDraft && u.Value.Text != w.Draft {
				return nil, fault("draft-conflict", "Accepted value differs from current draft")
			}
			w.Value = u.Value.Text
			w.Draft = w.Value
			w.Dirty = false
			w.ValueRevision++
			w.DraftRevision++
		case Enabled, Visible:
			if !validValue(u.Value, Boolean) {
				return nil, fault("property-type", "Boolean property required")
			}
			if len(s.panes) == 0 && u.Value.Bool && (u.Property == Enabled && !w.ancestorEnabled || u.Property == Visible && !w.ancestorVisible) {
				return nil, fault("inactive-ancestor", w.Handle.Path)
			}
			if u.Property == Enabled {
				w.Enabled = u.Value.Bool
			} else {
				w.Visible = u.Value.Bool
			}
		default:
			return nil, fault("property", fmt.Sprint(u.Property))
		}
	}

	candidate := s.copyState()
	for path, w := range next {
		if candidate.panes[path] != nil {
			candidate.panes[path] = w
		} else {
			candidate.widgets[path] = w
		}
		if len(candidate.panes) > 0 {
			a := candidate.intent[w.InstancePath]
			if seen[path+"/"+string(Enabled)] {
				a.enabled = w.Enabled
			}
			if seen[path+"/"+string(Visible)] {
				a.visible = w.Visible
			}
			candidate.intent[w.InstancePath] = a
		}
		if !w.Enabled || !w.Visible {
			if len(candidate.panes) == 0 && candidate.focused == path {
				candidate.focused = ""
			}
			if c := candidate.collections[path]; c != nil {
				c.cancel()
			}
		}
	}
	candidate.BatchRevision = batch
	return candidate, nil
}

func (s *Session) SnapshotRoot() *parser.Instance {
	root := clone(s.root)
	root.Walk(func(n *parser.Instance) {
		if len(s.panes) > 0 {
			if a, ok := s.active[n.Path]; ok {
				n.Layout["enabled"] = a.enabled
				n.Layout["visible"] = a.visible
			}
		}
		if n.Kind != "widget" && n.Kind != "composition" {
			return
		}
		w := s.currentControl(public(n.Path))
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
		if w.Handle.Kind != "split" {
			n.Arguments[key] = parser.Literal{Kind: "string", Value: w.Label, Span: n.Span}
		}
	})
	return root
}
