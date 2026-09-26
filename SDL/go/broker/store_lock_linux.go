//go:build linux

package broker

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

func (s *Store) lock() error {
	f, e := os.OpenFile(filepath.Join(s.root, "store.lock"), os.O_CREATE|os.O_RDWR, 0600)
	if e != nil {
		return e
	}
	if e = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); e != nil {
		f.Close()
		return fmt.Errorf("artifact store already in use")
	}
	s.lockFile = f
	return nil
}
func (s *Store) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.lockFile != nil {
		syscall.Flock(int(s.lockFile.Fd()), syscall.LOCK_UN)
		s.lockFile.Close()
		s.lockFile = nil
	}
}
func (s *Store) recover() error {
	known := map[string]bool{}
	for _, a := range s.items {
		known[a.Directory] = true
	}
	entries, e := os.ReadDir(s.root)
	if e != nil {
		return e
	}
	for _, d := range entries {
		if d.IsDir() && !known[d.Name()] && (safeToken(d.Name()) || strings.HasPrefix(d.Name(), ".sdl-publish-") || strings.HasSuffix(d.Name(), ".lock")) {
			if e = os.RemoveAll(filepath.Join(s.root, d.Name())); e != nil {
				return e
			}
		}
	}
	return nil
}
