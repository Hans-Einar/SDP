package governance

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
)

func decisionFixture(t *testing.T) (designFixture, DesignDecisionRequest, OwnerDesignInput) {
	t.Helper()
	f := providerDesignFixture(t)
	review, err := f.s.CheckAssignmentDesign(f.q)
	if err != nil || review.Status != "current" {
		t.Fatal(review, err)
	}
	state, _, err := f.s.load()
	if err != nil {
		t.Fatal(err)
	}
	id, err := identifier()
	if err != nil {
		t.Fatal(err)
	}
	return f, DesignDecisionRequest{Schema: DesignDecisionSchema, ProjectRef: f.s.project, CheckoutID: state.CheckoutID, OperationID: id, Review: f.q, ReviewRevision: review.ReviewRevision}, OwnerDesignInput{Actor: "fixture-owner", Source: "scripted-test-input", Decision: "approved", Comment: "Synthetic decision, not actual owner acceptance"}
}
func freshDesignQuery(t *testing.T, f designFixture) DesignCheckRequest {
	t.Helper()
	state, head, err := f.s.load()
	if err != nil {
		t.Fatal(err)
	}
	q := f.q
	q.ExpectedGeneration = head
	q.ExpectedRevision = state.Prepared[q.AssignmentID].Revision
	return q
}
func TestDesignDecisionPersistenceApplicabilityAndSupersession(t *testing.T) {
	f, q, input := decisionFixture(t)
	first, err := f.s.RecordDesignDecision(LocalController(), q, input)
	if err != nil {
		t.Fatal(err)
	}
	if first.Dispatch != "unavailable" || first.Record.Input != input {
		t.Fatal(first)
	}
	// Re-open the actual store and recover an immutable result, not fresh acceptance.
	reopened, err := New(f.s.root, f.s.area, f.s.project, f.s.state)
	if err != nil {
		t.Fatal(err)
	}
	f.s = reopened
	retry, err := f.s.RecordDesignDecision(LocalController(), q, input)
	if err != nil || !reflect.DeepEqual(first, retry) {
		t.Fatal(retry, err)
	}
	before := inventoryDesign(t, filepath.Dir(f.s.root))
	review, err := f.s.ReviewAssignmentDesign(freshDesignQuery(t, f))
	if err != nil || review.Acceptance != "approved" || review.Dispatch != "unavailable" {
		t.Fatal(review, err)
	}
	if !reflect.DeepEqual(before, inventoryDesign(t, filepath.Dir(f.s.root))) {
		t.Fatal("inspection wrote files")
	}
	target := filepath.Join(f.work, "System.design")
	original, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(target, append(original, []byte("\n// changed TARGET\n")...), 0600); err != nil {
		t.Fatal(err)
	}
	review, err = f.s.ReviewAssignmentDesign(freshDesignQuery(t, f))
	if err != nil || review.Acceptance != "stale" || review.Decision.Record.Input.Decision != "approved" {
		t.Fatal(review, err)
	}
	retry, err = f.s.RecordDesignDecision(LocalController(), q, input)
	if err != nil || !reflect.DeepEqual(first, retry) {
		t.Fatal(retry, err)
	}
	if err = os.WriteFile(target, original, 0600); err != nil {
		t.Fatal(err)
	}
	for _, decision := range []string{"changes-requested", "rejected"} {
		q.Review = freshDesignQuery(t, f)
		q.OperationID, _ = identifier()
		input.Decision = decision
		if _, err = f.s.RecordDesignDecision(LocalController(), q, input); err != nil {
			t.Fatal(err)
		}
		review, err = f.s.ReviewAssignmentDesign(freshDesignQuery(t, f))
		if err != nil || review.Acceptance != decision {
			t.Fatal(review, err)
		}
	}
	state, _, err := f.s.load()
	if err != nil || state.Schema != designDecisionStateSchema || len(state.DesignDecisionReceipts) != 3 {
		t.Fatal(state.Schema, err)
	}
	// Revision invalidates an equal-source basis and must not downgrade profile 0.4.
	change := revisionRequest(t, f.s, state.Prepared[f.q.AssignmentID])
	if _, err = f.s.ReviseAssignment(LocalController(), change); err != nil {
		t.Fatal(err)
	}
	review, err = f.s.ReviewAssignmentDesign(freshDesignQuery(t, f))
	if err != nil || review.Acceptance != "stale" {
		t.Fatal(review, err)
	}
	state, head, err := f.s.load()
	if err != nil || state.Schema != designDecisionStateSchema {
		t.Fatal(state.Schema, err)
	}
	retire := RetirePreparationRequest{Schema: PreparationLifecycleSchema, ProjectRef: f.s.project, CheckoutID: state.CheckoutID, ExpectedGeneration: head, AssignmentID: f.q.AssignmentID, ExpectedRevision: state.Prepared[f.q.AssignmentID].Revision, Reason: "Retire fixture"}
	retire.OperationID, _ = identifier()
	if _, err = f.s.RetireAssignment(LocalController(), retire); err != nil {
		t.Fatal(err)
	}
	review, err = f.s.ReviewAssignmentDesign(freshDesignQuery(t, f))
	if err != nil || review.Acceptance != "stale" || review.Check.Status != "blocked" {
		t.Fatal(review, err)
	}
	// A later prepare also preserves the decision profile and retained history.
	prepare := refreshPreparation(t, f.s, f.prepare)
	prepare.OperationID, _ = identifier()
	if _, err = f.s.PrepareAssignment(LocalController(), prepare); err != nil {
		t.Fatal(err)
	}
	state, _, err = f.s.load()
	if err != nil || state.Schema != designDecisionStateSchema || len(state.DesignDecisionReceipts) != 3 {
		t.Fatal(state.Schema, err)
	}
}

