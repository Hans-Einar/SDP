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

func TestSessionOwnerCompiledJourney(t *testing.T) {
	binary := os.Getenv("PG_OWNER_BINARY")
	if binary == "" {
		binary = filepath.Join(t.TempDir(), "sdptool")
		if b, e := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "build", "-o", binary, "../cmd/sdptool").CombinedOutput(); e != nil {
			t.Fatal(e, string(b))
		}
	}
	root, state := os.Getenv("PG_OWNER_LAB_ROOT"), os.Getenv("PG_OWNER_LAB_STATE")
	if root != "" && state == "" {
		t.Fatal("explicit private state required")
	}
	s, q, d := ownerFixtureAt(t, root, state)
	dir := t.TempDir()
	cases := []map[string]any{}
	call := func(op string, request any, def string, extra ...string) json.RawMessage {
		t.Helper()
		args := []string{s.root, "governance", op, "--state", s.state, "--json"}
		if request != nil {
			b, _ := json.Marshal(request)
			p := filepath.Join(dir, "request.json")
			os.WriteFile(p, b, 0600)
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
	var prepared, done SessionMutationResult
	json.Unmarshal(call("session-owner", q, "mutation"), &prepared)
	r := SessionResumeRequest{SessionSchema, q.ProjectRef, q.CheckoutID, newTestID(), q.InstanceID, prepared.Revision, "Commit owner review request"}
	json.Unmarshal(call("session-resume", r, "mutation"), &done)
	call("session-read", nil, "view", "--instance", q.InstanceID)
	accept := ownerDecision(q, done.Revision, "accept")
	file := filepath.Join(dir, "owner.json")
	b, _ := json.Marshal(accept)
	os.WriteFile(file, b, 0600)
	request, _ := json.Marshal(map[string]any{"schema": "sdptool.invoke/1", "requestId": "owner-gate", "action": "governance.session.owner", "context": map[string]string{"project": s.root}, "parameters": map[string]string{"state": s.state, "requestFile": file}})
	invoke := func(authority string) ([]byte, error) {
		cmd := exec.Command(binary, "invoke", "--request", "-", "--as", "fixture-principal", "--authority", authority)
		cmd.Stdin = bytes.NewReader(request)
		return cmd.CombinedOutput()
	}
	head := mustFile(t, filepath.Join(s.state, "HEAD"))
	if b, e := invoke("controller"); e == nil || !bytes.Contains(b, []byte("unauthorized")) {
		t.Fatal("controller decision", e, string(b))
	}
	if !bytes.Equal(head, mustFile(t, filepath.Join(s.state, "HEAD"))) {
		t.Fatal("unauthorized mutation")
	}
	b, e := invoke("owner")
	if e != nil {
		t.Fatal(e, string(b))
	}
	var env struct {
		Result SessionMutationResult `json:"result"`
	}
	if json.Unmarshal(b, &env) != nil {
		t.Fatal(string(b))
	}
	json.Unmarshal(call("session-owner", accept, "mutation"), &prepared)
	if !reflect.DeepEqual(prepared, env.Result) {
		t.Fatal("owner CLI/JSON receipt differs")
	}
	r.OperationID = newTestID()
	r.ExpectedRevision = prepared.Revision
	json.Unmarshal(call("session-resume", r, "mutation"), &done)
	d.ExpectedRevision = done.Revision
	json.Unmarshal(call("session-dispose", d, "mutation"), &prepared)
	call("session-resume", disposeResume(d, prepared), "mutation")
	call("session-read", nil, "view", "--instance", q.InstanceID)
	assertDisposed(t, s, d)
	cases = append(cases, map[string]any{"definition": "ownerRequest", "value": q}, map[string]any{"definition": "ownerRequest", "value": accept})
	raw, _ := json.MarshalIndent(cases, "", "  ")
	wire := filepath.Join(dir, "wire.json")
	if root != "" {
		wire = filepath.Join(root, "evidence/owner-wire.json")
	}
	os.WriteFile(wire, raw, 0600)
	if b, e := exec.Command("python3", "check_session_schema.py", wire).CombinedOutput(); e != nil {
		t.Fatal(e, string(b))
	} else {
		t.Log(string(b))
	}
}
