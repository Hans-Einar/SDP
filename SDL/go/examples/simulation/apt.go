// Package simulation is handwritten demo domain logic, never a Ponsse backend.
package simulation

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/bridge"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	sdl "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/runtime"
)

type Operation struct {
	Sequence           uint64 `json:"sequence"`
	Expected, Revision int64
	Value              string
	Accepted           bool
}
type Apt struct {
	mu       sync.Mutex
	value    string
	revision int64
	trace    []Operation
}

func NewApt() *Apt { return &Apt{value: "400"} }
func (a *Apt) Registry() sdl.Registry {
	return sdl.Registry{"GoEditAptCell": {Input: parser.RecordType{"CellKey": parser.TextType, "ExpectedRevision": parser.IntegerType, "Value": parser.TextType}, Output: parser.RecordType{"Revision": parser.IntegerType, "Value": parser.TextType}, Call: a.Edit}}
}
func (a *Apt) Edit(ctx context.Context, input sdl.Record) (sdl.Record, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	call, ok := sdl.CurrentInvocation(ctx)
	if !ok {
		return nil, fmt.Errorf("simulation-context: missing command identity")
	}
	value := input["Value"].Text
	expected := input["ExpectedRevision"].Integer
	operation := Operation{Sequence: call.Sequence, Expected: expected, Revision: a.revision, Value: value}
	var err error
	if input["CellKey"].Text != "LengthA1" {
		err = fmt.Errorf("apt-cell: unknown simulated cell")
	} else if expected != a.revision {
		err = fmt.Errorf("apt-revision: expected %d, current %d", expected, a.revision)
	} else {
		n, e := strconv.Atoi(value)
		if e != nil || n < 100 || n > 1000 {
			err = fmt.Errorf("apt-value: simulation requires integer 100–1000")
		}
	}
	if err != nil {
		a.trace = append(a.trace, operation)
		return nil, err
	}
	a.value = value
	a.revision++
	operation.Revision = a.revision
	operation.Accepted = true
	a.trace = append(a.trace, operation)
	return sdl.Record{"Value": sdl.Text(a.value), "Revision": sdl.Integer(a.revision)}, nil
}
func (a *Apt) Snapshot() (string, int64, []Operation) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.value, a.revision, append([]Operation{}, a.trace...)
}
func Plans() map[string]bridge.Plan {
	cell := sdl.Text("LengthA1")
	return map[string]bridge.Plan{"apt.EditAptCell": {Inputs: map[string]bridge.Source{"CellKey": {Literal: &cell}, "ExpectedRevision": {Context: "AptRevision"}, "Value": {Widget: "page/value"}}, OutputField: "Value", RevisionField: "Revision", RevisionContext: "AptRevision"}}
}
func InitialContext() map[string]sdl.Value {
	return map[string]sdl.Value{"AptRevision": sdl.Integer(0)}
}
