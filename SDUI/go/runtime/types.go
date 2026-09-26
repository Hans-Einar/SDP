// Package runtime owns SDUI UI state. Its methods run on one owner goroutine
// (the host UI goroutine). Domain state and native widget pointers never enter it.
package runtime

import (
	"fmt"
	"strings"

	"github.com/Hans-Einar/SDP/SDUI/go/parser"
)

type Handle struct {
	Session    string
	Path       string
	Generation uint64
	Kind       string
}
type ValueKind string

const (
	String  ValueKind = "string"
	Boolean ValueKind = "boolean"
)

type Value struct {
	Kind ValueKind
	Text string
	Bool bool
}

func Text(v string) Value { return Value{Kind: String, Text: v} }
func Bool(v bool) Value   { return Value{Kind: Boolean, Bool: v} }

type EventKind string

const (
	Activate EventKind = "activate"
	Commit   EventKind = "commit"
)

type Event struct {
	Handle                                 Handle
	ModelRevision, Sequence, DraftRevision uint64
	Kind                                   EventKind
	Value                                  Value
}
type Property string

const (
	Label         Property = "label"
	AcceptedValue Property = "value"
	Enabled       Property = "enabled"
	Visible       Property = "visible"
)

type Update struct {
	Handle                Handle
	Property              Property
	Value                 Value
	ExpectedValueRevision uint64
	AcceptDraft           bool
}
type Handler func(Event) ([]Update, error)
type Widget struct {
	ancestorEnabled, ancestorVisible bool
	Handle                           Handle
	InstancePath                     string
	Label, Value, Draft              string
	ValueRevision, DraftRevision     uint64
	Dirty, Enabled, Visible          bool
	Binding                          parser.Reference
}
type Fault struct{ Code, Message string }

func (f *Fault) Error() string         { return f.Code + ": " + f.Message }
func fault(code, message string) error { return &Fault{code, message} }
func public(path string) string {
	parts := strings.Split(path, "/")
	if strings.HasPrefix(parts[len(parts)-1], "$") {
		return "@" + path
	}
	out := []string{}
	for _, p := range parts {
		if !strings.HasPrefix(p, "$") {
			out = append(out, p)
		}
	}
	return strings.Join(out, "/")
}
func clone(n *parser.Instance) *parser.Instance {
	if n == nil {
		return nil
	}
	v := *n
	v.Layout = map[string]any{}
	for k, x := range n.Layout {
		if a, ok := x.([]float64); ok {
			x = append([]float64{}, a...)
		}
		v.Layout[k] = x
	}
	v.Arguments = map[string]any{}
	for k, x := range n.Arguments {
		v.Arguments[k] = x
	}
	v.Rows = make([][]*parser.Instance, len(n.Rows))
	for r, row := range n.Rows {
		for _, c := range row {
			v.Rows[r] = append(v.Rows[r], clone(c))
		}
	}
	v.Regions = nil
	for _, r := range n.Regions {
		v.Regions = append(v.Regions, parser.Region{Role: r.Role, Node: clone(r.Node)})
	}
	return &v
}
func validValue(v Value, kind ValueKind) bool {
	if v.Kind != kind {
		return false
	}
	if kind == String {
		return !v.Bool && len(v.Text) <= 32768
	}
	return v.Text == ""
}
func (s *Session) lookup(h Handle) (*Widget, error) {
	if s.closed {
		return nil, fault("closed", "Session is closed")
	}
	w := s.widgets[h.Path]
	if w == nil || w.Handle != h {
		return nil, fault("stale-handle", fmt.Sprintf("Handle %s is no longer valid", h.Path))
	}
	return w, nil
}
