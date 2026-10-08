package bridge

import (
	"fmt"
	"strconv"

	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	sdl "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/runtime"
)

func (b *Bridge) validateSource(source Source, widget ui.Widget, kind parser.ScalarType) error {
	if source.FieldPath != "" && source.EventField != DialogFieldValue {
		return fmt.Errorf("FieldPath is only valid with DialogFieldValue")
	}
	count := 0
	if source.EventField != "" {
		count++
		switch source.EventField {
		case ControlBoolean, ControlNumber, ControlText, ChoiceOptionID:
			if err := b.validateControlSource(source.EventField, widget, kind); err != nil {
				return err
			}
		case CollectionItemID:
			if (widget.Handle.Kind != "tree" && widget.Handle.Kind != "list") || kind != parser.TextType {
				return fmt.Errorf("event-field: collection.item-id requires a collection Activate and text destination")
			}
		case TabPageID, TabPreviousPageID:
			if widget.Handle.Kind != "tabs" || kind != parser.TextType {
				return fmt.Errorf("event-field: tab identity requires a tabs ActivatePage and text destination")
			}
		case CommandContextItemID:
			command, ok := b.UI.CommandState(widget.Handle)
			if !ok || command.Context != "item" || kind != parser.TextType {
				return fmt.Errorf("event-field: context item ID requires an item-context command and text")
			}
		case CommandChecked:
			command, ok := b.UI.CommandState(widget.Handle)
			if !ok || !command.Toggle || kind != parser.BooleanType {
				return fmt.Errorf("event-field: checked requires a toggle command and boolean")
			}
		case DialogFieldValue:
			if widget.Handle.Kind != "dialog" || kind != parser.TextType || source.FieldPath == "" {
				return fmt.Errorf("event-field: dialog field requires an Accept owner, named FieldPath and text")
			}
		default:
			return fmt.Errorf("event-field: unknown selector %q", source.EventField)
		}
	}
	if source.Widget != "" {
		if widget.Handle.Kind == "dialog" {
			return fmt.Errorf("dialog input requires captured DialogFieldValue, not a live Widget source")
		}
		count++
		w, ok := b.UI.Widget(source.Widget)
		if !ok || w.Handle.Kind != "input" {
			return fmt.Errorf("widget source must name an input")
		}
	}
	if source.Event {
		count++
		if widget.Handle.Kind != "input" {
			return fmt.Errorf("event value requires an input callback")
		}
	}
	if source.Literal != nil {
		count++
		if err := sdl.ValidateRecord(sdl.Record{"Value": *source.Literal}, parser.RecordType{"Value": kind}); err != nil {
			return err
		}
	}
	if source.Context != "" {
		count++
		v, ok := b.Context[source.Context]
		if !ok {
			return fmt.Errorf("missing context value %s", source.Context)
		}
		if err := sdl.ValidateRecord(sdl.Record{"Value": v}, parser.RecordType{"Value": kind}); err != nil {
			return err
		}
	}
	if count != 1 {
		return fmt.Errorf("exactly one explicit value source required")
	}
	return nil
}
func (b *Bridge) value(source Source, event ui.Event, kind parser.ScalarType) (sdl.Value, error) {
	if source.EventField != "" {
		switch source.EventField {
		case ControlBoolean, ControlNumber, ControlText, ChoiceOptionID:
			return b.controlValue(source.EventField, event, kind)
		case CommandContextItemID:
			if kind != parser.TextType || event.Kind != ui.InvokeCommand || event.Command == nil || event.Command.Context == nil || event.Command.Context.Item == nil {
				return sdl.Value{}, fmt.Errorf("event-field: missing command item context")
			}
			if err := b.UI.ValidateCollectionTarget(*event.Command.Context.Item); err != nil {
				return sdl.Value{}, err
			}
			return sdl.Text(string(event.Command.Context.Item.ItemID)), nil
		case CommandChecked:
			if kind != parser.BooleanType || event.Kind != ui.InvokeCommand || event.Command == nil || event.Command.Checked == nil {
				return sdl.Value{}, fmt.Errorf("event-field: missing proposed checked value")
			}
			return sdl.Boolean(*event.Command.Checked), nil
		case DialogFieldValue:
			if kind != parser.TextType || event.Kind != ui.Accept || event.Dialog == nil {
				return sdl.Value{}, fmt.Errorf("event-field: missing captured dialog request")
			}
			for _, field := range event.Dialog.Fields {
				if field.Handle == source.resolvedField {
					if field.Value.Kind != ui.String {
						return sdl.Value{}, fmt.Errorf("event-field: captured dialog field is not String")
					}
					return sdl.Text(field.Value.Text), nil
				}
			}
			return sdl.Value{}, fmt.Errorf("event-field: owned field absent from captured request")
		}
		if source.EventField == TabPageID || source.EventField == TabPreviousPageID {
			if kind != parser.TextType {
				return sdl.Value{}, fmt.Errorf("event-field: tab identity requires text")
			}
			if err := b.validatePageEvent(event); err != nil {
				return sdl.Value{}, err
			}
			if source.EventField == TabPreviousPageID {
				return sdl.Text(event.Page.PreviousID), nil
			}
			return sdl.Text(event.Page.PageID), nil
		}
		if source.EventField != CollectionItemID || kind != parser.TextType || event.Kind != ui.Activate || event.Collection == nil {
			return sdl.Value{}, fmt.Errorf("event-field: invalid collection activation")
		}
		if err := b.UI.ValidateCollectionTarget(*event.Collection); err != nil {
			return sdl.Value{}, err
		}
		return sdl.Text(string(event.Collection.ItemID)), nil
	}
	if source.Literal != nil {
		return *source.Literal, nil
	}
	if source.Context != "" {
		return b.Context[source.Context], nil
	}
	text := event.Value.Text
	if source.Widget != "" {
		w, ok := b.UI.Widget(source.Widget)
		if !ok {
			return sdl.Value{}, fmt.Errorf("stale-widget-source: %s", source.Widget)
		}
		text = w.Draft
	}
	switch kind {
	case parser.TextType:
		return sdl.Text(text), nil
	case parser.IntegerType:
		v, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			return sdl.Value{}, fmt.Errorf("input-integer: %w", err)
		}
		return sdl.Integer(v), nil
	case parser.BooleanType:
		if text == "true" {
			return sdl.Boolean(true), nil
		}
		if text == "false" {
			return sdl.Boolean(false), nil
		}
		return sdl.Value{}, fmt.Errorf("input-boolean: expected true or false")
	}
	return sdl.Value{}, fmt.Errorf("input-type: %s", kind)
}
func (b *Bridge) handler(engine *sdl.Engine, action string, plan Plan, target string) ui.Handler {
	boundRevision := engine.Revision()
	return func(event ui.Event) ([]ui.Update, error) {
		receiver, ok := b.UI.Widget(target)
		if !ok {
			return nil, fmt.Errorf("stale-result-target: %s", target)
		}
		var capture *ui.FieldTarget
		if b.extendedInput(event.Handle) {
			field, err := b.controlCapture(event)
			if err != nil {
				return nil, err
			}
			if receiver.Handle != event.Handle || event.Value != ui.Text(event.Value.Text) {
				return nil, fmt.Errorf("text-result: invalid self text capture")
			}
			capture = &field.Target
		}
		inputType, _, err := engine.Signature(action)
		if err != nil {
			return nil, err
		}
		record := sdl.Record{}
		for field, kind := range inputType {
			value, err := b.value(plan.Inputs[field], event, kind)
			if err != nil {
				return nil, err
			}
			record[field] = value
		}
		result, err := engine.Execute(b.ctx, sdl.Request{Action: action, Revision: boundRevision, Sequence: event.Sequence, Input: record})
		if err != nil {
			return nil, err
		}
		widget, ok := b.UI.Widget(target)
		if !ok || widget.Handle != receiver.Handle || widget.ValueRevision != receiver.ValueRevision || widget.DraftRevision != receiver.DraftRevision {
			return nil, fmt.Errorf("stale-result-target: %s changed during action", target)
		}
		if engine.Revision() != boundRevision {
			return nil, fmt.Errorf("stale-result: SDL revision changed")
		}
		if event.Collection != nil {
			if err := b.UI.ValidateCollectionTarget(*event.Collection); err != nil {
				return nil, err
			}
		}
		value := result.Output[plan.OutputField].Text
		if capture != nil {
			after, err := b.controlCapture(event)
			if err != nil {
				return nil, err
			}
			if after.Target != *capture || value != event.Value.Text {
				return nil, fmt.Errorf("text-result: result does not accept unchanged captured proposal")
			}
		}
		if plan.RevisionField != "" {
			b.Context[plan.RevisionContext] = result.Output[plan.RevisionField]
		}
		update := ui.Update{Handle: widget.Handle, Property: ui.AcceptedValue, Value: ui.Text(value), ExpectedValueRevision: widget.ValueRevision, AcceptDraft: value == widget.Draft}
		if b.extendedInput(widget.Handle) {
			update.ExpectedDraftRevision = receiver.DraftRevision
		}
		return []ui.Update{update}, nil
	}
}

// Runtime policy is present only for source-explicit extended inputs.
func (b *Bridge) extendedInput(h ui.Handle) bool {
	field, ok := b.UI.Field(h)
	return ok && h.Kind == "input" && field.Input != nil
}
