// Package governanceapp owns the shared typed governance operation inventory.
package governanceapp

import (
	"context"
	"fmt"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
)

type InputKind int

const (
	NoInput InputKind = iota
	RequestInput
	RunInput
	InstanceInput
)

// Action is transport-neutral metadata; adapters own their wire schemas.
type Action struct{ ID, Label, Description, Mutation, OutputSchema string }
type Operation struct {
	Name    string
	Input   InputKind
	Write   bool
	Action  *Action
	execute func(*governance.Service, governance.Caller, []byte, string) (any, error)
}

func (op Operation) Validate(input, run, instance string) error {
	if instance != "" && op.Input != InstanceInput {
		return fmt.Errorf("--instance is only valid for session-read")
	}
	switch op.Input {
	case NoInput:
		if input != "" || run != "" {
			if op.Name == "bind" {
				return fmt.Errorf("bind takes no request/run")
			}
			return fmt.Errorf("operation takes no request/run")
		}
	case InstanceInput:
		if input != "" || run != "" || instance == "" {
			return fmt.Errorf("session-read requires --instance only")
		}
	case RunInput:
		if input != "" || run == "" {
			return fmt.Errorf("read requires --run only")
		}
	case RequestInput:
		if input == "" || run != "" {
			return fmt.Errorf("operation requires --input only")
		}
	}
	return nil
}
func (op Operation) Execute(ctx context.Context, s *governance.Service, c governance.Caller, b []byte, identity string) (any, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	// Never reinterpret a committed mutation as canceled after the handler returns.
	return op.execute(s, c, b, identity)
}
func request[Q, R any](name string, write bool, call func(*governance.Service, governance.Caller, Q) (R, error)) Operation {
	return Operation{Name: name, Input: RequestInput, Write: write, execute: func(s *governance.Service, c governance.Caller, b []byte, _ string) (any, error) {
		var q Q
		if err := governance.Decode(b, &q); err != nil {
			return nil, err
		}
		return call(s, c, q)
	}}
}
func readRequest[Q, R any](name string, call func(*governance.Service, Q) (R, error)) Operation {
	return request(name, false, func(s *governance.Service, _ governance.Caller, q Q) (R, error) { return call(s, q) })
}
func noInput[R any](name string, write bool, call func(*governance.Service, governance.Caller) (R, error)) Operation {
	return Operation{Name: name, Input: NoInput, Write: write, execute: func(s *governance.Service, c governance.Caller, _ []byte, _ string) (any, error) { return call(s, c) }}
}
func identity[R any](name string, kind InputKind, call func(*governance.Service, string) (R, error)) Operation {
	return Operation{Name: name, Input: kind, execute: func(s *governance.Service, _ governance.Caller, _ []byte, id string) (any, error) { return call(s, id) }}
}

type CancelRequest struct {
	OperationID string `json:"operationId"`
}

