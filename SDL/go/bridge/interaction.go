package bridge

import (
	"fmt"

	ui "github.com/Hans-Einar/SDP/SDUI/go/runtime"
	sdl "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/runtime"
)

// The runtime owns normalized page identities and validates dispatch. The bridge
// only consumes that state; it never reconstructs definition scope from paths.
func (b *Bridge) validatePageEvent(event ui.Event) error {
	if event.Kind != ui.ActivatePage || event.Page == nil || event.Collection != nil || event.Split != nil || event.ModelRevision != b.UI.Revision {
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

// Unlike legacy Handler, InteractionReply carries the domain outcome through
// post-Execute adapter failures so a failed UI publication cannot hide success.
func (b *Bridge) interactionHandler(engine *sdl.Engine, action string, plan Plan, target string) ui.InteractionHandler {
	boundRevision := engine.Revision()
	return func(event ui.Event) (ui.InteractionReply, error) {
		reply := ui.InteractionReply{Domain: ui.DomainNotCalled}
		if err := b.validatePageEvent(event); err != nil {
			return reply, err
		}
		receiver, ok := b.UI.Widget(target)
		if !ok {
			return reply, fmt.Errorf("stale-result-target: %s", target)
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
		// Execute may reject before calling Go or fail after an effect; without a
		// stronger execution receipt an error is conservatively unknown, never retried.
		reply.Domain = ui.DomainUnknown
		result, err := engine.Execute(b.ctx, sdl.Request{Action: action, Revision: boundRevision, Sequence: event.Sequence, Input: record})
		if err != nil {
			return reply, err
		}
		reply.Domain = ui.DomainSucceeded
		widget, ok := b.UI.Widget(target)
		if !ok || widget.Handle != receiver.Handle || widget.ValueRevision != receiver.ValueRevision || widget.DraftRevision != receiver.DraftRevision {
			return reply, fmt.Errorf("stale-result-target: %s changed during action", target)
		}
		if engine.Revision() != boundRevision {
			return reply, fmt.Errorf("stale-result: SDL revision changed")
		}
		if err := b.validatePageEvent(event); err != nil {
			return reply, err
		}
		if plan.RevisionField != "" {
			b.Context[plan.RevisionContext] = result.Output[plan.RevisionField]
		}
		value := result.Output[plan.OutputField].Text
		reply.Updates = []ui.Update{{Handle: widget.Handle, Property: ui.AcceptedValue, Value: ui.Text(value), ExpectedValueRevision: widget.ValueRevision, AcceptDraft: value == widget.Draft}}
		return reply, nil
	}
}
