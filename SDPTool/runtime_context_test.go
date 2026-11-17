package sdptool

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestDocumentContextCheckCLIAndAction(t *testing.T) {
	dir := t.TempDir()
	source, fixtures := filepath.Join(dir, "model.sdui"), filepath.Join(dir, "context.json")
	model := `sdui 0.3; page=[settings=tabs("Settings")[one=page("One")[choice=select("Choice",value="a")]]];`
	data := `{"schema":"sdui-prototype-context/1","choices":{"page/settings/one/choice":[{"id":"a","label":"A","enabled":true}]}}`
	if err := os.WriteFile(source, []byte(model), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(fixtures, []byte(data), 0600); err != nil {
		t.Fatal(err)
	}
	params := map[string]any{"source": source, "entry": "page", "context": fixtures}
	invoke := func() InvokeResponse {
		return Invoke(context.Background(), bytes.NewReader(requestBytes(t, "", "sdui.source.check", params)), ActionCaller{})
	}
	r := invoke()
	if r.ExitCode != 0 {
		t.Fatalf("%+v", r)
	}
	wire, _ := json.Marshal(r.Result)
	var report map[string]any
	if err := json.Unmarshal(wire, &report); err != nil {
		t.Fatal(err)
	}
	if report["host"] != "document" || report["nativeLayout"] != "pending" || report["status"] != "preflight" {
		t.Fatal(report)
	}
	digest, ok := report["contextRevision"].(string)
	if !ok || len(digest) != 64 {
		t.Fatal(report)
	}
	var out, errs bytes.Buffer
	args := []string{"sdui", "check", "--source", source, "--entry", "page", "--context", fixtures, "--context-revision", digest, "--json"}
	if code := Run(context.Background(), args, &out, &errs); code != 0 {
		t.Fatal(errs.String())
	}
	var cli map[string]any
	json.Unmarshal(out.Bytes(), &cli)
	if !reflect.DeepEqual(cli, report) {
		t.Fatalf("CLI/action differ: %s / %s", out.String(), wire)
	}
	params["contextRevision"] = digest
	if err := os.WriteFile(fixtures, []byte(strings.Replace(data, `"A"`, `"Changed"`, 1)), 0600); err != nil {
		t.Fatal(err)
	}
	if r = invoke(); r.ExitCode == 0 || !strings.Contains(r.Error.Message, "stale") {
		t.Fatalf("stale context accepted: %+v", r)
	}
	delete(params, "context")
	delete(params, "contextRevision")
	if r = invoke(); r.ExitCode == 0 {
		t.Fatal("missing choice provider accepted")
	}
	out.Reset()
	errs.Reset()
	if code := Run(context.Background(), []string{"sdui", "preview", "--source", source, "--entry", "page", "--context", fixtures, "--output", t.TempDir()}, &out, &errs); code == 0 {
		t.Fatal("static preview accepted ignored context")
	}
}

func TestNativeLaunchForwardsContextAsData(t *testing.T) {
	old := bundledNativeHost
	t.Cleanup(func() { bundledNativeHost = old })
	var compressed bytes.Buffer
	z := gzip.NewWriter(&compressed)
	z.Write([]byte("#!/bin/sh\nprintf '%s\\000' \"$@\" > \"$SDP_NATIVE_ARGUMENTS\"\n"))
	z.Close()
	bundledNativeHost = compressed.Bytes()
	t.Setenv("XDG_CACHE_HOME", t.TempDir())
	record := filepath.Join(t.TempDir(), "args")
	t.Setenv("SDP_NATIVE_ARGUMENTS", record)
	watch := false
	p := sduiRunParameters{Source: "/model with spaces.sdui", Entry: "page", Revision: "sourcehash", Context: "/fixtures with spaces.json", ContextRevision: "contexthash", Watch: &watch}
	if _, _, err := performSduiRun(context.Background(), actionContext{}, p); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(record)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"-entry", "page", "-revision", "sourcehash", "-context", p.Context, "-context-revision", "contexthash", "-watch=false", p.Source, ""}
	if !reflect.DeepEqual(strings.Split(string(b), "\x00"), want) {
		t.Fatalf("%q", b)
	}
	for _, flag := range []string{"--context", "--context-revision"} {
		args, mode, err := outputArguments([]string{"sdui", "check", flag, "--json"})
		if err != nil || mode || args[3] != "--json" {
			t.Fatalf("flag value consumed as output mode: %v %v %v", args, mode, err)
		}
	}
}
