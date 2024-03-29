//go:build !windows

package osutil

import (
	"syscall"
)

// ErrDirectoryNotEmpty indicates that a directory is not empty.
var ErrDirectoryNotEmpty = syscall.ENOTEMPTY
