//go:build linux

package governance

import (
	"os"
	"path/filepath"
	"syscall"
)

func lock(dir string) (func(), error) {
	if e := safeAbsolute(dir); e != nil {
		return nil, e
	}
	if e := mkdirDurable(dir); e != nil {
		return nil, e
	}
	st, e := os.Stat(dir)
	if e != nil {
		return nil, e
	}
	if st.Mode().Perm()&0077 != 0 {
		return nil, fail("path", "state directory must be private (0700)")
	}
	f, e := os.OpenFile(filepath.Join(dir, "lock"), os.O_CREATE|os.O_RDWR|syscall.O_NOFOLLOW, 0600)
	if e != nil {
		return nil, e
	}
	if e = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); e != nil {
		f.Close()
		return nil, fail("busy", "another process owns the checkout state")
	}
	return func() { syscall.Flock(int(f.Fd()), syscall.LOCK_UN); f.Close() }, nil
}
