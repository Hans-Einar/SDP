package governance

import (
	"strings"
	"time"
	"unicode/utf8"
)

// OwnerInput is controller-observed input and exact outbound text. A sending or
// uncertain record must never be retransmitted automatically after a crash.
type OwnerInput struct {
	ContentOmitted bool   `json:"contentOmitted,omitempty"`
	ID             string `json:"id"`
	Kind           string `json:"kind"`
	Text           string `json:"text"`
	Payload        string `json:"payload"`
	AttemptID      string `json:"attemptId"`
	Epoch          string `json:"epoch"`
	ThreadID       string `json:"threadId"`
	TurnID         string `json:"turnId"`
	State          string `json:"state"`
	ReceivedAt     string `json:"receivedAt"`
}
type InputRequest struct {
	Schema           string     `json:"schema"`
	ProjectRef       string     `json:"projectRef"`
	CheckoutID       string     `json:"checkoutId"`
	OperationID      string     `json:"operationId"`
	RunID            string     `json:"runId"`
	ExpectedRevision uint64     `json:"expectedRevision"`
	InputRevision    string     `json:"inputRevision"`
	Action           string     `json:"action"`
	Input            OwnerInput `json:"input"`
}

// ControllerLease prevents two controller loops from consuming the same launch.
// The separate checkout mutation lock remains short-lived and shared with MCP.
func (s *Service) ControllerLease(c Caller) (func(), error) {
	if !c.controller {
		return nil, fail("unauthorized", "controller required")
	}
	return lock(s.state + "/controller")
}

func (s *Service) OwnerInput(c Caller, runID, id string) (OwnerInput, error) {
	if !c.controller {
		return OwnerInput{}, fail("unauthorized", "controller required")
	}
	state, _, e := s.load()
	if e != nil {
		return OwnerInput{}, e
	}
	for _, v := range state.Runs[runID].Inputs {
		if v.ID == id {
			return v, nil
		}
	}
	return OwnerInput{}, fail("missing_coverage", "owner input absent")
}

