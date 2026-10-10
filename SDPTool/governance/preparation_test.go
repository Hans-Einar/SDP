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
	"sync"
	"testing"
)

func preparationFixture(t *testing.T) (*Service, PrepareRequest, PrepareRequest) {
	t.Helper()
	s, _, legacy := fixture(t)
	ledger, err := os.ReadFile(filepath.Join(s.root, "SDP/ProjectManagement/Ledger.ndjson"))
	if err != nil {
		t.Fatal(err)
	}
	qs := []PrepareRequest{}
	for i, role := range []string{"SAD", "Master"} {
		name := strings.ToLower(role)
		planType := "DesignPlan"
		if role == "Master" {
			planType = "ImplementationPlan"
		}
		card := fmt.Sprintf("KB-TEST-%03d", i+2)
		plan := fmt.Sprintf("PLAN-TEST-%04d", i+1)
		sel := legacy.Selection
		sel.Card = "SDP/KanBan/active/" + name + ".md"
		sel.Plan = "SDP/Maintenance/" + name + ".md"
		files := map[string]string{sel.Card: fmt.Sprintf("# Work\n\n| Field | Value |\n| --- | --- |\n| id | %s |\n| PlanId | %s |\n| CardState | in-progress |\n", card, plan), sel.Plan: fmt.Sprintf("# Plan\n\n| Field | Value |\n| --- | --- |\n| id | %s |\n| state | active |\n| PlanType | %s |\n", plan, planType)}
		for p, b := range files {
			if err = os.WriteFile(filepath.Join(s.root, p), []byte(b), 0600); err != nil {
				t.Fatal(err)
			}
		}
		ledger = append(ledger, []byte(fmt.Sprintf(`{"schemaVersion":"1.0","eventId":"P%d","eventType":"x-kanban:created","subjectId":"%s","payload":{"schemaVersion":"0.2","previousEventId":null,"to":"active","toPath":"active/%s.md"}}
{"schemaVersion":"1.0","eventId":"Q%d","eventType":"x-management:created","subjectId":"%s","payload":{"schemaVersion":"0.2","previousEventId":null,"to":"active","toPath":"Maintenance/%s.md"}}
`, i, card, name, i, plan, name))...)
		kinds := []string{"now"}
		if role == "Master" {
			kinds = append(kinds, "target", "task", "blueprint", "design-disposition")
		}
		pins := []ArtifactPin{}
		for _, kind := range kinds {
			p := "fixture-" + kind + ".json"
			b := []byte(`{"fixture":"` + kind + `"}`)
			if err = os.WriteFile(filepath.Join(s.root, p), b, 0600); err != nil {
				t.Fatal(err)
			}
			pins = append(pins, ArtifactPin{kind, p, hash(b)})
		}
		op, _ := identifier()
		qs = append(qs, PrepareRequest{Schema: PreparationSchema, ProjectRef: s.project, OperationID: op, Role: role, Selection: sel, Assignment: Assignment{Outcome: "bounded " + role, WritablePaths: []string{"work/" + name}, ExcludedPaths: []string{"unrelated.txt"}}, Artifacts: pins, ReturnCriteria: []string{"Return candidate-bound evidence and unresolved findings"}})
	}
	if err = os.WriteFile(filepath.Join(s.root, "SDP/ProjectManagement/Ledger.ndjson"), ledger, 0600); err != nil {
		t.Fatal(err)
	}
	return s, qs[0], qs[1]
}
func refreshPreparation(t *testing.T, s *Service, q PrepareRequest) PrepareRequest {
	t.Helper()
	c, err := s.PreparationContext(PreparationContextRequest{q.Role, q.Selection})
	if err != nil || len(c.Context.Conditions) > 0 {
		t.Fatal(c, err)
	}
	q.CheckoutID = c.Context.CheckoutID
	q.ExpectedGeneration = c.Generation
	q.InputRevision = c.Context.Revision
	return q
}
func TestPreparationAuthorityPersistenceAndStaleTarget(t *testing.T) {
	s, sad, master := preparationFixture(t)
	c, err := s.PreparationContext(PreparationContextRequest{sad.Role, sad.Selection})
	if err != nil || c.Generation != "" || len(c.Context.Conditions) > 0 {
		t.Fatal(c, err)
	}
	_, err = s.PrepareAssignment(Caller{}, sad)
	code(t, err, "unauthorized")
	if _, err = os.Stat(s.state); !os.IsNotExist(err) {
		t.Fatal("untrusted call created state")
	}
	if _, err = s.Bind(LocalController()); err != nil {
		t.Fatal(err)
	}
	sad = refreshPreparation(t, s, sad)
	one, err := s.PrepareAssignment(LocalController(), sad)
	if err != nil {
		t.Fatal(err)
	}
	master = refreshPreparation(t, s, master)
	two, err := s.PrepareAssignment(LocalController(), master)
	if err != nil {
		t.Fatal(err)
	}
	if one.Assignment.ID == two.Assignment.ID || two.Dispatch != "unavailable" || two.DesignAcceptance != "not-assessed" {
		t.Fatal(one, two)
	}
	reopened, err := New(s.root, s.area, s.project, s.state)
	if err != nil {
		t.Fatal(err)
	}
	v, err := reopened.InspectAssignments("")
	if err != nil || len(v.Assignments) != 2 {
		t.Fatal(v, err)
	}
	before, _ := os.ReadFile(filepath.Join(s.state, "HEAD"))
	if err = os.WriteFile(filepath.Join(s.root, "fixture-target.json"), []byte("changed TARGET"), 0600); err != nil {
		t.Fatal(err)
	}
	v, err = reopened.InspectAssignments(two.Assignment.ID)
	if err != nil || v.Assignments[0].Readiness != "stale" || !strings.Contains(v.Assignments[0].Conditions[0], "target") {
		t.Fatal(v, err)
	}
	again, err := reopened.PrepareAssignment(LocalController(), master)
	if err != nil || !reflect.DeepEqual(again, two) {
		t.Fatal("retry lost original identity", again, err)
	}
	after, _ := os.ReadFile(filepath.Join(s.state, "HEAD"))
	if !bytes.Equal(before, after) {
		t.Fatal("inspect/retry changed state")
	}
	master.Assignment.Outcome = "changed meaning"
	_, err = reopened.PrepareAssignment(LocalController(), master)
	code(t, err, "idempotency_conflict")
	// Preparation does not widen the existing Maintenance route.
	old, err := s.Context(sad.Selection)
	if err != nil || len(old.Conditions) == 0 {
		t.Fatal("legacy route admitted DesignPlan", old, err)
	}
}
func TestPreparationOverlapAndStaleGeneration(t *testing.T) {
	s, sad, master := preparationFixture(t)
	if _, err := s.Bind(LocalController()); err != nil {
		t.Fatal(err)
	}
	sad = refreshPreparation(t, s, sad)
	master = refreshPreparation(t, s, master)
	if _, err := s.PrepareAssignment(LocalController(), sad); err != nil {
		t.Fatal(err)
	}
	_, err := s.PrepareAssignment(LocalController(), master)
	code(t, err, "stale_revision")
	for _, p := range []string{"work", "work/sad", "work/sad/nested"} {
		q := refreshPreparation(t, s, master)
		q.Assignment.WritablePaths = []string{p}
		_, err = s.PrepareAssignment(LocalController(), q)
		code(t, err, "conflict")
	}
	q := refreshPreparation(t, s, master)
	q.Assignment.WritablePaths = []string{"work/sad-other"}
	if _, err = s.PrepareAssignment(LocalController(), q); err != nil {
		t.Fatal("sibling scope falsely overlaps", err)
	}
}
func TestPreparationSaveFailureRetryAndCAS(t *testing.T) {
	for _, point := range []string{"before-head", "after-head"} {
		t.Run(point, func(t *testing.T) {
			s, q, _ := preparationFixture(t)
			s.Bind(LocalController())
			q = refreshPreparation(t, s, q)
			s.fault = func(where string) error {
				if where == point {
					return errors.New("simulated crash")
				}
				return nil
			}
			if _, err := s.PrepareAssignment(LocalController(), q); err == nil {
				t.Fatal("fault did not fire")
			}
			s.fault = nil
			got, err := s.PrepareAssignment(LocalController(), q)
			if err != nil {
				t.Fatal(err)
			}
			again, err := s.PrepareAssignment(LocalController(), q)
			if err != nil || again.Assignment.ID != got.Assignment.ID {
				t.Fatal(again, err)
			}
			v, err := s.InspectAssignments("")
			if err != nil || len(v.Assignments) != 1 {
				t.Fatal(v, err)
			}
		})
	}
	s, a, b := preparationFixture(t)
	s.Bind(LocalController())
	a = refreshPreparation(t, s, a)
	b = refreshPreparation(t, s, b)
	var wg sync.WaitGroup
	results := make(chan error, 2)
	for _, q := range []PrepareRequest{a, b} {
		wg.Add(1)
		go func(q PrepareRequest) {
			defer wg.Done()
			_, err := s.PrepareAssignment(LocalController(), q)
			results <- err
		}(q)
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
		t.Fatal("generation CAS admitted concurrent stale operation", success)
	}
}
func TestPreparationPathsAuthorityAndSharedReceiptNamespace(t *testing.T) {
	s, q, _ := preparationFixture(t)
	s.Bind(LocalController())
	q = refreshPreparation(t, s, q)
	for _, edit := range []func(*PrepareRequest){func(v *PrepareRequest) { v.Role = "ProjectLeader" }, func(v *PrepareRequest) { v.ReturnCriteria = nil }, func(v *PrepareRequest) { v.Artifacts[0].Path = "../outside" }, func(v *PrepareRequest) { v.Assignment.WritablePaths = []string{"work/*"} }} {
		b, _ := json.Marshal(q)
		var v PrepareRequest
		json.Unmarshal(b, &v)
		edit(&v)
		if _, err := s.PrepareAssignment(LocalController(), v); err == nil {
			t.Fatal("invalid request accepted", v)
		}
	}
	original := q.Artifacts[0].Path
	if err := os.Symlink(filepath.Join(s.root, original), filepath.Join(s.root, "alias.json")); err != nil {
		t.Fatal(err)
	}
	q.Artifacts[0].Path = "alias.json"
	if _, err := s.PrepareAssignment(LocalController(), q); err == nil {
		t.Fatal("symlink admitted")
	}
	q.Artifacts[0].Path = original
	q.CheckoutID, _ = identifier()
	_, err := s.PrepareAssignment(LocalController(), q)
	code(t, err, "unauthorized")
	q = refreshPreparation(t, s, q)
	one, err := s.PrepareAssignment(LocalController(), q)
	if err != nil {
		t.Fatal(err)
	}
	state, _, err := s.load()
	if err != nil {
		t.Fatal(err)
	}
	_, ok, err := retry(state, one.OperationID, "anything")
	if !ok {
		t.Fatal("legacy retry ignored preparation receipt")
	}
	code(t, err, "idempotency_conflict")
	state.Receipts = map[string]receipt{}
	// Exercise legacy-to-preparation collision through an actual valid UUID receipt.
	op, _ := identifier()
	state.Receipts[op] = receipt{}
	_, head, _ := s.load()
	if err = s.save(state, head); err != nil {
		t.Fatal(err)
	}
	q = refreshPreparation(t, s, q)
	q.OperationID = op
	_, err = s.PrepareAssignment(LocalController(), q)
	code(t, err, "idempotency_conflict")
}

