package runtime

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func textSession(t *testing.T, source string) *Session {
	t.Helper()
	s, err := New("text", paneRoot(t, source))
	requireOK(t, err)
	return s
}

func TestExtendedInputOptInProjectionAndPlaceholder(t *testing.T) {
	s := textSession(t, `sdui 0.3;Main=[old=input("Old");edit=input("Edit",multiline=false,value="saved");empty=input("Empty",placeholder="")];`)
	legacy := fieldOf(t, s, "Main/old")
	if legacy.Input != nil {
		t.Fatal("basic .3 opted in")
	}
	_, err := s.EditField(legacy.Target.Handle, s.Revision, Text("no"))
	code(t, err, "field")
	f := fieldOf(t, s, "Main/edit")
	if f.Input == nil || f.Input.Multiline || f.Input.Placeholder != "Edit" || f.Accepted != Text("saved") || f.Proposed != f.Accepted {
		t.Fatal(f)
	}
	// There is no second mutable text store inside scalar field metadata.
	if private := s.fields[f.Target.Handle.Path]; private.Accepted != (Value{}) || private.Proposed != (Value{}) || private.RawDraft != nil {
		t.Fatal("duplicated text authority", private)
	}
	f.Input.Placeholder = "mutated"
	*f.RawDraft = "mutated"
	snap := s.Snapshot()
	snap.Fields["Main/edit"].Input.Multiline = true
	if current := fieldOf(t, s, "Main/edit"); current.Input.Multiline || current.Input.Placeholder != "Edit" || *current.RawDraft != "saved" {
		t.Fatal("snapshot escaped", current)
	}
	requireOK(t, s.Apply(s.Revision, s.BatchRevision+1, []Update{{Handle: f.Target.Handle, Property: Label, Value: Text("Changed")}, {Handle: control(t, s, "Main/empty"), Property: Label, Value: Text("Changed too")}}))
	if fieldOf(t, s, "Main/edit").Input.Placeholder != "Changed" || fieldOf(t, s, "Main/empty").Input.Placeholder != "" {
		t.Fatal("placeholder presence lost")
	}
}

func TestExtendedInputChangeCommitExactTextAndLegacyEnvelope(t *testing.T) {
	s := textSession(t, `sdui 0.3;Main=[edit=input("Edit",multiline=true,value="saved");old=input("Old")];`)
	h := control(t, s, "Main/edit")
	changes, calls := 0, 0
	requireOK(t, s.ObserveChanges(h, func(c FieldChange) {
		changes++
		if c.Value != c.Field.Proposed || !c.Field.Dirty || c.Option != nil {
			t.Fatal(c)
		}
		c.Field.Input.Placeholder = "external mutation"
	}))
	requireOK(t, s.Bind(h, func(e Event) ([]Update, error) {
		calls++
		if e.Control == nil || e.Control.RawDraft != nil || e.Control.Option != nil || e.Value != Text("α\r\n🙂") {
			t.Fatal(e)
		}
		return []Update{acceptance(fieldOf(t, s, h.Path))}, nil
	}))
	requireOK(t, s.Draft(h, "α\r\n🙂"))
	f := fieldOf(t, s, h.Path)
	w, _ := s.Widget(h.Path)
	if changes != 1 || calls != 0 || w.Draft != "α\r\n🙂" || w.Value != "saved" || f.Input.Placeholder != "Edit" {
		t.Fatal(changes, calls, w, f)
	}
	legacy := Event{Handle: h, Kind: Commit, ModelRevision: s.Revision, Sequence: s.Sequence() + 1, DraftRevision: f.Target.DraftRevision, Value: f.Proposed}
	code(t, s.Dispatch(legacy), "event-type")
	e, err := s.CaptureCommit(f.Target)
	requireOK(t, err)
	requireOK(t, s.Dispatch(e))
	if calls != 1 || changes != 1 || fieldOf(t, s, h.Path).Accepted != Text("α\r\n🙂") {
		t.Fatal(calls, changes)
	}
	if err = s.Dispatch(e); err == nil || calls != 1 {
		t.Fatal("replay", err, calls)
	}
	old := fieldOf(t, s, "Main/old")
	_, err = s.CaptureCommit(old.Target)
	code(t, err, "field")
	requireOK(t, s.Draft(old.Target.Handle, "legacy"))
	requireOK(t, s.Bind(old.Target.Handle, func(Event) ([]Update, error) { return nil, nil }))
	old = fieldOf(t, s, "Main/old")
	requireOK(t, s.Dispatch(Event{Handle: old.Target.Handle, Kind: Commit, ModelRevision: s.Revision, Sequence: s.Sequence() + 1, DraftRevision: old.Target.DraftRevision, Value: Text("legacy")}))
}

