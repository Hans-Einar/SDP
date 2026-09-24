package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

func decoded(t *testing.T, value any) any {
	t.Helper()
	b, e := json.Marshal(value)
	if e != nil {
		t.Fatal(e)
	}
	var out any
	if e = json.Unmarshal(b, &out); e != nil {
		t.Fatal(e)
	}
	return out
}
func TestFrozenPythonPortCases(t *testing.T) {
	raw, e := os.ReadFile("testdata/python-port-cases.json")
	if e != nil {
		t.Fatal(e)
	}
	var cases []struct {
		Source      string
		AST         any
		Diagnostics []struct {
			Code string
			Span any
		}
		Canonical *string
	}
	if e = json.Unmarshal(raw, &cases); e != nil {
		t.Fatal(e)
	}
	for i, c := range cases {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			m, diagnostics := Check(c.Source)
			if !reflect.DeepEqual(decoded(t, Data(m)), c.AST) {
				t.Fatalf("AST differs for %s", c.Source)
			}
			if len(diagnostics) != len(c.Diagnostics) {
				t.Fatalf("diagnostics %v want %v source %s", diagnostics, c.Diagnostics, c.Source)
			}
			for j, d := range diagnostics {
				if d.Code != c.Diagnostics[j].Code || !reflect.DeepEqual(decoded(t, Data(d.Span)), c.Diagnostics[j].Span) {
					t.Fatalf("diagnostic %d got %v want %v", j, d, c.Diagnostics[j])
				}
			}
			if c.Canonical != nil {
				s, d := Canonical(m)
				if len(d) > 0 || s != *c.Canonical {
					t.Fatal("canonical mismatch", d)
				}
			}
		})
	}
}
func TestArchitectureModel(t *testing.T) {
	source, e := os.ReadFile("../../../SDUI/design/architecture.design")
	if e != nil {
		t.Fatal(e)
	}
	m, d := Check(string(source))
	if len(d) > 0 {
		t.Fatal(d)
	}
	if len(m.Declarations) != 463 || len(m.Statements) != 1424 {
		t.Fatalf("unexpected model size %d/%d", len(m.Declarations), len(m.Statements))
	}
	raw, e := os.ReadFile("../../../SDUI/design/architecture.ast.json")
	if e != nil {
		t.Fatal(e)
	}
	var report map[string]any
	if e = json.Unmarshal(raw, &report); e != nil {
		t.Fatal(e)
	}
	expected := report["ast"]
	if expected == nil {
		expected = report
	}
	if !reflect.DeepEqual(decoded(t, Data(m)), expected) {
		t.Fatal("full architecture source AST differs")
	}
}
func TestLimitsAndTruncations(t *testing.T) {
	source := "language design-core version 0.5.\nunit Host.\nfunctionality Run.\nHost owns Run.\n"
	for i := range source {
		_, _ = Check(source[:i])
	}
	_, e := Parse(strings.Repeat("x", MaxBytes+1))
	if e == nil {
		t.Fatal("unbounded source")
	}
	_, e = Parse("\xff")
	if e == nil {
		t.Fatal("invalid utf8")
	}
}
func FuzzParser(f *testing.F) {
	for _, s := range []string{"", "language design-core version 0.5.\n", "language design-core version 0.5.\nunit A.\nA contains A.\n"} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		if len(s) > 20000 {
			t.Skip()
		}
		_, _ = Check(s)
	})
}
