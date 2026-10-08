package bridge_test

import (
	"context"
	"fyne.io/fyne/v2/test"
	"reflect"
	"strings"
	"testing"

	uiparser "github.com/Hans-Einar/SDP/SDUI/go/parser"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/bridge"
	fixture "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/examples/commands"
	sdl "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/runtime"
)

func commandSession(t *testing.T, source string) (*uiparser.Document, *ui.Session, *fixture.Fixture) {
	t.Helper()
	d, roots, err := uiparser.Compile(source)
	if err != nil {
		t.Fatal(err)
	}
	s, err := ui.New("commands-test", roots["page"])
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(s.Close)
	f, err := fixture.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(f.Close)
	if err = s.BindProviders(f.Providers()); err != nil {
		t.Fatal(err)
	}
	f.Conflict = func(path string) error {
		w, ok := s.Widget(path)
		if !ok {
			t.Fatalf("missing conflict input %s", path)
		}
		return s.Draft(w.Handle, "newer draft during SDL")
	}
	return d, s, f
}
func bindCommands(t *testing.T, d *uiparser.Document, s *ui.Session, f *fixture.Fixture) {
	t.Helper()
	if _, err := bridge.Bind(context.Background(), s, d, map[string]*sdl.Engine{"actions": f.Engine}, fixture.Plans(), nil); err != nil {
		t.Fatal(err)
	}
	if len(f.Calls()) != 0 {
		t.Fatal("preflight called SDL")
	}
}
func origin(t *testing.T, s *ui.Session, path string) ui.Handle {
	t.Helper()
	if p, ok := s.Snapshot().Presentations[path]; ok {
		return p.Handle
	}
	if h, ok := s.Command(path); ok {
		return h
	}
	t.Fatalf("missing origin %s", path)
	return ui.Handle{}
}
func invoke(t *testing.T, s *ui.Session, path, via string) (ui.InteractionResult, error) {
	t.Helper()
	e, err := s.CaptureCommand(origin(t, s, path), via, nil)
	if err != nil {
		return ui.InteractionResult{}, err
	}
	return s.DispatchInteraction(e)
}
func openDialog(t *testing.T, s *ui.Session) ui.SurfaceTarget {
	t.Helper()
	h, ok := s.Surface(fixture.SettingsPath)
	if !ok {
		t.Fatal("missing surface")
	}
	token, err := s.OpenSurface(h, ui.ContextTarget{})
	if err != nil {
		t.Fatal(err)
	}
	// State-level test explicitly acknowledges publication; this is not native evidence.
	if err = s.ConfirmSurfacePublication(token); err != nil {
		t.Fatal(err)
	}
	return token
}
func dialogEvent(t *testing.T, s *ui.Session, token ui.SurfaceTarget, kind ui.EventKind) ui.Event {
	t.Helper()
	e, err := s.CaptureDialog(token, kind)
	if err != nil {
		t.Fatal(err)
	}
	return e
}

