package sdptool

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Hans-Einar/SDP/SDUI/go/preview"
)

type sduiSourceParameters struct {
	Source   string `json:"source"`
	Entry    string `json:"entry"`
	Revision string `json:"revision,omitempty"`
	Output   string `json:"output,omitempty"`
}
type sduiRunParameters struct {
	Source   string `json:"source"`
	Entry    string `json:"entry"`
	Revision string `json:"revision,omitempty"`
	Watch    *bool  `json:"watch,omitempty"`
}

func runtimeActions() []actionSpec {
	policy := actionPolicy{hostContext: true}
	return []actionSpec{
		defineAction(ActionDescriptor{ID: "runtime.setup", Label: "Register shared runtime", Description: "Explicitly register this installed engine for direct local consumers.", Group: "Runtime", IconKey: "settings", Mutation: "host-files", OutputSchema: Version}, performSetup, policy),
		defineAction(ActionDescriptor{ID: "sdui.source.preview", Label: "Preview SDUI source", Description: "Render a checked source into a caller-owned empty directory.", Group: "SDUI", IconKey: "preview", Mutation: "documents", OutputSchema: Version}, performSduiPreview, policy),
		defineAction(ActionDescriptor{ID: "sdui.source.check", Label: "Check SDUI prototype", Description: "Check local native prototype readiness without launching it.", Group: "SDUI", IconKey: "check", Mutation: "none", OutputSchema: Version}, performSduiCheck, policy),
		defineAction(ActionDescriptor{ID: "sdui.source.run", Label: "Run SDUI prototype", Description: "Run the native host bundled with this exact engine, offline.", Group: "SDUI", IconKey: "run", Mutation: "host-process", OutputSchema: Version}, performSduiRun, policy),
	}
}
func performSetup(ctx context.Context, _ actionContext, _ noActionParameters) (any, int, error) {
	r, e := RegisterRuntime(ctx)
	return r, 0, e
}
func performSduiPreview(ctx context.Context, _ actionContext, p sduiSourceParameters) (any, int, error) {
	return performSource(ctx, p, false)
}
func performSduiCheck(ctx context.Context, _ actionContext, p sduiSourceParameters) (any, int, error) {
	return performSource(ctx, p, true)
}
func performSource(ctx context.Context, p sduiSourceParameters, check bool) (any, int, error) {
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}
	if check && p.Output != "" {
		return nil, 0, fmt.Errorf("check does not accept an output directory")
	}
	args := []string{"-source", p.Source, "-entry", p.Entry, "-revision", p.Revision}
	if check {
		args = append(args, "-check")
	} else {
		args = append(args, "-output", p.Output)
	}
	var out bytes.Buffer
	if err := preview.Execute(args, &out); err != nil {
		return nil, 0, err
	}
	var result any
	err := json.Unmarshal(out.Bytes(), &result)
	return result, 0, err
}
func performSduiRun(ctx context.Context, _ actionContext, p sduiRunParameters) (any, int, error) {
	args := []string{"-entry", p.Entry, "-revision", p.Revision}
	watch := true
	if p.Watch != nil {
		watch = *p.Watch
	}
	args = append(args, fmt.Sprintf("-watch=%t", watch), p.Source)
	err := runNativeHost(ctx, args)
	return map[string]any{"schema": Version, "operation": "sdui-run", "version": BuildVersion, "revision": BuildRevision}, 0, err
}
