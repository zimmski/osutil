package osutil

import (
	"io/fs"
	"os"
	"path/filepath"
)

// DirectoryPermissionOfParent looks at parent directory of the given path and returns a directory permission based on the permission of the parent.
// The returned permission copies the read, write and execute permissions.
func DirectoryPermissionOfParent(path string) (permission fs.FileMode, err error) {
	path, err = filepath.Abs(path)
	if err != nil {
		return 0, err
	}
	var s fs.FileInfo
	for {
		if path != "/" {
			path = filepath.Dir(path)
		}

		s, err = os.Stat(path)
		if err != nil {
			continue
		}

		break
	}

	permission |= s.Mode() & 0777

	return permission, nil
}

// FilePermissionOfParent looks at parent directory of the given path and returns a file permission based on the permission of the parent.
// The returned permission copies the read and write permissions but not the execute permission.
func FilePermissionOfParent(path string) (permission fs.FileMode, err error) {
	path, err = filepath.Abs(path)
	if err != nil {
		return 0, err
	}
	var s fs.FileInfo
	for {
		if path != "/" {
			path = filepath.Dir(path)
		}

		s, err = os.Stat(path)
		if err != nil {
			continue
		}

		break
	}

	permission |= s.Mode() & 0666

	return permission, nil
}