func TestSharedCommandsUseOneActualSDLHandlerAndTypedValues(t *testing.T) {
	d, s, f := commandSession(t, fixture.Source)
	bindCommands(t, d, s, f)
	var toggleInput sdl.Value
	f.Log = func(kind string, v any) {
		if kind == "action" {
			r := v.(map[string]any)
			if r["action"] == "Toggle" {
				toggleInput = r["input"].(sdl.Record)["Checked"]
			}
		}
	}
	for _, route := range []struct{ path, via string }{{"page/view/header/toolbar/runButton", "button"}, {"page/view/header/actions/runItem", "menu"}, {"page/run", "key"}} {
		if route.via == "menu" {
			m, _ := s.Menu("page/view/header/actions")
			if _, err := s.OpenMenu(m, ui.ContextTarget{}); err != nil {
				t.Fatal(err)
			}
		}
		result, err := invoke(t, s, route.path, route.via)
		if err != nil || result.Domain != ui.DomainSucceeded {
			t.Fatal(route, result, err)
		}
	}
	if f.Calls()["Run"] != 3 {
		t.Fatal("shared handler duplicated/missing", f.Calls())
	}
	if result, err := invoke(t, s, "page/view/header/toolbar/flagButton", "button"); err != nil || result.Domain != ui.DomainSucceeded {
		t.Fatal(result, err)
	}
	h, _ := s.Command("page/flag")
	flag, _ := s.CommandState(h)
	if !flag.Checked || toggleInput != sdl.Boolean(true) || f.Calls()["Toggle"] != 1 {
		t.Fatal("proposed boolean was not shared", flag, toggleInput, f.Calls())
	}
}
func TestCommandContextUsesCapturedRowAndRejectsStaleMenu(t *testing.T) {
	d, s, f := commandSession(t, fixture.Source)
	bindCommands(t, d, s, f)
	w, _ := s.Widget(fixture.ItemsPath)
	a, _ := s.Target(w.Handle, "alpha")
	b, _ := s.Target(w.Handle, "beta")
	if err := s.Dispatch(ui.Event{Handle: w.Handle, ModelRevision: s.Revision, Sequence: s.Sequence() + 1, Kind: ui.Select, Collection: &a}); err != nil {
		t.Fatal(err)
	}
	menu, _ := s.Menu("page/itemMenu")
	context := ui.ContextTarget{Widget: w.Handle, ModelRevision: s.Revision, Item: &b}
	if _, err := s.OpenMenu(menu, context); err != nil {
		t.Fatal(err)
	}
	result, err := invoke(t, s, "page/itemMenu/inspectItem", "menu")
	if err != nil || result.Domain != ui.DomainSucceeded {
		t.Fatal(result, err)
	}
	preview, _ := s.Widget(fixture.ItemPreviewPath)
	if preview.Value != "Item: beta" {
		t.Fatal("retargeted current selection", preview.Value)
	}
	if _, err = s.OpenMenu(menu, context); err != nil {
		t.Fatal(err)
	}
	captured, err := s.CaptureCommand(origin(t, s, "page/itemMenu/inspectItem"), "menu", nil)
	if err != nil {
		t.Fatal(err)
	}
	target, _ := s.Target(w.Handle, "")
	if err = s.ReplaceCollection(target, f.Providers()[fixture.ItemsPath].Initial); err != nil {
		t.Fatal(err)
	}
	if _, err = s.DispatchInteraction(captured); err == nil || f.Calls()["Inspect"] != 1 {
		t.Fatal("stale context executed SDL", err, f.Calls())
	}
}
func TestDialogTypedAcceptCapturesFieldsAndPublishesOneResult(t *testing.T) {
	d, s, f := commandSession(t, fixture.Source)
	bindCommands(t, d, s, f)
	token := openDialog(t, s)
	name, _ := s.Widget(fixture.NamePath)
	if err := s.Draft(name.Handle, "captured name"); err != nil {
		t.Fatal(err)
	}
	event := dialogEvent(t, s, token, ui.Accept)
	result, err := s.DispatchInteraction(event)
	if err != nil || result.Domain != ui.DomainSucceeded {
		t.Fatal(result, err)
	}
	if f.Domain().Name != "captured name" || f.Domain().Commits != 1 {
		t.Fatal("SDL did not receive capture", f.Domain())
	}
	name, _ = s.Widget(fixture.NamePath)
	if name.Value != "captured name" || name.Dirty {
		t.Fatal("capture not accepted", name)
	}
	receipts := s.DrainDialogResults()
	if len(receipts) != 1 || receipts[0].Kind != "accept" || receipts[0].Domain != ui.DomainSucceeded || len(receipts[0].Fields) != 2 || receipts[0].Surface != token {
		t.Fatal(receipts)
	}
	if len(s.DrainDialogResults()) != 0 {
		t.Fatal("terminal result repeated")
	}
	if _, err = s.DispatchInteraction(event); err == nil || f.Calls()["Save"] != 1 {
		t.Fatal("Accept replayed", err, f.Calls())
	}
}
func TestDialogFalseUnknownAndPostDomainConflictRemainTruthful(t *testing.T) {
	for _, mode := range []string{"false", "error", "malformed", "draft-conflict", "resource-conflict"} {
		t.Run(mode, func(t *testing.T) {
			d, s, f := commandSession(t, fixture.Source)
			bindCommands(t, d, s, f)
			if mode == "resource-conflict" {
				test.NewTempApp(t)
				request, err := f.Request(fixture.Source, 1)
				if err != nil {
					t.Fatal(err)
				}
				if err = s.PreparePresentationWith(func(snapshot ui.Snapshot) (ui.PresentationTicket, error) {
					return ui.PresentationTicket{}, request.PrepareResources(snapshot)
				}); err != nil {
					t.Fatal(err)
				}
			}
			token := openDialog(t, s)
			name, _ := s.Widget(fixture.NamePath)
			if err := s.Draft(name.Handle, "captured"); err != nil {
				t.Fatal(err)
			}
			if err := f.NextAction("Save", mode); err != nil {
				t.Fatal(err)
			}
			event := dialogEvent(t, s, token, ui.Accept)
			result, err := s.DispatchInteraction(event)
			state, ok := s.SurfaceState(token.Handle)
			if !ok || !state.Open || len(s.DrainDialogResults()) != 0 {
				t.Fatal("rejected/failed Accept closed", state)
			}
			expected := ui.DomainUnknown
			if mode == "false" {
				expected = ui.DomainRejected
			}
			if mode == "draft-conflict" || mode == "resource-conflict" {
				expected = ui.DomainSucceeded
			}
			if result.Domain != expected || f.Calls()["Save"] != 1 {
				t.Fatal("outcome lost", result, err)
			}
			if mode == "false" {
				if state.AcceptBlocked || f.Domain().Commits != 0 {
					t.Fatal("false decision blocked recovery/committed", state)
				}
				if _, err = s.DispatchInteraction(dialogEvent(t, s, token, ui.Accept)); err != nil || f.Calls()["Save"] != 2 {
					t.Fatal("fresh true retry after false", err)
				}
				return
			}
			if err == nil || !state.AcceptBlocked {
				t.Fatal("uncertain/succeeded conflict not blocked", result, err, state)
			}
			if again, captureErr := s.CaptureDialog(token, ui.Accept); captureErr == nil {
				if _, dispatchErr := s.DispatchInteraction(again); dispatchErr == nil {
					t.Fatal("blocked opening accepted again")
				}
			}
			if f.Calls()["Save"] != 1 {
				t.Fatal("blocked attempt replayed")
			}
			// These new revisions must not make token-based Cancel impossible.
			name, _ = s.Widget(fixture.NamePath)
			if err = s.Draft(name.Handle, "edit after conflict"); err != nil {
				t.Fatal(err)
			}
			if _, err = s.DispatchInteraction(dialogEvent(t, s, token, ui.Cancel)); err != nil {
				t.Fatal("Cancel used obsolete capture revisions", err)
			}
			receipts := s.DrainDialogResults()
			if len(receipts) != 1 || receipts[0].Kind != "cancel" || receipts[0].Domain != expected || receipts[0].AcceptSequence != event.Sequence {
				t.Fatal("closing hid domain outcome", receipts)
			}
			if expected == ui.DomainSucceeded && f.Domain().Commits != 1 {
				t.Fatal("domain save lost", f.Domain())
			}
		})
	}
}
func TestM2BindingPlansRejectWithoutPartialInstallation(t *testing.T) {
	modes := []string{"unknown-mode", "dialog-text", "text-accept-fields", "accept-on-command", "accept-field", "message-field", "output", "revision", "setHandle", "path-empty", "path-parent", "path-missing", "live-field", "path-on-literal", "checked-on-run", "context-on-run", "capture-on-command"}
	for _, mode := range modes {
		t.Run(mode, func(t *testing.T) {
			source := fixture.Source
			if mode == "setHandle" {
				source += "\nactions.Save.setHandle(page.view.body.mainDraft);\n"
			}
			d, s, f := commandSession(t, source)
			plans := fixture.Plans()
			save := plans["actions.Save"]
			run := plans["actions.Run"]
			switch mode {
			case "unknown-mode":
				save.ResultMode = "unknown"
			case "dialog-text":
				save.ResultMode = bridge.TextResult
			case "text-accept-fields":
				run.AcceptField = "Accepted"
			case "accept-on-command":
				run.ResultMode = bridge.DialogAcceptResult
				run.AcceptField = "Accepted"
				run.MessageField = "Message"
			case "accept-field":
				save.AcceptField = "Missing"
			case "message-field":
				save.MessageField = "Accepted"
			case "output":
				save.OutputField = "Preview"
			case "revision":
				save.RevisionField = "Revision"
				save.RevisionContext = "Rev"
			case "path-empty":
				save.Inputs["Name"] = bridge.Source{EventField: bridge.DialogFieldValue}
			case "path-parent":
				save.Inputs["Name"] = bridge.Source{EventField: bridge.DialogFieldValue, FieldPath: "../name"}
			case "path-missing":
				save.Inputs["Name"] = bridge.Source{EventField: bridge.DialogFieldValue, FieldPath: "missing"}
			case "live-field":
				save.Inputs["Name"] = bridge.Source{Widget: fixture.NamePath}
			case "path-on-literal":
				src := run.Inputs["Token"]
				src.FieldPath = "name"
				run.Inputs["Token"] = src
			case "checked-on-run":
				run.Inputs["Token"] = bridge.Source{EventField: bridge.CommandChecked}
			case "context-on-run":
				run.Inputs["Token"] = bridge.Source{EventField: bridge.CommandContextItemID}
			case "capture-on-command":
				run.Inputs["Token"] = bridge.Source{EventField: bridge.DialogFieldValue, FieldPath: "form/body/name"}
			}
			plans["actions.Save"] = save
			plans["actions.Run"] = run
			if _, err := bridge.Bind(context.Background(), s, d, map[string]*sdl.Engine{"actions": f.Engine}, plans, nil); err == nil {
				t.Fatal("invalid plan accepted", mode)
			}
			if len(f.Calls()) != 0 {
				t.Fatal("preflight called SDL")
			}
			for _, owner := range s.CallbackOwners() {
				if s.HasInteractionBinding(owner.Handle) {
					t.Fatal("partial interaction binding", owner)
				}
			}
		})
	}
}
func TestDialogCapturePlanRejectsNestedFieldAndPreservesCallerPlans(t *testing.T) {
	source := strings.Replace(fixture.Source, `form=[`, `inner=dialog("Inner")[secret=input("Secret",value="")]; form=[`, 1)
	d, s, f := commandSession(t, source)
	plans := fixture.Plans()
	plans["actions.Save"].Inputs["Name"] = bridge.Source{EventField: bridge.DialogFieldValue, FieldPath: "inner/secret"}
	if _, err := bridge.Bind(context.Background(), s, d, map[string]*sdl.Engine{"actions": f.Engine}, plans, nil); err == nil {
		t.Fatal("nested-dialog field was accepted")
	}
	// Preflight resolution must not write opaque handles back into the caller's plan.
	plans = fixture.Plans()
	before := fixture.Plans()
	if _, err := bridge.Bind(context.Background(), s, d, map[string]*sdl.Engine{"actions": f.Engine}, plans, nil); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(plans, before) {
		t.Fatal("preflight mutated source plans")
	}
}

