package osutil

import (
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setAndGetSmallRLimit(t *testing.T) syscall.Rlimit {
	var limit syscall.Rlimit

	limit.Max = 10
	limit.Cur = 10

	assert.NoError(t, syscall.Setrlimit(syscall.RLIMIT_NOFILE, &limit))
	assert.NoError(t, syscall.Getrlimit(syscall.RLIMIT_NOFILE, &limit))

	return limit
}

func TestCapture(t *testing.T) {
	limit := setAndGetSmallRLimit(t)

	// Use at least one more file descriptor than our current limit, so we make sure that there are no file descriptor leaks.
	for i := 0; i <= int(limit.Cur); i++ {
		testCapture(t)
	}
}

func TestCaptureWithCGo(t *testing.T) {
	limit := setAndGetSmallRLimit(t)

	// Use at least one more file descriptor than our current limit, so we make sure that there are no file descriptor leaks.
	for i := 0; i <= int(limit.Cur); i++ {
		testCaptureWithCGo(t)
	}
}
