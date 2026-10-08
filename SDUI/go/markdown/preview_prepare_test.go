package markdown

import (
	"fmt"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"strings"
	"testing"
)

type previewRendererFunc func(string) (Resource, error)

func (f previewRendererFunc) Render(s string) (Resource, error) { return f(s) }
func previewRoot(t *testing.T, body string) *parser.Instance {
	t.Helper()
	d, e := parser.Parse(`sdui 0.3; ref: art "unopened"; Main=[` + body + `];`)
	if e != nil {
		t.Fatal(e)
	}
	roots, e := parser.Normalize(d)
	if e != nil {
		t.Fatal(e)
	}
	return roots["Main"]
}
func svgNodeSource(name, policy string) string {
	return name + `=svg(art.Chart.@resource,description="Resource description",fallback="` + policy + `")`
}
func mdNodeSource(name, source, description, policy string) string {
	return fmt.Sprintf(`%s=markdown(%q,description=%q,fallback=%q)`, name, source, description, policy)
}
func bindSVG(root *parser.Instance, r Resource) map[string]PreparedSVG {
	m := map[string]PreparedSVG{}
	root.Walk(func(n *parser.Instance) {
		if n.Widget == "svg" {
			m[n.Path] = PreparedSVG{Source: n.Arguments["source"].(parser.Reference), ProviderID: "art-v1", SHA256: digest(r.SVG), Resource: r}
		}
	})
	return m
}
func fullBackend() PreviewBackend {
	return PreviewBackend{SVG: func(Resource) error { return nil }, Markdown: func() error { return nil }, Mermaid: func(Resource) error { return nil }}
}

const diagramMD = "Before\n\n```mermaid\nflowchart LR\n A-->B\n```\n\nAfter"

