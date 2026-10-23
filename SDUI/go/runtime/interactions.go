package runtime

type PageActivation struct {
	PreviousID, PageID string
	Page               Handle
}
type SplitChange struct {
	Operation  string
	Proportion float64
}
type DomainOutcome string

const (
	DomainNotCalled DomainOutcome = "not-called"
	DomainSucceeded DomainOutcome = "succeeded"
	DomainRejected  DomainOutcome = "rejected"
	DomainUnknown   DomainOutcome = "unknown"
)

type InteractionReply struct {
	Accept  *AcceptDecision
	Updates []Update
	Domain  DomainOutcome
}
type InteractionHandler func(Event) (InteractionReply, error)
type InteractionResult struct {
	Sequence uint64
	Status   string
	Domain   DomainOutcome
}

func (s *Session) BindInteraction(h Handle, handler InteractionHandler) error {
	w, err := s.lookupControl(h)
	if err != nil {
		return err
	}
	if (w.Handle.Kind != "tabs" && s.commands[h.Path] == nil && s.surfaces[h.Path] == nil) || handler == nil {
		return fault("binding", "Interaction handler requires tabs")
	}
	if c := s.commands[h.Path]; c != nil && c.Effect != "" {
		return fault("binding", "Local effects do not own handlers")
	}
	s.interactions[h.Path] = handler
	s.StateRevision++
	return nil
}
func (s *Session) HasInteractionBinding(h Handle) bool {
	if _, err := s.lookupControl(h); err != nil {
		return false
	}
	return s.interactions[h.Path] != nil
}
func (s *Session) interactionCandidate(e Event) (*Session, bool, error) {
	w, err := s.pane(e.Handle)
	if err != nil {
		return nil, false, err
	}
	if !s.inputAllowed(w) {
		return nil, false, fault("inactive-widget", w.Handle.Path)
	}
	if e.ModelRevision != s.Revision || e.StateRevision != s.StateRevision {
		return nil, false, fault("stale-event", "Interaction state is no longer current")
	}
	if e.Sequence == 0 || e.Sequence <= s.sequence {
		return nil, false, fault("duplicate-event", "Event sequence already consumed")
	}
	if e.Control != nil || e.Command != nil || e.Dialog != nil || e.Collection != nil || e.Value != (Value{}) || e.DraftRevision != 0 {
		return nil, false, fault("event-type", "Unexpected legacy payload")
	}
	n := s.copyState()
	switch e.Kind {
	case ActivatePage:
		if e.Page == nil || e.Split != nil || w.Handle.Kind != "tabs" {
			return nil, false, fault("event-type", "ActivatePage requires only page payload")
		}
		t := s.tabs[w.Handle.Path]
		if e.Page.PreviousID != t.Selected {
			return nil, false, fault("page-selection", "Previous page differs")
		}
		found := false
		for _, p := range t.Pages {
			if p.ID == e.Page.PageID && p.Handle == e.Page.Page {
				found = true
			}
		}
		if !found {
			return nil, false, fault("page-selection", "Page target differs from named child")
		}
		if err = n.selectPage(e.Handle, e.Page.PageID, true); err != nil {
			return nil, false, err
		}
		return n, t.Selected == e.Page.PageID, nil
	case AdjustSplit:
		if e.Split == nil || e.Page != nil || w.Handle.Kind != "split" {
			return nil, false, fault("event-type", "AdjustSplit requires only split payload")
		}
		if err = n.changeSplit(e.Handle, *e.Split, false); err != nil {
			return nil, false, err
		}
		n.focused = w.Handle.Path
		return n, false, nil
	default:
		return nil, false, fault("event-type", "Unknown interaction kind")
	}
}

// DispatchInteraction stages pane state and callback updates as one publication.
// Its geometry probe never prepares or promotes a native presentation ticket.
func (s *Session) DispatchInteraction(e Event) (result InteractionResult, err error) {
	if e.Kind == InvokeCommand || e.Kind == Accept || e.Kind == Cancel || e.Kind == Close {
		return s.dispatchCommand(e)
	}
	result = InteractionResult{Sequence: e.Sequence, Status: "rejected", Domain: DomainNotCalled}
	if s.interacting {
		return result, fault("reentrant-interaction", "Interaction dispatch is already active")
	}
	s.interacting = true
	defer func() { s.interacting = false }()
	n, same, err := s.interactionCandidate(e)
	if err != nil {
		return result, err
	}
	if same {
		result.Status = "committed"
		return result, nil
	}
	handler := s.interactions[e.Handle.Path]
	if e.Kind == ActivatePage && handler == nil && s.panes[e.Handle.Path].Binding.Module != "" {
		return result, fault("unbound", e.Handle.Path)
	}
	if handler == nil || e.Kind == AdjustSplit {
		n.sequence = e.Sequence
		err = s.publish(n)
		if err == nil {
			result.Status = "committed"
		}
		return result, err
	}
	probe := n.copyState()
	probe.refreshActivity()
	if err = probe.gate(); err != nil {
		return result, err
	}
	if s.closed || s.Revision != e.ModelRevision || s.StateRevision != e.StateRevision {
		return result, fault("stale-event", "State changed during interaction probe")
	}
	captured := s.captureControls()
	s.sequence = e.Sequence
	s.StateRevision++
	baseline := s.StateRevision
	if e.Page != nil {
		page := *e.Page
		e.Page = &page
	}
	reply, callErr := handler(e)
	result.Domain = reply.Domain
	switch result.Domain {
	case "":
		result.Domain = DomainUnknown
	case DomainNotCalled, DomainSucceeded, DomainRejected, DomainUnknown:
	default:
		result.Domain = DomainUnknown
		result.Status = "ui-conflict"
		return result, fault("interaction-reply", "Unknown domain outcome")
	}
	fail := func(err error) (InteractionResult, error) {
		if result.Domain == DomainSucceeded || result.Domain == DomainUnknown {
			result.Status = "ui-conflict"
		}
		return result, err
	}
	if reply.Accept != nil {
		return fail(fault("interaction-reply", "Accept is only valid for dialog acceptance"))
	}
	if callErr != nil {
		return fail(callErr)
	}
	if result.Domain == DomainRejected {
		return fail(fault("interaction-rejected", "Handler rejected interaction"))
	}
	if s.closed || s.Revision != e.ModelRevision || s.StateRevision != baseline {
		return fail(fault("stale-result", "State changed while callback executed"))
	}

	for _, u := range reply.Updates {
		if u.Handle == e.Handle {
			return fail(fault("reserved-update", "Reply overlaps pane owner"))
		}
		if tabs := s.tabs[e.Handle.Path]; tabs != nil {
			for _, page := range tabs.Pages {
				if u.Handle == page.Handle {
					return fail(fault("reserved-update", "Reply overlaps activation pages"))
				}
			}
		}
		old, ok := captured[u.Handle]
		live := s.currentControl(u.Handle.Path)
		if !ok || live == nil || live.Handle != old.Handle || live.ValueRevision != old.ValueRevision || live.DraftRevision != old.DraftRevision {
			return fail(fault("stale-result", "Update target was not captured or changed"))
		}
	}
	n.sequence = e.Sequence
	n.StateRevision = baseline
	next, err := n.applyCandidate(s.Revision, s.BatchRevision+1, reply.Updates)
	if err != nil {
		return fail(err)
	}
	if s.closed || s.Revision != e.ModelRevision || s.StateRevision != baseline {
		return fail(fault("stale-result", "State changed during field validation"))
	}
	if err = s.publish(next); err != nil {
		return fail(err)
	}
	result.Status = "committed"
	return result, nil
}
