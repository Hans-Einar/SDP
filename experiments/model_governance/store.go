// Package probe is disposable design evidence, NOT the production ModelGovernance engine.
package probe

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type Files map[string]string

// Record is an executable subset of the proposed schema; timestamps, authority,
// profile selection, limits, and full metadata hashing are deliberately omitted.
type Record struct {
	ID        string            `yaml:"id"`
	Parents   []string          `yaml:"parents"`
	Message   string            `yaml:"message"`
	Digest    string            `yaml:"sourceDigest"`
	Inventory map[string]string `yaml:"inventory"`
	Deleted   []string          `yaml:"deleted"`
}
type Head struct {
	Schema string `yaml:"schema"`
	ID     string `yaml:"id"`
	Seq    int    `yaml:"seq"`
}
type Store struct{ Root string }

func sum(b []byte) string { h := sha256.Sum256(b); return hex.EncodeToString(h[:]) }
func inventory(f Files) map[string]string {
	out := map[string]string{}
	for p, b := range f {
		out[p] = sum([]byte(b))
	}
	return out
}
func digest(f Files) string {
	h := sha256.New()
	h.Write([]byte("SDP-model-sources/1\x00"))
	names := []string{}
	for p := range f {
		names = append(names, p)
	}
	sort.Strings(names)
	for _, p := range names {
		var n [8]byte
		binary.BigEndian.PutUint64(n[:], uint64(len(p)))
		h.Write(n[:])
		h.Write([]byte(p))
		b := sha256.Sum256([]byte(f[p]))
		h.Write(b[:])
	}
	return hex.EncodeToString(h.Sum(nil))
}
func writeYAML(path string, v any) error {
	b, e := yaml.Marshal(v)
	if e != nil {
		return e
	}
	if e = os.MkdirAll(filepath.Dir(path), 0700); e != nil {
		return e
	}
	return os.WriteFile(path, b, 0600)
}
func readYAML(path string, v any) error {
	b, e := os.ReadFile(path)
	if e != nil {
		return e
	}
	d := yaml.NewDecoder(bytes.NewReader(b))
	d.KnownFields(true)
	return d.Decode(v)
}
func (s Store) head() (Head, error) {
	var h Head
	e := readYAML(filepath.Join(s.Root, "head.yaml"), &h)
	return h, e
}
func (s Store) dir(n int) string { return filepath.Join(s.Root, ".commits", fmt.Sprintf("#%05d", n)) }
func writeFiles(root string, f Files) error {
	for p, b := range f {
		if p == "" || filepath.IsAbs(p) || strings.Contains(p, "\\") || filepath.ToSlash(filepath.Clean(p)) != p || p == ".." || strings.HasPrefix(p, "../") {
			return errors.New("unsafe path")
		}
		dest := filepath.Join(root, filepath.FromSlash(p))
		if e := os.MkdirAll(filepath.Dir(dest), 0700); e != nil {
			return e
		}
		if e := os.WriteFile(dest, []byte(b), 0600); e != nil {
			return e
		}
	}
	return nil
}
func New(root, id string, f Files) (Store, error) {
	s := Store{root}
	if _, e := os.Stat(root); !os.IsNotExist(e) {
		return s, errors.New("destination exists")
	}
	if e := os.MkdirAll(filepath.Join(root, ".merge"), 0700); e != nil {
		return s, e
	}
	if e := writeYAML(filepath.Join(root, "head.yaml"), Head{"mg-probe/0.1", id, -1}); e != nil {
		return s, e
	}
	_, e := s.Commit(f, "baseline", nil, false)
	return s, e
}
func (s Store) Read(n int) (Files, []Record, error) {
	f := Files{}
	records := []Record{}
	h, e := s.head()
	if e != nil {
		return nil, nil, e
	}
	if n < 0 || n > h.Seq {
		return nil, nil, errors.New("not a published commit")
	}
	for i := 0; i <= n; i++ {
		var r Record
		if e = readYAML(filepath.Join(s.dir(i), "commit.yaml"), &r); e != nil {
			return nil, nil, e
		}
		if r.ID != fmt.Sprintf("%s:%05d", h.ID, i) {
			return nil, nil, errors.New("record identity mismatch")
		}
		for _, p := range r.Deleted {
			delete(f, p)
		}
		for p, want := range r.Inventory {
			b, err := os.ReadFile(filepath.Join(s.dir(i), "files", filepath.FromSlash(p)))
			if err == nil {
				if sum(b) != want {
					return nil, nil, errors.New("corrupt after-image")
				}
				f[p] = string(b)
			} else if !os.IsNotExist(err) {
				return nil, nil, err
			} else {
				old, exists := f[p]
				if !exists || sum([]byte(old)) != want {
					return nil, nil, errors.New("missing content")
				}
			}
		}
		if len(f) != len(r.Inventory) || digest(f) != r.Digest {
			return nil, nil, errors.New("snapshot mismatch")
		}
		records = append(records, r)
	}
	return f, records, nil
}

