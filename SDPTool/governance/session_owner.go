package governance

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/Hans-Einar/SDP/SDPTool/projecthistory"
)

const sessionOwnerStateSchema = "sdp-project-governance-state/0.19"

func sessionDispositionProfile(schema string) bool {
	return schema == sessionDispositionStateSchema || schema == sessionOwnerStateSchema || schema == sessionLifecycleStateSchema
}

type SessionOwnerRequest struct {
	Schema              string `json:"schema"`
	ProjectRef          string `json:"projectRef"`
	CheckoutID          string `json:"checkoutId"`
	OperationID         string `json:"operationId"`
	InstanceID          string `json:"instanceId"`
	ExpectedRevision    uint64 `json:"expectedRevision"`
	ExpectedRunRevision uint64 `json:"expectedRunRevision"`
	Action              string `json:"action"`
	GateID              string `json:"gateId"`
	ResultID            string `json:"resultId"`
	Reason              string `json:"reason"`
}
type SessionOwnerGate struct {
	Request SessionOwnerRequest  `json:"request"`
	Effects SessionRecordEffects `json:"effects"`
	Before  map[string]string    `json:"before"`
}

func ownerGatePath(q SessionOwnerRequest) string {
	return "SDP/Sessions/execution/owner/" + q.OperationID + ".json"
}
func latestOwnerGate(i SessionInstance) *SessionOwnerGate {
	if len(i.OwnerGates) == 0 {
		return nil
	}
	return &i.OwnerGates[len(i.OwnerGates)-1]
}
func sessionOwnerAccepted(i SessionInstance) bool {
	g := latestOwnerGate(i)
	return g != nil && g.Effects.Status == "committed" && g.Request.Action == "accept" && len(i.Work) >= 3 && g.Request.ResultID == dispositionResult(i).Request.OperationID && g.Request.ExpectedRunRevision == dispositionResult(i).Request.ExpectedRunRevision
}
func dispositionPredecessor(i SessionInstance) SessionRecordEffects {
	if g := latestOwnerGate(i); g != nil {
		return g.Effects
	}
	return dispositionResult(i).Effects
}
func validSessionOwner(q SessionOwnerRequest) error {
	if q.Schema != SessionSchema || !uuidPattern.MatchString(q.OperationID) || !uuidPattern.MatchString(q.InstanceID) || !uuidPattern.MatchString(q.GateID) || !uuidPattern.MatchString(q.ResultID) || q.ExpectedRevision == 0 || q.ExpectedRunRevision == 0 || !validLeaderText(q.Reason, 2048) || !slices.Contains([]string{"request", "accept", "changes-requested"}, q.Action) {
		return fail("invalid_request", "exact owner gate identity, revision, action and reason required")
	}
	if (q.Action == "request") != (q.GateID == q.OperationID) {
		return fail("invalid_request", "request owns gate identity; decisions reference an existing gate")
	}
	return nil
}
func ownerGateEligibility(i SessionInstance, q SessionOwnerRequest) error {
	if i.Selection == nil || i.Disposition != nil || len(i.Work) < 3 || len(i.Work) > 4 || selectedSessionStep(i).OwnerAcceptance != "required" {
		return fail("missing_coverage", "verified selected step requiring owner acceptance required")
	}
	w := dispositionResult(i)
	if w.Effects.Status != "committed" || w.Request.OperationID != q.ResultID || w.Request.ExpectedRunRevision != q.ExpectedRunRevision || (selectedSessionStep(i).Review == "independent" && w.Kind != "review") {
		return fail("missing_coverage", "exact committed required review/test result required")
	}
	g := latestOwnerGate(i)
	if q.Action == "request" {
		if g != nil && (g.Effects.Status != "committed" || g.Request.Action == "request" || g.Request.ResultID == q.ResultID) {
			return fail("conflict", "previous owner gate requires a new reviewed rework result")
		}
	} else if g == nil || g.Effects.Status != "committed" || g.Request.Action != "request" || g.Request.GateID != q.GateID || g.Request.ResultID != q.ResultID || g.Request.ExpectedRunRevision != q.ExpectedRunRevision {
		return fail("conflict", "exact committed waiting owner gate required")
	}
	return nil
}
func (s *Service) ownerGateBasis(st snapshot, i SessionInstance, g SessionOwnerGate) error {
	if e := ownerGateEligibility(i, g.Request); e != nil {
		return e
	}
	r := st.Runs[i.Work[0].Request.OperationID]
	if r.Revision != g.Request.ExpectedRunRevision {
		return fail("stale_revision", "owner gate run changed")
	}
	return s.sessionResultSnapshot(st, i, dispositionResult(i))
}
func ownerGateArtifact(i SessionInstance, g SessionOwnerGate) string {
	actor := "local-owner"
	if g.Request.Action == "request" {
		actor = "sdp-controller"
	}
	b, _ := json.MarshalIndent(struct {
		Schema        string              `json:"schema"`
		Request       SessionOwnerRequest `json:"request"`
		Actor         string              `json:"actor"`
		CreatedAt     string              `json:"createdAt"`
		BindingDigest string              `json:"bindingDigest"`
		RunID         string              `json:"runId"`
		Candidate     string              `json:"candidate"`
	}{"sdp-session-owner/0.1", g.Request, actor, g.Effects.CreatedAt, i.Context.BindingDigest, i.Work[0].Request.OperationID, i.Work[2].Evidence.Request.Candidate.Digest}, "", "  ")
	return string(b) + "\n"
}
func ownerGateJournal(i SessionInstance, g SessionOwnerGate, p string) string {
	rel, _ := filepath.Rel(filepath.Dir(p), ownerGatePath(g.Request))
	return "\n## PG owner gate " + g.Request.Action + " " + g.Request.OperationID + "\n\nStep " + selectedSessionStep(i).ID + "; exact result " + g.Request.ResultID + ".\n[Owner gate record](" + filepath.ToSlash(rel) + ").\n\n" + publicationProse(g.Request.Reason)
}

