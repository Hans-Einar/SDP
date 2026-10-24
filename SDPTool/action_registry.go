package sdptool

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
)

type actionContext struct {
	Root, Area string
	Caller     ActionCaller
}
type actionSpec struct {
	descriptor ActionDescriptor
	decode     func(json.RawMessage) (any, error)
	run        func(context.Context, actionContext, any) (any, int, error)
}

func defineAction[T any](d ActionDescriptor, run func(context.Context, actionContext, T) (any, int, error)) actionSpec {
	s := actionSchema(reflect.TypeOf((*T)(nil)).Elem())
	if d.InputSchema != nil {
		s = d.InputSchema
	}
	d.InputSchema = s
	return actionSpec{d, func(raw json.RawMessage) (any, error) { return decodeAction[T](raw, s) }, func(ctx context.Context, c actionContext, v any) (any, int, error) { return run(ctx, c, v.(T)) }}
}

type noActionParameters struct{}

func actionRegistry() []actionSpec {
	return append(append([]actionSpec{discoveryAction(), modelTreeAction()}, blueprintActions()...), sessionActions()...)
}
func performDiscover(ctx context.Context, c actionContext, options DiscoverOptions) (any, int, error) {
	if err := ctx.Err(); err != nil {
		return nil, 130, err
	}
	r, e := DiscoverWithOptions(c.Root, options)
	return r, 0, e
}
func resolveActionContext(selected string, caller ActionCaller) (actionContext, error) {
	p, err := discover(selected, false)
	if err != nil {
		return actionContext{}, err
	}
	return actionContext{Root: p.Root, Area: p.Area, Caller: caller}, nil
}
func actionPath(root, path string) string {
	if filepath.IsAbs(path) {
		return filepath.Clean(path)
	}
	return filepath.Join(root, path)
}

// Actions returns only executable registry entries. Availability is advisory and read-only.
func Actions(ctx context.Context, selected string) (ActionCatalogue, error) {
	r := ActionCatalogue{Schema: ActionsSchema, Operation: "actions", InvocationVersions: []string{InvokeSchema}, Actions: []ActionDescriptor{}}
	r.Producer.Version = BuildVersion
	r.Producer.Revision = BuildRevision
	var contextErr error
	if selected != "" {
		_, contextErr = resolveActionContext(selected, ActionCaller{})
	}
	if err := ctx.Err(); err != nil {
		return r, err
	}
	specs := actionRegistry()
	sort.Slice(specs, func(i, j int) bool { return specs[i].descriptor.ID < specs[j].descriptor.ID })
	stable := []ActionDescriptor{}
	for _, s := range specs {
		d := s.descriptor
		stable = append(stable, d)
		d.Availability = ActionAvailability{State: "requires-input", Reason: "parameters-required", Message: "Supply the declared inputs; invocation checks current state."}
		if len(d.InputSchema.Required) == 0 && selected != "" && contextErr == nil {
			d.Availability = ActionAvailability{State: "available", Reason: "context-resolved", Message: "Project context resolved; invocation checks it again."}
		}
		if contextErr != nil {
			d.Availability = ActionAvailability{State: "unavailable", Reason: "invalid-context", Message: contextErr.Error()}
		}
		if d.ID == "blueprint.assignment.apply" && runtime.GOOS != "linux" {
			d.Availability = ActionAvailability{State: "unavailable", Reason: "unsupported-platform", Message: "Assignment writes currently require Linux."}
		}
		r.Actions = append(r.Actions, d)
	}
	b, err := json.Marshal(stable)
	if err != nil {
		return r, err
	}
	sum := sha256.Sum256(b)
	r.Revision = hex.EncodeToString(sum[:])
	if contextErr != nil {
		r.Diagnostics = []InvokeError{{Code: "invalid-context", Message: contextErr.Error()}}
	}
	return r, nil
}
func findAction(id string) (actionSpec, error) {
	for _, s := range actionRegistry() {
		if s.descriptor.ID == id {
			return s, nil
		}
	}
	return actionSpec{}, fmt.Errorf("unsupported action %q", id)
}
