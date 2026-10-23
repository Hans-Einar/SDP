package layout

import (
	"reflect"
	"testing"

	"github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

func TestCanvasGateRejectsInvalidResizeBeforeFinalTicketPublication(t *testing.T) {
	root := compileProfile(t, "0.3", `[opener=button("Open");d=dialog("Dialog",modal=false)[s=split(axis="horizontal")[a=button("A");b=button("B")] {scale=1}] {scale=1,gap=0}] {scale=1,gap=0}`)
	s, err := runtime.New("canvas-gate", root)
	if err != nil {
		t.Fatal(err)
	}
	e := &Engine{Measure: canvasMeasure()}
	mainSize := Size{400, 300}
	sizes := map[string]Size{"page/d": {210, 100}}
	if err = s.CheckPresentationWith(func(snapshot runtime.Snapshot) (runtime.PresentationState, error) {
		g, err := e.LayoutCanvases(snapshot, mainSize, sizes)
		if err != nil {
			return runtime.PresentationState{}, err
		}
		return g.PresentationState(), nil
	}); err != nil {
		t.Fatal(err)
	}
	var accepted *CanvasLayout
	prepared, published := 0, 0
	if err = s.PreparePresentationWith(func(snapshot runtime.Snapshot) (runtime.PresentationTicket, error) {
		prepared++
		candidate, err := e.LayoutCanvases(snapshot, mainSize, sizes)
		if err != nil {
			return runtime.PresentationTicket{}, err
		}
		return runtime.PresentationTicket{Publish: func() { accepted = candidate; published++ }}, nil
	}); err != nil {
		t.Fatal(err)
	}
	handle, ok := s.Surface("page/d")
	if !ok {
		t.Fatal("missing surface")
	}
	if _, err = s.OpenSurface(handle, runtime.ContextTarget{}); err != nil {
		t.Fatal(err)
	}
	if accepted.Surfaces["page/d"] == nil || len(accepted.PresentationState().Splits) != 1 {
		t.Fatal("final ticket omitted open-canvas bounds")
	}
	original := s.Snapshot()
	oldGeometry := accepted
	oldPrepared, oldPublished := prepared, published
	sizes["page/d"] = Size{20, 20}
	if err = s.SetViewports(nil); err == nil {
		t.Fatal("invalid resized surface published")
	}
	if !reflect.DeepEqual(original, s.Snapshot()) || accepted != oldGeometry || prepared != oldPrepared || published != oldPublished {
		t.Fatal("failed canvas reached preparation/publication or mutated live state")
	}
	sizes["page/d"] = Size{210, 100}
	mainSize = Size{800, 600}
	if err = s.SetViewports(nil); err != nil {
		t.Fatal(err)
	}
	if accepted.Surfaces["page/d"].Root.Rect != (Rect{0, 0, 210, 100}) {
		t.Fatal("parent resize changed accepted nonmodal reference")
	}
	if accepted.PresentationState().Splits["page/d/s"].Effective != s.Snapshot().Splits["page/d/s"].Proportion {
		t.Fatal("finalized state and surface geometry differ")
	}
}
