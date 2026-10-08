package runtime

import "github.com/Hans-Einar/SDP/SDUI/go/parser"

// Reload retains the legacy synchronous API. Connected hosts use Successor,
// bind/prepare it, then publish atomically instead of live reload.
func (s *Session) Reload(root *parser.Instance) error {
	if s.closed {
		return fault("closed", "Session is closed")
	}
	if _, err := parser.EffectiveProfile(root); err != nil {
		return err
	}
	baseline := s.StateRevision
	providers := map[string]CollectionProvider{}
	rootPaths := map[string]bool{}
	if root != nil {
		root.Walk(func(n *parser.Instance) {
			if n.Kind == "widget" && isCollection(n.Widget) {
				rootPaths[n.Path] = true
			}
		})
	}
	for _, c := range s.collections {
		if rootPaths[c.InstancePath] {
			providers[c.InstancePath] = c.provider
		}
	}
	n, err := s.Successor(root, providers)
	if err != nil {
		return err
	}
	for path, w := range n.widgets {
		if old := s.widgets[path]; old != nil && old.Handle == w.Handle && old.Binding == w.Binding {
			n.handlers[path] = s.handlers[path]
		}
	}
	for path, w := range n.panes {
		if old := s.panes[path]; old != nil && old.Handle == w.Handle && old.Binding == w.Binding {
			n.interactions[path] = s.interactions[path]
		}
	}
	for path, c := range n.commands {
		if old := s.commands[path]; old != nil && old.Handle == c.Handle && old.Binding == c.Binding {
			n.interactions[path] = s.interactions[path]
		}
	}
	for path := range n.surfaces {
		if old := s.aux[path]; old != nil && old.Handle == n.aux[path].Handle && old.Binding == n.aux[path].Binding {
			n.interactions[path] = s.interactions[path]
		}
	}
	old := s.copyState()
	_ = old.RevokeSurfaces(Handle{}, "reload")
	n.dialogResults = copyResults(old.dialogResults)
	for path := range n.fields {
		old := s.widgets[path]
		current := n.widgets[path]
		if old != nil && old.Handle == current.Handle {
			n.validators[path] = s.validators[path]
			n.changes[path] = s.changes[path]
			if err := n.validateField(path); err != nil {
				return err
			}
		}
	}
	if s.closed || s.StateRevision != baseline {
		return fault("stale-validation", "State changed while preparing reload")
	}
	n.presentationCheck = s.presentationCheck
	n.presentationPrepare = s.presentationPrepare
	n.check = s.check
	n.stateCheck = s.stateCheck
	return s.publish(n)
}
