package runtime

import "reflect"

func (s *Session) validateExclusive() error {
	seen := map[string]bool{}
	for _, c := range s.commands {
		if c.Checked && c.Exclusive != "" {
			key := c.ExclusiveScope + "\x00" + c.Exclusive
			if seen[key] {
				return fault("exclusive", "Multiple checked commands")
			}
			seen[key] = true
		}
	}
	return nil
}
func cloneInteraction(e Event) Event {
	if e.Command != nil {
		v := *e.Command
		v.Surface = copySurface(v.Surface)
		v.Context = copyContext(v.Context)
		if v.Checked != nil {
			x := *v.Checked
			v.Checked = &x
		}
		e.Command = &v
	}
	if e.Dialog != nil {
		v := *e.Dialog
		v.Fields = append([]DraftField(nil), v.Fields...)
		e.Dialog = &v
	}
	return e
}
func (s *Session) commandEnvelope(e Event) error {
	if s.closed {
		return fault("closed", "Session closed")
	}
	if e.ModelRevision != s.Revision || e.StateRevision != s.StateRevision {
		return fault("stale-event", "Interaction state changed")
	}
	if e.Sequence == 0 || e.Sequence <= s.sequence {
		return fault("duplicate-event", "Event sequence consumed")
	}
	if e.Page != nil || e.Split != nil || e.Collection != nil || e.Value != (Value{}) || e.DraftRevision != 0 {
		return fault("event-type", "Unexpected payload")
	}
	var expected Event
	var err error
	if e.Kind == InvokeCommand {
		if e.Command == nil || e.Dialog != nil {
			return fault("event-type", "Expected command payload only")
		}
		expected, err = s.CaptureCommand(e.Command.Origin, e.Command.Via, e.Command.Context)
	} else {
		if e.Dialog == nil || e.Command != nil {
			return fault("event-type", "Expected dialog payload only")
		}
		expected, err = s.CaptureDialog(e.Dialog.Surface, e.Kind)
	}
	if err != nil {
		return err
	}
	expected.Sequence = e.Sequence
	if !reflect.DeepEqual(e, expected) {
		return fault("event-type", "Payload differs from runtime capture")
	}
	return nil
}
func (s *Session) dispatchCommand(e Event) (result InteractionResult, err error) {
	result = InteractionResult{Sequence: e.Sequence, Status: "rejected", Domain: DomainNotCalled}
	if s.interacting {
		return result, fault("reentrant-interaction", "Interaction already active")
	}
	s.interacting = true
	defer func() { s.interacting = false }()
	var menu *MenuScope
	if e.Command != nil && e.Command.Via == "menu" {
		if m := s.menus[s.menuOwner[e.Command.Origin.Path]]; m != nil && m.Open && m.Scope.StateRevision == e.StateRevision && m.Scope.ModelRevision == e.ModelRevision {
			scope := m.Scope
			menu = &scope
		}
	}
	defer func() {
		if result.Status != "committed" && menu != nil {
			_ = s.CloseMenu(*menu)
		}
	}()
	if err = s.commandEnvelope(e); err != nil {
		return result, err
	}
	e = cloneInteraction(e)
	n := s.copyState()
	handlerEvent := e
	reserved := map[Handle]bool{}
	var handler InteractionHandler
	var accept *SurfaceTarget
	if e.Kind == InvokeCommand {
		c := n.commands[e.Handle.Path]
		if c.Toggle {
			if c.Exclusive != "" && c.Checked {
				if menu != nil {
					n.menus[menu.Handle.Path].Open = false
					n.sequence = e.Sequence
					err = s.publish(n)
				}
				if err == nil {
					result.Status = "committed"
				}
				return result, err
			}
			for _, peer := range n.commands {
				if peer.Handle == c.Handle || c.Exclusive != "" && peer.Exclusive == c.Exclusive && peer.ExclusiveScope == c.ExclusiveScope {
					reserved[peer.Handle] = true
					peer.Checked = peer.Handle == c.Handle && *e.Command.Checked
				}
			}
		}
		switch c.Effect {
		case "":
			handler = s.interactions[c.Handle.Path]
			if handler == nil && (c.Binding.Module != "" || !c.Toggle) {
				return result, fault("unbound", c.Handle.Path)
			}
		case "open":
			context := ContextTarget{}
			if e.Command.Context != nil {
				context = *e.Command.Context
			}
			opener := e.Command.Origin
			if menu != nil {
				if m := s.menus[menu.Handle.Path]; m != nil && m.Opener != (Handle{}) {
					opener = m.Opener
				}
			}
			_, err = n.openSurface(c.Target, opener, context)
		case "accept", "cancel", "close":
			d := s.surfaceFor(s.currentControl(e.Command.Origin.Path))
			if d == nil || !d.Open {
				return result, fault("surface", "Effect requires open dialog")
			}
			handlerEvent, err = s.CaptureDialog(d.Target, EventKind(c.Effect))
			handlerEvent.Sequence = e.Sequence
		default:
			return result, fault("command-effect", "Unknown effect")
		}
		if err != nil {
			return result, err
		}
	}
	if handlerEvent.Kind == Accept || handlerEvent.Kind == Cancel || handlerEvent.Kind == Close {
		t := handlerEvent.Dialog.Surface
		if handlerEvent.Kind == Accept {
			accept = &t
			handler = s.interactions[t.Handle.Path]
			if handler == nil && s.aux[t.Handle.Path].Binding.Module != "" {
				return result, fault("unbound", t.Handle.Path)
			}
			for _, f := range handlerEvent.Dialog.Fields {
				reserved[f.Handle] = true
			}
		} else {
			n.closeSurface(t, string(handlerEvent.Kind), "user", e.Sequence, nil)
			handler = nil
		}
	}
	if menu != nil {
		n.menus[menu.Handle.Path].Open = false
	}
	probe := n.copyState()
	probe.refreshActivity()
	probe.closeHiddenSurfaces()
	probe.refreshActivity()
	if err = probe.gate(); err != nil {
		return result, err
	}
	if s.closed || s.Revision != e.ModelRevision || s.StateRevision != e.StateRevision {
		return result, fault("stale-event", "State changed during probe")
	}
	captured := s.captureControls()
	var reply InteractionReply
	baseline := s.StateRevision
	if handler != nil {
		s.sequence = e.Sequence
		s.StateRevision++
		baseline = s.StateRevision
		if accept != nil {
			s.accepting = copySurface(accept)
			s.noteAccept(*accept, e.Sequence, DomainUnknown, false, "")
			defer func() { s.accepting = nil }()
		}
		reply, err = handler(cloneInteraction(handlerEvent))
		result.Domain = reply.Domain
		switch result.Domain {
		case DomainNotCalled, DomainSucceeded, DomainRejected, DomainUnknown:
		default:
			result.Domain = DomainUnknown
		}
	} else if accept != nil {
		reply.Accept = &AcceptDecision{Accepted: true}
		reply.Domain = DomainNotCalled
	}
	fail := func(cause error) (InteractionResult, error) {
		if result.Domain == DomainSucceeded || result.Domain == DomainUnknown {
			result.Status = "ui-conflict"
		}
		if accept != nil {
			s.noteAccept(*accept, e.Sequence, result.Domain, result.Domain == DomainSucceeded || result.Domain == DomainUnknown, boundedError(cause))
		}
		return result, cause
	}
	if err != nil {
		return fail(err)
	}
	if s.closed || s.Revision != e.ModelRevision || s.StateRevision != baseline {
		return fail(fault("stale-result", "State changed during callback"))
	}
	if accept == nil && reply.Accept != nil {
		return fail(fault("interaction-reply", "Unexpected acceptance decision"))
	}
	if handler != nil && (result.Domain == DomainUnknown || result.Domain == DomainRejected) && accept == nil {
		return fail(fault("interaction-rejected", "Command did not report success"))
	}
	if accept != nil {
		if reply.Accept == nil || !validDraft(reply.Accept.Message) {
			return fail(fault("interaction-reply", "Accept requires bounded decision"))
		}
		if !reply.Accept.Accepted {
			if len(reply.Updates) != 0 || result.Domain != DomainRejected {
				return fail(fault("interaction-reply", "False decision requires rejected outcome and no updates"))
			}
			s.noteAccept(*accept, e.Sequence, result.Domain, false, reply.Accept.Message)
			return result, nil
		}
		if handler != nil && result.Domain != DomainSucceeded {
			return fail(fault("interaction-reply", "Accepted requires succeeded outcome"))
		}
	}
	for _, u := range reply.Updates {
		old, ok := captured[u.Handle]
		live := s.currentControl(u.Handle.Path)
		if !ok || live == nil || live.Handle != old.Handle || live.ValueRevision != old.ValueRevision || live.DraftRevision != old.DraftRevision {
			return fail(fault("stale-result", "Receiver changed or was not captured"))
		}
		if reserved[u.Handle] {
			return fail(fault("reserved-update", "Reply overlaps reserved change"))
		}
	}
	updates := append([]Update(nil), reply.Updates...)
	if accept != nil {
		for _, f := range handlerEvent.Dialog.Fields {
			updates = append(updates, Update{Handle: f.Handle, Property: AcceptedValue, Value: f.Value, ExpectedValueRevision: f.ValueRevision, AcceptDraft: true})
		}
	}
	if len(updates) > 256 {
		return fail(fault("update-limit", "Fields and updates exceed 256 writes"))
	}
	n.sequence = e.Sequence
	n.StateRevision = baseline
	if len(updates) > 0 {
		n, err = n.applyCandidate(s.Revision, s.BatchRevision+1, updates)
		if err != nil {
			return fail(err)
		}
	}
	if accept != nil {
		n.noteAccept(*accept, e.Sequence, result.Domain, false, reply.Accept.Message)
		n.closeSurface(*accept, "accept", "user", e.Sequence, handlerEvent.Dialog.Fields)
	}
	if err = s.publish(n); err != nil {
		return fail(err)
	}
	result.Status = "committed"
	return result, nil
}

// Attempt diagnostics never prepare or promote native tickets.
func (s *Session) noteAccept(t SurfaceTarget, seq uint64, domain DomainOutcome, blocked bool, message string) {
	if d := s.surfaces[t.Handle.Path]; d != nil && d.Open && d.Target == t {
		d.AcceptSequence = seq
		d.Domain = domain
		d.AcceptBlocked = blocked
		d.Message = message
	}
	for i := range s.dialogResults {
		r := &s.dialogResults[i]
		if r.Surface == t {
			r.AcceptSequence = seq
			r.Domain = domain
		}
	}
}

// Includes concrete M1/M2 property owners as well as ordinary controls.
func (s *Session) captureControls() map[Handle]Widget {
	out := map[Handle]Widget{}
	for _, m := range []map[string]*Widget{s.widgets, s.panes, s.aux} {
		for _, w := range m {
			out[w.Handle] = *w
		}
	}
	return out
}
