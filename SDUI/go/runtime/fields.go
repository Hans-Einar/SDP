package runtime

import (
	"github.com/Hans-Einar/SDP/SDUI/go/numeric"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"unicode/utf8"
)

const ReadOnly Property = "read-only"

type field struct {
	FieldState
	grid        *numeric.Grid
	acceptedRaw string
	feedback    FieldValidation
}

func scalar(kind string) bool {
	return kind == "checkbox" || kind == "slider" || kind == "number" || kind == "select"
}
func copyString(p *string) *string {
	if p == nil {
		return nil
	}
	v := *p
	return &v
}
func copyFieldState(f FieldState) FieldState {
	f.RawDraft = copyString(f.RawDraft)
	if f.Numeric != nil {
		x := *f.Numeric
		f.Numeric = &x
	}
	f.Options = append([]ChoiceOption(nil), f.Options...)
	return f
}
func (s *Session) initFields() error {
	s.fields = map[string]*field{}
	s.validators = map[string]FieldValidator{}
	s.changes = map[string]ChangeHandler{}
	var first error
	s.root.Walk(func(n *parser.Instance) {
		if first != nil || n.Kind != "widget" || !scalar(n.Widget) {
			return
		}
		if n.Profile != "sdui/0.3" {
			first = fault("field-profile", "Typed controls require sdui/0.3")
			return
		}
		w := s.widgets[public(n.Path)]
		f := &field{FieldState: FieldState{InstancePath: n.Path, ReadOnly: boolArgument(n, "readOnly", false), Required: boolArgument(n, "required", false)}}
		switch n.Widget {
		case "checkbox":
			f.Accepted = Bool(boolArgument(n, "value", false))
		case "select":
			f.Accepted = Choice(ItemID(n.Argument("value")))
		case "number", "slider":
			a, b, d, v, err := parser.NumericArguments(n)
			if err != nil {
				first = err
				return
			}
			f.Numeric = &NumericConstraints{a, b, d}
			f.grid, err = numeric.NewGrid(a, b, d)
			if err != nil {
				first = err
				return
			}
			number, err := f.grid.Parse(v)
			if err != nil {
				first = err
				return
			}
			f.Accepted = Numeric(number)
			f.acceptedRaw = v
			if n.Widget == "number" {
				f.RawDraft = copyString(&v)
			}
		}
		f.Proposed = f.Accepted
		s.fields[w.Handle.Path] = f
		w.Value = ""
		w.Draft = ""
		f.Validation = s.fieldValidation(f, w.Handle.Kind)
	})
	return first
}
func (s *Session) Field(h Handle) (FieldState, bool) {
	w, err := s.lookup(h)
	if err != nil {
		return FieldState{}, false
	}
	f := s.fields[h.Path]
	var out FieldState
	if f != nil {
		out = copyFieldState(f.FieldState)
	} else if h.Kind == "input" {
		out = FieldState{InstancePath: w.InstancePath, Accepted: Text(w.Value), Proposed: Text(w.Draft), RawDraft: copyString(&w.Draft)}
	} else {
		return FieldState{}, false
	}
	out.Target = FieldTarget{h, s.Revision, s.StateRevision, w.ValueRevision, w.DraftRevision, out.Target.OptionGeneration}
	out.Dirty = w.Dirty
	return out, true
}
func (s *Session) copyFields(n *Session) {
	n.fields = map[string]*field{}
	for k, f := range s.fields {
		x := *f
		x.FieldState = copyFieldState(f.FieldState)
		n.fields[k] = &x
	}
	n.validators = map[string]FieldValidator{}
	for k, v := range s.validators {
		n.validators[k] = v
	}
	n.changes = map[string]ChangeHandler{}
	for k, v := range s.changes {
		n.changes[k] = v
	}
}
func (s *Session) snapshotFields(v *Snapshot) {
	v.Fields = map[string]FieldState{}
	for _, w := range s.widgets {
		if f, ok := s.Field(w.Handle); ok {
			v.Fields[w.InstancePath] = f
		}
	}
}
func validValidation(v FieldValidation) bool {
	return len(v.Code) <= 4096 && len(v.Message) <= 4096 && utf8.ValidString(v.Code) && utf8.ValidString(v.Message) && (v.Code != "" || v.Message == "")
}
func validation(code string) FieldValidation { return FieldValidation{code, code} }
func (s *Session) fieldValidation(f *field, kind string) FieldValidation {
	switch kind {
	case "checkbox":
		if !validValue(f.Proposed, Boolean) {
			return validation("field-type")
		}
	case "number":
		if f.RawDraft == nil || !validDraft(*f.RawDraft) {
			return validation("number-draft")
		}
		v, err := f.grid.Parse(*f.RawDraft)
		if err != nil {
			f.Proposed = Value{}
			return validation(err.Error())
		}
		f.Proposed = Numeric(v)
	case "slider":
		if !validValue(f.Proposed, Number) {
			return validation("field-type")
		}
		if _, err := f.grid.Tick(f.Proposed.Number); err != nil {
			return validation(err.Error())
		}
	case "select":
		if f.Target.OptionGeneration == 0 {
			return validation("choices-unsupplied")
		}
		if !validValue(f.Proposed, OptionID) {
			return validation("field-type")
		}
		if f.Proposed.OptionID == "" {
			if f.Required {
				return validation("choice-required")
			}
			if f.Accepted.OptionID != "" && !optionEnabled(f, f.Accepted.OptionID) {
				return validation("choice-removed")
			}
		} else if !optionEnabled(f, f.Proposed.OptionID) {
			return validation("choice-ineligible")
		}
	}
	if f.feedback.Code != "" {
		return f.feedback
	}
	return FieldValidation{}
}
func optionEnabled(f *field, id ItemID) bool {
	for _, o := range f.Options {
		if o.ID == id {
			return o.Enabled
		}
	}
	return false
}
func (s *Session) validateField(path string) error {
	f := s.fields[path]
	if f == nil {
		return nil
	}
	w := s.widgets[path]
	f.Validation = s.fieldValidation(f, w.Handle.Kind)
	if f.Validation.Code == "" && s.validators[path] != nil {
		v, _ := s.Field(w.Handle)
		result := s.validators[path](v)
		if !validValidation(result) {
			return fault("field-validation", "Validator returned malformed feedback")
		}
		f.Validation = result
	}
	return nil
}
func (s *Session) ObserveChanges(h Handle, fn ChangeHandler) error {
	if _, ok := s.Field(h); !ok || s.fields[h.Path] == nil {
		return fault("field", "Expected typed field")
	}
	s.changes[h.Path] = fn
	s.StateRevision++
	return nil
}
func (s *Session) ValidateFieldWith(h Handle, fn FieldValidator) error {
	if _, ok := s.Field(h); !ok || s.fields[h.Path] == nil {
		return fault("field", "Expected typed field")
	}
	state := s.StateRevision
	n := s.copyState()
	n.validators[h.Path] = fn
	if err := n.validateField(h.Path); err != nil {
		return err
	}
	if s.StateRevision != state {
		return fault("stale-validation", "State changed in validator")
	}
	return s.publish(n)
}
func (s *Session) DialogControls(h Handle) ([]FieldState, error) {
	if _, ok := s.SurfaceState(h); !ok {
		return nil, fault("surface", "Unknown dialog")
	}
	out := []FieldState{}
	s.root.Walk(func(node *parser.Instance) {
		w := s.widgets[public(node.Path)]
		if w != nil && s.ownerSurface[node.Path] == h.Path {
			if f, ok := s.Field(w.Handle); ok {
				out = append(out, f)
			}
		}
	})
	return out, nil
}
func (s *Session) resetField(path string) {
	w := s.widgets[path]
	f := s.fields[path]
	if f == nil {
		return
	}
	if w.Dirty {
		w.DraftRevision++
	}
	f.Proposed = f.Accepted
	if w.Handle.Kind == "number" {
		f.RawDraft = copyString(&f.acceptedRaw)
	}
	f.feedback = FieldValidation{}
	w.Dirty = false
	f.Validation = s.fieldValidation(f, w.Handle.Kind)
}
func (s *Session) RevertField(h Handle) error {
	if h.Kind == "input" {
		return s.Revert(h)
	}
	if _, ok := s.Field(h); !ok {
		return fault("field", "Unknown field")
	}
	n := s.copyState()
	n.resetField(h.Path)
	return s.publish(n)
}

// Formatting a point is display only. Decimal admission always uses its original
// lexeme; exact reconstructed ticks get their own finite decimal representation.
func rawPoint(g *numeric.Grid, k uint64) string { v, _ := g.Text(k); return v }
