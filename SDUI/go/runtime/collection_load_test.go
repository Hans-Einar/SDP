package runtime

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestLoadCancellationSupersessionAndFailureRecovery(t *testing.T) {
	s, h := collectionSession(t, "tree", CollectionData{[]CollectionItem{branch("a", "", false), branch("b", "", false)}}, true, true)
	gen := state(t, s, h).Generation
	requireOK(t, s.ExpandItem(target(t, s, h, "a")))
	first := *state(t, s, h).Request
	if first.RequestID != 1 || state(t, s, h).Status["a"].Phase != Loading || state(t, s, h).Generation != gen {
		t.Fatal(first, state(t, s, h))
	}
	requireOK(t, s.ExpandItem(target(t, s, h, "b")))
	second := *state(t, s, h).Request
	if second.RequestID <= first.RequestID || state(t, s, h).Status["a"].Phase != Canceled {
		t.Fatal("superseding did not revoke")
	}
	before := s.Snapshot()
	if e := s.CompleteLoad(first, CollectionData{[]CollectionItem{row("late", "a")}}, nil); e == nil {
		t.Fatal("old success accepted")
	}
	if e := s.CompleteLoad(first, CollectionData{}, errors.New("late")); e == nil {
		t.Fatal("old failure accepted")
	}
	unchanged(t, s, before, s.Revision, s.BatchRevision)
	requireOK(t, s.CompleteLoad(second, CollectionData{}, errors.New(strings.Repeat("😀", 1400))))
	c := state(t, s, h)
	if c.Request != nil || c.Status["b"].Phase != LoadError || len(c.Status["b"].Error) > 4096 || !utf8.ValidString(c.Status["b"].Error) || c.Generation != gen {
		t.Fatal(c)
	}
	requireOK(t, s.CollapseItem(target(t, s, h, "b")))
	requireOK(t, s.ExpandItem(target(t, s, h, "b")))
	if state(t, s, h).Request != nil || state(t, s, h).Status["b"].Phase != LoadError {
		t.Fatal("error expansion retried")
	}
	requireOK(t, s.RetryItem(target(t, s, h, "b")))
	third := *state(t, s, h).Request
	if third.RequestID <= second.RequestID {
		t.Fatal(third)
	}
	requireOK(t, s.RetryItem(target(t, s, h, "b")))
	if *state(t, s, h).Request != third {
		t.Fatal("duplicate Retry started a second load")
	}
	requireOK(t, s.CompleteLoad(third, CollectionData{}, nil))
	c = state(t, s, h)
	if c.Generation != gen+1 || c.Request != nil || c.Status["b"].Phase != Loaded {
		t.Fatal(c)
	}
	requireOK(t, s.CollapseItem(target(t, s, h, "b")))
	requireOK(t, s.ExpandItem(target(t, s, h, "b")))
	if state(t, s, h).Request != nil {
		t.Fatal("empty success automatically reloaded")
	}
	requireOK(t, s.CollapseItem(target(t, s, h, "a")))
	requireOK(t, s.ExpandItem(target(t, s, h, "a")))
	fourth := *state(t, s, h).Request
	if fourth.RequestID <= third.RequestID {
		t.Fatal("canceled branch did not get fresh request")
	}
	requireOK(t, s.CancelLoad(h))
	if state(t, s, h).Status["a"].Phase != Canceled {
		t.Fatal("cancel not visible")
	}
}

func TestRootAutoLoadAndExplicitRestart(t *testing.T) {
	s, h := collectionSession(t, "list", CollectionData{}, false, true)
	if !state(t, s, h).AutoLoadPending {
		t.Fatal("fresh unloaded root not pending")
	}
	root := target(t, s, h, "")
	first, err := s.BeginLoad(root)
	requireOK(t, err)
	if state(t, s, h).AutoLoadPending {
		t.Fatal("first load did not consume auto flag")
	}
	requireOK(t, s.CancelLoad(h))
	requireOK(t, s.Focus(h))
	if state(t, s, h).Request != nil || state(t, s, h).AutoLoadPending {
		t.Fatal("focus restarted canceled root")
	}
	requireOK(t, s.RetryItem(root))
	second := *state(t, s, h).Request
	if second.RequestID <= first.RequestID {
		t.Fatal("root retry reused token")
	}
	before := s.Snapshot()
	if e := s.CompleteLoad(first, CollectionData{}, nil); e == nil {
		t.Fatal("canceled completion accepted")
	}
	unchanged(t, s, before, s.Revision, s.BatchRevision)
	requireOK(t, s.CompleteLoad(second, CollectionData{}, nil))
	c := state(t, s, h)
	if !c.RootLoaded || c.AutoLoadPending || c.Request != nil || c.Status[""].Phase != Loaded {
		t.Fatal(c)
	}
	if e := s.RetryItem(target(t, s, h, "")); e == nil {
		t.Fatal("loaded root retried without refresh")
	}
}

