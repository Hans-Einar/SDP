package runtime

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func openDialog(t *testing.T, s *Session, published bool) SurfaceTarget {
	t.Helper()
	h := control(t, s, "Main/dialog")
	target, err := s.OpenSurfaceFrom(h, control(t, s, "Main/open"), ContextTarget{})
	requireOK(t, err)
	if published {
		requireOK(t, s.ConfirmSurfacePublication(target))
	}
	return target
}
func dialogEvent(t *testing.T, s *Session, target SurfaceTarget, kind EventKind) Event {
	t.Helper()
	e, err := s.CaptureDialog(target, kind)
	requireOK(t, err)
	return e
}
func TestM2LocalAcceptAndPublishedReceipt(t *testing.T) {
	for _, published := range []bool{false, true} {
		t.Run(fmt.Sprint(published), func(t *testing.T) {
			s := commandSession(t)
			target := openDialog(t, s, published)
			field := control(t, s, "Main/dialog/field")
			requireOK(t, s.Draft(field, "proposed"))
			e := dialogEvent(t, s, target, Accept)
			r, err := s.DispatchInteraction(e)
			requireOK(t, err)
			if r.Domain != DomainNotCalled || r.Status != "committed" {
				t.Fatal(r)
			}
			w, _ := s.Widget(field.Path)
			if w.Value != "proposed" || w.Dirty {
				t.Fatal(w)
			}
			if d, _ := s.SurfaceState(target.Handle); d.Open {
				t.Fatal(d)
			}
			results := s.DrainDialogResults()
			if !published {
				if len(results) != 0 {
					t.Fatal(results)
				}
				return
			}
			if len(results) != 1 || results[0].Kind != "accept" || results[0].AcceptSequence != e.Sequence || results[0].Sequence != e.Sequence || len(results[0].Fields) != 1 || results[0].Fields[0].Value.Text != "proposed" {
				t.Fatal(results)
			}
			if len(s.DrainDialogResults()) != 0 {
				t.Fatal("duplicate result")
			}
		})
	}
}
func TestM2AcceptFalseUnknownAndSucceededConflict(t *testing.T) {
	for _, mode := range []string{"false", "unknown", "succeeded-conflict"} {
		t.Run(mode, func(t *testing.T) {
			s := commandSession(t)
			calls := 0
			h := control(t, s, "Main/dialog")
			field := control(t, s, "Main/dialog/field")
			requireOK(t, s.BindInteraction(h, func(Event) (InteractionReply, error) {
				calls++
				switch mode {
				case "false":
					return InteractionReply{Domain: DomainRejected, Accept: &AcceptDecision{Message: "Fix it"}}, nil
				case "unknown":
					return InteractionReply{}, errors.New("connection lost")
				default:
					requireOK(t, s.Draft(field, "newer"))
					return InteractionReply{Domain: DomainSucceeded, Accept: &AcceptDecision{Accepted: true}}, nil
				}
			}))
			target := openDialog(t, s, true)
			requireOK(t, s.Draft(field, "attempt"))
			e := dialogEvent(t, s, target, Accept)
			r, err := s.DispatchInteraction(e)
			if mode == "false" {
				requireOK(t, err)
				if r.Status != "rejected" || r.Domain != DomainRejected {
					t.Fatal(r)
				}
			} else {
				if err == nil || r.Status != "ui-conflict" {
					t.Fatal(r, err)
				}
				_, err = s.CaptureDialog(target, Accept)
				code(t, err, "accept-blocked")
			}
			d, _ := s.SurfaceState(h)
			if !d.Open || d.AcceptSequence != e.Sequence || d.Domain != r.Domain {
				t.Fatal(d, r)
			}
			requireOK(t, s.Draft(field, "edit after attempt"))
			requireOK(t, s.Revert(field))
			requireOK(t, s.Draft(field, "discard"))
			_, err = s.DispatchInteraction(dialogEvent(t, s, target, Cancel))
			requireOK(t, err)
			results := s.DrainDialogResults()
			if calls != 1 || len(results) != 1 || results[0].Domain != r.Domain || results[0].AcceptSequence != e.Sequence || len(results[0].Fields) != 0 {
				t.Fatal(calls, results)
			}
			w, _ := s.Widget(field.Path)
			if w.Value != "saved" || w.Draft != "saved" {
				t.Fatal(w)
			}
		})
	}
}
func TestM2AcceptReentrantCloseFinalizesUndrainedReceipt(t *testing.T) {
	s := commandSession(t)
	h := control(t, s, "Main/dialog")
	var target SurfaceTarget
	calls := 0
	requireOK(t, s.BindInteraction(h, func(Event) (InteractionReply, error) {
		calls++
		requireOK(t, s.CloseSurface(target, "close"))
		if got := s.DrainDialogResults(); len(got) != 0 {
			t.Fatal("unfinished receipt released", got)
		}
		_, err := s.OpenSurface(h, ContextTarget{})
		code(t, err, "surface-result")
		return InteractionReply{Domain: DomainSucceeded, Accept: &AcceptDecision{Accepted: true}}, nil
	}))
	target = openDialog(t, s, true)
	r, err := s.DispatchInteraction(dialogEvent(t, s, target, Accept))
	code(t, err, "stale-result")
	if r.Domain != DomainSucceeded || calls != 1 {
		t.Fatal(r, calls)
	}
	receipts := s.DrainDialogResults()
	if len(receipts) != 1 || receipts[0].Domain != DomainSucceeded || receipts[0].AcceptSequence == 0 || receipts[0].Kind != "close" {
		t.Fatal(receipts)
	}
	next := openDialog(t, s, true)
	if next.OpenGeneration <= target.OpenGeneration {
		t.Fatal(next, target)
	}
	code(t, s.CloseSurface(target, "cancel"), "stale-surface")
}
func TestM2AcceptBudgetBeforeAndAfterExecution(t *testing.T) {
	for _, count := range []int{256, 257} {
		t.Run(fmt.Sprint(count), func(t *testing.T) {
			var source strings.Builder
			source.WriteString(`sdui 0.3; Main=[open=button("Open");out=input("Out");dialog=dialog("Dialog",modal=false)[`)
			for i := 0; i < count; i++ {
				if i > 0 {
					source.WriteString(";")
				}
				fmt.Fprintf(&source, "f%d=input(\"Field\",value=\"saved\")", i)
			}
			source.WriteString("]];")
			s, err := New("budget", paneRoot(t, source.String()))
			requireOK(t, err)
			calls := 0
			h := control(t, s, "Main/dialog")
			requireOK(t, s.BindInteraction(h, func(Event) (InteractionReply, error) {
				calls++
				return InteractionReply{Domain: DomainSucceeded, Accept: &AcceptDecision{Accepted: true}, Updates: []Update{{Handle: control(t, s, "Main/out"), Property: Label, Value: Text("changed")}}}, nil
			}))
			target := openDialog(t, s, true)
			e, err := s.CaptureDialog(target, Accept)
			if count == 257 {
				code(t, err, "update-limit")
				if calls != 0 {
					t.Fatal(calls)
				}
				return
			}
			requireOK(t, err)
			r, err := s.DispatchInteraction(e)
			code(t, err, "update-limit")
			if calls != 1 || r.Domain != DomainSucceeded || r.Status != "ui-conflict" {
				t.Fatal(r, calls)
			}
			w, _ := s.Widget("Main/out")
			if w.Label != "Out" {
				t.Fatal("partial update")
			}
			d, _ := s.SurfaceState(h)
			if !d.Open || !d.AcceptBlocked {
				t.Fatal(d)
			}
		})
	}
}
func TestM2FailedFinalGateNoPublicationAndCancelRetainsOutcome(t *testing.T) {
	s := commandSession(t)
	h := control(t, s, "Main/dialog")
	reject := false
	prepared, published := 0, 0
	requireOK(t, s.CheckPresentationWith(func(v Snapshot) (PresentationState, error) {
		if reject {
			return PresentationState{}, errors.New("geometry")
		}
		return PresentationState{Viewports: map[string]ViewportState{}}, nil
	}))
	requireOK(t, s.PreparePresentationWith(func(Snapshot) (PresentationTicket, error) {
		prepared++
		return PresentationTicket{Publish: func() { published++ }}, nil
	}))
	requireOK(t, s.BindInteraction(h, func(Event) (InteractionReply, error) {
		reject = true
		return InteractionReply{Domain: DomainSucceeded, Accept: &AcceptDecision{Accepted: true}}, nil
	}))
	target := openDialog(t, s, true)
	before := s.Snapshot()
	p := prepared
	r, err := s.DispatchInteraction(dialogEvent(t, s, target, Accept))
	if err == nil || r.Domain != DomainSucceeded || r.Status != "ui-conflict" {
		t.Fatal(r, err)
	}
	if prepared != p || published != p {
		t.Fatal(prepared, published, p)
	}
	after := s.Snapshot()
	if !reflect.DeepEqual(before.Root, after.Root) || !after.Surfaces["Main/dialog"].Open {
		t.Fatal("failed gate published UI")
	}
	reject = false
	requireOK(t, s.CloseSurface(target, "cancel"))
	results := s.DrainDialogResults()
	if len(results) != 1 || results[0].Domain != DomainSucceeded {
		t.Fatal(results)
	}
}
func TestM2AcceptHiddenPageCaptureNestedExclusionAndUTF8(t *testing.T) {
	source := `sdui 0.3;Main=[open=button("Open");dialog=dialog("D",modal=false)[pages=tabs("Tabs")[one=page("One")[a=input("A")];two=page("Two")[b=input("B")]];child=dialog("Child")[c=input("C")]]];`
	s, err := New("fields", paneRoot(t, source))
	requireOK(t, err)
	target := openDialog(t, s, true)
	fields, err := s.DialogFields(target.Handle)
	requireOK(t, err)
	if len(fields) != 2 {
		t.Fatal(fields)
	}
	_, err = s.DialogField(target.Handle, "child/c")
	if err == nil {
		t.Fatal("nested field accepted")
	}
	h, err := s.DialogField(target.Handle, "pages/two/b")
	requireOK(t, err)
	e := dialogEvent(t, s, target, Accept)
	if len(e.Dialog.Fields) != 2 || e.Dialog.Fields[1].Handle != h {
		t.Fatal(e)
	}
	active := control(t, s, "Main/dialog/pages/one/a")
	requireOK(t, s.Draft(active, string([]byte{0xff})))
	_, err = s.CaptureDialog(target, Accept)
	code(t, err, "dialog-field")
}
