package governance

import (
	"encoding/json"
	"path/filepath"
	"sort"
	"strings"
)

const PreparationSchema = "sdp-governance-assignment/1"
const preparationStateSchema = "sdp-project-governance-state/0.2"

// ArtifactPin references provider-owned records; it does not certify their meaning.
type ArtifactPin struct {
	Kind   string `json:"kind"`
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}
type PreparationContextRequest struct {
	Role      string    `json:"role"`
	Selection Selection `json:"selection"`
}
type PreparationContext struct {
	Schema     string  `json:"schema"`
	Generation string  `json:"generation"`
	Role       string  `json:"role"`
	Context    Context `json:"context"`
}
type PrepareRequest struct {
	Schema             string        `json:"schema"`
	ProjectRef         string        `json:"projectRef"`
	CheckoutID         string        `json:"checkoutId"`
	OperationID        string        `json:"operationId"`
	ExpectedGeneration string        `json:"expectedGeneration"`
	InputRevision      string        `json:"inputRevision"`
	Role               string        `json:"role"`
	Selection          Selection     `json:"selection"`
	Assignment         Assignment    `json:"assignment"`
	Artifacts          []ArtifactPin `json:"artifacts"`
	ReturnCriteria     []string      `json:"returnCriteria"`
}
type PreparedAssignment struct {
	ID             string        `json:"id"`
	Revision       uint64        `json:"revision"`
	State          string        `json:"state,omitempty"`
	ChangeReason   string        `json:"changeReason,omitempty"`
	Role           string        `json:"role"`
	Selection      Selection     `json:"selection"`
	Assignment     Assignment    `json:"assignment"`
	InputRevision  string        `json:"inputRevision"`
	Artifacts      []ArtifactPin `json:"artifacts"`
	ReturnCriteria []string      `json:"returnCriteria"`
}
type PreparationResult struct {
	Schema           string             `json:"schema"`
	OperationID      string             `json:"operationId"`
	Assignment       PreparedAssignment `json:"assignment"`
	Dispatch         string             `json:"dispatch"`
	DesignAcceptance string             `json:"designAcceptance"`
}
type PreparationView struct {
	Assignment       PreparedAssignment `json:"assignment"`
	Readiness        string             `json:"readiness"`
	Conditions       []string           `json:"conditions"`
	Dispatch         string             `json:"dispatch"`
	DesignAcceptance string             `json:"designAcceptance"`
}
type Preparations struct {
	Schema      string            `json:"schema"`
	Generation  string            `json:"generation"`
	CheckoutID  string            `json:"checkoutId"`
	Scope       string            `json:"scope"`
	Assignments []PreparationView `json:"assignments"`
}
type preparationReceipt struct {
	Digest string            `json:"digest"`
	Result PreparationResult `json:"result"`
}

