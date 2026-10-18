package governance

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"testing"
)

func preparedFixture(t *testing.T) (*Service, PrepareRequest, PrepareRequest, PreparationResult) {
	t.Helper()
	s, q, other := preparationFixture(t)
	if _, err := s.Bind(LocalController()); err != nil {
		t.Fatal(err)
	}
	q = refreshPreparation(t, s, q)
	r, err := s.PrepareAssignment(LocalController(), q)
	if err != nil {
		t.Fatal(err)
	}
	return s, q, other, r
}

func preparationChange(t *testing.T, s *Service, p PreparedAssignment) PreparationChange {
	t.Helper()
	v, err := s.InspectAssignments(p.ID)
	if err != nil {
		t.Fatal(err)
	}
	op, err := identifier()
	if err != nil {
		t.Fatal(err)
	}
	return PreparationChange{Schema: PreparationLifecycleSchema, ProjectRef: s.project,
		CheckoutID: v.CheckoutID, OperationID: op, ExpectedGeneration: v.Generation,
		AssignmentID: p.ID, ExpectedRevision: p.Revision, Reason: "Owner narrowed the prepared work"}
}

func revisionRequest(t *testing.T, s *Service, p PreparedAssignment) RevisePreparationRequest {
	t.Helper()
	c := preparationChange(t, s, p)
	return RevisePreparationRequest{Schema: c.Schema, ProjectRef: c.ProjectRef, CheckoutID: c.CheckoutID,
		OperationID: c.OperationID, ExpectedGeneration: c.ExpectedGeneration, AssignmentID: c.AssignmentID,
		ExpectedRevision: c.ExpectedRevision, Reason: c.Reason, InputRevision: p.InputRevision,
		Assignment: p.Assignment, Artifacts: append([]ArtifactPin{}, p.Artifacts...), ReturnCriteria: p.ReturnCriteria}
}

func TestPreparationRevisionRetirementAndImmutableHistory(t *testing.T) {
	s, original, other, first := preparedFixture(t)
	// Source drift blocks replacement; refreshing pins makes the same identity usable.
	q := revisionRequest(t, s, first.Assignment)
	path := filepath.Join(s.root, q.Artifacts[0].Path)
	if err := os.WriteFile(path, []byte("new NOW"), 0600); err != nil {
		t.Fatal(err)
	}
	_, err := s.ReviseAssignment(LocalController(), q)
	code(t, err, "stale_revision")
	q.Artifacts[0].SHA256 = hash([]byte("new NOW"))
	q.Assignment.WritablePaths = []string{"work/revised"}
	second, err := s.ReviseAssignment(LocalController(), q)
	if err != nil || second.Assignment.ID != first.Assignment.ID || second.Assignment.Revision != 2 || second.Assignment.State != "prepared" || second.Schema != PreparationLifecycleSchema {
		t.Fatal(second, err)
	}
	if second.Assignment.Role != first.Assignment.Role || second.Assignment.Selection != first.Assignment.Selection {
		t.Fatal("replacement changed identity's governing role/selection")
	}
	// New preparations may claim the released old scope without downgrading 0.3.
	other.Artifacts[0].SHA256 = q.Artifacts[0].SHA256
	other.Assignment.WritablePaths = first.Assignment.Assignment.WritablePaths
	other = refreshPreparation(t, s, other)
	if _, err = s.PrepareAssignment(LocalController(), other); err != nil {
		t.Fatal(err)
	}
	state, _, err := s.load()
	if err != nil || state.Schema != preparationLifecycleStateSchema {
		t.Fatal("preparation downgraded lifecycle store", state.Schema, err)
	}
	// Retirement requires no repaired coverage or artifact bytes.
	retire := RetirePreparationRequest(preparationChange(t, s, second.Assignment))
	for _, path := range []string{q.Artifacts[0].Path, original.Selection.Authority} {
		if err = os.Remove(filepath.Join(s.root, path)); err != nil {
			t.Fatal(err)
		}
	}
	third, err := s.RetireAssignment(LocalController(), retire)
	if err != nil || third.Assignment.Revision != 3 || third.Assignment.State != "retired" || third.Dispatch != "unavailable" || third.DesignAcceptance != "not-assessed" {
		t.Fatal(third, err)
	}
	// Reopen the actual persisted store; receipts never become current-state answers.
	reopened, err := New(s.root, s.area, s.project, s.state)
	if err != nil {
		t.Fatal(err)
	}
	a, e1 := reopened.PrepareAssignment(LocalController(), original)
	b, e2 := reopened.ReviseAssignment(LocalController(), q)
	c, e3 := reopened.RetireAssignment(LocalController(), retire)
	if e1 != nil || e2 != nil || e3 != nil || !reflect.DeepEqual(a, first) || !reflect.DeepEqual(b, second) || !reflect.DeepEqual(c, third) {
		t.Fatal("historical retry changed", a, b, c, e1, e2, e3)
	}
	v, err := reopened.InspectAssignments(third.Assignment.ID)
	if err != nil || v.Schema != PreparationLifecycleSchema || v.Assignments[0].Readiness != "retired" || len(v.Assignments[0].Conditions) != 0 {
		t.Fatal(v, err)
	}
	_, err = reopened.ReviseAssignment(LocalController(), revisionRequest(t, s, third.Assignment))
	code(t, err, "invalid_transition")
	_, err = reopened.RetireAssignment(LocalController(), RetirePreparationRequest(preparationChange(t, s, third.Assignment)))
	code(t, err, "invalid_transition")
}

