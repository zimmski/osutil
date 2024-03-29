package osutil

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopyFile(t *testing.T) {
	src := "copy.go"
	dst := "copy.go.tmp"

	assert.NoError(t, CopyFile(src, dst))

	s, err := os.ReadFile(src)
	assert.NoError(t, err)

	d, err := os.ReadFile(dst)
	assert.NoError(t, err)

	assert.Equal(t, s, d)

	assert.NoError(t, os.Remove(dst))
}
