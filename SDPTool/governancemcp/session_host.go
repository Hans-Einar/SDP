package governancemcp

import (
	"github.com/Hans-Einar/SDP/SDPTool/governance"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// NewSessionHostController adds only explicit host-request authority. It cannot
// execute a host RPC, grant owner acceptance or address another Session.
func NewSessionHostController(service *governance.Service, instance string, skill governance.SkillInput) (*mcp.Server, error) {
	server, err := NewSessionController(service, instance, skill)
	if err != nil {
		return nil, err
	}
	check := func() error { return service.CheckSessionLeader(instance, skill) }
	empty := map[string]any{"type": "object", "properties": map[string]any{}, "additionalProperties": false}
	addSessionTool(server, "session_host_read", "Read durable host requests for the launch-bound Session. Acknowledgement proves delivery only, not completion or acceptance.", empty, true, func(b []byte) (any, error) {
		var q struct{}
		if e := governance.Decode(b, &q); e != nil {
			return nil, e
		}
		if e := check(); e != nil {
			return nil, e
		}
		return service.ReadSessionHostCommands(instance)
	})
	addSessionTool(server, "session_host_request", "Queue exactly one revision-bound launch, attention continuation, interrupt or scope-change. Reuse operationId on retry. The trusted live controller executes; uncertain delivery requires recovery, never resend.", governance.SessionWireDefinition("hostRequest"), false, func(b []byte) (any, error) {
		var q governance.SessionHostRequest
		if e := governance.Decode(b, &q); e != nil {
			return nil, e
		}
		if q.InstanceID != instance || q.LeaderSkill != skill {
			return nil, &governance.Error{Code: "unauthorized", Message: "request outside launcher Session/skill grant"}
		}
		if e := check(); e != nil {
			return nil, e
		}
		return service.RequestSessionHost(governance.LocalController(), q)
	})
	addSessionTool(server, "session_host_cancel", "Cancel a request only when no durable host input exists and no controller holds its lease. Otherwise reconcile exact input delivery first.", governance.SessionWireDefinition("hostCancelRequest"), false, func(b []byte) (any, error) {
		var q struct {
			OperationID string `json:"operationId"`
		}
		if e := governance.Decode(b, &q); e != nil {
			return nil, e
		}
		if e := check(); e != nil {
			return nil, e
		}
		command, e := service.ReadSessionHostCommand(q.OperationID)
		if e != nil {
			return nil, e
		}
		if command.Request.InstanceID != instance {
			return nil, &governance.Error{Code: "unauthorized", Message: "request outside launcher Session grant"}
		}
		return service.CancelSessionHost(governance.LocalController(), q.OperationID)
	})
	return server, nil
}
