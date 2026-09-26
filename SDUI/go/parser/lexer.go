package parser

import (
	"math"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

const MaxBytes = 262144
const MaxTokens = 50000

var identifierRE = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*(?:-[A-Za-z_][A-Za-z0-9_]*)*`)
var numberRE = regexp.MustCompile(`^-?(?:0|[1-9][0-9]*)(?:\.[0-9]+)?(?:[eE][+-]?[0-9]+)?`)

type token struct {
	kind  string
	value any
	span  Span
}

func lex(source string) []token {
	origin := Span{0, 0, 1, 1}
	if !utf8.ValidString(source) {
		fail("encoding", "Source must be UTF-8", origin)
	}
	if len(source) > MaxBytes {
		fail("source-limit", "Source exceeds 256 KiB", origin)
	}
	positions := make([]Span, len(source)+1)
	line, col := 1, 1
	for i, r := range source {
		positions[i] = Span{i, i, line, col}
		if r == '\n' {
			line++
			col = 1
		} else {
			col++
		}
	}
	positions[len(source)] = Span{len(source), len(source), line, col}
	span := func(a, b int) Span { s := positions[a]; s.End = b; return s }
	tokens := []token{}
	i := 0
	for i < len(source) {
		c := source[i]
		if strings.ContainsRune(" \t\r\n", rune(c)) {
			i++
			continue
		}
		if c == '#' {
			for i < len(source) && source[i] != '\r' && source[i] != '\n' {
				_, n := utf8.DecodeRuneInString(source[i:])
				i += n
			}
			continue
		}
		start := i
		kind := ""
		var value any
		switch {
		case strings.HasPrefix(source[i:], `"""`):
			end := strings.Index(source[i+3:], `"""`)
			if end < 0 {
				fail("lexical", "Unterminated multiline string", span(start, len(source)))
			}
			end += i + 3
			raw := source[i+3 : end]
			for _, r := range raw {
				if r < 32 && r != '\n' && r != '\r' && r != '\t' || r == 127 {
					fail("lexical", "Control character in string", span(start, end))
				}
			}
			kind = "STRING"
			value = raw
			i = end + 3
		case c == '"' || c == '\'':
			quote := c
			i++
			var out strings.Builder
			for i < len(source) && source[i] != quote {
				r, n := utf8.DecodeRuneInString(source[i:])
				if r < 32 || r == 127 {
					fail("lexical", "Raw control character in string", span(i, i+n))
				}
				if r != '\\' {
					out.WriteRune(r)
					i += n
					continue
				}
				i++
				if i == len(source) {
					fail("lexical", "Unfinished escape", span(start, i))
				}
				escapes := map[byte]rune{'n': '\n', 'r': '\r', 't': '\t', '\\': '\\', '"': '"', '\'': '\''}
				if source[i] == 'u' {
					if i+5 > len(source) {
						fail("lexical", "Expected four hex digits", span(i, len(source)))
					}
					raw := source[i+1 : i+5]
					for _, h := range raw {
						if !strings.ContainsRune("0123456789abcdefABCDEF", h) {
							fail("lexical", "Invalid Unicode escape", span(i, i+5))
						}
					}
					v, e := strconv.ParseUint(raw, 16, 16)
					if e != nil || v == 0 || v >= 0xD800 && v <= 0xDFFF {
						fail("lexical", "Invalid Unicode scalar escape", span(i, i+5))
					}
					out.WriteRune(rune(v))
					i += 5
				} else if r, ok := escapes[source[i]]; ok {
					out.WriteRune(r)
					i++
				} else {
					fail("lexical", "Unknown escape", span(i-1, i+1))
				}
			}
			if i == len(source) {
				fail("lexical", "Unterminated string", span(start, i))
			}
			i++
			kind = "STRING"
			value = out.String()
		default:
			if raw := identifierRE.FindString(source[i:]); raw != "" {
				kind = "ID"
				value = raw
				i += len(raw)
			} else if raw := numberRE.FindString(source[i:]); raw != "" {
				v, e := strconv.ParseFloat(raw, 64)
				i += len(raw)
				if e != nil || math.IsInf(v, 0) || math.IsNaN(v) {
					fail("lexical", "Number must be finite", span(start, i))
				}
				kind = "NUMBER"
				value = v
			} else if strings.ContainsRune("[]{}(),;=.:@<>*^|-", rune(c)) {
				kind = string(c)
				value = kind
				i++
			} else {
				_, n := utf8.DecodeRuneInString(source[i:])
				fail("lexical", "Unexpected character", span(i, i+n))
			}
		}
		tokens = append(tokens, token{kind, value, span(start, i)})
		if len(tokens) > MaxTokens {
			fail("token-limit", "More than 50000 tokens", span(start, i))
		}
	}
	return append(tokens, token{"EOF", "", span(i, i)})
}