// Commit stages new payload first. Probe has a single writer and a fixed pending
// record; production needs writer locking, full recovery journal and source edits checks.
func (s Store) Commit(f Files, msg string, extra []string, failBeforeHead bool) (int, error) {
	if _, e := os.Stat(filepath.Join(s.Root, "pending.yaml")); e == nil {
		return 0, errors.New("recovery required")
	}
	h, e := s.head()
	if e != nil {
		return 0, e
	}
	before := Files{}
	parents := []string{}
	if h.Seq >= 0 {
		before, _, e = s.Read(h.Seq)
		if e != nil {
			return 0, e
		}
		parents = append(parents, fmt.Sprintf("%s:%05d", h.ID, h.Seq))
		if digest(f) == digest(before) && len(extra) == 0 {
			return h.Seq, nil
		}
	}
	parents = append(parents, extra...)
	n := h.Seq + 1
	changed := Files{}
	deleted := []string{}
	for p, b := range f {
		old, exists := before[p]
		if !exists || old != b {
			changed[p] = b
		}
	}
	for p := range before {
		if _, ok := f[p]; !ok {
			deleted = append(deleted, p)
		}
	}
	sort.Strings(deleted)
	if _, e = os.Stat(s.dir(n)); !os.IsNotExist(e) {
		return 0, errors.New("orphan payload requires recovery")
	}
	if e = writeFiles(filepath.Join(s.dir(n), "files"), changed); e != nil {
		return 0, e
	}
	r := Record{fmt.Sprintf("%s:%05d", h.ID, n), parents, msg, digest(f), inventory(f), deleted}
	if e = writeYAML(filepath.Join(s.dir(n), "commit.yaml"), r); e != nil {
		return 0, e
	}
	h.Seq = n
	if e = writeYAML(filepath.Join(s.Root, "pending.yaml"), h); e != nil {
		return 0, e
	}
	if failBeforeHead {
		return n, errors.New("injected stop before head publication")
	}
	return n, s.Resume()
}
func (s Store) Resume() error {
	var pending Head
	if e := readYAML(filepath.Join(s.Root, "pending.yaml"), &pending); e != nil {
		return e
	}
	h, e := s.head()
	if e != nil {
		return e
	}
	if pending.ID != h.ID || pending.Seq != h.Seq+1 {
		return errors.New("stale pending head")
	}
	var r Record
	if e = readYAML(filepath.Join(s.dir(pending.Seq), "commit.yaml"), &r); e != nil {
		return e
	}
	// Validate staged complete state before publishing, using a private copy of the store.
	tmp, e := os.MkdirTemp("", "mg-resume-")
	if e != nil {
		return e
	}
	defer os.RemoveAll(tmp)
	if e = os.CopyFS(tmp, os.DirFS(s.Root)); e != nil {
		return e
	}
	if e = writeYAML(filepath.Join(tmp, "head.yaml"), pending); e != nil {
		return e
	}
	if _, _, e = (Store{tmp}).Read(pending.Seq); e != nil {
		return e
	}
	return os.Rename(filepath.Join(s.Root, "pending.yaml"), filepath.Join(s.Root, "head.yaml"))
}
func (s Store) Restore(n int, dirty Files) (int, error) {
	target, _, e := s.Read(n)
	if e != nil {
		return 0, e
	}
	if _, e = s.Commit(dirty, "pre-restore checkpoint", nil, false); e != nil {
		return 0, e
	}
	return s.Commit(target, fmt.Sprintf("restore %d", n), nil, false)
}
func Merge(base, a, b Files) (Files, []string) {
	out := Files{}
	conflicts := []string{}
	keys := map[string]bool{}
	for _, f := range []Files{base, a, b} {
		for p := range f {
			keys[p] = true
		}
	}
	for p := range keys {
		v0, o0 := base[p]
		va, oa := a[p]
		vb, ob := b[p]
		switch {
		case va == vb && oa == ob:
			if oa {
				out[p] = va
			}
		case va == v0 && oa == o0:
			if ob {
				out[p] = vb
			}
		case vb == v0 && ob == o0:
			if oa {
				out[p] = va
			}
		default:
			conflicts = append(conflicts, p)
		}
	}
	sort.Strings(conflicts)
	return out, conflicts
}

type Frozen struct {
	Schema                  string   `yaml:"schema"`
	ID                      string   `yaml:"id"`
	Digest                  string   `yaml:"sourceDigest"`
	RestorePayloadAvailable bool     `yaml:"restorePayloadAvailable"`
	Lineage                 []Record `yaml:"lineage"`
}

func Promote(dest, id string, f Files, records []Record) error {
	if _, e := os.Stat(dest); !os.IsNotExist(e) {
		return errors.New("destination exists")
	}
	seen := map[string]Record{}
	unique := []Record{}
	for _, r := range records {
		if old, ok := seen[r.ID]; ok {
			a, _ := yaml.Marshal(old)
			b, _ := yaml.Marshal(r)
			if !bytes.Equal(a, b) {
				return errors.New("conflicting lineage identity")
			}
			continue
		}
		seen[r.ID] = r
		unique = append(unique, r)
	}
	for _, r := range unique {
		for _, parent := range r.Parents {
			if _, ok := seen[parent]; !ok {
				return errors.New("missing lineage parent")
			}
		}
	}
	if e := writeFiles(filepath.Join(dest, "sources"), f); e != nil {
		return e
	}
	return writeYAML(filepath.Join(dest, "artifact.yaml"), Frozen{"mg-probe/0.1", id, digest(f), false, unique})
}
