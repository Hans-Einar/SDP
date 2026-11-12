package governancemcp

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// responsePolicy preserves the two existing wire contracts. Capability and
// launch-bound authorization belong to each registered operation, not this codec.
type responsePolicy struct {
	schema, limitMessage string
}

func addTool(server *mcp.Server, name, description string, schema map[string]any, readOnly bool, policy responsePolicy, apply func([]byte) (any, error)) {
	server.AddTool(&mcp.Tool{Name: name, Description: description, InputSchema: schema, Annotations: &mcp.ToolAnnotations{ReadOnlyHint: readOnly}}, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return toolResponse(ctx, req.Params.Arguments, policy, apply)
	})
}

func toolResponse(ctx context.Context, arguments []byte, policy responsePolicy, apply func([]byte) (any, error)) (*mcp.CallToolResult, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	value, err := apply(arguments)
	// Core busy means no mutation occurred. Never retry stale requests or host
	// effects; retain the exact arguments and operation identity for each retry.
	for attempt := 0; attempt < 99; attempt++ {
		var busy *governance.Error
		if !errors.As(err, &busy) || busy.Code != "busy" {
			break
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(20 * time.Millisecond):
		}
		value, err = apply(arguments)
	}
	if err != nil {
		var domain *governance.Error
		if !errors.As(err, &domain) {
			domain = &governance.Error{Code: "invalid_request", Message: err.Error()}
		}
		value = map[string]any{"schema": policy.schema, "error": domain}
	}
	b, marshalErr := json.Marshal(value)
	if marshalErr != nil {
		return nil, marshalErr
	}
	if len(b) > 256<<10 {
		err = &governance.Error{Code: "limit", Message: policy.limitMessage}
		b, _ = json.Marshal(map[string]any{"schema": policy.schema, "error": err})
	}
	// No post-call cancellation check: a committed outcome remains committed.
	return &mcp.CallToolResult{IsError: err != nil, StructuredContent: json.RawMessage(b), Content: []mcp.Content{&mcp.TextContent{Text: string(b)}}}, nil
}
