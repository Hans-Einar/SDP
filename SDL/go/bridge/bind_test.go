package bridge_test

import (
	"context"
	"os"
	"strings"
	"testing"

	uiparser "github.com/Hans-Einar/SDP/SDUI/go/parser"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/bridge"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/examples/simulation"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	sdl "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/runtime"
)

func source(t *testing.T, name string) string {
	t.Helper()
	s, e := os.ReadFile("../examples/" + name)
	if e != nil {
		t.Fatal(e)
	}
	return string(s)
}
func TestEchoButtonAndInput(t *testing.T) {
	p, e := parser.CompileActions(source(t, "echo.sdl"))
	if e != nil {
		t.Fatal(e)
	}
	calls := 0
	engine, e := sdl.New(p, sdl.Registry{"GoEcho": {Input: parser.RecordType{"Value": parser.TextType}, Output: parser.RecordType{"Value": parser.TextType}, Call: func(_ context.Context, r sdl.Record) (sdl.Record, error) { calls++; return r, nil }}})
	if e != nil {
		t.Fatal(e)
	}
	d, roots, e := uiparser.Compile(source(t, "echo.sdui"))
	if e != nil {
		t.Fatal(e)
	}
	session, e := ui.New("test", roots["page"])
	if e != nil {
		t.Fatal(e)
	}
	binding, e := bridge.Bind(context.Background(), session, d, map[string]*sdl.Engine{"api": engine}, map[string]bridge.Plan{"api.Echo": {Inputs: map[string]bridge.Source{"Value": {Widget: "page/edit"}}, OutputField: "Value"}}, nil)
	if e != nil {
		t.Fatal(e)
	}
	button, _ := session.Widget("page/ok")
	edit, _ := session.Widget("page/edit")
	session.Draft(edit.Handle, "button draft")
	if e = session.Dispatch(ui.Event{Handle: button.Handle, Kind: ui.Activate, ModelRevision: 1, Sequence: 1}); e != nil {
		t.Fatal(e)
	}
	edit, _ = session.Widget("page/edit")
	if edit.Value != "button draft" || calls != 1 {
		t.Fatal(edit, calls)
	}
	session.Draft(edit.Handle, "input draft")
	edit, _ = session.Widget("page/edit")
	if e = session.Dispatch(ui.Event{Handle: edit.Handle, Kind: ui.Commit, ModelRevision: 1, Sequence: 2, DraftRevision: edit.DraftRevision, Value: ui.Text(edit.Draft)}); e != nil {
		t.Fatal(e)
	}
	if calls != 2 || len(binding.Links) != 2 || binding.Links[0].SDLSpan.Line == 0 || binding.Links[0].UISpan.Line == 0 {
		t.Fatal("missing dispatch or source map")
	}
}
func TestSimulatedEditAptCell(t *testing.T) {
	p, e := parser.CompileActions(source(t, "edit-apt-cell.sdl"))
	if e != nil {
		t.Fatal(e)
	}
	domain := simulation.NewApt()
	engine, e := sdl.New(p, domain.Registry())
	if e != nil {
		t.Fatal(e)
	}
	d, roots, e := uiparser.Compile(source(t, "edit-apt-cell.sdui"))
	if e != nil {
		t.Fatal(e)
	}
	session, e := ui.New("apt", roots["page"])
	if e != nil {
		t.Fatal(e)
	}
	binding, e := bridge.Bind(context.Background(), session, d, map[string]*sdl.Engine{"apt": engine}, simulation.Plans(), simulation.InitialContext())
	if e != nil {
		t.Fatal(e)
	}
	edit, _ := session.Widget("page/value")
	button, _ := session.Widget("page/apply")
	session.Draft(edit.Handle, "430")
	event := ui.Event{Handle: button.Handle, Kind: ui.Activate, ModelRevision: 1, Sequence: 1}
	if e = session.Dispatch(event); e != nil {
		t.Fatal(e)
	}
	value, revision, trace := domain.Snapshot()
	state, _ := session.Widget("page/value")
	if value != "430" || state.Value != "430" || state.Dirty || revision != 1 || len(trace) != 1 || !trace[0].Accepted || binding.Context["AptRevision"].Integer != 1 {
		t.Fatal(value, revision, trace, state)
	}
	if e = session.Dispatch(event); e == nil {
		t.Fatal("duplicate edit accepted")
	}
	session.Draft(edit.Handle, "not a number")
	event.Sequence = 2
	if e = session.Dispatch(event); e == nil {
		t.Fatal("invalid domain value accepted")
	}
	value, revision, trace = domain.Snapshot()
	state, _ = session.Widget("page/value")
	if value != "430" || revision != 1 || state.Value != "430" || !state.Dirty || len(trace) != 2 || trace[1].Accepted {
		t.Fatal("rejected edit changed accepted truth")
	}
	binding.Context["AptRevision"] = sdl.Integer(0)
	session.Draft(edit.Handle, "440")
	event.Sequence = 3
	if e = session.Dispatch(event); e == nil {
		t.Fatal("stale revision accepted")
	}
	value, revision, _ = domain.Snapshot()
	if value != "430" || revision != 1 {
		t.Fatal("stale edit mutated domain")
	}
}

func TestInvalidBindingPlansInstallNothing(t *testing.T) {
	for _, kind := range []string{"module", "member", "target", "source", "literal", "unused"} {
		t.Run(kind, func(t *testing.T) {
			p, e := parser.CompileActions(source(t, "echo.sdl"))
			if e != nil {
				t.Fatal(e)
			}
			calls := 0
			engine, e := sdl.New(p, sdl.Registry{"GoEcho": {Input: parser.RecordType{"Value": parser.TextType}, Output: parser.RecordType{"Value": parser.TextType}, Call: func(_ context.Context, r sdl.Record) (sdl.Record, error) { calls++; return r, nil }}})
			if e != nil {
				t.Fatal(e)
			}
			text := source(t, "echo.sdui")
			if kind == "member" {
				text = strings.ReplaceAll(text, "@invoke", "@unknown")
			}
			if kind == "target" {
				text = strings.Replace(text, "setHandle(page.edit)", "setHandle(page.ok)", 1)
			}
			d, roots, e := uiparser.Compile(text)
			if e != nil {
				t.Fatal(e)
			}
			session, e := ui.New("test", roots["page"])
			if e != nil {
				t.Fatal(e)
			}
			modules := map[string]*sdl.Engine{"api": engine}
			plan := bridge.Plan{Inputs: map[string]bridge.Source{"Value": {Widget: "page/edit"}}, OutputField: "Value"}
			plans := map[string]bridge.Plan{"api.Echo": plan}
			switch kind {
			case "module":
				modules = nil
			case "source":
				plan.Inputs["Value"] = bridge.Source{Widget: "page/missing"}
			case "literal":
				bad := sdl.Integer(1)
				plan.Inputs["Value"] = bridge.Source{Literal: &bad}
			case "unused":
				plans["api.Extra"] = plan
			}
			if _, e = bridge.Bind(context.Background(), session, d, modules, plans, nil); e == nil {
				t.Fatal("invalid binding accepted")
			}
			button, _ := session.Widget("page/ok")
			if e = session.Dispatch(ui.Event{Handle: button.Handle, Kind: ui.Activate, ModelRevision: 1, Sequence: 1}); e == nil || !strings.Contains(e.Error(), "unbound") {
				t.Fatal("partial handler installation", e)
			}
			if calls != 0 {
				t.Fatal("invalid binding called domain")
			}
		})
	}
}
