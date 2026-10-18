package governancehost

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
)

// Recovery is prepared from an actual mapped-thread read. Applying it never
// submits a turn or infers routine completion. Callers persist each request intent
// before calling RecordHost; on ambiguous writes they retry that exact request.
type Recovery struct {
	Binding    Binding
	Events     []Event
	Inputs     []governance.OwnerInput
	TurnID     string
	Status     string
	ReceivedAt time.Time
}

func (c *Client) PrepareRecovery(ctx context.Context, service *governance.Service, runID, attemptID string) (Recovery, error) {
	r, e := service.Read(runID)
	if e != nil {
		return Recovery{}, e
	}
	for _, a := range r.Run.Attempts {
		if a.ID != attemptID {
			continue
		}
		t, e := c.ReadThread(ctx, a.ThreadID)
		if e != nil {
			return Recovery{}, e
		}
		now := time.Now().UTC()
		pending := []governance.OwnerInput{}
		for _, summary := range r.Run.Inputs {
			if summary.AttemptID != a.ID || (summary.State != "sending" && summary.State != "uncertain") || (summary.Kind != "launch" && summary.Kind != "steer") {
				continue
			}
			input, err := service.OwnerInput(governance.LocalController(), runID, summary.ID)
			if err != nil {
				return Recovery{}, err
			}
			pending = append(pending, input)
		}
		turnID, inputs, e := correlatedInputs(t, a.TurnID, pending)
		if e != nil {
			return Recovery{}, e
		}
		events, status, e := History(t, turnID, c.Epoch, now)
		if e != nil {
			return Recovery{}, e
		}
		return Recovery{Binding: Binding{Result: r, AttemptID: a.ID, ThreadID: a.ThreadID, Epoch: c.Epoch}, Events: events, Inputs: inputs, TurnID: turnID, Status: status, ReceivedAt: now}, nil
	}
	return Recovery{}, fmt.Errorf("missing_coverage: mapped attempt absent")
}

// ReconcileRequest records the coverage gap and observed status before any replay.
// Replay each allowlisted completed item using Binding.Request("observe", event),
// refreshing Binding.Result after each receipt. Conflicting replay refuses.
func (r Recovery) ReconcileRequest() governance.HostRequest {
	turnID := r.TurnID
	for _, a := range r.Binding.Result.Run.Attempts {
		if a.ID == r.Binding.AttemptID && turnID == "" {
			turnID = a.TurnID
		}
	}
	return r.Binding.Request("reconcile", governance.Observation{Kind: "reconciled", TurnID: turnID, Epoch: r.Binding.Epoch, ReceivedAt: r.ReceivedAt.Format(time.RFC3339Nano), Status: r.Status})
}

// Exact unique correlation is required; payload similarity or an idle thread is
// not proof that a particular send was accepted.
func correlatedInputs(t Thread, turnID string, pending []governance.OwnerInput) (string, []governance.OwnerInput, error) {
	inputs := []governance.OwnerInput{}
	if len(pending) > 0 {
		for _, turn := range t.Turns {
			if turn.ItemsView != "" && turn.ItemsView != "full" {
				return "", nil, fmt.Errorf("observation_gap: full correlation history unavailable")
			}
		}
	}
	for _, input := range pending {
		matches := 0
		for _, turn := range t.Turns {
			for _, raw := range turn.Items {
				var item struct {
					Type     string          `json:"type"`
					ClientID string          `json:"clientId"`
					Content  json.RawMessage `json:"content"`
				}
				if err := json.Unmarshal(raw, &item); err != nil {
					return "", nil, err
				}
				if item.Type != "userMessage" || item.ClientID != input.ID {
					continue
				}
				if governance.ObservedInputDigest(item.Content) != input.WireDigest() || (turn.ItemsView != "" && turn.ItemsView != "full") || (turnID != "" && turnID != turn.ID) {
					return "", nil, fmt.Errorf("reconciliation_required: conflicting correlated owner input")
				}
				matches++
				input.ThreadID = t.ID
				input.TurnID = turn.ID
			}
		}
		if matches != 1 {
			return "", nil, fmt.Errorf("reconciliation_required: owner input lacks unique recorded acceptance")
		}
		turnID = input.TurnID
		inputs = append(inputs, input)
	}
	return turnID, inputs, nil
}
