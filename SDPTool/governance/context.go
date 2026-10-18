package governance

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Service is bound by the adapter to an explicitly discovered project. New does
// not create state, enabling status/context/route to remain genuinely read-only.
type Service struct {
	root, area, project, state string
	fault                      func(string) error
}

func New(root, area, project, state string) (*Service, error) {
	root, e := filepath.Abs(root)
	if e != nil {
		return nil, e
	}
	root, e = filepath.EvalSymlinks(root)
	if e != nil {
		return nil, e
	}
	area, e = filepath.Abs(area)
	if e != nil {
		return nil, e
	}
	if e = safeAbsolute(area); e != nil {
		return nil, e
	}
	rel, e := filepath.Rel(root, area)
	if e != nil || !relative(filepath.ToSlash(rel)) {
		return nil, fail("path", "SDP area must be inside selected project")
	}
	if project == "" {
		return nil, fail("missing_coverage", "project identity missing")
	}
	state, e = filepath.Abs(state)
	if e != nil {
		return nil, e
	}
	if e = safeAbsolute(state); e != nil {
		return nil, e
	}
	// State must not shadow project sources or make candidate hashes depend on runs.
	inside, e := filepath.Rel(root, state)
	if e != nil {
		return nil, e
	}
	if inside == "." || inside != ".." && !strings.HasPrefix(inside, ".."+string(filepath.Separator)) {
		return nil, fail("path", "operational state must be outside the project")
	}
	return &Service{root: root, area: area, project: project, state: state}, nil
}
func metadata(b []byte) (map[string]string, error) {
	m := map[string]string{}
	started := false
	for _, line := range strings.Split(string(b), "\n") {
		cells := strings.Split(line, "|")
		if !started {
			if len(cells) == 4 && strings.TrimSpace(cells[1]) == "Field" && strings.TrimSpace(cells[2]) == "Value" {
				started = true
			}
			continue
		}
		if !strings.HasPrefix(line, "|") {
			break
		}
		if len(cells) != 4 {
			return nil, fail("conflict", "ambiguous metadata table")
		}
		k, v := strings.TrimSpace(cells[1]), strings.TrimSpace(cells[2])
		if strings.HasPrefix(k, "---") {
			continue
		}
		if _, ok := m[k]; ok {
			return nil, fail("conflict", "duplicate metadata row: "+k)
		}
		m[k] = v
	}
	if !started {
		return nil, fail("missing_coverage", "Field/Value metadata missing")
	}
	return m, nil
}
func (s *Service) Context(sel Selection) (Context, error) {
	return s.contextFor(sel, []string{"MaintenancePlan"}, "MAINT-")
}

