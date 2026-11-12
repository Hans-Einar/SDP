package activities

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Hans-Einar/SDP/SDPTool/projecthistory"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

func oid(n int) string { return fmt.Sprintf("%08d-1234-4123-8123-%012d", n, n) }
func lifecycleFixture(t *testing.T) (*Service, Caller) {
	t.Helper()
	s, c := fixture(t)
	if _, e := s.Create(request(t, s, oid(1)), c); e != nil {
		t.Fatal(e)
	}
	cat, e := s.Read()
	if e != nil {
		t.Fatal(e)
	}
	if _, e = s.AdoptAdmission(AdoptionRequest{OperationID: oid(2), ExpectedHistoryRevision: cat.HistoryRevision, Stores: []ExecutionStore{{StatePath: t.TempDir()}}}, c, knownEmpty()); e != nil {
		t.Fatal(e)
	}
	return s, c
}
func transitionRequest(t *testing.T, s *Service, n int, id, action string) TransitionRequest {
	t.Helper()
	c, e := s.Read()
	if e != nil {
		t.Fatal(e)
	}
	a, ok := c.Find(id)
	if !ok {
		t.Fatal(id)
	}
	return TransitionRequest{OperationID: oid(n), ExpectedHistoryRevision: c.HistoryRevision, ActivityID: id, ExpectedActivityEventID: a.LastEventID, ExpectedDescriptorRevision: a.DescriptorRef.SHA256, Action: action, ResumeCondition: "Owner selects resumption after held need is ready", Reason: "Owner lifecycle disposition"}
}
func TestLifecycleTransitionsAndCompletionEvidence(t *testing.T) {
	s, c := lifecycleFixture(t)
	id := "ACTIVITY-TEST-0001"
	for i, action := range []string{"activate", "hold", "resume"} {
		if _, e := s.Transition(transitionRequest(t, s, 10+i, id, action), c, knownEmpty()); e != nil {
			t.Fatal(action, e)
		}
	}
	q := transitionRequest(t, s, 20, id, "complete")
	if _, e := s.Transition(q, c, knownEmpty()); e == nil {
		t.Fatal("completion without explicit evidence")
	}
	os.WriteFile(filepath.Join(s.Area, "evidence.md"), []byte("| id | EVIDENCE-TEST-0001 |\n"), 0600)
	q.Outcome = "Owner accepts evidenced outcome"
	q.Evidence = []Reference{{Role: "verification", ID: "EVIDENCE-TEST-0001", Path: "evidence.md"}}
	r, e := s.Transition(q, c, knownEmpty())
	if e != nil {
		t.Fatal(e)
	}
	if _, e = s.Transition(transitionRequest(t, s, 21, id, "reopen"), c, nil); e != nil {
		t.Fatal(e)
	}
	retry, e := s.Transition(q, c, nil)
	if e != nil || retry != r {
		t.Fatalf("%+v %+v %v", r, retry, e)
	}
	if _, e = s.Transition(transitionRequest(t, s, 22, id, "cancel"), c, nil); e != nil {
		t.Fatal(e)
	}
}
func TestLifecycleUnknownMissingAndPendingQuiescence(t *testing.T) {
	s, c := lifecycleFixture(t)
	id := "ACTIVITY-TEST-0001"
	if _, e := s.Transition(transitionRequest(t, s, 10, id, "activate"), c, nil); e != nil {
		t.Fatal(e)
	}
	q := transitionRequest(t, s, 11, id, "hold")
	for _, inspector := range []ExecutionInspector{nil, admissionInspector{snapshot: ExecutionSnapshot{Known: false}}} {
		if _, e := s.Transition(q, c, inspector); e == nil {
			t.Fatal("unknown treated as empty")
		}
	}
	busy := admissionInspector{snapshot: ExecutionSnapshot{Known: true, InFlight: []ExecutionRef{{RunID: "actual-running"}}}}
	r, e := s.Transition(q, c, busy)
	if e != nil || !r.Pending {
		t.Fatalf("%+v %v", r, e)
	}
	ledger := filepath.Join(s.Area, "ProjectManagement/Ledger.ndjson")
	snap, e := projecthistory.Read(ledger)
	if e != nil {
		t.Fatal("reservation held while agent runs", e)
	}
	cat, e := s.Read()
	if e != nil || cat.Activities[0].State != "active" {
		t.Fatal("pretended run stopped", e)
	}
	create := CreateRequest{OperationID: oid(99), ExpectedHistoryRevision: snap.Revision, Title: "blocked", Owner: "owner", Mandate: "need", Reason: "new need"}
	if _, e = s.Create(create, c); e == nil {
		t.Fatal("new mutation bypassed pending request")
	}
	if rec, e := s.Recover(q.OperationID, c); e != nil || !rec.Pending {
		t.Fatal(rec, e)
	}
	changed := q
	changed.Reason = "changed"
	if _, e = s.Transition(changed, c, knownEmpty()); !errors.Is(e, projecthistory.ErrConflict) {
		t.Fatal(e)
	}
	// Unrelated capture may advance global history while a lifecycle waits.
	ev := []byte(`{"schemaVersion":"1.0","eventId":"EVT-OTHER-1","eventType":"work-recorded","occurredAt":"2026-11-12T00:00:00Z","actor":"capture","commit":null,"payload":{}}`)
	if _, e = projecthistory.Append(ledger, snap.Revision, ev); e != nil {
		t.Fatal(e)
	}
	done, e := s.RecoverLifecycle(q.OperationID, c, knownEmpty())
	if e != nil || done.Pending {
		t.Fatal(done, e)
	}
	cat, e = s.Read()
	if e != nil || cat.Activities[0].State != "on-hold" {
		t.Fatal(cat, e)
	}
	tx, e := projecthistory.BeginTransaction(ledger)
	if e != nil {
		t.Fatal(e)
	}
	defer tx.Close()
	view, e := s.ReadAdmission(tx, "")
	if e != nil || view.Pending != nil {
		t.Fatal(view, e)
	}
}
func switchRequest(t *testing.T, s *Service, n int) TransitionRequest {
	t.Helper()
	q := transitionRequest(t, s, n, "ACTIVITY-TEST-0001", "switch")
	c, _ := s.Read()
	a, _ := c.Find("ACTIVITY-TEST-0002")
	q.TargetActivityID = a.Descriptor.ID
	q.ExpectedTargetEventID = a.LastEventID
	q.ExpectedTargetDescriptorRevision = a.DescriptorRef.SHA256
	return q
}
func TestAtomicSwitchCrashRecoveryAndIncompleteBatch(t *testing.T) {
	for _, point := range []string{"after-intent", "after-reserve", "before-history", "after-history", "after-release"} {
		t.Run(point, func(t *testing.T) {
			s, c := lifecycleFixture(t)
			if _, e := s.Create(request(t, s, oid(3)), c); e != nil {
				t.Fatal(e)
			}
			if _, e := s.Transition(transitionRequest(t, s, 10, "ACTIVITY-TEST-0001", "activate"), c, nil); e != nil {
				t.Fatal(e)
			}
			q := switchRequest(t, s, 11)
			s.fault = func(p string) error {
				if p == point {
					return errors.New("crash")
				}
				return nil
			}
			if _, e := s.Transition(q, c, knownEmpty()); e == nil {
				t.Fatal("fault not reached")
			}
			s.fault = nil
			if _, e := s.Recover(q.OperationID, c); e != nil {
				t.Fatal(e)
			}
			cat, e := s.Read()
			if e != nil || cat.Activities[0].State != "on-hold" || cat.Activities[1].State != "active" {
				t.Fatal(cat, e)
			}
			snap, e := projecthistory.Read(filepath.Join(s.Area, "ProjectManagement/Ledger.ndjson"))
			if e != nil {
				t.Fatal(e)
			}
			lines := bytes.Split(bytes.TrimSpace(snap.Bytes), []byte{'\n'})
			partial := bytes.Join(lines[:len(lines)-1], []byte{'\n'})
			partial = append(partial, '\n')
			if _, e := Replay(s.Area, s.ProjectID, projecthistory.Snapshot{Bytes: partial, Revision: digest(partial)}); e == nil {
				t.Fatal("accepted incomplete switch")
			}
			var last map[string]any
			json.Unmarshal(lines[len(lines)-1], &last)
			last["subjectId"] = "ACTIVITY-TEST-0001"
			lines[len(lines)-1], _ = json.Marshal(last)
			bad := append(bytes.Join(lines, []byte{'\n'}), '\n')
			if _, e := Replay(s.Area, s.ProjectID, projecthistory.Snapshot{Bytes: bad, Revision: digest(bad)}); e == nil {
				t.Fatal("accepted forged switch")
			}
		})
	}
}
func TestConcurrentActivationSingleWinner(t *testing.T) {
	s, c := lifecycleFixture(t)
	if _, e := s.Create(request(t, s, oid(3)), c); e != nil {
		t.Fatal(e)
	}
	q1 := transitionRequest(t, s, 10, "ACTIVITY-TEST-0001", "activate")
	q2 := transitionRequest(t, s, 11, "ACTIVITY-TEST-0002", "activate")
	var wg sync.WaitGroup
	errs := make(chan error, 2)
	for _, q := range []TransitionRequest{q1, q2} {
		wg.Add(1)
		go func(q TransitionRequest) { defer wg.Done(); _, e := s.Transition(q, c, nil); errs <- e }(q)
	}
	wg.Wait()
	close(errs)
	wins := 0
	for e := range errs {
		if e == nil {
			wins++
		}
	}
	if wins != 1 {
		t.Fatal(wins)
	}
	cat, e := s.Read()
	if e != nil {
		t.Fatal(e)
	}
	active := 0
	for _, a := range cat.Activities {
		if a.State == "active" {
			active++
		}
	}
	if active != 1 {
		t.Fatal(active)
	}
}
func TestUpdatePreservesDescriptorHistoryAndValidatesReferences(t *testing.T) {
	s, c := fixture(t)
	r, e := s.Create(request(t, s, oid(1)), c)
	if e != nil {
		t.Fatal(e)
	}
	cat, _ := s.Read()
	a := cat.Activities[0]
	q := UpdateRequest{OperationID: oid(2), ExpectedHistoryRevision: cat.HistoryRevision, ActivityID: r.ActivityID, ExpectedActivityEventID: a.LastEventID, ExpectedDescriptorRevision: a.DescriptorRef.SHA256, Title: "Changed", Owner: "owner", Reason: "Explicit scope", Scope: []string{"Source"}, References: []Reference{}, SourceCards: []string{}, Systems: []string{}}
	if _, e = s.Update(q, c); e != nil {
		t.Fatal(e)
	}
	cat, e = s.Read()
	if e != nil || cat.Activities[0].Descriptor.Title != "Changed" {
		t.Fatal(cat, e)
	}
	if _, e = s.Recover(oid(1), c); e != nil {
		t.Fatal("original receipt failed after update", e)
	}
	q = UpdateRequest{OperationID: oid(3), ExpectedHistoryRevision: cat.HistoryRevision, ActivityID: r.ActivityID, ExpectedActivityEventID: cat.Activities[0].LastEventID, ExpectedDescriptorRevision: cat.Activities[0].DescriptorRef.SHA256, Title: "Changed", Owner: "owner", Reason: "bad ref", Scope: []string{}, References: []Reference{{Role: "session", ID: "SESSION-TEST-0999", Path: "authority.md"}}, SourceCards: []string{}}
	if _, e = s.Update(q, c); e == nil {
		t.Fatal("wrong Session identity accepted")
	}
	q.References = []Reference{}
	q.SourceCards = []string{"KB-TEST-999"}
	if _, e = s.Update(q, c); e == nil {
		t.Fatal("invented card accepted")
	}
}