func TestPreparationReplacementClaimsAreAtomic(t *testing.T) {
	s, _, other, first := preparedFixture(t)
	other = refreshPreparation(t, s, other)
	if _, err := s.PrepareAssignment(LocalController(), other); err != nil {
		t.Fatal(err)
	}
	q := revisionRequest(t, s, first.Assignment)
	q.Assignment.WritablePaths = []string{"work/master/child"}
	_, head, _ := s.load()
	_, err := s.ReviseAssignment(LocalController(), q)
	code(t, err, "conflict")
	state, after, err := s.load()
	if err != nil || head != after || state.Prepared[first.Assignment.ID].Revision != 1 {
		t.Fatal("failed replacement released old claim", err)
	}
	if why := scopeConflict([]string{"work/sad/child"}, state); why == "" {
		t.Fatal("old scope no longer reserved")
	}
	// Self-overlap is allowed; all other prepared and legacy claims still matter.
	q.Assignment.WritablePaths = []string{"work/sad/child"}
	if _, err = s.ReviseAssignment(LocalController(), q); err != nil {
		t.Fatal(err)
	}
}

func TestPreparationChangesRefuseUntrustedStaleAndReusedOperations(t *testing.T) {
	s, _, _, first := preparedFixture(t)
	q := revisionRequest(t, s, first.Assignment)
	_, err := s.ReviseAssignment(Caller{runID: first.Assignment.ID}, q)
	code(t, err, "unauthorized")
	_, err = s.RetireAssignment(Caller{}, RetirePreparationRequest(q.change()))
	code(t, err, "unauthorized")
	bad := q
	bad.CheckoutID, _ = identifier()
	_, err = s.ReviseAssignment(LocalController(), bad)
	code(t, err, "unauthorized")
	bad = q
	bad.ExpectedRevision++
	_, err = s.ReviseAssignment(LocalController(), bad)
	code(t, err, "stale_revision")
	bad = q
	bad.OperationID = first.OperationID
	_, err = s.ReviseAssignment(LocalController(), bad)
	code(t, err, "idempotency_conflict")
	second, err := s.ReviseAssignment(LocalController(), q)
	if err != nil {
		t.Fatal(err)
	}
	_, err = s.RetireAssignment(LocalController(), RetirePreparationRequest(q.change()))
	code(t, err, "idempotency_conflict")
	bad = q
	bad.Reason = "different content"
	_, err = s.ReviseAssignment(LocalController(), bad)
	code(t, err, "idempotency_conflict")
	bad = q
	bad.OperationID, _ = identifier()
	bad.ExpectedRevision = second.Assignment.Revision
	_, err = s.ReviseAssignment(LocalController(), bad)
	code(t, err, "stale_revision")
}

