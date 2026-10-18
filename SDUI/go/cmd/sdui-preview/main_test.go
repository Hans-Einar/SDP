package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Hans-Einar/SDP/SDUI/go/parser"
	"github.com/Hans-Einar/SDP/SDUI/go/presentation"
)

func TestStaticPreviewProfileAndLegacyProtocol(t *testing.T) {
	for _, version := range []string{"0.2", "0.3"} {
		source := `sdui ` + version + `; page=[button("OK")];`
		if version == "0.3" {
			source = `sdui 0.3; ref: app "absent.sdl"; page=[items=tree("Navigation",callback=app.Activate.@invoke)];`
		}
		path := filepath.Join(t.TempDir(), "Page.sdui")
		if err := os.WriteFile(path, []byte(source), 0600); err != nil {
			t.Fatal(err)
		}
		out := t.TempDir()
		var stdout bytes.Buffer
		if err := execute([]string{"-source", path, "-entry", "page", "-output", out}, &stdout); err != nil {
			t.Fatal(err)
		}
		var result map[string]string
		if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
			t.Fatal(err)
		}
		manifest, err := os.ReadFile(filepath.Join(out, "sdui.json"))
		if err != nil {
			t.Fatal(err)
		}
		var metadata map[string]any
		if err := json.Unmarshal(manifest, &metadata); err != nil {
			t.Fatal(err)
		}
		if version == "0.3" {
			if result["profile"] != "sdui/0.3" || metadata["profile"] != "sdui/0.3" {
				t.Fatal(result, metadata)
			}
		} else {
			// The 0.2 wire objects stay exactly the existing shape, with no new keys.
			if len(result) != 5 || len(metadata) != 6 || result["profile"] != "" || metadata["profile"] != nil {
				t.Fatal("legacy protocol changed", result, metadata)
			}
		}
		doc, roots, err := parser.Compile(source)
		if err != nil {
			t.Fatal(err)
		}
		want, err := presentation.Combined(roots["page"], 160, doc)
		actual, readErr := os.ReadFile(filepath.Join(out, "entry.md"))
		if err != nil || readErr != nil || string(actual) != want {
			t.Fatal("helper diverged from SDUI presentation", err, readErr)
		}
		if version == "0.3" && !strings.Contains(string(actual), "provider data not supplied") {
			t.Fatal("collection preview implied live data", string(actual))
		}
	}
}

func TestCheckNeverClaimsCollectionProviderReadiness(t *testing.T) {
	for _, source := range []string{
		`sdui 0.3; page=[items=tree("Navigation")];`,
		`sdui 0.3; page=[items=list("Entries") {visible=false}];`,
	} {
		path := filepath.Join(t.TempDir(), "Page.sdui")
		if err := os.WriteFile(path, []byte(source), 0600); err != nil {
			t.Fatal(err)
		}
		out := t.TempDir()
		var stdout bytes.Buffer
		err := execute([]string{"-source", path, "-entry", "page", "-check", "-output", out}, &stdout)
		if err == nil || !strings.Contains(err.Error(), "unsupported-provider") || !strings.Contains(err.Error(), "page/items") || stdout.Len() != 0 {
			t.Fatal("collection readiness or unsourced rejection", err, stdout.String())
		}
		files, err := os.ReadDir(out)
		if err != nil || len(files) != 0 {
			t.Fatal("preflight wrote output", files, err)
		}
	}
}
