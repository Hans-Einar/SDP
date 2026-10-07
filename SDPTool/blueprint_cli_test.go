package sdptool

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

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
	if _, e = os.Stat(filepath.Join(area, ".git")); !os.IsNotExist(e) {
		t.Fatal("Git required")
	}
}
