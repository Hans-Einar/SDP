package governance

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"time"
)

const PublicationSchema = "sdp-governance-publication/1"
const publicationStateSchema = "sdp-project-governance-state/0.12"

type PublicationRequest struct {
	Schema           string `json:"schema"`
	ProjectRef       string `json:"projectRef"`
	CheckoutID       string `json:"checkoutId"`
	OperationID      string `json:"operationId"`
	RunID            string `json:"runId"`
	ExpectedRevision uint64 `json:"expectedRevision"`
	ReviewID         string `json:"reviewId"`
	System           string `json:"system"`
	Summary          string `json:"summary"`
	Remaining        string `json:"remaining"`
}
type PublicationChange struct {
	Path    string `json:"path"`
	Before  string `json:"before"`
	After   string `json:"after"`
	Content string `json:"content"`
}
type Publication struct {
	Request      PublicationRequest  `json:"request"`
	Status       string              `json:"status"`
	CreatedAt    string              `json:"createdAt"`
	PublishedAt  string              `json:"publishedAt"`
	Inputs       []Source            `json:"inputs"`
	Changes      []PublicationChange `json:"changes"`
	EvidencePath string              `json:"evidencePath"`
	EventIDs     []string            `json:"eventIds"`
	Digest       string              `json:"digest"`
}
type PublicationView struct {
	Schema      string      `json:"schema"`
	Publication Publication `json:"publication"`
	Applicable  bool        `json:"applicable"`
	Conditions  []string    `json:"conditions"`
}

