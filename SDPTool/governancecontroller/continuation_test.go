package governancecontroller

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
	"github.com/Hans-Einar/SDP/SDPTool/governancehost"
)

type continuationHost struct {
	*fakeHost
	threads      int
	beforeThread func()
}

func (h *continuationHost) StartWorkThread(ctx context.Context, root, sandbox, instructions, approvals string, config map[string]any) (governancehost.Thread, error) {
	if h.beforeThread != nil {
		h.beforeThread()
	}
	v, e := h.fakeHost.StartWorkThread(ctx, root, sandbox, instructions, approvals, config)
	h.threads++
	v.ID = fmt.Sprintf("thread-%d", h.threads)
	return v, e
}
func waitingFixture(t *testing.T) (*Controller, *continuationHost, string, governance.AttentionDisposition) {
	t.Helper()
	c, f, id := setup(t)
	c.Approvals = "never"
	h := &continuationHost{fakeHost: f}
	c.Host = h
	r, e := c.Launch(context.Background(), id, "Fix the selected fixture")
	if e != nil {
		t.Fatal(e)
	}
	agent, _ := governance.AssignedAgent(id, r.Run.Attempts[0].ID, c.Epoch)
	r, e = c.Service.RequestAttention(agent, governance.AttentionRequest{Schema: governance.AttentionSchema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, RunID: id, OperationID: governancehost.UUID(), InputRevision: r.Run.InputRevision, ExpectedRevision: r.Run.Revision, Prompt: "Which output label?", Blocking: true})
	if e != nil {
		t.Fatal(e)
	}
	q := governance.AttentionDisposition{Schema: governance.AttentionSchema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, RunID: id, OperationID: governancehost.UUID(), QuestionID: r.Run.Attention[0].ID, Basis: r.Run.Attention[0].Basis}
	_, e = c.Observe(id, governancehost.Event{ThreadID: "thread-1", Observation: governance.Observation{Kind: "turn.completed", TurnID: "turn", Epoch: c.Epoch, ReceivedAt: now(), Status: "completed"}})
	if e != nil {
		t.Fatal(e)
	}
	return c, h, id, q
}
func answerFixture(t *testing.T, c *Controller, q governance.AttentionDisposition) {
	t.Helper()
	_, e := c.Service.ReplyAttention(governance.LocalController(), q, governance.OwnerAttentionInput{Actor: "fixture-owner", Source: "test", Text: "Use the short label."})
	if e != nil {
		t.Fatal(e)
	}
}
func TestAttentionContinuationPreservesRunAndPreventsDuplicates(t *testing.T) {
	c, h, id, q := waitingFixture(t)
	waiting, ready, e := c.AttentionPending(id)
	if e != nil || !waiting || ready {
		t.Fatal(waiting, ready, e)
	}
	if _, e = c.ContinueAttention(context.Background(), id); e == nil {
		t.Fatal("unanswered dispatch")
	}
	answerFixture(t, c, q)
	waiting, ready, e = c.AttentionPending(id)
	if e != nil || waiting || !ready {
		t.Fatal(waiting, ready, e)
	}
	r, e := c.ContinueAttention(context.Background(), id)
	if e != nil {
		t.Fatal(e)
	}
	if r.Run.ID != id || r.Run.State != "running" || len(r.Run.Attempts) != 2 || h.starts != 2 || r.Run.Attempts[0].Status != "completed" {
		t.Fatal(r, h.starts)
	}
	input, e := c.Service.OwnerInput(governance.LocalController(), id, r.Run.Inputs[1].ID)
	if e != nil || input.State != "acknowledged" || len(input.AttentionReplies) != 1 || input.AttentionReplies[0].ReplyOperationID != q.OperationID {
		t.Fatal(input, e)
	}
	var payload map[string]any
	if e = json.Unmarshal([]byte(input.Payload), &payload); e != nil || len(payload["attentionContext"].([]any)) != 1 {
		t.Fatal(payload, e)
	}
	view, e := c.Service.AttentionInbox(id)
	if e != nil || view.Items[0].Consumption != "dispatch-acknowledged" || view.Items[0].Applicability != "current" {
		t.Fatal(view, e)
	}
	for i := 0; i < 2; i++ {
		if _, e = c.ContinueAttention(context.Background(), id); e == nil {
			t.Fatal("duplicate continuation")
		}
	}
	if h.starts != 2 {
		t.Fatal(h.starts)
	}
	if dir := os.Getenv("PGI_SCHEMA_EVIDENCE"); dir != "" {
		raw, _ := json.Marshal([]map[string]any{{"definition": "OwnerInput", "value": input}, {"definition": "Result", "value": r}, {"definition": "AttentionInbox", "value": view}})
		if e = os.WriteFile(filepath.Join(dir, "continuation.json"), raw, 0600); e != nil {
			t.Fatal(e)
		}
	}
}