func rolePlans(role string) ([]string, error) {
	switch role {
	case "SAD":
		return []string{"DesignPlan", "ArchitecturePlan"}, nil
	case "Master":
		return []string{"ImplementationPlan"}, nil
	default:
		return nil, fail("invalid_request", "preparation role must be SAD or Master")
	}
}
func (s *Service) PreparationContext(q PreparationContextRequest) (PreparationContext, error) {
	out := PreparationContext{Schema: PreparationSchema, Role: q.Role}
	plans, err := rolePlans(q.Role)
	if err != nil {
		return out, err
	}
	_, head, err := s.load()
	if err != nil {
		return out, err
	}
	out.Generation = head
	out.Context, err = s.contextFor(q.Selection, plans, "PLAN-")
	return out, err
}
func validatePreparation(q PrepareRequest) error {
	if q.Schema != PreparationSchema {
		return fail("incompatible_profile", "unsupported assignment preparation schema")
	}
	if _, err := rolePlans(q.Role); err != nil {
		return err
	}
	if !digestPattern.MatchString(q.ExpectedGeneration) || !digestPattern.MatchString(q.InputRevision) {
		return fail("invalid_request", "pinned generation and input revision required")
	}
	if err := validateAssignment(q.Assignment); err != nil {
		return err
	}
	for _, p := range append(append([]string{}, q.Assignment.WritablePaths...), q.Assignment.ExcludedPaths...) {
		if strings.ContainsAny(p, "*?[]") {
			return fail("path", "scope paths must be literal, not globs")
		}
	}
	if len(q.ReturnCriteria) == 0 || len(q.ReturnCriteria) > 16 {
		return fail("invalid_request", "one to sixteen return criteria required")
	}
	for _, v := range q.ReturnCriteria {
		if strings.TrimSpace(v) == "" || len(v) > 2048 {
			return fail("invalid_request", "empty or excessive return criterion")
		}
	}
	required := map[string]bool{"now": true}
	if q.Role == "Master" {
		for _, k := range []string{"target", "task", "blueprint", "design-disposition"} {
			required[k] = true
		}
	}
	if len(q.Artifacts) != len(required) {
		return fail("invalid_request", "exact role-specific artifact pins required")
	}
	for _, a := range q.Artifacts {
		if !required[a.Kind] || !relative(a.Path) || !digestPattern.MatchString(a.SHA256) {
			return fail("invalid_request", "invalid, repeated or unexpected artifact pin")
		}
		delete(required, a.Kind)
	}
	return nil
}
func overlaps(a, b string) bool {
	return a == b || strings.HasPrefix(a, b+"/") || strings.HasPrefix(b, a+"/")
}
func scopeConflict(scope []string, state snapshot) string {
	return scopeConflictExcept(scope, state, "")
}
func scopeConflictExcept(scope []string, state snapshot, except string) string {
	ids := make([]string, 0, len(state.Prepared))
	for id := range state.Prepared {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		if id == except || state.Prepared[id].State == "retired" {
			continue
		}
		for _, p := range scope {
			for _, other := range state.Prepared[id].Assignment.WritablePaths {
				if overlaps(p, other) {
					return "prepared assignment " + id + " overlaps " + p
				}
			}
		}
	}
	ids = ids[:0]
	for id := range state.Runs {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		r := state.Runs[id]
		if r.State == "completed" || r.State == "canceled" {
			continue
		}
		for _, p := range scope {
			for _, other := range r.Assignment.WritablePaths {
				if strings.ContainsAny(other, "*?[]") {
					return "legacy run " + id + " has nonliteral scope requiring reconciliation"
				}
				if overlaps(p, other) {
					return "legacy run " + id + " overlaps " + p
				}
			}
		}
	}
	return ""
}
func (s *Service) checkPreparation(role string, sel Selection, revision string, pins []ArtifactPin, scope Assignment) error {
	plans, err := rolePlans(role)
	if err != nil {
		return err
	}
	c, err := s.contextFor(sel, plans, "PLAN-")
	if err != nil {
		return err
	}
	if len(c.Conditions) > 0 {
		return fail("missing_coverage", strings.Join(c.Conditions, "; "))
	}
	if c.Revision != revision {
		return fail("stale_revision", "governing sources changed")
	}
	for _, a := range pins {
		b, err := readFile(filepath.Join(s.root, filepath.FromSlash(a.Path)), MaxRequest)
		if err != nil {
			return fail("stale_revision", "artifact unavailable: "+a.Kind+": "+err.Error())
		}
		if hash(b) != a.SHA256 {
			return fail("stale_revision", "artifact changed: "+a.Kind)
		}
	}
	for _, p := range append(append([]string{}, scope.WritablePaths...), scope.ExcludedPaths...) {
		if err := safeAbsolute(filepath.Join(s.root, filepath.FromSlash(p))); err != nil {
			return err
		}
	}
	return nil
}

