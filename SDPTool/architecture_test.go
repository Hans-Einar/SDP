package sdptool

import (
	"bytes"
	"context"
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

// These rules enforce the selected architecture, not inferred SDL semantics.
// Intentional process adapters (controller children, gh-sdp, renderers) remain
// allowed. Domain/application code may not depend on those transport adapters.
func architectureViolations(path string, source []byte) ([]string, error) {
	f, err := parser.ParseFile(token.NewFileSet(), path, source, 0)
	if err != nil {
		return nil, err
	}
	var violations []string
	parts := strings.Split(filepath.ToSlash(path), "/")
	domain := map[string]bool{"activities": true, "governance": true, "governanceapp": true, "model": true, "blueprints": true, "blueprintstate": true, "projecthistory": true}
	for _, im := range f.Imports {
		name, _ := strconv.Unquote(im.Path.Value)
		const prefix = "github.com/Hans-Einar/SDP/SDPTool"
		if domain[parts[0]] && (name == prefix || strings.HasPrefix(name, prefix+"/cmd/") || name == prefix+"/presentation" || name == prefix+"/governancemcp" || name == prefix+"/governancehost" || name == prefix+"/governancecontroller") {
			violations = append(violations, "domain imports adapter: "+name)
		}
		if path == "action_sessions.go" && (name == "bytes" || name == "encoding/json") {
			violations = append(violations, "Session adapter must retain typed results")
		}
		if parts[0] == "activities" && name == prefix+"/governance" {
			violations = append(violations, "Activity domain must use an execution port, not import governance")
		}
	}
	ast.Inspect(f, func(n ast.Node) bool {
		if call, ok := n.(*ast.CallExpr); ok {
			if id, ok := call.Fun.(*ast.Ident); ok && (id.Name == "activityCommand" || id.Name == "governanceCommand" || id.Name == "governanceCommandAs" || id.Name == "RunWithInput") && path != "cli.go" && path != "activity_cli.go" && path != "governance_cli.go" && !strings.HasPrefix(path, "cmd/") {
				violations = append(violations, "non-CLI code calls CLI: "+id.Name)
			}
			if sel, ok := call.Fun.(*ast.SelectorExpr); ok && (sel.Sel.Name == "AddTool" || sel.Sel.Name == "Marshal") && strings.HasPrefix(path, "governancemcp/") && path != "governancemcp/response.go" {
				violations = append(violations, "MCP framing must use shared response policy")
			}
		}
		if b, ok := n.(*ast.BinaryExpr); ok && path == "action_invoke.go" {
			check := func(e ast.Expr) bool { s, ok := e.(*ast.SelectorExpr); return ok && s.Sel.Name == "Action" }
			if check(b.X) || check(b.Y) {
				violations = append(violations, "Invoke must use registered policy, not an action-ID list")
			}
		}
		return true
	})
	return violations, nil
}

func TestArchitectureBoundaries(t *testing.T) {
	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if d.Name() == "testdata" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		v, err := architectureViolations(filepath.ToSlash(path), b)
		if err != nil {
			return err
		}
		if len(v) > 0 {
			t.Errorf("%s: %v", path, v)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestArchitectureChecksRejectViolatingSpecimens(t *testing.T) {
	for _, tc := range []struct{ path, code string }{
		{"governance/bad.go", `package governance; import _ "github.com/Hans-Einar/SDP/SDPTool/presentation"`},
		{"governanceapp/bad.go", `package governanceapp; import _ "github.com/Hans-Einar/SDP/SDPTool"`},
		{"activities/bad.go", `package activities; import _ "github.com/Hans-Einar/SDP/SDPTool/governance"`},
		{"action_activities.go", `package sdptool; func bad(){activityCommand()}`},
		{"action_sessions.go", `package sdptool; func bad(){governanceCommandAs()}`},
		{"action_sessions.go", `package sdptool; import "encoding/json"; func bad(){json.Unmarshal(nil,nil)}`},
		{"action_invoke.go", `package sdptool; func bad(){if req.Action == "privileged" {}}`},
		{"governancemcp/session.go", `package governancemcp; func bad(){s.AddTool(nil,nil)}`},
	} {
		v, e := architectureViolations(tc.path, []byte(tc.code))
		if e != nil || len(v) == 0 {
			t.Fatalf("missed specimen %s: %v %v", tc.path, v, e)
		}
	}
}

func TestMigrationPreservesExistingPublicActionContracts(t *testing.T) {
	b, err := os.ReadFile("testdata/actions-before-migration.json")
	if err != nil {
		t.Fatal(err)
	}
	var before ActionCatalogue
	if err = json.Unmarshal(b, &before); err != nil {
		t.Fatal(err)
	}
	after, err := Actions(context.Background(), "")
	if err != nil {
		t.Fatal(err)
	}
	// Build identity is intentionally separate from stable metadata and wire schemas.
	if runtime.GOOS != "linux" {
		for i := range before.Actions {
			if before.Actions[i].ID == "blueprint.assignment.apply" {
				before.Actions[i].Availability = ActionAvailability{State: "unavailable", Reason: "unsupported-platform", Message: "Assignment writes currently require Linux."}
			}
		}
	}
	// Activity0002 explicitly adds Activity operations. Keep the frozen obligations
	// for every pre-existing action instead of replacing the compatibility fixture.
	retained := []ActionDescriptor{}
	for _, d := range after.Actions {
		if !strings.HasPrefix(d.ID, "activity.") {
			retained = append(retained, d)
		}
	}
	beforeActions, _ := json.Marshal(before.Actions)
	afterActions, _ := json.Marshal(retained)
	if !bytes.Equal(beforeActions, afterActions) || before.Schema != after.Schema || !reflect.DeepEqual(before.InvocationVersions, after.InvocationVersions) {
		t.Fatal("public action contract drifted")
	}
}