func Operations() []Operation {
	out := []Operation{
		expose(noInput("session-definition", false, func(_ *governance.Service, _ governance.Caller) (any, error) {
			w := governance.SessionWorkflow()
			return map[string]any{"schema": governance.SessionSchema, "operation": "session-definition", "definition": w, "digest": w.Digest()}, nil
		}), Action{ID: "governance.session.definition", Label: "Session workflow definition", Description: "Read the supported pinned workflow; no execution.", Mutation: "none", OutputSchema: "sdp-session-execution/0.1"}),
		expose(noInput("session-list", false, func(s *governance.Service, _ governance.Caller) (any, error) { return s.ListSessions() }), Action{ID: "governance.session.list", Label: "Session instances", Description: "Inspect instances in one selected local store; no discovery of other stores.", Mutation: "none", OutputSchema: "sdp-session-execution/0.1"}),
		expose(identity("session-read", InstanceInput, (*governance.Service).ReadSession), Action{ID: "governance.session.read", Label: "Inspect Session instance", Description: "Read stages, unmet conditions and coverage without starting work.", Mutation: "none", OutputSchema: "sdp-session-execution/0.1"}),
		expose(identity("session-host-read", InstanceInput, (*governance.Service).ReadSessionHostCommands), Action{ID: "governance.session.host-read", Label: "Session host-read", Description: "Inspect, queue or cancel a durable host request; execution belongs to the trusted controller.", Mutation: "none", OutputSchema: "sdp-session-host/0.1"}),
		expose(request("session-host-request", true, func(s *governance.Service, c governance.Caller, q governance.SessionHostRequest) (any, error) {
			r, e := s.RequestSessionHost(c, q)
			return map[string]any{"schema": governance.SessionHostSchema, "operation": "session-host-request", "command": r}, e
		}), Action{ID: "governance.session.host-request", Label: "Session host-request", Description: "Inspect, queue or cancel a durable host request; execution belongs to the trusted controller.", Mutation: "state", OutputSchema: "sdp-session-host/0.1"}),
		expose(request("session-host-cancel", true, func(s *governance.Service, c governance.Caller, q CancelRequest) (any, error) {
			r, e := s.CancelSessionHost(c, q.OperationID)
			return map[string]any{"schema": governance.SessionHostSchema, "operation": "session-host-cancel", "command": r}, e
		}), Action{ID: "governance.session.host-cancel", Label: "Session host-cancel", Description: "Inspect, queue or cancel a durable host request; execution belongs to the trusted controller.", Mutation: "state", OutputSchema: "sdp-session-host/0.1"}),
		expose(readRequest("session-context", (*governance.Service).SessionContext), Action{ID: "governance.session.context", Label: "Session context", Description: "Read and validate the binding and canonical coverage.", Mutation: "none", OutputSchema: "sdp-session-execution/0.1"}),
		expose(readRequest("session-design-context", (*governance.Service).InspectSessionDesign), Action{ID: "governance.session.design-context", Label: "Session design-context", Description: "Inspect exact assessed SDL WORK, task and required Blueprint; no design acceptance or dispatch.", Mutation: "none", OutputSchema: "sdp-session-execution/0.1"}),
		expose(readRequest("session-rebind-context", (*governance.Service).SessionRebindContext), Action{ID: "governance.session.rebind-context", Label: "Session rebind-context", Description: "Read and validate the binding and canonical coverage.", Mutation: "none", OutputSchema: "sdp-session-execution/0.1"}),
		expose(readRequest("session-successor", (*governance.Service).SessionSuccessorContext), Action{ID: "governance.session.successor", Label: "Session successor", Description: "Read and validate the binding and canonical coverage.", Mutation: "none", OutputSchema: "sdp-session-execution/0.1"}),
		expose(request("session-rebind", true, (*governance.Service).RebindSession), Action{ID: "governance.session.rebind", Label: "Session rebind", Description: "Prepare an exact reviewed binding revision at a quiescent boundary; retain earlier evidence and require fresh selection.", Mutation: "state", OutputSchema: "sdp-session-execution/0.1"}),
		expose(request("session-adopt", true, (*governance.Service).AdoptSession), Action{ID: "governance.session.adopt", Label: "Session adopt", Description: "Persist adoption intent; no card activation or agent dispatch.", Mutation: "state", OutputSchema: "sdp-session-execution/0.1"}),
		expose(request("session-resume", true, (*governance.Service).ResumeSession), Action{ID: "governance.session.resume", Label: "Session resume", Description: "Reconcile Session registration and coverage; no agent dispatch or later-stage execution.", Mutation: "documents", OutputSchema: "sdp-session-execution/0.1"}),
		expose(request("session-select", true, (*governance.Service).SelectSession), Action{ID: "governance.session.select", Label: "Session select", Description: "Prepare checked Session step selection and ready-card effects.", Mutation: "state", OutputSchema: "sdp-session-execution/0.1"}),
		expose(request("session-assess", true, (*governance.Service).AssessSession), Action{ID: "governance.session.assess", Label: "Session assess", Description: "Prepare an attributed design assessment and checked canonical records; no host dispatch.", Mutation: "state", OutputSchema: "sdp-session-execution/0.1"}),
		expose(request("session-advance", true, (*governance.Service).AdvanceSession), Action{ID: "governance.session.advance", Label: "Session advance", Description: "Prepare checked rework, next-step continuation or overall Session closure with retained history.", Mutation: "state", OutputSchema: "sdp-session-execution/0.1"}),
		expose(request("session-owner", true, (*governance.Service).OwnerSession), Action{ID: "governance.session.owner", Label: "Session owner", Description: "Request concrete owner review or record an owner-authorized exact-result decision.", Mutation: "state", OutputSchema: "sdp-session-execution/0.1"}),
		expose(request("session-dispose", true, (*governance.Service).DisposeSession), Action{ID: "governance.session.dispose", Label: "Session dispose", Description: "Prepare checked selected-step and card disposition; resume publishes results and releases the claim.", Mutation: "state", OutputSchema: "sdp-session-execution/0.1"}),
		expose(Operation{Name: "session-start", Input: RequestInput, Write: true, execute: func(s *governance.Service, c governance.Caller, b []byte, _ string) (any, error) {
			q, e := governance.DecodeSessionWork(b, "start")
			if e != nil {
				return nil, e
			}
			return s.StartSession(c, q)
		}}, Action{ID: "governance.session.start", Label: "Session start", Description: "Prepare Session-bound in-progress or candidate progress records; resume commits them.", Mutation: "state", OutputSchema: "sdp-session-execution/0.1"}),
		expose(Operation{Name: "session-progress", Input: RequestInput, Write: true, execute: func(s *governance.Service, c governance.Caller, b []byte, _ string) (any, error) {
			q, e := governance.DecodeSessionWork(b, "progress")
			if e != nil {
				return nil, e
			}
			return s.ProgressSession(c, q)
		}}, Action{ID: "governance.session.progress", Label: "Session progress", Description: "Prepare Session-bound in-progress or candidate progress records; resume commits them.", Mutation: "state", OutputSchema: "sdp-session-execution/0.1"}),
		expose(Operation{Name: "session-verify", Input: RequestInput, Write: true, execute: func(s *governance.Service, c governance.Caller, b []byte, _ string) (any, error) {
			q, e := governance.DecodeSessionWork(b, "verify")
			if e != nil {
				return nil, e
			}
			return s.VerifySession(c, q)
		}}, Action{ID: "governance.session.verify", Label: "Session verify", Description: "Prepare exact candidate test or observed review intake; resume commits canonical records.", Mutation: "state", OutputSchema: "sdp-session-execution/0.1"}),
		expose(Operation{Name: "session-review", Input: RequestInput, Write: true, execute: func(s *governance.Service, c governance.Caller, b []byte, _ string) (any, error) {
			q, e := governance.DecodeSessionWork(b, "review")
			if e != nil {
				return nil, e
			}
			return s.ReviewSession(c, q)
		}}, Action{ID: "governance.session.review", Label: "Session review", Description: "Prepare exact candidate test or observed review intake; resume commits canonical records.", Mutation: "state", OutputSchema: "sdp-session-execution/0.1"}),
		readRequest("context", (*governance.Service).Context),
		readRequest("route", (*governance.Service).Route),
		noInput("bind", true, (*governance.Service).Bind),
		request("start", true, (*governance.Service).Start),
		identity("read", RunInput, (*governance.Service).Read),
		request("transition", true, (*governance.Service).Transition),
	}
	return out
}
func expose(op Operation, action Action) Operation { op.Action = &action; return op }

func Find(name string) (Operation, bool) {
	for _, op := range Operations() {
		if op.Name == name {
			return op, true
		}
	}
	return Operation{}, false
}
