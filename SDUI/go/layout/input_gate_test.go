package layout

import (
	"reflect"
	"testing"

	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

func TestInputGateRejectsBeforeDraftObserverOrTicket(t *testing.T) {
	s := fieldSession(t, `[n=input("Notes",multiline=true,required=true) {scale=1}] {scale=1,gap=0}`)
	w, _ := s.Widget("page/n")
	m := newInputMetrics()
	m.change = func(_ *parser.Instance, f runtime.FieldState, _ Size, out *FieldMetrics) {
		if f.Proposed.Text == "reject" {
			out.Minimum.H = 1000
		}
	}
	engine := &Engine{Measure: m}
	size := Size{160, 100}
	if err := s.CheckPresentationWith(func(snap runtime.Snapshot) (runtime.PresentationState, error) {
		g, err := engine.LayoutSnapshot(snap, size)
		if err != nil {
			return runtime.PresentationState{}, err
		}
		return g.PresentationState(), nil
	}); err != nil {
		t.Fatal(err)
	}
	prepared, published, changed := 0, 0, 0
	var live *SnapshotLayout
	if err := s.PreparePresentationWith(func(snap runtime.Snapshot) (runtime.PresentationTicket, error) {
		prepared++
		g, err := engine.LayoutSnapshot(snap, size)
		if err != nil {
			return runtime.PresentationTicket{}, err
		}
		return runtime.PresentationTicket{Publish: func() { published++; live = g }}, nil
	}); err != nil {
		t.Fatal(err)
	}
	if err := s.ObserveChanges(w.Handle, func(runtime.FieldChange) { changed++ }); err != nil {
		t.Fatal(err)
	}
	before := s.Snapshot()
	old := live
	p0, q0 := prepared, published
	if _, err := s.EditField(w.Handle, s.Revision, runtime.Text("reject")); err == nil {
		t.Fatal("bad geometry accepted")
	}
	if !reflect.DeepEqual(before, s.Snapshot()) || old != live || prepared != p0 || published != q0 || changed != 0 {
		t.Fatal("failed input gate mutated runtime/ticket")
	}
	initial := live.Fields["page/n"]
	for _, value := range []string{"accepted proposal 文\nsecond line", " "} {
		if _, err := s.EditField(w.Handle, s.Revision, runtime.Text(value)); err != nil {
			t.Fatal(err)
		}
		f, _ := s.Field(w.Handle)
		if f.Accepted != runtime.Text("") || f.Proposed != runtime.Text(value) || live.Fields["page/n"] != initial {
			t.Fatal("text proposal changed fixed geometry/value", f)
		}
	}
	f, _ := s.Field(w.Handle)
	if f.Validation.Code == "" || changed != 2 || prepared != p0+2 || published != q0+2 {
		t.Fatal("invalid editable draft not published once", f, changed, prepared, published)
	}
}

func TestInputHiddenPageOffsetsAndSplitRestore(t *testing.T) {
	s := fieldSession(t, `[s=split(axis="horizontal")[left=tabs("Tabs")[one=page("One")[scroller=[n=input("Notes",multiline=true) {x=fill,scale-y=2}] {scale=1,gap=0,overflow-y=scroll}];two=page("Two")[]] {scale=1};right=input("Other",multiline=false) {scale=1}] {scale=1}] {scale=1,gap=0}`)
	engine := &Engine{Measure: newInputMetrics()}
	size := Size{250, 200}
	var live *SnapshotLayout
	if err := s.CheckPresentationWith(func(snap runtime.Snapshot) (runtime.PresentationState, error) {
		g, err := engine.LayoutSnapshot(snap, size)
		if err != nil {
			return runtime.PresentationState{}, err
		}
		return g.PresentationState(), nil
	}); err != nil {
		t.Fatal(err)
	}
	if err := s.PreparePresentationWith(func(snap runtime.Snapshot) (runtime.PresentationTicket, error) {
		g, err := engine.LayoutSnapshot(snap, size)
		if err != nil {
			return runtime.PresentationTicket{}, err
		}
		return runtime.PresentationTicket{Publish: func() { live = g }}, nil
	}); err != nil {
		t.Fatal(err)
	}
	path := "page/s/left/one/scroller"
	entry := path + "/n"
	w, _ := s.Widget(entry)
	if _, err := s.EditField(w.Handle, s.Revision, runtime.Text("retained draft\n文")); err != nil {
		t.Fatal(err)
	}
	if err := s.SetViewports(map[string]runtime.ViewportState{path: {Y: 150}}); err != nil {
		t.Fatal(err)
	}
	tabs, _ := s.Pane("page/s/left")
	split, _ := s.Pane("page/s")
	if err := s.SelectPage(tabs, "two"); err != nil {
		t.Fatal(err)
	}
	if _, ok := live.Fields[entry]; ok {
		t.Fatal("inactive input measured")
	}
	if s.Snapshot().Viewports[path].Y != 150 {
		t.Fatal("inactive ancestor offset lost")
	}
	if err := s.CollapseSplit(split, runtime.SplitFirst); err != nil {
		t.Fatal(err)
	}
	size = Size{90, 120}
	if err := s.SetViewports(nil); err != nil {
		t.Fatal(err)
	}
	before := s.Snapshot()
	old := live
	if err := s.RestoreSplit(split); err == nil {
		t.Fatal("unfit pane restore admitted")
	}
	if !reflect.DeepEqual(before, s.Snapshot()) || old != live {
		t.Fatal("failed restore mutated state/geometry")
	}
	size = Size{250, 120}
	if err := s.RestoreSplit(split); err != nil {
		t.Fatal(err)
	}
	if err := s.SelectPage(tabs, "one"); err != nil {
		t.Fatal(err)
	}
	f, _ := s.Field(w.Handle)
	if f.Proposed != runtime.Text("retained draft\n文") || live.Fields[entry].Control.W != 108 || s.Snapshot().Viewports[path].Y != 96 {
		t.Fatal("reveal changed draft or failed finite clamp", f, live.Fields[entry], s.Snapshot().Viewports)
	}
}

func TestInputOpenCanvasResizeFailureIsAtomic(t *testing.T) {
	s := fieldSession(t, `[main=input("Main",multiline=false);hidden=input("Hidden",multiline=true) {visible=false};d=dialog("Editor",modal=false)[n=input("Notes",multiline=true,required=true) {scale=1};nested=dialog("Closed")[bad=input("Never",multiline=true)]] {scale=0.5,gap=0}] {scale=1,gap=0}`)
	m := newInputMetrics()
	engine := &Engine{Measure: m}
	mainSize := Size{320, 220}
	sizes := map[string]Size{"page/d": {160, 100}}
	var live *CanvasLayout
	measure := func(snap runtime.Snapshot) (*CanvasLayout, error) {
		return engine.LayoutCanvases(snap, mainSize, sizes)
	}
	if err := s.CheckPresentationWith(func(snap runtime.Snapshot) (runtime.PresentationState, error) {
		g, err := measure(snap)
		if err != nil {
			return runtime.PresentationState{}, err
		}
		return g.PresentationState(), nil
	}); err != nil {
		t.Fatal(err)
	}
	if err := s.PreparePresentationWith(func(snap runtime.Snapshot) (runtime.PresentationTicket, error) {
		g, err := measure(snap)
		if err != nil {
			return runtime.PresentationTicket{}, err
		}
		return runtime.PresentationTicket{Publish: func() { live = g }}, nil
	}); err != nil {
		t.Fatal(err)
	}
	if len(m.seen) != 1 || len(live.Main.Fields) != 1 {
		t.Fatal("hidden/closed field measured", m.seen)
	}
	h, _ := s.Surface("page/d")
	if _, err := s.OpenSurface(h, runtime.ContextTarget{}); err != nil {
		t.Fatal(err)
	}
	path := "page/d/n"
	w, _ := s.Widget(path)
	initial := live.Surfaces["page/d"].Fields[path]
	if initial.Control != (Rect{0, 10, 160, 70}) || initial.Feedback != (Rect{0, 80, 160, 20}) || len(live.Surfaces["page/d"].Fields) != 1 {
		t.Fatal("surface-local Entry geometry", initial)
	}
	if _, err := s.EditField(w.Handle, s.Revision, runtime.Text(" ")); err != nil {
		t.Fatal(err)
	}
	before := s.Snapshot()
	old := live
	sizes["page/d"] = Size{160, 89}
	if err := s.SetViewports(nil); err == nil {
		t.Fatal("exhausted Entry canvas accepted")
	}
	if !reflect.DeepEqual(before, s.Snapshot()) || old != live {
		t.Fatal("failed child resize mutated field/ticket")
	}
	sizes["page/d"] = Size{160, 100}
	mainSize = Size{640, 440}
	if err := s.SetViewports(nil); err != nil {
		t.Fatal(err)
	}
	if live.Surfaces["page/d"].Fields[path] != initial || len(live.PresentationState().Viewports) != 0 {
		t.Fatal("parent resize rescaled Entry or introduced text viewport")
	}
	if _, ok := m.seen["page/d/nested/bad"]; ok {
		t.Fatal("closed nested input measured")
	}
}
