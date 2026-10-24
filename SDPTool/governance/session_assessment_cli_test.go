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

// Real CLI/invoke processes share one store. Lab export requires an absent root;
// prior lab fixtures are preserved by the common fixture constructor.
func TestSessionAssessmentCompiledJourney(t *testing.T) {
	root, state := os.Getenv("PG_ASSESSMENT_LAB_ROOT"), os.Getenv("PG_ASSESSMENT_LAB_STATE")
	if root != "" && state == "" {
		t.Fatal("explicit lab state required")
	}
	binary := os.Getenv("PG_ASSESSMENT_BINARY")
	if binary == "" {
		binary = filepath.Join(t.TempDir(), "sdptool")
		if b, e := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "build", "-o", binary, "../cmd/sdptool").CombinedOutput(); e != nil {
			t.Fatal(e, string(b))
		}
	}
	s, q := assessmentFixture(t, root, state)
	q.Criteria[0].Impact = "unknown"
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
	raw := call("session-assess", q, "mutation")
	var prepared SessionMutationResult
	json.Unmarshal(raw, &prepared)
	if prepared.RecordEffects != "prepared" {
		t.Fatal(prepared)
	}
	call("session-read", nil, "view", "--instance", q.InstanceID)
	file := filepath.Join(reqdir, "assessment.json")
	request, _ := json.Marshal(q)
	os.WriteFile(file, request, 0600)
	inv, _ := json.Marshal(map[string]any{"schema": "sdptool.invoke/1", "requestId": "assessment-trial", "action": "governance.session.assess", "context": map[string]string{"project": s.root}, "parameters": map[string]string{"requestFile": file, "state": s.state}})
	cmd := exec.Command(binary, "invoke", "--request", "-", "--as", "local-owner", "--authority", "controller")
	cmd.Stdin = bytes.NewReader(inv)
	b, e := cmd.CombinedOutput()
	if e != nil {
		t.Fatal(e, string(b))
	}
	var result struct {
		Result SessionMutationResult `json:"result"`
	}
	if json.Unmarshal(b, &result) != nil || !reflect.DeepEqual(result.Result, prepared) {
		t.Fatal("invoke differs", string(b))
	}
	resume := assessResume(q, prepared)
	raw = call("session-resume", resume, "mutation")
	var done SessionMutationResult
	json.Unmarshal(raw, &done)
	if !bytes.Equal(raw, call("session-resume", resume, "mutation")) {
		t.Fatal("retry differs")
	}
	raw = call("session-read", nil, "view", "--instance", q.InstanceID)
	var view SessionView
	json.Unmarshal(raw, &view)
	if view.Steps[0].Stage != "admission" || view.Steps[0].Assessment.DesignImpact != "unknown" {
		t.Fatal(view)
	}
	cases = append(cases, map[string]any{"definition": "assessRequest", "value": q}, map[string]any{"definition": "resumeRequest", "value": resume})
	q.OperationID = newTestID()
	q.ExpectedRevision = done.Revision
	q.Criteria = append([]SessionDesignCriterion{}, q.Criteria...)
	q.Criteria[0].Impact = "unchanged"
	q.Reason = "Resolved fixture uncertainty; explanatory wording only"
	raw = call("session-assess", q, "mutation")
	json.Unmarshal(raw, &prepared)
	resume = assessResume(q, prepared)
	call("session-resume", resume, "mutation")
	raw = call("session-read", nil, "view", "--instance", q.InstanceID)
	json.Unmarshal(raw, &view)
	if view.Coverage != "current" || view.RecordEffects != "committed" || view.Steps[0].Stage != "implementation" || view.HostExecution != "not-dispatched" || len(view.Steps[0].Holds) == 0 {
		t.Fatal(view)
	}
	call("session-list", nil, "list")
	cases = append(cases, map[string]any{"definition": "assessRequest", "value": q})
	capturePath := filepath.Join(reqdir, "wire.json")
	if root != "" {
		capturePath = filepath.Join(s.root, "evidence/assessment-wire.json")
	}
	captures, _ := json.MarshalIndent(cases, "", "  ")
	if e = os.WriteFile(capturePath, captures, 0600); e != nil {
		t.Fatal(e)
	}
	if b, e := exec.Command("python3", "check_session_schema.py", capturePath).CombinedOutput(); e != nil {
		t.Fatal(e, string(b))
	} else {
		t.Log(string(b))
	}
	if _, e = s.publicationFiles(); e != nil {
		t.Fatal(e)
	}
}
