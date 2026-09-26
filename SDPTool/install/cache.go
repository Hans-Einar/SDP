package install

import (
	"bytes"
	"os"
	"path/filepath"
)

func inputCache() (string, error) {
	base := os.Getenv("SDP_CACHE_DIR")
	if base == "" {
		p, e := os.UserCacheDir()
		if e != nil {
			return "", e
		}
		base = filepath.Join(p, "sdptool", "distribution")
	}
	p, e := filepath.Abs(filepath.Join(base, "inputs"))
	if e != nil {
		return "", e
	}
	if e = SafeAbsolute(p); e != nil {
		return "", e
	}
	if e = os.MkdirAll(p, 0700); e != nil {
		return "", e
	}
	return p, nil
}
func createIdentical(p string, b []byte) error {
	if e := SafeAbsolute(p); e != nil {
		return e
	}
	if old, e := Read(p, RecordLimit); e == nil {
		if !bytes.Equal(old, b) {
			return fail("cache", 4, "unequal cache entry")
		}
		return nil
	} else if !os.IsNotExist(e) {
		return e
	}
	f, e := os.CreateTemp(filepath.Dir(p), ".pending-")
	if e != nil {
		return e
	}
	temp := f.Name()
	defer os.Remove(temp)
	if e = f.Chmod(0600); e == nil {
		_, e = f.Write(b)
	}
	if e == nil {
		e = f.Sync()
	}
	ce := f.Close()
	if e == nil {
		e = ce
	}
	if e != nil {
		return e
	}
	if e = os.Link(temp, p); e != nil {
		old, re := Read(p, RecordLimit)
		if re != nil || !bytes.Equal(old, b) {
			return fail("cache", 4, "cache publication conflict")
		}
	}
	return nil
}
func cacheInput(in Input) error {
	if _, e := descriptor(in); e != nil {
		return e
	}
	dir, e := inputCache()
	if e != nil {
		return e
	}
	in.Path = filepath.Join(dir, in.SHA256+".descriptor.json")
	if e = createIdentical(in.Path, in.Bytes); e != nil {
		return e
	}
	b, e := Canonical(in)
	if e != nil {
		return e
	}
	return createIdentical(filepath.Join(dir, in.SHA256+".input.json"), b)
}
func cachedInput(digest string) (Input, error) {
	var in Input
	if !digestRE.MatchString(digest) {
		return in, fail("digest", 4, "invalid receipt digest")
	}
	dir, e := inputCache()
	if e != nil {
		return in, e
	}
	b, e := Read(filepath.Join(dir, digest+".input.json"), RecordLimit)
	if e != nil {
		return in, e
	}
	if e = Decode(b, RecordLimit, &in); e != nil {
		return in, e
	}
	if in.SHA256 != digest {
		return in, fail("digest", 4, "cached input identity")
	}
	_, e = descriptor(in)
	return in, e
}
