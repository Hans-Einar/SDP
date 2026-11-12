package governance

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/Hans-Einar/SDP/SDPTool/activities"
	"github.com/Hans-Einar/SDP/SDPTool/projecthistory"
)

func adoptedGovernance(t *testing.T) (*Service, *activities.Service, StartRequest) {
	t.Helper()
	s, sel, q := fixture(t)
	// Legacy fixture intentionally omits generic fields. Adoption requires a valid
	// canonical stream; preserve the same domain payload and current projections.
	ledger := filepath.Join(s.root, managementLedger)
	raw, err := os.ReadFile(ledger)
	if err != nil {
		t.Fatal(err)
	}
	var valid []byte
	for _, line := range bytes.Split(raw, []byte{'\n'}) {
		if len(line) == 0 {
			continue
		}
		var ev map[string]any
		if err = json.Unmarshal(line, &ev); err != nil {
			t.Fatal(err)
		}
		ev["eventId"] = "EVT-" + ev["eventId"].(string)
		ev["occurredAt"] = "2026-11-12T00:00:00Z"
		ev["actor"] = "fixture"
		ev["commit"] = nil
		b, _ := json.Marshal(ev)
		valid = append(valid, b...)
		valid = append(valid, '\n')
	}
	if err = os.WriteFile(ledger, valid, 0600); err != nil {
		t.Fatal(err)
	}
	bind(t, s, &q)
	domain := activities.New(s.area, "TEST")
	caller := activities.Caller{Actor: "fixture-owner", Authority: activities.Authority{Role: "project-owner", Source: "authority.md"}}
	snap, err := projecthistory.Read(ledger)
	if err != nil {
		t.Fatal(err)
	}
	created, err := domain.Create(activities.CreateRequest{OperationID: newTestID(), ExpectedHistoryRevision: snap.Revision, Title: "Selected need", Owner: "owner", Mandate: "Need", Reason: "Owner need", Scope: []string{"target.txt"}, References: []activities.Reference{{Role: "session", ID: "SESSION-TEST-0001", Path: "Sessions/session.md"}}}, caller)
	if err != nil {
		t.Fatal(err)
	}
	from := "planned"
	prev := created.EventID
	event := activities.Event{SchemaVersion: "1.0", EventID: "EVT-ACTIVITY-TEST-900001", EventType: "x-activity:transitioned", OccurredAt: "2026-11-12T00:00:00Z", Actor: "fixture-owner", SubjectID: created.ActivityID, Payload: activities.Payload{SchemaVersion: "0.2", ProjectID: "TEST", OperationID: newTestID(), RecordedAt: "2026-11-12T00:00:00Z", Authority: caller.Authority, PreviousEventID: &prev, From: &from, To: "active", Reason: "Owner selects work", Descriptor: activities.DescriptorReference{Path: created.Path + "/activity.yaml", SHA256: created.DescriptorRevision}, Links: []string{}, RequestDigest: hash([]byte("activate"))}}
	b, _ := json.Marshal(event)
	appendResult, err := projecthistory.Append(ledger, created.HistoryRevision, b)
	if err != nil {
		t.Fatal(err)
	}
	_, err = domain.AdoptAdmission(activities.AdoptionRequest{OperationID: newTestID(), ExpectedHistoryRevision: appendResult.Revision, Stores: []activities.ExecutionStore{{StatePath: s.state}}}, caller, ActivityExecutionInspector(s.root, s.area, s.project))
	if err != nil {
		t.Fatal(err)
	}
	context, err := s.Context(sel)
	if err != nil {
		t.Fatal(err)
	}
	q.InputRevision = context.Revision
	q.OperationID = newTestID()
	q.RequestID = newTestID()
	return s, domain, q
}
func pendingActivity(t *testing.T, s *Service, domain *activities.Service) activities.PendingLifecycle {
	t.Helper()
	cat, err := domain.Read()
	if err != nil {
		t.Fatal(err)
	}
	pending := activities.PendingLifecycle{OperationID: newTestID(), ActivityID: cat.Activities[0].Descriptor.ID, Action: "hold"}
	tx, err := s.canonicalTransaction()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Close()
	if err = domain.SetPendingLifecycle(tx, "", pending); err != nil {
		t.Fatal(err)
	}
	return pending
}
func TestActivityAdmissionDirectStartAndScope(t *testing.T) {
	s, _, q := adoptedGovernance(t)
	bad := q
	bad.Assignment.WritablePaths = []string{"outside.txt"}
	if _, err := s.Start(LocalController(), bad); err == nil {
		t.Fatal("direct domain start bypassed scope")
	}
	r, err := s.Start(LocalController(), q)
	if err != nil || r.Run.Activity == nil {
		t.Fatal(r, err)
	}
	if r.Run.Activity.ActivityID != "ACTIVITY-TEST-0001" {
		t.Fatal(r.Run.Activity)
	}
	// Capture an actual adopted result, rather than validating a nil-binding
	// legacy fixture against the strict additionalProperties=false wire schema.
	capture, err := json.Marshal([]map[string]any{{"definition": "Run", "value": r.Run}, {"definition": "Result", "value": r}})
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "adopted-wire.json")
	if err = os.WriteFile(path, capture, 0600); err != nil {
		t.Fatal(err)
	}
	if output, err := exec.Command("python3", "check_schema.py", path).CombinedOutput(); err != nil {
		t.Fatalf("adopted Run wire schema: %v %s", err, output)
	}
	st, _, err := s.load()
	if err != nil || st.Runs[r.Run.ID].Activity == nil {
		t.Fatal("admission basis not durable", err)
	}
}
func TestActivityPendingPreservesExactSendingAndBlocksNewClaims(t *testing.T) {
	s, domain, q := adoptedGovernance(t)
	r, err := s.Start(LocalController(), q)
	if err != nil {
		t.Fatal(err)
	}
	v := OwnerInput{ID: newTestID(), AttemptID: newTestID(), Epoch: newTestID(), Kind: "launch", Text: "Owner work", Payload: "Exact payload", ReceivedAt: time.Now().UTC().Format(time.RFC3339Nano)}
	input := func(action string) (Result, error) {
		return s.RecordInput(LocalController(), InputRequest{Schema: Schema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, RunID: r.Run.ID, OperationID: newTestID(), ExpectedRevision: r.Run.Revision, InputRevision: r.Run.InputRevision, Action: action, Input: v})
	}
	r, err = input("prepare")
	if err != nil {
		t.Fatal(err)
	}
	v.State = "prepared"
	r, err = input("claim")
	if err != nil {
		t.Fatal(err)
	}
	v.State = "sending"
	pendingActivity(t, s, domain)
	if err = s.CheckActivityDispatch(LocalController(), r.Run.ID, v.ID); err != nil {
		t.Fatal("already-issued exact send lost settlement", err)
	}
	q.OperationID = newTestID()
	q.RequestID = newTestID()
	if _, err = s.Start(LocalController(), q); err == nil {
		t.Fatal("pending hold admitted new run")
	}
	if err = s.CheckActivityDispatch(LocalController(), r.Run.ID, newTestID()); err == nil {
		t.Fatal("pending permitted new host send")
	}
	observed, err := ActivityExecutionInspector(s.root, s.area, s.project).InspectExecutions([]activities.ExecutionStore{{StatePath: s.state}})
	if err != nil || !observed.Known || len(observed.InFlight) != 1 {
		t.Fatal(observed, err)
	}
	// Failure/uncertainty settlement must remain possible while the hold drains.
	r, err = input("uncertain")
	if err != nil {
		t.Fatal(err)
	}
	observed, err = ActivityExecutionInspector(s.root, s.area, s.project).InspectExecutions([]activities.ExecutionStore{{StatePath: s.state}})
	if err != nil || len(observed.InFlight) != 1 {
		t.Fatal(observed, err)
	}
}
func TestActivityPendingBeforeClaimAndAcrossStores(t *testing.T) {
	s, domain, q := adoptedGovernance(t)
	r, err := s.Start(LocalController(), q)
	if err != nil {
		t.Fatal(err)
	}
	v := OwnerInput{ID: newTestID(), AttemptID: newTestID(), Epoch: newTestID(), Kind: "launch", Text: "Owner work", Payload: "Exact payload", ReceivedAt: time.Now().UTC().Format(time.RFC3339Nano)}
	r, err = s.RecordInput(LocalController(), InputRequest{Schema: Schema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, RunID: r.Run.ID, OperationID: newTestID(), ExpectedRevision: r.Run.Revision, InputRevision: r.Run.InputRevision, Action: "prepare", Input: v})
	if err != nil {
		t.Fatal(err)
	}
	// Explicit second store uses the same canonical gate, not only first-store state.
	other, err := New(s.root, s.area, s.project, filepath.Join(t.TempDir(), "other-state"))
	if err != nil {
		t.Fatal(err)
	}
	bound, err := other.Bind(LocalController())
	if err != nil {
		t.Fatal(err)
	}
	otherRequest := q
	otherRequest.CheckoutID = bound.CheckoutID
	otherRequest.OperationID = newTestID()
	otherRequest.RequestID = newTestID()
	if _, err = other.Start(LocalController(), otherRequest); err != nil {
		t.Fatal(err)
	}
	tx, err := s.canonicalTransaction()
	if err != nil {
		t.Fatal(err)
	}
	view, err := domain.ReadAdmission(tx, "")
	tx.Close()
	if err != nil || len(view.Stores) != 2 {
		t.Fatal(view.Stores, err)
	}
	pendingActivity(t, s, domain)
	v.State = "prepared"
	_, err = s.RecordInput(LocalController(), InputRequest{Schema: Schema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, RunID: r.Run.ID, OperationID: newTestID(), ExpectedRevision: r.Run.Revision, InputRevision: r.Run.InputRevision, Action: "claim", Input: v})
	if err == nil {
		t.Fatal("claim crossed pending lifecycle")
	}
	otherRequest.OperationID = newTestID()
	otherRequest.RequestID = newTestID()
	if _, err = other.Start(LocalController(), otherRequest); err == nil {
		t.Fatal("other store bypassed pending lifecycle")
	}
	if err = os.Remove(filepath.Join(other.state, "HEAD")); err != nil {
		t.Fatal(err)
	}
	if _, err = ActivityExecutionInspector(s.root, s.area, s.project).InspectExecutions(view.Stores); err == nil {
		t.Fatal("lost store treated as empty")
	}
}

