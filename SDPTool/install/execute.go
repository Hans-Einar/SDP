package install

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Result struct {
	Status        string `json:"status"`
	OperationID   string `json:"operationId"`
	MaintenanceID string `json:"maintenanceId"`
	Report        string `json:"report"`
}

// Fault is a test-only callback. Production CLI never accepts executable hooks.
type Executor struct{ Fault func(string, int) error }

func journalHash(j Journal) string { j.Integrity = ""; b, _ := Canonical(j); return Hash(b) }
func (x Executor) Apply(ctx context.Context, root string, p Plan) (Result, error) {
	r := Result{}
	canonical, e := Root(root)
	if e != nil {
		return r, e
	}
	if p.ProjectRoot != canonical || p.SchemaVersion != PlanSchema || p.PlanDigest != planHash(p) {
		return r, fail("plan", 3, "saved plan does not match project/identity")
	}
	unlock, e := lock(root)
	if e != nil {
		return r, e
	}
	defer unlock()
	fresh, e := Build(root, p.Operation, p.Release, p.Previous, p.Adoption)
	if e != nil {
		return r, e
	}
	if fresh.PlanDigest != p.PlanDigest {
		return r, fail("drift", 3, "project or inputs changed since preview")
	}
	if !p.CanApply {
		return r, fail("conflict", 3, "plan contains conflicts")
	}
	if p.NoChange {
		return Result{Status: "no-change"}, nil
	}
	j, e := prepare(p)
	if e != nil {
		return r, e
	}
	r = result(j)
	if _, err := os.Lstat(filepath.Join(root, Operations, j.OperationID)); err == nil {
		return Result{}, fail("operation-exists", 5, "operation identity already exists; inspect or resume %s", j.OperationID)
	} else if !os.IsNotExist(err) {
		return Result{}, err
	}
	// Nothing in the reviewed inventory has changed until this first publication.
	// A caught preparation failure is retryable, not a fictitious resumable journal.
	if e = x.boundary("publication", 0); e == nil {
		e = saveJournal(root, &j)
	}
	if e != nil {
		discardPreparation(root, j)
		if problem, ok := e.(*Error); ok && problem.Exit == 2 {
			return Result{}, e
		}
		return Result{}, fail("preparation", 4, "initial journal not published; no file actions ran: %v", e)
	}
	if e = x.boundary("prepared", 0); e != nil {
		return r, e
	}
	return x.run(ctx, root, &j)
}
func result(j Journal) Result {
	return Result{j.Status, j.OperationID, j.MaintenanceID, "SDP/Maintenance/" + j.MaintenanceID + ".md"}
}
func (x Executor) Resume(ctx context.Context, root, id string) (Result, error) {
	r := Result{}
	if !operationRE.MatchString(id) {
		return r, fail("arguments", 2, "invalid operation id")
	}
	physical, e := Root(root)
	if e != nil {
		return r, e
	}
	root = physical
	unlock, e := lock(root)
	if e != nil {
		return r, e
	}
	defer unlock()
	b, e := Read(filepath.Join(root, Operations, id, "journal.json"), RecordLimit)
	if e != nil {
		return r, e
	}
	var j Journal
	if e = Decode(b, RecordLimit, &j); e != nil {
		return r, fail("journal", 2, "unsupported/invalid journal; legacy recovery needs its original engine: %v", e)
	}
	if j.SchemaVersion != JournalSchema || j.OperationID != id || j.Plan.ProjectRoot != root || j.Integrity != journalHash(j) || j.Plan.PlanDigest != planHash(j.Plan) || id != "install-"+j.Plan.PlanDigest[:24] || j.Next < 0 || j.Next > len(j.Steps) {
		return r, fail("journal", 4, "journal integrity/identity")
	}
	if j.Status != "active" && j.Status != "failed" && j.Status != "completed" {
		return r, fail("journal", 2, "unsupported status")
	}
	if _, e = descriptor(j.Plan.Release); e != nil {
		return r, e
	}
	// Saved exact bytes and derivable finalization protect against an edited action
	// list even when a local journal checksum was recomputed.
	expected, e := prepareAt(j.Plan, j.CreatedAt, j.MaintenanceID)
	if e != nil {
		return r, e
	}
	a, _ := Canonical(expected.Steps)
	c, _ := Canonical(j.Steps)
	if !bytes.Equal(a, c) {
		return r, fail("journal", 4, "journal steps differ from reserved plan")
	}
	if e = progress(root, j); e != nil {
		return result(j), e
	}
	if j.Status == "completed" {
		return result(j), nil
	}
	return x.run(ctx, root, &j)
}
func prepare(p Plan) (Journal, error) {
	// Allocate against the complete observed content, not only the ledger counter.
	var text strings.Builder
	for file, o := range p.Snapshot {
		if o.Type == "file" && (strings.HasPrefix(file, "SDP/Maintenance/") || strings.HasSuffix(file, "Ledger.ndjson")) {
			b, e := Read(filepath.Join(p.ProjectRoot, file), FileLimit)
			if e != nil {
				return Journal{}, e
			}
			text.Write(b)
		}
	}
	for _, a := range p.Actions {
		if strings.HasSuffix(a.Path, "Ledger.ndjson") {
			text.Write(a.Content)
		}
	}
	number := 1
	for {
		id := fmt.Sprintf("MAINT-%s-%04d", p.ProjectID, number)
		if !strings.Contains(text.String(), id) {
			if _, ok := p.Snapshot["SDP/Maintenance/"+id+".md"]; !ok {
				break
			}
		}
		number++
	}
	return prepareAt(p, time.Now().UTC().Format(time.RFC3339), fmt.Sprintf("MAINT-%s-%04d", p.ProjectID, number))
}
func prepareAt(p Plan, when, maint string) (Journal, error) {
	j := Journal{SchemaVersion: JournalSchema, OperationID: "install-" + p.PlanDigest[:24], Plan: p, CreatedAt: when, MaintenanceID: maint, Status: "active", Steps: []Action{}}
	if _, e := time.Parse(time.RFC3339, when); e != nil {
		return j, e
	}
	if !regexp.MustCompile(`^MAINT-` + regexp.QuoteMeta(p.ProjectID) + `-[0-9]{4,}$`).MatchString(maint) {
		return j, fail("journal", 2, "invalid maintenance identity")
	}
	const ledger = "SDP/ProjectManagement/Ledger.ndjson"
	history := []byte{}
	// Use observed baseline from plan: history before bytes need to be embedded in
	// the reserved plan even when unchanged. It is captured as a finalization input.
	if p.Snapshot[ledger].Type == "file" {
		b, e := Read(filepath.Join(p.ProjectRoot, ledger), FileLimit)
		if e != nil {
			return j, e
		}
		history = b
	}
	for _, a := range p.Actions {
		if a.Path == ledger && a.Action == "write" {
			history = a.Content
		}
		j.Steps = append(j.Steps, a)
	}
	// Resume cannot read an already appended ledger as the original baseline.
	// Locate the exact original prefix by its reviewed digest when finalizing.
	if h := observationHash(p.Snapshot, ledger); h != nil && Hash(history) != *h {
		found := false
		for i, c := range history {
			if c == '\n' && Hash(history[:i+1]) == *h {
				history = history[:i+1]
				found = true
				break
			}
		}
		if !found && *h == Hash(nil) {
			history = []byte{}
			found = true
		}
		if !found {
			return j, fail("history", 3, "baseline history changed")
		}
	}
	n := 1
	re := regexp.MustCompile(`"eventId"\s*:\s*"EVT-PM-` + regexp.QuoteMeta(p.ProjectID) + `-([0-9]+)"`)
	for _, m := range re.FindAllSubmatch(history, -1) {
		v, _ := strconv.Atoi(string(m[1]))
		if v >= n {
			n = v + 1
		}
	}
	reportPath := "Maintenance/" + maint + ".md"
	first := fmt.Sprintf("EVT-PM-%s-%06d", p.ProjectID, n)
	ledgerBefore := observationHash(p.Snapshot, ledger)
	for _, a := range p.Actions {
		if a.Path == ledger && a.Action == "write" {
			ledgerBefore = a.After
		}
	}
	for i := 0; i < 2; i++ {
		var prev, from, fromPath any
		to := "active"
		event := "x-management:created"
		if i == 1 {
			prev = first
			from = "active"
			fromPath = reportPath
			to = "completed"
			event = "x-management:completed"
		}
		b, _ := Canonical(map[string]any{"schemaVersion": "1.0", "eventId": fmt.Sprintf("EVT-PM-%s-%06d", p.ProjectID, n+i), "eventType": event, "occurredAt": when, "actor": "sdptool", "commit": nil, "subjectId": maint, "payload": map[string]any{"schemaVersion": "0.1", "kind": "Maintenance", "previousEventId": prev, "from": from, "to": to, "fromPath": fromPath, "toPath": reportPath, "reason": "Execute reviewed SDP installation " + j.OperationID, "links": []string{j.OperationID}}})
		history = append(history, b...)
	}
	var report strings.Builder
	fmt.Fprintf(&report, "# %s — SDP installation\n\n| Field | Value |\n| --- | --- |\n| id | %s |\n| project | %s |\n| state | completed |\n| operation | %s |\n\nBaseline: %s. Target descriptor: %s. Provenance: %s.\n\nCompletion is conditional on the operation journal reporting completed.\nBackups and forward recovery are retained in SDP/.sdp-operations/%s.\n\n## Changes\n\n", maint, maint, p.ProjectID, j.OperationID, p.Baseline, p.Release.SHA256, p.Release.Provenance, j.OperationID)
	for _, a := range p.Actions {
		fmt.Fprintf(&report, "- %s %s\n", a.Action, a.Path)
	}
	report.WriteString("\n## Preserved paths\n\n")
	for _, s := range p.Preserved {
		fmt.Fprintf(&report, "- %s\n", s)
	}
	report.WriteString("\n## Warnings\n\n")
	for _, s := range p.Warnings {
		fmt.Fprintf(&report, "- %s\n", s)
	}
	d, e := descriptor(p.Release)
	if e != nil {
		return j, e
	}
	receipt := Receipt{ReceiptSchema, d.Release, p.Release.SHA256, d.SourceCommit, d.ProcessProfile, d.ManagementProfile, d.Capabilities, p.Release.Provenance, j.OperationID, when}
	b, _ := Canonical(receipt)
	for _, entry := range []struct {
		path    string
		content []byte
	}{{"SDP/" + reportPath, []byte(report.String())}, {ledger, history}, {ReceiptPath, b}} {
		hash := Hash(entry.content)
		before := observationHash(p.Snapshot, entry.path)
		if entry.path == ledger {
			before = ledgerBefore
		}
		j.Steps = append(j.Steps, Action{"write", entry.path, before, &hash, entry.content})
	}
	return j, nil
}
func atomic(root, rel string, b []byte) error {
	if e := Relative(rel); e != nil {
		return e
	}
	if e := SafeAbsolute(filepath.Join(root, rel)); e != nil {
		return e
	}
	r, e := os.OpenRoot(root)
	if e != nil {
		return e
	}
	defer r.Close()
	if e = r.MkdirAll(filepath.FromSlash(path.Dir(rel)), 0700); e != nil {
		return e
	}
	// A stable operation-owned temporary lives under the journal directory. This
	// avoids leaving an untracked temporary inside the observed destination tree.
	temp := Operations + "/temp-" + Hash([]byte(rel))
	if e = r.MkdirAll(Operations, 0700); e != nil {
		return e
	}
	if e = SafeAbsolute(filepath.Join(root, temp)); e != nil {
		return e
	}
	f, e := r.OpenFile(temp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if e != nil {
		return e
	}
	_, e = f.Write(b)
	if e == nil {
		e = f.Sync()
	}
	ce := f.Close()
	if e == nil {
		e = ce
	}
	if e != nil {
		return e
	}
	return r.Rename(temp, rel)
}
func saveJournal(root string, j *Journal) error {
	j.Integrity = journalHash(*j)
	b, e := Canonical(j)
	if e != nil {
		return e
	}
	if len(b) > RecordLimit {
		return fail("limit", 2, "journal exceeds limit")
	}
	return atomic(root, Operations+"/"+j.OperationID+"/journal.json", b)
}
func (x Executor) boundary(name string, index int) error {
	if x.Fault != nil {
		if e := x.Fault(name, index); e != nil {
			return fail("interrupted", 6, "%s/%d: %v", name, index, e)
		}
	}
	return nil
}
func fileHash(root, rel string) (*string, error) {
	if e := Relative(rel); e != nil {
		return nil, e
	}
	p := filepath.Join(root, rel)
	if e := SafeAbsolute(p); e != nil {
		return nil, e
	}
	b, e := Read(p, FileLimit)
	if os.IsNotExist(e) {
		return nil, nil
	}
	if e != nil {
		return nil, e
	}
	h := Hash(b)
	return &h, nil
}
func progress(root string, j Journal) error {
	expected := map[string]Observation{}
	extra := []string{}
	for k, v := range j.Plan.Snapshot {
		expected[k] = v
		extra = append(extra, k)
	}
	apply := func(a Action) {
		if a.Action == "delete" {
			expected[a.Path] = Observation{Type: "absent"}
		} else if a.Action == "mkdir" {
			expected[a.Path] = Observation{Type: "directory"}
		} else {
			expected[a.Path] = Observation{Type: "file", SHA256: a.After}
		}
		if a.Action != "delete" {
			for d := path.Dir(a.Path); d != "."; d = path.Dir(d) {
				expected[d] = Observation{Type: "directory"}
			}
		}
	}
	for i := 0; i < j.Next; i++ {
		a := j.Steps[i]
		if a.Before != nil {
			h, e := fileHash(root, Operations+"/"+j.OperationID+"/backups/"+strconv.Itoa(i))
			if e != nil || !equalHash(h, a.Before) {
				return fail("backup", 4, "missing/corrupt completed backup %d", i)
			}
		}
		apply(a)
	}
	if j.Next < len(j.Steps) {
		a := j.Steps[j.Next]
		if a.Action == "mkdir" {
			s, e := os.Stat(filepath.Join(root, a.Path))
			if e == nil && s.IsDir() {
				apply(a)
			}
		} else {
			actual, e := fileHash(root, a.Path)
			if e != nil {
				return e
			}
			if equalHash(actual, a.After) {
				apply(a)
			} else if !equalHash(actual, a.Before) {
				return fail("drift", 3, "post-failure edit: %s", a.Path)
			}
		}
	}
	// Publishing the journal creates administrative ancestors of an empty project.
	expected["SDP"] = Observation{Type: "directory"}
	for _, a := range j.Steps {
		extra = append(extra, a.Path)
	}
	actual, e := Inspect(root, extra)
	if e != nil {
		return e
	}
	for k, v := range actual.Snapshot {
		want, ok := expected[k]
		if !ok && v.Type == "absent" {
			continue
		}
		if !ok || want.Type != v.Type || !equalHash(want.SHA256, v.SHA256) {
			return fail("drift", 3, "project changed after interruption: %s", k)
		}
	}
	for k, v := range expected {
		if got, ok := actual.Snapshot[k]; !ok || got.Type != v.Type || !equalHash(got.SHA256, v.SHA256) {
			return fail("drift", 3, "missing/changed observed path: %s", k)
		}
	}
	return nil
}
func (x Executor) run(ctx context.Context, root string, j *Journal) (Result, error) {
	abort := func(e error) (Result, error) {
		j.Status = "failed"
		saveErr := saveJournal(root, j)
		if saveErr != nil {
			e = fmt.Errorf("%v; journal checkpoint: %v", e, saveErr)
		}
		return result(*j), fail("mutation", 6, "operation %s: %v", j.OperationID, e)
	}
	if e := progress(root, *j); e != nil {
		return result(*j), e
	}
	for i := j.Next; i < len(j.Steps); i++ {
		select {
		case <-ctx.Done():
			return abort(ctx.Err())
		default:
		}
		a := j.Steps[i]
		if e := Relative(a.Path); e != nil {
			return abort(e)
		}
		if a.Action == "mkdir" {
			if e := SafeAbsolute(filepath.Join(root, a.Path)); e != nil {
				return abort(e)
			}
			r, e := os.OpenRoot(root)
			if e != nil {
				return abort(e)
			}
			e = r.MkdirAll(a.Path, 0700)
			r.Close()
			if e != nil {
				return abort(e)
			}
		} else {
			actual, e := fileHash(root, a.Path)
			if e != nil {
				return abort(e)
			}
			if !equalHash(actual, a.Before) && !equalHash(actual, a.After) {
				return abort(fail("drift", 3, "destination %s", a.Path))
			}
			backup := Operations + "/" + j.OperationID + "/backups/" + strconv.Itoa(i)
			if a.Before != nil {
				stored, e := fileHash(root, backup)
				if e != nil {
					return abort(e)
				}
				if stored == nil {
					if !equalHash(actual, a.Before) {
						return abort(fail("backup", 4, "missing backup: %s", a.Path))
					}
					b, e := Read(filepath.Join(root, a.Path), FileLimit)
					if e != nil {
						return abort(e)
					}
					if e = atomic(root, backup, b); e != nil {
						return abort(e)
					}
				} else if !equalHash(stored, a.Before) {
					return abort(fail("backup", 4, "corrupt backup: %s", a.Path))
				}
			}
			if e = x.boundary("backup", i); e != nil {
				return abort(e)
			}
			current, e := fileHash(root, a.Path)
			if e != nil || !equalHash(current, actual) {
				return abort(fail("drift", 3, "changed during backup: %s", a.Path))
			}
			if !equalHash(actual, a.After) {
				if a.Action == "write" {
					if Hash(a.Content) != *a.After {
						return abort(fail("digest", 4, "action payload"))
					}
					e = atomic(root, a.Path, a.Content)
				} else if a.Action == "delete" {
					r, openErr := os.OpenRoot(root)
					if openErr != nil {
						return abort(openErr)
					}
					e = r.Remove(a.Path)
					r.Close()
				} else {
					return abort(fail("action", 2, "unsupported action"))
				}
				if e != nil {
					return abort(e)
				}
			}
			if e = x.boundary("write", i); e != nil {
				return abort(e)
			}
			actual, e = fileHash(root, a.Path)
			if e != nil || !equalHash(actual, a.After) {
				return abort(fail("verification", 4, "output mismatch: %s", a.Path))
			}
		}
		j.Next = i + 1
		j.Status = "active"
		if e := saveJournal(root, j); e != nil {
			return abort(e)
		}
		if e := x.boundary("journal", i); e != nil {
			return abort(e)
		}
	}
	if e := progress(root, *j); e != nil {
		return abort(e)
	}
	if e := x.boundary("complete", j.Next); e != nil {
		return abort(e)
	}
	j.Status = "completed"
	if e := saveJournal(root, j); e != nil {
		return abort(e)
	}
	return result(*j), nil
}

func discardPreparation(root string, j Journal) {
	r, e := os.OpenRoot(root)
	if e != nil {
		return
	}
	defer r.Close()
	rel := Operations + "/" + j.OperationID + "/journal.json"
	// Remove only the temporary/empty directories belonging to this unstarted
	// operation. Remove fails safely if another writer populated a directory.
	if _, e := r.Lstat(rel); !os.IsNotExist(e) {
		return
	}
	_ = r.Remove(Operations + "/temp-" + Hash([]byte(rel)))
	_ = r.Remove(Operations + "/" + j.OperationID)
	_ = r.Remove(Operations)
	if j.Plan.Snapshot["SDP"].Type != "directory" {
		_ = r.Remove("SDP")
	}
}
