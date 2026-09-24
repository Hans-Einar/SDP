package presentation

import (
	"fmt"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"strings"
)

func quoted(lines []string) []string {
	out := make([]string, len(lines))
	for i, l := range lines {
		out[i] = ">"
		if l != "" {
			out[i] += " " + l
		}
	}
	return out
}
func Markdown(root *parser.Instance, columns int) (string, error) {
	overview, e := Dump(root, columns)
	if e != nil {
		return "", e
	}
	overview = strings.TrimRight(overview, "\n")
	fence := strings.Repeat("`", max(3, maxRun(overview, '`')+1))
	budget := 0
	var render func(*parser.Instance) ([]string, error)
	render = func(n *parser.Instance) ([]string, error) {
		if !visible(n) {
			return nil, nil
		}
		lines := []string{}
		if n.Kind == "markdown" {
			lines = quoted(contentLines(n.Text))
		} else if n.Kind == "widget" {
			text := ""
			switch n.Widget {
			case "button":
				text = "**Button:** " + codeSpan(n.Argument("label"))
			case "input":
				text = "**Input:** " + codeSpan(n.Argument("text")) + " — " + codeSpan(n.Argument("value"))
			default:
				label := n.Argument("label")
				if label == "" {
					label = n.Path
				}
				text = "**SVG placeholder:** " + codeSpan(label)
			}
			lines = []string{text}
		} else {
			kind := "Group"
			if n.Kind == "frame" {
				kind = "Frame"
			}
			if n.Variant == "box" {
				kind = "BoxUI-frame"
			}
			lines = []string{"**" + kind + ":** " + codeSpan(n.Path), ""}
			region := func(role string) error {
				if r := n.Region(role); visible(r) {
					block, e := render(r)
					if e != nil {
						return e
					}
					lines = append(lines, "**"+role+"**", "")
					lines = append(lines, block...)
					lines = append(lines, "")
				}
				return nil
			}
			for _, role := range []string{"header", "body"} {
				if e := region(role); e != nil {
					return nil, e
				}
			}
			rows := [][]*parser.Instance{}
			for _, row := range n.Rows {
				items := []*parser.Instance{}
				for _, c := range row {
					if visible(c) {
						items = append(items, c)
					}
				}
				if len(items) > 0 {
					rows = append(rows, items)
				}
			}
			for i, row := range rows {
				if len(rows) > 1 || len(row) > 1 {
					noun := "component"
					if len(row) != 1 {
						noun += "s"
					}
					lines = append(lines, fmt.Sprintf("**Row %d · %d %s from left to right**", i+1, len(row), noun), "")
				}
				for _, c := range row {
					block, e := render(c)
					if e != nil {
						return nil, e
					}
					lines = append(lines, block...)
					lines = append(lines, "", "---", "")
				}
			}
			if e := region("footer"); e != nil {
				return nil, e
			}
			lines = quoted(lines)
		}
		for _, line := range lines {
			budget += len(line) + 1
		}
		if budget > MaxCells {
			return nil, diagnostic("markdown-limit", "Markdown exceeds output budget", n)
		}
		return lines, nil
	}
	content, e := render(root)
	if e != nil {
		return "", e
	}
	lines := []string{"# SDUI — " + codeSpan(root.Path), "", "Static GUI dump. Buttons and fields are text labels; no callbacks execute.", "", "## Layout overview", "", "Row/column structure in terminal cells. Heights follow content; this is not measured GUI geometry.", "", fence + "text", overview, fence, "", "## Content", "", "Markdown is rendered as content. Nested blockquotes represent groups and frames. Horizontal siblings appear in reading order here; the overview shows placement. Mermaid diagrams are omitted.", ""}
	lines = append(lines, content...)
	text := strings.TrimRight(strings.Join(lines, "\n"), " \t\r\n") + "\n"
	if len(text) > MaxCells {
		return "", diagnostic("markdown-limit", "Markdown exceeds output budget", root)
	}
	return text, nil
}
