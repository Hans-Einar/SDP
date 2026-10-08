package runtime

import (
	"errors"
	"testing"
)

const contextSource = `sdui 0.3;Main=[items=tree("Items");run=command("Run",context="item",target="items",key="Primary+R");button=button(command="run");menu=menu("Context",mode="context",target="items")[entry=item(command="run")]];`

func TestM2MenuClickedRowAndSelectedButtonStayDistinct(t *testing.T) {
	s, err := New("context", paneRoot(t, contextSource))
	requireOK(t, err)
	h := control(t, s, "Main/items")
	requireOK(t, s.BindProviders(map[string]CollectionProvider{"Main/items": {ID: "items", Epoch: 1, RootLoaded: true, Initial: CollectionData{Items: []CollectionItem{{ID: "one", Kind: Row, ChildrenLoaded: true, Label: "One"}, {ID: "two", Kind: Row, ChildrenLoaded: true, Label: "Two"}}}}}))
	c := control(t, s, "Main/run")
	seen := []ItemID{}
	requireOK(t, s.BindInteraction(c, func(e Event) (InteractionReply, error) {
		seen = append(seen, e.Command.Context.Item.ItemID)
		e.Command.Context.Item.ItemID = "corrupted handler copy"
		return InteractionReply{Domain: DomainSucceeded}, nil
	}))
	if s.Snapshot().Presentations["Main/button"].Enabled {
		t.Fatal("no selection left button enabled")
	}
	one, err := s.Target(h, "one")
	requireOK(t, err)
	two, err := s.Target(h, "two")
	requireOK(t, err)
	requireOK(t, s.SelectItem(one))
	invoke(t, s, "Main/button")
	menu := control(t, s, "Main/menu")
	_, err = s.OpenMenu(menu, ContextTarget{Widget: h, ModelRevision: s.Revision, Item: &two})
	requireOK(t, err)
	event := capture(t, s, "Main/menu/entry", "menu")
	_, err = s.DispatchInteraction(event)
	requireOK(t, err)
	if len(seen) != 2 || seen[0] != "one" || seen[1] != "two" || s.collections[h.Path].Selected != "one" || event.Command.Context.Item.ItemID != "two" {
		t.Fatal(seen, event)
	}
	_, err = s.OpenMenu(menu, ContextTarget{Widget: h, ModelRevision: s.Revision, Item: &two})
	requireOK(t, err)
	event = capture(t, s, "Main/menu/entry", "menu")
	root, err := s.Target(h, "")
	requireOK(t, err)
	requireOK(t, s.ReplaceCollection(root, CollectionData{Items: []CollectionItem{{ID: "two", Kind: Row, ChildrenLoaded: true, Label: "Replacement"}}}))
	_, err = s.DispatchInteraction(event)
	code(t, err, "stale-event")
	if len(seen) != 2 || s.menus[menu.Path].Open {
		t.Fatal("stale row invoked or capture retained")
	}
}
func TestM2StaleMenuEventCannotDismissNewOpening(t *testing.T) {
	s := commandSession(t)
	h := control(t, s, "Main/menu")
	old, err := s.OpenMenu(h, ContextTarget{})
	requireOK(t, err)
	event := capture(t, s, "Main/menu/entry", "menu")
	requireOK(t, s.CloseMenu(old))
	fresh, err := s.OpenMenu(h, ContextTarget{})
	requireOK(t, err)
	_, err = s.DispatchInteraction(event)
	code(t, err, "stale-event")
	if !s.menus[h.Path].Open || s.menus[h.Path].Scope != fresh {
		t.Fatal("stale event dismissed replacement")
	}
}
func TestM2AcceptExactCaptureAndOutcomePairs(t *testing.T) {
	for _, outcome := range []DomainOutcome{DomainUnknown, DomainRejected, DomainNotCalled} {
		t.Run(string(outcome), func(t *testing.T) {
			s := commandSession(t)
			calls := 0
			h := control(t, s, "Main/dialog")
			requireOK(t, s.BindInteraction(h, func(Event) (InteractionReply, error) {
				calls++
				return InteractionReply{Domain: outcome, Accept: &AcceptDecision{Accepted: true}}, nil
			}))
			target := openDialog(t, s, true)
			e := dialogEvent(t, s, target, Accept)
			bad := cloneInteraction(e)
			bad.Dialog.Fields[0].DraftRevision++
			_, err := s.DispatchInteraction(bad)
			code(t, err, "event-type")
			if calls != 0 {
				t.Fatal("bad capture executed")
			}
			r, err := s.DispatchInteraction(e)
			code(t, err, "interaction-reply")
			if r.Domain != outcome || calls != 1 {
				t.Fatal(r, calls)
			}
			d, _ := s.SurfaceState(h)
			if !d.Open || d.AcceptBlocked != (outcome == DomainUnknown) {
				t.Fatal(d)
			}
		})
	}
}
func TestM2AcceptResourceFailureDiscardsTicketOnly(t *testing.T) {
	s := commandSession(t)
	h := control(t, s, "Main/dialog")
	fail := false
	discard, publish := 0, 0
	requireOK(t, s.PreparePresentationWith(func(Snapshot) (PresentationTicket, error) {
		ticket := PresentationTicket{Publish: func() { publish++ }, Discard: func() { discard++ }}
		if fail {
			return ticket, errors.New("native resource")
		}
		return ticket, nil
	}))
	requireOK(t, s.BindInteraction(h, func(Event) (InteractionReply, error) {
		fail = true
		return InteractionReply{Domain: DomainSucceeded, Accept: &AcceptDecision{Accepted: true}}, nil
	}))
	target := openDialog(t, s, true)
	before := publish
	r, err := s.DispatchInteraction(dialogEvent(t, s, target, Accept))
	if err == nil || r.Domain != DomainSucceeded || publish != before || discard != 1 {
		t.Fatal(r, err, publish, discard)
	}
	if len(s.DrainDialogResults()) != 0 {
		t.Fatal("unpublished acceptance receipt")
	}
	d, _ := s.SurfaceState(h)
	if !d.Open || !d.AcceptBlocked {
		t.Fatal(d)
	}
}
func TestM2ReentrantReloadRetainsOutcomeAndDiscardsDraft(t *testing.T) {
	s := commandSession(t)
	h := control(t, s, "Main/dialog")
	requireOK(t, s.BindInteraction(h, func(Event) (InteractionReply, error) {
		requireOK(t, s.Reload(paneRoot(t, commandSource)))
		if len(s.DrainDialogResults()) != 0 {
			t.Fatal("released unfinished reload receipt")
		}
		return InteractionReply{Domain: DomainSucceeded, Accept: &AcceptDecision{Accepted: true}}, nil
	}))
	target := openDialog(t, s, true)
	requireOK(t, s.Draft(control(t, s, "Main/dialog/field"), "old unaccepted"))
	r, err := s.DispatchInteraction(dialogEvent(t, s, target, Accept))
	code(t, err, "stale-result")
	if r.Domain != DomainSucceeded {
		t.Fatal(r)
	}
	results := s.DrainDialogResults()
	if len(results) != 1 || results[0].Surface != target || results[0].Reason != "reload" || results[0].Domain != DomainSucceeded {
		t.Fatal(results)
	}
	w, _ := s.Widget("Main/dialog/field")
	if w.Draft != "saved" {
		t.Fatal(w)
	}
}
