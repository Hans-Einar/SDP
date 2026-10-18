package governancecontroller

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
	"github.com/Hans-Einar/SDP/SDPTool/governancehost"
)

// TestMain supplies a recorded fake app-server subprocess only for this test.
// It never invokes Codex or a model, and never represents independent review.
func TestMain(m *testing.M) {
	if os.Getenv("SDP_PGI3_FAKE_HOST") == "1" {
		if len(os.Args) > 1 && os.Args[1] == "--version" {
			fmt.Println("codex-cli " + governancehost.HostVersion)
			return
		}
		encoder := json.NewEncoder(os.Stdout)
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Buffer(make([]byte, 4096), 1<<20)
		item := map[string]any{"id": "final", "type": "agentMessage", "phase": "final_answer", "text": "Recorded fixture outcome"}
		for scanner.Scan() {
			var q struct {
				ID     any    `json:"id"`
				Method string `json:"method"`
			}
			if json.Unmarshal(scanner.Bytes(), &q) != nil {
				return
			}
			if q.ID == nil {
				continue
			}
			if q.Method == "" && os.Getenv("SDP_PGI3_APPROVAL") != "" {
				encoder.Encode(map[string]any{"method": "serverRequest/resolved", "params": map[string]any{"threadId": "thread", "requestId": 901}})
				encoder.Encode(map[string]any{"method": "turn/completed", "params": map[string]any{"threadId": "thread", "turn": map[string]any{"id": "turn", "status": "completed"}}})
				continue
			}
			var result any
			switch q.Method {
			case "initialize":
				result = map[string]any{"userAgent": "fixture", "platformOs": "linux"}
			case "account/read":
				result = map[string]any{"account": map[string]any{"type": "chatgpt"}}
			case "thread/start":
				result = map[string]any{"thread": map[string]any{"id": "thread", "status": map[string]any{"type": "idle"}}}
			case "mcpServerStatus/list":
				result = map[string]any{"data": []any{map[string]any{"name": "sdp_governance", "authStatus": "unsupported", "tools": map[string]any{"context_read": map[string]any{}, "events_read": map[string]any{}, "route_evaluate": map[string]any{}, "run_read": map[string]any{}, "run_transition": map[string]any{}}}}}
			case "turn/start":
				result = map[string]any{"turn": map[string]any{"id": "turn", "status": "inProgress"}}
			case "thread/read":
				result = map[string]any{"thread": map[string]any{"id": "thread", "status": map[string]any{"type": "idle"}, "turns": []any{map[string]any{"id": "turn", "status": "completed", "items": []any{item}}}}}
			default:
				encoder.Encode(map[string]any{"id": q.ID, "error": map[string]any{"code": -32601}})
				continue
			}
			encoder.Encode(map[string]any{"id": q.ID, "result": result})
			if q.Method == "turn/start" {
				if root := os.Getenv("SDP_PGI3_APPROVAL"); root != "" {
					var wire struct {
						Params struct {
							Input []struct {
								Text string `json:"text"`
							} `json:"input"`
						} `json:"params"`
					}
					json.Unmarshal(scanner.Bytes(), &wire)
					var payload struct {
						RunID string `json:"runId"`
					}
					json.Unmarshal([]byte(wire.Params.Input[0].Text), &payload)
					service, _ := governance.New(root, filepath.Join(root, "SDP"), "trial", filepath.Join(filepath.Dir(root), "state"))
					r, _ := service.Read(payload.RunID)
					request := governance.TransitionRequest{Schema: governance.Schema, ProjectRef: r.ProjectRef, CheckoutID: r.CheckoutID, RunID: r.Run.ID, InputRevision: r.Run.InputRevision, OperationID: governancehost.UUID(), ExpectedRevision: r.Run.Revision, Action: "report", Reason: "compiled fixture"}
					encoder.Encode(map[string]any{"id": 901, "method": "mcpServer/elicitation/request", "params": map[string]any{"threadId": "thread", "turnId": "turn", "serverName": "sdp_governance", "mode": "form", "message": `Allow the sdp_governance MCP server to run tool "run_transition"?`, "requestedSchema": map[string]any{"type": "object", "properties": map[string]any{}}, "_meta": map[string]any{"codex_approval_kind": "mcp_tool_call", "tool_params": request}}})
					continue
				}
				encoder.Encode(map[string]any{"method": "item/completed", "params": map[string]any{"threadId": "thread", "turnId": "turn", "item": item}})
				encoder.Encode(map[string]any{"method": "turn/completed", "params": map[string]any{"threadId": "thread", "turn": map[string]any{"id": "turn", "status": "completed"}}})
			}
		}
		return
	}
	os.Exit(m.Run())
}

