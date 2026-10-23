package layout

import (
	"math"
	"reflect"
	"testing"

	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

type gutterMetrics struct {
	fixtureMetrics
	calls   map[string]int
	invalid *ViewportInsets
}

func (m *gutterMetrics) MeasureViewport(n *parser.Instance, _ float64) (ViewportInsets, error) {
	if m.calls != nil {
		m.calls[n.Path]++
	}
	if m.invalid != nil {
		return *m.invalid, nil
	}
	v := ViewportInsets{}
	if scrolls(n, "y") {
		v.Right = 12
	}
	if scrolls(n, "x") {
		v.Bottom = 10
	}
	return v, nil
}
func gutterLayout(t *testing.T, source string, size Size, offsets map[string]runtime.ViewportState) *SnapshotLayout {
	t.Helper()
	g, e := (&Engine{Measure: &gutterMetrics{}}).LayoutSnapshot(runtime.Snapshot{Root: compileProfile(t, "0.3", source), Viewports: offsets}, size)
	if e != nil {
		t.Fatal(e)
	}
	return g
}
func TestMeasuredGuttersPaddingCornerAndContentHits(t *testing.T) {
	g := gutterLayout(t, `[a=button("A") {scale=2}] {scale=1,padding=0.1,gap=0,overflow-x=scroll,overflow-y=scroll}`, Size{100, 80}, nil)
	v := g.Viewports["page"]
	if v.Rect != (Rect{10, 8, 68, 54}) || v.VerticalGutter != (Rect{78, 8, 12, 54}) || v.HorizontalGutter != (Rect{10, 62, 68, 10}) {
		t.Fatalf("padding/strips: %+v", v)
	}
	if v.VerticalGutter.Contains(78, 62) || v.HorizontalGutter.Contains(78, 62) {
		t.Fatal("corner belongs to a track")
	}
	if g.Root.Hit(79, 10) != nil || g.Root.Hit(10, 63) != nil || g.Root.Hit(77, 61) == nil {
		t.Fatal("children painted/hit ancestor gutter")
	}
	next, rest, e := g.RouteScroll(80, 10, 0, 12)
	if e != nil || next["page"].Y != 12 || rest.Y != 0 {
		t.Fatalf("own gutter routing: %+v %+v %v", next, rest, e)
	}
	next, rest, e = g.RouteScroll(20, 65, 10, 0)
	if e != nil || next["page"].X != 10 || rest.X != 0 {
		t.Fatalf("horizontal gutter routing: %+v %+v %v", next, rest, e)
	}
	_, rest, e = g.RouteScroll(80, 65, 10, 12)
	if e != nil || rest != (runtime.ViewportState{X: 10, Y: 12}) {
		t.Fatalf("corner routed: %+v %v", rest, e)
	}
}

func TestNestedGuttersDisjointAndClippedByAncestorContent(t *testing.T) {
	g := gutterLayout(t, `[inner=<a=button("A") {scale-y=2,x=fill}> {x=fill,scale-y=2,gap=0,overflow-y=scroll}] {scale=1,gap=0,overflow-y=scroll}`, Size{200, 100}, map[string]runtime.ViewportState{"page": {Y: 30}, "page/inner": {Y: 50}})
	outer, inner := g.Viewports["page"], g.Viewports["page/inner"]
	if outer.Rect != (Rect{0, 0, 188, 100}) || outer.VerticalGutter != (Rect{188, 0, 12, 100}) {
		t.Fatalf("outer %+v", outer)
	}
	if inner.Rect != (Rect{0, -30, 176, 200}) || inner.VerticalGutter != (Rect{176, 0, 12, 100}) {
		t.Fatalf("inner own offset moved/clipped strip: %+v", inner)
	}
	if inner.VerticalGutter.Intersect(outer.VerticalGutter).W != 0 {
		t.Fatal("overlapping gutters")
	}
	if got := boxes(g.Root)["page/inner/a"].Rect; got != (Rect{0, -80, 176, 400}) {
		t.Fatalf("finite remaining reference/double transform: %+v", got)
	}
	out, rest, e := g.RouteScroll(180, 20, 0, 10)
	if e != nil || out["page/inner"].Y != 60 || out["page"].Y != 30 || rest.Y != 0 {
		t.Fatalf("inner gutter hit outer: %+v %+v %v", out, rest, e)
	}
	out, rest, e = g.RouteScroll(194, 20, 0, 10)
	if e != nil || out["page/inner"].Y != 50 || out["page"].Y != 40 || rest.Y != 0 {
		t.Fatalf("outer gutter hit inner: %+v %+v %v", out, rest, e)
	}
	// Full track height stays 200 even though the visible strip is only 100.
	if inner.Rect.H != 200 || inner.VerticalGutter.H != 100 {
		t.Fatal("clipping shortened track geometry")
	}
}

func TestGuttersChargedOnceContentSizingAndStableZeroRange(t *testing.T) {
	m := &gutterMetrics{calls: map[string]int{}}
	root := compileProfile(t, "0.3", `[a=button("A")] {y=fill,gap=0,overflow-y=scroll}`)
	g, e := (&Engine{Measure: m}).LayoutSnapshot(runtime.Snapshot{Root: root}, Size{200, 100})
	if e != nil {
		t.Fatal(e)
	}
	if g.Root.Rect.W != 32 || g.Viewports["page"].Rect.W != 20 || m.calls["page"] != 1 {
		t.Fatalf("content width/call cost: %+v %v", g.Root.Rect, m.calls)
	}
	if g.Viewports["page"].Maximum.Y != 0 || g.Viewports["page"].VerticalGutter.W != 12 {
		t.Fatal("zero range discarded declared gutter")
	}
	g = gutterLayout(t, `[a=button("A")] {x=fill,gap=0,overflow-x=scroll}`, Size{200, 100}, nil)
	if g.Root.Rect.H != 30 || g.Viewports["page"].Rect.H != 20 {
		t.Fatalf("content height charged incorrectly %+v", g.Root.Rect)
	}
}

func TestGuttersRejectInvalidAxesAndExhaustedSpace(t *testing.T) {
	source := `[] {scale=1,overflow-y=scroll}`
	for _, v := range []ViewportInsets{{Right: -1}, {Right: math.NaN()}, {Right: math.Inf(1)}, {Right: 1e7 + 1}, {Bottom: 1}, {Right: 100}, {Right: 101}} {
		root := compileProfile(t, "0.3", source)
		if _, e := (&Engine{Measure: &gutterMetrics{invalid: &v}}).Layout(root, Size{100, 80}); e == nil {
			t.Errorf("invalid gutter accepted %+v", v)
		}
	}
	v := ViewportInsets{Right: 1}
	root := compileProfile(t, "0.3", `[] {scale=1,overflow-x=scroll}`)
	if _, e := (&Engine{Measure: &gutterMetrics{invalid: &v}}).Layout(root, Size{100, 80}); e == nil {
		t.Fatal("right on scroll-x accepted")
	}
}

func TestGuttersLegacyAndAbsentAdjunctCompatibility(t *testing.T) {
	m := &gutterMetrics{calls: map[string]int{}}
	root := compileProfile(t, "0.2", `[button("A")] {scale=1,gap=0}`)
	plain, e := (&Engine{Measure: fixtureMetrics{}}).Layout(root, Size{200, 100})
	if e != nil {
		t.Fatal(e)
	}
	measured, e := (&Engine{Measure: m}).Layout(root, Size{200, 100})
	if e != nil {
		t.Fatal(e)
	}
	if !reflect.DeepEqual(plain, measured) || len(m.calls) != 0 {
		t.Fatal("0.2 geometry called/changed by adjunct")
	}
	root = compileProfile(t, "0.2", `[button("A")] {scale=1,overflow-y=scroll}`)
	_, e = (&Engine{Measure: m}).Layout(root, Size{200, 100})
	if d, ok := e.(*parser.Diagnostic); !ok || d.Code != "unsupported-scroll" || len(m.calls) != 0 {
		t.Fatalf("0.2 scroll changed: %v", e)
	}
	g := snapshotGeometry(t, `[button("A") {scale=2}] {scale=1,overflow-x=scroll,overflow-y=scroll}`, nil)
	if g.Viewports["page"].Rect.W != 200 || g.Viewports["page"].VerticalGutter != (Rect{}) {
		t.Fatal("absent adjunct inferred native chrome")
	}
}

func TestCollectionGuttersUseExistingMetricsOnce(t *testing.T) {
	root := compileProfile(t, "0.3", `[rows=list("Rows") {scale=1,overflow-x=scroll,overflow-y=scroll}] {scale=1,gap=0}`)
	g, e := (&Engine{Measure: collectionOracle{content: Size{600, 600}}}).LayoutSnapshot(runtime.Snapshot{Root: root}, Size{200, 100})
	if e != nil {
		t.Fatal(e)
	}
	v := g.Viewports["page/rows"]
	if v.Rect != (Rect{0, 20, 190, 80}) || v.VerticalGutter != (Rect{190, 20, 10, 80}) || v.HorizontalGutter != (Rect{}) {
		t.Fatalf("collection title/chrome changed %+v", v)
	}
	out, rest, e := g.RouteScroll(195, 30, 0, 20)
	if e != nil || rest.Y != 0 || out["page/rows"].Y != 20 {
		t.Fatalf("collection gutter skipped own viewport: %+v %+v %v", out, rest, e)
	}
}

func TestGuttersResizeFailureAndContentShrinkStateGate(t *testing.T) {
	root := compileProfile(t, "0.3", `[a=button("A") {scale=2}] {scale=1,gap=0,overflow-x=scroll,overflow-y=scroll}`)
	session, err := runtime.New("gutter-gate", root)
	if err != nil {
		t.Fatal(err)
	}
	size := Size{200, 100}
	var prepared *SnapshotLayout
	if err = session.CheckStateWith(func(s runtime.Snapshot) (map[string]runtime.ViewportState, error) {
		candidate, e := (&Engine{Measure: &gutterMetrics{}}).LayoutSnapshot(s, size)
		if e != nil {
			return nil, e
		}
		prepared = candidate
		return candidate.EffectiveOffsets(), nil
	}); err != nil {
		t.Fatal(err)
	}
	if err = session.SetViewports(map[string]runtime.ViewportState{"page": {X: 999, Y: 999}}); err != nil {
		t.Fatal(err)
	}
	if got := session.Snapshot().Viewports["page"]; got != (runtime.ViewportState{X: 188, Y: 90}) {
		t.Fatalf("initial reduced viewport clamp: %+v", got)
	}
	before, oldGeometry := session.Snapshot(), prepared
	size = Size{12, 100} // The right gutter consumes the entire assigned width.
	if err = session.SetViewports(map[string]runtime.ViewportState{"page": {X: 10}}); err == nil {
		t.Fatal("exhausted candidate accepted")
	}
	if !reflect.DeepEqual(before, session.Snapshot()) || prepared != oldGeometry {
		t.Fatal("failed gutter measurement changed accepted state/geometry")
	}
	size = Size{120, 60}
	if err = session.SetViewports(before.Viewports); err != nil {
		t.Fatal(err)
	}
	v := prepared.Viewports["page"]
	if v.Rect != (Rect{0, 0, 108, 50}) || v.Offset != (runtime.ViewportState{X: 108, Y: 50}) || v.VerticalGutter != (Rect{108, 0, 12, 50}) {
		t.Fatalf("resize geometry/clamp: %+v", v)
	}
	w, _ := session.Widget("page/a")
	if err = session.Apply(session.Revision, session.BatchRevision+1, []runtime.Update{{Handle: w.Handle, Property: runtime.Visible, Value: runtime.Value{Kind: runtime.Boolean, Bool: false}}}); err != nil {
		t.Fatal(err)
	}
	v = prepared.Viewports["page"]
	if session.Snapshot().Viewports["page"] != (runtime.ViewportState{}) || v.Maximum != (runtime.ViewportState{}) || v.VerticalGutter != (Rect{108, 0, 12, 50}) || v.HorizontalGutter != (Rect{0, 50, 108, 10}) {
		t.Fatalf("content shrink changed reservation or failed clamp: %+v", v)
	}
	_, rest, err := prepared.RouteScroll(110, 10, 0, 20)
	if err != nil || rest.Y != 20 {
		t.Fatalf("zero range consumed delta: %+v %v", rest, err)
	}
}

func TestGuttersLeaveFiniteFrTracks(t *testing.T) {
	g := gutterLayout(t, `[a=button("A") {x=1fr,y=fill}, b=button("B") {x=3fr,y=fill}] {scale=1,gap=0,overflow-x=scroll,overflow-y=scroll}`, Size{200, 100}, nil)
	b := boxes(g.Root)
	if b["page/a"].Rect != (Rect{0, 0, 47, 90}) || b["page/b"].Rect != (Rect{47, 0, 141, 90}) {
		t.Fatalf("fr tracks used outer area: %+v %+v", b["page/a"].Rect, b["page/b"].Rect)
	}
}

type gutterCollectionMetrics struct{ gutterMetrics }

func (m *gutterCollectionMetrics) MeasureCollection(n *parser.Instance, font float64, outer Size) (CollectionMetrics, error) {
	return (collectionOracle{content: Size{600, 600}, alter: func(v *CollectionMetrics) {
		v.Viewport.H -= 10
	}}).MeasureCollection(n, font, outer)
}

func TestCollectionGuttersExcludeTitleCornerAndAncestorStrips(t *testing.T) {
	m := &gutterCollectionMetrics{gutterMetrics{calls: map[string]int{}}}
	root := compileProfile(t, "0.3", `[rows=list("Rows") {scale=2,overflow-x=scroll,overflow-y=scroll}] {scale=1,gap=0,overflow-x=scroll,overflow-y=scroll}`)
	g, err := (&Engine{Measure: m}).LayoutSnapshot(runtime.Snapshot{Root: root, Viewports: map[string]runtime.ViewportState{"page": {X: 188, Y: 90}, "page/rows": {X: 30, Y: 40}}}, Size{200, 100})
	if err != nil {
		t.Fatal(err)
	}
	v := g.Viewports["page/rows"]
	if v.Rect != (Rect{-188, -70, 366, 150}) || v.VerticalGutter != (Rect{178, 0, 10, 80}) || v.HorizontalGutter != (Rect{0, 80, 178, 10}) {
		t.Fatalf("collection strips/title/ancestor clipping: %+v", v)
	}
	if len(m.calls) != 1 || m.calls["page"] != 1 {
		t.Fatalf("collection invoked container adjunct: %v", m.calls)
	}
	if v.VerticalGutter.Contains(180, 85) || v.HorizontalGutter.Contains(180, 85) {
		t.Fatal("collection corner is a track")
	}
	out, rest, err := g.RouteScroll(100, 85, 10, 0)
	if err != nil || rest.X != 0 || out["page/rows"].X != 40 || out["page"].X != 188 {
		t.Fatalf("collection horizontal gutter routing: %+v %+v %v", out, rest, err)
	}
}
