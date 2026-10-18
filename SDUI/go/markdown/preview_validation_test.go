package markdown

import (
	"math"
	"strings"
	"testing"
)

func shapeResource(body string) Resource {
	return Resource{SVG: []byte(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 40 20">` + body + `</svg>`), Width: 40, Height: 20}
}
func TestPreviewSVGClosedGrammar(t *testing.T) {
	valid := []string{
		`<g transform="translate(1,2) scale(2) rotate(20 3 4)"><rect x="-1" y=".5" width="3" height="4" rx="1" fill="#aBc" opacity=".5"/></g>`,
		`<path d="M1 2 3 4 L5 6 h2 v-1 C1 2 3 4 5 6 s1 2 3 4 Q1 2 3 4 t2 3 A4 5 30 0 1 6 7 z" stroke="red" fill="none"/>`,
		`<circle cx="1" cy="2" r="3"/><ellipse rx="2" ry="1"/><line x1="-2" y1="0" x2="3" y2="4"/><polyline points="0,0 2,2"/><polygon points="0 0 1 1 2 0"/>`,
		`<title>A &amp; B</title><desc>description</desc><g transform="matrix(1 0 0 1 0 0) skewX(10) skewY(5)"><rect width="4" height="3" fill="rgb(1,2,3)"/></g>`,
	}
	for _, body := range valid {
		r := shapeResource(body)
		got, e := ValidateSVGResource(r)
		if e != nil {
			t.Errorf("valid %s: %v", body, e)
			continue
		}
		r.SVG[0] = '!'
		if got.SVG[0] != '<' {
			t.Fatal("retained caller bytes")
		}
	}
	invalid := []string{
		`<script/>`, `<image href="data:image/svg+xml;base64,abc"/>`, `<foreignObject/>`, `<animate/>`, `<style>path{fill:red}</style>`, `<text>unsupported</text>`, `<use href="#local"/>`, `<linearGradient/>`,
		`<rect style="fill:red"/>`, `<rect onload="x"/>`, `<rect class="x"/>`, `<rect fill="url(#x)"/>`, `<rect fill="url(https://x)"/>`, `<rect fill="currentColor"/>`, `<rect width="NaN"/>`, `<rect height="1e999"/>`, `<rect width="-1"/>`, `<rect opacity="2"/>`, `<rect width="2px"/>`, `<rect foo="0"/>`, `<rect x="1" x="2"/>`, `<x:rect xmlns:x="urn:other"/>`, `<rect xmlns="urn:other"/>`,
		`<path d="L0 0"/>`, `<path d="M0"/>`, `<path d="M0 0 C1 2"/>`, `<path d="M0 0 A1 1 0 2 0 1 1"/>`, `<path d="M0 0 A1 1 0 1.0 0 1 1"/>`, `<path d="M0 0 A1 1 0 1e0 0 1 1"/>`, `<path d="M0 0 A1 1 0 +1 0 1 1"/>`, `<path d="M0 0 A1 1 0 01 0 1 1"/>`, `<path d="M0 0 A-1 1 0 0 0 1 1"/>`, `<path d="M0 0 X1 2"/>`, `<path d="M0 0,"/>`, `<polygon points="1 2 3"/>`,
		`<g transform="translate(1 2 3)"/>`, `<g transform="matrix(1 2)"/>`, `<g transform="scale(1e308) scale(1e308)"/>`, `<g transform="skewX(90)"/>`, `<g transform="rotate(Infinity)"/>`, `<g transform="translate(1),"/>`, `<g transform="translate(1e308)"><g transform="translate(1e308)"/></g>`,
		`<path d="M1e308 1 l1e308 1"/>`, `<title><g/></title>`, `<svg viewBox="0 0 1 1"/>`, `unstructured text`,
	}
	for _, body := range invalid {
		if _, e := ValidateSVGResource(shapeResource(body)); e == nil {
			t.Errorf("accepted %s", body)
		}
	}
}
func TestPreviewCompleteDocumentAndDimensions(t *testing.T) {
	base := shapeResource(`<rect width="40" height="20"/>`)
	cases := []Resource{
		{SVG: append(append([]byte{}, base.SVG...), base.SVG...), Width: 40, Height: 20},
		{SVG: append(append([]byte{}, base.SVG...), []byte("junk")...), Width: 40, Height: 20},
		{SVG: []byte(`<!DOCTYPE svg [<!ENTITY x "boom">]>` + string(base.SVG)), Width: 40, Height: 20},
		{SVG: []byte(`<?probe x?>` + string(base.SVG)), Width: 40, Height: 20},
		{SVG: []byte(strings.Replace(string(base.SVG), "</svg>", "<?x y?></svg>", 1)), Width: 40, Height: 20},
		{SVG: base.SVG, Width: math.NaN(), Height: 20}, {SVG: base.SVG, Width: math.Inf(1), Height: 20}, {SVG: base.SVG, Width: 41, Height: 20},
		{SVG: []byte(`<svg viewBox="NaN 0 40 20"/>`), Width: 40, Height: 20},
		{SVG: []byte(`<svg viewBox="0 0 40 20" width="41"/>`), Width: 40, Height: 20},
		{SVG: []byte(`<svg viewBox="0 0 40 20" width="40px"/>`), Width: 40, Height: 20},
		{SVG: []byte(`<svg viewBox="0 0 40 20">`), Width: 40, Height: 20},
		{SVG: []byte(strings.Repeat(" ", resourceLimit+1)), Width: 40, Height: 20},
	}
	for i, r := range cases {
		for _, check := range []func(Resource) (Resource, error){ValidateSVGResource, ValidateMermaidResource} {
			if _, e := check(r); e == nil {
				t.Errorf("case %d accepted", i)
			}
		}
	}
	r := base
	r.SVG = []byte(`<?xml version="1.0" encoding="UTF-8"?>` + string(base.SVG) + "\n<!--end-->")
	if _, e := ValidateSVGResource(r); e != nil {
		t.Fatal(e)
	}
}
func TestPreviewMermaidProfileSeparate(t *testing.T) {
	r := shapeResource(`<defs><marker id="arrow" viewBox="0 0 10 10" refX="5" refY="5" markerWidth="8" markerHeight="8" orient="auto"><path d="M0 0 L10 5 L0 10z"/></marker></defs><path d="M0 0L20 10" marker-end="url(#arrow)" stroke-dasharray="1,0"/><text x="10" y="10" font-size="12" font-family="sans-serif"><tspan dy="0">Hello</tspan></text>`)
	if _, e := ValidateSVGResource(r); e == nil {
		t.Fatal("shape route admitted text/marker")
	}
	if _, e := ValidateMermaidResource(r); e != nil {
		t.Fatal(e)
	}
	for _, body := range []string{`<text style="font:url(https://x)">x</text>`, `<path marker-end="url(https://x)"/>`, `<text font-family="url(x)">x</text>`, `<a href="https://x"/>`, `<set/>`} {
		if _, e := ValidateMermaidResource(shapeResource(body)); e == nil {
			t.Errorf("accepted %s", body)
		}
	}
}

func FuzzPreviewSVGAdmission(f *testing.F) {
	for _, r := range []Resource{shapeResource(`<rect width="10" height="5"/>`), shapeResource(`<path d="M0 0 A1 2 0 0 1 3 4z"/>`), shapeResource(`<script/>`)} {
		f.Add(r.SVG)
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 16384 {
			return
		}
		for _, validate := range []func(Resource) (Resource, error){ValidateSVGResource, ValidateMermaidResource} {
			r, e := validate(Resource{SVG: data, Width: 40, Height: 20})
			if e != nil {
				continue
			}
			if string(r.SVG) != string(data) || r.Width != 40 || r.Height != 20 {
				t.Fatal("admission changed bytes/dimensions")
			}
			if len(data) > 0 {
				original := r.SVG[0]
				data[0] ^= 1
				if r.SVG[0] != original {
					t.Fatal("admission alias")
				}
				data[0] ^= 1
			}
		}
	})
}
