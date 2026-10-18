package sdptool

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"path/filepath"
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
	if e != nil || json.Unmarshal(b, &metadata) != nil || len(metadata.Actions) != 5 {
		t.Fatalf("catalogue %s %v", b, e)
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
