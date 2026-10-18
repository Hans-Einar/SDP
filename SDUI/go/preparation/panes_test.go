package preparation_test

import (
	"github.com/Hans-Einar/SDP/SDUI/go/host/fynehost/admission"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/preparation"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"strings"
	"testing"
)

func TestPanePreparationRequiresActualInteractionBinding(t *testing.T) {
	d, err := parser.Parse(`sdui 0.3;ref:api "unused.sdl";page=[tabs=tabs("Tabs",callback=api.Page.@invoke)[one=page("One")[]]];`)
	if err != nil {
		t.Fatal(err)
	}
	for _, bind := range []bool{false, true} {
		calls := 0
		r := preparation.Request{Document: d, Entry: "page", SessionID: "pane-prep", SourceRevision: "source", Mode: preparation.Connected, Capabilities: admission.PaneCapabilities(), ValidatePresentation: func(ui.Snapshot) (ui.PresentationState, error) { return ui.PresentationState{}, nil }, Bind: func(s *ui.Session, _ *parser.Document) error {
			if !bind {
				return nil
			}
			for _, w := range s.CallbackOwners() {
				if err := s.BindInteraction(w.Handle, func(ui.Event) (ui.InteractionReply, error) { calls++; return ui.InteractionReply{}, nil }); err != nil {
					return err
				}
			}
			return nil
		}}
		c, err := preparation.Prepare(r)
		if bind && err != nil || !bind && err == nil {
			t.Fatal("binding readiness", bind, err)
		}
		if calls != 0 {
			t.Fatal("preparation executed callback")
		}
		if c != nil {
			c.Close()
		}
	}
}
func TestHiddenPaneIconCapabilityAndM2RemainRejected(t *testing.T) {
	d, roots, err := parser.Compile(`sdui 0.3;page=[tabs=tabs("Tabs")[one=page("One")[];two=page("Two",icon="missing")[] {visible=false}]];`)
	if err != nil {
		t.Fatal(err)
	}
	err = preparation.Check(d.Profile, roots["page"], admission.PaneCapabilities())
	if err == nil || !strings.Contains(err.Error(), "provider icon/1") || !strings.Contains(err.Error(), "/two") {
		t.Fatal("hidden icon silently accepted", err)
	}
	n := roots["page"].Rows[0][0]
	n.Widget = "dialog"
	if err = preparation.Check(d.Profile, roots["page"], admission.PaneCapabilities()); err == nil {
		t.Fatal("M2 admission widened")
	}
}
