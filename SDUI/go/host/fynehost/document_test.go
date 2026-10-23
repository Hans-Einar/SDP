package fynehost

import (
	"context"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/preparation"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"testing"
)

const collectionSource = `sdui 0.3; page=[tree=tree("Tree") {x=fill,y=fill,overflow-x=scroll,overflow-y=scroll};edit=input("Preview") {x=fill}] {x=fill,y=fill};`

func documentRequest(t *testing.T, sequence uint64, p ui.CollectionProvider) DocumentRequest {
	t.Helper()
	d, err := parser.Parse(collectionSource)
	if err != nil {
		t.Fatal(err)
	}
	return DocumentRequest{Document: d, Entry: "page", SessionID: "native-test", Sequence: sequence, SourceRevision: "source", Mode: preparation.Prototype, Providers: map[string]ui.CollectionProvider{"page/tree": p}}
}
func documentHost(t *testing.T) *DocumentHost {
	t.Helper()
	a := test.NewTempApp(t)
	w := a.NewWindow("document")
	h := NewDocumentHost(w.Canvas(), layout.Size{W: 800, H: 500})
	t.Cleanup(func() { h.Close(); w.Close() })
	return h
}
func eagerProvider() ui.CollectionProvider {
	return ui.CollectionProvider{ID: "data", Epoch: 1, RootLoaded: true, Initial: ui.CollectionData{Items: []ui.CollectionItem{{ID: "row", Kind: ui.Row, Label: "Row", ChildrenLoaded: true}}}}
}
func TestDocumentPublicationPreservesLiveBundleOnFailure(t *testing.T) {
	h := documentHost(t)
	r := documentRequest(t, 1, eagerProvider())
	if e := h.Adopt(r); e != nil {
		t.Fatal(e)
	}
	old := h.Current()
	control := old.Controls()["page/tree"]
	r = documentRequest(t, 2, eagerProvider())
	r.PrepareResources = func(ui.Snapshot) error { return errors.New("native resource failure") }
	if _, e := h.Prepare(r); e == nil {
		t.Fatal("failed resource admitted")
	}
	if h.Current() != old || old.closed || old.Controls()["page/tree"] != control {
		t.Fatal("failure changed live bundle")
	}
	r = documentRequest(t, 3, eagerProvider())
	next, e := h.Prepare(r)
	if e != nil {
		t.Fatal(e)
	}
	edit, _ := old.Session.Widget("page/edit")
	if e = old.Session.Draft(edit.Handle, "new draft"); e != nil {
		t.Fatal(e)
	}
	if e = h.Commit(next); e == nil || !next.closed || h.Current() != old {
		t.Fatal("stale old state published", e)
	}
	r = documentRequest(t, 4, eagerProvider())
	next, e = h.Prepare(r)
	if e != nil {
		t.Fatal(e)
	}
	if e = h.Commit(next); e != nil {
		t.Fatal(e)
	}
	current, _ := next.Session.Widget("page/edit")
	if current.Draft != "new draft" || !old.closed {
		t.Fatal("successor lost draft or failed disposal")
	}
	if e = h.Commit(next); e == nil {
		t.Fatal("published twice")
	}
	control.(*CollectionControl).TypedRune('r') // retired controls must be inert
}
func TestProviderCancellationLateDeliveryAndExplicitRecovery(t *testing.T) {
	h := documentHost(t)
	type call struct {
		r     ui.LoadRequest
		ctx   context.Context
		reply chan ui.CollectionData
	}
	started := make(chan call, 4)
	posts := make(chan func(), 4)
	h.Post = func(fn func()) { posts <- fn }
	p := ui.CollectionProvider{ID: "lazy", Epoch: 1, Load: func(ctx context.Context, r ui.LoadRequest) (ui.CollectionData, error) {
		reply := make(chan ui.CollectionData, 1)
		started <- call{r, ctx, reply}
		return <-reply, nil
	}}
	prepared, e := h.Prepare(documentRequest(t, 1, p))
	if e != nil {
		t.Fatal(e)
	}
	select {
	case <-started:
		t.Fatal("provider called during preparation")
	default:
	}
	if e = h.Commit(prepared); e != nil {
		t.Fatal(e)
	}
	first := <-started
	if e = h.CancelCollection("page/tree"); e != nil {
		t.Fatal(e)
	}
	select {
	case <-first.ctx.Done():
	default:
		t.Fatal("provider context not canceled")
	}
	state, _ := prepared.Session.Collection(first.r.Target.Handle)
	if state.Status[""].Phase != ui.Canceled || state.AutoLoadPending {
		t.Fatal("cancel not paused", state)
	}
	first.reply <- ui.CollectionData{Items: []ui.CollectionItem{{ID: "late", Kind: ui.Row, Label: "Late", ChildrenLoaded: true}}}
	(<-posts)()
	state, _ = prepared.Session.Collection(first.r.Target.Handle)
	if len(state.Data.Items) != 0 {
		t.Fatal("late data published")
	}
	c := prepared.Controls()["page/tree"].(*CollectionControl)
	c.TypedRune('r')
	second := <-started
	if second.r.RequestID <= first.r.RequestID {
		t.Fatal("reused token")
	}
	c.TypedRune('r')
	state, _ = prepared.Session.Collection(second.r.Target.Handle)
	if state.Request == nil || state.Request.RequestID != second.r.RequestID {
		t.Fatal("duplicate retry changed token")
	}
	second.reply <- ui.CollectionData{Items: []ui.CollectionItem{{ID: "current", Kind: ui.Row, Label: "Current", ChildrenLoaded: true}}}
	(<-posts)()
	state, _ = prepared.Session.Collection(second.r.Target.Handle)
	if len(state.Data.Items) != 1 || state.Data.Items[0].ID != "current" || state.Request != nil {
		t.Fatal("valid completion missing", state)
	}
}
func TestCollectionLocalNavigationAndHiddenClipping(t *testing.T) {
	h := documentHost(t)
	if e := h.Adopt(documentRequest(t, 1, eagerProvider())); e != nil {
		t.Fatal(e)
	}
	b := h.Current()
	c := b.Controls()["page/tree"].(*CollectionControl)
	calls := 0
	if e := b.Session.Bind(c.state.Handle, func(ui.Event) ([]ui.Update, error) { calls++; return nil, nil }); e != nil {
		t.Fatal(e)
	}
	h.focusCollection(c)
	h.collectionEvent(c, "row", ui.Select)
	if calls != 0 {
		t.Fatal("selection called action")
	}
	h.collectionEvent(c, "row", ui.Activate)
	if calls != 1 {
		t.Fatal("activation missing")
	}
	if _, ok := c.rowAt(fyne.NewPos(4, 1)); ok {
		t.Fatal("title hit row")
	}
}

