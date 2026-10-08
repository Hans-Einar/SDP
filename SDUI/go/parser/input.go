package parser

// InputPolicy contains effective source options, without inserting AST defaults.
// Extended is true only when a .3 input explicitly supplies a new text option.
type InputPolicy struct {
	Extended, Multiline, ReadOnly, Required bool
	Placeholder                             string
	PlaceholderSet                          bool
}

var input03 = widgetSchema{"text", map[string]string{
	"text": "string", "value": "string", "callback": "reference",
	"multiline": "boolean", "readOnly": "boolean", "required": "boolean", "placeholder": "string",
}, "text"}

func extendedInputArguments(args map[string]any) bool {
	for _, key := range []string{"multiline", "readOnly", "required", "placeholder"} {
		if _, ok := args[key]; ok {
			return true
		}
	}
	return false
}

// InputOptions validates an input's closed schema and returns effective options.
// Presence (even false/empty), not profile or value, opts into extended behavior.
// Text content validation belongs to runtime; required-empty is not a syntax error.
func InputOptions(n *Instance) (InputPolicy, error) {
	var out InputPolicy
	bad := func(message string, span Span) (InputPolicy, error) {
		path := ""
		if n != nil {
			path = n.Path + ": "
		}
		return InputPolicy{}, &Diagnostic{Code: "input-argument", Message: path + message, Span: span}
	}
	if n == nil {
		return bad("Expected input widget", Span{})
	}
	profile, err := EffectiveProfile(n)
	if err != nil {
		return out, err
	}
	if n.Kind != "widget" || n.Widget != "input" || len(n.Rows) != 0 || len(n.Regions) != 0 {
		return bad("Expected leaf input widget", n.Span)
	}
	schema, _ := schemaFor(profile, "input")
	for key, arg := range n.Arguments {
		expected, ok := schema.fields[key]
		if !ok {
			span := n.Span
			switch v := arg.(type) {
			case Literal:
				span = v.Span
			case Reference:
				span = v.Span
			}
			return bad("Unknown input argument "+key, span)
		}
		switch v := arg.(type) {
		case Literal:
			valid := v.Kind == expected
			switch expected {
			case "string":
				_, ok = v.Value.(string)
			case "boolean":
				_, ok = v.Value.(bool)
			default:
				ok = false
			}
			if !valid || !ok {
				return bad("Invalid type for "+key, v.Span)
			}
		case Reference:
			if expected != "reference" {
				return bad("Unexpected reference for "+key, v.Span)
			}
		default:
			return bad("Invalid argument "+key, n.Span)
		}
	}
	if _, ok := n.Arguments["text"]; !ok {
		return bad("Missing text", n.Span)
	}
	out.Extended = extendedInputArguments(n.Arguments)
	out.Multiline = boolArg(n, "multiline")
	out.ReadOnly = boolArg(n, "readOnly")
	out.Required = boolArg(n, "required")
	_, out.PlaceholderSet = n.Arguments["placeholder"]
	out.Placeholder = n.Argument("text")
	if out.PlaceholderSet {
		out.Placeholder = n.Argument("placeholder")
	}
	return out, nil
}
