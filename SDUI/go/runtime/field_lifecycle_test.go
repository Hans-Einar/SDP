package runtime

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestChoiceGenerationReplacementNoLabelIndexFallback(t *testing.T) {
	s := scalarSession(t)
	h := control(t, s, "Main/mode")
	old, err := s.Option(h, "b")
	requireOK(t, err)
	f := fieldOf(t, s, h.Path)
	requireOK(t, s.ReplaceChoices(h, f.Target.OptionGeneration, []ChoiceOption{{ID: "b", Label: "Same", Enabled: true}}))
	f = fieldOf(t, s, h.Path)
	if f.Accepted != Choice("a") || f.Proposed != Choice("") || f.Validation.Code == "" {
		t.Fatal(f)
	}
	_, err = s.ChooseOption(old)
	code(t, err, "stale-option")
	fresh, err := s.Option(h, "b")
	requireOK(t, err)
	change, err := s.ChooseOption(fresh)
	requireOK(t, err)
	requireOK(t, fieldCommit(t, s, change.Field))
	f = fieldOf(t, s, h.Path)
	if f.Accepted != Choice("b") || f.Validation.Code != "" {
		t.Fatal(f)
	}
	requireOK(t, s.ReplaceChoices(h, f.Target.OptionGeneration, []ChoiceOption{{ID: "b", Label: "Same", Enabled: true}, {ID: "a", Label: "Same", Enabled: true}}))
	code(t, s.ValidateOptionTarget(fresh), "stale-option")
}
func TestChoiceInventoryEmptyRequiredAndAtomicRejection(t *testing.T) {
	src := `sdui 0.3;Main=[mode=select("Mode",required=true)];`
	s, err := New("choices", paneRoot(t, src))
	requireOK(t, err)
	f := fieldOf(t, s, "Main/mode")
	if f.Target.OptionGeneration != 0 {
		t.Fatal(f)
	}
	before := s.Snapshot()
	code(t, s.BindChoices(nil), "choices-unsupplied")
	if !reflect.DeepEqual(before, s.Snapshot()) {
		t.Fatal("partial inventory")
	}
	requireOK(t, s.BindChoices(map[string][]ChoiceOption{"Main/mode": {}}))
	f = fieldOf(t, s, "Main/mode")
	if f.Target.OptionGeneration == 0 || f.Validation.Code != "choice-required" {
		t.Fatal(f)
	}
	_, err = s.CaptureCommit(f.Target)
	code(t, err, "field-validation")
	before = s.Snapshot()
	code(t, s.ReplaceChoices(f.Target.Handle, f.Target.OptionGeneration, []ChoiceOption{{ID: "same", Label: "A", Enabled: true}, {ID: "same", Label: "B", Enabled: true}}), "choice-id")
	if !reflect.DeepEqual(before, s.Snapshot()) {
		t.Fatal("partial invalid option set")
	}
}

func TestRequiredEmptyChoiceSuccessorKeepsEditableAbsence(t *testing.T) {
	const source = `sdui 0.3;Main=[mode=select("Mode",required=true)];`
	for _, options := range [][]ChoiceOption{nil, {{ID: "a", Label: "Available", Enabled: true}}} {
		s, err := New("required", paneRoot(t, source))
		requireOK(t, err)
		choices := map[string][]ChoiceOption{"Main/mode": options}
		requireOK(t, s.BindChoices(choices))
		before := s.Snapshot()
		old := fieldOf(t, s, "Main/mode")
		n, err := s.SuccessorWithChoices(paneRoot(t, source), nil, choices)
		requireOK(t, err)
		f := fieldOf(t, n, "Main/mode")
		if f.Accepted != Choice("") || f.Proposed != Choice("") || f.Dirty || f.Validation.Code != "choice-required" || f.Target.ValueRevision != old.Target.ValueRevision || f.Target.OptionGeneration <= old.Target.OptionGeneration {
			t.Fatal("reload accepted or substituted missing choice", f)
		}
		_, err = n.CaptureCommit(f.Target)
		code(t, err, "field-validation")
		if !reflect.DeepEqual(before, s.Snapshot()) {
			t.Fatal("successor changed predecessor")
		}
		requireOK(t, s.Reload(paneRoot(t, source)))
		f = fieldOf(t, s, "Main/mode")
		if f.Accepted != Choice("") || f.Validation.Code != "choice-required" || f.Target.ValueRevision != old.Target.ValueRevision {
			t.Fatal("compatibility reload changed baseline", f)
		}
	}
}

