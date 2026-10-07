package codegen

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

// Compile and execute generated constructors, so dropping provenance fields or
// sharing their slices cannot pass through a text-only assertion.
func TestGeneratedReuseProvenance(t *testing.T) {
	const source = `sdui 0.2;
leaf = <ok=button("OK")>;
alias = leaf;
page = [left=alias, right=alias];
`
	generated, err := Generate(source, "page", "generated")
	if err != nil {
		t.Fatal(err)
	}
	module, err := filepath.Abs("..")
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	harness := fmt.Sprintf(`package generated
import (
 "reflect"
 "testing"
 ui "github.com/Hans-Einar/SDP/SDUI/go/parser"
)
func TestRoundTrip(t *testing.T) {
 doc, roots, err := ui.Compile(%q)
 if err != nil { t.Fatal(err) }
 root := Root()
 if !reflect.DeepEqual(doc, Document()) || !reflect.DeepEqual(roots["page"], root) {
  t.Fatal("generated constructors lost normalized model or source provenance")
 }
 if root.Declaration != "page" { t.Fatal("root declaration lost") }
 left, right := root.Rows[0][0], root.Rows[0][1]
 for _, node := range []*ui.Instance{left, right} {
  if node.Declaration != "alias" || len(node.Uses) != 2 ||
   node.Uses[0].Definition != "leaf" || node.Uses[1].Definition != "alias" {
   t.Fatalf("reuse chain lost: %%+v", node)
  }
  if node.Span != doc.Definitions[0].Root.Span ||
   node.Uses[0].Span != doc.Definitions[1].Root.Span {
   t.Fatal("declaration or intermediate use source span lost")
  }
 }
 if left.Path != "page/left" || right.Path != "page/right" ||
  left.Uses[1].Span != doc.Definitions[2].Root.Rows[0].Items[0].Span ||
  right.Uses[1].Span != doc.Definitions[2].Root.Rows[0].Items[1].Span ||
  left.Uses[1].Span == right.Uses[1].Span {
  t.Fatal("distinct use sites collapsed")
 }
 left.Uses[0].Definition = "mutated"
 left.Uses[1].Span.Line = -1
 if right.Uses[0].Definition != "leaf" || right.Uses[1].Span.Line < 1 ||
  !reflect.DeepEqual(roots["page"], Root()) {
  t.Fatal("constructors or sibling instances share mutable provenance")
 }
}
`, source)
	files := map[string][]byte{
		"model_gen.go":       generated,
		"provenance_test.go": []byte(harness),
		"go.mod":             []byte(fmt.Sprintf("module provenance-test\n\ngo 1.26.0\n\nrequire github.com/Hans-Einar/SDP/SDUI/go v0.0.0\nreplace github.com/Hans-Einar/SDP/SDUI/go => %q\n", module)),
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), content, 0600); err != nil {
			t.Fatal(err)
		}
	}
	cmd := exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), "test", "-count=1", "./...")
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOWORK=off")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("generated provenance round trip: %v\n%s", err, output)
	}
}
