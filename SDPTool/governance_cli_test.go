package sdptool

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
)

// Each command runs in a fresh process, so the journey exercises durable restart
// and the actual CLI adapter, independently of an in-memory Service instance.
func TestCompiledGovernanceJourney(t *testing.T) {
	binary := filepath.Join(t.TempDir(), "sdptool")
	build := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "build", "-o", binary, "./cmd/sdptool")
	if b, e := build.CombinedOutput(); e != nil {
		t.Fatalf("build: %v %s", e, b)
	}
	root, _ := projectFixture(t)
	files := map[string]string{
		"SDP/KanBan/active/card.md": "| Field | Value |\n| --- | --- |\n| id | KB-TEST-001 |\n| CardState | ready |\n| PlanId | MAINT-TEST-0001 |\n",
		"SDP/Maintenance/Plan.md":   "| Field | Value |\n| --- | --- |\n| id | MAINT-TEST-0001 |\n| PlanType | MaintenancePlan |\n| state | active |\n",
		"SDP/Sessions/session.md":   "| State | Step |\n| --- | --- |\n| next | S1 |\n\n| Field | Value |\n| --- | --- |\n| Session reference | SESSION-TEST-0001 |\n| Status | active |\n",
		"SDP/authority.md":          "Fixture owner selects bounded work.\n",
		"dirty.txt":                 "preserve owner edits\n",
		"SDP/ProjectManagement/Ledger.ndjson": `{"schemaVersion":"1.0","eventId":"E1","eventType":"x-kanban:created","subjectId":"KB-TEST-001","payload":{"schemaVersion":"0.2","previousEventId":null,"to":"active","toPath":"active/card.md"}}
{"schemaVersion":"1.0","eventId":"E2","eventType":"x-management:started","subjectId":"MAINT-TEST-0001","payload":{"schemaVersion":"0.2","previousEventId":null,"to":"active","toPath":"Maintenance/Plan.md"}}
`,
	}
	for p, b := range files {
		f := filepath.Join(root, p)
		if e := os.MkdirAll(filepath.Dir(f), 0700); e != nil {
			t.Fatal(e)
		}
		if e := os.WriteFile(f, []byte(b), 0600); e != nil {
			t.Fatal(e)
		}
	}
	fingerprints := func() map[string][32]byte {
		result := map[string][32]byte{}
		e := filepath.WalkDir(root, func(p string, d os.DirEntry, e error) error {
			if e != nil {
				return e
			}
			if d.IsDir() {
				return nil
			}
			b, e := os.ReadFile(p)
			if e != nil {
				return e
			}
			result[p] = sha256.Sum256(b)
			return nil
		})
		if e != nil {
			t.Fatal(e)
		}
		return result
	}
	before := fingerprints()
	state := filepath.Join(t.TempDir(), "state")
	requestDir := t.TempDir()
	sequence := 0
	cases := []map[string]any{}
	capture := func(name string, v any) { cases = append(cases, map[string]any{"definition": name, "value": v}) }
	call := func(op string, q any, wantCode int, extra ...string) []byte {
		t.Helper()
		args := []string{root, "governance", op, "--state", state, "--json"}
		if q != nil {
			sequence++
			p := filepath.Join(requestDir, fmt.Sprintf("%d.json", sequence))
			b, e := json.Marshal(q)
			if e != nil {
				t.Fatal(e)
			}
			if e = os.WriteFile(p, b, 0600); e != nil {
				t.Fatal(e)
			}
			args = append(args, "--input", p)
		}
		args = append(args, extra...)
		cmd := exec.Command(binary, args...)
		var out, errs bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errs
		e := cmd.Run()
		code := 0
		if e != nil {
			if x, ok := e.(*exec.ExitError); ok {
				code = x.ExitCode()
			} else {
				t.Fatal(e)
			}
		}
		if code != wantCode {
			t.Fatalf("%s code %d: %s %s", op, code, out.String(), errs.String())
		}
		if code != 0 {
			var v any
			if e = json.Unmarshal(errs.Bytes(), &v); e != nil {
				t.Fatal(e, errs.String())
			}
			capture("Error", v)
			return errs.Bytes()
		}
		return out.Bytes()
	}
	decode := func(b []byte, v any) {
		t.Helper()
		if e := json.Unmarshal(b, v); e != nil {
			t.Fatal(e, string(b))
		}
	}
	sel := governance.Selection{Card: "SDP/KanBan/active/card.md", Plan: "SDP/Maintenance/Plan.md", Session: "SDP/Sessions/session.md", Authority: "SDP/authority.md", Step: "S1"}
	capture("Selection", sel)
	var c governance.Context
	decode(call("context", sel, 0), &c)
	capture("Context", c)
	rq := governance.RouteRequest{Schema: governance.Schema, Categories: []string{"status"}}
	capture("RouteRequest", rq)
	var route governance.Route
	decode(call("route", rq, 0), &route)
	capture("Route", route)
	if route.Disposition != "no-work" {
		t.Fatal(route)
	}
	rq.Categories = []string{"maintenance"}
	rq.Selection = sel
	rq.RoutineID = "maintenance"
	rq.RoutineVersion = "0.1"
	decode(call("route", rq, 0), &route)
	capture("RouteRequest", rq)
	capture("Route", route)
	if route.Disposition != "matched" {
		t.Fatal(route)
	}
	if _, e := os.Stat(state); !os.IsNotExist(e) {
		t.Fatal("read created state")
	}
	var bound governance.Binding
	decode(call("bind", nil, 0), &bound)
	capture("Binding", bound)
	var again governance.Binding
	decode(call("bind", nil, 0), &again)
	if bound != again {
		t.Fatal("binding changed")
	}
	uuid := func(i int) string { return fmt.Sprintf("00000000-0000-4000-8000-%012d", i) }
	q := governance.StartRequest{Schema: governance.Schema, ProjectRef: bound.ProjectRef, CheckoutID: bound.CheckoutID, OperationID: uuid(1), RequestID: uuid(2), InputRevision: c.Revision, RoutineDigest: route.RoutineDigest, Selection: sel, Assignment: governance.Assignment{Outcome: "fixture fix", WritablePaths: []string{"target.txt"}, ExcludedPaths: []string{"dirty.txt"}}}
	capture("StartRequest", q)
	var run governance.Result
	decode(call("start", q, 0), &run)
	capture("Result", run)
	var replay governance.Result
	decode(call("start", q, 0), &replay)
	if !reflect.DeepEqual(run, replay) {
		t.Fatal("retry changed receipt")
	}
	for i, action := range []string{"start", "interrupt", "resume", "start", "report"} {
		tq := governance.TransitionRequest{Schema: governance.Schema, ProjectRef: bound.ProjectRef, CheckoutID: bound.CheckoutID, OperationID: uuid(i + 10), RunID: run.Run.ID, ExpectedRevision: run.Run.Revision, InputRevision: c.Revision, Action: action, Reason: "local fixture observation"}
		capture("TransitionRequest", tq)
		decode(call("transition", tq, 0), &run)
		capture("Result", run)
	}
	tq := governance.TransitionRequest{Schema: governance.Schema, ProjectRef: bound.ProjectRef, CheckoutID: bound.CheckoutID, OperationID: uuid(30), RunID: run.Run.ID, ExpectedRevision: run.Run.Revision, InputRevision: c.Revision, Action: "complete", Reason: "attempt unsupported completion"}
	if b := call("transition", tq, 2); !bytes.Contains(b, []byte("evidence_missing")) {
		t.Fatal(string(b))
	}
	decode(call("read", nil, 0, "--run", run.Run.ID), &replay)
	capture("Result", replay)
	if replay.Run.State != "verifying" || replay.Run.Revision != 6 {
		t.Fatal(replay)
	}
	cmd := exec.Command(binary, root, "governance", "read", "--state", state, "--run", run.Run.ID)
	if b, e := cmd.CombinedOutput(); e != nil || json.Valid(b) || !strings.Contains(string(b), "verifying") {
		t.Fatal(string(b), e)
	}
	if !reflect.DeepEqual(before, fingerprints()) {
		t.Fatal("project records or unrelated inputs changed")
	}
	if dir := os.Getenv("PGI_SCHEMA_EVIDENCE"); dir != "" {
		if e := os.MkdirAll(dir, 0700); e != nil {
			t.Fatal(e)
		}
		b, e := json.MarshalIndent(cases, "", "  ")
		if e != nil {
			t.Fatal(e)
		}
		if e = os.WriteFile(filepath.Join(dir, "cli-wire-cases.json"), b, 0600); e != nil {
			t.Fatal(e)
		}
	}
}
