package runtime

import (
	"errors"
	"math"
	"reflect"
	"strings"
	"testing"
)

func TestPanePresentationBoundsValidateAndFailedInstallRetainsGate(t *testing.T) {
	for name, geometry := range map[string]SplitGeometry{
		"nan": {math.NaN(), .8, .4}, "inverted": {.7, .3, .4}, "wrong-clamp": {.2, .8, .5}, "source-minimum": {0, .8, .4}, "outside": {.2, 1.1, .4},
	} {
		t.Run(name, func(t *testing.T) {
			s := paneSession(t)
			before := s.Snapshot()
			err := s.CheckPresentationWith(func(Snapshot) (PresentationState, error) {
				return PresentationState{Splits: map[string]SplitGeometry{"Main/pane": geometry}}, nil
			})
			if err == nil || !reflect.DeepEqual(before, s.Snapshot()) {
				t.Fatal(err)
			}
			requireOK(t, s.SetSplitProportion(paneHandle(t, s, "Main/pane"), .5))
		})
	}
	s := paneSession(t)
	before := s.Snapshot()
	err := s.CheckPresentationWith(func(v Snapshot) (PresentationState, error) {
		p, _ := paneGate(.2, .8)(v)
		p.Splits["unknown"] = SplitGeometry{.2, .8, .4}
		return p, nil
	})
	if err == nil || !reflect.DeepEqual(before, s.Snapshot()) {
		t.Fatal("extra geometry accepted")
	}
	err = s.CheckPresentationWith(func(v Snapshot) (PresentationState, error) {
		p, _ := paneGate(.2, .8)(v)
		p.Viewports["unknown"] = ViewportState{}
		return p, nil
	})
	if err == nil || !reflect.DeepEqual(before, s.Snapshot()) {
		t.Fatal("unknown viewport accepted")
	}
}
func TestPaneFocusFallbackAndCollectionMemory(t *testing.T) {
	s := paneSession(t)
	tabs := paneHandle(t, s, "Main/pane/left")
	edit, _ := s.Widget("Main/pane/left/first/edit")
	requireOK(t, s.Focus(edit.Handle))
	requireOK(t, s.Apply(s.Revision, s.BatchRevision+1, []Update{{Handle: edit.Handle, Property: Visible, Value: Bool(false)}}))
	if s.Focused() != tabs.Path {
		t.Fatal("invalid content focus did not return to header", s.Focused())
	}
	items, _ := s.Widget("Main/pane/left/first/items")
	requireOK(t, s.FocusItem(target(t, s, items.Handle, "")))
	requireOK(t, s.SelectPage(tabs, "second"))
	requireOK(t, s.SelectPage(tabs, "first"))
	requireOK(t, s.EnterPage(tabs))
	if s.Focused() != items.Handle.Path {
		t.Fatal("collection page focus not remembered", s.Focused())
	}
}
func TestPaneCallbackOwnerDiscoveryAndUnboundNavigation(t *testing.T) {
	source := strings.Replace(paneSource, `sdui 0.3;`, `sdui 0.3; ref: actions "actions.sdl";`, 1)
	source = strings.Replace(source, `tabs("Workspace")`, `tabs("Workspace",callback=actions.Page.@invoke)`, 1)
	s, err := New("bindings", paneRoot(t, source))
	requireOK(t, err)
	h := paneHandle(t, s, "Main/pane/left")
	owners := s.CallbackOwners()
	if len(owners) != 1 || owners[0].Handle != h || owners[0].Binding.Object != "Page" {
		t.Fatal(owners)
	}
	for _, w := range s.Widgets() {
		if w.Handle.Kind == "tabs" || w.Handle.Kind == "page" || w.Handle.Kind == "split" {
			t.Fatal("legacy Widgets changed")
		}
	}
	requireOK(t, s.CheckPresentationWith(paneGate(.2, .8)))
	before := s.Snapshot()
	_, err = s.DispatchInteraction(paneEvent(t, s, "second"))
	code(t, err, "unbound")
	if !reflect.DeepEqual(before, s.Snapshot()) {
		t.Fatal("unbound navigation published")
	}
	requireOK(t, s.SelectPage(h, "second")) // programmatic navigation does not require binding
	requireOK(t, s.BindInteraction(h, func(Event) (InteractionReply, error) { return InteractionReply{Domain: DomainSucceeded}, nil }))
	if !s.HasInteractionBinding(h) {
		t.Fatal("binding not discoverable")
	}
	s.Close()
	if s.HasInteractionBinding(h) {
		t.Fatal("closed binding exposed")
	}
}
func TestPaneReplyTargetAndDraftConflictAreAtomic(t *testing.T) {
	for _, kind := range []string{"reserved-pane", "dirty-target", "unknown-domain", "error-domain"} {
		t.Run(kind, func(t *testing.T) {
			s := paneSession(t)
			h := paneHandle(t, s, "Main/pane/left")
			out, _ := s.Widget("Main/out")
			if kind == "dirty-target" {
				requireOK(t, s.Draft(out.Handle, "unaccepted"))
			}
			requireOK(t, s.BindInteraction(h, func(Event) (InteractionReply, error) {
				r := InteractionReply{Domain: DomainSucceeded, Updates: []Update{{Handle: out.Handle, Property: AcceptedValue, Value: Text("reply"), ExpectedValueRevision: out.ValueRevision}}}
				switch kind {
				case "reserved-pane":
					r.Updates = append(r.Updates, Update{Handle: h, Property: Visible, Value: Bool(false)})
				case "unknown-domain":
					r.Domain = "invalid"
				case "error-domain":
					return InteractionReply{}, errors.New("unknown execution result")
				}
				return r, nil
			}))
			result, err := s.DispatchInteraction(paneEvent(t, s, "second"))
			if err == nil || result.Status != "ui-conflict" || paneSelected(s) != "first" {
				t.Fatal(kind, result, err)
			}
			w, _ := s.Widget(out.Handle.Path)
			if w.Value != "before" {
				t.Fatal("reply partially published")
			}
			if (kind == "unknown-domain" || kind == "error-domain") && result.Domain != DomainUnknown {
				t.Fatal(result)
			}
		})
	}
}
func TestPaneFailedPropertyBatchAndGatePreserveIntentAndLoad(t *testing.T) {
	s := paneSession(t)
	h := paneHandle(t, s, "Main/pane/left")
	tabs, _ := s.Tabs(h)
	items, _ := s.Widget("Main/pane/left/first/items")
	_, err := s.BeginLoad(target(t, s, items.Handle, ""))
	requireOK(t, err)
	requireOK(t, s.CheckPresentationWith(func(v Snapshot) (PresentationState, error) {
		if v.Tabs[h.Path].Selected == "second" {
			return PresentationState{}, errors.New("reject fallback")
		}
		return paneGate(.2, .8)(v)
	}))
	before := s.Snapshot()
	err = s.Apply(s.Revision, s.BatchRevision+1, []Update{{Handle: tabs.Pages[0].Handle, Property: Visible, Value: Bool(false)}})
	if err == nil || !reflect.DeepEqual(before, s.Snapshot()) {
		t.Fatal("failed property gate lost intent/request")
	}
	for _, label := range []string{"", "\xff"} {
		err = s.Apply(s.Revision, s.BatchRevision+1, []Update{{Handle: h, Property: Label, Value: Text(label)}})
		if err == nil {
			t.Fatal("invalid pane label accepted")
		}
	}
}
