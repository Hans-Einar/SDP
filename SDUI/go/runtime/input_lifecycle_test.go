package runtime

import (
	"reflect"
	"strings"
	"testing"
)

func TestExtendedInputSuccessorEmptyTighteningAndNewlines(t *testing.T) {
	for _, tc := range []struct {
		name, source, next, draft string
		wantError                 bool
	}{
		{"required-empty", `required=true`, `required=true`, "", false},
		{"new-required", `required=false`, `required=true`, "valid proposal", true},
		{"required-whitespace", `required=true,value=" "`, `required=true,value=" "`, "", true},
		{"invalid-draft", `required=true,value="saved"`, `required=true,value="saved"`, " ", false},
		{"multiline-draft-conversion", `multiline=true,value="saved"`, `multiline=false,value="saved"`, "a\r\nb", true},
		{"multiline-accepted-conversion", `multiline=true,value="a\nb"`, `multiline=false,value="saved"`, "", true},
		{"single-line-invalid-draft", `multiline=false,value="saved"`, `multiline=false,value="saved"`, "a\nb", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			source := `sdui 0.3;Main=[edit=input("Edit",` + tc.source + `)];`
			nextSource := `sdui 0.3;Main=[edit=input("Edit",` + tc.next + `)];`
			s := textSession(t, source)
			h := control(t, s, "Main/edit")
			if tc.draft != "" {
				requireOK(t, s.Draft(h, tc.draft))
			}
			before := s.Snapshot()
			n, err := s.Successor(paneRoot(t, nextSource), nil)
			if (err != nil) != tc.wantError || !reflect.DeepEqual(before, s.Snapshot()) {
				t.Fatal(err, tc.wantError, "live mutated")
			}
			if err != nil {
				return
			}
			old, got := fieldOf(t, s, h.Path), fieldOf(t, n, h.Path)
			if old.Accepted != got.Accepted || old.Proposed != got.Proposed || old.Target.ValueRevision != got.Target.ValueRevision || got.Validation.Code == "" {
				t.Fatal(old, got)
			}
			_, err = n.CaptureCommit(got.Target)
			code(t, err, "field-validation")
		})
	}
}

func TestExtendedInputClosedDialogResetBeforeConversionAndChildCommit(t *testing.T) {
	const source = `sdui 0.3;Main=[dialog=dialog("Form",modal=false)[edit=input("Edit",multiline=true,value="saved")]];`
	for _, childCommit := range []bool{false, true} {
		s := textSession(t, source)
		d, err := s.OpenSurface(control(t, s, "Main/dialog"), ContextTarget{})
		requireOK(t, err)
		h := control(t, s, "Main/dialog/edit")
		if childCommit {
			requireOK(t, s.Bind(h, func(Event) ([]Update, error) { return []Update{acceptance(fieldOf(t, s, h.Path))}, nil }))
		}
		change, err := s.EditField(h, s.Revision, Text("a\r\nb"))
		requireOK(t, err)
		requireOK(t, fieldCommit(t, s, change.Field))
		f := fieldOf(t, s, h.Path)
		if !childCommit && (!f.Dirty || f.Accepted != Text("saved")) {
			t.Fatal("unbound child accepted", f)
		}
		before := s.Snapshot()
		n, err := s.Successor(paneRoot(t, strings.Replace(source, "multiline=true", "multiline=false", 1)), nil)
		if childCommit {
			code(t, err, "input-accepted")
		} else {
			requireOK(t, err)
			f = fieldOf(t, n, h.Path)
			if f.Proposed != Text("saved") || f.Dirty || f.Validation.Code != "" {
				t.Fatal("closed form did not discard draft before line check", f)
			}
		}
		if !reflect.DeepEqual(before, s.Snapshot()) {
			t.Fatal("candidate modified original form")
		}
		requireOK(t, s.CloseSurface(d, "cancel"))
		if childCommit && fieldOf(t, s, h.Path).Accepted != Text("a\r\nb") {
			t.Fatal("Cancel undid explicit child acceptance")
		}
	}
}

