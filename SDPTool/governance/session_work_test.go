package governance

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func workFixture(t *testing.T, root, state string) (*Service, SessionWorkRequest) {
	return workConfigured(t, root, state, nil, nil)
}
func workConfigured(t *testing.T, root, state string, configure sessionTestConfig, assess func(*SessionAssessRequest)) (*Service, SessionWorkRequest) {
	t.Helper()
	s, a := assessmentConfigured(t, root, state, configure)
	if assess != nil {
		assess(&a)
	}
	done := commitAssessment(t, s, a)
	st, _, e := s.load()
	if e != nil {
		t.Fatal(e)
	}
	for _, r := range st.Runs {
		if r.State == "canceled" {
			continue
		}
		if _, e = s.Transition(LocalController(), transition(Result{ProjectRef: s.project, CheckoutID: st.CheckoutID, Run: r}, "cancel")); e != nil {
			t.Fatal(e)
		}
	}
	q := SessionWorkRequest{Schema: SessionSchema, ProjectRef: s.project, CheckoutID: a.CheckoutID, OperationID: newTestID(), InstanceID: a.InstanceID, ExpectedRevision: done.Revision, Reason: "Start the selected bounded no-design task", AssessmentID: a.OperationID, Outcome: "Correct the selected fixture text without changing system design"}
	for _, name := range []string{"sdp-master", "sdp-worker"} {
		p := "Skills/" + name + "/SKILL.md"
		body := "---\nname: " + name + "\ndescription: Synthetic bounded test procedure.\n---\n\nPreserve assigned scope; report exact candidate.\n"
		evidenceFile(t, s, p, body)
		q.Skills = append(q.Skills, SkillInput{name, filepath.Join(s.root, p), hash([]byte(body))})
	}
	return s, q
}
func workResume(q SessionWorkRequest, r SessionMutationResult) SessionResumeRequest {
	return SessionResumeRequest{SessionSchema, q.ProjectRef, q.CheckoutID, newTestID(), q.InstanceID, r.Revision, "Commit Session work effects"}
}
func startWork(t *testing.T, s *Service, q SessionWorkRequest) (SessionMutationResult, Result) {
	t.Helper()
	r, e := s.StartSession(LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	r, e = s.ResumeSession(LocalController(), workResume(q, r))
	if e != nil {
		t.Fatal(e)
	}
	run, e := s.Read(q.OperationID)
	if e != nil {
		t.Fatal(e)
	}
	return r, run
}
func reportWork(t *testing.T, s *Service, r Result) Result {
	t.Helper()
	at, epoch := newTestID(), newTestID()
	host := func(action, kind, status string) {
		q := HostRequest{Schema: Schema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, OperationID: newTestID(), RunID: r.Run.ID, ExpectedRevision: r.Run.Revision, InputRevision: r.Run.InputRevision, Action: action, AttemptID: at, HostVersion: "0.164.0", ThreadID: "session-fixture-thread-" + r.Run.ID, Observation: Observation{Kind: kind, Epoch: epoch, ReceivedAt: time.Now().UTC().Format(time.RFC3339Nano), Status: status}}
		if action != "bind" {
			q.Observation.TurnID = "session-fixture-turn"
		}
		var e error
		r, e = s.RecordHost(LocalController(), q)
		if e != nil {
			t.Fatal(e)
		}
	}
	host("bind", "bound", "idle")
	var e error
	r, e = s.Transition(LocalController(), transition(r, "start"))
	if e != nil {
		t.Fatal(e)
	}
	host("observe", "turn.started", "running")
	if _, e = s.Transition(LocalController(), transition(r, "report")); e == nil {
		t.Fatal("controller impersonated assigned report")
	}
	host("observe", "turn.completed", "completed")
	agent, e := AssignedAgent(r.Run.ID, at, epoch)
	if e != nil {
		t.Fatal(e)
	}
	r, e = s.Transition(agent, transition(r, "report"))
	if e != nil {
		t.Fatal(e)
	}
	return r
}
func progressRequest(t *testing.T, s *Service, q SessionWorkRequest, done SessionMutationResult, r Result) SessionWorkRequest {
	t.Helper()
	candidate, e := s.CaptureCandidate(r.Run.ID, nil)
	if e != nil {
		t.Fatal(e)
	}
	evidenceFile(t, s, "evidence/session-result.log", "Synthetic candidate check; no native execution or independent review claim.\n")
	eq := EvidenceRequest{Schema: EvidenceSchema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, OperationID: newTestID(), RunID: r.Run.ID, ExpectedRevision: r.Run.Revision, InputRevision: r.Run.InputRevision, ReportID: r.OperationID, Candidate: candidate, Kind: "test", Verdict: "passed", Issuer: "fixture verifier", Note: "Candidate-specific fixture evidence", Artifact: Source{"evidence/session-result.log", hash(mustFile(t, filepath.Join(s.root, "evidence/session-result.log")))}}
	if _, e = s.AttachEvidence(LocalController(), eq); e != nil {
		t.Fatal(e)
	}
	return SessionWorkRequest{Schema: SessionSchema, ProjectRef: q.ProjectRef, CheckoutID: q.CheckoutID, OperationID: newTestID(), InstanceID: q.InstanceID, ExpectedRevision: done.Revision, Reason: "Record submitted candidate for verification", EvidenceID: eq.OperationID, ExpectedRunRevision: r.Run.Revision}
}
func TestSessionWorkStartProgressJourney(t *testing.T) {
	s, q := workFixture(t, "", "")
	before := mustFile(t, filepath.Join(s.root, managementLedger))
	prepared, e := s.StartSession(LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = s.Read(q.OperationID); e == nil {
		t.Fatal("run exists before effects")
	}
	st, _, _ := s.load()
	if len(st.Runs) != 1 {
		t.Fatal("prepare created run")
	}
	retry, e := s.StartSession(LocalController(), q)
	if e != nil || retry != prepared {
		t.Fatal(retry, e)
	}
	resume := workResume(q, prepared)
	done, e := s.ResumeSession(LocalController(), resume)
	if e != nil {
		t.Fatal(e)
	}
	again, e := s.ResumeSession(LocalController(), resume)
	if e != nil || again != done {
		t.Fatal(again, e)
	}
	r, e := s.Read(q.OperationID)
	if e != nil || r.Run.State != "ready" || r.Coverage != "known" {
		t.Fatal(r, e)
	}
	card := mustFile(t, filepath.Join(s.root, "SDP/KanBan/active/#001--Change--fixture.md"))
	m, _ := metadata(card)
	if m["CardState"] != "in-progress" {
		t.Fatal(m)
	}
	v, e := s.ReadSession(q.InstanceID)
	if e != nil || v.HostExecution != "not-dispatched" || v.Steps[0].RunID != r.Run.ID {
		t.Fatal(v, e)
	}
	addUnrelatedCard(t, s)
	if e = s.CheckRunReady(r.Run.ID); e != nil {
		t.Fatal("unrelated history invalidated run", e)
	}
	// Existing engine/capability reports attach to this exact Session-created run.
	evidenceFile(t, s, "target.txt", "Implementation fixture update\n")
	reported := reportWork(t, s, r)
	pq := progressRequest(t, s, q, done, reported)
	p, e := s.ProgressSession(LocalController(), pq)
	if e != nil {
		t.Fatal(e)
	}
	progressResume := workResume(pq, p)
	_, e = s.ResumeSession(LocalController(), progressResume)
	if e != nil {
		t.Fatal(e)
	}
	v, e = s.ReadSession(q.InstanceID)
	if e != nil || v.Steps[0].Stage != "verification" || v.Steps[0].ProgressEvidence != pq.EvidenceID || v.HostExecution != "completed" || v.Coverage != "current" {
		t.Fatal(v, e)
	}
	if !bytes.HasPrefix(mustFile(t, filepath.Join(s.root, managementLedger)), before) {
		t.Fatal("history rewritten")
	}
	m, _ = metadata(mustFile(t, filepath.Join(s.root, "SDP/KanBan/active/#001--Change--fixture.md")))
	if m["CardState"] != "in-progress" {
		t.Fatal("report closed card")
	}
	if _, e = s.PreparePublication(LocalController(), PublicationRequest{Schema: PublicationSchema, ProjectRef: q.ProjectRef, CheckoutID: q.CheckoutID, OperationID: newTestID(), RunID: r.Run.ID, ExpectedRevision: reported.Run.Revision, ReviewID: newTestID(), System: "SDPTOOL", Summary: "Premature publication", Remaining: "none"}); e == nil {
		t.Fatal("legacy publication bypassed Session disposition")
	}
	if _, e = s.ProgressSession(LocalController(), pq); e != nil {
		t.Fatal("historical progress retry", e)
	}
	st, _, e = s.load()
	if e != nil || len(st.Runs) != 2 {
		t.Fatal(e)
	}
}
func TestSessionWorkCrashBoundaries(t *testing.T) {
	for _, kind := range []string{"start", "progress"} {
		for _, point := range []string{"before-head", "after-head", "after-claim", "after-file:0", "after-file:1", "after-file:2", "after-file:3", "after-file:4"} {
			t.Run(kind+point, func(t *testing.T) {
				s, q := workFixture(t, "", "")
				prepare := s.StartSession
				if kind == "progress" {
					done, r := startWork(t, s, q)
					r = reportWork(t, s, r)
					q = progressRequest(t, s, q, done, r)
					prepare = s.ProgressSession
				}
				p, e := prepare(LocalController(), q)
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
					t.Fatal("fault not reached", fault, e)
				}
				fresh, e := New(s.root, s.area, s.project, s.state)
				if e != nil {
					t.Fatal(e)
				}
				if _, e = fresh.ResumeSession(LocalController(), resume); e != nil {
					t.Fatal(e)
				}
				v, e := fresh.ReadSession(q.InstanceID)
				if e != nil || v.Coverage != "current" || v.RecordEffects != "committed" {
					t.Fatal(v, e)
				}
				st, _, e := fresh.load()
				if e != nil || len(st.Runs) != 2 {
					t.Fatal("duplicate or missing run", e)
				}
				if _, e = os.Stat(filepath.Join(s.root, managementLedger) + ".pg-pending"); !os.IsNotExist(e) {
					t.Fatal("pending reservation retained")
				}
			})
		}
	}
}
func TestSessionWorkAdmissionRefusals(t *testing.T) {
	for _, name := range []string{"unauthorized", "stale", "old-assessment", "missing-skills", "skill-drift", "extra-progress-fields", "changed-retry", "overlap"} {
		t.Run(name, func(t *testing.T) {
			s, q := workFixture(t, "", "")
			caller := LocalController()
			switch name {
			case "unauthorized":
				caller = Caller{}
			case "stale":
				q.ExpectedRevision--
			case "old-assessment":
				q.AssessmentID = newTestID()
			case "missing-skills":
				q.Skills = q.Skills[:1]
			case "skill-drift":
				os.WriteFile(q.Skills[0].Path, []byte("changed"), 0600)
			case "extra-progress-fields":
				q.EvidenceID = newTestID()
			case "changed-retry":
				if _, e := s.StartSession(caller, q); e != nil {
					t.Fatal(e)
				}
				q.Outcome = "different outcome"
			case "overlap":
				st, head, _ := s.load()
				for id, r := range st.Runs {
					r.State = "ready"
					st.Runs[id] = r
				}
				if e := s.save(st, head); e != nil {
					t.Fatal(e)
				}
			}
			_, before, _ := s.load()
			if _, e := s.StartSession(caller, q); e == nil {
				t.Fatal(name)
			}
			_, after, _ := s.load()
			if before != after {
				t.Fatal("refusal changed state")
			}
		})
	}
}
func TestSessionWorkProgressDriftPreservesCandidate(t *testing.T) {
	for _, phase := range []string{"prepare", "after-claim", "after-file:2"} {
		t.Run(phase, func(t *testing.T) {
			s, q := workFixture(t, "", "")
			done, r := startWork(t, s, q)
			r = reportWork(t, s, r)
			pq := progressRequest(t, s, q, done, r)
			if phase == "prepare" {
				evidenceFile(t, s, "target.txt", "external change")
				if _, e := s.ProgressSession(LocalController(), pq); e == nil {
					t.Fatal("stale candidate accepted")
				}
				return
			}
			p, e := s.ProgressSession(LocalController(), pq)
			if e != nil {
				t.Fatal(e)
			}
			resume := workResume(pq, p)
			s.fault = func(point string) error {
				if point == "session-progress-"+phase {
					return errors.New("crash")
				}
				return nil
			}
			if _, e = s.ResumeSession(LocalController(), resume); e == nil {
				t.Fatal("fault absent")
			}
			s.fault = nil
			evidenceFile(t, s, "target.txt", "external change")
			if _, e = s.ResumeSession(LocalController(), resume); e == nil {
				t.Fatal("candidate drift accepted")
			}
			if string(mustFile(t, filepath.Join(s.root, "target.txt"))) != "external change" {
				t.Fatal("external bytes lost")
			}
		})
	}
}
func TestSessionWorkRequiredLaunchSkillsAndClaims(t *testing.T) {
	s, q := workFixture(t, "", "")
	_, r := startWork(t, s, q)
	input := InputRequest{Schema: Schema, ProjectRef: q.ProjectRef, CheckoutID: q.CheckoutID, OperationID: newTestID(), RunID: r.Run.ID, ExpectedRevision: r.Run.Revision, InputRevision: r.Run.InputRevision, Action: "prepare", Input: OwnerInput{ID: newTestID(), Kind: "launch", Text: "Bounded input", Payload: "Selected Session", AttemptID: newTestID(), Epoch: newTestID(), ReceivedAt: time.Now().UTC().Format(time.RFC3339Nano)}}
	if _, e := s.RecordInput(LocalController(), input); e == nil {
		t.Fatal("missing Master/Worker input accepted")
	}
	input.Input.Skills = q.Skills
	if _, e := s.RecordInput(LocalController(), input); e != nil {
		t.Fatal(e)
	}
	legacy := StartRequest{Schema: Schema, ProjectRef: q.ProjectRef, CheckoutID: q.CheckoutID, OperationID: newTestID(), RequestID: newTestID(), Selection: r.Run.Selection, Assignment: r.Run.Assignment, RoutineDigest: Maintenance().Digest(), InputRevision: r.Run.InputRevision}
	if _, e := s.Start(LocalController(), legacy); e == nil {
		t.Fatal("legacy entry bypassed Session scope")
	}
}
func TestSessionWorkStoredBindingRefusesTamper(t *testing.T) {
	s, q := workFixture(t, "", "")
	startWork(t, s, q)
	st, _, _ := s.load()
	for _, name := range []string{"profile", "run-scope", "run-input", "record", "pins", "receipt"} {
		t.Run(name, func(t *testing.T) {
			b, _ := json.Marshal(st)
			var bad snapshot
			json.Unmarshal(b, &bad)
			i := bad.Sessions[q.InstanceID]
			switch name {
			case "profile":
				bad.Schema = sessionAssessmentStateSchema
			case "run-scope":
				r := bad.Runs[q.OperationID]
				r.Assignment.WritablePaths = []string{"other"}
				bad.Runs[r.ID] = r
			case "run-input":
				r := bad.Runs[q.OperationID]
				r.InputRevision = hash([]byte("other"))
				bad.Runs[r.ID] = r
			case "record":
				i.Work[0].Effects.Changes = i.Work[0].Effects.Changes[1:]
			case "pins":
				i.Work[0].Effects.Coverage = i.Work[0].Effects.Coverage[1:]
			case "receipt":
				delete(bad.Receipts, q.OperationID)
			}
			bad.Sessions[i.ID] = i
			if validateSessionWork(bad) == nil {
				t.Fatal(name)
			}
		})
	}
	if e := validateSessionWork(st); e != nil {
		t.Fatal(e)
	}
}

func TestSessionWorkDesignAndBlueprintGates(t *testing.T) {
	for _, tc := range []struct{ impact, blueprint string }{{"changed", "not-required"}, {"unknown", "not-required"}, {"unchanged", "required"}, {"unchanged", "unknown"}} {
		t.Run(tc.impact+tc.blueprint, func(t *testing.T) {
			s, q := workConfigured(t, "", "", nil, func(a *SessionAssessRequest) {
				a.Criteria[0].Impact = tc.impact
				a.Blueprint.Requirement = tc.blueprint
			})
			if _, e := s.StartSession(LocalController(), q); e == nil {
				t.Fatal("design gate bypassed")
			}
			st, _, _ := s.load()
			if len(st.Runs) != 1 {
				t.Fatal("gate created run")
			}
		})
	}
}
func multiCardImplementation(t *testing.T, s *Service, b *SessionBinding) {
	t.Helper()
	// This alters synthetic authored fixture records before adoption, not history
	// of an executed Session. The existing synthetic legacy run is canceled below.
	old := b.Plan.ID
	b.Plan.ID = "PLAN-SDP-0001"
	paths := []string{b.Plan.Path, "SDP/KanBan/backlog/#001--Change--fixture.md", managementLedger}
	for _, p := range paths {
		raw := string(mustFile(t, filepath.Join(s.root, p)))
		raw = strings.ReplaceAll(raw, old, b.Plan.ID)
		raw = strings.ReplaceAll(raw, "MaintenancePlan", "ImplementationPlan")
		raw = strings.ReplaceAll(raw, "\"Maintenance\"", "\"Plan\"")
		evidenceFile(t, s, p, raw)
	}
	addUnrelatedCard(t, s)
	p := "SDP/KanBan/backlog/#002--Change--unrelated.md"
	raw := string(mustFile(t, filepath.Join(s.root, p)))
	evidenceFile(t, s, p, raw+"| PlanId | PLAN-SDP-0001 |\n")
	b.Steps[0].Cards[0].Coverage = "contribution"
	b.Steps[0].Cards = append(b.Steps[0].Cards, SessionCardRef{"KB-SDP-002", "whole"})
	// The earlier synthetic run references a MaintenancePlan now replaced in this
	// disposable fixture. Retire that test run before the new typed-plan scenario.
	st, head, e := s.load()
	if e != nil {
		t.Fatal(e)
	}
	for id, r := range st.Runs {
		r.State = "canceled"
		st.Runs[id] = r
	}
	if e = s.save(st, head); e != nil {
		t.Fatal(e)
	}
}
func TestSessionWorkImplementationPlanAndAllContributions(t *testing.T) {
	s, q := workConfigured(t, "", "", multiCardImplementation, nil)
	done, r := startWork(t, s, q)
	v, e := s.RunSession(r.Run.ID)
	if e != nil || v.PlanID != "PLAN-SDP-0001" || v.Binding == nil || len(v.Binding.Steps[0].Cards) != 2 || v.Binding.Steps[0].Cards[0].Coverage != "contribution" {
		t.Fatal(v, e)
	}
	r = reportWork(t, s, r)
	pq := progressRequest(t, s, q, done, r)
	p, e := s.ProgressSession(LocalController(), pq)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = s.ResumeSession(LocalController(), workResume(pq, p)); e != nil {
		t.Fatal(e)
	}
	for _, name := range []string{"#001--Change--fixture.md", "#002--Change--unrelated.md"} {
		raw := mustFile(t, filepath.Join(s.root, "SDP/KanBan/active", name))
		m, e := metadata(raw)
		if e != nil || m["CardState"] != "in-progress" || !strings.Contains(string(raw), "PG Session progress") {
			t.Fatal(name, m, e)
		}
	}
}
func TestSessionWorkWireRejectsAliasesAndCrossVerbFields(t *testing.T) {
	_, q := workFixture(t, "", "")
	raw, _ := json.Marshal(q)
	if _, e := DecodeSessionWork(raw, "start"); e != nil {
		t.Fatal(e)
	}
	for _, suffix := range []string{`,"schema":"duplicate"}`, `,"Schema":"alias"}`, `,"evidenceId":""}`, `,"expectedRunRevision":0}`, `,"role":"owner"}`} {
		b := append(append([]byte{}, raw[:len(raw)-1]...), []byte(suffix)...)
		if _, e := DecodeSessionWork(b, "start"); e == nil {
			t.Fatal(suffix)
		}
	}
}
