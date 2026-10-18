package runtime

import "github.com/Hans-Einar/SDP/SDUI/go/parser"

// OpenSurface uses current logical focus; input adapters should supply their
// actual origin through OpenSurfaceFrom instead.
func (s *Session) OpenSurface(h Handle, context ContextTarget) (SurfaceTarget, error) {
	opener := s.rootOwner
	if w := s.currentControl(s.focused); w != nil {
		opener = w.Handle
	}
	return s.OpenSurfaceFrom(h, opener, context)
}
func (s *Session) OpenSurfaceFrom(h, opener Handle, context ContextTarget) (SurfaceTarget, error) {
	n := s.copyState()
	t, err := n.openSurface(h, opener, context)
	if err != nil {
		return SurfaceTarget{}, err
	}
	if err = s.publish(n); err != nil {
		return SurfaceTarget{}, err
	}
	return t, nil
}
func (s *Session) openSurface(h, opener Handle, context ContextTarget) (SurfaceTarget, error) {
	d, ok := s.SurfaceState(h)
	if !ok {
		return SurfaceTarget{}, fault("surface", "Unknown dialog")
	}
	for _, r := range s.dialogResults {
		if r.Surface.Handle == h {
			return SurfaceTarget{}, fault("surface-result", "Drain previous terminal result before reopening")
		}
	}
	if d.Open {
		for _, other := range s.surfaces {
			if other.Open && other.Modal && !s.surfaceDescends(s.surfaces[h.Path], other.Target) {
				return SurfaceTarget{}, fault("modal", "Dialog is blocked by modal child")
			}
		}
		s.activeSurface = copySurface(&d.Target)
		s.focused = s.surfaceFocus(h.Path, d.Focused)
		s.surfaces[h.Path].Focused = s.focused
		return d.Target, nil
	}
	var parent *SurfaceState
	if opener != s.rootOwner {
		w, err := s.lookupControl(opener)
		if err != nil {
			return SurfaceTarget{}, err
		}
		if !s.inputAllowed(w) {
			return SurfaceTarget{}, fault("surface-opener", "Opener is inactive")
		}
		parent = s.surfaceFor(w)
	}
	if err := s.validateContext(&context); err != nil {
		return SurfaceTarget{}, err
	}
	for p := parent; p != nil; {
		if p.Handle == h {
			return SurfaceTarget{}, fault("surface-parent", "Cyclic parent lifetime")
		}
		if p.ParentSurface == nil {
			break
		}
		p = s.surfaces[p.ParentSurface.Handle.Path]
	}
	live := s.surfaces[h.Path]
	if !d.Open {
		a := s.active[d.InstancePath]
		if !a.enabled || !a.visible {
			return SurfaceTarget{}, fault("surface", "Dialog declaration is inactive")
		}
		s.surfaceEpoch++
		live.Open = true
		live.Target = SurfaceTarget{h, s.Revision, s.surfaceEpoch}
		live.Parent = s.rootOwner
		live.ParentSurface = nil
		live.Opener = opener
		live.Context = nil
		if context.Widget != (Handle{}) {
			live.Context = copyContext(&context)
		}
		if parent != nil {
			live.Parent = parent.Handle
			live.ParentSurface = copySurface(&parent.Target)
		}
		live.AcceptSequence = 0
		live.Domain = DomainNotCalled
		live.AcceptBlocked = false
		live.Message = ""
		live.Focused = ""
	}
	s.activeSurface = copySurface(&live.Target)
	s.refreshActivity()
	if w := s.currentControl(live.Focused); w == nil || !s.inputAllowed(w) {
		live.Focused = s.firstSurfaceFocus(h.Path)
	}
	s.focused = live.Focused
	return live.Target, nil
}
func (s *Session) firstSurfaceFocus(path string) string {
	first := ""
	s.root.Walk(func(n *parser.Instance) {
		if first != "" || s.ownerSurface[n.Path] != path {
			return
		}
		w := s.currentControl(public(n.Path))
		if w != nil && s.inputAllowed(w) && (w.Handle.Kind == "input" || w.Handle.Kind == "button" || w.Handle.Kind == "tree" || w.Handle.Kind == "list" || w.Handle.Kind == "tabs" || w.Handle.Kind == "split" || scalar(w.Handle.Kind)) {
			first = w.Handle.Path
		}
	})
	return first
}
func (s *Session) FocusSurface(target *SurfaceTarget) error {
	if s.closed {
		return fault("closed", "Session is closed")
	}
	n := s.copyState()
	if target == nil {
		for _, d := range s.surfaces {
			if d.Open && d.Modal {
				return fault("modal", "Main canvas is blocked")
			}
		}
		n.activeSurface = nil
		n.focused = n.surfaceFocus("", n.mainFocus)
	} else {
		d, err := n.surfaceTarget(*target)
		if err != nil {
			return err
		}
		for _, other := range n.surfaces {
			if other.Open && other.Modal && !n.surfaceDescends(d, other.Target) {
				return fault("modal", "Canvas is blocked")
			}
		}
		n.activeSurface = copySurface(target)
		n.focused = d.Focused
		if w := n.currentControl(n.focused); w == nil || !n.inputAllowed(w) {
			n.focused = n.firstSurfaceFocus(d.Handle.Path)
			d.Focused = n.focused
		}
	}
	return s.publish(n)
}
func (s *Session) surfaceDescends(d *SurfaceState, t SurfaceTarget) bool {
	for d != nil {
		if d.Target == t {
			return true
		}
		if d.ParentSurface == nil {
			return false
		}
		d = s.surfaces[d.ParentSurface.Handle.Path]
	}
	return false
}
func (s *Session) ConfirmSurfacePublication(t SurfaceTarget) error {
	if _, err := s.surfaceTarget(t); err != nil {
		return err
	}
	s.publishedSurfaces[t] = true
	return nil
}
func (s *Session) DrainDialogResults() []DialogResult {
	if s.accepting != nil {
		return nil
	}
	out := copyResults(s.dialogResults)
	s.dialogResults = nil
	return out
}
func (s *Session) CloseSurface(t SurfaceTarget, kind string) error {
	if kind != "cancel" && kind != "close" {
		return fault("surface-kind", "Expected cancel or close")
	}
	if _, err := s.surfaceTarget(t); err != nil {
		return err
	}
	n := s.copyState()
	n.closeSurface(t, kind, "programmatic", 0, nil)
	return s.publish(n)
}

