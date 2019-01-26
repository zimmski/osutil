package osutil

import (
	"fmt"
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

func TestCaptureWithPanic(t *testing.T) {
	assert.Panics(t, func() {
		_, _ = Capture(func() {
			fmt.Println("abc")

			panic("stop")
		})
	})
}

func TestCaptureWithCGo(t *testing.T) {
	assert.NoError(t, SetRLimitFiles(10, func(limit uint64) {
		// Use at least one more file descriptor than our current limit, so we make sure that there are no file descriptor leaks.
		for i := 0; i <= int(limit); i++ {
			testCaptureWithCGo(t)
		}
	}))
}

func TestCaptureWithCGoWithPanic(t *testing.T) {
	assert.Panics(t, func() {
		_, _ = CaptureWithCGo(func() {
			fmt.Println("abc")

			panic("stop")
		})
	})
}