func TestExtendedInputMixedAcceptValidatesFieldsAndCopiesCapture(t *testing.T) {
	s := textSession(t, `sdui 0.3;Main=[dialog=dialog("Form",modal=false)[edit=input("Edit",required=true,multiline=true);flag=checkbox("Flag")]];`)
	d, err := s.OpenSurface(control(t, s, "Main/dialog"), ContextTarget{})
	requireOK(t, err)
	_, err = s.CaptureDialog(d, Accept)
	code(t, err, "field-validation")
	h := control(t, s, "Main/dialog/edit")
	requireOK(t, s.Draft(h, "日本語\r\n🙂"))
	requireOK(t, s.ConfirmSurfacePublication(d))
	calls := 0
	requireOK(t, s.BindInteraction(d.Handle, func(e Event) (InteractionReply, error) {
		calls++
		if len(e.Dialog.Fields) != 2 || e.Dialog.Fields[0].Value != Text("日本語\r\n🙂") || e.Dialog.Fields[0].RawDraft == nil {
			t.Fatal(e)
		}
		*e.Dialog.Fields[0].RawDraft = "tampered"
		return InteractionReply{Accept: &AcceptDecision{Accepted: true}, Domain: DomainSucceeded}, nil
	}))
	e, err := s.CaptureDialog(d, Accept)
	requireOK(t, err)
	_, err = s.DispatchInteraction(e)
	requireOK(t, err)
	results := s.DrainDialogResults()
	if calls != 1 || fieldOf(t, s, h.Path).Accepted != Text("日本語\r\n🙂") || len(results) != 1 || *results[0].Fields[0].RawDraft != "日本語\r\n🙂" {
		t.Fatal(calls, results)
	}
}

func TestExtendedInputApplyBaselineAndSilentRevert(t *testing.T) {
	for _, reverse := range []bool{false, true} {
		s := textSession(t, `sdui 0.3;Main=[edit=input("Edit",required=false,value="saved")];`)
		h := control(t, s, "Main/edit")
		changes := 0
		requireOK(t, s.ObserveChanges(h, func(FieldChange) { changes++ }))
		requireOK(t, s.Draft(h, "pending"))
		requireOK(t, s.RevertField(h))
		f := fieldOf(t, s, h.Path)
		if changes != 1 || f.Dirty || f.Proposed != Text("saved") {
			t.Fatal(changes, f)
		}
		value := acceptance(f)
		value.AcceptDraft, value.Value = false, Text("loaded")
		readonly := Update{Handle: h, Property: ReadOnly, Value: Bool(true), ExpectedValueRevision: f.Target.ValueRevision, ExpectedDraftRevision: f.Target.DraftRevision}
		updates := []Update{value, readonly}
		if reverse {
			updates = []Update{readonly, value}
		}
		requireOK(t, s.Apply(s.Revision, s.BatchRevision+1, updates))
		current := fieldOf(t, s, h.Path)
		if changes != 1 || current.Accepted != Text("loaded") || !current.ReadOnly || current.Target.ValueRevision != f.Target.ValueRevision+1 {
			t.Fatal(changes, current)
		}
		before := s.Snapshot()
		code(t, s.Apply(s.Revision, s.BatchRevision+1, []Update{value}), "field-conflict")
		if !reflect.DeepEqual(before, s.Snapshot()) {
			t.Fatal("stale input update published")
		}
	}
}

func TestExtendedInputPolicyTransitionAndReloadClosures(t *testing.T) {
	const source = `sdui 0.3;Main=[edit=input("Edit",multiline=true,value="saved")];`
	s := textSession(t, source)
	h := control(t, s, "Main/edit")
	changes := 0
	requireOK(t, s.ObserveChanges(h, func(FieldChange) { changes++ }))
	requireOK(t, s.ValidateFieldWith(h, func(f FieldState) FieldValidation {
		if f.Proposed.Text == "invalid" {
			return FieldValidation{"application", "Invalid"}
		}
		return FieldValidation{}
	}))
	requireOK(t, s.Draft(h, "draft"))
	requireOK(t, s.Reload(paneRoot(t, source)))
	requireOK(t, s.Draft(h, "invalid"))
	if changes != 2 || fieldOf(t, s, h.Path).Validation.Code != "application" {
		t.Fatal("compatibility Reload dropped closures")
	}
	n, err := s.Successor(paneRoot(t, source), nil)
	requireOK(t, err)
	if fieldOf(t, n, h.Path).Validation.Code != "" || fieldOf(t, n, h.Path).Proposed != Text("invalid") {
		t.Fatal("detached successor copied closures or lost draft")
	}
	legacy := strings.Replace(source, ",multiline=true", "", 1)
	n, err = s.Successor(paneRoot(t, legacy), nil)
	requireOK(t, err)
	if fieldOf(t, n, h.Path).Input != nil || fieldOf(t, n, h.Path).Proposed != Text("invalid") {
		t.Fatal("policy removal lost legacy text")
	}
	_, err = n.CaptureCommit(fieldOf(t, n, h.Path).Target)
	code(t, err, "field")
	back, err := n.Successor(paneRoot(t, source), nil)
	requireOK(t, err)
	if fieldOf(t, back, h.Path).Input == nil || fieldOf(t, back, h.Path).Accepted != Text("saved") {
		t.Fatal("legacy promotion lost accepted value")
	}
}
