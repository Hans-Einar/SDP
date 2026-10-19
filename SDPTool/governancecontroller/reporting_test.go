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

type reportHost struct {
	*fakeHost
	config    map[string]any
	approvals string
}

func (h *reportHost) StartWorkThread(ctx context.Context, root, sandbox, instructions, approvals string, config map[string]any) (governancehost.Thread, error) {
	h.config, h.approvals = config, approvals
	return h.fakeHost.StartWorkThread(ctx, root, sandbox, instructions, approvals, config)
}
func TestReportingPolicyIsExplicitAndToolScoped(t *testing.T) {
	for _, mode := range []string{"", "manual", "automatic", "invalid"} {
		t.Run(mode, func(t *testing.T) {
			c, f, id := setup(t)
			h := &reportHost{fakeHost: f}
			c.Host = h
			c.Reporting = mode
			c.Approvals = "never"
			r, err := c.Launch(context.Background(), id, "Fix the selected fixture")
			if mode == "invalid" {
				if err == nil || h.config != nil {
					t.Fatal("invalid policy dispatched")
				}
				r, _ = c.Service.Read(id)
				if len(r.Run.Inputs) != 0 {
					t.Fatal("invalid policy recorded")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			server := h.config["sdp_governance"].(map[string]any)
			policy, _ := json.Marshal(server["tools"])
			if mode == "automatic" {
				if string(policy) != `{"run_transition":{"approval_mode":"approve"}}` {
					t.Fatal(string(policy))
				}
			} else if string(policy) != "null" {
				t.Fatal("legacy policy changed", string(policy))
			}
			if len(h.config) != 1 || h.approvals != "never" {
				t.Fatal("broadened policy")
			}
			input, err := c.Service.OwnerInput(governance.LocalController(), id, r.Run.Inputs[0].ID)
			if err != nil {
				t.Fatal(err)
			}
			var payload map[string]any
			if json.Unmarshal([]byte(input.Payload), &payload) != nil {
				t.Fatal(input)
			}
			want := mode
			if want == "" {
				want = "manual"
			}
			if payload["reportingMode"] != want {
				t.Fatal(payload)
			}
		})
	}
}

func TestAutomaticReportAcrossOnlyObservations(t *testing.T) {
	for _, change := range []string{"observations", "controller", "steer", "scope-change", "disconnect", "authority", "approval", "fail"} {
		t.Run(change, func(t *testing.T) {
			c, _, id := setup(t)
			r, err := c.Launch(context.Background(), id, "Fix the selected fixture")
			if err != nil {
				t.Fatal(err)
			}
			agent, err := governance.AssignedAgent(id, r.Run.Attempts[0].ID, c.Epoch)
			if err != nil {
				t.Fatal(err)
			}
			request := governance.TransitionRequest{Schema: governance.Schema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, OperationID: governancehost.UUID(), RunID: id, ExpectedRevision: r.Run.Revision, InputRevision: r.Run.InputRevision, Action: "report", Reason: "Candidate and checks ready for review"}
			for _, kind := range []string{"tool.completed", "agentMessage"} {
				status, phase, text := "", "commentary", "observation"
				if kind == "tool.completed" {
					status, phase, text = "completed", "", ""
				}
				_, err = c.Observe(id, governancehost.Event{ThreadID: "thread", Observation: governance.Observation{Kind: kind, ItemID: kind, TurnID: "turn", Epoch: c.Epoch, Status: status, Text: text, Phase: phase, ReceivedAt: now()}})
				if err != nil {
					t.Fatal(err)
				}
			}
			switch change {
			case "controller":
				agent = governance.LocalController()
			case "steer":
				_, err = c.Control(context.Background(), id, "steer", "Owner adds a constraint")
			case "scope-change":
				_, err = c.Control(context.Background(), id, "scope-change", "Owner changes scope")
			case "disconnect":
				_, err = c.Disconnect(id)
			case "authority":
				coverage, _ := c.Service.Context(r.Run.Selection)
				err = os.WriteFile(filepath.Join(coverage.Root, r.Run.Selection.Authority), []byte("Owner changes authority"), 0600)
			case "fail":
				request.Action = "fail"
			case "approval":
				prompt := governance.ApprovalPrompt{RequestID: "approval", Epoch: c.Epoch, ThreadID: "thread", TurnID: "turn", Message: `Allow the sdp_governance MCP server to run tool "run_transition"?`, Request: request}
				data, _ := json.Marshal(prompt)
				_, err = c.Observe(id, governancehost.Event{ThreadID: "thread", Observation: governance.Observation{Kind: "approval.requested", ItemID: "approval", TurnID: "turn", Epoch: c.Epoch, Text: string(data), ReceivedAt: now()}})
			}
			if err != nil {
				t.Fatal(err)
			}
			got, err := c.Service.Transition(agent, request)
			if change != "observations" {
				if err == nil {
					t.Fatal("skipped governing change", change)
				}
				return
			}
			if err != nil || got.Run.State != "verifying" {
				t.Fatal(got, err)
			}
			again, err := c.Service.Transition(agent, request)
			if err != nil || again.Run.Revision != got.Run.Revision {
				t.Fatal(again, err)
			}
			inbox, err := c.Service.ReportInbox(id)
			if err != nil || len(inbox.Reports) != 1 || inbox.Reports[0].ID != request.OperationID {
				t.Fatal(inbox, err)
			}
		})
	}
}
