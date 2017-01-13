package osutil

import (
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCapture(t *testing.T) {
	var limit syscall.Rlimit
	assert.NoError(t, syscall.Getrlimit(syscall.RLIMIT_NOFILE, &limit))

	// Use at least one more file descriptor than our current limit, so we make sure that there are no file descriptor leaks.
	for i := 0; i <= int(limit.Cur); i++ {
		testCapture(t)
	}
}

func TestCaptureWithCGo(t *testing.T) {
	var limit syscall.Rlimit
	assert.NoError(t, syscall.Getrlimit(syscall.RLIMIT_NOFILE, &limit))

	// Use at least one more file descriptor than our current limit, so we make sure that there are no file descriptor leaks.
	for i := 0; i <= int(limit.Cur); i++ {
		testCaptureWithCGo(t)
	}
}