// Exercise Fyne's actual focus traversal over the mounted object tree, including
// empty collections: map iteration order must never become native Tab order.
func TestDocumentNativeTabOrderAndDisabledSkip(t *testing.T) {
	for _, disabled := range []bool{false, true} {
		t.Run(fmt.Sprint(disabled), func(t *testing.T) {
			a := test.NewTempApp(t)
			w := a.NewWindow("tab order")
			h := NewDocumentHost(w.Canvas(), layout.Size{W: 800, H: 500})
			defer w.Close()
			defer h.Close()
			enabled := "true"
			if disabled {
				enabled = "false"
			}
			d, err := parser.Parse(`sdui 0.3; page=[body=[nav=tree("Navigation") {x=fill,y=fill,enabled=` + enabled + `},entries=list("Entries") {x=fill,y=fill}] {x=fill,y=fill};footer=[preview=input("Preview") {x=fill}]] {x=fill,y=fill};`)
			if err != nil {
				t.Fatal(err)
			}
			r := DocumentRequest{Document: d, Entry: "page", SessionID: "tab", SourceRevision: "tab-source", Sequence: 1, Mode: preparation.Prototype, Providers: map[string]ui.CollectionProvider{
				"page/body/nav": {ID: "nav", Epoch: 1, RootLoaded: true}, "page/body/entries": {ID: "entries", Epoch: 1, RootLoaded: true},
			}}
			if err = h.Adopt(r); err != nil {
				t.Fatal(err)
			}
			w.SetPadded(false)
			w.SetContent(h.Container)
			w.Resize(fyne.NewSize(800, 500))
			w.Show()
			w.Canvas().Unfocus()
			paths := []string{"page/body/nav", "page/body/entries", "page/footer/preview"}
			if disabled {
				paths = paths[1:]
			}
			for _, path := range paths {
				w.Canvas().FocusNext()
				if got := w.Canvas().Focused(); any(got) != any(h.Current().Controls()[path]) {
					t.Fatalf("Tab expected %s, got %T %p", path, got, got)
				}
			}
		})
	}
}

