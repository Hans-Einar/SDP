package prototype

import (
	"errors"
	"github.com/Hans-Einar/SDP/SDUI/go/preparation"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStandaloneScalarCapabilityGapIsExplicit(t *testing.T) {
	for _, body := range []string{`checkbox("C")`, `slider("S",min=0,max=10,step=1,value=1)`, `number("N",min=0,max=1,step=0.1,value=0.3)`, `select("Choice",value="unverified")`} {
		file := filepath.Join(t.TempDir(), "values.sdui")
		if err := os.WriteFile(file, []byte(`sdui 0.3; Main=[field=`+body+` {visible=false}];`), 0600); err != nil {
			t.Fatal(err)
		}
		_, report, err := Check(file, "Main", "")
		var diagnostic *preparation.Diagnostic
		if err == nil || !errors.As(err, &diagnostic) || report.Status != "" || report.Profile != "sdui/0.3" || diagnostic.Path != "Main/field" || diagnostic.Span.Line == 0 {
			t.Fatal(report, err)
		}
		if strings.HasPrefix(body, "select") {
			if diagnostic.Capability.Dimension != preparation.Provider || diagnostic.Capability.ID != "choice-options" {
				t.Fatal(diagnostic)
			}
		} else if diagnostic.Capability.Dimension != preparation.Host || !strings.Contains(err.Error(), "unsupported-value") {
			t.Fatal(err)
		}
	}
}