func TestDesignDecisionRefusalsAndOperationNamespace(t *testing.T) {
	f, q, input := decisionFixture(t)
	before := inventoryDesign(t, filepath.Dir(f.s.root))
	for _, caller := range []Caller{{}, {runID: "worker"}} {
		_, err := f.s.RecordDesignDecision(caller, q, input)
		code(t, err, "unauthorized")
		_, err = f.s.DesignDecisionReceipt(caller, q)
		code(t, err, "unauthorized")
	}
	for _, mutate := range []func(*DesignDecisionRequest, *OwnerDesignInput){
		func(q *DesignDecisionRequest, i *OwnerDesignInput) { q.CheckoutID = "wrong" },
		func(q *DesignDecisionRequest, i *OwnerDesignInput) { q.ReviewRevision = strings.Repeat("0", 64) },
		func(q *DesignDecisionRequest, i *OwnerDesignInput) {
			q.Review.ExpectedGeneration = strings.Repeat("0", 64)
		},
		func(q *DesignDecisionRequest, i *OwnerDesignInput) { i.Decision = "" },
		func(q *DesignDecisionRequest, i *OwnerDesignInput) { i.Actor = "" },
		func(q *DesignDecisionRequest, i *OwnerDesignInput) { i.Source = "" },
		func(q *DesignDecisionRequest, i *OwnerDesignInput) { i.Comment = "" },
	} {
		bad, owner := q, input
		mutate(&bad, &owner)
		if _, err := f.s.RecordDesignDecision(LocalController(), bad, owner); err == nil {
			t.Fatal("bad decision accepted")
		}
	}
	if !reflect.DeepEqual(before, inventoryDesign(t, filepath.Dir(f.s.root))) {
		t.Fatal("refusal wrote project/state")
	}
	state, _, err := f.s.load()
	if err != nil {
		t.Fatal(err)
	}
	for id := range state.PreparationReceipts {
		bad := q
		bad.OperationID = id
		_, err = f.s.RecordDesignDecision(LocalController(), bad, input)
		code(t, err, "idempotency_conflict")
		break
	}
	if _, err = f.s.RecordDesignDecision(LocalController(), q, input); err != nil {
		t.Fatal(err)
	}
	input.Comment = "different"
	_, err = f.s.RecordDesignDecision(LocalController(), q, input)
	code(t, err, "idempotency_conflict")
	bad := q
	bad.ReviewRevision = strings.Repeat("0", 64)
	_, err = f.s.DesignDecisionReceipt(LocalController(), bad)
	code(t, err, "idempotency_conflict")
	state, _, err = f.s.load()
	if err != nil {
		t.Fatal(err)
	}
	_, ok, err := retry(state, q.OperationID, "anything")
	if !ok {
		t.Fatal("legacy retry missed decision collision")
	}
	code(t, err, "idempotency_conflict")
	prepare := refreshPreparation(t, f.s, f.prepare)
	prepare.OperationID = q.OperationID
	_, err = f.s.PrepareAssignment(LocalController(), prepare)
	code(t, err, "idempotency_conflict")
	revise := revisionRequest(t, f.s, state.Prepared[q.Review.AssignmentID])
	revise.OperationID = q.OperationID
	_, err = f.s.ReviseAssignment(LocalController(), revise)
	code(t, err, "idempotency_conflict")
}

