package sdptool

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func TestCompiledInvocationConsumer(t *testing.T) {
	binary := filepath.Join(t.TempDir(), "sdptool")
	if b, e := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "build", "-o", binary, "./cmd/sdptool").CombinedOutput(); e != nil {
		t.Fatal(e, string(b))
	}
	root, _ := projectFixture(t)
	var metadata ActionCatalogue
	b, e := exec.Command(binary, root, "actions", "--json").Output()
	if e != nil || json.Unmarshal(b, &metadata) != nil {
		t.Fatalf("catalogue %s %v", b, e)
	}
	var ids []string
	for _, action := range metadata.Actions {
		ids = append(ids, action.ID)
	}
	wantIDs := []string{"blueprint.assess", "blueprint.assignment.apply", "blueprint.assignment.list", "blueprint.generate", "governance.session.adopt", "governance.session.assess", "governance.session.context", "governance.session.definition", "governance.session.list", "governance.session.progress", "governance.session.read", "governance.session.resume", "governance.session.select", "governance.session.start", "project.discover", "project.model-tree"}
	if !reflect.DeepEqual(ids, wantIDs) {
		t.Fatalf("catalogue actions: got %v, want %v", ids, wantIDs)
	}
	for _, tc := range []struct {
		action string
		want   int
	}{{"project.discover", 0}, {"not-supported", 2}} {
		cmd := exec.Command(binary, "invoke", "--request", "-")
		cmd.Stdin = bytes.NewReader(requestBytes(t, root, tc.action, map[string]any{}))
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		output, err := cmd.Output()
		code := 0
		if err != nil {
			exit, ok := err.(*exec.ExitError)
			if !ok {
				t.Fatal(err)
			}
			code = exit.ExitCode()
		}
		var result InvokeResponse
		if code != tc.want || json.Unmarshal(output, &result) != nil || result.ExitCode != code || stderr.Len() != 0 {
			t.Fatalf("%s: code=%d %s %s", tc.action, code, output, &stderr)
		}
		// Even nonzero exit responses are retained and decoded by this reference consumer.
		if result.RequestID == nil || *result.RequestID != "trial-1" {
			t.Fatal("correlation lost")
		}
	}
}
