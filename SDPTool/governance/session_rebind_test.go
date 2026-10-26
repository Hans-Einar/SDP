package governance

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func rebindFixtureAt(t *testing.T, root, state string) (*Service, string) {
	t.Helper()
	s, _, q := sessionFixtureAt(t, root, state)
	st, _, e := s.load()
	if e != nil {
		t.Fatal(e)
	}
	for _, r := range st.Runs {
		if r.State != "completed" && r.State != "canceled" {
			if _, e = s.Transition(LocalController(), transition(Result{ProjectRef: s.project, CheckoutID: st.CheckoutID, Run: r}, "cancel")); e != nil {
				t.Fatal(e)
			}
		}
	}
	p, e := s.AdoptSession(LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = s.ResumeSession(LocalController(), sessionResume(p, q)); e != nil {
		t.Fatal(e)
	}
	return s, p.InstanceID
}
func editRebindPlan(t *testing.T, s *Service, id string, edit func(*SessionBinding)) {
	t.Helper()
	st, _, e := s.load()
	if e != nil {
		t.Fatal(e)
	}
	i := st.Sessions[id]
	b := i.Context.Binding
	if edit != nil {
		edit(&b)
	}
	raw, _ := json.Marshal(b)
	evidenceFile(t, s, i.Context.BindingPath, string(raw))
	p := b.Plan.Path
	evidenceFile(t, s, p, string(mustFile(t, filepath.Join(s.root, p)))+"\nReviewed clarification of remaining work.\n")
}
func rebindRequest(t *testing.T, s *Service, id string) SessionRebindRequest {
	t.Helper()
	c, e := s.SessionRebindContext(SessionRebindContextRequest{id})
	if e != nil {
		t.Fatal(e)
	}
	q := SessionRebindRequest{SessionSchema, s.project, c.Next.CheckoutID, newTestID(), id, c.ExpectedRevision, c.Next.InputRevision, Source{}, "Replan remaining work and retain earlier completed evidence"}
	review := SessionRebindReview{"sdp-session-rebind-review/0.1", id, c.ExpectedRevision, c.OldBindingDigest, c.Next.InputRevision, "fixture controller reviewer", "passed", "Synthetic attributed review of exact pending binding; no native review or owner decision."}
	b, _ := json.Marshal(review)
	q.Review = Source{"evidence/rebind-review-" + q.OperationID + ".json", hash(b)}
	evidenceFile(t, s, q.Review.Path, string(b))
	return q
}
func commitRebind(t *testing.T, s *Service, q SessionRebindRequest) SessionMutationResult {
	t.Helper()
	p, e := s.RebindSession(LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	fresh, e := New(s.root, s.area, s.project, s.state)
	if e != nil {
		t.Fatal(e)
	}
	resume := SessionResumeRequest{SessionSchema, q.ProjectRef, q.CheckoutID, newTestID(), q.InstanceID, p.Revision, "Resume reviewed binding revision"}
	done, e := fresh.ResumeSession(LocalController(), resume)
	if e != nil {
		t.Fatal(e)
	}
	again, e := s.ResumeSession(LocalController(), resume)
	if e != nil || !reflect.DeepEqual(done, again) {
		t.Fatal(again, e)
	}
	again, e = s.RebindSession(LocalController(), q)
	if e != nil || !reflect.DeepEqual(p, again) {
		t.Fatal(again, e)
	}
	return done
}
func TestSessionRebindBeforeSelectionAndRepeatedEpochs(t *testing.T) {
	s, id := rebindFixtureAt(t, "", "")
	before, _, _ := s.load()
	registration := before.Sessions[id].Registration
	for n := 0; n < 2; n++ {
		editRebindPlan(t, s, id, func(b *SessionBinding) {
			b.Steps[0].WritablePaths = append(b.Steps[0].WritablePaths, fmt.Sprintf("extra%d.txt", n))
		})
		q := rebindRequest(t, s, id)
		done := commitRebind(t, s, q)
		st, _, e := s.load()
		if e != nil {
			t.Fatal(e)
		}
		i := st.Sessions[id]
		if len(i.Rebindings) != n+1 || i.Selection != nil || len(i.Assessments) != 0 || i.Revision != q.ExpectedRevision+2 || done.Revision != i.Revision || !reflect.DeepEqual(i.Registration, registration) {
			t.Fatal("binding epoch not retained", i)
		}
		v, e := s.ReadSession(id)
		if e != nil || v.Coverage != "current" || len(v.BindingRevisions) != n+1 {
			t.Fatal(v, e)
		}
	}
	st, _, _ := s.load()
	i := st.Sessions[id]
	sq := SessionSelectRequest{SessionSchema, s.project, st.CheckoutID, newTestID(), id, i.Revision, "S1", "Fresh selection after reviewed rebind"}
	p, e := s.SelectSession(LocalController(), sq)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = s.ResumeSession(LocalController(), SessionResumeRequest{SessionSchema, s.project, st.CheckoutID, newTestID(), id, p.Revision, "Commit fresh selection"}); e != nil {
		t.Fatal(e)
	}
	if _, e = s.SessionRebindContext(SessionRebindContextRequest{id}); e == nil {
		t.Fatal("unfinished selection allowed rebind")
	}
}
func rebindTwoSteps(t *testing.T, s *Service, b *SessionBinding) {
	b.Steps[0].Cards[0].Coverage = "contribution"
	next := b.Steps[0]
	next.ID = "second"
	next.Cards = []SessionCardRef{{"KB-SDP-001", "whole"}}
	next.DependsOn = []string{b.Steps[0].ID}
	b.Steps = append(b.Steps, next)
	evidenceFile(t, s, b.Session.Path, string(mustFile(t, filepath.Join(s.root, b.Session.Path)))+"\n## second\n\nComplete remaining work.\n")
}
func TestSessionRebindRetainsCompletedWorkAndContinues(t *testing.T) {
	for _, advance := range []bool{false, true} {
		t.Run(fmt.Sprint("afterNext=", advance), func(t *testing.T) {
			s, d := dispositionFixture(t, "", "", rebindTwoSteps)
			completeFixtureDisposition(t, s, d)
			st, _, _ := s.load()
			oldRun := st.Sessions[d.InstanceID].Work[0].Request.OperationID
			if advance {
				commitAdvance(t, s, advanceRequest(t, s, d.InstanceID, "next", "second", "disposition", d.OperationID))
			}
			editRebindPlan(t, s, d.InstanceID, func(b *SessionBinding) { b.Steps[1].WritablePaths = append(b.Steps[1].WritablePaths, "extra.txt") })
			q := rebindRequest(t, s, d.InstanceID)
			done := commitRebind(t, s, q)
			v, e := s.ReadSession(d.InstanceID)
			if e != nil || v.Steps[0].Stage != "done" || v.Steps[0].RunID != oldRun {
				t.Fatal(v, e)
			}
			if _, e = s.RunSession(oldRun); e != nil {
				t.Fatal(e)
			}
			sq := SessionSelectRequest{SessionSchema, q.ProjectRef, q.CheckoutID, newTestID(), q.InstanceID, done.Revision, "second", "Continue remaining step with fresh assessment"}
			p, e := s.SelectSession(LocalController(), sq)
			if e != nil {
				t.Fatal(e)
			}
			if _, e = s.ResumeSession(LocalController(), SessionResumeRequest{SessionSchema, q.ProjectRef, q.CheckoutID, newTestID(), q.InstanceID, p.Revision, "Commit remaining selection"}); e != nil {
				t.Fatal(e)
			}
			d2 := finishSelectedStep(t, s, q.InstanceID)
			commitAdvance(t, s, advanceRequest(t, s, q.InstanceID, "complete", "", "disposition", d2.OperationID))
			v, e = s.ReadSession(q.InstanceID)
			if e != nil || v.State != "completed" || v.Steps[0].Stage != "done" || v.Steps[1].Stage != "done" {
				t.Fatal(v, e)
			}
		})
	}
}
func TestSessionRebindGuards(t *testing.T) {
	for _, kind := range []string{"review-failed", "review-stale", "review-bytes", "stale-context", "instructions", "journal", "policy", "protected", "missing-step", "wrong-caller"} {
		t.Run(kind, func(t *testing.T) {
			s, id := rebindFixtureAt(t, "", "")
			editRebindPlan(t, s, id, nil)
			q := rebindRequest(t, s, id)
			switch kind {
			case "review-failed", "review-stale":
				var r SessionRebindReview
				json.Unmarshal(mustFile(t, filepath.Join(s.root, q.Review.Path)), &r)
				if kind == "review-failed" {
					r.Verdict = "failed"
				} else {
					r.ExpectedRevision++
				}
				b, _ := json.Marshal(r)
				evidenceFile(t, s, q.Review.Path, string(b))
				q.Review.SHA256 = hash(b)
			case "review-bytes":
				evidenceFile(t, s, q.Review.Path, "{}")
			case "stale-context":
				editRebindPlan(t, s, id, nil)
			case "instructions":
				evidenceFile(t, s, "AGENTS.md", "Changed governance instructions.\n")
			case "journal":
				st, _, _ := s.load()
				p := st.Sessions[id].Context.Binding.Session.Path
				b := mustFile(t, filepath.Join(s.root, p))
				evidenceFile(t, s, p, strings.ReplaceAll(string(b), "This records adoption only", "Modified original journal"))
			case "policy", "protected", "missing-step":
				editRebindPlan(t, s, id, func(b *SessionBinding) {
					switch kind {
					case "policy":
						b.Steps[0].Review = "not-required"
					case "protected":
						b.Steps[0].ProtectedPaths = []string{}
					case "missing-step":
						b.Steps[0].ID = "replaced"
					}
				})
			}
			if kind == "instructions" || kind == "journal" {
				q = rebindRequest(t, s, id)
			}
			head := mustFile(t, filepath.Join(s.state, "HEAD"))
			caller := LocalController()
			if kind == "wrong-caller" {
				caller = Caller{}
			}
			if _, e := s.RebindSession(caller, q); e == nil {
				t.Fatal("unsafe rebind accepted")
			}
			if !bytes.Equal(head, mustFile(t, filepath.Join(s.state, "HEAD"))) {
				t.Fatal("refusal changed state")
			}
		})
	}
}
func TestSessionRebindCrashRecovery(t *testing.T) {
	for _, point := range []string{"after-claim", "after-file:0", "after-file:1", "after-file:2", "after-file:3", "after-file:4", "before-head", "after-head"} {
		t.Run(point, func(t *testing.T) {
			s, id := rebindFixtureAt(t, "", "")
			editRebindPlan(t, s, id, nil)
			q := rebindRequest(t, s, id)
			p, e := s.RebindSession(LocalController(), q)
			if e != nil {
				t.Fatal(e)
			}
			resume := SessionResumeRequest{SessionSchema, q.ProjectRef, q.CheckoutID, newTestID(), id, p.Revision, "Recover interrupted binding revision"}
			fault := point
			if strings.HasPrefix(point, "after-") && point != "after-head" {
				fault = "session-rebind-" + point
			}
			fired := false
			s.fault = func(at string) error {
				if at == fault && !fired {
					fired = true
					return errors.New("simulated interruption")
				}
				return nil
			}
			if _, e = s.ResumeSession(LocalController(), resume); e == nil || !fired {
				t.Fatal("fault absent", e)
			}
			fresh, e := New(s.root, s.area, s.project, s.state)
			if e != nil {
				t.Fatal(e)
			}
			done, e := fresh.ResumeSession(LocalController(), resume)
			if e != nil {
				t.Fatal(e)
			}
			again, e := fresh.ResumeSession(LocalController(), resume)
			if e != nil || !reflect.DeepEqual(done, again) {
				t.Fatal(again, e)
			}
			st, _, e := fresh.load()
			if e != nil || len(st.Sessions[id].Rebindings) != 1 {
				t.Fatal(e)
			}
			for _, path := range []string{st.Sessions[id].Context.Binding.Session.Path, st.Sessions[id].Context.Binding.Plan.Path} {
				if strings.Count(string(mustFile(t, filepath.Join(s.root, path))), "## PG Session rebind "+q.OperationID) != 1 {
					t.Fatal("duplicate journal")
				}
			}
		})
	}
}

func TestSessionRebindRetainedTamper(t *testing.T) {
	s, id := rebindFixtureAt(t, "", "")
	editRebindPlan(t, s, id, nil)
	q := rebindRequest(t, s, id)
	commitRebind(t, s, q)
	st, _, e := s.load()
	if e != nil {
		t.Fatal(e)
	}
	raw, _ := json.Marshal(st)
	for _, kind := range []string{"profile", "context", "digest", "preimage", "postimage", "receipt", "coverage", "recursive", "original", "source-set", "records", "epoch-revision", "review"} {
		t.Run(kind, func(t *testing.T) {
			var bad snapshot
			json.Unmarshal(raw, &bad)
			i := bad.Sessions[id]
			a := &i.Rebindings[0].Transition
			switch kind {
			case "profile":
				bad.Schema = sessionTransferStateSchema
			case "context":
				i.Context.Binding.Steps[0].WritablePaths = []string{"forged.txt"}
			case "digest":
				a.PriorDigest = strings.Repeat("0", 64)
			case "preimage":
				a.Before[i.Context.Binding.Plan.Path] += "Changed retained input.\n"
			case "postimage":
				a.Effects.Changes[0].Content += "Forged postimage.\n"
			case "receipt":
				delete(bad.Receipts, q.OperationID)
			case "coverage":
				a.Effects.Coverage = a.Effects.Coverage[1:]
			case "recursive":
				i.Rebindings[0].Before.Rebindings = []SessionRebindEpoch{{}}
			case "original":
				i.Rebindings[0].Before.Registration.Content += "Rewritten adoption.\n"
			case "source-set":
				a.Next.Sources = a.Next.Sources[1:]
			case "records":
				a.Records = a.Records[1:]
			case "epoch-revision":
				i.Rebindings[0].Before.Revision++
			case "review":
				a.Request.Review.SHA256 = strings.Repeat("0", 64)
			}
			bad.Sessions[id] = i
			if validateSessionRebinds(bad) == nil {
				t.Fatal("retained tamper accepted")
			}
		})
	}
}
func TestSessionRebindPreservesBoundary(t *testing.T) {
	for _, kind := range []string{"completed-step", "review-policy", "owner-policy", "protected-path", "removed-card", "old-artifact", "old-journal", "old-rebind-journal"} {
		t.Run(kind, func(t *testing.T) {
			s, d := dispositionFixture(t, "", "", rebindTwoSteps)
			completeFixtureDisposition(t, s, d)
			if kind == "old-rebind-journal" {
				editRebindPlan(t, s, d.InstanceID, nil)
				commitRebind(t, s, rebindRequest(t, s, d.InstanceID))
			}
			editRebindPlan(t, s, d.InstanceID, func(b *SessionBinding) {
				switch kind {
				case "completed-step":
					b.Steps[0].WritablePaths = append(b.Steps[0].WritablePaths, "other.txt")
				case "review-policy":
					b.Steps[1].Review = "not-required"
				case "protected-path":
					b.Steps[1].ProtectedPaths = []string{}
				case "removed-card":
					b.Steps[1].Cards = []SessionCardRef{{"KB-SDP-002", "whole"}}
				}
			})
			if kind == "owner-policy" { // Exact typed policy check independent of fixture owner workflow.
				st, _, _ := s.load()
				i := st.Sessions[d.InstanceID]
				i.Context.Binding.Steps[1].OwnerAcceptance = "required"
				var next SessionContext
				b, _ := json.Marshal(i.Context)
				json.Unmarshal(b, &next)
				next.Binding.Steps[1].OwnerAcceptance = "not-required"
				if validateRebindBinding(i, next) == nil {
					t.Fatal("weakened owner policy")
				}
				return
			}
			if kind == "old-artifact" || kind == "old-journal" || kind == "old-rebind-journal" {
				st, _, _ := s.load()
				i := st.Sessions[d.InstanceID]
				if kind == "old-artifact" {
					p := dispositionPath(SessionDisposition{Request: d})
					evidenceFile(t, s, p, string(mustFile(t, filepath.Join(s.root, p)))+"\nChanged prior evidence.\n")
				} else {
					p := i.Context.Binding.Plan.Path
					text := string(mustFile(t, filepath.Join(s.root, p)))
					if kind == "old-journal" {
						text = strings.ReplaceAll(text, "## PG ", "## Changed ")
					} else {
						text = strings.ReplaceAll(text, "## PG Session rebind ", "## Changed rebind ")
					}
					evidenceFile(t, s, p, text)
				}
				q := rebindRequest(t, s, d.InstanceID)
				if _, e := s.RebindSession(LocalController(), q); e == nil {
					t.Fatal("changed historical bytes accepted")
				}
				return
			}
			if _, e := s.SessionRebindContext(SessionRebindContextRequest{d.InstanceID}); e == nil {
				t.Fatal("unsafe binding accepted")
			}
		})
	}
}
func TestSessionRebindPreparedEditIsNotOverwritten(t *testing.T) {
	s, id := rebindFixtureAt(t, "", "")
	editRebindPlan(t, s, id, nil)
	q := rebindRequest(t, s, id)
	p, e := s.RebindSession(LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	st, _, _ := s.load()
	path := st.Sessions[id].Context.Binding.Plan.Path
	original := mustFile(t, filepath.Join(s.root, path))
	edited := string(original) + "\nExternal owner correction after review.\n"
	evidenceFile(t, s, path, edited)
	head := mustFile(t, filepath.Join(s.state, "HEAD"))
	resume := SessionResumeRequest{SessionSchema, q.ProjectRef, q.CheckoutID, newTestID(), id, p.Revision, "Resume exact reviewed rebind"}
	if _, e = s.ResumeSession(LocalController(), resume); e == nil {
		t.Fatal("stale preparation committed")
	}
	if string(mustFile(t, filepath.Join(s.root, path))) != edited || !bytes.Equal(head, mustFile(t, filepath.Join(s.state, "HEAD"))) {
		t.Fatal("external change overwritten")
	}
	evidenceFile(t, s, path, string(original))
	if _, e = s.ResumeSession(LocalController(), resume); e != nil {
		t.Fatal(e)
	}
}

func TestSessionRebindWithCompletedCard(t *testing.T) {
	s, d := dispositionFixture(t, "", "", func(t *testing.T, s *Service, b *SessionBinding) {
		multiCardImplementation(t, s, b)
		b.Steps[0].Cards = []SessionCardRef{{"KB-SDP-001", "whole"}}
		next := b.Steps[0]
		next.ID = "second"
		next.Cards = []SessionCardRef{{"KB-SDP-002", "whole"}}
		next.DependsOn = []string{b.Steps[0].ID}
		b.Steps = append(b.Steps, next)
		evidenceFile(t, s, b.Session.Path, string(mustFile(t, filepath.Join(s.root, b.Session.Path)))+"\n## second\n\nComplete unrelated remaining card.\n")
	})
	completeFixtureDisposition(t, s, d)
	path := "SDP/KanBan/completed/#001--Change--fixture.md"
	original := mustFile(t, filepath.Join(s.root, path))
	editRebindPlan(t, s, d.InstanceID, nil)
	q := rebindRequest(t, s, d.InstanceID)
	done := commitRebind(t, s, q)
	if !bytes.Equal(original, mustFile(t, filepath.Join(s.root, path))) {
		t.Fatal("completed card changed during rebind")
	}
	sq := SessionSelectRequest{SessionSchema, q.ProjectRef, q.CheckoutID, newTestID(), q.InstanceID, done.Revision, "S1", "Attempt to repeat completed work"}
	if _, e := s.SelectSession(LocalController(), sq); e == nil {
		t.Fatal("completed step selected")
	}
	sq.OperationID = newTestID()
	sq.StepID = "second"
	sq.Reason = "Select remaining card"
	p, e := s.SelectSession(LocalController(), sq)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = s.ResumeSession(LocalController(), SessionResumeRequest{SessionSchema, q.ProjectRef, q.CheckoutID, newTestID(), q.InstanceID, p.Revision, "Commit remaining selection"}); e != nil {
		t.Fatal(e)
	}
	d2 := finishSelectedStep(t, s, q.InstanceID)
	commitAdvance(t, s, advanceRequest(t, s, q.InstanceID, "complete", "", "disposition", d2.OperationID))
	// Overall closure may append its own audit journal; the completed content stays intact.
	if !bytes.HasPrefix(mustFile(t, filepath.Join(s.root, path)), original) {
		t.Fatal("completed card history rewritten during later work")
	}
}

func TestSessionRebindPreparationRecoveryAndConflicts(t *testing.T) {
	for _, point := range []string{"before-head", "after-head"} {
		t.Run(point, func(t *testing.T) {
			s, id := rebindFixtureAt(t, "", "")
			editRebindPlan(t, s, id, nil)
			q := rebindRequest(t, s, id)
			st, _, _ := s.load()
			path := st.Sessions[id].Context.Binding.Session.Path
			before := mustFile(t, filepath.Join(s.root, path))
			fired := false
			s.fault = func(at string) error {
				if at == point && !fired {
					fired = true
					return errors.New("preparation interrupted")
				}
				return nil
			}
			if _, e := s.RebindSession(LocalController(), q); e == nil || !fired {
				t.Fatal("fault absent", e)
			}
			if !bytes.Equal(before, mustFile(t, filepath.Join(s.root, path))) {
				t.Fatal("preparation wrote canonical records")
			}
			fresh, e := New(s.root, s.area, s.project, s.state)
			if e != nil {
				t.Fatal(e)
			}
			p, e := fresh.RebindSession(LocalController(), q)
			if e != nil {
				t.Fatal(e)
			}
			changed := q
			changed.Reason = "Changed retry"
			if _, e = fresh.RebindSession(LocalController(), changed); e == nil {
				t.Fatal("conflicting retry accepted")
			}
			if _, e = fresh.SessionRebindContext(SessionRebindContextRequest{id}); e == nil {
				t.Fatal("prepared rebind permitted another")
			}
			if _, e = fresh.ResumeSession(LocalController(), SessionResumeRequest{SessionSchema, q.ProjectRef, q.CheckoutID, newTestID(), id, p.Revision, "Recover exact rebind"}); e != nil {
				t.Fatal(e)
			}
		})
	}
	t.Run("shared-claim", func(t *testing.T) {
		s, id := rebindFixtureAt(t, "", "")
		editRebindPlan(t, s, id, nil)
		claim := sessionScopeClaim{Schema: "sdp-session-scope-claim/0.1", OperationID: newTestID(), SessionID: "SESSION-OTHER-0001", PlanID: "MAINT-SDP-0002", WritablePaths: []string{"target.txt"}, ProtectedPaths: []string{}, Cards: []string{"KB-SDP-002"}}
		b, _ := json.Marshal(claim)
		evidenceFile(t, s, "SDP/Sessions/execution/claims/"+claim.OperationID+".json", string(b))
		if _, e := s.SessionRebindContext(SessionRebindContextRequest{id}); e == nil {
			t.Fatal("shared claim conflict accepted")
		}
	})
}
