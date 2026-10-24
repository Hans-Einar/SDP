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
	"strings"
	"testing"
	"time"

	"github.com/Hans-Einar/SDP/SDPTool/governancehost"
)

// Explicit protocol fixture only; never execute the installed Codex/model.
func leaderProtocolFixture() {
	dir := os.Getenv("SDP_PGL3_C7_HOST")
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Println("codex-cli 0.165.0")
		return
	}
	root := os.Getenv("SDP_PGL3_C7_ROOT")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 4096), 1<<20)
	enc := json.NewEncoder(os.Stdout)
	for scanner.Scan() {
		var q struct {
			ID     any             `json:"id"`
			Method string          `json:"method"`
			Params json.RawMessage `json:"params"`
		}
		if json.Unmarshal(scanner.Bytes(), &q) != nil {
			return
		}
		if q.ID == nil {
			continue
		}
		log, e := os.OpenFile(filepath.Join(dir, "calls"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
		if e != nil {
			return
		}
		fmt.Fprintln(log, q.Method)
		log.Close()
		var result any
		switch q.Method {
		case "initialize":
			result = map[string]any{"userAgent": "sdp_governance/0.165.0 fixture", "platformOs": "linux"}
		case "account/read":
			result = map[string]any{"account": map[string]any{"type": "chatgpt"}}
		case "skills/list":
			result = map[string]any{"data": []any{map[string]any{"cwd": root, "errors": []any{}, "skills": []any{map[string]any{"name": "sdp", "path": filepath.Join(root, "Skills/sdp/SKILL.md"), "enabled": true}}}}}
		case "thread/start":
			var params struct {
				Cwd, Sandbox   string
				ApprovalPolicy string `json:"approvalPolicy"`
				Config         struct {
					MCP map[string]struct {
						Args []string `json:"args"`
					} `json:"mcp_servers"`
				} `json:"config"`
			}
			if json.Unmarshal(q.Params, &params) != nil || params.Cwd != root || params.Sandbox != "read-only" || params.ApprovalPolicy != "never" || !strings.Contains(strings.Join(params.Config.MCP["sdp_governance"].Args, " "), "--coordinator fixture-leader") {
				return
			}
			calls, _ := os.ReadFile(filepath.Join(dir, "calls"))
			threadID := fmt.Sprintf("leader-thread-%d", bytes.Count(calls, []byte("thread/start\n")))
			os.WriteFile(filepath.Join(dir, "thread"), []byte(threadID), 0600)
			result = map[string]any{"thread": map[string]any{"id": threadID, "cwd": root, "cliVersion": "0.165.0", "status": map[string]any{"type": "idle"}}}
		case "mcpServerStatus/list":
			result = map[string]any{"data": []any{map[string]any{"name": "sdp_governance", "authStatus": "unsupported", "tools": map[string]any{"coordination_read": map[string]any{}, "report_acknowledge": map[string]any{}}}}}
		case "turn/start":
			var p struct {
				ThreadID string          `json:"threadId"`
				ClientID string          `json:"clientUserMessageId"`
				Input    json.RawMessage `json:"input"`
			}
			if json.Unmarshal(q.Params, &p) != nil || !strings.HasPrefix(p.ThreadID, "leader-thread-") {
				return
			}
			user := map[string]any{"id": "input", "type": "userMessage", "clientId": p.ClientID, "content": p.Input}
			raw, _ := json.Marshal(user)
			if os.WriteFile(filepath.Join(dir, "input.json"), raw, 0600) != nil {
				return
			}
			if os.Getenv("SDP_PGL3_C7_LOST") == "1" {
				enc.Encode(map[string]any{"id": q.ID, "error": map[string]any{"code": -32000, "message": "recorded lost turn response"}})
				continue
			}
			result = map[string]any{"turn": map[string]any{"id": "leader-turn", "status": "inProgress"}}
		case "thread/read":
			raw, e := os.ReadFile(filepath.Join(dir, "input.json"))
			if e != nil {
				return
			}
			threadID, _ := os.ReadFile(filepath.Join(dir, "thread"))
			result = map[string]any{"thread": map[string]any{"id": string(threadID), "cwd": root, "cliVersion": "0.165.0", "status": map[string]any{"type": "idle"}, "turns": []any{map[string]any{"id": "leader-turn", "status": "completed", "itemsView": "full", "items": []any{json.RawMessage(raw), map[string]any{"id": "answer", "type": "agentMessage", "phase": "final_answer", "text": "Recorded Project Leader proposal; owner decision remains."}}}}}}
		default:
			enc.Encode(map[string]any{"id": q.ID, "error": map[string]any{"code": -32601}})
			continue
		}
		enc.Encode(map[string]any{"id": q.ID, "result": result})
	}
}
func TestCompiledLeaderQueueAndLostResponseRestart(t *testing.T) {
	s, sel, _ := fixture(t)
	coverage, e := s.Context(sel)
	if e != nil {
		t.Fatal(e)
	}
	root := coverage.Root
	state := filepath.Join(filepath.Dir(root), "state")
	path := filepath.Join(root, "Skills/sdp/SKILL.md")
	os.MkdirAll(filepath.Dir(path), 0700)
	if e = os.WriteFile(path, []byte("bounded fixture skill"), 0600); e != nil {
		t.Fatal(e)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	binary := filepath.Join(t.TempDir(), "sdp-controller")
	if b, e := exec.CommandContext(ctx, filepath.Join(runtime.GOROOT(), "bin", "go"), "build", "-o", binary, "../cmd/sdp-controller").CombinedOutput(); e != nil {
		t.Fatal(e, string(b))
	}
	host, _ := os.Executable()
	hostState := t.TempDir()
	invoke := func(action string, lost bool, stdin string, extra ...string) ([]byte, error) {
		args := []string{"leader", action, "--project", root, "--state", state}
		if action == "console" {
			args = []string{"coordinate", "--project", root, "--state", state, "--watch", "--json", "--owner", "owner", "--interval", "100ms", "--timeout", "2s"}
		}
		args = append(args, extra...)
		cmd := exec.CommandContext(ctx, binary, args...)
		cmd.Stdin = strings.NewReader(stdin)
		cmd.Env = append(os.Environ(), "SDP_PGL3_C7_HOST="+hostState, "SDP_PGL3_C7_ROOT="+root)
		if lost {
			cmd.Env = append(cmd.Env, "SDP_PGL3_C7_LOST=1")
		}
		return cmd.CombinedOutput()
	}
	capture := func(name string, b []byte) {
		t.Helper()
		if dir := os.Getenv("SDP_PGL3_C7_CAPTURE"); dir != "" {
			if e := os.MkdirAll(dir, 0700); e != nil {
				t.Fatal(e)
			}
			if e := os.WriteFile(filepath.Join(dir, name), b, 0600); e != nil {
				t.Fatal(e)
			}
		}
	}
	inspect, e := invoke("inspect", false, "")
	if e != nil {
		t.Fatal(string(inspect), e)
	}
	if _, e = os.Stat(state); !os.IsNotExist(e) {
		t.Fatal("inspection created state")
	}
	capture("empty.json", inspect)
	message := filepath.Join(t.TempDir(), "message.txt")
	os.WriteFile(message, []byte("Review current registered work"), 0600)
	id := governancehost.UUID()
	b, e := invoke("message", false, "", "--owner", "owner", "--operation", id, "--input", message)
	if e != nil {
		t.Fatal(string(b), e)
	}
	capture("message.json", b)
	retry, e := invoke("message", false, "", "--owner", "owner", "--operation", id, "--input", message)
	if e != nil || !bytes.Equal(b, retry) {
		t.Fatal("enqueue retry changed", string(retry), e)
	}
	flags := []string{"--coordinator", "fixture-leader", "--mcp-binary", "/fixture/mcp", "--codex", host, "--skill", "Skills/sdp/SKILL.md", "--interval", "100ms", "--timeout", "10s"}
	lost, e := invoke("run", true, "", flags...)
	if e == nil || !bytes.Contains(lost, []byte("host_rpc_error: turn/start")) {
		t.Fatal("fixture did not lose reply", string(lost), e)
	}
	v, e := s.LeaderRead()
	if e != nil || v.Leader.Dispatches[0].State != "sending" {
		t.Fatal(v, e)
	}
	// Owner input remains usable while the first send is unresolved.
	cmd := fmt.Sprintf("{\"action\":\"message\",\"message\":{\"operationId\":%q,\"text\":\"Owner correction queued during unresolved send\"}}\n{\"action\":\"quit\"}\n", governancehost.UUID())
	console, e := invoke("console", false, cmd)
	if e != nil || !bytes.Contains(console, []byte("leader-message")) {
		t.Fatal(string(console), e)
	}
	capture("console.ndjson", console)
	recovered, e := invoke("run", false, "", flags...)
	if e != nil {
		t.Fatal(string(recovered), e)
	}
	capture("recovered.ndjson", recovered)
	v, e = s.LeaderRead()
	if e != nil || len(v.Leader.Dispatches) != 1 || v.Leader.Dispatches[0].State != "completed" || len(v.Leader.Messages) != 2 {
		t.Fatal(v, e)
	}
	calls, _ := os.ReadFile(filepath.Join(hostState, "calls"))
	if bytes.Count(calls, []byte("turn/start\n")) != 1 || bytes.Count(calls, []byte("thread/start\n")) != 1 {
		t.Fatal("resubmitted after restart", string(calls))
	}
	capture("host-methods.log", calls)
	inspected, e := invoke("inspect", false, "")
	if e != nil {
		t.Fatal(e)
	}
	capture("completed.json", inspected)
	if old := os.Getenv("SDP_PGL3_C7_OLD_CONTROLLER"); old != "" {
		output, err := exec.CommandContext(ctx, old, "overview", "--project", root, "--state", state, "--json").CombinedOutput()
		if err == nil || (!bytes.Contains(output, []byte("reconciliation_required")) && !bytes.Contains(output, []byte("incompatible_profile"))) {
			t.Fatal("old reader did not refuse", string(output), err)
		}
		capture("old-reader.log", output)
	}

	// The original message has been delivered; the next turn selects only new input.
	next, e := invoke("run", false, "", flags...)
	if e != nil {
		t.Fatal(string(next), e)
	}
	capture("next.ndjson", next)
	v, e = s.LeaderRead()
	if e != nil || len(v.Leader.Dispatches) != 2 || len(v.Leader.Dispatches[1].MessageIDs) != 1 || v.Leader.Dispatches[1].MessageIDs[0] != v.Leader.Messages[1].ID {
		t.Fatal(v, e)
	}
	// Simulated history is isolated from canonical process documents and work runs.
	overview, e := s.Overview()
	if e != nil {
		t.Fatal(e)
	}
	if len(overview.Runs) != 0 {
		t.Fatal("leader created worker assignments")
	}
}
