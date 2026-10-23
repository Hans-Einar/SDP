package parser

import (
	"errors"
	"strings"

	"github.com/Hans-Einar/SDP/SDUI/go/numeric"
)

var scalars = map[string]widgetSchema{
	"checkbox": {"label", map[string]string{"label": "string", "value": "boolean", "readOnly": "boolean", "callback": "reference"}, "label"},
	"slider":   {"label", map[string]string{"label": "string", "min": "number-lexeme", "max": "number-lexeme", "step": "number-lexeme", "value": "number-lexeme", "readOnly": "boolean", "callback": "reference"}, "label"},
	"number":   {"label", map[string]string{"label": "string", "min": "number-lexeme", "max": "number-lexeme", "step": "number-lexeme", "value": "number-lexeme", "readOnly": "boolean", "placeholder": "string", "callback": "reference"}, "label"},
	"select":   {"label", map[string]string{"label": "string", "value": "string", "required": "boolean", "readOnly": "boolean", "callback": "reference"}, "label"},
}

func scalarKind(kind string) bool { _, ok := scalars[kind]; return ok }

// NumericArguments extracts exact source lexemes, never a formatted binary64 value.
// Grid admission is performed by the shared numeric package, not this accessor.
func NumericArguments(n *Instance) (min, max, step, value string, err error) {
	if n == nil {
		return "", "", "", "", &Diagnostic{Code: "numeric-arguments", Message: "Expected numeric widget"}
	}
	bad := func(message string, span Span) error {
		return &Diagnostic{Code: "numeric-arguments", Message: n.Path + ": " + message, Span: span}
	}
	if n.Profile != "sdui/0.3" || n.Kind != "widget" || !member(n.Widget, "slider number") {
		return "", "", "", "", bad("Expected .3 slider/number", n.Span)
	}
	values := [4]string{}
	for i, key := range []string{"min", "max", "step", "value"} {
		lit, ok := n.Arguments[key].(Literal)
		if !ok {
			return "", "", "", "", bad("Missing number-lexeme "+key, n.Span)
		}
		raw, ok := lit.Value.(string)
		if lit.Kind != "number-lexeme" || !ok || raw == "" {
			return "", "", "", "", bad("Expected number-lexeme "+key, lit.Span)
		}
		values[i] = raw
	}
	return values[0], values[1], values[2], values[3], nil
}

func validateScalar(n *Instance) {
	if n.Profile != "sdui/0.3" || !scalarKind(n.Widget) {
		return
	}
	if n.Kind != "widget" || len(n.Rows) > 0 || len(n.Regions) > 0 {
		fail("value-body", "Scalar controls are leaves", n.Span)
	}
	schema := scalars[n.Widget]
	for key, arg := range n.Arguments {
		expected, ok := schema.fields[key]
		if !ok {
			fail("value-argument", n.Path+": unknown argument "+key, n.Span)
		}
		switch v := arg.(type) {
		case Literal:
			valid := v.Kind == expected
			switch expected {
			case "string", "number-lexeme":
				_, ok := v.Value.(string)
				valid = valid && ok
			case "boolean":
				_, ok := v.Value.(bool)
				valid = valid && ok
			default:
				valid = false
			}
			if !valid {
				fail("value-argument", n.Path+": invalid type for "+key, v.Span)
			}
		case Reference:
			if expected != "reference" {
				fail("value-argument", "Unexpected reference for "+key, v.Span)
			}
		default:
			fail("value-argument", "Invalid argument "+key, n.Span)
		}
	}
	if strings.TrimSpace(n.Argument("label")) == "" {
		fail("value-label", "Scalar requires nonblank label", n.Span)
	}
	if member(n.Widget, "slider number") {
		min, max, step, value, err := NumericArguments(n)
		if err != nil {
			panic(err)
		}
		grid, err := numeric.NewGrid(min, max, step)
		if err != nil {
			scalarNumericError(n, err, "")
		}
		if _, err = grid.Parse(value); err != nil {
			scalarNumericError(n, err, "value")
		}
	}
}
func scalarNumericError(n *Instance, err error, argument string) {
	span := n.Span
	var detail *numeric.Error
	if errors.As(err, &detail) && argument == "" {
		argument = detail.Argument
	}
	if lit, ok := n.Arguments[argument].(Literal); ok {
		span = lit.Span
	}
	code := "value-number"
	if detail != nil && detail.Code != "" {
		code = detail.Code
	}
	fail(code, n.Path+": "+err.Error(), span)
}
