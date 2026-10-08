package runtime

import (
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
)

func copyContext(c *ContextTarget) *ContextTarget {
	if c == nil {
		return nil
	}
	v := *c
	if c.Item != nil {
		x := *c.Item
		v.Item = &x
	}
	return &v
}
func copySurface(t *SurfaceTarget) *SurfaceTarget {
	if t == nil {
		return nil
	}
	v := *t
	return &v
}
func copyMenu(m *MenuState) MenuState {
	v := *m
	v.Context = copyContext(m.Context)
	v.Surface = copySurface(m.Surface)
	return v
}
func copySurfaceState(t *SurfaceState) SurfaceState {
	v := *t
	v.Context = copyContext(t.Context)
	v.ParentSurface = copySurface(t.ParentSurface)
	return v
}
func copyResults(in []DialogResult) []DialogResult {
	out := append([]DialogResult(nil), in...)
	for i := range out {
		out[i].Fields = copyDraftFields(out[i].Fields)
	}
	return out
}
func (s *Session) copyCommands(n *Session) {
	n.commands = map[string]*CommandState{}
	for k, v := range s.commands {
		x := *v
		n.commands[k] = &x
	}
	n.presentations = map[string]*CommandPresentation{}
	for k, v := range s.presentations {
		x := *v
		n.presentations[k] = &x
	}
	n.menus = map[string]*MenuState{}
	for k, v := range s.menus {
		x := copyMenu(v)
		n.menus[k] = &x
	}
	n.surfaces = map[string]*SurfaceState{}
	for k, v := range s.surfaces {
		x := copySurfaceState(v)
		n.surfaces[k] = &x
	}
	n.aux = map[string]*Widget{}
	for k, v := range s.aux {
		x := *v
		n.aux[k] = &x
	}
	n.publishedSurfaces = map[SurfaceTarget]bool{}
	for k, v := range s.publishedSurfaces {
		n.publishedSurfaces[k] = v
	}
	n.dialogResults = copyResults(s.dialogResults)
	n.activeSurface = copySurface(s.activeSurface)
	n.accepting = copySurface(s.accepting)
}
func (s *Session) snapshotCommands(v *Snapshot) {
	v.Commands = map[string]CommandState{}
	for _, c := range s.commands {
		v.Commands[c.InstancePath] = *c
	}
	v.Presentations = map[string]CommandPresentation{}
	for _, p := range s.presentations {
		v.Presentations[p.InstancePath] = *p
	}
	v.Menus = map[string]MenuState{}
	for _, m := range s.menus {
		v.Menus[m.InstancePath] = copyMenu(m)
	}
	v.Surfaces = map[string]SurfaceState{}
	for _, t := range s.surfaces {
		v.Surfaces[t.InstancePath] = copySurfaceState(t)
	}
	v.ActiveSurface = copySurface(s.activeSurface)
}
func (s *Session) Command(path string) (Handle, bool) {
	c := s.commands[path]
	if s.closed || c == nil {
		return Handle{}, false
	}
	return c.Handle, true
}
func (s *Session) CommandState(h Handle) (CommandState, bool) {
	c := s.commands[h.Path]
	if s.closed || c == nil || c.Handle != h {
		return CommandState{}, false
	}
	return *c, true
}
func (s *Session) Menu(path string) (Handle, bool) {
	m := s.menus[path]
	if s.closed || m == nil {
		return Handle{}, false
	}
	return m.Handle, true
}
func (s *Session) Surface(path string) (Handle, bool) {
	d := s.surfaces[path]
	if s.closed || d == nil {
		return Handle{}, false
	}
	return d.Handle, true
}
func (s *Session) SurfaceState(h Handle) (SurfaceState, bool) {
	d := s.surfaces[h.Path]
	if s.closed || d == nil || d.Handle != h {
		return SurfaceState{}, false
	}
	return copySurfaceState(d), true
}
func (s *Session) surfaceTarget(t SurfaceTarget) (*SurfaceState, error) {
	if s.closed {
		return nil, fault("closed", "Session is closed")
	}
	d := s.surfaces[t.Handle.Path]
	if d == nil || !d.Open || d.Target != t {
		return nil, fault("stale-surface", "Surface opening is no longer current")
	}
	return d, nil
}
func (s *Session) DialogFields(h Handle) ([]Widget, error) {
	d, ok := s.SurfaceState(h)
	if !ok {
		return nil, fault("surface", "Dialog is not current")
	}
	out := []Widget{}
	s.root.Walk(func(n *parser.Instance) {
		if w := s.widgets[public(n.Path)]; w != nil && w.InstancePath == n.Path && w.Handle.Kind == "input" && s.ownerSurface[n.Path] == d.Handle.Path {
			out = append(out, *w)
		}
	})
	return out, nil
}
func (s *Session) DialogField(h Handle, path string) (Handle, error) {
	d, ok := s.SurfaceState(h)
	if !ok {
		return Handle{}, fault("surface", "Dialog is not current")
	}
	node, err := parser.ResolveDialogField(s.root, d.InstancePath, path)
	if err != nil {
		return Handle{}, err
	}
	w := s.widgets[public(node.Path)]
	if w == nil {
		return Handle{}, fault("dialog-field", path)
	}
	return w.Handle, nil
}

func (s *Session) surfaceFor(w *Widget) *SurfaceState {
	if w == nil {
		return nil
	}
	return s.surfaces[s.ownerSurface[w.InstancePath]]
}
func (s *Session) inputAllowed(w *Widget) bool {
	if w == nil || !w.Enabled || !w.Visible {
		return false
	}
	for _, d := range s.surfaces {
		if d.Open && d.Modal && !s.surfaceDescends(s.surfaceFor(w), d.Target) {
			return false
		}
	}
	return true
}

func copyDraftFields(fields []DraftField) []DraftField {
	out := append([]DraftField(nil), fields...)
	for i := range out {
		out[i].RawDraft = copyString(out[i].RawDraft)
		if out[i].OptionTarget != nil {
			v := *out[i].OptionTarget
			out[i].OptionTarget = &v
		}
	}
	return out
}
func optionGeneration(t *OptionTarget) uint64 {
	if t == nil {
		return 0
	}
	return t.OptionGeneration
}
