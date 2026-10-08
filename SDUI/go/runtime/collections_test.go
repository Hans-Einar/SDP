package runtime

import (
	"context"
	"errors"
	"fmt"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"math"
	"reflect"
	"strings"
	"testing"
)

func row(id, parent string) CollectionItem {
	return CollectionItem{ID: ItemID(id), Parent: ItemID(parent), Kind: Row, Label: "same label", ChildrenLoaded: true}
}
func branch(id, parent string, loaded bool) CollectionItem {
	x := row(id, parent)
	x.HasChildren = true
	x.ChildrenLoaded = loaded
	return x
}
func group(id, parent string, loaded bool) CollectionItem {
	x := branch(id, parent, loaded)
	x.Kind = Group
	return x
}
func collectionRoot(t *testing.T, kind string) *parser.Instance {
	t.Helper()
	_, roots, err := parser.Compile(`sdui 0.3; page=[<items=` + kind + `("Items"){overflow-x=scroll,overflow-y=scroll},edit=input("Preview",value="initial")>] {overflow-y=scroll};`)
	if err != nil {
		t.Fatal(err)
	}
	return roots["page"]
}
func collectionSession(t *testing.T, kind string, data CollectionData, loaded, lazy bool) (*Session, Handle) {
	t.Helper()
	s, err := New("collection-session", collectionRoot(t, kind))
	if err != nil {
		t.Fatal(err)
	}
	p := CollectionProvider{ID: "fixture", Epoch: 1, Initial: data, RootLoaded: loaded}
	if lazy {
		p.Load = func(context.Context, LoadRequest) (CollectionData, error) {
			t.Error("runtime invoked provider")
			return CollectionData{}, nil
		}
	}
	w, _ := s.Widget("page/items")
	if err = s.BindProviders(map[string]CollectionProvider{w.InstancePath: p}); err != nil {
		t.Fatal(err)
	}
	return s, w.Handle
}
func target(t *testing.T, s *Session, h Handle, id string) CollectionTarget {
	t.Helper()
	v, e := s.Target(h, ItemID(id))
	if e != nil {
		t.Fatal(e)
	}
	return v
}
func state(t *testing.T, s *Session, h Handle) CollectionState {
	t.Helper()
	v, ok := s.Collection(h)
	if !ok {
		t.Fatal("missing collection")
	}
	return v
}
func requireOK(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
func unchanged(t *testing.T, s *Session, before Snapshot, rev, batch uint64) {
	t.Helper()
	if !reflect.DeepEqual(before, s.Snapshot()) || s.Revision != rev || s.BatchRevision != batch {
		t.Fatal("failed operation changed live state")
	}
}
func eventFor(s *Session, t CollectionTarget, k EventKind) Event {
	return Event{Handle: t.Handle, ModelRevision: s.Revision, Sequence: s.Sequence() + 1, Kind: k, Collection: &t}
}

func TestCollectionDataValidationAndAtomicMixedBatch(t *testing.T) {
	valid := CollectionData{[]CollectionItem{branch("a", "", true), row("child", "a"), row("b", "")}}
	s, h := collectionSession(t, "tree", valid, true, true)
	requireOK(t, s.ExpandItem(target(t, s, h, "a")))
	requireOK(t, s.SelectItem(target(t, s, h, "child")))
	deep := []CollectionItem{}
	for i := 0; i < 65; i++ {
		p := ""
		if i > 0 {
			p = fmt.Sprint(i - 1)
		}
		deep = append(deep, branch(fmt.Sprint(i), p, true))
	}
	oversized := make([]CollectionItem, 4097)
	for i := range oversized {
		oversized[i] = row(fmt.Sprint(i), "")
	}
	cases := map[string]CollectionData{
		"duplicate":      {[]CollectionItem{row("x", ""), row("x", "")}},
		"empty id":       {[]CollectionItem{row("", "")}},
		"long id":        {[]CollectionItem{row(strings.Repeat("a", 1025), "")}},
		"bad utf8":       {[]CollectionItem{row("\xff", "")}},
		"control id":     {[]CollectionItem{row("bad\n", "")}},
		"missing parent": {[]CollectionItem{row("x", "missing")}},
		"cycle":          {[]CollectionItem{branch("a", "b", true), branch("b", "a", true)}},
		"deep":           {deep}, "many": {oversized},
		"long label":       {[]CollectionItem{{ID: "x", Kind: Row, Label: strings.Repeat("a", 32769), ChildrenLoaded: true}}},
		"newline label":    {[]CollectionItem{{ID: "x", Kind: Row, Label: "a\nb", ChildrenLoaded: true}}},
		"unicode line":     {[]CollectionItem{{ID: "x", Kind: Row, Label: "a\u2028b", ChildrenLoaded: true}}},
		"bad label utf8":   {[]CollectionItem{{ID: "x", Kind: Row, Label: "\xff", ChildrenLoaded: true}}},
		"empty label":      {[]CollectionItem{{ID: "x", Kind: Row, ChildrenLoaded: true}}},
		"unknown kind":     {[]CollectionItem{{ID: "x", Kind: "unknown", Label: "X", ChildrenLoaded: true}}},
		"unloaded leaf":    {[]CollectionItem{{ID: "x", Kind: Row, Label: "X"}}},
		"separator label":  {[]CollectionItem{{ID: "x", Kind: Separator, Label: "X", ChildrenLoaded: true}}},
		"separator branch": {[]CollectionItem{{ID: "x", Kind: Separator, HasChildren: true, ChildrenLoaded: true}}},
		"leaf parent":      {[]CollectionItem{row("a", ""), row("b", "a")}},
		"unloaded parent":  {[]CollectionItem{branch("a", "", false), row("b", "a")}},
		"mixed batch":      {[]CollectionItem{row("new", ""), row("broken", "missing")}},
	}
	for name, d := range cases {
		t.Run(name, func(t *testing.T) {
			before := s.Snapshot()
			rev, batch := s.Revision, s.BatchRevision
			if err := s.ReplaceCollection(target(t, s, h, ""), d); err == nil {
				t.Fatal("accepted invalid data")
			}
			unchanged(t, s, before, rev, batch)
		})
	}
	flat, fh := collectionSession(t, "list", CollectionData{[]CollectionItem{row("old", "")}}, true, false)
	if err := flat.ReplaceCollection(target(t, flat, fh, ""), valid); err == nil {
		t.Fatal("hierarchical list accepted")
	}
	requireOK(t, flat.ReplaceCollection(target(t, flat, fh, ""), CollectionData{[]CollectionItem{row("new", ""), {ID: "s", Kind: Separator, ChildrenLoaded: true}, {ID: "g", Kind: Group, Label: "Heading", ChildrenLoaded: true}}}))
}

func TestProviderAdmissionAllInstancesAndOwnedCopies(t *testing.T) {
	r := collectionRoot(t, "tree")
	s, e := New("s", r)
	requireOK(t, e)
	w, _ := s.Widget("page/items")
	original := CollectionData{[]CollectionItem{row("x", "")}}
	p := CollectionProvider{ID: "p", Epoch: 1, Initial: original, RootLoaded: true}
	for _, providers := range []map[string]CollectionProvider{nil, {"wrong": p}, {w.InstancePath: p, "unused": p}, {w.InstancePath: {ID: "", Epoch: 0, RootLoaded: true}}, {w.InstancePath: {ID: "p", Epoch: 1, RootLoaded: false}}, {w.InstancePath: {ID: "p", Epoch: 1, RootLoaded: true, Initial: CollectionData{[]CollectionItem{branch("x", "", false)}}}}} {
		before := s.Snapshot()
		if e := s.BindProviders(providers); e == nil {
			t.Fatal("accepted invalid provider map")
		}
		unchanged(t, s, before, s.Revision, s.BatchRevision)
	}
	requireOK(t, s.BindProviders(map[string]CollectionProvider{w.InstancePath: p}))
	original.Items[0].Label = "caller mutation"
	got := state(t, s, w.Handle)
	got.Data.Items[0].Label = "getter mutation"
	got.Expanded["x"] = true
	got.Status[""] = LoadStatus{Phase: LoadError}
	snap := s.Snapshot()
	snap.Collections[w.InstancePath].Data.Items[0].Label = "snapshot mutation"
	delete(snap.Collections, w.InstancePath)
	provider, _ := s.Provider(w.Handle)
	provider.Initial.Items[0].Label = "provider getter mutation"
	if c := state(t, s, w.Handle); c.Data.Items[0].Label != "same label" || c.Expanded["x"] || c.Status[""].Phase != Loaded {
		t.Fatal("aliased data", c)
	}
	// Rebinding with the same epoch must retain live state; provider identity is explicit.
	p.Initial = CollectionData{[]CollectionItem{row("new", "")}}
	requireOK(t, s.BindProviders(map[string]CollectionProvider{w.InstancePath: p}))
	if state(t, s, w.Handle).Data.Items[0].ID != "x" {
		t.Fatal("same binding reset")
	}
	p.Epoch++
	old := target(t, s, w.Handle, "x")
	requireOK(t, s.BindProviders(map[string]CollectionProvider{w.InstancePath: p}))
	code(t, s.ValidateCollectionTarget(old), "stale-collection")
}

func TestSubtreeReplacementIdentityAndSelection(t *testing.T) {
	s, h := collectionSession(t, "tree", CollectionData{[]CollectionItem{branch("a", "", true), row("a1", "a"), branch("b", "", true), row("b1", "b")}}, true, true)
	requireOK(t, s.ExpandItem(target(t, s, h, "a")))
	requireOK(t, s.ExpandItem(target(t, s, h, "b")))
	requireOK(t, s.SelectItem(target(t, s, h, "b1")))
	requireOK(t, s.FocusItem(target(t, s, h, "b1")))
	old := target(t, s, h, "a1")
	gen := state(t, s, h).Generation
	requireOK(t, s.ReplaceChildren(target(t, s, h, "a"), CollectionData{[]CollectionItem{row("a3", "a"), row("a2", "a")}}))
	c := state(t, s, h)
	if c.Generation != gen+1 || c.Selected != "b1" || c.Focused != "b1" || !c.Expanded["b"] {
		t.Fatal(c)
	}
	ids := []ItemID{}
	for _, x := range c.Data.Items {
		ids = append(ids, x.ID)
	}
	if !reflect.DeepEqual(ids, []ItemID{"a", "b", "b1", "a3", "a2"}) {
		t.Fatal(ids)
	}
	code(t, s.ValidateCollectionTarget(old), "stale-collection")
	before := s.Snapshot()
	if e := s.ReplaceChildren(target(t, s, h, "a"), CollectionData{[]CollectionItem{row("b1", "a")}}); e == nil {
		t.Fatal("colliding subtree accepted")
	}
	unchanged(t, s, before, s.Revision, s.BatchRevision)
	if e := s.ReplaceChildren(target(t, s, h, "a"), CollectionData{[]CollectionItem{row("bad", "b")}}); e == nil {
		t.Fatal("foreign subtree accepted")
	}
	unchanged(t, s, before, s.Revision, s.BatchRevision)
	requireOK(t, s.ReplaceChildren(target(t, s, h, "a"), CollectionData{[]CollectionItem{row("a1", "a")}}))
	code(t, s.ValidateCollectionTarget(old), "stale-collection")
	d := state(t, s, h).Data
	for i := range d.Items {
		if d.Items[i].ID == "b1" {
			d.Items[i].Kind = Group
		}
	}
	requireOK(t, s.ReplaceCollection(target(t, s, h, ""), d))
	if c := state(t, s, h); c.Selected != "" || c.Focused != "" {
		t.Fatal("row-to-group retained row state", c)
	}
	gen = state(t, s, h).Generation
	requireOK(t, s.ReplaceCollection(target(t, s, h, ""), state(t, s, h).Data))
	if state(t, s, h).Generation != gen+1 {
		t.Fatal("identical replacement did not invalidate")
	}
}

func TestPureGateAtomicDataAndClampedViewport(t *testing.T) {
	s, h := collectionSession(t, "list", CollectionData{[]CollectionItem{row("one", "")}}, true, false)
	w, _ := s.Widget(h.Path)
	requireOK(t, s.CheckStateWith(func(v Snapshot) (map[string]ViewportState, error) {
		c := v.Collections[w.InstancePath]
		if len(c.Data.Items) > 2 {
			return nil, errors.New("native measure fails")
		}
		offsets := v.Viewports
		x := offsets[w.InstancePath]
		x.X = math.Min(x.X, float64(len(c.Data.Items)*10))
		x.Y = math.Min(x.Y, 5)
		offsets[w.InstancePath] = x
		// Input ownership is independent of the live model.
		v.Root.Path = "not live"
		if len(c.Data.Items) > 0 {
			c.Data.Items[0].Label = "not live"
		}
		return offsets, nil
	}))
	vh, ok := s.Viewport(w.InstancePath)
	if !ok {
		t.Fatal("missing viewport identity")
	}
	requireOK(t, s.SetViewport(vh, s.Revision, ViewportState{X: 999, Y: 999}))
	if got := s.Snapshot().Viewports[w.InstancePath]; got != (ViewportState{10, 5}) {
		t.Fatal(got)
	}
	before := s.Snapshot()
	rev, batch := s.Revision, s.BatchRevision
	if err := s.ReplaceCollection(target(t, s, h, ""), CollectionData{[]CollectionItem{row("a", ""), row("b", ""), row("c", "")}}); err == nil {
		t.Fatal("gate failure ignored")
	}
	unchanged(t, s, before, rev, batch)
	for _, v := range []ViewportState{{X: math.NaN()}, {Y: math.Inf(1)}, {X: -1}} {
		if e := s.SetViewport(vh, s.Revision, v); e == nil {
			t.Fatal("bad offset accepted")
		}
		unchanged(t, s, before, rev, batch)
	}
	forged := vh
	forged.Generation++
	code(t, s.SetViewport(forged, s.Revision, ViewportState{}), "stale-handle")
	code(t, s.SetViewport(vh, s.Revision+1, ViewportState{}), "stale-event")
	requireOK(t, s.ReplaceCollection(target(t, s, h, ""), CollectionData{}))
	if got := s.Snapshot().Viewports[w.InstancePath]; got.X != 0 {
		t.Fatal("deletion did not clamp", got)
	}
}
