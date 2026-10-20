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

	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/parser"
)

// The merged Actions producer must expose executable overview targets through
// both discovery interfaces without dropping the existing detail targets.
func TestInvokedDiscoveryPreservesOverviewTargets(t *testing.T) {
	root, _ := projectFixture(t)
	m, err := parser.Parse(`language design-core version 0.5.
actor Reader.
usecase Browse.
unit Documents.
container Desktop.
feature Navigation.
functionality Display.
mode Online.
Reader pursues Browse.
Documents owns Display.
Display contributes-to Navigation.
Display allocated-to Desktop in mode Online.
`)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join(root, "SDP", "overview.design"), []byte(parser.FormatFile(m)), 0600); err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	var direct, invoked, diagnostics bytes.Buffer
	if code := Run(ctx, []string{root, "discover", "--json"}, &direct, &diagnostics); code != 0 {
		t.Fatal(code, diagnostics.String())
	}
	request := requestBytes(t, root, "project.discover", map[string]any{})
	if code := RunWithInput(ctx, []string{"invoke", "--request", "-"}, bytes.NewReader(request), &invoked, &diagnostics); code != 0 {
		t.Fatal(code, invoked.String(), diagnostics.String())
	}
	var envelope struct {
		Result json.RawMessage `json:"result"`
	}
	if err = json.Unmarshal(invoked.Bytes(), &envelope); err != nil {
		t.Fatal(err)
	}
	var a, b any
	if err = json.Unmarshal(direct.Bytes(), &a); err != nil {
		t.Fatal(err)
	}
	if err = json.Unmarshal(envelope.Result, &b); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a, b) {
		t.Fatal("invoke altered the discovery payload")
	}
	var discovery struct {
		Navigation Tree `json:"navigation"`
	}
	if err = json.Unmarshal(envelope.Result, &discovery); err != nil {
		t.Fatal(err)
	}
	p, err := Discover(root)
	if err != nil {
		t.Fatal(err)
	}
	count := 0
	for _, node := range discovery.Navigation.Nodes {
		target := node.OverviewTarget
		if target == nil && strings.Contains(node.ID, "/overview-mode/") {
			target = node.Target
		}
		if target == nil {
			continue
		}
		if node.OverviewTarget != nil && (node.Target == nil || node.Target.URI == target.URI) {
			t.Fatal("lost original detail target", node)
		}
		result, err := SelectProject(ctx, p, target.Model, PreviewOptions{URI: target.URI, Revision: target.Revision, Output: filepath.Join(t.TempDir(), "selected")})
		if err != nil {
			t.Fatal(target, err)
		}
		if _, err = os.Stat(result.Entry); err != nil {
			t.Fatal(err)
		}
		count++
	}
	if count != 12 {
		t.Fatalf("want VP01–VP09 overviews and VP05/VP07/VP08 Online targets, got %d", count)
	}
}
