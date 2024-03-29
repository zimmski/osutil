//go:build darwin

package osutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInfo(t *testing.T) {
	info, err := Info()
	assert.NoError(t, err)

	assert.NotEmpty(t, info[OperatingSystemVersionIdentifier])
	assert.NotEmpty(t, info[OperatingSystemIdentifier])
	assert.NotEmpty(t, info[OperatingSystemBuildIdentifier])
	assert.NotEmpty(t, info[KernelVersionIdentifier])
	assert.NotEmpty(t, info[EnvironmentPathIdentifier])
}
