package parser

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestSourceASTParity(t *testing.T) {
	files, err := filepath.Glob("../../examples/*.ast.json")
	if err != nil || len(files) == 0 {
		t.Fatal("No reference AST fixtures", err)
	}
	for _, file := range files {
		t.Run(filepath.Base(file), func(t *testing.T) {
			src := strings.TrimSuffix(file, ".ast.json") + ".sdui"
			source, e := os.ReadFile(src)
			if e != nil {
				t.Fatal(e)
			}
			doc, e := Parse(string(source))
			if e != nil {
				t.Fatal(e)
			}
			reference, e := os.ReadFile(file)
			if e != nil {
				t.Fatal(e)
			}
			var expected map[string]any
			if e = json.Unmarshal(reference, &expected); e != nil {
				t.Fatal(e)
			}
			encoded, _ := json.Marshal(Data(doc))
			var actual any
			json.Unmarshal(encoded, &actual)
			if !reflect.DeepEqual(expected["document"], actual) {
				t.Fatalf("AST differs from Python fixture %s", file)
			}
		})
	}
}
func TestPositionsAndBoundaries(t *testing.T) {
	src := "sdui 0.2;\r\nP=[\"\"\"# Æøå\ntext\"\"\"];"
	doc, e := Parse(src)
	if e != nil {
		t.Fatal(e)
	}
	l := doc.Definitions[0].Root.Rows[0].Items[0].Text
	if l.Span.Line != 2 || l.Span.Column != 4 || src[l.Span.Start:l.Span.End] != "\"\"\""+l.Value.(string)+"\"\"\"" {
		t.Fatal(l)
	}
	cases := map[string]string{"sdui 0.20; P=[];": "version", "sdui 2e-1; P=[];": "version", "sdui 0.2; P=[] {,};": "syntax", "sdui 0.2; P=[\"\\uD800\"];": "lexical", string([]byte{255}): "encoding", strings.Repeat(";", 50001): "token-limit", strings.Repeat("x", MaxBytes+1): "source-limit", "sdui 0.2; P=" + strings.Repeat("[", 65) + strings.Repeat("]", 65) + ";": "depth-limit"}
	for source, code := range cases {
		_, e := Parse(source)
		d, ok := e.(*Diagnostic)
		if !ok || d.Code != code {
			t.Errorf("Expected %s, got %v", code, e)
		}
	}
}
func TestTruncations(t *testing.T) {
	src := `sdui 0.2; Page = [<greeting="Hello">];`
	for i := 1; i < len(src); i++ {
		if _, e := Parse(src[:i]); e == nil {
			t.Fatalf("Accepted truncation %d", i)
		}
	}
}
func FuzzParser(f *testing.F) {
	for _, s := range []string{`sdui 0.2; P=[];`, `sdui 0.2; P=[<"Hi"; b=button("OK")>]*b {16:9,<->};`, "\xff", `sdui 0.2; P=["\uFFFF"];`} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) { _, _ = Parse(s) })
}
