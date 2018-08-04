package osutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirExists(t *testing.T) {
	assert.NoError(t, DirExists("/"))
	assert.NoError(t, DirExists("../osutil"))

	assert.Error(t, DirExists("hey"))

	assert.Equal(t, ErrNotADirectory, DirExists("exists.go"))
}
