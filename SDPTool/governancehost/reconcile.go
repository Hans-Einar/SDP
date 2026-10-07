package governancehost

import (
	"context"
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
		events, status, e := History(t, a.TurnID, c.Epoch, now)
		if e != nil {
			return Recovery{}, e
		}
		return Recovery{Binding: Binding{Result: r, AttemptID: a.ID, ThreadID: a.ThreadID, Epoch: c.Epoch}, Events: events, Status: status, ReceivedAt: now}, nil
	}
	return Recovery{}, fmt.Errorf("missing_coverage: mapped attempt absent")
}

// ReconcileRequest records the coverage gap and observed status before any replay.
// Replay each allowlisted completed item using Binding.Request("observe", event),
// refreshing Binding.Result after each receipt. Conflicting replay refuses.
func (r Recovery) ReconcileRequest() governance.HostRequest {
	turnID := ""
	for _, a := range r.Binding.Result.Run.Attempts {
		if a.ID == r.Binding.AttemptID {
			turnID = a.TurnID
		}
	}
	return r.Binding.Request("reconcile", governance.Observation{Kind: "reconciled", TurnID: turnID, Epoch: r.Binding.Epoch, ReceivedAt: r.ReceivedAt.Format(time.RFC3339Nano), Status: r.Status})
}
