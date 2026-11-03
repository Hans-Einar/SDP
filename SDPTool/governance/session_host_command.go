package governance

import (
	"path/filepath"
	"reflect"
	"slices"
	"sort"
	"time"
)

const SessionHostSchema = "sdp-session-host/0.1"
const sessionHostStateSchema = "sdp-project-governance-state/0.23"

func sessionRebindProfile(schema string) bool {
	return schema == sessionRebindStateSchema || schema == sessionHostStateSchema
}

type SessionHostRequest struct {
	Schema              string     `json:"schema"`
	ProjectRef          string     `json:"projectRef"`
	CheckoutID          string     `json:"checkoutId"`
	OperationID         string     `json:"operationId"`
	InstanceID          string     `json:"instanceId"`
	ExpectedRevision    uint64     `json:"expectedRevision"`
	RunID               string     `json:"runId"`
	ExpectedRunRevision uint64     `json:"expectedRunRevision"`
	Action              string     `json:"action"`
	Reason              string     `json:"reason"`
	LeaderSkill         SkillInput `json:"leaderSkill"`
}
type SessionHostCommand struct {
	Request       SessionHostRequest `json:"request"`
	InputRevision string             `json:"inputRevision"`
	AttemptID     string             `json:"attemptId,omitempty"`
	State         string             `json:"state"`
	CreatedAt     string             `json:"createdAt"`
	UpdatedAt     string             `json:"updatedAt"`
	Detail        string             `json:"detail"`
}
type SessionHostBasis struct {
	ProjectRef    string `json:"projectRef"`
	CheckoutID    string `json:"checkoutId"`
	InstanceID    string `json:"instanceId"`
	Revision      uint64 `json:"revision"`
	RunID         string `json:"runId"`
	RunRevision   uint64 `json:"runRevision"`
	RunState      string `json:"runState"`
	AttemptStatus string `json:"attemptStatus,omitempty"`
}
type SessionHostCommands struct {
	Current    *SessionHostBasis    `json:"current,omitempty"`
	Schema     string               `json:"schema"`
	Operation  string               `json:"operation"`
	InstanceID string               `json:"instanceId"`
	Commands   []SessionHostCommand `json:"commands"`
}

