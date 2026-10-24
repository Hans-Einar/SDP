package governance

import (
	"path/filepath"
	"reflect"
	"slices"
)

const sessionResultStateSchema = "sdp-project-governance-state/0.17"

func sessionResultProfile(schema string) bool {
	return schema == sessionResultStateSchema || schema == sessionDispositionStateSchema
}
func sessionWorkProfile(schema string) bool {
	return schema == sessionWorkStateSchema || sessionResultProfile(schema)
}

func (s *Service) VerifySession(c Caller, q SessionWorkRequest) (SessionMutationResult, error) {
	return s.prepareSessionWork(c, "verify", q)
}
func (s *Service) ReviewSession(c Caller, q SessionWorkRequest) (SessionMutationResult, error) {
	return s.prepareSessionWork(c, "review", q)
}
func sessionWorkExtraPins(root string, w SessionWork) []Source {
	pins := []Source{}
	if w.Evidence != nil {
		pins = append(pins, w.Evidence.Request.Artifact)
	}
	for _, t := range w.Tests {
		pins = append(pins, t.Request.Artifact)
	}
	if w.Review != nil {
		for _, skill := range w.Review.Spec.Skills {
			p, _ := filepath.Rel(root, skill.Path)
			pins = append(pins, Source{filepath.ToSlash(p), skill.Digest})
		}
	}
	return pins
}
func sameSessionCandidate(a, b EvidenceRecord) bool {
	return a.Request.RunID == b.Request.RunID && a.Request.ReportID == b.Request.ReportID && a.Request.Candidate.Digest == b.Request.Candidate.Digest
}
func (s *Service) prepareSessionResult(st snapshot, i SessionInstance, w *SessionWork) error {
	want := 2
	if w.Kind == "review" {
		want = 3
	}
	if len(i.Work) != want || i.Work[want-1].Effects.Status != "committed" {
		return fail("conflict", "committed preceding Session stage required")
	}
	if w.Kind == "verify" {
		ev := st.Receipts[w.Request.EvidenceID].Result.Evidence
		if ev == nil {
			return fail("evidence_missing", "exact test receipt required")
		}
		copy := *ev
		w.Evidence = &copy
	} else {
		d, ok := st.Reviews[w.Request.ReviewID]
		if !ok {
			return fail("evidence_missing", "observed review required")
		}
		w.Review = &d
		for _, id := range d.Spec.EvidenceIDs {
			ev := st.Receipts[id].Result.Evidence
			if ev == nil {
				return fail("evidence_missing", "review test receipt absent")
			}
			w.Tests = append(w.Tests, *ev)
		}
	}
	return s.sessionResultSnapshot(st, i, *w)
}

// Pure linkage validation also runs on retained generations. It does not treat
// historical review as current acceptance or require its files to stay unchanged.
func validateSessionResult(st snapshot, i SessionInstance, w SessionWork) error {
	if w.Kind != "review" && (w.Review != nil || len(w.Tests) != 0) {
		return fail("reconciliation_required", "review snapshot on another work kind")
	}
	if w.Kind == "verify" {
		if len(i.Work) < 2 || w.Evidence == nil || i.Work[1].Evidence == nil || w.Evidence.Request.Kind != "test" || w.Evidence.Request.Verdict != "passed" || !sameSessionCandidate(*w.Evidence, *i.Work[1].Evidence) {
			return fail("evidence_missing", "passing tests for progressed report/candidate required")
		}
	}
	if w.Kind == "review" {
		d := w.Review
		if len(i.Work) < 3 || d == nil || w.Evidence != nil || w.Request.ReviewID != d.ID || w.Request.ExpectedReviewRevision != d.Revision || !reflect.DeepEqual(st.Reviews[d.ID], *d) || d.State != "completed" || d.Outcome == nil || d.Outcome.Verdict != "passed" || d.ThreadID == "" || d.TurnID == "" {
			return fail("evidence_missing", "exact terminal passing observed review required")
		}
		verified := i.Work[2].Evidence
		if verified == nil || d.RunID != verified.Request.RunID || d.Spec.ReportID != verified.Request.ReportID || d.Spec.Candidate.Digest != verified.Request.Candidate.Digest || !slices.Contains(d.Spec.EvidenceIDs, i.Work[2].Request.EvidenceID) || len(w.Tests) != len(d.Spec.EvidenceIDs) {
			return fail("evidence_missing", "review must cover verified Session basis")
		}
		var payload reviewPayload
		r := st.Runs[d.RunID]
		expected := reviewPayload{ReviewSchema, r.Selection, r.Assignment, d.Spec, w.Tests, &i.Context.Binding}
		if Decode([]byte(d.Payload), &payload) != nil || !reflect.DeepEqual(payload, expected) {
			return fail("reconciliation_required", "review context does not match full Session binding and tests")
		}
		for n, id := range d.Spec.EvidenceIDs {
			ev := st.Receipts[id].Result.Evidence
			if ev == nil || !reflect.DeepEqual(*ev, w.Tests[n]) {
				return fail("reconciliation_required", "review test snapshot differs")
			}
		}
	}
	return nil
}
func failedSessionEvidence(st snapshot, r Run, ev EvidenceRecord) error {
	for _, receipt := range st.Receipts {
		if v := receipt.Result.Evidence; v != nil && sameSessionCandidate(*v, ev) && v.Request.Kind != "claim" && v.Request.Verdict == "failed" {
			return fail("evidence_missing", "contradictory failed evidence requires rework")
		}
	}
	for _, d := range st.Reviews {
		if d.RunID != r.ID {
			continue
		}
		if ReviewActive(d.State) {
			return fail("conflict", "active review prevents Session result intake")
		}
		if d.Spec.ReportID == ev.Request.ReportID && d.Spec.Candidate.Digest == ev.Request.Candidate.Digest && d.Outcome != nil && d.Outcome.Verdict != "passed" {
			return fail("evidence_missing", "unresolved review result requires reconciliation")
		}
	}
	return nil
}
func (s *Service) sessionResultSnapshot(st snapshot, i SessionInstance, w SessionWork) error {
	if e := validateSessionResult(st, i, w); e != nil {
		return e
	}
	if e := s.progressSnapshot(st, i, i.Work[1]); e != nil {
		return e
	}
	r := st.Runs[i.Work[0].Request.OperationID]
	if r.Revision != w.Request.ExpectedRunRevision {
		return fail("stale_revision", "result run revision differs")
	}
	if e := failedSessionEvidence(st, r, *i.Work[1].Evidence); e != nil {
		return e
	}
	if w.Kind == "verify" {
		return s.progressSnapshot(st, i, w)
	}
	// This check deliberately does not call runContext while own canonical effects
	// are pending. Exact pre/postimages protect authority; candidate/tests/skills
	// are checked before and after each effect through the existing review predicate.
	return s.checkedReviewEvidence(st, PublicationRequest{RunID: r.ID, ExpectedRevision: r.Revision, ReviewID: w.Request.ReviewID})
}

