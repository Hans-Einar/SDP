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

func TestSessionAdvanceCompiledJourney(t *testing.T) {
	root, state := os.Getenv("PG_ADVANCE_LAB_ROOT"), os.Getenv("PG_ADVANCE_LAB_STATE")
	if root != "" && state == "" {
		t.Fatal("explicit private state required")
	}
	binary := os.Getenv("PG_ADVANCE_BINARY")
	if binary == "" {
		binary = filepath.Join(t.TempDir(), "sdptool")
		if b, e := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "build", "-o", binary, "../cmd/sdptool").CombinedOutput(); e != nil {
			t.Fatal(e, string(b))
		}
	}
	s, d := dispositionFixture(t, root, state, nil)
	completeFixtureDisposition(t, s, d)
	q := advanceRequest(t, s, d.InstanceID, "complete", "", "disposition", d.OperationID)
	cases := []map[string]any{{"definition": "advanceRequest", "value": q}}
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
	var p SessionMutationResult
	raw := call("session-advance", q, "mutation")
	json.Unmarshal(raw, &p)
	if !bytes.Equal(raw, call("session-advance", q, "mutation")) {
		t.Fatal("retry differs")
	}
	file := filepath.Join(dir, "advance.json")
	b, _ := json.Marshal(q)
	os.WriteFile(file, b, 0600)
	invoke, _ := json.Marshal(map[string]any{"schema": "sdptool.invoke/1", "requestId": "advance", "action": "governance.session.advance", "context": map[string]string{"project": s.root}, "parameters": map[string]string{"state": s.state, "requestFile": file}})
	cmd := exec.Command(binary, "invoke", "--request", "-", "--as", "fixture-controller", "--authority", "controller")
	cmd.Stdin = bytes.NewReader(invoke)
	b, e := cmd.CombinedOutput()
	if e != nil {
		t.Fatal(e, string(b))
	}
	var env struct {
		Result SessionMutationResult `json:"result"`
	}
	json.Unmarshal(b, &env)
	if !reflect.DeepEqual(p, env.Result) {
		t.Fatal("CLI/JSON mismatch")
	}
	r := SessionResumeRequest{SessionSchema, q.ProjectRef, q.CheckoutID, newTestID(), q.InstanceID, p.Revision, "Commit exact plan and Session closure"}
	first := call("session-resume", r, "mutation")
	if !bytes.Equal(first, call("session-resume", r, "mutation")) {
		t.Fatal("resume replay differs")
	}
	var v SessionView
	json.Unmarshal(call("session-read", nil, "view", "--instance", q.InstanceID), &v)
	if v.State != "completed" || v.Coverage != "current" {
		t.Fatal(v)
	}
	if _, e = s.publicationFiles(); e != nil {
		t.Fatal(e)
	}
	raw, _ = json.MarshalIndent(cases, "", "  ")
	wire := filepath.Join(dir, "wire.json")
	if root != "" {
		wire = filepath.Join(root, "evidence/advance-wire.json")
	}
	if e = os.WriteFile(wire, raw, 0600); e != nil {
		t.Fatal(e)
	}
	if b, e = exec.Command("python3", "check_session_schema.py", wire).CombinedOutput(); e != nil {
		t.Fatal(e, string(b))
	} else {
		t.Log(string(b))
	}
}
