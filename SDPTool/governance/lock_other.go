//go:build !linux

package governance

func lock(dir string) (func(), error) {
	return nil, fail("incompatible_profile", "PGI1 writes require Linux; reads are available")
}

func sharedLock(dir string) (func(), error) {
	return lock(dir)
}
