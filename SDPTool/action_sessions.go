package sdptool

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
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

// Use exactly the same strict domain decoder and service as the trusted CLI.
// Authority is supplied by the invocation host, never by request JSON.
func sessionActionCall(ctx context.Context, c actionContext, op, state, input, instance string, write bool) (any, int, error) {
	if e := ctx.Err(); e != nil {
		return nil, 130, e
	}
	if write && (c.Caller.Principal.Authority != "controller" || strings.TrimSpace(c.Caller.Principal.Actor) == "") {
		return nil, 2, &actionFailure{code: "unauthorized", exit: 2, message: "trusted controller principal required"}
	}
	args := []string{op}
	if state != "" {
		args = append(args, "--state", actionPath(c.Root, state))
	}
	if input != "" {
		args = append(args, "--input", actionPath(c.Root, input))
	}
	if instance != "" {
		args = append(args, "--instance", instance)
	}
	var out, errs bytes.Buffer
	code := governanceCommand(c.Root, args, &out, &errs, true)
	if code != 0 {
		var envelope struct {
			Error governance.Error `json:"error"`
		}
		if e := json.Unmarshal(errs.Bytes(), &envelope); e == nil && envelope.Error.Code != "" {
			return nil, code, &actionFailure{code: envelope.Error.Code, exit: code, message: envelope.Error.Message}
		}
		return nil, code, fmt.Errorf("Session operation failed: %s", errs.String())
	}
	var result any
	if e := json.Unmarshal(out.Bytes(), &result); e != nil {
		return nil, 2, e
	}
	return result, 0, nil
}
func sessionActions() []actionSpec {
	descriptor := func(id, label, desc, mutation string) ActionDescriptor {
		return ActionDescriptor{ID: "governance.session." + id, Label: label, Description: desc, Group: "ProjectGovernance", IconKey: "session", Mutation: mutation, OutputSchema: governance.SessionSchema}
	}
	out := []actionSpec{
		defineAction(descriptor("definition", "Session workflow definition", "Read the supported pinned workflow; no execution.", "none"), func(ctx context.Context, c actionContext, _ noActionParameters) (any, int, error) {
			return sessionActionCall(ctx, c, "session-definition", "", "", "", false)
		}),
		defineAction(descriptor("list", "Session instances", "Inspect instances in one selected local store; no discovery of other stores.", "none"), func(ctx context.Context, c actionContext, q SessionListInput) (any, int, error) {
			return sessionActionCall(ctx, c, "session-list", q.State, "", "", false)
		}),
		defineAction(descriptor("read", "Inspect Session instance", "Read stages, unmet conditions and coverage without starting work.", "none"), func(ctx context.Context, c actionContext, q SessionReadInput) (any, int, error) {
			return sessionActionCall(ctx, c, "session-read", q.State, "", q.Instance, false)
		}),
	}
	for _, op := range []string{"context", "adopt", "resume", "select", "assess", "start", "progress"} {
		mutation := "none"
		desc := "Read and validate the binding and canonical coverage."
		if op == "adopt" {
			mutation = "state"
			desc = "Persist adoption intent; no card activation or agent dispatch."
		}
		if op == "start" || op == "progress" {
			mutation = "state"
			desc = "Prepare Session-bound in-progress or candidate progress records; resume commits them."
		}
		if op == "assess" {
			mutation = "state"
			desc = "Prepare an attributed design assessment and checked canonical records; no host dispatch."
		}
		if op == "select" {
			mutation = "state"
			desc = "Prepare checked Session step selection and ready-card effects."
		}
		if op == "resume" {
			mutation = "documents"
			desc = "Reconcile Session registration and coverage; no agent dispatch or later-stage execution."
		}
		out = append(out, defineAction(descriptor(op, "Session "+op, desc, mutation), func(ctx context.Context, c actionContext, q SessionActionInput) (any, int, error) {
			return sessionActionCall(ctx, c, "session-"+op, q.State, q.RequestFile, "", op != "context")
		}))
	}
	return out
}
