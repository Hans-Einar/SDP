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

func TestSessionSelectionCompiledJourney(t *testing.T) {
	root, state := os.Getenv("PG_SELECTION_LAB_ROOT"), os.Getenv("PG_SELECTION_LAB_STATE")
	if root != "" && state == "" {
		t.Fatal("explicit lab state required")
	}
	binary := os.Getenv("PG_SELECTION_BINARY")
	if binary == "" {
		binary = filepath.Join(t.TempDir(), "sdptool")
		if b, e := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "build", "-o", binary, "../cmd/sdptool").CombinedOutput(); e != nil {
			t.Fatal(e, string(b))
		}
	}
	s, q := selectionFixture(t, root, state)
	reqdir := t.TempDir()
	cases := []map[string]any{}
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
		if e != nil {
			t.Fatal(e, string(b))
		}
		if !json.Valid(b) {
			t.Fatal(string(b))
		}
		cases = append(cases, map[string]any{"definition": def, "value": json.RawMessage(b)})
		return b
	}
	raw := call("session-select", q, "mutation")
	var selected SessionMutationResult
	if e := json.Unmarshal(raw, &selected); e != nil || selected.RecordEffects != "prepared" {
		t.Fatal(e, string(raw))
	}
	call("session-read", nil, "view", "--instance", q.InstanceID)
	// Exact invocation retry uses trusted local attribution and shares the receipt.
	file := filepath.Join(reqdir, "select.json")
	request, _ := json.Marshal(q)
	os.WriteFile(file, request, 0600)
	inv, _ := json.Marshal(map[string]any{"schema": "sdptool.invoke/1", "requestId": "selection-trial", "action": "governance.session.select", "context": map[string]string{"project": s.root}, "parameters": map[string]string{"requestFile": file, "state": s.state}})
	cmd := exec.Command(binary, "invoke", "--request", "-", "--as", "local-owner", "--authority", "controller")
	cmd.Stdin = bytes.NewReader(inv)
	output, e := cmd.CombinedOutput()
	if e != nil {
		t.Fatal(e, string(output))
	}
	var result struct {
		Result SessionMutationResult `json:"result"`
	}
	if json.Unmarshal(output, &result) != nil || !reflect.DeepEqual(result.Result, selected) {
		t.Fatal("invocation retry differs", string(output))
	}
	resume := selectionResume(q, selected)
	raw = call("session-resume", resume, "mutation")
	var done SessionMutationResult
	json.Unmarshal(raw, &done)
	if done.RecordEffects != "committed" {
		t.Fatal(string(raw))
	}
	if !bytes.Equal(raw, call("session-resume", resume, "mutation")) {
		t.Fatal("resume retry differs")
	}
	raw = call("session-read", nil, "view", "--instance", q.InstanceID)
	var view SessionView
	json.Unmarshal(raw, &view)
	if view.Coverage != "current" || !view.Steps[0].Selected || view.Steps[0].Stage != "admission" || view.HostExecution != "not-dispatched" {
		t.Fatal(view)
	}
	call("session-list", nil, "list")
	if _, e = s.publicationFiles(); e != nil {
		t.Fatal(e)
	}
	card := mustFile(t, filepath.Join(s.root, "SDP/KanBan/active/#001--Change--fixture.md"))
	meta, _ := metadata(card)
	if meta["CardState"] != "ready" {
		t.Fatal(meta)
	}
	cases = append(cases, map[string]any{"definition": "selectRequest", "value": q}, map[string]any{"definition": "resumeRequest", "value": resume})
	captures, _ := json.MarshalIndent(cases, "", "  ")
	capturePath := filepath.Join(reqdir, "wire.json")
	if root != "" {
		capturePath = filepath.Join(s.root, "evidence/selection-wire.json")
	}
	if e = os.WriteFile(capturePath, captures, 0600); e != nil {
		t.Fatal(e)
	}
	if b, e := exec.Command("python3", "check_session_schema.py", capturePath).CombinedOutput(); e != nil {
		t.Fatal(e, string(b))
	} else {
		t.Log(string(b))
	}
}
