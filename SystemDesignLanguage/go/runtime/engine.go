package runtime

import (
	"context"
	"fmt"
	"sync"

	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
)

type Handler func(context.Context, Record) (Record, error)
type Binding struct {
	Input, Output parser.RecordType
	Call          Handler
}
type Registry map[string]Binding
type Request struct {
	Action   string
	Revision uint64
	Sequence uint64
	Input    Record
}
type Result struct {
	Action             string
	Revision, Sequence uint64
	Output             Record
}
type Engine struct {
	mu                 sync.Mutex
	program            *parser.Program
	registry           Registry
	revision, sequence uint64
	closed             bool
}

func New(program *parser.Program, registry Registry) (*Engine, error) {
	p, err := copyProgram(program)
	if err != nil {
		return nil, err
	}
	bindings := Registry{}
	for name, b := range registry {
		bindings[name] = Binding{cloneType(b.Input), cloneType(b.Output), b.Call}
	}
	if err = checkBindings(p, bindings); err != nil {
		return nil, err
	}
	return &Engine{program: p, registry: bindings, revision: 1}, nil
}
func copyProgram(p *parser.Program) (*parser.Program, error) {
	if p == nil || p.Model == nil {
		return nil, fmt.Errorf("action-model: missing model")
	}
	m := *p.Model
	m.Declarations = append([]parser.Declaration{}, m.Declarations...)
	m.Statements = append([]parser.ActionStatement{}, m.Statements...)
	return parser.ValidateActions(&m)
}
func checkBindings(p *parser.Program, r Registry) error {
	for _, d := range p.Model.Declarations {
		if d.Kind != "action" {
			continue
		}
		a := p.Actions[d.Name.Name]
		b, ok := r[a.GoSymbol]
		if !ok || b.Call == nil {
			return fmt.Errorf("unbound-action: %s invokes %s at %d:%d", a.Name, a.GoSymbol, a.Span.Line, a.Span.Column)
		}
		if !sameType(p.Records[a.Input], b.Input) || !sameType(p.Records[a.Output], b.Output) {
			return fmt.Errorf("binding-signature: %s at %d:%d", a.Name, a.Span.Line, a.Span.Column)
		}
	}
	return nil
}
func (e *Engine) Revision() uint64 { e.mu.Lock(); defer e.mu.Unlock(); return e.revision }
func (e *Engine) Signature(action string) (parser.RecordType, parser.RecordType, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	a, ok := e.program.Actions[action]
	if !ok {
		return nil, nil, fmt.Errorf("unknown-action: %s", action)
	}
	return cloneType(e.program.Records[a.Input]), cloneType(e.program.Records[a.Output]), nil
}

// Execute serializes invocations. A consumed command sequence is never retried
// automatically, including handler errors or invalid output. Callers must not
// re-enter this Engine from a handler; domain handlers own only domain state.
func (e *Engine) Execute(ctx context.Context, request Request) (result Result, err error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.closed {
		return Result{}, fmt.Errorf("closed: SDL runtime")
	}
	if request.Revision != e.revision {
		return Result{}, fmt.Errorf("stale-model: expected revision %d", e.revision)
	}
	if request.Sequence == 0 || request.Sequence <= e.sequence {
		return Result{}, fmt.Errorf("duplicate-command: %d", request.Sequence)
	}
	a, ok := e.program.Actions[request.Action]
	if !ok {
		return Result{}, fmt.Errorf("unknown-action: %s", request.Action)
	}
	if err = checkRecord(request.Input, e.program.Records[a.Input]); err != nil {
		return Result{}, err
	}
	if err = ctx.Err(); err != nil {
		return Result{}, err
	}
	e.sequence = request.Sequence
	defer func() {
		if v := recover(); v != nil {
			result = Result{}
			err = fmt.Errorf("handler-panic: %s: %v", a.Name, v)
		}
	}()
	output, err := e.registry[a.GoSymbol].Call(ctx, cloneRecord(request.Input))
	if err != nil {
		return Result{}, err
	}
	if err = checkRecord(output, e.program.Records[a.Output]); err != nil {
		return Result{}, fmt.Errorf("action-result %s: %w", a.Name, err)
	}
	return Result{Action: a.Name, Revision: e.revision, Sequence: request.Sequence, Output: cloneRecord(output)}, nil
}
func (e *Engine) Close() { e.mu.Lock(); defer e.mu.Unlock(); e.closed = true; e.revision++ }
