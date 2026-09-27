package install

import (
	"net/url"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var referenceLink = regexp.MustCompile(`(?m)^ {0,3}\[[^\]\r\n]+\]:[ \t]*`)

func rebaseMarkdown(b []byte, source, dest string, moves map[string]string, mappings []Move) []byte {
	if len(moves) == 0 {
		return b
	}
	s := string(b)
	type span struct{ a, b int }
	spans := []span{}
	readURL := func(start int) span {
		for start < len(s) && (s[start] == ' ' || s[start] == '\t') {
			start++
		}
		if start >= len(s) {
			return span{start, start}
		}
		if s[start] == '<' {
			end := strings.IndexByte(s[start+1:], '>')
			if end < 0 {
				return span{start, start}
			}
			return span{start + 1, start + 1 + end}
		}
		end := start
		depth := 0
		for end < len(s) {
			c := s[end]
			if c == '\\' && end+1 < len(s) {
				end += 2
				continue
			}
			if c == '(' {
				depth++
			}
			if c == ')' {
				if depth == 0 {
					break
				}
				depth--
			}
			if c == ' ' || c == '\t' || c == '\r' || c == '\n' {
				break
			}
			end++
		}
		return span{start, end}
	}
	for offset := 0; offset < len(s); {
		i := strings.Index(s[offset:], "](")
		if i < 0 {
			break
		}
		pos := offset + i + 2
		r := readURL(pos)
		if r.b > r.a {
			spans = append(spans, r)
		}
		offset = pos
	}
	for _, m := range referenceLink.FindAllStringIndex(s, -1) {
		r := readURL(m[1])
		if r.b > r.a {
			spans = append(spans, r)
		}
	}
	replacements := map[int]span{}
	for _, sp := range spans {
		replacements[sp.a] = sp
	}
	var out strings.Builder
	for i := 0; i < len(s); {
		sp, ok := replacements[i]
		if !ok {
			out.WriteByte(s[i])
			i++
			continue
		}
		raw := s[sp.a:sp.b]
		replacement := raw
		u, e := url.Parse(strings.ReplaceAll(strings.ReplaceAll(raw, `\(`, "("), `\)`, ")"))
		if e == nil && u.Scheme == "" && u.Host == "" && u.Path != "" && !strings.HasPrefix(u.Path, "/") {
			target := path.Clean(path.Join(path.Dir(source), u.Path))
			newTarget := target
			if v, ok := moves[target]; ok {
				newTarget = v
			} else {
				for _, m := range mappings {
					if target == m.From || strings.HasPrefix(target, m.From+"/") {
						newTarget = m.To + strings.TrimPrefix(target, m.From)
						break
					}
				}
			}
			if source != dest || target != newTarget {
				rel, e := filepath.Rel(filepath.FromSlash(path.Dir(dest)), filepath.FromSlash(newTarget))
				if e == nil {
					parts := strings.Split(filepath.ToSlash(rel), "/")
					for j, p := range parts {
						parts[j] = url.PathEscape(p)
					}
					replacement = strings.Join(parts, "/")
					if u.RawQuery != "" {
						replacement += "?" + u.RawQuery
					}
					if u.Fragment != "" {
						replacement += "#" + u.EscapedFragment()
					}
				}
			}
		}
		out.WriteString(replacement)
		i = sp.b
	}
	return []byte(out.String())
}
