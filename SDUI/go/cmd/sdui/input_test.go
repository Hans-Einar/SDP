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

func TestBasicInputFrozenOutputBytes(t *testing.T) {
	// Captured by running cmd/sdui from the delivered M1 commit ea49991f,
	// independently built from git archive; these are not new-output snapshots.
	for _, profile := range []string{"0.2", "0.3"} {
		want := map[string]string{
			"ast":      "5e19c07f93beeebb4f226eab98d3054306b3306a8a5ae484393ed26652ead26a",
			"dump":     "65648af7d49fcb263bb9bd02188e1dec5e11383c88c85a563a824b005fd96c22",
			"markdown": "33c177735e8a5f1ff12bf167183c5d06ea1a6139344353e7541f4e48aeaa5395",
			"svg":      "f87b0b02c283d410a83b6ec9cf1a8fce08f7f2a94f00f77efd6acfd57482a9d8",
		}
		if profile == "0.3" {
			want["ast"] = "e443199ff530bc7b630b6af88c4efce744e094361d0abb9bf7d1e32dcd53e5f1"
		}
		for format, hash := range want {
			var out, errs bytes.Buffer
			source := `sdui ` + profile + `; Main=[field=input("Label",value="å🙂")];`
			if code := execute([]string{"-", "--format", format, "--entry", "Main"}, strings.NewReader(source), &out, &errs); code != 0 {
				t.Fatal(code, errs.String())
			}
			if actual := fmt.Sprintf("%x", sha256.Sum256(out.Bytes())); actual != hash {
				t.Fatalf("%s %s changed: %s", profile, format, actual)
			}
		}
	}
}

func TestExtendedInputCLIAndAtomicExportRejection(t *testing.T) {
	const source = `sdui 0.3; Main=[field=input("Text",multiline=false,placeholder="") {visible=false}];`
	var out, errs bytes.Buffer
	if code := execute([]string{"-"}, strings.NewReader(source), &out, &errs); code != 0 || !strings.Contains(out.String(), `"sdui-ast/0.3"`) || !strings.Contains(out.String(), `"value": false`) {
		t.Fatal(code, out.String(), errs.String())
	}
	for _, format := range []string{"dump", "markdown"} {
		out.Reset()
		errs.Reset()
		if code := execute([]string{"-", "--format", format, "--entry", "Main"}, strings.NewReader(source), &out, &errs); code != 0 || !strings.Contains(out.String(), "Static input") || !strings.Contains(out.String(), `placeholder=""`) {
			t.Fatal(code, out.String(), errs.String())
		}
	}
	path := filepath.Join(t.TempDir(), "prior.svg")
	if err := os.WriteFile(path, []byte("prior artifact"), 0600); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	errs.Reset()
	if code := execute([]string{"-", "--format", "svg", "--entry", "Main", "-o", path}, strings.NewReader(source), &out, &errs); code == 0 || out.Len() != 0 || !strings.Contains(errs.String(), "unsupported-text-export") || !strings.Contains(errs.String(), "Main/field") {
		t.Fatal(code, out.String(), errs.String())
	}
	data, err := os.ReadFile(path)
	if err != nil || string(data) != "prior artifact" {
		t.Fatal("output overwritten", err)
	}
}
