package presentation

import (
	"fmt"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"strings"
)

func scalarDescription(n *parser.Instance) bool {
	return n.Profile == "sdui/0.3" && n.Kind == "widget" && (n.Widget == "checkbox" || n.Widget == "slider" || n.Widget == "number" || n.Widget == "select")
}
func scalarText(n *parser.Instance) string {
	parts := []string{"Static " + n.Widget, n.Path, "label=" + n.Argument("label"), "source initial values; not live state"}
	for _, key := range []string{"min", "max", "step", "value", "readOnly", "required", "placeholder"} {
		if v, ok := n.Arguments[key].(parser.Literal); ok {
			parts = append(parts, fmt.Sprintf("%s=%v", key, v.Value))
		}
	}
	if _, ok := n.Arguments["readOnly"]; !ok {
		parts = append(parts, "readOnly=false (default)")
	}
	if n.Widget == "checkbox" {
		if _, ok := n.Arguments["value"]; !ok {
			parts = append(parts, "value=false (default)")
		}
	}
	if n.Widget == "select" {
		if _, ok := n.Arguments["value"]; !ok {
			parts = append(parts, "value=empty option ID (default)")
		}
		parts = append(parts, "choice options not supplied; IDs are not labels")
	}
	if ref, ok := n.Arguments["callback"].(parser.Reference); ok {
		parts = append(parts, "Commit callback="+ref.Module+"."+ref.Object+".@"+ref.Member+" (not executed; binding not verified)")
	} else {
		parts = append(parts, "no source Commit binding; local acceptance requires runtime adapter")
	}
	for _, key := range []string{"enabled", "visible"} {
		if v, ok := n.Layout[key]; ok {
			parts = append(parts, fmt.Sprintf("%s=%v", key, v))
		}
	}
	parts = append(parts, fmt.Sprintf("source=%d:%d", n.Span.Line, n.Span.Column))
	return strings.Join(parts, " | ")
}
