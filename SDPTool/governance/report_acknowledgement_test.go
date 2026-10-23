package governance

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func reportFixture(t *testing.T) (*Service, Selection, Result, ReportAcknowledgementRequest) {
	t.Helper()
	s, sel, r := attentionFixture(t)
	r, e := s.Transition(LocalController(), transition(r, "report"))
	if e != nil {
		t.Fatal(e)
	}
	return s, sel, r, acknowledgementRequest(t, s, r)
}
func acknowledgementRequest(t *testing.T, s *Service, r Result) ReportAcknowledgementRequest {
	t.Helper()
	view, e := s.ReportInbox(r.Run.ID)
	if e != nil || len(view.Reports) != 1 {
		t.Fatal(view, e)
	}
	v := view.Reports[0]
	op, _ := identifier()
	return ReportAcknowledgementRequest{Schema: ReportAcknowledgementSchema, ProjectRef: view.ProjectRef, CheckoutID: view.CheckoutID, OperationID: op, RunID: v.RunID, ReportID: v.ID, ReportedRevision: v.ReportedRevision, InputRevision: v.InputRevision, ExpectedRevision: v.CurrentRevision, ExpectedGeneration: view.Generation, Coordinator: "fixture-leader", Note: "Report received; independent verification remains."}
}
func stateHead(t *testing.T, s *Service) string {
	t.Helper()
	b, e := os.ReadFile(filepath.Join(s.state, "HEAD"))
	if e != nil {
		t.Fatal(e)
	}
	return string(b)
}
func TestReportAcknowledgementPersistenceRetryAndNoAdvancement(t *testing.T) {
	s, sel, r, q := reportFixture(t)
	original, _, e := s.load()
	if e != nil {
		t.Fatal(e)
	}
	out, e := s.AcknowledgeReport(LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	if !reflect.DeepEqual(out.Run, r.Run) {
		t.Fatal("acknowledgement changed run", out, r)
	}
	state, _, e := s.load()
	if e != nil || !reflect.DeepEqual(state.Runs, original.Runs) {
		t.Fatal("run state changed", e)
	}
	s, e = New(s.root, s.area, s.project, s.state)
	if e != nil {
		t.Fatal(e)
	}
	view, e := s.ReportInbox("")
	if e != nil || view.Schema != "sdp-governance-inbox/2" || len(view.Reports) != 1 || view.Reports[0].AcknowledgementID != q.OperationID || !reflect.DeepEqual(view.Reports[0].Acknowledgement, out.ReportAcknowledgement) {
		t.Fatal(view, e)
	}
	head := stateHead(t, s)
	again, e := s.AcknowledgeReport(LocalController(), q)
	if e != nil || !reflect.DeepEqual(again, out) || stateHead(t, s) != head {
		t.Fatal("retry changed result/state", again, e)
	}
	altered := q
	altered.Note = "different"
	_, e = s.AcknowledgeReport(LocalController(), altered)
	code(t, e, "idempotency_conflict")
	altered = q
	altered.OperationID, _ = identifier()
	altered.ExpectedGeneration = head
	_, e = s.AcknowledgeReport(LocalController(), altered)
	code(t, e, "conflict")
	// Run progression stays guarded and is a separate operation.
	canceled, e := s.Transition(LocalController(), transition(r, "cancel"))
	if e != nil || canceled.Run.State != "canceled" {
		t.Fatal(canceled, e)
	}
	if e = os.WriteFile(filepath.Join(s.root, sel.Authority), []byte("new authority"), 0600); e != nil {
		t.Fatal(e)
	}
	head = stateHead(t, s)
	again, e = s.AcknowledgeReport(LocalController(), q)
	if e != nil || !reflect.DeepEqual(again, out) || stateHead(t, s) != head {
		t.Fatal("historical retry lost", e)
	}
	view, e = s.ReportInbox("")
	if e != nil || view.Reports[0].Coverage != "stale" || view.Reports[0].AcknowledgementID != q.OperationID || view.Reports[0].CurrentState != "canceled" {
		t.Fatal(view, e)
	}
}
func TestReportAcknowledgementRefusesWithoutMutation(t *testing.T) {
	for _, change := range []string{"untrusted", "agent", "checkout", "project", "generation", "revision", "source", "cancel", "report", "run", "reported-revision", "input", "operation-collision", "profile", "note"} {
		t.Run(change, func(t *testing.T) {
			s, sel, r, q := reportFixture(t)
			caller := LocalController()
			expected := "stale_revision"
			switch change {
			case "untrusted":
				caller = Caller{}
				expected = "unauthorized"
			case "agent":
				caller = Caller{runID: r.Run.ID}
				expected = "unauthorized"
			case "checkout":
				q.CheckoutID, _ = identifier()
				expected = "unauthorized"
			case "project":
				q.ProjectRef = "elsewhere"
				expected = "unauthorized"
			case "generation":
				q.ExpectedGeneration = hash([]byte("old"))
			case "revision":
				q.ExpectedRevision++
			case "source":
				if e := os.WriteFile(filepath.Join(s.root, sel.Authority), []byte("changed"), 0600); e != nil {
					t.Fatal(e)
				}
			case "cancel":
				var e error
				r, e = s.Transition(LocalController(), transition(r, "cancel"))
				if e != nil {
					t.Fatal(e)
				}
				q.ExpectedGeneration = stateHead(t, s)
				q.ExpectedRevision = r.Run.Revision
				expected = "conflict"
			case "report":
				q.ReportID, _ = identifier()
				expected = "missing_coverage"
			case "run":
				q.RunID, _ = identifier()
				expected = "missing_coverage"
			case "reported-revision":
				q.ReportedRevision--
				expected = "missing_coverage"
			case "input":
				q.InputRevision = hash([]byte("other"))
				expected = "missing_coverage"
			case "operation-collision":
				q.OperationID = q.ReportID
				expected = "idempotency_conflict"
			case "profile":
				q.Schema = "unknown"
				expected = "incompatible_profile"
			case "note":
				q.Note = "  "
				expected = "invalid_request"
			}
			head := stateHead(t, s)
			_, e := s.AcknowledgeReport(caller, q)
			code(t, e, expected)
			if stateHead(t, s) != head {
				t.Fatal("refused acknowledgement mutated state")
			}
		})
	}
}
func TestReportAcknowledgementPreservesAttentionAndContinuation(t *testing.T) {
	s, _, _, q := reportFixture(t)
	// An unrelated run starts before acknowledgement; refresh only the explicit test basis.
	r, v := continuationFixture(t, s)
	q.ExpectedGeneration = stateHead(t, s)
	if _, e := s.AcknowledgeReport(LocalController(), q); e != nil {
		t.Fatal(e)
	}
	if _, e := continuationInput(t, s, r, v, "prepare"); e != nil {
		t.Fatal(e)
	}
	state, _, e := s.load()
	if e != nil || state.Schema != reportAcknowledgementStateSchema {
		t.Fatal(state.Schema, e)
	}
	// New attention must also preserve the newer store profile.
	_, _, start := fixture(t)
	c, e := s.Context(r.Run.Selection)
	if e != nil {
		t.Fatal(e)
	}
	start.CheckoutID = c.CheckoutID
	start.InputRevision = c.Revision
	other, e := s.Start(LocalController(), start)
	if e != nil {
		t.Fatal(e)
	}
	other, e = s.Transition(LocalController(), transition(other, "start"))
	if e != nil {
		t.Fatal(e)
	}
	other, e = s.RequestAttention(LocalController(), questionRequest(other))
	if e != nil {
		t.Fatal(e)
	}
	if _, e = s.ReplyAttention(LocalController(), questionDisposition(other, 0), OwnerAttentionInput{Actor: "fixture", Source: "test", Text: "continue"}); e != nil {
		t.Fatal(e)
	}
	state, _, e = s.load()
	if e != nil || state.Schema != reportAcknowledgementStateSchema {
		t.Fatal(state.Schema, e)
	}
}
func TestReportAcknowledgementRejectsCorruptHistory(t *testing.T) {
	for _, change := range []string{"missing-record", "missing-report", "changed-note", "duplicate", "wrong-profile", "no-ack", "wrong-revision"} {
		t.Run(change, func(t *testing.T) {
			s, _, _, q := reportFixture(t)
			if _, e := s.AcknowledgeReport(LocalController(), q); e != nil {
				t.Fatal(e)
			}
			state, head, e := s.load()
			if e != nil {
				t.Fatal(e)
			}
			receipt := state.Receipts[q.OperationID]
			switch change {
			case "missing-record":
				receipt.Result.ReportAcknowledgement = nil
			case "missing-report":
				delete(state.Receipts, q.ReportID)
			case "changed-note":
				receipt.Result.ReportAcknowledgement.Request.Note = "tampered"
			case "duplicate":
				id, _ := identifier()
				state.Receipts[id] = receipt
			case "wrong-profile":
				state.Schema = Schema
			case "no-ack":
				delete(state.Receipts, q.OperationID)
			case "wrong-revision":
				receipt.Result.Run.Revision++
			}
			if change != "no-ack" {
				state.Receipts[q.OperationID] = receipt
			}
			if e = s.save(state, head); e != nil {
				t.Fatal(e)
			}
			if _, _, e = s.load(); e == nil {
				t.Fatal("accepted corrupt acknowledgement")
			}
		})
	}
}

func acknowledgeFixtureOnService(t *testing.T, s *Service) {
	t.Helper()
	r, _ := continuationFixture(t, s)
	r, e := s.Transition(LocalController(), transition(r, "report"))
	if e != nil {
		t.Fatal(e)
	}
	if _, e = s.AcknowledgeReport(LocalController(), acknowledgementRequest(t, s, r)); e != nil {
		t.Fatal(e)
	}
}
func TestAcknowledgementProfileSurvivesPreparationAndDesignWriters(t *testing.T) {
	t.Run("preparation", func(t *testing.T) {
		s, _, other, p := preparedFixture(t)
		acknowledgeFixtureOnService(t, s)
		if _, e := s.PrepareAssignment(LocalController(), refreshPreparation(t, s, other)); e != nil {
			t.Fatal(e)
		}
		updated, e := s.ReviseAssignment(LocalController(), revisionRequest(t, s, p.Assignment))
		if e != nil {
			t.Fatal(e)
		}
		if _, e = s.RetireAssignment(LocalController(), RetirePreparationRequest(preparationChange(t, s, updated.Assignment))); e != nil {
			t.Fatal(e)
		}
		state, _, e := s.load()
		if e != nil || state.Schema != reportAcknowledgementStateSchema {
			t.Fatal(state.Schema, e)
		}
	})
	t.Run("design-and-activation", func(t *testing.T) {
		f, q, owner := decisionFixture(t)
		acknowledgeFixtureOnService(t, f.s)
		q.Review = freshDesignQuery(t, f)
		d, e := f.s.RecordDesignDecision(LocalController(), q, owner)
		if e != nil {
			t.Fatal(e)
		}
		op, _ := identifier()
		req, _ := identifier()
		_, e = f.s.ActivateAssignment(LocalController(), ActivationRequest{Schema: ActivationSchema, ProjectRef: q.ProjectRef, CheckoutID: q.CheckoutID, OperationID: op, RequestID: req, Review: freshDesignQuery(t, f), ReviewRevision: q.ReviewRevision, DecisionOperationID: d.OperationID})
		if e != nil {
			t.Fatal(e)
		}
		state, _, e := f.s.load()
		if e != nil || state.Schema != reportAcknowledgementStateSchema {
			t.Fatal(state.Schema, e)
		}
	})
}
func TestReportAcknowledgementConcurrentDelivery(t *testing.T) {
	s, _, _, q := reportFixture(t)
	other := q
	other.OperationID, _ = identifier()
	results := make(chan error, 2)
	for _, request := range []ReportAcknowledgementRequest{q, other} {
		go func(request ReportAcknowledgementRequest) {
			_, e := s.AcknowledgeReport(LocalController(), request)
			results <- e
		}(request)
	}
	successes := 0
	for i := 0; i < 2; i++ {
		if <-results == nil {
			successes++
		}
	}
	if successes != 1 {
		t.Fatal("duplicate delivery", successes)
	}
	view, e := s.ReportInbox("")
	if e != nil || len(view.Reports) != 1 || view.Reports[0].Acknowledgement == nil {
		t.Fatal(view, e)
	}
}
