package governance

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
)

func dispositionFixture(t *testing.T, root, state string, configure sessionTestConfig) (*Service, SessionDisposeRequest) {
	t.Helper()
	s, q := resultFixture(t, root, state, configure)
	done := commitResult(t, s, "verify", q)
	st, _, e := s.load()
	if e != nil {
		t.Fatal(e)
	}
	if selectedSessionStep(st.Sessions[q.InstanceID]).Review == "independent" {
		q = passingSessionReview(t, s, q, done)
		done = commitResult(t, s, "review", q)
	}
	st, _, e = s.load()
	if e != nil {
		t.Fatal(e)
	}
	i := st.Sessions[q.InstanceID]
	d := SessionDisposeRequest{Schema: SessionSchema, ProjectRef: q.ProjectRef, CheckoutID: q.CheckoutID, OperationID: newTestID(), InstanceID: q.InstanceID, ExpectedRevision: done.Revision, ExpectedRunRevision: st.Runs[i.Work[0].Request.OperationID].Revision, ResultID: q.OperationID, System: "SDPTOOL", Reason: "Publish the verified selected outcome and release its claim"}
	for _, c := range selectedSessionStep(i).Cards {
		v := SessionCardDisposition{ID: c.ID, Outcome: "completed", Summary: "Selected synthetic fixture outcome verified"}
		if c.Coverage == "contribution" {
			v.Outcome = "contribution-completed"
			v.Remaining = "Other planned contributions remain open"
		}
		d.Cards = append(d.Cards, v)
	}
	return s, d
}
func disposeResume(q SessionDisposeRequest, p SessionMutationResult) SessionResumeRequest {
	return SessionResumeRequest{SessionSchema, q.ProjectRef, q.CheckoutID, newTestID(), q.InstanceID, p.Revision, "Commit exact disposition effects"}
}
func assertDisposed(t *testing.T, s *Service, q SessionDisposeRequest) {
	t.Helper()
	v, e := s.ReadSession(q.InstanceID)
	if e != nil || v.Coverage != "current" || v.Steps[0].Stage != "done" || v.Steps[0].RunState != "completed" || !reflect.DeepEqual(v.Steps[0].CardDispositions, q.Cards) {
		t.Fatalf("view %+v: %v", v, e)
	}
	st, _, e := s.load()
	if e != nil {
		t.Fatal(e)
	}
	i := st.Sessions[q.InstanceID]
	if _, e = os.Lstat(filepath.Join(s.root, dispositionClaim(i))); !os.IsNotExist(e) {
		t.Fatal("claim retained", e)
	}
	if e = s.checkExternalSessionClaims(st.Runs[i.Work[0].Request.OperationID].Assignment); e != nil {
		t.Fatal("legacy claim still blocks", e)
	}
	if _, e = s.publicationFiles(); e != nil {
		t.Fatal(e)
	}
	for _, p := range []string{i.Context.Binding.Plan.Path, i.Context.Binding.Session.Path} {
		b := mustFile(t, filepath.Join(s.root, p))
		if !bytes.Contains(b, []byte("plan/Session closure is separate")) {
			t.Fatal("closure confused", p)
		}
	}
	for _, c := range q.Cards {
		name := "#001--Change--fixture.md"
		if c.ID == "KB-SDP-002" {
			name = "#002--Change--unrelated.md"
		}
		folder, state := "completed", "completed"
		if c.Outcome == "contribution-completed" {
			folder, state = "active", "in-progress"
		}
		b := mustFile(t, filepath.Join(s.root, "SDP/KanBan", folder, name))
		m, _ := metadata(b)
		if m["CardState"] != state || !bytes.Contains(b, []byte(c.Summary)) || (c.Remaining != "" && !bytes.Contains(b, []byte(c.Remaining))) {
			t.Fatal("card", m, string(b))
		}
		if folder == "completed" {
			if _, e = os.Lstat(filepath.Join(s.root, "SDP/KanBan/active", name)); !os.IsNotExist(e) {
				t.Fatal("old card exists")
			}
		}
	}
	trace := mustFile(t, filepath.Join(s.root, dispositionTrace))
	if bytes.Count(trace, []byte("EVT-SDPTOOL-PG-"+strings.ToUpper(q.OperationID))) != 1 {
		t.Fatal("trace duplicate/missing")
	}
}
func TestSessionDispositionWholeAndContribution(t *testing.T) {
	s, q := dispositionFixture(t, "", "", multiCardImplementation)
	st, _, _ := s.load()
	i := st.Sessions[q.InstanceID]
	ledger := mustFile(t, filepath.Join(s.root, managementLedger))
	p, e := s.DisposeSession(LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	if !bytes.Equal(ledger, mustFile(t, filepath.Join(s.root, managementLedger))) {
		t.Fatal("preparation wrote records")
	}
	v, e := s.ReadSession(q.InstanceID)
	if e != nil || v.Steps[0].Stage == "done" {
		t.Fatal(v, e)
	}
	if _, e = os.Stat(filepath.Join(s.root, dispositionClaim(i))); e != nil {
		t.Fatal(e)
	}
	retry, e := s.DisposeSession(LocalController(), q)
	if e != nil || !reflect.DeepEqual(p, retry) {
		t.Fatal(retry, e)
	}
	resume := disposeResume(q, p)
	done, e := s.ResumeSession(LocalController(), resume)
	if e != nil {
		t.Fatal(e)
	}
	again, e := s.ResumeSession(LocalController(), resume)
	if e != nil || !reflect.DeepEqual(done, again) {
		t.Fatal(again, e)
	}
	retry, e = s.DisposeSession(LocalController(), q)
	if e != nil || !reflect.DeepEqual(p, retry) {
		t.Fatal(retry, e)
	}
	if !bytes.HasPrefix(mustFile(t, filepath.Join(s.root, managementLedger)), ledger) {
		t.Fatal("history rewritten")
	}
	assertDisposed(t, s, q)
	q.Reason = "different retry"
	if _, e = s.DisposeSession(LocalController(), q); e == nil {
		t.Fatal("operation reuse")
	}
}
func TestSessionDispositionPolicyAndRefusals(t *testing.T) {
	for _, mode := range []string{"no-review", "owner", "caller", "revision", "run-revision", "result", "card", "remaining", "missing-card", "other-step", "candidate", "test", "skill", "source", "destination", "outside-link", "claim"} {
		t.Run(mode, func(t *testing.T) {
			configure := func(t *testing.T, s *Service, b *SessionBinding) {
				multiCardImplementation(t, s, b)
				if mode == "no-review" {
					b.Steps[0].Review = "not-required"
				}
				if mode == "owner" {
					b.Steps[0].OwnerAcceptance = "required"
				}
				if mode == "other-step" {
					other := b.Steps[0]
					other.ID = "later"
					other.DependsOn = []string{b.Steps[0].ID}
					b.Steps = append(b.Steps, other)
					evidenceFile(t, s, b.Session.Path, string(mustFile(t, filepath.Join(s.root, b.Session.Path)))+"\n## later\n\nSubsequent selected scope.\n")
				}
			}
			s, q := dispositionFixture(t, "", "", configure)
			caller := LocalController()
			switch mode {
			case "caller":
				caller = Caller{}
			case "revision":
				q.ExpectedRevision++
			case "run-revision":
				q.ExpectedRunRevision++
			case "result":
				q.ResultID = newTestID()
			case "card":
				q.Cards[0].Outcome = "completed"
				q.Cards[0].Remaining = ""
			case "remaining":
				q.Cards[0].Remaining = ""
			case "missing-card":
				q.Cards = q.Cards[:1]
			case "candidate":
				evidenceFile(t, s, "target.txt", "external candidate edit")
			case "test":
				evidenceFile(t, s, "evidence/session-result.log", "external test edit")
			case "skill":
				evidenceFile(t, s, "Skills/sdp-reviewer/SKILL.md", "external skill edit")
			case "source":
				evidenceFile(t, s, "SDP/KanBan/active/#001--Change--fixture.md", "external card edit")
			case "destination":
				evidenceFile(t, s, "SDP/KanBan/completed/#002--Change--unrelated.md", "destination collision")
			case "outside-link":
				evidenceFile(t, s, "SDP/unselected.md", "[card](KanBan/active/%23002--Change--unrelated.md)\n")
			case "claim":
				st, _, _ := s.load()
				i := st.Sessions[q.InstanceID]
				raw := mustFile(t, filepath.Join(s.root, dispositionClaim(i)))
				var claim map[string]any
				json.Unmarshal(raw, &claim)
				claim["operationId"] = newTestID()
				b, _ := json.Marshal(claim)
				evidenceFile(t, s, "SDP/Sessions/execution/claims/"+newTestID()+".json", string(b))
			}
			head := mustFile(t, filepath.Join(s.state, "HEAD"))
			p, e := s.DisposeSession(caller, q)
			if mode == "no-review" {
				if e != nil {
					t.Fatal(e)
				}
				if _, e = s.ResumeSession(caller, disposeResume(q, p)); e != nil {
					t.Fatal(e)
				}
				assertDisposed(t, s, q)
				return
			}
			if e == nil {
				t.Fatal("accepted", mode)
			}
			if !bytes.Equal(head, mustFile(t, filepath.Join(s.state, "HEAD"))) {
				t.Fatal("refusal wrote state")
			}
		})
	}
}
func TestSessionDispositionRecoveryEveryEffect(t *testing.T) {
	// Discover the actual effect count so additions retain interruption coverage.
	s, q := dispositionFixture(t, "", "", multiCardImplementation)
	if _, e := s.DisposeSession(LocalController(), q); e != nil {
		t.Fatal(e)
	}
	st, _, e := s.load()
	if e != nil {
		t.Fatal(e)
	}
	points := []string{"before-head", "after-head", "session-disposition-after-claim"}
	for n := range st.Sessions[q.InstanceID].Disposition.Effects.Changes {
		points = append(points, fmt.Sprintf("session-disposition-after-file:%d", n))
	}
	for _, point := range points {
		t.Run(point, func(t *testing.T) {
			s, q := dispositionFixture(t, "", "", multiCardImplementation)
			p, e := s.DisposeSession(LocalController(), q)
			if e != nil {
				t.Fatal(e)
			}
			r := disposeResume(q, p)
			hit := false
			s.fault = func(p string) error {
				if p == point {
					hit = true
					return errors.New("interrupted")
				}
				return nil
			}
			if _, e = s.ResumeSession(LocalController(), r); e == nil || !hit {
				t.Fatal("fault missed", e)
			}
			if point != "after-head" {
				if e = s.checkExternalSessionClaims(Assignment{}); e == nil {
					t.Fatal("pending disposition exposed scope")
				}
			}
			fresh, e := New(s.root, s.area, s.project, s.state)
			if e != nil {
				t.Fatal(e)
			}
			if _, e = fresh.ResumeSession(LocalController(), r); e != nil {
				t.Fatal(e)
			}
			if _, e = fresh.ResumeSession(LocalController(), r); e != nil {
				t.Fatal(e)
			}
			assertDisposed(t, fresh, q)
		})
	}
}
func TestSessionDispositionTamper(t *testing.T) {
	s, q := dispositionFixture(t, "", "", multiCardImplementation)
	if _, e := s.DisposeSession(LocalController(), q); e != nil {
		t.Fatal(e)
	}
	st, _, e := s.load()
	if e != nil {
		t.Fatal(e)
	}
	raw, _ := json.Marshal(st)
	for _, mode := range []string{"profile", "postimage", "preimage", "extra", "pins", "outcome", "receipt", "work", "terminal", "malformed-history"} {
		t.Run(mode, func(t *testing.T) {
			var bad snapshot
			json.Unmarshal(raw, &bad)
			i := bad.Sessions[q.InstanceID]
			d := i.Disposition
			switch mode {
			case "profile":
				bad.Schema = sessionResultStateSchema
			case "postimage":
				d.Effects.Changes[0].Content += "tamper"
			case "preimage":
				d.Before[i.Context.Binding.Session.Path] += "tamper"
			case "extra":
				d.Before["unrelated.md"] = "tamper"
			case "pins":
				d.Effects.Inputs = d.Effects.Inputs[1:]
			case "outcome":
				d.Request.Cards[0].Outcome = "completed"
			case "receipt":
				delete(bad.Receipts, q.OperationID)
			case "work":
				i.Work[2].Evidence = nil
			case "malformed-history":
				events, _ := publicationEvents([]byte(d.Before[managementLedger]))
				for _, ev := range events {
					if ev["subjectId"] == q.Cards[0].ID {
						ev["payload"] = []string{"invalid"}
					}
				}
				body := ""
				for _, ev := range events {
					b, _ := json.Marshal(ev)
					body += string(b) + "\n"
				}
				d.Before[managementLedger] = body
			case "terminal":
				i.Disposition = nil
				r := bad.Runs[i.Work[0].Request.OperationID]
				r.State = "completed"
				bad.Runs[r.ID] = r
			}
			bad.Sessions[i.ID] = i
			if validateSessionWork(bad) == nil && validateSessionDispositions(bad) == nil {
				t.Fatal("tamper accepted", mode)
			}
		})
	}
	if !slices.Contains([]string{sessionDispositionStateSchema}, st.Schema) {
		t.Fatal(st.Schema)
	}
}

func TestSessionDispositionMidPublicationEdit(t *testing.T) {
	for _, mode := range []string{"candidate", "selected-card", "destination", "trace"} {
		t.Run(mode, func(t *testing.T) {
			s, q := dispositionFixture(t, "", "", multiCardImplementation)
			p, e := s.DisposeSession(LocalController(), q)
			if e != nil {
				t.Fatal(e)
			}
			path := "target.txt"
			switch mode {
			case "selected-card":
				path = "SDP/KanBan/active/#001--Change--fixture.md"
			case "destination":
				path = "SDP/KanBan/completed/#002--Change--unrelated.md"
			case "trace":
				path = dispositionTrace
			}
			s.fault = func(point string) error {
				if point == "session-disposition-after-file:0" {
					evidenceFile(t, s, path, "External edit must survive\n")
				}
				return nil
			}
			r := disposeResume(q, p)
			if _, e = s.ResumeSession(LocalController(), r); e == nil {
				t.Fatal("external edit accepted")
			}
			s.fault = nil
			if _, e = s.ResumeSession(LocalController(), r); e == nil {
				t.Fatal("retry overwrote edit")
			}
			if string(mustFile(t, filepath.Join(s.root, path))) != "External edit must survive\n" {
				t.Fatal("edit lost")
			}
			st, _, e := s.load()
			if e != nil {
				t.Fatal(e)
			}
			i := st.Sessions[q.InstanceID]
			if i.Disposition.Effects.Status != "prepared" || st.Runs[i.Work[0].Request.OperationID].State == "completed" {
				t.Fatal("uncommitted result advanced")
			}
			if _, e = os.Stat(filepath.Join(s.root, dispositionClaim(i))); e != nil {
				t.Fatal("scope claim lost", e)
			}
			if _, e = os.Stat(filepath.Join(s.root, managementLedger) + ".pg-pending"); e != nil {
				t.Fatal("reservation lost", e)
			}
		})
	}
}
func TestSessionDispositionPreparationLostReply(t *testing.T) {
	for _, point := range []string{"before-head", "after-head"} {
		t.Run(point, func(t *testing.T) {
			s, q := dispositionFixture(t, "", "", nil)
			ledger := mustFile(t, filepath.Join(s.root, managementLedger))
			s.fault = func(p string) error {
				if p == point {
					return errors.New("lost preparation reply")
				}
				return nil
			}
			if _, e := s.DisposeSession(LocalController(), q); e == nil {
				t.Fatal("fault missed")
			}
			s.fault = nil
			if !bytes.Equal(ledger, mustFile(t, filepath.Join(s.root, managementLedger))) {
				t.Fatal("prepare wrote canonical records")
			}
			p, e := s.DisposeSession(LocalController(), q)
			if e != nil {
				t.Fatal(e)
			}
			if _, e = s.ResumeSession(LocalController(), disposeResume(q, p)); e != nil {
				t.Fatal(e)
			}
			assertDisposed(t, s, q)
		})
	}
}
