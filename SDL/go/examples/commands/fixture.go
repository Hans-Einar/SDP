package commands

import (
	"context"
	"fmt"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"

	"github.com/Hans-Einar/SDP/SDUI/go/host/fynehost"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/preparation"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/bridge"
	sdlparser "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	sdl "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/runtime"
)

type DomainState struct {
	Name, Note string
	Commits    int
}
type Fixture struct {
	Engine                  *sdl.Engine
	mu                      sync.Mutex
	calls                   map[string]int
	outcomes                map[string]string
	domain                  DomainState
	pending                 map[string]pendingLoad
	resourceFailure, closed bool
	Log                     func(string, any)  // Install before invoking actions or providers.
	Conflict                func(string) error // UI-owner-only reentrant draft mutation by exact input path.
}

func New() (*Fixture, error) {
	f := &Fixture{calls: map[string]int{}, outcomes: map[string]string{}, pending: map[string]pendingLoad{}}
	p, err := sdlparser.CompileActions(Actions)
	if err != nil {
		return nil, err
	}
	registry := sdl.Registry{}
	inputs := map[string]sdlparser.RecordType{"Run": {"Token": sdlparser.TextType}, "Toggle": {"Checked": sdlparser.BooleanType}, "Inspect": {"ItemId": sdlparser.TextType}, "Save": {"Name": sdlparser.TextType, "Note": sdlparser.TextType}}
	for name, input := range inputs {
		output := sdlparser.RecordType{"Preview": sdlparser.TextType}
		if name == "Save" {
			output = sdlparser.RecordType{"Accepted": sdlparser.BooleanType, "Message": sdlparser.TextType}
		}
		registry["Go"+name] = sdl.Binding{Input: input, Output: output, Call: f.handler(name)}
	}
	f.Engine, err = sdl.New(p, registry)
	return f, err
}
func (f *Fixture) handler(name string) sdl.Handler {
	return func(ctx context.Context, r sdl.Record) (sdl.Record, error) {
		f.mu.Lock()
		f.calls[name]++
		count := f.calls[name]
		mode := f.outcomes[name]
		delete(f.outcomes, name)
		f.mu.Unlock()
		if mode == "" {
			mode = "success"
			if name == "Save" {
				mode = "true"
			}
		}
		invocation, _ := sdl.CurrentInvocation(ctx)
		if f.Log != nil {
			f.Log("action", map[string]any{"action": name, "input": r, "sequence": invocation.Sequence, "calls": count, "outcome": mode})
		}
		if mode == "error" {
			return nil, fmt.Errorf("injected %s domain error", name)
		}
		if mode == "malformed" {
			if name == "Save" {
				return sdl.Record{"Accepted": sdl.Text("invalid"), "Message": sdl.Text("Malformed acceptance")}, nil
			}
			return sdl.Record{"Preview": sdl.Integer(1)}, nil
		}
		if name == "Save" {
			if mode == "false" {
				return sdl.Record{"Accepted": sdl.Boolean(false), "Message": sdl.Text("Save rejected by fixture")}, nil
			}
			f.mu.Lock()
			f.domain = DomainState{Name: r["Name"].Text, Note: r["Note"].Text, Commits: f.domain.Commits + 1}
			f.mu.Unlock()
			if mode == "draft-conflict" {
				if err := f.conflict(NamePath); err != nil {
					return nil, err
				}
			}
			if mode == "resource-conflict" {
				f.mu.Lock()
				f.resourceFailure = true
				f.mu.Unlock()
			}
			return sdl.Record{"Accepted": sdl.Boolean(true), "Message": sdl.Text("Saved")}, nil
		}
		if mode == "draft-conflict" {
			if err := f.conflict(map[string]string{"Run": PreviewPath, "Toggle": FlagPreviewPath, "Inspect": ItemPreviewPath}[name]); err != nil {
				return nil, err
			}
		}
		value := "Run completed"
		if name == "Toggle" {
			value = fmt.Sprintf("Flag: %t", r["Checked"].Boolean)
		}
		if name == "Inspect" {
			value = "Item: " + r["ItemId"].Text
		}
		return sdl.Record{"Preview": sdl.Text(value)}, nil
	}
}
func (f *Fixture) conflict(path string) error {
	if f.Conflict == nil {
		return fmt.Errorf("fixture conflict hook missing")
	}
	return f.Conflict(path)
}
func (f *Fixture) Calls() map[string]int {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := map[string]int{}
	for k, v := range f.calls {
		out[k] = v
	}
	return out
}
func (f *Fixture) Domain() DomainState { f.mu.Lock(); defer f.mu.Unlock(); return f.domain }
func (f *Fixture) NextAction(name, mode string) error {
	switch name {
	case "Run", "Toggle", "Inspect":
		switch mode {
		case "success", "error", "malformed", "draft-conflict":
		default:
			return fmt.Errorf("unknown command outcome %q", mode)
		}
	case "Save":
		switch mode {
		case "true", "false", "error", "malformed", "draft-conflict", "resource-conflict":
		default:
			return fmt.Errorf("unknown acceptance outcome %q", mode)
		}
	default:
		return fmt.Errorf("unknown fixture action %q", name)
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.closed {
		return fmt.Errorf("fixture closed")
	}
	f.outcomes[name] = mode
	return nil
}
func Plans() map[string]bridge.Plan {
	token := sdl.Text("run")
	return map[string]bridge.Plan{
		"actions.Run":     {Inputs: map[string]bridge.Source{"Token": {Literal: &token}}, OutputField: "Preview"},
		"actions.Toggle":  {Inputs: map[string]bridge.Source{"Checked": {EventField: bridge.CommandChecked}}, OutputField: "Preview"},
		"actions.Inspect": {Inputs: map[string]bridge.Source{"ItemId": {EventField: bridge.CommandContextItemID}}, OutputField: "Preview"},
		"actions.Save":    {Inputs: map[string]bridge.Source{"Name": {EventField: bridge.DialogFieldValue, FieldPath: "form/body/name"}, "Note": {EventField: bridge.DialogFieldValue, FieldPath: "form/body/note"}}, ResultMode: bridge.DialogAcceptResult, AcceptField: "Accepted", MessageField: "Message"},
	}
}
func (f *Fixture) Request(source string, sequence uint64) (fynehost.DocumentRequest, error) {
	d, err := parser.Parse(source)
	if err != nil {
		return fynehost.DocumentRequest{}, err
	}
	revision := f.Engine.Revision()
	return fynehost.DocumentRequest{Document: d, Entry: "page", SessionID: "commands", SourceRevision: fmt.Sprintf("commands-%d", sequence), Sequence: sequence, Mode: preparation.Connected, Icons: map[string]fyne.Resource{"fixture-run": theme.MediaPlayIcon()}, Providers: f.Providers(), PrepareResources: func(ui.Snapshot) error {
		f.mu.Lock()
		fail := f.resourceFailure
		f.resourceFailure = false
		f.mu.Unlock()
		if fail {
			return fmt.Errorf("injected post-domain resource failure")
		}
		return nil
	}, Guard: func() error {
		if f.Engine.Revision() != revision {
			return fmt.Errorf("stale SDL revision")
		}
		return nil
	}, Bind: func(s *ui.Session, doc *parser.Document) error {
		_, err := bridge.Bind(context.Background(), s, doc, map[string]*sdl.Engine{"actions": f.Engine}, Plans(), nil)
		return err
	}}, nil
}
func (f *Fixture) Close() {
	f.mu.Lock()
	if f.closed {
		f.mu.Unlock()
		return
	}
	f.closed = true
	for key, p := range f.pending {
		p.reply <- loadReply{err: fmt.Errorf("fixture closed")}
		delete(f.pending, key)
	}
	f.mu.Unlock()
	f.Engine.Close()
}
