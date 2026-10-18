package blueprints

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Hans-Einar/SDP/SDPTool/model"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/blueprint"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/documents"
)

func assessmentFixture(t *testing.T) (string, string, Evidence) {
	t.Helper()
	o, _ := setup(t)
	o.Output = ""
	task := blueprint.Task{ID: "trial", Intent: "Extract calibration in bounded fixture", AllowedChanges: []blueprint.Permission{
		{Element: "MVP1InspectionPilot/unit/CalibrationProcessor", Change: "added"},
		{Fact: "MachineService owns CalibrateMeasurement.", Change: "removed"},
		{Fact: "CalibrationProcessor owns CalibrateMeasurement.", Change: "added"},
		{Fact: "MachineService contains CalibrationProcessor.", Change: "added"},
	}, Rules: []blueprint.Rule{{ID: "preserve", Source: "trial", Fact: "MachineService owns ReduceMachineState.", Now: true, Target: true}}}
	raw, _ := json.Marshal(task)
	if e := os.WriteFile(o.Task, raw, 0600); e != nil {
		t.Fatal(e)
	}
	r, e := GenerateRetained(context.Background(), o, filepath.Join(t.TempDir(), "Blueprints"))
	if e != nil {
		t.Fatal(e)
	}
	d, _, e := verify(r.Path, &budget{})
	if e != nil {
		t.Fatal(e)
	}
	ep := filepath.Join(t.TempDir(), "evidence.json")
	in := Evidence{Schema: EvidenceSchema, BlueprintRevision: d.Revision, RetainedRevision: r.RetainedRevision, TaskDigest: d.TaskDigest, Scope: "temporary model-only trial", Code: map[string]CodeSnapshot{}}
	return r.Path, ep, in
}
func writeEvidence(t *testing.T, p string, in Evidence) {
	t.Helper()
	b, e := json.Marshal(in)
	if e != nil {
		t.Fatal(e)
	}
	if e = os.WriteFile(p, b, 0600); e != nil {
		t.Fatal(e)
	}
}
func deferUnknowns(t *testing.T, bundle, ep string, in *Evidence) Assessment {
	t.Helper()
	in.Dispositions = nil
	writeEvidence(t, ep, *in)
	a, e := Assess(bundle, ep)
	if e != nil {
		t.Fatal(e)
	}
	for _, u := range a.Unknowns {
		in.Dispositions = append(in.Dispositions, Disposition{Unknown: u.ID, Actor: "test-reviewer", Role: "reviewer", Scope: in.Scope, Rationale: "Defer within synthetic model trial; no production permission", Traceability: "TEST-TRIAL"})
	}
	writeEvidence(t, ep, *in)
	a, e = Assess(bundle, ep)
	if e != nil {
		t.Fatal(e)
	}
	return a
}
func TestAssessmentReadinessAndPins(t *testing.T) {
	bundle, ep, in := assessmentFixture(t)
	writeEvidence(t, ep, in)
	a, e := Assess(bundle, ep)
	if e != nil || a.Status != "blocked" || len(a.Unknowns) == 0 {
		t.Fatal(a, e)
	}
	ready := deferUnknowns(t, bundle, ep, &in)
	if ready.Status != "ready" {
		t.Fatal(ready.Blockers)
	}
	again, e := Assess(bundle, ep)
	if e != nil || again.Identity != ready.Identity {
		t.Fatal(e, "non deterministic")
	}
	in.Scope = "different"
	writeEvidence(t, ep, in)
	if _, e = Assess(bundle, ep); e == nil {
		t.Fatal("wrong scope allowed")
	}
	in.Scope = "temporary model-only trial"
	in.BlueprintRevision = strings.Repeat("0", 64)
	writeEvidence(t, ep, in)
	if _, e = Assess(bundle, ep); e == nil {
		t.Fatal("stale input allowed")
	}
}
func TestAssessmentReceiptsAndCodeChanges(t *testing.T) {
	bundle, ep, in := assessmentFixture(t)
	data := []byte("package pilot\nfunc Reduce(a,b int) int {return a+b}\n")
	root := filepath.Dir(ep)
	os.Mkdir(filepath.Join(root, "code"), 0700)
	os.WriteFile(filepath.Join(root, "code/machine.go"), data, 0600)
	files := model.Files{"machine.go": data}
	digest := model.Digest(files)
	in.Code["TARGET"] = CodeSnapshot{Revision: "synthetic-target", Root: "code", Digest: digest, Files: map[string]string{"machine.go": documents.Hash(data)}}
	d, _, _ := verify(bundle, &budget{})
	in.RequiredChecks = []CheckRequirement{{ID: "behavior", Side: "TARGET", Obligation: "behavior:reduction"}}
	log := []byte("reported test passed\n")
	os.WriteFile(filepath.Join(root, "test.txt"), log, 0600)
	r := CheckReceipt{ID: "behavior", Side: "TARGET", Obligation: "behavior:reduction", SourceDigest: d.Target.SourceDigest, TaskDigest: d.TaskDigest, CodeDigest: digest, Command: "NEVER EXECUTE THIS", Environment: "synthetic", InputsDigest: model.Digest(model.Files{"fixture.txt": []byte("fixture")}), Inputs: map[string]string{"fixture.txt": documents.Hash([]byte("fixture"))}, Result: "pass", Artifact: "test.txt", ArtifactDigest: documents.Hash(log), Traceability: "TEST-1"}
	os.WriteFile(filepath.Join(root, "fixture.txt"), []byte("fixture"), 0600)
	in.Receipts = []CheckReceipt{r}
	writeEvidence(t, ep, in)
	ambiguous := strings.Replace(string(read(t, ep)), `"result":"pass"`, `"Result":"fail","result":"pass"`, 1)
	os.WriteFile(ep, []byte(ambiguous), 0600)
	if _, err := Assess(bundle, ep); err == nil {
		t.Fatal("case-aliased result override accepted")
	}

	in.Mappings = []Mapping{{Side: "TARGET", Element: "MVP1InspectionPilot/container/MachineService", Role: "implementation", Path: "machine.go", Symbol: "Reduce", FileDigest: documents.Hash(data)}}
	a := deferUnknowns(t, bundle, ep, &in)
	if a.Status != "ready" || a.Checks[0].Effective != "pass" {
		t.Fatal(a)
	}
	original := a.Identity
	if _, err := assess(bundle, ep, func() { os.WriteFile(filepath.Join(root, "fixture.txt"), []byte("edited while assessing"), 0600) }); err == nil {
		t.Fatal("mid-assessment edit not rejected")
	}
	os.WriteFile(filepath.Join(root, "fixture.txt"), []byte("fixture"), 0600)

	// A code-only change leaves SDL unchanged but stales the receipt and mapping.
	os.WriteFile(filepath.Join(root, "code/machine.go"), []byte("changed code"), 0600)
	in.Dispositions = nil
	writeEvidence(t, ep, in)
	a, e := Assess(bundle, ep)
	if e != nil || a.Status != "blocked" || a.Checks[0].Effective != "stale" || a.Identity == original {
		t.Fatal(a, e)
	}
	a = deferUnknowns(t, bundle, ep, &in)
	if a.Checks[0].Effective != "stale" {
		t.Fatal("disposition turned stale to pass")
	}
	// A recorded failure cannot be waived by an unknown disposition.
	in.Receipts[0].Result = "fail"
	a = deferUnknowns(t, bundle, ep, &in)
	if a.Status != "blocked" {
		t.Fatal("failure waived")
	}
	in.Receipts[0] = r
	os.WriteFile(filepath.Join(root, "code/machine.go"), data, 0600)
	os.WriteFile(filepath.Join(root, "test.txt"), []byte("modified"), 0600)
	a = deferUnknowns(t, bundle, ep, &in)
	if a.Checks[0].Effective != "stale" {
		t.Fatal("evidence hash ignored")
	}
	in.Receipts[0].Side = "NOW"
	writeEvidence(t, ep, in)
	if _, e = Assess(bundle, ep); e == nil {
		t.Fatal("NOW receipt reused for TARGET")
	}
}
func TestAssessmentUnknownMappingsAndMalformedInput(t *testing.T) {
	bundle, ep, in := assessmentFixture(t)
	in.Mappings = []Mapping{{Side: "TARGET", Element: "MVP1InspectionPilot/container/MachineService", Role: "implementation", Path: "missing.go", Symbol: "Machine", FileDigest: strings.Repeat("a", 64)}}
	in.Mappings = append(in.Mappings, in.Mappings[0])
	a := deferUnknowns(t, bundle, ep, &in)
	found := false
	for _, u := range a.Unknowns {
		if strings.HasPrefix(u.ID, "mapping/") && strings.Contains(u.Detail, "locator") {
			found = true
		}
	}
	if !found {
		t.Fatal("missing/ambiguous locator hidden")
	}
	in.RequiredChecks = []CheckRequirement{{ID: "missing", Side: "NOW", Obligation: "behavior:baseline"}}
	a = deferUnknowns(t, bundle, ep, &in)
	if a.Checks[0].Effective != "not-run" {
		t.Fatal(a)
	}
	in.Dispositions = append(in.Dispositions, in.Dispositions[0])
	writeEvidence(t, ep, in)
	if _, e := Assess(bundle, ep); e == nil {
		t.Fatal("duplicate disposition")
	}
	writeEvidence(t, ep, in)
	data, _ := os.ReadFile(ep)
	data = append([]byte(`{"schema":"bad",`), data[1:]...)
	os.WriteFile(ep, data, 0600)
	if _, e := Assess(bundle, ep); e == nil {
		t.Fatal("duplicate JSON key")
	}
	in.Dispositions = nil
	in.Code["NOW"] = CodeSnapshot{Revision: "x", Root: "../escape", Digest: strings.Repeat("a", 64), Files: map[string]string{"x": strings.Repeat("b", 64)}}
	writeEvidence(t, ep, in)
	if _, e := Assess(bundle, ep); e == nil {
		t.Fatal("path escape")
	}
}
