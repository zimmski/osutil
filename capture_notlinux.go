//go:build !linux || !cgo

package osutil

// Capture captures stderr and stdout of a given function call.
func Capture(call func()) (output []byte, err error) {
	panic("not implemented") // WORKAROUND Implement this function for MacOS and Windows when it is actual needed. Until then we can cross-compile even if the function is only mentioned in a package. https://$INTERNAL/symflower/symflower/-/issues/3575
}
