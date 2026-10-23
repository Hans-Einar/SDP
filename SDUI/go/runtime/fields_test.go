package runtime

import (
	"errors"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"reflect"
	"testing"
)

const scalarSource = `sdui 0.3;Main=[flag=checkbox("Flag");slide=slider("Slide",min=0,max=1,step=0.1,value=0.3);count=number("Count",min=-10,max=10,step=0.1,value=0.3);mode=select("Mode",value="a");old=input("Old",value="text");open=button("Open");dialog=dialog("Form",modal=false)[dflag=checkbox("Flag");dcount=number("Count",min=0,max=10,step=1,value=1);dmode=select("Mode");text=input("Text",value="saved")]];`

func scalarSession(t *testing.T) *Session {
	t.Helper()
	s, err := New("scalar", paneRoot(t, scalarSource))
	requireOK(t, err)
	requireOK(t, s.BindChoices(scalarChoices()))
	return s
}
func scalarChoices() map[string][]ChoiceOption {
	return map[string][]ChoiceOption{"Main/mode": {{ID: "a", Label: "Same", Enabled: true}, {ID: "b", Label: "Same", Enabled: true}, {ID: "disabled", Label: "Disabled"}}, "Main/dialog/dmode": {{ID: "x", Label: "X", Enabled: true}}}
}
func fieldOf(t *testing.T, s *Session, path string) FieldState {
	t.Helper()
	f, ok := s.Field(control(t, s, path))
	if !ok {
		t.Fatal(path)
	}
	return f
}
func fieldCommit(t *testing.T, s *Session, f FieldState) error {
	t.Helper()
	e, err := s.CaptureCommit(f.Target)
	if err != nil {
		return err
	}
	return s.Dispatch(e)
}
func acceptance(f FieldState) Update {
	return Update{Handle: f.Target.Handle, Property: AcceptedValue, Value: f.Proposed, ExpectedValueRevision: f.Target.ValueRevision, ExpectedDraftRevision: f.Target.DraftRevision, ExpectedOptionGeneration: f.Target.OptionGeneration, AcceptDraft: true}
}
func TestTypedFieldsSeparateLegacyStorageAndProjection(t *testing.T) {
	s := scalarSession(t)
	for _, p := range []string{"Main/flag", "Main/slide", "Main/count", "Main/mode"} {
		w, _ := s.Widget(p)
		if w.Value != "" || w.Draft != "" {
			t.Fatal("scalar stringified", w)
		}
	}
	f := fieldOf(t, s, "Main/count")
	if f.Accepted != Numeric(.3) || *f.RawDraft != "0.3" {
		t.Fatal(f)
	}
	snap := s.Snapshot()
	snap.Fields["Main/mode"].Options[0].Label = "corrupt"
	*snap.Fields["Main/count"].RawDraft = "corrupt"
	snap.Fields["Main/count"].Numeric.Min = "corrupt"
	if fieldOf(t, s, "Main/mode").Options[0].Label != "Same" || *fieldOf(t, s, "Main/count").RawDraft != "0.3" {
		t.Fatal("snapshot aliases")
	}
	_, err := parser.ResolveInteractions(s.SnapshotRoot())
	requireOK(t, err)
}
func TestTypedChangeThenCommitAndSilentUpdates(t *testing.T) {
	s := scalarSession(t)
	h := control(t, s, "Main/flag")
	changes, calls := 0, 0
	requireOK(t, s.ObserveChanges(h, func(change FieldChange) {
		changes++
		if change.Field.Proposed != Bool(true) || change.Field.Accepted != Bool(false) {
			t.Fatal(change)
		}
	}))
	requireOK(t, s.Bind(h, func(e Event) ([]Update, error) {
		calls++
		if e.Control == nil || e.Value != Bool(true) {
			t.Fatal(e)
		}
		f := fieldOf(t, s, h.Path)
		return []Update{acceptance(f)}, nil
	}))
	c, err := s.EditField(h, s.Revision, Bool(true))
	requireOK(t, err)
	if changes != 1 || calls != 0 || !c.Field.Dirty {
		t.Fatal(c, changes, calls)
	}
	requireOK(t, fieldCommit(t, s, c.Field))
	f := fieldOf(t, s, h.Path)
	if f.Accepted != Bool(true) || f.Dirty || calls != 1 {
		t.Fatal(f, calls)
	}
	requireOK(t, s.Apply(s.Revision, s.BatchRevision+1, []Update{{Handle: h, Property: AcceptedValue, Value: Bool(false), ExpectedValueRevision: f.Target.ValueRevision, ExpectedDraftRevision: f.Target.DraftRevision}}))
	if changes != 1 || calls != 1 {
		t.Fatal("programmatic events")
	}
}
func TestNumericInvalidDraftExactCommitAndStep(t *testing.T) {
	s := scalarSession(t)
	h := control(t, s, "Main/count")
	for _, raw := range []string{"-", "1e", "0.30000000000000000000000000001", "11", "1e999999999999999999999"} {
		c, err := s.EditField(h, s.Revision, Text(raw))
		requireOK(t, err)
		if c.Field.Validation.Code == "" || c.Field.Proposed.Kind != "" || *c.Field.RawDraft != raw {
			t.Fatal(c)
		}
		_, err = s.CaptureCommit(c.Field.Target)
		code(t, err, "field-validation")
	}
	c, err := s.EditField(h, s.Revision, Text("0.4"))
	requireOK(t, err)
	requireOK(t, fieldCommit(t, s, c.Field))
	if f := fieldOf(t, s, h.Path); f.Accepted != Numeric(.4) || f.Dirty {
		t.Fatal(f)
	}
	c, err = s.EditTick(h, s.Revision, 105)
	requireOK(t, err)
	if *c.Field.RawDraft != "0.5" || c.Field.Proposed != Numeric(.5) {
		t.Fatal(c)
	}
	requireOK(t, fieldCommit(t, s, c.Field))
}
func TestObserverReentranceStalesOriginalAutoCommit(t *testing.T) {
	s := scalarSession(t)
	h := control(t, s, "Main/flag")
	requireOK(t, s.ObserveChanges(h, func(FieldChange) { requireOK(t, s.Focus(control(t, s, "Main/old"))) }))
	c, err := s.EditField(h, s.Revision, Bool(true))
	requireOK(t, err)
	_, err = s.CaptureCommit(c.Field.Target)
	code(t, err, "stale-field")
	f := fieldOf(t, s, h.Path)
	if f.Accepted != Bool(false) || !f.Dirty {
		t.Fatal(f)
	}
}
func TestTypedReadOnlyValidationAndMixedBatchAtomicity(t *testing.T) {
	s := scalarSession(t)
	f := fieldOf(t, s, "Main/count")
	h := f.Target.Handle
	requireOK(t, s.Apply(s.Revision, s.BatchRevision+1, []Update{{Handle: h, Property: ReadOnly, Value: Bool(true), ExpectedValueRevision: f.Target.ValueRevision, ExpectedDraftRevision: f.Target.DraftRevision}}))
	_, err := s.EditField(h, s.Revision, Text("1"))
	code(t, err, "field-read-only")
	requireOK(t, s.Focus(h))
	f = fieldOf(t, s, h.Path)
	requireOK(t, s.Apply(s.Revision, s.BatchRevision+1, []Update{{Handle: h, Property: AcceptedValue, Value: Numeric(1), ExpectedValueRevision: f.Target.ValueRevision, ExpectedDraftRevision: f.Target.DraftRevision}}))
	before := s.Snapshot()
	flag := fieldOf(t, s, "Main/flag")
	f = fieldOf(t, s, h.Path)
	err = s.Apply(s.Revision, s.BatchRevision+1, []Update{{Handle: flag.Target.Handle, Property: AcceptedValue, Value: Bool(true), ExpectedValueRevision: flag.Target.ValueRevision, ExpectedDraftRevision: flag.Target.DraftRevision}, {Handle: h, Property: AcceptedValue, Value: Numeric(1.01), ExpectedValueRevision: f.Target.ValueRevision, ExpectedDraftRevision: f.Target.DraftRevision}})
	if err == nil || !reflect.DeepEqual(before, s.Snapshot()) {
		t.Fatal("mixed batch partially published", err)
	}
}
func TestFieldGateAndValidatorReentrance(t *testing.T) {
	s := scalarSession(t)
	h := control(t, s, "Main/count")
	reject := false
	requireOK(t, s.CheckPresentationWith(func(Snapshot) (PresentationState, error) {
		if reject {
			return PresentationState{}, errors.New("geometry")
		}
		return PresentationState{Viewports: map[string]ViewportState{}}, nil
	}))
	before := s.Snapshot()
	reject = true
	if _, err := s.EditField(h, s.Revision, Text("bad")); err == nil || !reflect.DeepEqual(before, s.Snapshot()) {
		t.Fatal("failed geometry changed draft")
	}
	reject = false
	requireOK(t, s.ValidateFieldWith(h, func(f FieldState) FieldValidation {
		if f.Proposed.Number > 1 {
			return FieldValidation{"limit", "Too large"}
		}
		return FieldValidation{}
	}))
	c, err := s.EditField(h, s.Revision, Text("2"))
	requireOK(t, err)
	if c.Field.Validation.Code != "limit" {
		t.Fatal(c)
	}
	_, err = s.CaptureCommit(c.Field.Target)
	code(t, err, "field-validation")
}
func TestTypedCommitRequiresOwnExactAcceptance(t *testing.T) {
	s := scalarSession(t)
	h := control(t, s, "Main/flag")
	calls := 0
	requireOK(t, s.Bind(h, func(Event) ([]Update, error) { calls++; return nil, nil }))
	c, err := s.EditField(h, s.Revision, Bool(true))
	requireOK(t, err)
	code(t, fieldCommit(t, s, c.Field), "commit-not-accepted")
	if calls != 1 || fieldOf(t, s, h.Path).Accepted != Bool(false) {
		t.Fatal(calls)
	}
	requireOK(t, s.Bind(h, func(e Event) ([]Update, error) {
		f := fieldOf(t, s, h.Path)
		requireOK(t, s.RevertField(h))
		return []Update{acceptance(f)}, nil
	}))
	f := fieldOf(t, s, h.Path)
	code(t, fieldCommit(t, s, f), "stale-result")
}

