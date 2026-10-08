package runtime

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/Hans-Einar/SDP/SDUI/go/parser"
)

func TestExtendedInputInitialContentAdmissionKeepsSourceDiagnostic(t *testing.T) {
	for _, value := range []string{"a\nb", "a\rb", string([]byte{0xff}), strings.Repeat("x", 32769)} {
		root := paneRoot(t, `sdui 0.3;Main=[edit=input("Edit",multiline=false,value="saved")];`)
		var expected parser.Span
		root.Walk(func(n *parser.Instance) {
			if n.Widget == "input" {
				literal := n.Arguments["value"].(parser.Literal)
				expected = literal.Span
				literal.Value = value
				n.Arguments["value"] = literal
			}
		})
		_, err := New("admission", root)
		var diagnostic *parser.Diagnostic
		if !errors.As(err, &diagnostic) || diagnostic.Code != "input-initial" || diagnostic.Span != expected || !strings.Contains(diagnostic.Message, "Main/edit") {
			t.Fatal(err)
		}
	}
}

func TestExtendedInputInactiveDraftRetentionAndDialogValidation(t *testing.T) {
	const source = `sdui 0.3;Main=[tabs=tabs("Tabs")[one=page("One")[edit=input("Edit",required=true,value="saved")];two=page("Two")[other=input("Other")]]];`
	s := textSession(t, source)
	h := control(t, s, "Main/tabs/one/edit")
	requireOK(t, s.Draft(h, " "))
	requireOK(t, s.SelectPage(control(t, s, "Main/tabs"), "two"))
	before := s.Snapshot()
	code(t, s.Draft(h, "hidden"), "field-read-only")
	code(t, s.Revert(h), "field-read-only")
	if !reflect.DeepEqual(before, s.Snapshot()) {
		t.Fatal("hidden draft mutation")
	}
	n, err := s.Successor(paneRoot(t, source), nil)
	requireOK(t, err)
	f := fieldOf(t, n, h.Path)
	if f.Proposed != Text(" ") || !f.Dirty || f.Validation.Code != "text-required" {
		t.Fatal("inactive invalid draft lost", f)
	}
	requireOK(t, n.SelectPage(control(t, n, "Main/tabs"), "one"))
	if fieldOf(t, n, h.Path).Proposed != Text(" ") {
		t.Fatal("revealing page overwrote draft")
	}
	dialog := textSession(t, `sdui 0.3;Main=[dialog=dialog("Form",modal=false)[tabs=tabs("Tabs")[one=page("One")[edit=input("Edit",required=true)];two=page("Two")[flag=checkbox("Flag")]]]];`)
	target, err := dialog.OpenSurface(control(t, dialog, "Main/dialog"), ContextTarget{})
	requireOK(t, err)
	requireOK(t, dialog.SelectPage(control(t, dialog, "Main/dialog/tabs"), "two"))
	_, err = dialog.CaptureDialog(target, Accept)
	code(t, err, "field-validation")
}

func TestExtendedInputCommitFailurePreservesDraftAndAcceptedTicket(t *testing.T) {
	s := textSession(t, `sdui 0.3;Main=[edit=input("Edit",multiline=true,value="saved")];`)
	h := control(t, s, "Main/edit")
	reject, prepares, publications := false, 0, 0
	var visible Snapshot
	requireOK(t, s.PreparePresentationWith(func(v Snapshot) (PresentationTicket, error) {
		prepares++
		if reject {
			return PresentationTicket{}, errors.New("resource failure")
		}
		return PresentationTicket{Publish: func() { publications++; visible = v }}, nil
	}))
	calls := 0
	requireOK(t, s.Bind(h, func(Event) ([]Update, error) {
		calls++
		reject = true
		return []Update{acceptance(fieldOf(t, s, h.Path))}, nil
	}))
	change, err := s.EditField(h, s.Revision, Text("pending"))
	requireOK(t, err)
	before := fieldOf(t, s, h.Path)
	last := visible
	oldPublications, oldPrepares := publications, prepares
	e, err := s.CaptureCommit(change.Field.Target)
	requireOK(t, err)
	if s.Dispatch(e) == nil {
		t.Fatal("failed preparation accepted")
	}
	after := fieldOf(t, s, h.Path)
	if calls != 1 || publications != oldPublications || prepares != oldPrepares+1 || !reflect.DeepEqual(last, visible) || after.Accepted != before.Accepted || after.Proposed != before.Proposed || after.Target.ValueRevision != before.Target.ValueRevision || after.Target.DraftRevision != before.Target.DraftRevision {
		t.Fatal(calls, prepares, publications, before, after)
	}
	if s.Dispatch(e) == nil || calls != 1 {
		t.Fatal("failed Commit replayed domain callback")
	}
}

func TestExtendedInputValidatorReentrantAcceptKeepsDomainBarrier(t *testing.T) {
	s := textSession(t, `sdui 0.3;Main=[dialog=dialog("Form",modal=false)[edit=input("Edit",multiline=true,value="saved");other=input("Other",multiline=true,value="before")]];`)
	target, err := s.OpenSurface(control(t, s, "Main/dialog"), ContextTarget{})
	requireOK(t, err)
	h := control(t, s, "Main/dialog/edit")
	other := control(t, s, "Main/dialog/other")
	armed := false
	requireOK(t, s.ValidateFieldWith(h, func(FieldState) FieldValidation {
		if armed {
			armed = false
			requireOK(t, s.Draft(other, "newer"))
		}
		return FieldValidation{}
	}))
	requireOK(t, s.Draft(h, "pending"))
	requireOK(t, s.BindInteraction(target.Handle, func(Event) (InteractionReply, error) {
		armed = true
		return InteractionReply{Accept: &AcceptDecision{Accepted: true}, Domain: DomainSucceeded}, nil
	}))
	e, err := s.CaptureDialog(target, Accept)
	requireOK(t, err)
	result, err := s.DispatchInteraction(e)
	code(t, err, "stale-result")
	if result.Domain != DomainSucceeded || fieldOf(t, s, other.Path).Proposed != Text("newer") || fieldOf(t, s, h.Path).Accepted != Text("saved") {
		t.Fatal("domain conflict lost newer draft", result)
	}
	_, err = s.CaptureDialog(target, Accept)
	code(t, err, "accept-blocked")
}
