package runtime

import "math"

func (s *Session) selectPage(h Handle, id string, user bool) error {
	t, ok := s.Tabs(h)
	if !ok {
		return fault("pane-kind", "Selection requires tabs")
	}
	if user && (!t.Enabled || !t.Visible) {
		return fault("inactive-widget", "Tabs is inactive")
	}
	found := false
	for _, p := range t.Pages {
		if p.ID == id && p.Enabled && p.Visible {
			found = true
		}
	}
	if !found {
		return fault("page-selection", "Page is unknown or ineligible")
	}
	s.tabs[h.Path].Selected = id
	if user {
		s.focused = h.Path
	}
	return nil
}
func (s *Session) SelectPage(h Handle, id string) error {
	n := s.copyState()
	if err := n.selectPage(h, id, false); err != nil {
		return err
	}
	if s.tabs[h.Path].Selected == id {
		return nil
	}
	return s.publish(n)
}
func (s *Session) changeSplit(h Handle, change SplitChange, strict bool) error {
	p, ok := s.Split(h)
	if !ok {
		return fault("pane-kind", "Adjustment requires split")
	}
	if s.presentationCheck == nil {
		return fault("presentation-gate", "Split adjustment requires measured geometry")
	}
	if change.Operation != "ratio" && change.Proportion != 0 {
		return fault("event-type", "Nonratio operation cannot carry a proportion")
	}
	target := s.splits[h.Path]
	switch change.Operation {
	case "ratio":
		if math.IsNaN(change.Proportion) || math.IsInf(change.Proportion, 0) {
			return fault("split-proportion", "Proportion must be finite")
		}
		if p.Collapsed != SplitNone {
			return fault("split-state", "Restore a collapsed split before changing its ratio")
		}
		if strict && (!finiteRatio(change.Proportion) || change.Proportion < p.MinFirst || change.Proportion > 1-p.MinSecond) {
			return fault("split-proportion", "Programmatic proportion exceeds declared bounds")
		}
		target.Proportion = clampRatio(change.Proportion, 0, 1)
		if strict {
			s.strictSplit = h.Path
		}
	case "collapse-first", "collapse-second":
		if !p.Collapsible {
			return fault("split-state", "Split is not collapsible")
		}
		if p.Collapsed == SplitNone {
			target.SavedProportion = p.Proportion
		}
		target.Collapsed = SplitFirst
		if change.Operation == "collapse-second" {
			target.Collapsed = SplitSecond
		}
	case "restore":
		if p.Collapsed != SplitNone {
			target.Collapsed = SplitNone
			target.Proportion = p.SavedProportion
		}
	default:
		return fault("event-type", "Unknown split operation")
	}
	return nil
}
func (s *Session) SetSplitProportion(h Handle, p float64) error {
	n := s.copyState()
	if err := n.changeSplit(h, SplitChange{"ratio", p}, true); err != nil {
		return err
	}
	return s.publish(n)
}
func (s *Session) CollapseSplit(h Handle, side SplitSide) error {
	op := "collapse-first"
	if side == SplitSecond {
		op = "collapse-second"
	} else if side != SplitFirst {
		return fault("split-state", "Collapse requires first or second")
	}
	n := s.copyState()
	if err := n.changeSplit(h, SplitChange{Operation: op}, false); err != nil {
		return err
	}
	return s.publish(n)
}
func (s *Session) RestoreSplit(h Handle) error {
	n := s.copyState()
	if err := n.changeSplit(h, SplitChange{Operation: "restore"}, false); err != nil {
		return err
	}
	return s.publish(n)
}
