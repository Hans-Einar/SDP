package previews

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/Hans-Einar/SDP/SDUI/go/host/fynehost"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	"github.com/Hans-Einar/SDP/SDUI/go/markdown"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/bridge"
	sdl "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/runtime"
)

func newHost(t *testing.T, nonmodal bool) (*Fixture, *Controller) {
	t.Helper()
	a := test.NewTempApp(t)
	w := a.NewWindow(WindowTitle)
	t.Cleanup(w.Close)
	f, err := New(nonmodal)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(f.Close)
	h := fynehost.NewDocumentHost(w.Canvas(), layout.Size{W: 1100, H: 850})
	t.Cleanup(h.Close)
	c := &Controller{Fixture: f, Host: h}
	c.InstallHooks()
	c.Resize = h.Resize
	t.Cleanup(c.Abandon)
	if err := c.Run("reload"); err != nil {
		t.Fatal(err)
	}
	return f, c
}
func zero(m map[string]int) bool {
	for _, v := range m {
		if v != 0 {
			return false
		}
	}
	return true
}
func TestRequestIdentityAndRealSDLPreflight(t *testing.T) {
	for _, nonmodal := range []bool{false, true} {
		f, err := New(nonmodal)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		r, err := f.Request(1)
		if err != nil {
			t.Fatal(err)
		}
		root, err := parser.Normalize(r.Document)
		if err != nil {
			t.Fatal(err)
		}
		s, err := ui.New("test", root["page"])
		if err != nil {
			t.Fatal(err)
		}
		defer s.Close()
		if _, err = bridge.Bind(context.Background(), s, r.Document, map[string]*sdl.Engine{"previews": f.Engine}, Plans(), nil); err != nil {
			t.Fatal(err)
		}
		for _, alias := range []string{"Hero", "ReportSVG", "DetailSVG"} {
			resource, ok := r.SVGResources[Targets[alias]]
			if !ok || resource.SHA256 != digest(resource.Resource.SVG) {
				t.Fatal(alias, resource)
			}
		}
		for _, alias := range []string{"Note", "LoadStatus", "MarkStatus", "DetailNote"} {
			w, ok := s.Widget(Targets[alias])
			if !ok {
				t.Fatal(alias)
			}
			field, ok := s.Field(w.Handle)
			if !ok || field.Input == nil {
				t.Fatal("expected extended", alias)
			}
		}
		for _, path := range []string{"page/view/header/toolbar/loadButton", "page/view/header/toolbar/markButton", "page/view/header/toolbar/openButton", DetailPath + "/detailActions/acceptButton", DetailPath + "/detailActions/cancelButton", DetailPath + "/detailActions/closeButton"} {
			if _, ok := s.Widget(path); !ok {
				t.Fatal(path)
			}
		}
		calls, changes, renderers := f.Counts()
		if !zero(calls) || !zero(changes) || !zero(renderers) {
			t.Fatal("preflight executed callback")
		}
	}
}
func TestNativeRequestAdmission(t *testing.T) {
	f, c := newHost(t, false)
	state := c.Inspect()
	if state["previews"] == nil {
		t.Fatal("missing actual host preview inspector")
	}
	if got := len(c.Host.Current().Session.Snapshot().Surfaces); got != 1 {
		t.Fatal(got)
	}
	calls, changes, before := f.Counts()
	if !zero(calls) || !zero(changes) {
		t.Fatal("admission executed domain")
	}
	for _, cmd := range []string{"state", "mutate-input Hero", "caption Hero long", "caption Hero empty", "resize 1200 900", "resize 1100 850"} {
		if err := c.Run(cmd); err != nil {
			t.Fatal(cmd, err)
		}
	}
	_, _, after := f.Counts()
	if !reflect.DeepEqual(before, after) {
		t.Fatal("renderer ran after preparation", before, after)
	}
	if err := c.Run("invoke Mark"); err == nil {
		t.Fatal("control channel simulates action")
	}
}
func TestRequestGuardAndLoanSeparation(t *testing.T) {
	f, err := New(false)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	r, err := f.Request(1)
	if err != nil {
		t.Fatal(err)
	}
	loan := r.SVGResources[Targets["Hero"]].Resource.SVG
	loan[0] = '!'
	if err = r.Guard(); err != nil {
		t.Fatal("loan mutation changed authority", err)
	}
	next, err := f.Request(2)
	if err != nil {
		t.Fatal(err)
	}
	if next.SVGResources[Targets["Hero"]].Resource.SVG[0] != '<' {
		t.Fatal("request loan aliased catalogue")
	}
	for _, change := range [][3]string{{"resource", "Hero", "alternate"}, {"renderer-revision", "ReportA", "r2"}, {"policy", "Hero", "label"}} {
		r, err = f.Request(3)
		if err != nil {
			t.Fatal(err)
		}
		if err = f.configure(change[0], change[1], change[2]); err != nil {
			t.Fatal(err)
		}
		if err = r.Guard(); err == nil || !strings.Contains(err.Error(), "stale") {
			t.Fatal("changed authority admitted", change, err)
		}
	}
}

