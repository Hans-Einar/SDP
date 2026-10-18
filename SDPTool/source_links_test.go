package sdptool

import (
	"context"
	"fmt"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/documents"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func detailFixture(t *testing.T) (string, string) {
	t.Helper()
	root := t.TempDir()
	sources := map[string]string{
		"System.design": "language design-core version 0.6.\nincludes \"Åctors.design\".\nsystem Demo.\nDemo contains Caller.\n",
		"Åctors.design": `language design-core version 0.6.
actor Reader.
usecase Browse.
usecase Isolated.
feature Reading.
functionality Read.
unit Caller.
unit Handler.
mode Normal.
channel Bus.
contract Events.
contract Body.
message Notice.
scenario Happy.
scenario Unrelated.
field Value.
Reader pursues Browse.
Reader pursues Isolated.
Reading supports Browse.
Read contributes-to Reading.
Caller owns Read.
Body has completeness = closed.
Body has-field Value.
Value has presence = required.
Value has value-type = text.
Events has completeness = closed.
Happy has completeness = closed.
Unrelated has completeness = closed.
Bus upholds Events.
Events permits Notice.
Notice upholds Body.
Notice has message-kind = event.
Caller uses Bus as sender of Notice in mode Normal.
Handler uses Bus as receiver of Notice in mode Normal.
Happy exercises Browse.
Happy runs-in Normal.
Happy step 1 sends Notice from Caller to Handler via Bus.
Unrelated runs-in Normal.
Unrelated step 1 sends Notice from Caller to Handler via Bus.
`,
	}
	for name, s := range sources {
		m, e := parser.Parse(s)
		if e != nil {
			t.Fatal(e)
		}
		if e = os.WriteFile(filepath.Join(root, name), []byte(parser.FormatFile(m)), 0600); e != nil {
			t.Fatal(e)
		}
	}
	return root, filepath.Join(root, "System.design")
}
func TestFocusedUseCaseAndExactSourceLinks(t *testing.T) {
	root, source := detailFixture(t)
	for _, query := range []string{"diagram=VP01-Browse", "focus=Browse"} {
		result, e := Preview(context.Background(), PreviewOptions{Source: source, URI: "sdl-view://trial/VP01?" + query, Output: filepath.Join(t.TempDir(), "view")})
		if e != nil {
			t.Fatal(e)
		}
		b, e := os.ReadFile(result.Entry)
		if e != nil {
			t.Fatal(e)
		}
		text := string(b)
		for _, want := range []string{"usecase-beta", "sequenceDiagram", "Scenario: Happy", "Preconditions: not specified", "Postconditions: not specified", "Read | contributes-to | Reading", "Reading | supports | Browse"} {
			if !strings.Contains(text, want) {
				t.Fatal(want, text)
			}
		}
		if strings.Contains(text, "Scenario: Unrelated") || strings.Contains(text, "Use case: Isolated") || strings.Contains(text, "(Isolated; mode=") {
			t.Fatal("inferred scenario link")
		}
		if len(result.SourceLinks) != 2 {
			t.Fatal(result.SourceLinks)
		}
		links := map[string]SourceLink{}
		for _, l := range result.SourceLinks {
			links[l.ID] = l
			saved, e := os.ReadFile(filepath.Join(result.Directory, l.Snapshot))
			if e != nil {
				t.Fatal(e)
			}
			original, e := os.ReadFile(l.Path)
			if e != nil || string(saved) != string(original) || documents.Hash(saved) != l.Revision {
				t.Fatal(l, e)
			}
		}
		found := false
		for _, m := range regexp.MustCompile(`\[([^]]+)\]\(sdl-source://(s[0-9]+)/([0-9]+)/([0-9]+)\)`).FindAllStringSubmatch(text, -1) {
			l, ok := links[m[2]]
			if !ok {
				t.Fatal(m)
			}
			b, _ := os.ReadFile(l.Path)
			start, _ := strconv.Atoi(m[3])
			end, _ := strconv.Atoi(m[4])
			if start < 0 || end > len(b) || end <= start {
				t.Fatal(m)
			}
			statement := string(b[start:end])
			if statement == "Reader pursues Browse." {
				found = true
				if filepath.Base(l.Path) != "Åctors.design" {
					t.Fatal(l)
				}
			}
			if !strings.HasSuffix(statement, ".") {
				t.Fatal("not an exact declaration/statement", statement)
			}
		}
		if !found {
			t.Fatal("missing clickable pursuit", text)
		}
	}
	v, _, e := loadModel(source)
	if e != nil {
		t.Fatal(e)
	}
	child := filepath.Join(root, "Åctors.design")
	f, _ := os.OpenFile(child, os.O_APPEND|os.O_WRONLY, 0600)
	f.WriteString("\n")
	f.Close()
	if _, _, e = bindSourceLinks(v, source); e == nil || !strings.Contains(e.Error(), "stale") {
		t.Fatal(e)
	}
}
func TestSingleFileSourceLinksAndIsolatedUseCase(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "Model.design")
	os.WriteFile(source, []byte("language design-core version 0.5.\nusecase Alone.\n"), 0600)
	result, e := Preview(context.Background(), PreviewOptions{Source: source, URI: "sdl-view://trial/VP01?focus=Alone", Output: filepath.Join(t.TempDir(), "view")})
	if e != nil {
		t.Fatal(e)
	}
	text, _ := os.ReadFile(result.Entry)
	if len(result.SourceLinks) != 1 || !strings.Contains(string(text), "usecase-beta") || !strings.Contains(string(text), "[Alone](sdl-source://s0/") || !strings.Contains(string(text), "No explicitly linked scenario") {
		t.Fatal(string(text), result.SourceLinks)
	}
}

// Owner #025 removes count-only selection limits; real bundle byte/file and
// individual diagram geometry bounds remain explicit.
func TestSelectionHasNoArbitraryDiagramCountCap(t *testing.T) {
	source := "language design-core version 0.5.\n"
	for i := 0; i < 50; i++ {
		source += fmt.Sprintf("usecase Goal%d.\n", i)
	}
	m, e := parser.Parse(source)
	if e != nil {
		t.Fatal(e)
	}
	file := filepath.Join(t.TempDir(), "Many.design")
	os.WriteFile(file, []byte(parser.FormatFile(m)), 0600)
	result, e := Preview(context.Background(), PreviewOptions{Source: file, Viewpoint: "VP01", Output: filepath.Join(t.TempDir(), "view")})
	if e != nil {
		t.Fatal(e)
	}
	text, _ := os.ReadFile(result.Entry)
	if strings.Count(string(text), "```mermaid") != 50 {
		t.Fatal("truncated selection")
	}
}
