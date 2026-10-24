package governance

import (
	"encoding/json"
	"path/filepath"
	"slices"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/Hans-Einar/SDP/SDPTool/codexcompat"
)

const LeaderSchema = "sdp-governance-leader/1"
const leaderStateSchema = "sdp-project-governance-state/0.9"

type LeaderMessageInput struct {
	Actor  string `json:"actor"`
	Source string `json:"source"`
	Text   string `json:"text"`
}
type LeaderMessage struct {
	ID         string             `json:"id"`
	Input      LeaderMessageInput `json:"input"`
	RecordedAt string             `json:"recordedAt"`
}
type LeaderConfig struct {
	Coordinator  string       `json:"coordinator"`
	HostVersion  string       `json:"hostVersion"`
	MCPBinary    string       `json:"mcpBinary"`
	Instructions string       `json:"instructions"`
	Skills       []SkillInput `json:"skills"`
}
type LeaderDispatch struct {
	ID         string       `json:"id"`
	State      string       `json:"state"`
	Config     LeaderConfig `json:"config"`
	MessageIDs []string     `json:"messageIds"`
	NoticeKeys []string     `json:"noticeKeys"`
	InputHead  string       `json:"inputHead"`
	Payload    string       `json:"payload"`
	ThreadID   string       `json:"threadId"`
	TurnID     string       `json:"turnId"`
	Answer     string       `json:"answer"`
	UpdatedAt  string       `json:"updatedAt"`
}
type LeaderState struct {
	Revision   uint64           `json:"revision"`
	Messages   []LeaderMessage  `json:"messages"`
	Dispatches []LeaderDispatch `json:"dispatches"`
}
type LeaderView struct {
	Schema     string      `json:"schema"`
	ProjectRef string      `json:"projectRef"`
	CheckoutID string      `json:"checkoutId"`
	Leader     LeaderState `json:"leader"`
}
type LeaderRequest struct {
	Schema           string              `json:"schema"`
	ProjectRef       string              `json:"projectRef"`
	CheckoutID       string              `json:"checkoutId"`
	OperationID      string              `json:"operationId"`
	ExpectedRevision uint64              `json:"expectedRevision"`
	Action           string              `json:"action"`
	Message          *LeaderMessageInput `json:"message,omitempty"`
	Config           *LeaderConfig       `json:"config,omitempty"`
	DispatchID       string              `json:"dispatchId,omitempty"`
	ThreadID         string              `json:"threadId,omitempty"`
	TurnID           string              `json:"turnId,omitempty"`
	Status           string              `json:"status,omitempty"`
	Answer           string              `json:"answer,omitempty"`
}
type LeaderResult struct {
	Schema      string          `json:"schema"`
	OperationID string          `json:"operationId"`
	Action      string          `json:"action"`
	Revision    uint64          `json:"revision"`
	Message     *LeaderMessage  `json:"message,omitempty"`
	Dispatch    *LeaderDispatch `json:"dispatch,omitempty"`
}
type leaderReceipt struct {
	Digest string       `json:"digest"`
	Result LeaderResult `json:"result"`
}
type leaderNotice struct {
	Key  string           `json:"key"`
	Item CoordinationItem `json:"item"`
}

func LeaderActive(state string) bool {
	return slices.Contains([]string{"prepared", "opening", "bound", "sending", "running"}, state)
}
func validLeaderText(s string, limit int) bool {
	return utf8.ValidString(s) && !strings.ContainsRune(s, 0) && boundedAttentionText(s, limit)
}
func validLeaderConfig(c LeaderConfig) bool {
	return validLeaderText(c.Coordinator, 256) && codexcompat.Supported(c.HostVersion) && filepath.IsAbs(c.MCPBinary) && filepath.Clean(c.MCPBinary) == c.MCPBinary && validLeaderText(c.MCPBinary, 4096) && validLeaderText(c.Instructions, 4096) && len(c.Skills) > 0 && ValidSkillInputs(c.Skills)
}
func (s *Service) LeaderRead() (LeaderView, error) {
	state, _, e := s.load()
	out := LeaderView{Schema: LeaderSchema, ProjectRef: s.project, CheckoutID: state.CheckoutID, Leader: LeaderState{Messages: []LeaderMessage{}, Dispatches: []LeaderDispatch{}}}
	if state.Leader != nil {
		out.Leader = *state.Leader
	}
	return out, e
}

