package blueprints

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Hans-Einar/SDP/SDPTool/blueprintstate"
	"github.com/Hans-Einar/SDP/SDPTool/model"
	"github.com/Hans-Einar/SDP/SDPTool/projecthistory"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/blueprint"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/documents"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
)

func lifecycleFixture(t *testing.T) (string, blueprintstate.Binding, string, Evidence) {
	t.Helper()
	o, _ := setup(t)
	area := filepath.Join(t.TempDir(), "SDP")
	for _, p := range []string{"Models", "ProjectManagement", "Traceability", "Evidence"} {
		if e := os.MkdirAll(filepath.Join(area, p), 0700); e != nil {
			t.Fatal(e)
		}
	}
	if e := os.CopyFS(filepath.Join(area, "Models"), os.DirFS(o.Area)); e != nil {
		t.Fatal(e)
	}
	o.Area = filepath.Join(area, "Models")
	o.Output = ""
	task := blueprint.Task{ID: "lifecycle", Intent: "Inspect unchanged bounded model"}
	// Identical model inputs keep the workflow fixture independent of permission scope.
	o.To = o.From
	raw, _ := json.Marshal(task)
	os.WriteFile(o.Task, raw, 0600)
	r, e := GenerateRetained(context.Background(), o, filepath.Join(area, "Blueprints"))
	if e != nil {
		t.Fatal(e)
	}
	d, _, e := verify(r.Path, &budget{})
	if e != nil {
		t.Fatal(e)
	}
	rel, _ := filepath.Rel(area, r.Path)
	b := blueprintstate.Binding{Bundle: filepath.ToSlash(rel), Task: d.TaskID, System: d.Analysis.System, Revision: d.Revision, Retained: r.RetainedRevision, Plan: "PLAN-SDP-0001", Milestone: "M1", ModelArea: "Models", From: o.From, To: o.To}
	plan := `{"schemaVersion":"1.0","eventId":"EVT-PM-SDP-000001","eventType":"x-management:created","occurredAt":"2026-10-09T00:00:00Z","actor":"owner","commit":null,"subjectId":"PLAN-SDP-0001","payload":{"schemaVersion":"0.2","kind":"Plan","planType":"ImplementationPlan","previousEventId":null,"from":null,"to":"active","fromPath":null,"toPath":"05--Implementation/Test.md","reason":"fixture","links":[]}}` + "\n"
	os.WriteFile(filepath.Join(area, "ProjectManagement/Ledger.ndjson"), []byte(plan), 0600)
	ep := filepath.Join(area, "Evidence/ready.json")
	evidence := Evidence{Schema: EvidenceSchema, BlueprintRevision: d.Revision, RetainedRevision: r.RetainedRevision, TaskDigest: d.TaskDigest, Scope: "bounded lifecycle fixture", Code: map[string]CodeSnapshot{}}
	if a := deferUnknowns(t, r.Path, ep, &evidence); a.Status != "ready" {
		t.Fatal(a)
	}
	return area, b, ep, evidence
}
func implementationEvidence(t *testing.T, area string, b blueprintstate.Binding, in Evidence, id string) (string, string) {
	t.Helper()
	ep := filepath.Join(area, "Evidence", id+".json")
	root := filepath.Dir(ep)
	code := []byte("package pilot\nfunc Reduce(a,b int) int {return a+b}\n")
	os.MkdirAll(filepath.Join(root, "code"), 0700)
	os.WriteFile(filepath.Join(root, "code/main.go"), code, 0600)
	digest := model.Digest(model.Files{"main.go": code})
	in.Code = map[string]CodeSnapshot{"TARGET": {Revision: "fixture", Root: "code", Digest: digest, Files: map[string]string{"main.go": documents.Hash(code)}}}
	// This is a real bounded Go behavior check; receipt content is still attributed.
	os.WriteFile(filepath.Join(root, "code/go.mod"), []byte("module pilot\ngo 1.24\n"), 0600)
	test := []byte("package pilot\nimport \"testing\"\nfunc TestReduce(t *testing.T){if Reduce(2,3)!=5 {t.Fatal(\"reduction\")}}\n")
	os.WriteFile(filepath.Join(root, "code/main_test.go"), test, 0600)
	var cmd *exec.Cmd
	// Use the running toolchain even when GOROOT is not exported.
	goexe, e := exec.LookPath("go")
	if e != nil {
		t.Fatal(e)
	}
	cmd = exec.Command(goexe, "test", "./...")
	cmd.Dir = filepath.Join(root, "code")
	log, e := cmd.CombinedOutput()
	if e != nil {
		t.Fatal(string(log), e)
	}
	os.WriteFile(filepath.Join(root, "test.log"), log, 0600)
	trace := "EVT-IMPL-" + strings.ToUpper(id) + "-1"
	d, _, _ := verify(filepath.Join(area, b.Bundle), &budget{})
	inputs := model.Files{"code/main_test.go": test, "code/go.mod": []byte("module pilot\ngo 1.24\n")}
	hashes := map[string]string{}
	for p, data := range inputs {
		hashes[p] = documents.Hash(data)
	}
	in.RequiredChecks = []CheckRequirement{{ID: "reduction", Side: "TARGET", Obligation: "behavior:reduction"}}
	in.Receipts = []CheckReceipt{{ID: "reduction", Side: "TARGET", Obligation: "behavior:reduction", SourceDigest: d.Target.SourceDigest, TaskDigest: d.TaskDigest, CodeDigest: digest, Command: "go test ./...", Environment: "bounded Go fixture", InputsDigest: model.Digest(inputs), Inputs: hashes, Result: "pass", Artifact: "test.log", ArtifactDigest: documents.Hash(log), Traceability: trace}}
	a := deferUnknowns(t, filepath.Join(area, b.Bundle), ep, &in)
	if a.Status != "ready" {
		t.Fatal(a)
	}
	receipt := ImplementationReceipt{Schema: "sdp-blueprint-implementation/1", Assignment: id, Retained: b.Retained, EvidenceDigest: a.EvidenceDigest, CodeDigest: digest, Checks: []string{"reduction"}}
	event := map[string]any{"schemaVersion": "1.0", "eventId": trace, "eventType": "x-verification:recorded", "occurredAt": "2026-10-09T01:00:00Z", "actor": "worker", "commit": nil, "subjectId": id, "payload": receipt}
	raw, _ := json.Marshal(event)
	p := filepath.Join(area, "Traceability/Ledger.ndjson")
	h, e := projecthistory.Read(p)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = projecthistory.Append(p, h.Revision, raw); e != nil {
		t.Fatal(e)
	}
	rel, _ := filepath.Rel(area, ep)
	return filepath.ToSlash(rel), trace
}
func TestAssignmentLifecycleCLI(t *testing.T) {
	area, b, _, in := lifecycleFixture(t)
	_, bundleBefore, e := verify(filepath.Join(area, b.Bundle), &budget{})
	if e != nil {
		t.Fatal(e)
	}
	retainedBefore := model.Digest(model.Files(bundleBefore.Files))
	binary := filepath.Join(t.TempDir(), "sdptool")
	cmd := exec.Command("go", "build", "-o", binary, "../cmd/sdptool")
	if out, e := cmd.CombinedOutput(); e != nil {
		t.Fatal(string(out), e)
	}
	seq := 0
	last := map[string]string{}
	request := func(id, action string) blueprintstate.Request {
		seq++
		return blueprintstate.Request{Schema: blueprintstate.Schema, ID: id, Action: action, EventID: fmt.Sprintf("EVT-BPA-TEST-%06d", seq), OccurredAt: fmt.Sprintf("2026-10-09T02:%02d:00Z", seq), Expected: last[id], Reason: "bounded lifecycle verification"}
	}
	run := func(r blueprintstate.Request, actor, authority string, want bool) AssignmentResult {
		t.Helper()
		raw, _ := json.Marshal(r)
		p := filepath.Join(t.TempDir(), "request.json")
		os.WriteFile(p, raw, 0600)
		args := []string{area, "model", "assignment", "apply", "--request", p, "--as", actor, "--authority", authority, "--json"}
		beforeHistory := read(t, filepath.Join(area, "ProjectManagement/Ledger.ndjson"))
		out, e := exec.Command(binary, args...).CombinedOutput()
		if !want && !bytes.Equal(beforeHistory, read(t, filepath.Join(area, "ProjectManagement/Ledger.ndjson"))) {
			t.Fatal("refused request mutated history")
		}
		if (e == nil) != want {
			t.Fatalf("%s: %s %v", r.Action, out, e)
		}
		var result AssignmentResult
		if want {
			if e = json.Unmarshal(out, &result); e != nil {
				t.Fatal(e)
			}
			last[r.ID] = result.Assignment.EventID
		}
		return result
	}
	id := "BPA-one"
	r := request(id, "create")
	r.Binding = &b
	run(r, "owner", "controller", true)
	r = request(id, "adopt-readiness")
	r.Evidence = "Evidence/ready.json"
	run(r, "worker", "assignee", false)
	run(r, "owner", "controller", true)
	r = request(id, "assign")
	r.Assignee = "worker"
	run(r, "owner", "controller", true)
	stale := request(id, "start")
	stale.Expected = "EVT-NOTCURRENT-1"
	run(stale, "worker", "assignee", false)
	r = request(id, "start")
	run(r, "intruder", "assignee", false)
	run(r, "worker", "assignee", true)
	start := r
	r = request(id, "hold")
	run(r, "owner", "controller", true)
	r = request(id, "resume")
	run(r, "owner", "controller", true)
	retry := run(start, "worker", "assignee", true)
	if retry.Appended || retry.Assignment.State != "in-progress" {
		t.Fatal(retry)
	}
	ep, trace := implementationEvidence(t, area, b, in, id)
	r = request(id, "submit")
	r.Evidence = ep
	r.Trace = trace
	run(r, "worker", "assignee", true)
	r = request(id, "accept-review")
	run(r, "worker", "reviewer", false)
	r = request(id, "complete")
	run(r, "owner", "controller", false)
	r = request(id, "reject")
	run(r, "reviewer", "reviewer", true)
	r = request(id, "submit")
	r.Evidence = ep
	r.Trace = trace
	run(r, "worker", "assignee", true)
	r = request(id, "accept-review")
	run(r, "reviewer", "reviewer", true)
	r = request(id, "complete")
	done := run(r, "owner", "controller", true)
	if done.Assignment.State != "completed" {
		t.Fatal(done)
	}
	r = request(id, "start")
	run(r, "worker", "assignee", false)
	// Another attempt against the same revision retains its own state.
	r = request("BPA-two", "create")
	r.Binding = &b
	run(r, "owner", "controller", true)
	r = request("BPA-three", "create")
	r.Binding = &b
	run(r, "owner", "controller", true)
	r = request("BPA-two", "supersede")
	r.Successor = "BPA-three"
	run(r, "owner", "controller", true)
	r = request("BPA-three", "cancel")
	run(r, "owner", "controller", true)
	// A new retained revision of the same task never inherits completion.
	_, e = model.CreateWork(filepath.Join(area, "Models"), "Later", "release:0.1.0", false)
	if e != nil {
		t.Fatal(e)
	}
	second, e := GenerateRetained(context.Background(), Options{Area: filepath.Join(area, "Models"), From: b.From, To: "work:Later", Entry: "System.design", Task: copiedTask(t, area, b)}, filepath.Join(area, "Blueprints"))
	if e != nil {
		t.Fatal(e)
	}
	b2 := b
	b2.To = "work:Later"
	b2.Revision = second.Revision
	b2.Retained = second.RetainedRevision
	rel, _ := filepath.Rel(area, second.Path)
	b2.Bundle = filepath.ToSlash(rel)
	r = request("BPA-four", "create")
	r.Binding = &b2
	run(r, "owner", "controller", true)
	views, e := AssignmentViews(area)
	if e != nil || len(views) != 4 {
		t.Fatal(views, e)
	}
	// Cross-language replay must accept the exact canonical chain.
	h, _ := projecthistory.Read(filepath.Join(area, "ProjectManagement/Ledger.ndjson"))
	cmd = exec.Command("python3", "-c", `import json,runpy,sys; f=runpy.run_path('../../SDP/ProjectManagement/validate.py'); f['history']([json.loads(x) for x in open(sys.argv[1])])`, filepath.Join(area, "ProjectManagement/Ledger.ndjson"))
	if out, e := cmd.CombinedOutput(); e != nil {
		t.Fatal(string(out), e)
	}
	// Discover -> groups -> open target, using the compiled consumer contract.
	out, e := exec.Command(binary, filepath.Dir(area), "discover", "--json").CombinedOutput()
	if e != nil {
		t.Fatal(string(out), e)
	}
	if !bytes.Contains(out, []byte(`"workState": "completed"`)) && !bytes.Contains(out, []byte(`"workState":"completed"`)) {
		t.Fatal("missing completed projection", string(out))
	}
	if !bytes.Contains(out, []byte("blueprints/assignments/completed")) {
		t.Fatal("missing group")
	}
	// Retained history is not rewritten on new inputs. Remove one bundle as a negative consumer control.
	index := filepath.Join(area, b.Bundle, "index.md")
	before := read(t, index)
	os.WriteFile(index, []byte("corrupt"), 0600)
	views, e = AssignmentViews(area)
	if e != nil {
		t.Fatal(e)
	}
	for _, v := range views {
		if v.ID == id && (v.State.State != "completed" || v.OpenPath != "" || v.Readiness != "changed-or-unavailable" || v.EvidenceStatus != "changed-or-unavailable") {
			t.Fatal(v)
		}
	}
	os.WriteFile(index, before, 0600)
	after, _ := projecthistory.Read(filepath.Join(area, "ProjectManagement/Ledger.ndjson"))
	if !bytes.Equal(h.Bytes, after.Bytes) {
		t.Fatal("read mutated history")
	}
	_, bundleAfter, e := verify(filepath.Join(area, b.Bundle), &budget{})
	if e != nil || model.Digest(model.Files(bundleAfter.Files)) != retainedBefore {
		t.Fatal("retained bundle modified", e)
	}
	// Optional reproducible trial export, only into an explicitly supplied empty destination.
	if dest := os.Getenv("SDP_BLUEPRINT_LIFECYCLE_EXPORT"); dest != "" {
		if _, e = os.Stat(dest); !os.IsNotExist(e) {
			t.Fatal("export destination exists")
		}
		if e = os.CopyFS(dest, os.DirFS(area)); e != nil {
			t.Fatal(e)
		}
	}
}
func TestAssignmentStaleAndEvidenceRefusal(t *testing.T) {
	area, b, ep, in := lifecycleFixture(t)
	n := 0
	state := blueprintstate.State{}
	step := func(action string) blueprintstate.Request {
		n++
		return blueprintstate.Request{Schema: blueprintstate.Schema, ID: "BPA-one", Action: action, EventID: fmt.Sprintf("EVT-TEST-%d", n), OccurredAt: "2026-10-09T03:00:00Z", Expected: state.EventID, Reason: "test"}
	}
	apply := func(r blueprintstate.Request, p blueprintstate.Principal, want bool) {
		t.Helper()
		result, e := ApplyAssignment(area, p, r)
		if (e == nil) != want {
			t.Fatal(r.Action, result, e)
		}
		if want {
			state = result.Assignment
		}
	}
	controller := blueprintstate.Principal{Actor: "owner", Authority: "controller"}
	worker := blueprintstate.Principal{Actor: "worker", Authority: "assignee"}
	r := step("create")
	r.Binding = &b
	apply(r, controller, true)
	// Blocked evidence cannot be adopted.
	saved := read(t, ep)
	in.Dispositions = nil
	writeEvidence(t, ep, in)
	r = step("adopt-readiness")
	r.Evidence = "Evidence/ready.json"
	apply(r, controller, false)
	os.WriteFile(ep, saved, 0600)
	apply(r, controller, true)
	r = step("assign")
	r.Assignee = "worker"
	apply(r, controller, true)
	r = step("start")
	apply(r, worker, true)
	evidence, trace := implementationEvidence(t, area, b, in, "BPA-one")
	r = step("submit")
	r.Evidence = evidence
	r.Trace = trace
	code := filepath.Join(area, "Evidence/code/main.go")
	old := read(t, code)
	os.WriteFile(code, []byte("wrong code"), 0600)
	apply(r, worker, false)
	os.WriteFile(code, old, 0600)
	apply(r, worker, true)
	r = step("accept-review")
	apply(r, blueprintstate.Principal{Actor: "independent", Authority: "reviewer"}, true)
	// Trace record tamper must be detected at closure.
	tp := filepath.Join(area, "Traceability/Ledger.ndjson")
	tr := read(t, tp)
	os.WriteFile(tp, bytes.Replace(tr, []byte(`"actor":"worker"`), []byte(`"actor":"other"`), 1), 0600)
	r = step("complete")
	apply(r, controller, false)
	os.WriteFile(tp, tr, 0600)
	// Request replay rejects changed content under an operation identity.
	changed := r
	changed.EventID = state.EventID
	apply(changed, controller, false)
	apply(r, controller, true)
	if _, e := os.Stat(filepath.Join(area, ".git")); !os.IsNotExist(e) {
		t.Fatal("unexpected Git dependency")
	}
	if !reflect.DeepEqual(state.Binding, b) {
		t.Fatal("binding mutated")
	}
}