func TestExtendedInputValidationReadOnlyAndAtomicApply(t *testing.T) {
	s := textSession(t, `sdui 0.3;Main=[edit=input("Edit",required=true);other=input("Other",readOnly=true,value="saved")];`)
	h := control(t, s, "Main/edit")
	f := fieldOf(t, s, h.Path)
	if f.Validation.Code != "text-required" || f.Dirty {
		t.Fatal("required-empty initial must be editable invalid", f)
	}
	for _, raw := range []string{" \t\u2003", "a\nb", "a\rb"} {
		change, err := s.EditField(h, s.Revision, Text(raw))
		requireOK(t, err)
		if change.Field.Validation.Code == "" || change.Value != Text(raw) || change.Field.Accepted != Text("") {
			t.Fatal(change)
		}
		_, err = s.CaptureCommit(change.Field.Target)
		code(t, err, "field-validation")
	}
	for _, raw := range []string{string([]byte{0xff}), strings.Repeat("x", 32769)} {
		before := s.Snapshot()
		_, err := s.EditField(h, s.Revision, Text(raw))
		if err == nil || !reflect.DeepEqual(before, s.Snapshot()) {
			t.Fatal("invalid admission mutated state", err)
		}
	}
	change, err := s.EditField(h, s.Revision, Text("  α\u2028β  "))
	requireOK(t, err)
	requireOK(t, fieldCommit(t, s, change.Field)) // only CR/LF are line breaks
	other := fieldOf(t, s, "Main/other")
	code(t, s.Draft(other.Target.Handle, "blocked"), "field-read-only")
	code(t, s.Revert(other.Target.Handle), "field-read-only")
	_, err = s.CaptureCommit(other.Target)
	code(t, err, "field-read-only")
	requireOK(t, s.Focus(other.Target.Handle))
	load := acceptance(fieldOf(t, s, "Main/other"))
	load.AcceptDraft, load.Value = false, Text("loaded")
	requireOK(t, s.Apply(s.Revision, s.BatchRevision+1, []Update{load}))
	before := s.Snapshot()
	f = fieldOf(t, s, h.Path)
	bad := acceptance(f)
	bad.AcceptDraft, bad.Value = false, Text("")
	code(t, s.Apply(s.Revision, s.BatchRevision+1, []Update{{Handle: other.Target.Handle, Property: Label, Value: Text("should not publish")}, bad}), "field-validation")
	if !reflect.DeepEqual(before, s.Snapshot()) {
		t.Fatal("mixed invalid batch partially applied")
	}
}

func TestExtendedInputObserverValidatorAndGateReentrance(t *testing.T) {
	s := textSession(t, `sdui 0.3;Main=[edit=input("Edit",multiline=true,value="saved")];`)
	h := control(t, s, "Main/edit")
	armed := false
	requireOK(t, s.ValidateFieldWith(h, func(f FieldState) FieldValidation {
		if armed {
			armed = false
			requireOK(t, s.Draft(h, "reentrant"))
		}
		if f.Proposed.Text == "invalid" {
			return FieldValidation{"app", "Invalid"}
		}
		return FieldValidation{}
	}))
	armed = true
	_, err := s.EditField(h, s.Revision, Text("outer"))
	code(t, err, "stale-validation")
	if fieldOf(t, s, h.Path).Proposed != Text("reentrant") {
		t.Fatal("validator overwrote newer work")
	}
	requireOK(t, s.ObserveChanges(h, func(FieldChange) { requireOK(t, s.Revert(h)) }))
	change, err := s.EditField(h, s.Revision, Text("observer"))
	requireOK(t, err)
	_, err = s.CaptureCommit(change.Field.Target)
	code(t, err, "stale-field")
	requireOK(t, s.ObserveChanges(h, nil))
	reject := false
	requireOK(t, s.CheckStateWith(func(v Snapshot) (map[string]ViewportState, error) {
		if reject {
			return nil, errors.New("geometry")
		}
		return v.Viewports, nil
	}))
	before := s.Snapshot()
	reject = true
	_, err = s.EditField(h, s.Revision, Text("rejected"))
	if err == nil || !reflect.DeepEqual(before, s.Snapshot()) {
		t.Fatal("gate rejection changed authoritative draft", err)
	}
}

func TestExtendedInputCallbackMustAcceptCaptureAndReentrantReplyRejects(t *testing.T) {
	for _, mode := range []string{"missing", "different", "reentrant", "success"} {
		t.Run(mode, func(t *testing.T) {
			s := textSession(t, `sdui 0.3;Main=[edit=input("Edit",multiline=true,value="saved")];`)
			h := control(t, s, "Main/edit")
			calls := 0
			requireOK(t, s.Bind(h, func(Event) ([]Update, error) {
				calls++
				u := acceptance(fieldOf(t, s, h.Path))
				switch mode {
				case "missing":
					return nil, nil
				case "different":
					u.Value = Text("normalized")
				case "reentrant":
					requireOK(t, s.Draft(h, "newer"))
				}
				return []Update{u}, nil
			}))
			change, err := s.EditField(h, s.Revision, Text("draft"))
			requireOK(t, err)
			e, err := s.CaptureCommit(change.Field.Target)
			requireOK(t, err)
			err = s.Dispatch(e)
			if (mode == "success") != (err == nil) || calls != 1 {
				t.Fatal(mode, err, calls)
			}
			f := fieldOf(t, s, h.Path)
			if mode != "success" && f.Accepted != Text("saved") || mode == "reentrant" && f.Proposed != Text("newer") {
				t.Fatal(f)
			}
			if s.Dispatch(e) == nil || calls != 1 {
				t.Fatal("failed callback was replayed")
			}
		})
	}
}
