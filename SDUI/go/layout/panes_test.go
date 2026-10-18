package layout

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

type paneMetrics struct {
	fixtureMetrics
	header   Size
	divider  float64
	sizes    map[string]Size
	measured map[string]int
}

func (m *paneMetrics) Measure(n *parser.Instance, font, width float64) (Size, error) {
	if m.measured != nil {
		m.measured[n.Path]++
	}
	if n.Argument("label") == "Poison" {
		return Size{}, fmt.Errorf("inactive content was measured")
	}
	if s, ok := m.sizes[n.Argument("label")]; ok {
		return s, nil
	}
	return m.fixtureMetrics.Measure(n, font, width)
}
func (m *paneMetrics) MeasureTabs(_ *parser.Instance, _ float64, _ Size) (TabsMetrics, error) {
	return TabsMetrics{Header: m.header}, nil
}
func (m *paneMetrics) MeasureSplit(_ *parser.Instance, _ float64) (float64, error) {
	return m.divider, nil
}
func defaultPaneMetrics() *paneMetrics { return &paneMetrics{header: Size{80, 24}, divider: 10} }
func paneSnapshot(t *testing.T, source string) runtime.Snapshot {
	t.Helper()
	root := compileProfile(t, "0.3", source)
	s, err := runtime.New("pane-geometry", root)
	if err != nil {
		t.Fatal(err)
	}
	return s.Snapshot()
}
func measuredPanes(t *testing.T, s runtime.Snapshot, size Size, m *paneMetrics) *SnapshotLayout {
	t.Helper()
	g, err := (&Engine{Measure: m}).LayoutSnapshot(s, size)
	if err != nil {
		t.Fatal(err)
	}
	return g
}
func TestTabsOnlyActiveBodyAndHeaderGeometry(t *testing.T) {
	s := paneSnapshot(t, `[t=tabs("Workspace")[one=page("One")[a=button("A") {scale=1}];two=page("Two")[bad=button("Poison")]] {scale=1}] {scale=1,gap=0}`)
	before := s.Root.Rows[0][0].Rows[1][0].Layout["visible"]
	m := defaultPaneMetrics()
	m.measured = map[string]int{}
	g := measuredPanes(t, s, Size{300, 200}, m)
	tab := g.Tabs["page/t"]
	if tab.Header != (Rect{0, 0, 300, 24}) || tab.HeaderClip != tab.Header || tab.Body != (Rect{0, 24, 300, 176}) || tab.Selected != "one" {
		t.Fatalf("header/body %+v", tab)
	}
	b := boxes(g.Root)
	if b["page/t/two"] != nil || b["page/t/two/bad"] != nil || m.measured["page/t/two/bad"] != 0 {
		t.Fatal("inactive subtree measured or laid out")
	}
	if b["page/t/one/a"].Rect != (Rect{0, 24, 300, 176}) || g.Root.Hit(10, 10) != nil || g.Root.Hit(10, 30).Path != "page/t/one/a" {
		t.Fatal("header/body hit separation")
	}
	if before != s.Root.Rows[0][0].Rows[1][0].Layout["visible"] {
		t.Fatal("geometry mutated snapshot activity")
	}
	if _, err := g.EnsureVisible("page/t/two/bad", Rect{}); err == nil {
		t.Fatal("inactive target revealed")
	}
	// A runtime-selected empty body is not replaced by a source/default selection.
	state := s.Tabs["page/t"]
	state.Selected = ""
	s.Tabs["page/t"] = state
	g = measuredPanes(t, s, Size{300, 200}, m)
	if len(boxes(g.Root)["page/t"].Children) != 0 || g.Tabs["page/t"].Selected != "" {
		t.Fatal("empty body fell back to source")
	}
}
func TestSplitMeasuredBoundsAndExactRects(t *testing.T) {
	s := paneSnapshot(t, `[s=split(axis="horizontal",proportion=0.1,minFirst=0.1,minSecond=0.2)[a=button("A");b=button("B")] {scale=1}] {scale=1,gap=0}`)
	m := defaultPaneMetrics()
	m.sizes = map[string]Size{"A": {60, 20}, "B": {40, 20}}
	before := s.Splits["page/s"]
	g := measuredPanes(t, s, Size{210, 100}, m)
	split := g.Splits["page/s"]
	if split.Geometry != (runtime.SplitGeometry{Lower: .3, Upper: .8, Effective: .3}) || split.First != (Rect{0, 0, 60, 100}) || split.Divider != (Rect{60, 0, 10, 100}) || split.Second != (Rect{70, 0, 140, 100}) {
		t.Fatalf("split %+v", split)
	}
	if g.PresentationState().Splits["page/s"] != split.Geometry || s.Splits["page/s"] != before {
		t.Fatal("gate result missing or snapshot mutated")
	}
	if g.Root.Hit(65, 50) != nil || g.Root.Hit(75, 50).Path != "page/s/b" {
		t.Fatal("divider hit reached child")
	}
	state := g.PresentationState()
	state.Splits["page/s"] = runtime.SplitGeometry{}
	if g.Splits["page/s"].Geometry.Effective != .3 {
		t.Fatal("gate result aliases geometry")
	}
}
func TestSplitCollapseSkipsMeasurementAndSuspendsRelativeMinima(t *testing.T) {
	s := paneSnapshot(t, `[s=split(axis="vertical",minFirst=0.4,minSecond=0.4)[a=button("Poison");b=button("B")] {scale=1}] {scale=1,gap=0}`)
	v := s.Splits["page/s"]
	v.Collapsed = runtime.SplitFirst
	v.SavedProportion = .65
	s.Splits["page/s"] = v
	m := defaultPaneMetrics()
	m.sizes = map[string]Size{"B": {20, 30}}
	for _, height := range []float64{40, 100} {
		g := measuredPanes(t, s, Size{100, height}, m)
		split := g.Splits["page/s"]
		if split.First.H != 0 || split.Second.H != height-10 || split.Divider != (Rect{0, 0, 100, 10}) || len(g.PresentationState().Splits) != 0 {
			t.Fatalf("collapsed %+v", split)
		}
		if boxes(g.Root)["page/s/a"] != nil || s.Splits["page/s"].SavedProportion != .65 {
			t.Fatal("collapsed body/state changed")
		}
	}
	if _, err := (&Engine{Measure: m}).LayoutSnapshot(s, Size{100, 39}); err == nil {
		t.Fatal("visible minimum ignored")
	}
}
func TestDeepNestedSplitAndTabsMinima(t *testing.T) {
	s := paneSnapshot(t, `[s=split(axis="horizontal",proportion=0.1)[a=tabs("Tabs")[p=page("Page")[deep=split(axis="vertical",minFirst=0.4,minSecond=0.4)[top=button("Tall");bottom=button("B")] {scale=1}]];b=button("B")] {scale=1}] {scale=1,gap=0}`)
	m := defaultPaneMetrics()
	m.sizes = map[string]Size{"Tall": {60, 60}, "B": {20, 20}}
	g := measuredPanes(t, s, Size{210, 134}, m)
	// Inner split needs U=100: max(.4U,60)+max(.4U,20)=100;
	// add divider 10 and native tabs header 24 => exact minimum height 134.
	outer := g.Splits["page/s"]
	inner := g.Splits["page/s/a/p/deep"]
	near(t, outer.Geometry.Lower, .4)
	near(t, outer.First.W, 80)
	near(t, inner.Geometry.Lower, .6)
	near(t, inner.Geometry.Upper, .6)
	if g.Tabs["page/s/a"].Body.H != 110 || inner.First.H != 60 || inner.Second.H != 40 {
		t.Fatalf("deep %+v %+v", g.Tabs, inner)
	}
	if _, err := (&Engine{Measure: m}).LayoutSnapshot(s, Size{210, 133}); err == nil {
		t.Fatal("deep measured minimum ignored")
	}
}
func TestPaneInvalidMetricsAndMissingAdapter(t *testing.T) {
	sources := []string{`[t=tabs("T")[p=page("P")[]] {scale=1}] {scale=1}`, `[s=split(axis="horizontal")[a=[];b=[]] {scale=1}] {scale=1}`}
	for _, source := range sources {
		s := paneSnapshot(t, source)
		if _, err := (&Engine{Measure: fixtureMetrics{}}).LayoutSnapshot(s, Size{200, 100}); err == nil {
			t.Fatal("guessed native chrome")
		}
		for _, value := range []float64{0, -1, math.NaN(), math.Inf(1), 1e7 + 1} {
			m := defaultPaneMetrics()
			m.header.H = value
			m.divider = value
			if _, err := (&Engine{Measure: m}).LayoutSnapshot(s, Size{200, 100}); err == nil {
				t.Errorf("accepted metric %v", value)
			}
		}
	}
}
func TestPaneAdjunctPreservesExactLegacySourceGeometry(t *testing.T) {
	sources := []string{`[a=button("A"), b=input("B")] {scale=1,gap=0}`, `[a=[] {x=1fr,min-x=0.2}, b=[] {x=3fr}] {scale=1,gap=0}`, `[a=button("A") {scale-x=2}] {scale=1,overflow-x=clip}`}
	for _, source := range sources {
		root := compileProfile(t, "0.2", source)
		want, err := (&Engine{Measure: fixtureMetrics{}}).Layout(root, Size{200, 100})
		if err != nil {
			t.Fatal(err)
		}
		got, err := (&Engine{Measure: defaultPaneMetrics()}).Layout(root, Size{200, 100})
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("legacy geometry changed for exact source %s", source)
		}
	}
}

