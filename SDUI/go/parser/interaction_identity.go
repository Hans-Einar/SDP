package parser

import "strings"

// InteractionIdentity contains exact frontend-resolved paths; empty Dialog means entry canvas.
type InteractionIdentity struct{ Scope, Command, Target, Dialog string }
type interactionIndex struct {
	nodes   map[string]*Instance
	named   map[string][]*Instance
	dialogs map[string]string
}

func namedPath(path string) string {
	parts := []string{}
	for _, p := range strings.Split(path, "/") {
		if !strings.HasPrefix(p, "$") {
			parts = append(parts, p)
		}
	}
	return strings.Join(parts, "/")
}
func indexInteractions(root *Instance) *interactionIndex {
	index := &interactionIndex{map[string]*Instance{}, map[string][]*Instance{}, map[string]string{}}
	var visit func(*Instance, string, int)
	visit = func(n *Instance, dialog string, depth int) {
		if n == nil || depth > 64 || len(index.nodes) >= 8192 {
			fail("interaction-tree", "Invalid or oversized selected tree", root.Span)
		}
		if _, ok := index.nodes[n.Path]; ok {
			fail("interaction-path", "Duplicate instance path "+n.Path, n.Span)
		}
		index.nodes[n.Path] = n
		index.dialogs[n.Path] = dialog
		parts := strings.Split(n.Path, "/")
		if !strings.HasPrefix(parts[len(parts)-1], "$") {
			key := namedPath(n.Path)
			index.named[key] = append(index.named[key], n)
		}
		if n.Kind == "composition" && n.Widget == "dialog" {
			dialog = n.Path
		}
		for _, region := range n.Regions {
			visit(region.Node, dialog, depth+1)
		}
		for _, row := range n.Rows {
			for _, child := range row {
				visit(child, dialog, depth+1)
			}
		}
	}
	visit(root, "", 1)
	return index
}
func referenceParts(path string, absolute bool, span Span) (string, bool) {
	isAbsolute := strings.HasPrefix(path, "/")
	if isAbsolute && !absolute {
		fail("interaction-reference", "Absolute field path forbidden", span)
	}
	raw := path
	if isAbsolute {
		raw = strings.TrimPrefix(raw, "/")
	}
	for _, p := range strings.Split(raw, "/") {
		if p == "" || p == "." || p == ".." {
			fail("interaction-reference", "Nonempty named path segments required: "+path, span)
		}
		for i, r := range p {
			if !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r == '_' || i > 0 && r >= '0' && r <= '9') {
				fail("interaction-reference", "Invalid named path: "+path, span)
			}
		}
	}
	return raw, isAbsolute
}
func (index *interactionIndex) resolve(root, owner *Instance, key, scope string) *Instance {
	v, ok := owner.Arguments[key].(Literal)
	if !ok || v.Kind != "string" {
		fail("interaction-reference", "Expected string "+key, owner.Span)
	}
	path, ok := v.Value.(string)
	if !ok {
		fail("interaction-reference", "Expected string "+key, v.Span)
	}
	relative, absolute := referenceParts(path, true, v.Span)
	base := scope
	if absolute {
		base = root.Path
	}
	targets := index.named[namedPath(base)+"/"+relative]
	if len(targets) != 1 {
		fail("interaction-reference", owner.Path+": missing or ambiguous "+key+" "+path, v.Span)
	}
	return targets[0]
}
func contextControl(n *Instance) bool {
	return n.Kind == "widget" && (member(n.Widget, "button input svg tree list") || scalarKind(n.Widget))
}

