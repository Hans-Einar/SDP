package prototype

import (
	"errors"
	"github.com/Hans-Einar/SDP/SDUI/go/preparation"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStandaloneExtendedTextExplicitlyUnsupported(t *testing.T) {
	for _, arg := range []string{`multiline=false`, `readOnly=false`, `required=false`, `placeholder=""`} {
		path := filepath.Join(t.TempDir(), "text.sdui")
		if err := os.WriteFile(path, []byte(`sdui 0.3; Field=input("Label",`+arg+`); Main=[note=Field {visible=false}];`), 0600); err != nil {
			t.Fatal(err)
		}
		_, report, err := Check(path, "Main", "")
		var d *preparation.Diagnostic
		if !errors.As(err, &d) || !strings.Contains(err.Error(), "unsupported-text") || report.Status != "" || report.Profile != "sdui/0.3" || d.Path != "Main/note" || d.Span.Line == 0 || len(d.Uses) != 1 || d.Capability != (preparation.Capability{Dimension: preparation.Host, ID: "input", Major: 1}) {
			t.Fatal(report, err)
		}
	}
}
