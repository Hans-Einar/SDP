package markdown

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"
)

type previewInput struct {
	outcome     PreviewOutcome
	text        string
	fingerprint string
}

func reason(e error) string {
	s := strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return ' '
		}
		return r
	}, strings.ToValidUTF8(e.Error(), "�"))
	if len(s) > 4096 {
		s = s[:4096]
		for !utf8.ValidString(s) {
			s = s[:len(s)-1]
		}
	}
	return s
}
func previewError(o PreviewOutcome, code, message string) error {
	return &parser.Diagnostic{Code: code, Span: o.Span, Message: fmt.Sprintf("%s (provider %q, source %s.%s.@%s): %s", o.Path, o.ProviderID, o.Source.Module, o.Source.Object, o.Source.Member, reason(fmt.Errorf("%s", message)))}
}
func identity(s string) bool { return utf8.ValidString(s) && strings.TrimSpace(s) != "" }
func digest(b []byte) string { return fmt.Sprintf("%x", sha256.Sum256(b)) }
func validDigest(s string) bool {
	if len(s) != 64 {
		return false
	}
	for _, c := range s {
		if !(c >= '0' && c <= '9' || c >= 'a' && c <= 'f') {
			return false
		}
	}
	return true
}
func sameReference(a, b parser.Reference) bool {
	return a.Module == b.Module && a.Object == b.Object && a.Member == b.Member
}
func nilRenderer(r Renderer) bool {
	if r == nil {
		return true
	}
	v := reflect.ValueOf(r)
	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return v.IsNil()
	}
	return false
}
func diagramSource(s string) error {
	f := strings.Fields(s)
	if len(s) > 12000 || len(f) == 0 || (f[0] != "flowchart" && f[0] != "graph") || strings.Contains(s, "%%{") {
		return fmt.Errorf("mermaid-profile: only bounded flowchart/graph without configuration")
	}
	return nil
}

// Syntax-only inventory deliberately does not stop on unsupported HTML/images.
func hasDiagram(source string) bool {
	raw := []byte(source)
	tree := goldmark.New().Parser().Parse(text.NewReader(raw))
	found := false
	_ = ast.Walk(tree, func(n ast.Node, enter bool) (ast.WalkStatus, error) {
		if enter {
			if f, ok := n.(*ast.FencedCodeBlock); ok && string(f.Language(raw)) == "mermaid" {
				found = true
			}
		}
		return ast.WalkContinue, nil
	})
	return found
}
func previewInventory(root *parser.Instance) ([]previewInput, error) {
	var out []previewInput
	var failure error
	root.Walk(func(n *parser.Instance) {
		if failure != nil || !(n.Kind == "markdown" || n.Widget == "svg" || n.Widget == "markdown") {
			return
		}
		policy, e := parser.PreviewOptions(n)
		if e != nil {
			failure = e
			return
		}
		if !policy.Explicit {
			return
		}
		kind := "markdown"
		source := parser.Reference{}
		if n.Widget == "svg" {
			kind = "svg"
			source = n.Arguments["source"].(parser.Reference)
		}
		o := PreviewOutcome{Path: n.Path, Kind: kind, Description: policy.Description, Fallback: policy.Fallback, Span: n.Span, Uses: append([]parser.UseSite(nil), n.Uses...), Source: source}
		key := []string{n.Profile, n.Path, kind, n.Text, source.Module, source.Object, source.Member, policy.Description, policy.Fallback}
		data, _ := json.Marshal(key)
		out = append(out, previewInput{o, n.Text, digest(data)})
	})
	return out, failure
}
