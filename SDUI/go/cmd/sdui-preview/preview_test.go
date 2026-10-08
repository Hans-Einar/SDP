package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/presentation"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAllFamilySourcePreviewDelegatesTruthfully(t *testing.T) {
	source, err := os.ReadFile("../../parser/testdata/wci4-all-families.sdui")
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "Families.sdui")
	if err := os.WriteFile(path, source, 0600); err != nil {
		t.Fatal(err)
	}
	out := t.TempDir()
	var stdout bytes.Buffer
	if err := execute([]string{"-source", path, "-entry", "Main", "-output", out}, &stdout); err != nil {
		t.Fatal(err)
	}
	var result map[string]string
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil || result["profile"] != "sdui/0.3" || result["schema"] != "sdptool/0.2" || result["revision"] != fmt.Sprintf("%x", sha256.Sum256(source)) {
		t.Fatal(result, err)
	}
	doc, roots, err := parser.Compile(string(source))
	if err != nil {
		t.Fatal(err)
	}
	want, err := presentation.Combined(roots["Main"], 160, doc)
	if err != nil {
		t.Fatal(err)
	}
	actual, err := os.ReadFile(filepath.Join(out, "entry.md"))
	if err != nil || string(actual) != want {
		t.Fatal("helper diverged", err)
	}
	for _, s := range []string{"resources not supplied; source intent only", "provider data not supplied", "Main/hiddenReading", "Main/amount", "Main/palette", "Main/modal", "Main/panes", "Main/pair", "Main/multi"} {
		if !strings.Contains(want, s) {
			t.Fatal("missing family", s)
		}
	}
	snapshot, err := os.ReadFile(filepath.Join(out, "source.sdui"))
	if err != nil || !bytes.Equal(snapshot, source) {
		t.Fatal("source snapshot changed", err)
	}
	manifest, err := os.ReadFile(filepath.Join(out, "sdui.json"))
	if err != nil {
		t.Fatal(err)
	}
	var metadata struct {
		Profile   string
		SpanUnits string
		Outputs   map[string]string
	}
	if err := json.Unmarshal(manifest, &metadata); err != nil || metadata.Profile != "sdui/0.3" || metadata.SpanUnits != "UTF-8 bytes; end exclusive" || metadata.Outputs["source.sdui"] != result["revision"] || metadata.Outputs["entry.md"] != fmt.Sprintf("%x", sha256.Sum256(actual)) {
		t.Fatal("artifact provenance", metadata, err)
	}
	if err := os.WriteFile(path, append(source, '\n'), 0600); err != nil {
		t.Fatal(err)
	}
	next := t.TempDir()
	stdout.Reset()
	err = execute([]string{"-source", path, "-entry", "Main", "-revision", result["revision"], "-output", next}, &stdout)
	if err == nil || !strings.Contains(err.Error(), "stale") || stdout.Len() != 0 {
		t.Fatal("stale source accepted", err)
	}
	files, err := os.ReadDir(next)
	if err != nil || len(files) != 0 {
		t.Fatal("stale request published", err)
	}
}
func TestPreviewCheckDoesNotClaimNativeAdapterReadiness(t *testing.T) {
	for _, body := range []string{`svg(art.X.@resource,description="Figure",fallback="label")`, `markdown("Text",description="Reading",fallback="label")`} {
		path := filepath.Join(t.TempDir(), "Preview.sdui")
		if err := os.WriteFile(path, []byte(`sdui 0.3; ref: art "unopened"; Main=[p=`+body+` {visible=false}];`), 0600); err != nil {
			t.Fatal(err)
		}
		out := t.TempDir()
		var stdout bytes.Buffer
		err := execute([]string{"-source", path, "-entry", "Main", "-check", "-output", out}, &stdout)
		if err == nil || stdout.Len() != 0 || !strings.Contains(err.Error(), "unsupported-preview") || !strings.Contains(err.Error(), "Main/p") {
			t.Fatal(err, stdout.String())
		}
		files, err := os.ReadDir(out)
		if err != nil || len(files) != 0 {
			t.Fatal("failed check published", files, err)
		}
	}
}
