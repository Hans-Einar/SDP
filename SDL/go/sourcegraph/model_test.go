package sourcegraph

import (
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const header = "language design-core version 0.6.\n"

func write(t *testing.T, root, name, text string) {
	t.Helper()
	p := filepath.Join(root, name)
	if e := os.MkdirAll(filepath.Dir(p), 0755); e != nil {
		t.Fatal(e)
	}
	if e := os.WriteFile(p, []byte(text), 0644); e != nil {
		t.Fatal(e)
	}
}
func fixture(t *testing.T) (string, *parser.SyntaxCache) {
	t.Helper()
	r := t.TempDir()
	write(t, r, "System.design", header+"system Demo.\nincludes \"Shared.design\".\nDemo contains Containers/Machine.\n")
	write(t, r, "Containers/Machine.design", header+"container Machine.\nincludes \"Shared.design\".\n")
	write(t, r, "Shared.design", header+"contract Shared.\nincludes \"System.design\".\nShared has completeness = open.\n")
	return r, &parser.SyntaxCache{}
}
func TestLateRootForestAndCycles(t *testing.T) {
	r, c := fixture(t)
	b, _ := os.ReadFile(filepath.Join(r, "Containers/Machine.design"))
	f, e := c.Parse("Containers/Machine.design", string(b))
	if e != nil {
		t.Fatal(e)
	}
	if f.Inspect()["validated"] != false {
		t.Fatal("fragment claimed validated")
	}
	f.Model().Declarations[0].Name.Name = "Mutated"
	l, e := Load(filepath.Join(r, "System.design"), c, true)
	if e != nil {
		t.Fatal(e)
	}
	hits, _ := c.Stats()
	if hits < 1 || l.Snapshot.FileCount() != 3 || l.Snapshot.System() != "Demo" {
		t.Fatal(hits, l.Snapshot.AST())
	}
	cold, e := Load(filepath.Join(r, "System.design"), nil, true)
	if e != nil || cold.Snapshot.Revision() != l.Snapshot.Revision() {
		t.Fatal(e)
	}
	if len(l.Snapshot.Edges()) != 4 {
		t.Fatal("lost repeated edges")
	}
	if l.Snapshot.Model().Statements[0].Object.Name != "Machine" {
		t.Fatal("syntax mutated across contexts")
	}
	for _, s := range l.Snapshot.Model().Statements {
		if s.Verb == "contains" && s.Span.Source != "System.design" {
			t.Fatal("lost source")
		}
	}
	if e = l.Fresh(); e != nil {
		t.Fatal(e)
	}
	write(t, r, "Unlisted.design", header+"system Other.\n")
	if e = l.Fresh(); e != nil {
		t.Fatal("unlisted changed revision", e)
	}
	write(t, r, "Shared.design", header+"contract Changed.\n")
	if e = l.Fresh(); e == nil {
		t.Fatal("non-entry edit was ignored")
	}
}
func TestContextIsolation(t *testing.T) {
	c := &parser.SyntaxCache{}
	child, _ := c.Parse("Child.design", header+"unit Child.\n")
	a, _ := c.Parse("A.design", header+"system A.\nincludes \"Child.design\".\nA contains Child.\n")
	b, _ := c.Parse("B.design", header+"system B.\nincludes \"Child.design\".\nB contains Child.\n")
	for _, test := range []struct {
		entry, system string
		root          *parser.File
	}{{"A.design", "A", a}, {"B.design", "B", b}} {
		s, ds := Compile(test.entry, []*parser.File{test.root, child}, true)
		if len(ds) > 0 || s.System() != test.system {
			t.Fatal(ds)
		}
	}
	bad, _ := c.Parse("A.design", header+"system A.\nunit Child.\nincludes \"Child.design\".\n")
	_, ds := Compile("A.design", []*parser.File{bad, child}, false)
	if len(ds) == 0 || len(ds[0].Related) == 0 {
		t.Fatal("duplicate lost provenance", ds)
	}
	missing, _ := c.Parse("A.design", header+"system A.\nA contains Child.\n")
	_, ds = Compile("A.design", []*parser.File{missing, child}, false)
	if len(ds) == 0 {
		t.Fatal("unreachable cache satisfied symbol")
	}
}
func TestNegativeInputs(t *testing.T) {
	tests := []struct{ name, file, content, code string }{
		{"duplicate", "Containers/Machine.design", header + "container Machine.\ncontract Shared.\n", "DUPLICATE_DECLARATION"},
		{"target", "Containers/Machine.design", header + "container Renamed.\n", "PATH_TARGET"},
		{"mixed", "Shared.design", "language design-core version 0.5.\ncontract Shared.\n", "PROFILE_MISMATCH"},
		{"second-system", "Shared.design", header + "system Other.\n", "SYSTEM"},
		{"missing", "System.design", header + "system Demo.\nincludes \"Missing.design\".\n", "SOURCE_READ"},
		{"escape", "System.design", header + "system Demo.\nincludes \"../Outside.design\".\n", "SOURCE_PATH"},
		{"cycle", "Shared.design", header + "unit A.\nunit B.\nA contains B.\nB contains A.\n", "STRUCTURE_CYCLE"},
		{"unit-kind", "Shared.design", header + "unit A.\ncontract Shared.\nA contains Shared.\n", "TYPE_MISMATCH"},
		{"missing-system", "System.design", header + "includes \"Shared.design\".\n", "SYSTEM_CARDINALITY"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r, _ := fixture(t)
			write(t, r, tc.file, tc.content)
			_, e := Load(filepath.Join(r, "System.design"), nil, false)
			if e == nil || !strings.Contains(e.Error(), tc.code) {
				t.Fatalf("want %s got %v", tc.code, e)
			}
		})
	}
}
func TestPathsFormattingAndPublication(t *testing.T) {
	r, _ := fixture(t)
	l, e := Load(filepath.Join(r, "System.design"), nil, true)
	if e != nil {
		t.Fatal(e)
	}
	for _, p := range []string{r, filepath.Join(r, "Containers"), filepath.Join(r, "Shared.design")} {
		if l.Outside(p) == nil {
			t.Fatal("output aliases source", p)
		}
	}
	if e = l.Outside(filepath.Join(r, "output")); e != nil {
		t.Fatal(e)
	}
	other := t.TempDir()
	for n, s := range l.Snapshot.Format() {
		write(t, other, n, s)
	}
	copy, e := Load(filepath.Join(other, "System.design"), nil, true)
	if e != nil || copy.Snapshot.Revision() != l.Snapshot.Revision() {
		t.Fatal("relocation", e)
	}
	os.Remove(filepath.Join(r, "Shared.design"))
	if e = os.Symlink(filepath.Join(other, "Shared.design"), filepath.Join(r, "Shared.design")); e != nil {
		t.Fatal(e)
	}
	if _, e = Load(filepath.Join(r, "System.design"), nil, true); e == nil {
		t.Fatal("symlink allowed")
	}
}

