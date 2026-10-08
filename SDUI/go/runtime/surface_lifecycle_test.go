package runtime

import (
	"context"
	"errors"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"reflect"
	"testing"
)

func TestM2ParentHideFailurePreservesSurfaceDraftRequestAndReceipt(t *testing.T) {
	src := `sdui 0.3;Main=[pages=tabs("Pages")[a=page("A")[open=button("Open")];b=page("B")[other=input("Other")]];dialog=dialog("D",modal=false)[field=input("Field",value="saved");items=tree("Items")]];`
	s, err := New("lifetime", paneRoot(t, src))
	requireOK(t, err)
	requireOK(t, s.BindProviders(map[string]CollectionProvider{"Main/dialog/items": {ID: "p", Epoch: 1, Load: func(context.Context, LoadRequest) (CollectionData, error) {
		t.Fatal("domain load called")
		return CollectionData{}, nil
	}}}))
	h := control(t, s, "Main/dialog")
	target, err := s.OpenSurfaceFrom(h, control(t, s, "Main/pages/a/open"), ContextTarget{})
	requireOK(t, err)
	requireOK(t, s.ConfirmSurfacePublication(target))
	requireOK(t, s.Draft(control(t, s, "Main/dialog/field"), "draft"))
	root, err := s.Target(control(t, s, "Main/dialog/items"), "")
	requireOK(t, err)
	request, err := s.BeginLoad(root)
	requireOK(t, err)
	reject := false
	requireOK(t, s.CheckPresentationWith(func(Snapshot) (PresentationState, error) {
		if reject {
			return PresentationState{}, errors.New("geometry")
		}
		return PresentationState{Viewports: map[string]ViewportState{}}, nil
	}))
	before := s.Snapshot()
	reject = true
	err = s.SelectPage(control(t, s, "Main/pages"), "b")
	if err == nil || !reflect.DeepEqual(before, s.Snapshot()) {
		t.Fatal("failed page change altered live state", err)
	}
	if len(s.DrainDialogResults()) != 0 {
		t.Fatal("failed transition emitted receipt")
	}
	reject = false
	requireOK(t, s.SelectPage(control(t, s, "Main/pages"), "b"))
	d, _ := s.SurfaceState(h)
	if d.Open {
		t.Fatal("opener hidden without closing dependent surface")
	}
	if err = s.CompleteLoad(request, CollectionData{}, nil); err == nil {
		t.Fatal("old load survived parent hide")
	}
	got := s.DrainDialogResults()
	if len(got) != 1 || got[0].Reason != "parent-hidden" {
		t.Fatal(got)
	}
	w, _ := s.Widget("Main/dialog/field")
	if w.Draft != "saved" {
		t.Fatal(w)
	}
}
func TestM2ForcedRevocationCannotBeVetoedAndCandidateDoesNotEmit(t *testing.T) {
	s := commandSession(t)
	target := openDialog(t, s, true)
	reject := false
	requireOK(t, s.CheckPresentationWith(func(Snapshot) (PresentationState, error) {
		if reject {
			return PresentationState{}, errors.New("unavailable geometry")
		}
		return PresentationState{Viewports: map[string]ViewportState{}}, nil
	}))
	reject = true
	requireOK(t, s.RevokeSurfaces(Handle{}, "parent-closed"))
	requireOK(t, s.RevokeSurfaces(Handle{}, "parent-closed"))
	results := s.DrainDialogResults()
	if len(results) != 1 || results[0].Surface != target || results[0].Reason != "parent-closed" {
		t.Fatal(results)
	}
	s.Close()
	if len(s.DrainDialogResults()) != 0 {
		t.Fatal("duplicate teardown receipt")
	}
	candidate := commandSession(t)
	openDialog(t, candidate, false)
	candidate.Close()
	if len(candidate.DrainDialogResults()) != 0 {
		t.Fatal("unpublished candidate emitted")
	}
}
func TestM2SuccessorClosedAndMonotonicTokens(t *testing.T) {
	s := commandSession(t)
	invoke(t, s, "Main/invoke")
	target := openDialog(t, s, true)
	requireOK(t, s.Draft(control(t, s, "Main/dialog/field"), "dirty"))
	before := s.Snapshot()
	n, err := s.Successor(paneRoot(t, commandSource), nil)
	requireOK(t, err)
	if !reflect.DeepEqual(before, s.Snapshot()) {
		t.Fatal("successor mutated predecessor")
	}
	d, _ := n.SurfaceState(control(t, n, "Main/dialog"))
	c, _ := n.CommandState(control(t, n, "Main/choice"))
	if d.Open || !c.Checked || n.Sequence() != s.Sequence() {
		t.Fatal(d, c)
	}
	next := openDialog(t, n, true)
	if next.OpenGeneration <= target.OpenGeneration || next.ModelRevision == target.ModelRevision {
		t.Fatal(next, target)
	}
	code(t, n.CloseSurface(target, "close"), "stale-surface")
	requireOK(t, s.RevokeSurfaces(Handle{}, "reload"))
	s.Close()
	results := s.DrainDialogResults()
	if len(results) != 1 || results[0].Reason != "reload" {
		t.Fatal(results)
	}
}
func TestM2ModalRoutingAndNonmodalFocus(t *testing.T) {
	s := commandSession(t)
	main := control(t, s, "Main/out")
	requireOK(t, s.Focus(main))
	target := openDialog(t, s, true)
	requireOK(t, s.FocusSurface(nil))
	requireOK(t, s.CloseSurface(target, "close"))
	if s.Focused() != main.Path {
		t.Fatal("unfocused dialog stole focus", s.Focused())
	}
	s.DrainDialogResults()
	src := `sdui 0.3;Main=[open=button("Open");out=input("Outside");dialog=dialog("Modal")[field=input("Field");childOpen=button("Child")];child=dialog("Child")[inside=input("Inside")]];`
	s, err := New("modal", paneRoot(t, src))
	requireOK(t, err)
	target = openDialog(t, s, true)
	code(t, s.Draft(control(t, s, "Main/out"), "blocked"), "draft")
	code(t, s.Focus(control(t, s, "Main/out")), "focus")
	code(t, s.FocusSurface(nil), "modal")
	child, err := s.OpenSurfaceFrom(control(t, s, "Main/child"), control(t, s, "Main/dialog/childOpen"), ContextTarget{})
	requireOK(t, err)
	requireOK(t, s.ConfirmSurfacePublication(child))
	requireOK(t, s.Draft(control(t, s, "Main/child/inside"), "allowed"))
	code(t, s.Draft(control(t, s, "Main/dialog/field"), "blocked"), "draft")
	requireOK(t, s.CloseSurface(target, "close"))
	results := s.DrainDialogResults()
	if len(results) != 2 || results[0].Surface != child || results[0].Reason != "parent-closed" {
		t.Fatal(results)
	}
}
func TestM2SnapshotDetachedAndReopenKeepsDraft(t *testing.T) {
	s := commandSession(t)
	target := openDialog(t, s, true)
	field := control(t, s, "Main/dialog/field")
	requireOK(t, s.Draft(field, "draft"))
	again, err := s.OpenSurfaceFrom(target.Handle, control(t, s, "Main/open"), ContextTarget{})
	requireOK(t, err)
	if again != target {
		t.Fatal("reopen minted new token")
	}
	before := s.Snapshot()
	snap := s.Snapshot()
	d := snap.Surfaces["Main/dialog"]
	d.Open = false
	snap.Surfaces["Main/dialog"] = d
	snap.ActiveSurface.OpenGeneration++
	snap.Commands["Main/choice"] = CommandState{}
	snap.Root.Uses = append(snap.Root.Uses, parser.UseSite{})
	if !reflect.DeepEqual(before, s.Snapshot()) {
		t.Fatal("snapshot leaked")
	}
	w, _ := s.Widget(field.Path)
	if w.Draft != "draft" {
		t.Fatal(w)
	}
}