// closeSurface runs only on a candidate, or on explicitly forced native teardown.
func (s *Session) closeSurface(t SurfaceTarget, kind, reason string, seq uint64, fields []DraftField) {
	d := s.surfaces[t.Handle.Path]
	if d == nil || !d.Open || d.Target != t {
		return
	}
	for _, child := range s.orderedSurfaces() {
		if child.Open && child.ParentSurface != nil && *child.ParentSurface == t {
			childReason := "parent-closed"
			if reason == "reload" || reason == "dispose" || reason == "parent-hidden" {
				childReason = reason
			}
			s.closeSurface(child.Target, "close", childReason, 0, nil)
		}
	}
	for path, w := range s.widgets {
		if s.ownerSurface[w.InstancePath] == d.Handle.Path && s.fields[path] != nil {
			s.resetField(path)
		}
	}
	owned, _ := s.DialogFields(d.Handle)
	for _, field := range owned {
		w := s.widgets[field.Handle.Path]
		if w.Dirty {
			w.Draft = w.Value
			w.Dirty = false
			w.DraftRevision++
		}
	}
	for path, c := range s.collections {
		w := s.widgets[path]
		if w != nil && s.ownerSurface[w.InstancePath] == d.Handle.Path {
			c.cancel()
		}
	}
	for _, m := range s.menus {
		if m.Surface != nil && *m.Surface == t {
			m.Open = false
		}
	}
	if s.publishedSurfaces[t] {
		s.dialogResults = append(s.dialogResults, DialogResult{Surface: t, Sequence: seq, AcceptSequence: d.AcceptSequence, Kind: kind, Reason: reason, Domain: d.Domain, Fields: copyDraftFields(fields)})
		delete(s.publishedSurfaces, t)
	}
	d.Open = false
	if s.activeSurface != nil && *s.activeSurface == t {
		s.activeSurface = copySurface(d.ParentSurface)
		s.focused = ""
		parentPath := ""
		remembered := s.mainFocus
		if d.ParentSurface != nil {
			parentPath = d.ParentSurface.Handle.Path
			remembered = s.surfaces[parentPath].Focused
		}
		if w := s.currentControl(d.Opener.Path); w != nil && w.Handle == d.Opener && s.inputAllowed(w) {
			s.focused = w.Handle.Path
		} else {
			s.focused = s.surfaceFocus(parentPath, remembered)
		}
		s.rememberFocus(s.focused)
	}
}

