package runtime

import "github.com/Hans-Einar/SDP/SDUI/go/parser"

func (n *Session) inheritCommands(s *Session) {
	n.surfaceEpoch = s.surfaceEpoch
	for _, w := range n.widgets {
		if n.ownerSurface[w.InstancePath] != "" && w.Handle.Kind == "input" && w.Dirty {
			w.Draft = w.Value
			w.Dirty = false
			w.DraftRevision++
		}
	}
	n.root.Walk(func(node *parser.Instance) {
		w := n.aux[public(node.Path)]
		if w == nil {
			return
		}
		old := s.aux[w.Handle.Path]
		if old != nil && old.Handle.Kind == w.Handle.Kind && !anonymousPath(w.InstancePath) {
			w.Handle = old.Handle
		} else {
			n.generation++
			w.Handle.Generation = n.generation
		}
	})
	if n.rootOwner != (Handle{}) {
		n.generation++
		n.rootOwner.Generation = n.generation
	}
	for path, c := range n.commands {
		c.Handle = n.currentControl(path).Handle
		if c.Target != (Handle{}) {
			c.Target = n.currentControl(c.Target.Path).Handle
		}
		old := s.commands[path]
		if old != nil && old.Handle == c.Handle && old.Toggle == c.Toggle && old.Exclusive == c.Exclusive && old.ExclusiveScope == c.ExclusiveScope && old.Context == c.Context && old.Effect == c.Effect && old.Target.Path == c.Target.Path && old.Binding == c.Binding {
			c.Checked = old.Checked
		}
	}
	// A source group change may collide with a retained choice; use the validated
	// declaration for the entire conflicting group, never pick a map-order winner.
	groups := map[string]int{}
	for _, c := range n.commands {
		if c.Checked && c.Exclusive != "" {
			groups[c.ExclusiveScope+"\x00"+c.Exclusive]++
		}
	}
	n.root.Walk(func(node *parser.Instance) {
		if c := n.commands[public(node.Path)]; c != nil && groups[c.ExclusiveScope+"\x00"+c.Exclusive] > 1 {
			c.Checked = boolArgument(node, "checked", false)
		}
	})
	for path, p := range n.presentations {
		p.Handle = n.currentControl(path).Handle
		p.Command = n.commands[p.Command.Path].Handle
	}
	for path, h := range n.menuTargets {
		n.menuTargets[path] = n.currentControl(h.Path).Handle
	}
	for path, m := range n.menus {
		m.Handle = n.aux[path].Handle
	}
	for path, d := range n.surfaces {
		d.Handle = n.aux[path].Handle
	}
}
