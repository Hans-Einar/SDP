package governance

import (
	"encoding/json"
	"reflect"
	"slices"
	"strings"
	"time"
)

const DesignDecisionSchema = "sdp-governance-design-decision/1"
const designDecisionStateSchema = "sdp-project-governance-state/0.4"

// The request names the reviewed package, never an owner identity or decision.
type DesignDecisionRequest struct {
	Schema         string             `json:"schema"`
	ProjectRef     string             `json:"projectRef"`
	CheckoutID     string             `json:"checkoutId"`
	OperationID    string             `json:"operationId"`
	Review         DesignCheckRequest `json:"review"`
	ReviewRevision string             `json:"reviewRevision"`
}

// OwnerDesignInput is supplied separately by a trusted controller adapter.
// It is an attributed observation, not proof of human identity. Same-user shell
// access is outside Caller authority, just as for existing local controller APIs.
type OwnerDesignInput struct {
	Actor    string `json:"actor"`
	Source   string `json:"source"`
	Decision string `json:"decision"`
	Comment  string `json:"comment"`
}

type DesignDecisionRecord struct {
	Input      OwnerDesignInput  `json:"input"`
	RecordedAt string            `json:"recordedAt"`
	Basis      DesignCheckResult `json:"basis"`
}

type DesignDecisionResult struct {
	Schema      string               `json:"schema"`
	OperationID string               `json:"operationId"`
	Record      DesignDecisionRecord `json:"record"`
	Dispatch    string               `json:"dispatch"`
}

type DesignReviewResult struct {
	Schema     string                `json:"schema"`
	Check      DesignCheckResult     `json:"check"`
	Decision   *DesignDecisionResult `json:"decision,omitempty"`
	Acceptance string                `json:"acceptance"`
	Dispatch   string                `json:"dispatch"`
}

type designDecisionReceipt struct {
	RequestDigest string               `json:"requestDigest"`
	InputDigest   string               `json:"inputDigest"`
	Result        DesignDecisionResult `json:"result"`
}

func decisionDigest(v any) string { b, _ := json.Marshal(v); return hash(b) }
func validOwnerDesignInput(input OwnerDesignInput) bool {
	return strings.TrimSpace(input.Actor) != "" && len(input.Actor) <= 256 &&
		strings.TrimSpace(input.Source) != "" && len(input.Source) <= 1024 &&
		validPreparationReason(input.Comment) && slices.Contains([]string{"approved", "changes-requested", "rejected"}, input.Decision)
}

func (s *Service) decisionRequest(caller Caller, q DesignDecisionRequest, state snapshot) error {
	if !caller.controller {
		return fail("unauthorized", "design decisions require trusted controller owner input")
	}
	if q.Schema != DesignDecisionSchema {
		return fail("incompatible_profile", "unsupported design decision schema")
	}
	if !digestPattern.MatchString(q.ReviewRevision) {
		return fail("invalid_request", "exact reviewed revision required")
	}
	return s.authorize(caller, Schema, q.ProjectRef, q.CheckoutID, q.OperationID, state)
}

// DesignDecisionReceipt returns only a historical result. It never asserts that
// the recorded decision still applies; use ReviewAssignmentDesign for that.
func (s *Service) DesignDecisionReceipt(caller Caller, q DesignDecisionRequest) (*DesignDecisionResult, error) {
	if !caller.controller {
		return nil, fail("unauthorized", "trusted controller required")
	}
	state, _, err := s.load()
	if err != nil {
		return nil, err
	}
	if err = s.decisionRequest(caller, q, state); err != nil {
		return nil, err
	}
	return findDesignDecision(state, q)
}
func findDesignDecision(state snapshot, q DesignDecisionRequest) (*DesignDecisionResult, error) {
	if _, ok := state.Receipts[q.OperationID]; ok {
		return nil, fail("idempotency_conflict", "operationId belongs to legacy operation")
	}
	if _, ok := state.PreparationReceipts[q.OperationID]; ok {
		return nil, fail("idempotency_conflict", "operationId belongs to preparation")
	}
	if old, ok := state.DesignDecisionReceipts[q.OperationID]; ok {
		if old.RequestDigest != decisionDigest(q) {
			return nil, fail("idempotency_conflict", "operationId has different review request")
		}
		return &old.Result, nil
	}
	return nil, nil
}

