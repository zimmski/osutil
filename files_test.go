package osutil

import (
	"io/ioutil"
	"os"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilesRecursive(t *testing.T) {
	path, err := ioutil.TempDir("", "os-util")
	assert.NoError(t, err)

	assert.NoError(t, os.MkdirAll(path+"/a/b", 0750))
	assert.NoError(t, ioutil.WriteFile(path+"/c.txt", []byte("foobar"), 0640))
	assert.NoError(t, ioutil.WriteFile(path+"/a/d.txt", []byte("foobar"), 0640))
	assert.NoError(t, ioutil.WriteFile(path+"/a/b/e.txt", []byte("foobar"), 0640))

	fs, err := FilesRecursive(path)
	assert.NoError(t, err)

	sort.Strings(fs)

	assert.Equal(
		t,
		[]string{
			path + "/a/b/e.txt",
			path + "/a/d.txt",
			path + "/c.txt",
		},
		fs,
	)
}
