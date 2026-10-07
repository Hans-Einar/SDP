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
)

func setup(t *testing.T) (Options, model.Result) {
	t.Helper()
	area := t.TempDir()
	work, e := model.CreateWork(area, "Before", "", true)
	if e != nil {
		t.Fatal(e)
	}
	put := func(dir, name string) {
		root := filepath.Join("../../experiments/blueprint_mvp1", name)
		e := os.CopyFS(dir, os.DirFS(root))
		if e != nil {
			t.Fatal(e)
		}
	}
	put(work.Path, "NOW")
	if _, e = model.Commit(area, "work:Before", "baseline", false); e != nil {
		t.Fatal(e)
	}
	if _, e = model.Freeze(area, "candidate", "Base", "work:Before", ""); e != nil {
		t.Fatal(e)
	}
	if _, e = model.Freeze(area, "release", "0.1.0", "candidate:Base", "model-only"); e != nil {
		t.Fatal(e)
	}
	next, e := model.CreateWork(area, "After", "release:0.1.0", false)
	if e != nil {
		t.Fatal(e)
	}
	b, e := os.ReadFile("../../experiments/blueprint_mvp1/TARGET/Containers/MachineService.design")
	if e != nil {
		t.Fatal(e)
	}
	if e = os.WriteFile(filepath.Join(next.Path, "Containers/MachineService.design"), b, 0600); e != nil {
		t.Fatal(e)
	}
	task := filepath.Join(t.TempDir(), "task.json")
	data, _ := json.Marshal(blueprint.Task{ID: "calibrate", Intent: "<script> [bad](outside) `ticks`"})
	if e = os.WriteFile(task, data, 0600); e != nil {
		t.Fatal(e)
	}
	return Options{area, "release:0.1.0", "work:After", "System.design", task, filepath.Join(t.TempDir(), "bundle")}, next
}
func read(t *testing.T, p string) []byte {
	t.Helper()
	b, e := os.ReadFile(p)
	if e != nil {
		t.Fatal(e)
	}
	return b
}
func TestGenerationAndFreshness(t *testing.T) {
	o, next := setup(t)
	ctx := context.Background()
	r, e := Generate(ctx, o)
	if e != nil {
		t.Fatal(e)
	}
	if !r.Preliminary || r.Status != "diagnostic-preview" {
		t.Fatal(r)
	}
	first := read(t, filepath.Join(o.Output, "manifest.json"))
	if _, e = Generate(ctx, o); e != nil {
		t.Fatal(e)
	}
	if string(first) != string(read(t, filepath.Join(o.Output, "manifest.json"))) {
		t.Fatal("not deterministic")
	}
	if strings.Contains(string(read(t, filepath.Join(o.Output, "index.md"))), "<script>") {
		t.Fatal("unescaped author text")
	}
	if e = os.WriteFile(filepath.Join(o.Output, "notes.txt"), []byte("keep"), 0600); e != nil {
		t.Fatal(e)
	}
	if _, e = Generate(ctx, o); e != nil {
		t.Fatal(e)
	}
	if string(read(t, filepath.Join(o.Output, "notes.txt"))) != "keep" {
		t.Fatal("notes lost")
	}
	_, e = generate(ctx, o, func() {
		p := filepath.Join(next.Path, "Containers/MachineService.design")
		b := read(t, p)
		os.WriteFile(p, append(b, '\n'), 0600)
	})
	if e == nil || !strings.Contains(e.Error(), "stale") {
		t.Fatal("expected stale", e)
	}
	if string(first) != string(read(t, filepath.Join(o.Output, "manifest.json"))) {
		t.Fatal("stale replaced output")
	}
	if e = os.WriteFile(filepath.Join(o.Output, "index.md"), []byte("edited"), 0600); e != nil {
		t.Fatal(e)
	}
	if _, e = Generate(ctx, o); e == nil {
		t.Fatal("overwrote edited generated file")
	}
	if string(read(t, filepath.Join(o.Output, "index.md"))) != "edited" {
		t.Fatal("edit lost")
	}
}
func TestCancelTaskAndOverlap(t *testing.T) {
	o, next := setup(t)
	ctx, cancel := context.WithCancel(context.Background())
	if _, e := generate(ctx, o, cancel); e != context.Canceled {
		t.Fatal(e)
	}
	if _, e := os.Stat(o.Output); !os.IsNotExist(e) {
		t.Fatal("canceled output exists")
	}
	_, e := generate(context.Background(), o, func() { os.WriteFile(o.Task, []byte("{}"), 0600) })
	if e == nil || !strings.Contains(e.Error(), "stale") {
		t.Fatal(e)
	}
	o, _ = setup(t)
	o.Output = filepath.Join(next.Path, "generated")
	// Use the current actual TARGET path for overlap.
	v, e := model.Snapshot(o.Area, o.To)
	if e != nil {
		t.Fatal(e)
	}
	o.Output = filepath.Join(v.Path, "generated")
	if _, e = Generate(context.Background(), o); e == nil {
		t.Fatal("output inside source")
	}
}
func TestInvalidAndSymlink(t *testing.T) {
	o, next := setup(t)
	p := filepath.Join(next.Path, "System.design")
	os.WriteFile(p, []byte("invalid"), 0600)
	if _, e := Generate(context.Background(), o); e == nil {
		t.Fatal("invalid sources accepted")
	}
	o, _ = setup(t)
	parent := t.TempDir()
	link := filepath.Join(t.TempDir(), "link")
	if e := os.Symlink(parent, link); e != nil {
		t.Skip(e)
	}
	o.Output = filepath.Join(link, "bundle")
	if _, e := Generate(context.Background(), o); e == nil {
		t.Fatal("symlink parent accepted")
	}
}

func TestSymlinkAreaCannotWriteIntoSource(t *testing.T) {
	o, target := setup(t)
	alias := filepath.Join(t.TempDir(), "alias")
	if e := os.Symlink(o.Area, alias); e != nil {
		t.Skip(e)
	}
	o.Area = alias
	o.Output = filepath.Join(target.Path, "generated")
	if _, e := Generate(context.Background(), o); e == nil {
		t.Fatal("aliased source overwrite allowed")
	}
	if _, e := os.Stat(o.Output); !os.IsNotExist(e) {
		t.Fatal("source output created")
	}
}
func TestEntryAndCompilerProvenance(t *testing.T) {
	o, _ := setup(t)
	r, e := Generate(context.Background(), o)
	if e != nil {
		t.Fatal(e)
	}
	var d Document
	if e = json.Unmarshal(read(t, filepath.Join(o.Output, "blueprint.json")), &d); e != nil {
		t.Fatal(e)
	}
	if d.Entry != "System.design" || len(d.Compiler.ExecutableSHA256) != 64 || d.Compiler.SourceGraph == "" || d.Compiler.Profile != "design-core/0.6" || d.Revision != r.Revision {
		t.Fatalf("%+v", d)
	}
}
