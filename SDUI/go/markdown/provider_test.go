package markdown

import (
	"os"
	"strings"
	"testing"

	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/svg"
)

func TestProfileAndResources(t *testing.T) {
	source := "# Heading\n\n**Bold** and `code`\n\n- item\n\n| A | B |\n|---|---|\n| 1 | 2 |\n\n```mermaid\nflowchart LR\n A-->B\n```"
	d, e := Parse(source)
	if e != nil {
		t.Fatal(e)
	}
	if len(d.Diagrams) != 1 || len(d.Blocks) != 5 {
		t.Fatalf("document %#v", d)
	}
	n := &parser.Instance{Kind: "markdown", Path: "page/text", Text: source, Layout: map[string]any{}}
	p, e := Prepare(n, nil)
	if e != nil {
		t.Fatal(e)
	}
	s, e := p.Measure(n, 14, 400)
	if e != nil || s.H <= 0 {
		t.Fatal(s, e)
	}
	b := &layout.Box{Instance: n, Path: n.Path, Rect: layout.Rect{W: 400, H: s.H}, Clip: layout.Rect{W: 400, H: s.H}, Font: 14}
	out, e := svg.Render(b, svg.Options{Width: 400, Height: s.H, Content: p})
	if e != nil {
		t.Fatal(e)
	}
	if strings.Contains(out, "**Bold**") || !strings.Contains(out, "Mermaid:") {
		t.Fatal("raw Markdown or missing resource status")
	}
	for _, invalid := range []string{"<script>alert(1)</script>", "# ![remote](https://example.com/a.png)", strings.Repeat("a", 32769)} {
		if _, e := Parse(invalid); e == nil {
			t.Fatalf("accepted unsupported %q", invalid[:min(len(invalid), 40)])
		}
	}
	if _, e := (Mmdr{}).Render("usecase-beta\n A"); e == nil {
		t.Fatal("unsupported diagram accepted")
	}
	for _, invalid := range []string{`<svg viewBox="0 0 10 10"><script/></svg>`, `<svg viewBox="0 0 10 10"><use href="https://remote"/></svg>`, `<svg viewBox="0 0 10 10"><foreignObject/></svg>`} {
		if _, e := validateResource([]byte(invalid)); e == nil {
			t.Fatal("active resource accepted")
		}
	}
}
func TestConfiguredRenderer(t *testing.T) {
	path := os.Getenv("SDUI_MMDR")
	if path == "" {
		t.Skip("set SDUI_MMDR for registered renderer acceptance")
	}
	n := &parser.Instance{Kind: "markdown", Text: "```mermaid\nflowchart LR\n Alpha-->Beta\n```"}
	p, e := Prepare(n, Mmdr{path})
	if e != nil {
		t.Fatal(e)
	}
	if len(p.Resources) != 1 {
		t.Fatal("missing resource")
	}
	for _, r := range p.Resources {
		if !strings.Contains(string(r.SVG), "Alpha") || !strings.Contains(string(r.SVG), "Beta") || !strings.Contains(string(r.SVG), "<rect") {
			t.Fatal("renderer lost labels or node geometry")
		}
	}
	dir := t.TempDir()
	if e = p.WriteResources(dir); e != nil {
		t.Fatal(e)
	}
	if e = p.WriteResources(dir); e != nil {
		t.Fatal(e)
	}
	files, e := os.ReadDir(dir)
	if e != nil || len(files) != 1 {
		t.Fatal(files, e)
	}
}
