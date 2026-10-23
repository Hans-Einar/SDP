package governancehost

import (
	"context"
	"fmt"
	"time"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
)

// PrepareAttachment subscribes to one already executing turn on the same measured
// server. It never starts a turn, supplies configuration overrides or retries input.
func (c *Client) PrepareAttachment(ctx context.Context, service *governance.Service, runID, cwd, approvals string) (Recovery, error) {
	r, err := service.Read(runID)
	if err != nil {
		return Recovery{}, err
	}
	if r.Coverage != "known" || r.Run.State != "running" || len(r.Run.Attempts) == 0 {
		return Recovery{}, fmt.Errorf("current running assignment required")
	}
	a := r.Run.Attempts[len(r.Run.Attempts)-1]
	if c.HostIdentity == "" || c.HostIdentity != a.HostIdentity || a.HostVersion != c.Version() || a.TurnID == "" || (a.Status != "running" && a.Status != "unknown") {
		return Recovery{}, fmt.Errorf("same active socket host and mapped turn required")
	}
	for _, input := range r.Run.Inputs {
		if input.State != "acknowledged" {
			return Recovery{}, fmt.Errorf("reconciliation_required: unresolved outbound input")
		}
	}
	pending := map[string]bool{}
	cursor := ""
	for {
		page, err := service.Events(runID, a.ID, cursor, 256)
		if err != nil {
			return Recovery{}, err
		}
		if page.Revision != r.Run.Revision {
			return Recovery{}, fmt.Errorf("stale_revision: observations changed during attachment preflight")
		}
		for _, event := range page.Events {
			if event.Epoch != a.Epoch {
				continue
			}
			if event.Kind == "approval.requested" {
				pending[event.ItemID] = true
			}
			if event.Kind == "approval.resolved" {
				delete(pending, event.ItemID)
			}
		}
		cursor = page.NextCursor
		if cursor == "" {
			break
		}
	}
	if len(pending) > 0 {
		return Recovery{}, fmt.Errorf("reconciliation_required: unresolved approval closure")
	}
	if approvals != "ask" && approvals != "never" {
		return Recovery{}, fmt.Errorf("unsupported approval mode")
	}
	c.approvalMu.Lock()
	c.ask = approvals == "ask"
	c.approvalMu.Unlock()
	valid := func(t Thread) bool {
		if t.ID != a.ThreadID || t.Cwd != cwd || t.CLIVersion != c.Version() || t.Status.Type != "active" {
			return false
		}
		found := false
		for _, turn := range t.Turns {
			if turn.Status == "inProgress" {
				if turn.ID != a.TurnID || found {
					return false
				}
				found = true
			}
		}
		return found
	}
	t, err := c.ReadThread(ctx, a.ThreadID)
	if err != nil {
		return Recovery{}, err
	}
	if !valid(t) {
		return Recovery{}, fmt.Errorf("reconciliation_required: mapped turn is not executing on this host")
	}
	var out struct {
		Thread Thread `json:"thread"`
	}
	if err = c.call(ctx, "thread/resume", map[string]string{"threadId": a.ThreadID}, &out); err != nil {
		return Recovery{}, err
	}
	if !valid(out.Thread) {
		return Recovery{}, fmt.Errorf("reconciliation_required: turn changed while subscribing")
	}
	now := time.Now().UTC()
	events, status, err := History(out.Thread, a.TurnID, c.Epoch, now)
	if err != nil {
		return Recovery{}, err
	}
	return Recovery{Binding: Binding{Result: r, AttemptID: a.ID, ThreadID: a.ThreadID, Epoch: c.Epoch}, Events: events, TurnID: a.TurnID, Status: status, ReceivedAt: now}, nil
}
func (r Recovery) ReattachRequest(hostIdentity string) governance.HostRequest {
	q := r.ReconcileRequest()
	q.Action = "reattach"
	q.HostIdentity = hostIdentity
	return q
}
