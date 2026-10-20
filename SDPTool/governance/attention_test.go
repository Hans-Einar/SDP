package governance

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func attentionFixture(t *testing.T) (*Service, Selection, Result) {
	t.Helper()
	s, sel, q := fixture(t)
	bind(t, s, &q)
	r, e := s.Start(LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	r, e = s.Transition(LocalController(), transition(r, "start"))
	if e != nil {
		t.Fatal(e)
	}
	return s, sel, r
}
func questionRequest(r Result) AttentionRequest {
	id, _ := identifier()
	return AttentionRequest{Schema: AttentionSchema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, RunID: r.Run.ID, OperationID: id, InputRevision: r.Run.InputRevision, ExpectedRevision: r.Run.Revision, Prompt: "Which label should the output use?", Options: []AttentionOption{{ID: "short", Label: "Short label"}, {ID: "long", Label: "Detailed label"}}, Blocking: true}
}
func questionDisposition(r Result, index int) AttentionDisposition {
	id, _ := identifier()
	q := r.Run.Attention[index]
	return AttentionDisposition{Schema: AttentionSchema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, RunID: r.Run.ID, OperationID: id, QuestionID: q.ID, Basis: q.Basis}
}
func TestAttentionDurableReplyAndImmutableRetry(t *testing.T) {
	s, _, r := attentionFixture(t)
	ask := questionRequest(r)
	submitted, e := s.RequestAttention(LocalController(), ask)
	if e != nil {
		t.Fatal(e)
	}
	_, e = s.Transition(LocalController(), transition(submitted, "report"))
	code(t, e, "attention_required")
	again, e := s.RequestAttention(LocalController(), ask)
	if e != nil || !reflect.DeepEqual(again, submitted) {
		t.Fatal(again, e)
	}
	s, e = New(s.root, s.area, s.project, s.state)
	if e != nil {
		t.Fatal(e)
	}
	view, e := s.AttentionInbox(r.Run.ID)
	if e != nil || len(view.Items) != 1 || view.Items[0].Status != "open" || view.Items[0].Applicability != "current" {
		t.Fatal(view, e)
	}
	q := questionDisposition(submitted, 0)
	input := OwnerAttentionInput{Actor: "fixture-owner", Source: "test-controller", OptionID: "short", Text: "Use the short label for this fixture."}
	replied, e := s.ReplyAttention(LocalController(), q, input)
	if e != nil {
		t.Fatal(e)
	}
	if replied.Run.Attention[0].Reply.Text != input.Text || submitted.Run.Attention[0].Reply != nil {
		t.Fatal("history changed")
	}
	again, e = s.ReplyAttention(LocalController(), q, input)
	if e != nil || !reflect.DeepEqual(again, replied) {
		t.Fatal(again, e)
	}
	input.Actor = "another-owner"
	_, e = s.ReplyAttention(LocalController(), q, input)
	code(t, e, "idempotency_conflict")
	q.OperationID, _ = identifier()
	_, e = s.ReplyAttention(LocalController(), q, input)
	code(t, e, "conflict")
	if _, e = s.Transition(LocalController(), transition(replied, "report")); e != nil {
		t.Fatal(e)
	}
	state, _, e := s.load()
	if e != nil || state.Schema != attentionStateSchema {
		t.Fatal(state.Schema, e)
	}
	old, e := s.RequestAttention(LocalController(), ask)
	if e != nil || old.Run.Attention[0].Reply != nil {
		t.Fatal("historical ask changed", e)
	}
}
func TestAttentionStaleRefusalAndExplicitWithdrawal(t *testing.T) {
	for _, change := range []string{"source", "cancel", "basis", "option", "wrong-run", "wrong-checkout", "untrusted"} {
		t.Run(change, func(t *testing.T) {
			s, sel, r := attentionFixture(t)
			r, e := s.RequestAttention(LocalController(), questionRequest(r))
			if e != nil {
				t.Fatal(e)
			}
			q := questionDisposition(r, 0)
			original := q
			input := OwnerAttentionInput{Actor: "fixture-owner", Source: "test", Text: "A free-text answer."}
			caller := LocalController()
			switch change {
			case "source":
				e = os.WriteFile(filepath.Join(s.root, sel.Authority), []byte("New owner scope"), 0600)
			case "cancel":
				_, e = s.Transition(LocalController(), transition(r, "cancel"))
			case "basis":
				q.Basis = hash([]byte("wrong"))
			case "option":
				input.OptionID = "missing"
			case "wrong-run":
				q.RunID, _ = identifier()
			case "wrong-checkout":
				q.CheckoutID, _ = identifier()
			case "untrusted":
				caller = Caller{}
			}
			if e != nil {
				t.Fatal(e)
			}
			before, _ := os.ReadFile(filepath.Join(s.state, "HEAD"))
			_, e = s.ReplyAttention(caller, q, input)
			if e == nil {
				t.Fatal("invalid reply accepted", change)
			}
			after, _ := os.ReadFile(filepath.Join(s.state, "HEAD"))
			if string(before) != string(after) {
				t.Fatal("refusal wrote state")
			}
			withdrawn, e := s.WithdrawAttention(LocalController(), original, "Obsolete fixture question; retain its history.")
			if e != nil || withdrawn.Run.Attention[0].Withdrawal == nil {
				t.Fatal(withdrawn, e)
			}
			view, e := s.AttentionInbox(r.Run.ID)
			if e != nil || view.Items[0].Status != "withdrawn" {
				t.Fatal(view, e)
			}
		})
	}
}
func TestAttentionSiblingAndOtherRunIndependence(t *testing.T) {
	s, _, r := attentionFixture(t)
	r, e := s.RequestAttention(LocalController(), questionRequest(r))
	if e != nil {
		t.Fatal(e)
	}
	r, e = s.RequestAttention(LocalController(), questionRequest(r))
	if e != nil {
		t.Fatal(e)
	}
	for i := 0; i < 2; i++ {
		r, e = s.ReplyAttention(LocalController(), questionDisposition(r, i), OwnerAttentionInput{Actor: "fixture-owner", Source: "test", Text: "Keep the current scope."})
		if e != nil {
			t.Fatal(e)
		}
	}
	if r, e = s.Transition(LocalController(), transition(r, "report")); e != nil {
		t.Fatal(e)
	}
	// A blocking question on the first run cannot stop another authorized run.
	r, e = s.Transition(LocalController(), transition(r, "fail"))
	if e != nil {
		t.Fatal(e)
	}
	r, e = s.Transition(LocalController(), transition(r, "retry"))
	if e != nil {
		t.Fatal(e)
	}
	r, e = s.Transition(LocalController(), transition(r, "start"))
	if e != nil {
		t.Fatal(e)
	}
	r, e = s.RequestAttention(LocalController(), questionRequest(r))
	if e != nil {
		t.Fatal(e)
	}
	_, _, start := fixture(t)
	start.Selection = r.Run.Selection
	start.Assignment = r.Run.Assignment
	start.InputRevision = r.Run.InputRevision
	start.ProjectRef = r.ProjectRef
	start.CheckoutID = r.CheckoutID
	other, e := s.Start(LocalController(), start)
	if e != nil {
		t.Fatal(e)
	}
	other, e = s.Transition(LocalController(), transition(other, "start"))
	if e != nil {
		t.Fatal(e)
	}
	if _, e = s.Transition(LocalController(), transition(other, "report")); e != nil {
		t.Fatal(e)
	}
	_, e = s.Transition(LocalController(), transition(r, "report"))
	code(t, e, "attention_required")
}
func TestAttentionEmptyAndValidation(t *testing.T) {
	s, _, _ := fixture(t)
	_, denied := s.RequestAttention(Caller{}, AttentionRequest{})
	code(t, denied, "unauthorized")
	v, e := s.AttentionInbox("")
	if e != nil || len(v.Items) != 0 {
		t.Fatal(v, e)
	}
	if _, e = os.Stat(s.state); !os.IsNotExist(e) {
		t.Fatal("empty read wrote state")
	}
	s, _, r := attentionFixture(t)
	for _, change := range []string{"empty", "duplicate", "many", "stale"} {
		q := questionRequest(r)
		switch change {
		case "empty":
			q.Prompt = " "
		case "duplicate":
			q.Options = append(q.Options, q.Options[0])
		case "many":
			q.Options = make([]AttentionOption, 6)
		case "stale":
			q.ExpectedRevision++
		}
		if _, e = s.RequestAttention(LocalController(), q); e == nil {
			t.Fatal("invalid question accepted", change)
		}
	}
}

// Add a maintenance run without modifying the shared source fixture. Exercise
// profile upgrades against real persisted preparation/design stores.
func addAttentionRun(t *testing.T, s *Service) Result {
	t.Helper()
	_, sel, q := fixture(t)
	c, e := s.Context(sel)
	if e != nil || len(c.Conditions) > 0 {
		t.Fatal(c, e)
	}
	q.ProjectRef = s.project
	q.CheckoutID = c.CheckoutID
	q.InputRevision = c.Revision
	r, e := s.Start(LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	r, e = s.Transition(LocalController(), transition(r, "start"))
	if e != nil {
		t.Fatal(e)
	}
	r, e = s.RequestAttention(LocalController(), questionRequest(r))
	if e != nil {
		t.Fatal(e)
	}
	return r
}
func TestAttentionProfileSurvivesOtherWriters(t *testing.T) {
	t.Run("preparation", func(t *testing.T) {
		s, q, other := preparationFixture(t)
		if _, e := s.Bind(LocalController()); e != nil {
			t.Fatal(e)
		}
		addAttentionRun(t, s)
		p, e := s.PrepareAssignment(LocalController(), refreshPreparation(t, s, q))
		if e != nil {
			t.Fatal(e)
		}
		p, e = s.ReviseAssignment(LocalController(), revisionRequest(t, s, p.Assignment))
		if e != nil {
			t.Fatal(e)
		}
		if _, e = s.RetireAssignment(LocalController(), RetirePreparationRequest(preparationChange(t, s, p.Assignment))); e != nil {
			t.Fatal(e)
		}
		if _, e = s.PrepareAssignment(LocalController(), refreshPreparation(t, s, other)); e != nil {
			t.Fatal(e)
		}
		state, _, e := s.load()
		if e != nil || state.Schema != attentionStateSchema {
			t.Fatal(state.Schema, e)
		}
	})
	t.Run("design-and-activation", func(t *testing.T) {
		f, q, input := decisionFixture(t)
		addAttentionRun(t, f.s)
		q.Review = freshDesignQuery(t, f)
		decision, e := f.s.RecordDesignDecision(LocalController(), q, input)
		if e != nil {
			t.Fatal(e)
		}
		op, _ := identifier()
		request, _ := identifier()
		a := ActivationRequest{Schema: ActivationSchema, ProjectRef: q.ProjectRef, CheckoutID: q.CheckoutID, OperationID: op, RequestID: request, Review: freshDesignQuery(t, f), ReviewRevision: q.ReviewRevision, DecisionOperationID: decision.OperationID}
		if _, e = f.s.ActivateAssignment(LocalController(), a); e != nil {
			t.Fatal(e)
		}
		state, _, e := f.s.load()
		if e != nil || state.Schema != attentionStateSchema {
			t.Fatal(state.Schema, e)
		}
	})
}

func TestAttentionBlockingAndInformationalQuestions(t *testing.T) {
	for _, blocking := range []bool{false, true} {
		s, _, r := attentionFixture(t)
		q := questionRequest(r)
		q.Blocking = blocking
		r, e := s.RequestAttention(LocalController(), q)
		if e != nil {
			t.Fatal(e)
		}
		r, e = s.Transition(LocalController(), transition(r, "interrupt"))
		if e != nil {
			t.Fatal(e)
		}
		resumed, e := s.Transition(LocalController(), transition(r, "resume"))
		if blocking {
			code(t, e, "attention_required")
		} else if e != nil || resumed.Run.State != "ready" {
			t.Fatal(resumed, e)
		}
	}
}
