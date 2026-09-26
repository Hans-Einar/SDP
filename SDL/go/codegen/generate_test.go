package codegen_test

import (
	"bytes"
	ui "github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/codegen"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/examples/application"
	model "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/examples/generatedmodel"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	"os"
	"reflect"
	"testing"
)

func TestGeneratedConstructors(t *testing.T) {
	a, e := os.ReadFile("../examples/edit-apt-cell.sdl")
	if e != nil {
		t.Fatal(e)
	}
	s, e := os.ReadFile("../examples/edit-apt-cell.sdui")
	if e != nil {
		t.Fatal(e)
	}
	p, e := parser.CompileActions(string(a))
	if e != nil {
		t.Fatal(e)
	}
	d, r, e := ui.Compile(string(s))
	if e != nil {
		t.Fatal(e)
	}
	if !reflect.DeepEqual(p, model.Program()) || !reflect.DeepEqual(d, model.Document()) || !reflect.DeepEqual(r["page"], model.Root()) {
		t.Fatal("generated constructor differs from frontend")
	}
	root := model.Root()
	root.Layout["x"] = "changed"
	doc := model.Document()
	doc.References[0].Alias = "changed"
	program := model.Program()
	delete(program.Actions, "EditAptCell")
	if !reflect.DeepEqual(p, model.Program()) || !reflect.DeepEqual(d, model.Document()) || !reflect.DeepEqual(r["page"], model.Root()) {
		t.Fatal("constructor shared mutable data")
	}
	app, e := application.New(model.Program(), model.Document(), model.Root())
	if e != nil {
		t.Fatal(e)
	}
	defer app.Close()
	if len(app.Bindings.Links) == 0 {
		t.Fatal("missing bridge provenance")
	}
	b, e := codegen.Bundle(string(a), string(s), "page", "generatedmodel")
	if e != nil {
		t.Fatal(e)
	}
	for name, v := range b.Files {
		old, e := os.ReadFile("../examples/generatedmodel/" + name)
		if e != nil || !bytes.Equal(v, old) {
			t.Fatal("generated example stale", name, e)
		}
	}
	for _, f := range []func() error{
		func() error { _, e := codegen.Bundle(string(a), string(s), "absent", "model"); return e },
		func() error { _, e := codegen.Generate("language design-core version 0.5.\n", "model"); return e },
		func() error { _, e := codegen.Generate(string(a), "package"); return e },
	} {
		if f() == nil {
			t.Fatal("unsupported profile/entry/package accepted")
		}
	}
}
