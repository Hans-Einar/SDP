package governance

import "path/filepath"

// RunControllerLease excludes another controller operation on this run while
// retaining exclusion with legacy checkout-wide controllers.
func (s *Service) RunControllerLease(c Caller, runID string) (func(), error) {
	return s.runLease(c, runID, "controller")
}

// RunLiveControllerLease owns one run's terminal session, not every run in the
// checkout. OS locks, not the presence of directories, establish ownership.
func (s *Service) RunLiveControllerLease(c Caller, runID string) (func(), error) {
	return s.runLease(c, runID, "live-controller")
}

func (s *Service) runLease(c Caller, runID, kind string) (func(), error) {
	if !c.controller {
		return nil, fail("unauthorized", "controller required")
	}
	if !uuidPattern.MatchString(runID) {
		return nil, fail("invalid_request", "valid run identity required")
	}
	state, _, err := s.load()
	if err != nil {
		return nil, err
	}
	if _, ok := state.Runs[runID]; !ok {
		return nil, fail("missing_coverage", "run absent")
	}
	dir := filepath.Join(s.state, kind)
	compatibility, err := sharedLock(dir)
	if err != nil {
		return nil, err
	}
	run, err := lock(filepath.Join(dir, "runs", runID))
	if err != nil {
		compatibility()
		return nil, err
	}
	return func() { run(); compatibility() }, nil
}