func validSessionHostRequest(q SessionHostRequest) error {
	if q.Schema != SessionHostSchema || !uuidPattern.MatchString(q.OperationID) || !uuidPattern.MatchString(q.InstanceID) || !uuidPattern.MatchString(q.RunID) || q.ExpectedRevision == 0 || q.ExpectedRunRevision == 0 || !validLeaderText(q.Reason, 4096) || !slices.Contains([]string{"launch", "continue", "interrupt", "scope-change"}, q.Action) || q.LeaderSkill.Name != "sdp-project-leader" || !ValidSkillInputs([]SkillInput{q.LeaderSkill}) {
		return fail("invalid_request", "exact Session/run revisions, action, reason and named leader skill required")
	}
	return nil
}
func sessionHostPending(state string) bool {
	return state == "queued" || state == "claimed" || state == "uncertain"
}
func sessionHostHold(st snapshot, id string) error {
	for _, c := range st.SessionHostCommands {
		if c.Request.InstanceID == id && sessionHostPending(c.State) {
			return fail("reconciliation_required", "pending Session host request must be reconciled first")
		}
	}
	return nil
}
func (s *Service) ReadSessionHostCommands(id string) (SessionHostCommands, error) {
	out := SessionHostCommands{Schema: SessionHostSchema, Operation: "session-host-read", InstanceID: id, Commands: []SessionHostCommand{}}
	st, _, e := s.load()
	if e != nil {
		return out, e
	}
	i, ok := st.Sessions[id]
	if !ok {
		return out, fail("missing_coverage", "Session absent")
	}
	if len(i.Work) > 0 && i.Work[0].Effects.Status == "committed" {
		run := st.Runs[i.Work[0].Request.OperationID]
		basis := SessionHostBasis{ProjectRef: st.Project, CheckoutID: st.CheckoutID, InstanceID: id, Revision: i.Revision, RunID: run.ID, RunRevision: run.Revision, RunState: run.State}
		if len(run.Attempts) > 0 {
			basis.AttemptStatus = run.Attempts[len(run.Attempts)-1].Status
		}
		out.Current = &basis
	}
	for _, c := range st.SessionHostCommands {
		if c.Request.InstanceID == id {
			out.Commands = append(out.Commands, c)
		}
	}
	sort.Slice(out.Commands, func(i, j int) bool {
		if out.Commands[i].CreatedAt == out.Commands[j].CreatedAt {
			return out.Commands[i].Request.OperationID < out.Commands[j].Request.OperationID
		}
		return out.Commands[i].CreatedAt < out.Commands[j].CreatedAt
	})
	return out, nil
}
func (s *Service) ReadSessionHostCommand(id string) (SessionHostCommand, error) {
	st, _, e := s.load()
	if e != nil {
		return SessionHostCommand{}, e
	}
	c, ok := st.SessionHostCommands[id]
	if !ok {
		return c, fail("missing_coverage", "Session host request absent")
	}
	return c, nil
}
func (s *Service) sessionHostBasis(st snapshot, q SessionHostRequest) (Run, error) {
	i, ok := st.Sessions[q.InstanceID]
	if !ok || i.Revision != q.ExpectedRevision {
		return Run{}, fail("stale_revision", "Session revision differs")
	}
	r, ok := st.Runs[q.RunID]
	if !ok || r.Revision != q.ExpectedRunRevision {
		return r, fail("stale_revision", "run revision differs")
	}
	if !i.Registered || i.Rebind != nil || i.Advance != nil || len(i.Work) == 0 || i.Work[0].Request.OperationID != q.RunID || sessionLatestEffects(i).Status != "committed" || i.Disposition != nil {
		return r, fail("conflict", "current committed unfinished Session assignment required")
	}
	if e := s.sessionCurrent(i); e != nil {
		return r, e
	}
	if e := s.CheckSessionLeader(i.ID, q.LeaderSkill); e != nil {
		return r, e
	}
	if e := s.sessionWorkClaims(st, i); e != nil {
		return r, e
	}
	switch q.Action {
	case "launch":
		if r.State != "ready" || len(r.Attempts) != 0 || len(r.Inputs) != 0 {
			return r, fail("conflict", "fresh ready assignment required")
		}
	case "continue":
		if r.State != "running" || len(r.Attempts) == 0 || r.Attempts[len(r.Attempts)-1].Status != "completed" {
			return r, fail("conflict", "completed waiting attempt required; controller also checks current replies")
		}
	default:
		if r.State != "running" || len(r.Attempts) == 0 || r.Attempts[len(r.Attempts)-1].Status != "running" {
			return r, fail("conflict", "mapped running attempt required")
		}
	}
	for _, a := range r.Attempts {
		if a.Gap || a.Status == "unknown" {
			return r, fail("reconciliation_required", "unresolved host observation")
		}
	}
	for _, v := range r.Inputs {
		if v.State != "acknowledged" {
			return r, fail("reconciliation_required", "unresolved outbound host input")
		}
	}
	return r, nil
}
func (s *Service) RequestSessionHost(c Caller, q SessionHostRequest) (SessionHostCommand, error) {
	if !c.controller {
		return SessionHostCommand{}, fail("unauthorized", "trusted controller required")
	}
	if e := validSessionHostRequest(q); e != nil {
		return SessionHostCommand{}, e
	}
	_, tx, unlock, e := s.lockRecordWriter()
	if e != nil {
		return SessionHostCommand{}, e
	}
	defer unlock()
	st, head, e := s.load()
	if e != nil {
		return SessionHostCommand{}, e
	}
	if e = s.authorize(c, Schema, q.ProjectRef, q.CheckoutID, q.OperationID, st); e != nil {
		return SessionHostCommand{}, e
	}
	if old, ok := st.SessionHostCommands[q.OperationID]; ok {
		if !reflect.DeepEqual(old.Request, q) {
			return old, fail("idempotency_conflict", "host request identity reused with different content")
		}
		return old, nil
	}
	if _, used, _ := retry(st, q.OperationID, ""); used {
		return SessionHostCommand{}, fail("idempotency_conflict", "operation belongs to another workflow")
	}
	if e = tx.Check(""); e != nil {
		return SessionHostCommand{}, e
	}
	if e = sessionHostHold(st, q.InstanceID); e != nil {
		return SessionHostCommand{}, e
	}
	r, e := s.sessionHostBasis(st, q)
	if e != nil {
		return SessionHostCommand{}, e
	}
	if len(st.SessionHostCommands) >= 128 {
		return SessionHostCommand{}, fail("limit", "Session host request capacity reached")
	}
	at := time.Now().UTC().Format(time.RFC3339Nano)
	out := SessionHostCommand{Request: q, InputRevision: r.InputRevision, State: "queued", CreatedAt: at, UpdatedAt: at, Detail: "Awaiting trusted controller; no host effect dispatched"}
	if len(r.Attempts) > 0 {
		out.AttemptID = r.Attempts[len(r.Attempts)-1].ID
	}
	if st.SessionHostCommands == nil {
		st.SessionHostCommands = map[string]SessionHostCommand{}
	}
	st.SessionHostCommands[q.OperationID] = out
	st.Schema = sessionHostStateSchema
	if e = s.save(st, head); e != nil {
		return SessionHostCommand{}, e
	}
	return out, nil
}