func TestPreparationLifecycleCrashRetryAndConcurrentCAS(t *testing.T) {
	for _, action := range []string{"revise", "retire"} {
		for _, point := range []string{"before-head", "after-head"} {
			t.Run(action+"/"+point, func(t *testing.T) {
				s, _, _, first := preparedFixture(t)
				q := revisionRequest(t, s, first.Assignment)
				q.Assignment.WritablePaths = []string{"new-scope"}
				call := func() (PreparationResult, error) {
					if action == "retire" {
						return s.RetireAssignment(LocalController(), RetirePreparationRequest(q.change()))
					}
					return s.ReviseAssignment(LocalController(), q)
				}
				s.fault = func(where string) error {
					if where == point {
						return errors.New("simulated crash")
					}
					return nil
				}
				if _, err := call(); err == nil {
					t.Fatal("fault did not fire")
				}
				s.fault = nil
				state, _, err := s.load()
				if err != nil {
					t.Fatal(err)
				}
				wantRevision := uint64(1)
				if point == "after-head" {
					wantRevision = 2
				}
				if state.Prepared[first.Assignment.ID].Revision != wantRevision {
					t.Fatal("mutation published at wrong failure boundary")
				}
				got, err := call()
				if err != nil || got.Assignment.Revision != 2 {
					t.Fatal(got, err)
				}
				again, err := call()
				if err != nil || !reflect.DeepEqual(got, again) {
					t.Fatal(again, err)
				}
			})
		}
	}
	s, _, _, first := preparedFixture(t)
	q := revisionRequest(t, s, first.Assignment)
	retire := RetirePreparationRequest(preparationChange(t, s, first.Assignment))
	var wg sync.WaitGroup
	results := make(chan error, 2)
	for _, call := range []func() error{
		func() error { _, err := s.ReviseAssignment(LocalController(), q); return err },
		func() error { _, err := s.RetireAssignment(LocalController(), retire); return err },
	} {
		wg.Add(1)
		go func(call func() error) { defer wg.Done(); results <- call() }(call)
	}
	wg.Wait()
	close(results)
	success := 0
	for err := range results {
		if err == nil {
			success++
		} else {
			var e *Error
			if !errors.As(err, &e) || (e.Code != "busy" && e.Code != "stale_revision") {
				t.Fatal(err)
			}
		}
	}
	if success != 1 {
		t.Fatal("concurrent revise/retire lost update", success)
	}
}

func TestRetiredClaimAllowsLegacyStartAndPreservesProfile(t *testing.T) {
	s, _, _, first := preparedFixture(t)
	sel := Selection{"SDP/KanBan/active/card.md", "SDP/Maintenance/Plan.md", "SDP/Sessions/session.md", "SDP/authority.md", "S1"}
	c, err := s.Context(sel)
	if err != nil {
		t.Fatal(err)
	}
	op, _ := identifier()
	req, _ := identifier()
	q := StartRequest{Schema: Schema, ProjectRef: s.project, CheckoutID: c.CheckoutID, OperationID: op,
		RequestID: req, InputRevision: c.Revision, RoutineDigest: Maintenance().Digest(), Selection: sel,
		Assignment: Assignment{Outcome: "legacy work", WritablePaths: []string{"work/sad"}}}
	_, err = s.Start(LocalController(), q)
	code(t, err, "conflict")
	retire := RetirePreparationRequest(preparationChange(t, s, first.Assignment))
	if _, err = s.RetireAssignment(LocalController(), retire); err != nil {
		t.Fatal(err)
	}
	if _, err = s.Start(LocalController(), q); err != nil {
		t.Fatal("retired claim blocks legacy work", err)
	}
	state, _, err := s.load()
	if err != nil || state.Schema != preparationLifecycleStateSchema || len(state.Prepared) != 1 || len(state.PreparationReceipts) != 2 {
		t.Fatal("legacy start dropped lifecycle history/profile", state.Schema, err)
	}
	_, _, err = retry(state, retire.OperationID, "anything")
	code(t, err, "idempotency_conflict")
}

