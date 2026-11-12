package sdptool

import (
	"context"
	"fmt"
	"strings"

	"github.com/Hans-Einar/SDP/SDPTool/activities"
)

type ActivityViewInput struct {
	URI      string `json:"uri" nonempty:"true"`
	Revision string `json:"revision,omitempty"`
	Output   string `json:"output" nonempty:"true" pathBase:"project"`
}
type ActivityTreeInput struct {
	URI      string `json:"uri" nonempty:"true"`
	Revision string `json:"revision" nonempty:"true"`
}
type ActivityRequestInput struct {
	RequestFile     string `json:"requestFile" nonempty:"true" pathBase:"project"`
	AuthoritySource string `json:"authoritySource" nonempty:"true"`
}
type ActivityReadInput struct {
	Activity string `json:"activity" nonempty:"true"`
}
type ActivityRecoveryInput struct {
	Operation       string `json:"operation" nonempty:"true"`
	AuthoritySource string `json:"authoritySource" nonempty:"true"`
}

func performActivity(ctx context.Context, c actionContext, op activities.Operation, in activities.Input, source string) (any, int, error) {
	caller := activities.Caller{}
	if op.Write {
		if c.Caller.Principal.Authority != "owner" || strings.TrimSpace(c.Caller.Principal.Actor) == "" {
			return nil, 2, &actionFailure{code: "unauthorized", exit: 2, message: "trusted owner invocation required"}
		}
		caller = activities.Caller{Actor: c.Caller.Principal.Actor, Authority: activities.Authority{Role: "project-owner", Source: source}}
	}
	s, e := activityService(c.Root)
	if e != nil {
		return nil, 2, e
	}
	r, e := op.Execute(ctx, s, caller, in, activityOperationContext(ctx, c.Root, s))
	if e != nil {
		return nil, 2, &actionFailure{code: "activity", exit: 2, message: e.Error()}
	}
	if op.Input == activities.ViewInput || op.Input == activities.TreeInput {
		return r.Data, 0, nil
	}
	return r, 0, nil
}

func activityActions() []actionSpec {
	var specs []actionSpec
	for _, op := range activities.Operations() {
		d := ActionDescriptor{ID: op.ActionID, Label: op.Label, Description: op.Description, Group: "Activities", IconKey: "session", Mutation: "none", OutputSchema: activities.ApplicationSchema}
		if op.Write {
			d.Mutation = "documents"
		}
		policy := actionPolicy{callerAttribution: op.Write}
		if op.Input == activities.ViewInput || op.Input == activities.TreeInput {
			d.OutputSchema = Version
		}
		if op.Input == activities.ViewInput {
			d.Mutation = "documents"
		}
		switch op.Input {
		case activities.ViewInput:
			specs = append(specs, defineAction(d, func(ctx context.Context, c actionContext, q ActivityViewInput) (any, int, error) {
				return performActivity(ctx, c, op, activities.Input{URI: q.URI, Revision: q.Revision, Output: actionPath(c.Root, q.Output)}, "")
			}, policy))
		case activities.TreeInput:
			specs = append(specs, defineAction(d, func(ctx context.Context, c actionContext, q ActivityTreeInput) (any, int, error) {
				return performActivity(ctx, c, op, activities.Input{URI: q.URI, Revision: q.Revision}, "")
			}, policy))
		case activities.NoInput:
			specs = append(specs, defineAction(d, func(ctx context.Context, c actionContext, _ noActionParameters) (any, int, error) {
				return performActivity(ctx, c, op, activities.Input{}, "")
			}, policy))
		case activities.ActivityInput:
			specs = append(specs, defineAction(d, func(ctx context.Context, c actionContext, q ActivityReadInput) (any, int, error) {
				return performActivity(ctx, c, op, activities.Input{ActivityID: q.Activity}, "")
			}, policy))
		case activities.RequestInput:
			specs = append(specs, defineAction(d, func(ctx context.Context, c actionContext, q ActivityRequestInput) (any, int, error) {
				b, e := boundedFile(actionPath(c.Root, q.RequestFile))
				if e != nil {
					return nil, 2, e
				}
				return performActivity(ctx, c, op, activities.Input{Request: b}, q.AuthoritySource)
			}, policy))
		case activities.RecoveryInput:
			specs = append(specs, defineAction(d, func(ctx context.Context, c actionContext, q ActivityRecoveryInput) (any, int, error) {
				return performActivity(ctx, c, op, activities.Input{OperationID: q.Operation}, q.AuthoritySource)
			}, policy))
		default:
			panic(fmt.Sprintf("unmapped Activity input kind %s", op.Input))
		}
	}
	return specs
}
