package governance

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

// Only the installed canonical SDP layout is selected by this bounded writer.
// Other layouts remain readable through the existing discovery/context service.
func (s *Service) publicationFiles() (map[string][]byte, error) {
	if s.area != filepath.Join(s.root, "SDP") {
		return nil, fail("incompatible_profile", "publication requires canonical SDP shared-management layout")
	}
	files := map[string][]byte{}
	load := func(p string) error {
		if !literalCandidatePath(p) {
			return fail("path", "invalid canonical record path")
		}
		b, e := readFile(filepath.Join(s.root, p), maxState)
		if e != nil {
			return e
		}
		files[p] = b
		return nil
	}
	if e := load("SDP/ProjectManagement/Ledger.ndjson"); e != nil {
		return nil, e
	}
	latest := map[string]map[string]any{}
	for _, line := range bytes.Split(files["SDP/ProjectManagement/Ledger.ndjson"], []byte("\n")) {
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		var ev map[string]any
		if e := decode(line, &ev); e != nil {
			return nil, e
		}
		typ, _ := ev["eventType"].(string)
		if strings.HasPrefix(typ, "x-kanban:") || strings.HasPrefix(typ, "x-management:") {
			id, _ := ev["subjectId"].(string)
			latest[id] = ev
		}
	}
	for _, ev := range latest {
		p, ok := ev["payload"].(map[string]any)
		if !ok {
			return nil, fail("conflict", "invalid management payload")
		}
		dest, ok := p["toPath"].(string)
		if !ok || !relative(dest) {
			return nil, fail("path", "invalid management destination")
		}
		prefix := "SDP/"
		if strings.HasPrefix(fmt.Sprint(ev["eventType"]), "x-kanban:") {
			prefix += "KanBan/"
		}
		if e := load(prefix + dest); e != nil {
			return nil, e
		}
	}
	e := filepath.WalkDir(filepath.Join(s.area, "KanBan"), func(path string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if d.Type()&os.ModeSymlink != 0 {
			return fail("path", "symlink in publication board")
		}
		if !d.IsDir() && strings.HasSuffix(path, ".md") {
			rel, _ := filepath.Rel(s.root, path)
			if e = load(filepath.ToSlash(rel)); e != nil {
				return e
			}
		}
		if len(files) > 1024 {
			return fail("limit", "publication record count exceeds 1024")
		}
		return nil
	})
	if e != nil {
		return nil, e
	}
	if e = validatePublicationHistory(files["SDP/ProjectManagement/Ledger.ndjson"], files); e != nil {
		return nil, fail("incompatible_profile", "canonical record validation: "+e.Error())
	}
	return files, nil
}
func publicationEvents(b []byte) ([]map[string]any, error) {
	out := []map[string]any{}
	seen := map[string]bool{}
	for _, line := range bytes.Split(b, []byte("\n")) {
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		var ev map[string]any
		if e := decode(line, &ev); e != nil {
			return nil, e
		}
		id, _ := ev["eventId"].(string)
		if id == "" || seen[id] {
			return nil, fail("conflict", "duplicate or absent canonical event identity")
		}
		seen[id] = true
		out = append(out, ev)
	}
	return out, nil
}
func nextPublicationEvent(events []map[string]any, prefix string) (string, error) {
	n := uint64(0)
	for _, ev := range events {
		id, _ := ev["eventId"].(string)
		if strings.HasPrefix(id, prefix) {
			v, e := strconv.ParseUint(strings.TrimPrefix(id, prefix), 10, 64)
			if e != nil {
				return "", fail("conflict", "invalid event sequence")
			}
			if v > n {
				n = v
			}
		}
	}
	if n == ^uint64(0) {
		return "", fail("limit", "event sequence exhausted")
	}
	return fmt.Sprintf("%s%06d", prefix, n+1), nil
}
func appendRecord(b []byte, tail string) string {
	if len(b) > 0 && b[len(b)-1] != '\n' {
		return string(b) + "\n" + tail
	}
	return string(b) + tail
}
func publicationProse(text string) string {
	return "    " + strings.ReplaceAll(text, "\n", "\n    ") + "\n"
}
func (s *Service) buildPublication(st snapshot, p Publication) (Publication, error) {
	files, e := s.publicationFiles()
	if e != nil {
		return p, e
	}
	q := p.Request
	r := st.Runs[q.RunID]
	d := st.Reviews[q.ReviewID]
	ctx, e := s.runContext(r)
	if e != nil {
		return p, e
	}
	if ctx.Revision != r.InputRevision || len(ctx.Conditions) > 0 {
		return p, fail("stale_revision", "publication work authority changed")
	}
	cm, e := metadata(files[r.Selection.Card])
	if e != nil {
		return p, e
	}
	pm, e := metadata(files[r.Selection.Plan])
	if e != nil {
		return p, e
	}
	match := regexp.MustCompile(`^KB-([A-Z][A-Z0-9]*)-[0-9]{3,}$`).FindStringSubmatch(ctx.CardID)
	if match == nil {
		return p, fail("incompatible_profile", "canonical card namespace required")
	}
	namespace := match[1]
	if cm["project"] != namespace || pm["project"] != namespace {
		return p, fail("incompatible_profile", "selected project namespaces must agree")
	}
	events, e := publicationEvents(files["SDP/ProjectManagement/Ledger.ndjson"])
	if e != nil {
		return p, e
	}
	latest := map[string]map[string]any{}
	for _, ev := range events {
		subject, _ := ev["subjectId"].(string)
		latest[subject] = ev
	}
	p.EvidencePath = "SDP/Verification/PG/" + q.OperationID + ".md"
	p.Inputs = []Source{}
	p.Changes = []PublicationChange{}
	p.EventIDs = []string{}
	for _, src := range ctx.Sources {
		if _, ok := files[src.Path]; !ok {
			if src.SHA256 == "absent" {
				p.Inputs = append(p.Inputs, src)
			} else {
				b, e := readFile(filepath.Join(s.root, src.Path), maxState)
				if e != nil {
					return p, e
				}
				files[src.Path] = b
			}
		}
	}
	tracePath := "SDP/Traceability/Ledger.ndjson"
	trace, e := readFile(filepath.Join(s.root, tracePath), maxState)
	if e != nil && !os.IsNotExist(e) {
		return p, e
	}
	traceAbsent := os.IsNotExist(e)
	files[tracePath] = trace
	if _, e = publicationEvents(trace); e != nil {
		return p, e
	}
	if _, e = readFile(filepath.Join(s.root, p.EvidencePath), maxState); !os.IsNotExist(e) {
		return p, fail("conflict", "publication evidence destination must be absent")
	}
	// No generated record may mutate candidate inputs or the selected authority.
	targets := []string{r.Selection.Session, r.Selection.Card, r.Selection.Plan, p.EvidencePath, "SDP/ProjectManagement/Ledger.ndjson", tracePath}
	if len(slices.Compact(slices.Sorted(slices.Values(targets)))) != 6 {
		return p, fail("conflict", "publication destinations overlap")
	}
	for _, path := range targets {
		if path == r.Selection.Authority {
			return p, fail("conflict", "publication cannot rewrite authority")
		}
		for _, scope := range append(append([]string{}, d.Spec.Candidate.WritablePaths...), d.Spec.Candidate.Dependencies...) {
			if within(path, scope) || within(scope, path) {
				return p, fail("conflict", "publication overlaps candidate scope")
			}
		}
		for _, sid := range d.Spec.EvidenceIDs {
			if path == st.Receipts[sid].Result.Evidence.Request.Artifact.Path {
				return p, fail("conflict", "publication overlaps test evidence")
			}
		}
	}
	for name, b := range files {
		before := hash(b)
		if name == tracePath && traceAbsent {
			before = "absent"
		}
		p.Inputs = append(p.Inputs, Source{Path: name, SHA256: before})
	}
	p.Inputs = append(p.Inputs, Source{Path: p.EvidencePath, SHA256: "absent"})
	slices.SortFunc(p.Inputs, func(a, b Source) int { return strings.Compare(a.Path, b.Path) })
	add := func(path, content string) {
		before := "absent"
		if b, ok := files[path]; ok {
			before = hash(b)
		}
		if path == tracePath && traceAbsent {
			before = "absent"
		}
		p.Changes = append(p.Changes, PublicationChange{Path: path, Before: before, After: hash([]byte(content)), Content: content})
	}
	tests := []EvidenceRecord{}
	for _, id := range d.Spec.EvidenceIDs {
		tests = append(tests, *st.Receipts[id].Result.Evidence)
	}
	artifact := struct {
		Schema    string             `json:"schema"`
		Request   PublicationRequest `json:"request"`
		Selection Selection          `json:"selection"`
		Review    Review             `json:"review"`
		Tests     []EvidenceRecord   `json:"tests"`
		Limits    []string           `json:"limitations"`
	}{PublicationSchema, q, r.Selection, d, tests, []string{"tests are controller-attributed external evidence", "fresh review context observed by controller; not owner acceptance", "records remain active; routine completion is separate"}}
	data, e := json.MarshalIndent(artifact, "", "  ")
	if e != nil {
		return p, e
	}
	add(p.EvidencePath, "# Published PG review result\n\nOperation: "+q.OperationID+"\n\n```json\n"+string(data)+"\n```\n")
	for _, path := range targets[:3] {
		rel, _ := filepath.Rel(filepath.Dir(path), p.EvidencePath)
		journal := "\n## PG result " + q.OperationID + "\n\nRecorded review of candidate `" + d.Spec.Candidate.Digest + "` for run `" + r.ID + "`.\nTests: attributed passed evidence; review: observed passed in a separate context.\nThis entry records a bounded result; lifecycle and owner acceptance remain unchanged.\n\nWork summary:\n\n" + publicationProse(q.Summary) + "\nRemaining work:\n\n" + publicationProse(q.Remaining) + "\n[Exact result and provenance](" + filepath.ToSlash(rel) + ").\n"
		add(path, appendRecord(files[path], journal))
	}
	newHistory := string(files["SDP/ProjectManagement/Ledger.ndjson"])
	for _, subject := range []string{ctx.PlanID, ctx.CardID} {
		old := latest[subject]
		if old == nil {
			return p, fail("missing_coverage", "selected history absent")
		}
		payload, ok := old["payload"].(map[string]any)
		if !ok {
			return p, fail("conflict", "invalid management payload")
		}
		np := map[string]any{}
		for k, v := range payload {
			np[k] = v
		}
		np["previousEventId"] = old["eventId"]
		np["from"] = payload["to"]
		np["fromPath"] = payload["toPath"]
		np["reason"] = "Recorded candidate-bound PG review result; lifecycle remains active. Publication " + q.OperationID
		np["links"] = []string{ctx.CardID, ctx.PlanID}
		typ, prefix := "x-management:updated", "EVT-PM-"+namespace+"-"
		if subject == ctx.CardID {
			typ = "x-kanban:reviewed"
			prefix = "EVT-KB-" + namespace + "-"
			np["links"] = []string{ctx.PlanID}
			delete(np, "lineage")
		}
		id, e := nextPublicationEvent(events, prefix)
		if e != nil {
			return p, e
		}
		ev := map[string]any{"schemaVersion": "1.0", "eventId": id, "eventType": typ, "occurredAt": p.CreatedAt, "actor": "sdp-controller", "commit": nil, "subjectId": subject, "payload": np}
		b, _ := json.Marshal(ev)
		newHistory = appendRecord([]byte(newHistory), string(b)+"\n")
		events = append(events, ev)
		p.EventIDs = append(p.EventIDs, id)
	}
	add("SDP/ProjectManagement/Ledger.ndjson", newHistory)
	traceID := "EVT-" + q.System + "-PG-" + strings.ToUpper(q.OperationID)
	subject := q.System + "-VER-PG-" + strings.ToUpper(q.OperationID)
	tev := map[string]any{"schemaVersion": "1.0", "eventId": traceID, "eventType": "x-verification:recorded", "occurredAt": p.CreatedAt, "actor": "sdp-controller", "commit": nil, "subjectId": subject, "payload": map[string]any{"plan": ctx.PlanID, "card": ctx.CardID, "session": ctx.SessionID, "managementEvent": p.EventIDs[0], "evidence": strings.TrimPrefix(p.EvidencePath, "SDP/"), "candidate": d.Spec.Candidate.Digest, "review": d.ID, "note": "Observed passed review and attributed passed tests; no owner acceptance or lifecycle completion."}}
	b, _ := json.Marshal(tev)
	add(tracePath, appendRecord(trace, string(b)+"\n"))
	p.EventIDs = append(p.EventIDs, traceID)
	size := 0
	for _, ch := range p.Changes {
		size += len(ch.Content)
	}
	if size > 4<<20 || len(p.Inputs) > 2048 {
		return p, fail("limit", "publication record set exceeds bound")
	}
	if e = s.validatePublicationRecords(p); e != nil {
		return p, e
	}
	return p, nil
}
func (s *Service) validatePublicationRecords(p Publication) error {
	files := map[string][]byte{}
	changes := map[string]PublicationChange{}
	for _, ch := range p.Changes {
		changes[ch.Path] = ch
	}
	for _, src := range p.Inputs {
		if ch, ok := changes[src.Path]; ok {
			files[src.Path] = []byte(ch.Content)
			continue
		}
		if src.SHA256 == "absent" {
			continue
		}
		b, e := readFile(filepath.Join(s.root, src.Path), maxState)
		if e != nil || hash(b) != src.SHA256 {
			return fail("reconciliation_required", "validation input changed: "+src.Path)
		}
		files[src.Path] = b
	}
	if e := validatePublicationHistory(files["SDP/ProjectManagement/Ledger.ndjson"], files); e != nil {
		return fail("conflict", "proposed canonical records invalid: "+e.Error())
	}
	for _, path := range []string{"SDP/ProjectManagement/Ledger.ndjson", "SDP/Traceability/Ledger.ndjson"} {
		if _, e := publicationEvents(files[path]); e != nil {
			return e
		}
	}
	return nil
}
