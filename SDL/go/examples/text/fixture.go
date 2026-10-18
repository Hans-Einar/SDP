package text

import (
	"context"
	"fmt"
	"github.com/Hans-Einar/SDP/SDUI/go/host/fynehost"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/preparation"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/bridge"
	sp "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	sdl "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/runtime"
	"sync"
)

type DomainState struct {
	Commits     map[string]int    `json:"commits"`
	Values      map[string]string `json:"values"`
	TextCommits int               `json:"textCommits"`
	TextValues  map[string]string `json:"textValues"`
}
type Fixture struct {
	Engine           *sdl.Engine
	mu               sync.Mutex
	calls, changes   map[string]int
	outcomes         map[string]string
	domain           DomainState
	load             string
	closed           bool
	Log              func(string, any)
	Conflict         func(string) error
	PrepareResources func(ui.Snapshot) error
}

func New() (*Fixture, error) {
	f := &Fixture{calls: map[string]int{}, changes: map[string]int{}, outcomes: map[string]string{}, load: "Loaded text", domain: DomainState{Commits: map[string]int{}, Values: map[string]string{}, TextValues: map[string]string{}}}
	for _, p := range Targets {
		f.changes[p] = 0
	}
	program, err := sp.CompileActions(Actions)
	if err != nil {
		return nil, err
	}
	registry := sdl.Registry{}
	for name, a := range program.Actions {
		f.calls[name] = 0
		f.domain.Commits[name] = 0
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
func actionPath(name string) string {
	if name == "Load" {
		return Targets["Preview"]
	}
	if name == "TextSave" {
		return Targets["FormName"]
	}
	return Targets[map[string]string{"SaveSingle": "Single", "SaveRequired": "Required", "SaveMulti": "Multi", "SaveOther": "Other", "SaveChild": "Child"}[name]]
}
func (f *Fixture) handler(name string) sdl.Handler {
	return func(ctx context.Context, r sdl.Record) (sdl.Record, error) {
		f.mu.Lock()
		mode := f.outcomes[name]
		delete(f.outcomes, name)
		if mode == "" {
			mode = "success"
		}
		f.calls[name]++
		count := f.calls[name]
		loaded := f.load
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
		value := r["Value"].Text
		if name == "Load" {
			value = loaded
		}
		f.mu.Lock()
		if name == "TextSave" {
			f.domain.TextCommits++
			f.domain.TextValues = map[string]string{"Name": r["Name"].Text, "Body": r["Body"].Text}
		} else if name != "Load" {
			f.domain.Commits[name]++
			f.domain.Values[name] = value
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
		if name == "TextSave" {
			return sdl.Record{"Accepted": sdl.Boolean(true), "Message": sdl.Text("Saved text")}, nil
		}
		if mode == "non-echo" {
			value += " changed by domain"
		}
		return sdl.Record{"Value": sdl.Text(value)}, nil
	}
}
func Plans() map[string]bridge.Plan {
	p := map[string]bridge.Plan{}
	for _, name := range []string{"SaveSingle", "SaveRequired", "SaveMulti", "SaveOther", "SaveChild"} {
		p["text."+name] = bridge.Plan{Inputs: map[string]bridge.Source{"Value": {EventField: bridge.ControlText}}, OutputField: "Value"}
	}
	token := sdl.Text("load")
	p["text.Load"] = bridge.Plan{Inputs: map[string]bridge.Source{"Value": {Literal: &token}}, OutputField: "Value"}
	p["text.TextSave"] = bridge.Plan{ResultMode: bridge.DialogAcceptResult, Inputs: map[string]bridge.Source{"Name": {EventField: bridge.DialogFieldValue, FieldPath: "formName"}, "Body": {EventField: bridge.DialogFieldValue, FieldPath: "formMulti"}}, AcceptField: "Accepted", MessageField: "Message"}
	return p
}
func (f *Fixture) Request(source string, sequence uint64) (fynehost.DocumentRequest, error) {
	d, err := parser.Parse(source)
	if err != nil {
		return fynehost.DocumentRequest{}, err
	}
	revision := f.Engine.Revision()
	return fynehost.DocumentRequest{Document: d, Entry: "page", SessionID: "text", SourceRevision: fmt.Sprintf("text-%d", sequence), Sequence: sequence, Mode: preparation.Connected, Guard: func() error {
		if f.Engine.Revision() != revision {
			return fmt.Errorf("stale SDL revision")
		}
		return nil
	}, PrepareResources: func(s ui.Snapshot) error {
		if f.PrepareResources != nil {
			return f.PrepareResources(s)
		}
		return nil
	}, Bind: func(s *ui.Session, doc *parser.Document) error {
		if _, err := bridge.Bind(context.Background(), s, doc, map[string]*sdl.Engine{"text": f.Engine}, Plans(), nil); err != nil {
			return err
		}
		for _, w := range s.Widgets() {
			field, ok := s.Field(w.Handle)
			if !ok || field.Input == nil {
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
	f.mu.Unlock()
	f.emit("change", map[string]any{"path": path, "count": count, "field": change.Field})
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
	d := DomainState{Commits: map[string]int{}, Values: map[string]string{}, TextValues: map[string]string{}, TextCommits: f.domain.TextCommits}
	for k, v := range f.domain.Commits {
		d.Commits[k] = v
	}
	for k, v := range f.domain.Values {
		d.Values[k] = v
	}
	for k, v := range f.domain.TextValues {
		d.TextValues[k] = v
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
