package governance

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func assessmentFixture(t *testing.T, root, state string, scopes ...string) (*Service, SessionAssessRequest) {
	return assessmentConfigured(t, root, state, nil, scopes...)
}
func assessmentConfigured(t *testing.T, root, state string, configure sessionTestConfig, scopes ...string) (*Service, SessionAssessRequest) {
	t.Helper()
	s, selectQ := selectionConfigured(t, root, state, configure, scopes...)
	selected, e := s.SelectSession(LocalController(), selectQ)
	if e != nil {
		t.Fatal(e)
	}
	committed, e := s.ResumeSession(LocalController(), selectionResume(selectQ, selected))
	if e != nil {
		t.Fatal(e)
	}
	q := SessionAssessRequest{Schema: SessionSchema, ProjectRef: selectQ.ProjectRef, CheckoutID: selectQ.CheckoutID, OperationID: newTestID(), InstanceID: committed.InstanceID, ExpectedRevision: committed.Revision, StepID: "S1", Reason: "Assess the selected documentation correction with explicit impact criteria", Blueprint: SessionBlueprintDecision{"not-required", "No model or implementation contract changes in this documentation-only fixture"}}
	for _, domain := range sessionDesignDomains {
		q.Criteria = append(q.Criteria, SessionDesignCriterion{domain, "unchanged", "Fixture changes explanatory wording only; this domain retains its current behavior"})
	}
	for _, name := range []string{"sdp-project-leader", "sdp-change-analysis"} {
		p := "Skills/" + name + "/SKILL.md"
		content := "Synthetic " + name + " procedure snapshot; no host consumption claim\n"
		evidenceFile(t, s, p, content)
		q.Skills = append(q.Skills, SkillInput{name, filepath.Join(s.root, p), hash([]byte(content))})
	}
	evidenceFile(t, s, "evidence/design-scope.md", "Fixture scope: correct explanatory wording; no contracts or implementation behavior change.\n")
	q.Evidence = []Source{{"evidence/design-scope.md", hash(mustFile(t, filepath.Join(s.root, "evidence/design-scope.md")))}}
	return s, q
}
func assessResume(q SessionAssessRequest, r SessionMutationResult) SessionResumeRequest {
	return SessionResumeRequest{SessionSchema, q.ProjectRef, q.CheckoutID, newTestID(), q.InstanceID, r.Revision, "Commit explicit assessment records"}
}
func commitAssessment(t *testing.T, s *Service, q SessionAssessRequest) SessionMutationResult {
	t.Helper()
	r, e := s.AssessSession(LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	r, e = s.ResumeSession(LocalController(), assessResume(q, r))
	if e != nil {
		t.Fatal(e)
	}
	return r
}
func TestSessionAssessmentRoutesAndRecords(t *testing.T) {
	for _, tc := range []struct{ impact, blueprint, stage string }{{"unchanged", "not-required", "implementation"}, {"changed", "required", "design"}, {"unknown", "not-required", "admission"}, {"unchanged", "unknown", "admission"}, {"unchanged", "required", "implementation"}} {
		t.Run(tc.impact+tc.blueprint, func(t *testing.T) {
			s, q := assessmentFixture(t, "", "")
			q.Criteria[0].Impact = tc.impact
			q.Blueprint.Requirement = tc.blueprint
			before := mustFile(t, filepath.Join(s.root, managementLedger))
			prepared, e := s.AssessSession(LocalController(), q)
			if e != nil {
				t.Fatal(e)
			}
			if _, e = os.Stat(filepath.Join(s.root, assessmentPath(q))); !os.IsNotExist(e) {
				t.Fatal("prepare published artifact")
			}
			v, e := s.ReadSession(q.InstanceID)
			if e != nil || v.RecordEffects != "prepared" || v.Steps[0].Stage != "admission" || v.Steps[0].Assessment != nil {
				t.Fatal(v, e)
			}
			resume := assessResume(q, prepared)
			done, e := s.ResumeSession(LocalController(), resume)
			if e != nil {
				t.Fatal(e)
			}
			retry, e := s.ResumeSession(LocalController(), resume)
			if e != nil || retry != done {
				t.Fatal(retry, e)
			}
			retry, e = s.AssessSession(LocalController(), q)
			if e != nil || retry != prepared {
				t.Fatal(retry, e)
			}
			v, e = s.ReadSession(q.InstanceID)
			if e != nil || v.Coverage != "current" || v.Steps[0].Stage != tc.stage || v.Steps[0].Assessment.DesignImpact != tc.impact || v.HostExecution != "not-dispatched" || len(v.Steps[0].Holds) == 0 {
				t.Fatal(v, e)
			}
			if !bytes.HasPrefix(mustFile(t, filepath.Join(s.root, managementLedger)), before) {
				t.Fatal("history rewritten")
			}
			st, _, e := s.load()
			if e != nil {
				t.Fatal(e)
			}
			a := st.Sessions[q.InstanceID].Assessments[0]
			if !bytes.Equal(mustFile(t, filepath.Join(s.root, assessmentPath(q))), assessmentArtifact(a, st.Sessions[q.InstanceID])) {
				t.Fatal("artifact mismatch")
			}
			for _, p := range []string{"SDP/KanBan/active/#001--Change--fixture.md", st.Sessions[q.InstanceID].Context.Binding.Plan.Path, st.Sessions[q.InstanceID].Context.Binding.Session.Path} {
				raw := mustFile(t, filepath.Join(s.root, p))
				if strings.Count(string(raw), "## PG design assessment "+q.OperationID) != 1 {
					t.Fatal("journal", p)
				}
				if strings.Contains(p, "KanBan/") {
					m, _ := metadata(raw)
					if m["CardState"] != "ready" {
						t.Fatal(m)
					}
				}
			}
			// Fixture has a pre-existing synthetic run. Assessment must not create another.
			if len(st.Runs) != 1 {
				t.Fatal("assessment changed run collection")
			}
		})
	}
}
func TestSessionAssessmentUnknownCanBeResolved(t *testing.T) {
	s, q := assessmentFixture(t, "", "")
	q.Criteria[0].Impact = "unknown"
	done := commitAssessment(t, s, q)
	first := mustFile(t, filepath.Join(s.root, assessmentPath(q)))
	next := q
	next.Criteria = append([]SessionDesignCriterion{}, q.Criteria...)
	next.Criteria[0].Impact = "unchanged"
	next.OperationID = newTestID()
	next.ExpectedRevision = done.Revision
	next.Reason = "Resolved uncertainty against the same pinned scope evidence"
	commitAssessment(t, s, next)
	v, e := s.ReadSession(q.InstanceID)
	if e != nil || v.Steps[0].Stage != "implementation" || v.Steps[0].Assessment.OperationID != next.OperationID {
		t.Fatal(v, e)
	}
	if !bytes.Equal(first, mustFile(t, filepath.Join(s.root, assessmentPath(q)))) {
		t.Fatal("previous decision rewritten")
	}
	st, _, e := s.load()
	if e != nil || len(st.Sessions[q.InstanceID].Assessments) != 2 {
		t.Fatal(e)
	}
}
func TestSessionAssessmentRefusalsPreserveState(t *testing.T) {
	for _, name := range []string{"unauthorized", "stale", "missing-domain", "duplicate-domain", "blank-rationale", "missing-skills", "skill-drift", "evidence-drift", "evidence-ledger-drift", "unknown-field", "wrong-step", "unselected", "changed-retry"} {
		t.Run(name, func(t *testing.T) {
			s, q := assessmentFixture(t, "", "")
			caller := LocalController()
			switch name {
			case "unauthorized":
				caller = Caller{}
			case "stale":
				q.ExpectedRevision--
			case "missing-domain":
				q.Criteria = q.Criteria[:6]
			case "duplicate-domain":
				q.Criteria[0].Domain = q.Criteria[1].Domain
			case "blank-rationale":
				q.Criteria[0].Rationale = " "
			case "missing-skills":
				q.Skills = q.Skills[:1]
			case "skill-drift":
				os.WriteFile(q.Skills[0].Path, []byte("changed"), 0600)
			case "evidence-drift":
				q.Evidence[0].SHA256 = hash([]byte("wrong"))
			case "evidence-ledger-drift":
				q.Evidence = []Source{{managementLedger, hash([]byte("wrong"))}}
			case "wrong-step":
				q.StepID = "S2"
			case "unselected":
				q.InstanceID = newTestID()
			case "changed-retry":
				if _, e := s.AssessSession(caller, q); e != nil {
					t.Fatal(e)
				}
				q.Reason = "different content under same operation"
			case "unknown-field":
				raw, _ := json.Marshal(q)
				raw = append(raw[:len(raw)-1], []byte(",\"role\":\"controller\"}")...)
				if Decode(raw, &q) == nil {
					t.Fatal("unknown field accepted")
				}
				return
			}
			_, head, _ := s.load()
			if _, e := s.AssessSession(caller, q); e == nil {
				t.Fatal("assessment accepted")
			}
			_, after, _ := s.load()
			if head != after {
				t.Fatal("refusal changed state")
			}
		})
	}
}
func TestSessionAssessmentCrashRecovery(t *testing.T) {
	for _, point := range []string{"session-assessment-after-claim", "session-assessment-after-file:0", "session-assessment-after-file:1", "session-assessment-after-file:2", "session-assessment-after-file:3", "session-assessment-after-file:4", "before-head", "after-head"} {
		t.Run(point, func(t *testing.T) {
			s, q := assessmentFixture(t, "", "")
			r, e := s.AssessSession(LocalController(), q)
			if e != nil {
				t.Fatal(e)
			}
			resume := assessResume(q, r)
			hit := false
			s.fault = func(p string) error {
				if p == point {
					hit = true
					return errors.New("simulated crash")
				}
				return nil
			}
			if _, e = s.ResumeSession(LocalController(), resume); e == nil || !hit {
				t.Fatal("fault not reached", point, e)
			}
			s.fault = nil
			s2, e := New(s.root, s.area, s.project, s.state)
			if e != nil {
				t.Fatal(e)
			}
			if _, e = s2.ResumeSession(LocalController(), resume); e != nil {
				t.Fatal(e)
			}
			v, e := s2.ReadSession(q.InstanceID)
			if e != nil || v.Coverage != "current" || v.RecordEffects != "committed" || v.Steps[0].Stage != "implementation" {
				t.Fatal(v, e)
			}
			if _, e = os.Stat(filepath.Join(s.root, managementLedger) + ".pg-pending"); !os.IsNotExist(e) {
				t.Fatal("reservation retained")
			}
			raw := mustFile(t, filepath.Join(s.root, "SDP/KanBan/active/#001--Change--fixture.md"))
			if strings.Count(string(raw), "## PG design assessment "+q.OperationID) != 1 {
				t.Fatal("duplicate journal")
			}
		})
	}
}
func TestSessionAssessmentDriftAndPendingRetry(t *testing.T) {
	s, q := assessmentFixture(t, "", "")
	st, _, _ := s.load()
	i := st.Sessions[q.InstanceID]
	var oldResume SessionMutationResult
	for _, r := range st.Receipts {
		if r.Result.Session != nil && r.Result.Session.Revision == i.Revision {
			oldResume = *r.Result.Session
		}
	}
	r, e := s.AssessSession(LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	resume := assessResume(q, r)
	s.fault = func(p string) error {
		if p == "session-assessment-after-claim" {
			return errors.New("crash")
		}
		return nil
	}
	s.ResumeSession(LocalController(), resume)
	s.fault = nil
	// A historical retry cannot release a newer pending assessment reservation.
	oldRequest := SessionResumeRequest{SessionSchema, q.ProjectRef, q.CheckoutID, oldResume.OperationID, q.InstanceID, oldResume.Revision - 1, "Apply recorded selection effects"}
	if _, e = s.ResumeSession(LocalController(), oldRequest); e != nil {
		t.Fatal(e)
	}
	if _, e = os.Stat(filepath.Join(s.root, managementLedger) + ".pg-pending"); e != nil {
		t.Fatal("historical retry released pending reservation", e)
	}
	tx, e := s.canonicalTransaction()
	if e != nil {
		t.Fatal(e)
	}
	if e = tx.Check(s.state + "\n" + q.OperationID + "\n"); e != nil {
		t.Fatal(e)
	}
	tx.Close()
	if oldResume.Operation != "session-resume" {
		t.Fatal("missing old receipt")
	}
	evidenceFile(t, s, q.Evidence[0].Path, "external drift")
	if _, e = s.ResumeSession(LocalController(), resume); e == nil {
		t.Fatal("overwrote evidence drift")
	}
	v, e := s.ReadSession(q.InstanceID)
	if e != nil || v.Coverage != "stale-or-unavailable" || v.Steps[0].Stage != "admission" {
		t.Fatal(v, e)
	}
	if string(mustFile(t, filepath.Join(s.root, q.Evidence[0].Path))) != "external drift" {
		t.Fatal("external bytes lost")
	}
}
func TestSessionAssessmentModelScopeCannotUseNoDesign(t *testing.T) {
	for _, path := range []string{"src/System.design", "SDP/SDL", "SDP/SDL/Example", "ui/window.sdui"} {
		t.Run(path, func(t *testing.T) {
			s, q := assessmentFixture(t, "", "", path)
			commitAssessment(t, s, q)
			v, e := s.ReadSession(q.InstanceID)
			if e != nil || v.Steps[0].Stage != "design" || v.Steps[0].Assessment.DesignImpact != "changed" {
				t.Fatal(v, e)
			}
		})
	}
}
func TestSessionAssessmentStoredTamperRefuses(t *testing.T) {
	s, q := assessmentFixture(t, "", "")
	commitAssessment(t, s, q)
	st, _, _ := s.load()
	for _, name := range []string{"old-profile", "missing-coverage", "extra-write", "missing-skill", "missing-artifact", "receipt"} {
		t.Run(name, func(t *testing.T) {
			raw, _ := json.Marshal(st)
			var bad snapshot
			json.Unmarshal(raw, &bad)
			i := bad.Sessions[q.InstanceID]
			a := &i.Assessments[0]
			switch name {
			case "old-profile":
				bad.Schema = sessionEffectsStateSchema
			case "missing-coverage":
				a.Effects.Coverage = a.Effects.Coverage[1:]
			case "extra-write":
				a.Effects.Changes = append(a.Effects.Changes, PublicationChange{"src/unowned.go", "absent", hash([]byte("bad")), "bad"})
				a.Effects.Inputs = append(a.Effects.Inputs, Source{"src/unowned.go", "absent"})
			case "missing-skill":
				a.Request.Skills = nil
			case "missing-artifact":
				for n, ch := range a.Effects.Changes {
					if ch.Path == assessmentPath(q) {
						a.Effects.Changes = append(a.Effects.Changes[:n], a.Effects.Changes[n+1:]...)
						break
					}
				}
			case "receipt":
				delete(bad.Receipts, q.OperationID)
			}
			bad.Sessions[q.InstanceID] = i
			if validateSessionAssessments(bad) == nil {
				t.Fatal(name)
			}
		})
	}
	if validateSessionAssessments(st) != nil {
		t.Fatal("original rejected")
	}
}

func TestSessionAssessmentCoverageReconciliation(t *testing.T) {
	for _, target := range []string{"evidence", "skill", "assessment-record", "selected-plan"} {
		t.Run(target, func(t *testing.T) {
			s, q := assessmentFixture(t, "", "")
			done := commitAssessment(t, s, q)
			addUnrelatedCard(t, s)
			resume := SessionResumeRequest{SessionSchema, q.ProjectRef, q.CheckoutID, newTestID(), q.InstanceID, done.Revision, "Observe unrelated valid history"}
			done, e := s.ResumeSession(LocalController(), resume)
			if e != nil {
				t.Fatal(e)
			}
			v, e := s.ReadSession(q.InstanceID)
			if e != nil || v.Coverage != "current" {
				t.Fatal(v, e)
			}
			path := q.Evidence[0].Path
			switch target {
			case "skill":
				path, _ = filepath.Rel(s.root, q.Skills[0].Path)
			case "assessment-record":
				path = assessmentPath(q)
			case "selected-plan":
				path = "SDP/Maintenance/Plan.md"
			}
			raw := mustFile(t, filepath.Join(s.root, path))
			evidenceFile(t, s, path, string(raw)+"\nExternal correction\n")
			v, e = s.ReadSession(q.InstanceID)
			if e != nil || v.Coverage == "current" {
				t.Fatal("drift accepted", v, e)
			}
			q.OperationID = newTestID()
			q.ExpectedRevision = done.Revision
			if _, e = s.AssessSession(LocalController(), q); e == nil {
				t.Fatal("assessment silently rebound edited coverage")
			}
		})
	}
}

func TestSessionAssessmentInterruptedPrepareAndLimit(t *testing.T) {
	for _, point := range []string{"before-head", "after-head"} {
		t.Run(point, func(t *testing.T) {
			s, q := assessmentFixture(t, "", "")
			s.fault = func(p string) error {
				if p == point {
					return errors.New("interrupted prepare")
				}
				return nil
			}
			if _, e := s.AssessSession(LocalController(), q); e == nil {
				t.Fatal("fault not reached")
			}
			s.fault = nil
			if _, e := s.AssessSession(LocalController(), q); e != nil {
				t.Fatal(e)
			}
			if _, e := os.Stat(filepath.Join(s.root, assessmentPath(q))); !os.IsNotExist(e) {
				t.Fatal("prepare wrote records")
			}
		})
	}
	s, q := assessmentFixture(t, "", "")
	for n := 0; n < maxSessionAssessments; n++ {
		done := commitAssessment(t, s, q)
		q.OperationID = newTestID()
		q.ExpectedRevision = done.Revision
	}
	_, before, _ := s.load()
	if _, e := s.AssessSession(LocalController(), q); e == nil {
		t.Fatal("assessment bound ignored")
	}
	_, after, _ := s.load()
	if before != after {
		t.Fatal("limit refusal changed state")
	}
	v, e := s.ReadSession(q.InstanceID)
	if e != nil {
		t.Fatal(e)
	}
	for _, action := range v.NextActions {
		if action == "session-assess" {
			t.Fatal("unavailable action advertised")
		}
	}
}
