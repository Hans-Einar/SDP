package documents

import (
	"context"
	"fmt"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/viewpoint"
	"strings"
	"testing"
)

type unusedTraceabilityRenderer struct{}

func (unusedTraceabilityRenderer) Identity() string { return "test-unused" }
func (unusedTraceabilityRenderer) Render(context.Context, string) ([]byte, error) {
	return nil, fmt.Errorf("VP11 must not invoke a diagram renderer")
}
func TestTraceabilityOverviewWithoutDiagramClaims(t *testing.T) {
	v, e := viewpoint.New("language design-core version 0.5.\nunit Solo.\n")
	if e != nil {
		t.Fatal(e)
	}
	b, e := Selected(context.Background(), v, viewpoint.Query{Viewpoint: "VP11", Projection: "overview"}, unusedTraceabilityRenderer{})
	if e != nil {
		t.Fatal(e)
	}
	text := string(b.Files["entry.md"])
	for _, unwanted := range []string{"```mermaid", "No diagrams match", "Symbol profile:"} {
		if strings.Contains(text, unwanted) {
			t.Fatal(unwanted, text)
		}
	}
	if !strings.Contains(text, "No facts are declared") || !strings.Contains(text, "Solo (unit)") || len(b.Manifest.Gaps) != len(v.Gaps) {
		t.Fatal(text)
	}
	for name := range b.Files {
		if strings.HasPrefix(name, "diagrams/") {
			t.Fatal("unexpected diagram asset", name)
		}
	}
}
