package governance

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

// Default runs use disposable directories. Explicit PG_SESSION_LAB_ROOT exports
// the same actual-process journey into an absent, owner-authorized lab fixture.
func TestSessionCompiledJourney(t *testing.T) {
	root, state := os.Getenv("PG_SESSION_LAB_ROOT"), os.Getenv("PG_SESSION_LAB_STATE")
	if root != "" && state == "" {
		t.Fatal("explicit lab state required")
	}
	binary := os.Getenv("PG_SESSION_BINARY")
	if binary == "" {
		binary = filepath.Join(t.TempDir(), "sdptool")
		if b, e := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "build", "-o", binary, "../cmd/sdptool").CombinedOutput(); e != nil {
			t.Fatal(e, string(b))
		}
	}
	s, c, q := sessionFixtureAt(t, root, state)
	original, _ := os.ReadFile(filepath.Join(s.root, c.Binding.Session.Path))
	reqdir := t.TempDir()
	cases := []map[string]any{}
	call := func(op string, req any, def string, extra ...string) json.RawMessage {
		t.Helper()
		args := []string{s.root, "governance", op, "--state", s.state, "--json"}
		if req != nil {
			data, _ := json.Marshal(req)
			file := filepath.Join(reqdir, "request.json")
			if e := os.WriteFile(file, data, 0600); e != nil {
				t.Fatal(e)
			}
			args = append(args, "--input", file)
		}
		args = append(args, extra...)
		cmd := exec.Command(binary, args...)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		b, e := cmd.Output()
		if e != nil {
			t.Fatal(op, e, stderr.String())
		}
		if !json.Valid(b) {
			t.Fatal("not JSON", string(b))
		}
		cases = append(cases, map[string]any{"definition": def, "value": json.RawMessage(b)})
		return b
	}
	call("session-definition", nil, "definition")
	raw := call("session-context", SessionContextRequest{q.BindingPath}, "context")
	var ctx SessionContext
	if e := json.Unmarshal(raw, &ctx); e != nil || !reflect.DeepEqual(c, ctx) {
		t.Fatal("CLI context differs", e)
	}
	raw = call("session-adopt", q, "mutation")
	var adopted SessionMutationResult
	json.Unmarshal(raw, &adopted)
	if adopted.Registration != "pending" {
		t.Fatal(adopted)
	}
	if got := call("session-adopt", q, "mutation"); !bytes.Equal(raw, got) {
		t.Fatal("adopt retry differs")
	}
	call("session-read", nil, "view", "--instance", adopted.InstanceID)
	resume := sessionResume(adopted, q)
	raw = call("session-resume", resume, "mutation")
	if got := call("session-resume", resume, "mutation"); !bytes.Equal(raw, got) {
		t.Fatal("resume retry differs")
	}
	viewRaw := call("session-read", nil, "view", "--instance", adopted.InstanceID)
	var view SessionView
	json.Unmarshal(viewRaw, &view)
	if view.Registration != "committed" || view.Revision != 2 || view.Coverage != "current" || view.HostExecution != "not-dispatched" {
		t.Fatal(view)
	}
	call("session-list", nil, "list")
	// Both adapters observe exactly the same instance/store.
	invoke, _ := json.Marshal(map[string]any{"schema": "sdptool.invoke/1", "requestId": "session-trial", "action": "governance.session.read", "context": map[string]string{"project": s.root}, "parameters": map[string]string{"state": s.state, "instance": adopted.InstanceID}})
	cmd := exec.Command(binary, "invoke", "--request", "-")
	cmd.Stdin = bytes.NewReader(invoke)
	b, e := cmd.CombinedOutput()
	if e != nil {
		t.Fatal(e, string(b))
	}
	var result struct {
		Result SessionView `json:"result"`
	}
	if json.Unmarshal(b, &result) != nil || !reflect.DeepEqual(view, result.Result) {
		t.Fatal("invocation/CLI differ", string(b))
	}
	final, _ := os.ReadFile(filepath.Join(s.root, c.Binding.Session.Path))
	if !bytes.HasPrefix(final, original) || strings.Count(string(final), "## PG Session registration ") != 1 {
		t.Fatal("journal changed/duplicated")
	}
	cases = append(cases, map[string]any{"definition": "binding", "value": c.Binding}, map[string]any{"definition": "adoptRequest", "value": q}, map[string]any{"definition": "resumeRequest", "value": resume}, map[string]any{"definition": "contextRequest", "value": SessionContextRequest{q.BindingPath}})
	captures, _ := json.MarshalIndent(cases, "", "  ")
	captureFile := filepath.Join(reqdir, "captures.json")
	if root != "" {
		captureFile = filepath.Join(s.root, "evidence/session-wire.json")
	}
	if e := os.WriteFile(captureFile, captures, 0600); e != nil {
		t.Fatal(e)
	}
	// A strict wire schema supports consumers independently of the Go structs.
	check := exec.Command("python3", "check_session_schema.py", captureFile)
	if b, e := check.CombinedOutput(); e != nil {
		t.Fatal(e, string(b))
	} else {
		t.Log(string(b))
	}
}