func evidence(t *testing.T, s *Service) Reference {
	t.Helper()
	r := Reference{Role: "verification", ID: "EVIDENCE-TEST-0001", Path: "evidence.md"}
	if e := os.WriteFile(filepath.Join(s.Area, r.Path), []byte("| id | "+r.ID+" |\n"), 0600); e != nil {
		t.Fatal(e)
	}
	return r
}
func updateRequest(t *testing.T, s *Service, n int, id string) UpdateRequest {
	t.Helper()
	cat, e := s.Read()
	if e != nil {
		t.Fatal(e)
	}
	a, ok := cat.Find(id)
	if !ok {
		t.Fatal(id)
	}
	return UpdateRequest{OperationID: oid(n), ExpectedHistoryRevision: cat.HistoryRevision, ActivityID: id, ExpectedActivityEventID: a.LastEventID, ExpectedDescriptorRevision: a.DescriptorRef.SHA256, Title: a.Descriptor.Title, Owner: a.Descriptor.Owner, Reason: "Explicit update", Scope: []string{}, References: []Reference{}, SourceCards: []string{}, Systems: []string{}, Dependencies: a.Descriptor.Dependencies}
}
func TestStableDependenciesCyclesAndExplicitSatisfaction(t *testing.T) {
	s, c := lifecycleFixture(t)
	first := "ACTIVITY-TEST-0001"
	if _, e := s.Transition(transitionRequest(t, s, 10, first, "activate"), c, nil); e != nil {
		t.Fatal(e)
	}
	q := request(t, s, oid(3))
	q.Dependencies = []Dependency{{ActivityID: first}}
	r, e := s.Create(q, c)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = s.Transition(transitionRequest(t, s, 11, r.ActivityID, "activate"), c, nil); e == nil {
		t.Fatal("unsatisfied dependency activated")
	}
	complete := transitionRequest(t, s, 12, first, "complete")
	complete.Outcome = "Accepted result"
	complete.Evidence = []Reference{evidence(t, s)}
	completed, e := s.Transition(complete, c, knownEmpty())
	if e != nil {
		t.Fatal(e)
	}
	if _, e = s.Transition(transitionRequest(t, s, 13, r.ActivityID, "activate"), c, nil); e == nil {
		t.Fatal("completion inferred explicit satisfaction")
	}
	u := updateRequest(t, s, 14, r.ActivityID)
	u.Dependencies = []Dependency{{ActivityID: first, SatisfiedByEventID: completed.EventID}}
	if _, e = s.Update(u, c); e != nil {
		t.Fatal(e)
	}
	if _, e = s.Transition(transitionRequest(t, s, 15, r.ActivityID, "activate"), c, nil); e != nil {
		t.Fatal(e)
	}
	u = updateRequest(t, s, 16, first)
	u.Title = "Renamed accepted need"
	if _, e = s.Update(u, c); e != nil {
		t.Fatal("title altered dependency identity", e)
	}
	if _, e = s.Transition(transitionRequest(t, s, 17, first, "reopen"), c, nil); e == nil {
		t.Fatal("reopened active dependency")
	}
	u = updateRequest(t, s, 18, first)
	u.Dependencies = []Dependency{{ActivityID: r.ActivityID}}
	if _, e = s.Update(u, c); e == nil {
		t.Fatal("cycle accepted")
	}
	for i, deps := range [][]Dependency{{{ActivityID: "ACTIVITY-TEST-9999"}}, {{ActivityID: first, SatisfiedByEventID: "EVT-UNKNOWN"}}, {{ActivityID: "ACTIVITY-TEST-0003"}}} {
		q = request(t, s, oid(30+i))
		q.Dependencies = deps
		if _, e = s.Create(q, c); e == nil {
			t.Fatal("invalid dependency", deps)
		}
	}
}
func TestContributionsUnavailableReferencesAndPinnedInputs(t *testing.T) {
	s, c := lifecycleFixture(t)
	if e := os.WriteFile(filepath.Join(s.Area, "card.md"), []byte("| id | KB-TEST-001 |\n| CardState | backlog |\n"), 0600); e != nil {
		t.Fatal(e)
	}
	u := updateRequest(t, s, 3, "ACTIVITY-TEST-0001")
	u.References = []Reference{{Role: "card", ID: "KB-TEST-001", Path: "card.md", Membership: "contribution"}, {Role: "external-design", ID: "DESIGN-OTHER-0001", Path: "Design/Missing.md", ProjectID: "OTHER", Availability: "unavailable", UnavailableReason: "Collaborating repository is not installed"}}
	if _, e := s.Update(u, c); e != nil {
		t.Fatal(e)
	}
	if _, e := s.Transition(transitionRequest(t, s, 10, u.ActivityID, "activate"), c, nil); e != nil {
		t.Fatal(e)
	}
	q := transitionRequest(t, s, 11, u.ActivityID, "complete")
	q.Outcome = "Selected contribution accepted"
	q.Evidence = []Reference{evidence(t, s)}
	if _, e := s.Transition(q, c, knownEmpty()); e == nil {
		t.Fatal("missing selected contribution disposition")
	}
	q.ContributionDispositions = []ContributionDisposition{{ReferenceID: "KB-TEST-001", Disposition: "completed", Outcome: "The selected contribution is accepted; shared card remains open"}}
	if _, e := s.Transition(q, c, knownEmpty()); e != nil {
		t.Fatal(e)
	}
	// Reviewed inputs require the exact referenced bytes, not only a matching path.
	b, _ := s.read("card.md")
	ref := Reference{Role: "card", ID: "KB-TEST-001", Path: "card.md", Revision: digest(b)}
	if e := validateRefs(s, []Reference{ref}); e != nil {
		t.Fatal(e)
	}
	os.WriteFile(filepath.Join(s.Area, "card.md"), append(b, []byte("changed\n")...), 0600)
	if e := validateRefs(s, []Reference{ref}); e == nil {
		t.Fatal("changed reviewed input accepted")
	}
	ref.ProjectID = "OTHER"
	ref.ID = "KB-OTHER-001"
	ref.Availability = "unavailable"
	ref.UnavailableReason = "missing"
	if e := validateRefs(s, []Reference{ref}); e != nil {
		t.Fatal(e)
	}
}
func TestDescriptorReplacementCrashRecovery(t *testing.T) {
	for _, point := range []string{"after-file-1", "after-file-2", "before-history", "after-history"} {
		t.Run(point, func(t *testing.T) {
			s, c := fixture(t)
			if _, e := s.Create(request(t, s, oid(1)), c); e != nil {
				t.Fatal(e)
			}
			q := updateRequest(t, s, 2, "ACTIVITY-TEST-0001")
			q.Title = "Updated"
			s.fault = func(p string) error {
				if p == point {
					return errors.New("crash")
				}
				return nil
			}
			if _, e := s.Update(q, c); e == nil {
				t.Fatal("fault not reached")
			}
			s.fault = nil
			if _, e := s.Recover(q.OperationID, c); e != nil {
				t.Fatal(e)
			}
			cat, e := s.Read()
			if e != nil || cat.Activities[0].Descriptor.Title != "Updated" || len(cat.Activities[0].Events) != 2 {
				t.Fatal(cat, e)
			}
		})
	}
}

