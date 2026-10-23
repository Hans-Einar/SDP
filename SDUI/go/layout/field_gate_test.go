package layout

import (
	"reflect"
	"testing"

	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

func TestScalarGateRejectsBeforeTicketOrTypedDraftPublication(t *testing.T) {
	s := fieldSession(t, `[n=checkbox("N") {scale=1}] {scale=1,gap=0}`)
	w, _ := s.Widget("page/n")
	m := newScalarMetrics()
	reject := true
	m.change = func(_ *parser.Instance, f runtime.FieldState, _ Size, out *FieldMetrics) {
		if reject && f.Proposed.Bool {
			out.Minimum.H = 1000
		}
	}
	engine := &Engine{Measure: m}
	const width = 160
	size := Size{width, 80}
	if err := s.CheckPresentationWith(func(snapshot runtime.Snapshot) (runtime.PresentationState, error) {
		g, err := engine.LayoutSnapshot(snapshot, size)
		if err != nil {
			return runtime.PresentationState{}, err
		}
		return g.PresentationState(), nil
	}); err != nil {
		t.Fatal(err)
	}
	prepared, published, changed := 0, 0, 0
	var live *SnapshotLayout
	if err := s.PreparePresentationWith(func(snapshot runtime.Snapshot) (runtime.PresentationTicket, error) {
		prepared++
		g, err := engine.LayoutSnapshot(snapshot, size)
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
	if _, err := s.EditField(w.Handle, s.Revision, runtime.Bool(true)); err == nil {
		t.Fatal("invalid scalar geometry accepted")
	}
	if !reflect.DeepEqual(before, s.Snapshot()) || live != old || prepared != p0 || published != q0 || changed != 0 {
		t.Fatal("failed pure geometry reached observer/ticket/state")
	}
	reject = false
	if _, err := s.EditField(w.Handle, s.Revision, runtime.Bool(true)); err != nil {
		t.Fatal(err)
	}
	f, ok := s.Field(w.Handle)
	if !ok || !f.Proposed.Bool || f.Accepted.Bool || !f.Dirty || live == old || prepared != p0+1 || published != q0+1 || changed != 1 || len(live.Fields) != 1 {
		t.Fatal("accepted typed proposal and final geometry diverged", f, prepared, published, changed)
	}
}

func TestScalarInvalidNumericFeedbackUsesSameCanvasWithoutSourceMutation(t *testing.T) {
	s := fieldSession(t, `[open=button("Open");d=dialog("D",modal=false)[n=number("N",min=0,max=10,step=1,value=1.00) {x=fill}] {scale=0.5,gap=0}] {scale=1,gap=0}`)
	engine := &Engine{Measure: newScalarMetrics()}
	mainSize := Size{300, 200}
	sizes := map[string]Size{"page/d": {160, 80}}
	measure := func(snapshot runtime.Snapshot) (*CanvasLayout, error) {
		return engine.LayoutCanvases(snapshot, mainSize, sizes)
	}
	if err := s.CheckPresentationWith(func(snapshot runtime.Snapshot) (runtime.PresentationState, error) {
		g, err := measure(snapshot)
		if err != nil {
			return runtime.PresentationState{}, err
		}
		return g.PresentationState(), nil
	}); err != nil {
		t.Fatal(err)
	}
	var live *CanvasLayout
	if err := s.PreparePresentationWith(func(snapshot runtime.Snapshot) (runtime.PresentationTicket, error) {
		g, err := measure(snapshot)
		if err != nil {
			return runtime.PresentationTicket{}, err
		}
		return runtime.PresentationTicket{Publish: func() { live = g }}, nil
	}); err != nil {
		t.Fatal(err)
	}
	handle, _ := s.Surface("page/d")
	if _, err := s.OpenSurface(handle, runtime.ContextTarget{}); err != nil {
		t.Fatal(err)
	}
	w, _ := s.Widget("page/d/n")
	initial := live.Surfaces["page/d"].Fields["page/d/n"]
	if _, err := s.EditField(w.Handle, s.Revision, runtime.Text("-")); err != nil {
		t.Fatal(err)
	}
	field, _ := s.Field(w.Handle)
	if field.Accepted != runtime.Numeric(1) || field.Proposed.Kind != "" || field.Validation.Code == "" || *field.RawDraft != "-" || live.Surfaces["page/d"].Fields["page/d/n"] != initial {
		t.Fatal("feedback mutated value or fixed native geometry", field)
	}
	before := s.Snapshot()
	old := live
	sizes["page/d"] = Size{40, 40}
	if err := s.SetViewports(nil); err == nil {
		t.Fatal("small invalid child canvas accepted")
	}
	if !reflect.DeepEqual(before, s.Snapshot()) || live != old {
		t.Fatal("failed child resize discarded invalid raw draft or live canvas")
	}
	sizes["page/d"] = Size{160, 80}
	mainSize = Size{600, 400}
	if err := s.SetViewports(nil); err != nil {
		t.Fatal(err)
	}
	if live.Surfaces["page/d"].Fields["page/d/n"] != initial || len(live.Main.Fields) != 0 {
		t.Fatal("parent resize moved separate field canvas")
	}
}
