package layout

import (
	"errors"
	"math"
	"reflect"
	"testing"

	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

func TestFitPreviewGeometry(t *testing.T) {
	for _, c := range []struct {
		name       string
		allocation Rect
		intrinsic  Size
		want       Rect
	}{
		{"landscape", Rect{10, 20, 100, 20}, Size{400, 200}, Rect{40, 20, 40, 20}},
		{"portrait", Rect{-30, -20, 40, 100}, Size{200, 400}, Rect{-30, -10, 40, 80}},
		{"upscale", Rect{10, 20, 100, 100}, Size{4, 2}, Rect{10, 45, 100, 50}},
		{"zero-width", Rect{10, 20, 0, 100}, Size{4, 2}, Rect{}},
		{"zero-height", Rect{10, 20, 100, 0}, Size{4, 2}, Rect{}},
		{"tiny-intrinsic", Rect{0, 0, 100, 100}, Size{math.SmallestNonzeroFloat64, math.SmallestNonzeroFloat64}, Rect{0, 0, 100, 100}},
	} {
		t.Run(c.name, func(t *testing.T) {
			got, err := FitPreview(c.allocation, c.intrinsic)
			if err != nil || got != c.want {
				t.Fatalf("got %+v, %v; want %+v", got, err, c.want)
			}
		})
	}
}

func TestFitPreviewInvalidGeometry(t *testing.T) {
	for _, v := range []float64{-1, math.NaN(), math.Inf(1), math.Inf(-1), 1e7 + 1} {
		for _, c := range []Rect{{0, 0, v, 10}, {0, 0, 10, v}} {
			if _, err := FitPreview(c, Size{10, 10}); err == nil {
				t.Fatalf("accepted allocation %+v", c)
			}
		}
		for _, s := range []Size{{v, 10}, {10, v}} {
			if _, err := FitPreview(Rect{0, 0, 10, 10}, s); err == nil {
				t.Fatalf("accepted intrinsic %+v", s)
			}
		}
	}
	for _, s := range []Size{{0, 10}, {10, 0}, {1e7, math.SmallestNonzeroFloat64}} {
		if _, err := FitPreview(Rect{0, 0, 10, 10}, s); err == nil {
			t.Fatalf("accepted unrepresentable aspect %+v", s)
		}
	}
	for _, v := range []float64{math.NaN(), math.Inf(1), math.Inf(-1)} {
		for _, r := range []Rect{{v, 0, 10, 10}, {0, v, 10, 10}} {
			if _, err := FitPreview(r, Size{10, 10}); err == nil {
				t.Fatal("accepted nonfinite position")
			}
		}
	}
	if _, err := FitPreview(Rect{}, Size{}); err == nil {
		t.Fatal("empty allocation bypassed intrinsic validation")
	}
}

const previewCall = `svg(art.Chart.@resource,description="Chart",fallback="label")`

func previewRoot(t *testing.T, body string) *parser.Instance {
	t.Helper()
	_, roots, err := parser.Compile(`sdui 0.3; ref: art "not-opened-by-layout"; page=` + body + `;`)
	if err != nil {
		t.Fatal(err)
	}
	return roots["page"]
}

type previewMetrics struct {
	TextMetrics
	natural, minimum Size
	err              error
	calls            map[string]int
	width, font      float64
	gutters          bool
}

func (m *previewMetrics) MeasurePreview(n *parser.Instance, font, width float64) (Size, Size, error) {
	if m.calls == nil {
		m.calls = map[string]int{}
	}
	m.calls[n.Path]++
	m.width, m.font = width, font
	return m.natural, m.minimum, m.err
}
func (*previewMetrics) MeasureTabs(*parser.Instance, float64, Size) (TabsMetrics, error) {
	return TabsMetrics{Header: Size{80, 24}}, nil
}
func (*previewMetrics) MeasureSplit(*parser.Instance, float64) (float64, error) { return 10, nil }
func (m *previewMetrics) MeasureViewport(n *parser.Instance, _ float64) (ViewportInsets, error) {
	v := ViewportInsets{}
	if m.gutters && scrolls(n, "x") {
		v.Bottom = 10
	}
	if m.gutters && scrolls(n, "y") {
		v.Right = 10
	}
	return v, nil
}

func TestPreviewNaturalAndSourceAllocation(t *testing.T) {
	m := &previewMetrics{natural: Size{400, 200}}
	e := &Engine{Measure: m}
	for _, c := range []struct {
		style string
		want  Rect
	}{
		{"", Rect{0, 0, 400, 200}},
		{"{scale-x=0.25,scale-y=0.1}", Rect{0, 0, 100, 20}},
		{"{x=fill,scale-y=0.1,max-x=0.25}", Rect{0, 0, 100, 20}},
		{"{min-x=0.5,min-y=0.5}", Rect{0, 0, 400, 200}},
	} {
		root := previewRoot(t, `[p=`+previewCall+c.style+`] {scale=1,gap=0,font=22}`)
		g, err := e.Layout(root, Size{400, 200})
		if err != nil {
			t.Fatal(err)
		}
		if got := boxes(g)["page/p"].Rect; got != c.want {
			t.Fatalf("%s: %+v", c.style, got)
		}
		if m.font != 22 || m.width != c.want.W {
			t.Fatal("wrong final inherited metrics", m.font, m.width)
		}
	}
	// Source minima remain authoritative even when a scalable image has none.
	m.natural = Size{1, 1}
	root := previewRoot(t, `[p=`+previewCall+` {min-x=0.5,min-y=0.5}] {scale=1,gap=0}`)
	g, err := e.Layout(root, Size{400, 200})
	if err != nil || boxes(g)["page/p"].Rect != (Rect{0, 0, 200, 100}) {
		t.Fatal(g, err)
	}
}

func TestPreviewRequiredTextMinimumAndErrors(t *testing.T) {
	m := &previewMetrics{natural: Size{100, 40}, minimum: Size{30, 40}}
	e := &Engine{Measure: m}
	for _, style := range []string{`{scale-x=0.1}`, `{scale-y=0.1}`, `{max-x=0.1}`, `{max-y=0.1}`} {
		root := previewRoot(t, `[p=`+previewCall+style+`] {scale=1,gap=0}`)
		if _, err := e.Layout(root, Size{200, 100}); err == nil {
			t.Fatal("lost required text minimum", style)
		}
	}
	root := previewRoot(t, `[p=`+previewCall+` {scale=1}] {scale=1,gap=0}`)
	if _, err := (&Engine{}).Layout(root, Size{200, 100}); err == nil {
		t.Fatal("placeholder licensed explicit preview")
	}
	for _, v := range []float64{-1, math.NaN(), math.Inf(1), 1e7 + 1} {
		for _, minimum := range []bool{false, true} {
			m.natural, m.minimum = Size{100, 40}, Size{30, 40}
			if minimum {
				m.minimum.H = v
			} else {
				m.natural.W = v
			}
			if _, err := e.Layout(root, Size{200, 100}); err == nil {
				t.Fatal("invalid metrics accepted", v, minimum)
			}
		}
	}
	m.err = errors.New("missing frozen path")
	if _, err := e.Layout(root, Size{200, 100}); !errors.Is(err, m.err) {
		t.Fatal("provider error lost", err)
	}
}

func TestPreviewSplitMinimumDoesNotUseIntrinsicImageSize(t *testing.T) {
	root := previewRoot(t, `[s=split(axis="horizontal",proportion=0.1)[a=`+previewCall+` {min-x=0.2};b=`+previewCall+`] {scale=1}] {scale=1,gap=0}`)
	s, err := runtime.New("preview-split", root)
	if err != nil {
		t.Fatal(err)
	}
	m := &previewMetrics{natural: Size{4000, 2000}, minimum: Size{30, 20}}
	g, err := (&Engine{Measure: m}).LayoutSnapshot(s.Snapshot(), Size{210, 100})
	if err != nil {
		t.Fatal(err)
	}
	sp := g.Splits["page/s"]
	if sp.First != (Rect{0, 0, 40, 100}) || sp.Geometry != (runtime.SplitGeometry{Lower: .2, Upper: .85, Effective: .2}) {
		t.Fatalf("intrinsic/source minima confused %+v", sp)
	}
	if g.PresentationState().Splits["page/s"] != sp.Geometry {
		t.Fatal("minimum omitted from pure gate projection")
	}
}

func TestPreviewLegacyGeometryUnchanged(t *testing.T) {
	for _, profile := range []string{"0.2", "0.3"} {
		_, roots, err := parser.Compile(`sdui ` + profile + `; ref: art "legacy"; page=[p=svg(art.Chart.@resource,label="Caption");"Prose"] {scale=1,gap=0};`)
		if err != nil {
			t.Fatal(err)
		}
		m := &previewMetrics{err: errors.New("legacy must not call adjunct")}
		before, err := (&Engine{}).Layout(roots["page"], Size{400, 200})
		if err != nil {
			t.Fatal(err)
		}
		after, err := (&Engine{Measure: m}).Layout(roots["page"], Size{400, 200})
		if err != nil || !reflect.DeepEqual(before, after) || len(m.calls) != 0 {
			t.Fatal("legacy geometry changed", profile, err)
		}
	}
}
