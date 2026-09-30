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
