package runtime

import (
	"errors"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"reflect"
	"strings"
	"testing"
)

func TestPaneSuccessorRetainsStateWithoutPublishingOrRevivingRequests(t *testing.T) {
	s := paneSession(t)
	tabs := paneHandle(t, s, "Main/pane/left")
	split := paneHandle(t, s, "Main/pane")
	edit, _ := s.Widget("Main/pane/left/first/edit")
	items, _ := s.Widget("Main/pane/left/first/items")
	requireOK(t, s.Draft(edit.Handle, "hidden draft"))
	requireOK(t, s.Focus(edit.Handle))
	vh, _ := s.Viewport(items.InstancePath)
	requireOK(t, s.SetViewport(vh, s.Revision, ViewportState{Y: 23}))
	request, err := s.BeginLoad(target(t, s, items.Handle, ""))
	requireOK(t, err)
	_, err = s.DispatchInteraction(paneEvent(t, s, "second"))
	requireOK(t, err)
	requireOK(t, s.SetSplitProportion(split, .6))
	requireOK(t, s.CollapseSplit(split, SplitSecond))
	promoted := 0
	requireOK(t, s.PreparePresentationWith(func(Snapshot) (PresentationTicket, error) {
		return PresentationTicket{Publish: func() { promoted++ }}, nil
	}))
	requireOK(t, s.BindInteraction(tabs, func(Event) (InteractionReply, error) {
		t.Error("inherited handler ran")
		return InteractionReply{}, nil
	}))
	before, count := s.Snapshot(), promoted
	n, err := s.Successor(paneRoot(t, paneSource), providersFor(s))
	requireOK(t, err)
	if !reflect.DeepEqual(before, s.Snapshot()) || count != promoted {
		t.Fatal("successor preparation changed predecessor")
	}
	if n.Revision != s.Revision+1 || n.StateRevision != s.StateRevision+1 || n.Sequence() != s.Sequence() || n.HasInteractionBinding(tabs) || n.presentationCheck != nil || n.presentationPrepare != nil {
		t.Fatal("successor epochs/hooks")
	}
	p, _ := n.Split(split)
	pages, _ := n.Tabs(tabs)
	w, _ := n.Widget(edit.Handle.Path)
	if pages.Selected != "second" || p.Collapsed != SplitSecond || p.Proportion != .6 || p.SavedProportion != .6 || w.Draft != "hidden draft" || w.Handle != edit.Handle || n.Snapshot().Viewports[items.InstancePath].Y != 23 {
		t.Fatal(p, pages, w)
	}
	if err = n.CompleteLoad(request, CollectionData{}, nil); err == nil {
		t.Fatal("successor accepted obsolete request")
	}
	if err = n.CheckStateWith(passGate); err == nil {
		t.Fatal("split accepted old untyped gate")
	}
	requireOK(t, n.CheckPresentationWith(paneGate(.2, .8)))
	requireOK(t, n.SelectPage(tabs, "first"))
	requireOK(t, n.EnterPage(tabs))
	if n.Focused() != edit.Handle.Path {
		t.Fatal("remembered focus lost", n.Focused())
	}
	requireOK(t, n.RetryItem(target(t, n, items.Handle, "")))
	if state(t, n, items.Handle).Request.RequestID <= request.RequestID {
		t.Fatal("request watermark reset")
	}
	n.Close()
	if !reflect.DeepEqual(before, s.Snapshot()) || count != promoted {
		t.Fatal("candidate disposal changed predecessor")
	}
}
func TestPaneSuccessorRemovalFallbackAndNewIdentity(t *testing.T) {
	s := paneSession(t)
	tabs := paneHandle(t, s, "Main/pane/left")
	requireOK(t, s.SelectPage(tabs, "second"))
	requireOK(t, s.EnterPage(tabs))
	old := s.Snapshot().Tabs[tabs.Path].Pages[1].Handle
	source := strings.Replace(paneSource, `second=page("Second")[other=input("Second",value="two")]`, `third=page("Third")[other=input("Third",value="three")]`, 1)
	n, err := s.Successor(paneRoot(t, source), providersFor(s))
	requireOK(t, err)
	if paneSelected(n) != "first" || n.Focused() != tabs.Path {
		t.Fatal("removed page fallback", paneSelected(n), n.Focused())
	}
	newPage := n.Snapshot().Tabs[tabs.Path].Pages[1].Handle
	if newPage.Generation <= s.generation {
		t.Fatal("new generation not monotonic", newPage)
	}
	n2, err := n.Successor(paneRoot(t, paneSource), providersFor(n))
	requireOK(t, err)
	replacement := n2.Snapshot().Tabs[tabs.Path].Pages[1].Handle
	if replacement == old || replacement.Generation <= newPage.Generation {
		t.Fatal("removed identity revived", old, replacement)
	}
	e := paneEvent(t, n2, "second")
	e.Page.Page = old
	if _, err = n2.DispatchInteraction(e); err == nil {
		t.Fatal("old page handle accepted")
	}
}
func TestPaneSuccessorSplitAxisResetAndCollapsedFalsePreparation(t *testing.T) {
	s := paneSession(t)
	h := paneHandle(t, s, "Main/pane")
	requireOK(t, s.SetSplitProportion(h, .7))
	requireOK(t, s.CollapseSplit(h, SplitFirst))
	root := paneRoot(t, strings.Replace(paneSource, `axis="horizontal"`, `axis="vertical"`, 1))
	n, err := s.Successor(root, providersFor(s))
	requireOK(t, err)
	p, _ := n.Split(h)
	if p.Axis != "vertical" || p.Collapsed != SplitNone || p.Proportion != .4 {
		t.Fatal("axis state retained", p)
	}
	source := strings.Replace(paneSource, `minSecond=0.2)`, `minSecond=0.2,collapsible=false)`, 1)
	n, err = s.Successor(paneRoot(t, source), providersFor(s))
	requireOK(t, err)
	p, _ = n.Split(h)
	if p.Collapsed != SplitNone || p.Collapsible || p.Proportion != .7 {
		t.Fatal("collapsible=false not expanded", p)
	}
	before := s.Snapshot()
	if err = n.CheckPresentationWith(paneGate(.9, .7)); err == nil {
		t.Fatal("impossible expansion prepared")
	}
	n.Close()
	if !reflect.DeepEqual(before, s.Snapshot()) {
		t.Fatal("failed successor modified predecessor")
	}
}
func TestPaneLegacyReloadPreservesInteractionGuard(t *testing.T) {
	s := paneSession(t)
	h := paneHandle(t, s, "Main/pane/left")
	requireOK(t, s.BindInteraction(h, func(Event) (InteractionReply, error) {
		requireOK(t, s.Reload(paneRoot(t, paneSource)))
		_, err := s.DispatchInteraction(paneEvent(t, s, "second"))
		code(t, err, "reentrant-interaction")
		return InteractionReply{Domain: DomainSucceeded}, nil
	}))
	result, err := s.DispatchInteraction(paneEvent(t, s, "second"))
	if err == nil || result.Status != "ui-conflict" || s.interacting {
		t.Fatal(result, err)
	}
}
func TestTabsOnlyLegacyGateRetainsInactiveOffsets(t *testing.T) {
	source := `sdui 0.3; Main=[tabs=tabs("Tabs")[a=page("A")[items=list("Items"){overflow-y=scroll}];b=page("B")[]]];`
	s, err := New("tabs", paneRoot(t, source))
	requireOK(t, err)
	requireOK(t, s.BindProviders(map[string]CollectionProvider{"Main/tabs/a/items": {ID: "items", RootLoaded: true}}))
	requireOK(t, s.CheckStateWith(func(v Snapshot) (map[string]ViewportState, error) {
		out := map[string]ViewportState{}
		v.Root.Walk(func(n *parser.Instance) {
			if n.Layout["visible"] != false && n.Layout["overflow-y"] == "scroll" {
				out[n.Path] = v.Viewports[n.Path]
			}
		})
		return out, nil
	}))
	requireOK(t, s.SetViewports(map[string]ViewportState{"Main/tabs/a/items": {Y: 19}}))
	h := paneHandle(t, s, "Main/tabs")
	requireOK(t, s.SelectPage(h, "b"))
	if s.Snapshot().Viewports["Main/tabs/a/items"].Y != 19 {
		t.Fatal("legacy gate adapter lost hidden offset")
	}
}
func TestPaneFailedProbeNeverCallsHandler(t *testing.T) {
	s := paneSession(t)
	h := paneHandle(t, s, "Main/pane/left")
	calls := 0
	requireOK(t, s.CheckPresentationWith(func(v Snapshot) (PresentationState, error) {
		if v.Tabs[h.Path].Selected == "second" {
			return PresentationState{}, errors.New("page too large")
		}
		return paneGate(.2, .8)(v)
	}))
	requireOK(t, s.BindInteraction(h, func(Event) (InteractionReply, error) { calls++; return InteractionReply{Domain: DomainSucceeded}, nil }))
	before := s.Snapshot()
	result, err := s.DispatchInteraction(paneEvent(t, s, "second"))
	if err == nil || calls != 0 || result.Domain != DomainNotCalled || !reflect.DeepEqual(before, s.Snapshot()) {
		t.Fatal(result, err, calls)
	}
}
