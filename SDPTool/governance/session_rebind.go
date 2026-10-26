package governance

import (
	"bytes"
	"crypto/sha256"
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

const sessionRebindStateSchema = "sdp-project-governance-state/0.22"

func sessionTransferProfile(schema string) bool {
	return schema == sessionTransferStateSchema || schema == sessionRebindStateSchema
}

type SessionRebindContextRequest struct {
	InstanceID string `json:"instanceId"`
}
type SessionRebindContext struct {
	Schema           string         `json:"schema"`
	Operation        string         `json:"operation"`
	InstanceID       string         `json:"instanceId"`
	ExpectedRevision uint64         `json:"expectedRevision"`
	OldBindingDigest string         `json:"oldBindingDigest"`
	Next             SessionContext `json:"next"`
}
type SessionRebindRequest struct {
	Schema           string `json:"schema"`
	ProjectRef       string `json:"projectRef"`
	CheckoutID       string `json:"checkoutId"`
	OperationID      string `json:"operationId"`
	InstanceID       string `json:"instanceId"`
	ExpectedRevision uint64 `json:"expectedRevision"`
	InputRevision    string `json:"inputRevision"`
	Review           Source `json:"review"`
	Reason           string `json:"reason"`
}
type SessionRebindReview struct {
	Schema           string `json:"schema"`
	InstanceID       string `json:"instanceId"`
	ExpectedRevision uint64 `json:"expectedRevision"`
	OldBindingDigest string `json:"oldBindingDigest"`
	InputRevision    string `json:"inputRevision"`
	Reviewer         string `json:"reviewer"`
	Verdict          string `json:"verdict"`
	Summary          string `json:"summary"`
}
type SessionRebind struct {
	Request     SessionRebindRequest `json:"request"`
	PriorDigest string               `json:"priorDigest"`
	Next        SessionContext       `json:"next"`
	Effects     SessionRecordEffects `json:"effects"`
	Before      map[string]string    `json:"before"`
	Records     []Source             `json:"records"`
}
type SessionRebindEpoch struct {
	Before     SessionInstance `json:"before"`
	Transition SessionRebind   `json:"transition"`
}
type SessionRebindView struct {
	OperationID      string `json:"operationId"`
	OldBindingDigest string `json:"oldBindingDigest"`
	NewBindingDigest string `json:"newBindingDigest"`
	Review           Source `json:"review"`
	Reason           string `json:"reason"`
}

func originalSessionBinding(i SessionInstance) SessionInstance {
	if len(i.Rebindings) > 0 {
		return i.Rebindings[0].Before
	}
	return i
}
func sessionBindingEpochs(st snapshot) []SessionInstance {
	out := []SessionInstance{}
	for _, i := range st.Sessions {
		for n, epoch := range i.Rebindings {
			old := epoch.Before
			old.Rebindings = i.Rebindings[:n]
			out = append(out, old)
		}
		out = append(out, i)
	}
	return out
}
func rebindBase(i SessionInstance) SessionInstance { i.Rebindings = nil; i.Rebind = nil; return i }
func rebindPath(q SessionRebindRequest) string {
	return "SDP/Sessions/execution/rebindings/" + q.OperationID + ".json"
}
func rebindEligible(i SessionInstance) error {
	if !i.Registered || i.Rebind != nil || i.Advance != nil {
		return fail("conflict", "registered quiescent nonterminal Session required")
	}
	if p := sessionLatestEffects(i); p != nil && p.Status != "committed" {
		return fail("reconciliation_required", "pending Session effects must reconcile before rebind")
	}
	if i.Selection != nil && (i.Disposition == nil || i.Disposition.Effects.Status != "committed") {
		return fail("conflict", "unfinished selected work must be disposed or retired before rebind")
	}
	return nil
}
func completedSessionCards(i SessionInstance) map[string]bool {
	out := map[string]bool{}
	record := func(d *SessionDisposition) {
		if d != nil && d.Effects.Status == "committed" {
			for _, c := range d.Request.Cards {
				if c.Outcome == "completed" {
					out[c.ID] = true
				}
			}
		}
	}
	for _, epoch := range append(append([]SessionRebindEpoch{}, i.Rebindings...), SessionRebindEpoch{Before: i}) {
		for _, p := range epoch.Before.Past {
			record(p.Execution.Disposition)
		}
		record(epoch.Before.Disposition)
	}
	return out
}
func validateRebindBinding(i SessionInstance, next SessionContext) error {
	old, b := i.Context.Binding, next.Binding
	if validateSessionBinding(b) != nil || b.Session != old.Session || b.Plan != old.Plan || b.Workflow != old.Workflow || next.BindingPath != i.Context.BindingPath || next.Operation != "session-context" || next.Schema != SessionSchema || next.ProjectRef != i.Context.ProjectRef || next.CheckoutID != i.Context.CheckoutID || next.InputRevision != sessionInputRevision(next.Sources) {
		return fail("conflict", "rebind must retain Session, plan and workflow identities with exact context")
	}
	byID := map[string]SessionStepBinding{}
	for _, step := range b.Steps {
		byID[step.ID] = step
	}
	for _, oldStep := range old.Steps {
		step, ok := byID[oldStep.ID]
		if !ok {
			return fail("conflict", "removing existing work requires explicit disposition")
		}
		if sessionStepDone(i, oldStep.ID) && !reflect.DeepEqual(step, oldStep) {
			return fail("conflict", "completed step binding is immutable")
		}
		if (oldStep.Review == "independent" && step.Review != "independent") || (oldStep.OwnerAcceptance == "required" && step.OwnerAcceptance != "required") {
			return fail("unauthorized", "rebind cannot weaken required review or owner acceptance")
		}
		for _, p := range oldStep.ProtectedPaths {
			if !slices.Contains(step.ProtectedPaths, p) {
				return fail("unauthorized", "rebind cannot remove protected paths")
			}
		}
		for _, c := range oldStep.Cards {
			if !slices.Contains(step.Cards, c) {
				return fail("conflict", "removing or reclassifying existing card coverage requires disposition")
			}
		}
	}
	pending := false
	for _, step := range b.Steps {
		if !sessionStepDone(i, step.ID) {
			pending = true
		}
	}
	if !pending {
		return fail("conflict", "no remaining step; complete or transfer the disposed Session")
	}
	return nil
}
func rebindPriorSources(i SessionInstance) []Source {
	if p := sessionLatestEffects(i); p != nil {
		return p.Coverage
	}
	out := append([]Source{}, i.Context.Sources...)
	for n, src := range out {
		if src.Path == i.Registration.Path {
			out[n].SHA256 = i.Registration.After
		}
	}
	return out
}
func rebindEditable(i SessionInstance, next SessionContext, before map[string]string) map[string]bool {
	out := map[string]bool{i.Context.BindingPath: true, i.Context.Binding.Session.Path: true, i.Context.Binding.Plan.Path: true, i.Context.Binding.Authority: true, next.Binding.Authority: true}
	completed := completedSessionCards(i)
	for _, src := range next.Sources {
		if strings.HasPrefix(src.Path, "SDP/KanBan/") {
			m, e := metadata([]byte(before[src.Path]))
			if e == nil && !completed[m["id"]] {
				out[src.Path] = true
			}
		}
	}
	return out
}
func ledgerPrefix(raw []byte, digest string) bool {
	h := sha256.New()
	for _, line := range bytes.SplitAfter(raw, []byte("\n")) {
		h.Write(line)
		if fmt.Sprintf("%x", h.Sum(nil)) == digest {
			return true
		}
	}
	return false
}
func retainedJournal(old, current string) bool {
	n := strings.Index(old, "\n## PG ")
	return n < 0 || strings.Contains(current, old[n:])
}
func validateRebindPreimages(i SessionInstance, a SessionRebind) error {
	files := map[string][]byte{}
	lastRecord := ""
	for _, src := range a.Records {
		if !literalCandidatePath(src.Path) || src.Path <= lastRecord || !digestPattern.MatchString(src.SHA256) || hash([]byte(a.Before[src.Path])) != src.SHA256 {
			return fail("reconciliation_required", "invalid canonical rebind record pin")
		}
		lastRecord = src.Path
		files[src.Path] = []byte(a.Before[src.Path])
	}
	if e := validatePublicationHistory(files[managementLedger], files); e != nil {
		return e
	}
	if e := validateRebindBinding(i, a.Next); e != nil {
		return e
	}
	expected := mergeSessionSources(a.Records, mergeSessionSources(rebindPriorSources(i), a.Next.Sources))
	expected = mergeSessionSources(expected, []Source{a.Request.Review})
	names := map[string]bool{}
	for _, src := range expected {
		names[src.Path] = true
		if src.SHA256 != "absent" {
			if _, ok := a.Before[src.Path]; !ok {
				return fail("missing_coverage", "rebind preimage absent")
			}
		}
	}
	for p := range a.Before {
		if !names[p] {
			return fail("reconciliation_required", "unselected rebind preimage")
		}
	}
	last := ""
	pins := map[string]string{}
	for _, src := range a.Next.Sources {
		if !literalCandidatePath(src.Path) || src.Path <= last {
			return fail("reconciliation_required", "invalid rebind context pins")
		}
		last = src.Path
		raw, ok := a.Before[src.Path]
		actual := "absent"
		if ok {
			actual = hash([]byte(raw))
		}
		if actual != src.SHA256 {
			return fail("reconciliation_required", "rebind context source differs")
		}
		pins[src.Path] = src.SHA256
	}
	b := a.Next.Binding
	var typed SessionBinding
	if decode([]byte(a.Before[a.Next.BindingPath]), &typed) != nil || !reflect.DeepEqual(typed, b) || pins[a.Next.BindingPath] != a.Next.BindingDigest {
		return fail("reconciliation_required", "rebind binding bytes differ")
	}
	for _, p := range []string{b.Session.Path, b.Plan.Path, b.Authority, managementLedger} {
		if !digestPattern.MatchString(pins[p]) {
			return fail("missing_coverage", "rebind canonical context absent")
		}
	}
	sm, e := metadata([]byte(a.Before[b.Session.Path]))
	if e != nil {
		return e
	}
	pm, e := metadata([]byte(a.Before[b.Plan.Path]))
	if e != nil {
		return e
	}
	if sm["Session reference"] != b.Session.ID || sm["Status"] != "active" || pm["id"] != b.Plan.ID || !slices.Contains([]string{"active", "planned"}, pm["state"]) || !strings.Contains(a.Before[b.Session.Path], filepath.Base(a.Next.BindingPath)) || strings.TrimSpace(a.Before[b.Authority]) == "" {
		return fail("missing_coverage", "live canonical rebind records required")
	}
	required := map[string]bool{a.Next.BindingPath: true, b.Session.Path: true, b.Plan.Path: true, b.Authority: true, managementLedger: true, "AGENTS.md": true, "SDP/SDP-project.manifest.yaml": true}
	cards := map[string]map[string]string{}
	cardPaths := map[string]string{}
	completed := completedSessionCards(i)
	for p, body := range a.Before {
		if strings.HasPrefix(p, "SDP/KanBan/") {
			m, e := metadata([]byte(body))
			if e == nil && sessionCardPattern.MatchString(m["id"]) {
				if cards[m["id"]] != nil {
					return fail("conflict", "duplicate rebind card")
				}
				cards[m["id"]] = m
				cardPaths[m["id"]] = p
			}
		}
	}
	for _, step := range b.Steps {
		if !sessionAnchor([]byte(a.Before[b.Session.Path]), step.ID) || !sessionAnchor([]byte(a.Before[b.Plan.Path]), step.Milestone) {
			return fail("missing_coverage", "rebind roadmap anchor absent")
		}
		for _, c := range step.Cards {
			m := cards[c.ID]
			required[cardPaths[c.ID]] = true
			if m["PlanId"] != b.Plan.ID || (!slices.Contains([]string{"backlog", "queued", "ready", "in-progress", "gate-review"}, m["CardState"]) && !(completed[c.ID] && m["CardState"] == "completed")) {
				return fail("missing_coverage", "rebind card ownership or state differs")
			}
		}
	}
	if len(pins) != len(required) {
		return fail("missing_coverage", "rebind context source set differs")
	}
	for p := range required {
		if pins[p] == "" {
			return fail("missing_coverage", "required rebind context source absent")
		}
	}
	editable := rebindEditable(i, a.Next, a.Before)
	for _, src := range rebindPriorSources(i) {
		body, ok := a.Before[src.Path]
		actual := "absent"
		if ok {
			actual = hash([]byte(body))
		}
		if src.Path == managementLedger {
			if !ledgerPrefix([]byte(body), src.SHA256) {
				return fail("stale_revision", "rebind cannot rewrite prior management history")
			}
			continue
		}
		if !editable[src.Path] && actual != src.SHA256 {
			return fail("stale_revision", "protected rebind source changed: "+src.Path)
		}
	}
	if !retainedJournal(i.Registration.Content, a.Before[i.Registration.Path]) {
		return fail("stale_revision", "original registration journal changed")
	}
	// Preserve every retained PG journal, including completed-step evidence links.
	for _, epoch := range append(append([]SessionRebindEpoch{}, i.Rebindings...), SessionRebindEpoch{Before: i}) {
		executions := []SessionInstance{epoch.Before}
		for _, p := range epoch.Before.Past {
			executions = append(executions, p.Execution)
		}
		effects := []*SessionRecordEffects{}
		if epoch.Transition.Request.OperationID != "" {
			effects = append(effects, &epoch.Transition.Effects)
		}
		for n := range epoch.Before.Past {
			effects = append(effects, &epoch.Before.Past[n].Transition.Effects)
		}
		for _, old := range executions {
			if old.Selection != nil {
				effects = append(effects, &old.Selection.SessionRecordEffects)
			}
			for n := range old.Assessments {
				effects = append(effects, &old.Assessments[n].Effects)
			}
			for n := range old.Work {
				effects = append(effects, &old.Work[n].Effects)
			}
			for n := range old.OwnerGates {
				effects = append(effects, &old.OwnerGates[n].Effects)
			}
			if old.Disposition != nil {
				effects = append(effects, &old.Disposition.Effects)
			}
			if old.Advance != nil {
				effects = append(effects, &old.Advance.Effects)
			}
		}
		for _, effect := range effects {
			for _, ch := range effect.Changes {
				if body, ok := a.Before[ch.Path]; ok && editable[ch.Path] && !retainedJournal(ch.Content, body) {
					return fail("stale_revision", "retained PG journal edited: "+ch.Path)
				}
			}
		}
	}
	return nil
}
func validRebindRequest(q SessionRebindRequest) error {
	if q.Schema != SessionSchema || !uuidPattern.MatchString(q.InstanceID) || !uuidPattern.MatchString(q.OperationID) || q.ExpectedRevision == 0 || !digestPattern.MatchString(q.InputRevision) || !literalCandidatePath(q.Review.Path) || !digestPattern.MatchString(q.Review.SHA256) || !validLeaderText(q.Reason, 2048) {
		return fail("invalid_request", "exact rebind identity, context, review and reason required")
	}
	return nil
}
func validateRebindReview(i SessionInstance, a SessionRebind) error {
	raw := []byte(a.Before[a.Request.Review.Path])
	var review SessionRebindReview
	if hash(raw) != a.Request.Review.SHA256 || decode(raw, &review) != nil || review.Schema != "sdp-session-rebind-review/0.1" || review.InstanceID != i.ID || review.ExpectedRevision != a.Request.ExpectedRevision || review.OldBindingDigest != i.Context.BindingDigest || review.InputRevision != a.Next.InputRevision || review.Verdict != "passed" || !validLeaderText(review.Reviewer, 256) || !validLeaderText(review.Summary, 4096) {
		return fail("missing_coverage", "exact passing attributed rebind review required")
	}
	return nil
}
func (s *Service) rebindClaims(st snapshot, i SessionInstance, next SessionContext) error {
	if i.Selection != nil {
		if e := s.sessionWorkClaims(st, i); e != nil {
			return e
		}
		if st.Runs[i.Work[0].Request.OperationID].State != "completed" {
			return fail("conflict", "previous selected run not completed")
		}
	}
	for _, step := range next.Binding.Steps {
		if sessionStepDone(i, step.ID) {
			continue
		}
		if e := s.checkSelectionClaims(SessionInstance{}, SessionSelectRequest{}, step); e != nil {
			return e
		}
		for _, r := range st.Runs {
			if r.State != "completed" && r.State != "canceled" && successorScopeConflict(step, r.Assignment) {
				return fail("conflict", "rebound scope overlaps active run")
			}
		}
		for _, p := range st.Prepared {
			if p.State != "retired" && successorScopeConflict(step, p.Assignment) {
				return fail("conflict", "rebound scope overlaps prepared assignment")
			}
		}
	}
	return nil
}
func (s *Service) SessionRebindContext(q SessionRebindContextRequest) (SessionRebindContext, error) {
	st, _, e := s.load()
	if e != nil {
		return SessionRebindContext{}, e
	}
	i, ok := st.Sessions[q.InstanceID]
	if !ok {
		return SessionRebindContext{}, fail("missing_coverage", "Session absent")
	}
	if e = rebindEligible(i); e != nil {
		return SessionRebindContext{}, e
	}
	next, e := s.sessionContextRetaining(i.Context.BindingPath, completedSessionCards(i))
	if e != nil {
		return SessionRebindContext{}, e
	}
	next.CheckoutID = st.CheckoutID
	if e = validateRebindBinding(i, next); e != nil {
		return SessionRebindContext{}, e
	}
	if e = s.rebindClaims(st, i, next); e != nil {
		return SessionRebindContext{}, e
	}
	return SessionRebindContext{SessionSchema, "session-rebind-context", i.ID, i.Revision, i.Context.BindingDigest, next}, nil
}
func rebindJournal(i SessionInstance, a SessionRebind, p string) string {
	rel, _ := filepath.Rel(filepath.Dir(p), rebindPath(a.Request))
	return "\n## PG Session rebind " + a.Request.OperationID + "\n\nExplicit reviewed binding revision; completed steps and earlier evidence retained.\nChanged remaining work requires fresh selection and assessment. Review is controller-attributed; no owner acceptance or native review is inferred.\n[Exact binding revision](" + filepath.ToSlash(rel) + ").\n\n" + publicationProse(a.Request.Reason)
}
func planSessionRebind(i SessionInstance, a SessionRebind) (map[string]string, error) {
	if e := validateRebindPreimages(i, a); e != nil {
		return nil, e
	}
	if e := validateRebindReview(i, a); e != nil {
		return nil, e
	}
	files := map[string]string{}
	for p, b := range a.Before {
		files[p] = b
	}
	subjects := map[string]string{a.Next.Binding.Plan.ID: a.Next.Binding.Plan.Path}
	completed := completedSessionCards(i)
	selected := map[string]bool{}
	for _, step := range a.Next.Binding.Steps {
		for _, c := range step.Cards {
			selected[c.ID] = true
		}
	}
	for p, b := range files {
		if strings.HasPrefix(p, "SDP/KanBan/") {
			m, e := metadata([]byte(b))
			if e == nil && selected[m["id"]] && !completed[m["id"]] {
				subjects[m["id"]] = p
			}
		}
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
		path := subjects[id]
		files[path] = appendRecord([]byte(files[path]), rebindJournal(i, a, path))
		old := latest[id]
		if old == nil {
			return nil, fail("missing_coverage", "rebind canonical history absent")
		}
		prev, ok := old["payload"].(map[string]any)
		if !ok {
			return nil, fail("reconciliation_required", "invalid rebind history")
		}
		payload := map[string]any{}
		for k, v := range prev {
			payload[k] = v
		}
		delete(payload, "lineage")
		payload["previousEventId"] = old["eventId"]
		payload["from"] = prev["to"]
		payload["fromPath"] = prev["toPath"]
		payload["reason"] = "Reviewed Session binding revision " + a.Request.OperationID + "; earlier evidence retained and remaining work requires fresh assessment."
		prefix, kind := "EVT-PM-"+strings.Split(id, "-")[1]+"-", "x-management:updated"
		if strings.HasPrefix(id, "KB-") {
			prefix = "EVT-KB-" + strings.Split(id, "-")[1] + "-"
			kind = "x-kanban:reviewed"
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
	path := a.Next.Binding.Session.Path
	files[path] = appendRecord([]byte(files[path]), rebindJournal(i, a, path))
	artifact, _ := json.MarshalIndent(struct {
		Schema      string               `json:"schema"`
		Request     SessionRebindRequest `json:"request"`
		PriorDigest string               `json:"priorDigest"`
		Next        SessionContext       `json:"next"`
		CreatedAt   string               `json:"createdAt"`
	}{"sdp-session-rebinding/0.1", a.Request, a.PriorDigest, a.Next, a.Effects.CreatedAt}, "", "  ")
	files[rebindPath(a.Request)] = string(artifact) + "\n"
	return files, nil
}
func rebindCoverage(i SessionInstance, a SessionRebind) []Source {
	previous := SessionRecordEffects{Coverage: mergeSessionSources(rebindPriorSources(i), a.Next.Sources)}
	previous.Coverage = mergeSessionSources(previous.Coverage, []Source{a.Request.Review})
	return ownerGateCoverage(previous, SessionOwnerGate{Effects: a.Effects})
}
func (s *Service) buildSessionRebind(i SessionInstance, q SessionRebindRequest, checkout string) (SessionRebind, error) {
	a := SessionRebind{Request: q, PriorDigest: requestDigest(rebindBase(i)), Effects: SessionRecordEffects{Status: "prepared", CreatedAt: time.Now().UTC().Format(time.RFC3339Nano)}, Before: map[string]string{}}
	next, e := s.sessionContextRetaining(i.Context.BindingPath, completedSessionCards(i))
	if e != nil {
		return a, e
	}
	next.CheckoutID = checkout
	a.Next = next
	if next.InputRevision != q.InputRevision {
		return a, fail("stale_revision", "reviewed rebind context changed")
	}
	files, e := s.publicationFiles()
	if e != nil {
		return a, e
	}
	for p, b := range files {
		a.Records = append(a.Records, Source{p, hash(b)})
	}
	slices.SortFunc(a.Records, func(a, b Source) int { return strings.Compare(a.Path, b.Path) })
	sources := mergeSessionSources(a.Records, mergeSessionSources(rebindPriorSources(i), next.Sources))
	sources = mergeSessionSources(sources, []Source{q.Review})
	for _, src := range sources {
		raw, e := readFile(filepath.Join(s.root, src.Path), maxState)
		if os.IsNotExist(e) && src.SHA256 == "absent" {
			continue
		}
		if e != nil {
			return a, e
		}
		a.Before[src.Path] = string(raw)
		files[src.Path] = raw
	}
	if _, e = os.Lstat(filepath.Join(s.root, rebindPath(q))); !os.IsNotExist(e) {
		return a, fail("conflict", "rebind artifact destination occupied")
	}
	after, e := planSessionRebind(i, a)
	if e != nil {
		return a, e
	}
	a.Effects.Changes = dispositionChanges(a.Before, after)
	for _, src := range sources {
		h := "absent"
		if b, ok := a.Before[src.Path]; ok {
			h = hash([]byte(b))
		}
		a.Effects.Inputs = append(a.Effects.Inputs, Source{src.Path, h})
	}
	a.Effects.Inputs = append(a.Effects.Inputs, Source{rebindPath(q), "absent"})
	slices.SortFunc(a.Effects.Inputs, func(a, b Source) int { return strings.Compare(a.Path, b.Path) })
	for p, b := range after {
		files[p] = []byte(b)
	}
	if e = validatePublicationHistory(files[managementLedger], files); e != nil {
		return a, e
	}
	a.Effects.Coverage = rebindCoverage(i, a)
	a.Effects.LedgerBytes = len(files[managementLedger])
	if e = validateRebindEffects(i, a); e != nil {
		return a, e
	}
	return a, nil
}
func (s *Service) RebindSession(c Caller, q SessionRebindRequest) (SessionMutationResult, error) {
	if !c.controller {
		return SessionMutationResult{}, fail("unauthorized", "trusted controller required")
	}
	if e := validRebindRequest(q); e != nil {
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
		return SessionMutationResult{}, fail("stale_revision", "Session revision changed")
	}
	if e = rebindEligible(i); e != nil {
		return SessionMutationResult{}, e
	}
	if i.Revision >= ^uint64(0)-1 || len(i.Rebindings) >= 8 || len(st.Receipts) >= 1024 {
		return SessionMutationResult{}, fail("limit", "rebind capacity reached")
	}
	a, e := s.buildSessionRebind(i, q, st.CheckoutID)
	if e != nil {
		return SessionMutationResult{}, e
	}
	if e = s.rebindClaims(st, i, a.Next); e != nil {
		return SessionMutationResult{}, e
	}
	i.Rebind = &a
	i.Revision++
	st.Sessions[i.ID] = i
	st.Schema = sessionRebindStateSchema
	out := SessionMutationResult{Schema: SessionSchema, Operation: "session-rebind", OperationID: q.OperationID, InstanceID: i.ID, Revision: i.Revision, Registration: "committed", RecordEffects: "prepared"}
	st.Receipts[q.OperationID] = receipt{requestDigest(q), Result{Schema: Schema, Operation: "session.rebind", OperationID: q.OperationID, ProjectRef: s.project, CheckoutID: st.CheckoutID, Session: &out}}
	if e = validateSessionRebinds(st); e != nil {
		return SessionMutationResult{}, e
	}
	if e = s.save(st, head); e != nil {
		return SessionMutationResult{}, e
	}
	return out, nil
}
func (s *Service) applySessionRebind(st *snapshot, i *SessionInstance, tx *projecthistory.Transaction) error {
	a := i.Rebind
	prior := *i
	prior.Rebind = nil
	prior.Revision = a.Request.ExpectedRevision
	if e := rebindEligible(prior); e != nil {
		return e
	}
	if e := s.rebindClaims(*st, prior, a.Next); e != nil {
		return e
	}
	if e := s.applySessionRecordEffects(&a.Effects, tx, s.sessionPendingOwner(*i), "session-rebind"); e != nil {
		return e
	}
	i.Rebindings = append(i.Rebindings, SessionRebindEpoch{rebindBase(prior), *a})
	i.Context = a.Next
	i.Rebind = nil
	i.Past = nil
	i.Selection = nil
	i.Assessments = nil
	i.Work = nil
	i.OwnerGates = nil
	i.Disposition = nil
	i.Advance = nil
	return nil
}
func validateRebindEffects(i SessionInstance, a SessionRebind) error {
	if e := validRebindRequest(a.Request); e != nil {
		return e
	}
	if a.Request.InputRevision != a.Next.InputRevision {
		return fail("reconciliation_required", "rebind reviewed revision differs")
	}
	after, e := planSessionRebind(i, a)
	if e != nil {
		return e
	}
	if !reflect.DeepEqual(a.Effects.Changes, dispositionChanges(a.Before, after)) {
		return fail("reconciliation_required", "rebind postimages differ")
	}
	expected := mergeSessionSources(a.Records, mergeSessionSources(rebindPriorSources(i), a.Next.Sources))
	expected = mergeSessionSources(expected, []Source{a.Request.Review, {rebindPath(a.Request), "absent"}})
	inputs := map[string]string{}
	for _, src := range a.Effects.Inputs {
		if !literalCandidatePath(src.Path) || inputs[src.Path] != "" {
			return fail("reconciliation_required", "invalid rebind pins")
		}
		inputs[src.Path] = src.SHA256
	}
	if len(inputs) != len(expected) {
		return fail("reconciliation_required", "rebind input set differs")
	}
	for _, src := range expected {
		h := "absent"
		if b, ok := a.Before[src.Path]; ok {
			h = hash([]byte(b))
		}
		if inputs[src.Path] != h {
			return fail("reconciliation_required", "rebind input pin differs")
		}
	}
	want := rebindCoverage(i, a)
	if len(want) != len(a.Effects.Coverage) {
		return fail("reconciliation_required", "rebind coverage set differs")
	}
	for n, src := range want {
		got := a.Effects.Coverage[n]
		if got.Path != src.Path || (src.Path != managementLedger && got.SHA256 != src.SHA256) || (!digestPattern.MatchString(got.SHA256) && got.SHA256 != "absent") {
			return fail("reconciliation_required", "rebind coverage differs")
		}
	}
	size := 0
	for _, b := range a.Before {
		size += len(b)
	}
	for _, ch := range a.Effects.Changes {
		size += len(ch.Content)
	}
	if size > 8<<20 || len(inputs) > 2048 || len(want) > 512 || len(a.Effects.Changes) > 128 || a.Effects.LedgerBytes <= 0 {
		return fail("limit", "rebind effect capacity exceeded")
	}
	return nil
}
func validateSessionRebinds(st snapshot) error {
	seen := map[string]bool{}
	for _, root := range st.Sessions {
		if len(root.Rebindings) > 8 || ((len(root.Rebindings) > 0 || root.Rebind != nil) && st.Schema != sessionRebindStateSchema) {
			return fail("incompatible_profile", "bounded binding epochs require profile 0.22")
		}
		validate := func(i SessionInstance, a SessionRebind, committed bool) error {
			q := a.Request
			prior := i
			prior.Rebind = nil
			prior.Revision = q.ExpectedRevision
			if rebindEligible(prior) != nil || a.PriorDigest != requestDigest(rebindBase(prior)) || q.InstanceID != root.ID || q.ProjectRef != st.Project || q.CheckoutID != st.CheckoutID || q.ExpectedRevision >= root.Revision || seen[q.OperationID] {
				return fail("reconciliation_required", "invalid rebind predecessor or identity")
			}
			seen[q.OperationID] = true
			if _, e := time.Parse(time.RFC3339Nano, a.Effects.CreatedAt); e != nil {
				return fail("reconciliation_required", "invalid rebind time")
			}
			if (committed && (a.Effects.Status != "committed" || root.Revision < q.ExpectedRevision+2)) || (!committed && (a.Effects.Status != "prepared" || root.Revision != q.ExpectedRevision+1)) {
				return fail("reconciliation_required", "invalid rebind commitment")
			}
			r, ok := st.Receipts[q.OperationID]
			if !ok || r.Digest != requestDigest(q) || r.Result.Session == nil || r.Result.Session.Operation != "session-rebind" || r.Result.Session.InstanceID != root.ID || r.Result.Session.Revision != q.ExpectedRevision+1 || r.Result.Session.RecordEffects != "prepared" {
				return fail("reconciliation_required", "rebind receipt differs")
			}
			return validateRebindEffects(prior, a)
		}
		var last *SessionRebind
		for n, epoch := range root.Rebindings {
			old := epoch.Before
			if old.ID != root.ID || len(old.Rebindings) != 0 || old.Rebind != nil || old.Revision != epoch.Transition.Request.ExpectedRevision || old.Reason != root.Reason || !reflect.DeepEqual(old.Registration, root.Registration) {
				return fail("reconciliation_required", "invalid flat binding epoch")
			}
			if last != nil && (!reflect.DeepEqual(old.Context, last.Next) || old.Revision < last.Request.ExpectedRevision+2) {
				return fail("reconciliation_required", "binding epoch chain differs")
			}
			old.Rebindings = root.Rebindings[:n]
			if e := validate(old, epoch.Transition, true); e != nil {
				return e
			}
			last = &root.Rebindings[n].Transition
		}
		if last != nil && (!reflect.DeepEqual(root.Context, last.Next) || root.Revision < last.Request.ExpectedRevision+2) {
			return fail("reconciliation_required", "current binding epoch differs")
		}
		if root.Rebind != nil {
			if e := validate(root, *root.Rebind, false); e != nil {
				return e
			}
		}
	}
	for _, r := range st.Receipts {
		if v := r.Result.Session; v != nil && v.Operation == "session-rebind" && !seen[v.OperationID] {
			return fail("reconciliation_required", "orphan rebind receipt")
		}
	}
	return nil
}
func projectSessionRebindings(i SessionInstance, v *SessionView) {
	for _, epoch := range i.Rebindings {
		a := epoch.Transition
		v.BindingRevisions = append(v.BindingRevisions, SessionRebindView{a.Request.OperationID, epoch.Before.Context.BindingDigest, a.Next.BindingDigest, a.Request.Review, a.Request.Reason})
		old := epoch.Before
		executions := []SessionInstance{old}
		for _, p := range old.Past {
			executions = append(executions, p.Execution)
		}
		for _, execution := range executions {
			if execution.Selection == nil || execution.Disposition == nil || execution.Disposition.Effects.Status != "committed" {
				continue
			}
			for n := range v.Steps {
				step := &v.Steps[n]
				if step.ID == execution.Selection.Request.StepID && !step.Selected {
					step.Stage = "done"
					step.Holds = []string{}
					step.RunID = execution.Work[0].Request.OperationID
					step.RunState = "completed"
					step.DispositionID = execution.Disposition.Request.OperationID
					step.CardDispositions = execution.Disposition.Request.Cards
				}
			}
		}
	}
	if i.Rebind != nil {
		v.RecordEffects = "prepared"
		v.NextActions = []string{"session-read", "session-resume"}
		return
	}
	if i.Selection == nil && len(i.Past) == 0 && len(i.Rebindings) > 0 {
		v.RecordEffects = "committed"
	}
	if rebindEligible(i) == nil {
		v.NextActions = append(v.NextActions, "session-rebind-context", "session-rebind")
	}
}
