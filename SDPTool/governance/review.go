package governance

import (
	"encoding/json"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/Hans-Einar/SDP/SDPTool/codexcompat"
)

const ReviewSchema = "sdp-governance-review/1"
const reviewStateSchema = "sdp-project-governance-state/0.11"

func evidenceProfile(schema string) bool {
	return schema == evidenceStateSchema || schema == reviewStateSchema || publicationProfile(schema)
}

type ReviewSpec struct {
	Instructions string       `json:"instructions"`
	Candidate    Candidate    `json:"candidate"`
	ReportID     string       `json:"reportId"`
	EvidenceIDs  []string     `json:"evidenceIds"`
	Focus        string       `json:"focus"`
	HostVersion  string       `json:"hostVersion"`
	MCPBinary    string       `json:"mcpBinary"`
	Skills       []SkillInput `json:"skills"`
}
type ReviewOutcome struct {
	Verdict  string   `json:"verdict"`
	Summary  string   `json:"summary"`
	Findings []string `json:"findings"`
}
type Review struct {
	RunRevision uint64         `json:"runRevision"`
	ID          string         `json:"id"`
	RunID       string         `json:"runId"`
	Revision    uint64         `json:"revision"`
	State       string         `json:"state"`
	Spec        ReviewSpec     `json:"spec"`
	Payload     string         `json:"payload"`
	ThreadID    string         `json:"threadId"`
	TurnID      string         `json:"turnId"`
	Answer      string         `json:"answer"`
	Outcome     *ReviewOutcome `json:"outcome,omitempty"`
	CreatedAt   string         `json:"createdAt"`
	UpdatedAt   string         `json:"updatedAt"`
}
type ReviewRequest struct {
	Schema              string      `json:"schema"`
	ProjectRef          string      `json:"projectRef"`
	CheckoutID          string      `json:"checkoutId"`
	OperationID         string      `json:"operationId"`
	RunID               string      `json:"runId"`
	ReviewID            string      `json:"reviewId"`
	ExpectedRevision    uint64      `json:"expectedRevision"`
	ExpectedRunRevision uint64      `json:"expectedRunRevision"`
	Action              string      `json:"action"`
	Spec                *ReviewSpec `json:"spec,omitempty"`
	ThreadID            string      `json:"threadId"`
	TurnID              string      `json:"turnId"`
	Status              string      `json:"status"`
	Answer              string      `json:"answer"`
}
type ReviewView struct {
	Schema            string   `json:"schema"`
	ProjectRef        string   `json:"projectRef"`
	CheckoutID        string   `json:"checkoutId"`
	Review            Review   `json:"review"`
	Applicable        bool     `json:"applicable"`
	Qualified         bool     `json:"qualified"`
	ContextSeparation string   `json:"contextSeparation"`
	Conditions        []string `json:"conditions"`
}

