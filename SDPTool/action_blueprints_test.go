package sdptool

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/Hans-Einar/SDP/SDPTool/blueprints"
	"github.com/Hans-Einar/SDP/SDPTool/blueprintstate"
	"github.com/Hans-Einar/SDP/SDPTool/model"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/blueprint"
)

func writeActionJSON(t *testing.T, path string, v any) {
	t.Helper()
	b, e := json.Marshal(v)
	if e != nil {
		t.Fatal(e)
	}
	if e = os.WriteFile(path, b, 0600); e != nil {
		t.Fatal(e)
	}
}
func actionBlueprintFixture(t *testing.T) (string, GenerateActionInput) {
	t.Helper()
	root, _ := projectFixture(t)
	area := filepath.Join(root, "SDP/Models")
	if err := os.MkdirAll(area, 0700); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"Before", "After"} {
		w, e := model.CreateWork(area, name, "", true)
		if e != nil {
			t.Fatal(e)
		}
		if e = os.CopyFS(w.Path, os.DirFS("../experiments/blueprint_mvp1/NOW")); e != nil {
			t.Fatal(e)
		}
	}
	writeActionJSON(t, filepath.Join(root, "task.json"), blueprint.Task{ID: "actions", Intent: "Inspect bounded unchanged model"})
	return root, GenerateActionInput{ModelArea: "SDP/Models", From: "work:Before", To: "work:After", Entry: "System.design", Task: "task.json", Catalogue: "SDP/Blueprints"}
}
func TestInvokeBlueprintWorkflow(t *testing.T) {
	root, input := actionBlueprintFixture(t)
	ctx := context.Background()
	invoke := func(action string, p any, caller ActionCaller) InvokeResponse {
		return Invoke(ctx, bytes.NewReader(requestBytes(t, root, action, p)), caller)
	}
	r := invoke("blueprint.generate", input, ActionCaller{})
	if r.ExitCode != 0 {
		t.Fatalf("%+v", r)
	}
	generated := r.Result.(blueprints.Result)
	// Repeated machine and human requests must preserve the immutable retained revision.
	again := invoke("blueprint.generate", input, ActionCaller{})
	if again.ExitCode != 0 || again.Result.(blueprints.Result).RetainedRevision != generated.RetainedRevision {
		t.Fatal(again)
	}
	args := []string{filepath.Join(root, input.ModelArea), "model", "create", "blueprint", "from", input.From, "to", input.To, "--entry", input.Entry, "--task", filepath.Join(root, input.Task), "--catalogue", filepath.Join(root, input.Catalogue), "--json"}
	var out, errout bytes.Buffer
	if code := Run(ctx, args, &out, &errout); code != 0 {
		t.Fatalf("CLI %d %s", code, &errout)
	}
	var legacy blueprints.Result
	if e := json.Unmarshal(out.Bytes(), &legacy); e != nil || legacy.RetainedRevision != generated.RetainedRevision {
		t.Fatal("CLI retention mismatch", e)
	}
	var d blueprints.Document
	raw, e := os.ReadFile(filepath.Join(generated.Path, "blueprint.json"))
	if e != nil {
		t.Fatal(e)
	}
	if e = json.Unmarshal(raw, &d); e != nil {
		t.Fatal(e)
	}
	evidence := blueprints.Evidence{Schema: blueprints.EvidenceSchema, BlueprintRevision: d.Revision, RetainedRevision: generated.RetainedRevision, TaskDigest: d.TaskDigest, Scope: "Bounded action protocol trial", Code: map[string]blueprints.CodeSnapshot{}}
	writeActionJSON(t, filepath.Join(root, "evidence.json"), evidence)
	assessment := invoke("blueprint.assess", AssessActionInput{Bundle: generated.Path, Evidence: "evidence.json"}, ActionCaller{})
	if assessment.ExitCode != 3 || assessment.Status != "blocked" || assessment.Error != nil {
		t.Fatalf("blocked result lost: %+v", assessment)
	}
	out.Reset()
	errout.Reset()
	if code := Run(ctx, []string{"model", "assess", "blueprint", "--bundle", generated.Path, "--evidence", filepath.Join(root, "evidence.json"), "--json"}, &out, &errout); code != 3 || !json.Valid(out.Bytes()) {
		t.Fatalf("CLI assessment %d %s %s", code, &out, &errout)
	}
	both := input
	both.Output = "preview"
	if bad := invoke("blueprint.generate", both, ActionCaller{}); bad.ExitCode != 2 {
		t.Fatalf("two outputs accepted: %+v", bad)
	}
	if runtime.GOOS != "linux" {
		return
	}
	area := filepath.Join(root, "SDP")
	if e = os.MkdirAll(filepath.Join(area, "ProjectManagement"), 0700); e != nil {
		t.Fatal(e)
	}
	history := filepath.Join(area, "ProjectManagement/Ledger.ndjson")
	plan := `{"schemaVersion":"1.0","eventId":"EVT-PM-SDP-000001","eventType":"x-management:created","occurredAt":"2026-10-09T00:00:00Z","actor":"owner","commit":null,"subjectId":"PLAN-SDP-0001","payload":{"schemaVersion":"0.2","kind":"Plan","planType":"ImplementationPlan","previousEventId":null,"from":null,"to":"active","fromPath":null,"toPath":"05--Implementation/Test.md","reason":"fixture","links":[]}}` + "\n"
	if e = os.WriteFile(history, []byte(plan), 0600); e != nil {
		t.Fatal(e)
	}
	rel, e := filepath.Rel(area, generated.Path)
	if e != nil {
		t.Fatal(e)
	}
	binding := &blueprintstate.Binding{Bundle: filepath.ToSlash(rel), Task: d.TaskID, System: d.Analysis.System, Revision: d.Revision, Retained: generated.RetainedRevision, Plan: "PLAN-SDP-0001", Milestone: "M1", ModelArea: "Models", From: input.From, To: input.To}
	req := blueprintstate.Request{Schema: blueprintstate.Schema, EventID: "EVT-BPA-ACT-000001", OccurredAt: "2026-10-18T01:00:00Z", ID: "BPA-actions", Action: "create", Expected: "", Binding: binding, Reason: "Bounded invocation trial"}
	caller := ActionCaller{Principal: blueprintstate.Principal{Actor: "controller", Authority: "controller"}}
	if bad := invoke("blueprint.assignment.apply", ApplyAssignmentInput{Request: req}, ActionCaller{}); bad.ExitCode != 1 || bad.Error.Code != "authority-denied" {
		t.Fatalf("missing principal accepted: %+v", bad)
	}
	result := invoke("blueprint.assignment.apply", ApplyAssignmentInput{Request: req}, caller)
	if result.ExitCode != 0 {
		t.Fatalf("create: %+v", result)
	}
	before, _ := os.ReadFile(history)
	// Simulate a lost response, then use the established event-ID retry contract.
	if retry := invoke("blueprint.assignment.apply", ApplyAssignmentInput{Request: req}, caller); retry.ExitCode != 0 {
		t.Fatal(retry)
	}
	reqFile := filepath.Join(root, "assignment.json")
	writeActionJSON(t, reqFile, req)
	out.Reset()
	errout.Reset()
	if code := Run(ctx, []string{area, "model", "assignment", "apply", "--request", reqFile, "--as", "controller", "--authority", "controller", "--json"}, &out, &errout); code != 0 {
		t.Fatalf("CLI retry: %s", &errout)
	}
	unchanged, _ := os.ReadFile(history)
	if !bytes.Equal(before, unchanged) {
		t.Fatal("retry appended duplicate")
	}
	stale := req
	stale.Action = "cancel"
	stale.Binding = nil
	stale.EventID = "EVT-BPA-ACT-000002"
	if refused := invoke("blueprint.assignment.apply", ApplyAssignmentInput{Request: stale}, caller); refused.ExitCode == 0 {
		t.Fatal("stale transition accepted")
	}
	after, _ := os.ReadFile(history)
	if !bytes.Equal(before, after) {
		t.Fatal("refused write mutated history")
	}
	stale.Expected = req.EventID
	if good := invoke("blueprint.assignment.apply", ApplyAssignmentInput{Request: stale}, caller); good.ExitCode != 0 {
		t.Fatal(good)
	}
	list := invoke("blueprint.assignment.list", noActionParameters{}, ActionCaller{})
	if list.ExitCode != 0 {
		t.Fatal(list)
	}
	b, _ := json.Marshal(list.Result)
	if !strings.Contains(string(b), `"workState":"canceled"`) {
		t.Fatalf("%s", b)
	}
	// A forged principal nested inside data is never decoded into host authority.
	forged := map[string]any{"request": req, "principal": map[string]string{"authority": "controller"}}
	if bad := invoke("blueprint.assignment.apply", forged, caller); bad.ExitCode != 2 {
		t.Fatalf("forged field accepted: %+v", bad)
	}
	if rawAfter, _ := os.ReadFile(filepath.Join(generated.Path, "blueprint.json")); !bytes.Equal(raw, rawAfter) {
		t.Fatal("assignment changed retained bundle")
	}
}
