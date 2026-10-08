package markdown

import (
	"fmt"
	"strings"
	"testing"
)

func paddedResource(size int, tag string) Resource {
	r := shapeResource("<!--" + tag + "-->")
	pad := size - len(r.SVG)
	if pad < 0 {
		panic("bad fixture size")
	}
	r.SVG = []byte(strings.Replace(string(r.SVG), "-->", strings.Repeat("x", pad)+"-->", 1))
	return r
}
func TestPreviewAdmissionBudgetBeforeBackend(t *testing.T) {
	bodies := []string{}
	for i := 0; i < 9; i++ {
		bodies = append(bodies, svgNodeSource(fmt.Sprintf("p%d", i), "label")+` {visible=false}`)
	}
	root := previewRoot(t, strings.Join(bodies, ";"))
	resources := bindSVG(root, paddedResource(resourceLimit, "shared"))
	// More than 32 MiB across paths, but exactly one distinct full digest.
	p, e := PreparePreviews(root, resources, nil, PreviewBackend{})
	if e != nil {
		t.Fatal("duplicate bytes charged repeatedly", e)
	}
	o, _ := p.Outcome("Main/p8")
	if o.Status != "label" || o.Resource != nil || o.SHA256 == "" {
		t.Fatal(o)
	}
	for i := 0; i < 8; i++ {
		path := fmt.Sprintf("Main/p%d", i)
		r := resources[path]
		r.Resource = paddedResource(resourceLimit, fmt.Sprintf("distinct%d", i))
		r.SHA256 = digest(r.Resource.SVG)
		resources[path] = r
	}
	// Eight exact 4 MiB resources plus a duplicate: exact 32 MiB admission passes,
	// despite hidden declarations and every backend outcome being a label.
	resources["Main/p8"] = resources["Main/p0"]
	if _, e = PreparePreviews(root, resources, nil, PreviewBackend{}); e != nil {
		t.Fatal("exact 32 MiB rejected", e)
	}
	r := resources["Main/p8"]
	r.Resource = shapeResource(`<rect/>`)
	r.SHA256 = digest(r.Resource.SVG)
	resources["Main/p8"] = r
	if _, e = PreparePreviews(root, resources, nil, PreviewBackend{}); e == nil || !strings.Contains(e.Error(), "preview-budget") {
		t.Fatal("label refunded admission tally", e)
	}
	// Invalid/malformed and per-resource oversized inputs never enter the aggregate.
	for _, bad := range []Resource{shapeResource(`<script/>`), paddedResource(resourceLimit+1, "oversize")} {
		r.Resource = bad
		r.SHA256 = digest(bad.SVG)
		resources["Main/p8"] = r
		if _, e = PreparePreviews(root, resources, nil, PreviewBackend{}); e != nil {
			t.Fatal("invalid resource charged aggregate", e)
		}
	}
}
func TestPreviewCustomRendererValidationAndLimits(t *testing.T) {
	root := previewRoot(t, mdNodeSource("a", diagramMD, "A", "label"))
	backend := fullBackend()
	checks := 0
	backend.Mermaid = func(Resource) error { checks++; return nil }
	for _, r := range []Resource{shapeResource(`<script/>`), {SVG: []byte(`<svg viewBox="0 0 40 20"/>`), Width: 0, Height: 20}, paddedResource(resourceLimit+1, "too-big")} {
		p, e := PreparePreviews(root, nil, map[string]MarkdownRenderer{"Main/a": {"p", "r", previewRendererFunc(func(string) (Resource, error) { return r, nil })}}, backend)
		if e != nil {
			t.Fatal(e)
		}
		o, _ := p.Outcome("Main/a")
		if o.Status != "partial" || o.Diagrams[0].Resource != nil || checks != 0 {
			t.Fatal("unvalidated custom output reached backend", o, checks)
		}
	}
	// Count across supplied and renderer resources with full-digest deduplication.
	root = previewRoot(t, svgNodeSource("svg", "reject")+";"+mdNodeSource("md", diagramMD, "D", "reject"))
	shared := shapeResource(`<rect/>`)
	p, e := PreparePreviews(root, bindSVG(root, shared), map[string]MarkdownRenderer{"Main/md": {"p", "r", previewRendererFunc(func(string) (Resource, error) { return shared, nil })}}, fullBackend())
	if e != nil {
		t.Fatal(e)
	}
	a, _ := p.Outcome("Main/svg")
	b, _ := p.Outcome("Main/md")
	if a.SHA256 != b.Diagrams[0].SHA256 {
		t.Fatal("cross-kind digest identity lost")
	}
}

func TestPreviewUnavailableWholeMarkdownStillPreparesRegisteredOutput(t *testing.T) {
	root := previewRoot(t, mdNodeSource("a", diagramMD, "A", "label"))
	calls := 0
	r := shapeResource(`<rect/>`)
	p, e := PreparePreviews(root, nil, map[string]MarkdownRenderer{"Main/a": {"p", "r", previewRendererFunc(func(string) (Resource, error) { calls++; return r, nil })}}, PreviewBackend{})
	if e != nil {
		t.Fatal(e)
	}
	o, _ := p.Outcome("Main/a")
	if calls != 1 || o.Status != "label" || o.Document != nil || o.Diagrams[0].Resource != nil || o.Diagrams[0].SHA256 != digest(r.SVG) {
		t.Fatal("unavailable backend skipped preparation", calls, o)
	}
}
