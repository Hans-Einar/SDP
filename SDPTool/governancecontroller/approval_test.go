package governancecontroller

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
	"github.com/Hans-Einar/SDP/SDPTool/governancehost"
)

type approvalHost struct {
	*fakeHost
	reply   func(governance.ApprovalReply) error
	replies int
}

func (f *approvalHost) ReplyApproval(v governance.ApprovalReply) error {
	f.replies++
	return f.reply(v)
}
func approvalFixture(t *testing.T) (*Controller, *approvalHost, string, governance.TransitionRequest, governance.ApprovalReply) {
	t.Helper()
	c, f, id := setup(t)
	h := &approvalHost{fakeHost: f}
	c.Host = h
	c.Approvals = "ask"
	r, e := c.Launch(context.Background(), id, "Fix the selected fixture")
	if e != nil {
		t.Fatal(e)
	}
	q := governance.TransitionRequest{Schema: governance.Schema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, RunID: id, OperationID: governancehost.UUID(), InputRevision: r.Run.InputRevision, ExpectedRevision: r.Run.Revision, Action: "report", Reason: "fixture report"}
	p := governance.ApprovalPrompt{RequestID: "0", Epoch: c.Epoch, ThreadID: "thread", TurnID: "turn", Message: `Allow the sdp_governance MCP server to run tool "run_transition"?`, Request: q}
	b, _ := json.Marshal(p)
	_, e = c.Observe(id, governancehost.Event{ThreadID: "thread", Observation: governance.Observation{Kind: "approval.requested", ItemID: "0", Epoch: c.Epoch, TurnID: "turn", Text: string(b), ReceivedAt: now()}})
	if e != nil {
		t.Fatal(e)
	}
	return c, h, id, q, governance.ApprovalReply{RequestID: "0", Epoch: c.Epoch, ThreadID: "thread", TurnID: "turn", Decision: "accept"}
}
func TestApprovalPreservesCASAndRequiresClosure(t *testing.T) {
	for _, scenario := range []string{"accept", "decline", "cancel", "changed-arguments", "steer", "scope-change", "drift", "write-loss", "wrong-epoch", "wrong-thread", "wrong-request", "resolved", "missing"} {
		t.Run(scenario, func(t *testing.T) {
			c, h, id, q, reply := approvalFixture(t)
			r, _ := c.Service.Read(id)
			agent, e := governance.AssignedAgent(id, r.Run.Attempts[0].ID, c.Epoch)
			if e != nil {
				t.Fatal(e)
			}
			// No implicit relaxation for either caller before actual owner input.
			if _, e = c.Service.Transition(agent, q); e == nil {
				t.Fatal("unapproved stale call advanced")
			}
			if _, e = c.Service.Transition(governance.LocalController(), q); e == nil {
				t.Fatal("controller CAS relaxed")
			}
			switch scenario {
			case "decline", "cancel":
				reply.Decision = scenario
			case "steer":
				if _, e = c.Control(context.Background(), id, "steer", "additional verification"); e != nil {
					t.Fatal(e)
				}
			case "scope-change":
				if _, e = c.Control(context.Background(), id, "scope-change", "separate new scope"); e != nil {
					t.Fatal(e)
				}
			case "drift":
				coverage, _ := c.Service.Context(r.Run.Selection)
				if e = os.WriteFile(filepath.Join(coverage.Root, r.Run.Selection.Authority), []byte("authority changed"), 0600); e != nil {
					t.Fatal(e)
				}
			case "wrong-epoch":
				reply.Epoch = governancehost.UUID()
			case "wrong-thread":
				reply.ThreadID = "other"
			case "wrong-request":
				reply.RequestID = "1"
			case "resolved":
				_, e = c.ResolveApproval(id, governancehost.Event{ThreadID: "thread", Observation: governance.Observation{Kind: "approval.resolved", ItemID: "0", TurnID: "turn", Epoch: c.Epoch, ReceivedAt: now()}})
				if e != nil {
					t.Fatal(e)
				}
			case "missing":
				return
			}
			h.reply = func(v governance.ApprovalReply) error {
				before, e := c.Service.Read(id)
				if e != nil {
					t.Fatal(e)
				}
				input := before.Run.Inputs[len(before.Run.Inputs)-1]
				if input.Kind != "approval" || input.State != "sending" {
					t.Fatal("effect before durable intent", input)
				}
				if scenario == "write-loss" {
					return errors.New("lost write")
				}
				sent := q
				if scenario == "changed-arguments" {
					sent.Reason = "changed"
				}
				result, e := c.Service.Transition(agent, sent)
				if scenario == "accept" {
					if e != nil || result.Run.State != "verifying" {
						t.Fatal("approved exact report refused", result, e)
					}
				} else if e == nil {
					t.Fatal("unapproved or invalidated report advanced", scenario)
				}
				return nil
			}
			result, e := c.Approve(id, reply, `{"action":"approval","scenario":"`+scenario+`"}`)
			refused := scenario == "drift" || scenario == "scope-change" || scenario == "wrong-epoch" || scenario == "wrong-thread" || scenario == "wrong-request" || scenario == "resolved"
			if refused {
				if e == nil || h.replies != 0 {
					t.Fatal("invalid response sent", result, e)
				}
				return
			}
			if scenario == "write-loss" {
				if e == nil {
					t.Fatal("write loss hidden")
				}
			} else if e != nil {
				t.Fatal(e)
			}
			result, _ = c.Service.Read(id)
			last := result.Run.Inputs[len(result.Run.Inputs)-1]
			if last.State != "uncertain" {
				t.Fatal("write mistaken for host acknowledgment", last)
			}
			if _, e = c.Approve(id, reply, "duplicate"); e == nil || h.replies != 1 {
				t.Fatal("response resent")
			}
			if scenario == "write-loss" {
				return
			}
			result, e = c.ResolveApproval(id, governancehost.Event{ThreadID: "thread", Observation: governance.Observation{Kind: "approval.resolved", ItemID: "0", TurnID: "turn", Epoch: c.Epoch, ReceivedAt: now()}})
			if e != nil || result.Run.Inputs[len(result.Run.Inputs)-1].State != "acknowledged" {
				t.Fatal(result, e)
			}
			if scenario == "accept" {
				retry, e := c.Service.Transition(agent, q)
				if e != nil || retry.Run.State != "verifying" {
					t.Fatal("receipt retry changed", e)
				}
			}
			if result.Run.State == "completed" {
				t.Fatal("approval granted SDP completion")
			}
		})
	}
}
