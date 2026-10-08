package runtime

import (
	"errors"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"reflect"
	"testing"
)

func providersFor(s *Session) map[string]CollectionProvider {
	p := map[string]CollectionProvider{}
	for _, w := range s.Widgets() {
		if x, ok := s.Provider(w.Handle); ok {
			p[w.InstancePath] = x
		}
	}
	return p
}
func passGate(v Snapshot) (map[string]ViewportState, error) { return v.Viewports, nil }
func TestDetachedSuccessorRetainsCompatibleStateAndNoHandlers(t *testing.T) {
	s, h := collectionSession(t, "tree", CollectionData{[]CollectionItem{branch("loaded", "", true), row("child", "loaded"), branch("lazy", "", false), branch("error", "", false), branch("canceled", "", false)}}, true, true)
	requireOK(t, s.ExpandItem(target(t, s, h, "loaded")))
	requireOK(t, s.SelectItem(target(t, s, h, "child")))
	requireOK(t, s.FocusItem(target(t, s, h, "child")))
	requireOK(t, s.ExpandItem(target(t, s, h, "error")))
	errReq := *state(t, s, h).Request
	requireOK(t, s.CompleteLoad(errReq, CollectionData{}, errors.New("stable error")))
	requireOK(t, s.ExpandItem(target(t, s, h, "canceled")))
	requireOK(t, s.CancelLoad(h))
	requireOK(t, s.ExpandItem(target(t, s, h, "lazy")))
	oldRequest := *state(t, s, h).Request
	requireOK(t, s.CheckStateWith(passGate))
	w, _ := s.Widget(h.Path)
	vh, _ := s.Viewport(w.InstancePath)
	requireOK(t, s.SetViewport(vh, s.Revision, ViewportState{X: 5, Y: 7}))
	edit, _ := s.Widget("page/edit")
	requireOK(t, s.Draft(edit.Handle, "user draft"))
	requireOK(t, s.Focus(edit.Handle))
	calls := 0
	requireOK(t, s.Bind(h, func(Event) ([]Update, error) { calls++; return nil, nil }))
	requireOK(t, s.Dispatch(eventFor(s, target(t, s, h, "child"), Activate)))
	before := s.Snapshot()
	oldState := state(t, s, h)
	n, err := s.Successor(collectionRoot(t, "tree"), providersFor(s))
	requireOK(t, err)
	unchanged(t, s, before, s.Revision, s.BatchRevision)
	if n.Revision != s.Revision+1 || n.StateRevision != s.StateRevision+1 || n.Sequence() != s.Sequence() || n.Focused() != s.Focused() {
		t.Fatal("successor guards not carried")
	}
	nw, _ := n.Widget("page/edit")
	if nw.Handle != edit.Handle || nw.Draft != "user draft" || !nw.Dirty || n.HasBinding(h) {
		t.Fatal("draft/handle/handler compatibility", nw)
	}
	nc := state(t, n, h)
	if nc.Generation != oldState.Generation+1 || nc.Selected != "child" || nc.Focused != "child" || !nc.Expanded["loaded"] || nc.Request != nil || nc.Status["lazy"].Phase != Unloaded || nc.Status["error"].Phase != LoadError || nc.Status["canceled"].Phase != Canceled {
		t.Fatal(nc)
	}
	nvh, _ := n.Viewport(w.InstancePath)
	if nvh != vh || n.Snapshot().Viewports[w.InstancePath] != (ViewportState{5, 7}) {
		t.Fatal("named viewport state not preserved")
	}
	if err = n.CompleteLoad(oldRequest, CollectionData{}, nil); err == nil {
		t.Fatal("successor accepted predecessor request")
	}
	requireOK(t, n.CheckStateWith(passGate))
	requireOK(t, n.ExpandItem(target(t, n, h, "lazy")))
	if state(t, n, h).Request.RequestID <= oldRequest.RequestID {
		t.Fatal("request watermark reset")
	}
	// Failed native preparation/disposal does not close or cancel predecessor.
	n.Close()
	unchanged(t, s, before, s.Revision, s.BatchRevision)
	requireOK(t, s.CompleteLoad(oldRequest, CollectionData{[]CollectionItem{row("new", "lazy")}}, nil))
	if calls != 1 {
		t.Fatal("preparing or loading executed callback")
	}
}

