package runtime

import (
	"strings"

	"github.com/Hans-Einar/SDP/SDUI/go/parser"
)

func (s *Session) initInput(node *parser.Instance) error {
	policy, err := parser.InputOptions(node)
	if err != nil || !policy.Extended {
		return err
	}
	w := s.widgets[public(node.Path)]
	f := &field{FieldState: FieldState{InstancePath: node.Path, Input: &InputState{Multiline: policy.Multiline, Placeholder: policy.Placeholder}, ReadOnly: policy.ReadOnly, Required: policy.Required}}
	f.placeholderSet = policy.PlaceholderSet
	if !validDraft(w.Value) || !policy.Multiline && strings.ContainsAny(w.Value, "\r\n") {
		span := node.Span
		if value, ok := node.Arguments["value"].(parser.Literal); ok {
			span = value.Span
		}
		return &parser.Diagnostic{Code: "input-initial", Message: node.Path + ": initial text must be bounded UTF-8 and obey single-line policy", Span: span}
	}
	s.fields[w.Handle.Path] = f
	f.Validation = inputValidation(f, w.Draft)
	return nil
}

// Text stays in Widget. The field contains only policy and feedback; every
// public accepted/proposed/raw value is a detached projection of that Widget.
func inputValidation(f *field, text string) FieldValidation {
	if !validDraft(text) {
		return validation("text-bound")
	}
	if !f.Input.Multiline && strings.ContainsAny(text, "\r\n") {
		return validation("text-single-line")
	}
	if f.Required && strings.TrimSpace(text) == "" {
		return validation("text-required")
	}
	return FieldValidation{}
}

func (s *Session) revertInput(w *Widget) error {
	if !s.inputAllowed(w) || s.fields[w.Handle.Path].ReadOnly {
		return fault("field-read-only", "Input is not editable")
	}
	baseline := s.StateRevision
	n := s.copyState()
	n.resetField(w.Handle.Path)
	if err := n.validateField(w.Handle.Path); err != nil {
		return err
	}
	if s.StateRevision != baseline {
		return fault("stale-validation", "State changed in validator")
	}
	return s.publish(n)
}

func (n *Session) inheritInput(s *Session, path string) error {
	w := n.widgets[path]
	f := n.fields[path]
	old := s.widgets[path]
	if old != nil && old.Handle == w.Handle {
		validation := inputValidation(f, w.Value)
		previous := s.fields[path]
		// Exact empty absence on an already-required field remains editable
		// invalid. Whitespace and every other invalid accepted value reject.
		emptyBaseline := previous != nil && previous.Input != nil && previous.Required && f.Required && w.Value == ""
		if validation.Code != "" && !(validation.Code == "text-required" && emptyBaseline) {
			return fault("input-accepted", "Retained accepted text violates new constraints")
		}
		// inheritCommands already discards unaccepted closed-dialog drafts.
		// A surviving multiline draft cannot be silently flattened on conversion.
		if !f.Input.Multiline && (previous == nil || previous.Input == nil || previous.Input.Multiline) && strings.ContainsAny(w.Draft, "\r\n") {
			return fault("input-single-line", "Retained draft contains CR or LF")
		}
	}
	return n.validateField(path)
}
