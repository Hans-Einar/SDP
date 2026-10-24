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

const sessionWorkStateSchema = "sdp-project-governance-state/0.16"

func sessionAssessmentProfile(schema string) bool {
	return schema == sessionAssessmentStateSchema || schema == sessionWorkStateSchema
}
func SessionImplementation() Routine { r := Maintenance(); r.ID = "session-implementation"; return r }

// Start and progress have disjoint additional fields. Adapters select the verb;
// request JSON cannot select a role, alternate run or broaden the binding scope.
type SessionWorkRequest struct {
	Schema              string       `json:"schema"`
	ProjectRef          string       `json:"projectRef"`
	CheckoutID          string       `json:"checkoutId"`
	OperationID         string       `json:"operationId"`
	InstanceID          string       `json:"instanceId"`
	ExpectedRevision    uint64       `json:"expectedRevision"`
	Reason              string       `json:"reason"`
	AssessmentID        string       `json:"assessmentId,omitempty"`
	Outcome             string       `json:"outcome,omitempty"`
	Skills              []SkillInput `json:"skills,omitempty"`
	EvidenceID          string       `json:"evidenceId,omitempty"`
	ExpectedRunRevision uint64       `json:"expectedRunRevision,omitempty"`
}
type SessionWork struct {
	Kind     string               `json:"kind"`
	Request  SessionWorkRequest   `json:"request"`
	Evidence *EvidenceRecord      `json:"evidence,omitempty"`
	Effects  SessionRecordEffects `json:"effects"`
}

func workDigest(kind string, q SessionWorkRequest) string {
	return requestDigest(struct {
		Kind    string
		Request SessionWorkRequest
	}{kind, q})
}
func workPath(w SessionWork) string {
	return "SDP/Sessions/execution/work/" + w.Request.OperationID + ".json"
}
func validSessionWork(kind string, q SessionWorkRequest) error {
	if q.Schema != SessionSchema || !uuidPattern.MatchString(q.InstanceID) || !uuidPattern.MatchString(q.OperationID) || q.ExpectedRevision == 0 || !validLeaderText(q.Reason, 2048) {
		return fail("invalid_request", "exact Session identity/revision and reason required")
	}
	switch kind {
	case "start":
		if !uuidPattern.MatchString(q.AssessmentID) || !validLeaderText(q.Outcome, 4096) || !ValidSkillInputs(q.Skills) || q.ExpectedRunRevision != 0 || q.EvidenceID != "" {
			return fail("invalid_request", "bounded start outcome, assessment and skills required")
		}
		master, worker := false, false
		for _, s := range q.Skills {
			if filepath.Base(s.Path) != "SKILL.md" || filepath.Base(filepath.Dir(s.Path)) != s.Name {
				return fail("invalid_request", "skill name/path mismatch")
			}
			master = master || s.Name == "sdp-master"
			worker = worker || s.Name == "sdp-worker"
		}
		if !master || !worker {
			return fail("missing_coverage", "Master and Worker skill snapshots required")
		}
	case "progress":
		if q.AssessmentID != "" || q.Outcome != "" || len(q.Skills) > 0 || !uuidPattern.MatchString(q.EvidenceID) || q.ExpectedRunRevision == 0 {
			return fail("invalid_request", "exact run revision and evidence receipt required")
		}
	default:
		return fail("invalid_request", "unsupported Session work verb")
	}
	return nil
}
func sessionRunBinding(st snapshot, id string) (SessionInstance, bool) {
	for _, i := range st.Sessions {
		if len(i.Work) > 0 && i.Work[0].Request.OperationID == id {
			return i, true
		}
	}
	return SessionInstance{}, false
}
func sessionRun(i SessionInstance) Run {
	w := i.Work[0]
	a := i.Assessments[len(i.Assessments)-1]
	step := selectedSessionStep(i)
	routine := SessionImplementation()
	sel := Selection{Plan: i.Context.Binding.Plan.Path, Session: i.Context.Binding.Session.Path, Authority: i.Context.Binding.Authority, Step: step.Milestone}
	// Selection.Card is the legacy primary-card projection only. Coverage/payload
	// and canonical effects retain every selected card and contribution.
	for _, ch := range i.Selection.Changes {
		if strings.HasPrefix(ch.Path, "SDP/KanBan/active/") {
			m, e := metadata([]byte(ch.Content))
			if e == nil && m["id"] == step.Cards[0].ID {
				sel.Card = ch.Path
			}
		}
	}
	revision := requestDigest(struct {
		Binding, Assessment, Routine string
		Request                      SessionWorkRequest
	}{i.Context.BindingDigest, requestDigest(a.Request), routine.Digest(), w.Request})
	return Run{ID: w.Request.OperationID, RequestID: w.Request.OperationID, Revision: 1, State: "ready", Step: "implement", InputRevision: revision, Selection: sel, Assignment: Assignment{Outcome: w.Request.Outcome, WritablePaths: append([]string{}, step.WritablePaths...), ExcludedPaths: append([]string{}, step.ProtectedPaths...)}, Routine: routine, RoutineDigest: routine.Digest(), Reason: "Session start records committed; explicit controller launch required"}
}
func (s *Service) sessionRunContext(r Run) (Context, error) {
	st, _, e := s.load()
	if e != nil {
		return Context{}, e
	}
	i, ok := sessionRunBinding(st, r.ID)
	if !ok {
		return Context{}, fail("reconciliation_required", "Session run binding absent")
	}
	c := Context{Schema: Schema, Operation: "context", ProjectRef: s.project, Root: s.root, CheckoutID: st.CheckoutID, Revision: sessionRun(i).InputRevision, Selection: r.Selection, Sources: append([]Source{}, sessionLatestEffects(i).Coverage...), Coverage: "known", Conditions: []string{}, CardID: selectedSessionStep(i).Cards[0].ID, PlanID: i.Context.Binding.Plan.ID, SessionID: i.Context.Binding.Session.ID}
	if e = s.sessionCurrent(i); e != nil {
		c.Coverage = "stale"
		c.Conditions = append(c.Conditions, e.Error())
	}
	return c, nil
}

