package governancehost

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
)

type Event struct {
	ThreadID    string                 `json:"threadId"`
	Observation governance.Observation `json:"observation"`
}

// Visible projects only completed visible messages and selected activity
// metadata. Reasoning, authentication, raw tool output and deltas are discarded.
func Visible(method string, params []byte, epoch string, received time.Time) (Event, bool, error) {
	if method != "item/completed" && method != "turn/started" && method != "turn/completed" {
		return Event{}, false, nil
	}
	if e := governance.CheckWire(params); e != nil {
		return Event{}, false, e
	}
	var p struct {
		ThreadID string          `json:"threadId"`
		TurnID   string          `json:"turnId"`
		Item     json.RawMessage `json:"item"`
		Turn     Turn            `json:"turn"`
	}
	if e := json.Unmarshal(params, &p); e != nil {
		return Event{}, false, e
	}
	o := governance.Observation{Epoch: epoch, ReceivedAt: received.Format(time.RFC3339Nano)}
	if p.ThreadID == "" {
		return Event{}, false, fmt.Errorf("missing thread identity")
	}
	if method != "item/completed" {
		o.TurnID = p.Turn.ID
		if o.TurnID == "" {
			return Event{}, false, fmt.Errorf("missing turn identity")
		}
		if method == "turn/started" {
			o.Kind = "turn.started"
			o.Status = "running"
		} else {
			o.Kind = "turn.completed"
			o.Status = p.Turn.Status
			if o.Status != "completed" && o.Status != "failed" && o.Status != "interrupted" {
				return Event{}, false, fmt.Errorf("unsupported turn status")
			}
		}
		return Event{p.ThreadID, o}, true, nil
	}
	var item struct {
		ClientID string  `json:"clientId"`
		ID       string  `json:"id"`
		Type     string  `json:"type"`
		Phase    *string `json:"phase"`
		Text     string  `json:"text"`
		Status   string  `json:"status"`
		Content  []struct {
			Type string `json:"type"`
			Text string `json:"text"`
			Path string `json:"path"`
		} `json:"content"`
	}
	if e := json.Unmarshal(p.Item, &item); e != nil {
		return Event{}, false, e
	}
	o.TurnID = p.TurnID
	o.ItemID = item.ID
	switch item.Type {
	case "agentMessage":
		o.Kind = item.Type
		o.Text = item.Text
		o.Phase = "unknown"
		if item.Phase != nil {
			o.Phase = *item.Phase
		}
		if o.Phase != "unknown" && o.Phase != "commentary" && o.Phase != "final_answer" {
			return Event{}, false, fmt.Errorf("unsupported message phase")
		}
	case "userMessage":
		o.Kind = item.Type
		o.ClientID = item.ClientID
		parts := []string{}
		for _, v := range item.Content {
			if v.Type == "text" {
				parts = append(parts, v.Text)
			} else {
				a := governance.Attachment{Kind: v.Type, Availability: "unavailable"}
				switch v.Type {
				case "localImage", "localAudio", "skill", "mention":
					a.Reference = v.Path
					if v.Path != "" {
						a.Availability = "unknown"
					}
				case "image", "audio":
				default:
					return Event{}, false, fmt.Errorf("unsupported attachment kind")
				}
				if len(a.Reference) > 4096 {
					return Event{}, false, fmt.Errorf("attachment reference limit")
				}
				o.Attachments = append(o.Attachments, a)
				if len(o.Attachments) > 64 {
					return Event{}, false, fmt.Errorf("attachment count limit")
				}
			}
		}
		o.Text = strings.Join(parts, "\n")
	case "commandExecution", "mcpToolCall", "collabAgentToolCall":
		o.Kind = "tool.completed"
		o.Status = item.Status
	default:
		return Event{}, false, nil
	}
	if o.ItemID == "" || o.TurnID == "" || len(o.Text) > 64<<10 {
		return Event{}, false, fmt.Errorf("invalid or oversized completed message")
	}
	return Event{p.ThreadID, o}, true, nil
}

// Binding is controller-owned. It holds durable SDP identities separately from
// the observed host thread/attempt. The caller retains each operation on failure
// and retries it unchanged through RecordHost, never by resubmitting a model turn.
type Binding struct {
	Result                     governance.Result
	AttemptID, ThreadID, Epoch string
}

func (b Binding) Request(action string, o governance.Observation) governance.HostRequest {
	return governance.HostRequest{Schema: governance.Schema, ProjectRef: b.Result.ProjectRef, CheckoutID: b.Result.CheckoutID, OperationID: UUID(), RunID: b.Result.Run.ID, ExpectedRevision: b.Result.Run.Revision, InputRevision: b.Result.Run.InputRevision, Action: action, AttemptID: b.AttemptID, HostVersion: HostVersion, ThreadID: b.ThreadID, Observation: o}
}

// History returns a conservative recovery result for one already mapped turn.
// Missing turn/history stays uncertain; nothing here launches or resumes work.
func History(t Thread, turnID, epoch string, received time.Time) ([]Event, string, error) {
	if turnID == "" {
		if len(t.Turns) == 0 && t.Status.Type == "idle" {
			return nil, "idle", nil
		}
		return nil, "unknown", fmt.Errorf("reconciliation_required: no mapped turn")
	}
	for _, turn := range t.Turns {
		if turn.ID != turnID {
			continue
		}
		if turn.ItemsView != "" && turn.ItemsView != "full" {
			return nil, "unknown", fmt.Errorf("observation_gap: full mapped-turn history unavailable")
		}
		if turn.Status == "inProgress" {
			return nil, "running", nil
		} // History can contain incomplete items; do not finalize them.
		events := []Event{}
		for _, item := range turn.Items {
			b, _ := json.Marshal(map[string]any{"threadId": t.ID, "turnId": turn.ID, "item": item})
			v, keep, e := Visible("item/completed", b, epoch, received)
			if e != nil {
				return nil, "unknown", e
			}
			if keep {
				events = append(events, v)
			}
		}
		status := turn.Status
		if status == "inProgress" {
			status = "running"
		}
		switch status {
		case "running", "completed", "failed", "interrupted":
			return events, status, nil
		}
		return nil, "unknown", fmt.Errorf("unsupported history status")
	}
	return nil, "unknown", fmt.Errorf("observation_gap: mapped turn absent from history")
}
