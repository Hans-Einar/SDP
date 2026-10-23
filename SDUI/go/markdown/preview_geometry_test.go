package markdown

import (
	"math"
	"strings"
	"sync"
	"testing"

	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
)

func TestPreviewSharedGeometryCaptionAndStatus(t *testing.T) {
	root := previewRoot(t, svgNodeSource("a", "label"))
	n := root.Rows[0][0]
	r := shapeResource(`<rect width="40" height="20"/>`)
	p, e := PreparePreviews(root, bindSVG(root, r), nil, fullBackend())
	if e != nil {
		t.Fatal(e)
	}
	b := &layout.Box{Instance: n, Path: n.Path, Font: 10, Rect: layout.Rect{X: 10, Y: 20, W: 100, H: 20}, Clip: layout.Rect{X: 45, Y: 25, W: 3, H: 5}}
	image, caption, status, e := p.SVGRects(b)
	if e != nil || image != (layout.Rect{X: 40, Y: 20, W: 40, H: 20}) || caption != (layout.Rect{}) || status != (layout.Rect{}) {
		t.Fatal(image, caption, status, e)
	}
	old := image
	b.Clip = layout.Rect{}
	image, _, _, e = p.SVGRects(b)
	if e != nil || image != old {
		t.Fatal("clip changed fit", e, image)
	}
	natural, minimum, e := p.MeasurePreview(n, 10, 100)
	if e != nil || natural != (layout.Size{W: 40, H: 20}) || minimum != (layout.Size{}) {
		t.Fatal(natural, minimum, e)
	}
	n.Arguments["label"] = parser.Literal{Kind: "string", Value: "A long caption deliberately requiring ellipsis"}
	b.Rect.W = 30
	b.Rect.H = 40
	image, caption, _, e = p.SVGRects(b)
	if e != nil || caption.Y != 46 || caption.H != 14 || image.Y < 20 || image.Y+image.H > 46+1e-7 {
		t.Fatal(image, caption, e)
	}
	var out strings.Builder
	if e = p.RenderSVGText(&out, b); e != nil {
		t.Fatal(e)
	}
	if strings.Contains(out.String(), "<image") || !strings.Contains(out.String(), "…") || !strings.Contains(out.String(), "<path") {
		t.Fatal(out.String())
	}
	// Missing resource label must retain full status independently of caption.
	p, e = PreparePreviews(root, nil, nil, fullBackend())
	if e != nil {
		t.Fatal(e)
	}
	natural, minimum, e = p.MeasurePreview(n, 10, 90)
	if e != nil || natural.H != minimum.H || minimum.H <= 14 {
		t.Fatal(natural, minimum, e)
	}
	b.Rect.W = 90
	b.Rect.H = minimum.H
	image, caption, status, e = p.SVGRects(b)
	if e != nil || image != (layout.Rect{}) || status.H+caption.H != minimum.H {
		t.Fatal(image, caption, status, e)
	}
	out.Reset()
	if e = p.RenderSVGText(&out, b); e != nil {
		t.Fatal(e)
	}
	if !strings.Contains(out.String(), "unavailable") || !strings.Contains(out.String(), "<title>") {
		t.Fatal(out.String())
	}
	b.Rect.H -= 1
	if _, _, _, e = p.SVGRects(b); e == nil {
		t.Fatal("truncated status admitted")
	}
	out.Reset()
	out.WriteString("keep")
	if e = p.RenderSVGText(&out, b); e == nil || out.String() != "keep" {
		t.Fatal("failed render mutated output")
	}
	b.Rect.H = minimum.H
	b.Rect.W = minimum.W / 2
	if _, _, _, e = p.SVGRects(b); e == nil {
		t.Fatal("too-narrow status admitted")
	}
	for _, v := range []float64{math.NaN(), math.Inf(1), -1, 1e8} {
		if _, _, e = p.MeasurePreview(n, 10, v); e == nil {
			t.Fatal("invalid width", v)
		}
	}
}
func TestPreviewReadOnlyConcurrentProjections(t *testing.T) {
	root := previewRoot(t, mdNodeSource("a", diagramMD, "diagram", "label"))
	p, e := PreparePreviews(root, nil, nil, fullBackend())
	if e != nil {
		t.Fatal(e)
	}
	n := root.Rows[0][0]
	var wg sync.WaitGroup
	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 10 {
				if e := p.Check(root); e != nil {
					t.Error(e)
				}
				o, _ := p.Outcome(n.Path)
				o.Document.Blocks[0].Text = "caller edit"
				size, e := p.Measure(n, 12, 240)
				if e != nil {
					t.Error(e)
					return
				}
				var out strings.Builder
				if e = p.Render(&out, &layout.Box{Instance: n, Path: n.Path, Font: 12, Rect: layout.Rect{W: 240, H: size.H}}); e != nil {
					t.Error(e)
				}
			}
		}()
	}
	wg.Wait()
}
func TestPreviewLegacyProjectionPreserved(t *testing.T) {
	for _, profile := range []string{"0.2", "0.3"} {
		d, e := parser.Parse(`sdui ` + profile + `; Main=["Before\n\n` + "```mermaid\\nflowchart LR\\n A-->B\\n```" + `"];`)
		if e != nil {
			t.Fatal(e)
		}
		roots, e := parser.Normalize(d)
		if e != nil {
			t.Fatal(e)
		}
		root := roots["Main"]
		old, e := Prepare(root, nil)
		if e != nil {
			t.Fatal(e)
		}
		newP, e := PreparePreviews(root, nil, nil, PreviewBackend{})
		if e != nil {
			t.Fatal(e)
		}
		n := root.Rows[0][0]
		a, e := old.Measure(n, 12, 150)
		if e != nil {
			t.Fatal(e)
		}
		b, e := newP.Measure(n, 12, 150)
		if e != nil || a != b {
			t.Fatal(a, b, e)
		}
		box := &layout.Box{Instance: n, Path: n.Path, Font: 12, Rect: layout.Rect{W: 150, H: a.H}}
		var x, y strings.Builder
		if e = old.Render(&x, box); e != nil {
			t.Fatal(e)
		}
		if e = newP.Render(&y, box); e != nil || x.String() != y.String() {
			t.Fatal("legacy bytes changed", e)
		}
	}
}

func TestPreviewMarkdownFallbackRejectsTruncatedAllocation(t *testing.T) {
	for _, source := range []string{"<b>unsupported</b>", diagramMD} {
		root := previewRoot(t, mdNodeSource("a", source, "Full unavailable explanation", "label"))
		n := root.Rows[0][0]
		p, e := PreparePreviews(root, nil, nil, fullBackend())
		if e != nil {
			t.Fatal(e)
		}
		size, e := p.Measure(n, 12, 180)
		if e != nil {
			t.Fatal(e)
		}
		var out strings.Builder
		out.WriteString("old")
		b := &layout.Box{Path: n.Path, Instance: n, Font: 12, Rect: layout.Rect{W: 180, H: size.H - 1}}
		if e = p.Render(&out, b); e == nil || out.String() != "old" {
			t.Fatal("truncated fallback published or output changed", e)
		}
		b.Rect.H = size.H
		out.Reset()
		if e = p.Render(&out, b); e != nil {
			t.Fatal("valid allocation rejected", e)
		}
	}
}
