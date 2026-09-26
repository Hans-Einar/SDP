//go:build !linux && !darwin && !freebsd

package install

func lock(root string) (func(), error) {
	return nil, fail("unsupported", 2, "installation locking has not been verified for this platform")
}