func TestAssignmentFreshnessAndPlanCleanup(t *testing.T) {
	area, b, _, _ := lifecycleFixture(t)
	_, e := model.CreateWork(filepath.Join(area, "Models"), "Live", "release:0.1.0", false)
	if e != nil {
		t.Fatal(e)
	}
	result, e := GenerateRetained(context.Background(), Options{Area: filepath.Join(area, "Models"), From: b.From, To: "work:Live", Entry: "System.design", Task: copiedTask(t, area, b)}, filepath.Join(area, "Blueprints"))
	if e != nil {
		t.Fatal(e)
	}
	b.To = "work:Live"
	b.Revision = result.Revision
	b.Retained = result.RetainedRevision
	rel, _ := filepath.Rel(area, result.Path)
	b.Bundle = filepath.ToSlash(rel)
	c := blueprintstate.Principal{Actor: "owner", Authority: "controller"}
	r := blueprintstate.Request{Schema: blueprintstate.Schema, ID: "BPA-cleanup", Action: "create", EventID: "EVT-FRESH-1", OccurredAt: "2026-10-09T04:00:00Z", Reason: "fixture", Binding: &b}
	created, e := ApplyAssignment(area, c, r)
	if e != nil {
		t.Fatal(e)
	}
	w, e := model.Snapshot(filepath.Join(area, "Models"), "work:Live")
	if e != nil {
		t.Fatal(e)
	}
	source := filepath.Join(w.Path, "System.design")
	old := read(t, source)
	os.WriteFile(source, append(old, '\n'), 0600)
	r.Binding = nil
	r.Action = "adopt-readiness"
	r.Evidence = "Evidence/ready.json"
	r.Expected = created.Assignment.EventID
	r.EventID = "EVT-FRESH-2"
	if _, e = ApplyAssignment(area, c, r); e == nil {
		t.Fatal("stale model adopted")
	}
	views, e := AssignmentViews(area)
	if e != nil || views[0].Freshness != "stale" {
		t.Fatal(views, e)
	}
	history := filepath.Join(area, "ProjectManagement/Ledger.ndjson")
	raw := read(t, history)
	// Externally closed plan: assignment cleanup must remain possible.
	var plan map[string]any
	if e = json.Unmarshal(bytes.Split(raw, []byte{'\n'})[0], &plan); e != nil {
		t.Fatal(e)
	}
	payload := plan["payload"].(map[string]any)
	payload["previousEventId"] = plan["eventId"]
	payload["from"] = "active"
	payload["to"] = "completed"
	payload["fromPath"] = payload["toPath"]
	payload["links"] = []string{"PLAN-SDP-0001"}
	plan["eventId"] = "EVT-PM-SDP-000002"
	plan["eventType"] = "x-management:completed"
	closed, _ := json.Marshal(plan)
	raw = append(append(raw, closed...), '\n')
	os.WriteFile(history, raw, 0600)
	r.Evidence = ""
	r.Action = "cancel"
	r.EventID = "EVT-FRESH-3"
	canceled, e := ApplyAssignment(area, c, r)
	if e != nil || canceled.Assignment.State != "canceled" {
		t.Fatal(canceled, e)
	}
}