func TestRequiredEmptyChoiceClosedSuccessorDiscardsOnlyProposal(t *testing.T) {
	const source = `sdui 0.3;Main=[open=button("Open");dialog=dialog("Form",modal=false)[mode=select("Mode",required=true)]];`
	s, err := New("required-dialog", paneRoot(t, source))
	requireOK(t, err)
	choices := map[string][]ChoiceOption{"Main/dialog/mode": {{ID: "a", Label: "A", Enabled: true}}}
	requireOK(t, s.BindChoices(choices))
	_, err = s.OpenSurface(control(t, s, "Main/dialog"), ContextTarget{})
	requireOK(t, err)
	old := fieldOf(t, s, "Main/dialog/mode")
	target, err := s.Option(old.Target.Handle, "a")
	requireOK(t, err)
	_, err = s.ChooseOption(target)
	requireOK(t, err)
	before := s.Snapshot()
	n, err := s.SuccessorWithChoices(paneRoot(t, source), nil, choices)
	requireOK(t, err)
	f := fieldOf(t, n, old.Target.Handle.Path)
	d, _ := n.SurfaceState(control(t, n, "Main/dialog"))
	if d.Open || f.Dirty || f.Accepted != Choice("") || f.Proposed != Choice("") || f.Validation.Code != "choice-required" || f.Target.ValueRevision != old.Target.ValueRevision {
		t.Fatal("closed form accepted uncommitted choice", d, f)
	}
	if !reflect.DeepEqual(before, s.Snapshot()) {
		t.Fatal("closed candidate reset live form")
	}
	reopened, err := n.OpenSurface(control(t, n, "Main/dialog"), ContextTarget{})
	requireOK(t, err)
	_, err = n.CaptureDialog(reopened, Accept)
	code(t, err, "field-validation")
	_, err = n.CaptureCommit(fieldOf(t, n, old.Target.Handle.Path).Target)
	code(t, err, "field-validation")
	newTarget, err := n.Option(fieldOf(t, n, old.Target.Handle.Path).Target.Handle, "a")
	requireOK(t, err)
	_, err = n.ChooseOption(newTarget)
	requireOK(t, err)
	if fieldOf(t, n, old.Target.Handle.Path).Validation.Code != "" {
		t.Fatal("missing choice stopped being editable")
	}
}