func pendingPublication(st snapshot) bool {
	for _, p := range st.Publications {
		if p.Status == "applying" {
			return true
		}
	}
	return false
}
func publicationDigest(p Publication) string {
	p.Status = ""
	p.PublishedAt = ""
	p.Digest = ""
	return requestDigest(p)
}
func validPublicationRequest(q PublicationRequest) bool {
	return q.Schema == PublicationSchema && uuidPattern.MatchString(q.OperationID) && uuidPattern.MatchString(q.RunID) && uuidPattern.MatchString(q.ReviewID) && uuidPattern.MatchString(q.CheckoutID) && q.ExpectedRevision > 0 && regexp.MustCompile(`^[A-Z][A-Z0-9]{0,31}$`).MatchString(q.System) && validLeaderText(q.Summary, 4096) && validLeaderText(q.Remaining, 4096)
}
func (s *Service) publicationEvidence(st snapshot, p Publication) error {
	q := p.Request
	r, ok := st.Runs[q.RunID]
	d, found := st.Reviews[q.ReviewID]
	if !ok || !found || r.Revision != q.ExpectedRevision || r.State != "verifying" || d.RunID != r.ID || d.State != "completed" || d.Outcome == nil || d.Outcome.Verdict != "passed" || latestReport(st, r) != d.Spec.ReportID {
		return fail("stale_revision", "publication requires the selected terminal passing review and unchanged run/report")
	}
	if hasBlockingAttention(r) {
		return fail("attention_required", "blocking questions prevent publication")
	}
	for _, a := range r.Attempts {
		if !slices.Contains([]string{"completed", "failed", "interrupted"}, a.Status) || a.Gap {
			return fail("reconciliation_required", "unresolved implementation attempt")
		}
	}
	for _, i := range r.Inputs {
		if i.State != "acknowledged" {
			return fail("reconciliation_required", "unresolved implementation input")
		}
	}
	for id, review := range st.Reviews {
		if id != d.ID && review.RunID == r.ID && ReviewActive(review.State) {
			return fail("conflict", "another review is active")
		}
		if id != d.ID && review.RunID == r.ID && review.Spec.ReportID == d.Spec.ReportID && review.Spec.Candidate.Digest == d.Spec.Candidate.Digest && review.Outcome != nil && review.Outcome.Verdict != "passed" {
			return fail("evidence_missing", "contradictory review remains for this report/candidate")
		}
	}
	c, e := s.captureCandidate(st.CheckoutID, r, d.Spec.Candidate.Dependencies)
	if e != nil {
		return e
	}
	if c.Digest != d.Spec.Candidate.Digest {
		return fail("stale_candidate", "review candidate changed")
	}
	if e = CheckSkillSnapshots(s.root, d.Spec.Skills); e != nil {
		return e
	}
	for _, id := range d.Spec.EvidenceIDs {
		v := st.Receipts[id].Result.Evidence
		if v == nil || v.Request.Verdict != "passed" || v.Request.Kind != "test" || v.Request.ReportID != d.Spec.ReportID || v.Request.Candidate.Digest != c.Digest {
			return fail("evidence_missing", "passing reviewed test evidence required")
		}
		b, e := readFile(filepath.Join(s.root, v.Request.Artifact.Path), 4<<20)
		if e != nil || hash(b) != v.Request.Artifact.SHA256 {
			return fail("stale_candidate", "test artifact changed")
		}
	}
	// Any current contradictory failed evidence must be reconciled first.
	for _, v := range st.Receipts {
		if a := v.Result.Evidence; a != nil && a.Request.RunID == r.ID && a.Request.ReportID == d.Spec.ReportID && a.Request.Candidate.Digest == c.Digest && a.Request.Kind != "claim" && a.Request.Verdict == "failed" {
			return fail("evidence_missing", "contradictory failed evidence remains for this report/candidate")
		}
	}
	return nil
}
func (s *Service) PreparePublication(c Caller, q PublicationRequest) (PublicationView, error) {
	if !c.controller {
		return PublicationView{}, fail("unauthorized", "publication requires trusted controller")
	}
	if !validPublicationRequest(q) {
		return PublicationView{}, fail("invalid_request", "exact review/run and bounded summary, remaining work and system required")
	}
	unlock, e := lock(s.state)
	if e != nil {
		return PublicationView{}, e
	}
	defer unlock()
	st, head, e := s.load()
	if e != nil {
		return PublicationView{}, e
	}
	if old, ok := st.Publications[q.OperationID]; ok {
		if requestDigest(q) != requestDigest(old.Request) {
			return PublicationView{}, fail("idempotency_conflict", "publication request changed")
		}
		return s.publicationView(st, old), nil
	}
	if e = s.authorize(c, Schema, q.ProjectRef, q.CheckoutID, q.OperationID, st); e != nil {
		return PublicationView{}, e
	}
	if _, ok, e := retry(st, q.OperationID, requestDigest(q)); ok {
		if e != nil {
			return PublicationView{}, e
		}
		return PublicationView{}, fail("idempotency_conflict", "operation belongs to another service")
	}
	if len(st.Publications) >= 16 || len(st.Receipts) >= 1024 {
		return PublicationView{}, fail("limit", "publication capacity reached")
	}
	for _, p := range st.Publications {
		if p.Request.RunID == q.RunID && p.Status != "abandoned" {
			return PublicationView{}, fail("conflict", "run already has a publication; inspect or abandon the unsent proposal")
		}
	}
	p := Publication{Request: q, Status: "prepared", CreatedAt: time.Now().UTC().Format(time.RFC3339Nano)}
	if e = s.publicationEvidence(st, p); e != nil {
		return PublicationView{}, e
	}
	r := st.Runs[q.RunID]
	if e = s.reviewBasis(st, r, st.Reviews[q.ReviewID].Spec); e != nil {
		return PublicationView{}, e
	}
	if p, e = s.buildPublication(st, p); e != nil {
		return PublicationView{}, e
	}
	if e = s.publicationInputs(p, false); e != nil {
		return PublicationView{}, e
	}
	if e = s.publicationEvidence(st, p); e != nil {
		return PublicationView{}, e
	}
	p.Digest = publicationDigest(p)
	if st.Schema != sessionStateSchema {
		st.Schema = publicationStateSchema
	}
	if st.Publications == nil {
		st.Publications = map[string]Publication{}
	}
	st.Publications[q.OperationID] = p
	st.Receipts[q.OperationID] = receipt{Digest: requestDigest(q), Result: Result{Schema: Schema, Operation: "publication.prepare", PublicationID: q.OperationID, OperationID: q.OperationID, ProjectRef: s.project, CheckoutID: st.CheckoutID, Run: summarizeRun(r), Coverage: "known"}}
	if e = s.save(st, head); e != nil {
		return PublicationView{}, e
	}
	return s.publicationView(st, p), nil
}
func (s *Service) publicationInputs(p Publication, allowAfter bool) error {
	changes := map[string]PublicationChange{}
	for _, c := range p.Changes {
		changes[c.Path] = c
	}
	for _, src := range p.Inputs {
		b, e := readFile(filepath.Join(s.root, src.Path), maxState)
		actual := "absent"
		if e == nil {
			actual = hash(b)
		} else if !os.IsNotExist(e) {
			return e
		}
		if actual == src.SHA256 {
			continue
		}
		if c, ok := changes[src.Path]; ok && allowAfter && actual == c.After {
			continue
		}
		return fail("reconciliation_required", "publication input changed: "+src.Path)
	}
	return nil
}
func (s *Service) publicationView(st snapshot, p Publication) PublicationView {
	v := PublicationView{Schema: PublicationSchema, Publication: p, Conditions: []string{}}
	if e := s.publicationInputs(p, p.Status != "prepared"); e != nil {
		v.Conditions = append(v.Conditions, e.Error())
	}
	if e := s.publicationEvidence(st, p); e != nil {
		v.Conditions = append(v.Conditions, e.Error())
	}
	if p.Status == "abandoned" {
		v.Conditions = append(v.Conditions, "publication proposal abandoned")
	}
	if p.Status == "published" {
		for _, c := range p.Changes {
			b, e := readFile(filepath.Join(s.root, c.Path), maxState)
			if e != nil || hash(b) != c.After {
				v.Conditions = append(v.Conditions, "published record changed: "+c.Path)
			}
		}
	}
	v.Applicable = len(v.Conditions) == 0
	return v
}
func (s *Service) PublicationRead(id string) (PublicationView, error) {
	st, head, e := s.load()
	if e != nil {
		return PublicationView{}, e
	}
	p, ok := st.Publications[id]
	if !ok || !uuidPattern.MatchString(id) {
		return PublicationView{}, fail("missing_coverage", "publication absent")
	}
	v := s.publicationView(st, p)
	_, latest, e := s.load()
	if e != nil {
		return PublicationView{}, e
	}
	if head != latest {
		return PublicationView{}, fail("stale_revision", "publication changed during inspection")
	}
	return v, nil
}