func TestPreviewIdentityBeforeFallbackAndCallbacks(t *testing.T) {
	root := previewRoot(t, svgNodeSource("a", "label")+";"+mdNodeSource("b", diagramMD, "diagram", "label"))
	calls := 0
	render := MarkdownRenderer{"m", "r", previewRendererFunc(func(string) (Resource, error) { calls++; return shapeResource(`<rect/>`), nil })}
	for _, mutate := range []func(map[string]PreparedSVG, map[string]MarkdownRenderer){
		func(r map[string]PreparedSVG, m map[string]MarkdownRenderer) {
			v := r["Main/a"]
			v.SHA256 = strings.Repeat("0", 64)
			r["Main/a"] = v
		},
		func(r map[string]PreparedSVG, m map[string]MarkdownRenderer) {
			v := r["Main/a"]
			v.Source.Object = "Other"
			r["Main/a"] = v
		},
		func(r map[string]PreparedSVG, m map[string]MarkdownRenderer) { r["extra"] = r["Main/a"] },
		func(r map[string]PreparedSVG, m map[string]MarkdownRenderer) { m["Main/a"] = render },
		func(r map[string]PreparedSVG, m map[string]MarkdownRenderer) {
			m["Main/b"] = MarkdownRenderer{"", "r", render.Renderer}
		},
		func(r map[string]PreparedSVG, m map[string]MarkdownRenderer) {
			m["Main/b"] = MarkdownRenderer{"m", "r", nil}
		},
	} {
		resources := bindSVG(root, shapeResource(`<script/>`))
		renderers := map[string]MarkdownRenderer{"Main/b": render}
		mutate(resources, renderers)
		if _, e := PreparePreviews(root, resources, renderers, PreviewBackend{}); e == nil || !strings.Contains(e.Error(), "preview-identity") {
			t.Fatalf("identity error: %v", e)
		}
		if calls != 0 {
			t.Fatal("renderer ran before identity pass")
		}
	}
	root = previewRoot(t, mdNodeSource("b", "<b>x</b>", "bad", "label"))
	if _, e := PreparePreviews(root, nil, map[string]MarkdownRenderer{"Main/b": render}, fullBackend()); e == nil {
		t.Fatal("unused renderer hidden by content")
	}
	root = previewRoot(t, mdNodeSource("b", "<b>x</b>\n\n"+diagramMD, "bad", "label"))
	p, e := PreparePreviews(root, nil, map[string]MarkdownRenderer{"Main/b": render}, fullBackend())
	if e != nil {
		t.Fatal(e)
	}
	o, _ := p.Outcome("Main/b")
	if o.Status != "label" || calls != 0 {
		t.Fatal(o, calls)
	}
}
func TestPreviewCopiesPerPathAndZeroRerenders(t *testing.T) {
	root := previewRoot(t, svgNodeSource("image", "reject")+";"+mdNodeSource("one", diagramMD, "first", "label")+";"+mdNodeSource("two", diagramMD, "second", "label"))
	data := shapeResource(`<rect width="30" height="15"/>`)
	resources := bindSVG(root, data)
	scratch := shapeResource(`<circle r="4"/>`)
	calls := 0
	renderer := previewRendererFunc(func(s string) (Resource, error) {
		calls++
		if !strings.HasPrefix(s, "flowchart") {
			t.Fatal("full Markdown supplied")
		}
		return scratch, nil
	})
	renderers := map[string]MarkdownRenderer{"Main/one": {"a", "v1", renderer}, "Main/two": {"b", "v2", previewRendererFunc(func(string) (Resource, error) { calls++; return Resource{}, fmt.Errorf("unavailable") })}}
	backend := fullBackend()
	backend.SVG = func(r Resource) error { r.SVG[0] = '!'; data.SVG[0] = '!'; return nil }
	backend.Mermaid = func(r Resource) error { r.SVG[0] = '!'; return nil }
	p, e := PreparePreviews(root, resources, renderers, backend)
	if e != nil {
		t.Fatal(e)
	}
	scratch.SVG[0] = '!'
	first, _ := p.Outcome("Main/one")
	second, _ := p.Outcome("Main/two")
	img, _ := p.Outcome("Main/image")
	if first.Status != "rendered" || second.Status != "partial" || first.Description == second.Description || img.Resource.SVG[0] != '<' || first.Diagrams[0].Resource.SVG[0] != '<' {
		t.Fatal(first, second, img)
	}
	first.Document.Blocks[0].Text = "mutated"
	first.Diagrams[0].Resource.SVG[0] = '!'
	img.Resource.SVG[0] = '!'
	again, _ := p.Outcome("Main/one")
	if again.Document.Blocks[0].Text == "mutated" || again.Diagrams[0].Resource.SVG[0] != '<' {
		t.Fatal("mutable accessor")
	}
	for range 5 {
		if e = p.Check(root); e != nil {
			t.Fatal(e)
		}
		for _, path := range []string{"Main/one", "Main/two"} {
			var n *parser.Instance
			root.Walk(func(v *parser.Instance) {
				if v.Path == path {
					n = v
				}
			})
			size, e := p.Measure(n, 12, 200)
			if e != nil {
				t.Fatal(e)
			}
			var out strings.Builder
			if e = p.Render(&out, &layout.Box{Instance: n, Path: path, Font: 12, Rect: layout.Rect{W: 200, H: size.H}}); e != nil {
				t.Fatal(e)
			}
			if path == "Main/two" && (!strings.Contains(out.String(), "Preview unavailable") || !strings.Contains(out.String(), "Before") || !strings.Contains(out.String(), "After")) {
				t.Fatal(out.String())
			}
		}
	}
	if calls != 2 {
		t.Fatal("rerendered", calls)
	}
}
func TestPreviewFingerprintAndPolicy(t *testing.T) {
	root := previewRoot(t, svgNodeSource("a", "label"))
	p, e := PreparePreviews(root, nil, nil, fullBackend())
	if e != nil {
		t.Fatal(e)
	}
	n := root.Rows[0][0]
	n.Arguments["label"] = parser.Literal{Kind: "string", Value: "caption"}
	n.Layout["visible"] = false
	if e = p.Check(root); e != nil {
		t.Fatal("caption invalidated resource", e)
	}
	n.Arguments["description"] = parser.Literal{Kind: "string", Value: "replacement"}
	if e = p.Check(root); e == nil {
		t.Fatal("changed description accepted")
	}
	if _, e = PreparePreviews(previewRoot(t, svgNodeSource("a", "reject")), nil, nil, fullBackend()); e == nil {
		t.Fatal("missing reject accepted")
	}
	root = previewRoot(t, mdNodeSource("a", diagramMD, "diagram", "label"))
	count := 0
	r := shapeResource(`<rect/>`)
	p, e = PreparePreviews(root, nil, map[string]MarkdownRenderer{"Main/a": {"p", "r", previewRendererFunc(func(string) (Resource, error) { count++; return r, nil })}}, PreviewBackend{Markdown: func() error { return nil }})
	if e != nil {
		t.Fatal(e)
	}
	o, _ := p.Outcome("Main/a")
	if count != 1 || o.Status != "partial" || o.Diagrams[0].Resource != nil || o.Diagrams[0].SHA256 != digest(r.SVG) {
		t.Fatal(o, count)
	}
}

