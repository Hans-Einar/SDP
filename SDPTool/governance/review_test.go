package governance

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func reviewFixture(t *testing.T) (*Service, Result, ReviewRequest) {
	t.Helper()
	s, r, ev := evidenceFixture(t)
	if _, e := s.AttachEvidence(LocalController(), ev); e != nil {
		t.Fatal(e)
	}
	evidenceFile(t, s, "Skills/sdp-reviewer/SKILL.md", "review independently")
	spec := ReviewSpec{Instructions: "Read-only independent review", Candidate: ev.Candidate, ReportID: ev.ReportID, EvidenceIDs: []string{ev.OperationID}, Focus: "Check exact owner outcome and regressions", HostVersion: "0.164.0", MCPBinary: "/fixture/sdp-mcp", Skills: []SkillInput{{Name: "sdp-reviewer", Path: filepath.Join(s.root, "Skills/sdp-reviewer/SKILL.md"), Digest: hash([]byte("review independently"))}}}
	id := newEvidenceID(t)
	return s, r, ReviewRequest{Schema: ReviewSchema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, OperationID: id, RunID: r.Run.ID, ReviewID: id, ExpectedRunRevision: r.Run.Revision, Action: "prepare", Spec: &spec}
}
func reviewAction(t *testing.T, s *Service, id, action string) ReviewRequest {
	t.Helper()
	v, e := s.ReviewRead(id)
	if e != nil {
		t.Fatal(e)
	}
	return ReviewRequest{Schema: ReviewSchema, ProjectRef: v.ProjectRef, CheckoutID: v.CheckoutID, OperationID: newEvidenceID(t), RunID: v.Review.RunID, ReviewID: id, ExpectedRevision: v.Review.Revision, Action: action}
}
func reviewApply(t *testing.T, s *Service, q ReviewRequest) Result {
	t.Helper()
	r, e := s.ReviewWrite(LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	return r
}
func reviewFinish(t *testing.T, s *Service, id, answer string) {
	t.Helper()
	reviewApply(t, s, reviewAction(t, s, id, "open"))
	q := reviewAction(t, s, id, "bind")
	q.ThreadID = "fresh-review"
	reviewApply(t, s, q)
	reviewApply(t, s, reviewAction(t, s, id, "send"))
	q = reviewAction(t, s, id, "observe")
	q.ThreadID = "fresh-review"
	q.TurnID = "review-turn"
	q.Status = "completed"
	q.Answer = answer
	reviewApply(t, s, q)
}
func TestReviewQualificationAndHistoricalDrift(t *testing.T) {
	s, r, q := reviewFixture(t)
	prepared := reviewApply(t, s, q)
	v, e := s.ReviewRead(q.ReviewID)
	if e != nil || v.Qualified || !v.Applicable || v.ContextSeparation != "not-observed" {
		t.Fatal(v, e)
	}
	if prepared.Run.Revision != r.Run.Revision || prepared.Run.State != "verifying" {
		t.Fatal("prepare advanced implementation")
	}
	reviewFinish(t, s, q.ReviewID, `{"verdict":"passed","summary":"Candidate satisfies scope","findings":[]}`)
	v, e = s.ReviewRead(q.ReviewID)
	if e != nil || !v.Qualified || v.ContextSeparation != "controller-observed-fresh-thread" {
		t.Fatal(v, e)
	}
	evidenceFile(t, s, "unrelated.txt", "preserved unrelated change")
	v, e = s.ReviewRead(q.ReviewID)
	if e != nil || !v.Qualified {
		t.Fatal(v, e)
	}
	evidenceFile(t, s, "target.txt", "changed relevant candidate")
	v, e = s.ReviewRead(q.ReviewID)
	if e != nil || v.Qualified || v.Applicable || v.Review.Outcome.Verdict != "passed" {
		t.Fatal(v, e)
	}
	repeated, e := s.ReviewWrite(LocalController(), q)
	if e != nil || !reflect.DeepEqual(prepared, repeated) {
		t.Fatal("historical retry", e)
	}
	_, e = s.Transition(LocalController(), transition(r, "complete"))
	code(t, e, "evidence_missing")
	if file := os.Getenv("PG_REVIEW_CORE_WIRE"); file != "" {
		cases := []map[string]any{{"definition": "ReviewRequest", "value": q}, {"definition": "Result", "value": prepared}, {"definition": "ReviewView", "value": v}}
		b, _ := json.MarshalIndent(cases, "", "  ")
		if e = os.WriteFile(file, b, 0600); e != nil {
			t.Fatal(e)
		}
	}
}
func TestReviewRejectsFalseIndependenceAndStalePreparation(t *testing.T) {
	for _, mode := range []string{"caller", "skill", "candidate", "evidence", "authority", "active", "thread"} {
		t.Run(mode, func(t *testing.T) {
			s, r, q := reviewFixture(t)
			switch mode {
			case "caller":
				_, e := s.ReviewWrite(Caller{}, q)
				code(t, e, "unauthorized")
				return
			case "skill":
				q.Spec.Skills[0].Name = "worker"
			case "candidate":
				evidenceFile(t, s, "target.txt", "changed")
			case "evidence":
				evidenceFile(t, s, "evidence/test.log", "changed")
			case "authority":
				evidenceFile(t, s, r.Run.Selection.Authority, "revoked")
			case "active":
				state, head, _ := s.load()
				run := state.Runs[r.Run.ID]
				run.Attempts = []Attempt{{Status: "running"}}
				state.Runs[run.ID] = run
				if e := s.save(state, head); e != nil {
					t.Fatal(e)
				}
			case "thread":
				reviewApply(t, s, q)
				reviewApply(t, s, reviewAction(t, s, q.ReviewID, "open"))
				state, head, _ := s.load()
				run := state.Runs[r.Run.ID]
				run.Attempts = []Attempt{{Status: "completed", ThreadID: "implementation-thread"}}
				state.Runs[run.ID] = run
				if e := s.save(state, head); e != nil {
					t.Fatal(e)
				}
				bind := reviewAction(t, s, q.ReviewID, "bind")
				bind.ThreadID = "implementation-thread"
				_, e := s.ReviewWrite(LocalController(), bind)
				code(t, e, "conflict")
				return
			}
			if _, e := s.ReviewWrite(LocalController(), q); e == nil {
				t.Fatal("accepted invalid review", mode)
			}
		})
	}
}
func TestReviewTerminalDoesNotInventPassingResult(t *testing.T) {
	for _, answer := range []string{"Looks good", `{"verdict":"passed","summary":"fine","findings":["bug"]}`, `{"verdict":"failed","summary":"bug","findings":["target.txt is wrong"]}`, `{"verdict":"inconclusive","summary":"missing evidence","findings":[]}`} {
		s, _, q := reviewFixture(t)
		reviewApply(t, s, q)
		reviewFinish(t, s, q.ReviewID, answer)
		v, e := s.ReviewRead(q.ReviewID)
		if e != nil || v.Qualified || v.Review.Answer != answer {
			t.Fatal(v, e)
		}
	}
}
func TestReviewCrashRetryAndProfilePreservation(t *testing.T) {
	for _, point := range []string{"before-head", "after-head"} {
		t.Run(point, func(t *testing.T) {
			s, _, q := reviewFixture(t)
			s.fault = func(p string) error {
				if p == point {
					return errors.New("crash")
				}
				return nil
			}
			if _, e := s.ReviewWrite(LocalController(), q); e == nil {
				t.Fatal("no crash")
			}
			s.fault = nil
			result := reviewApply(t, s, q)
			if result.Review.Revision != 1 {
				t.Fatal(result)
			}
			changed := q
			changed.Action = "open"
			_, e := s.ReviewWrite(LocalController(), changed)
			code(t, e, "idempotency_conflict")
			leaderApply(t, s, leaderRequest(t, s, "message", ""))
			state, _, e := s.load()
			if e != nil || state.Schema != reviewStateSchema {
				t.Fatal(state.Schema, e)
			}
			ev := state.Receipts[q.Spec.EvidenceIDs[0]].Result.Evidence.Request
			ev.OperationID = newEvidenceID(t)
			if _, e = s.AttachEvidence(LocalController(), ev); e != nil {
				t.Fatal(e)
			}
			state, _, e = s.load()
			if e != nil || state.Schema != reviewStateSchema {
				t.Fatal(state.Schema, e)
			}
			state.Schema = evidenceStateSchema
			if e = validateReviewState(state); e == nil {
				t.Fatal("old profile accepted review")
			}
		})
	}
}
func TestReviewUnsentSupersessionAndUncertainRefusal(t *testing.T) {
	s, _, q := reviewFixture(t)
	reviewApply(t, s, q)
	duplicate := q
	duplicate.OperationID = newEvidenceID(t)
	duplicate.ReviewID = duplicate.OperationID
	_, e := s.ReviewWrite(LocalController(), duplicate)
	code(t, e, "conflict")
	reviewApply(t, s, reviewAction(t, s, q.ReviewID, "supersede"))
	reviewApply(t, s, duplicate)
	reviewApply(t, s, reviewAction(t, s, duplicate.ReviewID, "open"))
	_, e = s.ReviewWrite(LocalController(), reviewAction(t, s, duplicate.ReviewID, "supersede"))
	code(t, e, "conflict")
}
