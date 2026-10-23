package governance

import "time"

const ReportAcknowledgementSchema = "sdp-governance-report-acknowledgement/1"
const reportAcknowledgementStateSchema = "sdp-project-governance-state/0.8"

// ReportAcknowledgementRequest pins the report and the inspected current state.
// Coordinator is local controller attribution, not authenticated human identity.
type ReportAcknowledgementRequest struct {
	Schema             string `json:"schema"`
	ProjectRef         string `json:"projectRef"`
	CheckoutID         string `json:"checkoutId"`
	OperationID        string `json:"operationId"`
	RunID              string `json:"runId"`
	ReportID           string `json:"reportId"`
	ReportedRevision   uint64 `json:"reportedRevision"`
	InputRevision      string `json:"inputRevision"`
	ExpectedRevision   uint64 `json:"expectedRevision"`
	ExpectedGeneration string `json:"expectedGeneration"`
	Coordinator        string `json:"coordinator"`
	Note               string `json:"note"`
}

type ReportAcknowledgement struct {
	Request    ReportAcknowledgementRequest `json:"request"`
	RecordedAt string                       `json:"recordedAt"`
}

func validReportAcknowledgement(q ReportAcknowledgementRequest) bool {
	return q.Schema == ReportAcknowledgementSchema && uuidPattern.MatchString(q.OperationID) &&
		uuidPattern.MatchString(q.RunID) && uuidPattern.MatchString(q.ReportID) && uuidPattern.MatchString(q.CheckoutID) &&
		digestPattern.MatchString(q.InputRevision) && digestPattern.MatchString(q.ExpectedGeneration) &&
		q.ReportedRevision > 0 && q.ExpectedRevision >= q.ReportedRevision &&
		boundedAttentionText(q.Coordinator, 256) && boundedAttentionText(q.Note, 4096)
}

func reportMatches(state snapshot, q ReportAcknowledgementRequest) bool {
	v, ok := state.Receipts[q.ReportID]
	r := v.Result
	return ok && r.Operation == "transition" && r.OperationID == q.ReportID && r.Run.ID == q.RunID &&
		r.Run.State == "verifying" && r.Run.Revision == q.ReportedRevision && r.Run.InputRevision == q.InputRevision
}

// AcknowledgeReport records processing only. It never advances a run, verifies
// evidence, accepts an assignment or launches a host turn.
func (s *Service) AcknowledgeReport(caller Caller, q ReportAcknowledgementRequest) (Result, error) {
	if !caller.controller {
		return Result{}, fail("unauthorized", "report acknowledgement requires the trusted controller")
	}
	if q.Schema != ReportAcknowledgementSchema {
		return Result{}, fail("incompatible_profile", "unsupported report acknowledgement schema")
	}
	if !validReportAcknowledgement(q) {
		return Result{}, fail("invalid_request", "exact report and inspected basis with bounded attribution/note required")
	}
	unlock, e := lock(s.state)
	if e != nil {
		return Result{}, e
	}
	defer unlock()
	state, head, e := s.load()
	if e != nil {
		return Result{}, e
	}
	if e = s.authorize(caller, Schema, q.ProjectRef, q.CheckoutID, q.OperationID, state); e != nil {
		return Result{}, e
	}
	digest := requestDigest(q)
	if old, ok, e := retry(state, q.OperationID, digest); ok {
		return old, e
	}
	if head != q.ExpectedGeneration {
		return Result{}, fail("stale_revision", "report inbox generation changed")
	}
	if !reportMatches(state, q) {
		return Result{}, fail("missing_coverage", "exact report receipt not found")
	}
	for _, receipt := range state.Receipts {
		if a := receipt.Result.ReportAcknowledgement; a != nil && a.Request.ReportID == q.ReportID {
			return Result{}, fail("conflict", "report already acknowledged")
		}
	}
	r, ok := state.Runs[q.RunID]
	if !ok {
		return Result{}, fail("missing_coverage", "report run absent")
	}
	if r.Revision != q.ExpectedRevision || r.InputRevision != q.InputRevision {
		return Result{}, fail("stale_revision", "report run revision changed")
	}
	if r.State != "verifying" {
		return Result{}, fail("conflict", "report is no longer awaiting verification")
	}
	if attentionBasis(summarizeRun(r)) != attentionBasis(state.Receipts[q.ReportID].Result.Run) {
		return Result{}, fail("stale_revision", "report decision context changed")
	}
	if e = s.currentAttentionInputs(r); e != nil {
		return Result{}, e
	}
	if len(state.Receipts) >= 1024 {
		return Result{}, fail("limit", "operational receipt capacity reached")
	}
	out := Result{Schema: Schema, Operation: "report.acknowledge", OperationID: q.OperationID, ProjectRef: s.project,
		CheckoutID: state.CheckoutID, Run: summarizeRun(r), Coverage: "known",
		ReportAcknowledgement: &ReportAcknowledgement{Request: q, RecordedAt: time.Now().UTC().Format(time.RFC3339Nano)}}
	state.Schema = reportAcknowledgementStateSchema
	state.Receipts[q.OperationID] = receipt{Digest: digest, Result: out}
	// As with other operations, source coverage is checked again before atomic
	// publication, but external filesystem edits are not locked by the run store.
	if e = s.currentAttentionInputs(r); e != nil {
		return Result{}, e
	}
	if e = s.save(state, head); e != nil {
		return Result{}, e
	}
	return out, nil
}

func validateReportAcknowledgements(state snapshot) error {
	seen := map[string]bool{}
	for id, receipt := range state.Receipts {
		v := receipt.Result
		a := v.ReportAcknowledgement
		if a == nil {
			if v.Operation == "report.acknowledge" {
				return fail("reconciliation_required", "acknowledgement receipt missing its record")
			}
			continue
		}
		q := a.Request
		_, timeErr := time.Parse(time.RFC3339Nano, a.RecordedAt)
		r, exists := state.Runs[q.RunID]
		if state.Schema != reportAcknowledgementStateSchema || !validReportAcknowledgement(q) || timeErr != nil ||
			q.ProjectRef != state.Project || q.CheckoutID != state.CheckoutID || id != q.OperationID || v.OperationID != id ||
			v.Schema != Schema || v.Operation != "report.acknowledge" || v.ProjectRef != state.Project || v.CheckoutID != state.CheckoutID ||
			v.Coverage != "known" || receipt.Digest != requestDigest(q) || !reportMatches(state, q) || seen[q.ReportID] ||
			!exists || r.Revision < q.ExpectedRevision || v.Run.ID != q.RunID || v.Run.Revision != q.ExpectedRevision ||
			v.Run.State != "verifying" || v.Run.InputRevision != q.InputRevision {
			return fail("reconciliation_required", "invalid retained report acknowledgement")
		}
		seen[q.ReportID] = true
	}
	if state.Schema == reportAcknowledgementStateSchema && len(seen) == 0 {
		return fail("reconciliation_required", "acknowledgement profile has no records")
	}
	return nil
}
