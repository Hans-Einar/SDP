package sdptool

import (
	"context"
	"errors"
	"github.com/Hans-Einar/SDP/SDPTool/governance"
	"github.com/Hans-Einar/SDP/SDPTool/governanceapp"
	"strings"
)

type SessionActionInput struct {
	State       string `json:"state,omitempty" pathBase:"project"`
	RequestFile string `json:"requestFile" nonempty:"true" pathBase:"project"`
}
type SessionReadInput struct {
	State    string `json:"state,omitempty" pathBase:"project"`
	Instance string `json:"instance" nonempty:"true"`
}
type SessionListInput struct {
	State string `json:"state,omitempty" pathBase:"project"`
}

// Authority is supplied by the invocation host, never by request JSON.
func sessionActionCall(ctx context.Context, c actionContext, op governanceapp.Operation, state, input, instance string) (any, int, error) {
	if e := ctx.Err(); e != nil {
		return nil, 130, e
	}
	if op.Write && ((c.Caller.Principal.Authority != "controller" && c.Caller.Principal.Authority != "owner") || strings.TrimSpace(c.Caller.Principal.Actor) == "") {
		return nil, 2, &actionFailure{code: "unauthorized", exit: 2, message: "trusted controller principal required"}
	}
	caller := governance.LocalController()
	if c.Caller.Principal.Authority == "owner" {
		caller = governance.LocalOwner()
	}
	if state != "" {
		state = actionPath(c.Root, state)
	}
	service, e := governanceService(c.Root, state)
	var b []byte
	if e == nil && input != "" {
		b, e = boundedFile(actionPath(c.Root, input))
	}
	var result any
	if e == nil {
		result, e = op.Execute(ctx, service, caller, b, instance)
	}
	if e != nil {
		if errors.Is(e, context.Canceled) || errors.Is(e, context.DeadlineExceeded) {
			return nil, 130, e
		}
		var domain *governance.Error
		if !errors.As(e, &domain) {
			domain = &governance.Error{Code: "unavailable", Message: e.Error()}
		}
		return nil, 2, &actionFailure{code: domain.Code, exit: 2, message: domain.Message}
	}
	return result, 0, nil
}
func sessionActions() []actionSpec {
	var out []actionSpec
	for _, op := range governanceapp.Operations() {
		if op.Action == nil {
			continue
		}
		a := op.Action
		descriptor := ActionDescriptor{ID: a.ID, Label: a.Label, Description: a.Description, Group: "ProjectGovernance", IconKey: "session", Mutation: a.Mutation, OutputSchema: a.OutputSchema}
		policy := actionPolicy{callerAttribution: op.Write}
		switch op.Input {
		case governanceapp.NoInput:
			out = append(out, defineAction(descriptor, func(ctx context.Context, c actionContext, _ noActionParameters) (any, int, error) {
				return sessionActionCall(ctx, c, op, "", "", "")
			}, policy))
		case governanceapp.StateInput:
			out = append(out, defineAction(descriptor, func(ctx context.Context, c actionContext, q SessionListInput) (any, int, error) {
				return sessionActionCall(ctx, c, op, q.State, "", "")
			}, policy))
		case governanceapp.InstanceInput:
			out = append(out, defineAction(descriptor, func(ctx context.Context, c actionContext, q SessionReadInput) (any, int, error) {
				return sessionActionCall(ctx, c, op, q.State, "", q.Instance)
			}, policy))
		case governanceapp.RequestInput:
			out = append(out, defineAction(descriptor, func(ctx context.Context, c actionContext, q SessionActionInput) (any, int, error) {
				return sessionActionCall(ctx, c, op, q.State, q.RequestFile, "")
			}, policy))
		}
	}
	return out
}
