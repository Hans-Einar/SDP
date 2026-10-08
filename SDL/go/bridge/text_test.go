package bridge_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	up "github.com/Hans-Einar/SDP/SDUI/go/parser"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/bridge"
	sp "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	sdl "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/runtime"
)

func textSession(t *testing.T, args, receiver string) (*up.Document, *ui.Session) {
	t.Helper()
	d, roots, err := up.Compile(`sdui 0.3;ref: api "api.sdl";page=[value=input("Text",value="Initial",callback=api.Echo.@invoke` + args + `);other=input("Other",readOnly=true)];api.Echo.setHandle(page.` + receiver + `);`)
	if err != nil {
		t.Fatal(err)
	}
	s, err := ui.New("text-test", roots["page"])
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(s.Close)
	return d, s
}
func textPlan(source bridge.Source) map[string]bridge.Plan {
	return map[string]bridge.Plan{"api.Echo": {Inputs: map[string]bridge.Source{"Value": source}, OutputField: "Value"}}
}
func TestExtendedTextSelfPreflightRegardlessSelector(t *testing.T) {
	literal := sdl.Text("literal")
	for _, args := range []string{`,multiline=false`, `,readOnly=false`, `,placeholder=""`, `,required=false`} {
		for _, source := range []bridge.Source{{EventField: bridge.ControlText}, {Event: true}, {Widget: "page/value"}, {Literal: &literal}, {Context: "Text"}} {
			d, s := textSession(t, args, "other")
			calls := 0
			e := scalarEngine(t, sp.TextType, func(_ context.Context, r sdl.Record) (sdl.Record, error) { calls++; return r, nil })
			_, err := bridge.Bind(context.Background(), s, d, map[string]*sdl.Engine{"api": e}, textPlan(source), map[string]sdl.Value{"Text": literal})
			if err == nil || !strings.Contains(err.Error(), "self") {
				t.Fatal(args, source, err)
			}
			w, _ := s.Widget("page/value")
			if calls != 0 || s.HasBinding(w.Handle) {
				t.Fatal("partial install/domain execution")
			}
		}
	}
}
func TestExtendedTextExactEchoAndNoReplay(t *testing.T) {
	for _, source := range []bridge.Source{{EventField: bridge.ControlText}, {Event: true}, {Widget: "page/value"}} {
		for _, value := range []string{"", "Blåbær 日本語 🙂 é", "A\r\nB\nC"} {
			d, s := textSession(t, `,multiline=true`, "value")
			calls := 0
			e := scalarEngine(t, sp.TextType, func(_ context.Context, r sdl.Record) (sdl.Record, error) {
				calls++
				if r["Value"].Text != value {
					t.Fatal(r)
				}
				return r, nil
			})
			if _, err := bridge.Bind(context.Background(), s, d, map[string]*sdl.Engine{"api": e}, textPlan(source), nil); err != nil {
				t.Fatal(err)
			}
			event := proposeScalar(t, s, ui.Text(value))
			if event.Control == nil || event.Control.RawDraft != nil || event.Control.Option != nil {
				t.Fatal(event)
			}
			if err := s.Dispatch(event); err != nil {
				t.Fatal(err)
			}
			w, _ := s.Widget("page/value")
			if w.Value != value || w.Draft != value || calls != 1 {
				t.Fatal(w, calls)
			}
			if err := s.Dispatch(event); err == nil || calls != 1 {
				t.Fatal("replay", err, calls)
			}
		}
	}
}
func TestExtendedTextFailedResultsDoNotReplayOrAccept(t *testing.T) {
	for _, mode := range []string{"error", "malformed", "non-echo", "draft-conflict", "state-conflict"} {
		t.Run(mode, func(t *testing.T) {
			d, s := textSession(t, `,multiline=false`, "value")
			calls := 0
			e := scalarEngine(t, sp.TextType, func(_ context.Context, r sdl.Record) (sdl.Record, error) {
				calls++
				switch mode {
				case "error":
					return nil, fmt.Errorf("domain failed")
				case "malformed":
					return sdl.Record{"Value": sdl.Boolean(true)}, nil
				case "non-echo":
					return sdl.Record{"Value": sdl.Text("normalized")}, nil
				case "draft-conflict":
					w, _ := s.Widget("page/value")
					if _, err := s.EditField(w.Handle, s.Revision, ui.Text("newer")); err != nil {
						t.Fatal(err)
					}
				case "state-conflict":
					w, _ := s.Widget("page/other")
					if err := s.Apply(s.Revision, s.BatchRevision+1, []ui.Update{{Handle: w.Handle, Property: ui.Label, Value: ui.Text("Changed")}}); err != nil {
						t.Fatal(err)
					}
				}
				return r, nil
			})
			if _, err := bridge.Bind(context.Background(), s, d, map[string]*sdl.Engine{"api": e}, textPlan(bridge.Source{Event: true}), nil); err != nil {
				t.Fatal(err)
			}
			event := proposeScalar(t, s, ui.Text("proposal"))
			if err := s.Dispatch(event); err == nil {
				t.Fatal("accepted failed result")
			}
			w, _ := s.Widget("page/value")
			if w.Value != "Initial" || calls != 1 {
				t.Fatal(w, calls)
			}
			if mode == "draft-conflict" && w.Draft != "newer" {
				t.Fatal("lost newer draft", w)
			}
			if err := s.Dispatch(event); err == nil || calls != 1 {
				t.Fatal("replay", err, calls)
			}
		})
	}
}
func TestBasicTextNonselfAndLoadExtendedReadOnly(t *testing.T) {
	for _, load := range []bool{false, true} {
		d, s := textSession(t, "", "other")
		if load {
			var err error
			d, roots, e := up.Compile(`sdui 0.3;ref: api "api.sdl";page=[value=button("Load",callback=api.Echo.@invoke);other=input("Other",readOnly=true,multiline=true)];api.Echo.setHandle(page.other);`)
			if e != nil {
				t.Fatal(e)
			}
			s, err = ui.New("load", roots["page"])
			if err != nil {
				t.Fatal(err)
			}
			defer s.Close()
			literal := sdl.Text("load")
			engine := scalarEngine(t, sp.TextType, func(_ context.Context, r sdl.Record) (sdl.Record, error) {
				return sdl.Record{"Value": sdl.Text("loaded\n日本語")}, nil
			})
			if _, err = bridge.Bind(context.Background(), s, d, map[string]*sdl.Engine{"api": engine}, textPlan(bridge.Source{Literal: &literal}), nil); err != nil {
				t.Fatal(err)
			}
			w, _ := s.Widget("page/value")
			if err = s.Dispatch(ui.Event{Kind: ui.Activate, Handle: w.Handle, ModelRevision: s.Revision, Sequence: 1}); err != nil {
				t.Fatal(err)
			}
			w, _ = s.Widget("page/other")
			if w.Value != "loaded\n日本語" {
				t.Fatal(w)
			}
		} else {
			engine := scalarEngine(t, sp.TextType, func(_ context.Context, r sdl.Record) (sdl.Record, error) { return r, nil })
			if _, err := bridge.Bind(context.Background(), s, d, map[string]*sdl.Engine{"api": engine}, textPlan(bridge.Source{EventField: bridge.ControlText}), nil); err != nil {
				t.Fatal(err)
			}
			w, _ := s.Widget("page/value")
			if err := s.Draft(w.Handle, "legacy"); err != nil {
				t.Fatal(err)
			}
			w, _ = s.Widget("page/value")
			if err := s.Dispatch(ui.Event{Kind: ui.Commit, Handle: w.Handle, ModelRevision: s.Revision, Sequence: 1, DraftRevision: w.DraftRevision, Value: ui.Text(w.Draft)}); err != nil {
				t.Fatal(err)
			}
			w, _ = s.Widget("page/other")
			if w.Value != "legacy" {
				t.Fatal(w)
			}
		}
	}
}
