package governance

import (
	"path/filepath"
	"slices"
	"time"
)

type EvidenceRequest struct {
	Schema           string    `json:"schema"`
	ProjectRef       string    `json:"projectRef"`
	CheckoutID       string    `json:"checkoutId"`
	OperationID      string    `json:"operationId"`
	RunID            string    `json:"runId"`
	ExpectedRevision uint64    `json:"expectedRevision"`
	InputRevision    string    `json:"inputRevision"`
	ReportID         string    `json:"reportId"`
	Candidate        Candidate `json:"candidate"`
	Kind             string    `json:"kind"`
	Verdict          string    `json:"verdict"`
	Issuer           string    `json:"issuer"`
	Note             string    `json:"note"`
	Artifact         Source    `json:"artifact"`
}

type EvidenceRecord struct {
	Request    EvidenceRequest `json:"request"`
	RecordedAt string          `json:"recordedAt"`
	Provenance string          `json:"provenance"`
}

type EvidenceItem struct {
	Record           EvidenceRecord `json:"record"`
	CandidateCurrent bool           `json:"candidateCurrent"`
	ArtifactCurrent  bool           `json:"artifactCurrent"`
	ReportCurrent    bool           `json:"reportCurrent"`
	Applicable       bool           `json:"applicable"`
	Conditions       []string       `json:"conditions"`
}

type EvidenceView struct {
	Schema      string         `json:"schema"`
	ProjectRef  string         `json:"projectRef"`
	CheckoutID  string         `json:"checkoutId"`
	RunID       string         `json:"runId"`
	Revision    uint64         `json:"revision"`
	State       string         `json:"state"`
	Items       []EvidenceItem `json:"items"`
	Limitations []string       `json:"limitations"`
}

func validEvidence(q EvidenceRequest) bool {
	c := q.Candidate
	return q.Schema == EvidenceSchema && uuidPattern.MatchString(q.OperationID) && uuidPattern.MatchString(q.RunID) &&
		uuidPattern.MatchString(q.CheckoutID) && uuidPattern.MatchString(q.ReportID) && q.ExpectedRevision > 0 &&
		digestPattern.MatchString(q.InputRevision) && slices.Contains([]string{"claim", "test", "review"}, q.Kind) &&
		slices.Contains([]string{"passed", "failed", "inconclusive"}, q.Verdict) && boundedAttentionText(q.Issuer, 256) && boundedAttentionText(q.Note, 4096) &&
		literalCandidatePath(q.Artifact.Path) && digestPattern.MatchString(q.Artifact.SHA256) &&
		c.Schema == EvidenceSchema && c.ProjectRef == q.ProjectRef && c.CheckoutID == q.CheckoutID && c.RunID == q.RunID && c.InputRevision == q.InputRevision &&
		len(c.Entries) <= 512 && len(c.Dependencies) <= 64 && digestPattern.MatchString(c.Digest) && candidateDigest(c) == c.Digest
}

func latestReport(state snapshot, r Run) string {
	id := ""
	var revision uint64
	for key, v := range state.Receipts {
		p := v.Result
		if p.Operation == "transition" && p.Run.ID == r.ID && p.Run.State == "verifying" && p.Run.InputRevision == r.InputRevision && p.Run.Revision > revision {
			id, revision = key, p.Run.Revision
		}
	}
	return id
}

func (s *Service) evidenceCurrent(state snapshot, r Run, record EvidenceRecord) EvidenceItem {
	v := EvidenceItem{Record: record, Conditions: []string{}}
	q := record.Request
	current, err := s.captureCandidate(state.CheckoutID, r, q.Candidate.Dependencies)
	v.CandidateCurrent = err == nil && current.Digest == q.Candidate.Digest
	if !v.CandidateCurrent {
		v.Conditions = append(v.Conditions, "candidate changed or unavailable")
	}
	b, err := readFile(filepath.Join(s.root, q.Artifact.Path), 4<<20)
	v.ArtifactCurrent = err == nil && hash(b) == q.Artifact.SHA256
	if !v.ArtifactCurrent {
		v.Conditions = append(v.Conditions, "evidence artifact changed or unavailable")
	}
	v.ReportCurrent = r.State == "verifying" && latestReport(state, r) == q.ReportID
	if !v.ReportCurrent {
		v.Conditions = append(v.Conditions, "report no longer awaiting this verification")
	}
	if err = s.currentAttentionInputs(r); err != nil {
		v.Conditions = append(v.Conditions, "authority or assignment coverage changed")
	}
	v.Applicable = len(v.Conditions) == 0
	return v
}