func TestCollectionPartiallyClippedTextHasPaintBounds(t *testing.T) {
	h := documentHost(t)
	if err := h.Adopt(documentRequest(t, 1, eagerProvider())); err != nil {
		t.Fatal(err)
	}
	c := h.Current().Controls()["page/tree"].(*CollectionControl)
	c.viewport.Offset.X = 18 // inspect rendering of a label crossing the left edge
	r := c.CreateRenderer().(*collectionRenderer)
	r.Refresh()
	defer r.Destroy()
	found := false
	for _, o := range r.box.Objects {
		clip, ok := o.(*container.Clip)
		if !ok {
			continue
		}
		for _, child := range clip.Content.(*fyne.Container).Objects {
			text, ok := child.(*canvas.Text)
			if !ok {
				continue
			}
			if text.Position().X >= 0 {
				t.Fatal("regression requires negative label origin")
			}
			if text.Size().Width <= -text.Position().X || text.Size().Height <= 0 {
				t.Fatalf("visible text has empty/cullable bounds: %v %v", text.Position(), text.Size())
			}
			found = true
		}
	}
	if !found {
		t.Fatal("no partial row text")
	}
}

func TestEmptyReloadPreservesNativeFocusAndExplicitLoadKey(t *testing.T) {
	a := test.NewTempApp(t)
	w := a.NewWindow("reload focus")
	h := NewDocumentHost(w.Canvas(), layout.Size{W: 800, H: 500})
	defer w.Close()
	defer h.Close()
	started := make(chan ui.LoadRequest, 4)
	posts := make(chan func(), 4)
	h.Post = func(f func()) { posts <- f }
	p := ui.CollectionProvider{ID: "lazy", Epoch: 1, Load: func(ctx context.Context, r ui.LoadRequest) (ui.CollectionData, error) {
		started <- r
		<-ctx.Done()
		return ui.CollectionData{}, ctx.Err()
	}}
	if err := h.Adopt(documentRequest(t, 1, p)); err != nil {
		t.Fatal(err)
	}
	<-started
	w.SetPadded(false)
	w.SetContent(h.Container)
	w.Resize(fyne.NewSize(800, 500))
	w.Show()
	w.Canvas().FocusNext()
	old := h.Current().Controls()["page/tree"]
	if any(w.Canvas().Focused()) != any(old) {
		t.Fatal("initial collection focus missing")
	}
	if err := h.Adopt(documentRequest(t, 2, p)); err != nil {
		t.Fatal(err)
	}
	current := h.Current().Controls()["page/tree"]
	if any(w.Canvas().Focused()) != any(current) {
		t.Fatalf("reload lost native focus: got %T, runtime=%s", w.Canvas().Focused(), h.Current().Session.Focused())
	}
	state := h.Current().Session.Snapshot().Collections["page/tree"]
	if state.Request != nil || state.AutoLoadPending || state.Status[""].Phase != ui.Unloaded {
		t.Fatal("reload must pause root load", state)
	}
	w.Canvas().Focused().TypedRune('r')
	state = h.Current().Session.Snapshot().Collections["page/tree"]
	if state.Request == nil || state.Status[""].Phase != ui.Loading {
		t.Fatal("focused R did not restart explicit load", state)
	}
	<-started
}
