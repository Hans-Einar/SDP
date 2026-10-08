package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLegacyPreviewFrozenOutputBytes(t *testing.T) {
	// Captured from a separately built git archive of delivered phase HEAD 90a94b5.
	for _, profile := range []string{"0.2", "0.3"} {
		want := map[string]string{
			"ast":      "98bbabd8fa4160cb8c7719e648ab94040a1c8423430a9bad927430e8b29abc76",
			"dump":     "544c362bf2ac1beb53c3a8af7a147248c262f542c0eae7a6310943e48588ddd9",
			"markdown": "d63e3f78beac2377aa69782af64b4ba8f064dbe6030234e430f2f846eb9e13a5",
			"svg":      "4f3be0dc2d90c5d9cb1847a14c4477d547a19a23fda7d7590adba95bd3140f43",
		}
		if profile == "0.3" {
			want["ast"] = "2f2694c46090820aa45cfb2bb79bee2681fb169994ded1f10f0a999b6ef7f0ef"
		}
		for format, hash := range want {
			var out, errs bytes.Buffer
			source := `sdui ` + profile + `; ref: art "unopened"; Main=[picture=svg(art.Image.@draw,label="Legacy");text="# Old prose"];`
			if code := execute([]string{"-", "--format", format, "--entry", "Main"}, strings.NewReader(source), &out, &errs); code != 0 {
				t.Fatal(code, errs.String())
			}
			if actual := fmt.Sprintf("%x", sha256.Sum256(out.Bytes())); actual != hash {
				t.Fatalf("%s %s changed: %s", profile, format, actual)
			}
		}
	}
}

func TestPreviewCLIAndAtomicRejection(t *testing.T) {
	for _, body := range []string{`svg(art.X.@resource,description="Figure",fallback="label")`, `markdown("# Text",description="Reading",fallback="reject")`} {
		source := `sdui 0.3; ref: art "unopened"; Main=[p=` + body + ` {visible=false}];`
		for _, format := range []string{"ast", "dump", "markdown"} {
			var out, errs bytes.Buffer
			if code := execute([]string{"-", "--format", format, "--entry", "Main"}, strings.NewReader(source), &out, &errs); code != 0 {
				t.Fatal(code, errs.String())
			}
			if format == "ast" {
				if !strings.Contains(out.String(), `"sdui-ast/0.3"`) {
					t.Fatal(out.String())
				}
			} else if !strings.Contains(out.String(), "resources not supplied; source intent only") || !strings.Contains(out.String(), "Main/p") {
				t.Fatal(out.String())
			}
		}
		path := filepath.Join(t.TempDir(), "prior.svg")
		if err := os.WriteFile(path, []byte("prior artifact"), 0600); err != nil {
			t.Fatal(err)
		}
		var out, errs bytes.Buffer
		if code := execute([]string{"-", "--format", "svg", "--entry", "Main", "-o", path}, strings.NewReader(source), &out, &errs); code == 0 || out.Len() != 0 || !strings.Contains(errs.String(), "unsupported-resource-export") || !strings.Contains(errs.String(), "Main/p") {
			t.Fatal(code, out.String(), errs.String())
		}
		after, err := os.ReadFile(path)
		if err != nil || string(after) != "prior artifact" {
			t.Fatal("failed export replaced artifact", err)
		}
	}
}
