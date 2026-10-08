package prototype

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Hans-Einar/SDP/SDUI/go/preparation"
)

func TestPanePreflightCannotClaimStandaloneReadiness(t *testing.T) {
	for _, source := range []string{
		`sdui 0.3; Main=[t=tabs("T")[p=page("P")[]]];`,
		`sdui 0.3; Main=[s=split("vertical")[a=[];b=[]] {visible=false}];`,
	} {
		path := filepath.Join(t.TempDir(), "Page.sdui")
		if err := os.WriteFile(path, []byte(source), 0600); err != nil {
			t.Fatal(err)
		}
		_, report, err := Check(path, "Main", "")
		var d *preparation.Diagnostic
		if err == nil || !strings.Contains(err.Error(), "unsupported-pane") || !errors.As(err, &d) || d.Capability.Dimension != preparation.Host || d.Span.Line == 0 || !strings.HasPrefix(d.Path, "Main/") || report.Status != "" || report.Profile != "sdui/0.3" {
			t.Fatal(report, err)
		}
	}
}