func TestPreviewRepeatedDiagramOccurrencesAndDetachedProvenance(t *testing.T) {
	source := "```mermaid\nflowchart LR\n A-->B\n```\n\n```mermaid\nflowchart LR\n A-->B\n```"
	root := previewRoot(t, mdNodeSource("a", source, "same fence", "label"))
	n := root.Rows[0][0]
	n.Uses = []parser.UseSite{{Definition: "Template", Span: parser.Span{Line: 5}}}
	calls := 0
	p, e := PreparePreviews(root, nil, map[string]MarkdownRenderer{n.Path: {"p", "r", previewRendererFunc(func(string) (Resource, error) {
		calls++
		if calls == 1 {
			return shapeResource(`<circle r="3"/>`), nil
		}
		return Resource{}, fmt.Errorf("second occurrence failed")
	})}}, fullBackend())
	if e != nil {
		t.Fatal(e)
	}
	o, _ := p.Outcome(n.Path)
	if len(o.Diagrams) != 2 || o.Diagrams[0].Resource == nil || o.Diagrams[1].Resource != nil {
		t.Fatal(o)
	}
	n.Uses[0].Definition = "source mutation"
	o.Uses[0].Definition = "accessor mutation"
	again, _ := p.Outcome(n.Path)
	if again.Uses[0].Definition != "Template" {
		t.Fatal("provenance aliased")
	}
	var out strings.Builder
	if e = p.Render(&out, &layout.Box{Path: n.Path, Instance: n, Font: 12, Rect: layout.Rect{W: 1000, H: 200}}); e != nil {
		t.Fatal(e)
	}
	if strings.Count(out.String(), "<image") != 1 || !strings.Contains(out.String(), "second occurrence") {
		t.Fatal("repeated ID merged occurrences", out.String())
	}
}
func TestPreviewFreezeBeforeApplicationCallbacks(t *testing.T) {
	root := previewRoot(t, svgNodeSource("a", "reject")+";"+svgNodeSource("b", "reject"))
	r := shapeResource(`<rect width="3"/>`)
	inputs := bindSVG(root, r)
	backend := fullBackend()
	calls := 0
	backend.SVG = func(Resource) error {
		calls++
		if calls == 1 {
			v := inputs["Main/b"]
			v.Resource.SVG[0] = '!'
			inputs["Main/b"] = PreparedSVG{}
			root.Rows[1][0].Arguments["description"] = parser.Literal{Kind: "string", Value: "new source"}
		}
		return nil
	}
	p, e := PreparePreviews(root, inputs, nil, backend)
	if e != nil {
		t.Fatal(e)
	}
	o, _ := p.Outcome("Main/b")
	if o.Resource.SVG[0] != '<' || o.Description != "Resource description" {
		t.Fatal("callback changed detached inputs", o)
	}
	if e = p.Check(root); e == nil {
		t.Fatal("changed source accepted after callback")
	}
}
func TestPreviewMarkdownBoundsAndTableCopies(t *testing.T) {
	root := previewRoot(t, mdNodeSource("a", "| A | B |\n|---|---|\n| one | two |", "table", "reject"))
	p, e := PreparePreviews(root, nil, nil, fullBackend())
	if e != nil {
		t.Fatal(e)
	}
	o, _ := p.Outcome("Main/a")
	o.Document.Blocks[0].Table[0][0] = "changed"
	again, _ := p.Outcome("Main/a")
	if again.Document.Blocks[0].Table[0][0] != "A" {
		t.Fatal("table accessor alias")
	}
	for _, source := range []string{strings.Repeat("text\n\n", 257), strings.Repeat("```mermaid\nflowchart LR\n A-->B\n```\n\n", 9)} {
		root = previewRoot(t, mdNodeSource("a", source, "bounded", "label"))
		p, e = PreparePreviews(root, nil, nil, fullBackend())
		if e != nil {
			t.Fatal(e)
		}
		o, _ = p.Outcome("Main/a")
		if o.Status != "label" {
			t.Fatal("document bound ignored", o.Status)
		}
	}
}
