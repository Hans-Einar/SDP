package governancecontroller

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Hans-Einar/SDP/SDPTool/governance"
)

// ContinueAttention launches a fresh attempt only for current recorded replies.
// No host request is retried; durable input claims fence uncertain delivery.
func (c *Controller) ContinueAttention(ctx context.Context, runID string) (governance.Result, error) {
	return c.launch(ctx, runID, "", true)
}

// AttentionPending is a read-only watcher decision, never dispatch authority.
// The core checks all references again during prepare and claim.
func (c *Controller) AttentionPending(runID string) (bool, bool, error) {
	r, e := c.Service.Read(runID)
	if e != nil {
		return false, false, e
	}
	if r.Run.State != "running" || len(r.Run.Attempts) == 0 {
		return false, false, nil
	}
	a := r.Run.Attempts[len(r.Run.Attempts)-1]
	if a.Status != "completed" {
		return false, false, nil
	}
	if a.Gap || r.Coverage != "known" {
		return false, false, fmt.Errorf("reconciliation_required: waiting assignment has stale or missing coverage")
	}
	for _, v := range r.Run.Inputs {
		if v.State != "acknowledged" {
			return false, false, fmt.Errorf("reconciliation_required: unresolved outbound input")
		}
	}
	view, e := c.Service.AttentionInbox(runID)
	if e != nil {
		return false, false, e
	}
	waiting, ready := false, false
	for _, item := range view.Items {
		if item.Status == "open" && item.Question.Blocking {
			if item.Applicability != "current" {
				return false, false, fmt.Errorf("stale_revision: blocking question no longer applies")
			}
			waiting = true
		}
		if item.Status == "answered" && item.Question.Blocking && item.Applicability == "current" && item.Consumption == "not-tracked" {
			ready = true
		}
	}
	return waiting, ready && !waiting, nil
}

type continuationContext struct {
	References []governance.AttentionReplyRef
	History    []governance.AttentionView
	Original   governance.OwnerInput
}

func (c *Controller) continuationContext(runID string) (continuationContext, error) {
	var out continuationContext
	waiting, ready, e := c.AttentionPending(runID)
	if e != nil {
		return out, e
	}
	if waiting || !ready {
		return out, fmt.Errorf("attention_required: no fully answered current waiting assignment")
	}
	r, e := c.Service.Read(runID)
	if e != nil {
		return out, e
	}
	view, e := c.Service.AttentionInbox(runID)
	if e != nil {
		return out, e
	}
	out.History = view.Items
	for _, item := range view.Items {
		if item.Status == "answered" && item.Applicability == "current" && item.Consumption == "not-tracked" {
			out.References = append(out.References, governance.AttentionReplyRef{QuestionID: item.Question.ID, ReplyOperationID: item.Question.Reply.OperationID, Basis: item.Question.Basis})
		}
	}
	for _, v := range r.Run.Inputs {
		if v.Kind == "launch" && len(v.AttentionReplies) == 0 {
			out.Original, e = c.Service.OwnerInput(governance.LocalController(), runID, v.ID)
			if e != nil {
				return out, e
			}
			break
		}
	}
	if out.Original.ID == "" {
		return out, fmt.Errorf("reconciliation_required: original launch input missing")
	}
	payload := out.Original.Payload
	if strings.HasPrefix(payload, "$") {
		_, payload, _ = strings.Cut(payload, "\n")
	}
	var policy struct {
		Sandbox       string `json:"sandbox"`
		ApprovalMode  string `json:"approvalMode"`
		ReportingMode string `json:"reportingMode"`
	}
	if e = json.Unmarshal([]byte(payload), &policy); e != nil {
		return out, fmt.Errorf("reconciliation_required: original launch policy unavailable")
	}
	mode := c.Reporting
	if mode == "" {
		mode = "manual"
	}
	if policy.ReportingMode == "" {
		policy.ReportingMode = "manual"
	}
	if policy.Sandbox != c.Sandbox || policy.ApprovalMode != c.Approvals || policy.ReportingMode != mode {
		return out, fmt.Errorf("conflict: continuation must preserve original launch policy")
	}
	if e = governance.CheckSkillSnapshots(c.Root, out.Original.Skills); e != nil {
		return out, e
	}
	return out, nil
}
func continuationSkillPaths(root string, skills []governance.SkillInput) ([]string, error) {
	var paths []string
	for _, v := range skills {
		p, e := filepath.Rel(root, v.Path)
		if e != nil {
			return nil, e
		}
		paths = append(paths, p)
	}
	return paths, nil
}
