package runtime

import (
	"reflect"
	"unicode/utf8"
)

func validDraft(v string) bool { return len(v) <= 32768 && utf8.ValidString(v) }

func (s *Session) validateContext(c *ContextTarget) error {
	if c == nil || *c == (ContextTarget{}) {
		return nil
	}
	w, err := s.lookup(c.Widget)
	if err != nil {
		return err
	}
	if c.ModelRevision != s.Revision || !s.inputAllowed(w) {
		return fault("command-context", "Context is inactive or stale")
	}
	if c.Item != nil {
		if c.Item.Handle != c.Widget || c.Item.ModelRevision != c.ModelRevision {
			return fault("command-context", "Item differs from context widget")
		}
		_, x, err := s.interactive(*c.Item)
		if err != nil {
			return err
		}
		if x.Kind != Row || c.Item.ItemID == "" {
			return fault("command-context", "Context requires visible row")
		}
	}
	return nil
}
func (s *Session) OpenMenu(h Handle, context ContextTarget) (MenuScope, error) {
	m := s.menus[h.Path]
	if s.closed || m == nil || m.Handle != h {
		return MenuScope{}, fault("menu", "Unknown root menu")
	}
	w := s.aux[h.Path]
	if !s.inputAllowed(w) {
		return MenuScope{}, fault("menu", "Menu is inactive")
	}
	if target, ok := s.menuTargets[h.Path]; ok {
		if context.Widget != target {
			return MenuScope{}, fault("menu-context", "Context differs from declared target")
		}
	} else if context != (ContextTarget{}) {
		return MenuScope{}, fault("menu-context", "Only context menus accept context")
	}
	if err := s.validateContext(&context); err != nil {
		return MenuScope{}, err
	}
	n := s.copyState()
	for _, other := range n.menus {
		other.Open = false
	}
	m = n.menus[h.Path]
	m.Open = true
	m.Context = nil
	m.Opener = Handle{}
	m.Surface = nil
	if context.Widget != (Handle{}) {
		m.Context = copyContext(&context)
		m.Opener = context.Widget
		if d := s.surfaceFor(s.currentControl(context.Widget.Path)); d != nil {
			m.Surface = copySurface(&d.Target)
		}
	} else if d := s.surfaceFor(w); d != nil {
		m.Surface = copySurface(&d.Target)
	}
	m.CaptureRevision = s.StateRevision + 1
	m.Scope = MenuScope{h, s.Revision, m.CaptureRevision}
	scope := m.Scope
	if err := s.publish(n); err != nil {
		return MenuScope{}, err
	}
	return scope, nil
}

