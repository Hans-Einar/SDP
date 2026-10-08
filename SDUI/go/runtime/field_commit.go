package runtime

import "reflect"

func (s *Session) CaptureCommit(target FieldTarget) (Event, error) {
	e := Event{}
	f, ok := s.Field(target.Handle)
	if !ok || s.fields[target.Handle.Path] == nil {
		return e, fault("field", "Expected typed field")
	}
	if f.Target != target {
		return e, fault("stale-field", "Field capture changed")
	}
	w := s.widgets[target.Handle.Path]
	if !s.inputAllowed(w) || f.ReadOnly {
		return e, fault("field-read-only", "Field cannot commit")
	}
	if f.Validation.Code != "" {
		return e, fault("field-validation", f.Validation.Message)
	}
	c := &ControlCommit{ValueRevision: target.ValueRevision}
	if target.Handle.Kind == "number" {
		c.RawDraft = copyString(f.RawDraft)
	}
	if target.Handle.Kind == "select" {
		t, err := s.Option(target.Handle, f.Proposed.OptionID)
		if err != nil {
			return e, err
		}
		c.Option = &t
	}
	return Event{Handle: target.Handle, ModelRevision: s.Revision, StateRevision: s.StateRevision, Sequence: s.sequence + 1, DraftRevision: target.DraftRevision, Kind: Commit, Value: f.Proposed, Control: c}, nil
}
func copyControlEvent(e Event) Event {
	if e.Control != nil {
		x := *e.Control
		x.RawDraft = copyString(x.RawDraft)
		if x.Option != nil {
			o := *x.Option
			x.Option = &o
		}
		e.Control = &x
	}
	return e
}
func (s *Session) dispatchField(e Event) error {
	if s.interacting {
		return fault("reentrant-interaction", "Interaction already active")
	}
	if e.Control == nil || e.Kind != Commit || e.Page != nil || e.Split != nil || e.Command != nil || e.Dialog != nil || e.Collection != nil {
		return fault("event-type", "Expected only typed Commit payload")
	}
	f, ok := s.Field(e.Handle)
	if !ok {
		return fault("field", "Unknown typed field")
	}
	target := f.Target
	target.ModelRevision = e.ModelRevision
	target.StateRevision = e.StateRevision
	target.ValueRevision = e.Control.ValueRevision
	target.DraftRevision = e.DraftRevision
	if e.Control.Option != nil {
		target.OptionGeneration = e.Control.Option.OptionGeneration
	}
	expected, err := s.CaptureCommit(target)
	if err != nil {
		return err
	}
	if e.Sequence == 0 || e.Sequence <= s.sequence {
		return fault("duplicate-event", "Sequence consumed")
	}
	expected.Sequence = e.Sequence
	if !reflect.DeepEqual(e, expected) {
		return fault("event-type", "Commit differs from exact capture")
	}
	handler := s.handlers[e.Handle.Path]
	w := s.widgets[e.Handle.Path]
	if handler == nil && w.Binding.Module != "" {
		return fault("unbound", e.Handle.Path)
	}
	if handler == nil {
		n := s.copyState()
		n.sequence = e.Sequence
		if d := s.surfaceFor(w); d != nil && d.Open {
			return s.publish(n)
		}
		update := Update{Handle: e.Handle, Property: AcceptedValue, Value: e.Value, ExpectedValueRevision: f.Target.ValueRevision, ExpectedDraftRevision: f.Target.DraftRevision, ExpectedOptionGeneration: f.Target.OptionGeneration, AcceptDraft: true}
		n, err = n.applyCandidate(s.Revision, s.BatchRevision+1, []Update{update})
		if err != nil {
			return err
		}
		if s.StateRevision != e.StateRevision || s.Revision != e.ModelRevision {
			return fault("stale-validation", "State changed during local acceptance")
		}
		return s.publish(n)
	}
	// Probe before execution never prepares resources. The proposal is already a
	// published draft; the accepting update is validated again on the final reply.
	probe := s.copyState()
	if err = probe.gate(); err != nil {
		return err
	}
	if s.StateRevision != e.StateRevision || s.Revision != e.ModelRevision {
		return fault("stale-event", "State changed during probe")
	}
	captured := s.captureControls()
	fieldTargets := s.captureFieldTargets()
	s.sequence = e.Sequence
	s.StateRevision++
	baseline := s.StateRevision
	s.interacting = true
	defer func() { s.interacting = false }()
	updates, err := handler(copyControlEvent(e))
	if err != nil {
		return err
	}
	if s.closed || s.Revision != e.ModelRevision || s.StateRevision != baseline {
		return fault("stale-result", "State changed during callback")
	}
	accepted := false
	for _, u := range updates {
		old, ok := captured[u.Handle]
		live := s.currentControl(u.Handle.Path)
		if !ok || live == nil || live.Handle != old.Handle || live.ValueRevision != old.ValueRevision || live.DraftRevision != old.DraftRevision {
			return fault("stale-result", "Receiver changed")
		}
		if before, ok := fieldTargets[u.Handle]; ok {
			now, _ := s.Field(u.Handle)
			if now.Target.OptionGeneration != before.OptionGeneration {
				return fault("stale-option", "Receiver options changed")
			}
		}
		if u.Handle == e.Handle && u.Property == AcceptedValue && u.AcceptDraft && u.Value == e.Value && u.ExpectedDraftRevision == e.DraftRevision {
			accepted = true
		}
	}
	if !accepted {
		return fault("commit-not-accepted", "Handler must explicitly accept captured source proposal")
	}
	n, err := s.applyCandidate(s.Revision, s.BatchRevision+1, updates)
	if err != nil {
		return err
	}
	if s.StateRevision != baseline {
		return fault("stale-result", "State changed during validation")
	}
	return s.publish(n)
}
func (s *Session) captureFieldTargets() map[Handle]FieldTarget {
	out := map[Handle]FieldTarget{}
	for _, w := range s.widgets {
		if f, ok := s.Field(w.Handle); ok {
			out[w.Handle] = f.Target
		}
	}
	return out
}
