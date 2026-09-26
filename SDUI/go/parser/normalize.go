package parser

import (
	"fmt"
	"strings"
)

type Region struct {
	Role string    `json:"role"`
	Node *Instance `json:"node"`
}
type Instance struct {
	Kind      string         `json:"kind"`
	Path      string         `json:"path"`
	Variant   string         `json:"variant,omitempty"`
	Widget    string         `json:"widget,omitempty"`
	Text      string         `json:"text,omitempty"`
	Arguments map[string]any `json:"arguments"`
	Layout    map[string]any `json:"layout"`
	Rows      [][]*Instance  `json:"rows"`
	Regions   []Region       `json:"regions"`
	Span      Span           `json:"span"`
}

func (i *Instance) Region(role string) *Instance {
	for _, r := range i.Regions {
		if r.Role == role {
			return r.Node
		}
	}
	return nil
}
func (i *Instance) Walk(fn func(*Instance)) {
	fn(i)
	for _, r := range i.Regions {
		r.Node.Walk(fn)
	}
	for _, row := range i.Rows {
		for _, c := range row {
			c.Walk(fn)
		}
	}
}
func (i *Instance) Argument(name string) string {
	if v, ok := i.Arguments[name].(Literal); ok {
		if text, ok := v.Value.(string); ok {
			return text
		}
	}
	return ""
}
func checkWrap(p map[string]any, rows [][]*Instance, s Span) {
	if p["wrap"] == "wrap" {
		for _, row := range rows {
			for _, child := range row {
				x := child.Layout["x"]
				if x == "fill" || strings.HasSuffix(fmt.Sprint(x), "fr") {
					fail("wrap-layout", "Wrapped child has fill/fr", s)
				}
			}
		}
	}
}

// Normalize returns bounded, independent instance trees. Mutating a tree never changes the AST.
func Normalize(d *Document) (roots map[string]*Instance, err error) {
	defer recoverDiagnostic(&err)
	validateLocal(d)
	defs := map[string]*Node{}
	for _, def := range d.Definitions {
		defs[def.Name] = def.Root
	}
	count := 0
	var expand func(*Node, string, int) *Instance
	expand = func(n *Node, path string, depth int) *Instance {
		count++
		if count > 8192 {
			fail("expansion-limit", "More than 8192 expanded components", n.Span)
		}
		if depth > 64 {
			fail("depth-limit", "Expanded depth exceeds 64", n.Span)
		}
		if n.Kind == "use" {
			i := expand(defs[val(n.Target)], path, depth+1)
			overlay := formatting(n, i.Kind, len(i.Rows))
			for k, v := range overlay {
				i.Layout[k] = v
			}
			checkCombination(i.Layout, i.Kind, len(i.Rows), n.Span)
			checkWrap(i.Layout, i.Rows, n.Span)
			return i
		}
		props := formatting(n, n.Kind, len(n.Rows))
		i := &Instance{Kind: n.Kind, Path: path, Widget: val(n.Widget), Layout: props, Arguments: map[string]any{}, Rows: [][]*Instance{}, Regions: []Region{}, Span: n.Span}
		if n.Variant != nil {
			i.Variant = "box"
		}
		if n.Text != nil {
			i.Text = n.Text.Value.(string)
		}
		if n.Kind == "widget" {
			i.Arguments = widgetArguments(n)
		}
		for r, row := range n.Rows {
			items := []*Instance{}
			for c, child := range row.Items {
				name := val(child.Role)
				if name == "" {
					name = val(child.Name)
				}
				if name == "" {
					name = val(child.Target)
				}
				if name == "" {
					name = fmt.Sprintf("$r%dc%d", r, c)
				}
				instance := expand(child, path+"/"+name, depth+1)
				if n.Kind == "group" && instance.Kind == "frame" {
					fail("group-content", "Frame reference in widget group", child.Span)
				}
				if child.Role != nil {
					i.Regions = append(i.Regions, Region{*child.Role, instance})
				} else {
					items = append(items, instance)
				}
			}
			if len(items) > 0 {
				i.Rows = append(i.Rows, items)
			}
		}
		checkWrap(props, i.Rows, n.Span)
		return i
	}
	roots = map[string]*Instance{}
	for _, def := range d.Definitions {
		roots[def.Name] = expand(def.Root, def.Name, 1)
	}
	validateConnections(d, roots)
	return roots, nil
}

// Compile shares exactly the Parse/Normalize path used by all hosts.
func Compile(source string) (*Document, map[string]*Instance, error) {
	d, err := Parse(source)
	if err != nil {
		return nil, nil, err
	}
	roots, err := Normalize(d)
	if err != nil {
		return nil, nil, err
	}
	return d, roots, nil
}