func TestActivityLocalSessionUsesCatalogAndExistingExecutionBinding(t *testing.T) {
	service, context, _ := sessionFixture(t)
	domain := activities.New(service.area, "TEST")
	caller := activities.Caller{Actor: "owner", Authority: activities.Authority{Role: "project-owner", Source: "authority.md"}}
	catalog, err := domain.Read()
	if err != nil {
		t.Fatal(err)
	}
	created, err := domain.Create(activities.CreateRequest{OperationID: newTestID(), ExpectedHistoryRevision: catalog.HistoryRevision, Title: "Local", Owner: "owner", Mandate: "Need", Reason: "Owner request", Scope: []string{"target.txt"}}, caller)
	if err != nil {
		t.Fatal(err)
	}
	session, err := domain.CreateSession(activities.CreateSessionRequest{OperationID: newTestID(), ExpectedHistoryRevision: created.HistoryRevision, ActivityID: created.ActivityID, ExpectedActivityEventID: created.EventID, ExpectedDescriptorRevision: created.DescriptorRevision, Title: "Execution", Objective: "Prepared work", Phase: "04-Design", Reason: "Owner preparation"}, caller)
	if err != nil {
		t.Fatal(err)
	}
	// Record creation itself did not construct an execution binding. This test now
	// authors the actual existing PG coverage before asking PG to inspect it.
	binding := context.Binding
	oldID := binding.Session.ID
	binding.Session.ID = session.SessionID
	binding.Session.Path = "SDP/" + session.Path
	old, err := os.ReadFile(filepath.Join(service.root, context.Binding.Session.Path))
	if err != nil {
		t.Fatal(err)
	}
	ready := bytes.ReplaceAll(old, []byte(oldID), []byte(session.SessionID))
	evidenceFile(t, service, binding.Session.Path, string(ready))
	raw, _ := json.Marshal(binding)
	bindingPath := "SDP/Sessions/execution/" + session.SessionID + ".json"
	evidenceFile(t, service, bindingPath, string(raw))
	got, err := service.SessionContext(SessionContextRequest{BindingPath: bindingPath})
	if err != nil || got.Binding.Session.Path != binding.Session.Path {
		t.Fatal(got, err)
	}
	binding.Session.Path = "SDP/Activities/Unregistered/Sessions/session.md"
	evidenceFile(t, service, binding.Session.Path, string(ready))
	raw, _ = json.Marshal(binding)
	evidenceFile(t, service, bindingPath, string(raw))
	if _, err = service.SessionContext(SessionContextRequest{BindingPath: bindingPath}); err == nil {
		t.Fatal("folder name bypassed catalog membership")
	}
}

