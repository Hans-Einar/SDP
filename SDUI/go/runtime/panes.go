package runtime

import (
	"sort"
	"strings"

	"github.com/Hans-Einar/SDP/SDUI/go/parser"
)

type SplitSide string

const (
	SplitNone   SplitSide = "none"
	SplitFirst  SplitSide = "first"
	SplitSecond SplitSide = "second"
)

type PageState struct {
	ID                                         string
	Handle                                     Handle
	InstancePath, Label, Icon, RememberedFocus string
	Enabled, Visible                           bool
}
type TabsState struct {
	Handle                 Handle
	InstancePath, Selected string
	Pages                  []PageState
	Enabled, Visible       bool
}
type SplitState struct {
	Handle                                           Handle
	InstancePath, Axis                               string
	Proportion, MinFirst, MinSecond, SavedProportion float64
	Collapsible                                      bool
	Collapsed                                        SplitSide
	Enabled, Visible                                 bool
}
type activity struct{ enabled, visible bool }

func (s *Session) initPanes() error {
	s.panes = map[string]*Widget{}
	s.tabs = map[string]*TabsState{}
	s.splits = map[string]*SplitState{}
	s.intent = map[string]activity{}
	s.interactions = map[string]InteractionHandler{}
	var first error
	s.root.Walk(func(n *parser.Instance) {
		if first != nil {
			return
		}
		if n.Kind == "widget" || n.Kind == "composition" {
			s.intent[n.Path] = activity{n.Layout["enabled"] != false, n.Layout["visible"] != false}
		}
		if n.Kind != "composition" {
			return
		}
		if n.Widget != "tabs" && n.Widget != "page" && n.Widget != "split" {
			if n.Widget == "dialog" || n.Widget == "menu" || n.Widget == "menuGroup" {
				return
			}
			first = fault("pane-kind", "Unsupported composition: "+n.Widget)
			return
		}
		path := public(n.Path)
		if s.panes[path] != nil || s.widgets[path] != nil {
			first = fault("ambiguous-instance", path)
			return
		}
		s.generation++
		w := &Widget{Handle: Handle{s.ID, path, s.generation, n.Widget}, InstancePath: n.Path, Label: n.Argument("label")}
		if ref, ok := n.Arguments["callback"].(parser.Reference); ok {
			w.Binding = ref
		}
		s.panes[path] = w
		switch n.Widget {
		case "tabs":
			children, err := parser.PaneChildren(n)
			if err != nil {
				first = err
				return
			}
			t := &TabsState{Handle: w.Handle, InstancePath: n.Path, Selected: n.Argument("selected")}
			for _, child := range children {
				t.Pages = append(t.Pages, PageState{ID: child.ID, InstancePath: child.Node.Path, Label: child.Node.Argument("label"), Icon: child.Node.Argument("icon")})
			}
			s.tabs[path] = t
		case "split":
			if _, err := parser.PaneChildren(n); err != nil {
				first = err
				return
			}
			p := numberArgument(n, "proportion", .5)
			v := &SplitState{Handle: w.Handle, InstancePath: n.Path, Axis: n.Argument("axis"), Proportion: p, SavedProportion: p, MinFirst: numberArgument(n, "minFirst", 0), MinSecond: numberArgument(n, "minSecond", 0), Collapsible: true, Collapsed: SplitNone}
			if a, ok := n.Arguments["collapsible"].(parser.Literal); ok {
				v.Collapsible, _ = a.Value.(bool)
			}
			if (v.Axis != "horizontal" && v.Axis != "vertical") || !finiteRatio(p) || !finiteRatio(v.MinFirst) || !finiteRatio(v.MinSecond) || v.MinFirst+v.MinSecond >= 1 || p < v.MinFirst || p > 1-v.MinSecond {
				first = fault("split-state", "Invalid split arguments")
				return
			}
			s.splits[path] = v
		}
	})
	if first != nil {
		return first
	}
	for _, t := range s.tabs {
		found := t.Selected == ""
		for i := range t.Pages {
			p := &t.Pages[i]
			p.Handle = s.panes[public(p.InstancePath)].Handle
			a := s.intent[p.InstancePath]
			if p.ID == t.Selected {
				found = a.enabled && a.visible
			}
		}
		if !found {
			return fault("page-selection", "Initial selected page is unknown or ineligible")
		}
	}
	s.refreshActivity()
	return nil
}
func numberArgument(n *parser.Instance, key string, fallback float64) float64 {
	if a, ok := n.Arguments[key].(parser.Literal); ok {
		if v, ok := a.Value.(float64); ok {
			return v
		}
	}
	return fallback
}
func copyTabs(t *TabsState) TabsState {
	v := *t
	v.Pages = append([]PageState(nil), t.Pages...)
	return v
}
func (s *Session) copyPanes(n *Session) {
	n.panes = map[string]*Widget{}
	for k, w := range s.panes {
		v := *w
		n.panes[k] = &v
	}
	n.tabs = map[string]*TabsState{}
	for k, t := range s.tabs {
		v := copyTabs(t)
		n.tabs[k] = &v
	}
	n.splits = map[string]*SplitState{}
	for k, p := range s.splits {
		v := *p
		n.splits[k] = &v
	}
	n.intent = map[string]activity{}
	for k, v := range s.intent {
		n.intent[k] = v
	}
	n.active = map[string]activity{}
	for k, v := range s.active {
		n.active[k] = v
	}
}
func (s *Session) Pane(path string) (Handle, bool) {
	w := s.panes[path]
	if s.closed || w == nil {
		return Handle{}, false
	}
	return w.Handle, true
}
func (s *Session) pane(h Handle) (*Widget, error) {
	if s.closed {
		return nil, fault("closed", "Session is closed")
	}
	w := s.panes[h.Path]
	if w == nil || w.Handle != h {
		return nil, fault("stale-handle", "Pane is no longer current")
	}
	return w, nil
}
func (s *Session) Tabs(h Handle) (TabsState, bool) {
	if _, err := s.pane(h); err != nil {
		return TabsState{}, false
	}
	t := s.tabs[h.Path]
	if t == nil {
		return TabsState{}, false
	}
	return copyTabs(t), true
}
func (s *Session) Split(h Handle) (SplitState, bool) {
	if _, err := s.pane(h); err != nil {
		return SplitState{}, false
	}
	p := s.splits[h.Path]
	if p == nil {
		return SplitState{}, false
	}
	return *p, true
}
func (s *Session) CallbackOwners() []Widget {
	owners := map[string]Widget{}
	for path, w := range s.panes {
		if w.Handle.Kind == "tabs" && w.Binding.Module != "" {
			owners[path] = *w
		}
	}
	for path, c := range s.commands {
		if c.Binding.Module != "" {
			w := *s.currentControl(path)
			w.Binding = c.Binding
			owners[path] = w
		}
	}
	for path := range s.surfaces {
		w := s.aux[path]
		if w.Binding.Module != "" {
			owners[path] = *w
		}
	}
	keys := []string{}
	for k := range owners {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	out := []Widget{}
	for _, k := range keys {
		out = append(out, owners[k])
	}
	return out
}
func (s *Session) currentControl(path string) *Widget {
	if w := s.widgets[path]; w != nil {
		return w
	}
	if w := s.panes[path]; w != nil {
		return w
	}
	return s.aux[path]
}
func within(path, parent string) bool { return path == parent || strings.HasPrefix(path, parent+"/") }
