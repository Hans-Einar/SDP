package bridge_test

import (
	"context"
	"github.com/Hans-Einar/SDP/SDUI/go/host/fynehost/admission"
	"github.com/Hans-Einar/SDP/SDUI/go/layout"
	uiparser "github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/preparation"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/bridge"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	sdl "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/runtime"
	"reflect"
	"testing"
)

func TestDetachedPreparationRealBridge(t *testing.T) {
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
	live, e := ui.New("live", roots["page"])
	if e != nil {
		t.Fatal(e)
	}
	defer live.Close()
	plans := func() map[string]bridge.Plan {
		return map[string]bridge.Plan{"api.Echo": {Inputs: map[string]bridge.Source{"Value": {Widget: "page/edit"}}, OutputField: "Value"}}
	}
	modules := map[string]*sdl.Engine{"api": engine}
	if _, e = bridge.Bind(context.Background(), live, d, modules, plans(), nil); e != nil {
		t.Fatal(e)
	}
	before := live.Widgets()
	revision := live.Revision
	for _, kind := range []string{"module", "signature", "result", "success"} {
		t.Run(kind, func(t *testing.T) {
			var detached *ui.Session
			r := preparation.Request{Document: d, Entry: "page", SessionID: "detached", SourceRevision: "source-1", Mode: preparation.Connected, Capabilities: admission.Capabilities(), ValidateLayout: func(n *uiparser.Instance) error { return admission.Layout(n, layout.Size{W: 1280, H: 800}) }}
			r.Bind = func(s *ui.Session, selected *uiparser.Document) error {
				detached = s
				m := modules
				bp := plans()
				switch kind {
				case "module":
					m = nil
				case "signature":
					v := sdl.Integer(2)
					bp["api.Echo"].Inputs["Value"] = bridge.Source{Literal: &v}
				case "result":
					plan := bp["api.Echo"]
					plan.OutputField = "Missing"
					bp["api.Echo"] = plan
				}
				_, err := bridge.Bind(context.Background(), s, selected, m, bp, nil)
				return err
			}
			c, err := preparation.Prepare(r)
			if calls != 0 || live.Revision != revision || !reflect.DeepEqual(before, live.Widgets()) {
				t.Fatal("preparation changed live state or invoked domain")
			}
			if kind != "success" {
				if err == nil || c != nil || detached == nil || !detached.Closed() {
					t.Fatal("failed candidate leaked", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			defer c.Close()
			if err = c.Admit("source-1"); err != nil {
				t.Fatal(err)
			}
			button, _ := c.Session.Widget("page/ok")
			if err = c.Session.Dispatch(ui.Event{Handle: button.Handle, Kind: ui.Activate, ModelRevision: c.Session.Revision, Sequence: 1}); err != nil {
				t.Fatal(err)
			}
			if calls != 1 {
				t.Fatal("prepared handler not usable")
			}
			calls = 0
		})
	}
	button, _ := live.Widget("page/ok")
	if e = live.Dispatch(ui.Event{Handle: button.Handle, Kind: ui.Activate, ModelRevision: live.Revision, Sequence: 2}); e != nil || calls != 1 {
		t.Fatal("previous live handler not retained", e, calls)
	}
}
