package governancecontroller

import (
	"context"
	"errors"
	"fmt"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
)

// ExecuteSessionHost consumes a launch-authorized request once. The caller owns
// the run's live-controller lease and observes its host until detach/completion.
// A lost response is recovered from the exact durable input; it is never resent.
func (c *Controller) ExecuteSessionHost(ctx context.Context, id string) (governance.SessionHostCommand, error) {
	caller := governance.LocalController()
	release, err := c.Service.SessionHostCommandLease(caller, id)
	if err != nil {
		return governance.SessionHostCommand{}, err
	}
	defer release()
	command, err := c.Service.ReadSessionHostCommand(id)
	if err != nil {
		return command, err
	}
	switch command.State {
	case "acknowledged":
		return command, nil
	case "canceled":
		return command, fmt.Errorf("conflict: Session host request canceled")
	case "claimed", "uncertain":
		command, err = c.Service.SettleSessionHost(caller, id)
		if err == nil && command.State != "acknowledged" {
			err = fmt.Errorf("reconciliation_required: inspect exact host input; request cannot be resent")
		}
		return command, err
	}
	command, err = c.Service.ClaimSessionHost(caller, id)
	if err != nil {
		return command, err
	}
	executor := *c
	executor.SessionHostInputID = id
	q := command.Request
	switch q.Action {
	case "launch":
		var selected governance.SessionHostSelection
		selected, err = c.Service.SessionHostSelection(q.InstanceID)
		if err == nil {
			executor.SkillPaths, err = continuationSkillPaths(c.Root, selected.Skills)
			if err == nil {
				_, err = executor.Launch(ctx, q.RunID, "Execute the committed Session assignment: "+selected.Outcome)
			}
		}
	case "continue":
		_, err = executor.ContinueAttention(ctx, q.RunID)
	case "interrupt", "scope-change":
		_, err = executor.Control(ctx, q.RunID, q.Action, q.Reason)
	default:
		err = fmt.Errorf("invalid_request: unsupported host action")
	}
	settled, settleErr := c.Service.SettleSessionHost(caller, id)
	if settleErr != nil {
		return command, errors.Join(err, settleErr)
	}
	if err == nil && settled.State != "acknowledged" {
		err = fmt.Errorf("reconciliation_required: host input acknowledgement absent")
	}
	return settled, err
}
