//go:build !bundled

package sdptool

// Source builds have no desktop payload. Release builds use package.sh.
var bundledNativeHost []byte
