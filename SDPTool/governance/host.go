package governance

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// Host observations are execution evidence only, never routine completion.
// Text is restricted to completed visible messages by the protocol adapter.
type Attachment struct {
	Kind         string `json:"kind"`
	Reference    string `json:"reference"`
	Availability string `json:"availability"`
}

type Observation struct {
	InputDigest string       `json:"inputDigest,omitempty"`
	ClientID    string       `json:"clientId,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
	Kind        string       `json:"kind"`
	TurnID      string       `json:"turnId"`
	ItemID      string       `json:"itemId"`
	Phase       string       `json:"phase"`
	Text        string       `json:"text"`
	Status      string       `json:"status"`
	Epoch       string       `json:"epoch"`
	ReceivedAt  string       `json:"receivedAt"`
}
type Attempt struct {
	HostIdentity   string        `json:"hostIdentity,omitempty"`
	AgentEpoch     string        `json:"agentEpoch,omitempty"`
	LastObservedAt string        `json:"lastObservedAt,omitempty"`
	EventCount     int           `json:"eventCount,omitempty"`
	ID             string        `json:"id"`
	HostVersion    string        `json:"hostVersion"`
	ThreadID       string        `json:"threadId"`
	TurnID         string        `json:"turnId"`
	InputRevision  string        `json:"inputRevision"`
	Epoch          string        `json:"epoch"`
	Status         string        `json:"status"`
	Gap            bool          `json:"gap"`
	Events         []Observation `json:"events"`
}
type HostRequest struct {
	HostIdentity     string      `json:"hostIdentity,omitempty"`
	Schema           string      `json:"schema"`
	ProjectRef       string      `json:"projectRef"`
	CheckoutID       string      `json:"checkoutId"`
	OperationID      string      `json:"operationId"`
	RunID            string      `json:"runId"`
	ExpectedRevision uint64      `json:"expectedRevision"`
	InputRevision    string      `json:"inputRevision"`
	Action           string      `json:"action"`
	AttemptID        string      `json:"attemptId"`
	HostVersion      string      `json:"hostVersion"`
	ThreadID         string      `json:"threadId"`
	Observation      Observation `json:"observation"`
}

// CapabilityEpoch preserves legacy state decoding; only trusted host operations change it.
func (a Attempt) CapabilityEpoch() string {
	if a.AgentEpoch != "" {
		return a.AgentEpoch
	}
	return a.Epoch
}

func opaque(s string) bool { return s != "" && len(s) <= 256 && !strings.ContainsAny(s, "\x00\r\n") }
func hostStatus(s string) bool {
	switch s {
	case "idle", "running", "completed", "failed", "interrupted", "unknown":
		return true
	}
	return false
}
func validObservation(o Observation) bool {
	encoded, err := json.Marshal(o)
	if err != nil || len(encoded) > 240<<10 {
		return false
	}
	if o.InputDigest != "" && (o.Kind != "userMessage" || (o.InputDigest != "unavailable" && !digestPattern.MatchString(o.InputDigest))) {
		return false
	}
	if o.ClientID != "" && (o.Kind != "userMessage" || !uuidPattern.MatchString(o.ClientID)) {
		return false
	}
	if len(o.Attachments) > 64 || (len(o.Attachments) > 0 && o.Kind != "userMessage") {
		return false
	}
	for _, a := range o.Attachments {
		if len(a.Reference) > 4096 || strings.ContainsAny(a.Reference, "\x00\r\n") {
			return false
		}
		switch a.Kind {
		case "image", "localImage", "audio", "localAudio", "skill", "mention":
		default:
			return false
		}
		if a.Availability != "unknown" && a.Availability != "unavailable" {
			return false
		}
	}
	if !uuidPattern.MatchString(o.Epoch) || len(o.Text) > 64<<10 {
		return false
	}
	if _, e := time.Parse(time.RFC3339Nano, o.ReceivedAt); e != nil {
		return false
	}
	if o.TurnID != "" && !opaque(o.TurnID) || o.ItemID != "" && !opaque(o.ItemID) {
		return false
	}
	switch o.Kind {
	case "bound", "disconnected", "reconciled":
		return o.ItemID == "" && o.Text == "" && o.Phase == "" && hostStatus(o.Status)
	case "turn.started":
		return opaque(o.TurnID) && o.ItemID == "" && o.Text == "" && o.Phase == "" && o.Status == "running"
	case "turn.completed":
		return opaque(o.TurnID) && o.ItemID == "" && o.Text == "" && o.Phase == "" && (o.Status == "completed" || o.Status == "failed" || o.Status == "interrupted")
	case "context.observed":
		var p ContextProvenance
		return o.TurnID == "" && o.ItemID == "launch" && o.Status == "" && o.Phase == "" && Decode([]byte(o.Text), &p) == nil && p.Valid()
	case "approval.requested":
		var p ApprovalPrompt
		return opaque(o.TurnID) && opaque(o.ItemID) && o.Phase == "" && o.Status == "" && Decode([]byte(o.Text), &p) == nil && p.Valid() && p.RequestID == o.ItemID && p.Epoch == o.Epoch && p.TurnID == o.TurnID
	case "approval.resolved":
		return opaque(o.TurnID) && opaque(o.ItemID) && o.Text == "" && o.Phase == "" && o.Status == ""
	case "userMessage":
		return opaque(o.TurnID) && opaque(o.ItemID) && o.Phase == "" && o.Status == ""
	case "agentMessage":
		return opaque(o.TurnID) && opaque(o.ItemID) && (o.Phase == "commentary" || o.Phase == "final_answer" || o.Phase == "unknown") && o.Status == ""
	case "tool.completed":
		return opaque(o.TurnID) && opaque(o.ItemID) && o.Text == "" && o.Phase == "" && (o.Status == "completed" || o.Status == "failed" || o.Status == "declined")
	}
	return false
}

// RecordHost accepts only observations supplied by the trusted controller.
// A new adapter epoch records a gap, even after a successful history read.
func (s *Service) RecordHost(caller Caller, q HostRequest) (Result, error) {
	if !caller.controller {
		return Result{}, fail("unauthorized", "host observations require controller provenance")
	}
	if !uuidPattern.MatchString(q.RunID) || !uuidPattern.MatchString(q.AttemptID) || !digestPattern.MatchString(q.InputRevision) || q.ExpectedRevision == 0 || !opaque(q.ThreadID) || !validObservation(q.Observation) {
		return Result{}, fail("invalid_request", "invalid bounded host observation")
	}
	if q.HostVersion != "0.160.1" && q.HostVersion != "0.161.0" {
		return Result{}, fail("incompatible_profile", "only measured Codex 0.160.1 and 0.161.0 are supported")
	}
	if q.HostIdentity != "" && !digestPattern.MatchString(q.HostIdentity) {
		return Result{}, fail("invalid_request", "invalid host identity")
	}
	o := q.Observation
	if q.Action == "bind" && o.Kind != "bound" || q.Action == "disconnect" && o.Kind != "disconnected" || (q.Action == "reconcile" || q.Action == "reattach") && o.Kind != "reconciled" {
		return Result{}, fail("invalid_request", "action/observation mismatch")
	}
	if q.Action == "observe" && (o.Kind == "bound" || o.Kind == "disconnected" || o.Kind == "reconciled") {
		return Result{}, fail("invalid_request", "use the matching host operation")
	}
	if q.Action != "bind" && q.Action != "observe" && q.Action != "disconnect" && q.Action != "reconcile" && q.Action != "reattach" {
		return Result{}, fail("invalid_request", "unsupported host operation")
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
	if e = s.authorize(caller, q.Schema, q.ProjectRef, q.CheckoutID, q.OperationID, state); e != nil {
		return Result{}, e
	}
	digest := requestDigest(q)
	if r, ok, e := retry(state, q.OperationID, digest); ok {
		return r, e
	}
	r, ok := state.Runs[q.RunID]
	if !ok {
		return Result{}, fail("missing_coverage", "run not found")
	}
	if r.Revision != q.ExpectedRevision {
		return Result{}, &Error{Code: "stale_revision", Message: "run revision changed", CurrentRevision: r.Revision}
	}
	if r.InputRevision != q.InputRevision {
		return Result{}, fail("stale_revision", "observation belongs to another assignment revision")
	}
	if len(state.Receipts) >= 1024 || r.Revision == ^uint64(0) {
		return Result{}, fail("limit", "state capacity reached")
	}
	a := -1
	for i, v := range r.Attempts {
		if v.ID == q.AttemptID {
			a = i
		}
	}
	if q.Action == "bind" {
		if a >= 0 {
			return Result{}, fail("conflict", "attempt already bound; retry original operation")
		}
		if r.State != "ready" && r.State != "running" {
			return Result{}, fail("conflict", "run not launchable")
		}
		c, e := s.runContext(r)
		if e != nil {
			return Result{}, e
		}
		if c.Revision != r.InputRevision || len(c.Conditions) > 0 {
			return Result{}, fail("stale_revision", "work coverage changed")
		}
		if len(r.Attempts) >= 16 {
			return Result{}, fail("limit", "attempt capacity reached")
		}
		if o.Status != "idle" || o.TurnID != "" {
			return Result{}, fail("invalid_request", "bind begins with observed idle thread")
		}
		for _, run := range state.Runs {
			for _, v := range run.Attempts {
				if v.ThreadID == q.ThreadID {
					return Result{}, fail("conflict", "thread already belongs to an attempt")
				}
			}
		}
		r.Attempts = append(r.Attempts, Attempt{HostIdentity: q.HostIdentity, AgentEpoch: o.Epoch, ID: q.AttemptID, HostVersion: q.HostVersion, ThreadID: q.ThreadID, InputRevision: r.InputRevision, Epoch: o.Epoch, Status: "idle", Events: []Observation{}})
		a = len(r.Attempts) - 1
	} else if a < 0 {
		return Result{}, fail("missing_coverage", "attempt not bound")
	}
	at := &r.Attempts[a]
	if at.ThreadID != q.ThreadID || at.HostVersion != q.HostVersion {
		return Result{}, fail("unauthorized", "observed host identity mismatch")
	}
	if o.Kind == "context.observed" {
		var p ContextProvenance
		if Decode([]byte(o.Text), &p) != nil || p.InputRevision != r.InputRevision || a != len(r.Attempts)-1 || at.Status != "idle" || at.TurnID != "" {
			return Result{}, fail("unauthorized", "context observation requires current launch binding")
		}
	}
	if o.Kind == "approval.requested" {
		var p ApprovalPrompt
		if Decode([]byte(o.Text), &p) != nil || p.ThreadID != q.ThreadID || p.Request.RunID != r.ID || p.Request.ProjectRef != s.project || p.Request.CheckoutID != state.CheckoutID || p.Request.InputRevision != r.InputRevision || at.Status != "running" || r.State != "running" {
			return Result{}, fail("unauthorized", "approval request is outside current running assignment")
		}
	}
	if q.Action == "observe" && at.Status == "unknown" {
		return Result{}, fail("reconciliation_required", "disconnected attempt requires observed history")
	}
	if q.Action == "reattach" {
		if a != len(r.Attempts)-1 || r.State != "running" || at.HostIdentity == "" || at.HostIdentity != q.HostIdentity || at.TurnID == "" || at.TurnID != o.TurnID || o.Status != "running" || (at.Status != "running" && at.Status != "unknown") {
			return Result{}, fail("unauthorized", "reattachment requires the same active host, attempt and turn")
		}
		for _, input := range r.Inputs {
			if input.State != "acknowledged" {
				return Result{}, fail("reconciliation_required", "outbound input remains unresolved")
			}
		}
		// A disconnected attempt may still have an unresolved prompt. Do not hide it
		// by changing epoch, even when no answer was prepared.
		pending := map[string]bool{}
		for _, event := range at.Events {
			if event.Epoch != at.Epoch {
				continue
			}
			if event.Kind == "approval.requested" {
				pending[event.ItemID] = true
			}
			if event.Kind == "approval.resolved" {
				delete(pending, event.ItemID)
			}
		}
		if len(pending) > 0 {
			return Result{}, fail("reconciliation_required", "approval closure remains unresolved")
		}
		coverage, err := s.runContext(r)
		if err != nil {
			return Result{}, err
		}
		if coverage.Revision != r.InputRevision || len(coverage.Conditions) > 0 {
			return Result{}, fail("stale_revision", "reattachment coverage changed")
		}
		at.AgentEpoch = at.CapabilityEpoch()
	}
	if q.Action == "reconcile" || q.Action == "reattach" {
		if o.Status == "unknown" {
			return Result{}, fail("observation_gap", "history does not establish host state")
		}
		if at.Epoch != o.Epoch {
			at.Gap = true
		}
		if q.Action == "reconcile" {
			at.AgentEpoch = o.Epoch
		}
		at.Epoch = o.Epoch
	} else if at.Epoch != o.Epoch {
		return Result{}, fail("reconciliation_required", "new connection epoch requires reconciliation")
	}
	if q.Action == "disconnect" {
		if o.Status != "unknown" {
			return Result{}, fail("invalid_request", "disconnect must retain uncertainty")
		}
		at.Gap = true
	}
	if o.TurnID != "" {
		if at.TurnID != "" && at.TurnID != o.TurnID {
			return Result{}, fail("conflict", "attempt is already bound to another turn")
		}
		at.TurnID = o.TurnID
	}
	duplicate := false
	if q.Action == "observe" {
		for _, old := range at.Events {
			if old.Kind == o.Kind && old.TurnID == o.TurnID && old.ItemID == o.ItemID && (!strings.HasPrefix(o.Kind, "approval.") || old.Epoch == o.Epoch) {
				if (old.ClientID != "" && old.ClientID != o.ClientID) || (old.InputDigest != "" && old.InputDigest != o.InputDigest) || old.Phase != o.Phase || old.Text != o.Text || old.Status != o.Status || requestDigest(old.Attachments) != requestDigest(o.Attachments) {
					return Result{}, fail("observation_gap", "conflicting repeated completed item")
				}
				duplicate = true
				break
			}
		}
	}
	if !duplicate && o.Kind == "turn.started" && (at.Status == "completed" || at.Status == "failed" || at.Status == "interrupted") {
		return Result{}, fail("observation_gap", "late start cannot regress terminal host status")
	}
	if !duplicate {
		if len(at.Events) >= 1024 {
			return Result{}, fail("limit", "observation capacity reached")
		}
		at.Events = append(at.Events, o)
		if o.Status != "" && o.Kind != "tool.completed" {
			at.Status = o.Status
		}
	}
	r.Revision++
	state.Runs[r.ID] = r
	out := Result{Schema: Schema, Operation: "host." + q.Action, ProjectRef: s.project, CheckoutID: state.CheckoutID, OperationID: q.OperationID, Run: summarizeRun(r), Coverage: "known"}
	c, e := s.runContext(r)
	if e != nil {
		return Result{}, e
	}
	if c.Revision != r.InputRevision || len(c.Conditions) > 0 {
		out.Coverage = "stale"
		if q.Action == "bind" {
			return Result{}, fail("stale_revision", "coverage changed during bind")
		}
	}
	state.Receipts[q.OperationID] = receipt{digest, out}
	if e = s.save(state, head); e != nil {
		return Result{}, e
	}
	return out, nil
}

type EventPage struct {
	RunID      string        `json:"runId"`
	AttemptID  string        `json:"attemptId"`
	Revision   uint64        `json:"revision"`
	Events     []Observation `json:"events"`
	NextCursor string        `json:"nextCursor"`
	Gap        bool          `json:"gap"`
}

// Events returns at most 256 records and 240 KiB of encoded observations. A cursor pins the exact run revision;
// callers restart paging after any concurrent change instead of silently losing data.
func (s *Service) Events(runID, attemptID, cursor string, limit int) (EventPage, error) {
	if limit < 1 || limit > 256 {
		return EventPage{}, fail("invalid_request", "page limit must be 1..256")
	}
	if !uuidPattern.MatchString(runID) {
		return EventPage{}, fail("invalid_request", "invalid runId")
	}
	state, _, e := s.load()
	if e != nil {
		return EventPage{}, e
	}
	r, ok := state.Runs[runID]
	if !ok {
		return EventPage{}, fail("missing_coverage", "run not found")
	}
	for _, a := range r.Attempts {
		if a.ID != attemptID {
			continue
		}
		offset := 0
		if cursor != "" {
			prefix := fmt.Sprintf("%s/%s/%d/", runID, attemptID, r.Revision)
			if !strings.HasPrefix(cursor, prefix) {
				return EventPage{}, fail("stale_revision", "cursor identity or revision changed")
			}
			tail := strings.TrimPrefix(cursor, prefix)
			offset, e = strconv.Atoi(tail)
			if e != nil || strconv.Itoa(offset) != tail || offset < 0 || offset > len(a.Events) {
				return EventPage{}, fail("invalid_request", "invalid cursor")
			}
		}
		end := offset
		size := 0
		for end < len(a.Events) && end-offset < limit {
			b, _ := json.Marshal(a.Events[end])
			if size+len(b) > 240<<10 {
				break
			}
			size += len(b)
			end++
		}
		if end == offset && offset < len(a.Events) {
			return EventPage{}, fail("limit", "one observation exceeds page byte limit")
		}
		page := EventPage{RunID: runID, AttemptID: attemptID, Revision: r.Revision, Events: append([]Observation{}, a.Events[offset:end]...), Gap: a.Gap}
		if end < len(a.Events) {
			page.NextCursor = fmt.Sprintf("%s/%s/%d/%d", runID, attemptID, r.Revision, end)
		}
		return page, nil
	}
	return EventPage{}, fail("missing_coverage", "attempt not found")
}

// CheckWire bounds and rejects duplicate keys in transport JSON before selective
// decoding. Unlike Decode, it permits additional host fields and nulls: adapters
// explicitly allowlist what they retain, without trusting unknown host metadata.
func CheckWire(b []byte) error {
	var v json.RawMessage
	if len(b) > MaxRequest || !utf8.Valid(b) {
		return fail("limit", "frame exceeds 1 MiB")
	}
	return decode(b, &v)
}

// Routine snapshots expose attempt metadata/counts. Message content is available
// only through byte- and count-bounded Events pages, including in receipts.
func summarizeRun(r Run) Run {
	r.Attention = append([]AttentionQuestion(nil), r.Attention...)
	for i := range r.Attention {
		r.Attention[i].Options = append([]AttentionOption{}, r.Attention[i].Options...)
		if v := r.Attention[i].Reply; v != nil {
			copy := *v
			r.Attention[i].Reply = &copy
		}
		if v := r.Attention[i].Withdrawal; v != nil {
			copy := *v
			r.Attention[i].Withdrawal = &copy
		}
	}
	r.Inputs = append([]OwnerInput(nil), r.Inputs...)
	for i := range r.Inputs {
		r.Inputs[i].ContentOmitted = true
		r.Inputs[i].Text = ""
		r.Inputs[i].Payload = ""
		r.Inputs[i].Skills = nil
	}
	r.Attempts = append([]Attempt(nil), r.Attempts...)
	for i := range r.Attempts {
		r.Attempts[i].EventCount = len(r.Attempts[i].Events)
		if n := len(r.Attempts[i].Events); n > 0 {
			r.Attempts[i].LastObservedAt = r.Attempts[i].Events[n-1].ReceivedAt
		}
		r.Attempts[i].Events = []Observation{}
	}
	return r
}