func TestEachCommandConflictPreservesItsOwnReceiverAndDomainOutcome(t *testing.T) {
	for _, tc := range []struct{ name, path, receiver string }{{"Run", "page/run", fixture.PreviewPath}, {"Toggle", "page/flag", fixture.FlagPreviewPath}, {"Inspect", "page/inspect", fixture.ItemPreviewPath}} {
		t.Run(tc.name, func(t *testing.T) {
			d, s, f := commandSession(t, fixture.Source)
			bindCommands(t, d, s, f)
			if tc.name == "Inspect" {
				w, _ := s.Widget(fixture.ItemsPath)
				target, _ := s.Target(w.Handle, "beta")
				if err := s.Dispatch(ui.Event{Handle: w.Handle, ModelRevision: s.Revision, Sequence: s.Sequence() + 1, Kind: ui.Select, Collection: &target}); err != nil {
					t.Fatal(err)
				}
			}
			before := map[string]ui.Widget{}
			for _, path := range []string{fixture.PreviewPath, fixture.FlagPreviewPath, fixture.ItemPreviewPath} {
				before[path], _ = s.Widget(path)
			}
			if err := f.NextAction(tc.name, "draft-conflict"); err != nil {
				t.Fatal(err)
			}
			button := map[string]string{"Run": "runButton", "Toggle": "flagButton", "Inspect": "inspectButton"}[tc.name]
			result, err := invoke(t, s, "page/view/header/toolbar/"+button, "button")
			if err == nil || result.Domain != ui.DomainSucceeded || f.Calls()[tc.name] != 1 {
				t.Fatal(result, err, f.Calls())
			}
			for path, old := range before {
				w, _ := s.Widget(path)
				if w.Value != old.Value {
					t.Fatal("failed reply published", path, w)
				}
				if path == tc.receiver {
					if w.Draft != "newer draft during SDL" || !w.Dirty {
						t.Fatal("own conflict missing", w)
					}
				} else if w.Draft != old.Draft {
					t.Fatal("wrong receiver changed", path, w)
				}
			}
			if tc.name == "Toggle" {
				h, _ := s.Command(tc.path)
				state, _ := s.CommandState(h)
				if state.Checked {
					t.Fatal("failed toggle published")
				}
			}
		})
	}
}

