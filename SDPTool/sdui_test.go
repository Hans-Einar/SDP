package sdptool

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSDUIProfileDiscoveryAndPreview(t *testing.T) {
	for _, version := range []string{"0.2", "0.3"} {
		t.Run(version, func(t *testing.T) {
			root, _ := projectFixture(t)
			source := `sdui ` + version + `; page=[button("OK")];`
			if version == "0.3" {
				source = `sdui 0.3; ref: app "absent.sdl"; page=[items=tree("Navigation",callback=app.Activate.@invoke)];`
			}
			path := filepath.Join(root, "SDP", "Page.sdui")
			if err := os.WriteFile(path, []byte(source), 0600); err != nil {
				t.Fatal(err)
			}
			p, err := Discover(root)
			if err != nil || len(p.Inventory.SDUI) != 1 || len(p.Sources) != 1 || p.Sources[0].State != "validated" {
				t.Fatal(p.Sources, err)
			}
			nodes, _, err := UINodes(p)
			if err != nil {
				t.Fatal(err)
			}
			var target *Target
			for _, n := range nodes {
				if n.Kind == "frame" {
					target = n.Target
				}
			}
			if target == nil || target.Operation != "sdui-preview" {
				t.Fatal("missing static delegation", nodes)
			}
			out := filepath.Join(t.TempDir(), "preview")
			result, err := UIPreview(context.Background(), p, target.Model, target.Entry, out, target.Revision)
			if err != nil || result.Profile != "sdui/"+version {
				t.Fatal(result, err)
			}
			metadata, err := os.ReadFile(filepath.Join(out, "sdptool.json"))
			if err != nil {
				t.Fatal(err)
			}
			var recorded Result
			if err := json.Unmarshal(metadata, &recorded); err != nil || recorded != result {
				t.Fatal("returned and published identities differ", err)
			}
			text, err := os.ReadFile(result.Entry)
			if err != nil || version == "0.3" && (!strings.Contains(string(text), "provider data not supplied") || !strings.Contains(string(text), "Navigation")) {
				t.Fatal(string(text), err)
			}
			// A stale inventory may not relabel new source or overwrite an artifact.
			p.Inventory.SDUI[0].Profile = "sdui/0.3"
			if version == "0.3" {
				p.Inventory.SDUI[0].Profile = "sdui/0.2"
			}
			if _, err := UIPreview(context.Background(), p, target.Model, target.Entry, out, ""); err == nil || !strings.Contains(err.Error(), "profile-mismatch") {
				t.Fatal("inventory/source mismatch accepted", err)
			}
			after, _ := os.ReadFile(result.Entry)
			if string(after) != string(text) {
				t.Fatal("failed preview changed existing artifact")
			}
			nodes, _, err = UINodes(p)
			if err != nil || len(nodes) != 2 || nodes[0].State != "invalid" || nodes[0].Target != nil || !strings.Contains(nodes[0].Diagnostic, "profile-mismatch") {
				t.Fatal(nodes, err)
			}
		})
	}
}

func TestSDUIInvalidVersusUnsupportedProfile(t *testing.T) {
	for _, tc := range []struct{ source, profile, state string }{
		{`sdui 0.3; page=[tree()];`, "sdui/0.3", "invalid"},
		{`sdui 0.2; page=[tree("Navigation")];`, "sdui/0.2", "invalid"},
		{`sdui 0.4; page=[];`, "sdui/0.4", "unsupported"},
	} {
		info := inspectSource("SDP/Page.sdui", "Page.sdui", []byte(tc.source), nil)
		if info.Profile != tc.profile || info.State != tc.state || info.Diagnostic == "" {
			t.Fatal(info)
		}
	}
}