func TestM2ReloadDiscardsOnlyUnacceptedDialogDrafts(t *testing.T) {
	s := commandSession(t)
	target := openDialog(t, s, true)
	field := control(t, s, "Main/dialog/field")
	main := control(t, s, "Main/out")
	requireOK(t, s.Draft(field, "explicitly committed"))
	w, _ := s.Widget(field.Path)
	requireOK(t, s.Apply(s.Revision, s.BatchRevision+1, []Update{{Handle: field, Property: AcceptedValue, Value: Text(w.Draft), ExpectedValueRevision: w.ValueRevision, AcceptDraft: true}}))
	requireOK(t, s.Draft(field, "discard on reload"))
	requireOK(t, s.Draft(main, "retain main draft"))
	before := s.Snapshot()
	n, err := s.Successor(paneRoot(t, commandSource), nil)
	requireOK(t, err)
	if !reflect.DeepEqual(before, s.Snapshot()) {
		t.Fatal("candidate changed live opening")
	}
	w, _ = n.Widget(field.Path)
	outside, _ := n.Widget(main.Path)
	if w.Value != "explicitly committed" || w.Draft != w.Value || w.Dirty || outside.Draft != "retain main draft" || !outside.Dirty {
		t.Fatal(w, outside)
	}
	next := openDialog(t, n, true)
	if next == target {
		t.Fatal("reused opening")
	}
	w, _ = n.Widget(field.Path)
	if w.Draft != "explicitly committed" {
		t.Fatal(w)
	}
	requireOK(t, s.RevokeSurfaces(Handle{}, "reload"))
	results := s.DrainDialogResults()
	if len(results) != 1 || results[0].Reason != "reload" || len(results[0].Fields) != 0 {
		t.Fatal(results)
	}
}