func TestNestedRelativeMinimumUsesOwnBody(t *testing.T) {
	for _, minimum := range []string{"0.8", "0.99", "1"} {
		s := paneSnapshot(t, `[s=split(axis="horizontal",proportion=0.1)[left=[a=button("A") {x=fill,min-x=`+minimum+`}] {gap=0};right=button("B")] {scale=1}] {scale=1,gap=0}`)
		g := measuredPanes(t, s, Size{210, 100}, defaultPaneMetrics())
		// The .8 child minimum refers to the left frame, not the entire 200-unit
		// split plane. A 20-unit native button fits a 20-unit left frame.
		near(t, g.Splits["page/s"].Geometry.Lower, .1)
		near(t, g.Splits["page/s"].First.W, 20)
		near(t, boxes(g.Root)["page/s/left/a"].Rect.W, 20)
	}
}

func TestNestedRelativePaddingMinimum(t *testing.T) {
	s := paneSnapshot(t, `[s=split(axis="horizontal",proportion=0.1)[left=[inner=[a=button("A")] {padding=0.1,gap=0}] {gap=0};right=button("B")] {scale=1}] {scale=1,gap=0}`)
	g := measuredPanes(t, s, Size{210, 100}, defaultPaneMetrics())
	// Left width W must accommodate 20 + .2*W, hence W=25.
	near(t, g.Splits["page/s"].Geometry.Lower, .125)
	near(t, g.Splits["page/s"].First.W, 25)
	near(t, boxes(g.Root)["page/s/left/inner/a"].Rect.X, 2.5)
}

