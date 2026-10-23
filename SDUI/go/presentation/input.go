package presentation

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Hans-Einar/SDP/SDUI/go/parser"
)

func inputDescription(n *parser.Instance) bool {
	if n.Widget != "input" {
		return false
	}
	p, err := parser.InputOptions(n)
	return err == nil && p.Extended
}

func inputText(n *parser.Instance) string {
	p, _ := parser.InputOptions(n) // check has validated the selected tree.
	parts := []string{"Static input", n.Path, "text=" + strconv.Quote(n.Argument("text")),
		"value=" + strconv.Quote(n.Argument("value")), "source initial text; not live state"}
	for _, key := range []string{"multiline", "readOnly", "required"} {
		value := "false (default)"
		if v, ok := n.Arguments[key].(parser.Literal); ok {
			value = fmt.Sprint(v.Value)
		}
		parts = append(parts, key+"="+value)
	}
	placeholder := "placeholder=" + strconv.Quote(p.Placeholder)
	if !p.PlaceholderSet {
		placeholder += " (text fallback)"
	}
	parts = append(parts, placeholder)
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