func TestDesignDecisionFaultRecoveryDriftAndConcurrentCAS(t *testing.T) {
	for _, point := range []string{"before-head", "after-head"} {
		t.Run(point, func(t *testing.T) {
			f, q, input := decisionFixture(t)
			f.s.fault = func(at string) error {
				if at == point {
					return errors.New("simulated write uncertainty")
				}
				return nil
			}
			if _, err := f.s.RecordDesignDecision(LocalController(), q, input); err == nil {
				t.Fatal("fault missing")
			}
			f.s.fault = nil
			receipt, err := f.s.DesignDecisionReceipt(LocalController(), q)
			if err != nil {
				t.Fatal(err)
			}
			if (receipt != nil) != (point == "after-head") {
				t.Fatal("incorrect commit boundary", receipt)
			}
			if _, err = f.s.RecordDesignDecision(LocalController(), q, input); err != nil {
				t.Fatal(err)
			}
			state, _, err := f.s.load()
			if err != nil || len(state.DesignDecisionReceipts) != 1 {
				t.Fatal(err)
			}
		})
	}
	t.Run("source-drift-before-publication", func(t *testing.T) {
		f, q, input := decisionFixture(t)
		_, head, _ := f.s.load()
		_, err := f.s.recordDesignDecision(LocalController(), q, input, func() {
			path := filepath.Join(f.work, "System.design")
			b, e := os.ReadFile(path)
			if e != nil {
				t.Fatal(e)
			}
			if e = os.WriteFile(path, append(b, []byte("\n// drift\n")...), 0600); e != nil {
				t.Fatal(e)
			}
		})
		code(t, err, "stale_revision")
		state, after, err := f.s.load()
		if err != nil || after != head || len(state.DesignDecisionReceipts) != 0 {
			t.Fatal("partial write", err)
		}
	})
	t.Run("concurrent-decisions", func(t *testing.T) {
		f, q, input := decisionFixture(t)
		other := q
		other.OperationID, _ = identifier()
		var wg sync.WaitGroup
		errs := make(chan error, 2)
		for _, request := range []DesignDecisionRequest{q, other} {
			wg.Add(1)
			go func(q DesignDecisionRequest) {
				defer wg.Done()
				_, e := f.s.RecordDesignDecision(LocalController(), q, input)
				errs <- e
			}(request)
		}
		wg.Wait()
		close(errs)
		passed := 0
		for err := range errs {
			if err == nil {
				passed++
			} else if !strings.HasPrefix(err.Error(), "busy:") {
				code(t, err, "stale_revision")
			}
		}
		if passed != 1 {
			t.Fatal(passed)
		}
	})
}
