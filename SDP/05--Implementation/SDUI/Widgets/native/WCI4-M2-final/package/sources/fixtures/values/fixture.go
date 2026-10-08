package values

import (
	"context"
	"fmt"
	"sync"

	"github.com/Hans-Einar/SDP/SDUI/go/host/fynehost"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/preparation"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/bridge"
	sp "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	sdl "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/runtime"
)

type DomainState struct {
	Scalars       map[string]sdl.Value `json:"scalars"`
	ScalarCommits map[string]int       `json:"scalarCommits"`
	Mixed         map[string]ui.Value  `json:"mixed"`
	MixedCommits  int                  `json:"mixedCommits"`
	Text          string               `json:"text"`
	TextCommits   int                  `json:"textCommits"`
}
type Fixture struct {
	Engine              *sdl.Engine
	mu                  sync.Mutex
	calls, changes      map[string]int
	outcomes, observers map[string]string
	domain              DomainState
	choices             map[string][]ui.ChoiceOption
	closed              bool
	Log                 func(string, any)
	Conflict            func(string) error
	Replace             func(string) error
}

func New() (*Fixture, error) {
	f := &Fixture{calls: map[string]int{}, changes: map[string]int{}, outcomes: map[string]string{}, observers: map[string]string{}, domain: DomainState{Scalars: map[string]sdl.Value{}, ScalarCommits: map[string]int{}, Mixed: map[string]ui.Value{}}, choices: map[string][]ui.ChoiceOption{Targets["Mode"]: OptionSet("initial"), Targets["FormMode"]: OptionSet("initial")}}
	for _, p := range Targets {
		f.changes[p] = 0
	}
	for _, p := range []string{FormNotePath, TextEditPath, TextNamePath, PreviewPath} {
		f.changes[p] = 0
	}
	program, err := sp.CompileActions(Actions)
	if err != nil {
		return nil, err
	}
	f.calls["mixed"] = 0
	registry := sdl.Registry{}
	for name, a := range program.Actions {
		f.calls[name] = 0
		registry[a.GoSymbol] = sdl.Binding{Input: program.Records[a.Input], Output: program.Records[a.Output], Call: f.handler(name)}
	}
	f.Engine, err = sdl.New(program, registry)
	return f, err
}
func (f *Fixture) emit(kind string, v any) {
	if f.Log != nil {
		f.Log(kind, v)
	}
}
func (f *Fixture) take(name string) string {
	f.mu.Lock()
	defer f.mu.Unlock()
	mode := f.outcomes[name]
	delete(f.outcomes, name)
	if mode == "" {
		return "success"
	}
	return mode
}
func actionPath(name string) string {
	switch name {
	case "AcceptFlag":
		return Targets["Flag"]
	case "AcceptLevel":
		return Targets["Level"]
	case "AcceptCount":
		return Targets["Count"]
	case "AcceptMode":
		return Targets["Mode"]
	case "ChildCount":
		return Targets["ChildCount"]
	case "SaveText":
		return TextEditPath
	case "Load":
		return PreviewPath
	case "TextSave":
		return TextNamePath
	}
	return ""
}
func (f *Fixture) handler(name string) sdl.Handler {
	return func(ctx context.Context, r sdl.Record) (sdl.Record, error) {
		mode := f.take(name)
		f.mu.Lock()
		f.calls[name]++
		count := f.calls[name]
		f.mu.Unlock()
		invocation, _ := sdl.CurrentInvocation(ctx)
		f.emit("action", map[string]any{"action": name, "input": r, "sequence": invocation.Sequence, "calls": count, "outcome": mode})
		if mode == "error" {
			return nil, fmt.Errorf("injected %s error", name)
		}
		if mode == "malformed" {
			return sdl.Record{"Wrong": sdl.Text("Malformed result")}, nil
		}
		if name == "TextSave" && mode == "false" {
			return sdl.Record{"Accepted": sdl.Boolean(false), "Message": sdl.Text("Text rejected")}, nil
		}
		v := r["Value"]
		if name == "Load" {
			v = sdl.Text("Loaded fixture text")
		}
		f.mu.Lock()
		if name == "TextSave" {
			f.domain.Text = v.Text
			f.domain.TextCommits++
		} else if name != "Load" {
			f.domain.Scalars[name] = v
			f.domain.ScalarCommits[name]++
		}
		f.mu.Unlock()
		if mode == "draft-conflict" {
			if f.Conflict == nil {
				return nil, fmt.Errorf("missing conflict hook")
			}
			if err := f.Conflict(actionPath(name)); err != nil {
				return nil, err
			}
		}
		if mode == "options-conflict" {
			if f.Replace == nil {
				return nil, fmt.Errorf("missing replacement hook")
			}
			if err := f.Replace(Targets["Mode"]); err != nil {
				return nil, err
			}
		}
		if name == "TextSave" {
			return sdl.Record{"Accepted": sdl.Boolean(true), "Message": sdl.Text("Saved text")}, nil
		}
		return sdl.Record{"Value": v}, nil
	}
}
func (f *Fixture) GoAccept(event ui.Event) (ui.InteractionReply, error) {
	reply := ui.InteractionReply{Domain: ui.DomainNotCalled}
	if event.Kind != ui.Accept || event.Dialog == nil || event.Handle.Path != MixedPath {
		return reply, fmt.Errorf("unexpected mixed Accept")
	}
	mode := f.take("mixed")
	f.mu.Lock()
	f.calls["mixed"]++
	count := f.calls["mixed"]
	f.mu.Unlock()
	f.emit("go-accept", map[string]any{"sequence": event.Sequence, "calls": count, "outcome": mode, "fields": event.Dialog.Fields})
	if mode == "error" {
		reply.Domain = ui.DomainUnknown
		return reply, fmt.Errorf("injected mixed persistence error")
	}
	reply.Accept = &ui.AcceptDecision{Accepted: mode != "false", Message: "Mixed accepted"}
	reply.Domain = ui.DomainRejected
	if mode == "false" {
		reply.Accept.Message = "Mixed rejected"
		return reply, nil
	}
	f.mu.Lock()
	for _, field := range event.Dialog.Fields {
		f.domain.Mixed[field.Handle.Path] = field.Value
	}
	f.domain.MixedCommits++
	f.mu.Unlock()
	reply.Domain = ui.DomainSucceeded
	if mode == "draft-conflict" {
		if err := f.Conflict(Targets["FormCount"]); err != nil {
			return reply, err
		}
	}
	if mode == "options-conflict" {
		if err := f.Replace(Targets["FormMode"]); err != nil {
			return reply, err
		}
	}
	return reply, nil
}
func Plans() map[string]bridge.Plan {
	p := map[string]bridge.Plan{}
	for name, selector := range map[string]bridge.EventField{"AcceptFlag": bridge.ControlBoolean, "AcceptLevel": bridge.ControlNumber, "AcceptCount": bridge.ControlNumber, "AcceptMode": bridge.ChoiceOptionID, "ChildCount": bridge.ControlNumber} {
		p["values."+name] = bridge.Plan{ResultMode: bridge.ScalarResult, Inputs: map[string]bridge.Source{"Value": {EventField: selector}}, OutputField: "Value"}
	}
	token := sdl.Text("load")
	p["values.Load"] = bridge.Plan{Inputs: map[string]bridge.Source{"Value": {Literal: &token}}, OutputField: "Value"}
	p["values.SaveText"] = bridge.Plan{Inputs: map[string]bridge.Source{"Value": {EventField: bridge.ControlText}}, OutputField: "Value"}
	p["values.TextSave"] = bridge.Plan{ResultMode: bridge.DialogAcceptResult, Inputs: map[string]bridge.Source{"Value": {EventField: bridge.DialogFieldValue, FieldPath: "textName"}}, AcceptField: "Accepted", MessageField: "Message"}
	return p
}
func (f *Fixture) Request(source string, sequence uint64) (fynehost.DocumentRequest, error) {
	d, err := parser.Parse(source)
	if err != nil {
		return fynehost.DocumentRequest{}, err
	}
	revision := f.Engine.Revision()
	return fynehost.DocumentRequest{Document: d, Entry: "page", SessionID: "values", SourceRevision: fmt.Sprintf("values-%d", sequence), Sequence: sequence, Mode: preparation.Connected, Choices: f.Choices(), Guard: func() error {
		if f.Engine.Revision() != revision {
			return fmt.Errorf("stale SDL revision")
		}
		return nil
	}, Bind: func(s *ui.Session, doc *parser.Document) error {
		if _, err := bridge.Bind(context.Background(), s, doc, map[string]*sdl.Engine{"values": f.Engine}, Plans(), nil); err != nil {
			return err
		}
		h, ok := s.Surface(MixedPath)
		if !ok {
			return fmt.Errorf("missing mixed form")
		}
		if err := s.BindInteraction(h, f.GoAccept); err != nil {
			return err
		}
		for _, w := range s.Widgets() {
			if w.Handle.Kind == "input" {
				continue
			} // Existing text Change delivery belongs to WCI3-M2.
			if _, ok := s.Field(w.Handle); !ok {
				continue
			}
			if err := s.ObserveChanges(w.Handle, f.changed); err != nil {
				return err
			}
		}
		return nil
	}}, nil
}
func (f *Fixture) changed(change ui.FieldChange) {
	path := change.Field.InstancePath
	f.mu.Lock()
	f.changes[path]++
	count := f.changes[path]
	armed := f.observers[path]
	delete(f.observers, path)
	f.mu.Unlock()
	f.emit("change", map[string]any{"path": path, "count": count, "field": change.Field})
	if armed != "" {
		err := f.Conflict(path)
		result := map[string]any{"path": path, "status": "ok"}
		if err != nil {
			result["status"] = "error"
			result["error"] = err.Error()
		}
		f.emit("observer-conflict", result)
	}
}
func (f *Fixture) Counts() (map[string]int, map[string]int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	a, c := map[string]int{}, map[string]int{}
	for k, v := range f.calls {
		a[k] = v
	}
	for k, v := range f.changes {
		c[k] = v
	}
	return a, c
}
func (f *Fixture) Domain() DomainState {
	f.mu.Lock()
	defer f.mu.Unlock()
	d := f.domain
	d.Scalars = map[string]sdl.Value{}
	d.ScalarCommits = map[string]int{}
	d.Mixed = map[string]ui.Value{}
	for k, v := range f.domain.Scalars {
		d.Scalars[k] = v
	}
	for k, v := range f.domain.ScalarCommits {
		d.ScalarCommits[k] = v
	}
	for k, v := range f.domain.Mixed {
		d.Mixed[k] = v
	}
	return d
}
func (f *Fixture) Close() {
	f.mu.Lock()
	if f.closed {
		f.mu.Unlock()
		return
	}
	f.closed = true
	f.mu.Unlock()
	f.Engine.Close()
}
