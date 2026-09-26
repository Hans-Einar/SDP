package bridge

import (
	"fmt"
	"strconv"

	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	sdl "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/runtime"
)

func (b *Bridge) validateSource(source Source, widget ui.Widget, kind parser.ScalarType) error {
	count := 0
	if source.Widget != "" {
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
		if !ok {
			return nil, fmt.Errorf("stale-result-target: %s", target)
		}
		if plan.RevisionField != "" {
			b.Context[plan.RevisionContext] = result.Output[plan.RevisionField]
		}
		value := result.Output[plan.OutputField].Text
		return []ui.Update{{Handle: widget.Handle, Property: ui.AcceptedValue, Value: ui.Text(value), ExpectedValueRevision: widget.ValueRevision, AcceptDraft: value == widget.Draft}}, nil
	}
}
