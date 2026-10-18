package sdptool

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func requestBytes(t *testing.T, root, action string, params any) []byte {
	t.Helper()
	b, e := json.Marshal(map[string]any{"schema": InvokeSchema, "requestId": "trial-1", "action": action, "context": map[string]string{"project": root}, "parameters": params})
	if e != nil {
		t.Fatal(e)
	}
	return b
}
func TestInvokeDiscoveryAndCatalogue(t *testing.T) {
	root, _ := projectFixture(t)
	ctx := context.Background()
	before, e := os.ReadFile(filepath.Join(root, "SDP/SDP-project.manifest.yaml"))
	if e != nil {
		t.Fatal(e)
	}
	a, e := Actions(ctx, "")
	if e != nil || len(a.Actions) == 0 {
		t.Fatalf("%+v %v", a, e)
	}
	b, e := Actions(ctx, root)
	if e != nil || a.Revision != b.Revision {
		t.Fatalf("context changed schema revision: %v", e)
	}
	bad, e := Actions(ctx, t.TempDir())
	if e != nil || len(bad.Diagnostics) == 0 || len(bad.Actions) != len(a.Actions) {
		t.Fatalf("invalid context hidden: %+v %v", bad, e)
	}
	seen := map[string]bool{}
	for _, d := range b.Actions {
		if seen[d.ID] {
			t.Fatal("duplicate action")
		}
		seen[d.ID] = true
		if _, e := findAction(d.ID); e != nil {
			t.Fatal(e)
		}
	}
	raw := requestBytes(t, root, "project.discover", map[string]any{})
	r := Invoke(ctx, bytes.NewReader(raw), ActionCaller{})
	if r.ExitCode != 0 {
		t.Fatalf("%+v", r)
	}
	var legacy, errs bytes.Buffer
	if Run(ctx, []string{root, "discover", "--json"}, &legacy, &errs) != 0 {
		t.Fatal(errs.String())
	}
	expected, _ := json.Marshal(r.Result)
	var x, y any
	json.Unmarshal(expected, &x)
	json.Unmarshal(legacy.Bytes(), &y)
	if !reflect.DeepEqual(x, y) {
		t.Fatal("CLI and invoke differ")
	}
	file := filepath.Join(t.TempDir(), "request.json")
	if e = os.WriteFile(file, raw, 0600); e != nil {
		t.Fatal(e)
	}
	for _, src := range []string{"-", file} {
		var out, errout bytes.Buffer
		code := RunWithInput(ctx, []string{"invoke", "--request", src}, bytes.NewReader(raw), &out, &errout)
		var got InvokeResponse
		if code != 0 || json.Unmarshal(out.Bytes(), &got) != nil || got.RequestID == nil || *got.RequestID != "trial-1" || errout.Len() != 0 {
			t.Fatalf("%s: %d %s %s", src, code, &out, &errout)
		}
	}
	after, _ := os.ReadFile(filepath.Join(root, "SDP/SDP-project.manifest.yaml"))
	if !bytes.Equal(before, after) {
		t.Fatal("read changed project")
	}
}
func TestInvokeRejectsMalformedBeforeDispatch(t *testing.T) {
	root, _ := projectFixture(t)
	valid := string(requestBytes(t, root, "project.discover", map[string]any{}))
	tests := []string{"", `null`, `[]`, valid + `{}`, strings.Replace(valid, `"schema":`, `"Schema":`, 1), strings.Replace(valid, `"parameters":{}`, `"parameters":null`, 1), strings.Replace(valid, `"parameters":{}`, `"parameters":{"unexpected":true}`, 1), strings.Replace(valid, `"context":`, `"caller":{},"context":`, 1), strings.Replace(valid, `"requestId":"trial-1"`, `"requestId":""`, 1), strings.Replace(valid, `"requestId":"trial-1"`, `"requestId":"a","requestId":"b"`, 1), strings.Replace(valid, `"context":{`, `"context":{"Project":"x",`, 1), strings.Replace(valid, `"schema":"sdptool.invoke/1"`, `"schema":"future"`, 1), strings.Replace(valid, `"action":"project.discover"`, `"action":"unknown"`, 1), strings.Repeat(" ", MaxInvokeBytes+1), strings.Repeat("[", MaxInvokeDepth+1) + strings.Repeat("]", MaxInvokeDepth+1), string([]byte{0xff})}
	for i, raw := range tests {
		r := Invoke(context.Background(), strings.NewReader(raw), ActionCaller{})
		if r.ExitCode != 2 || r.Error == nil {
			t.Fatalf("case %d accepted: %+v", i, r)
		}
	}
	for _, args := range [][]string{{"invoke"}, {"invoke", "--request", "-", "--json=false"}, {"invoke", "--request", "-", "--as", "x"}, {root, "invoke", "--request", "-"}} {
		var out, errs bytes.Buffer
		code := RunWithInput(context.Background(), args, strings.NewReader(valid), &out, &errs)
		if code != 2 || !json.Valid(out.Bytes()) {
			t.Fatalf("args %v: %d %s %s", args, code, &out, &errs)
		}
	}
}

type failedWriter struct{}

func (failedWriter) Write([]byte) (int, error) { return 0, errors.New("closed output") }

type shortWriter struct{}

func (shortWriter) Write(p []byte) (int, error) { return len(p) - 1, nil }
func TestInvokeCancellationAndOutputFailure(t *testing.T) {
	root, _ := projectFixture(t)
	raw := requestBytes(t, root, "project.discover", map[string]any{})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	r := Invoke(ctx, bytes.NewReader(raw), ActionCaller{})
	if r.ExitCode != 130 || r.Status != "canceled" {
		t.Fatalf("%+v", r)
	}
	for _, w := range []io.Writer{failedWriter{}, shortWriter{}} {
		var errs bytes.Buffer
		code := RunWithInput(context.Background(), []string{"invoke", "--request", "-"}, bytes.NewReader(raw), w, &errs)
		if code != 1 || errs.Len() == 0 {
			t.Fatalf("%d %s", code, &errs)
		}
	}
}