func (s *Service) RecordDesignDecision(caller Caller, q DesignDecisionRequest, input OwnerDesignInput) (DesignDecisionResult, error) {
	return s.recordDesignDecision(caller, q, input, nil)
}
func (s *Service) recordDesignDecision(caller Caller, q DesignDecisionRequest, input OwnerDesignInput, beforeRecheck func()) (DesignDecisionResult, error) {
	var out DesignDecisionResult
	if !caller.controller {
		return out, fail("unauthorized", "design decisions require trusted controller owner input")
	}
	if !validOwnerDesignInput(input) {
		return out, fail("invalid_request", "explicit decision, bounded actor/source and comment required")
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
	if err = s.decisionRequest(caller, q, state); err != nil {
		return out, err
	}
	old, err := findDesignDecision(state, q)
	if err != nil {
		return out, err
	}
	if old != nil {
		if state.DesignDecisionReceipts[q.OperationID].InputDigest != decisionDigest(input) {
			return out, fail("idempotency_conflict", "operationId has different owner input")
		}
		return *old, nil
	}
	if head != q.Review.ExpectedGeneration {
		return out, fail("stale_revision", "operational generation changed")
	}
	if len(state.DesignDecisionReceipts) >= 1024 {
		return out, fail("limit", "design decision history capacity reached; preserve receipts")
	}
	check, err := s.CheckAssignmentDesign(q.Review)
	if err != nil {
		return out, err
	}
	if check.Status != "current" || check.ReviewRevision != q.ReviewRevision {
		return out, fail("stale_revision", "reviewed design is no longer current")
	}
	if beforeRecheck != nil {
		beforeRecheck()
	}
	final, err := s.CheckAssignmentDesign(q.Review)
	if err != nil {
		return out, err
	}
	if !reflect.DeepEqual(check, final) {
		return out, fail("stale_revision", "design changed before decision publication")
	}
	out = DesignDecisionResult{Schema: DesignDecisionSchema, OperationID: q.OperationID,
		Record: DesignDecisionRecord{Input: input, RecordedAt: time.Now().UTC().Format(time.RFC3339Nano), Basis: final}, Dispatch: "unavailable"}
	if state.DesignDecisions == nil {
		state.DesignDecisions = map[string]string{}
	}
	if state.DesignDecisionReceipts == nil {
		state.DesignDecisionReceipts = map[string]designDecisionReceipt{}
	}
	if state.Schema != activationStateSchema {
		state.Schema = designDecisionStateSchema
	}
	state.DesignDecisions[q.Review.AssignmentID] = q.OperationID
	state.DesignDecisionReceipts[q.OperationID] = designDecisionReceipt{RequestDigest: decisionDigest(q), InputDigest: decisionDigest(input), Result: out}
	if err = s.save(state, head); err != nil {
		return DesignDecisionResult{}, err
	}
	return out, nil
}

// ReviewAssignmentDesign projects applicability without rewriting decision history.
func (s *Service) ReviewAssignmentDesign(q DesignCheckRequest) (DesignReviewResult, error) {
	out := DesignReviewResult{Schema: DesignDecisionSchema, Acceptance: "not-recorded", Dispatch: "unavailable"}
	check, err := s.CheckAssignmentDesign(q)
	if err != nil {
		return out, err
	}
	out.Check = check
	state, head, err := s.load()
	if err != nil {
		return out, err
	}
	if head != check.Generation {
		return out, fail("stale_revision", "generation changed during decision inspection")
	}
	if id, ok := state.DesignDecisions[q.AssignmentID]; ok {
		record := state.DesignDecisionReceipts[id].Result
		out.Decision = &record
		out.Acceptance = "stale"
		if check.Status == "current" && record.Record.Basis.ReviewRevision == check.ReviewRevision {
			out.Acceptance = record.Record.Input.Decision
		}
	}
	return out, nil
}

func validateDesignDecisionState(state snapshot) error {
	if state.Schema != designDecisionStateSchema && state.Schema != activationStateSchema {
		if state.DesignDecisions != nil || state.DesignDecisionReceipts != nil {
			return fail("incompatible_profile", "design decisions require profile 0.4")
		}
		return nil
	}
	if state.DesignDecisions == nil || state.DesignDecisionReceipts == nil || len(state.DesignDecisionReceipts) == 0 || len(state.DesignDecisionReceipts) > 1024 {
		return fail("reconciliation_required", "invalid design decision collections")
	}
	for id, r := range state.DesignDecisionReceipts {
		v := r.Result
		b := v.Record.Basis
		_, timeErr := time.Parse(time.RFC3339Nano, v.Record.RecordedAt)
		_, prepared := state.Prepared[b.Assignment.ID]
		if !uuidPattern.MatchString(id) || v.Schema != DesignDecisionSchema || v.OperationID != id || v.Dispatch != "unavailable" ||
			!digestPattern.MatchString(r.RequestDigest) || r.InputDigest != decisionDigest(v.Record.Input) || !validOwnerDesignInput(v.Record.Input) || timeErr != nil ||
			b.Schema != DesignCheckSchema || b.Status != "current" || !digestPattern.MatchString(b.ReviewRevision) || b.Provider == nil || !prepared {
			return fail("reconciliation_required", "invalid retained design decision")
		}
		if _, ok := state.Receipts[id]; ok {
			return fail("reconciliation_required", "decision operation identity collision")
		}
		if _, ok := state.PreparationReceipts[id]; ok {
			return fail("reconciliation_required", "decision operation identity collision")
		}
	}
	for assignment, id := range state.DesignDecisions {
		r, ok := state.DesignDecisionReceipts[id]
		if !ok || r.Result.Record.Basis.Assignment.ID != assignment {
			return fail("reconciliation_required", "invalid latest design decision reference")
		}
	}
	return nil
}
