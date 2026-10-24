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

func TestSessionDispositionCompiledJourney(t *testing.T) {
	root, state := os.Getenv("PG_DISPOSE_LAB_ROOT"), os.Getenv("PG_DISPOSE_LAB_STATE")
	if root != "" && state == "" {
		t.Fatal("explicit private state required")
	}
	binary := os.Getenv("PG_DISPOSE_BINARY")
	if binary == "" {
		binary = filepath.Join(t.TempDir(), "sdptool")
		if b, e := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "build", "-o", binary, "../cmd/sdptool").CombinedOutput(); e != nil {
			t.Fatal(e, string(b))
		}
	}
	s, q := dispositionFixture(t, root, state, multiCardImplementation)
	cases := []map[string]any{}
	dir := t.TempDir()
	call := func(op string, request any, def string, extra ...string) json.RawMessage {
		t.Helper()
		args := []string{s.root, "governance", op, "--state", s.state, "--json"}
		if request != nil {
			b, _ := json.Marshal(request)
			p := filepath.Join(dir, "request.json")
			if e := os.WriteFile(p, b, 0600); e != nil {
				t.Fatal(e)
			}
			args = append(args, "--input", p)
		}
		args = append(args, extra...)
		b, e := exec.Command(binary, args...).CombinedOutput()
		if e != nil || !json.Valid(b) {
			t.Fatal(op, e, string(b))
		}
		cases = append(cases, map[string]any{"definition": def, "value": json.RawMessage(b)})
		return b
	}
	raw := call("session-dispose", q, "mutation")
	var p SessionMutationResult
	json.Unmarshal(raw, &p)
	if p.RecordEffects != "prepared" {
		t.Fatal(string(raw))
	}
	if !bytes.Equal(raw, call("session-dispose", q, "mutation")) {
		t.Fatal("retry differs")
	}
	call("session-read", nil, "view", "--instance", q.InstanceID)
	file := filepath.Join(dir, "dispose.json")
	b, _ := json.Marshal(q)
	os.WriteFile(file, b, 0600)
	invoke, _ := json.Marshal(map[string]any{"schema": "sdptool.invoke/1", "requestId": "session-dispose-lab", "action": "governance.session.dispose", "context": map[string]string{"project": s.root}, "parameters": map[string]string{"state": s.state, "requestFile": file}})
	cmd := exec.Command(binary, "invoke", "--request", "-", "--as", "local-owner", "--authority", "controller")
	cmd.Stdin = bytes.NewReader(invoke)
	out, e := cmd.CombinedOutput()
	if e != nil {
		t.Fatal(e, string(out))
	}
	var result struct {
		Result SessionMutationResult `json:"result"`
	}
	if json.Unmarshal(out, &result) != nil || !reflect.DeepEqual(p, result.Result) {
		t.Fatal(string(out))
	}
	call("session-read", nil, "view", "--instance", q.InstanceID)
	rr := disposeResume(q, p)
	call("session-resume", rr, "mutation")
	raw = call("session-read", nil, "view", "--instance", q.InstanceID)
	var v SessionView
	json.Unmarshal(raw, &v)
	if v.Coverage != "current" || v.Steps[0].Stage != "done" || v.Steps[0].DispositionID != q.OperationID || len(v.Binding.Steps[0].Cards) != 2 {
		t.Fatal(v)
	}
	assertDisposed(t, s, q)
	call("session-list", nil, "list")
	cases = append(cases, map[string]any{"definition": "disposeRequest", "value": q}, map[string]any{"definition": "resumeRequest", "value": rr})
	captures, _ := json.MarshalIndent(cases, "", "  ")
	path := filepath.Join(dir, "wire.json")
	if root != "" {
		path = filepath.Join(s.root, "evidence/disposition-wire.json")
	}
	if e = os.WriteFile(path, captures, 0600); e != nil {
		t.Fatal(e)
	}
	if out, e = exec.Command("python3", "check_session_schema.py", path).CombinedOutput(); e != nil {
		t.Fatal(e, string(out))
	}
	t.Log(string(out))
	if _, e = s.publicationFiles(); e != nil {
		t.Fatal(e)
	}
}
