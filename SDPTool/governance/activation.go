package governance

import (
	"reflect"
	"strings"
)

const ActivationSchema = "sdp-governance-activation/1"
const activationStateSchema = "sdp-project-governance-state/0.5"

type ActivationRequest struct {
	Schema              string             `json:"schema"`
	ProjectRef          string             `json:"projectRef"`
	CheckoutID          string             `json:"checkoutId"`
	OperationID         string             `json:"operationId"`
	RequestID           string             `json:"requestId"`
	Review              DesignCheckRequest `json:"review"`
	ReviewRevision      string             `json:"reviewRevision"`
	DecisionOperationID string             `json:"decisionOperationId"`
}

type AssignmentRun struct {
	RunID               string `json:"runId"`
	PreparationRevision uint64 `json:"preparationRevision"`
	ReviewRevision      string `json:"reviewRevision"`
	DecisionOperationID string `json:"decisionOperationId"`
}

func Implementation() Routine {
	r := Maintenance()
	r.ID = "implementation"
	return r
}

func (s *Service) ActivateAssignment(caller Caller, q ActivationRequest) (Result, error) {
	var out Result
	if !caller.controller {
		return out, fail("unauthorized", "activation requires trusted controller")
	}
	if q.Schema != ActivationSchema {
		return out, fail("incompatible_profile", "unsupported activation schema")
	}
	if !uuidPattern.MatchString(q.RequestID) || !uuidPattern.MatchString(q.DecisionOperationID) || !digestPattern.MatchString(q.ReviewRevision) {
		return out, fail("invalid_request", "request identity and exact decision/review required")
	}
	unlock, err := lock(s.state)
	if err != nil {
		return out, err
	}
	defer unlock()
	state, head, err := s.load()
	if err != nil {
		return out, err
	}
	if err = s.authorize(caller, Schema, q.ProjectRef, q.CheckoutID, q.OperationID, state); err != nil {
		return out, err
	}
	digest := requestDigest(q)
	if prior, ok, e := retry(state, q.OperationID, digest); ok {
		return prior, e
	}
	if _, ok := state.AssignmentRuns[q.Review.AssignmentID]; ok {
		return out, fail("conflict", "assignment already has a run; inspect or retry original operation")
	}
	if q.Review.ExpectedGeneration != head {
		return out, fail("stale_revision", "operational generation changed")
	}
	if len(state.Runs) >= 128 || len(state.Receipts) >= 1024 {
		return out, fail("limit", "run/receipt capacity reached")
	}
	for _, r := range state.Runs {
		if r.RequestID == q.RequestID {
			return out, fail("conflict", "request already has a run")
		}
	}
	review, err := s.ReviewAssignmentDesign(q.Review)
	if err != nil {
		return out, err
	}
	if review.Acceptance != "approved" || review.Decision == nil || review.Decision.OperationID != q.DecisionOperationID || review.Check.ReviewRevision != q.ReviewRevision {
		return out, fail("stale_revision", "activation requires the current approved decision on this exact package")
	}
	p := review.Check.Assignment
	if why := scopeConflictExcept(p.Assignment.WritablePaths, state, p.ID); why != "" {
		return out, fail("conflict", why)
	}
	id, err := identifier()
	if err != nil {
		return out, err
	}
	routine := Implementation()
	r := Run{ID: id, Revision: 1, RequestID: q.RequestID, State: "ready", Step: "implement", InputRevision: p.InputRevision, Selection: p.Selection, Assignment: p.Assignment, Routine: routine, RoutineDigest: routine.Digest(), Reason: "trusted controller activated exact approved implementation package"}
	out = Result{Schema: Schema, Operation: "activate", ProjectRef: s.project, CheckoutID: state.CheckoutID, OperationID: q.OperationID, Run: summarizeRun(r), Coverage: "known"}
	final, err := s.ReviewAssignmentDesign(q.Review)
	if err != nil {
		return Result{}, err
	}
	if !reflect.DeepEqual(review, final) {
		return Result{}, fail("stale_revision", "design decision changed during activation")
	}
	if state.AssignmentRuns == nil {
		state.AssignmentRuns = map[string]AssignmentRun{}
	}
	state.AssignmentRuns[p.ID] = AssignmentRun{RunID: id, PreparationRevision: p.Revision, ReviewRevision: q.ReviewRevision, DecisionOperationID: q.DecisionOperationID}
	if state.Schema != attentionStateSchema {
		state.Schema = activationStateSchema
	}
	state.Runs[id] = r
	state.Receipts[q.OperationID] = receipt{digest, out}
	if err = s.save(state, head); err != nil {
		return Result{}, err
	}
	return out, nil
}

func assignmentRun(state snapshot, runID string) (PreparedAssignment, AssignmentRun, bool) {
	for id, b := range state.AssignmentRuns {
		if b.RunID == runID {
			return state.Prepared[id], b, true
		}
	}
	return PreparedAssignment{}, AssignmentRun{}, false
}

