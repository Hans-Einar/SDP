package parser

import "reflect"

// Data is the stable JSON AST compatible with the original design-core profile.
func Data(value any) any {
	switch v := value.(type) {
	case *Model:
		if v == nil {
			return nil
		}
		return Data(*v)
	case Model:
		d := []any{}
		for _, x := range v.Declarations {
			d = append(d, Data(x))
		}
		s := []any{}
		for _, x := range v.Statements {
			s = append(s, Data(x))
		}
		return map[string]any{"node": "Model", "header": Data(v.Header), "declarations": d, "statements": s, "span": Data(v.Span)}
	case Span:
		return map[string]any{"node": "Span", "start": v.Start, "end": v.End, "line": v.Line, "column": v.Column, "end_line": v.EndLine, "end_column": v.EndColumn}
	case Identifier:
		return map[string]any{"node": "Identifier", "name": v.Name, "span": Data(v.Span)}
	case Integer:
		return map[string]any{"node": "Integer", "value": v.Value, "span": Data(v.Span)}
	case Header:
		return map[string]any{"node": "Header", "language": v.Language, "version": v.Version, "span": Data(v.Span)}
	case Declaration:
		return map[string]any{"node": "Declaration", "kind": v.Kind, "name": Data(v.Name), "span": Data(v.Span)}
	case Diagnostic:
		return map[string]any{"node": "Diagnostic", "code": v.Code, "message": v.Message, "span": Data(v.Span)}
	case Statement:
		m := map[string]any{"node": v.Kind, "subject": Data(v.Subject), "span": Data(v.Span)}
		fields := map[string]any{}
		switch v.Kind {
		case "Relation":
			fields = map[string]any{"verb": v.Verb, "object": v.Object}
		case "Dependency":
			fields = map[string]any{"interface": v.Interface, "mode": v.Mode}
		case "Allocation":
			fields = map[string]any{"container": v.Container, "mode": v.Mode}
		case "PropertyAssignment":
			fields = map[string]any{"property": v.Property, "value": v.Value}
		case "Projection":
			fields = map[string]any{"dataset": v.Dataset, "datagram": v.Datagram}
		case "Placement":
			fields = map[string]any{"field": v.Field, "offset": v.Offset, "width": v.Width}
		case "Participation":
			fields = map[string]any{"channel": v.Channel, "role": v.Role, "message": v.Message, "mode": v.Mode}
		case "Step":
			fields = map[string]any{"ordinal": v.Ordinal, "message": v.Message, "variant": v.Variant, "sender": v.Sender, "receiver": v.Receiver, "channel": v.Channel, "reply_to": v.ReplyTo}
		}
		for key, x := range fields {
			m[key] = Data(x)
		}
		return m
	case *Identifier:
		if v == nil {
			return nil
		}
		return Data(*v)
	case *Integer:
		if v == nil {
			return nil
		}
		return Data(*v)
	}
	r := reflect.ValueOf(value)
	if r.IsValid() && r.Kind() == reflect.Slice {
		out := []any{}
		for i := 0; i < r.Len(); i++ {
			out = append(out, Data(r.Index(i).Interface()))
		}
		return out
	}
	return value
}
