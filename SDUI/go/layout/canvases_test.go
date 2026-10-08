package layout

import (
	"encoding/json"
	"math"
	"reflect"
	"testing"

	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

type canvasMetrics struct {
	*paneMetrics
	menu      Size
	menuCalls int
}

func (m *canvasMetrics) MeasureMenu(_ *parser.Instance, _ float64) (Size, error) {
	m.menuCalls++
	return m.menu, nil
}
func canvasMeasure() *canvasMetrics {
	return &canvasMetrics{paneMetrics: defaultPaneMetrics(), menu: Size{60, 22}}
}
func canvasSnapshot(t *testing.T, source string, open ...string) runtime.Snapshot {
	t.Helper()
	root := compileProfile(t, "0.3", source)
	s := runtime.Snapshot{Root: root, Surfaces: map[string]runtime.SurfaceState{}, Viewports: map[string]runtime.ViewportState{}}
	root.Walk(func(n *parser.Instance) {
		if dialog(n) {
			s.Surfaces[n.Path] = runtime.SurfaceState{InstancePath: n.Path}
		}
	})
	for _, path := range open {
		state := s.Surfaces[path]
		state.Open = true
		s.Surfaces[path] = state
	}
	return s
}
func TestAuxiliaryDeclarationsHaveNoTracksOrPopupBoxes(t *testing.T) {
	source := `[cmd=command("Toggle",toggle=true); a=button("A");ctx=menu("Context",mode="context",target="a")[item(command="cmd")];d=dialog("Dialog")[bad=button("Poison")];bar=menu("File")[item(command="cmd");separator();menuGroup("Group")[item(command="cmd")];menu("More",mode="submenu")[item(command="cmd")]];b=button("B")] {scale=1,gap=0.1}`
	s := canvasSnapshot(t, source)
	m := canvasMeasure()
	g, err := (&Engine{Measure: m}).LayoutSnapshot(s, Size{200, 100})
	if err != nil {
		t.Fatal(err)
	}
	got := boxes(g.Root)
	if len(got) != 4 || got["page/a"].Rect.Y != 0 || got["page/bar"].Rect != (Rect{0, 30, 60, 22}) || got["page/b"].Rect.Y != 62 {
		t.Fatalf("auxiliary declaration charged tracks/gaps %+v", got)
	}
	if len(got["page/bar"].Children) != 0 || m.menuCalls == 0 {
		t.Fatal("menu bar not native measured leaf")
	}
	if _, err = g.EnsureVisible("page/d/bad", Rect{}); err == nil {
		t.Fatal("closed surface entered main routing")
	}
}

const surfaceGeometrySource = `[main=button("Main");first=dialog("First",modal=false)[scroll=[a=button("A") {scale-y=2,x=fill}] {scale=1,gap=0,overflow-y=scroll};nested=dialog("Nested")[s=split(axis="horizontal")[left=button("A");right=button("B")] {scale=1}] {scale=0.5,gap=0}] {scale=0.5,gap=0,font=24}] {scale=1,gap=0,font=20}`

func TestSeparateSurfaceCoordinatesAggregateAndAcceptedResize(t *testing.T) {
	s := canvasSnapshot(t, surfaceGeometrySource, "page/first", "page/first/nested")
	s.Viewports["page/first/scroll"] = runtime.ViewportState{Y: 500}
	original, _ := json.Marshal(s)
	e := &Engine{Measure: canvasMeasure()}
	initial, err := e.SurfaceSize(s, "page/first", Size{400, 300})
	if err != nil {
		t.Fatal(err)
	}
	if initial != (Size{200, 150}) {
		t.Fatalf("initial relative size %+v", initial)
	}
	nested, err := e.SurfaceSize(s, "page/first/nested", initial)
	if err != nil {
		t.Fatal(err)
	}
	if nested != (Size{100, 75}) {
		t.Fatalf("nested opening reference %+v", nested)
	}
	accepted := map[string]Size{"page/first": initial, "page/first/nested": nested}
	g, err := e.LayoutCanvases(s, Size{400, 300}, accepted)
	if err != nil {
		t.Fatal(err)
	}
	if len(boxes(g.Main.Root)) != 2 || len(g.Surfaces) != 2 {
		t.Fatal("surface content leaked into main")
	}
	first, second := g.Surfaces["page/first"], g.Surfaces["page/first/nested"]
	if first.Root.Rect != (Rect{0, 0, 200, 150}) || second.Root.Rect != (Rect{0, 0, 100, 75}) || second.Root.Font != 24 {
		t.Fatalf("canvas scale/origin/font %+v %+v", first.Root, second.Root)
	}
	if boxes(first.Root)["page/first/nested"] != nil {
		t.Fatal("nested surface charged its parent's geometry")
	}
	v := first.Viewports["page/first/scroll"]
	if v.Parent != "" || v.Rect != (Rect{0, 0, 200, 150}) || v.Offset.Y != 150 {
		t.Fatalf("canvas viewport %+v", v)
	}
	p := g.PresentationState()
	if p.Viewports["page/first/scroll"].Y != 150 || len(p.Splits) != 1 {
		t.Fatalf("aggregate gate %+v", p)
	}
	p.Viewports["page/first/scroll"] = runtime.ViewportState{}
	if first.Viewports["page/first/scroll"].Offset.Y != 150 {
		t.Fatal("aggregate result aliases geometry")
	}
	// Parent resizing cannot recompute an accepted nonmodal size or its offset.
	resized, err := e.LayoutCanvases(s, Size{800, 600}, accepted)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(first, resized.Surfaces["page/first"]) {
		t.Fatal("parent resize changed accepted child canvas")
	}
	accepted["page/first"] = Size{260, 180}
	resized, err = e.LayoutCanvases(s, Size{800, 600}, accepted)
	if err != nil {
		t.Fatal(err)
	}
	if resized.Surfaces["page/first"].Root.Rect.W != 260 || resized.Surfaces["page/first"].Viewports["page/first/scroll"].Offset.Y != 180 || resized.Surfaces["page/first/nested"].Root.Rect.W != 100 {
		t.Fatal("actual child resize reapplied opening scale or resized sibling")
	}
	after, _ := json.Marshal(s)
	if string(original) != string(after) {
		t.Fatal("canvas measurement mutated snapshot")
	}
}
func TestInvalidSurfaceGeometryReturnsNoAggregate(t *testing.T) {
	s := canvasSnapshot(t, surfaceGeometrySource, "page/first", "page/first/nested")
	e := &Engine{Measure: canvasMeasure()}
	original, _ := json.Marshal(s)
	for _, size := range []Size{{}, {-1, 50}, {math.NaN(), 50}, {50, math.Inf(1)}, {32769, 50}, {20, 20}} {
		g, err := e.LayoutCanvases(s, Size{400, 300}, map[string]Size{"page/first": {200, 150}, "page/first/nested": size})
		if err == nil || g != nil {
			t.Fatalf("invalid surface published aggregate %+v", size)
		}
	}
	for _, sizes := range []map[string]Size{nil, {"unknown": {100, 100}}} {
		if g, err := e.LayoutCanvases(s, Size{400, 300}, sizes); err == nil || g != nil {
			t.Fatal("missing/unknown canvas admitted")
		}
	}
	after, _ := json.Marshal(s)
	if string(original) != string(after) {
		t.Fatal("failed geometry mutated accepted input")
	}
}
func TestOpeningBoundsAppliedOnceAndClosedSizesIgnored(t *testing.T) {
	s := canvasSnapshot(t, `[d=dialog("D")[button("A") {scale=1}] {x=fill,y=fill,max-x=0.8,max-y=0.5,gap=0}] {scale=1}`, "page/d")
	e := &Engine{Measure: canvasMeasure()}
	size, err := e.SurfaceSize(s, "page/d", Size{400, 300})
	if err != nil {
		t.Fatal(err)
	}
	if size != (Size{320, 150}) {
		t.Fatalf("opening bounds %+v", size)
	}
	g, err := e.LayoutCanvases(s, Size{400, 300}, map[string]Size{"page/d": size})
	if err != nil {
		t.Fatal(err)
	}
	if g.Surfaces["page/d"].Root.Rect != (Rect{0, 0, 320, 150}) {
		t.Fatal("source maximum applied twice")
	}
	state := s.Surfaces["page/d"]
	state.Open = false
	s.Surfaces["page/d"] = state
	g, err = e.LayoutCanvases(s, Size{400, 300}, map[string]Size{"page/d": {}})
	if err != nil || len(g.Surfaces) != 0 {
		t.Fatalf("closed cache entry measured %v", err)
	}
}

func TestEmptySurfaceNaturalSizeNeedsAcceptedNativeMinimum(t *testing.T) {
	s := canvasSnapshot(t, `[d=dialog("Empty")[]] {scale=1}`, "page/d")
	e := &Engine{Measure: canvasMeasure()}
	size, err := e.SurfaceSize(s, "page/d", Size{400, 300})
	if err != nil || size != (Size{}) {
		t.Fatalf("empty natural size %+v %v", size, err)
	}
	if _, err = e.LayoutCanvases(s, Size{400, 300}, map[string]Size{"page/d": size}); err == nil {
		t.Fatal("unallocated zero canvas accepted")
	}
	g, err := e.LayoutCanvases(s, Size{400, 300}, map[string]Size{"page/d": {1, 1}})
	if err != nil || g.Surfaces["page/d"].Root.Rect != (Rect{0, 0, 1, 1}) {
		t.Fatalf("host accepted minimum %v", err)
	}
}
func TestNativeMenuMetricsRejectWithoutAffectingLegacy(t *testing.T) {
	s := canvasSnapshot(t, `[c=command("C",toggle=true);bar=menu("File")[item(command="c")]] {scale=1}`)
	if _, err := (&Engine{Measure: defaultPaneMetrics()}).LayoutSnapshot(s, Size{200, 100}); err == nil {
		t.Fatal("guessed native menu metrics")
	}
	for _, size := range []Size{{0, 22}, {60, 0}, {-1, 22}, {math.NaN(), 22}, {60, math.Inf(1)}, {1e7 + 1, 22}} {
		m := canvasMeasure()
		m.menu = size
		if _, err := (&Engine{Measure: m}).LayoutSnapshot(s, Size{200, 100}); err == nil {
			t.Fatalf("invalid menu metrics %+v", size)
		}
	}
	tooSmall := canvasSnapshot(t, `[c=command("C",toggle=true);bar=menu("File")[item(command="c")] {max-x=0.1}] {scale=1}`)
	if _, err := (&Engine{Measure: canvasMeasure()}).LayoutSnapshot(tooSmall, Size{200, 100}); err == nil {
		t.Fatal("source maximum bypassed native bar minimum")
	}
	root := compileProfile(t, "0.2", `[button("A"),input("B")] {scale=1,gap=0}`)
	legacy, err := (&Engine{Measure: fixtureMetrics{}}).Layout(root, Size{200, 100})
	if err != nil {
		t.Fatal(err)
	}
	m := canvasMeasure()
	g, err := (&Engine{Measure: m}).LayoutCanvases(runtime.Snapshot{Root: root}, Size{200, 100}, nil)
	if err != nil || !reflect.DeepEqual(legacy, g.Main.Root) || m.menuCalls != 0 || len(g.Surfaces) != 0 {
		t.Fatalf("0.2 changed: %v", err)
	}
}
