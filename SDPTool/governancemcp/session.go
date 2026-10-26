package governancemcp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// NewSessionController grants only the explicitly launch-selected instance. It
// does not extend the ordinary coordinator/worker capability or dispatch hosts.
func NewSessionController(service *governance.Service, instance string, skill governance.SkillInput) (*mcp.Server, error) {
	if service == nil {
		return nil, fmt.Errorf("missing governance service")
	}
	if e := service.CheckSessionLeader(instance, skill); e != nil {
		return nil, e
	}
	server := mcp.NewServer(&mcp.Implementation{Name: "sdptool-session-controller", Version: "PGF2-C1a"}, &mcp.ServerOptions{SupportedProtocolVersions: []string{ProtocolVersion}})
	check := func() error { return service.CheckSessionLeader(instance, skill) }
	empty := map[string]any{"type": "object", "properties": map[string]any{}, "additionalProperties": false}
	addSessionTool(server, "session_read", "Read the launch-bound Session, stages, holds and history. No host dispatch or owner acceptance is inferred.", empty, true, func(b []byte) (any, error) {
		var q struct{}
		if e := governance.Decode(b, &q); e != nil {
			return nil, e
		}
		if e := check(); e != nil {
			return nil, e
		}
		return service.ReadSession(instance)
	})
	addSessionTool(server, "session_definition", "Read the pinned supported routine and mandatory launcher-pinned skill reference. Load and apply that skill; its presence does not prove consumption.", empty, true, func(b []byte) (any, error) {
		var q struct{}
		if e := governance.Decode(b, &q); e != nil {
			return nil, e
		}
		if e := check(); e != nil {
			return nil, e
		}
		return map[string]any{"workflow": governance.SessionWorkflow(), "skill": skill, "instanceId": instance}, nil
	})
	addSessionTool(server, "session_successor", "Inspect the exact manual successor context for remaining scope of the launch-bound predecessor. No adoption or execution.", governance.SessionWireDefinition("successorRequest"), true, func(b []byte) (any, error) {
		var q governance.SessionSuccessorRequest
		if e := governance.Decode(b, &q); e != nil {
			return nil, e
		}
		if q.InstanceID != instance {
			return nil, &governance.Error{Code: "unauthorized", Message: "request outside launch-bound Session"}
		}
		if e := check(); e != nil {
			return nil, e
		}
		return service.SessionSuccessorContext(q)
	})
	addSessionTool(server, "session_rebind_context", "Inspect the exact proposed binding revision at a quiescent boundary. Earlier evidence is retained; prepare an attributed exact-context review before rebind.", governance.SessionWireDefinition("rebindContextRequest"), true, func(b []byte) (any, error) {
		var q governance.SessionRebindContextRequest
		if e := governance.Decode(b, &q); e != nil {
			return nil, e
		}
		if q.InstanceID != instance {
			return nil, &governance.Error{Code: "unauthorized", Message: "request outside launch-bound Session"}
		}
		if e := check(); e != nil {
			return nil, e
		}
		return service.SessionRebindContext(q)
	})
	caller := governance.LocalController()
	addSessionMutation(server, service, instance, skill, "rebind", "rebindRequest", func(q governance.SessionRebindRequest) (any, error) { return service.RebindSession(caller, q) })
	addSessionMutation(server, service, instance, skill, "resume", "resumeRequest", func(q governance.SessionResumeRequest) (any, error) { return service.ResumeSession(caller, q) })
	addSessionMutation(server, service, instance, skill, "select", "selectRequest", func(q governance.SessionSelectRequest) (any, error) { return service.SelectSession(caller, q) })
	addSessionMutation(server, service, instance, skill, "assess", "assessRequest", func(q governance.SessionAssessRequest) (any, error) { return service.AssessSession(caller, q) })
	for _, op := range []string{"start", "progress", "verify", "review"} {
		method := map[string]func(governance.Caller, governance.SessionWorkRequest) (governance.SessionMutationResult, error){"start": service.StartSession, "progress": service.ProgressSession, "verify": service.VerifySession, "review": service.ReviewSession}[op]
		addSessionMutation(server, service, instance, skill, op, op+"Request", func(q governance.SessionWorkRequest) (any, error) { return method(caller, q) })
	}
	addSessionMutation(server, service, instance, skill, "dispose", "disposeRequest", func(q governance.SessionDisposeRequest) (any, error) { return service.DisposeSession(caller, q) })
	addSessionMutation(server, service, instance, skill, "owner_request", "ownerRequest", func(q governance.SessionOwnerRequest) (any, error) {
		if q.Action != "request" {
			return nil, &governance.Error{Code: "unauthorized", Message: "leader may request owner review; decisions require the owner adapter"}
		}
		return service.OwnerSession(caller, q)
	})
	addSessionMutation(server, service, instance, skill, "advance", "advanceRequest", func(q governance.SessionAdvanceRequest) (any, error) { return service.AdvanceSession(caller, q) })
	return server, nil
}
func sessionRequestInstance(q any) string {
	switch q := q.(type) {
	case governance.SessionRebindRequest:
		return q.InstanceID
	case governance.SessionResumeRequest:
		return q.InstanceID
	case governance.SessionSelectRequest:
		return q.InstanceID
	case governance.SessionAssessRequest:
		return q.InstanceID
	case governance.SessionWorkRequest:
		return q.InstanceID
	case governance.SessionDisposeRequest:
		return q.InstanceID
	case governance.SessionOwnerRequest:
		return q.InstanceID
	case governance.SessionAdvanceRequest:
		return q.InstanceID
	}
	return ""
}
func addSessionMutation[T any](server *mcp.Server, service *governance.Service, instance string, skill governance.SkillInput, op, definition string, apply func(T) (any, error)) {
	addSessionTool(server, "session_"+op, "Apply the supported "+op+" operation only to the launcher-selected Session. Core scope, coverage, stage and revision checks apply. Reuse operationId after an uncertain response; inspect/resume prepared effects. No implicit host dispatch or owner approval.", governance.SessionWireDefinition(definition), false, func(b []byte) (any, error) {
		var q T
		if e := governance.Decode(b, &q); e != nil {
			return nil, e
		}
		if sessionRequestInstance(q) != instance {
			return nil, &governance.Error{Code: "unauthorized", Message: "request outside launch-bound Session"}
		}
		if e := service.CheckSessionLeader(instance, skill); e != nil {
			return nil, e
		}
		return apply(q)
	})
}
func addSessionTool(server *mcp.Server, name, description string, schema map[string]any, readOnly bool, apply func([]byte) (any, error)) {
	server.AddTool(&mcp.Tool{Name: name, Description: description, InputSchema: schema, Annotations: &mcp.ToolAnnotations{ReadOnlyHint: readOnly}}, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		if e := ctx.Err(); e != nil {
			return nil, e
		}
		value, e := apply(req.Params.Arguments)
		// Retry only busy bookkeeping, never stale requests or host effects.
		for n := 0; n < 99; n++ {
			var ge *governance.Error
			if !errors.As(e, &ge) || ge.Code != "busy" {
				break
			}
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(20 * time.Millisecond):
			}
			value, e = apply(req.Params.Arguments)
		}
		if e != nil {
			var ge *governance.Error
			if !errors.As(e, &ge) {
				ge = &governance.Error{Code: "invalid_request", Message: e.Error()}
			}
			value = map[string]any{"schema": governance.SessionSchema, "error": ge}
		}
		b, err := json.Marshal(value)
		if err != nil {
			return nil, err
		}
		if len(b) > 256<<10 {
			e = fmt.Errorf("bounded Session response exceeded")
			b, _ = json.Marshal(map[string]any{"schema": governance.SessionSchema, "error": &governance.Error{Code: "limit", Message: e.Error()}})
		}
		return &mcp.CallToolResult{IsError: e != nil, StructuredContent: json.RawMessage(b), Content: []mcp.Content{&mcp.TextContent{Text: string(b)}}}, nil
	})
}