func TestSequenceAndGenerationAcrossMultiplePublishedSuccessors(t *testing.T) {
	s, h := collectionSession(t, "list", CollectionData{[]CollectionItem{row("x", "")}}, true, false)
	lastSeq := uint64(0)
	lastGeneration := state(t, s, h).Generation
	for i := 0; i < 3; i++ {
		oldEvent := eventFor(s, target(t, s, h, "x"), Activate)
		requireOK(t, s.Bind(h, func(e Event) ([]Update, error) {
			if e.Sequence <= lastSeq {
				t.Error("application sequence reset")
			}
			lastSeq = e.Sequence
			return nil, nil
		}))
		requireOK(t, s.Dispatch(oldEvent))
		n, err := s.Successor(collectionRoot(t, "list"), providersFor(s))
		requireOK(t, err)
		code(t, n.Dispatch(oldEvent), "stale-event")
		updated := oldEvent
		updated.ModelRevision = n.Revision
		ct := target(t, n, h, "x")
		updated.Collection = &ct
		code(t, n.Dispatch(updated), "duplicate-event")
		if state(t, n, h).Generation <= lastGeneration {
			t.Fatal("collection generation reset")
		}
		lastGeneration = state(t, n, h).Generation
		s.Close()
		s = n
	}
	// Remove and recreate an identically named widget: old handle never revives.
	old := h
	empty := root(t, `page=[];`)
	requireOK(t, s.Reload(empty))
	r := collectionRoot(t, "list")
	path := ""
	r.Walk(func(x *parser.Instance) {
		if x.Widget == "list" {
			path = x.Path
		}
	})
	n, err := s.Successor(r, map[string]CollectionProvider{path: {ID: "fixture", Epoch: 1, RootLoaded: true, Initial: CollectionData{[]CollectionItem{row("x", "")}}}})
	requireOK(t, err)
	nw, _ := n.Widget("page/items")
	if nw.Handle.Generation <= old.Generation {
		t.Fatal("widget generation reused")
	}
	code(t, n.ValidateCollectionTarget(CollectionTarget{Handle: old, ModelRevision: n.Revision, CollectionGeneration: 1, ItemID: "x"}), "stale-handle")
}

func TestSuccessorProviderChangeAndFailedAdmissionRetainOldBundle(t *testing.T) {
	s, h := collectionSession(t, "list", CollectionData{}, false, true)
	req, err := s.BeginLoad(target(t, s, h, ""))
	requireOK(t, err)
	before := s.Snapshot()
	providers := providersFor(s)
	for path, p := range providers {
		p.Initial = CollectionData{[]CollectionItem{row("invalid", "missing")}}
		p.RootLoaded = true
		providers[path] = p
	}
	if _, err = s.Successor(collectionRoot(t, "list"), providers); err == nil {
		t.Fatal("invalid detached provider admitted")
	}
	unchanged(t, s, before, s.Revision, s.BatchRevision)
	providers = providersFor(s)
	for path, p := range providers {
		p.Epoch++
		p.Initial = CollectionData{[]CollectionItem{row("different", "")}}
		p.RootLoaded = true
		providers[path] = p
	}
	n, err := s.Successor(collectionRoot(t, "list"), providers)
	requireOK(t, err)
	c := state(t, n, h)
	if c.Request != nil || c.Data.Items[0].ID != "different" || c.AutoLoadPending {
		t.Fatal(c)
	}
	unchanged(t, s, before, s.Revision, s.BatchRevision)
	requireOK(t, s.CompleteLoad(req, CollectionData{}, nil))
}

