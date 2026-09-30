package sdptool

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestOutputArguments(t *testing.T) {
	for _, tc := range []struct {
		args, want []string
		json       bool
	}{
		{[]string{"tree", "--json"}, []string{"tree"}, true},
		{[]string{"--json", "repo", "tree"}, []string{"repo", "tree"}, true},
		{[]string{"--json", "--version"}, []string{"--version"}, true},
		{[]string{"select", "--uri", "--json"}, []string{"select", "--uri", "--json"}, false},
		{[]string{"preview", "--json", "--output", "x"}, []string{"preview", "--json", "--output", "x"}, false},
		{[]string{"release-log", "--version", "--json"}, []string{"release-log", "--version", "--json"}, false},
		{[]string{"tree", "--json=false"}, []string{"tree"}, false},
		{[]string{"tree", "--", "--json"}, []string{"tree", "--", "--json"}, false},
	} {
		args, mode, err := outputArguments(tc.args)
		if err != nil || mode != tc.json || !reflect.DeepEqual(args, tc.want) {
			t.Fatalf("%v -> %v %v %v", tc.args, args, mode, err)
		}
	}
}
func TestCLIHumanAndJSON(t *testing.T) {
	root, _ := projectFixture(t)
	for _, op := range []string{"tree", "discover", "--version"} {
		for _, machine := range []bool{false, true} {
			args := []string{root, op}
			if op == "--version" {
				args = []string{op}
			}
			if machine {
				args = append(args, "--json")
			}
			var out, errs bytes.Buffer
			if code := Run(context.Background(), args, &out, &errs); code != 0 {
				t.Fatalf("%v: %d %s", args, code, &errs)
			}
			if json.Valid(out.Bytes()) != machine {
				t.Fatalf("%v: %s", args, &out)
			}
			if op == "tree" && !machine && !strings.Contains(out.String(), "── SDL") {
				t.Fatal(out.String())
			}
		}
	}
	for _, machine := range []bool{false, true} {
		args := []string{"bogus"}
		if machine {
			args = append(args, "--json")
		}
		var out, errs bytes.Buffer
		if code := Run(context.Background(), args, &out, &errs); code == 0 || out.Len() != 0 || json.Valid(errs.Bytes()) != machine {
			t.Fatalf("%v: %d %s %s", args, code, &out, &errs)
		}
	}
}

func TestReleaseLogAndInstallOutputModes(t *testing.T) {
	root := t.TempDir()
	notes := "# Notes\n\n## [1.0.0] - 2026-09-28\n\n### Added\n\n- Example.\n"
	if err := os.WriteFile(filepath.Join(root, "RELEASE-NOTES.md"), []byte(notes), 0600); err != nil {
		t.Fatal(err)
	}
	for _, machine := range []bool{false, true} {
		var out, errs bytes.Buffer
		args := []string{root, "release-log", "--version", "1.0.0"}
		if machine {
			args = append(args, "--json")
		}
		if code := Run(context.Background(), args, &out, &errs); code != 0 || json.Valid(out.Bytes()) != machine {
			t.Fatalf("%v: %d %s %s", args, code, &out, &errs)
		}
		if !strings.Contains(out.String(), "Example.") {
			t.Fatal(out.String())
		}
		out.Reset()
		errs.Reset()
		args = []string{root, "install", "--artifact", "missing"}
		if machine {
			args = append(args, "--json")
		}
		if code := Run(context.Background(), args, &out, &errs); code == 0 {
			t.Fatal("missing artifact accepted")
		}
		if machine {
			if !json.Valid(out.Bytes()) || errs.Len() != 0 {
				t.Fatalf("machine install: %s %s", &out, &errs)
			}
		} else if out.Len() != 0 || errs.Len() == 0 {
			t.Fatalf("human install: %s %s", &out, &errs)
		}
	}
}
