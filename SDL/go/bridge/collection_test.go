package bridge_test

import (
	"context"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/bridge"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/examples/collections"
	sdlparser "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	sdl "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/runtime"
	"strings"
	"testing"
)

func TestCollectionEventFieldAndCapturedReceiver(t *testing.T) {
	for _, mode := range []string{"valid", "selector", "integer", "mixed", "draft-during-action", "replace-receiver"} {
		t.Run(mode, func(t *testing.T) {
			d, roots, e := parser.Compile(collections.Source)
			if e != nil {
				t.Fatal(e)
			}
			s, e := ui.New("collection", roots["page"])
			if e != nil {
				t.Fatal(e)
			}
			defer s.Close()
			fixture, e := collections.New()
			if e != nil {
				t.Fatal(e)
			}
			defer fixture.Close()
			if e = s.BindProviders(fixture.Providers(false)); e != nil {
				t.Fatal(e)
			}
			actions := collections.Actions
			if mode == "integer" {
				actions = strings.Replace(actions, "ItemId as text", "ItemId as integer", 1)
			}
			program, e := sdlparser.CompileActions(actions)
			if e != nil {
				t.Fatal(e)
			}
			kind := sdlparser.TextType
			if mode == "integer" {
				kind = sdlparser.IntegerType
			}
			calls := 0
			engine, e := sdl.New(program, sdl.Registry{"GoActivate": {Input: sdlparser.RecordType{"ItemId": kind}, Output: sdlparser.RecordType{"Preview": sdlparser.TextType}, Call: func(context.Context, sdl.Record) (sdl.Record, error) {
				calls++
				receiver, _ := s.Widget("page/footer/preview")
				if mode == "draft-during-action" {
					if e := s.Draft(receiver.Handle, "user edit"); e != nil {
						t.Fatal(e)
					}
				}
				if mode == "replace-receiver" {
					_, other, e := parser.Compile(strings.Replace(collections.Source, `preview=input("Preview",value="No activation")`, `preview=button("Changed")`, 1))
					if e != nil {
						t.Fatal(e)
					}
					if e = s.Reload(other["page"]); e != nil {
						t.Fatal(e)
					}
					if e = s.Reload(roots["page"]); e != nil {
						t.Fatal(e)
					}
				}
				return sdl.Record{"Preview": sdl.Text("result")}, nil
			}}})
			if e != nil {
				t.Fatal(e)
			}
			source := bridge.Source{EventField: bridge.CollectionItemID}
			if mode == "selector" {
				source.EventField = "unknown"
			}
			if mode == "mixed" {
				source.Widget = "page/footer/preview"
			}
			_, e = bridge.Bind(context.Background(), s, d, map[string]*sdl.Engine{"nav": engine}, map[string]bridge.Plan{"nav.Activate": {Inputs: map[string]bridge.Source{"ItemId": source}, OutputField: "Preview"}}, nil)
			if mode == "selector" || mode == "integer" || mode == "mixed" {
				if e == nil || calls != 0 {
					t.Fatal("invalid bridge admitted", e, calls)
				}
				return
			}
			if e != nil {
				t.Fatal(e)
			}
			if calls != 0 {
				t.Fatal("preparation action")
			}
			w, _ := s.Widget("page/body/entries")
			target, e := s.Target(w.Handle, "entry-001")
			if e != nil {
				t.Fatal(e)
			}
			e = s.Dispatch(ui.Event{Handle: w.Handle, ModelRevision: s.Revision, Sequence: 1, Kind: ui.Activate, Collection: &target})
			if mode == "valid" {
				if e != nil {
					t.Fatal(e)
				}
			} else {
				if e == nil {
					t.Fatal("changed receiver accepted")
				}
				r, _ := s.Widget("page/footer/preview")
				if r.Value == "result" {
					t.Fatal("result retargeted")
				}
			}
			if calls != 1 {
				t.Fatal("action replayed", calls)
			}
		})
	}
}
