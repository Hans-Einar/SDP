package install

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Snapshot includes the process tree, managed skills and incoming Markdown link
// scope. Exclusions are deliberate and surfaced in every plan.
var Exclusions = []string{".git", "node_modules", ".venv", "vendor", "build", ".cache"}

type Tree struct {
	Snapshot map[string]Observation
	Files    map[string][]byte
}

func Inspect(root string, extra []string) (Tree, error) {
	t := Tree{map[string]Observation{}, map[string][]byte{}}
	seen := map[string]string{}
	total := 0
	add := func(rel string, info fs.FileInfo) error {
		if e := Relative(rel); e != nil {
			return e
		}
		if info.Mode()&os.ModeSymlink != 0 || (!info.Mode().IsRegular() && !info.IsDir()) {
			return fail("path", 2, "unsupported object: %s", rel)
		}
		fold := strings.ToLower(rel)
		if old, ok := seen[fold]; ok && old != rel {
			return fail("path", 2, "case collision: %s", rel)
		}
		seen[fold] = rel
		if info.IsDir() {
			t.Snapshot[rel] = Observation{Type: "directory"}
			return nil
		}
		b, e := Read(filepath.Join(root, filepath.FromSlash(rel)), FileLimit)
		if e != nil {
			return e
		}
		total += len(b)
		if total > PayloadLimit {
			return fail("limit", 2, "observed file bytes exceed limit")
		}
		h := Hash(b)
		t.Snapshot[rel] = Observation{Type: "file", SHA256: &h}
		t.Files[rel] = b
		return nil
	}
	e := filepath.WalkDir(root, func(p string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if p == root {
			return nil
		}
		rel, _ := filepath.Rel(root, p)
		rel = filepath.ToSlash(rel)
		if rel == Operations || rel == "SDP/.sdp-backups" {
			if d.Type()&os.ModeSymlink != 0 {
				return fail("path", 2, "linked administrative directory")
			}
			if !d.IsDir() {
				return fail("path", 2, "administrative directory required")
			}
			return filepath.SkipDir
		}
		process := rel == "SDP" || strings.HasPrefix(rel, "SDP/") || rel == ".codex/skills" || strings.HasPrefix(rel, ".codex/skills/")
		if !process {
			for _, excluded := range Exclusions {
				if d.Name() == excluded {
					if d.IsDir() {
						return filepath.SkipDir
					}
					return nil
				}
			}
		}
		selected := process || strings.HasSuffix(strings.ToLower(rel), ".md")
		for _, p := range extra {
			if rel == p || strings.HasPrefix(p, rel+"/") {
				selected = true
			}
		}
		if selected {
			info, e := d.Info()
			if e != nil {
				return e
			}
			if e = add(rel, info); e != nil {
				return e
			}
		}
		if len(t.Snapshot) > EntryLimit {
			return fail("limit", 2, "too many observed paths")
		}
		return nil
	})
	if e != nil {
		return t, e
	}
	// Include ancestors of every selected file, independently of whether it was
	// selected by Markdown scope or explicitly by an adoption inventory.
	selected := make([]string, 0, len(t.Snapshot))
	for rel := range t.Snapshot {
		selected = append(selected, rel)
	}
	for _, rel := range selected {
		for parent := filepath.ToSlash(filepath.Dir(rel)); parent != "."; parent = filepath.ToSlash(filepath.Dir(parent)) {
			if _, ok := t.Snapshot[parent]; ok {
				continue
			}
			info, e := os.Lstat(filepath.Join(root, parent))
			if e != nil {
				return t, e
			}
			if e = add(parent, info); e != nil {
				return t, e
			}
		}
	}
	for _, p := range extra {
		if e := Relative(p); e != nil {
			return t, e
		}
		if e := SafeAbsolute(filepath.Join(root, p)); e != nil {
			return t, e
		}
		if _, ok := t.Snapshot[p]; !ok {
			_, e := os.Lstat(filepath.Join(root, p))
			if !os.IsNotExist(e) {
				return t, fail("path", 2, "unobserved destination %s", p)
			}
			t.Snapshot[p] = Observation{Type: "absent"}
		}
	}
	return t, nil
}
func SameSnapshot(a, b map[string]Observation) bool {
	x, _ := Canonical(a)
	y, _ := Canonical(b)
	return string(x) == string(y)
}
func Pending(root string) ([]string, error) {
	p := filepath.Join(root, Operations)
	if e := SafeAbsolute(p); e != nil {
		return nil, e
	}
	entries, e := os.ReadDir(p)
	if os.IsNotExist(e) {
		return []string{}, nil
	}
	if e != nil {
		return nil, e
	}
	out := []string{}
	for _, entry := range entries {
		if entry.Name() == "install.lock" || strings.HasPrefix(entry.Name(), "temp-") {
			continue
		}
		if !operationRE.MatchString(entry.Name()) || !entry.IsDir() {
			return nil, fail("journal", 2, "unknown operation entry %s", entry.Name())
		}
		b, e := Read(filepath.Join(p, entry.Name(), "journal.json"), RecordLimit)
		if e != nil {
			return nil, fail("pending", 5, "unreadable operation %s: %v", entry.Name(), e)
		}
		// Legacy journals are recognized, but never resumed by this engine.
		var v map[string]any
		if e = decodeEnvelope(b, &v); e != nil {
			return nil, e
		}
		if v["operationId"] != entry.Name() {
			return nil, fail("journal", 2, "operation identity mismatch")
		}
		if v["schemaVersion"] != JournalSchema && v["schemaVersion"] != "2.0" {
			return nil, fail("unsupported", 2, "journal schema")
		}
		if v["status"] != "completed" {
			out = append(out, entry.Name())
		}
	}
	sort.Strings(out)
	return out, nil
}
func describeError(e error) error {
	if _, ok := e.(*Error); ok {
		return e
	}
	return fail("io", 4, "%v", e)
}
func observationHash(s map[string]Observation, p string) *string { return s[p].SHA256 }
func equalHash(a, b *string) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}
func decodeEnvelope(b []byte, v *map[string]any) error {
	// Full typed journal validation belongs to recovery. Discovery must also read
	// the legacy schema without pretending that its fields use the new contract.
	var node map[string]any
	if e := json.Unmarshal(b, &node); e != nil {
		return fail("journal", 2, "%v", e)
	}
	*v = node
	return nil
}
func badObservation(p string, o Observation) error {
	if e := Relative(p); e != nil {
		return e
	}
	switch o.Type {
	case "file":
		if o.SHA256 == nil || !digestRE.MatchString(*o.SHA256) {
			return fmt.Errorf("invalid observed hash")
		}
	case "directory", "absent":
		if o.SHA256 != nil {
			return fmt.Errorf("unexpected observed hash")
		}
	default:
		return fmt.Errorf("unknown observation type")
	}
	return nil
}
