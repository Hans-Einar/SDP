package runtime

import (
	"errors"
	"math"
	"reflect"
	"testing"
)

func TestPaneInteractionPublishesOneAtomicFinalizedTicket(t *testing.T) {
	s := paneSession(t)
	h := paneHandle(t, s, "Main/pane/left")
	out, _ := s.Widget("Main/out")
	var published Snapshot
	prepared, promoted, discarded, calls := 0, 0, 0, 0
	requireOK(t, s.PreparePresentationWith(func(v Snapshot) (PresentationTicket, error) {
		prepared++
		return PresentationTicket{Publish: func() {
			promoted++
			published = v
			if !reflect.DeepEqual(v, s.Snapshot()) {
				t.Error("ticket did not match accepted state")
			}
		}, Discard: func() { discarded++ }}, nil
	}))
	requireOK(t, s.BindInteraction(h, func(e Event) (InteractionReply, error) {
		calls++
		if paneSelected(s) != "first" || e.Page.PageID != "second" {
			t.Error("prospective selection escaped probe")
		}
		e.Page.PageID = "handler modified its copy"
		return InteractionReply{Domain: DomainSucceeded, Updates: []Update{{Handle: out.Handle, Property: AcceptedValue, Value: Text("page accepted"), ExpectedValueRevision: out.ValueRevision}}}, nil
	}))
	beforePrepared, beforePromoted := prepared, promoted
	result, err := s.DispatchInteraction(paneEvent(t, s, "second"))
	requireOK(t, err)
	value, _ := s.Widget(out.Handle.Path)
	if result.Status != "committed" || result.Domain != DomainSucceeded || calls != 1 || prepared != beforePrepared+1 || promoted != beforePromoted+1 || discarded != 0 || value.Value != "page accepted" || published.Tabs[h.Path].Selected != "second" {
		t.Fatal(result, calls, prepared, promoted, discarded, value, published.Tabs)
	}
	sequence := s.Sequence()
	result, err = s.DispatchInteraction(paneEvent(t, s, "second"))
	requireOK(t, err)
	if calls != 1 || prepared != beforePrepared+1 || s.Sequence() != sequence {
		t.Fatal("same page was not a no-op")
	}
	requireOK(t, s.SelectPage(h, "first"))
	if calls != 1 {
		t.Fatal("programmatic selection called handler")
	}
}
func TestPaneFailedCallbackProbeNeverPromotesOrCancelsLiveRequest(t *testing.T) {
	s := paneSession(t)
	h := paneHandle(t, s, "Main/pane/left")
	items, _ := s.Widget("Main/pane/left/first/items")
	request, err := s.BeginLoad(target(t, s, items.Handle, ""))
	requireOK(t, err)
	prepared, promoted := 0, 0
	requireOK(t, s.PreparePresentationWith(func(Snapshot) (PresentationTicket, error) {
		prepared++
		return PresentationTicket{Publish: func() { promoted++ }}, nil
	}))
	requireOK(t, s.BindInteraction(h, func(Event) (InteractionReply, error) {
		return InteractionReply{Domain: DomainRejected}, errors.New("application rejected page")
	}))
	count := promoted
	stateBefore := state(t, s, items.Handle)
	sequence := s.Sequence()
	result, err := s.DispatchInteraction(paneEvent(t, s, "second"))
	if err == nil || result.Status != "rejected" || result.Domain != DomainRejected || paneSelected(s) != "first" || s.Sequence() != sequence+1 || promoted != count || prepared != count {
		t.Fatal(result, err, prepared, promoted)
	}
	if !reflect.DeepEqual(stateBefore, state(t, s, items.Handle)) {
		t.Fatal("probe canceled live provider")
	}
	requireOK(t, s.CompleteLoad(request, CollectionData{Items: []CollectionItem{row("ok", "")}}, nil))
}
func TestPaneAcceptedLegacyReentranceSurvivesOuterFailure(t *testing.T) {
	for _, callbackError := range []bool{false, true} {
		t.Run(map[bool]string{false: "stale-result", true: "callback-error"}[callbackError], func(t *testing.T) {
			s := paneSession(t)
			h := paneHandle(t, s, "Main/pane/left")
			edit, _ := s.Widget("Main/pane/left/first/edit")
			promotions := 0
			var visible Snapshot
			requireOK(t, s.PreparePresentationWith(func(v Snapshot) (PresentationTicket, error) {
				return PresentationTicket{Publish: func() { promotions++; visible = v }}, nil
			}))
			requireOK(t, s.BindInteraction(h, func(Event) (InteractionReply, error) {
				_, nestedErr := s.DispatchInteraction(paneEvent(t, s, "second"))
				code(t, nestedErr, "reentrant-interaction")
				requireOK(t, s.Draft(edit.Handle, "accepted legacy draft"))
				if callbackError {
					return InteractionReply{Domain: DomainSucceeded}, errors.New("later adapter error")
				}
				return InteractionReply{Domain: DomainSucceeded}, nil
			}))
			count := promotions
			result, err := s.DispatchInteraction(paneEvent(t, s, "second"))
			if err == nil || result.Status != "ui-conflict" || result.Domain != DomainSucceeded || paneSelected(s) != "first" || promotions != count+1 || visible.Tabs[h.Path].Selected != "first" {
				t.Fatal(result, err, promotions)
			}
			w, _ := s.Widget(edit.Handle.Path)
			if w.Draft != "accepted legacy draft" {
				t.Fatal("accepted reentrant state rolled back")
			}
		})
	}
}
func TestPaneFinalPreparationFailurePreservesDomainOutcomeAndState(t *testing.T) {
	s := paneSession(t)
	h := paneHandle(t, s, "Main/pane/left")
	reject := false
	promoted, discarded := 0, 0
	requireOK(t, s.PreparePresentationWith(func(Snapshot) (PresentationTicket, error) {
		ticket := PresentationTicket{Publish: func() { promoted++ }, Discard: func() { discarded++ }}
		if reject {
			return ticket, errors.New("native preparation rejected")
		}
		return ticket, nil
	}))
	requireOK(t, s.BindInteraction(h, func(Event) (InteractionReply, error) {
		reject = true
		return InteractionReply{Domain: DomainSucceeded}, nil
	}))
	count := promoted
	result, err := s.DispatchInteraction(paneEvent(t, s, "second"))
	if err == nil || result.Status != "ui-conflict" || result.Domain != DomainSucceeded || paneSelected(s) != "first" || promoted != count || discarded != 1 {
		t.Fatal(result, err, promoted, discarded)
	}
}
func TestPresentationStalePreparationDiscardsOnlyItsOwnTicket(t *testing.T) {
	s := paneSession(t)
	h := paneHandle(t, s, "Main/pane/left")
	edit, _ := s.Widget("Main/pane/left/first/edit")
	nested := false
	promoted, discarded := 0, 0
	var accepted Snapshot
	// A violating adapter mutates runtime while preparing. The baseline guard must
	// discard its outer ticket without losing the independently accepted mutation.
	requireOK(t, s.PreparePresentationWith(func(v Snapshot) (PresentationTicket, error) {
		ticket := PresentationTicket{Publish: func() { promoted++; accepted = v }, Discard: func() { discarded++ }}
		if v.Tabs[h.Path].Selected == "second" && !nested {
			nested = true
			requireOK(t, s.Draft(edit.Handle, "independent"))
		}
		return ticket, nil
	}))
	count := promoted
	err := s.SelectPage(h, "second")
	code(t, err, "stale-presentation")
	if paneSelected(s) != "first" || discarded != 1 || promoted != count+1 || accepted.Tabs[h.Path].Selected != "first" {
		t.Fatal(discarded, promoted, accepted.Tabs)
	}
	w, _ := s.Widget(edit.Handle.Path)
	if w.Draft != "independent" {
		t.Fatal(w)
	}
}
func TestPaneInteractionRejectsMalformedAndStaleEnvelopes(t *testing.T) {
	s := paneSession(t)
	h := paneHandle(t, s, "Main/pane/left")
	calls := 0
	requireOK(t, s.BindInteraction(h, func(Event) (InteractionReply, error) { calls++; return InteractionReply{Domain: DomainSucceeded}, nil }))
	changes := []func(*Event){
		func(e *Event) { e.StateRevision-- }, func(e *Event) { e.ModelRevision++ }, func(e *Event) { e.Sequence = 0 },
		func(e *Event) { e.Page.Page.Generation++ }, func(e *Event) { e.Page.PreviousID = "wrong" }, func(e *Event) { e.Page.PageID = "wrong" },
		func(e *Event) { e.Split = &SplitChange{} }, func(e *Event) { e.Value = Text("extra") }, func(e *Event) { e.Page = nil }, func(e *Event) { e.Kind = Activate },
	}
	before := s.Snapshot()
	for _, change := range changes {
		e := paneEvent(t, s, "second")
		change(&e)
		if _, err := s.DispatchInteraction(e); err == nil {
			t.Fatal("bad event accepted", e)
		}
	}
	if calls != 0 || !reflect.DeepEqual(before, s.Snapshot()) {
		t.Fatal("invalid envelope affected state")
	}
	e := paneEvent(t, s, "second")
	code(t, s.Dispatch(e), "event-type")
	split := paneHandle(t, s, "Main/pane")
	for _, change := range []SplitChange{{"bad", 0}, {"restore", 1}, {"ratio", math.NaN()}, {"ratio", math.Inf(1)}} {
		e = Event{Handle: split, ModelRevision: s.Revision, StateRevision: s.StateRevision, Sequence: s.Sequence() + 1, Kind: AdjustSplit, Split: &change}
		if _, err := s.DispatchInteraction(e); err == nil {
			t.Fatal("bad split accepted", change)
		}
	}
}
func TestPaneFinalGateFailureDoesNotAcceptReplyUpdates(t *testing.T) {
	s := paneSession(t)
	h := paneHandle(t, s, "Main/pane/left")
	out, _ := s.Widget("Main/out")
	fail := false
	requireOK(t, s.CheckPresentationWith(func(v Snapshot) (PresentationState, error) {
		if fail {
			return PresentationState{}, errors.New("rejected geometry")
		}
		return paneGate(.2, .8)(v)
	}))
	requireOK(t, s.BindInteraction(h, func(Event) (InteractionReply, error) {
		fail = true
		return InteractionReply{Domain: DomainSucceeded, Updates: []Update{{Handle: out.Handle, Property: AcceptedValue, Value: Text("must not appear"), ExpectedValueRevision: out.ValueRevision}}}, nil
	}))
	result, err := s.DispatchInteraction(paneEvent(t, s, "second"))
	if err == nil || result.Status != "ui-conflict" || result.Domain != DomainSucceeded {
		t.Fatal(result, err)
	}
	w, _ := s.Widget(out.Handle.Path)
	if w.Value != "before" || paneSelected(s) != "first" {
		t.Fatal(w)
	}
}

func TestPaneProbeClampDoesNotReplaceFinalRequestedState(t *testing.T) {
	s := paneSession(t)
	h := paneHandle(t, s, "Main/pane/left")
	split := paneHandle(t, s, "Main/pane")
	requireOK(t, s.SetSplitProportion(split, .7))
	hi := .5
	requireOK(t, s.CheckPresentationWith(func(v Snapshot) (PresentationState, error) {
		upper := .8
		if v.Tabs[h.Path].Selected == "second" {
			upper = hi
		}
		return paneGate(.2, upper)(v)
	}))
	requireOK(t, s.BindInteraction(h, func(Event) (InteractionReply, error) { hi = .8; return InteractionReply{Domain: DomainSucceeded}, nil }))
	_, err := s.DispatchInteraction(paneEvent(t, s, "second"))
	requireOK(t, err)
	p, _ := s.Split(split)
	if p.Proportion != .7 {
		t.Fatal("discarded speculative clamp changed final request", p)
	}
}
