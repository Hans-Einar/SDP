//go:build !linux

package projecthistory

func lock(string) (func(), error) { return nil, ErrUnsupported }
