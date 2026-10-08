package runtime

import "github.com/Hans-Einar/SDP/SDUI/go/parser"

// refreshActivity never changes declaration intent. It runs on detached candidates.
func (s *Session) refreshActivity() {
	if len(s.panes) == 0 {
		return
	}
	for _, t := range s.tabs {
		eligible := false
		for i := range t.Pages {
			p := &t.Pages[i]
			w := s.panes[p.Handle.Path]
			a := s.intent[p.InstancePath]
			p.Enabled, p.Visible = a.enabled, a.visible
			p.Label = w.Label
			if p.ID == t.Selected && a.enabled && a.visible {
				eligible = true
			}
		}
		if !eligible {
			t.Selected = ""
			for _, p := range t.Pages {
				if p.Enabled && p.Visible {
					t.Selected = p.ID
					break
				}
			}
		}
	}
	s.active = map[string]activity{}
	var walk func(*parser.Instance, bool, bool)
	walk = func(n *parser.Instance, en, vis bool) {
		a := activity{n.Layout["enabled"] != false, n.Layout["visible"] != false}
		if intent, ok := s.intent[n.Path]; ok {
			a = intent
		}
		parentEn, parentVis := en, vis
		en, vis = en && a.enabled, vis && a.visible
		s.active[n.Path] = activity{en, vis}
		w := s.currentControl(public(n.Path))
		if w != nil && w.InstancePath == n.Path {
			w.Enabled, w.Visible = en, vis
			w.ancestorEnabled, w.ancestorVisible = parentEn, parentVis
			if !en || !vis {
				if c := s.collections[w.Handle.Path]; c != nil {
					c.cancel()
				}
			}
		}
		t := s.tabs[public(n.Path)]
		if t != nil {
			t.Enabled, t.Visible = en, vis
		}
		p := s.splits[public(n.Path)]
		if p != nil {
			p.Enabled, p.Visible = en, vis
		}
		for _, r := range n.Regions {
			walk(r.Node, en, vis)
		}
		i := 0
		for _, row := range n.Rows {
			for _, child := range row {
				active := true
				if t != nil {
					active = i < len(t.Pages) && t.Pages[i].ID == t.Selected
				}
				if p != nil {
					active = !(p.Collapsed == SplitFirst && i == 0 || p.Collapsed == SplitSecond && i == 1)
				}
				walk(child, en, vis && active)
				i++
			}
		}
	}
	walk(s.root, true, true)
	if w := s.currentControl(s.focused); w != nil && w.Enabled && w.Visible && w.Handle.Kind != "page" {
		return
	}
	old := s.currentControl(s.focused)
	s.focused = ""
	if old == nil {
		return
	}
	longest := 0
	for _, w := range s.panes {
		if w.Handle.Kind != "page" && w.Enabled && w.Visible && within(old.InstancePath, w.InstancePath) && len(w.InstancePath) > longest {
			s.focused = w.Handle.Path
			longest = len(w.InstancePath)
		}
	}
}
func (s *Session) rememberFocus(path string) {
	w := s.currentControl(path)
	if w == nil {
		return
	}
	for _, t := range s.tabs {
		for i := range t.Pages {
			if within(w.InstancePath, t.Pages[i].InstancePath) {
				t.Pages[i].RememberedFocus = path
			}
		}
	}
}
func (s *Session) EnterPage(h Handle) error {
	t, ok := s.Tabs(h)
	if !ok || !t.Enabled || !t.Visible {
		return fault("page-focus", "Tabs is not active")
	}
	for _, p := range t.Pages {
		if p.ID != t.Selected {
			continue
		}
		if w := s.currentControl(p.RememberedFocus); w != nil && within(w.InstancePath, p.InstancePath) && w.Enabled && w.Visible {
			return s.Focus(w.Handle)
		}
		var first *Widget
		s.root.Walk(func(n *parser.Instance) {
			if first != nil || !within(n.Path, p.InstancePath) {
				return
			}
			if w := s.currentControl(public(n.Path)); w != nil && w.Enabled && w.Visible && w.Handle.Kind != "page" && w.Handle.Kind != "svg" {
				first = w
			}
		})
		if first != nil {
			return s.Focus(first.Handle)
		}
	}
	return s.Focus(h)
}
func (s *Session) inactivePaneViewport(path string) bool {
	for _, t := range s.tabs {
		for _, p := range t.Pages {
			if within(path, p.InstancePath) && !s.active[path].visible {
				return true
			}
		}
	}
	for _, p := range s.splits {
		if within(path, p.InstancePath) && !s.active[path].visible {
			return true
		}
	}
	return false
}