func ReviewActive(state string) bool {
	return slices.Contains([]string{"prepared", "opening", "bound", "sending", "running"}, state)
}
func validReviewSpec(v ReviewSpec) bool {
	reviewer := false
	for _, s := range v.Skills {
		if s.Name == "sdp-reviewer" {
			reviewer = true
		}
	}
	if !boundedAttentionText(v.Instructions, 8192) || !reviewer || !ValidSkillInputs(v.Skills) || !codexcompat.Supported(v.HostVersion) || !filepath.IsAbs(v.MCPBinary) || len(v.MCPBinary) > 4096 || !boundedAttentionText(v.Focus, 4096) || !uuidPattern.MatchString(v.ReportID) || len(v.EvidenceIDs) < 1 || len(v.EvidenceIDs) > 8 || v.Candidate.Schema != EvidenceSchema || candidateDigest(v.Candidate) != v.Candidate.Digest {
		return false
	}
	seen := map[string]bool{}
	for _, id := range v.EvidenceIDs {
		if !uuidPattern.MatchString(id) || seen[id] {
			return false
		}
		seen[id] = true
	}
	return true
}
func parseReviewOutcome(answer string) *ReviewOutcome {
	var v ReviewOutcome
	if Decode([]byte(answer), &v) != nil || !slices.Contains([]string{"passed", "failed", "inconclusive"}, v.Verdict) || !boundedAttentionText(v.Summary, 8192) || v.Findings == nil || len(v.Findings) > 32 || (v.Verdict == "passed" && len(v.Findings) != 0) || (v.Verdict == "failed" && len(v.Findings) == 0) {
		return nil
	}
	for _, f := range v.Findings {
		if !boundedAttentionText(f, 4096) {
			return nil
		}
	}
	return &v
}
func (s *Service) reviewBasis(state snapshot, r Run, spec ReviewSpec) error {
	if !validReviewSpec(spec) {
		return fail("invalid_request", "bounded candidate, tests, reviewer skill and host configuration required")
	}
	c := spec.Candidate
	if c.ProjectRef != s.project || c.CheckoutID != state.CheckoutID || c.RunID != r.ID || c.InputRevision != r.InputRevision || r.State != "verifying" || latestReport(state, r) != spec.ReportID {
		return fail("stale_revision", "review assignment/report differs")
	}
	if e := s.currentAttentionInputs(r); e != nil {
		return e
	}
	if hasBlockingAttention(r) {
		return fail("attention_required", "blocking questions prevent review")
	}
	for _, a := range r.Attempts {
		if !slices.Contains([]string{"completed", "failed", "interrupted"}, a.Status) || a.Gap {
			return fail("reconciliation_required", "review requires terminal known implementation attempts")
		}
	}
	for _, in := range r.Inputs {
		if in.State != "acknowledged" {
			return fail("reconciliation_required", "unresolved input prevents review")
		}
	}
	actual, e := s.captureCandidate(state.CheckoutID, r, c.Dependencies)
	if e != nil {
		return e
	}
	if actual.Digest != c.Digest {
		return fail("stale_candidate", "review candidate changed")
	}
	if e = CheckSkillSnapshots(s.root, spec.Skills); e != nil {
		return e
	}
	for _, id := range spec.EvidenceIDs {
		v := state.Receipts[id].Result.Evidence
		if v == nil || v.Request.RunID != r.ID || v.Request.Kind != "test" || v.Request.Candidate.Digest != c.Digest || !s.evidenceCurrent(state, r, *v).Applicable {
			return fail("evidence_missing", "selected current candidate test evidence unavailable")
		}
	}
	return nil
}
func (s *Service) ReviewRead(id string) (ReviewView, error) {
	state, head, e := s.load()
	if e != nil {
		return ReviewView{}, e
	}
	d, ok := state.Reviews[id]
	if !ok || !uuidPattern.MatchString(id) {
		return ReviewView{}, fail("missing_coverage", "review absent")
	}
	v := ReviewView{Schema: ReviewSchema, ProjectRef: s.project, CheckoutID: state.CheckoutID, Review: d, Conditions: []string{}, ContextSeparation: "not-observed"}
	if d.ThreadID != "" {
		v.ContextSeparation = "controller-observed-fresh-thread"
	}
	if e = s.reviewBasis(state, state.Runs[d.RunID], d.Spec); e != nil {
		v.Conditions = append(v.Conditions, e.Error())
	}
	v.Applicable = len(v.Conditions) == 0
	v.Qualified = v.Applicable && d.State == "completed" && d.Outcome != nil && d.Outcome.Verdict == "passed"
	if d.State == "completed" && d.Outcome == nil {
		v.Conditions = append(v.Conditions, "final answer is not a valid structured review")
	}
	_, latest, e := s.load()
	if e != nil {
		return ReviewView{}, e
	}
	if latest != head {
		return ReviewView{}, fail("stale_revision", "review changed during inspection")
	}
	return v, nil
}
func reservedReviewThread(state snapshot, thread string) bool {
	for _, d := range state.Reviews {
		if d.ThreadID == thread {
			return true
		}
	}
	return false
}
func reviewThreadUsed(state snapshot, thread string) bool {
	for _, r := range state.Runs {
		for _, a := range r.Attempts {
			if a.ThreadID == thread {
				return true
			}
		}
	}
	if state.Leader != nil {
		for _, d := range state.Leader.Dispatches {
			if d.ThreadID == thread {
				return true
			}
		}
	}
	for _, d := range state.Reviews {
		if d.ThreadID == thread {
			return true
		}
	}
	return false
}
func (s *Service) ReviewWrite(c Caller, q ReviewRequest) (Result, error) {
	if !c.controller {
		return Result{}, fail("unauthorized", "review mutations require trusted controller")
	}
	if q.Schema != ReviewSchema {
		return Result{}, fail("incompatible_profile", "unsupported review schema")
	}
	if !uuidPattern.MatchString(q.RunID) || !uuidPattern.MatchString(q.ReviewID) || len(q.Answer) > 64<<10 {
		return Result{}, fail("invalid_request", "bounded review identity/result required")
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
	if e = s.authorize(c, Schema, q.ProjectRef, q.CheckoutID, q.OperationID, state); e != nil {
		return Result{}, e
	}
	digest := requestDigest(q)
	if old, ok, err := retry(state, q.OperationID, digest); ok {
		return old, err
	}
	if len(state.Receipts) >= 1024 {
		return Result{}, fail("limit", "receipt capacity reached")
	}
	r, ok := state.Runs[q.RunID]
	if !ok {
		return Result{}, fail("missing_coverage", "run absent")
	}
	d, exists := state.Reviews[q.ReviewID]
	if q.Action != "prepare" && (!exists || d.RunID != r.ID) {
		return Result{}, fail("missing_coverage", "review absent")
	}
	if q.ExpectedRevision != d.Revision {
		return Result{}, fail("stale_revision", "review revision changed")
	}
	if q.Action != "prepare" && (q.Spec != nil || q.ExpectedRunRevision != 0) {
		return Result{}, fail("invalid_request", "prepare fields on another action")
	}
	if q.Action != "bind" && q.Action != "observe" && q.ThreadID != "" || q.Action != "observe" && (q.TurnID != "" || q.Status != "" || q.Answer != "") {
		return Result{}, fail("invalid_request", "unrelated review action fields")
	}
	now := time.Now().UTC().Format(time.RFC3339Nano)
	switch q.Action {
	case "prepare":
		if exists || q.ReviewID != q.OperationID || q.Spec == nil || q.ExpectedRunRevision != r.Revision {
			return Result{}, fail("stale_revision", "new review and current run revision required")
		}
		if e = s.reviewBasis(state, r, *q.Spec); e != nil {
			return Result{}, e
		}
		count := 0
		for _, other := range state.Reviews {
			if other.RunID == r.ID {
				count++
				if ReviewActive(other.State) {
					return Result{}, fail("conflict", "another review is active for this run")
				}
			}
		}
		if count >= 8 || len(state.Reviews) >= 64 {
			return Result{}, fail("limit", "review history capacity reached")
		}
		evidence := []EvidenceRecord{}
		for _, id := range q.Spec.EvidenceIDs {
			evidence = append(evidence, *state.Receipts[id].Result.Evidence)
		}
		payload, _ := json.Marshal(struct {
			Schema     string           `json:"schema"`
			Selection  Selection        `json:"selection"`
			Assignment Assignment       `json:"assignment"`
			Spec       ReviewSpec       `json:"review"`
			Evidence   []EvidenceRecord `json:"tests"`
		}{ReviewSchema, r.Selection, r.Assignment, *q.Spec, evidence})
		if len(payload) > 64<<10 {
			return Result{}, fail("limit", "review payload exceeds 64 KiB")
		}
		d = Review{ID: q.ReviewID, RunID: r.ID, RunRevision: r.Revision, Spec: *q.Spec, Payload: string(payload), State: "prepared", CreatedAt: now}
	case "open":
		if d.State != "prepared" {
			return Result{}, fail("conflict", "review is not prepared")
		}
		d.State = "opening"
	case "bind":
		if d.State != "opening" || !opaque(q.ThreadID) || reviewThreadUsed(state, q.ThreadID) {
			return Result{}, fail("conflict", "new distinct review thread required")
		}
		d.ThreadID = q.ThreadID
		d.State = "bound"
	case "send":
		if d.State != "bound" {
			return Result{}, fail("conflict", "review is not bound")
		}
		d.State = "sending"
	case "supersede":
		if d.State != "prepared" && d.State != "bound" {
			return Result{}, fail("conflict", "only unsent review can be superseded")
		}
		d.State = "superseded"
	case "observe":
		if (d.State != "sending" && d.State != "running") || q.ThreadID != d.ThreadID || !opaque(q.TurnID) || (d.TurnID != "" && d.TurnID != q.TurnID) || !slices.Contains([]string{"running", "completed", "failed", "interrupted"}, q.Status) || (q.Status == "running" && q.Answer != "") {
			return Result{}, fail("conflict", "exact active review observation required")
		}
		d.TurnID = q.TurnID
		d.State = q.Status
		d.Answer = q.Answer
		if q.Status == "completed" {
			d.Outcome = parseReviewOutcome(q.Answer)
		}
	default:
		return Result{}, fail("invalid_request", "unknown review action")
	}
	if q.Action == "prepare" || q.Action == "open" || q.Action == "send" {
		if e = s.reviewBasis(state, r, d.Spec); e != nil {
			return Result{}, e
		}
	}
	if d.Revision == ^uint64(0) {
		return Result{}, fail("limit", "review revision exhausted")
	}
	d.Revision++
	d.UpdatedAt = now
	if state.Reviews == nil {
		state.Reviews = map[string]Review{}
	}
	state.Reviews[d.ID] = d
	if !publicationProfile(state.Schema) {
		state.Schema = reviewStateSchema
	}
	out := Result{Schema: Schema, Operation: "review." + q.Action, ProjectRef: s.project, CheckoutID: state.CheckoutID, OperationID: q.OperationID, Run: summarizeRun(r), Coverage: "known", Review: &d}
	state.Receipts[q.OperationID] = receipt{Digest: digest, Result: out}
	if e = s.save(state, head); e != nil {
		return Result{}, e
	}
	return out, nil
}
func validateReviewState(state snapshot) error {
	if state.Schema != reviewStateSchema && !publicationProfile(state.Schema) {
		if state.Reviews != nil {
			return fail("incompatible_profile", "reviews require profile 0.11")
		}
		for _, v := range state.Receipts {
			if v.Result.Review != nil || strings.HasPrefix(v.Result.Operation, "review.") {
				return fail("incompatible_profile", "review receipt outside profile")
			}
		}
		return nil
	}
	if (!sessionProfile(state.Schema) && len(state.Reviews) == 0) || len(state.Reviews) > 64 {
		return fail("reconciliation_required", "invalid review collection")
	}
	counts, active := map[string]int{}, map[string]int{}
	threads := map[string]bool{}
	for id, d := range state.Reviews {
		r, ok := state.Runs[d.RunID]
		_, created := time.Parse(time.RFC3339Nano, d.CreatedAt)
		_, updated := time.Parse(time.RFC3339Nano, d.UpdatedAt)
		counts[d.RunID]++
		if ReviewActive(d.State) {
			active[d.RunID]++
		}
		if id != d.ID || !uuidPattern.MatchString(id) || !ok || !validReviewSpec(d.Spec) || d.Spec.Candidate.RunID != r.ID || d.Spec.Candidate.CheckoutID != state.CheckoutID || d.Spec.Candidate.ProjectRef != state.Project || d.Revision == 0 || d.RunRevision == 0 || d.RunRevision > r.Revision || len(d.Payload) == 0 || len(d.Payload) > 64<<10 || len(d.Answer) > 64<<10 || created != nil || updated != nil || counts[d.RunID] > 8 || active[d.RunID] > 1 || !slices.Contains([]string{"prepared", "opening", "bound", "sending", "running", "completed", "failed", "interrupted", "superseded"}, d.State) {
			return fail("reconciliation_required", "invalid retained review")
		}
		if (d.State == "prepared" || d.State == "opening") && d.ThreadID != "" || (d.State == "bound" || d.State == "sending" || d.State == "running" || d.State == "completed" || d.State == "failed" || d.State == "interrupted") && !opaque(d.ThreadID) {
			return fail("reconciliation_required", "invalid review thread mapping")
		}
		if d.ThreadID != "" {
			if threads[d.ThreadID] {
				return fail("reconciliation_required", "review thread reused")
			}
			threads[d.ThreadID] = true
		}
		if (d.State == "running" || d.State == "completed" || d.State == "failed" || d.State == "interrupted") && !opaque(d.TurnID) {
			return fail("reconciliation_required", "review turn mapping missing")
		}
		if d.State == "completed" {
			if requestDigest(d.Outcome) != requestDigest(parseReviewOutcome(d.Answer)) {
				return fail("reconciliation_required", "review verdict differs from observed answer")
			}
		} else if d.Outcome != nil {
			return fail("reconciliation_required", "uncompleted review has verdict")
		}
		latest := false
		for _, v := range state.Receipts {
			if x := v.Result.Review; x != nil && x.ID == id && x.Revision == d.Revision && requestDigest(x) == requestDigest(d) {
				latest = true
			}
		}
		if !latest {
			return fail("reconciliation_required", "review receipt missing")
		}
	}
	for id, v := range state.Receipts {
		x := v.Result.Review
		if x == nil {
			if strings.HasPrefix(v.Result.Operation, "review.") {
				return fail("reconciliation_required", "review receipt missing record")
			}
			continue
		}
		current, ok := state.Reviews[x.ID]
		if !ok || x.RunID != current.RunID || x.Revision == 0 || x.Revision > current.Revision || v.Result.OperationID != id || v.Result.Schema != Schema || v.Result.ProjectRef != state.Project || v.Result.CheckoutID != state.CheckoutID || !strings.HasPrefix(v.Result.Operation, "review.") {
			return fail("reconciliation_required", "invalid review receipt")
		}
	}
	return nil
}
