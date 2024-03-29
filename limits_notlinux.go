//go:build !linux

package osutil

// EnforceProcessTreeLimits constrains the current process and all descendant processes to the specified limits. The current process exits when the limits are exceeded.
func EnforceProcessTreeLimits(limits ProcessTreeLimits) {
	// WORKAROUND Implement this function for MacOS and Windows when it is actual needed. Until then we can cross-compile even if the function is only mentioned in a package. https://$INTERNAL/symflower/symflower/-/issues/3592
}
