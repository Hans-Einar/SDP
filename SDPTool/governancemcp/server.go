// Package governancemcp exposes the shared process core over local MCP.
// Launch configuration, never tool arguments, owns project and assignment scope.
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

const SDKVersion = "v1.8.0"
const ProtocolVersion = "2025-11-25"

func New(service *governance.Service, runID, attemptID, epoch string) (*mcp.Server, error) {
	if service == nil {
		return nil, fmt.Errorf("missing governance service")
	}
	var caller governance.Caller
	var selection governance.Selection
	if runID != "" {
		r, e := service.Read(runID)
		if e != nil {
			return nil, e
		}
		selection = r.Run.Selection
		if attemptID != "" || epoch != "" {
			caller, e = governance.AssignedAgent(runID, attemptID, epoch)
			if e != nil {
				return nil, e
			}
		}
	} else if attemptID != "" || epoch != "" {
		return nil, fmt.Errorf("attempt binding requires a run")
	}
	s := mcp.NewServer(&mcp.Implementation{Name: "sdptool-governance", Version: "PGI2"},
		&mcp.ServerOptions{SupportedProtocolVersions: []string{ProtocolVersion}})
	add := func(name, description string, schema map[string]any, readOnly bool, f func([]byte) (any, error)) {
		s.AddTool(&mcp.Tool{Name: name, Description: description, InputSchema: schema,
			Annotations: &mcp.ToolAnnotations{ReadOnlyHint: readOnly}},
			func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				if e := ctx.Err(); e != nil {
					return nil, e
				}
				value, e := f(req.Params.Arguments)
				// Core busy means no mutation occurred. Retain exact tool arguments
				// and operation ID, and wait only for short local bookkeeping contention.
				for attempt := 0; attempt < 99; attempt++ {
					var busy *governance.Error
					if !errors.As(e, &busy) || busy.Code != "busy" {
						break
					}
					select {
					case <-ctx.Done():
						return nil, ctx.Err()
					case <-time.After(20 * time.Millisecond):
					}
					value, e = f(req.Params.Arguments)
				}
				if e != nil {
					var d *governance.Error
					if !errors.As(e, &d) {
						d = &governance.Error{Code: "invalid_request", Message: e.Error()}
					}
					value = map[string]any{"schema": governance.Schema, "error": d}
				}
				b, err := json.Marshal(value)
				if err != nil {
					return nil, err
				}
				if len(b) > 256<<10 {
					e = &governance.Error{Code: "limit", Message: "use bounded observation pages"}
					b, _ = json.Marshal(map[string]any{"schema": governance.Schema, "error": e})
				}
				return &mcp.CallToolResult{IsError: e != nil, StructuredContent: json.RawMessage(b), Content: []mcp.Content{&mcp.TextContent{Text: string(b)}}}, nil
			})
	}
	denied := func() error {
		return &governance.Error{Code: "unauthorized", Message: "request is outside launch-bound assignment"}
	}
	add("context_read", "Read project context; never creates work.", governance.WireDefinition("Selection"), true, func(b []byte) (any, error) {
		var q governance.Selection
		if e := governance.Decode(b, &q); e != nil {
			return nil, e
		}
		if runID != "" && q != selection {
			return nil, denied()
		}
		return service.Context(q)
	})
	add("route_evaluate", "Evaluate a proposed route without activating work.", governance.WireDefinition("RouteRequest"), true, func(b []byte) (any, error) {
		var q governance.RouteRequest
		if e := governance.Decode(b, &q); e != nil {
			return nil, e
		}
		if runID != "" && q.Selection != selection {
			return nil, denied()
		}
		return service.Route(q)
	})
	add("run_read", "Read only the launch-bound run.", map[string]any{"type": "object", "properties": map[string]any{}, "additionalProperties": false}, true, func(b []byte) (any, error) {
		var q struct{}
		if e := governance.Decode(b, &q); e != nil {
			return nil, e
		}
		if runID == "" {
			return nil, denied()
		}
		return service.Read(runID)
	})
	if attemptID != "" {
		add("events_read", "Read a bounded page of observations for the launch-bound attempt; cursor pins the run revision.", map[string]any{
			"type": "object", "additionalProperties": false, "properties": map[string]any{"cursor": map[string]any{"type": "string"}, "limit": map[string]any{"type": "integer", "minimum": 1, "maximum": 256}}, "required": []string{"cursor", "limit"}}, true, func(b []byte) (any, error) {
			var q struct {
				Cursor string `json:"cursor"`
				Limit  int    `json:"limit"`
			}
			if e := governance.Decode(b, &q); e != nil {
				return nil, e
			}
			return service.Events(runID, attemptID, q.Cursor, q.Limit)
		})
		schema := governance.WireDefinition("TransitionRequest")
		schema["properties"].(map[string]any)["action"] = map[string]any{"enum": []string{"report", "fail", "scope-change"}}
		add("run_transition", "Report, fail or request replanning for the assigned run; cannot start, resume, accept or complete work.", schema, false, func(b []byte) (any, error) {
			var q governance.TransitionRequest
			if e := governance.Decode(b, &q); e != nil {
				return nil, e
			}
			return service.Transition(caller, q)
		})
	}
	return s, nil
}
