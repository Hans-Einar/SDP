package runtime

func (s *Session) applyFieldUpdate(u Update) error {
	w := s.widgets[u.Handle.Path]
	f := s.fields[u.Handle.Path]

	switch u.Property {
	case ReadOnly:
		if !validValue(u.Value, Boolean) || u.Validation != nil || u.AcceptDraft {
			return fault("property-type", "ReadOnly requires Boolean")
		}
		f.ReadOnly = u.Value.Bool
		return nil
	case ValidationState:
		if u.Value != (Value{}) || u.Validation == nil || !validValidation(*u.Validation) || u.AcceptDraft {
			return fault("property-type", "Expected bounded validation metadata")
		}
		f.feedback = *u.Validation
		f.Validation = s.fieldValidation(f, w.Handle.Kind)
		return nil
	case AcceptedValue:
		if u.Validation != nil {
			return fault("property-type", "Unexpected validation payload")
		}
		if w.Dirty && !u.AcceptDraft {
			return fault("draft-conflict", w.Handle.Path)
		}
		if u.AcceptDraft && (f.Validation.Code != "" || u.Value != f.Proposed) {
			return fault("draft-conflict", "Value differs from valid captured proposal")
		}
		switch w.Handle.Kind {
		case "checkbox":
			if !validValue(u.Value, Boolean) {
				return fault("property-type", "Checkbox requires Boolean")
			}
		case "slider", "number":
			if !validValue(u.Value, Number) {
				return fault("property-type", "Numeric field requires Number")
			}
			k, err := f.grid.Tick(u.Value.Number)
			if err != nil {
				return err
			}
			if w.Handle.Kind == "number" {
				if u.AcceptDraft {
					f.acceptedRaw = *f.RawDraft
				} else {
					f.acceptedRaw = rawPoint(f.grid, k)
					f.RawDraft = copyString(&f.acceptedRaw)
				}
			}
		case "select":
			if !validValue(u.Value, OptionID) || u.Value.OptionID == "" && f.Required || u.Value.OptionID != "" && !optionEnabled(f, u.Value.OptionID) {
				return fault("choice-ineligible", "Invalid accepted option")
			}
		}
		f.Accepted = u.Value
		f.Proposed = u.Value
		f.feedback = FieldValidation{}
		w.Dirty = false
		w.ValueRevision++
		w.DraftRevision++
		return nil
	}
	return fault("property", "Unsupported typed property")
}

func (s *Session) checkFieldUpdate(u Update) error {
	w := s.widgets[u.Handle.Path]
	f := s.fields[u.Handle.Path]
	if u.ExpectedValueRevision != w.ValueRevision || u.ExpectedDraftRevision != w.DraftRevision {
		return fault("field-conflict", "Field revisions changed")
	}
	if w.Handle.Kind == "select" && (u.ExpectedOptionGeneration == 0 || u.ExpectedOptionGeneration != f.Target.OptionGeneration) {
		return fault("stale-option", "Receiver options changed")
	}
	return nil
}
