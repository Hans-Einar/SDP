package fynehost

import (
	"context"
	"errors"
	"math"
	"reflect"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/preparation"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

const lifecycleSource = `sdui 0.3; page=[tree=tree("Tree") {x=fill,y=fill,overflow-x=scroll,overflow-y=scroll};edit=input("Preview") {x=fill}] {x=fill,y=fill};`

func lifecycleHost(t *testing.T) (*DocumentHost, fyne.Window) {
	t.Helper()
	a := test.NewTempApp(t)
	w := a.NewWindow("lifecycle regression")
	w.SetPadded(false)
	h := NewDocumentHost(w.Canvas(), layout.Size{W: 800, H: 500})
	t.Cleanup(func() { h.Close(); w.Close() })
	return h, w
}
func lifecycleRequest(t *testing.T, sequence uint64, provider ui.CollectionProvider) DocumentRequest {
	t.Helper()
	d, err := parser.Parse(lifecycleSource)
	if err != nil {
		t.Fatal(err)
	}
	return DocumentRequest{Document: d, Entry: "page", SessionID: "lifecycle", SourceRevision: "source", Sequence: sequence, Mode: preparation.Prototype, Providers: map[string]ui.CollectionProvider{"page/tree": provider}}
}
func lifecycleData(id ui.ItemID) ui.CollectionData {
	return ui.CollectionData{Items: []ui.CollectionItem{{ID: id, Kind: ui.Row, Label: string(id), ChildrenLoaded: true}}}
}
func lifecycleOK(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
func lifecycleState(t *testing.T, h *DocumentHost) ui.CollectionState {
	t.Helper()
	w, ok := h.Current().Session.Widget("page/tree")
	if !ok {
		t.Fatal("missing tree")
	}
	c, ok := h.Current().Session.Collection(w.Handle)
	if !ok {
		t.Fatal("missing collection state")
	}
	return c
}
func lifecycleTarget(t *testing.T, s *ui.Session, id ui.ItemID) ui.CollectionTarget {
	t.Helper()
	w, ok := s.Widget("page/tree")
	if !ok {
		t.Fatal("missing tree")
	}
	target, e := s.Target(w.Handle, id)
	lifecycleOK(t, e)
	return target
}

type lifecycleCall struct {
	ctx     context.Context
	request ui.LoadRequest
	reply   chan ui.CollectionData
}

func lifecycleLoader(t *testing.T, h *DocumentHost) (ui.CollectionProvider, <-chan lifecycleCall, <-chan func()) {
	t.Helper()
	calls := make(chan lifecycleCall, 8)
	posts := make(chan func(), 8)
	stop := make(chan struct{})
	t.Cleanup(func() { close(stop) })
	h.Post = func(fn func()) { posts <- fn }
	p := ui.CollectionProvider{ID: "controlled", Epoch: 1, Load: func(ctx context.Context, r ui.LoadRequest) (ui.CollectionData, error) {
		reply := make(chan ui.CollectionData, 1)
		calls <- lifecycleCall{ctx, r, reply}
		// Deliberately ignore cancellation until explicit delivery: no timing race.
		select {
		case d := <-reply:
			return d, nil
		case <-stop:
			return ui.CollectionData{}, errors.New("test disposed")
		}
	}}
	return p, calls, posts
}
func lifecycleReceive[T any](t *testing.T, ch <-chan T) T {
	t.Helper()
	select {
	case value := <-ch:
		return value
	case <-time.After(5 * time.Second):
		t.Fatal("provider barrier was not reached")
	}
	var zero T
	return zero
}

func TestDocumentLifecycleReentrantResultErrorStillPublishesAcceptedState(t *testing.T) {
	h, _ := lifecycleHost(t)
	p := ui.CollectionProvider{ID: "data", Epoch: 1, RootLoaded: true, Initial: lifecycleData("old")}
	lifecycleOK(t, h.Adopt(lifecycleRequest(t, 1, p)))
	b := h.Current()
	s := b.Session
	target := lifecycleTarget(t, s, "old")
	lifecycleOK(t, s.Bind(target.Handle, func(ui.Event) ([]ui.Update, error) {
		return nil, s.ReplaceCollection(lifecycleTarget(t, s, ""), lifecycleData("new"))
	}))
	err := h.Mutate(func(s *ui.Session) error {
		return s.Dispatch(ui.Event{Handle: target.Handle, ModelRevision: target.ModelRevision, Sequence: s.Sequence() + 1, Kind: ui.Activate, Collection: &target})
	})
	if err == nil {
		t.Fatal("reentrant action result was not rejected")
	}
	if c := lifecycleState(t, h); len(c.Data.Items) != 1 || c.Data.Items[0].ID != "new" {
		t.Fatal("accepted runtime replacement missing", c)
	}
	rows := b.Inspect()["rows"].(map[string][]map[string]any)["page/tree"]
	if len(rows) != 1 || rows[0]["item"] != ui.ItemID("new") {
		t.Fatalf("error left stale native rows after accepted state: %v", rows)
	}
}

func TestDocumentLifecycleReentrantErrorReconcilesRevokedProvider(t *testing.T) {
	h, _ := lifecycleHost(t)
	p, calls, posts := lifecycleLoader(t, h)
	p.RootLoaded = true
	p.Initial = ui.CollectionData{Items: []ui.CollectionItem{{ID: "row", Kind: ui.Row, Label: "Row", ChildrenLoaded: true}, {ID: "lazy", Kind: ui.Row, Label: "Lazy", HasChildren: true}}}
	lifecycleOK(t, h.Adopt(lifecycleRequest(t, 1, p)))
	s := h.Current().Session
	lifecycleOK(t, h.Mutate(func(s *ui.Session) error { return s.ExpandItem(lifecycleTarget(t, s, "lazy")) }))
	flight := lifecycleReceive(t, calls)
	target := lifecycleTarget(t, s, "row")
	lifecycleOK(t, s.Bind(target.Handle, func(ui.Event) ([]ui.Update, error) {
		return nil, s.ReplaceCollection(lifecycleTarget(t, s, ""), lifecycleData("replacement"))
	}))
	if err := h.Mutate(func(s *ui.Session) error {
		return s.Dispatch(ui.Event{Handle: target.Handle, ModelRevision: s.Revision, Sequence: s.Sequence() + 1, Kind: ui.Activate, Collection: &target})
	}); err == nil {
		t.Fatal("stale result accepted")
	}
	if lifecycleState(t, h).Request != nil {
		t.Fatal("replacement did not revoke request")
	}
	select {
	case <-flight.ctx.Done():
	default:
		t.Error("accepted revocation did not cancel provider context after result error")
	}
	flight.reply <- ui.CollectionData{Items: []ui.CollectionItem{{ID: "late", Parent: "lazy", Kind: ui.Row, Label: "Late", ChildrenLoaded: true}}}
	lifecycleReceive(t, posts)()
	if c := lifecycleState(t, h); len(c.Data.Items) != 1 || c.Data.Items[0].ID != "replacement" {
		t.Fatal("revoked reply published", c)
	}
}

func TestDocumentLifecycleResourcesSeeAcceptedViewportOffsets(t *testing.T) {
	h, _ := lifecycleHost(t)
	p := ui.CollectionProvider{ID: "data", Epoch: 1, RootLoaded: true, Initial: lifecycleData("short")}
	r := lifecycleRequest(t, 1, p)
	var prepared ui.Snapshot
	r.PrepareResources = func(s ui.Snapshot) error { prepared = s; return nil }
	lifecycleOK(t, h.Adopt(r))
	b := h.Current()
	handle, ok := b.Session.Viewport("page/tree")
	if !ok {
		t.Fatal("missing viewport")
	}
	lifecycleOK(t, h.Mutate(func(s *ui.Session) error { return s.SetViewport(handle, s.Revision, ui.ViewportState{Y: 999}) }))
	accepted := b.Session.Snapshot().Viewports
	if accepted["page/tree"].Y != 0 {
		t.Fatal("fixture must have zero scroll range", accepted)
	}
	if !reflect.DeepEqual(prepared.Viewports, accepted) {
		t.Fatalf("resource snapshot differs from accepted offsets: resources=%v runtime=%v", prepared.Viewports, accepted)
	}
	if !reflect.DeepEqual(b.presentation.snapshot.Viewports, accepted) {
		t.Fatal("published snapshot retained unclamped offsets")
	}
	if !reflect.DeepEqual(b.Geometry().EffectiveOffsets(), accepted) {
		t.Fatal("runtime and native geometry disagree")
	}
}

func TestDocumentLifecycleFailedResizeStillAllowsCancellation(t *testing.T) {
	h, _ := lifecycleHost(t)
	p, calls, posts := lifecycleLoader(t, h)
	lifecycleOK(t, h.Adopt(lifecycleRequest(t, 1, p)))
	flight := lifecycleReceive(t, calls)
	old := h.Current()
	geometry := old.Geometry()
	before := old.Session.Snapshot()
	if err := h.Resize(layout.Size{W: 1, H: 1}); err == nil {
		t.Fatal("too-small presentation accepted")
	}
	if h.Current() != old || old.Geometry() != geometry || !reflect.DeepEqual(before, old.Session.Snapshot()) {
		t.Fatal("failed resize published partial state")
	}
	if err := h.CancelCollection("page/tree"); err != nil {
		t.Fatalf("failed window geometry blocked lifecycle cancellation: %v", err)
	}
	c := lifecycleState(t, h)
	if c.Request != nil || c.Status[""].Phase != ui.Canceled || c.AutoLoadPending {
		t.Fatal("canceled root not paused", c)
	}
	select {
	case <-flight.ctx.Done():
	default:
		t.Fatal("request revoked without context cancellation")
	}
	flight.reply <- lifecycleData("late")
	lifecycleReceive(t, posts)()
	if c = lifecycleState(t, h); len(c.Data.Items) != 0 || c.Request != nil {
		t.Fatal("late canceled reply published", c)
	}
	lifecycleOK(t, h.Resize(layout.Size{W: 800, H: 500}))
	if c = lifecycleState(t, h); c.Status[""].Phase != ui.Canceled || c.AutoLoadPending {
		t.Fatal("resize restarted canceled root", c)
	}
}

func TestDocumentLifecycleFailedResizeDoesNotStrandCompletedLoad(t *testing.T) {
	h, _ := lifecycleHost(t)
	p, calls, posts := lifecycleLoader(t, h)
	lifecycleOK(t, h.Adopt(lifecycleRequest(t, 1, p)))
	flight := lifecycleReceive(t, calls)
	if err := h.Resize(layout.Size{W: 1, H: 1}); err == nil {
		t.Fatal("too-small presentation accepted")
	}
	flight.reply <- lifecycleData("delivered")
	lifecycleReceive(t, posts)()
	lifecycleOK(t, h.Resize(layout.Size{W: 800, H: 500}))
	c := lifecycleState(t, h)
	if c.Request != nil || !c.RootLoaded || len(c.Data.Items) != 1 || c.Data.Items[0].ID != "delivered" {
		t.Fatalf("finished provider request stranded after geometry recovered: %+v", c)
	}
	select {
	case extra := <-calls:
		extra.reply <- ui.CollectionData{}
		t.Fatal("restoring geometry replayed provider load")
	default:
	}
}

// This exercises Fyne's real focus manager with the software canvas. Actual OS
// key delivery and renderer/cache timing still require the external XTest run.
func TestDocumentLifecycleReloadRetainsActualCanvasFocusForPausedRoot(t *testing.T) {
	h, w := lifecycleHost(t)
	p, calls, posts := lifecycleLoader(t, h)
	lifecycleOK(t, h.Adopt(lifecycleRequest(t, 1, p)))
	first := lifecycleReceive(t, calls)
	w.Resize(fyne.NewSize(800, 500))
	w.SetContent(h.Container)
	old := h.Current()
	control := old.Controls()["page/tree"].(*CollectionControl)
	w.Canvas().Focus(control)
	if w.Canvas().Focused() != control || old.Session.Focused() != "page/tree" {
		t.Fatal("fixture failed to establish native and runtime focus")
	}
	candidate, err := h.Prepare(lifecycleRequest(t, 2, p))
	lifecycleOK(t, err)
	stateRevision := candidate.Session.StateRevision
	var observedFocus fyne.Focusable
	h.OnChange = func(*Bundle) { observedFocus = w.Canvas().Focused() }
	lifecycleOK(t, h.Commit(candidate))
	newControl := candidate.Controls()["page/tree"].(*CollectionControl)
	if observedFocus != newControl {
		t.Fatal("publication observer saw predecessor native focus")
	}
	if newControl == control || w.Canvas().Focused() != newControl {
		t.Fatalf("native focus did not transfer to successor: actual=%T %p expected=%p", w.Canvas().Focused(), w.Canvas().Focused(), newControl)
	}
	if candidate.Session.StateRevision != stateRevision {
		t.Fatal("muted focus restoration synthesized runtime state mutation")
	}
	c := lifecycleState(t, h)
	if c.AutoLoadPending || c.Request != nil || c.Status[""].Phase != ui.Unloaded || candidate.Session.Focused() != "page/tree" {
		t.Fatal("reload did not preserve paused root focus/state", c)
	}
	select {
	case <-first.ctx.Done():
	default:
		t.Fatal("publication did not cancel predecessor flight")
	}
	first.reply <- lifecycleData("late")
	lifecycleReceive(t, posts)()
	// Route through actual canvas focus, not the known new control directly.
	w.Canvas().Focused().TypedRune('r')
	second := lifecycleReceive(t, calls)
	if second.request.RequestID <= first.request.RequestID || second.request.Target.ModelRevision != candidate.Session.Revision {
		t.Fatal("focused recovery did not issue fresh successor request")
	}
	w.Canvas().Focused().TypedRune('r')
	if c = lifecycleState(t, h); c.Request == nil || *c.Request != second.request {
		t.Fatal("repeat R changed the active request")
	}
	second.reply <- lifecycleData("current")
	lifecycleReceive(t, posts)()
	if c = lifecycleState(t, h); !c.RootLoaded || c.Data.Items[0].ID != "current" {
		t.Fatal("focused recovery failed", c)
	}
}

func TestDocumentLifecycleFailedPreparePreservesActiveFlight(t *testing.T) {
	h, _ := lifecycleHost(t)
	p, calls, posts := lifecycleLoader(t, h)
	lifecycleOK(t, h.Adopt(lifecycleRequest(t, 1, p)))
	flight := lifecycleReceive(t, calls)
	old := h.Current()
	before := old.Session.Snapshot()
	geometry := old.Geometry()
	control := old.Controls()["page/tree"]
	request := lifecycleRequest(t, 2, p)
	request.PrepareResources = func(ui.Snapshot) error { return errors.New("injected native resource failure") }
	if _, err := h.Prepare(request); err == nil {
		t.Fatal("failed candidate prepared")
	}
	if h.Current() != old || old.Session.Closed() || old.Geometry() != geometry || old.Controls()["page/tree"] != control || !reflect.DeepEqual(before, old.Session.Snapshot()) {
		t.Fatal("failed preparation changed published state/resources")
	}
	select {
	case <-flight.ctx.Done():
		t.Fatal("failed candidate canceled the predecessor load")
	default:
	}
	flight.reply <- lifecycleData("old-still-live")
	lifecycleReceive(t, posts)()
	if c := lifecycleState(t, h); !c.RootLoaded || c.Request != nil || c.Data.Items[0].ID != "old-still-live" {
		t.Fatal("failed prepare invalidated prior completion", c)
	}
}

func TestDocumentLifecycleInvalidResizePreservesLastValidState(t *testing.T) {
	h, _ := lifecycleHost(t)
	lifecycleOK(t, h.Adopt(lifecycleRequest(t, 1, ui.CollectionProvider{ID: "eager", RootLoaded: true, Initial: lifecycleData("row")})))
	for _, size := range []layout.Size{{W: math.NaN(), H: 500}, {W: 800, H: math.Inf(1)}, {W: 0, H: 500}, {W: 800, H: -1}} {
		b := h.Current()
		before := b.Session.Snapshot()
		geometry := b.Geometry()
		physical := h.size
		if err := h.Resize(size); err == nil {
			t.Fatal("invalid size accepted", size)
		}
		if h.size != physical || b.Geometry() != geometry || !reflect.DeepEqual(before, b.Session.Snapshot()) {
			t.Fatal("invalid native size changed accepted state")
		}
	}
}
