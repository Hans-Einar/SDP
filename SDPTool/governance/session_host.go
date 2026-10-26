package governance

// SessionHostSelection resolves the current committed execution only. It is
// inspection data, not permission to dispatch or proof that a host is running.
type SessionHostSelection struct {
	InstanceID string       `json:"instanceId"`
	Revision   uint64       `json:"revision"`
	RunID      string       `json:"runId"`
	Outcome    string       `json:"outcome"`
	Skills     []SkillInput `json:"skills"`
}

func (s *Service) SessionHostSelection(id string) (SessionHostSelection, error) {
	st, _, e := s.load()
	if e != nil {
		return SessionHostSelection{}, e
	}
	i, ok := st.Sessions[id]
	if !ok {
		return SessionHostSelection{}, fail("missing_coverage", "Session instance absent")
	}
	if !i.Registered || len(i.Work) == 0 || i.Work[0].Effects.Status != "committed" {
		return SessionHostSelection{}, fail("missing_coverage", "Session has no current committed assignment; reconcile its workflow records before host dispatch")
	}
	q := i.Work[0].Request
	if _, ok := st.Runs[q.OperationID]; !ok {
		return SessionHostSelection{}, fail("reconciliation_required", "current Session run absent")
	}
	return SessionHostSelection{i.ID, i.Revision, q.OperationID, q.Outcome, append([]SkillInput{}, q.Skills...)}, nil
}
