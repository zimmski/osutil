package osutil

import (
	"os"
)

// MkdirAll creates a directory named path, along with any necessary parents, and returns nil, or else returns an error.
func MkdirAll(path string) error {
	return os.MkdirAll(path, 0750)
}
