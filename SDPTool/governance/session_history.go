package governance

import (
	"reflect"
)

const sessionLifecycleStateSchema = "sdp-project-governance-state/0.20"

type SessionPast struct {
	Execution  SessionInstance `json:"execution"`
	Transition SessionAdvance  `json:"transition"`
}

// Stored history is flat. Prefixes are reconstructed only while validating or
// reading an earlier execution; old run/receipt identities remain addressable.
func sessionExecutions(st snapshot) []SessionInstance {
	out := []SessionInstance{}
	for _, i := range st.Sessions {
		for n, p := range i.Past {
			old := p.Execution
			old.Past = i.Past[:n]
			out = append(out, old)
		}
		out = append(out, i)
	}
	return out
}
func sessionSelectionSources(i SessionInstance) []Source {
	if len(i.Past) > 0 {
		return i.Past[len(i.Past)-1].Transition.Effects.Coverage
	}
	return i.Context.Sources
}
func sessionStepDone(i SessionInstance, id string) bool {
	if i.Selection != nil && i.Selection.Request.StepID == id && i.Disposition != nil && i.Disposition.Effects.Status == "committed" {
		return true
	}
	for _, p := range i.Past {
		if p.Execution.Selection != nil && p.Execution.Selection.Request.StepID == id && p.Execution.Disposition != nil && p.Execution.Disposition.Effects.Status == "committed" {
			return true
		}
	}
	return false
}
func sessionDependenciesSatisfied(i SessionInstance, step SessionStepBinding) bool {
	for _, id := range step.DependsOn {
		if !sessionStepDone(i, id) {
			return false
		}
	}
	return true
}
func validateSessionHistoryShape(st snapshot) error {
	for _, i := range st.Sessions {
		if len(i.Past) > 64 {
			return fail("reconciliation_required", "Session execution history capacity exceeded")
		}
		if (len(i.Past) > 0 || i.Advance != nil) && st.Schema != sessionLifecycleStateSchema {
			return fail("incompatible_profile", "Session lifecycle requires profile 0.20")
		}
		previous := uint64(0)
		for _, p := range i.Past {
			old := p.Execution
			q := p.Transition.Request
			if len(old.Past) != 0 || old.Advance != nil || old.ID != i.ID || old.Revision != q.ExpectedRevision || old.Revision <= previous || !old.Registered || !reflect.DeepEqual(old.Context, i.Context) || !reflect.DeepEqual(old.Registration, i.Registration) || old.Reason != i.Reason || p.Transition.Effects.Status != "committed" || q.Action == "complete" {
				return fail("reconciliation_required", "invalid flat Session execution history")
			}
			previous = q.ExpectedRevision + 2
		}
		if len(i.Past) > 0 && i.Revision < previous {
			return fail("reconciliation_required", "Session history exceeds current revision")
		}
	}
	return nil
}

type SessionHistoryView struct {
	StepID      string `json:"stepId"`
	RunID       string `json:"runId"`
	Action      string `json:"action"`
	OperationID string `json:"operationId"`
	Reason      string `json:"reason"`
}

func (s *Service) projectSessionHistory(i SessionInstance, st snapshot, v *SessionView) {
	v.State = "active"
	for _, p := range i.Past {
		old := p.Execution
		v.History = append(v.History, SessionHistoryView{selectedSessionStep(old).ID, old.Work[0].Request.OperationID, p.Transition.Request.Action, p.Transition.Request.OperationID, p.Transition.Request.Reason})
		if old.Disposition != nil && old.Disposition.Effects.Status == "committed" {
			for n := range v.Steps {
				step := &v.Steps[n]
				if step.ID == selectedSessionStep(old).ID && !step.Selected {
					step.Stage = "done"
					step.Holds = []string{}
					step.RunID = old.Work[0].Request.OperationID
					step.RunState = "completed"
					step.DispositionID = old.Disposition.Request.OperationID
					step.CardDispositions = old.Disposition.Request.Cards
				}
			}
		}
	}
	for n, step := range i.Context.Binding.Steps {
		if !v.Steps[n].Selected && v.Steps[n].Stage != "done" {
			v.Steps[n].Holds = append([]string{}, v.Conditions...)
			if !sessionDependenciesSatisfied(i, step) {
				v.Steps[n].Holds = append(v.Steps[n].Holds, "predecessor steps not completed")
			} else {
				v.Steps[n].Holds = append(v.Steps[n].Holds, "select this step before assessment and execution")
			}
		}
	}
	if i.Advance != nil {
		v.RecordEffects = i.Advance.Effects.Status
		v.NextActions = []string{"session-read", "session-resume"}
		if i.Advance.Effects.Status == "committed" && i.Advance.Request.Action == "complete" {
			v.State = "completed"
			v.NextActions = []string{"session-read"}
		}
		return
	}
	if i.Selection == nil && len(i.Past) > 0 {
		v.RecordEffects = i.Past[len(i.Past)-1].Transition.Effects.Status
	}
	if v.Coverage != "current" || i.Selection == nil {
		return
	}
	eligible := i.Disposition != nil && i.Disposition.Effects.Status == "committed"
	if g := latestOwnerGate(i); g != nil && g.Effects.Status == "committed" && g.Request.Action == "changes-requested" {
		eligible = true
	}
	if len(i.Work) >= 2 && i.Work[1].Effects.Status == "committed" {
		ev := i.Work[1].Evidence
		for _, r := range st.Receipts {
			if t := r.Result.Evidence; t != nil && t.Request.Kind == "test" && t.Request.Verdict == "failed" && sameSessionCandidate(*t, *ev) {
				eligible = true
			}
		}
		for _, r := range st.Reviews {
			if r.RunID == i.Work[0].Request.OperationID && r.State == "completed" && r.Outcome != nil && r.Outcome.Verdict != "passed" && r.Spec.Candidate.Digest == ev.Request.Candidate.Digest {
				eligible = true
			}
		}
	}
	if eligible {
		v.NextActions = append(v.NextActions, "session-advance")
	}
}
