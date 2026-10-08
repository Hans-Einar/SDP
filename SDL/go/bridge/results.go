package bridge

import (
	"fmt"
	"unicode/utf8"

	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	sdl "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/runtime"
)

func (b *Bridge) validateResult(plan Plan, owner ui.Widget, target string, output parser.RecordType) error {
	switch plan.ResultMode {
	case TextResult:
		if owner.Handle.Kind == "dialog" {
			return fmt.Errorf("dialog callback requires DialogAcceptResult")
		}
		if plan.AcceptField != "" || plan.MessageField != "" {
			return fmt.Errorf("accept fields require DialogAcceptResult")
		}
		if target == "" || plan.OutputField == "" || output[plan.OutputField] != parser.TextType {
			return fmt.Errorf("text result requires explicit setHandle and text OutputField")
		}
		receiver, ok := b.UI.Widget(target)
		if !ok || receiver.Handle.Kind != "input" {
			return fmt.Errorf("result-widget: %s is not an input", target)
		}
		if plan.RevisionField != "" {
			if output[plan.RevisionField] != parser.IntegerType || plan.RevisionContext == "" {
				return fmt.Errorf("revision-binding: invalid output/context")
			}
		} else if plan.RevisionContext != "" {
			return fmt.Errorf("revision-binding: missing result field")
		}
	case DialogAcceptResult:
		fields, err := b.UI.DialogFields(owner.Handle)
		if err != nil {
			return err
		}
		if len(fields) > 256 {
			return fmt.Errorf("dialog-fields: capture exceeds 256-write limit")
		}
		if owner.Handle.Kind != "dialog" {
			return fmt.Errorf("DialogAcceptResult requires a dialog Accept owner")
		}
		if target != "" || plan.OutputField != "" || plan.RevisionField != "" || plan.RevisionContext != "" {
			return fmt.Errorf("DialogAcceptResult forbids setHandle/output/revision receivers")
		}
		if plan.AcceptField == "" || plan.MessageField == "" || output[plan.AcceptField] != parser.BooleanType || output[plan.MessageField] != parser.TextType {
			return fmt.Errorf("DialogAcceptResult requires boolean AcceptField and text MessageField")
		}
	default:
		return fmt.Errorf("unknown result mode %q", plan.ResultMode)
	}
	return nil
}

// Engine.Execute validates the complete declared record. This conversion also
// enforces the adapter's closed result mapping before announcing a decision.
func acceptDecision(plan Plan, record sdl.Record) (*ui.AcceptDecision, error) {
	accepted, aok := record[plan.AcceptField]
	message, mok := record[plan.MessageField]
	if !aok || !mok || accepted.Type != parser.BooleanType || message.Type != parser.TextType || !utf8.ValidString(message.Text) || len(message.Text) > 32768 {
		return nil, fmt.Errorf("dialog-result: invalid typed acceptance record")
	}
	decision := &ui.AcceptDecision{Accepted: accepted.Boolean, Message: message.Text}
	if !decision.Accepted && decision.Message == "" {
		decision.Message = "Dialog acceptance rejected"
	}
	return decision, nil
}
