package governance

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/Hans-Einar/SDP/SDPTool/projecthistory"
)

type SessionAdvanceRequest struct {
	Successor           *SessionSuccessor `json:"successor,omitempty"`
	Schema              string            `json:"schema"`
	ProjectRef          string            `json:"projectRef"`
	CheckoutID          string            `json:"checkoutId"`
	OperationID         string            `json:"operationId"`
	InstanceID          string            `json:"instanceId"`
	ExpectedRevision    uint64            `json:"expectedRevision"`
	ExpectedRunRevision uint64            `json:"expectedRunRevision"`
	Action              string            `json:"action"`
	StepID              string            `json:"stepId"`
	Basis               string            `json:"basis"`
	ResultID            string            `json:"resultId"`
	Reason              string            `json:"reason"`
}
type SessionAdvance struct {
	Successor   *SessionContext       `json:"successor,omitempty"`
	Request     SessionAdvanceRequest `json:"request"`
	PriorDigest string                `json:"priorDigest"`
	Effects     SessionRecordEffects  `json:"effects"`
	Before      map[string]string     `json:"before"`
}

func advancePath(q SessionAdvanceRequest) string {
	return "SDP/Sessions/execution/transitions/" + q.OperationID + ".json"
}
func advanceExecution(i SessionInstance, revision uint64) SessionInstance {
	i.Past = nil
	i.Advance = nil
	i.Rebindings = nil
	i.Rebind = nil
	i.Revision = revision
	return i
}
func validSessionAdvance(q SessionAdvanceRequest) error {
	if q.Schema != SessionSchema || !uuidPattern.MatchString(q.OperationID) || !uuidPattern.MatchString(q.InstanceID) || !uuidPattern.MatchString(q.ResultID) || q.ExpectedRevision == 0 || q.ExpectedRunRevision == 0 || !validLeaderText(q.Reason, 2048) || !slices.Contains([]string{"next", "rework", "complete", "transfer"}, q.Action) {
		return fail("invalid_request", "exact Session continuation identity, revisions and reason required")
	}
	if terminalSessionAdvance(q.Action) {
		if q.StepID != "" || q.Basis != "disposition" {
			return fail("invalid_request", "completion requires disposition without a next step")
		}
	} else if !sessionStepPattern.MatchString(q.StepID) {
		return fail("invalid_request", "bounded target step required")
	}
	if (q.Action == "transfer") != (q.Successor != nil) {
		return fail("invalid_request", "successor required only for transfer")
	}
	if q.Successor != nil && (!digestPattern.MatchString(q.Successor.InputRevision) || !literalCandidatePath(q.Successor.BindingPath)) {
		return fail("invalid_request", "exact successor context required")
	}
	if q.Action == "rework" {
		if !slices.Contains([]string{"owner", "test", "review"}, q.Basis) {
			return fail("invalid_request", "exact failed result or owner changes request required")
		}
	} else if q.Basis != "disposition" {
		return fail("invalid_request", "committed disposition basis required")
	}
	return nil
}
func advanceEligibility(i SessionInstance, q SessionAdvanceRequest) error {
	if i.Selection == nil || len(i.Work) < 2 || sessionLatestEffects(i).Status != "committed" {
		return fail("missing_coverage", "committed selected work required")
	}
	if q.Action == "rework" {
		if i.Disposition != nil || q.StepID != selectedSessionStep(i).ID {
			return fail("conflict", "rework must retain unfinished selected step")
		}
		if q.Basis == "owner" {
			g := latestOwnerGate(i)
			if g == nil || g.Effects.Status != "committed" || g.Request.Action != "changes-requested" || g.Request.OperationID != q.ResultID {
				return fail("missing_coverage", "actual owner changes request required")
			}
		}
		return nil
	}
	if i.Disposition == nil || i.Disposition.Effects.Status != "committed" || i.Disposition.Request.OperationID != q.ResultID {
		return fail("missing_coverage", "exact completed disposition required")
	}
	if terminalSessionAdvance(q.Action) {
		for _, step := range i.Context.Binding.Steps {
			if !sessionStepDone(i, step.ID) {
				return fail("conflict", "all selected steps must be disposed before Session closure")
			}
		}
		return nil
	}
	for _, step := range i.Context.Binding.Steps {
		if step.ID == q.StepID {
			if sessionStepDone(i, step.ID) || !sessionDependenciesSatisfied(i, step) {
				return fail("conflict", "target step completed or dependencies incomplete")
			}
			return nil
		}
	}
	return fail("missing_coverage", "next step absent from selected binding")
}
func (s *Service) advanceBasis(st snapshot, i SessionInstance, a SessionAdvance) error {
	q := a.Request
	if e := advanceEligibility(i, q); e != nil {
		return e
	}
	r := st.Runs[i.Work[0].Request.OperationID]
	if r.Revision != q.ExpectedRunRevision || r.Revision == ^uint64(0) {
		return fail("stale_revision", "continuation run revision differs")
	}
	if q.Action != "rework" {
		if r.State != "completed" {
			return fail("conflict", "disposition run not completed")
		}
		return nil
	}
	if e := s.progressSnapshot(st, i, i.Work[1]); e != nil {
		return e
	}
	if e := advanceStoredBasis(st, i, q); e != nil {
		return e
	}
	if q.Basis == "review" {
		if e := CheckSkillSnapshots(s.root, st.Reviews[q.ResultID].Spec.Skills); e != nil {
			return e
		}
	}
	for _, src := range advanceBasisSources(st, q) {
		b, e := readFile(filepath.Join(s.root, src.Path), 4<<20)
		if e != nil || hash(b) != src.SHA256 {
			return fail("stale_candidate", "rework basis input changed")
		}
	}

	for _, d := range st.Reviews {
		if d.RunID == r.ID && ReviewActive(d.State) {
			return fail("conflict", "active review must reconcile before rework")
		}
	}
	return nil
}

