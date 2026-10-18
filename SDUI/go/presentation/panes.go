package presentation

import (
	"fmt"
	"strings"

	"github.com/Hans-Einar/SDP/SDUI/go/parser"
)

func paneText(n *parser.Instance) string {
	parts := []string{"Static " + n.Widget, n.Path}
	if n.Widget == "split" {
		parts = append(parts, "axis="+n.Argument("axis"))
		for _, key := range []string{"proportion", "minFirst", "minSecond", "collapsible"} {
			if v, ok := n.Arguments[key].(parser.Literal); ok {
				parts = append(parts, fmt.Sprintf("%s=%v", key, v.Value))
			}
		}
	} else {
		parts = append(parts, "label="+n.Argument("label"))
		if icon := n.Argument("icon"); icon != "" {
			parts = append(parts, "symbolic icon="+icon+" (not loaded)")
		}
		if n.Widget == "tabs" {
			selected := n.Argument("selected")
			if selected == "" {
				selected = "first eligible page"
			}
			parts = append(parts, "initial selection="+selected+" (not live state)")
		}
	}
	if ref, ok := n.Arguments["callback"].(parser.Reference); ok {
		parts = append(parts, "ActivatePage callback="+ref.Module+"."+ref.Object+".@"+ref.Member+" (not executed)")
	}
	for _, key := range []string{"enabled", "visible"} {
		if v, ok := n.Layout[key]; ok {
			parts = append(parts, fmt.Sprintf("%s=%v", key, v))
		}
	}
	parts = append(parts, fmt.Sprintf("source=%d:%d", n.Span.Line, n.Span.Column))
	return strings.Join(parts, " | ")
}
