package prototype

import (
	"errors"
	"github.com/Hans-Einar/SDP/SDUI/go/preparation"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStandaloneExplicitPreviewUnsupportedAdapter(t *testing.T) {
	for _, tc := range []struct{ body, id string }{{`svg(art.X.@resource,description="Figure",fallback="label")`, "svg-resource"}, {`markdown("Text",description="Reading",fallback="reject")`, "markdown"}} {
		path := filepath.Join(t.TempDir(), "preview.sdui")
		if err := os.WriteFile(path, []byte(`sdui 0.3; ref: art "unopened"; P=`+tc.body+`; Main=[p=P {visible=false}];`), 0600); err != nil {
			t.Fatal(err)
		}
		_, report, err := Check(path, "Main", "")
		var d *preparation.Diagnostic
		if !errors.As(err, &d) || !strings.Contains(err.Error(), "unsupported-preview") || report.Status != "" || report.Profile != "sdui/0.3" || d.Path != "Main/p" || d.Span.Line == 0 || len(d.Uses) != 1 || d.Capability != (preparation.Capability{Dimension: preparation.Host, ID: tc.id, Major: 1}) {
			t.Fatal(report, err)
		}
	}
}
