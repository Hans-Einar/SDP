package governancehost

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
)

type pendingApproval struct {
	prompt governance.ApprovalPrompt
	id     json.RawMessage
	sent   bool
}

func requestKey(id json.RawMessage) (string, error) {
	var value any
	decoder := json.NewDecoder(bytes.NewReader(id))
	decoder.UseNumber()
	if e := decoder.Decode(&value); e != nil {
		return "", e
	}
	switch v := value.(type) {
	case string:
		if v == "" {
			return "", fmt.Errorf("empty host request id")
		}
	case json.Number:
		if _, e := v.Int64(); e != nil {
			return "", e
		}
	default:
		return "", fmt.Errorf("unsupported host request id")
	}
	b, e := json.Marshal(value)
	if len(b) > 256 {
		return "", fmt.Errorf("host request id exceeds limit")
	}
	return string(b), e
}
func parseApproval(f frame, epoch string) (governance.ApprovalPrompt, error) {
	var p struct {
		ThreadID   string          `json:"threadId"`
		TurnID     string          `json:"turnId"`
		ServerName string          `json:"serverName"`
		Mode       string          `json:"mode"`
		Message    string          `json:"message"`
		Schema     json.RawMessage `json:"requestedSchema"`
		Meta       struct {
			Kind    string          `json:"codex_approval_kind"`
			Request json.RawMessage `json:"tool_params"`
		} `json:"_meta"`
	}
	if e := json.Unmarshal(f.Params, &p); e != nil {
		return governance.ApprovalPrompt{}, e
	}
	var form struct {
		Type       string                     `json:"type"`
		Properties map[string]json.RawMessage `json:"properties"`
	}
	if governance.Decode(p.Schema, &form) != nil || form.Type != "object" || form.Properties == nil || len(form.Properties) != 0 || p.ServerName != "sdp_governance" || p.Mode != "form" || p.Meta.Kind != "mcp_tool_call" {
		return governance.ApprovalPrompt{}, fmt.Errorf("unsupported approval form")
	}
	var q governance.TransitionRequest
	if e := governance.Decode(p.Meta.Request, &q); e != nil {
		return governance.ApprovalPrompt{}, e
	}
	key, e := requestKey(f.ID)
	if e != nil {
		return governance.ApprovalPrompt{}, e
	}
	result := governance.ApprovalPrompt{RequestID: key, Epoch: epoch, ThreadID: p.ThreadID, TurnID: p.TurnID, Message: p.Message, Request: q}
	if !result.Valid() {
		return result, fmt.Errorf("unsupported approval identity or operation")
	}
	return result, nil
}
func (c *Client) approvalRequest(f frame) (Event, bool, error) {
	c.approvalMu.Lock()
	defer c.approvalMu.Unlock()
	if !c.ask || f.Method != "mcpServer/elicitation/request" {
		return Event{}, false, nil
	}
	p, e := parseApproval(f, c.Epoch)
	if e != nil {
		return Event{}, false, nil
	} // Unsupported requests use an RPC error, never a grant.
	if _, ok := c.approvals[p.RequestID]; ok {
		return Event{}, false, fmt.Errorf("duplicate live approval request id")
	}
	if len(c.approvals) >= 16 {
		return Event{}, false, fmt.Errorf("approval queue limit")
	}
	c.approvals[p.RequestID] = pendingApproval{prompt: p, id: append(json.RawMessage(nil), f.ID...)}
	b, _ := json.Marshal(p)
	return Event{p.ThreadID, governance.Observation{Kind: "approval.requested", ItemID: p.RequestID, TurnID: p.TurnID, Epoch: c.Epoch, Text: string(b), ReceivedAt: time.Now().UTC().Format(time.RFC3339Nano)}}, true, nil
}
func (c *Client) resolveApproval(params []byte) (Event, bool, error) {
	var p struct {
		ThreadID  string          `json:"threadId"`
		RequestID json.RawMessage `json:"requestId"`
	}
	if e := json.Unmarshal(params, &p); e != nil {
		return Event{}, false, e
	}
	key, e := requestKey(p.RequestID)
	if e != nil {
		return Event{}, false, e
	}
	c.approvalMu.Lock()
	defer c.approvalMu.Unlock()
	pending, ok := c.approvals[key]
	if !ok {
		return Event{}, false, nil
	}
	if p.ThreadID != pending.prompt.ThreadID {
		return Event{}, false, fmt.Errorf("resolved approval thread mismatch")
	}
	delete(c.approvals, key)
	return Event{p.ThreadID, governance.Observation{Kind: "approval.resolved", ItemID: key, TurnID: pending.prompt.TurnID, Epoch: c.Epoch, ReceivedAt: time.Now().UTC().Format(time.RFC3339Nano)}}, true, nil
}

// ReplyApproval sends at most once on this transport. Only a later resolved
// notification establishes closure; writing does not prove the tool ran.
func (c *Client) ReplyApproval(v governance.ApprovalReply) error {
	c.approvalMu.Lock()
	defer c.approvalMu.Unlock()
	p, ok := c.approvals[v.RequestID]
	if !v.Valid() || !ok || p.sent || v.Epoch != c.Epoch || v.ThreadID != p.prompt.ThreadID || v.TurnID != p.prompt.TurnID {
		return fmt.Errorf("approval is absent, consumed or belongs to another attempt")
	}
	select {
	case <-c.done:
		return fmt.Errorf("approval transport closed")
	default:
	}
	p.sent = true
	c.approvals[v.RequestID] = p // Consume before the first outbound byte.
	var content any
	if v.Decision == "accept" {
		content = map[string]any{}
	}
	return c.send(map[string]any{"id": p.id, "result": map[string]any{"action": v.Decision, "content": content}})
}
