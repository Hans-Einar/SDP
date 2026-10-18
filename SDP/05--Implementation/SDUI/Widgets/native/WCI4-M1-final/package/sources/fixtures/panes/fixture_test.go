package panes

import (
	"context"
	"reflect"
	"testing"
	"time"

	"fyne.io/fyne/v2/test"
	"github.com/Hans-Einar/SDP/SDUI/go/host/fynehost"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

func newHost(t *testing.T) (*Fixture, *Controller) {
	t.Helper()
	a := test.NewTempApp(t)
	w := a.NewWindow("panes fixture")
	t.Cleanup(w.Close)
	f, err := New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(f.Close)
	h := fynehost.NewDocumentHost(w.Canvas(), layout.Size{W: 1000, H: 650})
	t.Cleanup(h.Close)
	c := &Controller{Fixture: f, Host: h, Source: Source, Sequence: 1}
	f.Conflict = func() error {
		b := h.Current()
		receiver, _ := b.Session.Widget(PreviewPath)
		return b.Session.Draft(receiver.Handle, "newer receiver draft")
	}
	r, err := f.Request(Source, 1)
	if err != nil {
		t.Fatal(err)
	}
	if err = h.Adopt(r); err != nil {
		t.Fatal(err)
	}
	return f, c
}
func activate(t *testing.T, c *Controller, id string) (ui.InteractionResult, error) {
	t.Helper()
	var result ui.InteractionResult
	err := c.Host.Mutate(func(s *ui.Session) error {
		h, ok := s.Pane(TabsPath)
		if !ok {
			t.Fatal("missing tabs")
		}
		tabs, _ := s.Tabs(h)
		for _, p := range tabs.Pages {
			if p.ID == id {
				var err error
				result, err = s.DispatchInteraction(ui.Event{Kind: ui.ActivatePage, Handle: h, ModelRevision: s.Revision, StateRevision: s.StateRevision, Sequence: s.Sequence() + 1, Page: &ui.PageActivation{PreviousID: tabs.Selected, PageID: id, Page: p.Handle}})
				return err
			}
		}
		t.Fatalf("missing page %s", id)
		return nil
	})
	return result, err
}
func TestConnectedPanesPreserveDraftScrollSplitAndReload(t *testing.T) {
	f, c := newHost(t)
	if f.Calls() != 0 {
		t.Fatal("preparation invoked SDL")
	}
	s := c.Host.Current().Session
	if err := c.Host.Mutate(func(s *ui.Session) error {
		edit, _ := s.Widget(TabsPath + "/overview/edit")
		if err := s.Draft(edit.Handle, "unsaved overview"); err != nil {
			return err
		}
		return s.SetViewports(map[string]ui.ViewportState{TreePath: {X: 20, Y: 120}})
	}); err != nil {
		t.Fatal(err)
	}
	offset := s.Snapshot().Viewports[TreePath]
	if offset.Y == 0 {
		t.Fatal("fixture tree did not scroll")
	}
	if _, err := activate(t, c, "notes"); err != nil {
		t.Fatal(err)
	}
	if s.Snapshot().Viewports[TreePath] != offset {
		t.Fatal("inactive scroll lost")
	}
	if err := c.Run("select overview"); err != nil {
		t.Fatal(err)
	}
	if f.Calls() != 1 {
		t.Fatal("programmatic selection called SDL")
	}
	edit, _ := s.Widget(TabsPath + "/overview/edit")
	if edit.Draft != "unsaved overview" || !edit.Dirty || s.Snapshot().Viewports[TreePath] != offset {
		t.Fatal("page state lost", edit)
	}
	if err := c.Run("collapse first"); err != nil {
		t.Fatal(err)
	}
	before := s.Snapshot().Splits[SplitPath]
	if before.Collapsed != ui.SplitFirst || before.MinFirst <= 0 || before.MinSecond <= 0 {
		t.Fatal(before)
	}
	if err := c.Run("reload"); err != nil {
		t.Fatal(err)
	}
	if !s.Closed() {
		t.Fatal("predecessor not revoked")
	}
	s = c.Host.Current().Session
	after := s.Snapshot().Splits[SplitPath]
	if after.Collapsed != ui.SplitFirst || after.SavedProportion != before.SavedProportion {
		t.Fatal("reload split state lost", after)
	}
	if err := c.Run("restore"); err != nil {
		t.Fatal(err)
	}
	edit, _ = s.Widget(TabsPath + "/overview/edit")
	if edit.Draft != "unsaved overview" || s.Snapshot().Viewports[TreePath] != offset || f.Calls() != 1 {
		t.Fatal("reload lost draft/scroll or invoked SDL")
	}
	if _, err := activate(t, c, "notes"); err != nil || f.Calls() != 2 {
		t.Fatal("postreload action", err, f.Calls())
	}
	preview, _ := s.Widget(PreviewPath)
	if preview.Value != "overview -> notes" {
		t.Fatal(preview)
	}
}
func TestFixtureActionAndFinalPreparationFailures(t *testing.T) {
	for _, mode := range []string{"error", "invalid", "draft-conflict", "resource"} {
		t.Run(mode, func(t *testing.T) {
			f, c := newHost(t)
			old := c.Host.Current()
			before := old.Geometry()
			if mode == "resource" {
				if err := c.Run("resource-error"); err != nil {
					t.Fatal(err)
				}
			} else if err := c.Run("action " + mode); err != nil {
				t.Fatal(err)
			}
			result, err := activate(t, c, "notes")
			if err == nil || f.Calls() != 1 {
				t.Fatal(result, err, f.Calls())
			}
			expected := ui.DomainUnknown
			if mode == "resource" || mode == "draft-conflict" {
				expected = ui.DomainSucceeded
			}
			if result.Domain != expected {
				t.Fatal("lost outcome", result)
			}
			if c.Host.Current() != old || old.Session.Snapshot().Tabs[TabsPath].Selected != "overview" {
				t.Fatal("failed selection published")
			}
			preview, _ := old.Session.Widget(PreviewPath)
			if preview.Value != "No page activation" {
				t.Fatal("failed result published")
			}
			if mode == "draft-conflict" {
				if preview.Draft != "newer receiver draft" {
					t.Fatal("reentrant accepted draft rolled back")
				}
			} else if old.Geometry() != before {
				t.Fatal("speculative presentation promoted")
			}
		})
	}
}
func TestFixtureFailureCommandsPreserveLiveBundle(t *testing.T) {
	f, c := newHost(t)
	for _, mode := range []string{"profile", "binding", "layout", "resource", "guard", "stale"} {
		t.Run(mode, func(t *testing.T) {
			old := c.Host.Current()
			if err := c.Run("fail " + mode); err == nil {
				t.Fatal("failure not injected")
			}
			if c.Host.Current() != old || old.Session.Closed() || f.Calls() != 0 {
				t.Fatal("failed candidate replaced live bundle or invoked SDL")
			}
		})
	}
	if err := c.Run("reload"); err != nil {
		t.Fatal("recovery reload", err)
	}
	before := c.Host.Current().Session.Snapshot()
	seq := c.Sequence
	for _, cmd := range []string{"activate notes", "select missing", "ratio NaN", "ratio -1", "collapse none", "hide missing", "complete unknown bogus", "resize Inf 2", "resource-error extra", "fail unknown", "action unknown"} {
		if err := c.Run(cmd); err == nil {
			t.Fatalf("accepted invalid command %q", cmd)
		}
	}
	if !reflect.DeepEqual(before, c.Host.Current().Session.Snapshot()) || seq != c.Sequence || f.Calls() != 0 {
		t.Fatal("invalid controls mutated state")
	}
}
func TestProviderBarriersLateReplyAndClose(t *testing.T) {
	f, err := New()
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	started := make(chan string, 2)
	f.Log = func(kind string, value any) {
		if kind == "load" {
			started <- value.(map[string]any)["key"].(string)
		}
	}
	provider := f.Providers()[TreePath]
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan loadReply, 1)
	request := ui.LoadRequest{Target: ui.CollectionTarget{Handle: ui.Handle{Path: TreePath}, ModelRevision: 1, ItemID: "lazy"}, RequestID: 1}
	go func() { data, err := provider.Load(ctx, request); done <- loadReply{data, err} }()
	var key string
	select {
	case key = <-started:
	case <-time.After(5 * time.Second):
		t.Fatal("provider did not start")
	}
	cancel()
	if err = f.Complete(key, "bogus"); err == nil || len(f.Pending()) != 1 {
		t.Fatal("invalid completion consumed barrier")
	}
	if err = f.Complete(key, "success"); err != nil {
		t.Fatal(err)
	}
	select {
	case reply := <-done:
		if reply.err != nil || len(reply.data.Items) != 24 || reply.data.Items[0].Parent != "lazy" {
			t.Fatal(reply)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("late result not delivered")
	}
	if err = f.Complete(key, "success"); err == nil {
		t.Fatal("barrier reused")
	}
	request.RequestID++
	go func() { data, err := provider.Load(context.Background(), request); done <- loadReply{data, err} }()
	select {
	case <-started:
	case <-time.After(5 * time.Second):
		t.Fatal("second provider did not start")
	}
	f.Close()
	select {
	case reply := <-done:
		if reply.err == nil {
			t.Fatal("close did not cancel barrier")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("close leaked provider")
	}
	if len(f.Pending()) != 0 {
		t.Fatal("pending barriers retained")
	}
}