func (s *Service) Evidence(runID string) (EvidenceView, error) {
	state, head, e := s.load()
	if e != nil {
		return EvidenceView{}, e
	}
	r, ok := state.Runs[runID]
	if !ok || !uuidPattern.MatchString(runID) {
		return EvidenceView{}, fail("missing_coverage", "run not found")
	}
	v := EvidenceView{Schema: EvidenceSchema, ProjectRef: s.project, CheckoutID: state.CheckoutID, RunID: r.ID, Revision: r.Revision, State: r.State,
		Items: []EvidenceItem{}, Limitations: []string{"controller-attributed claims; test execution not observed by PG", "review independence unverified", "checked publication and completion unavailable"}}
	for _, receipt := range state.Receipts {
		if a := receipt.Result.Evidence; a != nil && a.Request.RunID == r.ID {
			v.Items = append(v.Items, s.evidenceCurrent(state, r, *a))
		}
	}
	slices.SortFunc(v.Items, func(a, b EvidenceItem) int {
		if a.Record.RecordedAt < b.Record.RecordedAt {
			return -1
		}
		if a.Record.RecordedAt > b.Record.RecordedAt {
			return 1
		}
		if a.Record.Request.OperationID < b.Record.Request.OperationID {
			return -1
		}
		return 1
	})
	_, latest, e := s.load()
	if e != nil {
		return EvidenceView{}, e
	}
	if latest != head {
		return EvidenceView{}, fail("stale_revision", "run store changed during evidence inspection")
	}
	return v, nil
}

