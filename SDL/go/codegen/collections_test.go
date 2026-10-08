package codegen_test

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/codegen"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/documents"
)

func TestCollectionBundleProvenanceAndConstructors(t *testing.T) {
	actions, err := os.ReadFile("../examples/edit-apt-cell.sdl")
	if err != nil {
		t.Fatal(err)
	}
	const source = `sdui 0.3; page=[items=tree("Navigation")];`
	b, err := codegen.Bundle(string(actions), source, "page", "generated")
	if err != nil {
		t.Fatal(err)
	}
	const version = "sdl-go-model/1/sdui-go-model/2"
	if b.Manifest.Version != version {
		t.Fatal("incorrect in-memory bundle provenance", b.Manifest.Version)
	}
	var manifest documents.Manifest
	if err := json.Unmarshal(b.Files["manifest.json"], &manifest); err != nil || manifest.Version != version {
		t.Fatal("incorrect published bundle provenance", manifest.Version, err)
	}
	dir := t.TempDir()
	for name, data := range b.Files {
		if err := os.WriteFile(filepath.Join(dir, name), data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	sdl, err := filepath.Abs("..")
	if err != nil {
		t.Fatal(err)
	}
	sdui := filepath.Clean(filepath.Join(sdl, "..", "..", "SDUI", "go"))
	mod := fmt.Sprintf("module bundled-collections-test\n\ngo 1.26.0\nrequire (\n github.com/Hans-Einar/SDP/SDUI/go v0.0.0\n github.com/Hans-Einar/SDP/SystemDesignLanguage/go v0.0.0\n)\nreplace github.com/Hans-Einar/SDP/SDUI/go => %q\nreplace github.com/Hans-Einar/SDP/SystemDesignLanguage/go => %q\n", sdui, sdl)
	harness := fmt.Sprintf(`package generated
import (
 "reflect"
 "testing"
 ui "github.com/Hans-Einar/SDP/SDUI/go/parser"
 sdl "github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
)
func TestBundle(t *testing.T) {
 if UIProfile != "sdui/0.3" || UIASTFormat != "sdui-ast/0.3" { t.Fatal("metadata") }
 doc, roots, err := ui.Compile(%q)
 if err != nil { t.Fatal(err) }
 if !reflect.DeepEqual(doc,Document()) || !reflect.DeepEqual(roots["page"],Root()) { t.Fatal("UI constructor mismatch") }
 p,err:=sdl.CompileActions(%q)
 if err!=nil || !reflect.DeepEqual(p,Program()) {t.Fatal("SDL constructor mismatch",err)}
 if Root().Profile!=UIProfile || UISourceSHA256!=%q || SDLSourceSHA256!=%q {t.Fatal("source identity lost")}
}
`, source, string(actions), documents.Hash([]byte(source)), documents.Hash(actions))
	for name, content := range map[string]string{"go.mod": mod, "bundle_test.go": harness} {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0600); err != nil {
			t.Fatal(err)
		}
	}
	cmd := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "test", "-mod=mod", "-count=1", "./...")
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOWORK=off")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("generated bundle test: %v\n%s", err, output)
	}
}