func TestAttentionContinuationCurrentReplyBecomesStaleAfterChanges(t *testing.T) {
	for _, change := range []string{"source", "steer", "scope-change"} {
		t.Run(change, func(t *testing.T) {
			c, _, id, q := waitingFixture(t)
			answerFixture(t, c, q)
			if _, e := c.ContinueAttention(context.Background(), id); e != nil {
				t.Fatal(e)
			}
			if change == "source" {
				coverage, e := c.Service.RunContext(id)
				if e != nil {
					t.Fatal(e)
				}
				if e = os.WriteFile(filepath.Join(coverage.Root, coverage.Selection.Authority), []byte("Changed owner mandate"), 0600); e != nil {
					t.Fatal(e)
				}
			} else if _, e := c.Control(context.Background(), id, change, "Changed owner direction"); e != nil {
				t.Fatal(e)
			}
			view, e := c.Service.AttentionInbox(id)
			if e != nil || view.Items[0].Applicability != "stale" || view.Items[0].Consumption != "dispatch-acknowledged" {
				t.Fatal(view, e)
			}
		})
	}
}
func TestAttentionContinuationRefusesChangedOrUncertainWork(t *testing.T) {
	for _, change := range []string{"policy", "source", "gap", "scope", "lost", "ready-drift"} {
		t.Run(change, func(t *testing.T) {
			c, h, id, q := waitingFixture(t)
			answerFixture(t, c, q)
			switch change {
			case "policy":
				c.Sandbox = "workspace-write"
			case "source":
				coverage, _ := c.Service.RunContext(id)
				if e := os.WriteFile(filepath.Join(coverage.Root, coverage.Selection.Authority), []byte("changed authority"), 0600); e != nil {
					t.Fatal(e)
				}
			case "gap":
				if _, e := c.Disconnect(id); e != nil {
					t.Fatal(e)
				}
			case "scope":
				if _, e := c.transition(id, "scope-change", "new scope requested"); e != nil {
					t.Fatal(e)
				}
			case "lost":
				h.lost = true
			case "ready-drift":
				h.beforeThread = func() {
					coverage, _ := c.Service.RunContext(id)
					if e := os.WriteFile(filepath.Join(coverage.Root, coverage.Selection.Authority), []byte("changed after claim"), 0600); e != nil {
						t.Fatal(e)
					}
				}
			}
			if _, e := c.ContinueAttention(context.Background(), id); e == nil {
				t.Fatal("invalid continuation")
			}
			if _, e := c.ContinueAttention(context.Background(), id); e == nil {
				t.Fatal("invalid retry")
			}
			want := 1
			if change == "lost" {
				want = 2
			}
			if h.starts != want {
				t.Fatal("unexpected host turn", h.starts, want)
			}
			if change == "lost" || change == "ready-drift" {
				view, e := c.Service.AttentionInbox(id)
				if e != nil || view.Items[0].Consumption != "dispatch-uncertain" {
					t.Fatal(view, e)
				}
			}
		})
	}
}

func TestCompiledExplicitAttentionContinuation(t *testing.T) {
	c, _, id, q := waitingFixture(t)
	answerFixture(t, c, q)
	coverage, e := c.Service.RunContext(id)
	if e != nil {
		t.Fatal(e)
	}
	binary := filepath.Join(t.TempDir(), "sdp-controller")
	command := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "build", "-o", binary, "../cmd/sdp-controller")
	if out, e := command.CombinedOutput(); e != nil {
		t.Fatal(e, string(out))
	}
	args := []string{"continue", "--project", coverage.Root, "--state", filepath.Join(filepath.Dir(coverage.Root), "state"), "--run", id, "--codex", os.Args[0], "--mcp-binary", "fixture-mcp", "--sandbox", "read-only", "--approvals", "never", "--timeout", "20s"}
	command = exec.Command(binary, args...)
	command.Env = append(os.Environ(), "SDP_PGI3_FAKE_HOST=1")
	if out, e := command.CombinedOutput(); e != nil {
		t.Fatal(e, string(out))
	}
	r, e := c.Service.Read(id)
	if e != nil || len(r.Run.Attempts) != 2 || r.Run.Attempts[1].Status != "completed" {
		t.Fatal(r, e)
	}
	headPath := filepath.Join(filepath.Dir(coverage.Root), "state", "HEAD")
	beforeHead, e := os.ReadFile(headPath)
	if e != nil {
		t.Fatal(e)
	}
	command = exec.Command(binary, args...)
	command.Env = append(os.Environ(), "SDP_PGI3_FAKE_HOST=1")
	if out, e := command.CombinedOutput(); e == nil {
		t.Fatal("duplicate CLI dispatch", string(out))
	}
	after, e := c.Service.Read(id)
	if e != nil || len(after.Run.Inputs) != 2 || after.Run.Attempts[1].Status != "completed" {
		t.Fatal(after, e)
	}
	afterHead, e := os.ReadFile(headPath)
	if e != nil || string(beforeHead) != string(afterHead) {
		t.Fatal("refusal changed state", e)
	}
}
