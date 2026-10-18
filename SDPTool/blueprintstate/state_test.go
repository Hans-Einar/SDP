package blueprintstate

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestCanonicalReplayAndPythonParity(t *testing.T) {
	raw, e := os.ReadFile("testdata/lifecycle.ndjson")
	if e != nil {
		t.Fatal(e)
	}
	tests := []struct {
		name  string
		data  []byte
		valid bool
	}{{"full-chain", raw, true}}
	tests = append(tests, struct {
		name  string
		data  []byte
		valid bool
	}{"missing-required", bytes.Replace(raw, []byte(`"expectedEvent":"",`), nil, 1), false})
	tests = append(tests, struct {
		name  string
		data  []byte
		valid bool
	}{"null-required", bytes.Replace(raw, []byte(`"expectedEvent":""`), []byte(`"expectedEvent":null`), 1), false})
	tests = append(tests, struct {
		name  string
		data  []byte
		valid bool
	}{"case-alias", bytes.Replace(raw, []byte(`"assignmentId":`), []byte(`"AssignmentId":`), 1), false})
	tests = append(tests, struct {
		name  string
		data  []byte
		valid bool
	}{"duplicate-event", append(append([]byte{}, raw...), bytes.Split(raw, []byte{'\n'})[1]...), false})
	// Mutate typed event chains, preserving valid JSON/envelopes.
	for _, name := range []string{"self-review", "wrong-predecessor", "cross-system-successor", "time-reversal"} {
		var lines [][]byte
		changed := false
		for _, line := range bytes.Split(raw, []byte{'\n'}) {
			if len(line) == 0 {
				continue
			}
			var event Event
			if Strict(line, &event) != nil {
				lines = append(lines, line)
				continue
			}
			r := &event.Payload.Request
			switch name {
			case "self-review":
				if r.Action == "accept-review" && !changed {
					event.Actor = "worker"
					changed = true
				}
			case "wrong-predecessor":
				if r.Action == "assign" && !changed {
					r.Expected = "EVT-MISSING-1"
					changed = true
				}
			case "cross-system-successor":
				if r.Action == "create" && r.ID == "BPA-three" {
					r.Binding.System = "Other"
					changed = true
				}
			case "time-reversal":
				if r.Action == "assign" && !changed {
					r.OccurredAt = "2026-10-09T00:00:00.000000001Z"
					event.OccurredAt = r.OccurredAt
					changed = true
				}
			}
			b, _ := json.Marshal(event)
			lines = append(lines, b)
		}
		if !changed {
			t.Fatal("missing control", name)
		}
		tests = append(tests, struct {
			name  string
			data  []byte
			valid bool
		}{name, append(bytes.Join(lines, []byte{'\n'}), '\n'), false})
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, _, err := Replay(test.data)
			if (err == nil) != test.valid {
				t.Fatal(err)
			}
			p := filepath.Join(t.TempDir(), "history.ndjson")
			os.WriteFile(p, test.data, 0600)
			cmd := exec.Command("python3", "-c", `import json,runpy,sys; f=runpy.run_path('../../SDP/ProjectManagement/validate.py'); f['history']([json.loads(x) for x in open(sys.argv[1])])`, p)
			out, err := cmd.CombinedOutput()
			if (err == nil) != test.valid {
				t.Fatalf("Python parity: %s %v", out, err)
			}
		})
	}
}
func TestNanosecondOrdering(t *testing.T) {
	r := Request{Schema: Schema, ID: "BPA-nano", Action: "cancel", EventID: "EVT-NANO-2", Expected: "EVT-NANO-1", OccurredAt: "2026-10-09T05:00:00.000000001Z", Reason: "test"}
	event := Event{SchemaVersion: "1.0", EventID: r.EventID, EventType: "x-blueprint:cancel", OccurredAt: r.OccurredAt, Actor: "owner", SubjectID: r.ID, Payload: Payload{Schema: Schema, Authority: "controller", Request: r}}
	old := State{ID: r.ID, State: "draft", EventID: r.Expected, OccurredAt: "2026-10-09T05:00:00.000000002Z"}
	if _, e := Next(old, event); e == nil {
		t.Fatal("time reversal")
	}
	cmd := exec.Command("python3", "-c", `import runpy; f=runpy.run_path('../../SDP/ProjectManagement/blueprint_history.py'); assert f['instant']('2026-10-09T05:00:00.000000001Z') < f['instant']('2026-10-09T05:00:00.000000002Z')`)
	if out, e := cmd.CombinedOutput(); e != nil {
		t.Fatal(string(out), e)
	}
}
