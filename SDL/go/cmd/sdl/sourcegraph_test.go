package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestGraphCLI(t *testing.T) {
	r := t.TempDir()
	os.Mkdir(filepath.Join(r, "Containers"), 0755)
	os.WriteFile(filepath.Join(r, "System.design"), []byte("language design-core version 0.6.\nsystem Demo.\nDemo contains Containers/Child.\n"), 0644)
	child := filepath.Join(r, "Containers/Child.design")
	os.WriteFile(child, []byte("language design-core version 0.6.\ncontainer Child.\n"), 0644)
	for _, tc := range []struct {
		args     []string
		code     int
		contains string
	}{
		{[]string{"check", filepath.Join(r, "System.design")}, 0, `"system": "Demo"`},
		{[]string{"ast", filepath.Join(r, "System.design")}, 0, `"source": "Containers/Child.design"`},
		{[]string{"format", filepath.Join(r, "System.design")}, 1, "--file-map"},
		{[]string{"format", filepath.Join(r, "System.design"), "--file-map"}, 0, "sdl-formatted-files/1"},
		{[]string{"fragment", child}, 0, "parsed-context-required"},
		{[]string{"check", child}, 1, "SYSTEM_CARDINALITY"},
	} {
		var out, err bytes.Buffer
		code := execute(tc.args, bytes.NewReader(nil), &out, &err)
		if code != tc.code || !bytes.Contains(out.Bytes(), []byte(tc.contains)) {
			t.Fatalf("%v: %d %s %s", tc.args, code, &out, &err)
		}
		var obj any
		if json.Unmarshal(out.Bytes(), &obj) != nil {
			t.Fatal("not JSON")
		}
	}
}

func TestGraphViewConsumers(t *testing.T) {
	root := t.TempDir()
	os.Mkdir(filepath.Join(root, "Children"), 0700)
	entry := filepath.Join(root, "System.design")
	os.WriteFile(entry, []byte("language design-core version 0.6.\nsystem Demo.\nDemo contains Children/Child.\n"), 0600)
	child := filepath.Join(root, "Children", "Child.design")
	os.WriteFile(child, []byte("language design-core version 0.6.\nunit Child.\n"), 0600)
	var revisions []string
	for _, cmd := range []string{"viewpoints", "view"} {
		output := filepath.Join(t.TempDir(), "out")
		args := []string{cmd, entry, "--output", output}
		if cmd == "view" {
			args = append(args, "--uri", "sdl-view://demo/VP02")
		}
		var out, errs bytes.Buffer
		if code := execute(args, bytes.NewReader(nil), &out, &errs); code != 0 {
			t.Fatal(code, out.String(), errs.String())
		}
		b, e := os.ReadFile(filepath.Join(output, "sources.json"))
		if e != nil {
			t.Fatal(e)
		}
		var p struct {
			Revision string `json:"revision"`
		}
		if e = json.Unmarshal(b, &p); e != nil {
			t.Fatal(e)
		}
		revisions = append(revisions, p.Revision)
		args[3] = filepath.Join(root, "Children")
		if code := execute(args, bytes.NewReader(nil), &out, &errs); code == 0 {
			t.Fatal("source directory overwritten")
		}
	}
	if revisions[0] == "" || revisions[0] != revisions[1] {
		t.Fatal("consumer revisions differ", revisions)
	}
	if _, e := os.Stat(child); e != nil {
		t.Fatal("dependency lost", e)
	}
}
