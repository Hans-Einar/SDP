package sdptool

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
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
		var expected any
		if code != 0 || json.Unmarshal(out.Bytes(), &expected) != nil || r.ExitCode != 0 || !reflect.DeepEqual(expected, r.Result) {
			t.Fatalf("%s: %+v %s", op, r, errs.String())
		}
	}
	for _, op := range []string{"adopt", "resume", "select", "assess", "start", "progress", "verify", "review"} {
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