func copiedTask(t *testing.T, area string, b blueprintstate.Binding) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "task.json")
	if e := os.WriteFile(p, read(t, filepath.Join(area, b.Bundle, "assignment.json")), 0600); e != nil {
		t.Fatal(e)
	}
	return p
}

func TestAssignmentProjectionBudgetAndEvidenceDiagnostics(t *testing.T) {
	area, b, ep, _ := lifecycleFixture(t)
	hpath := filepath.Join(area, "ProjectManagement/Ledger.ndjson")
	h, _ := projecthistory.Read(hpath)
	raw := h.Bytes
	// Four assignments share one verified bundle. Distinct milestones force three
	// binding keys; the first two share one key and must reuse its source result.
	for i := 0; i < 4; i++ {
		binding := b
		if i > 1 {
			binding.Milestone = fmt.Sprintf("M%d", i)
		}
		r := blueprintstate.Request{Schema: blueprintstate.Schema, ID: fmt.Sprintf("BPA-%d", i), Action: "create", EventID: fmt.Sprintf("EVT-CACHE-%d", i), OccurredAt: "2026-10-09T05:00:00Z", Reason: "projection", Binding: &binding}
		event := blueprintstate.Event{SchemaVersion: "1.0", EventID: r.EventID, EventType: "x-blueprint:create", OccurredAt: r.OccurredAt, Actor: "owner", SubjectID: r.ID, Payload: blueprintstate.Payload{Schema: blueprintstate.Schema, Authority: "controller", Request: r}}
		bytes, _ := json.Marshal(event)
		raw = append(append(raw, bytes...), '\n')
	}
	os.WriteFile(hpath, raw, 0600)
	views, e := AssignmentViews(area)
	if e != nil {
		t.Fatal(e)
	}
	if len(views) != 4 || views[0].Freshness != "current" || views[1].Freshness != "current" || views[2].Freshness != "current" || views[3].Freshness != "unknown" || !strings.Contains(views[3].Diagnostic, "budget") {
		t.Fatal(views)
	}
	// A missing adopted evidence file must surface a reason, not erase adoption.
	r := blueprintstate.Request{Schema: blueprintstate.Schema, ID: "BPA-0", Action: "adopt-readiness", EventID: "EVT-CACHE-5", OccurredAt: "2026-10-09T05:00:00Z", Expected: "EVT-CACHE-0", Reason: "adopt", Evidence: "Evidence/ready.json"}
	if _, e = ApplyAssignment(area, blueprintstate.Principal{Actor: "owner", Authority: "controller"}, r); e != nil {
		t.Fatal(e)
	}
	os.Remove(ep)
	views, e = AssignmentViews(area)
	if e != nil || views[0].Readiness != "changed-or-unavailable" || !strings.Contains(views[0].Diagnostic, "readiness:") {
		t.Fatal(views, e)
	}
	// All bundle/evidence reads share one exhausted projection budget.
	reads := &budget{bytes: MaxCatalogueBytes}
	if _, _, e = boundCapture(area, b, reads); e == nil {
		t.Fatal("bundle ignored budget")
	}
	if _, e = assignmentProofBudget(area, b, "Evidence/ready.json", "", "BPA-0", reads); e == nil {
		t.Fatal("assessment ignored shared budget")
	}
}

