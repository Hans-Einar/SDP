package runtime

func (s *Session) Dispatch(event Event) error {
	if event.Command != nil || event.Dialog != nil || event.Page != nil || event.Split != nil || event.StateRevision != 0 {
		return fault("event-type", "Unexpected interaction payload")
	}
	w, err := s.lookup(event.Handle)
	if err != nil {
		return err
	}
	if event.ModelRevision != s.Revision {
		return fault("stale-event", "Event belongs to another model revision")
	}
	if event.Sequence == 0 || event.Sequence <= s.sequence {
		return fault("duplicate-event", "Event sequence already consumed")
	}
	if s.presentations[w.Handle.Path] != nil {
		return fault("event-type", "Command button requires InvokeCommand")
	}
	if !s.inputAllowed(w) {
		return fault("inactive-widget", w.Handle.Path)
	}
	var target CollectionTarget
	local := false
	if isCollection(w.Handle.Kind) {
		if event.Collection == nil || event.Value != (Value{}) || event.DraftRevision != 0 {
			return fault("event-type", "Collection event requires only a collection target")
		}
		target = *event.Collection
		if target.Handle != event.Handle || target.ModelRevision != event.ModelRevision {
			return fault("event-type", "Collection target differs from envelope")
		}
		_, x, err := s.interactive(target)
		if err != nil {
			return err
		}
		switch event.Kind {
		case Activate:
			if target.ItemID == "" || x.Kind != Row {
				return fault("event-type", "Activation requires a visible row")
			}
		case Select, Expand, Collapse, Retry:
			local = true
		default:
			return fault("event-type", string(event.Kind))
		}
	} else {
		if event.Collection != nil {
			return fault("event-type", "Unexpected collection payload")
		}
		switch event.Kind {
		case Activate:
			if w.Handle.Kind != "button" || event.Value != (Value{}) || event.DraftRevision != 0 {
				return fault("event-type", "Activate requires a button and no payload")
			}
		case Commit:
			if w.Handle.Kind != "input" || event.DraftRevision != w.DraftRevision || !validValue(event.Value, String) || event.Value.Text != w.Draft {
				return fault("event-type", "Commit requires the current input draft")
			}
		default:
			return fault("event-type", string(event.Kind))
		}
	}
	if local {
		n := s.copyState()
		// Gate once, with local semantics and consumed sequence in one candidate.
		n.check = nil
		n.stateCheck = nil
		n.presentationCheck = nil
		n.presentationPrepare = nil
		switch event.Kind {
		case Select:
			err = n.SelectItem(target)
		case Expand:
			err = n.ExpandItem(target)
		case Collapse:
			err = n.CollapseItem(target)
		case Retry:
			err = n.RetryItem(target)
		}
		if err != nil {
			return err
		}
		n.check = s.check
		n.stateCheck = s.stateCheck
		n.presentationCheck = s.presentationCheck
		n.presentationPrepare = s.presentationPrepare
		n.sequence = event.Sequence
		return s.publish(n)
	}
	s.sequence = event.Sequence
	s.StateRevision++
	handler := s.handlers[w.Handle.Path]
	if handler == nil {
		return fault("unbound", w.Handle.Path)
	}
	revision := s.Revision
	// Keep the reentrancy guard independent of caller/handler pointer mutation.
	if event.Collection != nil {
		copy := target
		event.Collection = &copy
	}
	updates, err := handler(event)
	if err != nil {
		return err
	}
	if s.closed || s.Revision != revision {
		return fault("stale-result", "Session changed while callback executed")
	}
	if event.Collection != nil {
		_, x, e := s.interactive(target)
		if e != nil || x.Kind != Row {
			return fault("stale-result", "Collection target changed while callback executed")
		}
	}
	return s.Apply(revision, s.BatchRevision+1, updates)
}