// The executing controller holds this lease across claim, host effect and settlement.
func (s *Service) SessionHostCommandLease(c Caller, id string) (func(), error) {
	if !c.controller || !uuidPattern.MatchString(id) {
		return nil, fail("unauthorized", "trusted exact host request required")
	}
	return lock(filepath.Join(s.state, "session-host-controller", id))
}
func (s *Service) ClaimSessionHost(c Caller, id string) (SessionHostCommand, error) {
	if !c.controller {
		return SessionHostCommand{}, fail("unauthorized", "trusted controller required")
	}
	_, tx, unlock, e := s.lockRecordWriter()
	if e != nil {
		return SessionHostCommand{}, e
	}
	defer unlock()
	st, head, e := s.load()
	if e != nil {
		return SessionHostCommand{}, e
	}
	cmd, ok := st.SessionHostCommands[id]
	if !ok {
		return cmd, fail("missing_coverage", "host request absent")
	}
	if cmd.State != "queued" {
		return cmd, fail("reconciliation_required", "host request already claimed; inspect/reconcile, never resend")
	}
	if e = tx.Check(""); e != nil {
		return cmd, e
	}
	if _, e = s.sessionHostBasis(st, cmd.Request); e != nil {
		return cmd, e
	}
	cmd.State = "claimed"
	cmd.UpdatedAt = time.Now().UTC().Format(time.RFC3339Nano)
	cmd.Detail = "Controller claimed request; inspect mapped host input before retry"
	st.SessionHostCommands[id] = cmd
	if e = s.save(st, head); e != nil {
		return SessionHostCommand{}, e
	}
	return cmd, nil
}

// Settlement derives acknowledgement from the existing exact host input, never
// from a caller-supplied success flag. Claimed work without proof stays uncertain.
func (s *Service) SettleSessionHost(c Caller, id string) (SessionHostCommand, error) {
	if !c.controller {
		return SessionHostCommand{}, fail("unauthorized", "trusted controller required")
	}
	unlock, e := lock(s.state)
	if e != nil {
		return SessionHostCommand{}, e
	}
	defer unlock()
	st, head, e := s.load()
	if e != nil {
		return SessionHostCommand{}, e
	}
	cmd, ok := st.SessionHostCommands[id]
	if !ok {
		return cmd, fail("missing_coverage", "host request absent")
	}
	if cmd.State == "acknowledged" || cmd.State == "canceled" {
		return cmd, nil
	}
	if cmd.State != "claimed" && cmd.State != "uncertain" {
		return cmd, fail("conflict", "request has not been claimed")
	}
	cmd.State = "uncertain"
	cmd.Detail = "Host delivery unresolved; recover exact run/input before further dispatch"
	for _, v := range st.Runs[cmd.Request.RunID].Inputs {
		if v.ID == id && v.State == "acknowledged" {
			cmd.State = "acknowledged"
			cmd.Detail = "Exact host input acknowledged; observed work completion is separate"
		}
	}
	cmd.UpdatedAt = time.Now().UTC().Format(time.RFC3339Nano)
	st.SessionHostCommands[id] = cmd
	if e = s.save(st, head); e != nil {
		return SessionHostCommand{}, e
	}
	return cmd, nil
}

