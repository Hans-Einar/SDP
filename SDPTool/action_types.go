package sdptool

import (
	"encoding/json"
	"github.com/Hans-Einar/SDP/SDPTool/blueprintstate"
)

const InvokeSchema = "sdptool.invoke/1"
const ActionsSchema = "sdptool.actions/1"
const MaxInvokeBytes = 1 << 20
const MaxInvokeDepth = 32

// InvokeRequest is one local subprocess request; it never conveys caller authority.
type InvokeRequest struct {
	Schema     string          `json:"schema"`
	RequestID  string          `json:"requestId,omitempty"`
	Action     string          `json:"action"`
	Context    InvokeContext   `json:"context"`
	Parameters json.RawMessage `json:"parameters"`
}
type InvokeContext struct {
	Project string `json:"project,omitempty"`
}
type InvokeError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
type InvokeResponse struct {
	Schema    string       `json:"schema"`
	RequestID *string      `json:"requestId"`
	Action    *string      `json:"action"`
	Status    string       `json:"status"`
	ExitCode  int          `json:"exitCode"`
	Result    any          `json:"result,omitempty"`
	Error     *InvokeError `json:"error,omitempty"`
}

// InputSchema is the supported JSON Schema subset. PathBase is a UI annotation.
type InputSchema struct {
	Type                 string                  `json:"type,omitempty"`
	Properties           map[string]*InputSchema `json:"properties,omitempty"`
	Required             []string                `json:"required,omitempty"`
	AdditionalProperties *bool                   `json:"additionalProperties,omitempty"`
	Enum                 []string                `json:"enum,omitempty"`
	MinLength            *int                    `json:"minLength,omitempty"`
	MaxLength            *int                    `json:"maxLength,omitempty"`
	OneOf                []*InputSchema          `json:"oneOf,omitempty"`
	PathBase             string                  `json:"x-pathBase,omitempty"`
	Format               string                  `json:"format,omitempty"`
}
type ActionAvailability struct {
	State   string `json:"state"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
}
type ActionDescriptor struct {
	ID           string             `json:"id"`
	Label        string             `json:"label"`
	Description  string             `json:"description"`
	Group        string             `json:"group"`
	IconKey      string             `json:"iconKey"`
	Mutation     string             `json:"mutation"`
	InputSchema  *InputSchema       `json:"inputSchema"`
	OutputSchema string             `json:"outputSchema"`
	Availability ActionAvailability `json:"availability"`
}
type ActionCatalogue struct {
	Diagnostics []InvokeError `json:"diagnostics,omitempty"`
	Schema      string        `json:"schema"`
	Operation   string        `json:"operation"`
	Producer    struct {
		Version  string `json:"version"`
		Revision string `json:"revision"`
	} `json:"producer"`
	InvocationVersions []string           `json:"invocationVersions"`
	Revision           string             `json:"revision"`
	Actions            []ActionDescriptor `json:"actions"`
}

// ActionCaller is injected by a trusted host or attributed by the local CLI.
// It is deliberately absent from InvokeRequest.
type ActionCaller struct{ Principal blueprintstate.Principal }
