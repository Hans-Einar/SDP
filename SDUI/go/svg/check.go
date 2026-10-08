package svg

import (
	"fmt"

	"github.com/Hans-Einar/SDP/SDUI/go/parser"
)

// Check validates the entire source tree, including hidden nodes, before output
// or measurement. NativeControls is a host-owned inventory, not an SVG fallback.
// It must describe controls actually prepared for this model revision.
func Check(root *parser.Instance, options Options) error {
	profile, err := parser.EffectiveProfile(root)
	if err != nil {
		return err
	}
	entry := root
	if options.InteractionRoot != nil {
		if !options.SkipControls {
			return &parser.Diagnostic{Code: "native-controls", Message: "InteractionRoot requires native SkipControls", Span: root.Span}
		}
		entry = options.InteractionRoot
		if _, err := parser.EffectiveProfile(entry); err != nil {
			return err
		}
		found := false
		entry.Walk(func(n *parser.Instance) {
			if n == root {
				found = true
			}
		})
		if !found {
			return &parser.Diagnostic{Code: "native-controls", Message: "Canvas root is not part of selected snapshot", Span: root.Span}
		}
	}
	if _, err := parser.ResolveInteractions(entry); err != nil {
		return err
	}
	if !options.SkipControls && len(options.NativeControls) > 0 {
		return &parser.Diagnostic{Code: "native-controls", Message: "Native control inventory requires SkipControls", Span: root.Span}
	}
	used := map[string]bool{}
	pages := map[*parser.Instance]bool{}
	menuChildren := map[*parser.Instance]bool{}
	var invalid error
	root.Walk(func(n *parser.Instance) {
		if invalid != nil {
			return
		}
		fail := func(code, msg string) {
			invalid = &parser.Diagnostic{Code: code, Message: n.Path + ": " + msg, Span: n.Span}
		}
		if !options.SkipControls && profile == "sdui/0.3" && n.Kind == "widget" && n.Widget == "button" && (n.Argument("icon") != "" || n.Argument("tooltip") != "") {
			fail("unsupported-interaction-export", "SVG export does not render button icon/tooltip decorations")
			return
		}
		m2 := n.Kind == "composition" && (n.Widget == "menu" || n.Widget == "menuGroup" || n.Widget == "dialog") || n.Kind == "widget" && (n.Widget == "command" || n.Widget == "item" || n.Widget == "separator") || parser.IsCommandButton(n)
		if m2 {
			if profile != "sdui/0.3" {
				fail("export-kind", "M2 requires sdui/0.3")
				return
			}
			if !options.SkipControls {
				fail("unsupported-interaction-export", "SVG export does not support "+n.Widget)
				return
			}
			switch n.Widget {
			case "command":
				return // Nonvisual declaration; strict entry resolution already succeeded.
			case "item", "separator", "menuGroup":
				if !menuChildren[n] {
					fail("native-controls", "Menu child requires enclosing prepared menu adapter")
					return
				}
				if n.Widget == "menuGroup" {
					for _, row := range n.Rows {
						for _, child := range row {
							menuChildren[child] = true
						}
					}
				}
				return
			case "menu", "dialog":
				if options.NativeControls[n.Path] != n.Widget {
					fail("native-controls", "Missing matching prepared native adapter for "+n.Widget)
					return
				}
				used[n.Path] = true
				if n.Widget == "menu" {
					for _, row := range n.Rows {
						for _, child := range row {
							menuChildren[child] = true
						}
					}
				}
				return
			}
		}
		switch n.Kind {
		case "frame", "group", "markdown":
			return
		case "composition":
			if profile != "sdui/0.3" || (n.Widget != "tabs" && n.Widget != "page" && n.Widget != "split") {
				fail("export-kind", "Unsupported composition "+n.Widget)
				return
			}
			if !options.SkipControls {
				fail("unsupported-pane-export", "SVG export does not support "+n.Widget)
				return
			}
			if n.Widget == "page" {
				if !pages[n] {
					fail("native-controls", "Page requires an enclosing prepared tabs adapter")
				}
				return
			}
			children, err := parser.PaneChildren(n)
			if err != nil {
				invalid = err
				return
			}
			if options.NativeControls[n.Path] != n.Widget {
				fail("native-controls", "Missing matching prepared native pane for "+n.Widget)
				return
			}
			used[n.Path] = true
			if n.Widget == "tabs" {
				for _, child := range children {
					pages[child.Node] = true
				}
			}
			return
		case "widget":
		default:
			fail("export-kind", "Unsupported node "+n.Kind)
			return
		}
		switch n.Widget {
		case "input":
			if profile != "sdui/0.3" {
				break // Preserve the legacy static-export path.
			}
			policy, err := parser.InputOptions(n)
			if err != nil {
				invalid = err
				return
			}
			if policy.Extended && !options.SkipControls {
				fail("unsupported-text-export", "SVG export does not support extended input properties")
				return
			}
		case "button", "svg":
		case "checkbox", "slider", "select", "number":
			if profile != "sdui/0.3" {
				fail("widget-kind", "Scalar control requires sdui/0.3")
				return
			}
			if !options.SkipControls {
				fail("unsupported-value-export", "SVG export does not support "+n.Widget)
				return
			}
		case "tree", "list":
			if profile != "sdui/0.3" {
				fail("widget-kind", "Collection requires sdui/0.3")
				return
			}
		default:
			fail("widget-kind", "Unsupported widget "+n.Widget)
			return
		}
		if options.SkipControls && n.Widget != "svg" {
			if options.NativeControls[n.Path] != n.Widget {
				fail("native-controls", "Missing matching prepared native control for "+n.Widget)
				return
			}
			used[n.Path] = true
			return
		}
		if n.Widget == "tree" || n.Widget == "list" {
			fail("unsupported-collection-export", "SVG export does not support "+n.Widget)
		}
	})
	if invalid != nil {
		return invalid
	}
	for path := range options.NativeControls {
		if !used[path] {
			return &parser.Diagnostic{Code: "native-controls", Message: fmt.Sprintf("%s: unused native control inventory entry", path), Span: root.Span}
		}
	}
	return nil
}