func TestRequiredEmptyExceptionDoesNotRelaxNewOrMissingChoices(t *testing.T) {
	for _, tc := range []struct {
		name, source, next string
		accepted           ItemID
		options            []ChoiceOption
	}{
		{"new-required", `sdui 0.3;Main=[mode=select("Mode")];`, `sdui 0.3;Main=[mode=select("Mode",required=true)];`, "", []ChoiceOption{{ID: "a", Label: "A", Enabled: true}}},
		{"removed", `sdui 0.3;Main=[mode=select("Mode",required=true)];`, `sdui 0.3;Main=[mode=select("Mode",required=true)];`, "a", nil},
		{"disabled", `sdui 0.3;Main=[mode=select("Mode",required=true)];`, `sdui 0.3;Main=[mode=select("Mode",required=true)];`, "a", []ChoiceOption{{ID: "a", Label: "A", Enabled: false}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s, err := New("required-reject", paneRoot(t, tc.source))
			requireOK(t, err)
			requireOK(t, s.BindChoices(map[string][]ChoiceOption{"Main/mode": {{ID: "a", Label: "A", Enabled: true}}}))
			if tc.accepted != "" {
				target, err := s.Option(control(t, s, "Main/mode"), tc.accepted)
				requireOK(t, err)
				change, err := s.ChooseOption(target)
				requireOK(t, err)
				requireOK(t, fieldCommit(t, s, change.Field))
			}
			before := s.Snapshot()
			_, err = s.SuccessorWithChoices(paneRoot(t, tc.next), nil, map[string][]ChoiceOption{"Main/mode": tc.options})
			code(t, err, "choice-accepted")
			if !reflect.DeepEqual(before, s.Snapshot()) {
				t.Fatal("rejected constraints changed live state")
			}
		})
	}
}
func TestMixedDialogProposalsAcceptAndCancelChildCommit(t *testing.T) {
	s := scalarSession(t)
	target := openDialog(t, s, true)
	flag := control(t, s, "Main/dialog/dflag")
	change, err := s.EditField(flag, s.Revision, Bool(true))
	requireOK(t, err)
	requireOK(t, fieldCommit(t, s, change.Field))
	f := fieldOf(t, s, flag.Path)
	if !f.Dirty || f.Accepted != Bool(false) {
		t.Fatal("unbound auto commit persisted form", f)
	}
	texts, err := s.DialogFields(target.Handle)
	requireOK(t, err)
	controls, err := s.DialogControls(target.Handle)
	requireOK(t, err)
	if len(texts) != 1 || len(controls) != 4 {
		t.Fatal(len(texts), len(controls))
	}
	calls := 0
	requireOK(t, s.BindInteraction(target.Handle, func(e Event) (InteractionReply, error) {
		calls++
		if len(e.Dialog.Fields) != 4 || e.Dialog.Fields[0].Value != Bool(true) || e.Dialog.Fields[1].RawDraft == nil || e.Dialog.Fields[2].OptionTarget == nil {
			t.Fatal(e.Dialog)
		}
		*e.Dialog.Fields[1].RawDraft = "malicious copy"
		return InteractionReply{Accept: &AcceptDecision{Accepted: true}, Domain: DomainSucceeded}, nil
	}))
	e := dialogEvent(t, s, target, Accept)
	_, err = s.DispatchInteraction(e)
	requireOK(t, err)
	if calls != 1 || fieldOf(t, s, flag.Path).Accepted != Bool(true) {
		t.Fatal(calls)
	}
	results := s.DrainDialogResults()
	if len(results) != 1 || len(results[0].Fields) != 4 || *results[0].Fields[1].RawDraft != "1" {
		t.Fatal(results)
	}
	target = openDialog(t, s, true)
	requireOK(t, s.Bind(flag, func(e Event) ([]Update, error) { return []Update{acceptance(fieldOf(t, s, flag.Path))}, nil }))
	change, err = s.EditField(flag, s.Revision, Bool(false))
	requireOK(t, err)
	requireOK(t, fieldCommit(t, s, change.Field))
	_, err = s.EditField(flag, s.Revision, Bool(true))
	requireOK(t, err)
	requireOK(t, s.CloseSurface(target, "cancel"))
	f = fieldOf(t, s, flag.Path)
	if f.Accepted != Bool(false) || f.Proposed != Bool(false) || f.Dirty {
		t.Fatal("Cancel rolled back child commit", f)
	}
}
func TestMixedAcceptInvalidAndPostDomainOptionsConflict(t *testing.T) {
	s := scalarSession(t)
	target := openDialog(t, s, true)
	number := control(t, s, "Main/dialog/dcount")
	_, err := s.EditField(number, s.Revision, Text("-"))
	requireOK(t, err)
	_, err = s.CaptureDialog(target, Accept)
	code(t, err, "field-validation")
	_, err = s.EditTick(number, s.Revision, 1)
	code(t, err, "field-validation")
	requireOK(t, s.RevertField(number))
	h := control(t, s, "Main/dialog/dmode")
	calls := 0
	requireOK(t, s.BindInteraction(target.Handle, func(Event) (InteractionReply, error) {
		calls++
		f := fieldOf(t, s, h.Path)
		requireOK(t, s.ReplaceChoices(h, f.Target.OptionGeneration, []ChoiceOption{{ID: "new", Label: "New", Enabled: true}}))
		return InteractionReply{Domain: DomainSucceeded, Accept: &AcceptDecision{Accepted: true}}, nil
	}))
	result, err := s.DispatchInteraction(dialogEvent(t, s, target, Accept))
	code(t, err, "stale-result")
	if result.Domain != DomainSucceeded || calls != 1 {
		t.Fatal(result, calls)
	}
	_, err = s.CaptureDialog(target, Accept)
	code(t, err, "accept-blocked")
	requireOK(t, s.CloseSurface(target, "cancel"))
	got := s.DrainDialogResults()
	if len(got) != 1 || got[0].Domain != DomainSucceeded {
		t.Fatal(got)
	}
}
func TestTypedSuccessorRetainsMainDraftResetsClosedDialogAndRejectsConstraints(t *testing.T) {
	s := scalarSession(t)
	target := openDialog(t, s, true)
	_, err := s.EditField(control(t, s, "Main/count"), s.Revision, Text("-"))
	requireOK(t, err)
	dh := control(t, s, "Main/dialog/dcount")
	requireOK(t, s.Bind(dh, func(Event) ([]Update, error) { return []Update{acceptance(fieldOf(t, s, dh.Path))}, nil }))
	change, err := s.EditField(dh, s.Revision, Text("2"))
	requireOK(t, err)
	requireOK(t, fieldCommit(t, s, change.Field))
	_, err = s.EditField(dh, s.Revision, Text("3"))
	requireOK(t, err)
	before := s.Snapshot()
	next, err := s.SuccessorWithChoices(paneRoot(t, scalarSource), nil, scalarChoices())
	requireOK(t, err)
	if !reflect.DeepEqual(before, s.Snapshot()) {
		t.Fatal("candidate mutated predecessor")
	}
	main := fieldOf(t, next, "Main/count")
	dialog := fieldOf(t, next, dh.Path)
	if *main.RawDraft != "-" || main.Validation.Code == "" || !main.Dirty || dialog.Accepted != Numeric(2) || dialog.Proposed != Numeric(2) || dialog.Dirty || *dialog.RawDraft != "2" {
		t.Fatal(main, dialog)
	}
	if d, _ := next.SurfaceState(control(t, next, "Main/dialog")); d.Open {
		t.Fatal(d)
	}
	oldOption, _ := s.Option(control(t, s, "Main/mode"), "a")
	code(t, next.ValidateOptionTarget(oldOption), "stale-option")
	bad := strings.Replace(scalarSource, "min=0,max=10,step=1,value=1", "min=0,max=1,step=1,value=1", 1)
	_, err = s.SuccessorWithChoices(paneRoot(t, bad), nil, scalarChoices())
	if err == nil || !reflect.DeepEqual(before, s.Snapshot()) {
		t.Fatal("invalid retained accepted value changed live", err)
	}
	requireOK(t, s.CloseSurface(target, "cancel"))
}
func TestTypedCommitPostHandlerGateFailureNoAcceptedUpdate(t *testing.T) {
	s := scalarSession(t)
	h := control(t, s, "Main/flag")
	reject := false
	requireOK(t, s.CheckPresentationWith(func(Snapshot) (PresentationState, error) {
		if reject {
			return PresentationState{}, errors.New("size")
		}
		return PresentationState{Viewports: map[string]ViewportState{}}, nil
	}))
	calls := 0
	requireOK(t, s.Bind(h, func(Event) ([]Update, error) {
		calls++
		reject = true
		return []Update{acceptance(fieldOf(t, s, h.Path))}, nil
	}))
	change, err := s.EditField(h, s.Revision, Bool(true))
	requireOK(t, err)
	if err = fieldCommit(t, s, change.Field); err == nil {
		t.Fatal("gate accepted")
	}
	f := fieldOf(t, s, h.Path)
	if calls != 1 || f.Accepted != Bool(false) || !f.Dirty {
		t.Fatal(calls, f)
	}
}

