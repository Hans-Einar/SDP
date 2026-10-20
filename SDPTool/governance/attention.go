package governance

import (
	"sort"
	"strings"
	"time"
)

const AttentionSchema = "sdp-governance-attention/1"
const attentionStateSchema = "sdp-project-governance-state/0.6"

type AttentionOption struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}
type AttentionQuestion struct {
	ID             string               `json:"id"`
	Basis          string               `json:"basis"`
	InputRevision  string               `json:"inputRevision"`
	RaisedRevision uint64               `json:"raisedRevision"`
	RaisedAt       string               `json:"raisedAt"`
	Origin         string               `json:"origin"`
	AttemptID      string               `json:"attemptId,omitempty"`
	Prompt         string               `json:"prompt"`
	Options        []AttentionOption    `json:"options"`
	Blocking       bool                 `json:"blocking"`
	Reply          *AttentionReply      `json:"reply,omitempty"`
	Withdrawal     *AttentionWithdrawal `json:"withdrawal,omitempty"`
}
type AttentionReply struct {
	OperationID string `json:"operationId"`
	Actor       string `json:"actor"`
	Source      string `json:"source"`
	OptionID    string `json:"optionId,omitempty"`
	Text        string `json:"text"`
	RecordedAt  string `json:"recordedAt"`
}
type AttentionWithdrawal struct {
	OperationID string `json:"operationId"`
	Reason      string `json:"reason"`
	RecordedAt  string `json:"recordedAt"`
}
type AttentionRequest struct {
	Schema           string            `json:"schema"`
	ProjectRef       string            `json:"projectRef"`
	CheckoutID       string            `json:"checkoutId"`
	OperationID      string            `json:"operationId"`
	RunID            string            `json:"runId"`
	InputRevision    string            `json:"inputRevision"`
	ExpectedRevision uint64            `json:"expectedRevision"`
	Prompt           string            `json:"prompt"`
	Options          []AttentionOption `json:"options,omitempty"`
	Blocking         bool              `json:"blocking"`
}
type AttentionDisposition struct {
	Schema      string `json:"schema"`
	ProjectRef  string `json:"projectRef"`
	CheckoutID  string `json:"checkoutId"`
	OperationID string `json:"operationId"`
	RunID       string `json:"runId"`
	QuestionID  string `json:"questionId"`
	Basis       string `json:"basis"`
}

// Supplied by the trusted controller separately from request JSON. Attribution
// follows local controller authority; it does not authenticate a human identity.
type OwnerAttentionInput struct {
	Actor    string `json:"actor"`
	Source   string `json:"source"`
	OptionID string `json:"optionId"`
	Text     string `json:"text"`
}
type AttentionView struct {
	RunID           string            `json:"runId"`
	Question        AttentionQuestion `json:"question"`
	Status          string            `json:"status"`
	Coverage        string            `json:"coverage"`
	Applicability   string            `json:"applicability"`
	CurrentRevision uint64            `json:"currentRevision"`
	Consumption     string            `json:"consumption"`
}
type AttentionInbox struct {
	Schema     string          `json:"schema"`
	ProjectRef string          `json:"projectRef"`
	CheckoutID string          `json:"checkoutId"`
	Generation string          `json:"generation"`
	Items      []AttentionView `json:"items"`
}

func boundedAttentionText(s string, limit int) bool {
	return strings.TrimSpace(s) != "" && len(s) <= limit
}
func validAttentionQuestion(prompt string, options []AttentionOption) bool {
	if !boundedAttentionText(prompt, 4096) || len(options) > 5 {
		return false
	}
	seen := map[string]bool{}
	for _, o := range options {
		if !boundedAttentionText(o.ID, 64) || !boundedAttentionText(o.Label, 1024) || seen[o.ID] {
			return false
		}
		seen[o.ID] = true
	}
	return true
}

