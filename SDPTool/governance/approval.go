package governance

import "encoding/json"

// ApprovalPrompt is a bounded, observed host request, not authority from model text.
// RequestID retains the JSON type of the host ID (a canonical JSON string or integer).
type ApprovalPrompt struct {
	RequestID string            `json:"requestId"`
	Epoch     string            `json:"epoch"`
	ThreadID  string            `json:"threadId"`
	TurnID    string            `json:"turnId"`
	Message   string            `json:"message"`
	Request   TransitionRequest `json:"request"`
}
type ApprovalReply struct {
	RequestID string `json:"requestId"`
	Epoch     string `json:"epoch"`
	ThreadID  string `json:"threadId"`
	TurnID    string `json:"turnId"`
	Decision  string `json:"decision"`
}

func (v ApprovalReply) Valid() bool {
	return opaque(v.RequestID) && uuidPattern.MatchString(v.Epoch) && opaque(v.ThreadID) && opaque(v.TurnID) && (v.Decision == "accept" || v.Decision == "decline" || v.Decision == "cancel")
}
func (p ApprovalPrompt) Valid() bool {
	return opaque(p.RequestID) && uuidPattern.MatchString(p.Epoch) && opaque(p.ThreadID) && opaque(p.TurnID) && p.Message == `Allow the sdp_governance MCP server to run tool "run_transition"?` && (p.Request.Action == "report" || p.Request.Action == "fail" || p.Request.Action == "scope-change") && uuidPattern.MatchString(p.Request.OperationID)
}
func approvalAt(a Attempt, id string) (ApprovalPrompt, bool) {
	var p ApprovalPrompt
	found := false
	if a.Status != "running" {
		return p, false
	}
	for _, o := range a.Events {
		if o.Epoch != a.Epoch || o.ItemID != id {
			continue
		}
		if o.Kind == "approval.requested" {
			if Decode([]byte(o.Text), &p) != nil {
				return p, false
			}
			found = true
		}
		if o.Kind == "approval.resolved" {
			return p, false
		}
	}
	return p, found
}
func (s *Service) PendingApprovals(runID string) ([]ApprovalPrompt, error) {
	state, _, e := s.load()
	if e != nil {
		return nil, e
	}
	r, ok := state.Runs[runID]
	if !ok {
		return nil, fail("missing_coverage", "run absent")
	}
	out := []ApprovalPrompt{}
	if len(r.Attempts) == 0 {
		return out, nil
	}
	a := r.Attempts[len(r.Attempts)-1]
	for _, o := range a.Events {
		if o.Kind != "approval.requested" || o.Epoch != a.Epoch {
			continue
		}
		p, ok := approvalAt(a, o.ItemID)
		if !ok {
			continue
		}
		answered := false
		for _, v := range r.Inputs {
			var reply ApprovalReply
			if v.Kind == "approval" && Decode([]byte(v.Payload), &reply) == nil && reply.RequestID == p.RequestID && reply.Epoch == p.Epoch {
				answered = true
			}
		}
		if !answered {
			out = append(out, p)
		}
	}
	return out, nil
}
func checkApprovalInput(r Run, v OwnerInput) error {
	if r.State != "running" {
		return fail("conflict", "routine no longer awaiting a running-work approval")
	}
	var reply ApprovalReply
	if Decode([]byte(v.Payload), &reply) != nil || !reply.Valid() || reply.Epoch != v.Epoch || reply.ThreadID != v.ThreadID || reply.TurnID != v.TurnID {
		return fail("invalid_request", "invalid approval reply")
	}
	a := r.Attempts[len(r.Attempts)-1]
	p, ok := approvalAt(a, reply.RequestID)
	if !ok || p.Epoch != reply.Epoch || p.ThreadID != reply.ThreadID || p.TurnID != reply.TurnID || p.Request.RunID != r.ID || p.Request.InputRevision != r.InputRevision {
		return fail("reconciliation_required", "approval is not pending in this attempt")
	}
	for _, old := range r.Inputs {
		var prior ApprovalReply
		if old.ID != v.ID && old.Kind == "approval" && Decode([]byte(old.Payload), &prior) == nil && prior.RequestID == reply.RequestID && prior.Epoch == reply.Epoch {
			return fail("conflict", "approval already consumed")
		}
	}
	return nil
}

// Approval bookkeeping must not invalidate the exact queued agent call. Every
// intervening revision must be an observation or this reply's bookkeeping; no
// steering, lifecycle change, reconnect or authority drift may be skipped.
func approvedRevision(state snapshot, r Run, caller Caller, q TransitionRequest) bool {
	if caller.controller || q.ExpectedRevision >= r.Revision || len(r.Attempts) == 0 {
		return false
	}
	a := r.Attempts[len(r.Attempts)-1]
	if a.ID != caller.attemptID || a.Epoch != caller.epoch || a.Status != "running" {
		return false
	}
	inputID := ""
	for _, v := range r.Inputs {
		var reply ApprovalReply
		if v.Kind != "approval" || v.AttemptID != a.ID || v.Epoch != a.Epoch || (v.State != "sending" && v.State != "uncertain" && v.State != "acknowledged") || Decode([]byte(v.Payload), &reply) != nil || reply.Decision != "accept" {
			continue
		}
		for _, o := range a.Events {
			var p ApprovalPrompt
			if o.Kind == "approval.requested" && o.Epoch == a.Epoch && o.ItemID == reply.RequestID && Decode([]byte(o.Text), &p) == nil && requestDigest(p.Request) == requestDigest(q) {
				inputID = v.ID
			}
		}
	}
	if inputID == "" {
		return false
	}
	revisions := map[uint64]Result{}
	for _, receipt := range state.Receipts {
		v := receipt.Result
		if v.Run.ID == r.ID {
			revisions[v.Run.Revision] = v
		}
	}
	base, ok := revisions[q.ExpectedRevision]
	if !ok {
		return false
	}
	for rev := q.ExpectedRevision + 1; rev <= r.Revision; rev++ {
		next, ok := revisions[rev]
		if !ok {
			return false
		}
		switch next.Operation {
		case "host.observe":
		case "input.prepare", "input.claim", "input.uncertain", "input.resolve":
			if len(next.Run.Inputs) == 0 || next.Run.Inputs[len(next.Run.Inputs)-1].ID != inputID {
				return false
			}
		default:
			return false
		}
		// Keep all decision and outbound intent fields, excluding only the one reply.
		normalize := func(v Run) string {
			v.Revision = 0
			v.Inputs = append([]OwnerInput(nil), v.Inputs...)
			filtered := v.Inputs[:0]
			for _, i := range v.Inputs {
				if i.ID != inputID {
					filtered = append(filtered, i)
				}
			}
			v.Inputs = filtered
			v.Attempts = append([]Attempt(nil), v.Attempts...)
			for i := range v.Attempts {
				v.Attempts[i].Events = nil
				v.Attempts[i].EventCount = 0
				v.Attempts[i].LastObservedAt = ""
			}
			b, _ := json.Marshal(v)
			return string(b)
		}
		if normalize(base.Run) != normalize(next.Run) {
			return false
		}
	}
	return true
}