func TestTypedReloadRetainsCompatibleObserversValidatorsButSuccessorDoesNot(t *testing.T) {
	s := scalarSession(t)
	h := control(t, s, "Main/count")
	calls := 0
	requireOK(t, s.ObserveChanges(h, func(FieldChange) { calls++ }))
	requireOK(t, s.ValidateFieldWith(h, func(f FieldState) FieldValidation {
		if f.Proposed.Kind == Number && f.Proposed.Number > 1 {
			return FieldValidation{"limit", "Too large"}
		}
		return FieldValidation{}
	}))
	n, err := s.SuccessorWithChoices(paneRoot(t, scalarSource), nil, scalarChoices())
	requireOK(t, err)
	c, err := n.EditField(control(t, n, h.Path), n.Revision, Text("2"))
	requireOK(t, err)
	if c.Field.Validation.Code != "" || calls != 0 {
		t.Fatal("successor retained closures")
	}
	requireOK(t, s.Reload(paneRoot(t, scalarSource)))
	c, err = s.EditField(control(t, s, h.Path), s.Revision, Text("2"))
	requireOK(t, err)
	if c.Field.Validation.Code != "limit" || calls != 1 {
		t.Fatal("Reload lost compatible closures", c, calls)
	}
}

func TestTypedMixedAcceptBudgetAndInactiveDraftRetention(t *testing.T) {
	for _, count := range []int{256, 257} {
		t.Run(fmt.Sprint(count), func(t *testing.T) {
			var src strings.Builder
			src.WriteString(`sdui 0.3;Main=[open=button("Open");dialog=dialog("D",modal=false)[`)
			for i := 0; i < count; i++ {
				if i > 0 {
					src.WriteString(";")
				}
				fmt.Fprintf(&src, "f%d=checkbox(\"Flag\")", i)
			}
			src.WriteString("]];")
			s, err := New("budget", paneRoot(t, src.String()))
			requireOK(t, err)
			target := openDialog(t, s, true)
			e, err := s.CaptureDialog(target, Accept)
			if count > 256 {
				code(t, err, "update-limit")
				return
			}
			requireOK(t, err)
			_, err = s.DispatchInteraction(e)
			requireOK(t, err)
			got := s.DrainDialogResults()
			if len(got) != 1 || len(got[0].Fields) != 256 {
				t.Fatal(len(got))
			}
		})
	}
	src := `sdui 0.3;Main=[tabs=tabs("Tabs")[one=page("One")[num=number("N",min=0,max=10,step=1,value=1)];two=page("Two")[flag=checkbox("Flag")]]];`
	s, err := New("inactive", paneRoot(t, src))
	requireOK(t, err)
	_, err = s.EditField(control(t, s, "Main/tabs/one/num"), s.Revision, Text("-"))
	requireOK(t, err)
	requireOK(t, s.SelectPage(control(t, s, "Main/tabs"), "two"))
	n, err := s.SuccessorWithChoices(paneRoot(t, src), nil, nil)
	requireOK(t, err)
	f := fieldOf(t, n, "Main/tabs/one/num")
	if *f.RawDraft != "-" || !f.Dirty || f.Validation.Code == "" {
		t.Fatal(f)
	}
	requireOK(t, n.SelectPage(control(t, n, "Main/tabs"), "one"))
	if *fieldOf(t, n, "Main/tabs/one/num").RawDraft != "-" {
		t.Fatal("hidden draft lost")
	}
}
