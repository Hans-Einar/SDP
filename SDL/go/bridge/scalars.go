package bridge

import (
	"fmt"

	"github.com/Hans-Einar/SDP/SDUI/go/numeric"
	uiparser "github.com/Hans-Einar/SDP/SDUI/go/parser"
	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	sdl "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/runtime"
)

func scalarKind(kind string) bool {
	return kind == "checkbox" || kind == "slider" || kind == "number" || kind == "select"
}
func scalarType(kind string) parser.ScalarType {
	switch kind {
	case "checkbox":
		return parser.BooleanType
	case "slider", "number":
		return parser.IntegerType
	case "select":
		return parser.TextType
	}
	return ""
}

// Source lexemes are frontend facts. No reparse, formatting-based provenance or
// second decimal/grid implementation belongs in the bridge.
func (b *Bridge) scalarGrid(w ui.Widget) (*numeric.Grid, error) {
	var node *uiparser.Instance
	b.UI.SnapshotRoot().Walk(func(n *uiparser.Instance) {
		if n.Path == w.InstancePath {
			node = n
		}
	})
	min, max, step, initial, err := uiparser.NumericArguments(node)
	if err != nil {
		return nil, err
	}
	wrap := func(err error) error {
		if err == nil {
			return nil
		}
		return fmt.Errorf("scalar-capability %s at %d:%d: %w", w.InstancePath, node.Span.Line, node.Span.Column, err)
	}
	grid, err := numeric.NewGrid(min, max, step)
	if err != nil {
		return nil, wrap(err)
	}
	if err = grid.ValidateSDL(); err != nil {
		return nil, wrap(err)
	}
	value, err := grid.Parse(initial)
	if err != nil {
		return nil, wrap(err)
	}
	integer, err := numeric.ExactInteger(initial)
	if err != nil {
		return nil, wrap(err)
	}
	if float64(integer) != value {
		return nil, wrap(fmt.Errorf("initial integer round trip failed"))
	}
	return grid, nil
}
func (b *Bridge) validateControlSource(selector EventField, w ui.Widget, kind parser.ScalarType) error {
	switch selector {
	case ControlBoolean:
		if w.Handle.Kind == "checkbox" && kind == parser.BooleanType {
			return nil
		}
	case ControlNumber:
		if (w.Handle.Kind == "number" || w.Handle.Kind == "slider") && kind == parser.IntegerType {
			_, err := b.scalarGrid(w)
			return err
		}
	case ChoiceOptionID:
		if w.Handle.Kind == "select" && kind == parser.TextType {
			field, ok := b.UI.Field(w.Handle)
			if !ok || field.Target.OptionGeneration == 0 {
				return fmt.Errorf("choice-options: unsupplied select %s", w.InstancePath)
			}
			return nil
		}
	case ControlText:
		if w.Handle.Kind == "input" && kind == parser.TextType {
			return nil
		}
	}
	return fmt.Errorf("event-field: %s incompatible with %s/%s", selector, w.Handle.Kind, kind)
}
func (b *Bridge) validateScalarResult(plan Plan, owner ui.Widget, target string, output parser.RecordType) error {
	if !scalarKind(owner.Handle.Kind) {
		return fmt.Errorf("ScalarResult requires typed scalar Commit")
	}
	receiver, ok := b.UI.Widget(target)
	if !ok || receiver.Handle != owner.Handle {
		return fmt.Errorf("ScalarResult requires explicit self setHandle to accept captured source")
	}
	if plan.OutputField == "" || output[plan.OutputField] != scalarType(receiver.Handle.Kind) {
		return fmt.Errorf("ScalarResult output type must match receiver")
	}
	if plan.AcceptField != "" || plan.MessageField != "" {
		return fmt.Errorf("ScalarResult forbids AcceptField/MessageField")
	}
	if plan.RevisionField != "" {
		if output[plan.RevisionField] != parser.IntegerType || plan.RevisionContext == "" {
			return fmt.Errorf("revision-binding: invalid scalar output/context")
		}
	} else if plan.RevisionContext != "" {
		return fmt.Errorf("revision-binding: missing scalar result field")
	}
	if receiver.Handle.Kind == "number" || receiver.Handle.Kind == "slider" {
		_, err := b.scalarGrid(receiver)
		return err
	}
	if receiver.Handle.Kind == "select" {
		f, ok := b.UI.Field(receiver.Handle)
		if !ok || f.Target.OptionGeneration == 0 {
			return fmt.Errorf("choice-options: missing receiver options")
		}
	}
	return nil
}
func equalRaw(a, b *string) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return *a == *b
}
func (b *Bridge) controlCapture(event ui.Event) (ui.FieldState, error) {
	field, ok := b.UI.Field(event.Handle)
	if !ok || event.Kind != ui.Commit || event.Control == nil || event.ModelRevision != b.UI.Revision || event.Collection != nil || event.Page != nil || event.Split != nil || event.Dialog != nil || event.Command != nil {
		return ui.FieldState{}, fmt.Errorf("control-event: invalid typed Commit")
	}
	if field.Target.Handle != event.Handle || field.Target.ValueRevision != event.Control.ValueRevision || field.Target.DraftRevision != event.DraftRevision || field.Proposed != event.Value || field.Validation.Code != "" {
		return ui.FieldState{}, fmt.Errorf("control-event: stale or invalid proposal")
	}
	if event.Handle.Kind == "number" {
		if event.Control.RawDraft == nil || !equalRaw(field.RawDraft, event.Control.RawDraft) {
			return ui.FieldState{}, fmt.Errorf("control-event: missing/stale numeric provenance")
		}
	} else if event.Control.RawDraft != nil {
		return ui.FieldState{}, fmt.Errorf("control-event: unexpected raw numeric draft")
	}
	if event.Handle.Kind == "select" {
		t := event.Control.Option
		if t == nil || t.Handle != event.Handle || t.ModelRevision != event.ModelRevision || t.OptionGeneration != field.Target.OptionGeneration || t.OptionID != event.Value.OptionID {
			return ui.FieldState{}, fmt.Errorf("control-event: mismatched option capture")
		}
		if err := b.UI.ValidateOptionTarget(*t); err != nil {
			return ui.FieldState{}, err
		}
	} else if event.Control.Option != nil {
		return ui.FieldState{}, fmt.Errorf("control-event: unexpected option target")
	}
	return field, nil
}
func (b *Bridge) controlValue(selector EventField, event ui.Event, kind parser.ScalarType) (sdl.Value, error) {
	if selector == ControlText {
		w, ok := b.UI.Widget(event.Handle.Path)
		if !ok || w.Handle.Kind != "input" || event.Kind != ui.Commit || kind != parser.TextType || event.ModelRevision != b.UI.Revision || event.DraftRevision != w.DraftRevision || event.Value != ui.Text(w.Draft) || event.Control != nil || event.Command != nil || event.Dialog != nil || event.Page != nil || event.Split != nil || event.Collection != nil {
			return sdl.Value{}, fmt.Errorf("control-text: invalid basic input Commit")
		}
		return sdl.Text(event.Value.Text), nil
	}
	if _, err := b.controlCapture(event); err != nil {
		return sdl.Value{}, err
	}
	w, ok := b.UI.Widget(event.Handle.Path)
	if !ok {
		return sdl.Value{}, fmt.Errorf("control-event: missing widget")
	}
	if err := b.validateControlSource(selector, w, kind); err != nil {
		return sdl.Value{}, err
	}
	switch selector {
	case ControlBoolean:
		if event.Value == ui.Bool(event.Value.Bool) {
			return sdl.Boolean(event.Value.Bool), nil
		}
	case ChoiceOptionID:
		if event.Value == ui.Choice(event.Value.OptionID) {
			return sdl.Text(string(event.Value.OptionID)), nil
		}
	case ControlNumber:
		if event.Value != ui.Numeric(event.Value.Number) {
			break
		}
		grid, err := b.scalarGrid(w)
		if err != nil {
			return sdl.Value{}, err
		}
		if _, err = grid.Tick(event.Value.Number); err != nil {
			return sdl.Value{}, err
		}
		integer, err := numeric.Integer(event.Value.Number)
		if err != nil {
			return sdl.Value{}, err
		}
		if event.Handle.Kind == "number" {
			raw := *event.Control.RawDraft
			exact, err := numeric.ExactInteger(raw)
			if err != nil {
				return sdl.Value{}, err
			}
			parsed, err := grid.Parse(raw)
			if err != nil {
				return sdl.Value{}, err
			}
			if exact != integer || parsed != event.Value.Number {
				return sdl.Value{}, fmt.Errorf("control-number: inconsistent raw provenance")
			}
		}
		return sdl.Integer(integer), nil
	}
	return sdl.Value{}, fmt.Errorf("control-value: wrong payload")
}
func (b *Bridge) scalarOutput(w ui.Widget, value sdl.Value) (ui.Value, error) {
	if err := sdl.ValidateRecord(sdl.Record{"Value": value}, parser.RecordType{"Value": scalarType(w.Handle.Kind)}); err != nil {
		return ui.Value{}, err
	}
	switch w.Handle.Kind {
	case "checkbox":
		return ui.Bool(value.Boolean), nil
	case "number", "slider":
		// Test int64 bounds BEFORE converting; ±2^53 can round an out-of-range integer.
		const safe int64 = 9007199254740991
		if value.Integer < -safe || value.Integer > safe {
			return ui.Value{}, fmt.Errorf("scalar-result: integer exceeds safe53")
		}
		number := float64(value.Integer)
		integer, err := numeric.Integer(number)
		if err != nil || integer != value.Integer {
			return ui.Value{}, fmt.Errorf("scalar-result: invalid integer round trip")
		}
		grid, err := b.scalarGrid(w)
		if err != nil {
			return ui.Value{}, err
		}
		if _, err = grid.Tick(number); err != nil {
			return ui.Value{}, err
		}
		return ui.Numeric(number), nil
	case "select":
		target, err := b.UI.Option(w.Handle, ui.ItemID(value.Text))
		if err != nil {
			return ui.Value{}, err
		}
		if err = b.UI.ValidateOptionTarget(target); err != nil {
			return ui.Value{}, err
		}
		return ui.Choice(target.OptionID), nil
	}
	return ui.Value{}, fmt.Errorf("scalar-result: invalid receiver")
}
func (b *Bridge) scalarHandler(engine *sdl.Engine, action string, plan Plan, target string) ui.Handler {
	boundRevision := engine.Revision()
	return func(event ui.Event) ([]ui.Update, error) {
		before, err := b.controlCapture(event)
		if err != nil {
			return nil, err
		}
		receiver, ok := b.UI.Widget(target)
		if !ok || receiver.Handle != event.Handle {
			return nil, fmt.Errorf("scalar-result: stale receiver")
		}
		input, _, err := engine.Signature(action)
		if err != nil {
			return nil, err
		}
		record := sdl.Record{}
		for name, kind := range input {
			v, err := b.value(plan.Inputs[name], event, kind)
			if err != nil {
				return nil, err
			}
			record[name] = v
		}
		result, err := engine.Execute(b.ctx, sdl.Request{Action: action, Revision: boundRevision, Sequence: event.Sequence, Input: record})
		if err != nil {
			return nil, err
		}
		if engine.Revision() != boundRevision {
			return nil, fmt.Errorf("scalar-result: stale SDL revision")
		}
		after, err := b.controlCapture(event)
		if err != nil {
			return nil, err
		}
		if after.Target != before.Target {
			return nil, fmt.Errorf("scalar-result: receiver/state/options changed during action")
		}
		value, err := b.scalarOutput(receiver, result.Output[plan.OutputField])
		if err != nil {
			return nil, err
		}
		if value != before.Proposed {
			return nil, fmt.Errorf("scalar-result: result does not accept captured proposal")
		}
		if plan.RevisionField != "" {
			b.Context[plan.RevisionContext] = result.Output[plan.RevisionField]
		}
		return []ui.Update{{Handle: receiver.Handle, Property: ui.AcceptedValue, Value: value, ExpectedValueRevision: before.Target.ValueRevision, ExpectedDraftRevision: before.Target.DraftRevision, ExpectedOptionGeneration: before.Target.OptionGeneration, AcceptDraft: true}}, nil
	}
}
