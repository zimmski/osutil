package osutil

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirectoryPermissionOfParent(t *testing.T) {
	type testCase struct {
		// Name holds the name of the test case.
		Name string

		// Path holds the relative directory path that should be created and used.
		Path string
		// Permission holds the directory permission that should be applied and that we should read back.
		Permission fs.FileMode
	}

	validate := func(t *testing.T, tc testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			path, err := os.MkdirTemp("", strings.ReplaceAll(t.Name(), "/", "-"))
			assert.NoError(t, err)
			defer func() {
				assert.NoError(t, os.RemoveAll(path))
			}()

			p := path + "/" + tc.Path
			assert.NoError(t, os.MkdirAll(p, tc.Permission))
			assert.NoError(t, os.Chmod(filepath.Dir(p), tc.Permission)) // The "umask" usually removes some bits of the permission, we do want to have all of them.
			permission, err := DirectoryPermissionOfParent(p)
			assert.NoError(t, err)
			if IsWindows() {
				assert.Equal(t, fmt.Sprintf("%o", 0777), fmt.Sprintf("%o", permission)) // TODO Implement file permission handling for Windows. https://$INTERNAL/symflower/symflower/-/issues/3637
			} else {
				assert.Equal(t, fmt.Sprintf("%o", tc.Permission), fmt.Sprintf("%o", permission))
			}
		})
	}

	validate(t, testCase{
		Name: "User",

		Path:       "user/child",
		Permission: 0700,
	})

	validate(t, testCase{
		Name: "Group",

		Path:       "group/child",
		Permission: 0770,
	})

	validate(t, testCase{
		Name: "All",

		Path:       "all/child",
		Permission: 0777,
	})
}

func TestFilePermissionOfParent(t *testing.T) {
	type testCase struct {
		// Name holds the name of the test case.
		Name string

		// Path holds the relative file path that should be created and used.
		Path string
		// Permission holds the directory permission that should be applied.
		Permission fs.FileMode

		// FilePermission holds the file permission that we want to read back.
		FilePermission fs.FileMode
	}

	validate := func(t *testing.T, tc testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			path, err := os.MkdirTemp("", strings.ReplaceAll(t.Name(), "/", "-"))
			assert.NoError(t, err)
			defer func() {
				assert.NoError(t, os.RemoveAll(path))
			}()

			directoryPath := filepath.Dir(path + "/" + tc.Path)
			assert.NoError(t, os.MkdirAll(directoryPath, tc.Permission))
			assert.NoError(t, os.Chmod(directoryPath, tc.Permission)) // The "umask" usually removes some bits of the permission, we do want to have all of them.
			permission, err := FilePermissionOfParent(path + "/" + tc.Path)
			assert.NoError(t, err)
			if IsWindows() {
				assert.Equal(t, fmt.Sprintf("%o", 0666), fmt.Sprintf("%o", permission)) // TODO Implement file permission handling for Windows. https://$INTERNAL/symflower/symflower/-/issues/3637
			} else {
				assert.Equal(t, fmt.Sprintf("%o", tc.FilePermission), fmt.Sprintf("%o", permission))
			}
		})
	}

	validate(t, testCase{
		Name: "User",

		Path:       "user/child",
		Permission: 0700,

		FilePermission: 0600,
	})

	validate(t, testCase{
		Name: "Group",

		Path:       "group/child",
		Permission: 0770,

		FilePermission: 0660,
	})

	validate(t, testCase{
		Name: "All",

		Path:       "all/child",
		Permission: 0777,

		FilePermission: 0666,
	})
}