// Cancel only a queued request, or a claimed request with no durable host input.
// Acquires the request lease; no external effect precedes OwnerInput.
func (s *Service) CancelSessionHost(c Caller, id string) (SessionHostCommand, error) {
	if !c.controller {
		return SessionHostCommand{}, fail("unauthorized", "trusted controller required")
	}
	release, e := s.SessionHostCommandLease(c, id)
	if e != nil {
		return SessionHostCommand{}, e
	}
	defer release()
	unlock, e := lock(s.state)
	if e != nil {
		return SessionHostCommand{}, e
	}
	defer unlock()
	st, head, e := s.load()
	if e != nil {
		return SessionHostCommand{}, e
	}
	cmd, ok := st.SessionHostCommands[id]
	if !ok {
		return cmd, fail("missing_coverage", "host request absent")
	}
	if cmd.State == "canceled" {
		return cmd, nil
	}
	if !sessionHostPending(cmd.State) {
		return cmd, fail("conflict", "only undispatched requests can be canceled")
	}
	for _, v := range st.Runs[cmd.Request.RunID].Inputs {
		if v.ID == id {
			return cmd, fail("reconciliation_required", "durable host input exists; reconcile that input instead of canceling")
		}
	}
	cmd.State = "canceled"
	cmd.UpdatedAt = time.Now().UTC().Format(time.RFC3339Nano)
	cmd.Detail = "Canceled with no durable host input or host dispatch"
	st.SessionHostCommands[id] = cmd
	if e = s.save(st, head); e != nil {
		return SessionHostCommand{}, e
	}
	return cmd, nil
}
func validateSessionHostCommands(st snapshot) error {
	if st.SessionHostCommands == nil {
		return nil
	}
	if st.Schema != sessionHostStateSchema || len(st.SessionHostCommands) > 128 {
		return fail("incompatible_profile", "Session host requests require bounded profile 0.23")
	}
	pending := map[string]bool{}
	for id, c := range st.SessionHostCommands {
		q := c.Request
		_, a := st.Receipts[id]
		_, b := st.LeaderReceipts[id]
		_, d := st.DesignDecisionReceipts[id]
		_, f := st.PreparationReceipts[id]
		if a || b || d || f {
			return fail("reconciliation_required", "host request identity overlaps another operation")
		}
		i, ok := st.Sessions[q.InstanceID]
		r, runOK := st.Runs[q.RunID]
		if !ok || !runOK || id != q.OperationID || validSessionHostRequest(q) != nil || q.ProjectRef != st.Project || q.CheckoutID != st.CheckoutID || q.ExpectedRevision > i.Revision || q.ExpectedRunRevision > r.Revision || c.InputRevision != r.InputRevision || !slices.Contains([]string{"queued", "claimed", "acknowledged", "uncertain", "canceled"}, c.State) || !validLeaderText(c.Detail, 4096) {
			return fail("reconciliation_required", "invalid retained Session host request")
		}
		bound := false
		for _, old := range sessionExecutions(snapshot{Sessions: map[string]SessionInstance{i.ID: i}}) {
			if len(old.Work) > 0 && old.Work[0].Request.OperationID == q.RunID {
				bound = true
			}
		}
		if !bound {
			return fail("reconciliation_required", "host request outside retained Session run")
		}
		created, e := time.Parse(time.RFC3339Nano, c.CreatedAt)
		if e != nil {
			return fail("reconciliation_required", "invalid host request time")
		}
		updated, e := time.Parse(time.RFC3339Nano, c.UpdatedAt)
		if e != nil || updated.Before(created) {
			return fail("reconciliation_required", "invalid host request update time")
		}
		if c.AttemptID != "" && !slices.ContainsFunc(r.Attempts, func(a Attempt) bool { return a.ID == c.AttemptID }) {
			return fail("reconciliation_required", "host request attempt absent")
		}
		found := false
		for _, v := range r.Inputs {
			if v.ID == id {
				found = true
				kind := q.Action
				if kind == "continue" {
					kind = "launch"
				}
				if v.Kind != kind || (c.State == "acknowledged" && v.State != "acknowledged") {
					return fail("reconciliation_required", "host request input proof differs")
				}
			}
		}
		if c.State == "acknowledged" && !found || (c.State == "canceled" || c.State == "queued") && found {
			return fail("reconciliation_required", "host request delivery proof absent or contradictory")
		}
		if sessionHostPending(c.State) {
			if pending[q.InstanceID] {
				return fail("reconciliation_required", "multiple pending host requests for one Session")
			}
			pending[q.InstanceID] = true
		}
	}
	return nil
}
