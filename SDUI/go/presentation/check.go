package presentation

import "github.com/Hans-Einar/SDP/SDUI/go/parser"

// check rejects unsupported hidden nodes too, before returning any artifact.
func check(root *parser.Instance) error {
	profile, err := parser.EffectiveProfile(root)
	if err != nil {
		return err
	}
	if _, err := parser.ResolveInteractions(root); err != nil {
		return err
	}
	root.Walk(func(n *parser.Instance) {
		if err != nil {
			return
		}
		switch n.Kind {
		case "frame", "group", "markdown":
			return
		case "composition":
			if profile != "sdui/0.3" || (n.Widget != "tabs" && n.Widget != "page" && n.Widget != "split" && n.Widget != "menu" && n.Widget != "menuGroup" && n.Widget != "dialog") {
				err = diagnostic("export-kind", n.Path+": unsupported composition "+n.Widget, n)
			} else if n.Widget == "tabs" || n.Widget == "split" {
				_, err = parser.PaneChildren(n)
			}
			return
		case "widget":
		default:
			err = diagnostic("export-kind", n.Path+": unsupported node "+n.Kind, n)
			return
		}
		switch n.Widget {
		case "button", "input", "svg":
		case "tree", "list", "command", "item", "separator", "checkbox", "slider", "select", "number":
			if profile != "sdui/0.3" {
				err = diagnostic("widget-kind", n.Path+": widget requires sdui/0.3", n)
			}
		default:
			err = diagnostic("widget-kind", n.Path+": unsupported widget "+n.Widget, n)
		}
	})
	return err
}
