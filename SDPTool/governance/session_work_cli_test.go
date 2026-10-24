package governance

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func TestSessionWorkCompiledJourney(t *testing.T) {
	root, state := os.Getenv("PG_WORK_LAB_ROOT"), os.Getenv("PG_WORK_LAB_STATE")
	if root != "" && state == "" {
		t.Fatal("explicit lab state required")
	}
	binary := os.Getenv("PG_WORK_BINARY")
	if binary == "" {
		binary = filepath.Join(t.TempDir(), "sdptool")
		if b, e := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "build", "-o", binary, "../cmd/sdptool").CombinedOutput(); e != nil {
			t.Fatal(e, string(b))
		}
	}
	host := filepath.Join(t.TempDir(), "session-controller")
	if b, e := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "build", "-o", host, "./testdata/session_controller").CombinedOutput(); e != nil {
		t.Fatal(e, string(b))
	}
	s, q := workFixture(t, root, state)
	cases := []map[string]any{}
	reqdir := t.TempDir()
	call := func(op string, request any, def string, extra ...string) json.RawMessage {
		t.Helper()
		args := []string{s.root, "governance", op, "--state", s.state, "--json"}
		if request != nil {
			raw, _ := json.Marshal(request)
			path := filepath.Join(reqdir, "request.json")
			if e := os.WriteFile(path, raw, 0600); e != nil {
				t.Fatal(e)
			}
			args = append(args, "--input", path)
		}
		args = append(args, extra...)
		b, e := exec.Command(binary, args...).CombinedOutput()
		if e != nil || !json.Valid(b) {
			t.Fatal(op, e, string(b))
		}
		cases = append(cases, map[string]any{"definition": def, "value": json.RawMessage(b)})
		return b
	}
	raw := call("session-start", q, "mutation")
	var prepared SessionMutationResult
	if json.Unmarshal(raw, &prepared) != nil || prepared.RecordEffects != "prepared" {
		t.Fatal("missing prepared start", string(raw))
	}
	call("session-read", nil, "view", "--instance", q.InstanceID)
	// Invocation retries share the same controller-attributed operation receipt.
	file := filepath.Join(reqdir, "start.json")
	request, _ := json.Marshal(q)
	os.WriteFile(file, request, 0600)
	inv, _ := json.Marshal(map[string]any{"schema": "sdptool.invoke/1", "requestId": "session-start-trial", "action": "governance.session.start", "context": map[string]string{"project": s.root}, "parameters": map[string]string{"requestFile": file, "state": s.state}})
	cmd := exec.Command(binary, "invoke", "--request", "-", "--as", "local-owner", "--authority", "controller")
	cmd.Stdin = bytes.NewReader(inv)
	output, e := cmd.CombinedOutput()
	if e != nil {
		t.Fatal(e, string(output))
	}
	var result struct {
		Result SessionMutationResult `json:"result"`
	}
	if json.Unmarshal(output, &result) != nil || !reflect.DeepEqual(result.Result, prepared) {
		t.Fatal("invoke receipt differs", string(output))
	}
	resume := workResume(q, prepared)
	raw = call("session-resume", resume, "mutation")
	var done SessionMutationResult
	json.Unmarshal(raw, &done)
	if !bytes.Equal(raw, call("session-resume", resume, "mutation")) {
		t.Fatal("retry differs")
	}
	call("session-read", nil, "view", "--instance", q.InstanceID)
	// Real Controller.Launch, existing skill inputs and report capability with a
	// deterministic protocol host; explicitly not a native Codex/model trial.
	output, e = exec.Command(host, s.root, s.state, q.OperationID, "success").CombinedOutput()
	if e != nil {
		t.Fatal(e, string(output))
	}
	var reported Result
	if json.Unmarshal(output, &reported) != nil || reported.Run.State != "verifying" {
		t.Fatal(string(output))
	}
	call("session-read", nil, "view", "--instance", q.InstanceID)
	pq := progressRequest(t, s, q, done, reported)
	raw = call("session-progress", pq, "mutation")
	json.Unmarshal(raw, &prepared)
	call("session-read", nil, "view", "--instance", q.InstanceID)
	pr := workResume(pq, prepared)
	call("session-resume", pr, "mutation")
	raw = call("session-read", nil, "view", "--instance", q.InstanceID)
	var view SessionView
	json.Unmarshal(raw, &view)
	if view.Steps[0].Stage != "verification" || view.HostExecution != "completed" || view.Steps[0].RunID != q.OperationID || view.Binding == nil {
		t.Fatal(view)
	}
	call("session-list", nil, "list")
	cases = append(cases, map[string]any{"definition": "startRequest", "value": q}, map[string]any{"definition": "progressRequest", "value": pq}, map[string]any{"definition": "resumeRequest", "value": pr})
	captures, _ := json.MarshalIndent(cases, "", "  ")
	path := filepath.Join(reqdir, "wire.json")
	if root != "" {
		path = filepath.Join(s.root, "evidence/work-wire.json")
		os.WriteFile(filepath.Join(s.root, "evidence/controller-report.json"), output, 0600)
	}
	if e = os.WriteFile(path, captures, 0600); e != nil {
		t.Fatal(e)
	}
	if b, e := exec.Command("python3", "check_session_schema.py", path).CombinedOutput(); e != nil {
		t.Fatal(e, string(b))
	} else {
		t.Log(string(b))
	}
	// A separate Session verifies uncertain host send is never repeated.
	lost, lq := workFixture(t, "", "")
	startWork(t, lost, lq)
	output, e = exec.Command(host, lost.root, lost.state, lq.OperationID, "lost").CombinedOutput()
	if e != nil {
		t.Fatal(e, string(output))
	}
	var unknown SessionView
	if json.Unmarshal(output, &unknown) != nil || unknown.HostExecution != "unknown" {
		t.Fatal(string(output))
	}
	if _, e = s.publicationFiles(); e != nil {
		t.Fatal(e)
	}
}