func TestPreparationAndLegacyRunScopeCoexist(t *testing.T) {
	s, sad, _ := preparationFixture(t)
	binding, err := s.Bind(LocalController())
	if err != nil {
		t.Fatal(err)
	}
	sel := Selection{"SDP/KanBan/active/card.md", "SDP/Maintenance/Plan.md", "SDP/Sessions/session.md", "SDP/authority.md", "S1"}
	c, err := s.Context(sel)
	if err != nil {
		t.Fatal(err)
	}
	op, _ := identifier()
	req, _ := identifier()
	q := StartRequest{Schema: Schema, ProjectRef: s.project, CheckoutID: binding.CheckoutID, OperationID: op, RequestID: req, InputRevision: c.Revision, RoutineDigest: Maintenance().Digest(), Selection: sel, Assignment: Assignment{Outcome: "legacy scoped work", WritablePaths: []string{"work/sad"}, ExcludedPaths: []string{"unrelated.txt"}}}
	run, err := s.Start(LocalController(), q)
	if err != nil {
		t.Fatal(err)
	}
	sad = refreshPreparation(t, s, sad)
	_, err = s.PrepareAssignment(LocalController(), sad)
	code(t, err, "conflict")
	if _, err = s.Transition(LocalController(), transition(run, "cancel")); err != nil {
		t.Fatal(err)
	}
	sad = refreshPreparation(t, s, sad)
	if _, err = s.PrepareAssignment(LocalController(), sad); err != nil {
		t.Fatal(err)
	}
	q.OperationID, _ = identifier()
	q.RequestID, _ = identifier()
	_, err = s.Start(LocalController(), q)
	code(t, err, "conflict")
	q.Assignment.WritablePaths = []string{"work/*"}
	_, err = s.Start(LocalController(), q)
	code(t, err, "conflict")
	q.Assignment.WritablePaths = []string{"work/other"}
	if _, err = s.Start(LocalController(), q); err != nil {
		t.Fatal("disjoint legacy run blocked", err)
	}
	// Both old and prepared records survive another legacy write in profile 0.2.
	v, err := s.InspectAssignments("")
	if err != nil || len(v.Assignments) != 1 {
		t.Fatal(v, err)
	}
	if err = os.WriteFile(filepath.Join(s.root, sel.Authority), []byte("new instruction"), 0600); err != nil {
		t.Fatal(err)
	}
	v, err = s.InspectAssignments("")
	if err != nil || v.Assignments[0].Readiness != "stale" {
		t.Fatal(v, err)
	}
}
