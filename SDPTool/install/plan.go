package install

import (
	"bytes"
	"encoding/json"
	"github.com/Hans-Einar/SDP/SDPTool/bootstrap"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type Options struct {
	Root, Operation, Artifact, PreviousArtifact, Manifest string
	AllowUnreleased                                       bool
}

func LocalInput(p string, allow bool) (Input, error) {
	if !allow {
		return Input{}, fail("unreleased", 2, "local artifacts require --allow-unreleased")
	}
	abs, e := filepath.Abs(p)
	if e != nil {
		return Input{}, e
	}
	b, e := Read(abs, MetadataLimit)
	if e != nil {
		return Input{}, describeError(e)
	}
	return Input{Path: abs, SHA256: Hash(b), Bytes: b, Provenance: "local-development"}, nil
}
func descriptor(in Input) (Descriptor, error) {
	var d Descriptor
	if Hash(in.Bytes) != in.SHA256 {
		return d, fail("digest", 4, "input digest mismatch")
	}
	if in.Provenance != "local-development" {
		if e := bootstrap.Verify(in.Bytes, in.Signature, in.KeyID, in.Provenance); e != nil {
			return d, fail("trust", 4, "%v", e)
		}
	}
	if e := Decode(in.Bytes, MetadataLimit, &d); e != nil {
		return d, e
	}
	return d, ValidateDescriptor(d)
}
func Preview(o Options) (Plan, error) {
	var p Plan
	root, e := Root(o.Root)
	if e != nil {
		return p, e
	}
	r, e := LocalInput(o.Artifact, o.AllowUnreleased)
	if e != nil {
		return p, e
	}
	return previewInput(o, root, r)
}
func previewInput(o Options, root string, r Input) (Plan, error) {
	var p Plan
	var old, adopt *Input
	if o.PreviousArtifact != "" {
		v, e := LocalInput(o.PreviousArtifact, o.AllowUnreleased)
		if e != nil {
			return p, e
		}
		old = &v
	}
	if o.Manifest != "" {
		v, e := LocalInput(o.Manifest, true)
		if e != nil {
			return p, e
		}
		v.Provenance = "adoption"
		adopt = &v
	}
	if e := cacheInput(r); e != nil {
		return p, e
	}
	if old == nil {
		if b, e := Read(filepath.Join(root, ReceiptPath), MetadataLimit); e == nil {
			var receipt Receipt
			if Decode(b, MetadataLimit, &receipt) == nil && receipt.DescriptorDigest != r.SHA256 {
				v, e := cachedInput(receipt.DescriptorDigest)
				if e == nil {
					old = &v
				}
			}
		}
	}
	return Build(root, o.Operation, r, old, adopt)
}
func targetPaths(d Descriptor) []string {
	p := []string{"AGENTS-project.md", ReceiptPath, "SDP/SDP-project.manifest.yaml", "SDP/ProjectManagement/Ledger.ndjson", "SDP/Traceability/Ledger.ndjson", "SDP/KanBan/board.json"}
	for _, f := range d.Files {
		p = append(p, f.Path)
	}
	return p
}
func Build(root, operation string, release Input, previous, adoption *Input) (Plan, error) {
	p := Plan{SchemaVersion: PlanSchema, Operation: operation, ProjectRoot: root, Release: release, Previous: previous, Adoption: adoption, Policy: "ordinal-writes-then-removals/1", Actions: []Action{}, Preserved: []string{}, Conflicts: []string{}, Warnings: []string{"Scan excludes .git, node_modules, .venv, vendor, build and .cache; non-Markdown references are not rewritten."}}
	if operation != "install" && operation != "upgrade" {
		return p, fail("arguments", 2, "install or upgrade required")
	}
	physical, e := Root(root)
	if e != nil || physical != root {
		return p, fail("root", 2, "noncanonical project root")
	}
	pending, e := Pending(root)
	if e != nil {
		return p, e
	}
	if len(pending) > 0 {
		return p, fail("pending", 5, "resume pending operation %s with its matching engine", strings.Join(pending, ","))
	}
	d, e := descriptor(release)
	if e != nil {
		return p, e
	}
	extra := targetPaths(d)
	var ad Adoption
	if adoption != nil {
		if Hash(adoption.Bytes) != adoption.SHA256 {
			return p, fail("digest", 4, "adoption digest")
		}
		if e = Decode(adoption.Bytes, MetadataLimit, &ad); e != nil {
			return p, e
		}
		if ad.SchemaVersion != AdoptionSchema || ad.ProjectRoot != root || ad.TargetDigest != release.SHA256 || (ad.Baseline != "manual" && ad.Baseline != "known") {
			return p, fail("adoption", 2, "adoption identity/root/target mismatch")
		}
		for k, v := range ad.Inventory {
			if e = badObservation(k, v); e != nil {
				return p, fail("adoption", 2, "%v", e)
			}
			extra = append(extra, k)
		}
	}
	before, e := Inspect(root, extra)
	if e != nil {
		return p, e
	}
	p.Snapshot = before.Snapshot
	after := map[string][]byte{}
	for k, v := range before.Files {
		after[k] = v
	}
	p.Baseline = "clean"
	if s := p.Snapshot["SDP"]; s.Type == "directory" {
		p.Baseline = "manual"
	}
	var old Descriptor
	if b, ok := before.Files[ReceiptPath]; ok {
		var r Receipt
		e = Decode(b, MetadataLimit, &r)
		if e == nil {
			e = ValidateReceipt(r)
		}
		if e != nil {
			if adoption == nil {
				return p, fail("receipt", 2, "legacy/invalid receipt needs explicit adoption and original inventory: %v", e)
			}
			p.Warnings = append(p.Warnings, "Historical receipt is not authoritative inventory; prior release remains unknown.")
		} else {
			p.OldReceipt = &r
			p.Baseline = "known"
			if previous == nil && r.DescriptorDigest == release.SHA256 {
				previous = &release
				p.Previous = previous
			}
			if previous == nil {
				return p, fail("previous-release", 4, "original descriptor required: %s (use --previous-artifact for development)", r.DescriptorDigest)
			}
			old, e = descriptor(*previous)
			if e != nil {
				return p, e
			}
			if previous.SHA256 != r.DescriptorDigest || old.Release != r.Release || old.ProcessProfile != r.ProcessProfile || old.SourceCommit != r.SourceCommit || previous.Provenance != r.Provenance {
				return p, fail("receipt", 4, "old receipt does not match verified descriptor")
			}
			if release.SHA256 != previous.SHA256 && !contains(d.UpgradesFrom, previous.SHA256) {
				p.Conflicts = append(p.Conflicts, "unsupported-release-transition")
			}
		}
	}
	if operation == "install" && p.Baseline != "clean" {
		p.Conflicts = append(p.Conflicts, "existing-SDP-requires-upgrade-and-adoption")
	}
	if operation == "upgrade" && p.Baseline != "known" && adoption == nil {
		p.Conflicts = append(p.Conflicts, "unknown-baseline-requires-adoption-manifest")
	}
	if adoption != nil {
		actual := map[string]Observation{}
		for k, v := range before.Snapshot {
			if v.Type != "absent" {
				actual[k] = v
			}
		}
		expected := map[string]Observation{}
		for k, v := range ad.Inventory {
			if v.Type != "absent" {
				expected[k] = v
			} else if before.Snapshot[k].Type != "absent" {
				p.Conflicts = append(p.Conflicts, "adoption-absence-drift: "+k)
			}
		}
		if !SameSnapshot(actual, expected) {
			p.Conflicts = append(p.Conflicts, "adoption-inventory-drift")
		}
		if ad.Baseline != p.Baseline {
			p.Conflicts = append(p.Conflicts, "adoption-baseline-kind-mismatch")
		}
	}
	moves := map[string]string{}
	for i, m := range ad.Moves {
		if e = Relative(m.From); e != nil {
			return p, e
		}
		if e = Relative(m.To); e != nil {
			return p, e
		}
		boardMove := m.From == "SDP/Agents/KanBan" && m.To == "SDP/KanBan"
		if !strings.HasPrefix(m.From, "SDP/") || !strings.HasPrefix(m.To, "SDP/") || Reserved(m.From) || Reserved(m.To) && !boardMove || overlap(m.From, m.To) {
			return p, fail("move", 2, "unsupported preserving move %s to %s", m.From, m.To)
		}
		for j, n := range ad.Moves {
			if j != i && (overlap(m.From, n.From) || overlap(m.To, n.To) || overlap(m.From, n.To)) {
				return p, fail("move", 2, "overlapping/cyclic mappings")
			}
		}
		found := false
		for k, v := range before.Files {
			if k == m.From || strings.HasPrefix(k, m.From+"/") {
				found = true
				dest := m.To + strings.TrimPrefix(k, m.From)
				if _, ok := before.Snapshot[dest]; ok && before.Snapshot[dest].Type != "absent" {
					p.Conflicts = append(p.Conflicts, "move-destination-exists: "+dest)
					continue
				}
				after[dest] = v
				delete(after, k)
				moves[k] = dest
			}
		}
		if !found {
			p.Conflicts = append(p.Conflicts, "missing-move-source: "+m.From)
		}
	}
	if len(moves) > 0 && !ad.AllowReferenceWarnings {
		p.Conflicts = append(p.Conflicts, "relocation-requires-explicit-reference-warning-disposition")
	}
	projectID, e := board(after, moves)
	if e != nil {
		return p, e
	}
	p.ProjectID = projectID
	for source := range before.Files {
		dest := source
		if moved, ok := moves[source]; ok {
			dest = moved
		}
		b, ok := after[dest]
		if !ok || !strings.HasSuffix(strings.ToLower(source), ".md") {
			continue
		}
		after[dest] = rebaseMarkdown(b, source, dest, moves, ad.Moves)
	}
	oldManaged := map[string]File{}
	for _, f := range old.Files {
		if f.Ownership == "managed" {
			oldManaged[f.Path] = f
		}
	}
	for _, s := range ad.RefreshManaged {
		if !managed(s) || Reserved(s) {
			return p, fail("ownership", 2, "invalid refresh %s", s)
		}
	}
	for _, f := range d.Files {
		current, exists := after[f.Path]
		if f.Type == "directory" {
			if before.Snapshot[f.Path].Type == "file" {
				p.Conflicts = append(p.Conflicts, "directory-is-file: "+f.Path)
			} else if before.Snapshot[f.Path].Type == "absent" {
				p.Actions = append(p.Actions, Action{Action: "mkdir", Path: f.Path})
			}
			continue
		}
		if before.Snapshot[f.Path].Type == "directory" {
			p.Conflicts = append(p.Conflicts, "file-is-directory: "+f.Path)
			continue
		}
		if !exists {
			after[f.Path] = f.Content
			continue
		}
		if f.Ownership == "initialize-if-missing" || bytes.Equal(current, f.Content) {
			continue
		}
		prev, known := oldManaged[f.Path]
		unchanged := known && prev.SHA256 != nil && Hash(current) == *prev.SHA256
		if !unchanged && !contains(ad.RefreshManaged, f.Path) && !(f.Path == "AGENTS.md" && !known && p.Baseline == "clean") {
			p.Conflicts = append(p.Conflicts, "managed-edit-requires-named-refresh: "+f.Path)
			continue
		}
		if f.Path == "AGENTS.md" && !known {
			if saved, ok := after["AGENTS-project.md"]; ok && !bytes.Equal(saved, current) {
				p.Conflicts = append(p.Conflicts, "agents-preservation-collision")
				continue
			}
			after["AGENTS-project.md"] = current
		}
		after[f.Path] = f.Content
	}
	for _, s := range d.Retired {
		f, ok := oldManaged[s]
		if !ok {
			return p, fail("retirement", 2, "retired file not in old managed inventory: %s", s)
		}
		if current, exists := after[s]; exists {
			if f.SHA256 == nil || Hash(current) != *f.SHA256 {
				p.Conflicts = append(p.Conflicts, "modified-retired-file: "+s)
			} else {
				delete(after, s)
			}
		}
	}
	if _, ok := after["SDP/Traceability/Ledger.ndjson"]; !ok {
		after["SDP/Traceability/Ledger.ndjson"] = []byte{}
	}
	// Legacy navigation.json is project-owned historical data, never a discovery registry.
	manifest := "SDP/SDP-project.manifest.yaml"
	if _, ok := after[manifest]; !ok {
		after[manifest] = []byte("schemaVersion: \"1.0\"\ninstalled:\n  manifestPath: Framework/installed-toolkit.manifest.yaml\n")
	}
	for dest, content := range after {
		if dest == ReceiptPath {
			continue
		}
		hash := Hash(content)
		prev := observationHash(before.Snapshot, dest)
		if equalHash(prev, &hash) {
			p.Preserved = append(p.Preserved, dest)
			continue
		}
		p.Actions = append(p.Actions, Action{"write", dest, prev, &hash, content})
	}
	sort.Slice(p.Actions, func(i, j int) bool { return p.Actions[i].Path < p.Actions[j].Path })
	deleted := []string{}
	for s := range before.Files {
		if _, ok := after[s]; !ok {
			deleted = append(deleted, s)
		}
	}
	sort.Strings(deleted)
	for _, s := range deleted {
		p.Actions = append(p.Actions, Action{"delete", s, observationHash(before.Snapshot, s), nil, nil})
	}
	sort.Strings(p.Preserved)
	sort.Strings(p.Conflicts)
	// Bind destinations and their ancestor types, including paths absent initially.
	for _, a := range p.Actions {
		extra = append(extra, a.Path)
	}
	final, e := Inspect(root, extra)
	if e != nil {
		return p, e
	}
	for k, v := range p.Snapshot {
		if b, ok := final.Snapshot[k]; !ok || b.Type != v.Type || !equalHash(b.SHA256, v.SHA256) {
			return p, fail("drift", 3, "changed during planning: %s", k)
		}
	}
	// Destination expansion may add absent paths, never new observed content.
	for k, v := range final.Snapshot {
		if v.Type != "absent" {
			if old, ok := p.Snapshot[k]; !ok || old.Type != v.Type || !equalHash(old.SHA256, v.SHA256) {
				return p, fail("drift", 3, "new observation during planning: %s", k)
			}
		}
	}
	// Repeat the full observation to detect additions as well as changed existing paths.
	check, e := Inspect(root, extra)
	if e != nil {
		return p, e
	}
	if !SameSnapshot(final.Snapshot, check.Snapshot) {
		return p, fail("drift", 3, "changed during planning")
	}
	p.Snapshot = check.Snapshot
	for _, in := range []*Input{&release, previous, adoption} {
		if in == nil {
			continue
		}
		b, e := Read(in.Path, MetadataLimit)
		if e != nil || Hash(b) != in.SHA256 {
			return p, fail("drift", 3, "input changed during planning: %s", in.Path)
		}
	}
	p.CanApply = len(p.Conflicts) == 0
	p.NoChange = len(p.Actions) == 0 && p.OldReceipt != nil && p.OldReceipt.DescriptorDigest == release.SHA256
	p.PlanDigest = planHash(p)
	return p, nil
}
func contains(a []string, s string) bool {
	for _, x := range a {
		if x == s {
			return true
		}
	}
	return false
}
func board(after map[string][]byte, moves map[string]string) (string, error) {
	const bp = "SDP/KanBan/board.json"
	const lp = "SDP/ProjectManagement/Ledger.ndjson"
	project := "PROJECT"
	if b, ok := after[bp]; ok {
		var v map[string]any
		if e := json.Unmarshal(b, &v); e != nil {
			return "", fail("board", 2, "%v", e)
		}
		project, _ = v["projectId"].(string)
		if v["schemaVersion"] == "0.1" && v["ledger"] == "Ledger.ndjson" {
			source := "SDP/KanBan/Ledger.ndjson"
			history, ok := after[source]
			if !ok {
				return "", fail("history", 3, "missing board ledger")
			}
			if _, ok = after[lp]; ok {
				return "", fail("history", 3, "multiple management histories")
			}
			after[lp] = history
			delete(after, source)
			original := source
			for k, v := range moves {
				if v == source {
					original = k
					delete(moves, k)
					break
				}
			}
			moves[original] = lp
			after[bp], _ = Canonical(map[string]any{"schemaVersion": "0.2", "projectId": project, "namespaces": []string{project}, "ledger": "../ProjectManagement/Ledger.ndjson", "profile": "sdp-project-management/0.2"})
		} else if v["schemaVersion"] != "0.2" || v["ledger"] != "../ProjectManagement/Ledger.ndjson" {
			return "", fail("board", 2, "unsupported board")
		} else {
			if v["profile"] != "sdp-project-management/0.2" {
				return "", fail("board", 2, "unsupported management profile")
			}
			if _, ok := after[lp]; !ok {
				return "", fail("history", 3, "missing shared ledger")
			}
		}
	} else {
		after[bp], _ = Canonical(map[string]any{"schemaVersion": "0.2", "projectId": project, "namespaces": []string{project}, "ledger": "../ProjectManagement/Ledger.ndjson", "profile": "sdp-project-management/0.2"})
	}
	if !regexp.MustCompile(`^[A-Z][A-Z0-9]*$`).MatchString(project) {
		return "", fail("board", 2, "invalid project identity")
	}
	if _, ok := after[lp]; !ok {
		after[lp] = []byte{}
	}
	history := after[lp]
	if len(history) > 0 && history[len(history)-1] != '\n' {
		return "", fail("history", 3, "history must end with LF")
	}
	ids := map[string]bool{}
	for _, line := range bytes.Split(bytes.TrimSuffix(history, []byte{'\n'}), []byte{'\n'}) {
		if len(line) == 0 {
			continue
		}
		var event map[string]any
		if json.Unmarshal(line, &event) != nil {
			return "", fail("history", 2, "invalid event")
		}
		id, ok := event["eventId"].(string)
		if !ok || id == "" || ids[id] {
			return "", fail("history", 2, "duplicate/absent event identity")
		}
		ids[id] = true
	}
	if e := validateHistory(history, after); e != nil {
		return "", fail("history", 3, "%v", e)
	}
	return project, nil
}
func SavePlan(p Plan, destination string) error {
	q, e := filepath.Abs(destination)
	if e != nil {
		return e
	}
	rel, e := filepath.Rel(p.ProjectRoot, q)
	if e == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return fail("plan-output", 2, "save plan outside the project inventory")
	}
	if e = SafeAbsolute(q); e != nil {
		return e
	}
	b, e := Canonical(p)
	if e != nil {
		return e
	}
	if len(b) > RecordLimit {
		return fail("limit", 2, "plan too large")
	}
	f, e := os.OpenFile(q, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if e != nil {
		return e
	}
	_, e = f.Write(b)
	ce := f.Close()
	if e != nil {
		return e
	}
	return ce
}
func LoadPlan(file string) (Plan, error) {
	var p Plan
	b, e := Read(file, RecordLimit)
	if e != nil {
		return p, e
	}
	if e = Decode(b, RecordLimit, &p); e != nil {
		return p, e
	}
	if p.SchemaVersion != PlanSchema || p.PlanDigest != planHash(p) {
		return p, fail("plan", 4, "plan identity mismatch")
	}
	return p, nil
}
