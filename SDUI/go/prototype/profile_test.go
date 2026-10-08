package prototype

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Hans-Einar/SDP/SDUI/go/preparation"
)

func TestCollectionPreflightRequiresActualProvider(t *testing.T) {
	for _, body := range []string{`tree("Navigation")`, `list("Entries") {visible=false}`} {
		path := filepath.Join(t.TempDir(), "Page.sdui")
		if err := os.WriteFile(path, []byte(`sdui 0.3; collection=`+body+`; page=[items=collection];`), 0600); err != nil {
			t.Fatal(err)
		}
		candidate, report, err := Check(path, "page", "")
		var diagnostic *preparation.Diagnostic
		if err == nil || !strings.Contains(err.Error(), "unsupported-provider") || !errors.As(err, &diagnostic) {
			t.Fatal("missing typed provider diagnostic", err)
		}
		if diagnostic.Path != "page/items" || diagnostic.Span.Line == 0 || len(diagnostic.Uses) == 0 || diagnostic.Capability.Dimension != preparation.Provider {
			t.Fatal("lost source/reuse provenance", diagnostic)
		}
		if report.Profile != "sdui/0.3" || report.Status != "" || candidate.Document.Profile != report.Profile {
			t.Fatal("incorrect profile/readiness", report)
		}
	}
}

func TestLegacyCheckWireShapeAndBoundaries(t *testing.T) {
	path := filepath.Join(t.TempDir(), "Page.sdui")
	if err := os.WriteFile(path, []byte(`sdui 0.2; ref: app "absent.sdl"; page=[save=button("Save",callback=app.Save.@invoke)];`), 0600); err != nil {
		t.Fatal(err)
	}
	_, report, err := Check(path, "page", "")
	if err != nil || report.Status != "prototype-unbound" || report.Profile != "" {
		t.Fatal(report, err)
	}
	encoded, err := json.Marshal(report)
	if err != nil || strings.Contains(string(encoded), `"profile"`) || strings.Contains(string(encoded), "connected") {
		t.Fatal(string(encoded), err)
	}
	// An application-capable document host does not license standalone preflight.
	if err := os.WriteFile(path, []byte(`sdui 0.3; page=[button("OK")];`), 0600); err != nil {
		t.Fatal(err)
	}
	_, report, err = Check(path, "page", "")
	if err == nil || report.Profile != "sdui/0.3" || report.Status != "" {
		t.Fatal("claimed native readiness beyond standalone admission", report, err)
	}
}