func TestM2ConvenienceReopenFocusedModalKeepsOpening(t *testing.T) {
	s, err := New("modal", paneRoot(t, `sdui 0.3;Main=[open=button("Open");dialog=dialog("D")[field=input("Field")]];`))
	requireOK(t, err)
	first := openDialog(t, s, true)
	requireOK(t, s.Draft(control(t, s, "Main/dialog/field"), "dirty"))
	again, err := s.OpenSurface(first.Handle, ContextTarget{})
	requireOK(t, err)
	if first != again {
		t.Fatal(first, again)
	}
	w, _ := s.Widget("Main/dialog/field")
	if w.Draft != "dirty" {
		t.Fatal(w)
	}
}

// Regressions reproduced independently during the M2 review.
func TestM2CloseDisabledOpenerRestoresFirstParentControl(t *testing.T) {
	s := commandSession(t)
	target := openDialog(t, s, true)
	opener := control(t, s, "Main/open")
	requireOK(t, s.Apply(s.Revision, s.BatchRevision+1, []Update{{Handle: opener, Property: Enabled, Value: Bool(false)}}))
	requireOK(t, s.CloseSurface(target, "cancel"))
	if got := s.Focused(); got != "Main/invoke" {
		t.Fatalf("focus=%q, want first eligible main control Main/invoke", got)
	}
}
func TestM2KeyFromInactiveNonmodalSurface(t *testing.T) {
	s := commandSession(t)
	openDialog(t, s, true)
	if _, err := s.CaptureCommand(control(t, s, "Main/choice"), "key", nil); err == nil {
		t.Fatal("captured key from inactive main canvas while dialog is active")
	}
}
func TestM2ReopenSurfaceAfterFocusedControlDisabled(t *testing.T) {
	s := commandSession(t)
	target := openDialog(t, s, true)
	field := control(t, s, "Main/dialog/field")
	requireOK(t, s.Apply(s.Revision, s.BatchRevision+1, []Update{{Handle: field, Property: Enabled, Value: Bool(false)}}))
	_, err := s.OpenSurface(target.Handle, ContextTarget{})
	requireOK(t, err)
	if got := s.Focused(); got != "Main/dialog/ok" {
		t.Fatalf("focus=%q, want first eligible dialog control Main/dialog/ok", got)
	}
}

func TestM2TabCallbackUpdatesCommand(t *testing.T) {
	src := `sdui 0.3;Main=[flag=command("Flag",toggle=true);tabs=tabs("Tabs")[one=page("One")[input("A")];two=page("Two")[input("B")]]];`
	s, err := New("review", paneRoot(t, src))
	requireOK(t, err)
	flag := control(t, s, "Main/flag")
	h := control(t, s, "Main/tabs")
	requireOK(t, s.BindInteraction(h, func(Event) (InteractionReply, error) {
		return InteractionReply{Domain: DomainSucceeded, Updates: []Update{{Handle: flag, Property: Checked, Value: Bool(true)}}}, nil
	}))
	tabs, _ := s.Tabs(h)
	e := Event{Handle: h, Kind: ActivatePage, ModelRevision: s.Revision, StateRevision: s.StateRevision, Sequence: s.Sequence() + 1, Page: &PageActivation{PreviousID: "one", PageID: "two", Page: tabs.Pages[1].Handle}}
	r, err := s.DispatchInteraction(e)
	if err != nil {
		t.Fatalf("valid tab callback command update rejected: %+v %v", r, err)
	}
}

func TestM2NativeExactSurfaceLossBypassesGateButRejectsReplacement(t *testing.T) {
	s := commandSession(t)
	target := openDialog(t, s, true)
	reject := false
	requireOK(t, s.CheckPresentationWith(func(Snapshot) (PresentationState, error) {
		if reject {
			return PresentationState{}, errors.New("geometry unavailable after native close")
		}
		return PresentationState{Viewports: map[string]ViewportState{}}, nil
	}))
	reject = true
	if err := s.CloseSurface(target, "close"); err == nil {
		t.Fatal("ordinary close unexpectedly passed gate")
	}
	requireOK(t, s.RevokeSurface(target, "parent-closed"))
	requireOK(t, s.RevokeSurface(target, "parent-closed"))
	got := s.DrainDialogResults()
	if len(got) != 1 || got[0].Surface != target || got[0].Kind != "close" || got[0].Reason != "parent-closed" {
		t.Fatal(got)
	}
	reject = false
	replacement := openDialog(t, s, true)
	code(t, s.RevokeSurface(target, "parent-closed"), "stale-surface")
	d, _ := s.SurfaceState(replacement.Handle)
	if !d.Open || d.Target != replacement {
		t.Fatal("old native close revoked replacement", d)
	}
}
