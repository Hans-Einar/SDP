package application

import (
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	"os"
	"strings"
	"testing"
)

func TestReloadSDLPreservesDomainAndInvalidatesUIEvents(t *testing.T) {
	s, e := os.ReadFile("../edit-apt-cell.sdl")
	if e != nil {
		t.Fatal(e)
	}
	u, e := os.ReadFile("../edit-apt-cell.sdui")
	if e != nil {
		t.Fatal(e)
	}
	app, e := Load(string(s), string(u))
	if e != nil {
		t.Fatal(e)
	}
	edit, _ := app.UI.Widget("page/value")
	button, _ := app.UI.Widget("page/apply")
	app.UI.Draft(edit.Handle, "430")
	event := ui.Event{Handle: button.Handle, Kind: ui.Activate, ModelRevision: 1, Sequence: 1}
	if e = app.UI.Dispatch(event); e != nil {
		t.Fatal(e)
	}
	app.UI.Draft(edit.Handle, "440")
	program, e := parser.CompileActions(string(s))
	if e != nil {
		t.Fatal(e)
	}
	if e = app.ReloadSDL(program); e != nil {
		t.Fatal(e)
	}
	state, _ := app.UI.Widget("page/value")
	value, rev, trace := app.Domain.Snapshot()
	if state.Draft != "440" || value != "430" || rev != 1 || len(trace) != 1 || app.SDL.Revision() != 2 {
		t.Fatal("reload reset or replayed")
	}
	event.Sequence = 2
	if e = app.UI.Dispatch(event); e == nil {
		t.Fatal("old UI epoch survived SDL reload")
	}
	event.ModelRevision = app.UI.Revision
	if e = app.UI.Dispatch(event); e != nil {
		t.Fatal(e)
	}
	value, rev, trace = app.Domain.Snapshot()
	if value != "440" || rev != 2 || len(trace) != 2 {
		t.Fatal(value, rev, trace)
	}
	bad, e := parser.CompileActions(strings.Replace(string(s), "GoEditAptCell", "MissingGo", 1))
	if e != nil {
		t.Fatal(e)
	}
	epoch := app.UI.Revision
	if e = app.ReloadSDL(bad); e == nil || app.SDL.Revision() != 2 || app.UI.Revision != epoch {
		t.Fatal("failed SDL reload partially published")
	}
}