func TestTypedUpdateMergesPropertiesAndRejectsStaleDraft(t *testing.T) {
	s := scalarSession(t)
	h := control(t, s, "Main/flag")
	change, err := s.EditField(h, s.Revision, Bool(true))
	requireOK(t, err)
	u := acceptance(change.Field)
	u.ExpectedDraftRevision--
	before := s.Snapshot()
	code(t, s.Apply(s.Revision, s.BatchRevision+1, []Update{u}), "field-conflict")
	if !reflect.DeepEqual(before, s.Snapshot()) {
		t.Fatal("stale write published")
	}
	u = acceptance(change.Field)
	requireOK(t, s.Apply(s.Revision, s.BatchRevision+1, []Update{{Handle: h, Property: Label, Value: Text("Renamed")}, u}))
	f := fieldOf(t, s, h.Path)
	w, _ := s.Widget(h.Path)
	if f.Accepted != Bool(true) || w.Label != "Renamed" {
		t.Fatal(f, w)
	}
	f = fieldOf(t, s, h.Path)
	u = Update{Handle: h, Property: AcceptedValue, Value: Bool(false), ExpectedValueRevision: f.Target.ValueRevision, ExpectedDraftRevision: f.Target.DraftRevision}
	requireOK(t, s.Apply(s.Revision, s.BatchRevision+1, []Update{u, {Handle: h, Property: Label, Value: Text("Again")}}))
	f = fieldOf(t, s, h.Path)
	w, _ = s.Widget(h.Path)
	if f.Accepted != Bool(false) || w.Label != "Again" {
		t.Fatal(f, w)
	}
}
func TestTypedValidatorCannotOverwriteReentrantAcceptedWork(t *testing.T) {
	s := scalarSession(t)
	h := control(t, s, "Main/count")
	mutate := false
	requireOK(t, s.ValidateFieldWith(h, func(FieldState) FieldValidation {
		if mutate {
			mutate = false
			requireOK(t, s.Draft(control(t, s, "Main/old"), "preserved reentrant"))
		}
		return FieldValidation{}
	}))
	mutate = true
	_, err := s.EditField(h, s.Revision, Text("2"))
	code(t, err, "stale-validation")
	w, _ := s.Widget("Main/old")
	if w.Draft != "preserved reentrant" || *fieldOf(t, s, h.Path).RawDraft != "0.3" {
		t.Fatal(w)
	}
	f := fieldOf(t, s, h.Path)
	mutate = true
	code(t, s.Apply(s.Revision, s.BatchRevision+1, []Update{{Handle: h, Property: AcceptedValue, Value: Numeric(2), ExpectedValueRevision: f.Target.ValueRevision, ExpectedDraftRevision: f.Target.DraftRevision}}), "stale-validation")
	if fieldOf(t, s, h.Path).Accepted != Numeric(.3) {
		t.Fatal("stale validator result published")
	}
}
func TestTypedMalformedValueAndMetadataRejectAtomically(t *testing.T) {
	s := scalarSession(t)
	h := control(t, s, "Main/flag")
	f := fieldOf(t, s, h.Path)
	for _, v := range []Value{{Kind: Boolean, Bool: true, Text: "notzero"}, {Kind: Boolean, Bool: true, Number: 1}, {Kind: Boolean, Bool: true, OptionID: "x"}, {Kind: Number, Number: 1}} {
		before := s.Snapshot()
		err := s.Apply(s.Revision, s.BatchRevision+1, []Update{{Handle: h, Property: AcceptedValue, Value: v, ExpectedValueRevision: f.Target.ValueRevision, ExpectedDraftRevision: f.Target.DraftRevision}})
		if err == nil || !reflect.DeepEqual(before, s.Snapshot()) {
			t.Fatal(v, err)
		}
	}
}

