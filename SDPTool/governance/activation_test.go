package governance

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func activationFixture(t *testing.T) (designFixture, ActivationRequest) {
	t.Helper()
	f, q, input := decisionFixture(t)
	d, err := f.s.RecordDesignDecision(LocalController(), q, input)
	if err != nil {
		t.Fatal(err)
	}
	operation, _ := identifier()
	request, _ := identifier()
	return f, ActivationRequest{Schema: ActivationSchema, ProjectRef: q.ProjectRef, CheckoutID: q.CheckoutID, OperationID: operation, RequestID: request, Review: freshDesignQuery(t, f), ReviewRevision: q.ReviewRevision, DecisionOperationID: d.OperationID}
}
func runChange(t *testing.T, s *Service, r Result, action string) (Result, error) {
	t.Helper()
	id, _ := identifier()
	return s.Transition(LocalController(), TransitionRequest{Schema: Schema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, OperationID: id, RunID: r.Run.ID, ExpectedRevision: r.Run.Revision, InputRevision: r.Run.InputRevision, Action: action, Reason: "Fixture " + action})
}
func TestActivationUsesImplementationContextAndPreservesHistory(t *testing.T) {
	f, q := activationFixture(t)
	before := inventoryDesign(t, filepath.Dir(f.s.root))
	_, err := f.s.ActivateAssignment(Caller{}, q)
	code(t, err, "unauthorized")
	if !reflect.DeepEqual(before, inventoryDesign(t, filepath.Dir(f.s.root))) {
		t.Fatal("worker mutation")
	}
	result, err := f.s.ActivateAssignment(LocalController(), q)
	if err != nil {
		t.Fatal(err)
	}
	if result.Run.State != "ready" || result.Run.Routine.ID != "implementation" {
		t.Fatal(result)
	}
	state, _, err := f.s.load()
	if err != nil || state.Schema != activationStateSchema || len(state.Runs) != 1 {
		t.Fatal(err)
	}
	again, err := f.s.ActivateAssignment(LocalController(), q)
	if err != nil || !reflect.DeepEqual(result, again) {
		t.Fatal(again, err)
	}
	other := q
	other.OperationID, _ = identifier()
	_, err = f.s.ActivateAssignment(LocalController(), other)
	code(t, err, "conflict")
	ctx, err := f.s.RunContext(result.Run.ID)
	if err != nil || ctx.PlanID != "PLAN-TEST-0002" || len(ctx.Conditions) > 0 {
		t.Fatal(ctx, err)
	}
	decision, err := f.s.RunDesign(result.Run.ID)
	if err != nil || decision.OperationID != q.DecisionOperationID {
		t.Fatal(decision, err)
	}
	if err = f.s.CheckRunReady(result.Run.ID); err != nil {
		t.Fatal(err)
	}
	revision := revisionRequest(t, f.s, state.Prepared[q.Review.AssignmentID])
	_, err = f.s.ReviseAssignment(LocalController(), revision)
	code(t, err, "invalid_transition")
	retire := RetirePreparationRequest(revision.change())
	_, err = f.s.RetireAssignment(LocalController(), retire)
	code(t, err, "invalid_transition")
	result, err = runChange(t, f.s, result, "start")
	if err != nil {
		t.Fatal(err)
	}
	// Live model drift must reach run/MCP context, not only the design-check command.
	path := filepath.Join(f.work, "System.design")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(path, append(raw, []byte("\n// changed target during work\n")...), 0600); err != nil {
		t.Fatal(err)
	}
	read, err := f.s.Read(result.Run.ID)
	if err != nil || read.Coverage != "stale" {
		t.Fatal(read, err)
	}
	if err = f.s.CheckRunReady(result.Run.ID); err == nil {
		t.Fatal("stale run ready")
	}
	_, err = runChange(t, f.s, read, "report")
	code(t, err, "stale_revision")
	result, err = runChange(t, f.s, read, "interrupt")
	if err != nil {
		t.Fatal(err)
	}
	_, err = runChange(t, f.s, result, "resume")
	code(t, err, "stale_revision")
	if err = os.WriteFile(path, raw, 0600); err != nil {
		t.Fatal(err)
	}
	result, err = runChange(t, f.s, result, "resume")
	if err != nil {
		t.Fatal(err)
	}
	// A superseding decision cannot silently alter the existing execution authority.
	qDecision := DesignDecisionRequest{Schema: DesignDecisionSchema, ProjectRef: q.ProjectRef, CheckoutID: q.CheckoutID, Review: freshDesignQuery(t, f), ReviewRevision: q.ReviewRevision}
	qDecision.OperationID, _ = identifier()
	_, err = f.s.RecordDesignDecision(LocalController(), qDecision, OwnerDesignInput{Actor: "fixture-owner", Source: "test", Decision: "rejected", Comment: "Synthetic revocation"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = runChange(t, f.s, result, "start")
	code(t, err, "stale_revision")
	if _, err = runChange(t, f.s, result, "cancel"); err != nil {
		t.Fatal(err)
	}
	state, _, err = f.s.load()
	if err != nil || state.Schema != activationStateSchema || len(state.AssignmentRuns) != 1 {
		t.Fatal(err)
	}
	// Historical retry never re-approves or creates another run.
	again, err = f.s.ActivateAssignment(LocalController(), q)
	if err != nil || again.Run.State != "ready" {
		t.Fatal(again, err)
	}
	current, err := f.s.Read(again.Run.ID)
	if err != nil || current.Run.State != "canceled" {
		t.Fatal(current, err)
	}
}
func TestActivationRefusesChangedApprovalAndRecoversWriteFaults(t *testing.T) {
	for _, kind := range []string{"decision", "review", "generation", "source", "scope", "before-head", "after-head"} {
		t.Run(kind, func(t *testing.T) {
			f, q := activationFixture(t)
			_, head, _ := f.s.load()
			switch kind {
			case "decision":
				q.DecisionOperationID, _ = identifier()
			case "review":
				q.ReviewRevision = strings.Repeat("0", 64)
			case "generation":
				q.Review.ExpectedGeneration = strings.Repeat("0", 64)
			case "source":
				path := filepath.Join(f.work, "System.design")
				b, _ := os.ReadFile(path)
				if err := os.WriteFile(path, append(b, []byte("\n// drift\n")...), 0600); err != nil {
					t.Fatal(err)
				}
			case "scope":
				state, h, err := f.s.load()
				if err != nil {
					t.Fatal(err)
				}
				id, _ := identifier()
				state.Runs[id] = Run{ID: id, State: "running", Routine: Maintenance(), Assignment: Assignment{WritablePaths: []string{f.prepare.Assignment.WritablePaths[0]}}}
				if err = f.s.save(state, h); err != nil {
					t.Fatal(err)
				}
				q.Review = freshDesignQuery(t, f)
				_, head, _ = f.s.load()
			default:
				f.s.fault = func(at string) error {
					if at == kind {
						return errors.New("uncertain activation write")
					}
					return nil
				}
			}
			if _, err := f.s.ActivateAssignment(LocalController(), q); err == nil {
				t.Fatal("invalid/fault activation accepted")
			}
			f.s.fault = nil
			state, after, err := f.s.load()
			if err != nil {
				t.Fatal(err)
			}
			if kind == "after-head" {
				if len(state.AssignmentRuns) != 1 {
					t.Fatal("committed activation lost")
				}
			} else if after != head || len(state.AssignmentRuns) > 0 {
				t.Fatal("partial activation")
			}
			if kind == "before-head" || kind == "after-head" {
				if _, err = f.s.ActivateAssignment(LocalController(), q); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}