func TestHardlinkMissingAndAggregateLimits(t *testing.T) {
	r, _ := fixture(t)
	os.Remove(filepath.Join(r, "Shared.design"))
	if e := os.Link(filepath.Join(r, "Containers/Machine.design"), filepath.Join(r, "Shared.design")); e != nil {
		t.Fatal(e)
	}
	if _, e := Load(filepath.Join(r, "System.design"), nil, true); e == nil || !strings.Contains(e.Error(), "SOURCE_ALIAS") {
		t.Fatal(e)
	}
	os.Remove(filepath.Join(r, "Shared.design"))
	os.Mkdir(filepath.Join(r, "Shared.design"), 0755)
	if _, e := Load(filepath.Join(r, "System.design"), nil, true); e == nil {
		t.Fatal("directory source")
	}
	os.Remove(filepath.Join(r, "Shared.design"))
	write(t, r, "Shared.design", header+strings.Repeat(" ", parser.MaxBytes))
	if _, e := Load(filepath.Join(r, "System.design"), nil, false); e == nil || !strings.Contains(e.Error(), "SOURCE_LIMIT") {
		t.Fatal(e)
	}
}

func TestPureCompileLimitsAndCrossFileDiagnostics(t *testing.T) {
	c := &parser.SyntaxCache{}
	parse := func(name, text string) *parser.File {
		t.Helper()
		f, e := c.Parse(name, header+text)
		if e != nil {
			t.Fatal(e)
		}
		return f
	}
	a := parse("A.design", "system Demo.\n"+strings.Repeat("includes \"B.design\".\n", 42000))
	b := parse("B.design", strings.Repeat("includes \"A.design\".\n", 42000))
	if s, ds := Compile("A.design", []*parser.File{a, b}, false); s != nil || len(ds) == 0 || ds[0].Code != "TOKEN_LIMIT" {
		t.Fatal("pure API bypassed aggregate token budget", ds)
	}
	a = parse("A.design", "unit A.\nunit B.\nsystem Demo.\nincludes \"B.design\".\nA contains B.\n")
	b = parse("B.design", "B contains A.\n")
	_, ds := Compile("A.design", []*parser.File{a, b}, false)
	found := false
	for _, d := range ds {
		if d.Code == "STRUCTURE_CYCLE" {
			for _, r := range d.Related {
				if r.Source != d.Span.Source {
					found = true
				}
			}
		}
	}
	if !found {
		t.Fatal("cycle lost cross-file evidence", ds)
	}
	a = parse("A.design", "unit A.\nunit B.\nsystem Demo.\nincludes \"B.design\".\n"+strings.Repeat("A contains B.\n", 2000))
	b = parse("B.design", strings.Repeat("A contains B.\n", 2000))
	_, ds = Compile("A.design", []*parser.File{a, b}, false)
	if len(ds) != MaxDiagnostics {
		t.Fatal("unbounded/missing diagnostics", len(ds))
	}
	for _, d := range ds {
		if len(d.Related) > MaxRelated {
			t.Fatal("unbounded related evidence")
		}
	}
}
