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

func ownerFixture(t *testing.T) (*Service, SessionOwnerRequest, SessionDisposeRequest) {
	t.Helper()
	return ownerFixtureAt(t, "", "")
}
func ownerFixtureAt(t *testing.T, root, state string) (*Service, SessionOwnerRequest, SessionDisposeRequest) {
	t.Helper()
	s, d := dispositionFixture(t, root, state, func(t *testing.T, s *Service, b *SessionBinding) {
		multiCardImplementation(t, s, b)
		b.Steps[0].OwnerAcceptance = "required"
	})
	q := SessionOwnerRequest{Schema: SessionSchema, ProjectRef: d.ProjectRef, CheckoutID: d.CheckoutID, OperationID: newTestID(), InstanceID: d.InstanceID, ExpectedRevision: d.ExpectedRevision, ExpectedRunRevision: d.ExpectedRunRevision, Action: "request", ResultID: d.ResultID, Reason: "Review this exact verified outcome before completing the selected work"}
	q.GateID = q.OperationID
	return s, q, d
}
func commitOwner(t *testing.T, s *Service, c Caller, q SessionOwnerRequest) SessionMutationResult {
	t.Helper()
	p, e := s.OwnerSession(c, q)
	if e != nil {
		t.Fatal(e)
	}
	r := SessionResumeRequest{SessionSchema, q.ProjectRef, q.CheckoutID, newTestID(), q.InstanceID, p.Revision, "Commit exact owner gate effects"}
	done, e := s.ResumeSession(LocalController(), r)
	if e != nil {
		t.Fatal(e)
	}
	retry, e := s.OwnerSession(c, q)
	if e != nil || !reflect.DeepEqual(p, retry) {
		t.Fatal(retry, e)
	}
	return done
}
func ownerDecision(q SessionOwnerRequest, revision uint64, action string) SessionOwnerRequest {
	q.OperationID = newTestID()
	q.ExpectedRevision = revision
	q.Action = action
	q.Reason = "Synthetic owner decision for this exact fixture result"
	return q
}
func TestSessionOwnerAcceptAndReject(t *testing.T) {
	for _, action := range []string{"accept", "changes-requested"} {
		t.Run(action, func(t *testing.T) {
			s, q, d := ownerFixture(t)
			done := commitOwner(t, s, LocalController(), q)
			v, e := s.ReadSession(q.InstanceID)
			if e != nil || v.Steps[0].Stage != "waiting-owner" {
				t.Fatal(v, e)
			}
			for _, name := range []string{"#001--Change--fixture.md", "#002--Change--unrelated.md"} {
				b := mustFile(t, filepath.Join(s.root, "SDP/KanBan/active", name))
				m, _ := metadata(b)
				if m["CardState"] != "gate-review" {
					t.Fatal(m)
				}
			}
			decision := ownerDecision(q, done.Revision, action)
			head := mustFile(t, filepath.Join(s.state, "HEAD"))
			if _, e = s.OwnerSession(LocalController(), decision); e == nil {
				t.Fatal("controller impersonated owner")
			}
			if !bytes.Equal(head, mustFile(t, filepath.Join(s.state, "HEAD"))) {
				t.Fatal("unauthorized mutation")
			}
			done = commitOwner(t, s, LocalOwner(), decision)
			v, e = s.ReadSession(q.InstanceID)
			if e != nil {
				t.Fatal(e)
			}
			d.ExpectedRevision = done.Revision
			if action == "accept" {
				if v.Steps[0].Stage != "disposition" {
					t.Fatal(v)
				}
				p, e := s.DisposeSession(LocalController(), d)
				if e != nil {
					t.Fatal(e)
				}
				if _, e = s.ResumeSession(LocalController(), disposeResume(d, p)); e != nil {
					t.Fatal(e)
				}
				assertDisposed(t, s, d)
			} else {
				if v.Steps[0].Stage != "rework" {
					t.Fatal(v)
				}
				if _, e = s.DisposeSession(LocalController(), d); e == nil {
					t.Fatal("rejected result completed")
				}
			}
			for _, name := range []string{"#001--Change--fixture.md"} {
				m, _ := metadata(mustFile(t, filepath.Join(s.root, "SDP/KanBan/active", name)))
				if m["CardState"] != "in-progress" {
					t.Fatal(m)
				}
			}
		})
	}
}
func TestSessionOwnerRefusals(t *testing.T) {
	for _, mode := range []string{"candidate", "skill", "revision", "run", "result", "gate", "decision-before-gate", "same-result-again"} {
		t.Run(mode, func(t *testing.T) {
			s, q, _ := ownerFixture(t)
			done := commitOwner(t, s, LocalController(), q)
			d := ownerDecision(q, done.Revision, "accept")
			switch mode {
			case "candidate":
				evidenceFile(t, s, "target.txt", "changed")
			case "skill":
				evidenceFile(t, s, "Skills/sdp-reviewer/SKILL.md", "changed")
			case "revision":
				d.ExpectedRevision++
			case "run":
				d.ExpectedRunRevision++
			case "result":
				d.ResultID = newTestID()
			case "gate":
				d.GateID = newTestID()
			case "decision-before-gate":
				d.GateID = d.OperationID
			case "same-result-again":
				r := commitOwner(t, s, LocalOwner(), d)
				d.Action = "request"
				d.OperationID = newTestID()
				d.GateID = d.OperationID
				d.ExpectedRevision = r.Revision
			}
			head := mustFile(t, filepath.Join(s.state, "HEAD"))
			if _, e := s.OwnerSession(LocalOwner(), d); e == nil {
				t.Fatal(mode)
			}
			if !bytes.Equal(head, mustFile(t, filepath.Join(s.state, "HEAD"))) {
				t.Fatal("refusal wrote state")
			}
		})
	}
}
func TestSessionOwnerRecovery(t *testing.T) {
	for _, action := range []string{"request", "accept"} {
		for _, point := range []string{"after-claim", "after-file:0", "after-file:1", "after-file:2", "after-file:3", "after-file:4", "before-head", "after-head"} {
			t.Run(action+point, func(t *testing.T) {
				s, q, _ := ownerFixture(t)
				caller := LocalController()
				if action == "accept" {
					done := commitOwner(t, s, caller, q)
					q = ownerDecision(q, done.Revision, action)
					caller = LocalOwner()
				}
				p, e := s.OwnerSession(caller, q)
				if e != nil {
					t.Fatal(e)
				}
				resume := SessionResumeRequest{SessionSchema, q.ProjectRef, q.CheckoutID, newTestID(), q.InstanceID, p.Revision, "Resume exact owner record publication"}
				fault := point
				if strings.HasPrefix(point, "after-") && point != "after-head" {
					fault = "session-owner-" + point
				}
				hit := false
				s.fault = func(p string) error {
					if p == fault {
						hit = true
						return errors.New("interrupted")
					}
					return nil
				}
				if _, e = s.ResumeSession(LocalController(), resume); e == nil || !hit {
					t.Fatal("fault missed", e)
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
					t.Fatal("reservation retained")
				}
			})
		}
	}
}
func TestSessionOwnerStoredTamper(t *testing.T) {
	s, q, _ := ownerFixture(t)
	commitOwner(t, s, LocalController(), q)
	st, _, e := s.load()
	if e != nil {
		t.Fatal(e)
	}
	raw, _ := json.Marshal(st)
	for _, mode := range []string{"profile", "outcome", "run", "postimage", "preimage", "coverage", "receipt"} {
		t.Run(mode, func(t *testing.T) {
			var bad snapshot
			json.Unmarshal(raw, &bad)
			i := bad.Sessions[q.InstanceID]
			g := &i.OwnerGates[0]
			switch mode {
			case "profile":
				bad.Schema = sessionDispositionStateSchema
			case "outcome":
				g.Request.Action = "accept"
			case "run":
				g.Request.ExpectedRunRevision++
			case "postimage":
				g.Effects.Changes[0].Content += "tamper"
			case "preimage":
				g.Before[managementLedger] = fmt.Sprintf("%s\n", g.Before[managementLedger])
			case "coverage":
				g.Effects.Coverage = g.Effects.Coverage[1:]
			case "receipt":
				delete(bad.Receipts, q.OperationID)
			}
			bad.Sessions[i.ID] = i
			if validateSessionOwnerGates(bad) == nil {
				t.Fatal(mode)
			}
		})
	}
}
