package governance

import (
	"strings"
)

func (s *Service) authorize(c Caller, schema, project, checkout, operation string, state snapshot) error {
	if !c.controller && c.runID == "" {
		return fail("unauthorized", "trusted controller invocation required")
	}
	if schema != Schema {
		return fail("incompatible_profile", "unsupported application schema")
	}
	if project != s.project || checkout == "" || checkout != state.CheckoutID {
		return fail("unauthorized", "checkout binding does not match")
	}
	if !uuidPattern.MatchString(operation) {
		return fail("invalid_request", "operationId must be a version 4 UUID")
	}
	return nil
}
func retry(state snapshot, id, digest string) (Result, bool, error) {
	if _, ok := state.LeaderReceipts[id]; ok {
		return Result{}, true, fail("idempotency_conflict", "operationId belongs to leader delivery")
	}
	if _, ok := state.DesignDecisionReceipts[id]; ok {
		return Result{}, true, fail("idempotency_conflict", "operationId belongs to design decision")
	}
	if _, ok := state.PreparationReceipts[id]; ok {
		return Result{}, true, fail("idempotency_conflict", "operationId belongs to assignment preparation")
	}
	old, ok := state.Receipts[id]
	if !ok {
		return Result{}, false, nil
	}
	if old.Digest != digest {
		return Result{}, true, fail("idempotency_conflict", "operationId has different content")
	}
	return old.Result, true, nil
}
func (s *Service) Start(caller Caller, q StartRequest) (Result, error) {
	var out Result
	if q.Assignment.ExcludedPaths == nil {
		q.Assignment.ExcludedPaths = []string{}
	}
	// Refuse untrusted callers before even acquiring a write lock/creating a directory.
	if !caller.controller {
		return out, fail("unauthorized", "controller invocation required")
	}
	if e := validateAssignment(q.Assignment); e != nil {
		return out, e
	}
	if !uuidPattern.MatchString(q.RequestID) || !digestPattern.MatchString(q.InputRevision) {
		return out, fail("invalid_request", "invalid requestId/inputRevision")
	}
	unlock, e := lock(s.state)
	if e != nil {
		return out, e
	}
	defer unlock()
	state, head, e := s.load()
	if e != nil {
		return out, e
	}
	if e = s.authorize(caller, q.Schema, q.ProjectRef, q.CheckoutID, q.OperationID, state); e != nil {
		return out, e
	}
	digest := requestDigest(q)
	if r, ok, e := retry(state, q.OperationID, digest); ok {
		return r, e
	}
	if len(state.Receipts) >= 1024 || len(state.Runs) >= 128 {
		return out, fail("limit", "pilot state capacity reached; retain historical receipts")
	}
	routine := Maintenance()
	if q.RoutineDigest != routine.Digest() {
		return out, fail("incompatible_profile", "routine definition does not match installed version")
	}
	c, e := s.Context(q.Selection)
	if e != nil {
		return out, e
	}
	if q.InputRevision != c.Revision {
		return out, fail("stale_revision", "work inputs changed")
	}
	if len(c.Conditions) > 0 {
		return out, fail("missing_coverage", strings.Join(c.Conditions, "; "))
	}
	for _, r := range state.Runs {
		if r.RequestID == q.RequestID {
			return out, fail("conflict", "request already has a run; retry the original operation or inspect it")
		}
	}
	if hasPreparedClaim(state) {
		for _, p := range q.Assignment.WritablePaths {
			if strings.ContainsAny(p, "*?[]") {
				return out, fail("conflict", "nonliteral legacy scope requires reconciliation with prepared work")
			}
		}
	}
	preparedOnly := snapshot{Prepared: state.Prepared}
	if why := scopeConflict(q.Assignment.WritablePaths, preparedOnly); why != "" {
		return out, fail("conflict", why)
	}
	id, e := identifier()
	if e != nil {
		return out, e
	}
	r := Run{ID: id, Revision: 1, RequestID: q.RequestID, State: "ready", Step: "implement", InputRevision: c.Revision, Selection: q.Selection, Assignment: q.Assignment, Routine: routine, RoutineDigest: routine.Digest(), Reason: "coverage resolved; trusted local controller selected bounded work"}
	out = Result{Schema: Schema, Operation: "start", ProjectRef: s.project, CheckoutID: state.CheckoutID, OperationID: q.OperationID, Run: summarizeRun(r), Coverage: "known"}
	state.Runs[id] = r
	state.Receipts[q.OperationID] = receipt{digest, out}
	// Recheck source revision immediately before publishing the new state.
	latest, e := s.Context(q.Selection)
	if e != nil {
		return Result{}, e
	}
	if latest.Revision != c.Revision {
		return Result{}, fail("stale_revision", "inputs changed during preparation")
	}
	if e = s.save(state, head); e != nil {
		return Result{}, e
	}
	return out, nil
}
func validateAssignment(a Assignment) error {
	if strings.TrimSpace(a.Outcome) == "" || len(a.Outcome) > 64<<10 || len(a.WritablePaths) == 0 || len(a.WritablePaths)+len(a.ExcludedPaths) > 64 {
		return fail("invalid_request", "bounded outcome and 1..64 scope references required")
	}
	paths := map[string]bool{}
	for _, p := range append(append([]string{}, a.WritablePaths...), a.ExcludedPaths...) {
		if !relative(p) || paths[p] {
			return fail("path", "invalid or duplicate assignment path")
		}
		paths[p] = true
	}
	return nil
}
func (s *Service) Read(runID string) (Result, error) {
	var out Result
	if !uuidPattern.MatchString(runID) {
		return out, fail("invalid_request", "invalid runId")
	}
	state, _, e := s.load()
	if e != nil {
		return out, e
	}
	r, ok := state.Runs[runID]
	if !ok {
		return out, fail("missing_coverage", "run not found")
	}
	c, e := s.runContext(r)
	if e != nil {
		return out, e
	}
	coverage := "known"
	if c.Revision != r.InputRevision || len(c.Conditions) > 0 {
		coverage = "stale"
	}
	return Result{Schema: Schema, Operation: "read", ProjectRef: s.project, CheckoutID: state.CheckoutID, Run: summarizeRun(r), Coverage: coverage}, nil
}
func (s *Service) Transition(caller Caller, q TransitionRequest) (Result, error) {
	var out Result
	if !caller.controller && (caller.runID == "" || caller.runID != q.RunID ||
		(q.Action != "report" && q.Action != "fail" && q.Action != "scope-change")) {
		return out, fail("unauthorized", "action requires controller or matching assignment capability")
	}
	if !uuidPattern.MatchString(q.RunID) || !digestPattern.MatchString(q.InputRevision) || q.ExpectedRevision == 0 || strings.TrimSpace(q.Reason) == "" || len(q.Reason) > 64<<10 {
		return out, fail("invalid_request", "run, revision and bounded reason required")
	}
	unlock, e := lock(s.state)
	if e != nil {
		return out, e
	}
	defer unlock()
	state, head, e := s.load()
	if e != nil {
		return out, e
	}
	if e = s.authorize(caller, q.Schema, q.ProjectRef, q.CheckoutID, q.OperationID, state); e != nil {
		return out, e
	}
	if !caller.controller {
		bound := false
		r, ok := state.Runs[q.RunID]
		if ok && len(r.Attempts) > 0 {
			at := r.Attempts[len(r.Attempts)-1]
			bound = at.ID == caller.attemptID && at.CapabilityEpoch() == caller.epoch && at.InputRevision == r.InputRevision && at.Status != "unknown"
		}
		if !bound {
			return out, fail("unauthorized", "agent attempt is absent, superseded or disconnected")
		}
	}
	digest := requestDigest(q)
	if r, ok, e := retry(state, q.OperationID, digest); ok {
		return r, e
	}
	if len(state.Receipts) >= 1024 {
		return out, fail("limit", "receipt capacity reached")
	}
	r, ok := state.Runs[q.RunID]
	if !ok {
		return out, fail("missing_coverage", "run not found")
	}
	if q.ExpectedRevision != r.Revision && !approvedRevision(state, r, caller, q) && !(q.Action == "report" && observedAgentRevision(state, r, caller, q.ExpectedRevision)) {
		return out, &Error{Code: "stale_revision", Message: "run revision changed", CurrentRevision: r.Revision}
	}
	c, e := s.runContext(r)
	if e != nil {
		return out, e
	}
	// Cancellation/interruption/scope notices must remain possible after authority
	// input drift; advancing work requires the pinned current coverage.
	stopping := q.Action == "interrupt" || q.Action == "cancel" || q.Action == "scope-change" || q.Action == "fail"
	if q.InputRevision != r.InputRevision {
		return out, fail("stale_revision", "request does not identify run inputs")
	}
	if !stopping && (c.Revision != r.InputRevision || len(c.Conditions) > 0) {
		return out, fail("stale_revision", "authority/work changed; reconcile before advancing")
	}
	if (q.Action == "report" || q.Action == "start" || q.Action == "resume") && hasBlockingAttention(r) {
		return out, fail("attention_required", "answer or withdraw blocking questions before advancing this run")
	}
	next := ""
	switch q.Action {
	case "start":
		if r.State == "ready" {
			next = "running"
		}
	case "interrupt":
		if r.State == "running" {
			next = "interrupted"
		}
	case "resume":
		if r.State == "interrupted" {
			next = "ready"
		}
	case "report":
		if r.State == "running" {
			next = "verifying"
		}
	case "fail":
		if r.State == "running" || r.State == "verifying" {
			next = "failed"
		}
	case "retry":
		if r.State == "failed" {
			next = "ready"
		}
	case "scope-change":
		if r.State != "completed" && r.State != "canceled" {
			next = "replan-required"
		}
	case "cancel":
		if r.State != "completed" && r.State != "canceled" {
			next = "canceled"
		}
	case "complete":
		return out, fail("evidence_missing", "completion requires independently verified review and checked publication; attached evidence alone cannot complete work")
	default:
		return out, fail("invalid_request", "unsupported routine action")
	}
	if next == "" {
		return out, fail("conflict", "action is not permitted from "+r.State)
	}
	if r.Revision == ^uint64(0) {
		return out, fail("limit", "run revision exhausted")
	}
	r.State = next
	r.Revision++
	r.Reason = q.Reason
	out = Result{Schema: Schema, Operation: "transition", ProjectRef: s.project, CheckoutID: state.CheckoutID, OperationID: q.OperationID, Run: summarizeRun(r), Coverage: "known"}
	if c.Revision != r.InputRevision {
		out.Coverage = "stale"
	}
	state.Runs[r.ID] = r
	state.Receipts[q.OperationID] = receipt{digest, out}
	latest, e := s.runContext(r)
	if e != nil {
		return Result{}, e
	}
	if !stopping && (latest.Revision != c.Revision || len(latest.Conditions) > 0) {
		return Result{}, fail("stale_revision", "inputs changed during preparation")
	}
	if e = s.save(state, head); e != nil {
		return Result{}, e
	}
	return out, nil
}
