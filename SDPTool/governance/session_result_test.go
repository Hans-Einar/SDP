package governance

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func resultFixture(t *testing.T, root, state string, configure sessionTestConfig) (*Service, SessionWorkRequest) {
	t.Helper()
	s, q := workConfigured(t, root, state, configure, nil)
	done, r := startWork(t, s, q)
	r = reportWork(t, s, r)
	pq := progressRequest(t, s, q, done, r)
	p, e := s.ProgressSession(LocalController(), pq)
	if e != nil {
		t.Fatal(e)
	}
	done, e = s.ResumeSession(LocalController(), workResume(pq, p))
	if e != nil {
		t.Fatal(e)
	}
	pq.OperationID = newTestID()
	pq.ExpectedRevision = done.Revision
	pq.Reason = "Verify exact submitted candidate tests"
	return s, pq
}
func commitResult(t *testing.T, s *Service, kind string, q SessionWorkRequest) SessionMutationResult {
	t.Helper()
	p, e := s.prepareSessionWork(LocalController(), kind, q)
	if e != nil {
		t.Fatal(e)
	}
	done, e := s.ResumeSession(LocalController(), workResume(q, p))
	if e != nil {
		t.Fatal(e)
	}
	return done
}
func sessionReviewRequest(t *testing.T, s *Service, q SessionWorkRequest) ReviewRequest {
	t.Helper()
	st, _, e := s.load()
	if e != nil {
		t.Fatal(e)
	}
	ev := *st.Receipts[q.EvidenceID].Result.Evidence
	r := st.Runs[ev.Request.RunID]
	skill := Source{"Skills/sdp-reviewer/SKILL.md", hash(mustFile(t, filepath.Join(s.root, "Skills/sdp-reviewer/SKILL.md")))}
	spec := ReviewSpec{Instructions: "Review the full selected Session binding independently; synthetic fixture observations", Candidate: ev.Request.Candidate, ReportID: ev.Request.ReportID, EvidenceIDs: []string{q.EvidenceID}, Focus: "All selected whole/contribution outcomes and regressions", HostVersion: "0.164.0", MCPBinary: "/fixture/sdp-mcp", Skills: []SkillInput{{Name: "sdp-reviewer", Path: filepath.Join(s.root, skill.Path), Digest: skill.SHA256}}}
	id := newTestID()
	return ReviewRequest{Schema: ReviewSchema, ProjectRef: q.ProjectRef, CheckoutID: q.CheckoutID, OperationID: id, RunID: r.ID, ReviewID: id, ExpectedRunRevision: r.Revision, Action: "prepare", Spec: &spec}
}
func passingSessionReview(t *testing.T, s *Service, q SessionWorkRequest, done SessionMutationResult) SessionWorkRequest {
	t.Helper()
	req := sessionReviewRequest(t, s, q)
	reviewApply(t, s, req)
	sessionResultReviewFinish(t, s, req.ReviewID, `{"verdict":"passed","summary":"Synthetic independently mapped fixture review","findings":[]}`)
	view, e := s.ReviewRead(req.ReviewID)
	if e != nil || !view.Qualified {
		t.Fatal(view, e)
	}
	q.OperationID = newTestID()
	q.ExpectedRevision = done.Revision
	q.ReviewID = req.ReviewID
	q.ExpectedReviewRevision = view.Review.Revision
	q.EvidenceID = ""
	q.Reason = "Record exact terminal passing observed review"
	return q
}
func TestSessionResultFullBindingAndStages(t *testing.T) {
	s, q := resultFixture(t, "", "", multiCardImplementation)
	early := sessionReviewRequest(t, s, q)
	if _, e := s.ReviewWrite(LocalController(), early); e == nil {
		t.Fatal("review before verification")
	}
	p, e := s.VerifySession(LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = s.ReviewWrite(LocalController(), early); e == nil {
		t.Fatal("review before canonical verification effects")
	}
	v, _ := s.ReadSession(q.InstanceID)
	if v.Steps[0].Stage != "verification" {
		t.Fatal(v)
	}
	done, e := s.ResumeSession(LocalController(), workResume(q, p))
	if e != nil {
		t.Fatal(e)
	}
	retry, e := s.VerifySession(LocalController(), q)
	if e != nil || !reflect.DeepEqual(p, retry) {
		t.Fatal(retry, e)
	}
	rq := passingSessionReview(t, s, q, done)
	st, _, _ := s.load()
	d := st.Reviews[rq.ReviewID]
	var payload struct {
		Session SessionBinding `json:"sessionBinding"`
	}
	json.Unmarshal([]byte(d.Payload), &payload)
	if len(payload.Session.Steps) != 1 || len(payload.Session.Steps[0].Cards) != 2 || payload.Session.Steps[0].Cards[0].Coverage != "contribution" {
		t.Fatal(d.Payload)
	}
	p, e = s.ReviewSession(LocalController(), rq)
	if e != nil {
		t.Fatal(e)
	}
	v, _ = s.ReadSession(q.InstanceID)
	if v.Steps[0].Stage != "review" {
		t.Fatal(v)
	}
	if _, e = s.ResumeSession(LocalController(), workResume(rq, p)); e != nil {
		t.Fatal(e)
	}
	v, e = s.ReadSession(q.InstanceID)
	if e != nil || v.Coverage != "current" || v.Steps[0].Stage != "disposition" || v.Steps[0].ReviewID != rq.ReviewID {
		t.Fatal(v, e)
	}
	review, e := s.ReviewRead(rq.ReviewID)
	if e != nil || !review.Qualified {
		t.Fatal(review, e)
	}
	for _, name := range []string{"#001--Change--fixture.md", "#002--Change--unrelated.md"} {
		raw := mustFile(t, filepath.Join(s.root, "SDP/KanBan/active", name))
		m, _ := metadata(raw)
		if m["CardState"] != "in-progress" || !strings.Contains(string(raw), "PG Session review") {
			t.Fatal(m)
		}
	}
	r, e := s.Read(d.RunID)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = s.Transition(LocalController(), transition(r, "complete")); e == nil {
		t.Fatal("bypassed disposition")
	}
	evidenceFile(t, s, "target.txt", "changed after passing review")
	v, e = s.ReadSession(q.InstanceID)
	if e != nil || len(v.Steps[0].Holds) < 2 {
		t.Fatal(v, e)
	}
}
func TestSessionResultRefusesBadBasis(t *testing.T) {
	for _, mode := range []string{"caller", "revision", "run-revision", "claim", "failed", "contradiction", "artifact", "candidate", "report", "wrong-verb"} {
		t.Run(mode, func(t *testing.T) {
			s, q := resultFixture(t, "", "", nil)
			caller := LocalController()
			switch mode {
			case "caller":
				caller = Caller{}
			case "revision":
				q.ExpectedRevision--
			case "run-revision":
				q.ExpectedRunRevision++
			case "artifact":
				evidenceFile(t, s, "evidence/session-result.log", "changed")
			case "candidate":
				evidenceFile(t, s, "target.txt", "changed")
			case "wrong-verb":
				q.ReviewID = newTestID()
			default:
				st, _, _ := s.load()
				ev := st.Receipts[q.EvidenceID].Result.Evidence.Request
				ev.OperationID = newTestID()
				if mode == "claim" {
					ev.Kind = "claim"
				}
				if mode == "failed" || mode == "contradiction" {
					ev.Verdict = "failed"
				}
				if mode == "report" {
					ev.ReportID = newTestID()
					if _, e := s.AttachEvidence(LocalController(), ev); e == nil {
						t.Fatal("other report")
					}
					return
				}
				if _, e := s.AttachEvidence(LocalController(), ev); e != nil {
					t.Fatal(e)
				}
				if mode != "contradiction" {
					q.EvidenceID = ev.OperationID
				}
			}
			head := mustFile(t, filepath.Join(s.state, "HEAD"))
			ledger := mustFile(t, filepath.Join(s.root, managementLedger))
			if _, e := s.VerifySession(caller, q); e == nil {
				t.Fatal(mode)
			}
			if !bytes.Equal(head, mustFile(t, filepath.Join(s.state, "HEAD"))) || !bytes.Equal(ledger, mustFile(t, filepath.Join(s.root, managementLedger))) {
				t.Fatal("refusal mutated state")
			}
		})
	}
}
func TestSessionResultReviewRefusals(t *testing.T) {
	for _, mode := range []string{"failed", "inconclusive", "unfinished", "review-revision", "candidate", "skill", "test", "caller", "other-review-active"} {
		t.Run(mode, func(t *testing.T) {
			s, q := resultFixture(t, "", "", nil)
			done := commitResult(t, s, "verify", q)
			req := sessionReviewRequest(t, s, q)
			reviewApply(t, s, req)
			answer := `{"verdict":"passed","summary":"Fixture","findings":[]}`
			if mode == "failed" {
				answer = `{"verdict":"failed","summary":"Fixture flaw","findings":["Incorrect result"]}`
			}
			if mode == "inconclusive" {
				answer = `{"verdict":"inconclusive","summary":"Missing coverage","findings":[]}`
			}
			if mode != "unfinished" {
				sessionResultReviewFinish(t, s, req.ReviewID, answer)
			}
			st, _, _ := s.load()
			d := st.Reviews[req.ReviewID]
			q.OperationID = newTestID()
			q.ExpectedRevision = done.Revision
			q.EvidenceID = ""
			q.ReviewID = d.ID
			q.ExpectedReviewRevision = d.Revision
			caller := LocalController()
			switch mode {
			case "caller":
				caller = Caller{}
			case "review-revision":
				q.ExpectedReviewRevision++
			case "candidate":
				evidenceFile(t, s, "target.txt", "changed")
			case "skill":
				evidenceFile(t, s, "Skills/sdp-reviewer/SKILL.md", "changed")
			case "test":
				evidenceFile(t, s, "evidence/session-result.log", "changed")
			case "other-review-active":
				req.OperationID = newTestID()
				req.ReviewID = req.OperationID
				reviewApply(t, s, req)
			}
			head := mustFile(t, filepath.Join(s.state, "HEAD"))
			if _, e := s.ReviewSession(caller, q); e == nil {
				t.Fatal(mode)
			}
			if !bytes.Equal(head, mustFile(t, filepath.Join(s.state, "HEAD"))) {
				t.Fatal("mutation")
			}
		})
	}
}
func TestSessionResultRecoveryEveryEffect(t *testing.T) {
	for _, kind := range []string{"verify", "review"} {
		for _, point := range []string{"before-head", "after-head", "after-claim", "after-file:0", "after-file:1", "after-file:2", "after-file:3", "after-file:4"} {
			t.Run(kind+point, func(t *testing.T) {
				s, q := resultFixture(t, "", "", nil)
				if kind == "review" {
					done := commitResult(t, s, "verify", q)
					q = passingSessionReview(t, s, q, done)
				}
				p, e := s.prepareSessionWork(LocalController(), kind, q)
				if e != nil {
					t.Fatal(e)
				}
				resume := workResume(q, p)
				fault := point
				if strings.HasPrefix(point, "after-") && point != "after-head" {
					fault = "session-" + kind + "-" + point
				}
				hit := false
				s.fault = func(p string) error {
					if p == fault {
						hit = true
						return errors.New("crash")
					}
					return nil
				}
				if _, e = s.ResumeSession(LocalController(), resume); e == nil || !hit {
					t.Fatal("fault not hit", e)
				}
				fresh, e := New(s.root, s.area, s.project, s.state)
				if e != nil {
					t.Fatal(e)
				}
				if _, e = fresh.ResumeSession(LocalController(), resume); e != nil {
					t.Fatal(e)
				}
				if _, e = fresh.ResumeSession(LocalController(), resume); e != nil {
					t.Fatal(e)
				}
				v, e := fresh.ReadSession(q.InstanceID)
				if e != nil || v.Coverage != "current" || v.RecordEffects != "committed" {
					t.Fatal(v, e)
				}
				if _, e = os.Stat(filepath.Join(s.root, managementLedger) + ".pg-pending"); !os.IsNotExist(e) {
					t.Fatal("reservation")
				}
			})
		}
	}
}
func TestSessionResultPolicyAndTamper(t *testing.T) {
	s, q := resultFixture(t, "", "", func(t *testing.T, s *Service, b *SessionBinding) {
		b.Steps[0].Review = "not-required"
		b.Steps[0].OwnerAcceptance = "required"
	})
	commitResult(t, s, "verify", q)
	v, e := s.ReadSession(q.InstanceID)
	if e != nil || v.Steps[0].Stage != "disposition" || !strings.Contains(strings.Join(v.Steps[0].Holds, " "), "owner acceptance required") {
		t.Fatal(v, e)
	}
	st, _, _ := s.load()
	for _, mode := range []string{"profile", "pins", "evidence", "receipt", "predecessor", "review-on-verify"} {
		t.Run(mode, func(t *testing.T) {
			b, _ := json.Marshal(st)
			var bad snapshot
			json.Unmarshal(b, &bad)
			i := bad.Sessions[q.InstanceID]
			switch mode {
			case "profile":
				bad.Schema = sessionWorkStateSchema
			case "pins":
				i.Work[2].Effects.Coverage = i.Work[2].Effects.Coverage[1:]
			case "evidence":
				i.Work[2].Evidence.Request.Verdict = "failed"
			case "receipt":
				delete(bad.Receipts, q.OperationID)
			case "predecessor":
				i.Work[1].Effects.Status = "prepared"
			case "review-on-verify":
				i.Work[2].Review = &Review{}
			}
			bad.Sessions[i.ID] = i
			if validateSessionWork(bad) == nil {
				t.Fatal(mode)
			}
		})
	}
}

func TestSessionResultStrictWireAndPreparedRecovery(t *testing.T) {
	for _, kind := range []string{"verify", "review"} {
		s, q := resultFixture(t, "", "", nil)
		if kind == "review" {
			q = passingSessionReview(t, s, q, commitResult(t, s, "verify", q))
		}
		raw, _ := json.Marshal(q)
		for _, extra := range []string{`,"role":"controller"}`, `,"skills":[]}`, `,"outcome":""}`, `,"OperationID":"` + q.OperationID + `"}`, `,"reason":"duplicate"}`} {
			if _, e := DecodeSessionWork(append(append([]byte{}, raw[:len(raw)-1]...), []byte(extra)...), kind); e == nil {
				t.Fatal("accepted wire", kind, extra)
			}
		}
		wrong := `,"evidenceId":""}`
		if kind == "verify" {
			wrong = `,"reviewId":""}`
		}
		if _, e := DecodeSessionWork(append(append([]byte{}, raw[:len(raw)-1]...), []byte(wrong)...), kind); e == nil {
			t.Fatal("accepted cross verb")
		}
		ledger := mustFile(t, filepath.Join(s.root, managementLedger))
		s.fault = func(point string) error {
			if point == "after-head" {
				return errors.New("lost preparation reply")
			}
			return nil
		}
		if _, e := s.prepareSessionWork(LocalController(), kind, q); e == nil {
			t.Fatal("fault")
		}
		fresh, e := New(s.root, s.area, s.project, s.state)
		if e != nil {
			t.Fatal(e)
		}
		p, e := fresh.prepareSessionWork(LocalController(), kind, q)
		if e != nil {
			t.Fatal(e)
		}
		if !bytes.Equal(ledger, mustFile(t, filepath.Join(s.root, managementLedger))) {
			t.Fatal("prepare published records")
		}
		if _, e = fresh.ResumeSession(LocalController(), workResume(q, p)); e != nil {
			t.Fatal(e)
		}
	}
}
func TestSessionResultPartialPublicationPreservesDrift(t *testing.T) {
	for _, kind := range []string{"verify", "review"} {
		s, q := resultFixture(t, "", "", nil)
		if kind == "review" {
			q = passingSessionReview(t, s, q, commitResult(t, s, "verify", q))
		}
		p, e := s.prepareSessionWork(LocalController(), kind, q)
		if e != nil {
			t.Fatal(e)
		}
		resume := workResume(q, p)
		s.fault = func(point string) error {
			if point == "session-"+kind+"-after-file:0" {
				evidenceFile(t, s, "target.txt", "external change during effects")
			}
			return nil
		}
		if _, e = s.ResumeSession(LocalController(), resume); e == nil {
			t.Fatal("advanced changed candidate")
		}
		if string(mustFile(t, filepath.Join(s.root, "target.txt"))) != "external change during effects" {
			t.Fatal("overwrote external bytes")
		}
		st, _, e := s.load()
		if e != nil {
			t.Fatal(e)
		}
		if sessionLatestEffects(st.Sessions[q.InstanceID]).Status != "prepared" {
			t.Fatal("committed stale result")
		}
	}
}

func sessionResultReviewFinish(t *testing.T, s *Service, id, answer string) {
	t.Helper()
	reviewApply(t, s, reviewAction(t, s, id, "open"))
	q := reviewAction(t, s, id, "bind")
	q.ThreadID = "session-review-" + id
	reviewApply(t, s, q)
	reviewApply(t, s, reviewAction(t, s, id, "send"))
	q = reviewAction(t, s, id, "observe")
	q.ThreadID = "session-review-" + id
	q.TurnID = "session-review-turn"
	q.Status = "completed"
	q.Answer = answer
	reviewApply(t, s, q)
}

func TestSessionResultReviewSnapshotCannotDropContext(t *testing.T) {
	s, q := resultFixture(t, "", "", multiCardImplementation)
	rq := passingSessionReview(t, s, q, commitResult(t, s, "verify", q))
	commitResult(t, s, "review", rq)
	st, _, e := s.load()
	if e != nil {
		t.Fatal(e)
	}
	for _, mode := range []string{"binding", "tests", "review-revision", "verdict"} {
		t.Run(mode, func(t *testing.T) {
			b, _ := json.Marshal(st)
			var bad snapshot
			json.Unmarshal(b, &bad)
			i := bad.Sessions[q.InstanceID]
			w := &i.Work[3]
			d := *w.Review
			switch mode {
			case "binding":
				var payload reviewPayload
				json.Unmarshal([]byte(d.Payload), &payload)
				payload.Session.Steps[0].Cards = payload.Session.Steps[0].Cards[:1]
				b, _ := json.Marshal(payload)
				d.Payload = string(b)
			case "tests":
				w.Tests = nil
			case "review-revision":
				d.Revision++
			case "verdict":
				d.Outcome.Verdict = "failed"
			}
			w.Review = &d
			bad.Reviews[d.ID] = d
			bad.Sessions[i.ID] = i
			if validateSessionWork(bad) == nil {
				t.Fatal(mode)
			}
		})
	}
}
