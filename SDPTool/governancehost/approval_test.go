package governancehost

import (
	"bufio"
	"encoding/json"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
)

func approvalFrame(id any) map[string]any {
	return map[string]any{"id": id, "method": "mcpServer/elicitation/request", "params": map[string]any{"threadId": "thread", "turnId": "turn", "serverName": "sdp_governance", "mode": "form", "message": `Allow the sdp_governance MCP server to run tool "run_transition"?`, "requestedSchema": map[string]any{"type": "object", "properties": map[string]any{}}, "_meta": map[string]any{"codex_approval_kind": "mcp_tool_call", "persist": []string{"session", "always"}, "tool_params": governance.TransitionRequest{OperationID: UUID(), Action: "report"}}}}
}
func TestApprovalWireRoundTripAndInvalidation(t *testing.T) {
	var cases []map[string]any
	for _, decision := range []string{"accept", "decline", "cancel"} {
		t.Run(decision, func(t *testing.T) {
			left, right := net.Pipe()
			c := newClient(left, left)
			c.ask = true
			defer c.Close()
			defer right.Close()
			right.SetDeadline(time.Now().Add(5 * time.Second))
			responses := make(chan frame, 1)
			request := approvalFrame("host-0")
			go func() {
				json.NewEncoder(right).Encode(request)
				s := bufio.NewScanner(right)
				if s.Scan() {
					var f frame
					json.Unmarshal(s.Bytes(), &f)
					responses <- f
				}
			}()
			event := <-c.Notifications
			var p governance.ApprovalPrompt
			if e := json.Unmarshal([]byte(event.Observation.Text), &p); e != nil {
				t.Fatal(e)
			}
			reply := governance.ApprovalReply{RequestID: p.RequestID, ThreadID: p.ThreadID, TurnID: p.TurnID, Epoch: c.Epoch, Decision: decision}
			wrong := reply
			wrong.Epoch = UUID()
			if c.ReplyApproval(wrong) == nil {
				t.Fatal("old epoch replied")
			}
			if e := c.ReplyApproval(reply); e != nil {
				t.Fatal(e)
			}
			f := <-responses
			var result map[string]any
			json.Unmarshal(f.Result, &result)
			if string(f.ID) != `"host-0"` || result["action"] != decision || result["_meta"] != nil {
				t.Fatal("wrong or persistent response", f)
			}
			if c.ReplyApproval(reply) == nil {
				t.Fatal("duplicate response")
			}
			cases = append(cases, map[string]any{"definition": "McpServerElicitationRequestParams", "value": request["params"]}, map[string]any{"definition": "McpServerElicitationRequestResponse", "value": f.Result})
			go json.NewEncoder(right).Encode(map[string]any{"method": "serverRequest/resolved", "params": map[string]any{"requestId": "host-0", "threadId": "thread"}})
			resolved := <-c.Notifications
			if resolved.Observation.Kind != "approval.resolved" {
				t.Fatal(resolved)
			}
			if c.ReplyApproval(reply) == nil {
				t.Fatal("resolved response")
			}
		})
	}
	if dir := os.Getenv("PGI_SCHEMA_EVIDENCE"); dir != "" {
		b, _ := json.MarshalIndent(cases, "", "  ")
		if e := os.WriteFile(filepath.Join(dir, "approval-rpc-cases.json"), b, 0600); e != nil {
			t.Fatal(e)
		}
	}
}
func TestUnsupportedApprovalShapesRefuse(t *testing.T) {
	for _, scenario := range []string{"never", "server", "url", "form-fields", "tool", "owner-action", "bad-id", "bad-operation"} {
		t.Run(scenario, func(t *testing.T) {
			request := approvalFrame(0)
			params := request["params"].(map[string]any)
			switch scenario {
			case "server":
				params["serverName"] = "other"
			case "url":
				params["mode"] = "url"
			case "form-fields":
				params["requestedSchema"] = map[string]any{"type": "object", "properties": map[string]any{"secret": map[string]string{"type": "string"}}}
			case "tool":
				params["message"] = "Allow another tool"
			case "owner-action":
				params["_meta"].(map[string]any)["tool_params"] = governance.TransitionRequest{OperationID: UUID(), Action: "complete"}
			case "bad-id":
				request["id"] = nil
			case "bad-operation":
				params["_meta"].(map[string]any)["tool_params"] = governance.TransitionRequest{Action: "report"}
			}
			b, _ := json.Marshal(request)
			var f frame
			json.Unmarshal(b, &f)
			c := &Client{Epoch: UUID(), ask: scenario != "never", approvals: map[string]pendingApproval{}}
			if _, keep, _ := c.approvalRequest(f); keep {
				t.Fatal("unsupported prompt accepted")
			}
		})
	}
}
