package presentation

import (
	"fmt"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"strings"
)

func interactionDescription(n *parser.Instance) bool {
	return n.Profile == "sdui/0.3" && (n.Kind == "widget" && (n.Widget == "command" || n.Widget == "item" || n.Widget == "separator" || n.Widget == "button" && (parser.IsCommandButton(n) || n.Argument("icon") != "" || n.Argument("tooltip") != "")) || n.Kind == "composition" && (n.Widget == "menu" || n.Widget == "menuGroup" || n.Widget == "dialog"))
}
func interactionText(n *parser.Instance) string {
	parts := []string{"Static " + n.Widget, n.Path, "not live state"}
	for _, key := range []string{"label", "command", "mode", "modal", "toggle", "checked", "exclusive", "key", "context", "target", "effect", "tooltip"} {
		if v, ok := n.Arguments[key].(parser.Literal); ok {
			parts = append(parts, fmt.Sprintf("%s=%v", key, v.Value))
		}
	}
	if n.Widget == "dialog" {
		parts = append(parts, "initially closed; opening/Accept/Cancel/Close not simulated")
	}
	if n.Widget == "menu" {
		parts = append(parts, "popup/context capture not simulated")
	}
	if n.Widget == "button" && !parser.IsCommandButton(n) {
		parts = append(parts, "legacy Activate")
	}
	if icon := n.Argument("icon"); icon != "" {
		parts = append(parts, "symbolic icon="+icon+" (not loaded)")
	}
	if ref, ok := n.Arguments["callback"].(parser.Reference); ok {
		parts = append(parts, "callback="+ref.Module+"."+ref.Object+".@"+ref.Member+" (not executed)")
	}
	if scope := n.Argument("$scope"); scope != "" {
		parts = append(parts, "definition instance="+scope)
	}
	for _, key := range []string{"enabled", "visible"} {
		if v, ok := n.Layout[key]; ok {
			parts = append(parts, fmt.Sprintf("%s=%v", key, v))
		}
	}
	parts = append(parts, fmt.Sprintf("source=%d:%d", n.Span.Line, n.Span.Column))
	return strings.Join(parts, " | ")
}
