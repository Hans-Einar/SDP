package prototype

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStandaloneInteractionPreflightIsNotReadiness(t *testing.T) {
	for _, body := range []string{`c=command("C",toggle=true) {visible=false}`, `b=button("B",toggle=false)`, `m=menu("M")[] {visible=false}`, `d=dialog("D")[]`} {
		file := filepath.Join(t.TempDir(), "source.sdui")
		if err := os.WriteFile(file, []byte(`sdui 0.3; Main=[`+body+`];`), 0600); err != nil {
			t.Fatal(err)
		}
		_, report, err := Check(file, "Main", "")
		if err == nil || !strings.Contains(err.Error(), "unsupported-interaction") || report.Status != "" || report.Profile != "sdui/0.3" {
			t.Fatal(report, err)
		}
	}
	file := filepath.Join(t.TempDir(), "source.sdui")
	if err := os.WriteFile(file, []byte(`sdui 0.3; Main=[button(command="missing")];`), 0600); err != nil {
		t.Fatal(err)
	}
	if _, _, err := Check(file, "Main", ""); err == nil || !strings.Contains(err.Error(), "interaction-reference") {
		t.Fatal("unresolved entry not checked", err)
	}
}
