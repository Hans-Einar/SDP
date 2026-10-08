package text

import (
	"fyne.io/fyne/v2/test"
	"github.com/Hans-Einar/SDP/SDUI/go/host/fynehost"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"reflect"
	"testing"
)

func newHost(t *testing.T, nonmodal bool) (*Fixture, *Controller) {
	t.Helper()
	return newHostSource(t, FixtureSource(nonmodal))
}
func newHostSource(t *testing.T, source string) (*Fixture, *Controller) {
	t.Helper()
	a := test.NewTempApp(t)
	w := a.NewWindow(WindowTitle)
	t.Cleanup(w.Close)
	f, err := New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(f.Close)
	h := fynehost.NewDocumentHost(w.Canvas(), layout.Size{W: 1100, H: 850})
	t.Cleanup(h.Close)
	c := &Controller{Fixture: f, Host: h, Source: source, Sequence: 1}
	c.InstallHooks()
	r, err := f.Request(c.Source, 1)
	if err != nil {
		t.Fatal(err)
	}
	if err = h.Adopt(r); err != nil {
		t.Fatal(err)
	}
	return f, c
}
func zero(m map[string]int) bool {
	for _, n := range m {
		if n != 0 {
			return false
		}
	}
	return true
}
func edit(c *Controller, path, value string, commit bool) error {
	return c.Host.Mutate(func(s *ui.Session) error {
		w, _ := s.Widget(path)
		change, err := s.EditField(w.Handle, s.Revision, ui.Text(value))
		if err != nil {
			return err
		}
		if !commit {
			return nil
		}
		e, err := s.CaptureCommit(change.Field.Target)
		if err != nil {
			return err
		}
		return s.Dispatch(e)
	})
}
func open(t *testing.T, c *Controller) ui.SurfaceTarget {
	t.Helper()
	var target ui.SurfaceTarget
	if err := c.Host.Mutate(func(s *ui.Session) error {
		h, _ := s.Surface(FormPath)
		var err error
		target, err = s.OpenSurface(h, ui.ContextTarget{})
		return err
	}); err != nil {
		t.Fatal(err)
	}
	return target
}
func terminal(c *Controller, target ui.SurfaceTarget, kind ui.EventKind) error {
	return c.Host.Mutate(func(s *ui.Session) error {
		e, err := s.CaptureDialog(target, kind)
		if err != nil {
			return err
		}
		_, err = s.DispatchInteraction(e)
		return err
	})
}
func TestActualRequestPathsAndTrustedConditions(t *testing.T) {
	for _, nonmodal := range []bool{false, true} {
		t.Run(map[bool]string{false: "modal", true: "nonmodal"}[nonmodal], func(t *testing.T) {
			f, c := newHost(t, nonmodal)
			s := c.Host.Current().Session
			for _, path := range Targets {
				w, ok := s.Widget(path)
				if !ok {
					t.Fatal(path)
				}
				field, ok := s.Field(w.Handle)
				if !ok || field.Input == nil {
					t.Fatal("not extended", path)
				}
			}
			for _, path := range []string{"page/view/header/loadButton", "page/view/header/dialogButton", FormPath + "/formActions/formAccept", FormPath + "/formActions/formCancel", FormPath + "/formActions/formClose"} {
				if _, ok := s.Widget(path); !ok {
					t.Fatal("frozen path mismatch", path)
				}
			}
			for _, cmd := range []string{"set Single unicode", "set Multi lines", "set Preview long", "set Single same", "readonly Single", "writable Single", "disable Single", "enable Single", "load crlf", "constraint Single multi", "constraint Single single", "reload"} {
				if err := c.Run(cmd); err != nil {
					t.Fatal(cmd, err)
				}
			}
			a, changes := f.Counts()
			if !zero(a) || !zero(changes) {
				t.Fatal("trusted mutation emitted action/change", a, changes)
			}
			before := c.Host.Current().Session.Snapshot()
			for _, cmd := range []string{"invoke SaveSingle", "Commit Single", "draft Single text", "clipboard unicode", "set single short", "set Single unknown", "constraint Single nonsense", "reject-edit Missing", "action Missing error"} {
				if err := c.Run(cmd); err == nil {
					t.Fatal("bad command admitted", cmd)
				}
			}
			if !reflect.DeepEqual(before, c.Host.Current().Session.Snapshot()) {
				t.Fatal("bad commands changed session")
			}
		})
	}
}
func TestEchoFailuresPreserveTruthfulDomain(t *testing.T) {
	for _, mode := range []string{"success", "error", "malformed", "non-echo", "draft-conflict"} {
		t.Run(mode, func(t *testing.T) {
			f, c := newHost(t, false)
			if err := c.Run("action SaveSingle " + mode); err != nil {
				t.Fatal(err)
			}
			err := edit(c, Targets["Single"], "Exact 日本語", true)
			if (err == nil) != (mode == "success") {
				t.Fatal(mode, err)
			}
			w, _ := c.Host.Current().Session.Widget(Targets["Single"])
			calls, _ := f.Counts()
			if calls["SaveSingle"] != 1 {
				t.Fatal(calls)
			}
			if mode == "success" && w.Value != "Exact 日本語" {
				t.Fatal(w)
			}
			if mode != "success" && w.Value != "Initial single" {
				t.Fatal(w)
			}
			if mode == "draft-conflict" && w.Draft != "Exact 日本語 newer" {
				t.Fatal(w)
			}
			want := 1
			if mode == "error" || mode == "malformed" {
				want = 0
			}
			if f.Domain().Commits["SaveSingle"] != want {
				t.Fatal("domain accounting", f.Domain())
			}
		})
	}
}
func TestPublicationDiagnosticOneShotAndFailedReload(t *testing.T) {
	f, c := newHost(t, false)
	if err := c.Run("reject-edit Single"); err != nil {
		t.Fatal(err)
	}
	if err := c.Run("disable Required"); err != nil {
		t.Fatal(err)
	}
	before := c.Host.Current()
	if err := edit(c, Targets["Single"], "rejected", false); err == nil {
		t.Fatal("missing gate rejection")
	}
	w, _ := before.Session.Widget(Targets["Single"])
	if w.Draft != "Initial single" {
		t.Fatal("failed publication leaked", w)
	}
	if err := edit(c, Targets["Single"], "admitted", false); err != nil {
		t.Fatal(err)
	}
	for _, failure := range []string{"profile", "binding", "resource", "guard", "stale"} {
		if err := c.Run("fail " + failure); err == nil {
			t.Fatal("failure succeeded", failure)
		}
		if before != c.Host.Current() {
			t.Fatal("failed reload replaced bundle")
		}
		w, _ = before.Session.Widget(Targets["Single"])
		if w.Draft != "admitted" {
			t.Fatal("failed reload lost draft", w)
		}
	}
	calls, changes := f.Counts()
	if !zero(calls) || changes[Targets["Single"]] != 1 {
		t.Fatal(calls, changes)
	}
}
func TestTextDialogCancelRetainsChildCommitAndAcceptCaptures(t *testing.T) {
	f, c := newHost(t, false)
	target := open(t, c)
	if err := edit(c, Targets["FormMulti"], "proposal\n日本語", false); err != nil {
		t.Fatal(err)
	}
	if err := edit(c, Targets["Child"], "Saved child", true); err != nil {
		t.Fatal(err)
	}
	if err := terminal(c, target, ui.Cancel); err != nil {
		t.Fatal(err)
	}
	child, _ := c.Host.Current().Session.Widget(Targets["Child"])
	body, _ := c.Host.Current().Session.Widget(Targets["FormMulti"])
	if child.Value != "Saved child" || body.Draft != "Initial body" || f.Domain().TextCommits != 0 || f.Domain().Commits["SaveChild"] != 1 {
		t.Fatal(child, body, f.Domain())
	}
	target = open(t, c)
	if err := edit(c, Targets["FormName"], "日本語", false); err != nil {
		t.Fatal(err)
	}
	if err := edit(c, Targets["FormMulti"], "one\r\ntwo", false); err != nil {
		t.Fatal(err)
	}
	if err := terminal(c, target, ui.Accept); err != nil {
		t.Fatal(err)
	}
	if f.Domain().TextCommits != 1 || f.Domain().TextValues["Name"] != "日本語" || f.Domain().TextValues["Body"] != "one\r\ntwo" {
		t.Fatal(f.Domain())
	}
}
func TestSourceTighteningLeavesFailedCandidateAndDraftIntact(t *testing.T) {
	_, c := newHost(t, false)
	if err := c.Run("set Multi empty"); err != nil {
		t.Fatal(err)
	}
	if err := edit(c, Targets["Multi"], "draft\nline", false); err != nil {
		t.Fatal(err)
	}
	source, bundle := c.Source, c.Host.Current()
	if err := c.Run("constraint Multi single"); err == nil {
		t.Fatal("multiline draft silently converted")
	}
	if source != c.Source || bundle != c.Host.Current() {
		t.Fatal("failed constraint replaced source")
	}
	if err := c.Run("constraint Multi required"); err == nil {
		t.Fatal("nonempty draft repaired invalid accepted candidate")
	}
}

func TestRequiredEmptyStartupAndReload(t *testing.T) {
	f, c := newHostSource(t, RequiredEmptySource(false))
	for i := 0; i < 2; i++ {
		w, _ := c.Host.Current().Session.Widget(Targets["Required"])
		field, _ := c.Host.Current().Session.Field(w.Handle)
		if w.Value != "" || w.Draft != "" || !field.Required || field.Validation.Code == "" {
			t.Fatal("required absence lost", field)
		}
		if i == 0 {
			if err := c.Run("reload"); err != nil {
				t.Fatal(err)
			}
		}
	}
	actions, _ := f.Counts()
	if !zero(actions) {
		t.Fatal("absence auto-committed", actions)
	}
	if err := edit(c, Targets["Required"], "explicit repair", true); err != nil {
		t.Fatal(err)
	}
	if f.Domain().Commits["SaveRequired"] != 1 {
		t.Fatal(f.Domain())
	}
}