func (s *Service) AttachEvidence(caller Caller, q EvidenceRequest) (Result, error) {
	if !caller.controller {
		return Result{}, fail("unauthorized", "evidence capture requires trusted controller")
	}
	if q.Schema != EvidenceSchema {
		return Result{}, fail("incompatible_profile", "unsupported evidence schema")
	}
	if !validEvidence(q) {
		return Result{}, fail("invalid_request", "exact candidate/report, artifact digest and bounded attribution required")
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
	if old, ok, err := retry(state, q.OperationID, digest); ok {
		return old, err
	}
	r, ok := state.Runs[q.RunID]
	if !ok {
		return Result{}, fail("missing_coverage", "run not found")
	}
	if r.Revision != q.ExpectedRevision || r.InputRevision != q.InputRevision {
		return Result{}, fail("stale_revision", "run or inputs changed")
	}
	if r.State != "verifying" || latestReport(state, r) != q.ReportID {
		return Result{}, fail("conflict", "latest report awaiting verification required")
	}
	count := 0
	for _, v := range state.Receipts {
		if v.Result.Evidence != nil && v.Result.Run.ID == r.ID {
			count++
		}
	}
	if count >= 16 || len(state.Receipts) >= 1024 {
		return Result{}, fail("limit", "evidence receipt capacity reached")
	}
	for _, path := range append(append([]string{}, q.Candidate.WritablePaths...), q.Candidate.Dependencies...) {
		if within(q.Artifact.Path, path) {
			return Result{}, fail("conflict", "evidence artifact must be outside candidate scope")
		}
	}
	record := EvidenceRecord{Request: q, RecordedAt: time.Now().UTC().Format(time.RFC3339Nano), Provenance: "controller-attributed"}
	if !s.evidenceCurrent(state, r, record).Applicable {
		return Result{}, fail("stale_candidate", "candidate, artifact, report or authority changed")
	}
	out := Result{Schema: Schema, Operation: "evidence.attach", ProjectRef: s.project, CheckoutID: state.CheckoutID, OperationID: q.OperationID,
		Run: summarizeRun(r), Coverage: "known", Evidence: &record}
	if state.Schema != reviewStateSchema {
		state.Schema = evidenceStateSchema
	}
	state.Receipts[q.OperationID] = receipt{Digest: digest, Result: out}
	if !s.evidenceCurrent(state, r, record).Applicable {
		return Result{}, fail("stale_candidate", "evidence basis changed before publication")
	}
	if e = s.save(state, head); e != nil {
		return Result{}, e
	}
	return out, nil
}

type EvidenceReworkRequest struct {
	Schema           string `json:"schema"`
	ProjectRef       string `json:"projectRef"`
	CheckoutID       string `json:"checkoutId"`
	OperationID      string `json:"operationId"`
	RunID            string `json:"runId"`
	ExpectedRevision uint64 `json:"expectedRevision"`
	InputRevision    string `json:"inputRevision"`
	EvidenceID       string `json:"evidenceId"`
	Reason           string `json:"reason"`
}

func (s *Service) ReworkEvidence(caller Caller, q EvidenceReworkRequest) (Result, error) {
	if !caller.controller {
		return Result{}, fail("unauthorized", "rework requires trusted controller")
	}
	if q.Schema != EvidenceSchema {
		return Result{}, fail("incompatible_profile", "unsupported evidence schema")
	}
	if !uuidPattern.MatchString(q.RunID) || !uuidPattern.MatchString(q.EvidenceID) || !digestPattern.MatchString(q.InputRevision) || q.ExpectedRevision == 0 || !boundedAttentionText(q.Reason, 4096) {
		return Result{}, fail("invalid_request", "exact failed evidence and bounded rework reason required")
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
	if old, ok, err := retry(state, q.OperationID, digest); ok {
		return old, err
	}
	r, ok := state.Runs[q.RunID]
	if !ok {
		return Result{}, fail("missing_coverage", "run not found")
	}
	if r.Revision != q.ExpectedRevision || r.InputRevision != q.InputRevision {
		return Result{}, fail("stale_revision", "run or inputs changed")
	}
	v := state.Receipts[q.EvidenceID].Result.Evidence
	if v == nil || v.Request.RunID != r.ID || v.Request.Verdict != "failed" || v.Request.Kind == "claim" {
		return Result{}, fail("evidence_missing", "failed test or review evidence required")
	}
	if !s.evidenceCurrent(state, r, *v).Applicable {
		return Result{}, fail("stale_candidate", "failed evidence is no longer current")
	}
	if hasBlockingAttention(r) {
		return Result{}, fail("attention_required", "blocking questions remain open")
	}
	for _, at := range r.Attempts {
		if !slices.Contains([]string{"completed", "failed", "interrupted"}, at.Status) || at.Gap {
			return Result{}, fail("reconciliation_required", "rework requires terminal observed attempts without gaps")
		}
	}
	for _, in := range r.Inputs {
		if in.State != "acknowledged" {
			return Result{}, fail("reconciliation_required", "unresolved controller input prevents rework")
		}
	}
	if len(state.Receipts) >= 1024 || r.Revision == ^uint64(0) {
		return Result{}, fail("limit", "run/receipt capacity reached")
	}
	if len(r.Attempts) > 0 {
		epoch, err := identifier()
		if err != nil {
			return Result{}, err
		}
		r.Attempts[len(r.Attempts)-1].AgentEpoch = epoch
	}
	r.State, r.Step, r.Reason = "ready", "implement", q.Reason
	r.Revision++
	out := Result{Schema: Schema, Operation: "evidence.rework", EvidenceRework: &q, ProjectRef: s.project, CheckoutID: state.CheckoutID, OperationID: q.OperationID, Run: summarizeRun(r), Coverage: "known"}
	if !s.evidenceCurrent(state, state.Runs[r.ID], *v).Applicable {
		return Result{}, fail("stale_candidate", "evidence basis changed before rework")
	}
	state.Runs[r.ID] = r
	state.Receipts[q.OperationID] = receipt{Digest: digest, Result: out}
	if e = s.save(state, head); e != nil {
		return Result{}, e
	}
	return out, nil
}

func validateEvidenceState(state snapshot) error {
	counts := map[string]int{}
	for id, receipt := range state.Receipts {
		v := receipt.Result
		if q := v.EvidenceRework; q != nil {
			a := state.Receipts[q.EvidenceID].Result.Evidence
			r, exists := state.Runs[q.RunID]
			if !evidenceProfile(state.Schema) || v.Operation != "evidence.rework" || v.Evidence != nil || q.Schema != EvidenceSchema ||
				id != q.OperationID || v.OperationID != id || receipt.Digest != requestDigest(q) ||
				q.ProjectRef != state.Project || q.CheckoutID != state.CheckoutID || v.ProjectRef != state.Project || v.CheckoutID != state.CheckoutID ||
				!exists || q.ExpectedRevision == 0 || q.ExpectedRevision == ^uint64(0) || r.Revision <= q.ExpectedRevision ||
				v.Run.ID != q.RunID || v.Run.State != "ready" || v.Run.Revision != q.ExpectedRevision+1 || v.Run.InputRevision != q.InputRevision ||
				a == nil || a.Request.RunID != q.RunID || a.Request.Kind == "claim" || a.Request.Verdict != "failed" {
				return fail("reconciliation_required", "invalid retained evidence rework")
			}
			continue
		}
		a := v.Evidence
		if a == nil {
			if v.Operation == "evidence.attach" || v.Operation == "evidence.rework" {
				return fail("reconciliation_required", "invalid evidence receipt")
			}
			continue
		}
		q := a.Request
		r, exists := state.Runs[q.RunID]
		report, found := state.Receipts[q.ReportID]
		_, err := time.Parse(time.RFC3339Nano, a.RecordedAt)
		counts[q.RunID]++
		if !evidenceProfile(state.Schema) || !validEvidence(q) || a.Provenance != "controller-attributed" || err != nil ||
			id != q.OperationID || receipt.Digest != requestDigest(q) || v.Operation != "evidence.attach" || v.Schema != Schema || v.OperationID != id ||
			q.ProjectRef != state.Project || q.CheckoutID != state.CheckoutID || v.ProjectRef != state.Project || v.CheckoutID != state.CheckoutID ||
			!exists || r.Revision < q.ExpectedRevision || v.Run.ID != q.RunID || v.Run.Revision != q.ExpectedRevision || v.Run.State != "verifying" || v.Run.InputRevision != q.InputRevision ||
			!found || report.Result.Operation != "transition" || report.Result.Run.ID != q.RunID || report.Result.Run.State != "verifying" || report.Result.Run.InputRevision != q.InputRevision || counts[q.RunID] > 16 {
			return fail("reconciliation_required", "invalid retained evidence")
		}
	}
	if state.Schema == evidenceStateSchema && len(counts) == 0 {
		return fail("reconciliation_required", "evidence profile has no records")
	}
	return nil
}
