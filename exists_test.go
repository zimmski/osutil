package osutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDirExists(t *testing.T) {
	assert.Nil(t, DirExists("/"))
	assert.Nil(t, DirExists("../osutil"))

	assert.NotNil(t, DirExists("hey"))

	assert.Equal(t, ErrNotADirectory, DirExists("exists.go"))
}
