package documents

import (
	"context"
	"fmt"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/viewpoint"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func model(t *testing.T) *viewpoint.Views {
	t.Helper()
	b, e := os.ReadFile("../../../SDUI/design/architecture.design")
	if e != nil {
		t.Fatal(e)
	}
	v, e := viewpoint.New(string(b))
	if e != nil {
		t.Fatal(e)
	}
	return v
}

type failRenderer struct{}

func (failRenderer) Identity() string { return "failure" }
func (failRenderer) Render(context.Context, string) ([]byte, error) {
	return nil, fmt.Errorf("render failure")
}
func TestNavigatorStaticAndPublication(t *testing.T) {
	v := model(t)
	ctx := context.Background()
	nav, e := Build(ctx, v, Options{Navigator: true, Renderer: failRenderer{}})
	if e != nil {
		t.Fatal(e)
	}
	if len(nav.Manifest.Diagrams) != 0 {
		t.Fatal("navigator rendered details")
	}
	for n := range nav.Files {
		if strings.HasSuffix(n, ".svg") || strings.HasSuffix(n, ".mmd") {
			t.Fatal(n)
		}
	}
	o := Options{Monolithic: true}
	a, e := Build(ctx, v, o)
	if e != nil {
		t.Fatal(e)
	}
	b, e := Build(ctx, v, o)
	if e != nil || !reflect.DeepEqual(a.Files, b.Files) {
		t.Fatal("not deterministic", e)
	}
	if len(a.Manifest.Diagrams) != 170 {
		t.Fatal(len(a.Manifest.Diagrams))
	}
	dir := filepath.Join(t.TempDir(), "views")
	if e = a.Publish(dir); e != nil {
		t.Fatal(e)
	}
	os.WriteFile(filepath.Join(dir, "notes.md"), []byte("private notes"), 0600)
	if e = nav.Publish(dir); e != nil {
		t.Fatal(e)
	}
	if _, e = os.Stat(filepath.Join(dir, "diagrams/VP02-roots.mmd")); !os.IsNotExist(e) {
		t.Fatal("stale diagram")
	}
	n, _ := os.ReadFile(filepath.Join(dir, "notes.md"))
	if string(n) != "private notes" {
		t.Fatal("notes lost")
	}
	before, _ := os.ReadFile(filepath.Join(dir, "manifest.json"))
	_, e = Build(ctx, v, Options{Renderer: failRenderer{}})
	if e == nil {
		t.Fatal("renderer failure ignored")
	}
	after, _ := os.ReadFile(filepath.Join(dir, "manifest.json"))
	if string(before) != string(after) {
		t.Fatal("failed render changed publication")
	}
	os.WriteFile(filepath.Join(dir, "index.md"), []byte("edited"), 0600)
	if e = nav.Publish(dir); e == nil {
		t.Fatal("overwrote local edit")
	}
}
func TestUnsafePathsAndSVG(t *testing.T) {
	for _, s := range []string{"../escape", "/tmp/out", "a/../../b"} {
		if safe(s) {
			t.Fatal(s)
		}
	}
	for _, s := range []string{"", `<svg xmlns="http://www.w3.org/2000/svg"><script/></svg>`, `<svg xmlns="http://www.w3.org/2000/svg"/><svg xmlns="http://www.w3.org/2000/svg"/>`} {
		if ValidateSVG([]byte(s)) == nil {
			t.Fatal(s)
		}
	}
}