// ProjectPublicationLock is shared by all PG state stores for this canonical root.
// The durable intent stays in the operational store, never in the project files.
func (s *Service) publicationLockDir() (string, error) {
	home := os.Getenv("XDG_STATE_HOME")
	if home == "" {
		v, e := os.UserHomeDir()
		if e != nil {
			return "", e
		}
		home = filepath.Join(v, ".local", "state")
	}
	if !filepath.IsAbs(home) {
		return "", fail("path", "absolute XDG state home required")
	}
	return filepath.Join(home, "sdptool", "project-publication", hash([]byte(s.root))), nil
}
func (s *Service) ApplyPublication(c Caller, id string) (PublicationView, error) {
	if !c.controller {
		return PublicationView{}, fail("unauthorized", "publication requires trusted controller")
	}
	if !uuidPattern.MatchString(id) {
		return PublicationView{}, fail("invalid_request", "publication UUID required")
	}
	unlock, e := lock(s.state)
	if e != nil {
		return PublicationView{}, e
	}
	defer unlock()
	st, head, e := s.load()
	if e != nil {
		return PublicationView{}, e
	}
	p, ok := st.Publications[id]
	if !ok {
		return PublicationView{}, fail("missing_coverage", "publication absent")
	}
	dir, e := s.publicationLockDir()
	if e != nil {
		return PublicationView{}, e
	}
	unlockProject, e := lock(dir)
	if e != nil {
		return PublicationView{}, e
	}
	defer unlockProject()
	marker := filepath.Join(dir, "pending")
	identity := s.state + "\n" + id + "\n"
	held, e := readFile(marker, 8192)
	if e != nil && !os.IsNotExist(e) {
		return PublicationView{}, e
	}
	if e == nil && string(held) != identity {
		return PublicationView{}, fail("reconciliation_required", "another PG store has pending project publication")
	}
	if p.Status == "published" {
		if string(held) == identity {
			if e = os.Remove(marker); e != nil {
				return PublicationView{}, e
			}
			if e = syncDir(dir); e != nil {
				return PublicationView{}, e
			}
		}
		return s.publicationView(st, p), nil
	}
	if p.Status == "abandoned" {
		return PublicationView{}, fail("conflict", "publication abandoned")
	}
	if e = s.publicationInputs(p, p.Status == "applying"); e != nil {
		return PublicationView{}, e
	}
	if e = s.publicationEvidence(st, p); e != nil {
		return PublicationView{}, e
	}
	if e = s.validatePublicationRecords(p); e != nil {
		return PublicationView{}, e
	}
	// Reserve enough room for the final timestamp before any project file effects.
	encoded, e := json.Marshal(st)
	if e != nil {
		return PublicationView{}, e
	}
	if len(encoded) > maxState-256 {
		return PublicationView{}, fail("limit", "insufficient state capacity for publication receipt")
	}
	// Claim project scope before making intent active. A lost response may leave the
	// marker, but this exact proposal can always reclaim it; another store cannot.
	if held == nil {
		if e = writeNew(marker, []byte(identity)); e != nil {
			return PublicationView{}, e
		}
		if e = syncDir(dir); e != nil {
			return PublicationView{}, e
		}
	}
	if p.Status == "prepared" {
		p.Status = "applying"
		st.Publications[id] = p
		if e = s.save(st, head); e != nil {
			return PublicationView{}, e
		}
		st, head, e = s.load()
		if e != nil {
			return PublicationView{}, e
		}
	}
	if s.fault != nil {
		if e = s.fault("publication-after-intent"); e != nil {
			return PublicationView{}, e
		}
	}
	for i, change := range p.Changes {
		if e = s.publicationInputs(p, true); e != nil {
			return PublicationView{}, e
		}
		if e = s.publicationEvidence(st, p); e != nil {
			return PublicationView{}, e
		}
		if e = s.applyPublicationFile(change); e != nil {
			return PublicationView{}, e
		}
		if s.fault != nil {
			if e = s.fault(fmt.Sprintf("publication-after-file:%d", i)); e != nil {
				return PublicationView{}, e
			}
		}
	}
	if e = s.publicationInputs(p, true); e != nil {
		return PublicationView{}, e
	}
	if e = s.publicationEvidence(st, p); e != nil {
		return PublicationView{}, e
	}
	if e = s.validatePublicationRecords(p); e != nil {
		return PublicationView{}, e
	}
	for _, ch := range p.Changes {
		b, e := readFile(filepath.Join(s.root, ch.Path), maxState)
		if e != nil || hash(b) != ch.After {
			return PublicationView{}, fail("reconciliation_required", "publication postimage differs: "+ch.Path)
		}
	}
	p.Status = "published"
	p.PublishedAt = time.Now().UTC().Format(time.RFC3339Nano)
	st.Publications[id] = p
	if e = s.save(st, head); e != nil {
		return PublicationView{}, e
	}
	if e = os.Remove(marker); e != nil {
		return PublicationView{}, e
	}
	if e = syncDir(dir); e != nil {
		return PublicationView{}, e
	}
	return s.publicationView(st, p), nil
}
func (s *Service) AbandonPublication(c Caller, id string) (PublicationView, error) {
	if !c.controller {
		return PublicationView{}, fail("unauthorized", "controller required")
	}
	unlock, e := lock(s.state)
	if e != nil {
		return PublicationView{}, e
	}
	defer unlock()
	st, head, e := s.load()
	if e != nil {
		return PublicationView{}, e
	}
	p, ok := st.Publications[id]
	if !ok {
		return PublicationView{}, fail("missing_coverage", "publication absent")
	}
	if p.Status != "prepared" && p.Status != "abandoned" {
		return PublicationView{}, fail("reconciliation_required", "only an unapplied proposal can be abandoned")
	}
	dir, e := s.publicationLockDir()
	if e != nil {
		return PublicationView{}, e
	}
	u, e := lock(dir)
	if e != nil {
		return PublicationView{}, e
	}
	defer u()
	b, e := readFile(filepath.Join(dir, "pending"), 8192)
	if e == nil && string(b) == s.state+"\n"+id+"\n" {
		// Durable state is still prepared/abandoned, so no project file effect
		// could have started. Release an orphan claim from a failed intent save.
		if e = os.Remove(filepath.Join(dir, "pending")); e != nil {
			return PublicationView{}, e
		}
		if e = syncDir(dir); e != nil {
			return PublicationView{}, e
		}
	}
	if e != nil && !os.IsNotExist(e) {
		return PublicationView{}, e
	}
	if p.Status == "prepared" {
		p.Status = "abandoned"
		st.Publications[id] = p
		if e = s.save(st, head); e != nil {
			return PublicationView{}, e
		}
	}
	return s.publicationView(st, p), nil
}
func (s *Service) applyPublicationFile(c PublicationChange) error {
	path := filepath.Join(s.root, c.Path)
	b, e := readFile(path, maxState)
	actual := "absent"
	if e == nil {
		actual = hash(b)
	} else if !os.IsNotExist(e) {
		return e
	}
	if actual == c.After {
		return nil
	}
	if actual != c.Before {
		return fail("reconciliation_required", "record preimage conflict: "+c.Path)
	}
	if e = mkdirDurable(filepath.Dir(path)); e != nil {
		return e
	}
	f, e := os.CreateTemp(filepath.Dir(path), ".pg-publication-")
	if e != nil {
		return e
	}
	name := f.Name()
	defer os.Remove(name)
	mode := os.FileMode(0600)
	if actual != "absent" {
		info, e := os.Stat(path)
		if e != nil {
			f.Close()
			return e
		}
		mode = info.Mode().Perm()
	}
	if e = f.Chmod(mode); e == nil {
		_, e = f.WriteString(c.Content)
	}
	if e == nil {
		e = f.Sync()
	}
	closeErr := f.Close()
	if e != nil {
		return e
	}
	if closeErr != nil {
		return closeErr
	}
	// Recheck immediately before rename. External editors do not participate in
	// the project lock, so this is optimistic exclusion, not filesystem CAS.
	b, e = readFile(path, maxState)
	now := "absent"
	if e == nil {
		now = hash(b)
	} else if !os.IsNotExist(e) {
		return e
	}
	if now != c.Before {
		return fail("reconciliation_required", "record changed before replacement: "+c.Path)
	}
	if c.Before == "absent" {
		if e = os.Link(name, path); e != nil {
			return e
		}
	} else {
		if e = os.Rename(name, path); e != nil {
			return e
		}
	}
	return syncDir(filepath.Dir(path))
}
func validatePublicationState(st snapshot) error {
	if !publicationProfile(st.Schema) {
		if len(st.Publications) > 0 {
			return fail("incompatible_profile", "publication requires profile 0.12")
		}
		return nil
	}
	if (st.Schema == publicationStateSchema && len(st.Publications) == 0) || len(st.Publications) > 16 {
		return fail("reconciliation_required", "invalid publication collection")
	}
	active := 0
	for id, p := range st.Publications {
		q := p.Request
		r, ok := st.Runs[q.RunID]
		receipt, exists := st.Receipts[id]
		_, timeErr := time.Parse(time.RFC3339Nano, p.CreatedAt)
		if !validPublicationRequest(q) || id != q.OperationID || !ok || !exists || receipt.Digest != requestDigest(q) || receipt.Result.Operation != "publication.prepare" || receipt.Result.PublicationID != id || receipt.Result.OperationID != id || receipt.Result.Schema != Schema || receipt.Result.ProjectRef != st.Project || receipt.Result.CheckoutID != st.CheckoutID || receipt.Result.Run.ID != q.RunID || receipt.Result.Run.Revision != q.ExpectedRevision || q.CheckoutID != st.CheckoutID || q.ProjectRef != st.Project || r.Revision < q.ExpectedRevision || p.Digest != publicationDigest(p) || timeErr != nil || !slices.Contains([]string{"prepared", "applying", "published", "abandoned"}, p.Status) || len(p.Changes) != 6 || len(p.Inputs) > 2048 || len(p.EventIDs) != 3 {
			return fail("reconciliation_required", "invalid publication intent")
		}
		if p.Status == "applying" {
			active++
		}
		if p.Status == "published" {
			if _, e := time.Parse(time.RFC3339Nano, p.PublishedAt); e != nil {
				return fail("reconciliation_required", "invalid publication timestamp")
			}
		} else if p.PublishedAt != "" {
			return fail("reconciliation_required", "unpublished timestamp")
		}
		seen := map[string]string{}
		for _, src := range p.Inputs {
			if !literalCandidatePath(src.Path) || seen[src.Path] != "" || src.SHA256 != "absent" && !digestPattern.MatchString(src.SHA256) {
				return fail("reconciliation_required", "invalid publication pins")
			}
			seen[src.Path] = src.SHA256
		}
		paths := map[string]bool{}
		size := 0
		for _, ch := range p.Changes {
			size += len(ch.Content)
			if paths[ch.Path] || seen[ch.Path] != ch.Before || ch.After != hash([]byte(ch.Content)) || ch.After == ch.Before {
				return fail("reconciliation_required", "invalid publication change")
			}
			paths[ch.Path] = true
		}
		if size > 4<<20 {
			return fail("reconciliation_required", "publication exceeds bound")
		}
	}
	if active > 1 {
		return fail("reconciliation_required", "multiple applying publications")
	}
	return nil
}

func publicationProfile(schema string) bool {
	return schema == publicationStateSchema || schema == sessionStateSchema
}
