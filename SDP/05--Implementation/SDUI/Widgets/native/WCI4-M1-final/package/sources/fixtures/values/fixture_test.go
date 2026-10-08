package values

import (
	"reflect"
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/Hans-Einar/SDP/SDUI/go/host/fynehost"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
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
func countsZero(m map[string]int) bool {
	for _, v := range m {
		if v != 0 {
			return false
		}
	}
	return true
}
func TestValuesActualRequestAdmissionAndTrustedMute(t *testing.T) {
	for _, nonmodal := range []bool{false, true} {
		t.Run(map[bool]string{false: "modal", true: "nonmodal"}[nonmodal], func(t *testing.T) {
			f, c := newHost(t, nonmodal)
			s := c.Host.Current().Session
			for _, path := range Targets {
				w, ok := s.Widget(path)
				if !ok {
					t.Fatal("missing path", path)
				}
				if _, ok = s.Field(w.Handle); !ok {
					t.Fatal("missing typed state", path)
				}
			}
			for _, path := range []string{TextNamePath, TextPath + "/textActions/textAccept", TextPath + "/textActions/textCancel", TextPath + "/textActions/textClose"} {
				if _, ok := s.Widget(path); !ok {
					t.Fatal("wrong frozen path", path)
				}
			}
			for _, cmd := range []string{"set Flag true", "set Level 30", "set Mode beta", "set Count 20", "readonly Count", "set Count 21", "writable Count", "disable Flag", "enable Flag", "options Mode reordered"} {
				if err := c.Run(cmd); err != nil {
					t.Fatal(cmd, err)
				}
			}
			actions, changes := f.Counts()
			if !countsZero(actions) || !countsZero(changes) {
				t.Fatal("trusted changes invoked callbacks", actions, changes)
			}
			before := s.Snapshot()
			for _, cmd := range []string{"set Count 2.5", "set Flag TRUE", "set Mode missing", "set count 3", "options Mode unknown", "observe Flag success", "accept mixed malformed", "action ChildCount options-conflict", "Commit Count", "invoke AcceptFlag"} {
				if err := c.Run(cmd); err == nil {
					t.Fatal("invalid control accepted", cmd)
				}
			}
			if !reflect.DeepEqual(before, s.Snapshot()) {
				t.Fatal("invalid conditions mutated state")
			}
		})
	}
}
func editCommit(c *Controller, path string, value ui.Value) error {
	return c.Host.Mutate(func(s *ui.Session) error {
		w, _ := s.Widget(path)
		var change ui.FieldChange
		var err error
		if w.Handle.Kind == "select" {
			target, e := s.Option(w.Handle, value.OptionID)
			if e != nil {
				return e
			}
			change, err = s.ChooseOption(target)
		} else {
			change, err = s.EditField(w.Handle, s.Revision, value)
		}
		if err != nil {
			return err
		}
		event, err := s.CaptureCommit(change.Field.Target)
		if err != nil {
			return err
		}
		return s.Dispatch(event)
	})
}
func TestObserverReentryStalesOriginalGestureWithoutDomainAction(t *testing.T) {
	f, c := newHost(t, false)
	if err := c.Run("observe Flag draft-conflict"); err != nil {
		t.Fatal(err)
	}
	if err := editCommit(c, Targets["Flag"], ui.Bool(true)); err == nil {
		t.Fatal("original gesture survived observer's newer edit")
	}
	actions, changes := f.Counts()
	if !countsZero(actions) || changes[Targets["Flag"]] != 2 {
		t.Fatal(actions, changes)
	}
	s := c.Host.Current().Session
	w, _ := s.Widget(Targets["Flag"])
	field, _ := s.Field(w.Handle)
	if field.Accepted != ui.Bool(false) || field.Proposed != ui.Bool(false) {
		t.Fatal("observer reentry state lost", field)
	}
	if err := editCommit(c, Targets["Flag"], ui.Bool(true)); err != nil {
		t.Fatal(err)
	}
	actions, _ = f.Counts()
	if actions["AcceptFlag"] != 1 {
		t.Fatal("observer was not one-shot", actions)
	}
}
func TestEchoDomainFailuresAndOptionConflict(t *testing.T) {
	for _, mode := range []string{"error", "malformed", "draft-conflict", "options-conflict"} {
		t.Run(mode, func(t *testing.T) {
			f, c := newHost(t, false)
			name, path, value := "AcceptCount", Targets["Count"], ui.Text("3")
			if mode == "options-conflict" {
				name, path, value = "AcceptMode", Targets["Mode"], ui.Choice("beta")
			}
			if err := c.Run("action " + name + " " + mode); err != nil {
				t.Fatal(err)
			}
			s := c.Host.Current().Session
			w, _ := s.Widget(path)
			before, _ := s.Field(w.Handle)
			if err := editCommit(c, path, value); err == nil {
				t.Fatal("failed delivery succeeded")
			}
			after, _ := s.Field(w.Handle)
			actions, _ := f.Counts()
			if after.Accepted != before.Accepted || actions[name] != 1 {
				t.Fatal(after, actions)
			}
			if (mode == "draft-conflict" || mode == "options-conflict") && f.Domain().ScalarCommits[name] != 1 {
				t.Fatal("domain commit hidden", f.Domain())
			}
		})
	}
}
func open(t *testing.T, c *Controller, path string) ui.SurfaceTarget {
	t.Helper()
	var target ui.SurfaceTarget
	if err := c.Host.Mutate(func(s *ui.Session) error {
		h, _ := s.Surface(path)
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
func TestMixedProposalCancelChildCommitAndAccept(t *testing.T) {
	f, c := newHost(t, false)
	target := open(t, c, MixedPath)
	if err := editCommit(c, Targets["FormFlag"], ui.Bool(true)); err != nil {
		t.Fatal(err)
	}
	if err := editCommit(c, Targets["ChildCount"], ui.Text("3")); err != nil {
		t.Fatal(err)
	}
	if err := terminal(c, target, ui.Cancel); err != nil {
		t.Fatal(err)
	}
	s := c.Host.Current().Session
	flag, _ := s.Widget(Targets["FormFlag"])
	child, _ := s.Widget(Targets["ChildCount"])
	a, _ := s.Field(flag.Handle)
	b, _ := s.Field(child.Handle)
	if a.Accepted != ui.Bool(false) || a.Proposed != ui.Bool(false) || b.Accepted != ui.Numeric(3) || b.Proposed != ui.Numeric(3) || f.Domain().ScalarCommits["ChildCount"] != 1 || f.Domain().MixedCommits != 0 {
		t.Fatal(a, b, f.Domain())
	}
	target = open(t, c, MixedPath)
	if err := editCommit(c, Targets["FormFlag"], ui.Bool(true)); err != nil {
		t.Fatal(err)
	}
	if err := terminal(c, target, ui.Accept); err != nil {
		t.Fatal(err)
	}
	if f.Domain().MixedCommits != 1 || f.Domain().Mixed[Targets["FormFlag"]] != ui.Bool(true) {
		t.Fatal(f.Domain())
	}
}
func TestMixedConflictAndTextOnlyAccept(t *testing.T) {
	for _, form := range []string{"mixed", "text"} {
		t.Run(form, func(t *testing.T) {
			f, c := newHost(t, false)
			path := MixedPath
			if form == "text" {
				path = TextPath
			}
			target := open(t, c, path)
			if err := c.Run("accept " + form + " draft-conflict"); err != nil {
				t.Fatal(err)
			}
			if err := terminal(c, target, ui.Accept); err == nil {
				t.Fatal("post-domain conflict missing")
			}
			state, _ := c.Host.Current().Session.SurfaceState(target.Handle)
			if !state.AcceptBlocked || state.Domain != ui.DomainSucceeded {
				t.Fatal(state)
			}
			if err := terminal(c, target, ui.Cancel); err != nil {
				t.Fatal(err)
			}
			if form == "mixed" && f.Domain().MixedCommits != 1 || form == "text" && f.Domain().TextCommits != 1 {
				t.Fatal("domain state rolled back", f.Domain())
			}
		})
	}
}
func TestFailedReloadRetainsLiveDraftAndOptions(t *testing.T) {
	f, c := newHost(t, false)
	if err := c.Host.Mutate(func(s *ui.Session) error {
		w, _ := s.Widget(Targets["Count"])
		_, err := s.EditField(w.Handle, s.Revision, ui.Text("-"))
		return err
	}); err != nil {
		t.Fatal(err)
	}
	for _, mode := range []string{"profile", "binding", "resource", "guard", "stale"} {
		old := c.Host.Current()
		if err := c.Run("fail " + mode); err == nil {
			t.Fatal("candidate did not fail", mode)
		}
		if c.Host.Current() != old || old.Session.Closed() {
			t.Fatal("failed candidate replaced live")
		}
	}
	if err := c.Run("reload"); err != nil {
		t.Fatal(err)
	}
	s := c.Host.Current().Session
	w, _ := s.Widget(Targets["Count"])
	field, _ := s.Field(w.Handle)
	if field.RawDraft == nil || *field.RawDraft != "-" || field.Validation.Code == "" {
		t.Fatal("invalid retained raw draft lost", field)
	}
	calls, _ := f.Counts()
	if !countsZero(calls) {
		t.Fatal("reload invoked domain", calls)
	}
}

func TestRequiredEmptyVariantAdmissionReloadAndExplicitSDLRepair(t *testing.T) {
	for _, nonmodal := range []bool{false, true} {
		t.Run(map[bool]string{false: "modal", true: "nonmodal"}[nonmodal], func(t *testing.T) {
			f, c := newHostSource(t, RequiredEmptySource(nonmodal))
			checkEmpty := func() {
				t.Helper()
				s := c.Host.Current().Session
				w, _ := s.Widget(Targets["Mode"])
				field, _ := s.Field(w.Handle)
				if !field.Required || field.Accepted != ui.Choice("") || field.Proposed != ui.Choice("") || field.Validation.Code == "" || !reflect.DeepEqual(field.Options, OptionSet("initial")) {
					t.Fatal("required absence/options changed", field)
				}
				if _, err := s.CaptureCommit(field.Target); err == nil {
					t.Fatal("invalid absence can commit")
				}
				form, _ := s.Widget(Targets["FormMode"])
				formField, _ := s.Field(form.Handle)
				if formField.Required || formField.Accepted != ui.Choice("alpha") || formField.Proposed != ui.Choice("alpha") {
					t.Fatal("startup variant changed FormMode", formField)
				}
				actions, changes := f.Counts()
				if !countsZero(actions) || !countsZero(changes) {
					t.Fatal("startup/reload dispatched user activity", actions, changes)
				}
			}
			checkEmpty()
			old := c.Host.Current()
			if err := c.Run("reload"); err != nil {
				t.Fatal("compatible required absence did not reload", err)
			}
			if !old.Session.Closed() {
				t.Fatal("predecessor not closed")
			}
			checkEmpty()
			// Runtime/SDL workflow evidence; the coordinator separately supplies native input.
			if err := editCommit(c, Targets["Mode"], ui.Choice("beta")); err != nil {
				t.Fatal(err)
			}
			s := c.Host.Current().Session
			w, _ := s.Widget(Targets["Mode"])
			field, _ := s.Field(w.Handle)
			actions, _ := f.Counts()
			if field.Accepted != ui.Choice("beta") || field.Validation.Code != "" || field.Dirty || actions["AcceptMode"] != 1 || f.Domain().ScalarCommits["AcceptMode"] != 1 {
				t.Fatal("explicit SDL repair failed", field, actions, f.Domain())
			}
		})
	}
}
