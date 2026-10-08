package runtime

import "github.com/Hans-Einar/SDP/SDUI/go/parser"

func (s *Session) SuccessorWithChoices(root *parser.Instance, providers map[string]CollectionProvider, choices map[string][]ChoiceOption) (*Session, error) {
	return s.successor(root, providers, choices, true)
}
func (n *Session) inheritFields(s *Session, choices map[string][]ChoiceOption, explicit bool) error {
	count := 0
	for path, f := range n.fields {
		w := n.widgets[path]
		old := s.fields[path]
		compatible := old != nil && s.widgets[path] != nil && s.widgets[path].Handle == w.Handle
		if w.Handle.Kind == "select" {
			count++
			options, supplied := choices[w.InstancePath]
			if !explicit && compatible && old.Target.OptionGeneration > 0 {
				options = old.Options
				supplied = true
			}
			if explicit && !supplied {
				return fault("choices-unsupplied", w.InstancePath)
			}
			if supplied {
				if err := validateOptions(options); err != nil {
					return err
				}
				f.Options = append([]ChoiceOption(nil), options...)
				f.Target.OptionGeneration = 1
				if old != nil {
					f.Target.OptionGeneration = old.Target.OptionGeneration + 1
				}
				if f.Accepted.OptionID != "" && !optionEnabled(f, f.Accepted.OptionID) {
					return fault("choice-initial", w.InstancePath)
				}
			}
		}
		if compatible {
			f.Accepted = old.Accepted
			f.Proposed = old.Proposed
			f.RawDraft = copyString(old.RawDraft)
			f.acceptedRaw = old.acceptedRaw
			if err := n.validAcceptedField(f, w.Handle.Kind, old.Required); err != nil {
				return err
			}
			if n.ownerSurface[w.InstancePath] != "" {
				n.resetField(path)
			}
		}
		if err := n.validateField(path); err != nil {
			return err
		}
	}
	if explicit && count != len(choices) {
		return fault("choices-unused", "Choice inventory differs")
	}
	return nil
}
func (s *Session) validAcceptedField(f *field, kind string, wasRequired bool) error {
	switch kind {
	case "checkbox":
		if !validValue(f.Accepted, Boolean) {
			return fault("field-type", "Invalid accepted checkbox")
		}
	case "number", "slider":
		if !validValue(f.Accepted, Number) {
			return fault("field-type", "Invalid accepted number")
		}
		if _, err := f.grid.Tick(f.Accepted.Number); err != nil {
			return err
		}
		if kind == "number" {
			v, err := f.grid.Parse(f.acceptedRaw)
			if err != nil {
				return err
			}
			if v != f.Accepted.Number {
				return fault("number-grid", "Accepted provenance differs")
			}
		}
	case "select":
		if f.Target.OptionGeneration == 0 {
			return fault("choices-unsupplied", f.InstancePath)
		}
		// An already-required empty baseline remains editable invalid absence.
		// Only newly requiring an empty value tightens this constraint; nonempty
		// accepted IDs must still exist and be enabled in the replacement set.
		if f.Accepted.OptionID != "" && !optionEnabled(f, f.Accepted.OptionID) || f.Required && !wasRequired && f.Accepted.OptionID == "" {
			return fault("choice-accepted", "Retained accepted choice violates new constraints")
		}
	}
	return nil
}