// RunContext resolves the run's actual plan and, for implementation, its current
// design disposition. Conditions block advancement, not stop/observation records.
func (s *Service) RunContext(runID string) (Context, error) {
	state, _, err := s.load()
	if err != nil {
		return Context{}, err
	}
	r, ok := state.Runs[runID]
	if !ok {
		return Context{}, fail("missing_coverage", "run not found")
	}
	return s.runContext(r)
}
func (s *Service) runContext(r Run) (Context, error) {
	if r.Routine.ID == "maintenance" {
		return s.Context(r.Selection)
	}
	c, err := s.contextFor(r.Selection, []string{"ImplementationPlan"}, "PLAN-")
	if err != nil {
		return c, err
	}
	state, head, err := s.load()
	if err != nil {
		return c, err
	}
	p, b, ok := assignmentRun(state, r.ID)
	if !ok {
		return c, fail("reconciliation_required", "implementation run lacks approved assignment binding")
	}
	recorded := state.DesignDecisionReceipts[b.DecisionOperationID].Result.Record.Basis
	q := DesignCheckRequest{Schema: DesignCheckSchema, AssignmentID: p.ID, ExpectedRevision: p.Revision, ExpectedGeneration: head, BlueprintAssignment: recorded.Provider.AssignmentID, BlueprintEvent: recorded.Provider.EventID}
	review, err := s.ReviewAssignmentDesign(q)
	if err != nil {
		c.Coverage = "unavailable"
		c.Conditions = append(c.Conditions, "design review unavailable: "+err.Error())
	} else if review.Acceptance != "approved" || review.Decision == nil || review.Decision.OperationID != b.DecisionOperationID || review.Check.ReviewRevision != b.ReviewRevision {
		c.Coverage = "stale"
		c.Conditions = append(c.Conditions, "bound design decision is no longer current and approved")
		c.Conditions = append(c.Conditions, review.Check.Conditions...)
	}
	return c, nil
}

// RunDesign returns the immutable controller decision sent to the implementation
// host. Its presence is historical; RunContext must pass again before dispatch.
func (s *Service) RunDesign(runID string) (*DesignDecisionResult, error) {
	state, _, err := s.load()
	if err != nil {
		return nil, err
	}
	if _, ok := state.Runs[runID]; !ok {
		return nil, fail("missing_coverage", "run not found")
	}
	_, b, ok := assignmentRun(state, runID)
	if !ok {
		return nil, nil
	}
	result := state.DesignDecisionReceipts[b.DecisionOperationID].Result
	return &result, nil
}

func validateAssignmentRuns(state snapshot) error {
	if state.Schema != activationStateSchema && state.Schema != attentionStateSchema {
		if state.AssignmentRuns != nil {
			return fail("incompatible_profile", "assignment activation requires profile 0.5")
		}
	} else if state.Schema == activationStateSchema && len(state.AssignmentRuns) == 0 {
		return fail("reconciliation_required", "activation profile has no bindings")
	}
	seen := map[string]bool{}
	for id, b := range state.AssignmentRuns {
		p, ok := state.Prepared[id]
		r, exists := state.Runs[b.RunID]
		d, decision := state.DesignDecisionReceipts[b.DecisionOperationID]
		if !ok || !exists || !decision || seen[b.RunID] || p.Role != "Master" || p.State == "retired" || p.Revision != b.PreparationRevision ||
			d.Result.Record.Input.Decision != "approved" || d.Result.Record.Basis.Assignment.ID != id || d.Result.Record.Basis.ReviewRevision != b.ReviewRevision ||
			!reflect.DeepEqual(d.Result.Record.Basis.Assignment, p) || !reflect.DeepEqual(r.Assignment, p.Assignment) || r.Selection != p.Selection || r.InputRevision != p.InputRevision ||
			r.RoutineDigest != Implementation().Digest() || r.Routine.Digest() != r.RoutineDigest {
			return fail("reconciliation_required", "invalid immutable assignment/run binding")
		}
		seen[b.RunID] = true
	}
	for id, r := range state.Runs {
		if r.Routine.ID != "maintenance" && (!seen[id] || r.Routine.ID != "implementation") {
			return fail("reconciliation_required", "unsupported or unbound run routine")
		}
	}
	return nil
}

// CheckRunReady is a final optimistic check at the controller's host-send boundary.
func (s *Service) CheckRunReady(runID string) error {
	result, err := s.Read(runID)
	if err != nil {
		return err
	}
	c, err := s.RunContext(runID)
	if err != nil {
		return err
	}
	if c.Revision != result.Run.InputRevision || len(c.Conditions) > 0 {
		return fail("stale_revision", "run sources or authority changed: "+strings.Join(c.Conditions, "; "))
	}
	return nil
}