// LeaderLease serializes dispatch effects, without excluding worker controllers or owner enqueue.
func (s *Service) LeaderLease(c Caller) (func(), error) {
	if !c.controller {
		return nil, fail("unauthorized", "trusted controller required")
	}
	return lock(s.state + "/leader-controller")
}
func (s *Service) leaderNotices() ([]leaderNotice, error) {
	out := []leaderNotice{}
	cursor := ""
	revision := ""
	for {
		v, e := s.Coordination(CoordinationRequest{Cursor: cursor, Limit: 50})
		if e != nil {
			return nil, e
		}
		if revision != "" && v.Revision != revision {
			return nil, fail("stale_revision", "coordination changed during leader selection")
		}
		revision = v.Revision
		for _, item := range v.Items {
			normalized := item
			if item.Report != nil {
				if item.Report.Acknowledgement != nil {
					continue
				}
				r := *item.Report
				r.CurrentRevision = 0
				normalized.Report = &r
			} else if item.Attention != nil {
				a := *item.Attention
				a.CurrentRevision = 0
				a.Consumption = ""
				normalized.Attention = &a
			} else {
				continue
			}
			out = append(out, leaderNotice{Key: requestDigest(normalized), Item: item})
		}
		if v.NextCursor == "" {
			return out, nil
		}
		cursor = v.NextCursor
	}
}
func leaderInputHead(l LeaderState) string {
	if len(l.Messages) == 0 {
		return ""
	}
	return l.Messages[len(l.Messages)-1].ID
}
func leaderDelivered(l LeaderState) (map[string]bool, map[string]bool) {
	messages, notices := map[string]bool{}, map[string]bool{}
	for _, d := range l.Dispatches {
		if d.State == "superseded" {
			continue
		}
		for _, id := range d.MessageIDs {
			messages[id] = true
		}
		for _, key := range d.NoticeKeys {
			notices[key] = true
		}
	}
	return messages, notices
}
func (s *Service) prepareLeader(l LeaderState, id string, c LeaderConfig) (LeaderDispatch, error) {
	d := LeaderDispatch{ID: id, State: "prepared", Config: c, InputHead: leaderInputHead(l), MessageIDs: []string{}, NoticeKeys: []string{}}
	if e := CheckSkillSnapshots(s.root, c.Skills); e != nil {
		return d, e
	}
	messages, notices := leaderDelivered(l)
	payload := struct {
		Schema  string          `json:"schema"`
		Project string          `json:"project"`
		Owner   []LeaderMessage `json:"ownerMessages"`
		Notices []leaderNotice  `json:"notices"`
	}{"sdp-leader-context/1", s.project, []LeaderMessage{}, []leaderNotice{}}
	for _, m := range l.Messages {
		if !messages[m.ID] && len(payload.Owner) < 8 {
			payload.Owner = append(payload.Owner, m)
			b, _ := json.Marshal(payload)
			if len(b) > 64<<10 {
				payload.Owner = payload.Owner[:len(payload.Owner)-1]
				break
			}
			d.MessageIDs = append(d.MessageIDs, m.ID)
		}
	}
	available, e := s.leaderNotices()
	if e != nil {
		return d, e
	}
	pendingNotice := false
	for _, n := range available {
		if !notices[n.Key] {
			pendingNotice = true
		}
		if notices[n.Key] || len(payload.Notices) >= 8 {
			continue
		}
		payload.Notices = append(payload.Notices, n)
		b, _ := json.Marshal(payload)
		if len(b) > 64<<10 {
			payload.Notices = payload.Notices[:len(payload.Notices)-1]
			continue
		}
		d.NoticeKeys = append(d.NoticeKeys, n.Key)
	}
	if len(d.MessageIDs)+len(d.NoticeKeys) == 0 {
		if pendingNotice {
			return d, fail("limit", "pending leader context exceeds payload bound")
		}
		return d, fail("no_work", "no pending owner input or changed notices")
	}
	b, e := json.Marshal(payload)
	if e != nil {
		return d, e
	}
	if len(b) > 64<<10 {
		return d, fail("limit", "leader context exceeds 64 KiB")
	}
	d.Payload = string(b)
	return d, nil
}
func (s *Service) leaderBasis(l LeaderState, d LeaderDispatch) error {
	if leaderInputHead(l) != d.InputHead {
		return fail("stale_revision", "new owner input precedes unsent leader work")
	}
	if e := CheckSkillSnapshots(s.root, d.Config.Skills); e != nil {
		return e
	}
	notices, e := s.leaderNotices()
	if e != nil {
		return e
	}
	current := map[string]bool{}
	for _, n := range notices {
		current[n.Key] = true
	}
	for _, key := range d.NoticeKeys {
		if !current[key] {
			return fail("stale_revision", "selected leader notice changed")
		}
	}
	return nil
}
func (s *Service) LeaderWrite(c Caller, q LeaderRequest) (LeaderResult, error) {
	var out LeaderResult
	if !c.controller {
		return out, fail("unauthorized", "leader writes require trusted controller")
	}
	if q.Schema != LeaderSchema {
		return out, fail("incompatible_profile", "unsupported leader request")
	}
	if !slices.Contains([]string{"message", "prepare", "open", "bind", "send", "observe", "supersede"}, q.Action) {
		return out, fail("invalid_request", "unknown leader action")
	}
	if (q.Action == "message") != (q.Message != nil) || (q.Action == "prepare") != (q.Config != nil) {
		return out, fail("invalid_request", "invalid action fields")
	}
	if q.Action == "message" || q.Action == "prepare" {
		if q.DispatchID != "" {
			return out, fail("invalid_request", "unexpected dispatch identity")
		}
	} else if !uuidPattern.MatchString(q.DispatchID) {
		return out, fail("invalid_request", "exact dispatch required")
	}
	if q.Action != "bind" && q.Action != "observe" && q.ThreadID != "" || q.Action != "observe" && (q.TurnID != "" || q.Status != "" || q.Answer != "") {
		return out, fail("invalid_request", "unexpected host observation fields")
	}
	if q.Message != nil && (!validLeaderText(q.Message.Actor, 256) || !validLeaderText(q.Message.Source, 1024) || !validLeaderText(q.Message.Text, 4096)) {
		return out, fail("invalid_request", "bounded attributed message required")
	}
	if q.Config != nil && !validLeaderConfig(*q.Config) {
		return out, fail("invalid_request", "explicit host/coordinator/skill configuration required")
	}
	unlock, e := lock(s.state)
	if e != nil {
		return out, e
	}
	defer unlock()
	state, head, e := s.load()
	if e != nil {
		return out, e
	}
	if e = s.authorize(c, Schema, q.ProjectRef, q.CheckoutID, q.OperationID, state); e != nil {
		return out, e
	}
	digest := requestDigest(q)
	if old, ok := state.LeaderReceipts[q.OperationID]; ok {
		if old.Digest != digest {
			return out, fail("idempotency_conflict", "leader operation content changed")
		}
		return old.Result, nil
	}
	if _, ok, e := retry(state, q.OperationID, digest); ok {
		if e != nil {
			return out, e
		}
		return out, fail("idempotency_conflict", "operation belongs to another service")
	}
	l := LeaderState{Messages: []LeaderMessage{}, Dispatches: []LeaderDispatch{}}
	if state.Leader != nil {
		l = *state.Leader
	}
	if (q.Action != "message" && q.ExpectedRevision != l.Revision) || (q.Action == "message" && q.ExpectedRevision != 0) {
		return out, fail("stale_revision", "leader revision changed")
	}
	if len(state.LeaderReceipts) >= 256 || l.Revision == ^uint64(0) {
		return out, fail("limit", "leader receipt/revision capacity reached")
	}
	out = LeaderResult{Schema: LeaderSchema, OperationID: q.OperationID, Action: q.Action, Revision: l.Revision + 1}
	timestamp := time.Now().UTC().Format(time.RFC3339Nano)
	switch q.Action {
	case "message":
		if len(l.Messages) >= 64 {
			return out, fail("limit", "leader message capacity reached")
		}
		m := LeaderMessage{ID: q.OperationID, Input: *q.Message, RecordedAt: timestamp}
		l.Messages = append(l.Messages, m)
		out.Message = &m
	case "prepare":
		for _, d := range l.Dispatches {
			if LeaderActive(d.State) {
				return out, fail("conflict", "a leader dispatch is already active")
			}
		}
		if len(l.Dispatches) >= 16 {
			return out, fail("limit", "leader dispatch capacity reached")
		}
		d, e := s.prepareLeader(l, q.OperationID, *q.Config)
		if e != nil {
			return out, e
		}
		d.UpdatedAt = timestamp
		l.Dispatches = append(l.Dispatches, d)
		out.Dispatch = &d
	default:
		i := slices.IndexFunc(l.Dispatches, func(d LeaderDispatch) bool { return d.ID == q.DispatchID })
		if i < 0 {
			return out, fail("missing_coverage", "leader dispatch absent")
		}
		d := l.Dispatches[i]
		switch q.Action {
		case "open":
			if d.State != "prepared" {
				return out, fail("conflict", "leader thread creation already claimed")
			}
			if e = s.leaderBasis(l, d); e != nil {
				return out, e
			}
			d.State = "opening"
		case "bind":
			if d.State != "opening" || !validLeaderText(q.ThreadID, 256) {
				return out, fail("conflict", "binding needs claimed thread creation and observed identity")
			}
			for _, existing := range l.Dispatches {
				if existing.ID != d.ID && existing.ThreadID == q.ThreadID {
					return out, fail("reconciliation_required", "leader thread already belongs to another dispatch")
				}
			}
			d.ThreadID = q.ThreadID
			d.State = "bound"
		case "send":
			if d.State != "bound" {
				return out, fail("conflict", "leader turn send already claimed or thread missing")
			}
			if e = s.leaderBasis(l, d); e != nil {
				return out, e
			}
			d.State = "sending"
		case "supersede":
			if d.State != "prepared" && d.State != "bound" {
				return out, fail("reconciliation_required", "cannot supersede possibly submitted leader work")
			}
			d.State = "superseded"
		case "observe":
			if (d.State != "sending" && d.State != "running") || q.ThreadID != d.ThreadID || !validLeaderText(q.TurnID, 256) || (d.TurnID != "" && q.TurnID != d.TurnID) || !slices.Contains([]string{"running", "completed", "failed", "interrupted"}, q.Status) || !utf8.ValidString(q.Answer) || len(q.Answer) > 64<<10 || (q.Status == "running" && q.Answer != "") {
				return out, fail("reconciliation_required", "leader observation does not match pending dispatch")
			}
			d.TurnID = q.TurnID
			d.State = q.Status
			d.Answer = q.Answer
		}
		d.UpdatedAt = timestamp
		l.Dispatches[i] = d
		out.Dispatch = &d
	}
	l.Revision++
	state.Leader = &l
	state.Schema = leaderStateSchema
	if state.LeaderReceipts == nil {
		state.LeaderReceipts = map[string]leaderReceipt{}
	}
	state.LeaderReceipts[q.OperationID] = leaderReceipt{Digest: digest, Result: out}
	if e = s.save(state, head); e != nil {
		return LeaderResult{}, e
	}
	return out, nil
}
func validateLeaderState(state snapshot) error {
	if state.Schema != leaderStateSchema {
		if state.Leader != nil || len(state.LeaderReceipts) > 0 {
			return fail("incompatible_profile", "leader records require operational profile 0.9")
		}
		return nil
	}
	if state.Leader == nil || state.Leader.Revision == 0 || len(state.Leader.Messages) > 64 || len(state.Leader.Dispatches) > 16 || len(state.LeaderReceipts) == 0 || len(state.LeaderReceipts) > 256 {
		return fail("reconciliation_required", "invalid leader collections")
	}
	l := state.Leader
	ids := map[string]bool{}
	messageIDs := map[string]bool{}
	active := 0
	for _, m := range l.Messages {
		r, ok := state.LeaderReceipts[m.ID]
		_, e := time.Parse(time.RFC3339Nano, m.RecordedAt)
		if !ok || !uuidPattern.MatchString(m.ID) || ids[m.ID] || e != nil || r.Result.Message == nil || requestDigest(*r.Result.Message) != requestDigest(m) || !validLeaderText(m.Input.Actor, 256) || !validLeaderText(m.Input.Source, 1024) || !validLeaderText(m.Input.Text, 4096) {
			return fail("reconciliation_required", "invalid retained owner message")
		}
		ids[m.ID] = true
		messageIDs[m.ID] = true
	}
	for _, d := range l.Dispatches {
		receipt, ok := state.LeaderReceipts[d.ID]
		_, e := time.Parse(time.RFC3339Nano, d.UpdatedAt)
		if !ok || receipt.Result.Action != "prepare" || !uuidPattern.MatchString(d.ID) || ids[d.ID] || !validLeaderConfig(d.Config) || !validLeaderText(d.Payload, 64<<10) || len(d.Answer) > 64<<10 || e != nil || !slices.Contains([]string{"prepared", "opening", "bound", "sending", "running", "completed", "failed", "interrupted", "superseded"}, d.State) {
			return fail("reconciliation_required", "invalid leader dispatch")
		}
		original := receipt.Result.Dispatch
		if original == nil || original.State != "prepared" || original.ThreadID != "" || original.TurnID != "" || original.Answer != "" || !sameLeaderBasis(*original, d) {
			return fail("reconciliation_required", "leader preparation changed")
		}
		latest := receipt.Result
		for _, r := range state.LeaderReceipts {
			if r.Result.Dispatch != nil && r.Result.Dispatch.ID == d.ID && r.Result.Revision > latest.Revision {
				latest = r.Result
			}
		}
		if requestDigest(latest.Dispatch) != requestDigest(d) {
			return fail("reconciliation_required", "leader state differs from latest receipt")
		}
		if len(d.MessageIDs)+len(d.NoticeKeys) == 0 || len(d.MessageIDs) > 8 || len(d.NoticeKeys) > 8 {
			return fail("reconciliation_required", "invalid leader inputs")
		}
		used := map[string]bool{}
		for _, id := range d.MessageIDs {
			if !messageIDs[id] || used[id] {
				return fail("reconciliation_required", "invalid leader message references")
			}
			used[id] = true
		}
		for _, key := range d.NoticeKeys {
			if !digestPattern.MatchString(key) || used[key] {
				return fail("reconciliation_required", "invalid leader notice references")
			}
			used[key] = true
		}
		if d.InputHead != "" && !messageIDs[d.InputHead] {
			return fail("reconciliation_required", "leader input head absent")
		}
		if slices.Contains([]string{"bound", "sending", "running", "completed", "failed", "interrupted"}, d.State) && strings.TrimSpace(d.ThreadID) == "" {
			return fail("reconciliation_required", "mapped leader thread absent")
		}
		if slices.Contains([]string{"running", "completed", "failed", "interrupted"}, d.State) && strings.TrimSpace(d.TurnID) == "" {
			return fail("reconciliation_required", "mapped leader turn absent")
		}
		if LeaderActive(d.State) {
			active++
		}
		ids[d.ID] = true
	}
	if active > 1 {
		return fail("reconciliation_required", "concurrent leader dispatches")
	}
	if uint64(len(state.LeaderReceipts)) != l.Revision {
		return fail("reconciliation_required", "leader revision coverage differs")
	}
	revisions := map[uint64]bool{}
	for id, r := range state.LeaderReceipts {
		_, a := state.Receipts[id]
		_, b := state.PreparationReceipts[id]
		_, c := state.DesignDecisionReceipts[id]
		if a || b || c || revisions[r.Result.Revision] || !uuidPattern.MatchString(id) || r.Result.OperationID != id || r.Result.Schema != LeaderSchema || r.Result.Revision == 0 || r.Result.Revision > l.Revision || !digestPattern.MatchString(r.Digest) {
			return fail("reconciliation_required", "invalid leader receipt")
		}
		revisions[r.Result.Revision] = true
	}
	return nil
}

func sameLeaderBasis(a, b LeaderDispatch) bool {
	return a.ID == b.ID && requestDigest(a.Config) == requestDigest(b.Config) && slices.Equal(a.MessageIDs, b.MessageIDs) && slices.Equal(a.NoticeKeys, b.NoticeKeys) && a.InputHead == b.InputHead && a.Payload == b.Payload
}