func TestImmutableProviderInputsAndExactRoot(t *testing.T) {
	f, err := New(false)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	r, err := f.Request(1)
	if err != nil {
		t.Fatal(err)
	}
	roots, err := parser.Normalize(r.Document)
	if err != nil {
		t.Fatal(err)
	}
	// A deterministic component-only backend permits actual renderer-output copy
	// checks. This is explicitly not the unsupported native diagram representation.
	backend := markdown.PreviewBackend{SVG: func(markdown.Resource) error { return nil }, Markdown: func() error { return nil }, Mermaid: func(markdown.Resource) error { return nil }}
	prepared, err := markdown.PreparePreviews(roots["page"], r.SVGResources, r.MarkdownRenderers, backend)
	if err != nil {
		t.Fatal(err)
	}
	c := &Controller{Fixture: f}
	if err = c.Run("mutate-input Hero"); err != nil {
		t.Fatal(err)
	}
	if err = c.Run("mutate-return ReportA"); err != nil {
		t.Fatal(err)
	}
	hero, ok := prepared.Outcome(Targets["Hero"])
	if !ok || hero.Resource == nil || hero.Resource.SVG[0] != '<' {
		t.Fatal("caller resource leaked", hero)
	}
	hero.Resource.SVG[0] = '!'
	hero, _ = prepared.Outcome(Targets["Hero"])
	if hero.Resource.SVG[0] != '<' {
		t.Fatal("accessor resource leaked")
	}
	diagram, _ := prepared.Outcome(Targets["ReportA"])
	if len(diagram.Diagrams) != 1 || diagram.Diagrams[0].Resource == nil || diagram.Diagrams[0].Resource.SVG[0] != '<' {
		t.Fatal("renderer return buffer leaked", diagram)
	}
	_, _, baseline := f.Counts()
	if baseline[Targets["ReportA"]] != 1 || baseline[Targets["ReportB"]] != 1 {
		t.Fatal("copy proof did not execute both renderers", baseline)
	}
	if err = prepared.Check(roots["page"]); err != nil {
		t.Fatal(err)
	}
	_, foreign, err := parser.Compile(strings.Replace(f.Source(), "Output by site", "Changed source description", 1))
	if err != nil {
		t.Fatal(err)
	}
	if err = prepared.Check(foreign["page"]); err == nil {
		t.Fatal("foreign source fingerprint accepted")
	}
	_, _, after := f.Counts()
	if !reflect.DeepEqual(baseline, after) {
		t.Fatal("pure check rendered")
	}
}

func TestFailedCandidateRetainsPublishedBundle(t *testing.T) {
	f, c := newHost(t, false)
	before := c.Host.Current()
	for _, commands := range [][]string{
		{"binding DetailSVG digest", "reload"},
		{"binding ReportA nil", "reload"},
		{"resource DetailSVG malformed", "reload"},
	} {
		if err := c.Run(commands[0]); err != nil {
			t.Fatal(err)
		}
		if err := c.Run(commands[1]); err == nil {
			t.Fatal("bad candidate accepted", commands)
		}
		if c.Host.Current() != before || before.Session.Closed() {
			t.Fatal("failed candidate replaced live bundle")
		}
	}
	if err := c.Run("resource DetailSVG good"); err != nil {
		t.Fatal(err)
	}
	if err := c.Run("prepare"); err != nil {
		t.Fatal(err)
	}
	if err := c.Run("renderer-revision ReportA r2"); err != nil {
		t.Fatal(err)
	}
	if err := c.Run("publish"); err == nil || !strings.Contains(err.Error(), "stale") {
		t.Fatal("stale renderer revision published", err)
	}
	if c.held != nil || c.Host.Current() != before {
		t.Fatal("failed publish retained candidate or replaced live")
	}
	if err := c.Run("reload"); err != nil {
		t.Fatal(err)
	}
	if !before.Session.Closed() {
		t.Fatal("successful replacement did not revoke old session")
	}
	before = c.Host.Current()
	if err := c.Run("fail resource-ticket"); err != nil {
		t.Fatal(err)
	}
	if err := c.Run("reload"); err == nil {
		t.Fatal("ticket failure not consumed")
	}
	if c.Host.Current() != before || before.Session.Closed() {
		t.Fatal("failed ticket replaced live")
	}
	calls, changes, _ := f.Counts()
	if !zero(calls) || !zero(changes) {
		t.Fatal("conditions invoked domain", calls, changes)
	}
}

