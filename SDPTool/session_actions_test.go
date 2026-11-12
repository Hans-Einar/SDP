package sdptool

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Hans-Einar/SDP/SDPTool/governanceapp"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/Hans-Einar/SDP/SDPTool/blueprintstate"
)

func TestSessionActionsReadAndAuthority(t *testing.T) {
	root, _ := projectFixture(t)
	state := filepath.Join(t.TempDir(), "state")
	for _, op := range []string{"definition", "list"} {
		params := map[string]any{}
		if op == "list" {
			params["state"] = state
		}
		r := Invoke(context.Background(), bytes.NewReader(requestBytes(t, root, "governance.session."+op, params)), ActionCaller{})
		var out, errs bytes.Buffer
		code := governanceCommand(root, []string{"session-" + op, "--state", state}, &out, &errs, true)
		var expected, actual any
		resultJSON, marshalErr := json.Marshal(r.Result)
		decodeErr := json.Unmarshal(resultJSON, &actual)
		if code != 0 || json.Unmarshal(out.Bytes(), &expected) != nil || r.ExitCode != 0 || marshalErr != nil || decodeErr != nil || !reflect.DeepEqual(expected, actual) {
			t.Fatalf("%s: %+v %s", op, r, errs.String())
		}
	}
	for _, op := range []string{"adopt", "resume", "select", "assess", "start", "progress", "verify", "review", "dispose", "owner", "advance"} {
		raw := requestBytes(t, root, "governance.session."+op, SessionActionInput{State: state, RequestFile: "absent.json"})
		for _, caller := range []ActionCaller{{}, {Principal: blueprintstate.Principal{Actor: "agent", Authority: "assignee"}}, {Principal: blueprintstate.Principal{Authority: "controller"}}} {
			r := Invoke(context.Background(), bytes.NewReader(raw), caller)
			if r.Error == nil || r.Error.Code != "unauthorized" {
				t.Fatalf("unexpected authority: %+v", r.Error)
			}
		}
		r := Invoke(context.Background(), bytes.NewReader(raw), ActionCaller{Principal: blueprintstate.Principal{Actor: "local-owner", Authority: "controller"}})
		if r.Error == nil || r.Error.Code == "unauthorized" {
			t.Fatalf("controller did not reach domain request read: %+v", r)
		}
	}
	if _, e := os.Stat(state); !os.IsNotExist(e) {
		t.Fatal("read or refused operation wrote state")
	}
}

func TestSessionActionCLIErrorParity(t *testing.T) {
	root, _ := projectFixture(t)
	state := filepath.Join(t.TempDir(), "state")
	request := filepath.Join(root, "request.json")
	if err := os.WriteFile(request, []byte(`{"unexpected":true}`), 0600); err != nil {
		t.Fatal(err)
	}
	for _, op := range governanceapp.Operations() {
		if op.Action == nil || op.Input != governanceapp.RequestInput {
			continue
		}
		t.Run(op.Name, func(t *testing.T) {
			var out, errs bytes.Buffer
			code := governanceCommand(root, []string{op.Name, "--state", state, "--input", request}, &out, &errs, true)
			var cli struct {
				Error *InvokeError `json:"error"`
			}
			if err := json.Unmarshal(errs.Bytes(), &cli); err != nil {
				t.Fatal(err)
			}
			caller := ActionCaller{}
			if op.Write {
				caller.Principal = blueprintstate.Principal{Actor: "local-owner", Authority: "owner"}
			}
			result := Invoke(context.Background(), bytes.NewReader(requestBytes(t, root, op.Action.ID, SessionActionInput{State: state, RequestFile: "request.json"})), caller)
			if code != 2 || result.ExitCode != code || !reflect.DeepEqual(result.Error, cli.Error) {
				t.Fatalf("CLI=%d %+v, action=%+v", code, cli.Error, result)
			}
		})
	}
	if _, err := os.Stat(state); !os.IsNotExist(err) {
		t.Fatal("strict request rejection created state")
	}
}

func TestActionPolicyMatchesRegisteredMutation(t *testing.T) {
	for _, spec := range actionRegistry() {
		descriptor := spec.descriptor
		want := descriptor.ID == "blueprint.assignment.apply" || (strings.HasPrefix(descriptor.ID, "governance.session.") && descriptor.Mutation != "none")
		if spec.policy.callerAttribution != want {
			t.Fatalf("%s attribution=%v want %v", descriptor.ID, spec.policy.callerAttribution, want)
		}
		if want {
			continue
		}
		// Attribution must be refused before context resolution, for every other registration.
		result := Invoke(context.Background(), bytes.NewReader(requestBytes(t, "/absent", descriptor.ID, validPolicyParameters(descriptor.ID))), ActionCaller{Principal: blueprintstate.Principal{Actor: "owner", Authority: "owner"}})
		if result.Error == nil || result.Error.Code != "invalid-parameters" || result.Error.Message != "caller attribution applies only to assignment or Session mutation" {
			t.Fatalf("%s accepted attribution: %+v", descriptor.ID, result)
		}
	}
}
func validPolicyParameters(id string) any {
	switch id {
	case "project.model-tree":
		return ModelTreeOptions{Model: "model", Revision: "revision"}
	case "blueprint.generate":
		return GenerateActionInput{ModelArea: "model", From: "from", To: "to", Entry: "Main.design", Task: "task.json", Output: "output"}
	case "blueprint.assess":
		return AssessActionInput{Bundle: "bundle", Evidence: "evidence"}
	}
	if strings.HasPrefix(id, "governance.session.") {
		if strings.HasSuffix(id, ".list") {
			return SessionListInput{}
		}
		if strings.HasSuffix(id, ".read") || strings.HasSuffix(id, ".host-read") {
			return SessionReadInput{Instance: "instance"}
		}
		if !strings.HasSuffix(id, ".definition") {
			return SessionActionInput{RequestFile: "request.json"}
		}
	}
	return map[string]any{}
}