// CloseMenu compares the stored opening, not the now possibly newer global
// revision. Native dismissal must revoke stale captures as well as current ones.
func (s *Session) CloseMenu(scope MenuScope) error {
	m := s.menus[scope.Handle.Path]
	if m == nil || !m.Open || m.Scope != scope {
		return nil
	}
	n := s.copyState()
	n.menus[scope.Handle.Path].Open = false
	if err := s.publish(n); err != nil { // A genuine dismissal cannot be vetoed by geometry.
		if live := s.menus[scope.Handle.Path]; live != nil && live.Open && live.Scope == scope {
			live.Open = false
			s.StateRevision++
		}
		return err
	}
	return nil
}
func (s *Session) CaptureCommand(origin Handle, via string, context *ContextTarget) (Event, error) {
	e := Event{Kind: InvokeCommand, ModelRevision: s.Revision, StateRevision: s.StateRevision, Sequence: s.sequence + 1}
	w, err := s.lookupControl(origin)
	if err != nil {
		return e, err
	}
	if !s.inputAllowed(w) {
		return e, fault("command-origin", "Origin is inactive")
	}
	var c *CommandState
	var surface *SurfaceTarget
	switch via {
	case "button", "menu":
		p := s.presentations[origin.Path]
		if p == nil || p.Handle != origin || !p.Enabled || !p.Visible {
			return e, fault("command-origin", "Unknown or inactive presentation")
		}
		if via == "button" && origin.Kind != "button" || via == "menu" && origin.Kind != "item" {
			return e, fault("command-origin", "Origin kind differs")
		}
		c = s.commands[p.Command.Path]
	case "key":
		owner := s.surfaceFor(w)
		if s.activeSurface == nil {
			if owner != nil {
				return e, fault("inactive-key", "Command belongs to inactive canvas")
			}
		} else if owner == nil || owner.Target != *s.activeSurface {
			return e, fault("inactive-key", "Command belongs to inactive canvas")
		}
		c = s.commands[origin.Path]
		if c == nil || c.Handle != origin || c.Key == "" {
			return e, fault("command-origin", "Unknown key command")
		}
	default:
		return e, fault("command-origin", "Unknown invocation route")
	}
	if c == nil || !c.Enabled || !c.Visible {
		return e, fault("command", "Command is inactive")
	}
	if d := s.surfaceFor(w); d != nil {
		if !d.Open {
			return e, fault("surface", "Origin surface closed")
		}
		surface = copySurface(&d.Target)
	}
	var captured *ContextTarget
	if via == "menu" {
		m := s.menus[s.menuOwner[origin.Path]]
		if m == nil || !m.Open || m.CaptureRevision != s.StateRevision || m.Scope.ModelRevision != s.Revision {
			return e, fault("stale-menu", "Menu capture is stale")
		}
		captured = copyContext(m.Context)
		surface = copySurface(m.Surface)
	}
	if c.Context != "none" {
		if via != "menu" {
			v := ContextTarget{Widget: c.Target, ModelRevision: s.Revision}
			if c.Context == "item" {
				coll := s.collections[c.Target.Path]
				if coll == nil || coll.Selected == "" {
					return e, fault("command-context", "No selected item")
				}
				item, err := s.Target(c.Target, coll.Selected)
				if err != nil {
					return e, err
				}
				v.Item = &item
			}
			captured = &v
		}
		if captured == nil || captured.Widget != c.Target || (c.Context == "item") != (captured.Item != nil) {
			return e, fault("command-context", "Declared context differs from capture")
		}
		if err := s.validateContext(captured); err != nil {
			return e, err
		}
	} else {
		captured = nil
	}
	if context != nil && !reflect.DeepEqual(context, captured) {
		return e, fault("command-context", "Caller context differs from runtime capture")
	}
	e.Handle = c.Handle
	e.Command = &CommandInvocation{Origin: origin, Via: via, Surface: surface, Context: captured}
	if c.Toggle {
		checked := !c.Checked
		if c.Exclusive != "" {
			checked = true
		}
		e.Command.Checked = &checked
	}
	return e, nil
}
func (s *Session) CaptureDialog(t SurfaceTarget, kind EventKind) (Event, error) {
	e := Event{Handle: t.Handle, Kind: kind, ModelRevision: s.Revision, StateRevision: s.StateRevision, Sequence: s.sequence + 1, Dialog: &DialogRequest{Surface: t}}
	d, err := s.surfaceTarget(t)
	if err != nil {
		return e, err
	}
	if kind != Accept && kind != Cancel && kind != Close {
		return e, fault("event-type", "Expected dialog interaction")
	}
	for _, other := range s.surfaces {
		if other.Open && other.Modal && !s.surfaceDescends(d, other.Target) {
			return e, fault("modal", "Dialog input is blocked")
		}
	}
	if kind == Accept {
		if d.AcceptBlocked {
			return e, fault("accept-blocked", "Previous attempt cannot be replayed")
		}
		fields, _ := s.DialogControls(t.Handle)
		if len(fields) > 256 {
			return e, fault("update-limit", "Captured fields exceed 256 writes")
		}
		for _, f := range fields {
			if f.Validation.Code != "" {
				return e, fault("field-validation", f.Validation.Message)
			}
			if f.RawDraft != nil && !validDraft(*f.RawDraft) {
				return e, fault("dialog-field", "Draft exceeds valid UTF-8 bound")
			}
			d := DraftField{Handle: f.Target.Handle, Value: f.Proposed, ValueRevision: f.Target.ValueRevision, DraftRevision: f.Target.DraftRevision}
			if scalar(d.Handle.Kind) {
				d.RawDraft = copyString(f.RawDraft)
				d.FieldValidation = f.Validation
				if d.Handle.Kind == "select" {
					option, err := s.Option(d.Handle, f.Proposed.OptionID)
					if err != nil {
						return e, err
					}
					d.OptionTarget = &option
				}
			}
			e.Dialog.Fields = append(e.Dialog.Fields, d)
		}
	}
	return e, nil
}
