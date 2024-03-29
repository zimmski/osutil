package osutil

import (
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilesRecursive(t *testing.T) {
	path, err := os.MkdirTemp("", "os-util")
	assert.NoError(t, err)

	assert.NoError(t, os.MkdirAll(filepath.Join(path, "a", "b"), 0750))
	assert.NoError(t, os.WriteFile(filepath.Join(path, "c.txt"), []byte("foobar"), 0640))
	assert.NoError(t, os.WriteFile(filepath.Join(path, "a", "d.txt"), []byte("foobar"), 0640))
	assert.NoError(t, os.WriteFile(filepath.Join(path, "a", "b", "e.txt"), []byte("foobar"), 0640))

	fs, err := FilesRecursive(path)
	assert.NoError(t, err)

	sort.Strings(fs)

	assert.Equal(
		t,
		[]string{
			filepath.Join(path, "a", "b", "e.txt"),
			filepath.Join(path, "a", "d.txt"),
			filepath.Join(path, "c.txt"),
		},
		fs,
	)
}

func TestFilePathsByHierarchy(t *testing.T) {
	type testCase struct {
		Name string

		ExpectedSortedFilePaths []string
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			for i := 0; i < 1000; i++ {
				actualSortedFilePaths := make([]string, len(tc.ExpectedSortedFilePaths))
				copy(actualSortedFilePaths, tc.ExpectedSortedFilePaths)

				rand.Shuffle(len(actualSortedFilePaths), func(i, j int) {
					actualSortedFilePaths[i], actualSortedFilePaths[j] = actualSortedFilePaths[j], actualSortedFilePaths[i]
				})
				sort.Sort(FilePathsByHierarchy(actualSortedFilePaths))

				assert.Equal(t, tc.ExpectedSortedFilePaths, actualSortedFilePaths)
			}
		})
	}

	validate(t, &testCase{
		ExpectedSortedFilePaths: []string{
			"a",
			"a b",
			"ab",
			"b",
			filepath.Join("a", "b"),
			filepath.Join("b c", "c"),
			filepath.Join("a", "b", " "),
			filepath.Join("a", "b", "c"),
		},
	})
}

func TestAppendtoFile(t *testing.T) {
	path, err := os.MkdirTemp("", "os-util")
	assert.NoError(t, err)

	file := filepath.Join(path, "test.log")
	f, err := AppendToFile(file)
	assert.NoError(t, err)

	_, err = f.WriteString("Test")
	assert.NoError(t, err)

	assert.NoError(t, f.Close())

	f, err = AppendToFile(file)
	assert.NoError(t, err)

	_, err = f.WriteString("Blub")
	assert.NoError(t, err)

	assert.NoError(t, f.Close())

	actual, err := os.ReadFile(file)
	assert.NoError(t, err)

	assert.Equal(t, "TestBlub", string(actual))
}
