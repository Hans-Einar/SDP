package codegen_test

import (
	"bytes"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/examples/application"
	model "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/examples/generatedmodel"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestCompiledExecutableParity(t *testing.T) {
	dir := t.TempDir()
	module, e := filepath.Abs("..")
	if e != nil {
		t.Fatal(e)
	}
	build := func(name string) string {
		exe := filepath.Join(dir, name)
		cmd := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "build", "-o", exe, "./cmd/"+name)
		cmd.Dir = module
		if b, e := cmd.CombinedOutput(); e != nil {
			t.Fatalf("build: %v %s", e, b)
		}
		return exe
	}
	source, compiled := build("sdl-simulate"), build("sdl-compiled")
	run := func(exe string, args ...string) []byte {
		cmd := exec.Command(exe, args...)
		cmd.Dir = dir
		b, e := cmd.CombinedOutput()
		if e != nil {
			t.Fatalf("run: %v %s", e, b)
		}
		return b
	}
	for _, values := range []string{"430,invalid,440", "100,1001,1000,0,600", "not-number"} {
		a := run(source, "-actions", filepath.Join(module, "examples/edit-apt-cell.sdl"), "-ui", filepath.Join(module, "examples/edit-apt-cell.sdui"), "-values", values)
		b := run(compiled, "-values", values)
		if !bytes.Equal(a, b) {
			t.Fatalf("event/state/links parity failure\n%s\n%s", a, b)
		}
	}
	// The generated executable runs from an empty cwd and takes no source paths.
	p := model.Program()
	for i := range p.Model.Statements {
		if p.Model.Statements[i].Verb == "invokes" {
			p.Model.Statements[i].Object.Name = "MissingImplementation"
		}
	}
	if app, e := application.New(p, model.Document(), model.Root()); e == nil {
		app.Close()
		t.Fatal("unimplemented Go action silently accepted")
	}
}
