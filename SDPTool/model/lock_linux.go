//go:build linux

package model

import (
	"os"
	"path/filepath"
	"syscall"
)

func lockArea(area string) (func(), error) {
	p := filepath.Join(area, operations)
	if e := os.MkdirAll(p, 0700); e != nil {
		return nil, e
	}
	i, e := os.Lstat(p)
	if e != nil {
		return nil, e
	}
	if !i.IsDir() || i.Mode()&os.ModeSymlink != 0 {
		return nil, fail("path", "unsafe operation directory")
	}
	f, e := os.OpenFile(filepath.Join(p, "lock"), os.O_CREATE|os.O_RDWR|syscall.O_NOFOLLOW, 0600)
	if e != nil {
		return nil, e
	}
	if e = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); e != nil {
		f.Close()
		return nil, fail("busy", "another model writer owns area")
	}
	return func() { syscall.Flock(int(f.Fd()), syscall.LOCK_UN); f.Close() }, nil
}
