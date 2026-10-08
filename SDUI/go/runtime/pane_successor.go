package runtime

import (
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"strings"
)

func (n *Session) inheritPanes(s *Session) {
	n.root.Walk(func(node *parser.Instance) {
		w := n.panes[public(node.Path)]
		if w == nil || w.InstancePath != node.Path {
			return
		}
		old := s.panes[w.Handle.Path]
		if old != nil && old.Handle.Kind == w.Handle.Kind && !strings.HasPrefix(w.Handle.Path, "@") {
			w.Handle = old.Handle
		} else {
			n.generation++
			w.Handle.Generation = n.generation
		}
	})
	for path, t := range n.tabs {
		t.Handle = n.panes[path].Handle
		old := s.tabs[path]
		compatible := old != nil && old.Handle == t.Handle
		if compatible {
			t.Selected = old.Selected
		}
		for i := range t.Pages {
			p := &t.Pages[i]
			p.Handle = n.panes[public(p.InstancePath)].Handle
			if !compatible {
				continue
			}
			for _, before := range old.Pages {
				if before.Handle != p.Handle {
					continue
				}
				a, b := s.currentControl(before.RememberedFocus), n.currentControl(before.RememberedFocus)
				if a != nil && b != nil && a.Handle == b.Handle && within(b.InstancePath, p.InstancePath) {
					p.RememberedFocus = before.RememberedFocus
				}
			}
		}
	}
	for path, p := range n.splits {
		p.Handle = n.panes[path].Handle
		old := s.splits[path]
		if old == nil || old.Handle != p.Handle || old.Axis != p.Axis {
			continue
		}
		p.Proportion, p.SavedProportion, p.Collapsed = old.Proportion, old.SavedProportion, old.Collapsed
		if !p.Collapsible && p.Collapsed != SplitNone {
			p.Collapsed = SplitNone
			p.Proportion = p.SavedProportion
		}
	}
	n.refreshActivity()
}
func (n *Session) inheritPaneFocus(s *Session) {
	old := s.currentControl(s.focused)
	if old == nil {
		return
	}
	if w := n.currentControl(s.focused); w != nil && w.Handle == old.Handle && w.Enabled && w.Visible {
		n.focused = s.focused
		return
	}
	longest := 0
	for _, w := range n.panes {
		if w.Handle.Kind != "page" && w.Enabled && w.Visible && within(old.InstancePath, w.InstancePath) && len(w.InstancePath) > longest {
			n.focused = w.Handle.Path
			longest = len(w.InstancePath)
		}
	}
}
