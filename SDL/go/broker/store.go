package broker

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/Hans-Einar/SDP/SystemDesignLanguage/go/documents"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

type Artifact struct {
	Key, Directory string
	Bytes          int64
	Touched        int64
	Leases         map[string]bool
}
type Store struct {
	lockFile  *os.File
	mu        sync.Mutex
	root      string
	quota     int64
	limit     int
	items     map[string]*Artifact
	sequences map[string]uint64
}

func token() string {
	var b [16]byte
	if _, e := rand.Read(b[:]); e != nil {
		panic(e)
	}
	return hex.EncodeToString(b[:])
}
func NewStore(root string, quota int64, limit int) (*Store, error) {
	if quota <= 0 || limit <= 0 {
		return nil, fmt.Errorf("invalid store limits")
	}
	if e := os.MkdirAll(root, 0700); e != nil {
		return nil, e
	}
	info, e := os.Lstat(root)
	if e != nil || !info.IsDir() || info.Mode().Perm() != 0700 {
		return nil, fmt.Errorf("artifact root must be private 0700 directory")
	}
	s := &Store{root: root, quota: quota, limit: limit, items: map[string]*Artifact{}}
	if e = s.lock(); e != nil {
		return nil, e
	}
	success := false
	defer func() {
		if !success {
			s.Close()
		}
	}()
	if b, e := os.ReadFile(filepath.Join(root, "store.json")); e == nil {
		var index struct {
			Artifacts map[string]*Artifact
			Sequences map[string]uint64
		}
		if e = json.Unmarshal(b, &index); e != nil {
			return nil, e
		}
		s.items = index.Artifacts
		s.sequences = index.Sequences
		if s.items == nil || s.sequences == nil {
			return nil, fmt.Errorf("invalid artifact index version")
		}
		for k, a := range s.items {
			if a.Key != k || !safeToken(k) || !safeToken(a.Directory) {
				return nil, fmt.Errorf("invalid artifact index")
			}
			if _, e = os.Stat(filepath.Join(root, a.Directory, "entry.md")); e != nil {
				return nil, fmt.Errorf("missing persisted artifact: %w", e)
			}
		}
	} else if !os.IsNotExist(e) {
		return nil, e
	}
	if s.sequences == nil {
		s.sequences = map[string]uint64{}
	}
	if e = s.recover(); e != nil {
		return nil, e
	}
	success = true
	return s, nil
}
func safeToken(s string) bool {
	if len(s) != 32 && len(s) != 64 {
		return false
	}
	_, e := hex.DecodeString(s)
	return e == nil
}
func (s *Store) save() error {
	b, e := json.MarshalIndent(struct {
		Artifacts map[string]*Artifact
		Sequences map[string]uint64
	}{s.items, s.sequences}, "", "  ")
	if e != nil {
		return e
	}
	p := filepath.Join(s.root, "store.tmp")
	if e = os.WriteFile(p, b, 0600); e != nil {
		return e
	}
	return os.Rename(p, filepath.Join(s.root, "store.json"))
}
func (s *Store) Lookup(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	a, ok := s.items[key]
	if !ok {
		return "", false
	}
	return filepath.Join(s.root, a.Directory, "entry.md"), true
}
func (s *Store) Put(key string, b *documents.Bundle) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !safeToken(key) {
		return "", fmt.Errorf("invalid artifact key")
	}
	if a, ok := s.items[key]; ok {
		return filepath.Join(s.root, a.Directory, "entry.md"), nil
	}
	var size int64
	for _, data := range b.Files {
		size += int64(len(data))
	}
	if size > s.quota {
		return "", fmt.Errorf("artifact exceeds byte quota")
	}
	if e := s.evict(size); e != nil {
		return "", e
	}
	a := &Artifact{key, token(), size, time.Now().UnixNano(), map[string]bool{}}
	dir := filepath.Join(s.root, a.Directory)
	if e := b.Publish(dir); e != nil {
		return "", e
	}
	s.items[key] = a
	if e := s.save(); e != nil {
		delete(s.items, key)
		os.RemoveAll(dir)
		return "", e
	}
	return filepath.Join(dir, "entry.md"), nil
}
func (s *Store) evict(extra int64) error {
	var bytes int64
	for _, a := range s.items {
		bytes += a.Bytes
	}
	candidates := []*Artifact{}
	for _, a := range s.items {
		if len(a.Leases) == 0 {
			candidates = append(candidates, a)
		}
	}
	sort.Slice(candidates, func(i, j int) bool { return candidates[i].Touched < candidates[j].Touched })
	for (bytes+extra > s.quota || len(s.items) >= s.limit) && len(candidates) > 0 {
		a := candidates[0]
		candidates = candidates[1:]
		if e := os.RemoveAll(filepath.Join(s.root, a.Directory)); e != nil {
			return e
		}
		delete(s.items, a.Key)
		bytes -= a.Bytes
	}
	if bytes+extra > s.quota || len(s.items) >= s.limit {
		return fmt.Errorf("artifact quota: all remaining bundles are leased")
	}
	return s.save()
}
func (s *Store) Acquire(key string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	a, ok := s.items[key]
	if !ok {
		return "", fmt.Errorf("artifact no longer cached")
	}
	total := 0
	for _, item := range s.items {
		total += len(item.Leases)
	}
	if total >= 1024 {
		return "", fmt.Errorf("reader lease limit")
	}
	lease := token()
	a.Leases[lease] = true
	a.Touched = time.Now().UnixNano()
	if e := s.save(); e != nil {
		delete(a.Leases, lease)
		return "", e
	}
	return lease, nil
}
func (s *Store) Release(lease string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, a := range s.items {
		if a.Leases[lease] {
			delete(a.Leases, lease)
			if e := s.save(); e != nil {
				a.Leases[lease] = true
				return e
			}
			return nil
		}
	}
	return fmt.Errorf("unknown lease")
}

// Sweep removes only released artifacts. Leases survive daemon restart and have
// no implicit timeout: an idle document can still reload its images.
func (s *Store) Sweep() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for k, a := range s.items {
		if len(a.Leases) == 0 {
			if e := os.RemoveAll(filepath.Join(s.root, a.Directory)); e != nil {
				return e
			}
			delete(s.items, k)
		}
	}
	return s.save()
}

func (s *Store) Begin(route string, seq uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	old := s.sequences[route]
	if seq <= old {
		return fmt.Errorf("stale request")
	}
	if old == 0 && len(s.sequences) >= 1024 {
		return fmt.Errorf("request route limit")
	}
	s.sequences[route] = seq
	if e := s.save(); e != nil {
		if old == 0 {
			delete(s.sequences, route)
		} else {
			s.sequences[route] = old
		}
		return e
	}
	return nil
}
func (s *Store) Latest(route string, seq uint64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.sequences[route] == seq
}