func (s *Service) contextFor(sel Selection, planTypes []string, idPrefix string) (Context, error) {
	c := Context{Schema: Schema, Operation: "context", Root: s.root, ProjectRef: s.project, Selection: sel, Sources: []Source{}, Conditions: []string{}, Coverage: "known"}
	state, _, e := s.load()
	if e != nil {
		return c, e
	}
	c.CheckoutID = state.CheckoutID
	// Store binding is not created by context or route reads.

	bodies := map[string][]byte{}
	paths := map[string]bool{}
	relArea, _ := filepath.Rel(s.root, s.area)
	areaPrefix := filepath.ToSlash(relArea) + "/"
	for _, p := range []string{sel.Card, sel.Plan, sel.Session, sel.Authority} {
		if p == "" {
			continue
		}
		if !relative(p) || !strings.HasPrefix(p, areaPrefix) || !strings.HasSuffix(p, ".md") {
			return c, fail("path", "selection must reference Markdown inside the selected SDP area")
		}
		paths[p] = true
	}
	ledgerPath := areaPrefix + "ProjectManagement/Ledger.ndjson"
	paths[ledgerPath] = true
	paths["AGENTS.md"] = true
	paths[areaPrefix+"SDP-project.manifest.yaml"] = true
	names := []string{}
	for p := range paths {
		names = append(names, p)
	}
	sort.Strings(names)
	for _, p := range names {
		limit := MaxRequest
		if p == ledgerPath {
			limit = maxState
		}
		b, e := readFile(filepath.Join(s.root, filepath.FromSlash(p)), int64(limit))
		if os.IsNotExist(e) {
			c.Sources = append(c.Sources, Source{p, "absent"})
			if p != "AGENTS.md" && p != areaPrefix+"SDP-project.manifest.yaml" {
				c.Coverage = "unavailable"
				c.Conditions = append(c.Conditions, "missing source: "+p)
			}
			continue
		}
		if e != nil {
			return c, e
		}
		bodies[p] = b
		c.Sources = append(c.Sources, Source{p, hash(b)})
	}
	add := func(message string) { c.Conditions = append(c.Conditions, message) }
	meta := func(p string) map[string]string {
		if p == "" {
			add("work reference missing")
			return nil
		}
		if bodies[p] == nil {
			return nil
		}
		m, e := metadata(bodies[p])
		if e != nil {
			add(e.Error())
			c.Coverage = "conflict"
			return nil
		}
		return m
	}
	card, plan, session := meta(sel.Card), meta(sel.Plan), meta(sel.Session)
	c.CardID = card["id"]
	c.PlanID = plan["id"]
	c.SessionID = session["Session reference"]
	if sel.Authority == "" {
		add("authority reference missing")
	}
	if !strings.HasPrefix(c.CardID, "KB-") {
		add("card identity missing")
	}
	if !strings.HasPrefix(c.PlanID, idPrefix) {
		add("required plan identity prefix: " + idPrefix)
	}
	if !strings.HasPrefix(c.SessionID, "SESSION-") {
		add("Session identity missing")
	}
	if card["CardState"] != "ready" && card["CardState"] != "in-progress" {
		add("card must be ready or in-progress")
	}
	if !strings.HasPrefix(sel.Card, areaPrefix+"KanBan/active/") {
		add("card must be in active board directory")
	}
	if card["PlanId"] != c.PlanID {
		add("card and plan identities disagree")
	}
	matchedPlan := false
	for _, typ := range planTypes {
		if plan["PlanType"] == typ {
			matchedPlan = true
		}
	}
	if !matchedPlan || plan["state"] != "active" {
		add("active " + strings.Join(planTypes, " or ") + " required")
	}
	if session["Status"] != "active" {
		add("active Session required")
	}
	if sel.Step == "" || !strings.Contains(string(bodies[sel.Session]), "| "+sel.Step+" |") {
		add("Session step is missing")
	}
	if bodies[ledgerPath] != nil && c.CardID != "" && c.PlanID != "" {
		e := ledgerCoverage(bodies[ledgerPath], c.CardID, c.PlanID, strings.TrimPrefix(sel.Card, areaPrefix+"KanBan/"), strings.TrimPrefix(sel.Plan, areaPrefix))
		if e != nil {
			add(e.Error())
			c.Coverage = "conflict"
		}
	}
	revision, _ := json.Marshal(struct {
		Project, Root string
		Selection     Selection
		Sources       []Source
	}{s.project, s.root, sel, c.Sources})
	c.Revision = hash(revision)
	return c, nil
}
func (s *Service) Route(q RouteRequest) (Route, error) {
	r := Route{Schema: Schema, Operation: "route", Disposition: "missing", Routine: Maintenance(), Basis: "caller-proposed category; no automatic language classification"}
	r.RoutineDigest = r.Routine.Digest()
	if q.Schema != Schema {
		return r, fail("incompatible_profile", "unsupported application schema")
	}
	if len(q.Intent) > 64<<10 || len(q.Categories) == 0 || len(q.Categories) > 16 {
		return r, fail("invalid_request", "one to sixteen categories and at most 64 KiB intent required")
	}
	var e error
	r.Context, e = s.Context(q.Selection)
	if e != nil {
		return r, e
	}
	if q.InputRevision != "" && q.InputRevision != r.Context.Revision {
		r.Disposition = "conflict"
		return r, nil
	}
	if len(q.Categories) != 1 {
		r.Disposition = "ambiguous"
		return r, nil
	}
	if q.Categories[0] == "status" {
		r.Disposition = "no-work"
		return r, nil
	}
	if q.Categories[0] != "maintenance" || q.RoutineID != "maintenance" {
		return r, nil
	}
	if q.RoutineVersion != "0.1" {
		r.Disposition = "incompatible"
		return r, nil
	}
	switch r.Context.Coverage {
	case "unavailable", "conflict":
		r.Disposition = r.Context.Coverage
		return r, nil
	}
	if len(r.Context.Conditions) > 0 {
		r.Disposition = "missing"
		return r, nil
	}
	r.Disposition = "matched"
	return r, nil
}
