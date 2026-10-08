package runtime

import (
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"math"
)

type ViewportState struct{ X, Y float64 }
type Snapshot struct {
	Commands                               map[string]CommandState
	Presentations                          map[string]CommandPresentation
	Menus                                  map[string]MenuState
	Surfaces                               map[string]SurfaceState
	ActiveSurface                          *SurfaceTarget
	Tabs                                   map[string]TabsState
	Splits                                 map[string]SplitState
	Focused                                string
	Root                                   *parser.Instance
	ModelRevision, StateRevision, Sequence uint64
	Collections                            map[string]CollectionState
	Viewports                              map[string]ViewportState
	ViewportHandles                        map[string]Handle
}

// StateGate validates prospective state and returns all effective viewport offsets.
// It is pure: no live mutation, provider loading or user callbacks are permitted.
type StateGate func(Snapshot) (map[string]ViewportState, error)

func copyViewports(in map[string]ViewportState) map[string]ViewportState {
	out := make(map[string]ViewportState, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
func (s *Session) Snapshot() Snapshot {
	v := Snapshot{Tabs: map[string]TabsState{}, Splits: map[string]SplitState{}, Focused: s.focused, Root: s.SnapshotRoot(), ModelRevision: s.Revision, StateRevision: s.StateRevision, Sequence: s.sequence, Collections: map[string]CollectionState{}, Viewports: copyViewports(s.viewports), ViewportHandles: map[string]Handle{}}
	s.snapshotCommands(&v)
	for _, t := range s.tabs {
		v.Tabs[t.InstancePath] = copyTabs(t)
	}
	for _, p := range s.splits {
		v.Splits[p.InstancePath] = *p
	}
	for k, h := range s.viewportHandles {
		v.ViewportHandles[k] = h
	}
	for _, c := range s.collections {
		v.Collections[c.InstancePath] = copyCollection(c).CollectionState
	}
	return v
}
func (s *Session) copyState() *Session {
	n := *s
	n.widgets = make(map[string]*Widget, len(s.widgets))
	for k, v := range s.widgets {
		w := *v
		n.widgets[k] = &w
	}
	n.collections = make(map[string]*collection, len(s.collections))
	for k, v := range s.collections {
		n.collections[k] = copyCollection(v)
	}
	n.viewports = copyViewports(s.viewports)
	s.copyPanes(&n)
	s.copyCommands(&n)
	return &n
}
func (s *Session) validateViewports(offsets map[string]ViewportState) error {
	owners := map[string]*parser.Instance{}
	s.root.Walk(func(n *parser.Instance) {
		if n.Layout["overflow-x"] == "scroll" || n.Layout["overflow-y"] == "scroll" {
			owners[n.Path] = n
		}
	})
	for path, v := range offsets {
		n := owners[path]
		if n == nil {
			return fault("viewport", "Unknown scroll owner: "+path)
		}
		if math.IsNaN(v.X) || math.IsNaN(v.Y) || math.IsInf(v.X, 0) || math.IsInf(v.Y, 0) || v.X < 0 || v.Y < 0 {
			return fault("viewport", "Offsets must be finite and nonnegative")
		}
		if n.Layout["overflow-x"] != "scroll" && v.X != 0 || n.Layout["overflow-y"] != "scroll" && v.Y != 0 {
			return fault("viewport", "Offset on non-scroll axis")
		}
	}
	return nil
}
func (s *Session) gate() error {
	if s.check != nil {
		if err := s.check(s.SnapshotRoot()); err != nil {
			return err
		}
	}
	if s.presentationCheck != nil {
		p, err := s.presentationCheck(s.Snapshot())
		if err != nil {
			return err
		}
		return s.measuredState(p)
	}
	if s.stateCheck != nil {
		offsets, err := s.stateCheck(s.Snapshot())
		if err != nil {
			return err
		}
		if err = s.validateViewports(offsets); err != nil {
			return err
		}
		return s.measuredState(PresentationState{Viewports: offsets})
	}
	return nil
}
func (s *Session) publish(n *Session) error {
	if s.closed {
		return fault("closed", "Session is closed")
	}
	revision, state := s.Revision, s.StateRevision
	n.StateRevision = state + 1
	n.refreshActivity()
	n.closeHiddenSurfaces()
	n.refreshActivity()
	if err := n.gate(); err != nil {
		return err
	}
	var ticket PresentationTicket
	if n.presentationPrepare != nil {
		var err error
		ticket, err = n.presentationPrepare(n.Snapshot())
		if err != nil {
			if ticket.Discard != nil {
				ticket.Discard()
			}
			return err
		}
	}
	if s.closed || s.Revision != revision || s.StateRevision != state {
		if ticket.Discard != nil {
			ticket.Discard()
		}
		return fault("stale-presentation", "State changed during presentation preparation")
	}
	n.strictSplit = ""
	n.interacting = s.interacting
	n.accepting = copySurface(s.accepting)
	*s = *n
	if ticket.Publish != nil {
		ticket.Publish()
	}
	return nil
}
func (s *Session) CheckStateWith(check StateGate) error {
	if len(s.splits) > 0 {
		return fault("presentation-gate", "Splits require a typed presentation gate")
	}
	n := s.copyState()
	n.presentationCheck = nil
	n.stateCheck = check
	return s.publish(n)
}

// SetViewports merges requested offsets. The gate returns the complete clamped map.
func (s *Session) SetViewports(offsets map[string]ViewportState) error {
	if err := s.validateViewports(offsets); err != nil {
		return err
	}
	if s.stateCheck == nil && s.presentationCheck == nil && len(offsets) > 0 {
		return fault("viewport", "Viewport updates require a state-aware geometry gate")
	}
	n := s.copyState()
	for k, v := range offsets {
		n.viewports[k] = v
	}
	return s.publish(n)
}
func (s *Session) Sequence() uint64 { return s.sequence }

// Viewport returns a separate owner identity keyed by exact normalized path.
func (s *Session) Viewport(path string) (Handle, bool) {
	if s.closed {
		return Handle{}, false
	}
	h, ok := s.viewportHandles[path]
	return h, ok
}

// SetViewport is the checked native operation. SetViewports is trusted bulk input.
func (s *Session) SetViewport(h Handle, modelRevision uint64, offset ViewportState) error {
	if s.closed {
		return fault("closed", "Session is closed")
	}
	if modelRevision != s.Revision {
		return fault("stale-event", "Viewport belongs to another model revision")
	}
	if current, ok := s.viewportHandles[h.Path]; !ok || current != h {
		return fault("stale-handle", "Viewport owner is no longer current")
	}
	if len(s.panes) > 0 || len(s.aux) > 0 || len(s.commands) > 0 {
		if a := s.active[h.Path]; !a.visible || !a.enabled {
			return fault("inactive-widget", "Viewport is inactive")
		}
	}
	for _, d := range s.surfaces {
		if d.Open && d.Modal && !s.surfaceDescends(s.surfaces[s.ownerSurface[h.Path]], d.Target) {
			return fault("modal", "Viewport is blocked")
		}
	}
	return s.SetViewports(map[string]ViewportState{h.Path: offset})
}
