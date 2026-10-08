package bridge

import (
	"fmt"

	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	sdl "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/runtime"
)

// The runtime owns normalized page identities and validates dispatch. The bridge
// only consumes that state; it never reconstructs definition scope from paths.
func (b *Bridge) validatePageEvent(event ui.Event) error {
	if event.Kind != ui.ActivatePage || event.Page == nil || event.Collection != nil || event.Split != nil || event.Command != nil || event.Dialog != nil || event.ModelRevision != b.UI.Revision {
		return fmt.Errorf("event-field: invalid tab activation")
	}
	tabs, ok := b.UI.Tabs(event.Handle)
	if !ok || tabs.Selected != event.Page.PreviousID {
		return fmt.Errorf("event-field: stale tabs owner/previous page")
	}
	for _, page := range tabs.Pages {
		if page.ID == event.Page.PageID && page.Handle == event.Page.Page && page.Enabled && page.Visible {
			return nil
		}
	}
	return fmt.Errorf("event-field: invalid direct page target")
}

// Only the dispatcher supplies these envelopes. Recheck the concrete source
// identities before and after Execute; never recapture context or dialog drafts.
func (b *Bridge) validateInteractionEvent(event ui.Event, mode ResultMode) error {
	if event.ModelRevision != b.UI.Revision {
		return fmt.Errorf("stale-event: model changed")
	}
	switch event.Kind {
	case ui.ActivatePage:
		if mode != TextResult {
			return fmt.Errorf("result-mode: tabs require TextResult")
		}
		return b.validatePageEvent(event)
	case ui.InvokeCommand:
		if mode != TextResult || event.Command == nil || event.Page != nil || event.Split != nil || event.Dialog != nil || event.Collection != nil {
			return fmt.Errorf("event-field: invalid command invocation")
		}
		command, ok := b.UI.CommandState(event.Handle)
		if !ok {
			return fmt.Errorf("event-field: stale canonical command")
		}
		context := event.Command.Context
		if command.Context == "none" {
			if context != nil {
				return fmt.Errorf("event-field: unexpected command context")
			}
		} else {
			if context == nil || context.Widget != command.Target || context.ModelRevision != event.ModelRevision {
				return fmt.Errorf("event-field: mismatched command context")
			}
			if command.Context == "item" {
				if context.Item == nil || context.Item.Handle != context.Widget {
					return fmt.Errorf("event-field: missing captured item")
				}
				if err := b.UI.ValidateCollectionTarget(*context.Item); err != nil {
					return err
				}
			} else if context.Item != nil {
				return fmt.Errorf("event-field: unexpected item context")
			}
		}
		if command.Toggle != (event.Command.Checked != nil) {
			return fmt.Errorf("event-field: checked payload does not match command")
		}
	case ui.Accept:
		if mode != DialogAcceptResult || event.Dialog == nil || event.Command != nil || event.Page != nil || event.Split != nil || event.Collection != nil {
			return fmt.Errorf("event-field: invalid dialog acceptance")
		}
		surface, ok := b.UI.SurfaceState(event.Handle)
		if !ok || !surface.Open || surface.Target != event.Dialog.Surface || event.Dialog.Surface.Handle != event.Handle {
			return fmt.Errorf("event-field: stale dialog opening")
		}
	default:
		return fmt.Errorf("event-field: unsupported interaction callback %s", event.Kind)
	}
	return nil
}

// InteractionReply keeps the domain outcome even when a subsequent UI/adapter
// check fails. Basic buttons continue using the unchanged legacy Handler.
func (b *Bridge) interactionHandler(engine *sdl.Engine, action string, plan Plan, target string) ui.InteractionHandler {
	boundRevision := engine.Revision()
	return func(event ui.Event) (ui.InteractionReply, error) {
		reply := ui.InteractionReply{Domain: ui.DomainNotCalled}
		if err := b.validateInteractionEvent(event, plan.ResultMode); err != nil {
			return reply, err
		}
		var receiver ui.Widget
		if plan.ResultMode == TextResult {
			var ok bool
			receiver, ok = b.UI.Widget(target)
			if !ok {
				return reply, fmt.Errorf("stale-result-target: %s", target)
			}
		}
		inputType, _, err := engine.Signature(action)
		if err != nil {
			return reply, err
		}
		record := sdl.Record{}
		for field, kind := range inputType {
			value, err := b.value(plan.Inputs[field], event, kind)
			if err != nil {
				return reply, err
			}
			record[field] = value
		}
		// Execute's errors cannot distinguish an effect followed by failure from
		// rejection before Go entry. Preserve uncertainty; never replay automatically.
		reply.Domain = ui.DomainUnknown
		result, err := engine.Execute(b.ctx, sdl.Request{Action: action, Revision: boundRevision, Sequence: event.Sequence, Input: record})
		if err != nil {
			return reply, err
		}
		if plan.ResultMode == DialogAcceptResult {
			decision, err := acceptDecision(plan, result.Output)
			if err != nil {
				return reply, err
			}
			reply.Accept = decision
			reply.Domain = ui.DomainRejected
			if decision.Accepted {
				reply.Domain = ui.DomainSucceeded
			}
		} else {
			reply.Domain = ui.DomainSucceeded
		}
		// Classify the validated domain record BEFORE any post-Execute state checks.
		if engine.Revision() != boundRevision {
			return reply, fmt.Errorf("stale-result: SDL revision changed")
		}
		if err := b.validateInteractionEvent(event, plan.ResultMode); err != nil {
			return reply, err
		}
		if plan.ResultMode == DialogAcceptResult {
			return reply, nil
		}
		widget, ok := b.UI.Widget(target)
		if !ok || widget.Handle != receiver.Handle || widget.ValueRevision != receiver.ValueRevision || widget.DraftRevision != receiver.DraftRevision {
			return reply, fmt.Errorf("stale-result-target: %s changed during action", target)
		}
		if plan.RevisionField != "" {
			b.Context[plan.RevisionContext] = result.Output[plan.RevisionField]
		}
		value := result.Output[plan.OutputField].Text
		reply.Updates = []ui.Update{{Handle: widget.Handle, Property: ui.AcceptedValue, Value: ui.Text(value), ExpectedValueRevision: widget.ValueRevision, AcceptDraft: value == widget.Draft}}
		return reply, nil
	}
}
