package fynehost

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"image"
	"strings"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/markdown"
)

// Rasterize through the same public canvas.Image decoder as mounted previews,
// not the layout inspector or an independently configured SVG renderer.
func previewRaster(t *testing.T, resource fyne.Resource, width float32) image.Image {
	t.Helper()
	im := canvas.NewImageFromResource(resource)
	im.Resize(fyne.NewSize(width, width/2))
	im.Refresh()
	if im.Image == nil {
		t.Fatal("native SVG decoder produced no pixels")
	}
	return im.Image
}

func TestPreparedPreviewOriginResourceLifetime(t *testing.T) {
	h := documentHost(t)
	source := []byte(`<?xml version="1.0"?><svg xmlns="http://www.w3.org/2000/svg" width="400" height="200" viewBox="40 20 400 200" fill="#238248"><!-- retained child bytes --><title>Chart &amp; details</title><rect x="40" y="20" width="400" height="200"/></svg>`)
	r := previewRequest(t, previewSource, 1)
	bindPreview(&r, "page/picture", source)
	b, err := h.Prepare(r)
	if err != nil {
		t.Fatal(err)
	}
	resource := b.previewResources["page/picture"]
	if !bytes.Contains(resource.Content(), []byte(`<!-- retained child bytes --><title>Chart &amp; details</title>`)) {
		t.Fatal("child content rewritten")
	}
	outcome, _ := b.previews.Outcome("page/picture")
	if !bytes.Equal(outcome.Resource.SVG, source) {
		t.Fatal("original prepared provider resource changed")
	}
	if err := h.Commit(b); err != nil {
		t.Fatal(err)
	}
	c := b.Controls()["page/picture"].(*previewControl)
	if err := h.Resize(layout.Size{W: 600, H: 450}); err != nil {
		t.Fatal(err)
	}
	if c.image.Resource != resource || b.previewResources["page/picture"] != resource {
		t.Fatal("resize recreated frozen derivative")
	}
	copy := resource.Content()
	copy[0] = '!'
	if resource.Content()[0] != '<' {
		t.Fatal("derived native bytes mutable")
	}
	b.Close()
	if c.image.Resource != nil || b.previewResources != nil {
		t.Fatal("close retained native derivative")
	}

	zero := []byte(`<?xml version="1.0"?><!-- prolog --><svg viewBox="0 0 400 200"><rect width="400" height="200"/></svg>`)
	data, err := nativePreviewSVG(markdown.Resource{SVG: zero, Width: 400, Height: 200})
	if err != nil || !bytes.Equal(data, zero) {
		t.Fatalf("zero-origin source bytes changed: %v", err)
	}
}

func TestPreparedPreviewOriginDerivativeFailureBeforeOutcome(t *testing.T) {
	h := documentHost(t)
	// Source-only coordinates are finite; combining the required origin shift
	// with the child position overflows. Never freeze a rendered outcome.
	source := []byte(`<svg viewBox="-1e308 0 400 200"><rect x="1e308" width="1" height="1"/></svg>`)
	if _, err := markdown.ValidateSVGResource(markdown.Resource{SVG: source, Width: 400, Height: 200}); err != nil {
		t.Fatalf("repro source invalid: %v", err)
	}
	for _, policy := range []string{"label", "reject"} {
		r := previewRequest(t, strings.Replace(previewSource, `fallback="reject"`, `fallback="`+policy+`"`, 1), 1)
		bindPreview(&r, "page/picture", source)
		b, err := h.Prepare(r)
		if policy == "reject" {
			if err == nil {
				b.Close()
				t.Fatal("invalid derivative admitted")
			}
			continue
		}
		if err != nil {
			t.Fatal(err)
		}
		if b.previewOutcomes["page/picture"].Status != "label" || b.previewResources["page/picture"] != nil {
			t.Fatal("invalid derivative froze rendered resource")
		}
		for _, cap := range b.capabilities() {
			if cap.ID == "svg-resource" {
				t.Fatal("invalid derivative advertised native SVG capability")
			}
		}
		b.Close()
	}
}

func TestPreparedPreviewOriginDerivativeByteBound(t *testing.T) {
	h := documentHost(t)
	const prefix = `<svg viewBox="40 20 400 200"><!--`
	const suffix = `--><rect x="40" y="20" width="400" height="200"/></svg>`
	source := []byte(prefix + strings.Repeat("a", (4<<20)-len(prefix)-len(suffix)) + suffix)
	if _, err := markdown.ValidateSVGResource(markdown.Resource{SVG: source, Width: 400, Height: 200}); err != nil {
		t.Fatal("original resource budget changed", err)
	}
	for _, policy := range []string{"label", "reject"} {
		r := previewRequest(t, strings.Replace(previewSource, `fallback="reject"`, `fallback="`+policy+`"`, 1), 1)
		bindPreview(&r, "page/picture", source)
		b, err := h.Prepare(r)
		if policy == "reject" {
			if err == nil {
				b.Close()
				t.Fatal("over-bound derivative admitted")
			}
			continue
		}
		if err != nil {
			t.Fatal(err)
		}
		o := b.previewOutcomes["page/picture"]
		if o.Status != "label" || !strings.Contains(o.Diagnostic, "4 MiB") || b.previewResources["page/picture"] != nil {
			t.Fatalf("untruthful backend byte outcome: %+v", o)
		}
		b.Close()
	}
}

