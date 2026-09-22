package runtime

import (
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"testing"
)

func session(t *testing.T) *Session {
	t.Helper()
	_, roots, e := parser.Compile(`sdui 0.2; page=[<go=button("Go"),edit=input("Input",value="1")>] {scale=1};`)
	if e != nil {
		t.Fatal(e)
	}
	s, e := New("test", roots["page"])
	if e != nil {
		t.Fatal(e)
	}
	return s
}
func code(t *testing.T, err error, want string) {
	t.Helper()
	e, ok := err.(*Fault)
	if !ok || e.Code != want {
		t.Fatalf("got %v want %s", err, want)
	}
}
func TestEventsAndAtomicUpdates(t *testing.T) {
	s := session(t)
	goButton, _ := s.Widget("page/go")
	edit, _ := s.Widget("page/edit")
	event := Event{Handle: goButton.Handle, ModelRevision: s.Revision, Sequence: 1, Kind: Activate}
	code(t, s.Dispatch(event), "unbound")
	calls := 0
	s.Bind(goButton.Handle, func(Event) ([]Update, error) {
		calls++
		return []Update{{Handle: edit.Handle, Property: AcceptedValue, Value: Text("2"), ExpectedValueRevision: 1}}, nil
	})
	code(t, s.Dispatch(event), "duplicate-event")
	event.Sequence = 2
	if e := s.Dispatch(event); e != nil {
		t.Fatal(e)
	}
	current, _ := s.Widget("page/edit")
	if current.Value != "2" || calls != 1 {
		t.Fatal(current, calls)
	}
	before := s.BatchRevision
	err := s.Apply(s.Revision, before+1, []Update{{Handle: goButton.Handle, Property: Label, Value: Text("changed")}, {Handle: edit.Handle, Property: Enabled, Value: Text("bad")}})
	code(t, err, "property-type")
	after, _ := s.Widget("page/go")
	if after.Label != "Go" || s.BatchRevision != before {
		t.Fatal("partial batch publication")
	}
	event.ModelRevision = 0
	event.Sequence = 3
	code(t, s.Dispatch(event), "stale-event")
	s.Close()
	code(t, s.Dispatch(event), "closed")
	if calls != 1 {
		t.Fatal("stale callback executed")
	}
}
func TestDraftConflictAndCommit(t *testing.T) {
	s := session(t)
	edit, _ := s.Widget("page/edit")
	s.Draft(edit.Handle, "3")
	s.Focus(edit.Handle)
	code(t, s.Apply(1, 1, []Update{{Handle: edit.Handle, Property: AcceptedValue, Value: Text("server"), ExpectedValueRevision: 1}}), "draft-conflict")
	s.Bind(edit.Handle, func(e Event) ([]Update, error) {
		return []Update{{Handle: e.Handle, Property: AcceptedValue, Value: e.Value, ExpectedValueRevision: 1, AcceptDraft: true}}, nil
	})
	w, _ := s.Widget("page/edit")
	if e := s.Dispatch(Event{Handle: w.Handle, ModelRevision: 1, Sequence: 1, DraftRevision: w.DraftRevision, Kind: Commit, Value: Text(w.Draft)}); e != nil {
		t.Fatal(e)
	}
	w, _ = s.Widget("page/edit")
	if w.Value != "3" || w.Dirty {
		t.Fatal(w)
	}
	snapshot := s.SnapshotRoot()
	snapshot.Rows[0][0].Rows[0][0].Arguments["label"] = parser.Literal{Value: "mutated"}
	unchanged, _ := s.Widget("page/go")
	if unchanged.Label != "Go" {
		t.Fatal("snapshot mutated session")
	}
}
