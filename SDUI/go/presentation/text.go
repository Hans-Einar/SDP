// Package presentation exports static previews without executing callbacks.
package presentation

import (
	"fmt"
	"golang.org/x/text/unicode/norm"
	"golang.org/x/text/width"
	"regexp"
	"strings"
	"unicode"
)

var fencePattern = regexp.MustCompile("^ {0,3}(`{3,}|~{3,})([^\\r\\n]*)$")

func splitLines(s string) []string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	if strings.HasSuffix(s, "\n") {
		s = strings.TrimSuffix(s, "\n")
	}
	return strings.Split(s, "\n")
}
func closing(line, mark string) bool {
	s := strings.TrimLeft(line, " ")
	if len(line)-len(s) > 3 {
		return false
	}
	s = strings.TrimRight(s, " \t")
	return len(s) >= len(mark) && strings.Trim(s, string(mark[0])) == ""
}
func MarkdownLines(s string) []string {
	out := []string{}
	fence := ""
	hidden := false
	for _, line := range splitLines(s) {
		m := fencePattern.FindStringSubmatch(line)
		if fence != "" {
			if !hidden {
				out = append(out, line)
			}
			if closing(line, fence) {
				fence = ""
				hidden = false
			}
		} else if m != nil && !(m[1][0] == '`' && strings.Contains(m[2], "`")) {
			fence = m[1]
			info := strings.Fields(m[2])
			hidden = len(info) > 0 && strings.EqualFold(info[0], "mermaid")
			if hidden {
				out = append(out, "[Mermaid utelatt]")
			} else {
				out = append(out, line)
			}
		} else {
			out = append(out, line)
		}
	}
	if len(out) == 0 {
		return []string{""}
	}
	return out
}
func Safe(s string) string {
	var out strings.Builder
	for _, r := range s {
		if r == '\t' {
			out.WriteString("    ")
		} else if unicode.Is(unicode.C, r) {
			fmt.Fprintf(&out, "\\u%04x", r)
		} else {
			out.WriteRune(r)
		}
	}
	return out.String()
}
func runeWidth(r rune) int {
	if norm.NFD.PropertiesString(string(r)).CCC() != 0 {
		return 0
	}
	kind := width.LookupRune(r).Kind()
	if kind == width.EastAsianWide || kind == width.EastAsianFullwidth {
		return 2
	}
	return 1
}
func CellWidth(s string) int {
	n := 0
	for _, r := range s {
		n += runeWidth(r)
	}
	return n
}
func wrapLine(s string, w int) []string {
	out := []string{}
	var line strings.Builder
	used := 0
	for _, r := range Safe(s) {
		size := runeWidth(r)
		if used+size > w {
			out = append(out, line.String())
			line.Reset()
			used = 0
		}
		line.WriteRune(r)
		used += size
	}
	return append(out, line.String())
}
func fit(s string, w int) string {
	line := wrapLine(s, w)[0]
	return line + strings.Repeat(" ", max(0, w-CellWidth(line)))
}
func codeSpan(s string) string {
	ticks := strings.Repeat("`", maxRun(s, '`')+1)
	return ticks + " " + strings.ReplaceAll(Safe(s), "\n", " ") + " " + ticks
}
func maxRun(s string, c rune) int {
	maximum, n := 0, 0
	for _, r := range s {
		if r == c {
			n++
			maximum = max(maximum, n)
		} else {
			n = 0
		}
	}
	return maximum
}
func contentLines(s string) []string {
	lines := MarkdownLines(s)
	fence := ""
	for i, line := range lines {
		lines[i] = Safe(line)
		if fence != "" {
			if closing(lines[i], fence) {
				fence = ""
			}
		} else if m := fencePattern.FindStringSubmatch(lines[i]); m != nil && !(m[1][0] == '`' && strings.Contains(m[2], "`")) {
			fence = m[1]
		}
	}
	if fence != "" {
		lines = append(lines, fence)
	}
	return lines
}
