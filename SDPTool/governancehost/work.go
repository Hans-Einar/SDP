package governancehost

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Hans-Einar/SDP/SDPTool/governance"
)

// StartWorkThread creates a fresh persistent thread with explicitly selected
// permissions. The controller does not infer permission to widen its sandbox.
func (c *Client) StartWorkThread(ctx context.Context, cwd, sandbox, instructions, approvals string, config map[string]any) (Thread, error) {
	if sandbox != "read-only" && sandbox != "workspace-write" {
		return Thread{}, fmt.Errorf("unsupported pilot sandbox")
	}
	var policy any = "never"
	if approvals != "" && approvals != "never" && approvals != "ask" {
		return Thread{}, fmt.Errorf("unsupported approval mode")
	}
	c.approvalMu.Lock()
	c.ask = approvals == "ask"
	c.approvalMu.Unlock()
	if approvals == "ask" {
		policy = "on-request" // Stable API; the adapter still refuses every unsupported request.
	}
	var out struct {
		Thread             Thread    `json:"thread"`
		InstructionSources *[]string `json:"instructionSources"`
	}
	e := c.call(ctx, "thread/start", map[string]any{"cwd": cwd, "ephemeral": false, "sandbox": sandbox, "approvalPolicy": policy, "approvalsReviewer": "user", "developerInstructions": instructions, "config": map[string]any{"mcp_servers": config}}, &out)
	if e == nil && (out.Thread.ID == "" || out.Thread.Status.Type != "idle") {
		e = fmt.Errorf("thread is not observed idle")
	}
	if e == nil {
		source := ""
		if len(out.Thread.Source) > 0 && string(out.Thread.Source) != "null" {
			var value any
			if err := json.Unmarshal(out.Thread.Source, &value); err != nil {
				return out.Thread, err
			}
			b, _ := json.Marshal(value)
			source = string(b)
		}
		out.Thread.Context = governance.HostContext{Status: "observed", SessionID: out.Thread.SessionID, ParentThreadID: out.Thread.ParentThreadID, NativeRole: out.Thread.AgentRole, Source: source, InstructionStatus: "unavailable", InstructionSources: []string{}}
		if out.InstructionSources != nil {
			out.Thread.Context.InstructionStatus = "observed"
			out.Thread.Context.InstructionSources = *out.InstructionSources
		}
	}
	return out.Thread, e
}
func (c *Client) StartTurn(ctx context.Context, threadID, messageID, text string) (Turn, error) {
	return c.StartSkillTurn(ctx, threadID, messageID, text, nil)
}

func (c *Client) StartSkillTurn(ctx context.Context, threadID, messageID, text string, skills []governance.SkillInput) (Turn, error) {
	if !governance.ValidSkillInputs(skills) {
		return Turn{}, fmt.Errorf("invalid skill input")
	}
	var out struct {
		Turn Turn `json:"turn"`
	}
	e := c.call(ctx, "turn/start", map[string]any{"threadId": threadID, "clientUserMessageId": messageID, "input": governance.TurnInputs(text, skills)}, &out)
	if e == nil && (out.Turn.ID == "" || out.Turn.Status != "inProgress") {
		e = fmt.Errorf("turn start not observed active")
	}
	return out.Turn, e
}
func (c *Client) SteerTurn(ctx context.Context, threadID, turnID, messageID, text string) error {
	var out struct {
		TurnID string `json:"turnId"`
	}
	if e := c.call(ctx, "turn/steer", map[string]any{"threadId": threadID, "expectedTurnId": turnID, "clientUserMessageId": messageID, "input": []any{map[string]any{"type": "text", "text": text}}}, &out); e != nil {
		return e
	}
	if out.TurnID != turnID {
		return fmt.Errorf("steer response changed turn identity")
	}
	return nil
}
func (c *Client) InterruptTurn(ctx context.Context, threadID, turnID string) error {
	var out struct{}
	return c.call(ctx, "turn/interrupt", map[string]any{"threadId": threadID, "turnId": turnID}, &out)
}