func TestTypedBatchChecksOneBaselineForEveryPropertyOrder(t *testing.T) {
	var first Snapshot
	for _, reverse := range []bool{false, true} {
		s := scalarSession(t)
		f := fieldOf(t, s, "Main/flag")
		value := Update{Handle: f.Target.Handle, Property: AcceptedValue, Value: Bool(true), ExpectedValueRevision: f.Target.ValueRevision, ExpectedDraftRevision: f.Target.DraftRevision}
		readonly := Update{Handle: f.Target.Handle, Property: ReadOnly, Value: Bool(true), ExpectedValueRevision: f.Target.ValueRevision, ExpectedDraftRevision: f.Target.DraftRevision}
		updates := []Update{value, readonly}
		if reverse {
			updates = []Update{readonly, value}
		}
		requireOK(t, s.Apply(s.Revision, s.BatchRevision+1, updates))
		f = fieldOf(t, s, "Main/flag")
		if f.Accepted != Bool(true) || !f.ReadOnly || f.Target.ValueRevision != 2 || f.Target.DraftRevision != 2 {
			t.Fatal(reverse, f)
		}
		if !reverse {
			first = s.Snapshot()
		} else if !reflect.DeepEqual(first, s.Snapshot()) {
			t.Fatal("batch order changed result")
		}
	}
}
