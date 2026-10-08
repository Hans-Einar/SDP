package layout

import (
	"reflect"
	"testing"

	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

func TestPreviewNestedScrollFitAndContentClip(t *testing.T) {
	root := previewRoot(t, `[inner=[p=`+previewCall+` {scale=2}] {scale=2,gap=0,overflow-x=scroll,overflow-y=scroll}] {scale=1,gap=0,overflow-x=scroll,overflow-y=scroll}`)
	offsets := map[string]runtime.ViewportState{"page": {X: 20, Y: 30}, "page/inner": {X: 40, Y: 50}}
	snapshot := runtime.Snapshot{Root: root, Viewports: offsets}
	m := &previewMetrics{natural: Size{400, 200}, gutters: true}
	e := &Engine{Measure: m}
	g, err := e.LayoutSnapshot(snapshot, Size{200, 100})
	if err != nil {
		t.Fatal(err)
	}
	p := boxes(g.Root)["page/inner/p"]
	if p.Rect != (Rect{-60, -80, 740, 340}) || p.Clip != (Rect{0, 0, 190, 90}) {
		t.Fatalf("translation/content clip %+v", p)
	}
	image, err := FitPreview(p.Rect, Size{400, 200})
	if err != nil || image != (Rect{-30, -80, 680, 340}) || image.Intersect(p.Clip) != (Rect{0, 0, 190, 90}) {
		t.Fatal("fit/clipping", image, err)
	}
	if g.Viewports["page/inner"].Parent != "page" || g.Viewports["page"].Maximum != (runtime.ViewportState{X: 190, Y: 90}) || g.Viewports["page/inner"].Maximum != (runtime.ViewportState{X: 370, Y: 170}) {
		t.Fatal("preview changed outer extents/ancestry", g.Viewports)
	}
	if hit := g.Root.Hit(50, 50); hit == nil || hit.Path != p.Path {
		t.Fatal("outer box hit lost")
	}
	if g.Root.Hit(195, 50) != nil || g.Root.Hit(50, 95) != nil {
		t.Fatal("child entered ancestor gutter")
	}
	snapshot.Viewports = nil
	atOrigin, err := e.LayoutSnapshot(snapshot, Size{200, 100})
	if err != nil {
		t.Fatal(err)
	}
	origin, err := FitPreview(boxes(atOrigin.Root)[p.Path].Rect, Size{400, 200})
	if err != nil || origin.W != image.W || origin.H != image.H || origin.X-image.X != 60 || origin.Y-image.Y != 80 {
		t.Fatal("scroll refitted or double-translated image", origin, image, err)
	}
	if offsets["page/inner"] != (runtime.ViewportState{X: 40, Y: 50}) {
		t.Fatal("geometry mutated requested offsets")
	}
}

func TestPreviewInactiveDeclarationsAndHiddenSchema(t *testing.T) {
	root := previewRoot(t, `[t=tabs("Tabs")[one=page("One")[a=button("A")];two=page("Two")[p=`+previewCall+`]] {scale=1};d=dialog("Closed")[q=`+previewCall+`];hidden=`+previewCall+` {visible=false}] {scale=1,gap=0}`)
	s, err := runtime.New("hidden-preview", root)
	if err != nil {
		t.Fatal(err)
	}
	m := &previewMetrics{}
	g, err := (&Engine{Measure: m}).LayoutSnapshot(s.Snapshot(), Size{400, 200})
	if err != nil {
		t.Fatal(err)
	}
	if len(m.calls) != 0 || boxes(g.Root)["page/t/two/p"] != nil || boxes(g.Root)["page/d/q"] != nil || boxes(g.Root)["page/hidden"] != nil {
		t.Fatal("hidden preview measured or exposed")
	}
	// Layout checks schema, but preparing hidden resource outcomes is the
	// provider's responsibility. No hidden geometry or renderer is requested.
	root.Walk(func(n *parser.Instance) {
		if n.Path == "page/d/q" {
			delete(n.Arguments, "fallback")
		}
	})
	if _, err = (&Engine{Measure: m}).Layout(root, Size{400, 200}); err == nil {
		t.Fatal("hidden malformed policy escaped layout admission")
	}
}

func TestPreviewCanvasGateFailedResizeIsPure(t *testing.T) {
	root := previewRoot(t, `[p=`+previewCall+` {scale=1};d=dialog("Preview",modal=false)[q=`+previewCall+` {scale=1}] {scale=0.5,gap=0}] {scale=1,gap=0}`)
	s, err := runtime.New("preview-canvases", root)
	if err != nil {
		t.Fatal(err)
	}
	m := &previewMetrics{natural: Size{800, 400}, minimum: Size{80, 40}}
	e := &Engine{Measure: m}
	mainSize := Size{200, 100}
	sizes := map[string]Size{"page/d": {120, 80}}
	if err = s.CheckPresentationWith(func(snapshot runtime.Snapshot) (runtime.PresentationState, error) {
		g, err := e.LayoutCanvases(snapshot, mainSize, sizes)
		if err != nil {
			return runtime.PresentationState{}, err
		}
		return g.PresentationState(), nil
	}); err != nil {
		t.Fatal(err)
	}
	var live *CanvasLayout
	prepared, published := 0, 0
	if err = s.PreparePresentationWith(func(snapshot runtime.Snapshot) (runtime.PresentationTicket, error) {
		prepared++
		g, err := e.LayoutCanvases(snapshot, mainSize, sizes)
		if err != nil {
			return runtime.PresentationTicket{}, err
		}
		return runtime.PresentationTicket{Publish: func() { live = g; published++ }}, nil
	}); err != nil {
		t.Fatal(err)
	}
	h, _ := s.Surface("page/d")
	if _, err = s.OpenSurface(h, runtime.ContextTarget{}); err != nil {
		t.Fatal(err)
	}
	image, err := FitPreview(boxes(live.Surfaces["page/d"].Root)["page/d/q"].Rect, Size{400, 200})
	if err != nil || image != (Rect{0, 10, 120, 60}) {
		t.Fatal("nonmodal local fit", image, err)
	}
	before, old, p0, q0 := s.Snapshot(), live, prepared, published
	sizes["page/d"] = Size{50, 80}
	if err = s.SetViewports(nil); err == nil {
		t.Fatal("undersized text allocation admitted")
	}
	if !reflect.DeepEqual(before, s.Snapshot()) || live != old || prepared != p0 || published != q0 {
		t.Fatal("failed pure geometry changed state/ticket/publication")
	}
	sizes["page/d"] = Size{120, 80}
	mainSize = Size{400, 200}
	if err = s.SetViewports(nil); err != nil {
		t.Fatal(err)
	}
	if live.Surfaces["page/d"].Root.Rect != (Rect{0, 0, 120, 80}) || live.Main.Root.Rect != (Rect{0, 0, 400, 200}) {
		t.Fatal("parent resize changed accepted child canvas")
	}
}
