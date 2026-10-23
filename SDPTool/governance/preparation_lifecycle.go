package governance

import (
	"encoding/json"
	"strings"
)

const PreparationLifecycleSchema = "sdp-governance-assignment/2"
const preparationLifecycleStateSchema = "sdp-project-governance-state/0.3"

// PreparationChange identifies exactly the unlaunched work the controller changes.
type PreparationChange struct {
	Schema             string `json:"schema"`
	ProjectRef         string `json:"projectRef"`
	CheckoutID         string `json:"checkoutId"`
	OperationID        string `json:"operationId"`
	ExpectedGeneration string `json:"expectedGeneration"`
	AssignmentID       string `json:"assignmentId"`
	ExpectedRevision   uint64 `json:"expectedRevision"`
	Reason             string `json:"reason"`
}

// RevisePreparationRequest replaces the package, preserving identity/role/selection.
type RevisePreparationRequest struct {
	Schema             string        `json:"schema"`
	ProjectRef         string        `json:"projectRef"`
	CheckoutID         string        `json:"checkoutId"`
	OperationID        string        `json:"operationId"`
	ExpectedGeneration string        `json:"expectedGeneration"`
	AssignmentID       string        `json:"assignmentId"`
	ExpectedRevision   uint64        `json:"expectedRevision"`
	Reason             string        `json:"reason"`
	InputRevision      string        `json:"inputRevision"`
	Assignment         Assignment    `json:"assignment"`
	Artifacts          []ArtifactPin `json:"artifacts"`
	ReturnCriteria     []string      `json:"returnCriteria"`
}

type RetirePreparationRequest PreparationChange

func (q RevisePreparationRequest) change() PreparationChange {
	return PreparationChange{Schema: q.Schema, ProjectRef: q.ProjectRef, CheckoutID: q.CheckoutID,
		OperationID: q.OperationID, ExpectedGeneration: q.ExpectedGeneration, AssignmentID: q.AssignmentID,
		ExpectedRevision: q.ExpectedRevision, Reason: q.Reason}
}

func validPreparationReason(reason string) bool {
	return strings.TrimSpace(reason) != "" && len(reason) <= 2048
}

func hasPreparedClaim(state snapshot) bool {
	for _, p := range state.Prepared {
		if p.State != "retired" {
			return true
		}
	}
	return false
}

func (s *Service) ReviseAssignment(caller Caller, q RevisePreparationRequest) (PreparationResult, error) {
	if q.Assignment.ExcludedPaths == nil {
		q.Assignment.ExcludedPaths = []string{}
	}
	return s.changePreparation(caller, q.change(), q, &q)
}

// RetireAssignment releases only an unlaunched preparation claim, never a run.
func (s *Service) RetireAssignment(caller Caller, q RetirePreparationRequest) (PreparationResult, error) {
	return s.changePreparation(caller, PreparationChange(q), q, nil)
}

func (s *Service) changePreparation(caller Caller, q PreparationChange, request any, revise *RevisePreparationRequest) (PreparationResult, error) {
	var out PreparationResult
	if !caller.controller {
		return out, fail("unauthorized", "assignment changes require trusted controller")
	}
	if q.Schema != PreparationLifecycleSchema {
		return out, fail("incompatible_profile", "unsupported assignment lifecycle schema")
	}
	if !uuidPattern.MatchString(q.AssignmentID) || !digestPattern.MatchString(q.ExpectedGeneration) || q.ExpectedRevision == 0 || !validPreparationReason(q.Reason) {
		return out, fail("invalid_request", "assignment UUID, generation, positive revision and bounded reason required")
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
	// Include operation type: absent replacement fields must never alias retirement.
	action := "retire"
	if revise != nil {
		action = "revise"
	}
	b, err := json.Marshal(struct {
		Action  string `json:"action"`
		Request any    `json:"request"`
	}{action, request})
	if err != nil {
		return out, err
	}
	digest := hash(b)
	if old, ok := state.PreparationReceipts[q.OperationID]; ok {
		if old.Digest != digest {
			return out, fail("idempotency_conflict", "operationId has different assignment change content")
		}
		return old.Result, nil
	}
	if _, ok := state.DesignDecisionReceipts[q.OperationID]; ok {
		return out, fail("idempotency_conflict", "operationId belongs to design decision")
	}
	if _, ok := state.Receipts[q.OperationID]; ok {
		return out, fail("idempotency_conflict", "operationId belongs to a legacy operation")
	}
	if head != q.ExpectedGeneration {
		return out, fail("stale_revision", "operational generation changed")
	}
	p, ok := state.Prepared[q.AssignmentID]
	if !ok {
		return out, fail("missing_coverage", "prepared assignment not found")
	}
	if p.Revision != q.ExpectedRevision {
		return out, fail("stale_revision", "assignment revision changed")
	}
	if _, ok := state.AssignmentRuns[p.ID]; ok {
		return out, fail("invalid_transition", "activated preparation is immutable; use run interruption/replanning")
	}
	if p.State == "retired" {
		return out, fail("invalid_transition", "retired preparation is terminal; prepare a new assignment")
	}
	if len(state.PreparationReceipts) >= 1024 || p.Revision == ^uint64(0) {
		return out, fail("limit", "assignment history capacity reached; preserve receipts")
	}
	if revise != nil {
		candidate := PrepareRequest{Schema: PreparationSchema, ExpectedGeneration: q.ExpectedGeneration,
			InputRevision: revise.InputRevision, Role: p.Role, Selection: p.Selection,
			Assignment: revise.Assignment, Artifacts: revise.Artifacts, ReturnCriteria: revise.ReturnCriteria}
		if err = validatePreparation(candidate); err != nil {
			return out, err
		}
		if err = s.checkPreparation(p.Role, p.Selection, revise.InputRevision, revise.Artifacts, revise.Assignment); err != nil {
			return out, err
		}
		if why := scopeConflictExcept(revise.Assignment.WritablePaths, state, p.ID); why != "" {
			return out, fail("conflict", why)
		}
		p.InputRevision, p.Assignment = revise.InputRevision, revise.Assignment
		p.Artifacts, p.ReturnCriteria = revise.Artifacts, revise.ReturnCriteria
		p.State = "prepared"
	} else {
		p.State = "retired"
	}
	p.Revision++
	p.ChangeReason = q.Reason
	out = PreparationResult{Schema: PreparationLifecycleSchema, OperationID: q.OperationID,
		Assignment: p, Dispatch: "unavailable", DesignAcceptance: "not-assessed"}
	if state.Schema != designDecisionStateSchema && state.Schema != activationStateSchema && !attentionProfile(state.Schema) {
		state.Schema = preparationLifecycleStateSchema
	}
	state.Prepared[p.ID] = p
	state.PreparationReceipts[q.OperationID] = preparationReceipt{Digest: digest, Result: out}
	if revise != nil {
		if err = s.checkPreparation(p.Role, p.Selection, p.InputRevision, p.Artifacts, p.Assignment); err != nil {
			return PreparationResult{}, err
		}
	}
	if err = s.save(state, head); err != nil {
		return PreparationResult{}, err
	}
	return out, nil
}