func TestPreparationLifecycleWireAndRefusalPreserveState(t *testing.T) {
	s, _, _, first := preparedFixture(t)
	q := revisionRequest(t, s, first.Assignment)
	b, err := json.Marshal(q)
	if err != nil {
		t.Fatal(err)
	}
	var decoded RevisePreparationRequest
	if err = Decode(b, &decoded); err != nil || !reflect.DeepEqual(q, decoded) {
		t.Fatal("revision wire request does not decode", err)
	}
	retire := RetirePreparationRequest(q.change())
	b, err = json.Marshal(retire)
	var decodedRetire RetirePreparationRequest
	if err != nil {
		t.Fatal(err)
	}
	if err = Decode(b, &decodedRetire); err != nil || !reflect.DeepEqual(retire, decodedRetire) {
		t.Fatal("retirement wire request does not decode", err)
	}
	b = append(b[:len(b)-1], []byte(`,"role":"Master"}`)...)
	code(t, Decode(b, &decodedRetire), "invalid_request")
	_, head, _ := s.load()
	for _, edit := range []func(*RevisePreparationRequest){
		func(v *RevisePreparationRequest) { v.Schema = PreparationSchema },
		func(v *RevisePreparationRequest) { v.Reason = " " },
		func(v *RevisePreparationRequest) { v.ExpectedRevision = 0 },
		func(v *RevisePreparationRequest) { v.Assignment.WritablePaths = []string{"../escape"} },
		func(v *RevisePreparationRequest) { v.Assignment.WritablePaths = []string{"work/*"} },
		func(v *RevisePreparationRequest) { v.ReturnCriteria = nil },
	} {
		bad := q
		edit(&bad)
		if _, err = s.ReviseAssignment(LocalController(), bad); err == nil {
			t.Fatal("invalid replacement accepted", bad)
		}
		_, after, err := s.load()
		if err != nil || head != after {
			t.Fatal("refusal changed state", err)
		}
	}
	// Do not let revision bypass a currently executing legacy claim.
	sel := Selection{"SDP/KanBan/active/card.md", "SDP/Maintenance/Plan.md", "SDP/Sessions/session.md", "SDP/authority.md", "S1"}
	c, err := s.Context(sel)
	if err != nil {
		t.Fatal(err)
	}
	op, _ := identifier()
	req, _ := identifier()
	start := StartRequest{Schema: Schema, ProjectRef: s.project, CheckoutID: c.CheckoutID, OperationID: op,
		RequestID: req, InputRevision: c.Revision, RoutineDigest: Maintenance().Digest(), Selection: sel,
		Assignment: Assignment{Outcome: "legacy work", WritablePaths: []string{"legacy-only"}}}
	if _, err = s.Start(LocalController(), start); err != nil {
		t.Fatal(err)
	}
	q = revisionRequest(t, s, first.Assignment)
	q.Assignment.WritablePaths = []string{"legacy-only/nested"}
	_, err = s.ReviseAssignment(LocalController(), q)
	code(t, err, "conflict")
	q.OperationID = op
	_, err = s.ReviseAssignment(LocalController(), q)
	code(t, err, "idempotency_conflict")
}

func TestPreparationLifecycleWorkerRefusalCreatesNothing(t *testing.T) {
	s, _, _ := preparationFixture(t)
	if _, err := s.ReviseAssignment(Caller{}, RevisePreparationRequest{}); err == nil {
		t.Fatal("worker revised assignment")
	}
	if _, err := s.RetireAssignment(Caller{}, RetirePreparationRequest{}); err == nil {
		t.Fatal("worker retired assignment")
	}
	if _, err := os.Stat(s.state); !os.IsNotExist(err) {
		t.Fatal("worker call created state", err)
	}
}