func TestActivityNoActivePreservesReadOnlyContextRouteAndSourceGuards(t *testing.T) {
	service, domain, q := adoptedGovernance(t)
	run, err := service.Start(LocalController(), q)
	if err != nil {
		t.Fatal(err)
	}
	evidenceFile(t, service, "target.txt", "candidate inspection input")
	catalog, err := domain.Read()
	if err != nil {
		t.Fatal(err)
	}
	a := catalog.Activities[0]
	_, err = domain.Transition(activities.TransitionRequest{OperationID: newTestID(), ExpectedHistoryRevision: catalog.HistoryRevision, ActivityID: a.Descriptor.ID, ExpectedActivityEventID: a.LastEventID, ExpectedDescriptorRevision: a.DescriptorRef.SHA256, Action: "hold", Reason: "Owner selected other discussion", ResumeCondition: "Owner resumes this need"}, activities.Caller{Actor: "owner", Authority: activities.Authority{Role: "project-owner", Source: "authority.md"}}, ActivityExecutionInspector(service.root, service.area, service.project))
	if err != nil {
		t.Fatal(err)
	}
	context, err := service.Context(q.Selection)
	if err != nil || len(context.Conditions) > 0 {
		t.Fatal("no-active blocked read-only context", context, err)
	}
	route, err := service.Route(RouteRequest{Schema: Schema, Categories: []string{"status"}, Selection: q.Selection})
	if err != nil || route.Disposition != "no-work" {
		t.Fatal(route, err)
	}
	if _, err = service.Read(run.Run.ID); err != nil {
		t.Fatal(err)
	}
	// Existing candidate inspection retains its source-pinning rule: any ledger
	// change can stale a run. Activity policy adds no separate read restriction.
	if _, err = service.CaptureCandidate(run.Run.ID, nil); err != nil {
		code(t, err, "stale_revision")
	}
	q.OperationID = newTestID()
	q.RequestID = newTestID()
	q.InputRevision = context.Revision
	if _, err = service.Start(LocalController(), q); err == nil {
		t.Fatal("no-active admitted implementation")
	}
}