func TestCompiledOwnerCommandAndRecovery(t *testing.T) {
	s, sel, q := fixture(t)
	bind(t, s, &q)
	r, e := s.Start(governance.LocalController(), q)
	if e != nil {
		t.Fatal(e)
	}
	coverage, e := s.Context(sel)
	if e != nil {
		t.Fatal(e)
	}
	project := coverage.Root
	state := filepath.Join(filepath.Dir(project), "state")
	before := map[string][]byte{}
	e = filepath.WalkDir(project, func(p string, d os.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if !d.IsDir() {
			b, e := os.ReadFile(p)
			if e != nil {
				return e
			}
			before[p] = b
		}
		return nil
	})
	if e != nil {
		t.Fatal(e)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	binary := filepath.Join(t.TempDir(), "sdp-controller")
	build := exec.CommandContext(ctx, filepath.Join(runtime.GOROOT(), "bin", "go"), "build", "-o", binary, "../cmd/sdp-controller")
	if b, e := build.CombinedOutput(); e != nil {
		t.Fatal(e, string(b))
	}
	host, e := os.Executable()
	if e != nil {
		t.Fatal(e)
	}
	message := filepath.Join(t.TempDir(), "message.txt")
	if e = os.WriteFile(message, []byte("Fix the selected fixture"), 0600); e != nil {
		t.Fatal(e)
	}
	invoke := func(action string, extra ...string) []byte {
		t.Helper()
		args := []string{action, "--project", project, "--state", state, "--run", r.Run.ID, "--codex", host}
		args = append(args, extra...)
		cmd := exec.CommandContext(ctx, binary, args...)
		cmd.Env = append(os.Environ(), "SDP_PGI3_FAKE_HOST=1")
		b, e := cmd.CombinedOutput()
		if e != nil {
			t.Fatal(action, e, string(b))
		}
		return b
	}
	out := invoke("run", "--message", message, "--mcp-binary", "/fixture/mcp")
	if !bytes.Contains(out, []byte("Recorded fixture outcome")) {
		t.Fatal(string(out))
	}
	r, e = s.Read(r.Run.ID)
	if e != nil {
		t.Fatal(e)
	}
	if r.Run.State != "running" || r.Run.Attempts[0].Status != "completed" || r.Run.Inputs[0].State != "acknowledged" {
		t.Fatal(r)
	}
	// An observer read does not change revision or start a host.
	invoke("status")
	out = invoke("provenance", "--codex", "/missing-codex-binary")
	var provenance governance.ProvenanceView
	if e = json.Unmarshal(out, &provenance); e != nil || provenance.Status != "observed" || provenance.Context == nil || provenance.Context.Skills.Status != "unavailable" || provenance.Context.Skills.UseEvidence != "not-observed" {
		t.Fatal(string(out), e)
	}
	after, e := s.Read(r.Run.ID)
	if e != nil || after.Run.Revision != r.Run.Revision {
		t.Fatal(after, e)
	}
	invoke("recover")
	after, e = s.Read(r.Run.ID)
	if e != nil || after.Run.ID != r.Run.ID || after.Run.Selection.Session != sel.Session || !after.Run.Attempts[0].Gap {
		t.Fatal(after, e)
	}
	if after.Run.Attempts[0].EventCount != r.Run.Attempts[0].EventCount+1 {
		t.Fatal("replayed final duplicated", r, after)
	}

	// Closing the output observer must not send SIGPIPE termination to the owner loop.
	s2, sel2, q2 := fixture(t)
	bind(t, s2, &q2)
	r2, e := s2.Start(governance.LocalController(), q2)
	if e != nil {
		t.Fatal(e)
	}
	cx2, e := s2.Context(sel2)
	if e != nil {
		t.Fatal(e)
	}
	rd, wr, e := os.Pipe()
	if e != nil {
		t.Fatal(e)
	}
	rd.Close()
	defer wr.Close()
	cmd := exec.CommandContext(ctx, binary, "run", "--project", cx2.Root, "--state", filepath.Join(filepath.Dir(cx2.Root), "state"), "--run", r2.Run.ID, "--codex", host, "--message", message, "--mcp-binary", "/fixture/mcp")
	cmd.Env = append(os.Environ(), "SDP_PGI3_FAKE_HOST=1")
	cmd.Stdout = wr
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if e = cmd.Run(); e != nil {
		t.Fatal("observer closure stopped controller", e, stderr.String())
	}
	r2, e = s2.Read(r2.Run.ID)
	if e != nil || r2.Run.Attempts[0].Status != "completed" {
		t.Fatal(r2, e)
	}
	for path, b := range before {
		current, e := os.ReadFile(path)
		if e != nil || !bytes.Equal(b, current) {
			t.Fatal("project fixture changed", path, e)
		}
	}
}

func TestCompiledApprovalInteraction(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	binary := filepath.Join(t.TempDir(), "sdp-controller")
	if b, e := exec.CommandContext(ctx, filepath.Join(runtime.GOROOT(), "bin", "go"), "build", "-o", binary, "../cmd/sdp-controller").CombinedOutput(); e != nil {
		t.Fatal(e, string(b))
	}
	host, _ := os.Executable()
	for _, decision := range []string{"accept", "decline"} {
		t.Run(decision, func(t *testing.T) {
			s, sel, q := fixture(t)
			bind(t, s, &q)
			r, e := s.Start(governance.LocalController(), q)
			if e != nil {
				t.Fatal(e)
			}
			coverage, _ := s.Context(sel)
			state := filepath.Join(filepath.Dir(coverage.Root), "state")
			message := filepath.Join(t.TempDir(), "owner.txt")
			os.WriteFile(message, []byte("Fix the selected fixture"), 0600)
			cmd := exec.CommandContext(ctx, binary, "run", "--project", coverage.Root, "--state", state, "--run", r.Run.ID, "--codex", host, "--message", message, "--mcp-binary", "/fixture/mcp", "--approvals", "ask")
			cmd.Env = append(os.Environ(), "SDP_PGI3_FAKE_HOST=1", "SDP_PGI3_APPROVAL="+coverage.Root)
			in, _ := cmd.StdinPipe()
			out, _ := cmd.StdoutPipe()
			var stderr bytes.Buffer
			cmd.Stderr = &stderr
			if e = cmd.Start(); e != nil {
				t.Fatal(e)
			}
			scanner := bufio.NewScanner(out)
			scanner.Buffer(make([]byte, 4096), 1<<20)
			answered := false
			for scanner.Scan() {
				var envelope struct {
					Approval *governance.ApprovalPrompt `json:"approval"`
				}
				json.Unmarshal(scanner.Bytes(), &envelope)
				if envelope.Approval == nil {
					continue
				}
				p := envelope.Approval
				// A separate observer can inspect the pending decision without taking a lease.
				status, e := exec.CommandContext(ctx, binary, "status", "--project", coverage.Root, "--state", state, "--run", r.Run.ID).CombinedOutput()
				if e != nil || !bytes.Contains(status, []byte("pendingApprovals")) {
					t.Fatal(string(status), e)
				}
				json.NewEncoder(in).Encode(map[string]any{"action": "approval", "approval": governance.ApprovalReply{RequestID: p.RequestID, Epoch: p.Epoch, ThreadID: p.ThreadID, TurnID: p.TurnID, Decision: decision}})
				answered = true
			}
			in.Close()
			if e = cmd.Wait(); e != nil {
				t.Fatal(e, stderr.String())
			}
			if !answered {
				t.Fatal("no observed prompt")
			}
			r, e = s.Read(r.Run.ID)
			if e != nil || len(r.Run.Inputs) != 2 || r.Run.Inputs[1].State != "acknowledged" || r.Run.State != "running" {
				t.Fatal(r, e)
			}
			input, e := s.OwnerInput(governance.LocalController(), r.Run.ID, r.Run.Inputs[1].ID)
			if e != nil || !bytes.Contains([]byte(input.Text), []byte(decision)) {
				t.Fatal(input, e)
			}
		})
	}
}
