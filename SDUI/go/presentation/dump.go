package presentation

import (
	"fmt"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"math"
	"sort"
	"strconv"
	"strings"
)

const MaxCells = 2000000

func diagnostic(code, msg string, n *parser.Instance) error {
	return &parser.Diagnostic{Code: code, Message: msg, Span: n.Span}
}
func visible(n *parser.Instance) bool { return n != nil && n.Layout["visible"] != false }
func widths(items []*parser.Instance, available int) ([]int, error) {
	sizes := make([]int, len(items))
	weights := make([]float64, len(items))
	flex := []int{}
	remaining := available
	maxWeight := 0.0
	for i, n := range items {
		scale, ok := n.Layout["scale-x"].(float64)
		if !ok {
			scale, ok = n.Layout["scale"].(float64)
		}
		if ok {
			raw := float64(available) * scale
			if raw > float64(available) || math.IsInf(raw, 0) {
				return nil, fmt.Errorf("Relative widths exceed available cells")
			}
			sizes[i] = max(3, int(raw))
			remaining -= sizes[i]
		} else {
			flex = append(flex, i)
		}
		weights[i] = 1
		if x, ok := n.Layout["x"].(string); ok && strings.HasSuffix(x, "fr") {
			weights[i], _ = strconv.ParseFloat(strings.TrimSuffix(x, "fr"), 64)
		}
		if sizes[i] == 0 {
			maxWeight = max(maxWeight, weights[i])
		}
	}
	if remaining < 3*len(flex) {
		return nil, fmt.Errorf("Row needs more columns")
	}
	if len(flex) > 0 {
		total := 0.0
		for _, i := range flex {
			total += weights[i] / maxWeight
		}
		extra := remaining - 3*len(flex)
		fractions := map[int]float64{}
		used := 0
		for _, i := range flex {
			fraction := float64(extra) * (weights[i] / maxWeight) / total
			sizes[i] = 3 + int(math.Floor(fraction))
			fractions[i] = fraction - math.Floor(fraction)
		}
		for _, w := range sizes {
			used += w
		}
		sort.SliceStable(flex, func(a, b int) bool { return fractions[flex[a]] > fractions[flex[b]] })
		for j := 0; j < available-used && j < len(flex); j++ {
			sizes[flex[j]]++
		}
	}
	return sizes, nil
}
func widgetText(n *parser.Instance) string {
	switch n.Widget {
	case "button":
		return "[ " + n.Argument("label") + " ]"
	case "input":
		return n.Argument("text") + ": [" + n.Argument("value") + "]"
	default:
		label := n.Argument("label")
		if label == "" {
			label = n.Path
		}
		return "[SVG plassholder: " + label + "]"
	}
}

// Dump produces a bounded terminal-cell structural preview, not measured GUI geometry.
func Dump(root *parser.Instance, columns int) (string, error) {
	if columns < 20 || columns > 400 {
		return "", diagnostic("dump-width", "Columns must be between 20 and 400", root)
	}
	budget := 0
	var render func(*parser.Instance, int) ([]string, error)
	render = func(n *parser.Instance, w int) ([]string, error) {
		if !visible(n) {
			return nil, nil
		}
		if w < 3 {
			return nil, diagnostic("dump-space", "Not enough columns", n)
		}
		lines := []string{}
		if n.Kind == "markdown" {
			for _, line := range MarkdownLines(n.Text) {
				lines = append(lines, wrapLine(line, w)...)
			}
		} else if n.Kind == "widget" {
			lines = wrapLine(widgetText(n), w)
		} else {
			boxed := n.Kind == "frame" && n.Variant == "box"
			inner := w
			if boxed {
				inner -= 2
			}
			for _, role := range []string{"header", "body"} {
				if r := n.Region(role); r != nil {
					block, e := render(r, inner)
					if e != nil {
						return nil, e
					}
					lines = append(lines, block...)
				}
			}
			for _, row := range n.Rows {
				items := []*parser.Instance{}
				for _, c := range row {
					if visible(c) {
						items = append(items, c)
					}
				}
				if len(items) == 0 {
					continue
				}
				ws, e := widths(items, inner-len(items)+1)
				if e != nil {
					return nil, diagnostic("dump-space", e.Error(), n)
				}
				blocks := make([][]string, len(items))
				height := 0
				for i, c := range items {
					blocks[i], e = render(c, ws[i])
					if e != nil {
						return nil, e
					}
					height = max(height, len(blocks[i]))
				}
				for y := 0; y < height; y++ {
					parts := []string{}
					for i, b := range blocks {
						text := ""
						if y < len(b) {
							text = b[y]
						}
						parts = append(parts, fit(text, ws[i]))
					}
					lines = append(lines, strings.Join(parts, " "))
				}
			}
			if r := n.Region("footer"); r != nil {
				block, e := render(r, inner)
				if e != nil {
					return nil, e
				}
				lines = append(lines, block...)
			}
			if boxed {
				parts := strings.Split(n.Path, "/")
				label := " " + parts[len(parts)-1] + " "
				if v, ok := n.Layout["y"]; ok {
					label += fmt.Sprint(v) + " "
				}
				out := []string{"+" + strings.ReplaceAll(fit(label, inner), " ", "-") + "+"}
				for _, line := range lines {
					out = append(out, "|"+fit(line, inner)+"|")
				}
				lines = append(out, "+"+strings.Repeat("-", inner)+"+")
			}
		}
		for _, line := range lines {
			budget += CellWidth(line)
		}
		if budget > MaxCells {
			return nil, diagnostic("dump-limit", "Dump exceeds cell budget", n)
		}
		return lines, nil
	}
	lines, e := render(root, columns)
	if e != nil {
		return "", e
	}
	for i := range lines {
		lines[i] = strings.TrimRight(lines[i], " \t\r\n")
	}
	return fmt.Sprintf("SDUI GUI dump | %s | %d columns | structural preview\n", Safe(root.Path), columns) + strings.Join(lines, "\n") + "\n", nil
}
