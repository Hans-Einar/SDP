package bridge_test

import (
	"context"
	"strings"
	"testing"

	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/bridge"
	fixture "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/examples/commands"
	sp "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	sdl "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/runtime"
)

func TestCommandTextResultExtendedReceiver(t *testing.T) {
	for _, policy := range []string{"", ",readOnly=true", ",multiline=false"} {
		t.Run(policy, func(t *testing.T) {
			source := strings.Replace(fixture.Source, `preview=input("Run result",value="No command")`, `preview=input("Run result",value="No command"`+policy+`)`, 1)
			d, s, f := commandSession(t, source)
			bindCommands(t, d, s, f)
			result, err := invoke(t, s, "page/view/header/toolbar/runButton", "button")
			w, _ := s.Widget(fixture.PreviewPath)
			if err != nil || result.Domain != ui.DomainSucceeded || result.Status != "committed" || w.Value != "Run completed" || w.Draft != w.Value || f.Calls()["Run"] != 1 {
				t.Fatal(result, err, w, f.Calls())
			}
		})
	}
}

func TestTabTextResultExtendedReceiver(t *testing.T) {
	for _, conflict := range []bool{false, true} {
		name := "readonly-success"
		policy := ",readOnly=true"
		if conflict {
			name = "writable-draft-conflict"
			policy = ",multiline=false"
		}
		t.Run(name, func(t *testing.T) {
			d, s := paneSession(t, strings.Replace(paneSource, `preview=input("Result",value="initial")`, `preview=input("Result",value="initial"`+policy+`)`, 1))
			calls := 0
			e := paneEngine(t, paneActions, sp.TextType, func(_ context.Context, r sdl.Record) (sdl.Record, error) {
				calls++
				if conflict {
					w, _ := s.Widget("page/preview")
					if err := s.Draft(w.Handle, "newer draft"); err != nil {
						return nil, err
					}
				}
				return sdl.Record{"Preview": sdl.Text("loaded")}, nil
			})
			if _, err := bridge.Bind(context.Background(), s, d, map[string]*sdl.Engine{"api": e}, panePlans(), nil); err != nil {
				t.Fatal(err)
			}
			result, err := s.DispatchInteraction(pageEvent(t, s, "second"))
			w, _ := s.Widget("page/preview")
			h, _ := s.Pane("page/tabs")
			tabs, _ := s.Tabs(h)
			if calls != 1 || result.Domain != ui.DomainSucceeded {
				t.Fatal(result, err, calls)
			}
			if conflict {
				if err == nil || !strings.Contains(err.Error(), "stale-result-target") || w.Value != "initial" || w.Draft != "newer draft" || tabs.Selected != "first" {
					t.Fatal(result, err, w, tabs)
				}
			} else if err != nil || result.Status != "committed" || w.Value != "loaded" || tabs.Selected != "second" {
				t.Fatal(result, err, w, tabs)
			}
		})
	}
}

func TestCommandExtendedReceiverConflictIsNotReplayed(t *testing.T) {
	source := strings.Replace(fixture.Source, `preview=input("Run result",value="No command")`, `preview=input("Run result",value="No command",multiline=false)`, 1)
	d, s, f := commandSession(t, source)
	bindCommands(t, d, s, f)
	if err := f.NextAction("Run", "draft-conflict"); err != nil {
		t.Fatal(err)
	}
	result, err := invoke(t, s, "page/view/header/toolbar/runButton", "button")
	w, _ := s.Widget(fixture.PreviewPath)
	if err == nil || !strings.Contains(err.Error(), "stale-result-target") || result.Domain != ui.DomainSucceeded || w.Value != "No command" || w.Draft != "newer draft during SDL" || f.Calls()["Run"] != 1 {
		t.Fatal(result, err, w, f.Calls())
	}
}