// PrepareAssignment persists controller-selected work, never approval or dispatch.
// Retry receipts are immutable results, not current readiness assessments.
func (s *Service) PrepareAssignment(caller Caller, q PrepareRequest) (PreparationResult, error) {
	var out PreparationResult
	if q.Assignment.ExcludedPaths == nil {
		q.Assignment.ExcludedPaths = []string{}
	}
	if !caller.controller {
		return out, fail("unauthorized", "assignment preparation requires trusted controller")
	}
	if err := validatePreparation(q); err != nil {
		return out, err
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
	b, _ := json.Marshal(q)
	digest := hash(b)
	if old, ok := state.PreparationReceipts[q.OperationID]; ok {
		if old.Digest != digest {
			return out, fail("idempotency_conflict", "operationId has different preparation content")
		}
		return old.Result, nil
	}
	if _, ok := state.DesignDecisionReceipts[q.OperationID]; ok {
		return out, fail("idempotency_conflict", "operationId belongs to design decision")
	}
	if _, ok := state.Receipts[q.OperationID]; ok {
		return out, fail("idempotency_conflict", "operationId already belongs to a legacy operation")
	}
	if head != q.ExpectedGeneration {
		return out, fail("stale_revision", "operational generation changed")
	}
	if len(state.Prepared) >= 128 || len(state.PreparationReceipts) >= 1024 {
		return out, fail("limit", "assignment preparation capacity reached")
	}
	if err = s.checkPreparation(q.Role, q.Selection, q.InputRevision, q.Artifacts, q.Assignment); err != nil {
		return out, err
	}
	if why := scopeConflict(q.Assignment.WritablePaths, state); why != "" {
		return out, fail("conflict", why)
	}
	id, err := identifier()
	if err != nil {
		return out, err
	}
	prepared := PreparedAssignment{ID: id, Revision: 1, Role: q.Role, Selection: q.Selection, Assignment: q.Assignment, InputRevision: q.InputRevision, Artifacts: q.Artifacts, ReturnCriteria: q.ReturnCriteria}
	out = PreparationResult{Schema: PreparationSchema, OperationID: q.OperationID, Assignment: prepared, Dispatch: "unavailable", DesignAcceptance: "not-assessed"}
	if state.Prepared == nil {
		state.Prepared = map[string]PreparedAssignment{}
	}
	if state.PreparationReceipts == nil {
		state.PreparationReceipts = map[string]preparationReceipt{}
	}
	if state.Schema == Schema {
		state.Schema = preparationStateSchema
	}
	state.Prepared[id] = prepared
	state.PreparationReceipts[q.OperationID] = preparationReceipt{Digest: digest, Result: out}
	if err = s.checkPreparation(q.Role, q.Selection, q.InputRevision, q.Artifacts, q.Assignment); err != nil {
		return PreparationResult{}, err
	}
	if err = s.save(state, head); err != nil {
		return PreparationResult{}, err
	}
	return out, nil
}
func (s *Service) InspectAssignments(id string) (Preparations, error) {
	out := Preparations{Schema: PreparationSchema, Scope: "selected-checkout-only", Assignments: []PreparationView{}}
	if id != "" && !uuidPattern.MatchString(id) {
		return out, fail("invalid_request", "invalid assignment UUID")
	}
	state, head, err := s.load()
	if err != nil {
		return out, err
	}
	out.Generation = head
	if state.Schema == preparationLifecycleStateSchema || state.Schema == designDecisionStateSchema {
		out.Schema = PreparationLifecycleSchema
	}
	out.CheckoutID = state.CheckoutID
	if id != "" {
		if _, ok := state.Prepared[id]; !ok {
			return out, fail("missing_coverage", "prepared assignment not found")
		}
	}
	ids := make([]string, 0, len(state.Prepared))
	for key := range state.Prepared {
		if id == "" || id == key {
			ids = append(ids, key)
		}
	}
	sort.Strings(ids)
	for _, key := range ids {
		p := state.Prepared[key]
		v := PreparationView{Assignment: p, Readiness: "prepared", Conditions: []string{}, Dispatch: "unavailable", DesignAcceptance: "not-assessed"}
		if p.State == "retired" {
			v.Readiness = "retired"
		} else if err = s.checkPreparation(p.Role, p.Selection, p.InputRevision, p.Artifacts, p.Assignment); err != nil {
			v.Readiness = "stale"
			v.Conditions = append(v.Conditions, err.Error())
		}
		out.Assignments = append(out.Assignments, v)
	}
	return out, nil
}
