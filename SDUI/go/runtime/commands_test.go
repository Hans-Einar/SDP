package runtime

import (
	"errors"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"reflect"
	"testing"
)

const commandSource = `sdui 0.3;
Main=[
 choice=command("Choice",toggle=true,key="Primary+K");
 a=command("A",toggle=true,exclusive="group",checked=true);
 b=command("B",toggle=true,exclusive="group");
 invoke=button(command="choice"); aa=button(command="a"); bb=button(command="b");
 legacy=button("Legacy",tooltip="Still legacy"); out=input("Output",value="old");
 menu=menu("Menu")[entry=item(command="choice");separator()];
 open=button("Open",effect="open",target="dialog");
 dialog=dialog("Dialog",modal=false)[field=input("Field",value="saved");ok=button("OK",effect="accept");cancel=button("Cancel",effect="cancel")]
];`

func commandSession(t *testing.T) *Session {
	t.Helper()
	s, err := New("commands", paneRoot(t, commandSource))
	requireOK(t, err)
	return s
}
func control(t *testing.T, s *Session, path string) Handle {
	t.Helper()
	w := s.currentControl(path)
	if w == nil {
		t.Fatalf("missing %s", path)
	}
	return w.Handle
}
func capture(t *testing.T, s *Session, path, via string) Event {
	t.Helper()
	e, err := s.CaptureCommand(control(t, s, path), via, nil)
	requireOK(t, err)
	return e
}
func invoke(t *testing.T, s *Session, path string) InteractionResult {
	t.Helper()
	r, err := s.DispatchInteraction(capture(t, s, path, "button"))
	requireOK(t, err)
	return r
}
func TestM2SharedCommandAndLegacyOptIn(t *testing.T) {
	s := commandSession(t)
	c := control(t, s, "Main/choice")
	calls := 0
	requireOK(t, s.BindInteraction(c, func(e Event) (InteractionReply, error) {
		calls++
		if e.Command == nil || e.Handle != c || e.Command.Checked == nil {
			t.Fatal(e)
		}
		return InteractionReply{Domain: DomainSucceeded}, nil
	}))
	invoke(t, s, "Main/invoke")
	state, _ := s.CommandState(c)
	if !state.Checked {
		t.Fatal("button did not toggle")
	}
	scope, err := s.OpenMenu(control(t, s, "Main/menu"), ContextTarget{})
	requireOK(t, err)
	e := capture(t, s, "Main/menu/entry", "menu")
	_, err = s.DispatchInteraction(e)
	requireOK(t, err)
	if s.menus[scope.Handle.Path].Open {
		t.Fatal("selection did not dismiss")
	}
	_, err = s.DispatchInteraction(capture(t, s, "Main/choice", "key"))
	requireOK(t, err)
	if calls != 3 {
		t.Fatal(calls)
	}
	legacy := control(t, s, "Main/legacy")
	requireOK(t, s.Bind(legacy, func(Event) ([]Update, error) { calls++; return nil, nil }))
	requireOK(t, s.Dispatch(Event{Handle: legacy, ModelRevision: s.Revision, Sequence: s.Sequence() + 1, Kind: Activate}))
	if calls != 4 {
		t.Fatal(calls)
	}
	code(t, s.Bind(control(t, s, "Main/invoke"), func(Event) ([]Update, error) { return nil, nil }), "binding")
	code(t, s.Dispatch(Event{Handle: control(t, s, "Main/invoke"), ModelRevision: s.Revision, Sequence: s.Sequence() + 1, Kind: Activate}), "event-type")
}
func TestM2ExclusiveAtomicAndSilentProgrammatic(t *testing.T) {
	s := commandSession(t)
	invoke(t, s, "Main/bb")
	a, _ := s.CommandState(control(t, s, "Main/a"))
	b, _ := s.CommandState(control(t, s, "Main/b"))
	if a.Checked || !b.Checked {
		t.Fatal(a, b)
	}
	before := s.Snapshot()
	invoke(t, s, "Main/bb")
	if !reflect.DeepEqual(before, s.Snapshot()) {
		t.Fatal("exclusive selected member should be noop")
	}
	err := s.Apply(s.Revision, s.BatchRevision+1, []Update{{Handle: a.Handle, Property: Checked, Value: Bool(true)}})
	code(t, err, "exclusive")
	requireOK(t, s.Apply(s.Revision, s.BatchRevision+1, []Update{{Handle: a.Handle, Property: Checked, Value: Bool(true)}, {Handle: b.Handle, Property: Checked, Value: Bool(false)}}))
}
func TestM2MenuStoredScopeAndNoRestamp(t *testing.T) {
	s := commandSession(t)
	h := control(t, s, "Main/menu")
	scope, err := s.OpenMenu(h, ContextTarget{})
	requireOK(t, err)
	e := capture(t, s, "Main/menu/entry", "menu")
	requireOK(t, s.Draft(control(t, s, "Main/out"), "new"))
	_, err = s.DispatchInteraction(e)
	code(t, err, "stale-event")
	if s.menus[h.Path].Open {
		t.Fatal("stale invocation retained menu")
	}
	replacement, err := s.OpenMenu(h, ContextTarget{})
	requireOK(t, err)
	requireOK(t, s.CloseMenu(scope))
	if !s.menus[h.Path].Open {
		t.Fatal("old cleanup dismissed replacement")
	}
	requireOK(t, s.Draft(control(t, s, "Main/out"), "newer"))
	requireOK(t, s.CloseMenu(replacement))
	if s.menus[h.Path].Open {
		t.Fatal("stale scope could not dismiss own opening")
	}
}
func TestM2MenuFailurePreservesDomainAndReplacement(t *testing.T) {
	s := commandSession(t)
	h := control(t, s, "Main/menu")
	var fresh MenuScope
	requireOK(t, s.BindInteraction(control(t, s, "Main/choice"), func(Event) (InteractionReply, error) {
		var err error
		fresh, err = s.OpenMenu(h, ContextTarget{})
		requireOK(t, err)
		return InteractionReply{Domain: DomainSucceeded}, errors.New("after execution")
	}))
	_, err := s.OpenMenu(h, ContextTarget{})
	requireOK(t, err)
	r, err := s.DispatchInteraction(capture(t, s, "Main/menu/entry", "menu"))
	if err == nil || r.Domain != DomainSucceeded || r.Status != "ui-conflict" {
		t.Fatal(r, err)
	}
	if !s.menus[h.Path].Open || s.menus[h.Path].Scope != fresh {
		t.Fatal("old cleanup revoked new opening")
	}
}
func TestM2InteractionPayloadAndReentrancy(t *testing.T) {
	s := commandSession(t)
	h := control(t, s, "Main/choice")
	calls := 0
	requireOK(t, s.BindInteraction(h, func(e Event) (InteractionReply, error) {
		calls++
		_, err := s.DispatchInteraction(e)
		code(t, err, "reentrant-interaction")
		requireOK(t, s.Draft(control(t, s, "Main/out"), "accepted legacy mutation"))
		return InteractionReply{Domain: DomainSucceeded}, nil
	}))
	e := capture(t, s, "Main/invoke", "button")
	extra := e
	extra.Dialog = &DialogRequest{}
	_, err := s.DispatchInteraction(extra)
	code(t, err, "event-type")
	r, err := s.DispatchInteraction(e)
	code(t, err, "stale-result")
	if calls != 1 || r.Domain != DomainSucceeded {
		t.Fatal(calls, r)
	}
	c, _ := s.CommandState(h)
	if c.Checked {
		t.Fatal("stale toggle published")
	}
	w, _ := s.Widget("Main/out")
	if w.Draft != "accepted legacy mutation" {
		t.Fatal(w)
	}
}

func TestM2SharedLabelUpdatesReachPresentations(t *testing.T) {
	s := commandSession(t)
	c := control(t, s, "Main/choice")
	requireOK(t, s.Apply(s.Revision, s.BatchRevision+1, []Update{{Handle: c, Property: Label, Value: Text("Updated")}}))
	v := s.Snapshot()
	if v.Presentations["Main/invoke"].Label != "Updated" || v.Presentations["Main/menu/entry"].Label != "Updated" {
		t.Fatal(v.Presentations)
	}
}

func TestM2SnapshotRetainsClosedArgumentSchema(t *testing.T) {
	s := commandSession(t)
	snapshot := s.Snapshot()
	_, err := parser.ResolveInteractions(snapshot.Root)
	requireOK(t, err)
	snapshot.Root.Walk(func(n *parser.Instance) {
		if n.Widget == "item" || n.Widget == "separator" {
			if _, ok := n.Arguments["label"]; ok {
				t.Fatal("runtime injected illegal menu label", n.Path)
			}
		}
	})
}
