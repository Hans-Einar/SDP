package runtime

func validateOptions(options []ChoiceOption) error {
	if len(options) > 4096 {
		return fault("choice-limit", "More than 4096 options")
	}
	ids := map[ItemID]bool{}
	for _, o := range options {
		if !validPlain(string(o.ID), 1024, false) || ids[o.ID] {
			return fault("choice-id", "Invalid or duplicate option ID")
		}
		if !validPlain(o.Label, 32768, false) {
			return fault("choice-label", "Invalid single-line label")
		}
		ids[o.ID] = true
	}
	return nil
}

func (s *Session) BindChoices(choices map[string][]ChoiceOption) error {
	if s.closed {
		return fault("closed", "Session is closed")
	}
	count := 0
	n := s.copyState()
	baseline := s.StateRevision
	for path, f := range n.fields {
		w := n.widgets[path]
		if w.Handle.Kind != "select" {
			continue
		}
		count++
		options, ok := choices[w.InstancePath]
		if !ok {
			return fault("choices-unsupplied", w.InstancePath)
		}
		if err := validateOptions(options); err != nil {
			return err
		}
		wasBound := f.Target.OptionGeneration > 0
		f.Options = append([]ChoiceOption(nil), options...)
		f.Target.OptionGeneration++
		if !wasBound && f.Accepted.OptionID != "" && !optionEnabled(f, f.Accepted.OptionID) {
			return fault("choice-initial", "Accepted option missing or disabled")
		}
		if wasBound {
			f.feedback = FieldValidation{}
			if f.Accepted.OptionID != "" && !optionEnabled(f, f.Accepted.OptionID) || f.Proposed.OptionID != "" && !optionEnabled(f, f.Proposed.OptionID) {
				f.Proposed = Choice("")
				w.DraftRevision++
			}
			w.Dirty = f.Proposed != f.Accepted
		}
		if err := n.validateField(path); err != nil {
			return err
		}
	}
	if count != len(choices) {
		return fault("choices-unused", "Choice inventory differs")
	}
	if s.StateRevision != baseline {
		return fault("stale-validation", "State changed in validator")
	}
	return s.publish(n)
}
func (s *Session) ReplaceChoices(h Handle, generation uint64, options []ChoiceOption) error {
	_, err := s.lookup(h)
	if err != nil {
		return err
	}
	f := s.fields[h.Path]
	if f == nil || h.Kind != "select" || generation == 0 || f.Target.OptionGeneration != generation {
		return fault("stale-option", "Option set changed")
	}
	if err = validateOptions(options); err != nil {
		return err
	}
	baseline := s.StateRevision
	n := s.copyState()
	f = n.fields[h.Path]
	f.Options = append([]ChoiceOption(nil), options...)
	f.Target.OptionGeneration++
	f.feedback = FieldValidation{}
	if f.Accepted.OptionID != "" && !optionEnabled(f, f.Accepted.OptionID) || f.Proposed.OptionID != "" && !optionEnabled(f, f.Proposed.OptionID) {
		f.Proposed = Choice("")
		n.widgets[h.Path].DraftRevision++
	}
	n.widgets[h.Path].Dirty = f.Proposed != f.Accepted
	if err = n.validateField(h.Path); err != nil {
		return err
	}
	if s.StateRevision != baseline {
		return fault("stale-validation", "State changed in validator")
	}
	return s.publish(n)
}
func (s *Session) Option(h Handle, id ItemID) (OptionTarget, error) {
	f, ok := s.Field(h)
	if !ok || h.Kind != "select" {
		return OptionTarget{}, fault("option-target", "Expected select")
	}
	t := OptionTarget{h, s.Revision, f.Target.OptionGeneration, id}
	return t, s.ValidateOptionTarget(t)
}
func (s *Session) ValidateOptionTarget(t OptionTarget) error {
	_, err := s.lookup(t.Handle)
	if err != nil {
		return err
	}
	f := s.fields[t.Handle.Path]
	if f == nil || t.Handle.Kind != "select" || t.ModelRevision != s.Revision || t.OptionGeneration == 0 || f.Target.OptionGeneration != t.OptionGeneration {
		return fault("stale-option", "Option identity changed")
	}
	if t.OptionID != "" && !optionEnabled(f, t.OptionID) {
		return fault("option-ineligible", "Option missing or disabled")
	}
	return nil
}