func (s *Service) RecordInput(c Caller, q InputRequest) (Result, error) {
	if !c.controller {
		return Result{}, fail("unauthorized", "owner input requires controller provenance")
	}
	v := q.Input
	if !uuidPattern.MatchString(q.RunID) || !uuidPattern.MatchString(v.ID) || !uuidPattern.MatchString(v.AttemptID) || !uuidPattern.MatchString(v.Epoch) || q.ExpectedRevision == 0 || !digestPattern.MatchString(q.InputRevision) {
		return Result{}, fail("invalid_request", "invalid input identity")
	}
	if v.ContentOmitted || !utf8.ValidString(v.Text) || !utf8.ValidString(v.Payload) || len(v.Text) > 64<<10 || len(v.Payload) > 64<<10 || strings.TrimSpace(v.Text) == "" || strings.TrimSpace(v.Payload) == "" {
		return Result{}, fail("invalid_request", "bounded exact input and payload required")
	}
	if _, e := time.Parse(time.RFC3339Nano, v.ReceivedAt); e != nil {
		return Result{}, fail("invalid_request", "invalid receive timestamp")
	}
	switch v.Kind {
	case "launch", "steer", "interrupt", "scope-change":
	default:
		return Result{}, fail("invalid_request", "unsupported owner input kind")
	}
	if v.ThreadID != "" && !opaque(v.ThreadID) || v.TurnID != "" && !opaque(v.TurnID) {
		return Result{}, fail("invalid_request", "invalid host identity")
	}
	unlock, e := lock(s.state)
	if e != nil {
		return Result{}, e
	}
	defer unlock()
	state, head, e := s.load()
	if e != nil {
		return Result{}, e
	}
	if e = s.authorize(c, q.Schema, q.ProjectRef, q.CheckoutID, q.OperationID, state); e != nil {
		return Result{}, e
	}
	digest := requestDigest(q)
	if r, ok, e := retry(state, q.OperationID, digest); ok {
		return r, e
	}
	r, ok := state.Runs[q.RunID]
	if !ok {
		return Result{}, fail("missing_coverage", "run absent")
	}
	if r.Revision != q.ExpectedRevision || r.InputRevision != q.InputRevision {
		return Result{}, fail("stale_revision", "input or run revision changed")
	}
	if len(state.Receipts) >= 1024 || r.Revision == ^uint64(0) {
		return Result{}, fail("limit", "state capacity reached")
	}
	index := -1
	for i, old := range r.Inputs {
		if old.ID == v.ID {
			index = i
		}
	}
	context, e := s.Context(r.Selection)
	if e != nil {
		return Result{}, e
	}
	if q.Action == "prepare" {
		if index >= 0 {
			return Result{}, fail("conflict", "input already exists")
		}
		if len(r.Inputs) >= 64 {
			return Result{}, fail("limit", "owner input capacity reached")
		}
		if v.State != "" {
			return Result{}, fail("invalid_request", "caller cannot supply input state")
		}
		if v.Kind == "launch" || v.Kind == "steer" {
			if context.Revision != r.InputRevision || len(context.Conditions) > 0 {
				return Result{}, fail("stale_revision", "coverage changed")
			}
			for _, old := range r.Inputs {
				if old.State != "acknowledged" {
					return Result{}, fail("reconciliation_required", "unresolved outbound intent")
				}
			}
		}
		if v.Kind == "launch" {
			if len(r.Attempts) > 0 {
				a := r.Attempts[len(r.Attempts)-1]
				if a.Status != "completed" && a.Status != "failed" && a.Status != "interrupted" {
					return Result{}, fail("reconciliation_required", "prior host attempt is not terminal")
				}
			}
			if r.State != "ready" || v.ThreadID != "" || v.TurnID != "" {
				return Result{}, fail("conflict", "launch needs a ready run and new thread")
			}
		} else {
			if len(r.Attempts) == 0 {
				return Result{}, fail("missing_coverage", "no mapped attempt")
			}
			a := r.Attempts[len(r.Attempts)-1]
			if a.ID != v.AttemptID || a.Epoch != v.Epoch || a.ThreadID != v.ThreadID || a.TurnID != v.TurnID || v.TurnID == "" {
				return Result{}, fail("unauthorized", "input outside current host attempt")
			}
			if v.Kind == "steer" && (r.State != "running" || a.Status != "running") {
				return Result{}, fail("conflict", "steering needs observed active work")
			}
		}
		v.State = "prepared"
		r.Inputs = append(r.Inputs, v)
	} else {
		if index < 0 {
			return Result{}, fail("missing_coverage", "input absent")
		}
		old := r.Inputs[index]
		// Acknowledgment adds observed host identity; content and source provenance stay immutable.
		if old.Kind != v.Kind || old.Text != v.Text || old.Payload != v.Payload || old.AttemptID != v.AttemptID || old.Epoch != v.Epoch || old.ReceivedAt != v.ReceivedAt || old.State != v.State {
			return Result{}, fail("conflict", "input content or state changed")
		}
		if old.ThreadID != "" && old.ThreadID != v.ThreadID || old.TurnID != "" && old.TurnID != v.TurnID {
			return Result{}, fail("conflict", "host identity changed")
		}
		switch q.Action {
		case "claim":
			if v.Kind == "launch" && r.State != "ready" || v.Kind == "steer" && r.State != "running" {
				return Result{}, fail("conflict", "routine state changed before send")
			}
			if old.State != "prepared" || old.ThreadID != v.ThreadID || old.TurnID != v.TurnID {
				return Result{}, fail("reconciliation_required", "intent already consumed")
			}
			if v.Kind == "launch" || v.Kind == "steer" {
				if context.Revision != r.InputRevision || len(context.Conditions) > 0 {
					return Result{}, fail("stale_revision", "coverage changed before dispatch")
				}
			}
			if v.Kind == "steer" {
				a := r.Attempts[len(r.Attempts)-1]
				if a.ID != v.AttemptID || a.Epoch != v.Epoch || a.Status != "running" {
					return Result{}, fail("reconciliation_required", "attempt changed before steering")
				}
			}
			v.State = "sending"
		case "ack":
			if len(r.Attempts) == 0 {
				return Result{}, fail("missing_coverage", "acknowledgment needs a mapped attempt")
			}
			a := r.Attempts[len(r.Attempts)-1]
			if a.ID != v.AttemptID || a.Epoch != v.Epoch || a.ThreadID != v.ThreadID || a.TurnID != v.TurnID {
				return Result{}, fail("conflict", "acknowledgment differs from observed attempt")
			}
			if old.State != "sending" || !opaque(v.ThreadID) || !opaque(v.TurnID) {
				return Result{}, fail("reconciliation_required", "acknowledgment requires an observed response")
			}
			v.State = "acknowledged"
		case "resolve":
			if old.Kind != "launch" && old.Kind != "steer" {
				return Result{}, fail("reconciliation_required", "interrupt outcomes need their own observed resolution")
			}
			if old.State != "sending" && old.State != "uncertain" {
				return Result{}, fail("conflict", "input is not awaiting acknowledgment")
			}
			found := false
			for _, a := range r.Attempts {
				if a.ID != old.AttemptID || a.ThreadID != v.ThreadID || a.TurnID != v.TurnID || a.Status == "unknown" {
					continue
				}
				for _, event := range a.Events {
					if event.Kind == "userMessage" && event.ClientID == old.ID && event.TurnID == v.TurnID && event.Text == old.Payload && len(event.Attachments) == 0 {
						found = true
					}
				}
			}
			if !found {
				return Result{}, fail("reconciliation_required", "no exact correlated host message proves acceptance")
			}
			v.State = "acknowledged"
		case "uncertain":
			if old.State != "sending" {
				return Result{}, fail("conflict", "input not in flight")
			}
			v.State = "uncertain"
		default:
			return Result{}, fail("invalid_request", "unsupported input operation")
		}
		r.Inputs[index] = v
	}
	r.Revision++
	state.Runs[r.ID] = r
	out := Result{Schema: Schema, Operation: "input." + q.Action, ProjectRef: s.project, CheckoutID: state.CheckoutID, OperationID: q.OperationID, Run: summarizeRun(r), Coverage: "known"}
	if context.Revision != r.InputRevision || len(context.Conditions) > 0 {
		out.Coverage = "stale"
	}
	state.Receipts[q.OperationID] = receipt{digest, out}
	if e = s.save(state, head); e != nil {
		return Result{}, e
	}
	return out, nil
}

// LiveControllerLease is held for the complete terminal control session. Observer
// reads do not acquire it. A second terminal cannot take over an active owner.
func (s *Service) LiveControllerLease(c Caller) (func(), error) {
	if !c.controller {
		return nil, fail("unauthorized", "controller required")
	}
	return lock(s.state + "/live-controller")
}
