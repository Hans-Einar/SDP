package runtime

import (
	"fmt"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
)

// Reload waits for an active invocation's commit boundary, validates the new
// model against the existing Go registry, and then publishes one revision.
// Go-owned domain state is retained; no command is replayed. Go implementation
// changes require a rebuilt process, not this model operation.
func (e *Engine) Reload(program *parser.Program) error {
	p, err := copyProgram(program)
	if err != nil {
		return err
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.closed {
		return fmt.Errorf("closed: SDL runtime")
	}
	if err = checkBindings(p, e.registry); err != nil {
		return err
	}
	e.program = p
	e.revision++
	return nil
}
func (e *Engine) ReloadSource(source string) error {
	p, err := parser.CompileActions(source)
	if err != nil {
		return err
	}
	return e.Reload(p)
}

// Preview checks a proposed model without changing revision or domain state.
func (e *Engine) Preview(program *parser.Program) (*Engine, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.closed {
		return nil, fmt.Errorf("closed: SDL runtime")
	}
	return New(program, e.registry)
}