func TestAssignmentInterruptedHelper(t *testing.T) {
	area := os.Getenv("SDP_ASSIGNMENT_CRASH_AREA")
	if area == "" {
		return
	}
	req, e := LoadAssignmentRequest(os.Getenv("SDP_ASSIGNMENT_CRASH_REQUEST"))
	if e != nil {
		os.Exit(72)
	}
	if _, e = ApplyAssignment(area, blueprintstate.Principal{Actor: "owner", Authority: "controller"}, req); e != nil {
		os.Exit(72)
	}
	os.Exit(73) // committed operation, no response delivered
}
func TestAssignmentRetryAfterLostResponse(t *testing.T) {
	area, b, _, _ := lifecycleFixture(t)
	r := blueprintstate.Request{Schema: blueprintstate.Schema, ID: "BPA-crash", Action: "create", EventID: "EVT-CRASH-1", OccurredAt: "2026-10-09T06:00:00Z", Reason: "lost response", Binding: &b}
	raw, _ := json.Marshal(r)
	p := filepath.Join(t.TempDir(), "request.json")
	os.WriteFile(p, raw, 0600)
	cmd := exec.Command(os.Args[0], "-test.run=^TestAssignmentInterruptedHelper$")
	cmd.Env = append(os.Environ(), "SDP_ASSIGNMENT_CRASH_AREA="+area, "SDP_ASSIGNMENT_CRASH_REQUEST="+p)
	if err := cmd.Run(); err == nil || cmd.ProcessState.ExitCode() != 73 {
		t.Fatal("helper", err)
	}
	result, err := ApplyAssignment(area, blueprintstate.Principal{Actor: "owner", Authority: "controller"}, r)
	if err != nil || result.Appended || result.Assignment.State != "draft" {
		t.Fatal(result, err)
	}
	h, err := projecthistory.Read(filepath.Join(area, "ProjectManagement/Ledger.ndjson"))
	if err != nil {
		t.Fatal(err)
	}
	_, events, err := blueprintstate.Replay(h.Bytes)
	if err != nil || len(events) != 1 {
		t.Fatal(events, err)
	}
}

