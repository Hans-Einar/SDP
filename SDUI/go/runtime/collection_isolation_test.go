package runtime

import (
	"context"
	"errors"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"reflect"
	"testing"
)

func TestReusedInstancesOwnProvidersRequestsStateAndProvenance(t *testing.T) {
	_, roots, err := parser.Compile(`sdui 0.3; part=<items=list("Same label")>; alias=part; page=[left=alias,right=alias];`)
	requireOK(t, err)
	s, err := New("reused", roots["page"])
	requireOK(t, err)
	p := CollectionProvider{ID: "shared source", Epoch: 1, RootLoaded: false, Load: func(context.Context, LoadRequest) (CollectionData, error) {
		t.Error("runtime executed loader")
		return CollectionData{}, nil
	}}
	providers := map[string]CollectionProvider{}
	for _, w := range s.Widgets() {
		providers[w.InstancePath] = p
	}
	requireOK(t, s.BindProviders(providers))
	left, _ := s.Widget("page/left/items")
	right, _ := s.Widget("page/right/items")
	a, err := s.BeginLoad(target(t, s, left.Handle, ""))
	requireOK(t, err)
	b, err := s.BeginLoad(target(t, s, right.Handle, ""))
	requireOK(t, err)
	if a.Target.Handle == b.Target.Handle || a.Target.Handle.Path == b.Target.Handle.Path {
		t.Fatal("reuse shared identity")
	}
	c := state(t, s, left.Handle)
	c.Request.RequestID++
	if state(t, s, left.Handle).Request.RequestID != a.RequestID {
		t.Fatal("request getter aliases live state")
	}
	requireOK(t, s.CompleteLoad(a, CollectionData{[]CollectionItem{row("same", "")}}, nil))
	requireOK(t, s.SelectItem(target(t, s, left.Handle, "same")))
	if state(t, s, right.Handle).Request == nil || state(t, s, right.Handle).Selected != "" || state(t, s, right.Handle).RootLoaded {
		t.Fatal("left mutation changed right")
	}
	requireOK(t, s.CompleteLoad(b, CollectionData{[]CollectionItem{row("same", "")}}, nil))
	if state(t, s, right.Handle).Selected != "" {
		t.Fatal("matching IDs coupled selection")
	}
	// Reused UseSite slices must remain independent at every cloned node.
	baseline := s.SnapshotRoot()
	mutated := s.SnapshotRoot()
	mutated.Walk(func(n *parser.Instance) {
		for i := range n.Uses {
			n.Uses[i].Definition = "mutated"
		}
	})
	if !reflect.DeepEqual(s.SnapshotRoot(), baseline) {
		t.Fatal("nested provenance aliased")
	}
}

func TestFailedDataAndLoadingGateRetainActiveRequest(t *testing.T) {
	s, h := collectionSession(t, "tree", CollectionData{[]CollectionItem{branch("a", "", false), branch("b", "", false)}}, true, true)
	requireOK(t, s.ExpandItem(target(t, s, h, "a")))
	request := *state(t, s, h).Request
	before := s.Snapshot()
	if e := s.ReplaceCollection(target(t, s, h, ""), CollectionData{[]CollectionItem{row("bad", "missing")}}); e == nil {
		t.Fatal("invalid data accepted")
	}
	unchanged(t, s, before, s.Revision, s.BatchRevision)
	requireOK(t, s.CheckStateWith(func(v Snapshot) (map[string]ViewportState, error) {
		for _, c := range v.Collections {
			if c.Request != nil && c.Request.Target.ItemID == "b" {
				return nil, errors.New("loading affordance cannot mount")
			}
		}
		return v.Viewports, nil
	}))
	before = s.Snapshot()
	if e := s.ExpandItem(target(t, s, h, "b")); e == nil {
		t.Fatal("loading geometry failed but accepted")
	}
	unchanged(t, s, before, s.Revision, s.BatchRevision)
	if *state(t, s, h).Request != request {
		t.Fatal("failed loading candidate canceled live request")
	}
	requireOK(t, s.CompleteLoad(request, CollectionData{}, nil))
}

func TestViewportGateRejectsInvalidEffectiveOffsetsAtomically(t *testing.T) {
	s, h := collectionSession(t, "list", CollectionData{}, true, false)
	w, _ := s.Widget(h.Path)
	before := s.Snapshot()
	if e := s.CheckStateWith(func(Snapshot) (map[string]ViewportState, error) {
		return map[string]ViewportState{"missing": {X: 1}}, nil
	}); e == nil {
		t.Fatal("unknown gate viewport accepted")
	}
	unchanged(t, s, before, s.Revision, s.BatchRevision)
	if e := s.CheckStateWith(func(Snapshot) (map[string]ViewportState, error) {
		return map[string]ViewportState{w.InstancePath: {Y: -1}}, nil
	}); e == nil {
		t.Fatal("invalid gate offset accepted")
	}
	unchanged(t, s, before, s.Revision, s.BatchRevision)
	if e := s.SetViewports(map[string]ViewportState{w.InstancePath: {Y: 1}}); e == nil {
		t.Fatal("unchecked offset accepted without gate")
	}
	unchanged(t, s, before, s.Revision, s.BatchRevision)
}
