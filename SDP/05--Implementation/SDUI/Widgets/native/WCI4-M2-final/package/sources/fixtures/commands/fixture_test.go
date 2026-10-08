package commands

import (
	"context"
	"reflect"
	"strings"
	"testing"
	"time"

	"fyne.io/fyne/v2/test"
	"github.com/Hans-Einar/SDP/SDUI/go/host/fynehost"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
)

func newHost(t *testing.T, nonmodal bool) (*Fixture, *Controller) {
	t.Helper()
	a := test.NewTempApp(t)
	w := a.NewWindow(WindowTitle)
	t.Cleanup(w.Close)
	f, err := New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(f.Close)
	h := fynehost.NewDocumentHost(w.Canvas(), layout.Size{W: 1100, H: 750})
	t.Cleanup(h.Close)
	c := &Controller{Fixture: f, Host: h, Source: FixtureSource(nonmodal), Sequence: 1}
	f.Conflict = func(path string) error {
		s := h.Current().Session
		w, ok := s.Widget(path)
		if !ok {
			t.Fatalf("missing input %s", path)
		}
		return s.Draft(w.Handle, "newer draft during SDL")
	}
	r, err := f.Request(c.Source, 1)
	if err != nil {
		t.Fatal(err)
	}
	if err = h.Adopt(r); err != nil {
		t.Fatal(err)
	}
	return f, c
}
func TestTrustedEnabledConditionsAndStrictNames(t *testing.T) {
	for _, nonmodal := range []bool{false, true} {
		t.Run(map[bool]string{false: "modal", true: "nonmodal"}[nonmodal], func(t *testing.T) {
			f, c := newHost(t, nonmodal)
			for _, name := range []string{"Run", "Flag"} {
				path := map[string]string{"Run": "page/run", "Flag": "page/flag"}[name]
				for _, op := range []string{"disable", "enable"} {
					if err := c.Run(op + " " + name); err != nil {
						t.Fatal(err)
					}
					s := c.Host.Current().Session
					h, _ := s.Command(path)
					state, _ := s.CommandState(h)
					if state.Enabled != (op == "enable") {
						t.Fatal(state)
					}
					if len(f.Calls()) != 0 {
						t.Fatal("trusted condition invoked SDL", f.Calls())
					}
				}
			}
			before := c.Host.Current().Session.Snapshot()
			seq := c.Sequence
			for _, cmd := range []string{"enable run", "disable flag", "enable Toggle", "enable page/run", "disable Run extra", "invoke Run", "Accept", "accept unknown", "action Save true", "resize NaN 2", "fail unknown", "complete unknown bogus"} {
				if err := c.Run(cmd); err == nil {
					t.Fatalf("accepted invalid/simulated command %q", cmd)
				}
			}
			if !reflect.DeepEqual(before, c.Host.Current().Session.Snapshot()) || seq != c.Sequence || len(f.Calls()) != 0 {
				t.Fatal("invalid control mutated state")
			}
		})
	}
}
func TestFailedReloadPreservesLiveBundleAndDraft(t *testing.T) {
	f, c := newHost(t, false)
	if err := c.Host.Mutate(func(s *ui.Session) error {
		w, _ := s.Widget("page/view/body/mainDraft")
		return s.Draft(w.Handle, "unsaved parent")
	}); err != nil {
		t.Fatal(err)
	}
	for _, mode := range []string{"profile", "binding", "resource", "guard", "stale"} {
		old := c.Host.Current()
		if err := c.Run("fail " + mode); err == nil {
			t.Fatalf("%s did not fail", mode)
		}
		if c.Host.Current() != old || old.Session.Closed() || len(f.Calls()) != 0 {
			t.Fatalf("%s changed bundle/invoked SDL", mode)
		}
		w, _ := old.Session.Widget("page/view/body/mainDraft")
		if w.Draft != "unsaved parent" || !w.Dirty {
			t.Fatal("failed candidate lost draft", w)
		}
	}
	old := c.Host.Current()
	if err := c.Run("reload"); err != nil {
		t.Fatal(err)
	}
	if !old.Session.Closed() || len(f.Calls()) != 0 {
		t.Fatal("successor failed to revoke/invoked SDL")
	}
	w, _ := c.Host.Current().Session.Widget("page/view/body/mainDraft")
	if w.Draft != "unsaved parent" || !w.Dirty {
		t.Fatal("compatible reload lost draft", w)
	}
}
func TestProviderBarriersDeliverLateOnceAndReleaseOnClose(t *testing.T) {
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
	p := f.Providers()[ItemsPath]
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan loadReply, 1)
	req := ui.LoadRequest{Target: ui.CollectionTarget{Handle: ui.Handle{Path: ItemsPath}, ModelRevision: 1, ItemID: "group"}, RequestID: 1}
	start := func() { go func() { data, err := p.Load(ctx, req); done <- loadReply{data, err} }() }
	key := func() string {
		select {
		case k := <-started:
			return k
		case <-time.After(5 * time.Second):
			t.Fatal("provider did not start")
			return ""
		}
	}
	start()
	k := key()
	cancel()
	if err = f.Complete(k, "bogus"); err == nil || len(f.Pending()) != 1 {
		t.Fatal("invalid completion consumed barrier")
	}
	if err = f.Complete(k, "success"); err != nil {
		t.Fatal(err)
	}
	select {
	case r := <-done:
		if r.err != nil || len(r.data.Items) != 12 || r.data.Items[0].Parent != "group" {
			t.Fatal(r)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("late response missing")
	}
	if err = f.Complete(k, "success"); err == nil {
		t.Fatal("barrier replayed")
	}
	req.RequestID++
	start()
	key()
	f.Close()
	select {
	case r := <-done:
		if r.err == nil {
			t.Fatal("closed provider succeeded")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("close leaked waiter")
	}
	if len(f.Pending()) != 0 {
		t.Fatal("close retained pending")
	}
}

func TestDialogContextMenuActualRequestAdmission(t *testing.T) {
	for _, nonmodal := range []bool{false, true} {
		t.Run(map[bool]string{false: "modal", true: "nonmodal"}[nonmodal], func(t *testing.T) {
			f, c := newHost(t, nonmodal)
			s := c.Host.Current().Session
			identities, err := parser.ResolveInteractions(s.SnapshotRoot())
			if err != nil {
				t.Fatal(err)
			}
			menuPath := SettingsPath + "/editMenu"
			markPath := SettingsPath + "/mark"
			if _, ok := s.Menu(menuPath); !ok {
				t.Fatal("dialog context menu not admitted")
			}
			if id := identities[menuPath]; id.Dialog != SettingsPath || id.Target != NamePath {
				t.Fatal("wrong context target/surface", id)
			}
			if id := identities[menuPath+"/markItem"]; id.Command != markPath || id.Dialog != SettingsPath {
				t.Fatal("wrong local command", id)
			}
			h, ok := s.Command(markPath)
			if !ok {
				t.Fatal("missing local toggle")
			}
			mark, _ := s.CommandState(h)
			if !mark.Toggle || mark.Checked || mark.Binding.Module != "" || mark.Effect != "" {
				t.Fatal("Mark must remain a local unchecked toggle", mark)
			}
			if len(f.Calls()) != 0 {
				t.Fatal("admission invoked SDL", f.Calls())
			}
		})
	}
}

func TestRootSiblingChildActualRequestAdmission(t *testing.T) {
	for _, nonmodal := range []bool{false, true} {
		t.Run(map[bool]string{false: "modal", true: "nonmodal"}[nonmodal], func(t *testing.T) {
			f, c := newHost(t, nonmodal)
			s := c.Host.Current().Session
			ids, err := parser.ResolveInteractions(s.SnapshotRoot())
			if err != nil {
				t.Fatal(err)
			}
			if id := ids[SettingsPath+"/openChild"]; id.Dialog != SettingsPath || id.Target != "page/x" {
				t.Fatal("child opener identity", id)
			}
			if id := ids[SettingsPath+"/editMenu/childItem"]; id.Dialog != SettingsPath || id.Command != SettingsPath+"/openChild" {
				t.Fatal("child menu identity", id)
			}
			h, ok := s.Surface("page/x")
			if !ok {
				t.Fatal("root sibling child not admitted")
			}
			child, ok := s.SurfaceState(h)
			if !ok || !child.Modal || child.Open || child.Label != "SDUI WCI2 Child" {
				t.Fatal("child must stay initially closed/modal", child)
			}
			fields, err := s.DialogFields(h)
			if err != nil || len(fields) != 1 || fields[0].InstancePath != "page/x/childInput" {
				t.Fatal("ordinary child input", fields, err)
			}
			if id := ids["page/x/childClose"]; id.Dialog != "page/x" {
				t.Fatal("Close ownership", id)
			}
			settingsControls := 0
			for _, w := range s.Widgets() {
				if strings.HasPrefix(w.InstancePath, SettingsPath+"/") && (w.Handle.Kind == "input" || w.Handle.Kind == "button") {
					settingsControls++
				}
			}
			if settingsControls != 5 {
				t.Fatal("settings focus-control inventory changed", settingsControls)
			}
			if len(f.Calls()) != 0 {
				t.Fatal("admission invoked SDL", f.Calls())
			}
		})
	}
}