// ResolveInteractions strictly resolves one selected tree, including hidden content.
// It performs no I/O or mutation. Reusable templates may need an enclosing entry.
func ResolveInteractions(root *Instance) (out map[string]InteractionIdentity, err error) {
	defer recoverDiagnostic(&err)
	profile, e := EffectiveProfile(root)
	if e != nil {
		return nil, e
	}
	if profile == "sdui/0.2" {
		return map[string]InteractionIdentity{}, nil
	}
	validateInteractionPlacement(root, nil)
	index := indexInteractions(root)
	out = map[string]InteractionIdentity{}
	exclusive := map[string]string{}
	keys := map[string]string{}
	root.Walk(func(n *Instance) {
		validateNormalizedInteraction(n)
		validateScalar(n)
		if n.Kind == "markdown" || n.Widget == "svg" || n.Widget == "markdown" {
			if _, err := PreviewOptions(n); err != nil {
				panic(err)
			}
		}
		if n.Widget == "input" {
			if _, err := InputOptions(n); err != nil {
				panic(err)
			}
		}
		relevant := isInteractionNode(n) || n.Kind == "widget" && (n.Widget == "input" || scalarKind(n.Widget)) && index.dialogs[n.Path] != ""
		if !relevant {
			return
		}
		id := InteractionIdentity{Scope: n.Argument("$scope"), Dialog: index.dialogs[n.Path]}
		if isInteractionNode(n) && (id.Scope == "" || index.nodes[id.Scope] == nil || n.Path != id.Scope && !strings.HasPrefix(n.Path, id.Scope+"/")) {
			fail("interaction-scope", n.Path+": missing or invalid definition-instance scope", n.Span)
		}
		if n.Widget == "command" || IsCommandButton(n) || n.Widget == "item" {
			if hasArg(n, "command") {
				target := index.resolve(root, n, "command", id.Scope)
				if target.Kind != "widget" || target.Widget != "command" {
					fail("command-reference", n.Path+": expected command declaration", n.Span)
				}
				id.Command = target.Path
			} else {
				id.Command = n.Path
			}
		}
		if n.Widget == "command" || IsCommandButton(n) && !hasArg(n, "command") {
			validateCommandSchema(n)
			context := n.Argument("context")
			effect := n.Argument("effect")
			if hasArg(n, "target") {
				target := index.resolve(root, n, "target", id.Scope)
				id.Target = target.Path
				if effect == "open" {
					if target.Kind != "composition" || target.Widget != "dialog" {
						fail("command-target", n.Path+": open requires dialog", n.Span)
					}
				} else if !contextControl(target) || context == "item" && !member(target.Widget, "tree list") {
					fail("command-target", n.Path+": wrong context target kind", n.Span)
				}
			}
			if member(effect, "accept cancel close") && id.Dialog == "" {
				fail("command-dialog", n.Path+": local effect requires enclosing dialog", n.Span)
			}
			if group := n.Argument("exclusive"); group != "" && boolArg(n, "checked") {
				key := id.Scope + "\x00" + group
				if exclusive[key] != "" {
					fail("command-exclusive", "Multiple initially checked exclusive members", n.Span)
				}
				exclusive[key] = n.Path
			}
			if key := n.Argument("key"); key != "" {
				k := id.Dialog + "\x00" + key
				if keys[k] != "" {
					fail("command-key", "Duplicate key in same surface", n.Span)
				}
				keys[k] = n.Path
			}
		}
		if n.Kind == "composition" && n.Widget == "menu" && n.Argument("mode") == "context" {
			target := index.resolve(root, n, "target", id.Scope)
			if !contextControl(target) {
				fail("menu-target", n.Path+": context menu requires control", n.Span)
			}
			id.Target = target.Path
		}
		out[n.Path] = id
	})
	var menuVisit func(*Instance, string)
	menuVisit = func(n *Instance, contextTarget string) {
		if n.Kind == "composition" && n.Widget == "menu" && n.Argument("mode") != "submenu" {
			contextTarget = out[n.Path].Target
		}
		if n.Kind == "widget" && n.Widget == "item" && contextTarget != "" {
			command := out[n.Path].Command
			owner := index.nodes[command]
			id := out[command]
			if owner != nil && member(owner.Argument("context"), "widget item") && id.Target != contextTarget {
				fail("menu-context", n.Path+": command target differs from captured menu target", n.Span)
			}
		}
		for _, r := range n.Regions {
			menuVisit(r.Node, contextTarget)
		}
		for _, row := range n.Rows {
			for _, child := range row {
				menuVisit(child, contextTarget)
			}
		}
	}
	menuVisit(root, "")
	return out, nil
}

// ResolveDialogField resolves a named relative input, excluding nested dialogs.
func ResolveDialogField(root *Instance, dialogPath, fieldPath string) (field *Instance, err error) {
	defer recoverDiagnostic(&err)
	if _, e := ResolveInteractions(root); e != nil {
		return nil, e
	}
	index := indexInteractions(root)
	dialog := index.nodes[dialogPath]
	if dialog == nil || dialog.Kind != "composition" || dialog.Widget != "dialog" {
		fail("dialog-field", "Unknown dialog "+dialogPath, root.Span)
	}
	relative, _ := referenceParts(fieldPath, false, dialog.Span)
	targets := index.named[namedPath(dialogPath)+"/"+relative]
	if len(targets) != 1 || targets[0].Kind != "widget" || targets[0].Widget != "input" || index.dialogs[targets[0].Path] != dialogPath {
		fail("dialog-field", "Expected an input owned by "+dialogPath+": "+fieldPath, dialog.Span)
	}
	return targets[0], nil
}
