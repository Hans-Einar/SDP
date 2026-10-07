package runtime

import (
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"math"
)

type ViewportState struct{ X, Y float64 }
type Snapshot struct {
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
	v := Snapshot{Root: s.SnapshotRoot(), ModelRevision: s.Revision, StateRevision: s.StateRevision, Sequence: s.sequence, Collections: map[string]CollectionState{}, Viewports: copyViewports(s.viewports), ViewportHandles: map[string]Handle{}}
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
	if s.stateCheck != nil {
		offsets, err := s.stateCheck(s.Snapshot())
		if err != nil {
			return err
		}
		if err = s.validateViewports(offsets); err != nil {
			return err
		}
		s.viewports = copyViewports(offsets)
	}
	return nil
}
func (s *Session) publish(n *Session) error {
	if s.closed {
		return fault("closed", "Session is closed")
	}
	n.StateRevision = s.StateRevision + 1
	if err := n.gate(); err != nil {
		return err
	}
	*s = *n
	return nil
}
func (s *Session) CheckStateWith(check StateGate) error {
	n := s.copyState()
	n.stateCheck = check
	return s.publish(n)
}

// SetViewports merges requested offsets. The gate returns the complete clamped map.
func (s *Session) SetViewports(offsets map[string]ViewportState) error {
	if err := s.validateViewports(offsets); err != nil {
		return err
	}
	if s.stateCheck == nil && len(offsets) > 0 {
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
	return s.SetViewports(map[string]ViewportState{h.Path: offset})
}