func TestManual01DescriptorUpgradeRetainsExactBytes(t *testing.T) {
	realArea := filepath.Join("..", "..", "SDP")
	s, c := fixture(t)
	s.ProjectID = "SDP"
	ledger, e := os.ReadFile(filepath.Join(realArea, "ProjectManagement/Ledger.ndjson"))
	if e != nil {
		t.Fatal(e)
	}
	if e = s.write("ProjectManagement/Ledger.ndjson", ledger); e != nil {
		t.Fatal(e)
	}
	for _, home := range []string{"Activity-0001--SDPTool-Restructuring", "Activity-0002--Activity-Management"} {
		p := "Activities/" + home + "/activity.yaml"
		b, e := os.ReadFile(filepath.Join(realArea, filepath.FromSlash(p)))
		if e != nil {
			t.Fatal(e)
		}
		if e = s.write(p, b); e != nil {
			t.Fatal(e)
		}
	}
	before, e := s.Read()
	if e != nil {
		t.Fatal(e)
	}
	old := before.Activities[1]
	oldBytes, _ := s.read(old.DescriptorRef.Path)
	q := updateRequest(t, s, 80, old.Descriptor.ID)
	q.Title = "Explicitly adopted manual need metadata"
	q.Scope = []string{"SDPTool"}
	if _, e = s.Update(q, c); e != nil {
		t.Fatal(e)
	}
	after, e := s.Read()
	if e != nil || after.Activities[1].Descriptor.Schema != "sdp-activity/0.2" || after.Activities[1].State != "active" {
		t.Fatal(after, e)
	}
	retained, e := s.read("Activities/.descriptors/" + old.DescriptorRef.SHA256 + ".yaml")
	if e != nil || !bytes.Equal(oldBytes, retained) {
		t.Fatal("manual descriptor bytes not preserved", e)
	}
	stream, _ := s.read("ProjectManagement/Ledger.ndjson")
	if !bytes.HasPrefix(stream, ledger) {
		t.Fatal("manual history bytes rewritten")
	}
}
func TestWholeScopeDoesNotAcceptOpenRecordAsCompleted(t *testing.T) {
	s, c := lifecycleFixture(t)
	os.WriteFile(filepath.Join(s.Area, "plan.md"), []byte("| id | PLAN-TEST-0001 |\n| state | active |\n"), 0600)
	q := updateRequest(t, s, 3, "ACTIVITY-TEST-0001")
	q.References = []Reference{{Role: "plan", ID: "PLAN-TEST-0001", Path: "plan.md", Membership: "whole-scope"}}
	if _, e := s.Update(q, c); e != nil {
		t.Fatal(e)
	}
	if _, e := s.Transition(transitionRequest(t, s, 10, q.ActivityID, "activate"), c, nil); e != nil {
		t.Fatal(e)
	}
	done := transitionRequest(t, s, 11, q.ActivityID, "complete")
	done.Outcome = "Need disposition"
	done.Evidence = []Reference{evidence(t, s)}
	done.ContributionDispositions = []ContributionDisposition{{ReferenceID: "PLAN-TEST-0001", Disposition: "completed", Outcome: "Claim complete"}}
	if _, e := s.Transition(done, c, knownEmpty()); e == nil {
		t.Fatal("open whole-scope record accepted")
	}
	os.WriteFile(filepath.Join(s.Area, "plan.md"), []byte("| id | PLAN-TEST-0001 |\n| state | completed |\n"), 0600)
	if _, e := s.Transition(done, c, knownEmpty()); e != nil {
		t.Fatal(e)
	}
}

