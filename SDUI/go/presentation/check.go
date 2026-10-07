package presentation

import "github.com/Hans-Einar/SDP/SDUI/go/parser"

// check rejects unsupported hidden nodes too, before returning any artifact.
func check(root *parser.Instance) error {
	profile, err := parser.EffectiveProfile(root)
	if err != nil {
		return err
	}
	root.Walk(func(n *parser.Instance) {
		if err != nil {
			return
		}
		switch n.Kind {
		case "frame", "group", "markdown":
			return
		case "widget":
		default:
			err = diagnostic("export-kind", n.Path+": unsupported node "+n.Kind, n)
			return
		}
		switch n.Widget {
		case "button", "input", "svg":
		case "tree", "list":
			if profile != "sdui/0.3" {
				err = diagnostic("widget-kind", n.Path+": collection requires sdui/0.3", n)
			}
		default:
			err = diagnostic("widget-kind", n.Path+": unsupported widget "+n.Widget, n)
		}
	})
	return err
}
