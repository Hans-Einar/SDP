package layout

import (
	"math"
	"reflect"
	"strings"
	"testing"

	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

// These synthetic adapter metrics are independent geometry oracles. They do
// not claim native row/font/scrollbar evidence, which belongs to the host lane.
type collectionOracle struct {
	fixtureMetrics
	content Size
	alter   func(*CollectionMetrics)
}

func (m collectionOracle) MeasureCollection(_ *parser.Instance, _ float64, outer Size) (CollectionMetrics, error) {
	metrics := CollectionMetrics{Minimum: Size{40, 40}, Content: m.content, Viewport: Rect{0, 20, outer.W - 10, outer.H - 20}}
	if m.alter != nil {
		m.alter(&metrics)
	}
	return metrics, nil
}
func TestCollectionOuterGeometryAndSingleOffsetAuthority(t *testing.T) {
	root := compileProfile(t, "0.3", `[rows=list("Entries") {scale=1,overflow-x=scroll,overflow-y=scroll}] {scale=1,gap=0}`)
	snapshot := runtime.Snapshot{Root: root, Viewports: map[string]runtime.ViewportState{"page/rows": {X: 900, Y: 900}}}
	g, err := (&Engine{Measure: collectionOracle{content: Size{600, 600}}}).LayoutSnapshot(snapshot, Size{200, 100})
	if err != nil {
		t.Fatal(err)
	}
	v := g.Viewports["page/rows"]
	if v.Rect != (Rect{0, 20, 190, 80}) || v.Clip != v.Rect || v.Offset != (runtime.ViewportState{X: 410, Y: 520}) {
		t.Fatalf("viewport/chrome clamp %+v", v)
	}
	if boxes(g.Root)["page/rows"].Rect != (Rect{0, 0, 200, 100}) {
		t.Fatal("collection outer box moved by its own offset")
	}
	// A native row's point receives only the offset returned above. Shared
	// layout has already moved the outer control by all ancestors, not by self.
	row := Rect{v.Rect.X + 420 - v.Offset.X, v.Rect.Y + 540 - v.Offset.Y, 30, 20}
	if row != (Rect{10, 40, 30, 20}) {
		t.Fatalf("row transform oracle %+v", row)
	}
	_, rest, err := g.RouteScroll(10, 10, 10, 10) // fixed title, outside inner viewport
	if err != nil || rest != (runtime.ViewportState{X: 10, Y: 10}) {
		t.Fatalf("title scrolled %+v %v", rest, err)
	}
	out, rest, err := g.RouteScroll(10, 40, -20, -30)
	if err != nil || rest != (runtime.ViewportState{}) || out["page/rows"] != (runtime.ViewportState{X: 390, Y: 490}) {
		t.Fatalf("row routing %+v %+v %v", out, rest, err)
	}
}

func TestCollectionNestedOuterExtentAndAncestorClip(t *testing.T) {
	root := compileProfile(t, "0.3", `[rows=tree("Tree") {scale-x=1,scale-y=2,overflow-x=scroll,overflow-y=scroll}] {scale=1,gap=0,overflow-y=scroll}`)
	g, err := (&Engine{Measure: collectionOracle{content: Size{600, 1000}}}).LayoutSnapshot(runtime.Snapshot{Root: root, Viewports: map[string]runtime.ViewportState{"page": {Y: 50}, "page/rows": {Y: 400}}}, Size{200, 100})
	if err != nil {
		t.Fatal(err)
	}
	if g.Viewports["page"].Content.H != 200 || g.Viewports["page/rows"].Content.H != 1000 {
		t.Fatalf("nested extent %+v", g.Viewports)
	}
	v := g.Viewports["page/rows"]
	if v.Parent != "page" || v.Rect != (Rect{0, -30, 190, 180}) || v.Clip != (Rect{0, 0, 190, 100}) || boxes(g.Root)["page/rows"].Rect.Y != -50 {
		t.Fatalf("nested collection transform %+v", v)
	}
}

func TestCollectionMetricsAndOverflowPolicies(t *testing.T) {
	for _, policy := range []string{"error", "clip", "scroll"} {
		t.Run(policy, func(t *testing.T) {
			root := compileProfile(t, "0.3", `[rows=list("Rows") {scale=1,overflow-x=`+policy+`,overflow-y=`+policy+`}] {scale=1,gap=0}`)
			_, err := (&Engine{Measure: collectionOracle{content: Size{600, 600}}}).LayoutSnapshot(runtime.Snapshot{Root: root}, Size{200, 100})
			if policy == "error" {
				if d, ok := err.(*parser.Diagnostic); !ok || d.Code != "overflow-x" {
					t.Fatalf("error policy: %v", err)
				}
			} else if err != nil {
				t.Fatal(err)
			}
		})
	}
	root := compileProfile(t, "0.3", `[rows=list("Rows") {scale=1,overflow-x=scroll,overflow-y=scroll}] {scale=1,gap=0}`)
	if _, err := (&Engine{}).Layout(root, Size{200, 100}); err == nil || !strings.Contains(err.Error(), "collection-measurement") {
		t.Fatalf("silent collection fallback: %v", err)
	}
	if _, err := (&Engine{Measure: collectionOracle{content: Size{600, 600}}}).Layout(root, Size{30, 100}); err == nil || !strings.Contains(err.Error(), "native-minimum") {
		t.Fatalf("minimum: %v", err)
	}
	for _, alter := range []func(*CollectionMetrics){
		func(m *CollectionMetrics) { m.Content.H = math.Inf(1) },
		func(m *CollectionMetrics) { m.Content.W = math.NaN() },
		func(m *CollectionMetrics) { m.Minimum.W = -1 },
		func(m *CollectionMetrics) { m.Viewport.X = -1 },
		func(m *CollectionMetrics) { m.Viewport.H = 10000 },
		func(m *CollectionMetrics) { m.Content.H = 1e7 + 1 },
	} {
		if _, err := (&Engine{Measure: collectionOracle{content: Size{60, 60}, alter: alter}}).Layout(root, Size{200, 100}); err == nil {
			t.Fatal("invalid native metrics accepted")
		}
	}
	// Auto/content sizing on a non-scroll axis includes chrome exactly once.
	root = compileProfile(t, "0.3", `[rows=list("Rows")] {scale=1,gap=0}`)
	g, err := (&Engine{Measure: collectionOracle{content: Size{60, 40}}}).LayoutSnapshot(runtime.Snapshot{Root: root}, Size{200, 100})
	if err != nil {
		t.Fatal(err)
	}
	if got := boxes(g.Root)["page/rows"].Rect; got.W != 70 || got.H != 60 {
		t.Fatalf("content size with chrome: %+v", got)
	}
}

func TestSnapshotPureStateGatePublishesClampsAtomically(t *testing.T) {
	root := compileProfile(t, "0.3", `[rows=list("Rows") {scale=1,overflow-x=scroll,overflow-y=scroll}] {scale=1,gap=0}`)
	session, err := runtime.New("layout-gate", root)
	if err != nil {
		t.Fatal(err)
	}
	items := make([]runtime.CollectionItem, 20)
	for i := range items {
		items[i] = runtime.CollectionItem{ID: runtime.ItemID(strings.Repeat("x", i+1)), Kind: runtime.Row, Label: "Entry", ChildrenLoaded: true}
	}
	if err = session.BindProviders(map[string]runtime.CollectionProvider{"page/rows": {ID: "rows", Epoch: 1, Initial: runtime.CollectionData{Items: items}, RootLoaded: true}}); err != nil {
		t.Fatal(err)
	}
	var prepared *SnapshotLayout
	reject := false
	err = session.CheckStateWith(func(s runtime.Snapshot) (map[string]runtime.ViewportState, error) {
		extent := Size{100, float64(len(s.Collections["page/rows"].Data.Items)) * 20}
		if reject {
			extent.H = math.Inf(1)
		}
		candidate, err := (&Engine{Measure: collectionOracle{content: extent}}).LayoutSnapshot(s, Size{200, 100})
		if err != nil {
			return nil, err
		}
		prepared = candidate // stage only; this test has no native publication
		return candidate.EffectiveOffsets(), nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if err = session.SetViewports(map[string]runtime.ViewportState{"page/rows": {Y: 999}}); err != nil {
		t.Fatal(err)
	}
	if got := session.Snapshot().Viewports["page/rows"].Y; got != 320 {
		t.Fatalf("published clamp %v", got)
	}
	before := session.Snapshot()
	rejectedGeometry := prepared
	reject = true
	if err = session.SetViewports(map[string]runtime.ViewportState{"page/rows": {Y: 100}}); err == nil {
		t.Fatal("failed gate accepted")
	}
	if !reflect.DeepEqual(before, session.Snapshot()) || prepared != rejectedGeometry {
		t.Fatal("failed geometry partially published")
	}
	reject = false
	target, err := session.Target(before.Collections["page/rows"].Handle, "")
	if err != nil {
		t.Fatal(err)
	}
	if err = session.ReplaceCollection(target, runtime.CollectionData{Items: items[:2]}); err != nil {
		t.Fatal(err)
	}
	if session.Snapshot().Viewports["page/rows"].Y != 0 || prepared.Viewports["page/rows"].Maximum.Y != 0 {
		t.Fatal("data shrink did not publish matching clamped geometry")
	}
}
