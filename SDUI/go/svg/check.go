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
	if !options.SkipControls && len(options.NativeControls) > 0 {
		return &parser.Diagnostic{Code: "native-controls", Message: "Native control inventory requires SkipControls", Span: root.Span}
	}
	used := map[string]bool{}
	var invalid error
	root.Walk(func(n *parser.Instance) {
		if invalid != nil {
			return
		}
		fail := func(code, msg string) {
			invalid = &parser.Diagnostic{Code: code, Message: n.Path + ": " + msg, Span: n.Span}
		}
		switch n.Kind {
		case "frame", "group", "markdown":
			return
		case "widget":
		default:
			fail("export-kind", "Unsupported node "+n.Kind)
			return
		}
		switch n.Widget {
		case "button", "input", "svg":
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
