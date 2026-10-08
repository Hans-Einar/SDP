package runtime

import "github.com/Hans-Einar/SDP/SDUI/go/parser"

func boolArgument(n *parser.Instance, key string, fallback bool) bool {
	if a, ok := n.Arguments[key].(parser.Literal); ok {
		if b, ok := a.Value.(bool); ok {
			return b
		}
	}
	return fallback
}
func (s *Session) initCommands() error {
	identities, err := parser.ResolveInteractions(s.root)
	if err != nil {
		return err
	}
	s.commands = map[string]*CommandState{}
	s.presentations = map[string]*CommandPresentation{}
	s.menus = map[string]*MenuState{}
	s.surfaces = map[string]*SurfaceState{}
	s.aux = map[string]*Widget{}
	s.ownerSurface = map[string]string{}
	s.menuOwner = map[string]string{}
	s.menuTargets = map[string]Handle{}
	s.publishedSurfaces = map[SurfaceTarget]bool{}
	nodes := map[string]*parser.Instance{}
	var walk func(*parser.Instance, string, string)
	walk = func(n *parser.Instance, dialog, menu string) {
		nodes[n.Path] = n
		path := public(n.Path)
		s.ownerSurface[n.Path] = dialog
		switch n.Widget {
		case "command", "item", "separator", "menu", "menuGroup", "dialog":
			s.generation++
			w := &Widget{Handle: Handle{s.ID, path, s.generation, n.Widget}, InstancePath: n.Path, Label: n.Argument("label")}
			if r, ok := n.Arguments["callback"].(parser.Reference); ok {
				w.Binding = r
			}
			s.aux[path] = w
		}
		if n.Widget == "dialog" {
			s.surfaces[path] = &SurfaceState{Handle: s.aux[path].Handle, InstancePath: n.Path, Label: n.Argument("label"), Modal: boolArgument(n, "modal", true)}
			dialog = path
		}
		if n.Widget == "menu" && n.Argument("mode") != "submenu" {
			menu = path
			s.menus[path] = &MenuState{Handle: s.aux[path].Handle, InstancePath: n.Path}
		}
		if menu != "" {
			s.menuOwner[path] = menu
		}
		for _, r := range n.Regions {
			walk(r.Node, dialog, menu)
		}
		for _, row := range n.Rows {
			for _, c := range row {
				walk(c, dialog, menu)
			}
		}
	}
	walk(s.root, "", "")
	if len(s.aux) > 0 || len(identities) > 0 {
		s.generation++
		s.rootOwner = Handle{s.ID, s.root.Path, s.generation, "surface-root"}
	}
	for path, id := range identities {
		if nodes[path].Widget == "menu" && id.Target != "" {
			s.menuTargets[public(path)] = s.currentControl(public(id.Target)).Handle
		}
	}
	for path, id := range identities {
		if id.Command != path {
			continue
		}
		n := nodes[path]
		w := s.currentControl(public(path))
		if w == nil {
			return fault("command", "Missing command owner")
		}
		c := &CommandState{Handle: w.Handle, InstancePath: path, Label: n.Argument("label"), Icon: n.Argument("icon"), Tooltip: n.Argument("tooltip"), Key: n.Argument("key"), Context: n.Argument("context"), Effect: n.Argument("effect"), Toggle: boolArgument(n, "toggle", false), Checked: boolArgument(n, "checked", false), Exclusive: n.Argument("exclusive"), ExclusiveScope: id.Scope, Binding: w.Binding}
		if c.Context == "" {
			c.Context = "none"
		}
		if id.Target != "" {
			target := s.currentControl(public(id.Target))
			if target == nil {
				return fault("command-target", id.Target)
			}
			c.Target = target.Handle
		}
		s.commands[w.Handle.Path] = c
		if n.Widget == "button" {
			w.Binding = parser.Reference{}
		}
	}
	for path, id := range identities {
		n := nodes[path]
		if id.Command == "" || n.Widget == "command" {
			continue
		}
		c := s.commands[public(id.Command)]
		w := s.currentControl(public(path))
		if c == nil || w == nil {
			return fault("command", "Missing presentation owner")
		}
		p := &CommandPresentation{Handle: w.Handle, Command: c.Handle, InstancePath: path, Label: c.Label, Icon: c.Icon, Tooltip: c.Tooltip}
		if _, ok := n.Arguments["label"]; ok {
			p.Label = n.Argument("label")
		}
		if _, ok := n.Arguments["icon"]; ok {
			p.Icon = n.Argument("icon")
		}
		if _, ok := n.Arguments["tooltip"]; ok {
			p.Tooltip = n.Argument("tooltip")
		}
		s.presentations[w.Handle.Path] = p
	}
	s.refreshActivity()
	return nil
}
func (s *Session) refreshCommands() {
	for path, c := range s.commands {
		w := s.currentControl(path)
		c.Enabled, c.Visible = w.Enabled, w.Visible
		c.Label = w.Label
	}
	for _, p := range s.presentations {
		w := s.currentControl(p.Handle.Path)
		c := s.commands[p.Command.Path]
		p.Enabled, p.Visible = w.Enabled && c.Enabled, w.Visible && c.Visible
		if p.Handle.Kind == "button" && c.Context != "none" {
			target := s.currentControl(c.Target.Path)
			eligible := target != nil && target.Enabled && target.Visible
			if c.Context == "item" {
				coll := s.collections[c.Target.Path]
				eligible = eligible && coll != nil && coll.Selected != ""
				if eligible {
					x, ok := coll.item(coll.Selected)
					eligible = ok && x.Kind == Row && coll.visible(coll.Selected)
				}
			}
			p.Enabled = p.Enabled && eligible
		}
	}
	for _, p := range s.presentations {
		w := s.currentControl(p.Handle.Path)
		if w.Label == "" || p.Handle == p.Command {
			p.Label = s.commands[p.Command.Path].Label
		} else {
			p.Label = w.Label
		}
	}

	for path, d := range s.surfaces {
		d.Label = s.aux[path].Label
	}
}