type paneCollectionMetrics struct{ *paneMetrics }

func (m paneCollectionMetrics) MeasureCollection(n *parser.Instance, font float64, outer Size) (CollectionMetrics, error) {
	return (collectionOracle{content: Size{600, 1000}}).MeasureCollection(n, font, outer)
}
func TestPaneCollectionUsesNativeMinimumNotLoadedExtent(t *testing.T) {
	root := compileProfile(t, "0.3", `[t=tabs("Tabs")[p=page("Page")[rows=list("Rows") {scale=1,overflow-x=scroll,overflow-y=scroll}]] {scale=1}] {scale=1,gap=0}`)
	g, err := (&Engine{Measure: paneCollectionMetrics{defaultPaneMetrics()}}).LayoutSnapshot(runtime.Snapshot{Root: root}, Size{100, 64})
	if err != nil {
		t.Fatal(err)
	}
	if g.Tabs["page/t"].Body.H != 40 || g.Viewports["page/t/p/rows"].Maximum.Y != 980 {
		t.Fatalf("collection minimum/extent %+v", g.Viewports)
	}
}

func TestNestedFrMinimumBothAxes(t *testing.T) {
	for _, axis := range []string{"horizontal", "vertical"} {
		rows := `a=button("A") {x=1fr},b=button("B") {x=3fr}`
		if axis == "vertical" {
			rows = `a=button("A") {y=1fr};b=button("B") {y=3fr}`
		}
		s := paneSnapshot(t, `[s=split(axis="`+axis+`",proportion=0.1)[left=[`+rows+`] {gap=0};right=button("B")] {scale=1}] {scale=1,gap=0}`)
		g := measuredPanes(t, s, Size{210, 210}, defaultPaneMetrics())
		near(t, g.Splits["page/s"].Geometry.Lower, .4)
		first := boxes(g.Root)["page/s/left/a"].Rect
		if axis == "horizontal" {
			near(t, first.W, 20)
		} else {
			near(t, first.H, 20)
		}
	}
}