func TestAssignmentConcurrentCancellation(t *testing.T) {
	area, b, _, _ := lifecycleFixture(t)
	c := blueprintstate.Principal{Actor: "owner", Authority: "controller"}
	r := blueprintstate.Request{Schema: blueprintstate.Schema, ID: "BPA-race", Action: "create", EventID: "EVT-RACE-1", OccurredAt: "2026-10-09T06:00:00Z", Reason: "concurrency", Binding: &b}
	created, e := ApplyAssignment(area, c, r)
	if e != nil {
		t.Fatal(e)
	}
	r.Binding = nil
	r.Action = "cancel"
	r.Expected = created.Assignment.EventID
	results := make(chan error, 2)
	var group sync.WaitGroup
	for i := 0; i < 2; i++ {
		group.Add(1)
		go func(i int) {
			defer group.Done()
			req := r
			req.EventID = fmt.Sprintf("EVT-RACE-%d", i+2)
			_, err := ApplyAssignment(area, c, req)
			results <- err
		}(i)
	}
	group.Wait()
	close(results)
	success := 0
	for err := range results {
		if err == nil {
			success++
		}
	}
	if success != 1 {
		t.Fatal("concurrent update winners", success)
	}
	h, e := projecthistory.Read(filepath.Join(area, "ProjectManagement/Ledger.ndjson"))
	if e != nil {
		t.Fatal(e)
	}
	states, events, e := blueprintstate.Replay(h.Bytes)
	if e != nil || len(events) != 2 || states[r.ID].State != "canceled" {
		t.Fatal(states, events, e)
	}
}
