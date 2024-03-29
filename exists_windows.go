package osutil

import (
	"syscall"
)

// ErrDirectoryNotEmpty indicates that a directory is not empty.
var ErrDirectoryNotEmpty = syscall.ERROR_DIR_NOT_EMPTY
