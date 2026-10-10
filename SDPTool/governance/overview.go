package governance

import (
	"sort"
	"time"
)

// Overview is a projection of registered work in one checkout store. It does not
// enumerate unregistered project work or probe hosts. Generation pins operational
// data; Coverage is a separate per-run filesystem check, not an atomic snapshot.
type Overview struct {
	Schema      string     `json:"schema"`
	ProjectRef  string     `json:"projectRef"`
	Root        string     `json:"root"`
	CheckoutID  string     `json:"checkoutId"`
	Generation  string     `json:"generation"`
	InspectedAt string     `json:"inspectedAt"`
	Scope       string     `json:"scope"`
	Runs        []WorkView `json:"runs"`
}

type WorkView struct {
	RunID           string    `json:"runId"`
	Revision        uint64    `json:"revision"`
	Outcome         string    `json:"outcome"`
	Selection       Selection `json:"selection"`
	InputRevision   string    `json:"inputRevision"`
	RoutineState    string    `json:"routineState"`
	RoutineStep     string    `json:"routineStep"`
	Coverage        string    `json:"coverage"`
	Conditions      []string  `json:"conditions"`
	AssignedRole    string    `json:"assignedRole"`
	RoleEvidence    string    `json:"roleEvidence"`
	AttemptID       string    `json:"attemptId"`
	ThreadID        string    `json:"threadId"`
	TurnID          string    `json:"turnId"`
	LastHostStatus  string    `json:"lastHostStatus"`
	LastObservedAt  string    `json:"lastObservedAt"`
	ObservationGap  bool      `json:"observationGap"`
	LiveState       string    `json:"liveState"`
	ToolCompletions int       `json:"toolCompletions"`
	LastToolItemID  string    `json:"lastToolItemId"`
	LastToolStatus  string    `json:"lastToolStatus"`
	ChildMapping    string    `json:"childMapping"`
	RequestedSkills int       `json:"requestedSkills"`
	SkillAttemptID  string    `json:"skillAttemptId"`
	SkillInputID    string    `json:"skillInputId"`
	SkillAcceptance string    `json:"skillAcceptance"`
	SkillCompliance string    `json:"skillCompliance"`
	OwnerAcceptance string    `json:"ownerAcceptance"`
}

func (s *Service) Overview() (Overview, error) {
	out := Overview{Schema: "sdp-governance-overview/1", ProjectRef: s.project, Root: s.root,
		InspectedAt: time.Now().UTC().Format(time.RFC3339Nano), Scope: "registered-runs-in-selected-checkout", Runs: []WorkView{}}
	state, head, err := s.load()
	if err != nil {
		return out, err
	}
	out.CheckoutID, out.Generation = state.CheckoutID, head
	ids := make([]string, 0, len(state.Runs))
	for id := range state.Runs {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		r := state.Runs[id]
		v, err := workView(r)
		if err != nil {
			return out, err
		}
		c, err := s.Context(r.Selection)
		switch {
		case err != nil:
			v.Coverage = "unavailable"
			v.Conditions = append(v.Conditions, err.Error())
		case c.Revision != r.InputRevision || len(c.Conditions) > 0:
			v.Coverage = "stale"
			v.Conditions = append(v.Conditions, c.Conditions...)
			if c.Revision != r.InputRevision {
				v.Conditions = append(v.Conditions, "governing sources differ from pinned input revision")
			}
		default:
			v.Coverage = "known"
		}
		out.Runs = append(out.Runs, v)
	}
	return out, nil
}

func workView(r Run) (WorkView, error) {
	v := WorkView{RunID: r.ID, Revision: r.Revision, Outcome: r.Assignment.Outcome,
		Selection: r.Selection, InputRevision: r.InputRevision, RoutineState: r.State,
		RoutineStep: r.Step, Conditions: []string{}, AssignedRole: "unknown", RoleEvidence: "unavailable",
		LastHostStatus: "not-observed", LiveState: "unknown", ChildMapping: "not-observed",
		SkillAcceptance: "not-observed", SkillCompliance: "not-observed", OwnerAcceptance: "not-observed"}
	p, err := runProvenance(r)
	if err != nil {
		return v, err
	}
	if p.Context != nil {
		v.AssignedRole, v.RoleEvidence = p.Context.AssignedRole, p.Status
	}
	if p.Invocation != nil {
		v.RequestedSkills, v.SkillAcceptance = len(p.Invocation.Skills), p.Invocation.Acceptance
		v.SkillAttemptID, v.SkillInputID = p.Invocation.AttemptID, p.Invocation.InputID
	}
	if len(r.Attempts) == 0 {
		return v, nil
	}
	a := r.Attempts[len(r.Attempts)-1]
	v.AttemptID, v.ThreadID, v.TurnID = a.ID, a.ThreadID, a.TurnID
	v.LastHostStatus, v.ObservationGap = a.Status, a.Gap
	for _, e := range a.Events {
		v.LastObservedAt = e.ReceivedAt
		if e.Kind == "tool.completed" {
			v.ToolCompletions++
			v.LastToolItemID, v.LastToolStatus = e.ItemID, e.Status
		}
	}
	return v, nil
}
