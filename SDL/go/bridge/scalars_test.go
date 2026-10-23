package bridge_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	up "github.com/Hans-Einar/SDP/SDUI/go/parser"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/bridge"
	commands "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/examples/commands"
	sp "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	sdl "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/runtime"
)

func scalarEngine(t *testing.T, kind sp.ScalarType, handler sdl.Handler) *sdl.Engine {
	t.Helper()
	source := fmt.Sprintf("language action-core version 0.1.\naction Echo.\nrecord ValueRecord.\nEcho invokes GoEcho.\nEcho returns ValueRecord.\nEcho takes ValueRecord.\nValueRecord field Value as %s.\n", kind)
	p, err := sp.CompileActions(source)
	if err != nil {
		t.Fatal(err)
	}
	e, err := sdl.New(p, sdl.Registry{"GoEcho": {Input: sp.RecordType{"Value": kind}, Output: sp.RecordType{"Value": kind}, Call: handler}})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(e.Close)
	return e
}
func scalarSession(t *testing.T, call string) (*up.Document, *ui.Session) {
	t.Helper()
	source := `sdui 0.3;ref: api "api.sdl";page=[value=` + call + `;other=checkbox("Other")];api.Echo.setHandle(page.value);`
	d, roots, err := up.Compile(source)
	if err != nil {
		t.Fatal(err)
	}
	s, err := ui.New("scalar-test", roots["page"])
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(s.Close)
	if strings.HasPrefix(call, "select(") {
		if err = s.BindChoices(map[string][]ui.ChoiceOption{"page/value": {{ID: "a", Label: "Repeated", Enabled: true}, {ID: "b", Label: "Repeated", Enabled: true}, {ID: "blocked", Label: "Disabled"}}}); err != nil {
			t.Fatal(err)
		}
	}
	return d, s
}
func scalarPlan(selector bridge.EventField) map[string]bridge.Plan {
	return map[string]bridge.Plan{"api.Echo": {ResultMode: bridge.ScalarResult, Inputs: map[string]bridge.Source{"Value": {EventField: selector}}, OutputField: "Value"}}
}
func bindScalar(t *testing.T, d *up.Document, s *ui.Session, e *sdl.Engine, selector bridge.EventField) {
	t.Helper()
	if _, err := bridge.Bind(context.Background(), s, d, map[string]*sdl.Engine{"api": e}, scalarPlan(selector), nil); err != nil {
		t.Fatal(err)
	}
}
func proposeScalar(t *testing.T, s *ui.Session, value ui.Value) ui.Event {
	t.Helper()
	w, _ := s.Widget("page/value")
	var change ui.FieldChange
	var err error
	if w.Handle.Kind == "select" {
		target, e := s.Option(w.Handle, value.OptionID)
		if e != nil {
			t.Fatal(e)
		}
		change, err = s.ChooseOption(target)
	} else {
		change, err = s.EditField(w.Handle, s.Revision, value)
	}
	if err != nil {
		t.Fatal(err)
	}
	event, err := s.CaptureCommit(change.Field.Target)
	if err != nil {
		t.Fatal(err)
	}
	return event
}
func TestTypedScalarEchoActualSDLAndNoReplay(t *testing.T) {
	tests := []struct {
		name, call     string
		kind           sp.ScalarType
		selector       bridge.EventField
		edit, accepted ui.Value
	}{
		{"checkbox", `checkbox("Flag",callback=api.Echo.@invoke)`, sp.BooleanType, bridge.ControlBoolean, ui.Bool(true), ui.Bool(true)},
		{"slider", `slider("Level",min=0,max=100,step=5,value=25,callback=api.Echo.@invoke)`, sp.IntegerType, bridge.ControlNumber, ui.Numeric(30), ui.Numeric(30)},
		{"number", `number("Count",min=-10,max=100,step=1,value=1,callback=api.Echo.@invoke)`, sp.IntegerType, bridge.ControlNumber, ui.Text("2e1"), ui.Numeric(20)},
		{"choice", `select("Mode",value="a",callback=api.Echo.@invoke)`, sp.TextType, bridge.ChoiceOptionID, ui.Choice("b"), ui.Choice("b")},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			d, s := scalarSession(t, tc.call)
			calls := 0
			var input sdl.Record
			e := scalarEngine(t, tc.kind, func(_ context.Context, r sdl.Record) (sdl.Record, error) { calls++; input = r; return r, nil })
			bindScalar(t, d, s, e, tc.selector)
			if calls != 0 {
				t.Fatal("preflight invoked SDL")
			}
			event := proposeScalar(t, s, tc.edit)
			if err := s.Dispatch(event); err != nil {
				t.Fatal(err)
			}
			f, _ := s.Field(event.Handle)
			if calls != 1 || f.Accepted != tc.accepted || f.Dirty || input["Value"].Type != tc.kind {
				t.Fatal(f, calls, input)
			}
			if err := s.Dispatch(event); err == nil || calls != 1 {
				t.Fatal("replay", err, calls)
			}
		})
	}
}
func TestScalarSafe53ExactNumericProvenance(t *testing.T) {
	call := `number("N",min=-9007199254740991,max=9007199254740991,step=1,value=0,callback=api.Echo.@invoke)`
	d, s := scalarSession(t, call)
	calls := 0
	e := scalarEngine(t, sp.IntegerType, func(_ context.Context, r sdl.Record) (sdl.Record, error) { calls++; return r, nil })
	bindScalar(t, d, s, e, bridge.ControlNumber)
	for _, raw := range []string{"9007199254740991", "-9007199254740991"} {
		event := proposeScalar(t, s, ui.Text(raw))
		if err := s.Dispatch(event); err != nil {
			t.Fatal(err)
		}
	}
	for _, raw := range []string{"9007199254740991.1", "9007199254740992", "1.00000000000000000001", "-", "1e999999999"} {
		w, _ := s.Widget("page/value")
		change, err := s.EditField(w.Handle, s.Revision, ui.Text(raw))
		if err != nil {
			t.Fatal("invalid intermediate draft must remain visible", err)
		}
		if change.Field.Validation.Code == "" || change.Field.RawDraft == nil || *change.Field.RawDraft != raw {
			t.Fatal("invalid raw lost", change)
		}
		if _, err = s.CaptureCommit(change.Field.Target); err == nil {
			t.Fatal("invalid decimal captured", raw)
		}
	}
	if calls != 2 {
		t.Fatal("invalid drafts reached SDL", calls)
	}
}
func TestScalarPreflightRejectsFractionalGridAndWrongModes(t *testing.T) {
	d, s := scalarSession(t, `number("N",min=0,max=1,step=0.1,value=0.3,callback=api.Echo.@invoke)`)
	calls := 0
	e := scalarEngine(t, sp.IntegerType, func(_ context.Context, r sdl.Record) (sdl.Record, error) { calls++; return r, nil })
	if _, err := bridge.Bind(context.Background(), s, d, map[string]*sdl.Engine{"api": e}, scalarPlan(bridge.ControlNumber), nil); err == nil || !strings.Contains(err.Error(), "page/value") {
		t.Fatal("fractional SDL preparation admitted/lost provenance", err)
	}
	for _, mode := range []string{"text", "unknown", "acceptfield", "wrongselector", "wrongreceiver", "missingoutput", "revision"} {
		t.Run(mode, func(t *testing.T) {
			d, s := scalarSession(t, `checkbox("Flag",callback=api.Echo.@invoke)`)
			e := scalarEngine(t, sp.BooleanType, func(_ context.Context, r sdl.Record) (sdl.Record, error) { calls++; return r, nil })
			plans := scalarPlan(bridge.ControlBoolean)
			p := plans["api.Echo"]
			switch mode {
			case "text":
				p.ResultMode = bridge.TextResult
			case "unknown":
				p.ResultMode = "other"
			case "acceptfield":
				p.AcceptField = "Value"
			case "wrongselector":
				p.Inputs["Value"] = bridge.Source{EventField: bridge.ControlNumber}
			case "wrongreceiver":
				d.Connections[0].Path = []string{"other"}
			case "missingoutput":
				p.OutputField = ""
			case "revision":
				p.RevisionField = "Value"
				p.RevisionContext = "Revision"
			}
			plans["api.Echo"] = p
			if _, err := bridge.Bind(context.Background(), s, d, map[string]*sdl.Engine{"api": e}, plans, nil); err == nil {
				t.Fatal("bad scalar plan admitted", mode)
			}
			w, _ := s.Widget("page/value")
			if s.HasBinding(w.Handle) || calls != 0 {
				t.Fatal("partial install/domain call")
			}
		})
	}
}
func TestScalarErrorsConflictsAndResultBounds(t *testing.T) {
	for _, mode := range []string{"error", "malformed", "unsafe", "different", "draft", "options"} {
		t.Run(mode, func(t *testing.T) {
			call := `number("N",min=0,max=100,step=1,value=1,callback=api.Echo.@invoke)`
			kind := sp.IntegerType
			selector := bridge.ControlNumber
			if mode == "options" {
				call = `select("Mode",value="a",callback=api.Echo.@invoke)`
				kind = sp.TextType
				selector = bridge.ChoiceOptionID
			}
			d, s := scalarSession(t, call)
			calls := 0
			e := scalarEngine(t, kind, func(_ context.Context, r sdl.Record) (sdl.Record, error) {
				calls++
				switch mode {
				case "error":
					return nil, fmt.Errorf("failed")
				case "malformed":
					return sdl.Record{"Value": sdl.Boolean(true)}, nil
				case "unsafe":
					return sdl.Record{"Value": sdl.Integer(9007199254740992)}, nil
				case "different":
					return sdl.Record{"Value": sdl.Integer(3)}, nil
				case "draft":
					w, _ := s.Widget("page/value")
					if _, err := s.EditField(w.Handle, s.Revision, ui.Text("3")); err != nil {
						t.Fatal(err)
					}
				case "options":
					w, _ := s.Widget("page/value")
					f, _ := s.Field(w.Handle)
					if err := s.ReplaceChoices(w.Handle, f.Target.OptionGeneration, f.Options); err != nil {
						t.Fatal(err)
					}
				}
				return r, nil
			})
			bindScalar(t, d, s, e, selector)
			edit := ui.Text("2")
			if mode == "options" {
				edit = ui.Choice("b")
			}
			event := proposeScalar(t, s, edit)
			before, _ := s.Field(event.Handle)
			if err := s.Dispatch(event); err == nil {
				t.Fatal("invalid/stale result accepted", mode)
			}
			after, _ := s.Field(event.Handle)
			if after.Accepted != before.Accepted || calls != 1 {
				t.Fatal("failed result changed accepted", after, calls)
			}
			if err := s.Dispatch(event); err == nil || calls != 1 {
				t.Fatal("failed command replay", err, calls)
			}
		})
	}
}

func TestSDLDialogAcceptRejectsMixedTypedOwner(t *testing.T) {
	source := strings.Replace(commands.Source, `  form=[`, ` typed=checkbox("Typed field"); form=[`, 1)
	d, s, f := commandSession(t, source)
	_, err := bridge.Bind(context.Background(), s, d, map[string]*sdl.Engine{"actions": f.Engine}, commands.Plans(), nil)
	if err == nil || !strings.Contains(err.Error(), "text-only") {
		t.Fatal("mixed SDL Accept was not rejected explicitly", err)
	}
	if len(f.Calls()) != 0 {
		t.Fatal("preflight invoked domain")
	}
}
