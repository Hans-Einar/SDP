// Package parser implements the bounded design-core structural language. It
// never executes design statements, loads modules or infers runtime behavior.
package parser

import (
	"fmt"
	"sort"
	"strings"
)

type Span struct {
	Start     int `json:"start"`
	End       int `json:"end"`
	Line      int `json:"line"`
	Column    int `json:"column"`
	EndLine   int `json:"end_line"`
	EndColumn int `json:"end_column"`
}
type Identifier struct {
	Name string
	Span Span
}
type Integer struct {
	Value int
	Span  Span
}
type Header struct {
	Language, Version string
	Span              Span
}
type Declaration struct {
	Kind string
	Name Identifier
	Span Span
}

// Statement is a closed sum: Kind selects its typed fields. Data emits the
// original named AST node form, preserving argument source spans.
type Statement struct {
	Kind                                                                                                      string
	Subject, Object, Interface, Mode, Container, Dataset, Datagram, Field, Channel, Message, Sender, Receiver Identifier
	Verb, Property, Value, Role                                                                               string
	Offset, Width, Ordinal                                                                                    Integer
	Variant                                                                                                   *Identifier
	ReplyTo                                                                                                   *Integer
	Span                                                                                                      Span
}
type Model struct {
	Header       Header
	Declarations []Declaration
	Statements   []Statement
	Span         Span
}
type Diagnostic struct {
	Code, Message string
	Span          Span
}

func (d Diagnostic) Error() string {
	return fmt.Sprintf("%s at %d:%d: %s", d.Code, d.Span.Line, d.Span.Column, d.Message)
}
func Symbols(m *Model) map[string]Declaration {
	out := map[string]Declaration{}
	for _, d := range m.Declarations {
		if _, ok := out[d.Name.Name]; !ok {
			out[d.Name.Name] = d
		}
	}
	return out
}
func (s Statement) Sentence() string {
	a := s.Subject.Name
	switch s.Kind {
	case "Relation":
		return a + " " + s.Verb + " " + s.Object.Name + "."
	case "Dependency":
		return a + " requires " + s.Interface.Name + " in mode " + s.Mode.Name + "."
	case "Allocation":
		return a + " allocated-to " + s.Container.Name + " in mode " + s.Mode.Name + "."
	case "PropertyAssignment":
		return a + " has " + s.Property + " = " + s.Value + "."
	case "Projection":
		return a + " projects " + s.Dataset.Name + " into " + s.Datagram.Name + "."
	case "Placement":
		return fmt.Sprintf("%s places %s at %d bits %d.", a, s.Field.Name, s.Offset.Value, s.Width.Value)
	case "Participation":
		return fmt.Sprintf("%s uses %s as %s of %s in mode %s.", a, s.Channel.Name, s.Role, s.Message.Name, s.Mode.Name)
	case "Step":
		v := fmt.Sprintf("%s step %d sends %s", a, s.Ordinal.Value, s.Message.Name)
		if s.Variant != nil {
			v += " variant " + s.Variant.Name
		}
		v += fmt.Sprintf(" from %s to %s via %s", s.Sender.Name, s.Receiver.Name, s.Channel.Name)
		if s.ReplyTo != nil {
			v += fmt.Sprintf(" reply-to %d", s.ReplyTo.Value)
		}
		return v + "."
	}
	panic("unknown statement kind")
}
func Canonical(m *Model) (string, []Diagnostic) {
	d := Validate(m)
	if len(d) > 0 {
		return "", d
	}
	decl := append([]Declaration{}, m.Declarations...)
	sort.Slice(decl, func(i, j int) bool { return decl[i].Name.Name < decl[j].Name.Name })
	lines := []string{"language design-core version 0.5."}
	for _, d := range decl {
		lines = append(lines, d.Kind+" "+d.Name.Name+".")
	}
	facts := []string{}
	for _, s := range m.Statements {
		facts = append(facts, s.Sentence())
	}
	sort.Strings(facts)
	lines = append(lines, facts...)
	return strings.Join(lines, "\n") + "\n", nil
}
