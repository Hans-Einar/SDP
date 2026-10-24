package governance

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func evidenceFile(t *testing.T, s *Service, path, text string) {
	t.Helper()
	p := filepath.Join(s.root, path)
	if e := os.MkdirAll(filepath.Dir(p), 0700); e != nil {
		t.Fatal(e)
	}
	if e := os.WriteFile(p, []byte(text), 0600); e != nil {
		t.Fatal(e)
	}
}

func evidenceFixture(t *testing.T) (*Service, Result, EvidenceRequest) {
	t.Helper()
	s, _, start := fixture(t)
	bind(t, s, &start)
	r, e := s.Start(LocalController(), start)
	if e != nil {
		t.Fatal(e)
	}
	r, e = s.Transition(LocalController(), transition(r, "start"))
	if e != nil {
		t.Fatal(e)
	}
	r, e = s.Transition(LocalController(), transition(r, "report"))
	if e != nil {
		t.Fatal(e)
	}
	evidenceFile(t, s, "target.txt", "tested candidate\n")
	evidenceFile(t, s, "dependencies/config.txt", "dependency\n")
	evidenceFile(t, s, "evidence/test.log", "all checks passed\n")
	c, e := s.CaptureCandidate(r.Run.ID, []string{"dependencies"})
	if e != nil {
		t.Fatal(e)
	}
	id, _ := identifier()
	q := EvidenceRequest{Schema: EvidenceSchema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, RunID: r.Run.ID, OperationID: id,
		ExpectedRevision: r.Run.Revision, InputRevision: r.Run.InputRevision, ReportID: r.OperationID, Candidate: c,
		Kind: "test", Verdict: "passed", Issuer: "local verifier", Note: "external fixture checks; no independent review claim",
		Artifact: Source{Path: "evidence/test.log", SHA256: hash([]byte("all checks passed\n"))}}
	return s, r, q
}