func TestImplicitCommandAndLegacyButtonInstallOnlyTheirOwnRoute(t *testing.T) {
	source := strings.Replace(fixture.Source, `flagButton=button(command="flag")`, `flagButton=button("Local flag",toggle=true,callback=actions.Toggle.@invoke),legacyRun=button("Legacy run",callback=actions.Run.@invoke)`, 1)
	source = strings.Replace(source, `toggle=true,callback=actions.Toggle.@invoke,key="Primary+K"`, `toggle=true,key="Primary+K"`, 1)
	d, s, f := commandSession(t, source)
	bindCommands(t, d, s, f)
	implicit, _ := s.Widget("page/view/header/toolbar/flagButton")
	legacy, _ := s.Widget("page/view/header/toolbar/legacyRun")
	if implicit.Binding.Module != "" || !s.HasInteractionBinding(implicit.Handle) || s.HasBinding(implicit.Handle) {
		t.Fatal("promoted button has duplicate/wrong route", implicit)
	}
	if legacy.Binding.Module == "" || !s.HasBinding(legacy.Handle) || s.HasInteractionBinding(legacy.Handle) {
		t.Fatal("legacy button changed route", legacy)
	}
	if _, err := invoke(t, s, implicit.InstancePath, "button"); err != nil {
		t.Fatal(err)
	}
	if err := s.Dispatch(ui.Event{Handle: legacy.Handle, ModelRevision: s.Revision, Sequence: s.Sequence() + 1, Kind: ui.Activate}); err != nil {
		t.Fatal(err)
	}
	if f.Calls()["Toggle"] != 1 || f.Calls()["Run"] != 1 {
		t.Fatal("duplicate/missing domain execution", f.Calls())
	}
}
