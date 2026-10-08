package runtime

func (s *Session) editField(h Handle, revision uint64, value Value, tick *uint64, option *OptionTarget) (FieldChange, error) {
	w, err := s.lookup(h)
	if err != nil {
		return FieldChange{}, err
	}
	f := s.fields[h.Path]
	if f == nil {
		return FieldChange{}, fault("field", "Not a typed field")
	}
	if revision != s.Revision {
		return FieldChange{}, fault("stale-event", "Model changed")
	}
	if !s.inputAllowed(w) || f.ReadOnly {
		return FieldChange{}, fault("field-read-only", "Field is not editable")
	}
	baseline := s.StateRevision
	n := s.copyState()
	f = n.fields[h.Path]
	f.feedback = FieldValidation{}
	switch h.Kind {
	case "number":
		if tick != nil {
			if f.Validation.Code != "" {
				return FieldChange{}, fault("field-validation", "Invalid numeric draft cannot step")
			}
			if _, err = f.grid.At(*tick); err != nil {
				return FieldChange{}, err
			}
			raw := rawPoint(f.grid, *tick)
			f.RawDraft = &raw
		} else {
			if !validValue(value, String) || !validDraft(value.Text) {
				return FieldChange{}, fault("field-type", "Number edit requires bounded raw String")
			}
			f.RawDraft = copyString(&value.Text)
		}
	case "slider":
		if tick != nil {
			v, e := f.grid.At(*tick)
			if e != nil {
				return FieldChange{}, e
			}
			value = Numeric(v)
		}
		if !validValue(value, Number) {
			return FieldChange{}, fault("field-type", "Slider requires Number")
		}
		if _, e := f.grid.Tick(value.Number); e != nil {
			return FieldChange{}, e
		}
		f.Proposed = value
	case "checkbox":
		if tick != nil || !validValue(value, Boolean) {
			return FieldChange{}, fault("field-type", "Checkbox requires Boolean")
		}
		f.Proposed = value
	case "select":
		if option == nil {
			return FieldChange{}, fault("option-target", "Selection requires option token")
		}
		if err = s.ValidateOptionTarget(*option); err != nil {
			return FieldChange{}, err
		}
		f.Proposed = Choice(option.OptionID)
	}
	nw := n.widgets[h.Path]
	nw.DraftRevision++
	if err = n.validateField(h.Path); err != nil {
		return FieldChange{}, err
	}
	nw.Dirty = f.Proposed != f.Accepted || f.RawDraft != nil && *f.RawDraft != f.acceptedRaw
	if s.StateRevision != baseline {
		return FieldChange{}, fault("stale-validation", "State changed in validator")
	}
	if err = s.publish(n); err != nil {
		return FieldChange{}, err
	}
	out, _ := s.Field(h)
	change := FieldChange{Field: out, Value: out.Proposed}
	if h.Kind == "number" {
		change.Value = Text(*out.RawDraft)
	}
	if h.Kind == "select" {
		v := OptionTarget{h, s.Revision, out.Target.OptionGeneration, out.Proposed.OptionID}
		change.Option = &v
	}
	if observer := s.changes[h.Path]; observer != nil {
		copy := change
		copy.Field = copyFieldState(out)
		if copy.Option != nil {
			o := *copy.Option
			copy.Option = &o
		}
		observer(copy)
	}
	return change, nil
}
func (s *Session) EditField(h Handle, revision uint64, value Value) (FieldChange, error) {
	return s.editField(h, revision, value, nil, nil)
}
func (s *Session) EditTick(h Handle, revision, tick uint64) (FieldChange, error) {
	if h.Kind != "slider" && h.Kind != "number" {
		return FieldChange{}, fault("field-type", "Ticks require numeric field")
	}
	return s.editField(h, revision, Value{}, &tick, nil)
}
func (s *Session) ChooseOption(t OptionTarget) (FieldChange, error) {
	return s.editField(t.Handle, t.ModelRevision, Value{}, nil, &t)
}
