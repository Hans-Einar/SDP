package governance

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func continuationFixture(t *testing.T, s *Service) (Result, OwnerInput) {
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
	id, _ := identifier()
	attempt, _ := identifier()
	epoch, _ := identifier()
	v := OwnerInput{ID: id, Kind: "launch", Text: "Original fixture outcome", Payload: "Exact fixture payload", AttemptID: attempt, Epoch: epoch, ReceivedAt: time.Now().UTC().Format(time.RFC3339Nano)}
	input := func(action string) {
		t.Helper()
		op, _ := identifier()
		r, e = s.RecordInput(LocalController(), InputRequest{Schema: Schema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, OperationID: op, RunID: r.Run.ID, InputRevision: r.Run.InputRevision, ExpectedRevision: r.Run.Revision, Action: action, Input: v})
		if e != nil {
			t.Fatal(e)
		}
	}
	input("prepare")
	v.State = "prepared"
	input("claim")
	v.State = "sending"
	host := func(action, kind, status, turn string) {
		t.Helper()
		op, _ := identifier()
		r, e = s.RecordHost(LocalController(), HostRequest{Schema: Schema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, OperationID: op, RunID: r.Run.ID, InputRevision: r.Run.InputRevision, ExpectedRevision: r.Run.Revision, Action: action, AttemptID: attempt, HostVersion: "0.161.0", ThreadID: "fixture-thread", Observation: Observation{Kind: kind, Status: status, TurnID: turn, Epoch: epoch, ReceivedAt: time.Now().UTC().Format(time.RFC3339Nano)}})
		if e != nil {
			t.Fatal(e)
		}
	}
	host("bind", "bound", "idle", "")
	r, e = s.Transition(LocalController(), transition(r, "start"))
	if e != nil {
		t.Fatal(e)
	}
	host("observe", "turn.started", "running", "fixture-turn")
	v.ThreadID = "fixture-thread"
	v.TurnID = "fixture-turn"
	input("ack")
	r, e = s.RequestAttention(LocalController(), questionRequest(r))
	if e != nil {
		t.Fatal(e)
	}
	r, e = s.ReplyAttention(LocalController(), questionDisposition(r, 0), OwnerAttentionInput{Actor: "fixture-owner", Source: "test", Text: "Choose short"})
	if e != nil {
		t.Fatal(e)
	}
	host("observe", "turn.completed", "completed", "fixture-turn")
	v.ID, _ = identifier()
	v.AttemptID, _ = identifier()
	v.State = ""
	v.ThreadID = ""
	v.TurnID = ""
	question := r.Run.Attention[0]
	v.AttentionReplies = []AttentionReplyRef{{QuestionID: question.ID, ReplyOperationID: question.Reply.OperationID, Basis: question.Basis}}
	return r, v
}
func continuationInput(t *testing.T, s *Service, r Result, v OwnerInput, action string) (Result, error) {
	t.Helper()
	op, _ := identifier()
	return s.RecordInput(LocalController(), InputRequest{Schema: Schema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, OperationID: op, RunID: r.Run.ID, InputRevision: r.Run.InputRevision, ExpectedRevision: r.Run.Revision, Action: action, Input: v})
}
func TestAttentionContinuationCoreClaimsAndReopen(t *testing.T) {
	for _, change := range []string{"valid", "reply", "basis", "question", "duplicate", "text"} {
		t.Run(change, func(t *testing.T) {
			s, _, _ := fixture(t)
			if _, e := s.Bind(LocalController()); e != nil {
				t.Fatal(e)
			}
			r, v := continuationFixture(t, s)
			switch change {
			case "reply":
				v.AttentionReplies[0].ReplyOperationID, _ = identifier()
			case "basis":
				v.AttentionReplies[0].Basis = hash([]byte("wrong"))
			case "question":
				v.AttentionReplies[0].QuestionID, _ = identifier()
			case "duplicate":
				v.AttentionReplies = append(v.AttentionReplies, v.AttentionReplies[0])
			case "text":
				v.Text = "Different mandate"
			}
			before, _ := os.ReadFile(filepath.Join(s.state, "HEAD"))
			prepared, e := continuationInput(t, s, r, v, "prepare")
			if change != "valid" {
				if e == nil {
					t.Fatal("bad references accepted")
				}
				after, _ := os.ReadFile(filepath.Join(s.state, "HEAD"))
				if string(before) != string(after) {
					t.Fatal("refusal wrote")
				}
				return
			}
			if e != nil || prepared.Run.State != "ready" {
				t.Fatal(prepared, e)
			}
			reopened, e := New(s.root, s.area, s.project, s.state)
			if e != nil {
				t.Fatal(e)
			}
			duplicate := v
			duplicate.ID, _ = identifier()
			if _, e = continuationInput(t, reopened, prepared, duplicate, "prepare"); e == nil {
				t.Fatal("prepared delivery replayed")
			}
			v.State = "prepared"
			claimed, e := continuationInput(t, reopened, prepared, v, "claim")
			if e != nil {
				t.Fatal(e)
			}
			if _, e = continuationInput(t, reopened, claimed, v, "claim"); e == nil {
				t.Fatal("claimed delivery replayed")
			}
			view, e := reopened.AttentionInbox(r.Run.ID)
			if e != nil || view.Items[0].Consumption != "dispatch-sending" {
				t.Fatal(view, e)
			}
		})
	}
}
func TestContinuationClaimRechecksAuthority(t *testing.T) {
	s, _, _ := fixture(t)
	if _, e := s.Bind(LocalController()); e != nil {
		t.Fatal(e)
	}
	r, v := continuationFixture(t, s)
	r, e := continuationInput(t, s, r, v, "prepare")
	if e != nil {
		t.Fatal(e)
	}
	if e = os.WriteFile(filepath.Join(s.root, r.Run.Selection.Authority), []byte("changed after prepare"), 0600); e != nil {
		t.Fatal(e)
	}
	v.State = "prepared"
	_, e = continuationInput(t, s, r, v, "claim")
	code(t, e, "stale_revision")
}
func TestContinuationProfileSurvivesOtherWriters(t *testing.T) {
	t.Run("preparation", func(t *testing.T) {
		s, _, other, p := preparedFixture(t)
		r, v := continuationFixture(t, s)
		if _, e := continuationInput(t, s, r, v, "prepare"); e != nil {
			t.Fatal(e)
		}
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
		if e != nil || state.Schema != attentionContinuationStateSchema {
			t.Fatal(state.Schema, e)
		}
	})
	t.Run("design-activation-attention", func(t *testing.T) {
		f, q, owner := decisionFixture(t)
		r, v := continuationFixture(t, f.s)
		if _, e := continuationInput(t, f.s, r, v, "prepare"); e != nil {
			t.Fatal(e)
		}
		q.Review = freshDesignQuery(t, f)
		d, e := f.s.RecordDesignDecision(LocalController(), q, owner)
		if e != nil {
			t.Fatal(e)
		}
		op, _ := identifier()
		req, _ := identifier()
		work, e := f.s.ActivateAssignment(LocalController(), ActivationRequest{Schema: ActivationSchema, ProjectRef: q.ProjectRef, CheckoutID: q.CheckoutID, OperationID: op, RequestID: req, Review: freshDesignQuery(t, f), ReviewRevision: q.ReviewRevision, DecisionOperationID: d.OperationID})
		if e != nil {
			t.Fatal(e)
		}
		work, e = runChange(t, f.s, work, "start")
		if e != nil {
			t.Fatal(e)
		}
		if _, e = f.s.RequestAttention(LocalController(), questionRequest(work)); e != nil {
			t.Fatal(e)
		}
		state, _, e := f.s.load()
		if e != nil || state.Schema != attentionContinuationStateSchema {
			t.Fatal(state.Schema, e)
		}
	})
}
