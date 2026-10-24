package governance

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const EvidenceSchema = "sdp-governance-evidence/1"
const evidenceStateSchema = "sdp-project-governance-state/0.10"

type CandidateEntry struct {
	Path       string `json:"path"`
	Kind       string `json:"kind"`
	SHA256     string `json:"sha256"`
	Executable uint32 `json:"executable"`
}

// Candidate identifies selected bytes, not the completeness of the selected
// dependency boundary. No Git repository or model is required.
type Candidate struct {
	Schema        string           `json:"schema"`
	ProjectRef    string           `json:"projectRef"`
	CheckoutID    string           `json:"checkoutId"`
	RunID         string           `json:"runId"`
	InputRevision string           `json:"inputRevision"`
	WritablePaths []string         `json:"writablePaths"`
	ExcludedPaths []string         `json:"excludedPaths"`
	Dependencies  []string         `json:"dependencies"`
	Entries       []CandidateEntry `json:"entries"`
	Digest        string           `json:"digest"`
}

func candidateDigest(c Candidate) string { c.Digest = ""; return requestDigest(c) }
func within(path, root string) bool      { return path == root || strings.HasPrefix(path, root+"/") }
func literalCandidatePath(p string) bool {
	return relative(p) && len(p) <= 4096 && !strings.ContainsAny(p, "*?[]")
}

func (s *Service) CaptureCandidate(runID string, dependencies []string) (Candidate, error) {
	state, head, e := s.load()
	if e != nil {
		return Candidate{}, e
	}
	r, ok := state.Runs[runID]
	if !ok || !uuidPattern.MatchString(runID) {
		return Candidate{}, fail("missing_coverage", "run not found")
	}
	if e = s.currentAttentionInputs(r); e != nil {
		return Candidate{}, e
	}
	c, e := s.captureCandidate(state.CheckoutID, r, dependencies)
	if e != nil {
		return Candidate{}, e
	}
	if e = s.currentAttentionInputs(r); e != nil {
		return Candidate{}, e
	}
	_, latest, e := s.load()
	if e != nil {
		return Candidate{}, e
	}
	if latest != head {
		return Candidate{}, fail("stale_revision", "run store changed during candidate inspection")
	}
	return c, nil
}

func (s *Service) captureCandidate(checkout string, r Run, dependencies []string) (Candidate, error) {
	c := Candidate{Schema: EvidenceSchema, ProjectRef: s.project, CheckoutID: checkout, RunID: r.ID, InputRevision: r.InputRevision,
		WritablePaths: append([]string{}, r.Assignment.WritablePaths...), ExcludedPaths: append([]string{}, r.Assignment.ExcludedPaths...),
		Dependencies: append([]string{}, dependencies...), Entries: []CandidateEntry{}}
	if len(dependencies) > 64 {
		return c, fail("limit", "at most 64 explicit dependencies")
	}
	slices.Sort(c.Dependencies)
	for i, p := range c.Dependencies {
		if i > 0 && p == c.Dependencies[i-1] {
			return c, fail("invalid_request", "duplicate dependency")
		}
	}
	for _, p := range append(append(append([]string{}, c.WritablePaths...), c.ExcludedPaths...), c.Dependencies...) {
		if !literalCandidatePath(p) {
			return c, fail("path", "candidate paths must be bounded literal project paths")
		}
	}
	seen := map[string]bool{}
	var total int
	visit := func(path string, info fs.FileInfo) error {
		if seen[path] {
			return nil
		}
		seen[path] = true
		if len(seen) > 512 {
			return fail("limit", "candidate exceeds 512 entries")
		}
		if !literalCandidatePath(path) {
			return fail("path", "unsupported candidate path")
		}
		entry := CandidateEntry{Path: path, Kind: "missing"}
		if info != nil {
			switch {
			case info.IsDir():
				entry.Kind = "directory"
			case info.Mode().IsRegular():
				b, err := readFile(filepath.Join(s.root, path), 4<<20)
				if err != nil {
					return err
				}
				total += len(b)
				if total > 32<<20 {
					return fail("limit", "candidate exceeds 32 MiB")
				}
				entry.Kind, entry.SHA256, entry.Executable = "file", hash(b), uint32(info.Mode().Perm()&0111)
			default:
				return fail("path", "candidate symlinks and special files are unsupported")
			}
		}
		c.Entries = append(c.Entries, entry)
		return nil
	}
	excluded := func(p string) bool {
		// A write exclusion is not a declaration of behavioral irrelevance.
		// Explicit read-only dependencies participate without granting writes.
		for _, dep := range c.Dependencies {
			if within(p, dep) {
				return false
			}
		}
		for _, ex := range c.ExcludedPaths {
			if within(p, ex) {
				return true
			}
		}
		return false
	}
	// Read directories in bounded batches. filepath.Walk reads a whole directory
	// before invoking its callback, which defeats the inventory's memory bound.
	var walk func(string, fs.FileInfo) error
	walk = func(path string, info fs.FileInfo) error {
		if excluded(path) || seen[path] {
			return nil
		}
		if err := visit(path, info); err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}
		abs := filepath.Join(s.root, path)
		if err := safeAbsolute(abs); err != nil {
			return err
		}
		f, err := os.Open(abs)
		if err != nil {
			return err
		}
		defer f.Close()
		opened, err := f.Stat()
		if err != nil {
			return err
		}
		if !os.SameFile(info, opened) {
			return fail("stale_revision", "directory changed while opening")
		}
		for {
			batch, err := f.ReadDir(64)
			if err != nil && err != io.EOF {
				return err
			}
			for _, entry := range batch {
				child := path + "/" + entry.Name()
				if excluded(child) {
					continue
				}
				childInfo, e := entry.Info()
				if e != nil {
					return e
				}
				if e = walk(child, childInfo); e != nil {
					return e
				}
			}
			if err == io.EOF {
				break
			}
		}
		return nil
	}
	for _, path := range append(append([]string{}, c.WritablePaths...), c.Dependencies...) {
		if excluded(path) {
			continue
		}
		abs := filepath.Join(s.root, path)
		if err := safeAbsolute(abs); err != nil {
			return c, err
		}
		info, err := os.Lstat(abs)
		if os.IsNotExist(err) {
			if err = visit(path, nil); err != nil {
				return c, err
			}
			continue
		} else if err != nil {
			return c, err
		}
		err = walk(path, info)
		if err != nil {
			return c, err
		}
	}
	slices.SortFunc(c.Entries, func(a, b CandidateEntry) int { return strings.Compare(a.Path, b.Path) })
	c.Digest = candidateDigest(c)
	return c, nil
}