func TestCompletionRejectsInvalidDataAndLayoutWithoutPartialPublication(t *testing.T) {
	s, h := collectionSession(t, "tree", CollectionData{[]CollectionItem{branch("a", "", false), row("keep", "")}}, true, true)
	requireOK(t, s.SelectItem(target(t, s, h, "keep")))
	requireOK(t, s.CheckStateWith(func(v Snapshot) (map[string]ViewportState, error) {
		for _, c := range v.Collections {
			for _, item := range c.Data.Items {
				if item.ID == "too-wide" {
					return nil, errors.New("content overflow")
				}
			}
		}
		return v.Viewports, nil
	}))
	for _, d := range []CollectionData{{[]CollectionItem{row("bad", "wrong-parent")}}, {[]CollectionItem{row("too-wide", "a")}}} {
		if state(t, s, h).Status["a"].Phase == LoadError {
			requireOK(t, s.RetryItem(target(t, s, h, "a")))
		} else {
			requireOK(t, s.ExpandItem(target(t, s, h, "a")))
		}
		req := *state(t, s, h).Request
		before := state(t, s, h)
		offsets := s.Snapshot().Viewports
		requireOK(t, s.CompleteLoad(req, d, nil))
		now := state(t, s, h)
		if !reflect.DeepEqual(before.Data, now.Data) || now.Generation != before.Generation || now.Selected != "keep" || now.Request != nil || now.Status["a"].Phase != LoadError || !reflect.DeepEqual(offsets, s.Snapshot().Viewports) {
			t.Fatal("failure partially published", now)
		}
	}
	requireOK(t, s.RetryItem(target(t, s, h, "a")))
	req := *state(t, s, h).Request
	requireOK(t, s.CompleteLoad(req, CollectionData{[]CollectionItem{row("child", "a")}}, nil))
	c := state(t, s, h)
	if c.Status["a"].Phase != Loaded || c.Selected != "keep" || len(c.Data.Items) != 3 {
		t.Fatal(c)
	}
}

func TestLoadRevocationForAncestorHideRefreshProviderAndClose(t *testing.T) {
	for _, mode := range []string{"ancestor-collapse", "hidden", "disabled", "replace", "provider", "reload", "close"} {
		t.Run(mode, func(t *testing.T) {
			s, h := collectionSession(t, "tree", CollectionData{[]CollectionItem{branch("outer", "", true), branch("inner", "outer", false)}}, true, true)
			requireOK(t, s.ExpandItem(target(t, s, h, "outer")))
			requireOK(t, s.ExpandItem(target(t, s, h, "inner")))
			req := *state(t, s, h).Request
			switch mode {
			case "ancestor-collapse":
				requireOK(t, s.CollapseItem(target(t, s, h, "outer")))
			case "hidden", "disabled":
				p := Visible
				if mode == "disabled" {
					p = Enabled
				}
				requireOK(t, s.Apply(s.Revision, s.BatchRevision+1, []Update{{Handle: h, Property: p, Value: Bool(false)}}))
				requireOK(t, s.Apply(s.Revision, s.BatchRevision+1, []Update{{Handle: h, Property: p, Value: Bool(true)}}))
			case "replace":
				requireOK(t, s.ReplaceCollection(target(t, s, h, ""), state(t, s, h).Data))
			case "provider":
				p, _ := s.Provider(h)
				p.Epoch++
				w, _ := s.Widget(h.Path)
				requireOK(t, s.BindProviders(map[string]CollectionProvider{w.InstancePath: p}))
			case "reload":
				requireOK(t, s.Reload(collectionRoot(t, "tree")))
			case "close":
				s.Close()
			}
			before := s.Snapshot()
			if e := s.CompleteLoad(req, CollectionData{[]CollectionItem{row("late", "inner")}}, nil); e == nil {
				t.Fatal("revoked request accepted")
			}
			if e := s.CompleteLoad(req, CollectionData{}, errors.New("late error")); e == nil {
				t.Fatal("revoked error accepted")
			}
			unchanged(t, s, before, s.Revision, s.BatchRevision)
		})
	}
}

func TestLoaderBarrierDeliveryCannotPublishAfterCancellation(t *testing.T) {
	s, h := collectionSession(t, "list", CollectionData{}, false, true)
	entered := make(chan LoadRequest)
	release := make(chan struct{})
	result := make(chan CollectionData)
	p, _ := s.Provider(h)
	p.Epoch++
	p.Load = func(_ context.Context, r LoadRequest) (CollectionData, error) {
		entered <- r
		<-release
		return CollectionData{[]CollectionItem{row("late", "")}}, nil
	}
	w, _ := s.Widget(h.Path)
	requireOK(t, s.BindProviders(map[string]CollectionProvider{w.InstancePath: p}))
	req, err := s.BeginLoad(target(t, s, h, ""))
	requireOK(t, err)
	// Host simulation owns goroutine and delivery; provider cannot touch Session.
	go func() { d, _ := p.Load(context.Background(), req); result <- d }()
	if got := <-entered; got != req {
		t.Fatal(got)
	}
	requireOK(t, s.CancelLoad(h))
	before := s.Snapshot()
	close(release)
	d := <-result
	if err = s.CompleteLoad(req, d, nil); err == nil {
		t.Fatal("provider ignoring cancellation published")
	}
	unchanged(t, s, before, s.Revision, s.BatchRevision)
}
