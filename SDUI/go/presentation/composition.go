package presentation

import (
	"encoding/json"
	"fmt"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"html"
	"strings"
)

// Combined preserves the static text renderer and adds source composition.
// Links refer to UTF-8 byte spans in the caller's revision-pinned source snapshot.
func Combined(root *parser.Instance, columns int, documents ...*parser.Document) (string, error) {
	text, err := Markdown(root, columns)
	if err != nil {
		return "", err
	}
	var diagrams, table strings.Builder
	table.WriteString("| Instance | Declaration / source | Reuse sites | Layout and bindings |\n| --- | --- | --- | --- |\n")
	cell := func(s string) string { return strings.ReplaceAll(codeSpan(s), "|", "&#124;") }
	link := func(s parser.Span) string {
		return fmt.Sprintf("[L%d:%d](sdui-source://%d/%d)", s.Line, s.Column, s.Start, s.End)
	}
	caption := func(n *parser.Instance) string {
		parts := strings.Split(n.Path, "/")
		s := parts[len(parts)-1] + " · " + n.Kind
		if n.Widget != "" {
			s = parts[len(parts)-1] + " · " + n.Widget
		}
		if n.Layout["visible"] == false {
			s += " · hidden"
		}
		if len(n.Uses) > 0 {
			s += " · reuse " + n.Declaration
		}
		return s
	}
	label := func(s string) string {
		return html.EscapeString(strings.ReplaceAll(strings.ReplaceAll(s, "\n", " "), "\r", " "))
	}
	var walk func(*parser.Instance)
	walk = func(n *parser.Instance) {
		uses := []string{}
		for _, u := range n.Uses {
			uses = append(uses, cell(u.Definition)+" "+link(u.Span))
		}
		layout, _ := json.Marshal(n.Layout)
		args, _ := json.Marshal(n.Arguments)
		fmt.Fprintf(&table, "| %s | %s %s | %s | %s %s |\n", cell(n.Path), cell(n.Declaration), link(n.Span), strings.Join(uses, ", "), cell(string(layout)), cell(string(args)))
		if n == root || len(n.Rows) > 0 || len(n.Regions) > 0 {
			fmt.Fprintf(&diagrams, "### %s\n\n```mermaid\nmindmap\n  root((%s))\n", codeSpan(n.Path), label(caption(n)))
			serial := 0
			for _, r := range n.Regions {
				serial++
				fmt.Fprintf(&diagrams, "    region%d[%s · %s]\n", serial, label(r.Role), label(caption(r.Node)))
			}
			for i, row := range n.Rows {
				serial++
				fmt.Fprintf(&diagrams, "    row%d[Row %d · left to right]\n", serial, i+1)
				for _, child := range row {
					serial++
					fmt.Fprintf(&diagrams, "      item%d[%s]\n", serial, label(caption(child)))
				}
			}
			diagrams.WriteString("```\n\n")
		}
		for _, r := range n.Regions {
			walk(r.Node)
		}
		for _, row := range n.Rows {
			for _, child := range row {
				walk(child)
			}
		}
	}
	walk(root)
	if len(documents) > 0 && documents[0] != nil {
		for _, ref := range documents[0].References {
			fmt.Fprintf(&table, "| %s | %s | — | %s |\n", cell("module "+ref.Alias), link(ref.Span), cell(ref.Path))
		}
		for _, con := range documents[0].Connections {
			if con.Definition == root.Path {
				fmt.Fprintf(&table, "| %s | %s | — | %s |\n", cell(con.Definition+"/"+strings.Join(con.Path, "/")), link(con.Span), cell("connection: "+con.Module+"."+con.Object+" (not executed)"))
			}
		}
	}
	return "# SDUI composition — " + codeSpan(root.Path) + "\n\nSource composition maps, not live UI or measured geometry. The overview shows regions and ordered rows; nested containers have their own detail maps below. Hidden and reused instances remain visible here. The text prototype follows visibility rules.\n\n" + diagrams.String() + "## Source and bindings\n\nSource links select exact UTF-8 byte spans from this revision. Reused instances link their declaration and use sites. Layout includes effective overrides; symbolic callbacks are not executed.\n\n" + table.String() + "\n---\n\n" + text, nil
}
