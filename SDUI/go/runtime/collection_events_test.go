package runtime

import (
	"errors"
	"testing"
)

func TestTypedCollectionEventsLocalNavigationAndSequence(t *testing.T) {
	s, h := collectionSession(t, "tree", CollectionData{[]CollectionItem{group("g", "", true), row("child", "g"), {ID: "sep", Kind: Separator, ChildrenLoaded: true}, row("row", "")}}, true, false)
	calls := 0
	requireOK(t, s.Bind(h, func(e Event) ([]Update, error) {
		calls++
		if e.Collection.ItemID != "row" {
			t.Error(e.Collection)
		}
		return nil, nil
	}))
	targetRow := target(t, s, h, "row")
	base := eventFor(s, targetRow, Activate)
	tests := map[string]Event{}
	bad := base
	bad.Collection = nil
	tests["missing payload"] = bad
	bad = base
	bad.DraftRevision = 1
	tests["draft payload"] = bad
	bad = base
	bad.Value = Text("row")
	tests["text payload"] = bad
	bad = base
	badTarget := targetRow
	badTarget.ModelRevision++
	bad.Collection = &badTarget
	tests["wrong envelope"] = bad
	bad = base
	badTarget2 := targetRow
	badTarget2.CollectionGeneration++
	bad.Collection = &badTarget2
	tests["stale generation"] = bad
	tests["group activate"] = eventFor(s, target(t, s, h, "g"), Activate)
	tests["separator select"] = eventFor(s, target(t, s, h, "sep"), Select)
	tests["hidden child"] = eventFor(s, target(t, s, h, "child"), Activate)
	tests["root select"] = eventFor(s, target(t, s, h, ""), Select)
	tests["leaf expand"] = eventFor(s, targetRow, Expand)
	tests["leaf collapse"] = eventFor(s, targetRow, Collapse)
	tests["loaded retry"] = eventFor(s, targetRow, Retry)
	for name, e := range tests {
		t.Run(name, func(t *testing.T) {
			before := s.Snapshot()
			if err := s.Dispatch(e); err == nil {
				t.Fatal("invalid event accepted")
			}
			unchanged(t, s, before, s.Revision, s.BatchRevision)
		})
	}
	if calls != 0 {
		t.Fatal("invalid event ran handler")
	}
	requireOK(t, s.Dispatch(eventFor(s, target(t, s, h, "g"), Expand)))
	requireOK(t, s.Dispatch(eventFor(s, target(t, s, h, "child"), Select)))
	requireOK(t, s.FocusItem(target(t, s, h, "child")))
	requireOK(t, s.Dispatch(eventFor(s, target(t, s, h, "g"), Collapse)))
	c := state(t, s, h)
	if c.Selected != "child" || c.Focused != "g" || calls != 0 || s.Sequence() != 3 {
		t.Fatal(c, calls, s.Sequence())
	}
	requireOK(t, s.FocusItem(target(t, s, h, "g")))
	if state(t, s, h).Selected != "child" {
		t.Fatal("group focus selected group")
	}
	ev := eventFor(s, targetRow, Activate)
	requireOK(t, s.Dispatch(ev))
	code(t, s.Dispatch(ev), "duplicate-event")
	if calls != 1 {
		t.Fatal(calls)
	}
	// Root/collection local navigation works without SDL handler.
	s2, h2 := collectionSession(t, "list", CollectionData{[]CollectionItem{row("x", "")}}, true, false)
	requireOK(t, s2.Dispatch(eventFor(s2, target(t, s2, h2, "x"), Select)))
	code(t, s2.Dispatch(eventFor(s2, target(t, s2, h2, "x"), Activate)), "unbound")
	if s2.Sequence() != 2 {
		t.Fatal("unbound event did not consume sequence")
	}
}

func TestCollectionResultGuardCoversReentrantMutations(t *testing.T) {
	for _, mode := range []string{"identical-data", "remove-recreate", "row-to-group", "collapse", "disable", "hide", "reload", "close", "mutate-event-pointer"} {
		t.Run(mode, func(t *testing.T) {
			s, h := collectionSession(t, "tree", CollectionData{[]CollectionItem{branch("p", "", true), row("x", "p")}}, true, true)
			requireOK(t, s.ExpandItem(target(t, s, h, "p")))
			preview, _ := s.Widget("page/edit")
			requireOK(t, s.Bind(h, func(e Event) ([]Update, error) {
				d := state(t, s, h).Data
				switch mode {
				case "identical-data":
					requireOK(t, s.ReplaceCollection(target(t, s, h, ""), d))
				case "remove-recreate":
					requireOK(t, s.ReplaceCollection(target(t, s, h, ""), CollectionData{}))
					requireOK(t, s.ReplaceCollection(target(t, s, h, ""), d))
				case "row-to-group":
					d.Items[1].Kind = Group
					requireOK(t, s.ReplaceCollection(target(t, s, h, ""), d))
				case "collapse":
					requireOK(t, s.CollapseItem(target(t, s, h, "p")))
				case "disable":
					requireOK(t, s.Apply(s.Revision, s.BatchRevision+1, []Update{{Handle: h, Property: Enabled, Value: Bool(false)}}))
				case "hide":
					requireOK(t, s.Apply(s.Revision, s.BatchRevision+1, []Update{{Handle: h, Property: Visible, Value: Bool(false)}}))
				case "reload":
					requireOK(t, s.Reload(collectionRoot(t, "tree")))
				case "close":
					s.Close()
				case "mutate-event-pointer":
					requireOK(t, s.ReplaceCollection(target(t, s, h, ""), d))
					*e.Collection = target(t, s, h, "x")
				}
				return []Update{{Handle: preview.Handle, Property: AcceptedValue, Value: Text("must not publish"), ExpectedValueRevision: preview.ValueRevision}}, nil
			}))
			code(t, s.Dispatch(eventFor(s, target(t, s, h, "x"), Activate)), "stale-result")
			if now, _ := s.Widget("page/edit"); now.Value != "initial" {
				t.Fatal("stale result published", now)
			}
		})
	}
}

func TestLocalEventGateFailureConsumesNeitherSequenceNorState(t *testing.T) {
	s, h := collectionSession(t, "list", CollectionData{[]CollectionItem{row("x", "")}}, true, false)
	requireOK(t, s.CheckStateWith(func(v Snapshot) (map[string]ViewportState, error) {
		for _, c := range v.Collections {
			if c.Selected != "" {
				return nil, errors.New("native preparation failed")
			}
		}
		return v.Viewports, nil
	}))
	before := s.Snapshot()
	if err := s.Dispatch(eventFor(s, target(t, s, h, "x"), Select)); err == nil {
		t.Fatal("local event ignored gate")
	}
	unchanged(t, s, before, s.Revision, s.BatchRevision)
}

func TestRetryEventsAreLocalAndOneRequest(t *testing.T) {
	s, h := collectionSession(t, "list", CollectionData{}, false, true)
	req, e := s.BeginLoad(target(t, s, h, ""))
	requireOK(t, e)
	requireOK(t, s.CompleteLoad(req, CollectionData{}, errors.New("failure")))
	calls := 0
	requireOK(t, s.Bind(h, func(Event) ([]Update, error) { calls++; return nil, nil }))
	ev := eventFor(s, target(t, s, h, ""), Retry)
	requireOK(t, s.Dispatch(ev))
	request := *state(t, s, h).Request
	requireOK(t, s.Dispatch(eventFor(s, target(t, s, h, ""), Retry)))
	if *state(t, s, h).Request != request || calls != 0 || s.Sequence() != 2 {
		t.Fatal("Retry executed domain or duplicated load")
	}
}