// RevokeSurfaces is for real native owner teardown, never a speculative change.
func (s *Session) RevokeSurfaces(parent Handle, reason string) error {
	if reason != "parent-closed" && reason != "parent-hidden" && reason != "reload" && reason != "dispose" {
		return fault("surface-reason", "Invalid native lifetime reason")
	}
	if parent != (Handle{}) && parent != s.rootOwner {
		if _, err := s.lookupControl(parent); err != nil {
			return err
		}
	}
	changed := false
	for _, d := range s.orderedSurfaces() {
		if d.Open && (parent == (Handle{}) || d.Parent == parent) {
			s.closeSurface(d.Target, "close", reason, 0, nil)
			changed = true
		}
	}
	if changed {
		s.StateRevision++
		s.refreshActivity()
	}
	return nil
}
func (s *Session) closeHiddenSurfaces() {
	for _, d := range s.orderedSurfaces() {
		if !d.Open {
			continue
		}
		a := s.active[d.InstancePath]
		opener := s.currentControl(d.Opener.Path)
		hidden := !a.visible
		if opener != nil && opener.Handle == d.Opener && !opener.Visible {
			hidden = true
		}
		if d.ParentSurface != nil {
			p := s.surfaces[d.ParentSurface.Handle.Path]
			if p == nil || !p.Open || p.Target != *d.ParentSurface {
				hidden = true
			}
		}
		if hidden {
			s.closeSurface(d.Target, "close", "parent-hidden", 0, nil)
		}
	}
}

func (s *Session) orderedSurfaces() []*SurfaceState {
	out := []*SurfaceState{}
	s.root.Walk(func(n *parser.Instance) {
		if d := s.surfaces[public(n.Path)]; d != nil && d.InstancePath == n.Path {
			out = append(out, d)
		}
	})
	return out
}

func (s *Session) surfaceFocus(owner, remembered string) string {
	if w := s.currentControl(remembered); w != nil && s.ownerSurface[w.InstancePath] == owner && s.inputAllowed(w) {
		return remembered
	}
	return s.firstSurfaceFocus(owner)
}

// RevokeSurface handles an already lost native surface. Unlike CloseSurface it
// cannot be vetoed by geometry; unlike RevokeSurfaces it includes the exact owner
// opening itself. A delayed native callback cannot revoke a replacement opening.
func (s *Session) RevokeSurface(target SurfaceTarget, reason string) error {
	if reason != "parent-closed" && reason != "parent-hidden" && reason != "reload" && reason != "dispose" {
		return fault("surface-reason", "Invalid native lifetime reason")
	}
	d := s.surfaces[target.Handle.Path]
	if d == nil || d.Target != target {
		return fault("stale-surface", "Native surface opening is no longer current")
	}
	if !d.Open {
		return nil
	}
	if s.closed {
		return fault("closed", "Session is closed")
	}
	s.closeSurface(target, "close", reason, 0, nil)
	s.StateRevision++
	s.refreshActivity()
	return nil
}
