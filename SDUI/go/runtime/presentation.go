package runtime

import "math"

type SplitGeometry struct{ Lower, Upper, Effective float64 }
type PresentationState struct {
	Viewports map[string]ViewportState
	Splits    map[string]SplitGeometry
}
type PresentationGate func(Snapshot) (PresentationState, error)

// PresentationTicket belongs to one finalized mutation. Publish cannot fail or
// reenter runtime; Discard releases only this ticket, never accepted resources.
type PresentationTicket struct {
	Publish func()
	Discard func()
}
type PresentationPrepare func(Snapshot) (PresentationTicket, error)

func finiteRatio(v float64) bool           { return !math.IsNaN(v) && !math.IsInf(v, 0) && v >= 0 && v <= 1 }
func clampRatio(v, lo, hi float64) float64 { return math.Max(lo, math.Min(hi, v)) }
func (s *Session) CheckPresentationWith(gate PresentationGate) error {
	if len(s.splits) > 0 && gate == nil {
		return fault("presentation-gate", "Splits require a typed presentation gate")
	}
	n := s.copyState()
	n.presentationCheck = gate
	n.stateCheck = nil
	return s.publish(n)
}
func (s *Session) PreparePresentationWith(prepare PresentationPrepare) error {
	n := s.copyState()
	n.presentationPrepare = prepare
	return s.publish(n)
}
func (s *Session) measuredState(p PresentationState) error {
	if err := s.validateViewports(p.Viewports); err != nil {
		return err
	}
	expected := 0
	for _, v := range s.splits {
		if !v.Visible || v.Collapsed != SplitNone {
			continue
		}
		expected++
		g, ok := p.Splits[v.InstancePath]
		if !ok || !finiteRatio(g.Lower) || !finiteRatio(g.Upper) || !finiteRatio(g.Effective) || g.Lower > g.Upper || g.Lower < v.MinFirst || g.Upper > 1-v.MinSecond {
			return fault("split-geometry", "Missing or invalid split bounds")
		}
		if g.Effective != clampRatio(v.Proportion, g.Lower, g.Upper) {
			return fault("split-geometry", "Effective ratio is not the measured clamp")
		}
		if s.strictSplit == v.Handle.Path && g.Effective != v.Proportion {
			return fault("split-proportion", "Programmatic proportion exceeds measured bounds")
		}
		v.Proportion = g.Effective
		v.SavedProportion = g.Effective
	}
	if len(p.Splits) != expected {
		return fault("split-geometry", "Unexpected split geometry")
	}
	offsets := copyViewports(p.Viewports)
	for path, v := range s.viewports {
		if s.inactivePaneViewport(path) {
			offsets[path] = v
		}
	}
	s.viewports = offsets
	return nil
}
