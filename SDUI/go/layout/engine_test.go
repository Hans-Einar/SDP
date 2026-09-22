package layout

import (
	"math"
	"os"
	"testing"

	"github.com/Hans-Einar/SDP/SDUI/go/parser"
)

func build(t *testing.T, source string, w, h float64) (*Box, error) {
	t.Helper()
	_, r, e := parser.Compile("sdui 0.2; page=" + source + ";")
	if e != nil {
		t.Fatal(e)
	}
	return (&Engine{}).Layout(r["page"], Size{w, h})
}
func near(t *testing.T, a, b float64) {
	t.Helper()
	if math.Abs(a-b) > .02 {
		t.Fatalf("got %g want %g", a, b)
	}
}
func boxes(b *Box) map[string]*Box {
	m := map[string]*Box{}
	if b != nil {
		b.Walk(func(b *Box) { m[b.Path] = b })
	}
	return m
}
func TestWeightsAndAncestor(t *testing.T) {
	b, e := build(t, `[a=[] {x=1fr,min-x=0.2}, b=[] {x=3fr}] {scale=1,gap=0}`, 1000, 600)
	if e != nil {
		t.Fatal(e)
	}
	m := boxes(b)
	near(t, m["page/a"].Rect.W, 250)
	near(t, m["page/b"].Rect.W, 750)
	b, e = build(t, `[a=[] {scale-x=0.25}, b=[] {x=1fr}] {scale=1,gap=0.01}`, 1000, 600)
	if e != nil {
		t.Fatal(e)
	}
	m = boxes(b)
	near(t, m["page/a"].Rect.W, 250)
	near(t, m["page/b"].Rect.W, 740)
	_, e = build(t, `[[] {scale-x=0.5}, [] {scale-x=0.5}] {scale=1,gap=0.01}`, 1000, 600)
	if e == nil {
		t.Fatal("scale siblings silently shrank")
	}
}
func TestRatioFontAndClip(t *testing.T) {
	b, e := build(t, `["Hello" {font=10}] {16:9,<->}`, 1600, 1000)
	if e != nil {
		t.Fatal(e)
	}
	near(t, b.Rect.H, 900)
	near(t, b.Children[0].Font, 10)
	b, e = build(t, `["Hello" {font=10}] {16:9,<->}`, 800, 1000)
	if e != nil {
		t.Fatal(e)
	}
	near(t, b.Rect.H, 450)
	near(t, b.Children[0].Font, 10)
	_, e = build(t, `[] {16:9,<->}`, 1600, 500)
	if e == nil {
		t.Fatal("ratio silently contained")
	}
	b, e = build(t, `[go=button("Go") {scale-x=2}] {scale=1,overflow-x=clip}`, 100, 80)
	if e != nil {
		t.Fatal(e)
	}
	if b.Hit(150, 10) != nil || b.Hit(10, 10) == nil {
		t.Fatal("clipped hit test")
	}
}
func TestRegionsAndDependencies(t *testing.T) {
	b, e := build(t, `[header="Head", body=[] {scale=1}, footer="Foot"] {scale=1,gap=0}`, 400, 300)
	if e != nil {
		t.Fatal(e)
	}
	m := boxes(b)
	near(t, m["page/body"].Rect.H, 300-2*19.6)
	near(t, m["page/body"].Rect.Y, 19.6)
	_, e = build(t, `[[] {scale-y=0.5}] {<->}`, 400, 300)
	if e == nil {
		t.Fatal("circular content height accepted")
	}
	_, e = build(t, `[button("Go")] {scale=1,overflow-y=scroll}`, 400, 300)
	if e == nil {
		t.Fatal("unsupported scroll accepted")
	}
}
func TestConcept1Geometry(t *testing.T) {
	s, e := os.ReadFile("../../examples/concept1-bucking.sdui")
	if e != nil {
		t.Fatal(e)
	}
	_, roots, e := parser.Compile(string(s))
	if e != nil {
		t.Fatal(e)
	}
	b, e := (&Engine{}).Layout(roots["bucking"], Size{1920, 1200})
	if e != nil {
		t.Fatal(e)
	}
	m := boxes(b)
	near(t, b.Rect.H, 1080)
	near(t, m["bucking/top/length"].Rect.H, m["bucking/top"].Rect.H)
	near(t, m["bucking/middle/selection"].Rect.H, m["bucking/middle"].Rect.H)
	near(t, m["bucking/top/length"].Rect.W, m["bucking/top/diameter"].Rect.W)
	near(t, m["bucking/middle/suggestions"].Rect.W, 2*m["bucking/middle/selection"].Rect.W)
}
func TestSaturatedTracks(t *testing.T) {
	v := distribute([]track{{weight: 1, min: 40, max: 50}, {weight: 3, min: 0, max: 1000}}, 200, 0)
	near(t, v[0], 50)
	near(t, v[1], 150)
	v = distribute([]track{{weight: 1, min: 80, max: 100}, {weight: 1, min: 80, max: 100}}, 100, 0)
	near(t, v[0], 80)
	near(t, v[1], 80)
}