func TestRootRecoverySurvivesCompatibleReloadWithoutAutomaticRestart(t *testing.T) {
	for _, mode := range []string{"fresh", "canceled", "error", "loading"} {
		t.Run(mode, func(t *testing.T) {
			s, h := collectionSession(t, "list", CollectionData{}, false, true)
			if mode != "fresh" {
				r, err := s.BeginLoad(target(t, s, h, ""))
				requireOK(t, err)
				if mode == "canceled" {
					requireOK(t, s.CancelLoad(h))
				}
				if mode == "error" {
					requireOK(t, s.CompleteLoad(r, CollectionData{}, errors.New("failure")))
				}
			}
			n, err := s.Successor(collectionRoot(t, "list"), providersFor(s))
			requireOK(t, err)
			c := state(t, n, h)
			if c.AutoLoadPending != (mode == "fresh") || c.Request != nil {
				t.Fatal(c)
			}
			if mode == "error" && c.Status[""].Phase != LoadError || mode == "canceled" && c.Status[""].Phase != Canceled || mode == "loading" && c.Status[""].Phase != Unloaded {
				t.Fatal(c)
			}
			if mode == "canceled" || mode == "error" || mode == "loading" {
				requireOK(t, n.RetryItem(target(t, n, h, "")))
				if state(t, n, h).Request == nil {
					t.Fatal("lost recovery")
				}
			}
		})
	}
}

func TestViewportIdentityRemovalKindChangeAndAxisPreservation(t *testing.T) {
	r := collectionRoot(t, "list")
	s, h := collectionSession(t, "list", CollectionData{[]CollectionItem{row("x", "")}}, true, false)
	requireOK(t, s.CheckStateWith(passGate))
	w, _ := s.Widget(h.Path)
	vh, _ := s.Viewport(w.InstancePath)
	requireOK(t, s.SetViewport(vh, s.Revision, ViewportState{8, 9}))
	r.Walk(func(n *parser.Instance) {
		if n.Path == w.InstancePath {
			n.Layout["overflow-x"] = "clip"
		}
	})
	n, err := s.Successor(r, providersFor(s))
	requireOK(t, err)
	if v := n.Snapshot().Viewports[w.InstancePath]; v.X != 0 || v.Y != 9 {
		t.Fatal(v)
	}
	requireOK(t, n.CheckStateWith(passGate))
	code(t, n.SetViewport(vh, s.Revision, ViewportState{}), "stale-event")
	oldRevision := n.Revision
	requireOK(t, n.Reload(root(t, `page=[];`)))
	n2, err := n.Successor(collectionRoot(t, "list"), providersFor(s))
	requireOK(t, err)
	requireOK(t, n2.CheckStateWith(passGate))
	code(t, n2.SetViewport(vh, n2.Revision, ViewportState{}), "stale-handle")
	if n2.Revision <= oldRevision {
		t.Fatal("model revision reset")
	}
	// Same normalized path with different owner kind receives a new identity.
	changed := collectionRoot(t, "tree")
	n3, err := s.Successor(changed, providersFor(s))
	requireOK(t, err)
	requireOK(t, n3.CheckStateWith(passGate))
	code(t, n3.SetViewport(vh, n3.Revision, ViewportState{}), "stale-handle")
	if v := n3.Snapshot().Viewports[w.InstancePath]; v != (ViewportState{}) {
		t.Fatal("changed kind preserved offsets")
	}
}