func TestAdoptedActiveUpdateRequiresObservedHoldForContractChanges(t *testing.T) {
	for _, kind := range []string{"scope", "references", "dependencies"} {
		t.Run(kind, func(t *testing.T) {
			s, c := lifecycleFixture(t)
			id := "ACTIVITY-TEST-0001"
			if _, e := s.Create(request(t, s, oid(3)), c); e != nil {
				t.Fatal(e)
			}
			if _, e := s.Transition(transitionRequest(t, s, 10, id, "activate"), c, nil); e != nil {
				t.Fatal(e)
			}
			s.write("Selected.md", []byte("| id | PLAN-TEST-0001 |\n"))
			q := updateRequest(t, s, 11, id)
			cat, _ := s.Read()
			a, _ := cat.Find(id)
			q.Scope = append([]string{}, a.Descriptor.Scope...)
			q.References = append([]Reference{}, a.Descriptor.References...)
			switch kind {
			case "scope":
				q.Scope = append(q.Scope, "Another")
			case "references":
				q.References = append(q.References, Reference{Role: "plan", ID: "PLAN-TEST-0001", Path: "Selected.md"})
			case "dependencies":
				q.Dependencies = []Dependency{{ActivityID: "ACTIVITY-TEST-0002"}}
			}
			before, _ := s.read("ProjectManagement/Ledger.ndjson")
			if _, e := s.Update(q, c); e == nil || !strings.Contains(e.Error(), "hold the active Activity") {
				t.Fatal("adopted active contract change was not refused with hold instruction", e)
			}
			after, _ := s.read("ProjectManagement/Ledger.ndjson")
			if !bytes.Equal(before, after) {
				t.Fatal("refused contract change appended history")
			}
			busy := admissionInspector{snapshot: ExecutionSnapshot{Known: true, InFlight: []ExecutionRef{{RunID: "actual-running"}}}}
			hold := transitionRequest(t, s, 12, id, "hold")
			if r, e := s.Transition(hold, c, busy); e != nil || !r.Pending {
				t.Fatal(r, e)
			}
			if _, e := s.Update(q, c); e == nil {
				t.Fatal("pending hold permitted contract rewrite")
			}
			if r, e := s.Transition(hold, c, knownEmpty()); e != nil || r.Pending {
				t.Fatal(r, e)
			}
			cat, _ = s.Read()
			a, _ = cat.Find(id)
			q.ExpectedHistoryRevision = cat.HistoryRevision
			q.ExpectedActivityEventID = a.LastEventID
			q.ExpectedDescriptorRevision = a.DescriptorRef.SHA256
			if _, e := s.Update(q, c); e != nil {
				t.Fatal("observed-quiescent held contract update refused", e)
			}
		})
	}
}
func TestAdoptedActiveUpdateAllowsMetadataAndSetReordering(t *testing.T) {
	s, c := lifecycleFixture(t)
	id := "ACTIVITY-TEST-0001"
	for _, name := range []string{"First", "Second"} {
		if e := s.write(name+".md", []byte("| id | PLAN-TEST-"+map[string]string{"First": "0001", "Second": "0002"}[name]+" |\n")); e != nil {
			t.Fatal(e)
		}
	}
	q := updateRequest(t, s, 3, id)
	q.Scope = []string{"Source", "Another"}
	q.References = []Reference{{Role: "plan", ID: "PLAN-TEST-0001", Path: "First.md"}, {Role: "plan", ID: "PLAN-TEST-0002", Path: "Second.md"}}
	if _, e := s.Update(q, c); e != nil {
		t.Fatal(e)
	}
	if _, e := s.Transition(transitionRequest(t, s, 10, id, "activate"), c, nil); e != nil {
		t.Fatal(e)
	}
	cat, _ := s.Read()
	a, _ := cat.Find(id)
	q = updateRequest(t, s, 11, id)
	q.Scope = append([]string{}, a.Descriptor.Scope...)
	q.References = append([]Reference{}, a.Descriptor.References...)
	q.Title = "Metadata rename"
	q.Owner = "New accountable owner"
	if _, e := s.Update(q, c); e != nil {
		t.Fatal("metadata update refused", e)
	}
	q = updateRequest(t, s, 12, id)
	q.Scope = []string{"Another", "Source"}
	q.References = []Reference{{Role: "plan", ID: "PLAN-TEST-0002", Path: "Second.md"}, {Role: "plan", ID: "PLAN-TEST-0001", Path: "First.md"}}
	if _, e := s.Update(q, c); e != nil {
		t.Fatal("semantically unchanged reordered contract refused", e)
	}
	before := Descriptor{Dependencies: []Dependency{{ActivityID: "ACTIVITY-TEST-0002", SatisfiedByEventID: "EVT-TWO"}, {ActivityID: "ACTIVITY-TEST-0003", SatisfiedByEventID: "EVT-THREE"}}}
	after := UpdateRequest{Scope: []string{}, References: []Reference{}, Dependencies: []Dependency{before.Dependencies[1], before.Dependencies[0]}}
	if executionContractChanged(before, after) {
		t.Fatal("dependency order or nil/empty representation changed execution semantics")
	}
	after.Dependencies[0].SatisfiedByEventID = "EVT-DIFFERENT"
	if !executionContractChanged(before, after) {
		t.Fatal("changed dependency satisfaction was ignored")
	}
}
