package runtime

import (
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"strings"
)

// Successor creates a detached candidate without handlers or presentation gates.
// It never cancels old requests, executes providers or mutates the live Session.
func (s *Session) Successor(root *parser.Instance, providers map[string]CollectionProvider) (*Session, error) {
	if s.closed {
		return nil, fault("closed", "Session is closed")
	}
	n, err := New(s.ID, root)
	if err != nil {
		return nil, err
	}
	n.Revision = s.Revision + 1
	n.StateRevision = s.StateRevision + 1
	n.BatchRevision = s.BatchRevision
	n.sequence = s.sequence
	n.generation = s.generation
	for _, item := range n.Widgets() {
		w := n.widgets[item.Handle.Path]
		old := s.widgets[item.Handle.Path]
		if old != nil && old.Handle.Kind == w.Handle.Kind && !strings.HasPrefix(w.Handle.Path, "@") {
			w.Handle = old.Handle
			w.Value = old.Value
			w.Draft = old.Draft
			w.Dirty = old.Dirty
			w.ValueRevision = old.ValueRevision
			w.DraftRevision = old.DraftRevision
		} else {
			n.generation++
			w.Handle.Generation = n.generation
		}
	}
	// Source order makes generation allocation deterministic.
	n.root.Walk(func(node *parser.Instance) {
		h, ok := n.viewportHandles[node.Path]
		if !ok {
			return
		}
		old, ok := s.viewportHandles[node.Path]
		if ok && old.Kind == h.Kind && !anonymousPath(node.Path) {
			h = old
			v := s.viewports[node.Path]
			if node.Layout["overflow-x"] != "scroll" {
				v.X = 0
			}
			if node.Layout["overflow-y"] != "scroll" {
				v.Y = 0
			}
			n.viewports[node.Path] = v
		} else {
			n.generation++
			h.Generation = n.generation
		}
		n.viewportHandles[node.Path] = h
	})
	if old := s.widgets[s.focused]; old != nil {
		if w := n.widgets[s.focused]; w != nil && w.Handle == old.Handle && w.Enabled && w.Visible {
			n.focused = s.focused
		}
	}
	if err = n.BindProviders(providers); err != nil {
		return nil, err
	}
	for path, c := range n.collections {
		old := s.collections[path]
		if old == nil {
			continue
		}
		c.requestID = old.requestID
		c.Generation = old.Generation + 1
		if c.Handle != old.Handle || c.ProviderID != old.ProviderID || c.ProviderEpoch != old.ProviderEpoch {
			continue
		}
		p := c.provider
		handle := c.Handle
		instancePath := c.InstancePath
		c = copyCollection(old)
		c.provider = p
		c.Handle = handle
		c.InstancePath = instancePath
		c.Generation++
		c.Request = nil
		for id, status := range c.Status {
			if status.Phase == Loading {
				c.Status[id] = LoadStatus{Phase: Unloaded}
			}
		}
		n.collections[path] = c
	}
	n.StateRevision = s.StateRevision + 1
	return n, nil
}
func anonymousPath(path string) bool {
	return strings.HasPrefix(path[strings.LastIndex(path, "/")+1:], "$")
}
