package parser

import "strings"

type widgetSchema struct {
	position string
	fields   map[string]string
	required string
}

var widgets = map[string]widgetSchema{
	"button": {"label", map[string]string{"label": "string", "callback": "reference"}, "label"},
	"input":  {"text", map[string]string{"text": "string", "value": "string", "callback": "reference"}, "text"},
	"svg":    {"source", map[string]string{"source": "reference", "label": "string"}, "source"},
}

func widgetArguments(n *Node) map[string]any {
	out := map[string]any{}
	for _, a := range n.Arguments {
		k := val(a.Name)
		if a.Name == nil {
			k = widgets[val(n.Widget)].position
		}
		out[k] = a.Value
	}
	return out
}
func validateWidget(n *Node, modules map[string]bool) {
	schema, ok := widgets[val(n.Widget)]
	if !ok {
		fail("widget-kind", "Unsupported widget "+val(n.Widget), n.Span)
	}
	seen := map[string]bool{}
	for _, a := range n.Arguments {
		if a.Name == nil && len(seen) > 0 {
			fail("argument-order", "Only one positional argument before named arguments", a.Span)
		}
		k := val(a.Name)
		if a.Name == nil {
			k = schema.position
		}
		if seen[k] {
			fail("duplicate-argument", "Duplicate "+k, a.Span)
		}
		expected, ok := schema.fields[k]
		if !ok {
			fail("widget-property", "Unknown widget property "+k, a.Span)
		}
		switch v := a.Value.(type) {
		case Reference:
			if expected != "reference" {
				fail("argument-type", "Expected "+expected, a.Span)
			}
			if !modules[v.Module] {
				fail("unknown-module", "Unknown module "+v.Module, v.Span)
			}
		case Literal:
			if v.Kind != expected {
				fail("argument-type", "Expected "+expected, a.Span)
			}
		default:
			fail("argument-type", "Invalid argument", a.Span)
		}
		seen[k] = true
	}
	if seen["callback"] && n.Name == nil {
		fail("widget-name", "A callback requires a named widget", n.Span)
	}
	if !seen[schema.required] {
		fail("missing-argument", "Missing "+schema.required, n.Span)
	}
}
func validateLocal(d *Document) {
	if d == nil {
		fail("document", "Nil document", Span{})
	}
	modules := map[string]bool{}
	definitions := map[string]Definition{}
	for _, r := range d.References {
		if modules[r.Alias] {
			fail("duplicate-module", "Duplicate module", r.Span)
		}
		bad := r.Path == ""
		for _, c := range r.Path {
			bad = bad || c < 32 || c == 127
		}
		if bad {
			fail("module-path", "Invalid module path", r.Span)
		}
		modules[r.Alias] = true
	}
	for _, def := range d.Definitions {
		if _, ok := definitions[def.Name]; ok || modules[def.Name] {
			fail("duplicate-definition", "Ambiguous definition "+def.Name, def.Span)
		}
		definitions[def.Name] = def
	}
	deps := map[string][]string{}
	for _, def := range d.Definitions {
		names := map[string]bool{}
		var visit func(*Node)
		visit = func(n *Node) {
			name := val(n.Name)
			if name == "" && n.Kind == "use" {
				name = val(n.Target)
			}
			if name != "" {
				if names[name] {
					fail("duplicate-node", "Duplicate name "+name, n.Span)
				}
				names[name] = true
			}
			if n.Kind == "use" {
				if _, ok := definitions[val(n.Target)]; !ok {
					fail("unknown-definition", "Unknown component "+val(n.Target), n.Span)
				}
				deps[def.Name] = append(deps[def.Name], val(n.Target))
			} else {
				formatting(n, n.Kind, len(n.Rows))
			}
			if n.Variant != nil && (n.Kind != "frame" || !member(*n.Variant, "b box")) {
				fail("variant", "Only frame variants box/b supported", n.Span)
			}
			if n.Kind == "widget" {
				validateWidget(n, modules)
			}
			roles := map[string]bool{}
			body := false
			for _, row := range n.Rows {
				for _, child := range row.Items {
					if n.Kind == "group" && child.Kind == "frame" {
						fail("group-content", "Frame in widget group", child.Span)
					}
					if child.Role != nil {
						if roles[*child.Role] {
							fail("duplicate-region", "Duplicate region", child.Span)
						}
						roles[*child.Role] = true
					} else {
						body = true
					}
					visit(child)
				}
			}
			if roles["body"] && body {
				fail("body-conflict", "Explicit and implicit body cannot mix", n.Span)
			}
		}
		visit(def.Root)
	}
	done := map[string]bool{}
	stack := map[string]bool{}
	var cycle func(string, int)
	cycle = func(name string, depth int) {
		if stack[name] {
			fail("reference-cycle", "Recursive reference "+name, definitions[name].Span)
		}
		if depth >= 64 {
			fail("depth-limit", "Reference chain exceeds 64", definitions[name].Span)
		}
		if done[name] {
			return
		}
		stack[name] = true
		for _, child := range deps[name] {
			cycle(child, depth+1)
		}
		delete(stack, name)
		done[name] = true
	}
	for _, def := range d.Definitions {
		cycle(def.Name, 0)
	}
}
func validateConnections(d *Document, roots map[string]*Instance) {
	modules := map[string]bool{}
	for _, r := range d.References {
		modules[r.Alias] = true
	}
	paths := map[string]bool{}
	var visit func(*Instance)
	visit = func(n *Instance) {
		parts := strings.Split(n.Path, "/")
		if n.Kind == "widget" && !strings.HasPrefix(parts[len(parts)-1], "$") {
			public := []string{}
			for _, p := range parts {
				if !strings.HasPrefix(p, "$") {
					public = append(public, p)
				}
			}
			path := strings.Join(public, "/")
			if paths[path] {
				fail("duplicate-instance", "Ambiguous path "+path, n.Span)
			}
			paths[path] = true
		}
		for _, row := range n.Rows {
			for _, child := range row {
				visit(child)
			}
		}
		for _, reg := range n.Regions {
			visit(reg.Node)
		}
	}
	for _, def := range d.Definitions {
		visit(roots[def.Name])
	}
	objects, targets := map[string]bool{}, map[string]bool{}
	for _, con := range d.Connections {
		if !modules[con.Module] {
			fail("unknown-module", "Unknown module "+con.Module, con.Span)
		}
		if roots[con.Definition] == nil {
			fail("unknown-definition", "Unknown definition", con.Span)
		}
		target := con.Definition + "/" + strings.Join(con.Path, "/")
		if !paths[target] {
			fail("unknown-widget", "Unknown widget "+target, con.Span)
		}
		obj := con.Module + "." + con.Object
		if objects[obj] || targets[target] {
			fail("duplicate-connection", "One connection per object/widget", con.Span)
		}
		objects[obj] = true
		targets[target] = true
	}
}

// Validate checks local semantics including bounded instance expansion.
func Validate(d *Document) error { _, err := Normalize(d); return err }
