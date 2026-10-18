package parser

import (
	"strings"
	"unicode/utf8"
)

// PreviewPolicy is source intent only, never prepared provider readiness.
type PreviewPolicy struct {
	Explicit              bool
	Description, Fallback string
}

var previews = map[string]widgetSchema{
	"svg":      {"source", map[string]string{"source": "reference", "label": "string", "description": "string", "fallback": "string"}, "source"},
	"markdown": {"text", map[string]string{"text": "string", "description": "string", "fallback": "string"}, "text"},
}

// PreviewOptions validates normalized SVG/Markdown source facts without I/O or mutation.
// Bare Markdown and SVG without either policy argument remain legacy content.
func PreviewOptions(n *Instance) (PreviewPolicy, error) {
	var out PreviewPolicy
	bad := func(msg string, span Span) (PreviewPolicy, error) {
		if n != nil {
			msg = n.Path + ": " + msg
		}
		return PreviewPolicy{}, &Diagnostic{Code: "preview-argument", Message: msg, Span: span}
	}
	if n == nil {
		return bad("Expected SVG or Markdown", Span{})
	}
	profile, err := EffectiveProfile(n)
	if err != nil {
		return out, err
	}
	markdown := n.Kind == "markdown" && n.Widget == ""
	if !markdown && !(n.Kind == "widget" && n.Widget == "svg") || len(n.Rows) > 0 || len(n.Regions) > 0 {
		return bad("Expected leaf SVG or normalized Markdown", n.Span)
	}
	if markdown && len(n.Arguments) == 0 {
		return out, nil
	}
	out.Explicit = markdown || hasArg(n, "description") || hasArg(n, "fallback")
	if out.Explicit && profile != "sdui/0.3" {
		return bad("Explicit previews require sdui/0.3", n.Span)
	}
	kind := "svg"
	if markdown {
		kind = "markdown"
	}
	schema := previews[kind]
	if !out.Explicit {
		schema = widgets["svg"]
	}
	for key, arg := range n.Arguments {
		span := n.Span
		switch v := arg.(type) {
		case Literal:
			span = v.Span
		case Reference:
			span = v.Span
		}
		expected, ok := schema.fields[key]
		if !ok {
			return bad("Unknown preview argument "+key, span)
		}
		switch v := arg.(type) {
		case Literal:
			_, ok := v.Value.(string)
			if expected != "string" || v.Kind != "string" || !ok {
				return bad("Expected string "+key, span)
			}
		case Reference:
			if expected != "reference" {
				return bad("Unexpected reference "+key, span)
			}
			if out.Explicit && (!previewIdentifier(v.Module) || !previewIdentifier(v.Object) || v.Member != "resource") {
				return bad("Expected module.object.@resource", span)
			}
		default:
			return bad("Invalid argument "+key, span)
		}
	}
	if !hasArg(n, schema.required) {
		return bad("Missing "+schema.required, n.Span)
	}
	if !out.Explicit {
		return out, nil
	}
	for _, key := range []string{"description", "fallback"} {
		if !hasArg(n, key) {
			return bad("Missing "+key, n.Span)
		}
	}
	out.Description = n.Argument("description")
	out.Fallback = n.Argument("fallback")
	if !utf8.ValidString(out.Description) || len(out.Description) > 4096 || strings.TrimSpace(out.Description) == "" {
		return bad("Description must be nonblank UTF-8, at most 4096 bytes", n.Arguments["description"].(Literal).Span)
	}
	if out.Fallback != "label" && out.Fallback != "reject" {
		return bad("Fallback must be label or reject", n.Arguments["fallback"].(Literal).Span)
	}
	if markdown && (n.Text != n.Argument("text") || !utf8.ValidString(n.Text) || len(n.Text) > 32768) {
		return bad("Markdown text must match its argument and be UTF-8, at most 32768 bytes", n.Arguments["text"].(Literal).Span)
	}
	return out, nil
}

func previewIdentifier(s string) bool {
	if s == "" {
		return false
	}
	for i, r := range s {
		if !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r == '_' || i > 0 && r >= '0' && r <= '9') {
			return false
		}
	}
	return !member(s, "sdui ref true false null setHandle")
}

func sourceKind(n *Node, profile string) string {
	if profile == "sdui/0.3" && n.Kind == "widget" && val(n.Widget) == "markdown" {
		return "markdown"
	}
	return n.Kind
}

func validatePreviewSource(n *Node, profile string) {
	if profile != "sdui/0.3" || !member(val(n.Widget), "svg markdown") {
		return
	}
	if n.Kind != "widget" || len(n.Rows) > 0 {
		fail("preview-argument", "Preview calls are leaves", n.Span)
	}
	i := &Instance{Profile: profile, Kind: n.Kind, Widget: val(n.Widget), Arguments: widgetArguments(n, profile), Span: n.Span}
	if i.Widget == "markdown" {
		i.Kind = "markdown"
		i.Widget = ""
		i.Text = i.Argument("text")
	}
	if _, err := PreviewOptions(i); err != nil {
		panic(err)
	}
}
