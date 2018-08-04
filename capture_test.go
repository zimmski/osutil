package osutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCapture(t *testing.T) {
	assert.NoError(t, SetRLimitFiles(10, func(limit uint64) {
		// Use at least one more file descriptor than our current limit, so we make sure that there are no file descriptor leaks.
		for i := 0; i <= int(limit); i++ {
			testCapture(t)
		}
	}))
}

func TestCaptureWithCGo(t *testing.T) {
	assert.NoError(t, SetRLimitFiles(10, func(limit uint64) {
		// Use at least one more file descriptor than our current limit, so we make sure that there are no file descriptor leaks.
		for i := 0; i <= int(limit); i++ {
			testCaptureWithCGo(t)
		}
	}))
}
