package layout

import (
	"reflect"
	"testing"

	"github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

func TestPaneGateRetainsInactiveScrollAndClampsOnReveal(t *testing.T) {
	root := compileProfile(t, "0.3", `[t=tabs("Tabs")[one=page("One")[scroller=[a=button("A") {scale-y=2,x=fill}] {scale=1,gap=0,overflow-y=scroll}];two=page("Two")[]] {scale=1}] {scale=1,gap=0}`)
	s, err := runtime.New("inactive-scroll", root)
	if err != nil {
		t.Fatal(err)
	}
	size := Size{200, 200}
	m := defaultPaneMetrics()
	var measured *SnapshotLayout
	gate := func(snapshot runtime.Snapshot) (runtime.PresentationState, error) {
		g, err := (&Engine{Measure: m}).LayoutSnapshot(snapshot, size)
		if err != nil {
			return runtime.PresentationState{}, err
		}
		measured = g
		return g.PresentationState(), nil
	}
	if err = s.CheckPresentationWith(gate); err != nil {
		t.Fatal(err)
	}
	path := "page/t/one/scroller"
	if err = s.SetViewports(map[string]runtime.ViewportState{path: {Y: 150}}); err != nil {
		t.Fatal(err)
	}
	handle, ok := s.Pane("page/t")
	if !ok {
		t.Fatal("missing tabs")
	}
	if err = s.SelectPage(handle, "two"); err != nil {
		t.Fatal(err)
	}
	if s.Snapshot().Viewports[path].Y != 150 {
		t.Fatal("runtime lost inactive offset")
	}
	if _, ok = measured.Viewports[path]; ok {
		t.Fatal("layout returned inactive viewport")
	}
	if _, err = measured.EnsureVisible("page/t/one/scroller/a", Rect{}); err == nil {
		t.Fatal("inactive target eligible")
	}
	size = Size{200, 100}
	if err = s.SelectPage(handle, "one"); err != nil {
		t.Fatal(err)
	}
	if s.Snapshot().Viewports[path].Y != 76 || measured.Viewports[path].Offset.Y != 76 {
		t.Fatalf("reveal clamp diverged %+v %+v", s.Snapshot().Viewports, measured.Viewports)
	}
	if measured.Viewports[path].Rect != (Rect{0, 24, 200, 76}) {
		t.Fatalf("header charged twice: %+v", measured.Viewports[path])
	}
}

func TestPaneGateCollapseResizeRestoreAndStrictRatio(t *testing.T) {
	root := compileProfile(t, "0.3", `[s=split(axis="horizontal",minFirst=0.1,minSecond=0.2)[a=button("A");b=button("B")] {scale=1}] {scale=1,gap=0}`)
	s, err := runtime.New("split-gate", root)
	if err != nil {
		t.Fatal(err)
	}
	size := Size{210, 100}
	m := defaultPaneMetrics()
	m.sizes = map[string]Size{"A": {60, 20}, "B": {40, 20}}
	var prepared *SnapshotLayout
	if err = s.CheckPresentationWith(func(snapshot runtime.Snapshot) (runtime.PresentationState, error) {
		candidate, err := (&Engine{Measure: m}).LayoutSnapshot(snapshot, size)
		if err != nil {
			return runtime.PresentationState{}, err
		}
		prepared = candidate
		return candidate.PresentationState(), nil
	}); err != nil {
		t.Fatal(err)
	}
	handle, ok := s.Pane("page/s")
	if !ok {
		t.Fatal("missing split")
	}
	before := s.Snapshot()
	if err = s.SetSplitProportion(handle, .2); err == nil {
		t.Fatal("strict proportion silently clamped")
	}
	if !reflect.DeepEqual(before, s.Snapshot()) {
		t.Fatal("strict rejection mutated accepted state")
	}
	if err = s.SetSplitProportion(handle, .7); err != nil {
		t.Fatal(err)
	}
	if err = s.CollapseSplit(handle, runtime.SplitFirst); err != nil {
		t.Fatal(err)
	}
	if prepared.Splits["page/s"].First.W != 0 || len(prepared.PresentationState().Splits) != 0 {
		t.Fatal("collapsed expanded geometry remains")
	}
	size = Size{90, 100} // B fits alone; 60+40+10 cannot fit expanded.
	if err = s.SetViewports(nil); err != nil {
		t.Fatal(err)
	}
	collapsed := s.Snapshot()
	if collapsed.Splits["page/s"].SavedProportion != .7 {
		t.Fatal("collapsed resize overwrote saved ratio")
	}
	if err = s.RestoreSplit(handle); err == nil {
		t.Fatal("unfit restore accepted")
	}
	if !reflect.DeepEqual(collapsed, s.Snapshot()) {
		t.Fatal("failed restore mutated collapsed state/focus")
	}
	size = Size{120, 100} // U=110, legal upper is 1-40/110, less than saved .7.
	if err = s.RestoreSplit(handle); err != nil {
		t.Fatal(err)
	}
	v := s.Snapshot().Splits["page/s"]
	near(t, v.Proportion, 1-40./110)
	if v.Collapsed != runtime.SplitNone || prepared.Splits["page/s"].Geometry.Effective != v.Proportion {
		t.Fatal("restore final state and geometry diverged")
	}
	// Visible disabled content still needs geometry; only input is ineligible.
	if err = s.Apply(s.Revision, s.BatchRevision+1, []runtime.Update{{Handle: handle, Property: runtime.Enabled, Value: runtime.Value{Kind: runtime.Boolean, Bool: false}}}); err != nil {
		t.Fatal(err)
	}
	if len(prepared.PresentationState().Splits) != 1 || prepared.Root.Hit(10, 10) != nil {
		t.Fatal("disabled split vanished or accepted input")
	}
}

func TestPaneChromeClipsAndNestedReveal(t *testing.T) {
	s := paneSnapshot(t, `[t=tabs("Tabs")[one=page("One")[scroller=[a=button("A") {scale-y=2,x=fill}] {scale=1,gap=0,overflow-y=scroll}]] {x=fill,scale-y=2}] {scale=1,gap=0,overflow-y=scroll}`)
	s.Viewports["page"] = runtime.ViewportState{Y: 20}
	g := measuredPanes(t, s, Size{200, 100}, defaultPaneMetrics())
	tab := g.Tabs["page/t"]
	if tab.Header != (Rect{0, -20, 200, 24}) || tab.HeaderClip != (Rect{0, 0, 200, 4}) || tab.Body != (Rect{0, 4, 200, 176}) {
		t.Fatalf("header transformed/clipped incorrectly %+v", tab)
	}
	if g.Root.Hit(10, 2) != nil || g.Root.Hit(10, 5).Path != "page/t/one/scroller/a" {
		t.Fatal("pane clip hit mismatch")
	}
	next, err := g.EnsureVisible("page/t/one/scroller/a", Rect{0, 304, 10, 20})
	if err != nil {
		t.Fatal(err)
	}
	if next["page/t/one/scroller"].Y != 144 || next["page"].Y != 100 {
		t.Fatalf("nested reveal %+v", next)
	}
	if s.Viewports["page"].Y != 20 {
		t.Fatal("reveal mutated state")
	}
}