func TestEvidenceCandidateDriftAndHistoricalRetry(t *testing.T) {
	s, r, q := evidenceFixture(t)
	before, _ := os.ReadFile(filepath.Join(s.state, "HEAD"))
	if v, e := s.Evidence(r.Run.ID); e != nil || len(v.Items) != 0 {
		t.Fatal(v, e)
	}
	after, _ := os.ReadFile(filepath.Join(s.state, "HEAD"))
	if string(before) != string(after) {
		t.Fatal("read mutated state")
	}
	attached, e := s.AttachEvidence(LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	if attached.Run.State != "verifying" || attached.Run.Revision != r.Run.Revision {
		t.Fatal("attachment advanced work")
	}
	check := func(applicable, candidate, artifact bool) {
		t.Helper()
		v, err := s.Evidence(r.Run.ID)
		if err != nil || len(v.Items) != 1 || v.Items[0].Applicable != applicable || v.Items[0].CandidateCurrent != candidate || v.Items[0].ArtifactCurrent != artifact {
			t.Fatal(v, err)
		}
	}
	check(true, true, true)
	evidenceFile(t, s, "unrelated.txt", "unrelated dirty edit")
	check(true, true, true)
	evidenceFile(t, s, "dependencies/new.txt", "new relevant input")
	check(false, false, true)
	other, e := New(s.root, s.area, s.project, s.state)
	if e != nil {
		t.Fatal(e)
	}
	repeated, e := other.AttachEvidence(LocalController(), q)
	if e != nil || !reflect.DeepEqual(attached, repeated) {
		t.Fatal("retry lost historical receipt", e)
	}
	q.OperationID, _ = identifier()
	_, e = s.AttachEvidence(LocalController(), q)
	code(t, e, "stale_candidate")
	if e = os.Remove(filepath.Join(s.root, "dependencies/new.txt")); e != nil {
		t.Fatal(e)
	}
	check(true, true, true)
	evidenceFile(t, s, "evidence/test.log", "edited evidence")
	check(false, true, false)
	_, e = s.Transition(LocalController(), transition(r, "complete"))
	code(t, e, "evidence_missing")
	evidenceFile(t, s, "evidence/test.log", "all checks passed\n")
	evidenceFile(t, s, r.Run.Selection.Authority, "different authority")
	check(false, true, true)
}

func TestCandidateCapturesAbsenceDeletionModeAndDirectoryMembership(t *testing.T) {
	s, r, _ := evidenceFixture(t)
	initial, e := s.CaptureCandidate(r.Run.ID, []string{"future.txt"})
	if e != nil {
		t.Fatal(e)
	}
	if initial.Entries[0].Kind != "missing" {
		t.Fatal(initial)
	}
	evidenceFile(t, s, "future.txt", "new")
	created, e := s.CaptureCandidate(r.Run.ID, []string{"future.txt"})
	if e != nil || created.Digest == initial.Digest {
		t.Fatal(e)
	}
	if e = os.Remove(filepath.Join(s.root, "future.txt")); e != nil {
		t.Fatal(e)
	}
	restored, e := s.CaptureCandidate(r.Run.ID, []string{"future.txt"})
	if e != nil || restored.Digest != initial.Digest {
		t.Fatal(e)
	}
	if e = os.Chmod(filepath.Join(s.root, "target.txt"), 0700); e != nil {
		t.Fatal(e)
	}
	mode, e := s.CaptureCandidate(r.Run.ID, []string{"future.txt"})
	if e != nil || mode.Digest == initial.Digest {
		t.Fatal(e)
	}
	if e = os.Chmod(filepath.Join(s.root, "target.txt"), 0610); e != nil {
		t.Fatal(e)
	}
	otherMode, e := s.CaptureCandidate(r.Run.ID, []string{"future.txt"})
	if e != nil || otherMode.Digest == mode.Digest {
		t.Fatal("executable bit positions were collapsed", e)
	}
	if e = os.Remove(filepath.Join(s.root, "target.txt")); e != nil {
		t.Fatal(e)
	}
	deleted, e := s.CaptureCandidate(r.Run.ID, nil)
	if e != nil || deleted.Entries[0].Kind != "missing" {
		t.Fatal(deleted, e)
	}
}

func TestCandidateAndEvidenceBoundaries(t *testing.T) {
	s, r, q := evidenceFixture(t)
	for _, paths := range [][]string{{"../outside"}, {"dependencies", "dependencies"}, {"target*.txt"}, {".git"}} {
		if _, e := s.CaptureCandidate(r.Run.ID, paths); e == nil {
			t.Fatal("accepted invalid dependency", paths)
		}
	}
	if e := os.Symlink(filepath.Join(s.root, "target.txt"), filepath.Join(s.root, "dependencies/link")); e != nil {
		t.Fatal(e)
	}
	_, e := s.CaptureCandidate(r.Run.ID, []string{"dependencies"})
	code(t, e, "path")
	os.Remove(filepath.Join(s.root, "dependencies/link"))
	evidenceFile(t, s, "large.txt", strings.Repeat("x", (4<<20)+1))
	_, e = s.CaptureCandidate(r.Run.ID, []string{"large.txt"})
	code(t, e, "limit")
	for _, mutate := range []func(*EvidenceRequest){
		func(v *EvidenceRequest) { v.Candidate.Entries = nil; v.Candidate.Digest = candidateDigest(v.Candidate) },
		func(v *EvidenceRequest) { v.ReportID, _ = identifier() },
		func(v *EvidenceRequest) { v.ExpectedRevision++ },
		func(v *EvidenceRequest) {
			v.Candidate.WritablePaths = []string{"unrelated.txt"}
			v.Candidate.Digest = candidateDigest(v.Candidate)
		},
		func(v *EvidenceRequest) { v.Kind = "independent-review" },
		func(v *EvidenceRequest) { v.Artifact.Path = "target.txt" },
	} {
		v := q
		mutate(&v)
		if _, e := s.AttachEvidence(LocalController(), v); e == nil {
			t.Fatal("accepted forged/stale evidence", v)
		}
	}
	_, e = s.AttachEvidence(Caller{}, q)
	code(t, e, "unauthorized")
	q.Kind = "review"
	v, e := s.AttachEvidence(LocalController(), q)
	if e != nil || v.Evidence.Provenance != "controller-attributed" {
		t.Fatal(v, e)
	}
	view, e := s.Evidence(r.Run.ID)
	if e != nil || !strings.Contains(strings.Join(view.Limitations, " "), "review independence unverified") {
		t.Fatal(view, e)
	}
}

func TestCandidateExplicitReadOnlyDependencyAndLimits(t *testing.T) {
	s, r, _ := evidenceFixture(t)
	before, e := s.CaptureCandidate(r.Run.ID, []string{"unrelated.txt"})
	if e != nil {
		t.Fatal(e)
	}
	evidenceFile(t, s, "unrelated.txt", "explicitly selected read-only dependency changed")
	after, e := s.CaptureCandidate(r.Run.ID, []string{"unrelated.txt"})
	if e != nil || before.Digest == after.Digest {
		t.Fatal(e)
	}
	for i := 0; i < 513; i++ {
		evidenceFile(t, s, "many/"+strings.Repeat("a", i%12)+fmt.Sprint(i), "x")
	}
	_, e = s.CaptureCandidate(r.Run.ID, []string{"many"})
	code(t, e, "limit")
	for i := 0; i < 9; i++ {
		evidenceFile(t, s, "bulk/"+fmt.Sprint(i), strings.Repeat("x", 4<<20))
	}
	_, e = s.CaptureCandidate(r.Run.ID, []string{"bulk"})
	code(t, e, "limit")
}

func TestEvidenceProfileSurvivesExistingWriters(t *testing.T) {
	s, r, q := evidenceFixture(t)
	if _, e := s.AttachEvidence(LocalController(), q); e != nil {
		t.Fatal(e)
	}
	leaderApply(t, s, leaderRequest(t, s, "message", ""))
	state, head, e := s.load()
	if e != nil || state.Schema != evidenceStateSchema {
		t.Fatal(state.Schema, e)
	}
	ack := ReportAcknowledgementRequest{Schema: ReportAcknowledgementSchema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID,
		OperationID: newEvidenceID(t), RunID: r.Run.ID, ReportID: r.OperationID, ReportedRevision: r.Run.Revision,
		InputRevision: r.Run.InputRevision, ExpectedRevision: r.Run.Revision, ExpectedGeneration: head, Coordinator: "leader", Note: "processed without accepting"}
	if _, e = s.AcknowledgeReport(LocalController(), ack); e != nil {
		t.Fatal(e)
	}
	if _, e = s.Transition(LocalController(), transition(r, "cancel")); e != nil {
		t.Fatal(e)
	}
	state, _, e = s.load()
	if e != nil || state.Schema != evidenceStateSchema {
		t.Fatal(state.Schema, e)
	}
	v, e := s.Evidence(r.Run.ID)
	if e != nil || len(v.Items) != 1 || v.Items[0].Applicable {
		t.Fatal(v, e)
	}
}

func newEvidenceID(t *testing.T) string {
	t.Helper()
	id, e := identifier()
	if e != nil {
		t.Fatal(e)
	}
	return id
}

func TestEvidenceAtomicFaultsAndOperationNamespace(t *testing.T) {
	for _, fault := range []string{"before-head", "after-head"} {
		t.Run(fault, func(t *testing.T) {
			s, r, q := evidenceFixture(t)
			s.fault = func(point string) error {
				if point == fault {
					return errors.New("injected crash")
				}
				return nil
			}
			if _, e := s.AttachEvidence(LocalController(), q); e == nil {
				t.Fatal("fault not injected")
			}
			s.fault = nil
			other, e := New(s.root, s.area, s.project, s.state)
			if e != nil {
				t.Fatal(e)
			}
			if _, e = other.AttachEvidence(LocalController(), q); e != nil {
				t.Fatal(e)
			}
			v, e := other.Evidence(r.Run.ID)
			if e != nil || len(v.Items) != 1 {
				t.Fatal(v, e)
			}
			q.Note = "different request bytes"
			_, e = other.AttachEvidence(LocalController(), q)
			code(t, e, "idempotency_conflict")
			tq := transition(r, "cancel")
			tq.OperationID = q.OperationID
			_, e = other.Transition(LocalController(), tq)
			code(t, e, "idempotency_conflict")
		})
	}
}

func reworkRequest(r Result, evidenceID string) EvidenceReworkRequest {
	op, _ := identifier()
	return EvidenceReworkRequest{Schema: EvidenceSchema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, OperationID: op, RunID: r.Run.ID,
		ExpectedRevision: r.Run.Revision, InputRevision: r.Run.InputRevision, EvidenceID: evidenceID, Reason: "correct the failed case within original scope"}
}

func TestEvidenceReworkPreservesHistoryAndRevokesPriorAgent(t *testing.T) {
	s, r, q := evidenceFixture(t)
	state, head, e := s.load()
	if e != nil {
		t.Fatal(e)
	}
	at, _ := identifier()
	epoch, _ := identifier()
	actual := state.Runs[r.Run.ID]
	actual.Attempts = []Attempt{{ID: at, Epoch: epoch, Status: "completed", InputRevision: r.Run.InputRevision}}
	state.Runs[r.Run.ID] = actual
	if e = s.save(state, head); e != nil {
		t.Fatal(e)
	}
	q.Verdict = "failed"
	attached, e := s.AttachEvidence(LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	rework := reworkRequest(attached, q.OperationID)
	agent, _ := AssignedAgent(r.Run.ID, at, epoch)
	_, e = s.ReworkEvidence(agent, rework)
	code(t, e, "unauthorized")
	result, e := s.ReworkEvidence(LocalController(), rework)
	if e != nil {
		t.Fatal(e)
	}
	if result.Run.State != "ready" || result.Run.Revision != r.Run.Revision+1 || result.Run.Selection != r.Run.Selection {
		t.Fatal(result)
	}
	if result.Run.Attempts[0].CapabilityEpoch() == epoch {
		t.Fatal("prior agent capability survived rework")
	}
	v, e := s.Evidence(r.Run.ID)
	if e != nil || len(v.Items) != 1 || v.Items[0].Applicable || !v.Items[0].CandidateCurrent {
		t.Fatal(v, e)
	}
	retry, e := s.ReworkEvidence(LocalController(), rework)
	if e != nil || !reflect.DeepEqual(result, retry) {
		t.Fatal(retry, e)
	}
	result, e = s.Transition(LocalController(), transition(result, "start"))
	if e != nil {
		t.Fatal(e)
	}
	_, e = s.Transition(agent, transition(result, "report"))
	code(t, e, "unauthorized")
	result, e = s.Transition(LocalController(), transition(result, "report"))
	if e != nil {
		t.Fatal(e)
	}
	v, e = s.Evidence(r.Run.ID)
	if e != nil || v.Items[0].ReportCurrent {
		t.Fatal("old report reused for new work", v, e)
	}
}

func TestEvidenceReworkRefusesUncertainOrPassedWork(t *testing.T) {
	for _, status := range []string{"passed", "active", "gap", "input", "changed", "claim"} {
		t.Run(status, func(t *testing.T) {
			s, r, q := evidenceFixture(t)
			q.Verdict = "failed"
			if status == "passed" {
				q.Verdict = "passed"
			}
			if status == "claim" {
				q.Kind = "claim"
			}
			v, e := s.AttachEvidence(LocalController(), q)
			if e != nil {
				t.Fatal(e)
			}
			state, head, _ := s.load()
			actual := state.Runs[r.Run.ID]
			if status == "active" || status == "gap" {
				actual.Attempts = []Attempt{{Status: "running", Gap: status == "gap"}}
			}
			if status == "input" {
				actual.Inputs = []OwnerInput{{State: "uncertain"}}
			}
			state.Runs[r.Run.ID] = actual
			if e = s.save(state, head); e != nil {
				t.Fatal(e)
			}
			if status == "changed" {
				evidenceFile(t, s, "target.txt", "untested candidate")
			}
			if _, e = s.ReworkEvidence(LocalController(), reworkRequest(v, q.OperationID)); e == nil {
				t.Fatal("unsafe rework accepted")
			}
		})
	}
}

func TestEvidenceStateProfileAndWireCapture(t *testing.T) {
	s, r, q := evidenceFixture(t)
	attached, e := s.AttachEvidence(LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	state, head, e := s.load()
	if e != nil || state.Schema != evidenceStateSchema {
		t.Fatal(state.Schema, e)
	}
	state.Schema = leaderStateSchema
	if e = validateEvidenceState(state); e == nil {
		t.Fatal("old profile accepted evidence")
	}
	state.Schema = evidenceStateSchema
	if e = validateLeaderState(state); e != nil {
		t.Fatal(e)
	}
	state.Receipts[q.OperationID].Result.Evidence.Request.Note = "corrupt retained request"
	if e = s.save(state, head); e != nil {
		t.Fatal(e)
	}
	_, _, e = s.load()
	code(t, e, "reconciliation_required")
	// Capture actual operation objects for JSON Schema validation when requested.
	if path := os.Getenv("PG_EVIDENCE_WIRE"); path != "" {
		s, r, q = evidenceFixture(t)
		q.Verdict = "failed"
		attached, e = s.AttachEvidence(LocalController(), q)
		if e != nil {
			t.Fatal(e)
		}
		view, e := s.Evidence(r.Run.ID)
		if e != nil {
			t.Fatal(e)
		}
		rq := reworkRequest(attached, q.OperationID)
		reworked, e := s.ReworkEvidence(LocalController(), rq)
		if e != nil {
			t.Fatal(e)
		}
		cases := []map[string]any{}
		for _, v := range []struct {
			name  string
			value any
		}{{"Candidate", q.Candidate}, {"EvidenceRequest", q}, {"Result", attached}, {"EvidenceView", view}, {"EvidenceReworkRequest", rq}, {"Result", reworked}} {
			cases = append(cases, map[string]any{"definition": v.name, "value": v.value})
		}
		b, _ := json.MarshalIndent(cases, "", "  ")
		if e = os.WriteFile(path, b, 0600); e != nil {
			t.Fatal(e)
		}
	}
}