func TestSnapshotProvenanceProfilesAndStateRevision(t *testing.T) {
	r := collectionRoot(t, "list")
	r.Uses = []parser.UseSite{{Definition: "original", Span: parser.Span{Line: 2}}}
	s, err := New("s", r)
	requireOK(t, err)
	r.Uses[0].Definition = "caller mutation"
	snapshot := s.Snapshot()
	snapshot.Root.Uses[0].Definition = "snapshot mutation"
	snapshot.Root.Profile = "bad"
	if got := s.SnapshotRoot(); got.Uses[0].Definition != "original" || got.Profile != "sdui/0.3" {
		t.Fatal(got)
	}
	for _, kind := range []string{"unknown", "mixed", "cycle", "nil child"} {
		r := collectionRoot(t, "list")
		switch kind {
		case "unknown":
			r.Profile = "sdui/0.4"
		case "mixed":
			r.Rows[0][0].Profile = ""
		case "cycle":
			r.Rows[0] = append(r.Rows[0], r)
		case "nil child":
			r.Rows[0] = append(r.Rows[0], nil)
		}
		if _, err := New("bad", r); err == nil {
			t.Fatal("accepted invalid tree", kind)
		}
	}
	legacy := session(t)
	edit, _ := legacy.Widget("page/edit")
	mutations := []func() error{func() error { return legacy.Draft(edit.Handle, "draft") }, func() error { return legacy.Revert(edit.Handle) }, func() error { return legacy.Focus(edit.Handle) }, func() error { return legacy.InvalidateEvents() }, func() error {
		return legacy.Apply(legacy.Revision, legacy.BatchRevision+1, []Update{{Handle: edit.Handle, Property: Label, Value: Text("new")}})
	}}
	for _, f := range mutations {
		old := legacy.StateRevision
		requireOK(t, f())
		if legacy.StateRevision <= old {
			t.Fatal("state revision did not advance")
		}
	}
	before := legacy.Snapshot()
	if err = legacy.Draft(edit.Handle, stringsOfLength(32769)); err == nil {
		t.Fatal("invalid draft accepted")
	}
	unchanged(t, legacy, before, legacy.Revision, legacy.BatchRevision)
	clone := legacy.SnapshotRoot()
	if clone.Profile != "" {
		t.Fatal("legacy profile encoding changed")
	}
	next, err := legacy.Successor(clone, nil)
	requireOK(t, err)
	if !reflect.DeepEqual(legacy.SnapshotRoot(), next.SnapshotRoot()) {
		t.Fatal("legacy successor changed state")
	}
}
func stringsOfLength(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = 'x'
	}
	return string(b)
}

func TestExplicitRootLoadAfterSuccessorAllocatesOnceWithoutDomainCall(t *testing.T) {
	s, h := collectionSession(t, "list", CollectionData{}, false, true)
	old, err := s.BeginLoad(target(t, s, h, ""))
	requireOK(t, err)
	n, err := s.Successor(collectionRoot(t, "list"), providersFor(s))
	requireOK(t, err)
	c := state(t, n, h)
	if c.RootLoaded || c.AutoLoadPending || c.Request != nil || c.Status[""].Phase != Unloaded {
		t.Fatal("root did not enter paused unloaded state", c)
	}
	calls := 0
	requireOK(t, n.Bind(h, func(Event) ([]Update, error) { calls++; return nil, nil }))
	retry := eventFor(n, target(t, n, h, ""), Retry)
	requireOK(t, n.Dispatch(retry))
	current := *state(t, n, h).Request
	if current.RequestID != old.RequestID+1 {
		t.Fatal("explicit Load did not allocate exactly one new request")
	}
	requireOK(t, n.Dispatch(eventFor(n, target(t, n, h, ""), Retry)))
	if *state(t, n, h).Request != current || n.Sequence() != 2 || calls != 0 {
		t.Fatal("repeated Load started another request or invoked domain")
	}
	code(t, n.Dispatch(retry), "duplicate-event")
	before := n.Snapshot()
	if err = n.CompleteLoad(old, CollectionData{[]CollectionItem{row("late", "")}}, nil); err == nil {
		t.Fatal("predecessor completion accepted")
	}
	unchanged(t, n, before, n.Revision, n.BatchRevision)
	requireOK(t, n.CompleteLoad(current, CollectionData{[]CollectionItem{row("fresh", "")}}, nil))
	if c = state(t, n, h); !c.RootLoaded || c.Data.Items[0].ID != "fresh" || calls != 0 {
		t.Fatal(c, calls)
	}
}
