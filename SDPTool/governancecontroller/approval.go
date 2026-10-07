package governancecontroller

import (
	"encoding/json"
	"fmt"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
	"github.com/Hans-Einar/SDP/SDPTool/governancehost"
)

func (c *Controller) Approve(runID string, reply governance.ApprovalReply, exactInput string) (result governance.Result, err error) {
	host, ok := c.Host.(interface {
		ReplyApproval(governance.ApprovalReply) error
	})
	if !ok || c.Approvals != "ask" || reply.Epoch != c.Epoch {
		return result, fmt.Errorf("explicit approval mode and current epoch required")
	}
	unlock, e := c.Service.ControllerLease(governance.LocalController())
	if e != nil {
		return result, e
	}
	defer unlock()
	r, e := c.Service.Read(runID)
	if e != nil {
		return r, e
	}
	if len(r.Run.Attempts) == 0 {
		return r, fmt.Errorf("no mapped attempt")
	}
	a := r.Run.Attempts[len(r.Run.Attempts)-1]
	b, _ := json.Marshal(reply)
	v := governance.OwnerInput{ID: governancehost.UUID(), Kind: "approval", Text: exactInput, Payload: string(b), AttemptID: a.ID, Epoch: reply.Epoch, ThreadID: reply.ThreadID, TurnID: reply.TurnID, ReceivedAt: now()}
	if r, e = c.input(runID, "prepare", v); e != nil {
		return r, e
	}
	v.State = "prepared"
	if r, e = c.input(runID, "claim", v); e != nil {
		return r, e
	}
	v.State = "sending"
	sendErr := host.ReplyApproval(reply)
	// Even a successful write has no host acknowledgment yet. Never retransmit.
	r, e = c.input(runID, "uncertain", v)
	if sendErr != nil {
		return r, sendErr
	}
	return r, e
}
func (c *Controller) ResolveApproval(runID string, event governancehost.Event) (governance.Result, error) {
	r, e := c.Observe(runID, event)
	if e != nil {
		return r, e
	}
	for _, summary := range r.Run.Inputs {
		if summary.Kind != "approval" || summary.Epoch != event.Observation.Epoch || (summary.State != "sending" && summary.State != "uncertain") {
			continue
		}
		input, e := c.Service.OwnerInput(governance.LocalController(), runID, summary.ID)
		if e != nil {
			return r, e
		}
		var reply governance.ApprovalReply
		if json.Unmarshal([]byte(input.Payload), &reply) != nil || reply.RequestID != event.Observation.ItemID {
			continue
		}
		r, e = c.input(runID, "resolve", input)
		if e != nil {
			return r, e
		}
	}
	return r, nil
}