// RunSession is read-only host context, preserving every selected contribution.
func (s *Service) RunSession(id string) (*SessionView, error) {
	st, _, e := s.load()
	if e != nil {
		return nil, e
	}
	i, ok := sessionRunBinding(st, id)
	if !ok {
		return nil, nil
	}
	v := s.sessionView(i, st)
	return &v, nil
}
func checkSessionLaunchSkills(st snapshot, r Run, skills []SkillInput) error {
	if r.Routine.ID != SessionImplementation().ID {
		return nil
	}
	i, ok := sessionRunBinding(st, r.ID)
	if !ok {
		return fail("reconciliation_required", "Session run binding absent")
	}
	if len(i.Work) > 1 {
		return fail("conflict", "Session progress must be reconciled before a new launch")
	}
	for _, required := range i.Work[0].Request.Skills {
		if !slices.Contains(skills, required) {
			return fail("missing_coverage", "launch lacks exact Session skill snapshot: "+required.Name)
		}
	}
	return nil
}
func sessionWorkArtifact(i SessionInstance, w SessionWork) []byte {
	b, _ := json.MarshalIndent(struct {
		Schema        string             `json:"schema"`
		Kind          string             `json:"kind"`
		BindingDigest string             `json:"bindingDigest"`
		RunID         string             `json:"runId"`
		CreatedAt     string             `json:"createdAt"`
		Request       SessionWorkRequest `json:"request"`
		Evidence      *EvidenceRecord    `json:"evidence,omitempty"`
	}{"sdp-session-work/0.1", w.Kind, i.Context.BindingDigest, i.Work[0].Request.OperationID, w.Effects.CreatedAt, w.Request, w.Evidence}, "", "  ")
	return append(b, '\n')
}
func sessionWorkJournal(i SessionInstance, w SessionWork, path string) string {
	rel, _ := filepath.Rel(filepath.Dir(path), workPath(w))
	out := "\n## PG Session " + w.Kind + " " + w.Request.OperationID + "\n\nStep " + selectedSessionStep(i).ID + "; run " + i.Work[0].Request.OperationID + ".\n[Exact assignment and evidence](" + filepath.ToSlash(rel) + ").\n"
	if w.Kind == "start" {
		out += "Assignment in-progress; host launch is a separate observed operation.\n"
	} else {
		out += "Result submitted for verification; no acceptance, card completion or model release.\n"
	}
	return out + "\nReason:\n\n" + publicationProse(w.Request.Reason)
}
func (s *Service) buildSessionWork(i SessionInstance, w SessionWork) (SessionWork, error) {
	files, e := s.publicationFiles()
	if e != nil {
		return w, e
	}
	previous := &i.Assessments[len(i.Assessments)-1].Effects
	if w.Kind == "progress" {
		previous = &i.Work[0].Effects
	}
	pins := append([]Source{}, previous.Coverage...)
	for _, skill := range w.Request.Skills {
		p, _ := filepath.Rel(s.root, skill.Path)
		pins = append(pins, Source{filepath.ToSlash(p), skill.Digest})
	}
	if w.Evidence != nil {
		pins = append(pins, w.Evidence.Request.Artifact)
	}
	for _, pin := range pins {
		b, err := readFile(filepath.Join(s.root, pin.Path), maxState)
		if os.IsNotExist(err) && pin.SHA256 == "absent" {
			continue
		}
		if err != nil {
			return w, err
		}
		if pin.Path != managementLedger && hash(b) != pin.SHA256 {
			return w, fail("stale_revision", "Session work input changed: "+pin.Path)
		}
		files[pin.Path] = b
	}
	before := map[string][]byte{}
	for p, b := range files {
		before[p] = b
	}
	if _, e = os.Lstat(filepath.Join(s.root, workPath(w))); !os.IsNotExist(e) {
		return w, fail("conflict", "work record destination exists or unavailable")
	}
	files[workPath(w)] = sessionWorkArtifact(i, w)
	subjects := map[string]string{i.Context.Binding.Plan.ID: i.Context.Binding.Plan.Path}
	for _, card := range selectedSessionStep(i).Cards {
		for p, b := range files {
			if strings.HasPrefix(p, "SDP/KanBan/active/") {
				m, err := metadata(b)
				if err == nil && m["id"] == card.ID {
					subjects[card.ID] = p
				}
			}
		}
		if subjects[card.ID] == "" {
			return w, fail("missing_coverage", "selected active card absent")
		}
	}
	events, e := publicationEvents(files[managementLedger])
	if e != nil {
		return w, e
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
		if strings.HasPrefix(id, "KB-") {
			m, _ := metadata(files[p])
			if !slices.Contains([]string{"ready", "in-progress"}, m["CardState"]) {
				return w, fail("conflict", "card cannot enter Session work")
			}
			files[p], e = replaceMetadata(files[p], "CardState", "in-progress")
			if e != nil {
				return w, e
			}
		}
		files[p] = []byte(appendRecord(files[p], sessionWorkJournal(i, w, p)))
		old := latest[id]
		if old == nil {
			return w, fail("missing_coverage", "selected history absent")
		}
		payload := map[string]any{}
		for k, v := range old["payload"].(map[string]any) {
			payload[k] = v
		}
		delete(payload, "lineage")
		payload["previousEventId"] = old["eventId"]
		payload["from"] = payload["to"]
		payload["fromPath"] = payload["toPath"]
		payload["reason"] = "Session " + w.Kind + " " + w.Request.OperationID + "; selected contribution in-progress, not completed."
		prefix, typ := "EVT-PM-"+strings.Split(id, "-")[1]+"-", "x-management:updated"
		if strings.HasPrefix(id, "KB-") {
			prefix = "EVT-KB-" + strings.Split(id, "-")[1] + "-"
			typ = "x-kanban:reviewed"
		}
		eid, err := nextPublicationEvent(events, prefix)
		if err != nil {
			return w, err
		}
		ev := map[string]any{"schemaVersion": "1.0", "eventId": eid, "eventType": typ, "occurredAt": w.Effects.CreatedAt, "actor": "sdp-controller", "commit": nil, "subjectId": id, "payload": payload}
		b, _ := json.Marshal(ev)
		files[managementLedger] = []byte(appendRecord(files[managementLedger], string(b)+"\n"))
		events = append(events, ev)
	}
	sp := i.Context.Binding.Session.Path
	files[sp] = []byte(appendRecord(files[sp], sessionWorkJournal(i, w, sp)))
	if e = validatePublicationHistory(files[managementLedger], files); e != nil {
		return w, e
	}
	names := []string{}
	for p := range files {
		names = append(names, p)
	}
	for _, pin := range pins {
		names = append(names, pin.Path)
	}
	sort.Strings(names)
	names = slices.Compact(names)
	covered := map[string]bool{workPath(w): true}
	for _, pin := range pins {
		covered[pin.Path] = true
	}
	size := 0
	for _, p := range names {
		pre, post := "absent", "absent"
		if b, ok := before[p]; ok {
			pre = hash(b)
		}
		if b, ok := files[p]; ok {
			post = hash(b)
		}
		w.Effects.Inputs = append(w.Effects.Inputs, Source{p, pre})
		if pre != post {
			w.Effects.Changes = append(w.Effects.Changes, PublicationChange{p, pre, post, string(files[p])})
			size += len(files[p])
		}
		if covered[p] {
			w.Effects.Coverage = append(w.Effects.Coverage, Source{p, post})
		}
	}
	w.Effects.LedgerBytes = len(files[managementLedger])
	if len(w.Effects.Inputs) > 2048 || size > 4<<20 {
		return w, fail("limit", "Session work effect capacity exceeded")
	}
	return w, nil
}
func (s *Service) sessionWorkClaims(st snapshot, i SessionInstance) error {
	step := selectedSessionStep(i)
	if e := s.checkSelectionClaims(i, i.Selection.Request, step); e != nil {
		return e
	}
	runID := ""
	if len(i.Work) > 0 {
		runID = i.Work[0].Request.OperationID
	}
	assignments := []Assignment{}
	for id, r := range st.Runs {
		if id != runID && r.State != "completed" && r.State != "canceled" {
			assignments = append(assignments, r.Assignment)
		}
	}
	for _, p := range st.Prepared {
		if p.State != "retired" {
			assignments = append(assignments, p.Assignment)
		}
	}
	for _, a := range assignments {
		for _, p := range append(append([]string{}, a.WritablePaths...), a.ExcludedPaths...) {
			if !literalCandidatePath(p) {
				return fail("conflict", "nonliteral existing scope requires reconciliation")
			}
		}
		for _, p := range a.WritablePaths {
			if !literalCandidatePath(p) {
				return fail("conflict", "nonliteral existing work requires scope reconciliation")
			}
			for _, b := range append(append([]string{}, step.WritablePaths...), step.ProtectedPaths...) {
				if overlaps(p, b) {
					return fail("conflict", "Session overlaps current run or prepared assignment")
				}
			}
		}
		for _, p := range a.ExcludedPaths {
			for _, b := range step.WritablePaths {
				if overlaps(p, b) {
					return fail("conflict", "Session writer overlaps existing protected scope")
				}
			}
		}
	}
	return nil
}
func (s *Service) progressSnapshot(st snapshot, i SessionInstance, w SessionWork) error {
	r, ok := st.Runs[i.Work[0].Request.OperationID]
	if !ok || r.Revision != w.Request.ExpectedRunRevision || r.State != "verifying" || w.Evidence == nil {
		return fail("stale_revision", "exact reporting run required")
	}
	ev := w.Evidence
	receipt, ok := st.Receipts[w.Request.EvidenceID]
	if !ok || receipt.Result.Evidence == nil || !reflect.DeepEqual(*receipt.Result.Evidence, *ev) || ev.Request.RunID != r.ID || latestReport(st, r) != ev.Request.ReportID {
		return fail("stale_revision", "bound report/evidence changed")
	}
	if len(r.Attempts) == 0 {
		return fail("missing_coverage", "assigned host attempt required")
	}
	for _, a := range r.Attempts {
		if a.Gap || !slices.Contains([]string{"completed", "failed", "interrupted"}, a.Status) {
			return fail("reconciliation_required", "terminal known host attempts required")
		}
	}
	for _, input := range r.Inputs {
		if input.State != "acknowledged" {
			return fail("reconciliation_required", "unresolved host intent")
		}
	}
	if hasBlockingAttention(r) {
		return fail("attention_required", "unresolved attention blocks progress")
	}
	c, e := s.captureCandidate(st.CheckoutID, r, ev.Request.Candidate.Dependencies)
	if e != nil || c.Digest != ev.Request.Candidate.Digest {
		return fail("stale_candidate", "reported candidate changed")
	}
	b, e := readFile(filepath.Join(s.root, ev.Request.Artifact.Path), 4<<20)
	if e != nil || hash(b) != ev.Request.Artifact.SHA256 {
		return fail("stale_candidate", "reported artifact changed")
	}
	return nil
}
func (s *Service) StartSession(c Caller, q SessionWorkRequest) (SessionMutationResult, error) {
	return s.prepareSessionWork(c, "start", q)
}
func (s *Service) ProgressSession(c Caller, q SessionWorkRequest) (SessionMutationResult, error) {
	return s.prepareSessionWork(c, "progress", q)
}
func (s *Service) prepareSessionWork(c Caller, kind string, q SessionWorkRequest) (SessionMutationResult, error) {
	if !c.controller {
		return SessionMutationResult{}, fail("unauthorized", "trusted controller required")
	}
	if e := validSessionWork(kind, q); e != nil {
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
	if r, ok, e := sessionRetry(st, q.OperationID, workDigest(kind, q)); ok {
		return r, e
	}
	if e = tx.Check(""); e != nil {
		return SessionMutationResult{}, e
	}
	if _, e = readFile(filepath.Join(dir, "pending"), 8192); !os.IsNotExist(e) {
		return SessionMutationResult{}, fail("reconciliation_required", "canonical publication pending")
	}
	i, ok := st.Sessions[q.InstanceID]
	if !ok || i.Selection == nil || len(i.Assessments) == 0 {
		return SessionMutationResult{}, fail("missing_coverage", "assessed selected Session required")
	}
	if i.Revision != q.ExpectedRevision {
		return SessionMutationResult{}, fail("stale_revision", "Session revision changed")
	}
	if i.Revision == ^uint64(0) || len(st.Receipts) >= 1024 {
		return SessionMutationResult{}, fail("limit", "Session capacity reached")
	}
	if e = s.sessionCurrent(i); e != nil {
		return SessionMutationResult{}, e
	}
	a := i.Assessments[len(i.Assessments)-1]
	if a.Effects.Status != "committed" || assessmentImpact(a.Request, selectedSessionStep(i)) != "unchanged" || a.Request.Blueprint.Requirement != "not-required" {
		return SessionMutationResult{}, fail("missing_coverage", "current explicit no-design/no-Blueprint route required")
	}
	w := SessionWork{Kind: kind, Request: q, Effects: SessionRecordEffects{Status: "prepared", CreatedAt: time.Now().UTC().Format(time.RFC3339Nano)}}
	if kind == "start" {
		if len(i.Work) != 0 || q.AssessmentID != a.Request.OperationID {
			return SessionMutationResult{}, fail("conflict", "unstarted Session and exact latest assessment required")
		}
		if len(st.Runs) >= 128 {
			return SessionMutationResult{}, fail("limit", "run capacity reached")
		}
		if _, exists := st.Runs[q.OperationID]; exists {
			return SessionMutationResult{}, fail("conflict", "run identity collision")
		}
		if e = CheckSkillSnapshots(s.root, q.Skills); e != nil {
			return SessionMutationResult{}, e
		}
		// Read-only governance inputs must not also be assigned implementation writes.
		for _, p := range selectedSessionStep(i).WritablePaths {
			for _, src := range sessionLatestEffects(i).Coverage {
				if overlaps(p, src.Path) {
					return SessionMutationResult{}, fail("conflict", "writable scope overlaps pinned governance inputs")
				}
			}
			for _, skill := range q.Skills {
				rel, _ := filepath.Rel(s.root, skill.Path)
				if overlaps(p, filepath.ToSlash(rel)) {
					return SessionMutationResult{}, fail("conflict", "writable scope overlaps required skills")
				}
			}
		}
	} else {
		if len(i.Work) != 1 || i.Work[0].Effects.Status != "committed" {
			return SessionMutationResult{}, fail("conflict", "one committed Session start required")
		}
		record := st.Receipts[q.EvidenceID].Result.Evidence
		if record == nil {
			return SessionMutationResult{}, fail("missing_coverage", "candidate-bound evidence receipt required")
		}
		copy := *record
		w.Evidence = &copy
		if e = s.progressSnapshot(st, i, w); e != nil {
			return SessionMutationResult{}, e
		}
	}
	if e = s.sessionWorkClaims(st, i); e != nil {
		return SessionMutationResult{}, e
	}
	i.Work = append(i.Work, w)
	w, e = s.buildSessionWork(i, w)
	if e != nil {
		return SessionMutationResult{}, e
	}
	i.Work[len(i.Work)-1] = w
	i.Revision++
	st.Sessions[i.ID] = i
	st.Schema = sessionWorkStateSchema
	out := SessionMutationResult{Schema: SessionSchema, Operation: "session-" + kind, OperationID: q.OperationID, InstanceID: i.ID, Revision: i.Revision, Registration: "committed", RecordEffects: "prepared"}
	st.Receipts[q.OperationID] = receipt{workDigest(kind, q), Result{Schema: Schema, Operation: "session." + kind, OperationID: q.OperationID, ProjectRef: s.project, CheckoutID: st.CheckoutID, Session: &out}}
	if e = s.save(st, head); e != nil {
		return SessionMutationResult{}, e
	}
	return out, nil
}
func (s *Service) applySessionWork(st *snapshot, i *SessionInstance, tx *projecthistory.Transaction) error {
	w := &i.Work[len(i.Work)-1]
	if e := s.sessionWorkClaims(*st, *i); e != nil {
		return e
	}
	check := func() error {
		if w.Kind == "progress" {
			return s.progressSnapshot(*st, *i, *w)
		}
		if len(st.Runs) >= 128 {
			return fail("limit", "run capacity reached before effects")
		}
		return nil
	}
	if e := s.applySessionRecordEffects(&w.Effects, tx, s.sessionPendingOwner(*i), "session-"+w.Kind, check); e != nil {
		return e
	}
	if w.Kind == "start" {
		r := sessionRun(*i)
		if _, ok := st.Runs[r.ID]; ok {
			return fail("conflict", "Session run identity collision")
		}
		st.Runs[r.ID] = r
	}
	return nil
}
func (s *Service) projectSessionWork(i SessionInstance, st snapshot, v *SessionView, step *SessionStepView) {
	if !step.Selected {
		return
	}
	if len(i.Work) == 0 {
		if step.Assessment != nil && step.Assessment.DesignImpact == "unchanged" && step.Assessment.Blueprint.Requirement == "not-required" && v.Coverage == "current" {
			v.NextActions = append(v.NextActions, "session-start")
			step.Holds = []string{"explicit Session start required"}
		}
		return
	}
	v.Binding = &i.Context.Binding
	w := i.Work[0]
	step.RunID = w.Request.OperationID
	if w.Effects.Status != "committed" {
		step.Holds = append(step.Holds, "start effects pending; no run created")
		return
	}
	r := st.Runs[step.RunID]
	step.RunState = r.State
	step.Holds = append([]string{}, v.Conditions...)
	v.HostExecution = "not-dispatched"
	if len(r.Attempts) > 0 {
		v.HostExecution = r.Attempts[len(r.Attempts)-1].Status
	} else if len(r.Inputs) > 0 {
		v.HostExecution = "launch-pending"
	}
	for _, input := range r.Inputs {
		if input.State == "uncertain" {
			v.HostExecution = "unknown"
			break
		}
		if input.State == "prepared" || input.State == "sending" {
			v.HostExecution = "launch-pending"
		}
	}
	if len(i.Work) == 2 && i.Work[1].Effects.Status == "committed" {
		step.Stage = "verification"
		step.ProgressEvidence = i.Work[1].Request.EvidenceID
		step.Holds = append(step.Holds, "verification, review and Session disposition required")
		if e := s.progressSnapshot(st, i, i.Work[1]); e != nil {
			step.Holds = append(step.Holds, e.Error())
		}
		return
	}
	step.Stage = "implementation"
	step.Holds = append(step.Holds, "assigned run report and current candidate evidence required")
	if len(i.Work) == 1 && r.State == "verifying" && v.Coverage == "current" {
		v.NextActions = append(v.NextActions, "session-progress")
	}
}

func validateSessionWork(st snapshot) error {
	seen := map[string]bool{}
	for _, i := range st.Sessions {
		if len(i.Work) == 0 {
			continue
		}
		if st.Schema != sessionWorkStateSchema || len(i.Work) > 2 || len(i.Assessments) == 0 || i.Selection == nil || !i.Registered {
			return fail("reconciliation_required", "invalid Session work owner/profile")
		}
		a := i.Assessments[len(i.Assessments)-1]
		if a.Effects.Status != "committed" || assessmentImpact(a.Request, selectedSessionStep(i)) != "unchanged" || a.Request.Blueprint.Requirement != "not-required" {
			return fail("reconciliation_required", "Session work lacks no-design disposition")
		}
		previous := a.Effects
		revision := a.Request.ExpectedRevision + 1
		for n, w := range i.Work {
			q, p := w.Request, w.Effects
			kind := []string{"start", "progress"}[n]
			if w.Kind != kind || validSessionWork(kind, q) != nil || q.InstanceID != i.ID || q.ProjectRef != st.Project || q.CheckoutID != st.CheckoutID || q.ExpectedRevision <= revision || q.ExpectedRevision >= i.Revision {
				return fail("reconciliation_required", "invalid Session work request")
			}
			revision = q.ExpectedRevision + 1
			if !slices.Contains([]string{"prepared", "committed"}, p.Status) || (p.Status == "prepared" && (n != len(i.Work)-1 || i.Revision != revision)) || (p.Status == "committed" && i.Revision <= revision) || len(p.Inputs) > 2048 || len(p.Coverage) > 256 || len(p.Changes) < 4 || len(p.Changes) > 20 || p.LedgerBytes <= 0 {
				return fail("reconciliation_required", "invalid Session work effects")
			}
			if _, e := time.Parse(time.RFC3339Nano, p.CreatedAt); e != nil {
				return fail("reconciliation_required", "invalid work time")
			}
			receipt, ok := st.Receipts[q.OperationID]
			if !ok || seen[q.OperationID] || receipt.Digest != workDigest(kind, q) || receipt.Result.Session == nil || receipt.Result.Session.Operation != "session-"+kind || receipt.Result.Session.InstanceID != i.ID || receipt.Result.Session.Revision != revision || receipt.Result.Session.RecordEffects != "prepared" {
				return fail("reconciliation_required", "Session work receipt mismatch")
			}
			seen[q.OperationID] = true
			if n == 0 {
				if q.AssessmentID != a.Request.OperationID || w.Evidence != nil {
					return fail("reconciliation_required", "Session start assessment differs")
				}
			} else {
				r := st.Receipts[q.EvidenceID].Result.Evidence
				if w.Evidence == nil || r == nil || !reflect.DeepEqual(r, w.Evidence) || r.Request.RunID != i.Work[0].Request.OperationID {
					return fail("reconciliation_required", "Session progress evidence differs")
				}
			}
			required := map[string]string{}
			for _, src := range previous.Coverage {
				required[src.Path] = src.SHA256
			}
			for _, skill := range q.Skills {
				p, e := filepath.Rel(st.Root, skill.Path)
				p = filepath.ToSlash(p)
				if e != nil || !literalCandidatePath(p) {
					return fail("reconciliation_required", "work skill outside project")
				}
				if old := required[p]; old != "" && old != skill.Digest {
					return fail("reconciliation_required", "work skill pin conflict")
				}
				required[p] = skill.Digest
			}
			if w.Evidence != nil {
				required[w.Evidence.Request.Artifact.Path] = w.Evidence.Request.Artifact.SHA256
			}
			inputs := map[string]string{}
			for _, src := range p.Inputs {
				if !literalCandidatePath(src.Path) || inputs[src.Path] != "" || (!digestPattern.MatchString(src.SHA256) && src.SHA256 != "absent") {
					return fail("reconciliation_required", "invalid work preimage")
				}
				inputs[src.Path] = src.SHA256
			}
			for path, digest := range required {
				if path != managementLedger && inputs[path] != digest {
					return fail("reconciliation_required", "work coverage preimage differs")
				}
			}
			cards := map[string]bool{}
			for _, c := range selectedSessionStep(i).Cards {
				cards[c.ID] = true
			}
			changed := map[string]bool{}
			size := 0
			for _, ch := range p.Changes {
				size += len(ch.Content)
				if changed[ch.Path] || inputs[ch.Path] != ch.Before || !digestPattern.MatchString(ch.After) || hash([]byte(ch.Content)) != ch.After {
					return fail("reconciliation_required", "invalid work postimage")
				}
				changed[ch.Path] = true
				switch ch.Path {
				case workPath(w):
					if ch.Before != "absent" || ch.Content != string(sessionWorkArtifact(i, w)) {
						return fail("reconciliation_required", "work artifact differs")
					}
				case managementLedger:
				default:
					prefix, ok := strings.CutSuffix(ch.Content, sessionWorkJournal(i, w, ch.Path))
					allowed := ch.Path == i.Context.Binding.Plan.Path || ch.Path == i.Context.Binding.Session.Path
					before := hash([]byte(prefix))
					if strings.HasPrefix(ch.Path, "SDP/KanBan/active/") && required[ch.Path] != "" {
						m, e := metadata([]byte(prefix))
						allowed = e == nil && cards[m["id"]] && m["CardState"] == "in-progress"
						if allowed {
							delete(cards, m["id"])
							if kind == "start" && before != ch.Before {
								b, _ := replaceMetadata([]byte(prefix), "CardState", "ready")
								before = hash(b)
							}
						}
					}
					if !ok || !allowed || before != ch.Before {
						return fail("reconciliation_required", "work journal changes unselected content")
					}
				}
				required[ch.Path] = ch.After
			}
			if size > 4<<20 || len(cards) > 0 || !changed[workPath(w)] || !changed[managementLedger] || !changed[i.Context.Binding.Session.Path] || !changed[i.Context.Binding.Plan.Path] || len(required) != len(p.Coverage) {
				return fail("reconciliation_required", "incomplete work effects")
			}
			pins := map[string]bool{}
			for _, src := range p.Coverage {
				if pins[src.Path] || required[src.Path] == "" || (!digestPattern.MatchString(src.SHA256) && src.SHA256 != "absent") || (src.Path != managementLedger && src.SHA256 != required[src.Path]) {
					return fail("reconciliation_required", "work coverage differs")
				}
				pins[src.Path] = true
			}
			previous = p
		}
		expected := sessionRun(i)
		r, exists := st.Runs[expected.ID]
		if i.Work[0].Effects.Status == "prepared" {
			if exists {
				return fail("reconciliation_required", "run exists before Session start commit")
			}
		} else if !exists || r.RequestID != expected.RequestID || r.InputRevision != expected.InputRevision || r.Selection != expected.Selection || !reflect.DeepEqual(r.Assignment, expected.Assignment) || r.RoutineDigest != expected.RoutineDigest || r.Routine.Digest() != expected.RoutineDigest {
			return fail("reconciliation_required", "Session run binding differs")
		}
	}
	for _, r := range st.Runs {
		if r.Routine.ID == SessionImplementation().ID {
			if _, ok := sessionRunBinding(st, r.ID); !ok {
				return fail("reconciliation_required", "orphan Session run")
			}
		}
	}
	for _, r := range st.Receipts {
		if v := r.Result.Session; v != nil && (v.Operation == "session-start" || v.Operation == "session-progress") && !seen[v.OperationID] {
			return fail("reconciliation_required", "orphan work receipt")
		}
	}
	return nil
}

// Current legacy admission/dispatch also observes persistent Session scope claims.
// This is an optimistic cooperative check, not interception of external writers.
func (s *Service) checkExternalSessionClaims(a Assignment) error {
	entries, _ := os.ReadDir(filepath.Join(s.root, "SDP/Sessions/execution/claims"))
	if len(entries) > 0 {
		for _, p := range append(append([]string{}, a.WritablePaths...), a.ExcludedPaths...) {
			if !literalCandidatePath(p) {
				return fail("conflict", "nonliteral scope cannot bypass Session claims")
			}
		}
	}
	return s.checkSelectionClaims(SessionInstance{}, SessionSelectRequest{}, SessionStepBinding{WritablePaths: a.WritablePaths, ProtectedPaths: a.ExcludedPaths})
}

// DecodeSessionWork keeps each verb's wire fields strict despite shared core storage.
func DecodeSessionWork(b []byte, kind string) (SessionWorkRequest, error) {
	var q SessionWorkRequest
	if e := Decode(b, &q); e != nil {
		return q, e
	}
	var fields map[string]json.RawMessage
	if e := json.Unmarshal(b, &fields); e != nil {
		return q, e
	}
	forbidden := []string{"evidenceId", "expectedRunRevision"}
	if kind == "progress" {
		forbidden = []string{"assessmentId", "outcome", "skills"}
	}
	for _, key := range forbidden {
		if _, exists := fields[key]; exists {
			return q, fail("invalid_request", "field is not supported by Session "+kind+": "+key)
		}
	}
	return q, validSessionWork(kind, q)
}