func TestPreparedPreviewOriginEmptyDocument(t *testing.T) {
	_ = documentHost(t)
	for _, source := range []string{
		`<svg viewBox="0 0 400 200"/>`,
		`<svg viewBox="40 20 400 200"/>`,
		`<svg viewBox="-40 -20 400 200"></svg>`,
		`<svg viewBox="0 0 400 200" transform="scale(.5)"/>`,
		`<svg viewBox="40 20 400 200" transform="scale(.5)"/>`,
	} {
		r := markdown.Resource{SVG: []byte(source), Width: 400, Height: 200}
		if _, err := markdown.ValidateSVGResource(r); err != nil {
			t.Fatal("invalid empty repro", err)
		}
		data, err := nativePreviewSVG(r)
		if err != nil {
			t.Fatalf("valid empty SVG rejected %s: %v", source, err)
		}
		if ink := previewInk(previewRaster(t, newPreviewResource(data), 200)); !ink.Empty() {
			t.Fatal("empty SVG gained pixels", ink)
		}
		if source == `<svg viewBox="0 0 400 200"/>` && string(data) != source {
			t.Fatal("ordinary empty bytes changed")
		}
	}
}

func TestPreparedPreviewOriginNativePixels(t *testing.T) {
	h := documentHost(t)
	cases := []struct{ name, box, attrs, child, expected string }{
		{"negative", "-40 -20 400 200", "", `<rect x="-40" y="-20" width="400" height="200" fill="#238248"/>`, `<rect width="400" height="200" fill="#238248"/>`},
		{"positive", "40 20 400 200", "", `<rect x="40" y="20" width="400" height="200" fill="#238248"/>`, `<rect width="400" height="200" fill="#238248"/>`},
		{"mixed", "-40 20 400 200", "", `<rect x="-40" y="20" width="400" height="200" fill="#238248"/>`, `<rect width="400" height="200" fill="#238248"/>`},
		{"root-transform", "40 20 400 200", `transform="scale(2 2)" fill="#238248" opacity="0.5"`, `<rect x="50" y="25" width="20" height="10"/>`, `<rect x="20" y="10" width="40" height="20" fill="#238248" opacity="0.5"/>`},
		{"root-single-scale", "40 20 400 200", `transform="scale(.5)"`, `<rect x="120" y="60" width="160" height="80" fill="#238248"/>`, `<rect x="40" y="20" width="80" height="40" fill="#238248"/>`},
		{"zero-root-single-scale", "0 0 400 200", `transform="scale(2)" fill="#238248"`, `<rect x="10" y="5" width="20" height="10"/>`, `<rect x="20" y="10" width="40" height="20" fill="#238248"/>`},
		{"zero-child-single-scale", "0 0 400 200", "", `<g transform="translate(10 5) scale(2e0)"><rect transform="scale(.5)" x="20" y="10" width="40" height="20" fill="#238248"/></g>`, `<rect x="30" y="15" width="40" height="20" fill="#238248"/>`},
		{"zero-single-scale-whitespace", "0 0 400 200", `transform="scale ( +2.0e0 )" fill="#238248"`, `<rect x="10" y="5" width="20" height="10"/>`, `<rect x="20" y="10" width="40" height="20" fill="#238248"/>`},
		{"negative-root-child-transform", "-40 -20 400 200", `transform="translate(10 5) scale(2)" fill="#238248"`, `<g transform="translate(5 3)"><rect x="-30" y="-15" width="20" height="10"/></g>`, `<rect x="40" y="21" width="40" height="20" fill="#238248"/>`},
		{"zero-ordinary", "0 0 400 200", "", `<circle cx="100" cy="70" r="30" fill="#238248"/>`, `<circle cx="100" cy="70" r="30" fill="#238248"/>`},
	}
	for i, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			source := []byte(fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="%s" %s>%s</svg>`, tc.box, tc.attrs, tc.child))
			r := previewRequest(t, previewSource, uint64(i+1))
			bindPreview(&r, "page/picture", source)
			b, err := h.Prepare(r)
			if err != nil {
				t.Fatal(err)
			}
			defer b.Close()
			resource := b.previewResources["page/picture"]
			if b.previewOutcomes["page/picture"].SHA256 != fmt.Sprintf("%x", sha256.Sum256(source)) {
				t.Fatal("source digest replaced by derivative identity")
			}
			if resource.Name() != fmt.Sprintf("preview-%x.svg", sha256.Sum256(resource.Content())) {
				t.Fatal("native cache name does not identify actual bytes")
			}
			if tc.name == "zero-ordinary" && !bytes.Equal(source, resource.Content()) {
				t.Fatal("zero origin bytes changed")
			}
			oracle := newPreviewResource([]byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 400 200">` + tc.expected + `</svg>`))
			for _, width := range []float32{200, 400, 754} {
				got, want := previewRaster(t, resource, width), previewRaster(t, oracle, width)
				if got.Bounds() != want.Bounds() {
					t.Fatal("raster size mismatch")
				}
				mismatch := 0
				for y := got.Bounds().Min.Y; y < got.Bounds().Max.Y; y++ {
					for x := got.Bounds().Min.X; x < got.Bounds().Max.X; x++ {
						gr, gg, gb, ga := got.At(x, y).RGBA()
						wr, wg, wb, wa := want.At(x, y).RGBA()
						if gr != wr || gg != wg || gb != wb || ga != wa {
							mismatch++
						}
					}
				}
				if mismatch != 0 {
					t.Errorf("width %g: %d native pixels differ from zero-origin geometry; got %v want %v; native %s", width, mismatch, previewInk(got), previewInk(want), resource.Content())
				}
			}
		})
	}
}

func previewInk(im image.Image) image.Rectangle {
	r := image.Rectangle{}
	for y := im.Bounds().Min.Y; y < im.Bounds().Max.Y; y++ {
		for x := im.Bounds().Min.X; x < im.Bounds().Max.X; x++ {
			_, _, _, a := im.At(x, y).RGBA()
			if a != 0 {
				r = r.Union(image.Rect(x, y, x+1, y+1))
			}
		}
	}
	return r
}