// Pure transformation: validation reuses the same selected-record plan as preparation.
func planOwnerGate(i SessionInstance, g SessionOwnerGate) (map[string]string, error) {
	files := map[string]string{}
	for p, b := range g.Before {
		files[p] = b
	}
	events, e := publicationEvents([]byte(files[managementLedger]))
	if e != nil {
		return nil, e
	}
	latest := map[string]map[string]any{}
	for _, ev := range events {
		latest[fmt.Sprint(ev["subjectId"])] = ev
	}
	_, cards, e := dispositionInputs(i)
	if e != nil {
		return nil, e
	}
	subjects := map[string]string{i.Context.Binding.Plan.ID: i.Context.Binding.Plan.Path}
	for id, p := range cards {
		subjects[id] = p
	}
	ids := []string{}
	for id := range subjects {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		p := subjects[id]
		if _, ok := files[p]; !ok {
			return nil, fail("missing_coverage", "owner gate selected preimage absent")
		}
		if strings.HasPrefix(id, "KB-") {
			from, to := "in-progress", "gate-review"
			if g.Request.Action != "request" {
				from, to = to, from
			}
			m, e := metadata([]byte(files[p]))
			if e != nil || m["CardState"] != from {
				return nil, fail("conflict", "owner gate card state differs")
			}
			b, e := replaceMetadata([]byte(files[p]), "CardState", to)
			if e != nil {
				return nil, e
			}
			files[p] = string(b)
		}
		files[p] = appendRecord([]byte(files[p]), ownerGateJournal(i, g, p))
		old := latest[id]
		if old == nil {
			return nil, fail("missing_coverage", "owner gate history absent")
		}
		prev, ok := old["payload"].(map[string]any)
		if !ok {
			return nil, fail("reconciliation_required", "owner history malformed")
		}
		payload := map[string]any{}
		for k, v := range prev {
			payload[k] = v
		}
		delete(payload, "lineage")
		payload["previousEventId"] = old["eventId"]
		payload["from"] = prev["to"]
		payload["fromPath"] = prev["toPath"]
		payload["reason"] = "Owner gate " + g.Request.Action + " " + g.Request.OperationID + "; exact result " + g.Request.ResultID
		prefix, kind := "EVT-PM-"+strings.Split(id, "-")[1]+"-", "x-management:updated"
		if strings.HasPrefix(id, "KB-") {
			prefix = "EVT-KB-" + strings.Split(id, "-")[1] + "-"
			kind = "x-kanban:reviewed"
		}
		eid, e := nextPublicationEvent(events, prefix)
		if e != nil {
			return nil, e
		}
		actor := "local-owner"
		if g.Request.Action == "request" {
			actor = "sdp-controller"
		}
		ev := map[string]any{"schemaVersion": "1.0", "eventId": eid, "eventType": kind, "occurredAt": g.Effects.CreatedAt, "actor": actor, "commit": nil, "subjectId": id, "payload": payload}
		raw, _ := json.Marshal(ev)
		files[managementLedger] = appendRecord([]byte(files[managementLedger]), string(raw)+"\n")
		events = append(events, ev)
	}
	sp := i.Context.Binding.Session.Path
	if _, ok := files[sp]; !ok {
		return nil, fail("missing_coverage", "owner gate Session absent")
	}
	files[sp] = appendRecord([]byte(files[sp]), ownerGateJournal(i, g, sp))
	files[ownerGatePath(g.Request)] = ownerGateArtifact(i, g)
	return files, nil
}
func ownerGateCoverage(previous SessionRecordEffects, g SessionOwnerGate) []Source {
	pins := map[string]string{}
	for _, src := range previous.Coverage {
		pins[src.Path] = src.SHA256
	}
	for _, ch := range g.Effects.Changes {
		pins[ch.Path] = ch.After
	}
	names := []string{}
	for p := range pins {
		names = append(names, p)
	}
	sort.Strings(names)
	out := []Source{}
	for _, p := range names {
		out = append(out, Source{p, pins[p]})
	}
	return out
}
func (s *Service) buildOwnerGate(i SessionInstance, g SessionOwnerGate) (SessionOwnerGate, error) {
	files, e := s.publicationFiles()
	if e != nil {
		return g, e
	}
	previous := dispositionPredecessor(i)
	for _, src := range previous.Coverage {
		b, err := readFile(filepath.Join(s.root, src.Path), maxState)
		if os.IsNotExist(err) && src.SHA256 == "absent" {
			continue
		}
		if err != nil {
			return g, err
		}
		if src.Path != managementLedger && hash(b) != src.SHA256 {
			return g, fail("stale_revision", "owner gate input changed")
		}
		files[src.Path] = b
	}
	_, cards, e := dispositionInputs(i)
	if e != nil {
		return g, e
	}
	paths := []string{i.Context.Binding.Session.Path, i.Context.Binding.Plan.Path, managementLedger}
	for _, p := range cards {
		paths = append(paths, p)
	}
	g.Before = map[string]string{}
	for _, p := range paths {
		b, ok := files[p]
		if !ok {
			return g, fail("missing_coverage", "owner gate record absent")
		}
		g.Before[p] = string(b)
	}
	if _, e = os.Lstat(filepath.Join(s.root, ownerGatePath(g.Request))); !os.IsNotExist(e) {
		return g, fail("conflict", "owner record destination occupied")
	}
	after, e := planOwnerGate(i, g)
	if e != nil {
		return g, e
	}
	g.Effects.Changes = dispositionChanges(g.Before, after)
	inputs := map[string]string{}
	for p, b := range files {
		inputs[p] = hash(b)
	}
	for _, src := range previous.Coverage {
		if _, ok := inputs[src.Path]; !ok {
			inputs[src.Path] = src.SHA256
		}
	}
	inputs[ownerGatePath(g.Request)] = "absent"
	names := []string{}
	for p := range inputs {
		names = append(names, p)
	}
	sort.Strings(names)
	for _, p := range names {
		g.Effects.Inputs = append(g.Effects.Inputs, Source{p, inputs[p]})
	}
	for p, b := range after {
		files[p] = []byte(b)
	}
	if e = validatePublicationHistory(files[managementLedger], files); e != nil {
		return g, e
	}
	g.Effects.Coverage = ownerGateCoverage(previous, g)
	g.Effects.LedgerBytes = len(files[managementLedger])
	size := 0
	for _, b := range g.Before {
		size += len(b)
	}
	for _, ch := range g.Effects.Changes {
		size += len(ch.Content)
	}
	if size > 8<<20 || len(g.Effects.Inputs) > 2048 || len(g.Effects.Coverage) > 512 {
		return g, fail("limit", "owner gate capacity exceeded")
	}
	return g, nil
}
func (s *Service) OwnerSession(c Caller, q SessionOwnerRequest) (SessionMutationResult, error) {
	if !c.controller || (q.Action != "request" && !c.owner) {
		return SessionMutationResult{}, fail("unauthorized", "owner decisions require trusted local owner authority")
	}
	if e := validSessionOwner(q); e != nil {
		return SessionMutationResult{}, e
	}
	dir, tx, unlock, e := s.lockRecordWriter()
	if e != nil {
		return SessionMutationResult{}, e
	}
	defer unlock()
	st, head, e := s.load()
	if e != nil {
		return SessionMutationResult{}, e
	}
	if e = s.authorize(c, Schema, q.ProjectRef, q.CheckoutID, q.OperationID, st); e != nil {
		return SessionMutationResult{}, e
	}
	if out, ok, e := sessionRetry(st, q.OperationID, requestDigest(q)); ok {
		return out, e
	}
	if e = tx.Check(""); e != nil {
		return SessionMutationResult{}, e
	}
	if _, e = readFile(filepath.Join(dir, "pending"), 8192); !os.IsNotExist(e) {
		return SessionMutationResult{}, fail("reconciliation_required", "canonical publication pending")
	}
	i, ok := st.Sessions[q.InstanceID]
	if !ok || i.Revision != q.ExpectedRevision {
		return SessionMutationResult{}, fail("stale_revision", "Session revision differs")
	}
	if i.Revision == ^uint64(0) || len(st.Receipts) >= 1024 || len(i.OwnerGates) >= 8 {
		return SessionMutationResult{}, fail("limit", "owner gate capacity reached")
	}
	if e = s.sessionCurrent(i); e != nil {
		return SessionMutationResult{}, e
	}
	g := SessionOwnerGate{Request: q, Effects: SessionRecordEffects{Status: "prepared", CreatedAt: time.Now().UTC().Format(time.RFC3339Nano)}}
	if e = s.ownerGateBasis(st, i, g); e != nil {
		return SessionMutationResult{}, e
	}
	if e = s.sessionWorkClaims(st, i); e != nil {
		return SessionMutationResult{}, e
	}
	g, e = s.buildOwnerGate(i, g)
	if e != nil {
		return SessionMutationResult{}, e
	}
	i.OwnerGates = append(i.OwnerGates, g)
	i.Revision++
	st.Sessions[i.ID] = i
	if st.Schema != sessionLifecycleStateSchema {
		st.Schema = sessionOwnerStateSchema
	}
	out := SessionMutationResult{Schema: SessionSchema, Operation: "session-owner", OperationID: q.OperationID, InstanceID: i.ID, Revision: i.Revision, Registration: "committed", RecordEffects: "prepared"}
	st.Receipts[q.OperationID] = receipt{requestDigest(q), Result{Schema: Schema, Operation: "session.owner", OperationID: q.OperationID, ProjectRef: s.project, CheckoutID: st.CheckoutID, Session: &out}}
	if e = validateSessionOwnerGates(st); e != nil {
		return SessionMutationResult{}, e
	}
	if e = s.save(st, head); e != nil {
		return SessionMutationResult{}, e
	}
	return out, nil
}
func (s *Service) applyOwnerGate(st snapshot, i *SessionInstance, tx *projecthistory.Transaction) error {
	n := len(i.OwnerGates) - 1
	g := &i.OwnerGates[n]
	prior := *i
	prior.OwnerGates = i.OwnerGates[:n]
	if e := s.sessionWorkClaims(st, prior); e != nil {
		return e
	}
	return s.applySessionRecordEffects(&g.Effects, tx, s.sessionPendingOwner(*i), "session-owner", func() error { return s.ownerGateBasis(st, prior, *g) })
}
func validateSessionOwnerGates(st snapshot) error {
	seen := map[string]bool{}
	for _, i := range sessionExecutions(st) {
		if len(i.OwnerGates) > 8 {
			return fail("reconciliation_required", "owner gate capacity exceeded")
		}
		for n, g := range i.OwnerGates {
			q, p := g.Request, g.Effects
			prior := i
			prior.Disposition = nil
			prior.OwnerGates = i.OwnerGates[:n]
			if (st.Schema != sessionOwnerStateSchema && st.Schema != sessionLifecycleStateSchema) || validSessionOwner(q) != nil || ownerGateEligibility(prior, q) != nil || len(i.Work) < 3 || i.Work[2].Evidence == nil || q.ProjectRef != st.Project || q.CheckoutID != st.CheckoutID || q.InstanceID != i.ID || q.ExpectedRevision >= i.Revision || q.ExpectedRevision <= dispositionResult(i).Request.ExpectedRevision || seen[q.OperationID] {
				return fail("reconciliation_required", "invalid owner gate identity/order")
			}
			if n > 0 && q.ExpectedRevision <= i.OwnerGates[n-1].Request.ExpectedRevision+1 {
				return fail("reconciliation_required", "owner gate revision order differs")
			}
			if _, e := time.Parse(time.RFC3339Nano, p.CreatedAt); e != nil {
				return fail("reconciliation_required", "invalid owner gate time")
			}
			if (p.Status == "prepared" && (n != len(i.OwnerGates)-1 || i.Revision != q.ExpectedRevision+1 || i.Disposition != nil)) || (p.Status == "committed" && i.Revision <= q.ExpectedRevision+1) || !slices.Contains([]string{"prepared", "committed"}, p.Status) {
				return fail("reconciliation_required", "invalid owner gate commitment")
			}
			r, ok := st.Receipts[q.OperationID]
			if !ok || r.Digest != requestDigest(q) || r.Result.Session == nil || r.Result.Session.Operation != "session-owner" || r.Result.Session.Revision != q.ExpectedRevision+1 || r.Result.Session.InstanceID != i.ID || r.Result.Session.RecordEffects != "prepared" {
				return fail("reconciliation_required", "owner gate receipt differs")
			}
			seen[q.OperationID] = true
			_, cards, e := dispositionInputs(prior)
			if e != nil {
				return e
			}
			allowed := map[string]bool{i.Context.Binding.Session.Path: true, i.Context.Binding.Plan.Path: true, managementLedger: true}
			for _, p := range cards {
				allowed[p] = true
			}
			for p := range g.Before {
				if !allowed[p] {
					return fail("reconciliation_required", "unselected owner gate preimage")
				}
			}
			after, e := planOwnerGate(prior, g)
			if e != nil {
				return e
			}
			if !reflect.DeepEqual(dispositionChanges(g.Before, after), p.Changes) {
				return fail("reconciliation_required", "owner gate transformation differs")
			}
			previous := dispositionPredecessor(prior)
			inputs := map[string]string{}
			for _, src := range p.Inputs {
				if !literalCandidatePath(src.Path) || inputs[src.Path] != "" || (!digestPattern.MatchString(src.SHA256) && src.SHA256 != "absent") {
					return fail("reconciliation_required", "invalid owner gate inputs")
				}
				inputs[src.Path] = src.SHA256
			}
			size := 0
			for path, b := range g.Before {
				size += len(b)
				if inputs[path] != hash([]byte(b)) {
					return fail("reconciliation_required", "owner preimage digest differs")
				}
			}
			for _, src := range previous.Coverage {
				if src.Path != managementLedger && inputs[src.Path] != src.SHA256 {
					return fail("reconciliation_required", "owner gate lost coverage")
				}
			}
			for _, ch := range p.Changes {
				size += len(ch.Content)
				if inputs[ch.Path] != ch.Before {
					return fail("reconciliation_required", "owner effect preimage differs")
				}
			}
			want := ownerGateCoverage(previous, g)
			if size > 8<<20 || len(p.Inputs) > 2048 || len(p.Coverage) > 512 || p.LedgerBytes <= 0 || len(want) != len(p.Coverage) {
				return fail("reconciliation_required", "owner gate limits/coverage differ")
			}
			for n, src := range want {
				got := p.Coverage[n]
				if src.Path != got.Path || (src.Path != managementLedger && src.SHA256 != got.SHA256) || (!digestPattern.MatchString(got.SHA256) && got.SHA256 != "absent") {
					return fail("reconciliation_required", "owner coverage differs")
				}
			}
		}
	}
	for _, r := range st.Receipts {
		if v := r.Result.Session; v != nil && v.Operation == "session-owner" && !seen[v.OperationID] {
			return fail("reconciliation_required", "orphan owner receipt")
		}
	}
	return nil
}
func projectOwnerGate(i SessionInstance, v *SessionView, step *SessionStepView) {
	if !step.Selected || i.Disposition != nil {
		return
	}
	g := latestOwnerGate(i)
	if g == nil {
		if step.Stage == "disposition" && selectedSessionStep(i).OwnerAcceptance == "required" && v.Coverage == "current" {
			v.NextActions = append(v.NextActions, "session-owner")
		}
		return
	}
	step.OwnerGate = &g.Request
	v.NextActions = []string{"session-read", "session-resume"}
	if g.Effects.Status != "committed" {
		step.Holds = append(step.Holds, "owner gate effects pending")
		return
	}
	step.Holds = append([]string{}, v.Conditions...)
	switch g.Request.Action {
	case "request":
		step.Stage = "waiting-owner"
		step.Holds = append(step.Holds, "actual owner decision required for this exact result")
		v.NextActions = append(v.NextActions, "session-owner")
	case "accept":
		step.Stage = "disposition"
		if v.Coverage == "current" {
			v.NextActions = append(v.NextActions, "session-dispose")
		}
	case "changes-requested":
		step.Stage = "rework"
		step.Holds = append(step.Holds, "owner requested changes; a new implementation and verified result are required")
	}
}
