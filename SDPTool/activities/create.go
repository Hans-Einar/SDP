package activities

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/Hans-Einar/SDP/SDPTool/projecthistory"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//go:embed activity-readme.md
var activityReadme []byte

func (s *Service) validateCaller(c Caller) error {
	if !namespaceRE.MatchString(s.ProjectID) || strings.TrimSpace(c.Actor) == "" || c.Authority.Role != "project-owner" {
		return fmt.Errorf("trusted project-owner authority required")
	}
	if e := portable(c.Authority.Source); e != nil {
		return e
	}
	_, e := s.read(c.Authority.Source)
	return e
}
func (s *Service) begin(id, expected string, request any, caller Caller) (*projecthistory.Transaction, string, *Result, error) {
	if e := s.validateCaller(caller); e != nil {
		return nil, "", nil, e
	}
	if !uuidRE.MatchString(id) || !hashRE.MatchString(expected) {
		return nil, "", nil, fmt.Errorf("operation UUID and expected revision required")
	}
	b, _ := json.Marshal(struct {
		Request any
		Caller  Caller
	}{request, caller})
	dg := digest(b)
	p, e := s.path("ProjectManagement/Ledger.ndjson")
	if e != nil {
		return nil, "", nil, e
	}
	tx, e := projecthistory.BeginTransaction(p)
	if e != nil {
		return nil, "", nil, e
	}
	in, e := s.load(id)
	if e == nil {
		if in.Digest != dg {
			tx.Close()
			return nil, "", nil, projecthistory.ErrConflict
		}
		r, e := s.apply(tx, in)
		tx.Close()
		return nil, dg, &r, e
	}
	if !os.IsNotExist(e) {
		tx.Close()
		return nil, "", nil, e
	}
	snap, e := tx.Read("activity:" + id)
	if e == nil {
		view, ve := s.ReadAdmission(tx, "activity:"+id)
		if ve != nil {
			e = ve
		} else if view.Pending != nil {
			e = fmt.Errorf("Activity lifecycle request pending")
		}
	}
	if e == nil && snap.Revision != expected {
		e = projecthistory.ErrStale
	}
	if e != nil {
		tx.Close()
		return nil, "", nil, e
	}
	return tx, dg, nil, nil
}
func slug(s string) string {
	s = regexp.MustCompile(`[^A-Za-z0-9]+`).ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		s = "Untitled"
	}
	if len(s) > 80 {
		s = s[:80]
	}
	return s
}
func (s *Service) allocate(kind string, c Catalog) (int, error) {
	max := 0
	re := activityRE
	if kind == "SESSION" {
		re = sessionRE
	}
	consume := func(id string) {
		m := re.FindStringSubmatch(id)
		if len(m) > 0 && m[1] == s.ProjectID {
			n, _ := strconv.Atoi(m[2])
			if n > max {
				max = n
			}
		}
	}
	for _, a := range c.Activities {
		consume(a.Descriptor.ID)
		for _, r := range a.Descriptor.References {
			consume(r.ID)
		}
		for _, v := range a.Sessions {
			consume(v.ID)
		}
	}
	e := filepath.WalkDir(s.Area, func(p string, d os.DirEntry, e error) error {
		if e != nil {
			return e
		}
		rel, e := filepath.Rel(s.Area, p)
		if e != nil {
			return e
		}
		rel = filepath.ToSlash(rel)
		if rel == "." {
			return nil
		}
		if d.Type()&os.ModeSymlink != 0 {
			if strings.HasPrefix(rel, "Activities/") || strings.HasPrefix(rel, "Sessions/") {
				return fmt.Errorf("symlink in identity allocation")
			}
			return nil
		}
		if d.IsDir() {
			if strings.HasPrefix(d.Name(), ".") {
				return filepath.SkipDir
			}
			if kind == "ACTIVITY" && strings.HasPrefix(d.Name(), "Activity-") {
				v := strings.Split(strings.TrimPrefix(d.Name(), "Activity-"), "--")[0]
				n, _ := strconv.Atoi(v)
				if n > max {
					max = n
				}
			}
			return nil
		}
		if kind == "SESSION" && strings.HasSuffix(strings.ToLower(d.Name()), ".md") {
			b, e := s.read(rel)
			if e != nil {
				return e
			}
			for _, line := range strings.Split(string(b), "\n") {
				cells := strings.Split(line, "|")
				if len(cells) >= 4 {
					key := strings.ToLower(strings.TrimSpace(cells[1]))
					if key == "id" || key == "session reference" {
						consume(strings.TrimSpace(cells[2]))
					}
				}
			}
		}
		if kind == "SESSION" {
			m := regexp.MustCompile(`^session-#([0-9]+)--`).FindStringSubmatch(d.Name())
			if len(m) > 0 {
				n, _ := strconv.Atoi(m[1])
				if n > max {
					max = n
				}
			}
		}
		return nil
	})
	if e != nil {
		return 0, e
	}
	if max >= 99999999 {
		return 0, fmt.Errorf("identity capacity exceeded")
	}
	return max + 1, nil
}
func (s *Service) event(id, kind, subject, reason string, caller Caller, dg string, ref DescriptorReference, prev *string, from *string, to string, session *Session) Event {
	now := s.Now
	if now == nil {
		now = time.Now
	}
	t := now().UTC().Format(time.RFC3339Nano)
	return Event{SchemaVersion: "1.0", EventID: "EVT-ACTIVITY-" + s.ProjectID + "-" + strings.ToUpper(strings.ReplaceAll(id, "-", "")), EventType: kind, OccurredAt: t, Actor: caller.Actor, SubjectID: subject, Payload: Payload{SchemaVersion: "0.2", ProjectID: s.ProjectID, OperationID: id, RecordedAt: t, Authority: caller.Authority, PreviousEventID: prev, From: from, To: to, Reason: reason, Descriptor: ref, Links: []string{}, RequestDigest: dg, Session: session}}
}
func (s *Service) publish(tx *projecthistory.Transaction, in intent) (Result, error) {
	snap, e := tx.Read("activity:" + in.OperationID)
	if e != nil {
		return Result{}, e
	}
	next := append([]byte{}, snap.Bytes...)
	for _, ev := range in.Events {
		next = append(next, ev...)
		next = append(next, '\n')
	}
	in.Result.HistoryRevision = digest(next)
	raw, e := json.Marshal(in)
	if len(raw) > projecthistory.MaxBytes {
		return Result{}, fmt.Errorf("operation intent capacity exceeded")
	}
	if e != nil {
		return Result{}, e
	}
	if e = s.write(s.intentPath(in.OperationID), raw); e != nil {
		return Result{}, e
	}
	if e = s.point("after-intent"); e != nil {
		return Result{}, e
	}
	return s.apply(tx, in)
}
func makeIntent(project, id, dg, expected string, files []image, event Event, r Result) intent {
	b, _ := json.Marshal(event)
	r.EventID = event.EventID
	return intent{Schema: "sdp-activity-operation/0.2", ProjectID: project, OperationID: id, Digest: dg, Expected: expected, Files: files, Events: [][]byte{b}, Result: r}
}
func (s *Service) Create(req CreateRequest, caller Caller) (Result, error) {
	tx, dg, prior, e := s.begin(req.OperationID, req.ExpectedHistoryRevision, req, caller)
	if e != nil {
		return Result{}, e
	}
	if prior != nil {
		return *prior, nil
	}
	defer tx.Close()
	if strings.TrimSpace(req.Title) == "" || strings.TrimSpace(req.Owner) == "" || strings.TrimSpace(req.Mandate) == "" || strings.TrimSpace(req.Reason) == "" {
		return Result{}, fmt.Errorf("title, owner, mandate and reason required")
	}
	for _, scope := range req.Scope {
		if _, e = s.path(scope); e != nil {
			return Result{}, e
		}
	}
	if e = s.validateCards(req.SourceCards); e != nil {
		return Result{}, e
	}
	if e = validateRefs(s, req.References); e != nil {
		return Result{}, e
	}
	snap, e := tx.Read("activity:" + req.OperationID)
	if e != nil {
		return Result{}, e
	}
	c, e := Replay(s.Area, s.ProjectID, snap)
	if e != nil {
		return Result{}, e
	}
	n, e := s.allocate("ACTIVITY", c)
	if e != nil {
		return Result{}, e
	}
	id := fmt.Sprintf("ACTIVITY-%s-%04d", s.ProjectID, n)
	dir := fmt.Sprintf("Activities/Activity-%04d--%s", n, slug(req.Title))
	cards := req.SourceCards
	if cards == nil {
		cards = []string{}
	}
	refs := req.References
	if refs == nil {
		refs = []Reference{}
	}
	d := Descriptor{Schema: "sdp-activity/0.2", ID: id, ProjectID: s.ProjectID, Title: req.Title, Owner: req.Owner, Mandate: "01-Mandate/Mandate.md", SourceCards: cards, Systems: req.Systems, References: refs, Scope: req.Scope, Dependencies: req.Dependencies}
	if e = validateDependencies(c, d); e != nil {
		return Result{}, e
	}
	b, e := yaml.Marshal(d)
	if e != nil {
		return Result{}, e
	}
	if _, e = parseDescriptor(b); e != nil {
		return Result{}, e
	}
	ref := DescriptorReference{Path: dir + "/activity.yaml", SHA256: digest(b)}
	ev := s.event(req.OperationID, "x-activity:created", id, req.Reason, caller, dg, ref, nil, nil, "planned", nil)
	files := []image{{Path: "Activities/.descriptors/" + ref.SHA256 + ".yaml", After: b}, {Path: ref.Path, After: b}, {Path: dir + "/01-Mandate/Mandate.md", After: []byte(req.Mandate + "\n")}, {Path: dir + "/README.md", After: activityReadme}}
	for _, phase := range []string{"02-Study", "03-Requirements", "04-Design", "05-Implementation", "06-Verification"} {
		files = append(files, image{Path: dir + "/" + phase + "/README.md", After: []byte("# " + phase + "\n\nThis phase contains the Activity's relevant records and references. Readiness follows evidenced outputs and disposition.\n")})
	}
	files = append(files, image{Path: dir + "/Sessions/README.md", After: []byte("# Sessions\n\nActivity Sessions share this home. Creating a Session does not dispatch work or establish an execution binding.\n")})
	return s.publish(tx, makeIntent(s.ProjectID, req.OperationID, dg, snap.Revision, files, ev, Result{OperationID: req.OperationID, ActivityID: id, Path: dir, DescriptorRevision: ref.SHA256}))
}
func (s *Service) CreateSession(req CreateSessionRequest, caller Caller) (Result, error) {
	tx, dg, prior, e := s.begin(req.OperationID, req.ExpectedHistoryRevision, req, caller)
	if e != nil {
		return Result{}, e
	}
	if prior != nil {
		return *prior, nil
	}
	defer tx.Close()
	if strings.TrimSpace(req.Title) == "" || strings.TrimSpace(req.Objective) == "" || strings.TrimSpace(req.Reason) == "" {
		return Result{}, fmt.Errorf("title, objective and reason required")
	}
	phases, e := validateSessionPhases(req.Phase, req.PhaseContributions)
	if e != nil {
		return Result{}, e
	}
	refs := append([]Reference{}, req.MilestoneRefs...)
	for _, ref := range req.MilestoneRefs {
		if ref.Role != "milestone" {
			return Result{}, fmt.Errorf("milestoneRefs require milestone role")
		}
	}
	if req.Plan != nil {
		if req.Plan.Role != "plan" && req.Plan.Role != "implementation-plan" {
			return Result{}, fmt.Errorf("Session plan requires plan reference role")
		}
		refs = append(refs, *req.Plan)
	}
	if req.Predecessor != nil {
		if req.Predecessor.Role != "session" && req.Predecessor.Role != "predecessor" {
			return Result{}, fmt.Errorf("predecessor requires Session reference")
		}
		refs = append(refs, *req.Predecessor)
	}
	if req.Handoff != nil {
		refs = append(refs, *req.Handoff)
	}
	for _, ref := range refs {
		if ref.Availability == "unavailable" {
			return Result{}, fmt.Errorf("Session links require available records")
		}
	}
	if e = validateRefs(s, refs); e != nil {
		return Result{}, e
	}
	snap, e := tx.Read("activity:" + req.OperationID)
	if e != nil {
		return Result{}, e
	}
	c, e := Replay(s.Area, s.ProjectID, snap)
	if e != nil {
		return Result{}, e
	}
	a, ok := c.Find(req.ActivityID)
	if !ok {
		return Result{}, fmt.Errorf("unknown Activity")
	}
	if a.LastEventID != req.ExpectedActivityEventID || a.DescriptorRef.SHA256 != req.ExpectedDescriptorRevision {
		return Result{}, projecthistory.ErrStale
	}
	if a.State == "completed" || a.State == "canceled" {
		return Result{}, fmt.Errorf("Activity is closed")
	}
	n, e := s.allocate("SESSION", c)
	if e != nil {
		return Result{}, e
	}
	id := fmt.Sprintf("SESSION-%s-%04d", s.ProjectID, n)
	p := a.Path + fmt.Sprintf("/Sessions/session-#%04d--%s.md", n, slug(req.Title))
	for _, ref := range refs {
		if ref.ID == id {
			return Result{}, fmt.Errorf("Session cannot refer to its allocated identity")
		}
		if strings.HasPrefix(ref.ID, "SESSION-") && ref.Path == p {
			return Result{}, fmt.Errorf("Session cannot be its own predecessor or handoff")
		}
	}
	primaryPhase := req.Phase
	if primaryPhase == "" {
		primaryPhase = phases[0]
	}
	text := fmt.Sprintf("# %s\n\n| Field | Value |\n| --- | --- |\n| Session reference | %s |\n| Status | planned |\n| Activity | %s |\n| Phase | %s |\n| Phase contributions | %s |\n", req.Title, id, a.Descriptor.ID, primaryPhase, strings.Join(phases, ", "))
	if req.Plan != nil {
		text += "| Plan | " + req.Plan.ID + " |\n"
	}
	if req.Predecessor != nil {
		text += "| Predecessor | " + req.Predecessor.ID + " |\n"
	}
	if req.Handoff != nil {
		text += "| Handoff | " + req.Handoff.ID + " |\n"
	}
	if len(req.MilestoneRefs) > 0 {
		ids := []string{}
		for _, ref := range req.MilestoneRefs {
			ids = append(ids, ref.ID)
		}
		text += "| Milestones | " + strings.Join(ids, ", ") + " |\n"
	}
	text += "\n## Goal\n\n" + req.Objective + "\n"
	if len(refs) > 0 {
		text += "\n## Linked records\n\n"
		for _, ref := range refs {
			rel, _ := filepath.Rel(filepath.Dir(filepath.FromSlash(p)), filepath.FromSlash(ref.Path))
			text += "- " + ref.Role + ": [" + ref.ID + "](" + filepath.ToSlash(rel) + ")\n"
		}
	}
	b := []byte(text)
	session := Session{ID: id, Path: p, SHA256: digest(b)}
	ev := s.event(req.OperationID, "x-activity:session-created", a.Descriptor.ID, req.Reason, caller, dg, a.DescriptorRef, &a.LastEventID, &a.State, a.State, &session)
	old, e := s.read(a.DescriptorRef.Path)
	if e != nil || digest(old) != a.DescriptorRef.SHA256 {
		return Result{}, fmt.Errorf("live descriptor differs from current pin")
	}
	files := []image{{Path: "Activities/.descriptors/" + a.DescriptorRef.SHA256 + ".yaml", After: old}, {Path: "Activities/.sessions/" + session.SHA256 + ".md", After: b}, {Path: p, After: b}}
	return s.publish(tx, makeIntent(s.ProjectID, req.OperationID, dg, snap.Revision, files, ev, Result{OperationID: req.OperationID, ActivityID: a.Descriptor.ID, SessionID: id, Path: p, DescriptorRevision: a.DescriptorRef.SHA256}))
}

func validateSessionPhases(primary string, contributions []string) ([]string, error) {
	phases := append([]string{}, contributions...)
	if len(phases) == 0 && primary != "" {
		phases = []string{primary}
	}
	if len(phases) == 0 {
		return nil, fmt.Errorf("Session phase or phaseContributions required")
	}
	seen := map[string]bool{}
	for _, phase := range phases {
		valid := false
		for _, allowed := range []string{"01-Mandate", "02-Study", "03-Requirements", "04-Design", "05-Implementation", "06-Verification"} {
			if phase == allowed {
				valid = true
			}
		}
		if !valid || seen[phase] {
			return nil, fmt.Errorf("invalid or duplicate Session phase contribution")
		}
		seen[phase] = true
	}
	if primary != "" && !seen[primary] {
		return nil, fmt.Errorf("primary phase must be a declared contribution")
	}
	return phases, nil
}
