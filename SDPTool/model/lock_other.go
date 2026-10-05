//go:build !linux

package model

func lockArea(area string) (func(), error) {
	return nil, fail("platform", "model mutations currently verified on Linux only; reads remain available")
}
