package osutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDirExists(t *testing.T) {
	assert.NoError(t, DirExists("/"))
	assert.NoError(t, DirExists("../osutil"))

	assert.Error(t, DirExists("hey"))

	assert.Equal(t, ErrNotADirectory, DirExists("exists.go"))
}

func TestRemoveFileIfExists(t *testing.T) {
	t.Run("File", func(t *testing.T) {
		temporaryPath := t.TempDir()

		filePath := filepath.Join(temporaryPath, "plain.txt")
		require.NoError(t, os.WriteFile(filePath, nil, 0600))

		assert.NoError(t, RemoveFileIfExists(filePath))
		assert.NoFileExists(t, filePath)
	})

	t.Run("Empty directory", func(t *testing.T) {
		temporaryPath := t.TempDir()

		filePath := filepath.Join(temporaryPath, "plain")
		require.NoError(t, os.MkdirAll(filePath, 0700))

		assert.NoError(t, RemoveFileIfExists(filePath))
		assert.NoDirExists(t, filePath)
	})

	t.Run("Non-empty directory", func(t *testing.T) {
		temporaryPath := t.TempDir()

		filePath := filepath.Join(temporaryPath, "plain/subdir")
		require.NoError(t, os.MkdirAll(filePath, 0700))

		assert.ErrorIs(t, RemoveFileIfExists(filepath.Dir(filePath)), ErrDirectoryNotEmpty)
		assert.DirExists(t, filePath)
	})

	t.Run("Symlink", func(t *testing.T) {
		temporaryPath := t.TempDir()

		filePath := filepath.Join(temporaryPath, "plain.txt")
		require.NoError(t, os.WriteFile(filePath, nil, 0600))

		linkFilePath := filepath.Join(temporaryPath, "symlink")
		require.NoError(t, os.Symlink(filePath, linkFilePath))

		assert.NoError(t, RemoveFileIfExists(linkFilePath))
		assert.FileExists(t, filePath)
		assert.NoFileExists(t, linkFilePath)
	})
}
