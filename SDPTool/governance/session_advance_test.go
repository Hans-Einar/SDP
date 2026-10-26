package governance

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func advanceRequest(t *testing.T, s *Service, id, action, step, basis, result string) SessionAdvanceRequest {
	t.Helper()
	st, _, e := s.load()
	if e != nil {
		t.Fatal(e)
	}
	i := st.Sessions[id]
	r := st.Runs[i.Work[0].Request.OperationID]
	return SessionAdvanceRequest{SessionSchema, s.project, st.CheckoutID, newTestID(), id, i.Revision, r.Revision, action, step, basis, result, "Continue within the selected Session and preserve earlier evidence"}
}
func commitAdvance(t *testing.T, s *Service, q SessionAdvanceRequest) SessionMutationResult {
	t.Helper()
	p, e := s.AdvanceSession(LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	resume := SessionResumeRequest{SessionSchema, q.ProjectRef, q.CheckoutID, newTestID(), q.InstanceID, p.Revision, "Commit exact continuation effects"}
	done, e := s.ResumeSession(LocalController(), resume)
	if e != nil {
		t.Fatal(e)
	}
	again, e := s.ResumeSession(LocalController(), resume)
	if e != nil || !reflect.DeepEqual(done, again) {
		t.Fatal(again, e)
	}
	if _, e = s.AdvanceSession(LocalController(), q); e != nil {
		t.Fatal(e)
	}
	return done
}
func completeFixtureDisposition(t *testing.T, s *Service, d SessionDisposeRequest) {
	t.Helper()
	p, e := s.DisposeSession(LocalController(), d)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = s.ResumeSession(LocalController(), disposeResume(d, p)); e != nil {
		t.Fatal(e)
	}
}
func TestSessionAdvanceCompletesWholePlan(t *testing.T) {
	s, d := dispositionFixture(t, "", "", nil)
	completeFixtureDisposition(t, s, d)
	q := advanceRequest(t, s, d.InstanceID, "complete", "", "disposition", d.OperationID)
	commitAdvance(t, s, q)
	st, _, e := s.load()
	if e != nil {
		t.Fatal(e)
	}
	i := st.Sessions[q.InstanceID]
	m, _ := metadata(mustFile(t, filepath.Join(s.root, i.Context.Binding.Plan.Path)))
	if m["state"] != "completed" {
		t.Fatal(m)
	}
	m, _ = metadata(mustFile(t, filepath.Join(s.root, i.Context.Binding.Session.Path)))
	if m["Status"] != "completed" {
		t.Fatal(m)
	}
	if _, e = s.publicationFiles(); e != nil {
		t.Fatal(e)
	}
	if _, e = s.ReadSession(q.InstanceID); e != nil {
		t.Fatal(e)
	}
}
func TestSessionAdvanceNextStep(t *testing.T) {
	for _, shared := range []bool{true, false} {
		t.Run(fmt.Sprint("shared=", shared), func(t *testing.T) {
			s, d := dispositionFixture(t, "", "", func(t *testing.T, s *Service, b *SessionBinding) {
				multiCardImplementation(t, s, b)
				b.Steps[0].Cards = b.Steps[0].Cards[:1]
				next := b.Steps[0]
				next.ID = "second"
				next.Cards = []SessionCardRef{{"KB-SDP-001", "whole"}, {"KB-SDP-002", "whole"}}
				if !shared {
					b.Steps[0].Cards[0].Coverage = "whole"
					next.Cards = next.Cards[1:]
				}
				next.DependsOn = []string{b.Steps[0].ID}
				b.Steps = append(b.Steps, next)
				evidenceFile(t, s, b.Session.Path, string(mustFile(t, filepath.Join(s.root, b.Session.Path)))+"\n## second\n\nComplete remaining selected scope.\n")
			})
			m, _ := metadata(mustFile(t, filepath.Join(s.root, "SDP/KanBan/active/#002--Change--unrelated.md")))
			if m["CardState"] != "ready" {
				t.Fatal("future bound card not ready", m)
			}
			completeFixtureDisposition(t, s, d)
			q := advanceRequest(t, s, d.InstanceID, "next", "second", "disposition", d.OperationID)
			done := commitAdvance(t, s, q)
			st, _, e := s.load()
			if e != nil {
				t.Fatal(e)
			}
			i := st.Sessions[q.InstanceID]
			if len(i.Past) != 1 || i.Selection != nil || len(i.Work) != 0 {
				t.Fatal("execution not archived")
			}
			oldrun := i.Past[0].Execution.Work[0].Request.OperationID
			if _, ok := sessionRunBinding(st, oldrun); !ok {
				t.Fatal("lost historical run")
			}
			selectQ := SessionSelectRequest{SessionSchema, q.ProjectRef, q.CheckoutID, newTestID(), q.InstanceID, done.Revision, "second", "Continue the same plan after its dependency completed"}
			p, e := s.SelectSession(LocalController(), selectQ)
			if e != nil {
				t.Fatal(e)
			}
			if _, e = s.ResumeSession(LocalController(), SessionResumeRequest{SessionSchema, q.ProjectRef, q.CheckoutID, newTestID(), q.InstanceID, p.Revision, "Commit next-step selection"}); e != nil {
				t.Fatal(e)
			}
			st, _, e = s.load()
			if e != nil {
				t.Fatal(e)
			}
			if selectedSessionStep(st.Sessions[q.InstanceID]).ID != "second" {
				t.Fatal("wrong step")
			}
			if _, e = s.RunSession(oldrun); e != nil {
				t.Fatal(e)
			}
			d2 := finishSelectedStep(t, s, q.InstanceID)
			commitAdvance(t, s, advanceRequest(t, s, q.InstanceID, "complete", "", "disposition", d2.OperationID))
			v, e := s.ReadSession(q.InstanceID)
			if e != nil || v.State != "completed" || len(v.History) != 1 || v.Steps[0].Stage != "done" || v.Steps[1].Stage != "done" {
				t.Fatal(v, e)
			}
		})
	}
}

func TestSessionAdvanceOwnerRework(t *testing.T) {
	s, g, _ := ownerFixture(t)
	done := commitOwner(t, s, LocalController(), g)
	decision := ownerDecision(g, done.Revision, "changes-requested")
	commitOwner(t, s, LocalOwner(), decision)
	st, _, _ := s.load()
	step := selectedSessionStep(st.Sessions[g.InstanceID]).ID
	q := advanceRequest(t, s, g.InstanceID, "rework", step, "owner", decision.OperationID)
	done = commitAdvance(t, s, q)
	st, _, e := s.load()
	if e != nil {
		t.Fatal(e)
	}
	i := st.Sessions[g.InstanceID]
	old := i.Past[0].Execution
	if st.Runs[old.Work[0].Request.OperationID].State != "canceled" {
		t.Fatal("old run not retired")
	}
	if _, e = os.Stat(filepath.Join(s.root, dispositionClaim(old))); !os.IsNotExist(e) {
		t.Fatal("claim retained")
	}
	sq := SessionSelectRequest{SessionSchema, q.ProjectRef, q.CheckoutID, newTestID(), q.InstanceID, done.Revision, step, "Reassess the same selected scope after owner feedback"}
	p, e := s.SelectSession(LocalController(), sq)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = s.ResumeSession(LocalController(), SessionResumeRequest{SessionSchema, q.ProjectRef, q.CheckoutID, newTestID(), q.InstanceID, p.Revision, "Commit rework selection"}); e != nil {
		t.Fatal(e)
	}
	b := string(mustFile(t, filepath.Join(s.root, "SDP/KanBan/active/#001--Change--fixture.md")))
	if !strings.Contains(b, "Prior implementation") || !strings.Contains(b, "changes-requested") {
		t.Fatal("lost rework history")
	}
	finishSelectedStep(t, s, q.InstanceID)
}

// Execute a newly selected step with fresh assessment, run, report and review.
// Earlier evidence is reused only as an input to the assessment, never as a result.
func finishSelectedStep(t *testing.T, s *Service, id string) SessionDisposeRequest {
	t.Helper()
	st, _, e := s.load()
	if e != nil {
		t.Fatal(e)
	}
	i := st.Sessions[id]
	old := i.Past[len(i.Past)-1].Execution
	a := old.Assessments[len(old.Assessments)-1].Request
	a.OperationID = newTestID()
	a.ExpectedRevision = i.Revision
	a.StepID = selectedSessionStep(i).ID
	done := commitAssessment(t, s, a)
	q := old.Work[0].Request
	q.OperationID = newTestID()
	q.ExpectedRevision = done.Revision
	q.AssessmentID = a.OperationID
	done, r := startWork(t, s, q)
	r = reportWork(t, s, r)
	pq := progressRequest(t, s, q, done, r)
	done = commitResult(t, s, "progress", pq)
	pq.OperationID = newTestID()
	pq.ExpectedRevision = done.Revision
	done = commitResult(t, s, "verify", pq)
	if selectedSessionStep(i).Review == "independent" {
		pq = passingSessionReview(t, s, pq, done)
		done = commitResult(t, s, "review", pq)
	}
	st, _, e = s.load()
	if e != nil {
		t.Fatal(e)
	}
	i = st.Sessions[id]
	d := SessionDisposeRequest{Schema: SessionSchema, ProjectRef: s.project, CheckoutID: st.CheckoutID, OperationID: newTestID(), InstanceID: id, ExpectedRevision: done.Revision, ExpectedRunRevision: st.Runs[q.OperationID].Revision, ResultID: pq.OperationID, System: "SDPTOOL", Reason: "Complete freshly assessed and verified next execution"}
	for _, c := range selectedSessionStep(i).Cards {
		cd := SessionCardDisposition{ID: c.ID, Outcome: "completed", Summary: "Fresh execution verified"}
		if c.Coverage == "contribution" {
			cd.Outcome = "contribution-completed"
			cd.Remaining = "Other contributions remain"
		}
		d.Cards = append(d.Cards, cd)
	}
	if selectedSessionStep(i).OwnerAcceptance == "required" {
		op := newTestID()
		g := SessionOwnerRequest{Schema: SessionSchema, ProjectRef: s.project, CheckoutID: st.CheckoutID, OperationID: op, InstanceID: id, ExpectedRevision: done.Revision, ExpectedRunRevision: d.ExpectedRunRevision, Action: "request", GateID: op, ResultID: d.ResultID, Reason: "Request fresh synthetic owner decision"}
		done = commitOwner(t, s, LocalController(), g)
		done = commitOwner(t, s, LocalOwner(), ownerDecision(g, done.Revision, "accept"))
		d.ExpectedRevision = done.Revision
	}
	completeFixtureDisposition(t, s, d)
	return d
}

func TestSessionAdvanceFailedResultRework(t *testing.T) {
	for _, basis := range []string{"test", "review"} {
		t.Run(basis, func(t *testing.T) {
			s, q := resultFixture(t, "", "", multiCardImplementation)
			result := ""
			if basis == "test" {
				st, _, _ := s.load()
				ev := st.Receipts[q.EvidenceID].Result.Evidence.Request
				ev.OperationID = newTestID()
				ev.Verdict = "failed"
				evidenceFile(t, s, "evidence/failure.log", "Synthetic failing test for the exact submitted candidate.\n")
				ev.Artifact = Source{"evidence/failure.log", hash(mustFile(t, filepath.Join(s.root, "evidence/failure.log")))}
				if _, e := s.AttachEvidence(LocalController(), ev); e != nil {
					t.Fatal(e)
				}
				result = ev.OperationID
			} else {
				commitResult(t, s, "verify", q)
				rq := sessionReviewRequest(t, s, q)
				reviewApply(t, s, rq)
				sessionResultReviewFinish(t, s, rq.ReviewID, `{"verdict":"failed","summary":"Synthetic defect","findings":["Correction required"]}`)
				result = rq.ReviewID
			}
			aq := advanceRequest(t, s, q.InstanceID, "rework", "S1", basis, result)
			done := commitAdvance(t, s, aq)
			st, _, e := s.load()
			if e != nil {
				t.Fatal(e)
			}
			old := st.Sessions[q.InstanceID].Past[0].Execution
			if st.Runs[old.Work[0].Request.OperationID].State != "canceled" {
				t.Fatal("failure run not retired")
			}
			sq := SessionSelectRequest{SessionSchema, aq.ProjectRef, aq.CheckoutID, newTestID(), aq.InstanceID, done.Revision, "S1", "Repair the checked failed result"}
			p, e := s.SelectSession(LocalController(), sq)
			if e != nil {
				t.Fatal(e)
			}
			if _, e = s.ResumeSession(LocalController(), selectionResume(sq, p)); e != nil {
				t.Fatal(e)
			}
			finishSelectedStep(t, s, aq.InstanceID)
		})
	}
}

func TestSessionAdvanceRecovery(t *testing.T) {
	for _, action := range []string{"rework", "complete"} {
		points := []string{"after-claim", "after-file:0", "after-file:1", "after-file:2", "after-file:3", "after-file:4", "before-head", "after-head"}
		if action == "rework" {
			points = append(points, "after-file:5", "after-file:6")
		}
		for _, point := range points {
			t.Run(action+point, func(t *testing.T) {
				var s *Service
				var q SessionAdvanceRequest
				if action == "complete" {
					svc, d := dispositionFixture(t, "", "", nil)
					s = svc
					completeFixtureDisposition(t, s, d)
					q = advanceRequest(t, s, d.InstanceID, action, "", "disposition", d.OperationID)
				} else {
					svc, g, _ := ownerFixture(t)
					s = svc
					done := commitOwner(t, s, LocalController(), g)
					d := ownerDecision(g, done.Revision, "changes-requested")
					commitOwner(t, s, LocalOwner(), d)
					q = advanceRequest(t, s, d.InstanceID, action, "S1", "owner", d.OperationID)
				}
				p, e := s.AdvanceSession(LocalController(), q)
				if e != nil {
					t.Fatal(e)
				}
				r := SessionResumeRequest{SessionSchema, q.ProjectRef, q.CheckoutID, newTestID(), q.InstanceID, p.Revision, "Recover exact continuation after interrupted publication"}
				fault := point
				if strings.HasPrefix(point, "after-") && point != "after-head" {
					fault = "session-advance-" + point
				}
				hit := false
				s.fault = func(p string) error {
					if p == fault {
						hit = true
						return errors.New("interrupted")
					}
					return nil
				}
				if _, e = s.ResumeSession(LocalController(), r); e == nil || !hit {
					t.Fatal("fault not exercised", e)
				}
				fresh, e := New(s.root, s.area, s.project, s.state)
				if e != nil {
					t.Fatal(e)
				}
				got, e := fresh.ResumeSession(LocalController(), r)
				if e != nil {
					t.Fatal(e)
				}
				again, e := fresh.ResumeSession(LocalController(), r)
				if e != nil || got != again {
					t.Fatal(again, e)
				}
				v, e := fresh.ReadSession(q.InstanceID)
				if e != nil || v.Coverage != "current" {
					t.Fatal(v, e)
				}
				if _, e = os.Stat(filepath.Join(s.root, managementLedger) + ".pg-pending"); !os.IsNotExist(e) {
					t.Fatal("reservation retained")
				}
				if _, e = fresh.publicationFiles(); e != nil {
					t.Fatal(e)
				}
			})
		}
	}
}

func TestSessionAdvanceRetainedTamper(t *testing.T) {
	s, d := dispositionFixture(t, "", "", nil)
	completeFixtureDisposition(t, s, d)
	q := advanceRequest(t, s, d.InstanceID, "complete", "", "disposition", d.OperationID)
	commitAdvance(t, s, q)
	st, _, e := s.load()
	if e != nil {
		t.Fatal(e)
	}
	raw, _ := json.Marshal(st)
	for _, mode := range []string{"basis", "digest", "preimage", "postimage", "receipt", "coverage"} {
		t.Run(mode, func(t *testing.T) {
			var bad snapshot
			json.Unmarshal(raw, &bad)
			i := bad.Sessions[q.InstanceID]
			a := i.Advance
			switch mode {
			case "basis":
				a.Request.ResultID = newTestID()
			case "digest":
				a.PriorDigest = strings.Repeat("0", 64)
			case "preimage":
				a.Before[managementLedger] += "\n"
			case "postimage":
				a.Effects.Changes[0].Content += "corrupt"
			case "receipt":
				delete(bad.Receipts, q.OperationID)
			case "coverage":
				a.Effects.Coverage = a.Effects.Coverage[1:]
			}
			bad.Sessions[i.ID] = i
			if validateSessionAdvances(bad) == nil {
				t.Fatal("tamper accepted", mode)
			}
		})
	}
}

func TestSessionAdvanceExternalEditAndRemainingScope(t *testing.T) {
	t.Run("external-edit", func(t *testing.T) {
		s, d := dispositionFixture(t, "", "", nil)
		completeFixtureDisposition(t, s, d)
		q := advanceRequest(t, s, d.InstanceID, "complete", "", "disposition", d.OperationID)
		p, e := s.AdvanceSession(LocalController(), q)
		if e != nil {
			t.Fatal(e)
		}
		r := SessionResumeRequest{SessionSchema, q.ProjectRef, q.CheckoutID, newTestID(), q.InstanceID, p.Revision, "Recover interrupted closure"}
		s.fault = func(p string) error {
			if p == "session-advance-after-file:0" {
				return errors.New("interrupted")
			}
			return nil
		}
		if _, e = s.ResumeSession(LocalController(), r); e == nil {
			t.Fatal("interruption absent")
		}
		s.fault = nil
		path := filepath.Join(s.root, "SDP/Sessions/session.md")
		original := mustFile(t, path)
		edited := append(original, []byte("\nExternal owner edit; preserve.\n")...)
		if e = os.WriteFile(path, edited, 0600); e != nil {
			t.Fatal(e)
		}
		head := mustFile(t, filepath.Join(s.state, "HEAD"))
		if _, e = s.ResumeSession(LocalController(), r); e == nil {
			t.Fatal("external edit overwritten")
		}
		if !bytes.Equal(edited, mustFile(t, path)) || !bytes.Equal(head, mustFile(t, filepath.Join(s.state, "HEAD"))) {
			t.Fatal("refusal mutated external state")
		}
		// An explicit test restoration supplies the exact pinned preimage; no force overwrite.
		if e = os.WriteFile(path, original, 0600); e != nil {
			t.Fatal(e)
		}
		if _, e = s.ResumeSession(LocalController(), r); e != nil {
			t.Fatal(e)
		}
	})
	t.Run("remaining-contribution", func(t *testing.T) {
		s, d := dispositionFixture(t, "", "", multiCardImplementation)
		completeFixtureDisposition(t, s, d)
		q := advanceRequest(t, s, d.InstanceID, "complete", "", "disposition", d.OperationID)
		head := mustFile(t, filepath.Join(s.state, "HEAD"))
		if _, e := s.AdvanceSession(LocalController(), q); e == nil {
			t.Fatal("unfinished broader card closed")
		}
		if !bytes.Equal(head, mustFile(t, filepath.Join(s.state, "HEAD"))) {
			t.Fatal("refusal wrote state")
		}
	})
}
