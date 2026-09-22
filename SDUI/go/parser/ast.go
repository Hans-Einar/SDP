// Package parser implements SDUI 0.2 without I/O, GUI dependencies or callback execution.
package parser

import (
	"fmt"
	"reflect"
	"strings"
)

type Span struct {
	Start  int `json:"start"`
	End    int `json:"end"`
	Line   int `json:"line"`
	Column int `json:"column"`
}
type Diagnostic struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Span    Span   `json:"span"`
}

func (d *Diagnostic) Error() string {
	return fmt.Sprintf("%s at %d:%d: %s", d.Code, d.Span.Line, d.Span.Column, d.Message)
}

type Literal struct {
	Kind  string `json:"kind"`
	Value any    `json:"value"`
	Span  Span   `json:"span"`
}
type Reference struct {
	Module string `json:"module"`
	Object string `json:"object"`
	Member string `json:"member"`
	Span   Span   `json:"span"`
}
type Argument struct {
	Name  *string `json:"name"`
	Value any     `json:"value"`
	Span  Span    `json:"span"`
}
type LayoutRule struct {
	Name     string  `json:"name"`
	Value    Literal `json:"value"`
	Spelling string  `json:"spelling"`
	Span     Span    `json:"span"`
}
type Row struct {
	Items []*Node `json:"items"`
	Span  Span    `json:"span"`
}
type Node struct {
	Kind      string       `json:"kind"`
	Name      *string      `json:"name"`
	Role      *string      `json:"role"`
	Rows      []Row        `json:"rows"`
	Widget    *string      `json:"widget"`
	Arguments []Argument   `json:"arguments"`
	Text      *Literal     `json:"text"`
	Target    *string      `json:"target"`
	Variant   *string      `json:"variant"`
	Layout    []LayoutRule `json:"layout"`
	Span      Span         `json:"span"`
}
type ModuleRef struct {
	Alias string `json:"alias"`
	Path  string `json:"path"`
	Span  Span   `json:"span"`
}
type Definition struct {
	Name string `json:"name"`
	Root *Node  `json:"root"`
	Span Span   `json:"span"`
}
type Connection struct {
	Module     string   `json:"module"`
	Object     string   `json:"object"`
	Definition string   `json:"definition"`
	Path       []string `json:"path"`
	Span       Span     `json:"span"`
}
type Document struct {
	Profile     string       `json:"profile"`
	References  []ModuleRef  `json:"references"`
	Definitions []Definition `json:"definitions"`
	Connections []Connection `json:"connections"`
	Span        Span         `json:"span"`
}

// Data returns the stable tagged sdui-ast/0.2 JSON representation.
func Data(x any) any { return data(reflect.ValueOf(x)) }
func data(v reflect.Value) any {
	if !v.IsValid() {
		return nil
	}
	if v.Kind() == reflect.Interface || v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil
		}
		return data(v.Elem())
	}
	switch v.Kind() {
	case reflect.Struct:
		out := map[string]any{"type": v.Type().Name()}
		for i := 0; i < v.NumField(); i++ {
			name := strings.Split(v.Type().Field(i).Tag.Get("json"), ",")[0]
			if name == "" {
				name = v.Type().Field(i).Name
			}
			if name != "-" {
				out[name] = data(v.Field(i))
			}
		}
		return out
	case reflect.Slice, reflect.Array:
		out := make([]any, v.Len())
		for i := range out {
			out[i] = data(v.Index(i))
		}
		return out
	default:
		return v.Interface()
	}
}
func str(s string) *string { return &s }
func val(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
func fail(code, msg string, s Span) { panic(&Diagnostic{code, msg, s}) }
func recoverDiagnostic(err *error) {
	if r := recover(); r != nil {
		if d, ok := r.(*Diagnostic); ok {
			*err = d
		} else {
			panic(r)
		}
	}
}