// Ignore only observation/progress bookkeeping. Owner input content and identity,
// source revision, routine state and attempt replacement remain decision inputs.
func attentionBasis(r Run) string {
	r.Revision = 0
	r.Attention = nil
	r.Inputs = append([]OwnerInput(nil), r.Inputs...)
	for i := range r.Inputs {
		r.Inputs[i].State = ""
	}
	r.Attempts = append([]Attempt(nil), r.Attempts...)
	for i := range r.Attempts {
		r.Attempts[i].Events = nil
		r.Attempts[i].EventCount = 0
		r.Attempts[i].LastObservedAt = ""
		r.Attempts[i].Status = ""
	}
	return requestDigest(r)
}
func attentionCaller(c Caller, r Run) error {
	if c.controller {
		return nil
	}
	if c.runID != r.ID || len(r.Attempts) == 0 {
		return fail("unauthorized", "question is outside assigned run")
	}
	a := r.Attempts[len(r.Attempts)-1]
	if a.ID != c.attemptID || a.CapabilityEpoch() != c.epoch || a.Status != "running" || a.Gap {
		return fail("unauthorized", "question requires current running attempt")
	}
	return nil
}
func hasBlockingAttention(r Run) bool {
	for _, q := range r.Attention {
		if q.Blocking && q.Reply == nil && q.Withdrawal == nil {
			return true
		}
	}
	return false
}
func (s *Service) currentAttentionInputs(r Run) error {
	c, e := s.runContext(r)
	if e != nil {
		return e
	}
	if c.Revision != r.InputRevision || len(c.Conditions) > 0 {
		return fail("stale_revision", "question sources or authority changed")
	}
	return nil
}
func (s *Service) saveAttention(state snapshot, head string, r Run, operation, op, digest string) (Result, error) {
	if len(state.Receipts) >= 1024 || r.Revision == ^uint64(0) {
		return Result{}, fail("limit", "operational receipt or revision capacity reached")
	}
	r.Revision++
	state.Schema = attentionStateSchema
	if state.Prepared == nil {
		state.Prepared = map[string]PreparedAssignment{}
	}
	if state.PreparationReceipts == nil {
		state.PreparationReceipts = map[string]preparationReceipt{}
	}
	if state.DesignDecisions == nil {
		state.DesignDecisions = map[string]string{}
	}
	if state.DesignDecisionReceipts == nil {
		state.DesignDecisionReceipts = map[string]designDecisionReceipt{}
	}
	state.Runs[r.ID] = r
	out := Result{Schema: Schema, Operation: operation, OperationID: op, ProjectRef: s.project, CheckoutID: state.CheckoutID, Run: summarizeRun(r), Coverage: "known"}
	if e := s.currentAttentionInputs(r); e != nil {
		if operation != "attention.withdraw" {
			return Result{}, e
		}
		out.Coverage = "stale"
	}
	state.Receipts[op] = receipt{Digest: digest, Result: out}
	if e := s.save(state, head); e != nil {
		return Result{}, e
	}
	return out, nil
}
func (s *Service) RequestAttention(caller Caller, q AttentionRequest) (Result, error) {
	if !caller.controller && caller.runID == "" {
		return Result{}, fail("unauthorized", "trusted controller or assigned agent required")
	}
	if q.Schema != AttentionSchema {
		return Result{}, fail("incompatible_profile", "unsupported attention schema")
	}
	if !uuidPattern.MatchString(q.RunID) || !digestPattern.MatchString(q.InputRevision) || q.ExpectedRevision == 0 || !validAttentionQuestion(q.Prompt, q.Options) {
		return Result{}, fail("invalid_request", "bounded question and current run revision required")
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
	if e = s.authorize(caller, Schema, q.ProjectRef, q.CheckoutID, q.OperationID, state); e != nil {
		return Result{}, e
	}
	r, ok := state.Runs[q.RunID]
	if !ok {
		return Result{}, fail("missing_coverage", "run absent")
	}
	if e = attentionCaller(caller, r); e != nil {
		return Result{}, e
	}
	digest := requestDigest(struct {
		Operation  string
		Request    AttentionRequest
		Attempt    string
		Controller bool
	}{"attention.request", q, caller.attemptID, caller.controller})
	if old, ok, e := retry(state, q.OperationID, digest); ok {
		return old, e
	}
	if r.State != "running" {
		return Result{}, fail("conflict", "questions require running work")
	}
	if q.InputRevision != r.InputRevision || (q.ExpectedRevision != r.Revision && !observedAgentRevision(state, r, caller, q.ExpectedRevision)) {
		return Result{}, fail("stale_revision", "run or input revision changed")
	}
	if e = s.currentAttentionInputs(r); e != nil {
		return Result{}, e
	}
	if len(r.Attention) >= 32 {
		return Result{}, fail("limit", "question capacity reached")
	}
	origin := "agent"
	if caller.controller {
		origin = "controller"
	}
	question := AttentionQuestion{ID: q.OperationID, Basis: attentionBasis(r), InputRevision: r.InputRevision, RaisedRevision: r.Revision + 1, RaisedAt: time.Now().UTC().Format(time.RFC3339Nano), Origin: origin, AttemptID: caller.attemptID, Prompt: q.Prompt, Options: append([]AttentionOption{}, q.Options...), Blocking: q.Blocking}
	r.Attention = append(r.Attention, question)
	if e = s.currentAttentionInputs(r); e != nil {
		return Result{}, e
	}
	return s.saveAttention(state, head, r, "attention.request", q.OperationID, digest)
}
func (s *Service) ReplyAttention(caller Caller, q AttentionDisposition, input OwnerAttentionInput) (Result, error) {
	if !boundedAttentionText(input.Actor, 256) || !boundedAttentionText(input.Source, 1024) || !boundedAttentionText(input.Text, 4096) || len(input.OptionID) > 64 {
		return Result{}, fail("invalid_request", "attributed bounded reply required")
	}
	return s.resolveAttention(caller, q, &input, "")
}
func (s *Service) WithdrawAttention(caller Caller, q AttentionDisposition, reason string) (Result, error) {
	if !boundedAttentionText(reason, 4096) {
		return Result{}, fail("invalid_request", "withdrawal reason required")
	}
	return s.resolveAttention(caller, q, nil, reason)
}
func (s *Service) resolveAttention(caller Caller, q AttentionDisposition, input *OwnerAttentionInput, reason string) (Result, error) {
	if !caller.controller {
		return Result{}, fail("unauthorized", "only the controller may reply or withdraw")
	}
	if q.Schema != AttentionSchema {
		return Result{}, fail("incompatible_profile", "unsupported attention schema")
	}
	if !uuidPattern.MatchString(q.RunID) || !uuidPattern.MatchString(q.QuestionID) || !digestPattern.MatchString(q.Basis) {
		return Result{}, fail("invalid_request", "exact question and basis required")
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
	if e = s.authorize(caller, Schema, q.ProjectRef, q.CheckoutID, q.OperationID, state); e != nil {
		return Result{}, e
	}
	operation := "attention.reply"
	if input == nil {
		operation = "attention.withdraw"
	}
	digest := requestDigest(struct {
		Operation string
		Request   AttentionDisposition
		Input     *OwnerAttentionInput
		Reason    string
	}{operation, q, input, reason})
	if old, ok, e := retry(state, q.OperationID, digest); ok {
		return old, e
	}
	r, ok := state.Runs[q.RunID]
	if !ok {
		return Result{}, fail("missing_coverage", "run absent")
	}
	for i := range r.Attention {
		question := &r.Attention[i]
		if question.ID != q.QuestionID {
			continue
		}
		if question.Basis != q.Basis {
			return Result{}, fail("stale_revision", "question basis differs")
		}
		if question.Reply != nil || question.Withdrawal != nil {
			return Result{}, fail("conflict", "question already resolved")
		}
		now := time.Now().UTC().Format(time.RFC3339Nano)
		if input != nil {
			if attentionBasis(r) != q.Basis {
				return Result{}, fail("stale_revision", "question is no longer applicable")
			}
			if e = s.currentAttentionInputs(r); e != nil {
				return Result{}, e
			}
			found := input.OptionID == ""
			for _, option := range question.Options {
				if option.ID == input.OptionID {
					found = true
				}
			}
			if !found {
				return Result{}, fail("invalid_request", "unknown question option")
			}
			question.Reply = &AttentionReply{OperationID: q.OperationID, Actor: input.Actor, Source: input.Source, OptionID: input.OptionID, Text: input.Text, RecordedAt: now}
			if e = s.currentAttentionInputs(r); e != nil {
				return Result{}, e
			}
		} else {
			question.Withdrawal = &AttentionWithdrawal{OperationID: q.OperationID, Reason: reason, RecordedAt: now}
		}
		return s.saveAttention(state, head, r, operation, q.OperationID, digest)
	}
	return Result{}, fail("missing_coverage", "question absent from selected run")
}
func (s *Service) AttentionInbox(runID string) (AttentionInbox, error) {
	out := AttentionInbox{Schema: AttentionSchema, ProjectRef: s.project, Items: []AttentionView{}}
	if runID != "" && !uuidPattern.MatchString(runID) {
		return out, fail("invalid_request", "invalid runId")
	}
	state, head, e := s.load()
	if e != nil {
		return out, e
	}
	out.CheckoutID, out.Generation = state.CheckoutID, head
	if runID != "" {
		if _, ok := state.Runs[runID]; !ok {
			return out, fail("missing_coverage", "run absent")
		}
	}
	for id, r := range state.Runs {
		if runID != "" && id != runID {
			continue
		}
		coverage := "known"
		if e = s.currentAttentionInputs(r); e != nil {
			coverage = "stale"
		}
		basis := attentionBasis(r)
		for _, q := range r.Attention {
			v := AttentionView{RunID: id, Question: q, Status: "open", Coverage: coverage, Applicability: "current", CurrentRevision: r.Revision, Consumption: "not-tracked"}
			if q.Basis != basis || coverage != "known" {
				v.Applicability = "stale"
			}
			if q.Reply != nil {
				v.Status = "answered"
			}
			if q.Withdrawal != nil {
				v.Status = "withdrawn"
			}
			out.Items = append(out.Items, v)
		}
	}
	sort.Slice(out.Items, func(i, j int) bool {
		a, b := out.Items[i], out.Items[j]
		if a.RunID != b.RunID {
			return a.RunID < b.RunID
		}
		return a.Question.RaisedRevision < b.Question.RaisedRevision
	})
	_, latest, e := s.load()
	if e != nil {
		return out, e
	}
	if latest != head {
		return out, fail("stale_revision", "attention changed during inspection")
	}
	return out, nil
}

func validateAttentionState(state snapshot) error {
	count := 0
	seen := map[string]bool{}
	for _, r := range state.Runs {
		if len(r.Attention) > 32 {
			return fail("reconciliation_required", "too many retained questions")
		}
		for _, q := range r.Attention {
			count++
			_, timeErr := time.Parse(time.RFC3339Nano, q.RaisedAt)
			if state.Schema != attentionStateSchema || !uuidPattern.MatchString(q.ID) || seen[q.ID] || !digestPattern.MatchString(q.Basis) || q.InputRevision != r.InputRevision || q.RaisedRevision == 0 || q.RaisedRevision > r.Revision || timeErr != nil || !validAttentionQuestion(q.Prompt, q.Options) || (q.Origin != "agent" && q.Origin != "controller") || (q.Origin == "agent" && !uuidPattern.MatchString(q.AttemptID)) || (q.Reply != nil && q.Withdrawal != nil) {
				return fail("reconciliation_required", "invalid retained question")
			}
			seen[q.ID] = true
			receipt, ok := state.Receipts[q.ID]
			if !ok || receipt.Result.Operation != "attention.request" || receipt.Result.Run.ID != r.ID {
				return fail("reconciliation_required", "question receipt missing")
			}
			if q.Reply != nil {
				a := q.Reply
				_, timeErr = time.Parse(time.RFC3339Nano, a.RecordedAt)
				receipt, ok := state.Receipts[a.OperationID]
				found := a.OptionID == ""
				for _, o := range q.Options {
					if o.ID == a.OptionID {
						found = true
					}
				}
				if !ok || receipt.Result.Operation != "attention.reply" || receipt.Result.Run.ID != r.ID || !found || timeErr != nil || !boundedAttentionText(a.Actor, 256) || !boundedAttentionText(a.Source, 1024) || !boundedAttentionText(a.Text, 4096) {
					return fail("reconciliation_required", "invalid retained reply")
				}
			}
			if q.Withdrawal != nil {
				a := q.Withdrawal
				_, timeErr = time.Parse(time.RFC3339Nano, a.RecordedAt)
				receipt, ok := state.Receipts[a.OperationID]
				if !ok || receipt.Result.Operation != "attention.withdraw" || receipt.Result.Run.ID != r.ID || timeErr != nil || !boundedAttentionText(a.Reason, 4096) {
					return fail("reconciliation_required", "invalid retained withdrawal")
				}
			}
		}
	}
	if state.Schema == attentionStateSchema && count == 0 {
		return fail("reconciliation_required", "attention profile has no questions")
	}
	return nil
}
