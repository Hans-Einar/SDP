package governancecontroller

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
	"github.com/Hans-Einar/SDP/SDPTool/governancehost"
)

func TestAttentionAcrossObservationsAndSteering(t *testing.T) {
	for _, change := range []string{"observations", "completion", "steer", "scope-change", "disconnect", "wrong-epoch", "wrong-run"} {
		t.Run(change, func(t *testing.T) {
			c, _, id := setup(t)
			r, e := c.Launch(context.Background(), id, "Fix the selected fixture")
			if e != nil {
				t.Fatal(e)
			}
			agent, _ := governance.AssignedAgent(id, r.Run.Attempts[0].ID, c.Epoch)
			ask := governance.AttentionRequest{Schema: governance.AttentionSchema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, OperationID: governancehost.UUID(), RunID: id, InputRevision: r.Run.InputRevision, ExpectedRevision: r.Run.Revision, Prompt: "Which output label should this task use?", Blocking: true}
			for _, kind := range []string{"tool.completed", "agentMessage"} {
				o := governance.Observation{Kind: kind, ItemID: kind, TurnID: "turn", Epoch: c.Epoch, ReceivedAt: now()}
				if kind == "tool.completed" {
					o.Status = "completed"
				} else {
					o.Text = "Waiting for clarification"
					o.Phase = "commentary"
				}
				if _, e = c.Observe(id, governancehost.Event{ThreadID: "thread", Observation: o}); e != nil {
					t.Fatal(e)
				}
			}
			caller := agent
			if change == "wrong-epoch" {
				caller, _ = governance.AssignedAgent(id, r.Run.Attempts[0].ID, governancehost.UUID())
			}
			if change == "wrong-run" {
				caller, _ = governance.AssignedAgent(governancehost.UUID(), r.Run.Attempts[0].ID, c.Epoch)
			}
			r, e = c.Service.RequestAttention(caller, ask)
			if change == "wrong-epoch" || change == "wrong-run" {
				if e == nil {
					t.Fatal("foreign capability accepted")
				}
				return
			}
			if e != nil {
				t.Fatal(e)
			}
			q := governance.AttentionDisposition{Schema: governance.AttentionSchema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, OperationID: governancehost.UUID(), RunID: id, QuestionID: r.Run.Attention[0].ID, Basis: r.Run.Attention[0].Basis}
			input := governance.OwnerAttentionInput{Actor: "fixture-owner", Source: "test", Text: "Use a short label."}
			if _, denied := c.Service.ReplyAttention(agent, q, input); denied == nil {
				t.Fatal("agent impersonated owner")
			}
			switch change {
			case "completion":
				_, e = c.Observe(id, governancehost.Event{ThreadID: "thread", Observation: governance.Observation{Kind: "turn.completed", TurnID: "turn", Epoch: c.Epoch, Status: "completed", ReceivedAt: now()}})
			case "steer", "scope-change":
				_, e = c.Control(context.Background(), id, change, "Owner changes the task constraint")
			case "disconnect":
				_, e = c.Disconnect(id)
			}
			if e != nil {
				t.Fatal(e)
			}
			view, e := c.Service.AttentionInbox(id)
			if e != nil {
				t.Fatal(e)
			}
			current := change == "observations" || change == "completion"
			if (view.Items[0].Applicability == "current") != current {
				t.Fatal(view)
			}
			replied, e := c.Service.ReplyAttention(governance.LocalController(), q, input)
			if (e == nil) != current {
				t.Fatal(replied, e)
			}
			if dir := os.Getenv("PGI_SCHEMA_EVIDENCE"); dir != "" && current {
				view, e = c.Service.AttentionInbox(id)
				if e != nil {
					t.Fatal(e)
				}
				raw, _ := json.Marshal([]map[string]any{{"definition": "AttentionRequest", "value": ask}, {"definition": "AttentionDisposition", "value": q}, {"definition": "AttentionInbox", "value": view}, {"definition": "Result", "value": replied}})
				if e = os.WriteFile(filepath.Join(dir, "attention-"+change+".json"), raw, 0600); e != nil {
					t.Fatal(e)
				}
			}
		})
	}
}
