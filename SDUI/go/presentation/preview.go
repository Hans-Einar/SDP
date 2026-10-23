package presentation

import (
	"fmt"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"strconv"
	"strings"
)

func previewDescription(n *parser.Instance) bool {
	if n.Kind != "markdown" && n.Widget != "svg" {
		return false
	}
	p, err := parser.PreviewOptions(n)
	return err == nil && p.Explicit
}

func previewText(n *parser.Instance) string {
	p, _ := parser.PreviewOptions(n) // check already validated the selected root.
	kind := "SVG"
	content := ""
	if n.Kind == "markdown" {
		kind = "Markdown"
		content = "text=" + strconv.Quote(n.Text)
	} else {
		r := n.Arguments["source"].(parser.Reference)
		content = "source=" + r.Module + "." + r.Object + ".@" + r.Member
	}
	parts := []string{"Static " + kind + " preview", n.Path, content, "description=" + strconv.Quote(p.Description), "fallback=" + strconv.Quote(p.Fallback), "resources not supplied; source intent only"}
	if label, ok := n.Arguments["label"].(parser.Literal); ok {
		parts = append(parts, "label="+strconv.Quote(label.Value.(string)))
	}
	for _, key := range []string{"enabled", "visible"} {
		if v, ok := n.Layout[key]; ok {
			parts = append(parts, fmt.Sprintf("%s=%v", key, v))
		}
	}
	parts = append(parts, fmt.Sprintf("source=%d:%d", n.Span.Line, n.Span.Column))
	return strings.Join(parts, " | ")
}