// Review admission is separate from result intake. The existing review engine
// still owns fresh-thread provenance, observations, outcomes and retries.
func sessionReviewBasis(st snapshot, r Run, spec ReviewSpec) error {
	if r.Routine.ID != SessionImplementation().ID {
		return nil
	}
	i, ok := sessionRunBinding(st, r.ID)
	if !ok || len(i.Work) < 3 || i.Work[2].Effects.Status != "committed" {
		return fail("conflict", "Session verification must commit before review")
	}
	ev := i.Work[2].Evidence
	if ev == nil || spec.ReportID != ev.Request.ReportID || spec.Candidate.Digest != ev.Request.Candidate.Digest || !slices.Contains(spec.EvidenceIDs, i.Work[2].Request.EvidenceID) {
		return fail("evidence_missing", "review must include verified Session tests/candidate")
	}
	for _, id := range spec.EvidenceIDs {
		v := st.Receipts[id].Result.Evidence
		if v == nil || v.Request.Verdict != "passed" {
			return fail("evidence_missing", "passing Session tests required before review")
		}
	}
	return nil
}
func (s *Service) projectSessionResult(i SessionInstance, st snapshot, v *SessionView, step *SessionStepView) {
	if len(i.Work) == 2 {
		if v.Coverage == "current" {
			v.NextActions = append(v.NextActions, "session-verify")
		}
		step.Holds = append(step.Holds, "explicit passing test verification required")
		return
	}
	verified := i.Work[2]
	if verified.Effects.Status != "committed" {
		step.Holds = append(step.Holds, "verification effects pending")
		return
	}
	step.Stage = "review"
	if selectedSessionStep(i).Review == "not-required" {
		step.Stage = "disposition"
	}
	step.VerificationEvidence = verified.Request.EvidenceID
	step.Holds = append([]string{}, v.Conditions...)
	if e := s.sessionResultSnapshot(st, i, verified); e != nil {
		step.Holds = append(step.Holds, e.Error())
	}
	if len(i.Work) == 3 && selectedSessionStep(i).Review == "not-required" {
		if selectedSessionStep(i).OwnerAcceptance == "required" {
			step.Holds = append(step.Holds, "owner acceptance required; decision adapter remains unimplemented")
		}
		step.Holds = append(step.Holds, "review explicitly not required; checked disposition and claim release required")
		return
	}
	if len(i.Work) == 3 {
		step.Holds = append(step.Holds, "current passing independent review required")
		if v.Coverage == "current" {
			v.NextActions = append(v.NextActions, "session-review")
		}
		return
	}
	reviewed := i.Work[3]
	if reviewed.Effects.Status != "committed" {
		step.Holds = append(step.Holds, "review effects pending")
		return
	}
	step.Stage = "disposition"
	step.ReviewID = reviewed.Request.ReviewID
	if e := s.sessionResultSnapshot(st, i, reviewed); e != nil {
		step.Holds = append(step.Holds, e.Error())
	}
	if selectedSessionStep(i).OwnerAcceptance == "required" {
		step.Holds = append(step.Holds, "owner acceptance required; decision adapter remains unimplemented")
	}
	step.Holds = append(step.Holds, "checked whole/contribution disposition and claim release required")
}