// Historical validation checks exact report/candidate identity without requiring
// historical candidates or skill files to remain the current working tree.
func advanceStoredBasis(st snapshot, i SessionInstance, q SessionAdvanceRequest) error {
	if q.Action != "rework" {
		return nil
	}
	if len(i.Work) < 2 || i.Work[1].Evidence == nil {
		return fail("missing_coverage", "reported candidate required")
	}
	switch q.Basis {
	case "test":
		ev := st.Receipts[q.ResultID].Result.Evidence
		if ev == nil || ev.Request.Kind != "test" || ev.Request.Verdict != "failed" || !sameSessionCandidate(*ev, *i.Work[1].Evidence) {
			return fail("missing_coverage", "exact failed test for reported candidate required")
		}
	case "review":
		d, ok := st.Reviews[q.ResultID]
		if !ok || d.State != "completed" || d.Outcome == nil || d.Outcome.Verdict == "passed" || d.RunID != i.Work[0].Request.OperationID || d.Spec.ReportID != i.Work[1].Evidence.Request.ReportID || d.Spec.Candidate.Digest != i.Work[1].Evidence.Request.Candidate.Digest {
			return fail("missing_coverage", "exact completed nonpassing review required")
		}
	}
	return nil
}
func advanceBasisSources(st snapshot, q SessionAdvanceRequest) []Source {
	if q.Action != "rework" {
		return nil
	}
	if q.Basis == "test" {
		if ev := st.Receipts[q.ResultID].Result.Evidence; ev != nil {
			return []Source{ev.Request.Artifact}
		}
	}
	return nil
}
func advanceArtifact(i SessionInstance, a SessionAdvance) string {
	b, _ := json.MarshalIndent(struct {
		Schema        string                `json:"schema"`
		Request       SessionAdvanceRequest `json:"request"`
		PriorDigest   string                `json:"priorDigest"`
		CreatedAt     string                `json:"createdAt"`
		BindingDigest string                `json:"bindingDigest"`
		RunID         string                `json:"runId"`
		Successor     *SessionContext       `json:"successor,omitempty"`
	}{"sdp-session-transition/0.1", a.Request, a.PriorDigest, a.Effects.CreatedAt, i.Context.BindingDigest, i.Work[0].Request.OperationID, a.Successor}, "", "  ")
	return string(b) + "\n"
}
func advanceJournal(i SessionInstance, a SessionAdvance, p string) string {
	rel, _ := filepath.Rel(filepath.Dir(p), advancePath(a.Request))
	description := "Completed step retained; next selected step requires fresh assessment and assignment."
	if a.Request.Action == "rework" {
		description = "Prior implementation and failed/owner result retained; fresh implementation required. Prior owner gate is no longer current; no acceptance inferred."
	} else if terminalSessionAdvance(a.Request.Action) {
		description = "All selected steps and whole-card outcomes disposed; selected plan and Session completed."
	}
	if a.Successor != nil {
		rel, _ := filepath.Rel(filepath.Dir(p), a.Successor.Binding.Session.Path)
		description = "Selected predecessor work completed; remaining cards transferred to [" + a.Successor.Binding.Session.ID + "](" + (&url.URL{Path: filepath.ToSlash(rel)}).EscapedPath() + ") and plan " + a.Successor.Binding.Plan.ID + ". Cards remain open; fresh adoption and assessment required. Predecessor: " + i.Context.Binding.Session.ID + "."
	}
	return "\n## PG Session " + a.Request.Action + " " + a.Request.OperationID + "\n\n" + description + "\n[Exact continuation record](" + filepath.ToSlash(rel) + ").\n\n" + publicationProse(a.Request.Reason)
}
func advanceCardPaths(i SessionInstance, files map[string]string, q SessionAdvanceRequest) (map[string]string, error) {
	required := map[string]bool{}
	steps := []SessionStepBinding{selectedSessionStep(i)}
	if terminalSessionAdvance(q.Action) {
		steps = i.Context.Binding.Steps
	}
	for _, step := range steps {
		for _, c := range step.Cards {
			required[c.ID] = true
		}
	}
	cards := map[string]string{}
	for p, b := range files {
		if strings.HasPrefix(p, "SDP/KanBan/") {
			m, e := metadata([]byte(b))
			if e == nil && required[m["id"]] {
				if cards[m["id"]] != "" {
					return nil, fail("conflict", "duplicate selected card")
				}
				cards[m["id"]] = p
			}
		}
	}
	if len(cards) != len(required) {
		return nil, fail("missing_coverage", "all continuation cards must exist")
	}
	return cards, nil
}
func planSessionAdvance(i SessionInstance, a SessionAdvance) (map[string]string, error) {
	if e := validateSuccessorRecord(i, a); e != nil {
		return nil, e
	}
	files := map[string]string{}
	for p, b := range a.Before {
		files[p] = b
	}
	cards, e := advanceCardPaths(i, files, a.Request)
	if e != nil {
		return nil, e
	}
	subjects := map[string]string{i.Context.Binding.Plan.ID: i.Context.Binding.Plan.Path}
	if a.Request.Action == "transfer" {
		subjects[a.Successor.Binding.Plan.ID] = a.Successor.Binding.Plan.Path
	}
	for id, p := range cards {
		subjects[id] = p
	}
	events, e := publicationEvents([]byte(files[managementLedger]))
	if e != nil {
		return nil, e
	}
	latest := map[string]map[string]any{}
	for _, ev := range events {
		latest[fmt.Sprint(ev["subjectId"])] = ev
	}
	ids := []string{}
	for id := range subjects {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		p := subjects[id]
		body, ok := files[p]
		if !ok {
			return nil, fail("missing_coverage", "continuation preimage absent")
		}
		m, e := metadata([]byte(body))
		if e != nil {
			return nil, e
		}
		if strings.HasPrefix(id, "KB-") {
			if a.Request.Action == "complete" && m["CardState"] != "completed" {
				return nil, fail("missing_coverage", "remaining card scope needs explicit successor disposition before overall closure")
			}
			if a.Request.Action == "transfer" && m["CardState"] != "completed" {
				if m["CardState"] != "in-progress" {
					return nil, fail("conflict", "remaining card must be in-progress")
				}
				b, e := replaceMetadata([]byte(body), "PlanId", a.Successor.Binding.Plan.ID)
				if e != nil {
					return nil, e
				}
				b, e = replaceMetadata(b, "CardState", "ready")
				if e != nil {
					return nil, e
				}
				body = string(b)
			}
			if a.Request.Action == "rework" {
				if !slices.Contains([]string{"in-progress", "gate-review"}, m["CardState"]) {
					return nil, fail("conflict", "rework card must be active")
				}
				b, e := replaceMetadata([]byte(body), "CardState", "in-progress")
				if e != nil {
					return nil, e
				}
				body = string(b)
			}
		} else if terminalSessionAdvance(a.Request.Action) && id == i.Context.Binding.Plan.ID {
			if m["state"] != "active" {
				return nil, fail("conflict", "active selected plan required")
			}
			b, e := replaceMetadata([]byte(body), "state", "completed")
			if e != nil {
				return nil, e
			}
			body = string(b)
		}
		files[p] = appendRecord([]byte(body), advanceJournal(i, a, p))
		old := latest[id]
		if old == nil {
			return nil, fail("missing_coverage", "continuation history absent")
		}
		prev, ok := old["payload"].(map[string]any)
		if !ok {
			return nil, fail("reconciliation_required", "continuation history malformed")
		}
		payload := map[string]any{}
		for k, v := range prev {
			payload[k] = v
		}
		delete(payload, "lineage")
		payload["previousEventId"] = old["eventId"]
		payload["from"] = prev["to"]
		payload["fromPath"] = prev["toPath"]
		payload["reason"] = "Session " + a.Request.Action + " " + a.Request.OperationID + "; prior evidence retained."
		if a.Successor != nil {
			links := []string{}
			if prior, ok := payload["links"].([]any); ok {
				for _, value := range prior {
					if id, ok := value.(string); ok {
						links = append(links, id)
					}
				}
			}
			links = append(links, i.Context.Binding.Plan.ID, a.Successor.Binding.Plan.ID)
			sort.Strings(links)
			payload["links"] = slices.Compact(links)
			payload["reason"] = fmt.Sprint(payload["reason"]) + " Remaining cards transferred to " + a.Successor.Binding.Plan.ID + "; fresh successor adoption required."
		}
		prefix, kind := "EVT-PM-"+strings.Split(id, "-")[1]+"-", "x-management:updated"
		if strings.HasPrefix(id, "KB-") {
			prefix = "EVT-KB-" + strings.Split(id, "-")[1] + "-"
			kind = "x-kanban:reviewed"
		} else if terminalSessionAdvance(a.Request.Action) && id == i.Context.Binding.Plan.ID {
			payload["to"] = "completed"
			kind = "x-management:completed"
		}
		eid, e := nextPublicationEvent(events, prefix)
		if e != nil {
			return nil, e
		}
		ev := map[string]any{"schemaVersion": "1.0", "eventId": eid, "eventType": kind, "occurredAt": a.Effects.CreatedAt, "actor": "sdp-controller", "commit": nil, "subjectId": id, "payload": payload}
		raw, _ := json.Marshal(ev)
		files[managementLedger] = appendRecord([]byte(files[managementLedger]), string(raw)+"\n")
		events = append(events, ev)
	}
	sp := i.Context.Binding.Session.Path
	body, ok := files[sp]
	if !ok {
		return nil, fail("missing_coverage", "continuation Session absent")
	}
	if terminalSessionAdvance(a.Request.Action) {
		b, e := replaceMetadata([]byte(body), "Status", "completed")
		if e != nil {
			return nil, e
		}
		body = string(b)
	}
	files[sp] = appendRecord([]byte(body), advanceJournal(i, a, sp))
	if a.Successor != nil {
		p := a.Successor.Binding.Session.Path
		files[p] = appendRecord([]byte(files[p]), advanceJournal(i, a, p))
	}
	files[advancePath(a.Request)] = advanceArtifact(i, a)
	if a.Request.Action == "rework" {
		if _, ok := files[dispositionClaim(i)]; !ok {
			return nil, fail("missing_coverage", "rework claim absent")
		}
		delete(files, dispositionClaim(i))
	}
	return files, nil
}
func advanceCoverage(previous SessionRecordEffects, a SessionAdvance, basis []Source) []Source {
	previous.Coverage = append(append([]Source{}, previous.Coverage...), basis...)
	if a.Successor != nil {
		previous.Coverage = mergeSessionSources(previous.Coverage, a.Successor.Sources)
	}
	return ownerGateCoverage(previous, SessionOwnerGate{Effects: a.Effects})
}
func (s *Service) buildSessionAdvance(st snapshot, i SessionInstance, a SessionAdvance) (SessionAdvance, error) {
	files, e := s.publicationFiles()
	if e != nil {
		return a, e
	}
	previous := *sessionLatestEffects(i)
	if a.Request.Action == "transfer" {
		ctx, e := s.readSuccessorContext(i, a.Request.Successor.BindingPath, files)
		if e != nil {
			return a, e
		}
		ctx.CheckoutID = st.CheckoutID
		if ctx.InputRevision != a.Request.Successor.InputRevision {
			return a, fail("stale_revision", "successor context changed")
		}
		a.Successor = &ctx
	}
	for _, src := range previous.Coverage {
		b, err := readFile(filepath.Join(s.root, src.Path), maxState)
		if os.IsNotExist(err) && src.SHA256 == "absent" {
			continue
		}
		if err != nil {
			return a, err
		}
		if src.Path != managementLedger && hash(b) != src.SHA256 {
			return a, fail("stale_revision", "continuation input changed")
		}
		files[src.Path] = b
	}
	textFiles := map[string]string{}
	for p, b := range files {
		textFiles[p] = string(b)
	}
	cards, e := advanceCardPaths(i, textFiles, a.Request)
	if e != nil {
		return a, e
	}
	paths := []string{i.Context.Binding.Session.Path, i.Context.Binding.Plan.Path, managementLedger}
	for _, p := range cards {
		paths = append(paths, p)
	}
	if a.Request.Action == "rework" {
		paths = append(paths, dispositionClaim(i))
	}
	if a.Successor != nil {
		for _, src := range a.Successor.Sources {
			paths = append(paths, src.Path)
		}
	}
	a.Before = map[string]string{}
	for _, p := range paths {
		b, ok := files[p]
		if !ok {
			return a, fail("missing_coverage", "continuation record absent")
		}
		a.Before[p] = string(b)
	}
	if _, e = os.Lstat(filepath.Join(s.root, advancePath(a.Request))); !os.IsNotExist(e) {
		return a, fail("conflict", "continuation destination occupied")
	}
	after, e := planSessionAdvance(i, a)
	if e != nil {
		return a, e
	}
	a.Effects.Changes = dispositionChanges(a.Before, after)
	inputs := map[string]string{}
	for p, b := range files {
		inputs[p] = hash(b)
	}
	for _, src := range previous.Coverage {
		if _, ok := inputs[src.Path]; !ok {
			inputs[src.Path] = src.SHA256
		}
	}
	inputs[advancePath(a.Request)] = "absent"
	// Preserve exact failed evidence through continuation and later selections.
	for _, src := range advanceBasisSources(st, a.Request) {
		inputs[src.Path] = src.SHA256
	}

	names := []string{}
	for p := range inputs {
		names = append(names, p)
	}
	sort.Strings(names)
	for _, p := range names {
		a.Effects.Inputs = append(a.Effects.Inputs, Source{p, inputs[p]})
	}
	for p := range a.Before {
		delete(files, p)
	}
	for p, b := range after {
		files[p] = []byte(b)
	}
	if e = validatePublicationHistory(files[managementLedger], files); e != nil {
		return a, e
	}
	a.Effects.Coverage = advanceCoverage(previous, a, advanceBasisSources(st, a.Request))
	a.Effects.LedgerBytes = len(files[managementLedger])
	size := 0
	for _, b := range a.Before {
		size += len(b)
	}
	for _, ch := range a.Effects.Changes {
		size += len(ch.Content)
	}
	if size > 8<<20 || len(a.Effects.Inputs) > 2048 || len(a.Effects.Coverage) > 512 {
		return a, fail("limit", "continuation capacity exceeded")
	}
	return a, nil
}
func (s *Service) AdvanceSession(c Caller, q SessionAdvanceRequest) (SessionMutationResult, error) {
	if !c.controller {
		return SessionMutationResult{}, fail("unauthorized", "trusted controller required")
	}
	if e := validSessionAdvance(q); e != nil {
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
	if i.Advance != nil {
		return SessionMutationResult{}, fail("conflict", "continuation already prepared or Session completed")
	}
	if i.Revision == ^uint64(0) || len(st.Receipts) >= 1024 || len(i.Past) >= 64 {
		return SessionMutationResult{}, fail("limit", "Session continuation capacity reached")
	}
	if e = s.sessionCurrent(i); e != nil {
		return SessionMutationResult{}, e
	}
	a := SessionAdvance{Request: q, PriorDigest: requestDigest(advanceExecution(i, q.ExpectedRevision)), Effects: SessionRecordEffects{Status: "prepared", CreatedAt: time.Now().UTC().Format(time.RFC3339Nano)}}
	if e = s.advanceBasis(st, i, a); e != nil {
		return SessionMutationResult{}, e
	}
	if e = s.sessionWorkClaims(st, i); e != nil {
		return SessionMutationResult{}, e
	}
	a, e = s.buildSessionAdvance(st, i, a)
	if e != nil {
		return SessionMutationResult{}, e
	}
	if e = s.transferClaims(st, i, a); e != nil {
		return SessionMutationResult{}, e
	}
	i.Advance = &a
	i.Revision++
	st.Sessions[i.ID] = i
	if a.Successor != nil && st.Schema != sessionRebindStateSchema {
		st.Schema = sessionTransferStateSchema
	} else if !sessionTransferProfile(st.Schema) {
		st.Schema = sessionLifecycleStateSchema
	}
	out := SessionMutationResult{Schema: SessionSchema, Operation: "session-advance", OperationID: q.OperationID, InstanceID: i.ID, Revision: i.Revision, Registration: "committed", RecordEffects: "prepared"}
	st.Receipts[q.OperationID] = receipt{requestDigest(q), Result{Schema: Schema, Operation: "session.advance", OperationID: q.OperationID, ProjectRef: s.project, CheckoutID: st.CheckoutID, Session: &out}}
	if e = validateSessionAdvances(st); e != nil {
		return SessionMutationResult{}, e
	}
	if e = s.save(st, head); e != nil {
		return SessionMutationResult{}, e
	}
	return out, nil
}
func (s *Service) applySessionAdvance(st *snapshot, i *SessionInstance, tx *projecthistory.Transaction) error {
	a := i.Advance
	prior := *i
	prior.Advance = nil
	if e := s.transferClaims(*st, prior, *a); e != nil {
		return e
	}
	if e := s.sessionWorkClaims(*st, prior); e != nil {
		return e
	}
	if e := s.applySessionRecordEffects(&a.Effects, tx, s.sessionPendingOwner(*i), "session-advance", func() error { return s.advanceBasis(*st, prior, *a) }); e != nil {
		return e
	}
	if terminalSessionAdvance(a.Request.Action) {
		return nil
	}
	if a.Request.Action == "rework" {
		r := st.Runs[i.Work[0].Request.OperationID]
		r.State = "canceled"
		r.Revision++
		r.Reason = "Superseded by checked Session rework; original report and failures retained"
		st.Runs[r.ID] = r
	}
	old := advanceExecution(prior, a.Request.ExpectedRevision)
	i.Past = append(i.Past, SessionPast{Execution: old, Transition: *a})
	i.Advance = nil
	i.Selection = nil
	i.Assessments = nil
	i.Work = nil
	i.OwnerGates = nil
	i.Disposition = nil
	return nil
}
func validateSessionAdvances(st snapshot) error {
	seen := map[string]bool{}
	for _, root := range sessionBindingEpochs(st) {
		validate := func(i SessionInstance, a SessionAdvance, current bool) error {
			q, p := a.Request, a.Effects
			prior := i
			prior.Advance = nil
			if !sessionLifecycleProfile(st.Schema) || (q.Action == "transfer" && !sessionTransferProfile(st.Schema)) || validSessionAdvance(q) != nil || advanceEligibility(prior, q) != nil || advanceStoredBasis(st, prior, q) != nil || a.PriorDigest != requestDigest(advanceExecution(prior, q.ExpectedRevision)) || q.ProjectRef != st.Project || q.CheckoutID != st.CheckoutID || q.InstanceID != root.ID || q.ExpectedRevision >= root.Revision || seen[q.OperationID] {
				return fail("reconciliation_required", "invalid Session continuation identity/basis")
			}
			seen[q.OperationID] = true
			if _, e := time.Parse(time.RFC3339Nano, p.CreatedAt); e != nil {
				return fail("reconciliation_required", "continuation time invalid")
			}
			if (p.Status == "prepared" && (!current || root.Revision != q.ExpectedRevision+1)) || (p.Status == "committed" && root.Revision <= q.ExpectedRevision+1) || !slices.Contains([]string{"prepared", "committed"}, p.Status) {
				return fail("reconciliation_required", "continuation commitment invalid")
			}
			r, ok := st.Receipts[q.OperationID]
			if !ok || r.Digest != requestDigest(q) || r.Result.Session == nil || r.Result.Session.Operation != "session-advance" || r.Result.Session.InstanceID != root.ID || r.Result.Session.Revision != q.ExpectedRevision+1 || r.Result.Session.RecordEffects != "prepared" {
				return fail("reconciliation_required", "continuation receipt differs")
			}
			cards, e := advanceCardPaths(prior, a.Before, q)
			if e != nil {
				return e
			}
			allowed := map[string]bool{i.Context.Binding.Session.Path: true, i.Context.Binding.Plan.Path: true, managementLedger: true}
			for _, p := range cards {
				allowed[p] = true
			}
			if q.Action == "rework" {
				allowed[dispositionClaim(prior)] = true
			}
			if a.Successor != nil {
				for _, src := range a.Successor.Sources {
					allowed[src.Path] = true
				}
			}
			for p := range a.Before {
				if !allowed[p] {
					return fail("reconciliation_required", "unselected continuation preimage")
				}
			}
			after, e := planSessionAdvance(prior, a)
			if e != nil {
				return e
			}
			if !reflect.DeepEqual(dispositionChanges(a.Before, after), p.Changes) {
				return fail("reconciliation_required", "continuation postimages differ")
			}
			inputs := map[string]string{}
			for _, src := range p.Inputs {
				if !literalCandidatePath(src.Path) || inputs[src.Path] != "" || (!digestPattern.MatchString(src.SHA256) && src.SHA256 != "absent") {
					return fail("reconciliation_required", "invalid continuation pins")
				}
				inputs[src.Path] = src.SHA256
			}
			if a.Successor != nil {
				for _, src := range a.Successor.Sources {
					if inputs[src.Path] != src.SHA256 {
						return fail("reconciliation_required", "successor input pin differs")
					}
				}
			}
			for _, src := range advanceBasisSources(st, q) {
				if inputs[src.Path] != src.SHA256 {
					return fail("reconciliation_required", "continuation basis pin differs")
				}
			}
			size := 0
			for path, b := range a.Before {
				size += len(b)
				if inputs[path] != hash([]byte(b)) {
					return fail("reconciliation_required", "continuation preimage differs")
				}
			}
			for _, ch := range p.Changes {
				size += len(ch.Content)
				if inputs[ch.Path] != ch.Before {
					return fail("reconciliation_required", "continuation effect pin differs")
				}
			}
			previous := *sessionLatestEffects(prior)
			for _, src := range previous.Coverage {
				if src.Path != managementLedger && inputs[src.Path] != src.SHA256 {
					return fail("reconciliation_required", "continuation lost prior coverage")
				}
			}
			want := advanceCoverage(previous, a, advanceBasisSources(st, a.Request))
			if size > 8<<20 || len(p.Inputs) > 2048 || len(p.Coverage) > 512 || len(want) != len(p.Coverage) || p.LedgerBytes <= 0 {
				return fail("reconciliation_required", "continuation limits/coverage invalid")
			}
			for n, src := range want {
				got := p.Coverage[n]
				if got.Path != src.Path || (src.Path != managementLedger && got.SHA256 != src.SHA256) || (!digestPattern.MatchString(got.SHA256) && got.SHA256 != "absent") {
					return fail("reconciliation_required", "continuation coverage differs")
				}
			}
			run := st.Runs[i.Work[0].Request.OperationID]
			if p.Status == "committed" && q.Action == "rework" && (run.State != "canceled" || run.Revision <= q.ExpectedRunRevision) {
				return fail("reconciliation_required", "reworked run not retired")
			}
			return nil
		}
		for n, p := range root.Past {
			i := p.Execution
			i.Past = root.Past[:n]
			i.Rebindings = root.Rebindings
			if e := validate(i, p.Transition, false); e != nil {
				return e
			}
		}
		if root.Advance != nil {
			if e := validate(root, *root.Advance, true); e != nil {
				return e
			}
		}
	}
	for _, r := range st.Receipts {
		if v := r.Result.Session; v != nil && v.Operation == "session-advance" && !seen[v.OperationID] {
			return fail("reconciliation_required", "orphan continuation receipt")
		}
	}
	return nil
}
