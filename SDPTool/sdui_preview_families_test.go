package sdptool

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/presentation"
)

func TestSDUIAllFamilyDiscoveryAndSourceDelegation(t *testing.T) {
	source, err := os.ReadFile("../SDUI/go/parser/testdata/wci4-all-families.sdui")
	if err != nil {
		t.Fatal(err)
	}
	root, _ := projectFixture(t)
	path := filepath.Join(root, "SDP", "Families.sdui")
	if err := os.WriteFile(path, source, 0600); err != nil {
		t.Fatal(err)
	}
	p, err := Discover(root)
	if err != nil || len(p.Inventory.SDUI) != 1 || p.Inventory.SDUI[0].Profile != "sdui/0.3" || len(p.Sources) != 1 || p.Sources[0].State != "validated" {
		t.Fatal(p.Sources, err)
	}
	nodes, _, err := UINodes(p)
	if err != nil {
		t.Fatal(err)
	}
	var target *Target
	for _, n := range nodes {
		if n.Kind == "frame" && n.Label == "Main" {
			target = n.Target
		}
	}
	if target == nil || target.Operation != "sdui-preview" || target.Revision == "" {
		t.Fatal("source delegation absent", nodes)
	}
	out := filepath.Join(t.TempDir(), "preview")
	result, err := UIPreview(context.Background(), p, target.Model, target.Entry, out, target.Revision)
	if err != nil || result.Profile != "sdui/0.3" {
		t.Fatal(result, err)
	}
	actual, err := os.ReadFile(result.Entry)
	if err != nil {
		t.Fatal(err)
	}
	_, roots, err := parser.Compile(string(source))
	if err != nil {
		t.Fatal(err)
	}
	want, err := presentation.Markdown(roots["Main"], 160)
	if err != nil || string(actual) != want {
		t.Fatal("SDPTool diverged from SDUI-owned source preview", err)
	}
	for _, s := range []string{"resources not supplied; source intent only", "provider data not supplied", "Main/hiddenReading", "Main/amount", "Main/modal", "Main/palette", "Main/panes", "Main/multi", "Main/context"} {
		if !strings.Contains(want, s) {
			t.Fatal("missing family", s)
		}
	}
	// Actual content revision, not a changed inventory hint: stale source must
	// preserve the old artifact before any new provider/readiness claim.
	if err := os.WriteFile(path, append(source, []byte("\n\n")...), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := UIPreview(context.Background(), p, target.Model, target.Entry, out, target.Revision); err == nil || !strings.Contains(err.Error(), "stale") {
		t.Fatal("stale source accepted", err)
	}
	after, err := os.ReadFile(result.Entry)
	if err != nil || string(after) != string(actual) {
		t.Fatal("rejected request replaced artifact", err)
	}
	// A malformed opt-in is invalid discovery, not a source-rendered fallback.
	if err := os.WriteFile(path, []byte(`sdui 0.3; Main=[markdown("x",description="D")];`), 0600); err != nil {
		t.Fatal(err)
	}
	p, err = Discover(root)
	if err != nil || len(p.Sources) != 1 || p.Sources[0].State != "invalid" {
		t.Fatal(p.Sources, err)
	}
	nodes, _, err = UINodes(p)
	if err != nil {
		t.Fatal(err)
	}
	for _, n := range nodes {
		if n.Target != nil {
			t.Fatal("invalid schema has launch target", n)
		}
	}
}
