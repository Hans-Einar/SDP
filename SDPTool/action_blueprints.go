package sdptool

import (
	"context"
	"reflect"

	"github.com/Hans-Einar/SDP/SDPTool/blueprints"
	"github.com/Hans-Einar/SDP/SDPTool/blueprintstate"
)

type GenerateActionInput struct {
	ModelArea string `json:"modelArea" nonempty:"true" pathBase:"project"`
	From      string `json:"from" nonempty:"true"`
	To        string `json:"to" nonempty:"true"`
	Entry     string `json:"entry" nonempty:"true" pathBase:"model-snapshot"`
	Task      string `json:"task" nonempty:"true" pathBase:"project"`
	Output    string `json:"output,omitempty" nonempty:"true" pathBase:"project"`
	Catalogue string `json:"catalogue,omitempty" nonempty:"true" pathBase:"project"`
}
type AssessActionInput struct {
	Bundle   string `json:"bundle" nonempty:"true" pathBase:"project"`
	Evidence string `json:"evidence" nonempty:"true" pathBase:"project"`
}
type ApplyAssignmentInput struct {
	Request blueprintstate.Request `json:"request"`
}

func blueprintActions() []actionSpec {
	generation := actionSchema(reflect.TypeOf(GenerateActionInput{}))
	generation.OneOf = []*InputSchema{{Required: []string{"output"}}, {Required: []string{"catalogue"}}}
	assignment := actionSchema(reflect.TypeOf(ApplyAssignmentInput{}))
	binding := assignment.Properties["request"].Properties["binding"]
	for _, key := range []string{"bundle", "modelArea"} {
		binding.Properties[key].PathBase = "sdp-area"
	}
	assignment.Properties["request"].Properties["evidence"].PathBase = "sdp-area"
	return []actionSpec{
		defineAction(ActionDescriptor{ID: "blueprint.generate", Label: "Generate blueprint", Description: "Capture model inputs and generate a preview or immutable retained revision.", Group: "Blueprints", IconKey: "blueprint-generate", Mutation: "documents", InputSchema: generation, OutputSchema: "sdp-blueprint/1"}, performGenerate),
		defineAction(ActionDescriptor{ID: "blueprint.assess", Label: "Assess blueprint", Description: "Assess pinned evidence without executing checks or changing workflow state.", Group: "Blueprints", IconKey: "blueprint-assess", Mutation: "none", OutputSchema: "sdp-blueprint-assessment/1"}, performAssess),
		defineAction(ActionDescriptor{ID: "blueprint.assignment.list", Label: "List assignments", Description: "Read revision-bound assignments from canonical project history.", Group: "Blueprints", IconKey: "blueprint-assignments", Mutation: "none", OutputSchema: blueprintstate.Schema}, performAssignmentList),
		defineAction(ActionDescriptor{ID: "blueprint.assignment.apply", Label: "Apply assignment transition", Description: "Apply one attributed revision-bound transition; requires a trusted local principal.", Group: "Blueprints", IconKey: "blueprint-transition", Mutation: "history", InputSchema: assignment, OutputSchema: blueprintstate.Schema}, performAssignmentApply),
	}
}
func performGenerate(ctx context.Context, c actionContext, p GenerateActionInput) (any, int, error) {
	if err := ctx.Err(); err != nil {
		return nil, 130, err
	}
	if p.ModelArea == "" || p.From == "" || p.To == "" || p.Entry == "" || p.Task == "" || (p.Output == "") == (p.Catalogue == "") {
		return nil, 2, &actionFailure{code: "invalid-parameters", exit: 2, message: "modelArea, from, to, entry, task and exactly one output/catalogue are required"}
	}
	opts := blueprints.Options{Area: actionPath(c.Root, p.ModelArea), From: p.From, To: p.To, Entry: p.Entry, Task: actionPath(c.Root, p.Task)}
	if p.Catalogue != "" {
		r, e := blueprints.GenerateRetained(ctx, opts, actionPath(c.Root, p.Catalogue))
		return r, 0, e
	}
	opts.Output = actionPath(c.Root, p.Output)
	r, e := blueprints.Generate(ctx, opts)
	return r, 0, e
}
func performAssess(ctx context.Context, c actionContext, p AssessActionInput) (any, int, error) {
	if err := ctx.Err(); err != nil {
		return nil, 130, err
	}
	r, e := blueprints.Assess(actionPath(c.Root, p.Bundle), actionPath(c.Root, p.Evidence))
	if e != nil {
		return nil, 1, e
	}
	if r.Status != "ready" {
		return r, 3, nil
	}
	return r, 0, nil
}
func performAssignmentList(ctx context.Context, c actionContext, _ noActionParameters) (any, int, error) {
	if err := ctx.Err(); err != nil {
		return nil, 130, err
	}
	views, e := blueprints.AssignmentViews(c.Area)
	return map[string]any{"schema": blueprintstate.Schema, "operation": "assignment-list", "assignments": views}, 0, e
}
func performAssignmentApply(ctx context.Context, c actionContext, p ApplyAssignmentInput) (any, int, error) {
	if err := blueprintstate.ValidateRequest(p.Request); err != nil {
		return nil, 2, &actionFailure{code: "invalid-parameters", exit: 2, message: err.Error()}
	}
	principal := c.Caller.Principal
	if principal.Actor == "" || (principal.Authority != "controller" && principal.Authority != "assignee" && principal.Authority != "reviewer") {
		return nil, 1, &actionFailure{code: "authority-denied", exit: 1, message: "assignment mutation requires a separately supplied local principal"}
	}
	if err := ctx.Err(); err != nil {
		return nil, 130, err
	}
	// No post-write cancellation check: a committed event must never be reported as uncommitted.
	r, e := blueprints.ApplyAssignment(c.Area, principal, p.Request)
	return r, 0, e
}
