package runtime

import (
	"context"
	"errors"
	"math"
	"reflect"
	"testing"

	"github.com/Hans-Einar/SDP/SDUI/go/parser"
)

const paneSource = `sdui 0.3;
Main=[
 pane=split(axis="horizontal",proportion=0.4,minFirst=0.1,minSecond=0.2)[
  left=tabs("Workspace")[
   first=page("First")[edit=input("Draft",value="saved"); items=tree("Items"){x=fill,y=fill,overflow-y=scroll}];
   second=page("Second")[other=input("Second",value="two")]
  ];
  right=input("Right",value="right")
 ];
 out=input("Output",value="before")
];`

func paneRoot(t *testing.T, source string) *parser.Instance {
	t.Helper()
	_, roots, err := parser.Compile(source)
	requireOK(t, err)
	return roots["Main"]
}
func paneSession(t *testing.T) *Session {
	t.Helper()
	s, err := New("panes", paneRoot(t, paneSource))
	requireOK(t, err)
	requireOK(t, s.BindProviders(map[string]CollectionProvider{"Main/pane/left/first/items": {ID: "items", Epoch: 1, RootLoaded: false, Load: func(context.Context, LoadRequest) (CollectionData, error) {
		t.Error("runtime invoked provider")
		return CollectionData{}, nil
	}}}))
	requireOK(t, s.CheckPresentationWith(paneGate(.2, .8)))
	return s
}
func paneGate(lo, hi float64) PresentationGate {
	return func(v Snapshot) (PresentationState, error) {
		p := PresentationState{Viewports: map[string]ViewportState{}, Splits: map[string]SplitGeometry{}}
		v.Root.Walk(func(n *parser.Instance) {
			if n.Layout["visible"] == false {
				return
			}
			if n.Layout["overflow-x"] == "scroll" || n.Layout["overflow-y"] == "scroll" {
				p.Viewports[n.Path] = v.Viewports[n.Path]
			}
		})
		for path, s := range v.Splits {
			if !s.Visible || s.Collapsed != SplitNone {
				continue
			}
			if lo > hi {
				return p, errors.New("minimum extents do not fit")
			}
			p.Splits[path] = SplitGeometry{lo, hi, math.Max(lo, math.Min(hi, s.Proportion))}
		}
		return p, nil
	}
}
func paneHandle(t *testing.T, s *Session, path string) Handle {
	t.Helper()
	h, ok := s.Pane(path)
	if !ok {
		t.Fatal("pane missing", path)
	}
	return h
}
func paneEvent(t *testing.T, s *Session, id string) Event {
	t.Helper()
	h := paneHandle(t, s, "Main/pane/left")
	tabs, _ := s.Tabs(h)
	e := Event{Handle: h, ModelRevision: s.Revision, StateRevision: s.StateRevision, Sequence: s.Sequence() + 1, Kind: ActivatePage}
	for _, p := range tabs.Pages {
		if p.ID == id {
			e.Page = &PageActivation{tabs.Selected, id, p.Handle}
			return e
		}
	}
	t.Fatal("page missing", id)
	return e
}
func paneSelected(s *Session) string { return s.Snapshot().Tabs["Main/pane/left"].Selected }
func TestPaneNavigationRetainsDraftOffsetAndCancelsOnlyPublishedLoads(t *testing.T) {
	s := paneSession(t)
	tabs := paneHandle(t, s, "Main/pane/left")
	edit, _ := s.Widget("Main/pane/left/first/edit")
	items, _ := s.Widget("Main/pane/left/first/items")
	requireOK(t, s.Draft(edit.Handle, "retained draft"))
	requireOK(t, s.Focus(edit.Handle))
	vh, _ := s.Viewport(items.InstancePath)
	requireOK(t, s.SetViewport(vh, s.Revision, ViewportState{Y: 17}))
	r, err := s.BeginLoad(target(t, s, items.Handle, ""))
	requireOK(t, err)
	before := s.Snapshot()
	result, err := s.DispatchInteraction(paneEvent(t, s, "second"))
	requireOK(t, err)
	if result.Status != "committed" || result.Domain != DomainNotCalled || s.Focused() != tabs.Path {
		t.Fatal(result, s.Focused())
	}
	after := s.Snapshot()
	c := after.Collections[items.InstancePath]
	if c.Request != nil || c.Status[""].Phase != Canceled || c.AutoLoadPending || after.Viewports[items.InstancePath].Y != 17 {
		t.Fatal(c, after.Viewports)
	}
	if err = s.CompleteLoad(r, CollectionData{}, nil); err == nil {
		t.Fatal("hidden request accepted")
	}
	if err = s.Draft(edit.Handle, "hidden edit"); err == nil {
		t.Fatal("hidden input accepted")
	}
	if err = s.SetViewport(vh, s.Revision, ViewportState{Y: 2}); err == nil {
		t.Fatal("hidden native scroll accepted")
	}
	requireOK(t, s.SelectPage(tabs, "first"))
	requireOK(t, s.EnterPage(tabs))
	restored, _ := s.Widget(edit.Handle.Path)
	if restored.Draft != "retained draft" || s.Focused() != edit.Handle.Path || s.Snapshot().Viewports[items.InstancePath].Y != 17 {
		t.Fatal(restored, s.Focused())
	}
	c = state(t, s, items.Handle)
	if c.Request != nil || c.AutoLoadPending {
		t.Fatal("reveal restarted request", c)
	}
	requireOK(t, s.RetryItem(target(t, s, items.Handle, "")))
	next := state(t, s, items.Handle).Request
	requireOK(t, s.RetryItem(target(t, s, items.Handle, "")))
	if next == nil || next.RequestID <= r.RequestID || *next != *state(t, s, items.Handle).Request {
		t.Fatal("retry generation changed twice")
	}
	// Detached pane maps and slices cannot modify live runtime state.
	before.Tabs[tabs.Path].Pages[0].Label = "corrupt"
	before.Splits["Main/pane"] = SplitState{}
	if s.Snapshot().Tabs[tabs.Path].Pages[0].Label == "corrupt" || s.Snapshot().Splits["Main/pane"].Axis == "" {
		t.Fatal("shared pane snapshot")
	}
}
func TestPaneIntentFallbackAndEmptyBody(t *testing.T) {
	s := paneSession(t)
	h := paneHandle(t, s, "Main/pane/left")
	tabs, _ := s.Tabs(h)
	set := func(h Handle, p Property, v bool) {
		requireOK(t, s.Apply(s.Revision, s.BatchRevision+1, []Update{{Handle: h, Property: p, Value: Bool(v)}}))
	}
	requireOK(t, s.EnterPage(h))
	set(tabs.Pages[0].Handle, Enabled, false)
	if paneSelected(s) != "second" || s.Focused() != h.Path {
		t.Fatal("fallback failed", s.Snapshot().Tabs)
	}
	set(h, Visible, false)
	if paneSelected(s) != "second" {
		t.Fatal("hiding erased selection")
	}
	set(tabs.Pages[0].Handle, Enabled, true)
	set(h, Visible, true)
	if paneSelected(s) != "second" {
		t.Fatal("revealing reset retained selection")
	}
	set(tabs.Pages[1].Handle, Visible, false)
	if paneSelected(s) != "first" {
		t.Fatal("hidden selected page survived")
	}
	set(tabs.Pages[0].Handle, Visible, false)
	if paneSelected(s) != "" {
		t.Fatal("no eligible page did not empty body")
	}
	set(tabs.Pages[1].Handle, Visible, true)
	if paneSelected(s) != "second" {
		t.Fatal("eligible fallback missing")
	}
}
func TestSplitMeasuredClampCollapseRestoreAndDisabledGeometry(t *testing.T) {
	s := paneSession(t)
	h := paneHandle(t, s, "Main/pane")
	before := s.Snapshot()
	if err := s.SetSplitProportion(h, .1); err == nil {
		t.Fatal("strict measured ratio clamped")
	}
	if !reflect.DeepEqual(before, s.Snapshot()) {
		t.Fatal("failed strict ratio changed state")
	}
	e := Event{Handle: h, Kind: AdjustSplit, ModelRevision: s.Revision, StateRevision: s.StateRevision, Sequence: s.Sequence() + 1, Split: &SplitChange{Operation: "ratio", Proportion: 2}}
	_, err := s.DispatchInteraction(e)
	requireOK(t, err)
	p, _ := s.Split(h)
	if p.Proportion != .8 {
		t.Fatal(p)
	}
	requireOK(t, s.CollapseSplit(h, SplitFirst))
	p, _ = s.Split(h)
	if p.Collapsed != SplitFirst || p.SavedProportion != .8 {
		t.Fatal(p)
	}
	requireOK(t, s.CheckPresentationWith(paneGate(.9, .7))) // collapsed: expanded minima are suspended
	before = s.Snapshot()
	if err = s.RestoreSplit(h); err == nil {
		t.Fatal("impossible restore accepted")
	}
	if !reflect.DeepEqual(before, s.Snapshot()) {
		t.Fatal("failed restore changed state")
	}
	requireOK(t, s.CheckPresentationWith(paneGate(.3, .6)))
	requireOK(t, s.SetViewports(nil))
	p, _ = s.Split(h)
	if p.SavedProportion != .8 || p.Collapsed != SplitFirst {
		t.Fatal("collapsed resize changed saved ratio", p)
	}
	requireOK(t, s.RestoreSplit(h))
	p, _ = s.Split(h)
	if p.Proportion != .6 || p.SavedProportion != .6 {
		t.Fatal(p)
	}
	requireOK(t, s.Apply(s.Revision, s.BatchRevision+1, []Update{{Handle: h, Property: Enabled, Value: Bool(false)}}))
	p, _ = s.Split(h)
	if p.Enabled || !p.Visible {
		t.Fatal(p)
	}
	// Disabled visible split still requires its geometry entry.
	if err = s.CheckPresentationWith(func(v Snapshot) (PresentationState, error) { return PresentationState{}, nil }); err == nil {
		t.Fatal("disabled visible geometry omitted")
	}
}
