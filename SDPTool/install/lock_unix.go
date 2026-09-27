//go:build linux || darwin || freebsd

package install

import (
	"os"
	"path/filepath"
	"syscall"
)

func lock(root string) (func(), error) {
	cache, e := os.UserCacheDir()
	if e != nil {
		return nil, e
	}
	dir := filepath.Join(cache, "sdptool", "locks")
	if e = SafeAbsolute(dir); e != nil {
		return nil, e
	}
	if e = os.MkdirAll(dir, 0700); e != nil {
		return nil, e
	}
	p := filepath.Join(dir, Hash([]byte(root))+".lock")
	if e = SafeAbsolute(p); e != nil {
		return nil, e
	}
	f, e := os.OpenFile(p, os.O_CREATE|os.O_RDWR, 0600)
	if e != nil {
		return nil, e
	}
	if e = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); e != nil {
		f.Close()
		return nil, fail("locked", 5, "project locked by another installer")
	}
	return func() { syscall.Flock(int(f.Fd()), syscall.LOCK_UN); f.Close() }, nil
}
