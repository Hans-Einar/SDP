package bridge_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	uiparser "github.com/Hans-Einar/SDP/SDUI/go/parser"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/bridge"
	sdlparser "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	sdl "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/runtime"
)

const paneActions = `language action-core version 0.1.
action Page.
record PageInput.
record PageOutput.
Page invokes GoPage.
Page returns PageOutput.
Page takes PageInput.
PageInput field PageId as text.
PageInput field PreviousId as text.
PageOutput field Preview as text.
`
const paneSource = `sdui 0.3;
ref: api "unused.sdl";
page=[
 tabs=tabs("Pages",callback=api.Page.@invoke)[
  first=page("First")[edit=input("Draft",value="accepted")];
  second=page("Second")["Other content"]
 ];
 preview=input("Result",value="initial")
];
api.Page.setHandle(page.preview);
`

func panePlans() map[string]bridge.Plan {
	return map[string]bridge.Plan{"api.Page": {Inputs: map[string]bridge.Source{"PageId": {EventField: bridge.TabPageID}, "PreviousId": {EventField: bridge.TabPreviousPageID}}, OutputField: "Preview"}}
}
func paneSession(t *testing.T, source string) (*uiparser.Document, *ui.Session) {
	t.Helper()
	d, roots, err := uiparser.Compile(source)
	if err != nil {
		t.Fatal(err)
	}
	s, err := ui.New("panes-test", roots["page"])
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(s.Close)
	return d, s
}
func paneEngine(t *testing.T, actions string, kind sdlparser.ScalarType, call sdl.Handler) *sdl.Engine {
	t.Helper()
	p, err := sdlparser.CompileActions(actions)
	if err != nil {
		t.Fatal(err)
	}
	e, err := sdl.New(p, sdl.Registry{"GoPage": {Input: sdlparser.RecordType{"PageId": kind, "PreviousId": sdlparser.TextType}, Output: sdlparser.RecordType{"Preview": sdlparser.TextType}, Call: call}})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(e.Close)
	return e
}
func pageEvent(t *testing.T, s *ui.Session, id string) ui.Event {
	t.Helper()
	h, ok := s.Pane("page/tabs")
	if !ok {
		t.Fatal("missing tabs")
	}
	state, ok := s.Tabs(h)
	if !ok {
		t.Fatal("missing state")
	}
	for _, p := range state.Pages {
		if p.ID == id {
			return ui.Event{Kind: ui.ActivatePage, Handle: h, ModelRevision: s.Revision, StateRevision: s.StateRevision, Sequence: s.Sequence() + 1, Page: &ui.PageActivation{PreviousID: state.Selected, PageID: id, Page: p.Handle}}
		}
	}
	t.Fatalf("missing page %s", id)
	return ui.Event{}
}
func TestTabCallbackTypedIDsAndSilentOperations(t *testing.T) {
	d, s := paneSession(t, paneSource)
	var records []sdl.Record
	e := paneEngine(t, paneActions, sdlparser.TextType, func(_ context.Context, r sdl.Record) (sdl.Record, error) {
		records = append(records, r)
		return sdl.Record{"Preview": sdl.Text(r["PreviousId"].Text + " -> " + r["PageId"].Text)}, nil
	})
	b, err := bridge.Bind(context.Background(), s, d, map[string]*sdl.Engine{"api": e}, panePlans(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 0 || len(b.Links) != 1 || b.Links[0].UISpan.Line == 0 || b.Links[0].SDLSpan.Line == 0 {
		t.Fatal("preflight called SDL or lost source provenance", b.Links)
	}
	for _, w := range s.Widgets() {
		if w.Handle.Kind == "tabs" {
			t.Fatal("legacy Widgets expanded")
		}
	}
	edit, _ := s.Widget("page/tabs/first/edit")
	if err = s.Draft(edit.Handle, "unsaved page draft"); err != nil {
		t.Fatal(err)
	}
	event := pageEvent(t, s, "second")
	result, err := s.DispatchInteraction(event)
	if err != nil || result.Domain != ui.DomainSucceeded || result.Status != "committed" {
		t.Fatal(result, err)
	}
	if len(records) != 1 || records[0]["PageId"].Text != "second" || records[0]["PreviousId"].Text != "first" {
		t.Fatal(records)
	}
	resultWidget, _ := s.Widget("page/preview")
	if resultWidget.Value != "first -> second" {
		t.Fatal(resultWidget)
	}
	if _, err = s.DispatchInteraction(pageEvent(t, s, "second")); err != nil {
		t.Fatal(err)
	}
	h, _ := s.Pane("page/tabs")
	if err = s.SelectPage(h, "first"); err != nil {
		t.Fatal(err)
	}
	if len(records) != 1 {
		t.Fatal("same-page/programmatic selection called SDL")
	}
	edit, _ = s.Widget(edit.Handle.Path)
	if edit.Draft != "unsaved page draft" || !edit.Dirty {
		t.Fatal("hidden draft lost", edit)
	}
	// Compatible reload keeps identities and the shared command sequence; binding is fresh.
	_, roots, err := uiparser.Compile(paneSource)
	if err != nil {
		t.Fatal(err)
	}
	if err = s.Reload(roots["page"]); err != nil {
		t.Fatal(err)
	}
	if _, err = bridge.Bind(context.Background(), s, d, map[string]*sdl.Engine{"api": e}, panePlans(), nil); err != nil {
		t.Fatal(err)
	}
	if len(records) != 1 {
		t.Fatal("reload invoked SDL")
	}
	if _, err = s.DispatchInteraction(pageEvent(t, s, "second")); err != nil {
		t.Fatal(err)
	}
	if len(records) != 2 {
		t.Fatal("reload action sequence not continued", records)
	}
}
func TestTabBridgePreflightRejectsInvalidSelectorsWithoutInstallation(t *testing.T) {
	for _, mode := range []string{"unknown", "collection", "integer", "mixed", "legacy-event", "module", "result", "hidden-module", "wrong-owner"} {
		t.Run(mode, func(t *testing.T) {
			text := paneSource
			if mode == "hidden-module" {
				text = strings.Replace(text, "\n ];", "\n ] {visible=false};", 1)
			}
			if mode == "wrong-owner" {
				text = strings.Replace(text, `preview=input("Result",value="initial")`, `button=button("Bad owner",callback=api.Page.@invoke); preview=input("Result",value="initial")`, 1)
			}
			d, s := paneSession(t, text)
			calls := 0
			actions := paneActions
			kind := sdlparser.TextType
			if mode == "integer" {
				actions = strings.Replace(actions, "PageId as text", "PageId as integer", 1)
				kind = sdlparser.IntegerType
			}
			e := paneEngine(t, actions, kind, func(context.Context, sdl.Record) (sdl.Record, error) {
				calls++
				return sdl.Record{"Preview": sdl.Text("bad")}, nil
			})
			plans := panePlans()
			p := plans["api.Page"]
			src := p.Inputs["PageId"]
			modules := map[string]*sdl.Engine{"api": e}
			switch mode {
			case "unknown":
				src.EventField = "tab.unknown"
			case "collection":
				src.EventField = bridge.CollectionItemID
			case "mixed":
				src.Widget = "page/preview"
			case "legacy-event":
				src = bridge.Source{Event: true}
			case "module", "hidden-module":
				modules = nil
			case "result":
				p.OutputField = "Missing"
			}
			p.Inputs["PageId"] = src
			plans["api.Page"] = p
			if _, err := bridge.Bind(context.Background(), s, d, modules, plans, nil); err == nil {
				t.Fatal("invalid preflight accepted")
			}
			if calls != 0 {
				t.Fatal("preflight called SDL")
			}
			for _, owner := range s.CallbackOwners() {
				if s.HasInteractionBinding(owner.Handle) {
					t.Fatal("partial tabs binding")
				}
			}
			for _, w := range s.Widgets() {
				if s.HasBinding(w.Handle) {
					t.Fatal("partial legacy binding")
				}
			}
		})
	}
}
func TestTabOutcomeAndReceiverDraftConflict(t *testing.T) {
	for _, mode := range []string{"error", "malformed", "draft-conflict", "geometry-conflict"} {
		t.Run(mode, func(t *testing.T) {
			d, s := paneSession(t, paneSource)
			calls := 0
			e := paneEngine(t, paneActions, sdlparser.TextType, func(context.Context, sdl.Record) (sdl.Record, error) {
				calls++
				switch mode {
				case "error":
					return nil, fmt.Errorf("domain failure")
				case "malformed":
					return sdl.Record{"Preview": sdl.Integer(2)}, nil
				case "draft-conflict":
					w, _ := s.Widget("page/preview")
					if err := s.Draft(w.Handle, "newer draft"); err != nil {
						t.Fatal(err)
					}
				}
				return sdl.Record{"Preview": sdl.Text("domain succeeded")}, nil
			})
			if _, err := bridge.Bind(context.Background(), s, d, map[string]*sdl.Engine{"api": e}, panePlans(), nil); err != nil {
				t.Fatal(err)
			}
			if mode == "geometry-conflict" {
				if err := s.CheckPresentationWith(func(snapshot ui.Snapshot) (ui.PresentationState, error) {
					var value string
					snapshot.Root.Walk(func(n *uiparser.Instance) {
						if n.Path == "page/preview" {
							if v, ok := n.Arguments["value"].(uiparser.Literal); ok {
								value, _ = v.Value.(string)
							}
						}
					})
					if value == "domain succeeded" {
						return ui.PresentationState{}, fmt.Errorf("injected final geometry failure")
					}
					return ui.PresentationState{Viewports: snapshot.Viewports, Splits: map[string]ui.SplitGeometry{}}, nil
				}); err != nil {
					t.Fatal(err)
				}
			}
			event := pageEvent(t, s, "second")
			result, err := s.DispatchInteraction(event)
			if err == nil || calls != 1 {
				t.Fatal("failure not observed once", result, err, calls)
			}
			expected := ui.DomainUnknown
			if mode == "draft-conflict" || mode == "geometry-conflict" {
				expected = ui.DomainSucceeded
			}
			if result.Domain != expected {
				t.Fatal("domain outcome lost", result)
			}
			h, _ := s.Pane("page/tabs")
			tabs, _ := s.Tabs(h)
			w, _ := s.Widget("page/preview")
			if tabs.Selected != "first" || w.Value != "initial" {
				t.Fatal("partial UI publication", tabs, w)
			}
			if mode == "draft-conflict" && w.Draft != "newer draft" {
				t.Fatal("independent draft rolled back")
			}
			if _, err = s.DispatchInteraction(event); err == nil || calls != 1 {
				t.Fatal("consumed callback replayed", err, calls)
			}
		})
	}
}

func TestReusedTabsHaveIndependentOwnersAndFailedRebindPreservesHandlers(t *testing.T) {
	source := `sdui 0.3;
ref: api "unused.sdl";
Pane=[tabs=tabs("Shared",callback=api.Page.@invoke)[first=page("First")[]; second=page("Second")[]]];
page=[left=Pane; right=Pane; preview=input("Result",value="initial")];
api.Page.setHandle(page.preview);
`
	d, s := paneSession(t, source)
	calls := 0
	engine := paneEngine(t, paneActions, sdlparser.TextType, func(_ context.Context, r sdl.Record) (sdl.Record, error) {
		calls++
		if r["PreviousId"].Text != "first" || r["PageId"].Text != "second" {
			t.Fatal("bad reused-page identity", r)
		}
		return sdl.Record{"Preview": sdl.Text("reused")}, nil
	})
	b, err := bridge.Bind(context.Background(), s, d, map[string]*sdl.Engine{"api": engine}, panePlans(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(b.Links) != 2 || len(s.CallbackOwners()) != 2 {
		t.Fatal("missing reused composition owners", b.Links)
	}
	bad := panePlans()["api.Page"]
	bad.Inputs["PageId"] = bridge.Source{EventField: "unknown"}
	b.Plans["api.Page"] = bad
	if err = b.Rebind(d); err == nil {
		t.Fatal("invalid rebind accepted")
	}
	for i, path := range []string{"page/left/tabs", "page/right/tabs"} {
		h, ok := s.Pane(path)
		if !ok {
			t.Fatal(path)
		}
		tabs, _ := s.Tabs(h)
		if tabs.Selected != "first" {
			t.Fatal("reuse shared selection", tabs)
		}
		p := tabs.Pages[1]
		if _, err = s.DispatchInteraction(ui.Event{Kind: ui.ActivatePage, Handle: h, ModelRevision: s.Revision, StateRevision: s.StateRevision, Sequence: s.Sequence() + 1, Page: &ui.PageActivation{PreviousID: tabs.Selected, PageID: p.ID, Page: p.Handle}}); err != nil {
			t.Fatal("prior handler lost after failed rebind", err)
		}
		if calls != i+1 {
			t.Fatal("reused callback called wrong count", calls)
		}
	}
}