func TestActualSDLCommandAndDialogComposition(t *testing.T) {
	f, c := newHost(t, true)
	s := c.Host.Current().Session
	_, _, baseline := f.Counts()
	// Component-level actual dispatch; the native CLI never exposes these gestures.
	if err := c.Host.Mutate(func(s *ui.Session) error {
		w, _ := s.Widget("page/view/header/toolbar/loadButton")
		return s.Dispatch(ui.Event{Kind: ui.Activate, Handle: w.Handle, ModelRevision: s.Revision, Sequence: s.Sequence() + 1})
	}); err != nil {
		t.Fatal(err)
	}
	if err := c.Host.Mutate(func(s *ui.Session) error {
		p := s.Snapshot().Presentations["page/view/header/toolbar/markButton"]
		event, err := s.CaptureCommand(p.Handle, "button", nil)
		if err != nil {
			return err
		}
		result, err := s.DispatchInteraction(event)
		if err == nil && result.Domain != ui.DomainSucceeded {
			t.Fatalf("unexpected result %+v", result)
		}
		return err
	}); err != nil {
		t.Fatal(err)
	}
	w, _ := s.Widget(Targets["MarkStatus"])
	if w.Value != "Marked: true" {
		t.Fatal("extended command receiver", w.Value)
	}
	var target ui.SurfaceTarget
	if err := c.Host.Mutate(func(s *ui.Session) error {
		h, _ := s.Surface(DetailPath)
		var err error
		target, err = s.OpenSurface(h, ui.ContextTarget{})
		return err
	}); err != nil {
		t.Fatal(err)
	}
	if err := c.Host.Mutate(func(s *ui.Session) error {
		w, _ := s.Widget(Targets["DetailNote"])
		change, err := s.EditField(w.Handle, s.Revision, ui.Text("Saved child note"))
		if err != nil {
			return err
		}
		event, err := s.CaptureCommit(change.Field.Target)
		if err != nil {
			return err
		}
		return s.Dispatch(event)
	}); err != nil {
		t.Fatal(err)
	}
	if err := c.Host.Mutate(func(s *ui.Session) error {
		event, err := s.CaptureDialog(target, ui.Cancel)
		if err != nil {
			return err
		}
		_, err = s.DispatchInteraction(event)
		return err
	}); err != nil {
		t.Fatal(err)
	}
	if f.Domain().Values["SaveDetail"] != "Saved child note" {
		t.Fatal("Cancel undid actual child domain commit")
	}
	if err := c.Host.Mutate(func(s *ui.Session) error {
		h, _ := s.Surface(DetailPath)
		var err error
		target, err = s.OpenSurface(h, ui.ContextTarget{})
		return err
	}); err != nil {
		t.Fatal(err)
	}
	if err := c.Run("accept false"); err != nil {
		t.Fatal(err)
	}
	for _, want := range []bool{false, true} {
		if err := c.Host.Mutate(func(s *ui.Session) error {
			event, err := s.CaptureDialog(target, ui.Accept)
			if err != nil {
				return err
			}
			_, err = s.DispatchInteraction(event)
			return err
		}); err != nil {
			t.Fatal(err)
		}
		if (f.Domain().DetailCommits == 1) != want {
			t.Fatal("Accept domain decision", f.Domain())
		}
	}
	calls, _, after := f.Counts()
	if calls["Load"] != 1 || calls["Mark"] != 1 || calls["SaveDetail"] != 1 || calls["AcceptDetail"] != 2 {
		t.Fatal(calls)
	}
	if !reflect.DeepEqual(baseline, after) {
		t.Fatal("interaction rerendered provider", baseline, after)
	}
}
