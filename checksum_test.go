package osutil

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func checkTearError(t *testing.T, err error) {
	if !assert.NoError(t, err) {
		t.FailNow()
	}
}

func TestChecksumForPath(t *testing.T) {
	if IsWindows() {
		t.SkipNow() // TODO Implement symlink handling under Windows or make this test case compatible with Windows. https://$INTERNAL/symflower/symflower/-/issues/3637
	}

	temporaryPath, err := os.MkdirTemp("", "TestChecksumForPath")
	checkTearError(t, err)
	defer func() {
		checkTearError(t, os.RemoveAll(temporaryPath))
	}()

	checkTearError(t, os.Symlink(".", filepath.Join(temporaryPath, "symlinkToDirectory")))
	digestForEmptyDirectory, err := ChecksumForPath(temporaryPath)
	assert.NoError(t, err)
	assert.NotEmpty(t, digestForEmptyDirectory)

	checkTearError(t, os.Symlink("broken target", filepath.Join(temporaryPath, "symlink with broken target")))
	_, err = ChecksumForPath(temporaryPath)
	assert.NoError(t, err)

	checkTearError(t, os.WriteFile(filepath.Join(temporaryPath, "some file"), []byte{}, 0600))
	checkTearError(t, os.WriteFile(filepath.Join(temporaryPath, "symlink or file"), []byte{}, 0600))
	digestWithFile, err := ChecksumForPath(temporaryPath)
	assert.NoError(t, err)

	checkTearError(t, os.Remove(filepath.Join(temporaryPath, "symlink or file")))
	checkTearError(t, os.Symlink("some file", filepath.Join(temporaryPath, "symlink or file")))
	digestWithSymlink, err := ChecksumForPath(temporaryPath)
	assert.NoError(t, err)

	assert.NotEqual(t, digestForEmptyDirectory, digestWithFile)
	assert.Equal(t, digestWithFile, digestWithSymlink)
}
