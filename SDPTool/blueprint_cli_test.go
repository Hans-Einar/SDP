package sdptool

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/Hans-Einar/SDP/SDPTool/blueprints"
	"github.com/Hans-Einar/SDP/SDPTool/model"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/blueprint"
)

func TestCompiledBlueprint(t *testing.T) {
	binary := filepath.Join(t.TempDir(), "sdptool")
	if b, e := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "build", "-o", binary, "./cmd/sdptool").CombinedOutput(); e != nil {
		t.Fatal(e, string(b))
	}
	area := t.TempDir()
	for _, s := range []struct{ name, fixture string }{{"Before", "NOW"}, {"After", "TARGET"}} {
		w, e := model.CreateWork(area, s.name, "", true)
		if e != nil {
			t.Fatal(e)
		}
		if e = os.CopyFS(w.Path, os.DirFS(filepath.Join("../experiments/blueprint_mvp1", s.fixture))); e != nil {
			t.Fatal(e)
		}
	}
	task := filepath.Join(t.TempDir(), "task.json")
	data, _ := json.Marshal(blueprint.Task{ID: "pilot", Intent: "Extract calibration"})
	if e := os.WriteFile(task, data, 0600); e != nil {
		t.Fatal(e)
	}
	out := filepath.Join(t.TempDir(), "blueprint")
	args := []string{area, "model", "create", "blueprint", "from", "work:Before", "to", "work:After", "--entry", "System.design", "--task", task, "--output", out}
	b, e := exec.Command(binary, append(args, "--json")...).CombinedOutput()
	if e != nil {
		t.Fatal(e, string(b))
	}
	var r map[string]any
	if e = json.Unmarshal(b, &r); e != nil || r["operation"] != "create-blueprint" {
		t.Fatal(e, string(b))
	}
	b, e = exec.Command(binary, args...).CombinedOutput()
	if e != nil || json.Valid(b) {
		t.Fatal("human output", e, string(b))
	}
	for _, name := range []string{"index.md", "context.md", "changes.md", "context.mmd", "obligations.md", "evidence.md", "blueprint.json", "assignment.json", "manifest.json", "sources/TARGET/System.design"} {
		if _, e = os.Stat(filepath.Join(out, name)); e != nil {
			t.Fatal(name, e)
		}
	}
	project := t.TempDir()
	catalogue := filepath.Join(project, "SDP", "Blueprints")
	retainedArgs := append([]string{}, args...)
	retainedArgs[len(retainedArgs)-2] = "--catalogue"
	retainedArgs[len(retainedArgs)-1] = catalogue
	b, e = exec.Command(binary, append(retainedArgs, "--json")...).CombinedOutput()
	if e != nil {
		t.Fatal(e, string(b))
	}
	if e = json.Unmarshal(b, &r); e != nil || r["retainedRevision"] == nil {
		t.Fatal(e, string(b))
	}

	retainedPath := r["path"].(string)
	raw, err := os.ReadFile(filepath.Join(retainedPath, "blueprint.json"))
	if err != nil {
		t.Fatal(err)
	}
	var doc blueprints.Document
	if err = json.Unmarshal(raw, &doc); err != nil {
		t.Fatal(err)
	}
	evidence := blueprints.Evidence{Schema: blueprints.EvidenceSchema, BlueprintRevision: doc.Revision, RetainedRevision: r["retainedRevision"].(string), TaskDigest: doc.TaskDigest, Scope: "compiled diagnostic fixture"}
	raw, _ = json.Marshal(evidence)
	evidencePath := filepath.Join(t.TempDir(), "evidence.json")
	os.WriteFile(evidencePath, raw, 0600)
	assessArgs := []string{"model", "assess", "blueprint", "--bundle", retainedPath, "--evidence", evidencePath}
	assessed, err := exec.Command(binary, append(assessArgs, "--json")...).CombinedOutput()
	exit, ok := err.(*exec.ExitError)
	if !ok || exit.ExitCode() != 3 {
		t.Fatal("blocked exit", err, string(assessed))
	}
	var assessment blueprints.Assessment
	if err = json.Unmarshal(assessed, &assessment); err != nil || assessment.Status != "blocked" {
		t.Fatal(err, string(assessed))
	}
	human, err := exec.Command(binary, assessArgs...).CombinedOutput()
	if err == nil || json.Valid(human) {
		t.Fatal("human assessment output", err, string(human))
	}
	b, e = exec.Command(binary, project, "discover", "--json").CombinedOutput()
	if e != nil {
		t.Fatal(e, string(b))
	}
	var discovery Project
	if e = json.Unmarshal(b, &discovery); e != nil {
		t.Fatal(e)
	}
	found := false
	for _, n := range discovery.Navigation.Nodes {
		if n.ID == "blueprints" && n.Kind == "tab" && n.State == "available" {
			found = true
		}
	}
	if !found {
		t.Fatal("compiled discovery missing tab")
	}
	if _, e = os.Stat(filepath.Join(area, ".git")); !os.IsNotExist(e) {
		t.Fatal("Git required")
	}
}
