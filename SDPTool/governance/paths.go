package governance

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

func relative(p string) bool {
	if p == "" || p == "." || filepath.IsAbs(p) || filepath.ToSlash(filepath.Clean(p)) != p || strings.ContainsAny(p, "\\:") || strings.IndexFunc(p, unicode.IsControl) >= 0 {
		return false
	}
	for _, s := range strings.Split(p, "/") {
		if s == ".." || strings.HasPrefix(s, ".") {
			return false
		}
	}
	return true
}

// Reject symlinks in every existing component, including the state directory.
func safeAbsolute(p string) error {
	if !filepath.IsAbs(p) {
		return fail("path", "absolute path required")
	}
	cur := string(filepath.Separator)
	for _, part := range strings.Split(strings.TrimPrefix(filepath.Clean(p), cur), cur) {
		cur = filepath.Join(cur, part)
		st, e := os.Lstat(cur)
		if os.IsNotExist(e) {
			continue
		}
		if e != nil {
			return e
		}
		if st.Mode()&os.ModeSymlink != 0 {
			return fail("path", "symlink is not permitted: "+cur)
		}
	}
	return nil
}
func readFile(p string, limit int64) ([]byte, error) {
	if e := safeAbsolute(p); e != nil {
		return nil, e
	}
	st, e := os.Lstat(p)
	if e != nil {
		return nil, e
	}
	if !st.Mode().IsRegular() {
		return nil, fail("path", "regular file required")
	}
	f, e := os.Open(p)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	actual, e := f.Stat()
	if e != nil {
		return nil, e
	}
	if !os.SameFile(st, actual) {
		return nil, fail("stale_revision", "file changed while opening")
	}
	b, e := io.ReadAll(io.LimitReader(f, limit+1))
	if e != nil {
		return nil, e
	}
	if int64(len(b)) > limit {
		return nil, fail("limit", "file exceeds input limit")
	}
	return b, nil
}
