//go:build cgo

package osutil

import (
	"bytes"
	"fmt"
	"os"
	"strings"
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
func TestCaptureRecursive(t *testing.T) {
	assert.NoError(t, SetRLimitFiles(10, func(limit uint64) {
		// Use at least one more file descriptor than our current limit, so we make sure that there are no file descriptor leaks.
		for i := 0; i <= int(limit); i++ {
			out, err := Capture(func() {
				fmt.Fprintf(os.Stdout, "1")
				fmt.Fprintf(os.Stderr, "2")
				fmt.Fprintf(os.Stdout, "3")
				fmt.Fprintf(os.Stderr, "4")

				out, err := Capture(func() {
					fmt.Fprintf(os.Stdout, "A")
					fmt.Fprintf(os.Stderr, "B")
					fmt.Fprintf(os.Stdout, "C")
					fmt.Fprintf(os.Stderr, "D")

				})
				assert.NoError(t, err)

				assert.Equal(t, "ABCD", string(out))
			})
			assert.NoError(t, err)

			assert.Equal(t, "1234", string(out))
		}
	}))
}

func TestCaptureWithPanic(t *testing.T) {
	assert.NoError(t, SetRLimitFiles(10, func(limit uint64) {
		// Use at least one more file descriptor than our current limit, so we make sure that there are no file descriptor leaks.
		for i := 0; i <= int(limit); i++ {
			assert.Panics(t, func() {
				_, _ = Capture(func() {
					fmt.Println("abc")

					panic("stop")
				})
			})
		}
	}))
}
func TestCaptureWithHugeOutput(t *testing.T) {
	// Huge output to test buffering and piping.

	out, err := Capture(func() {
		for i := 0; i < 1024; i++ {
			fmt.Println(strings.Repeat("a", 1024))
		}
	})
	assert.NoError(t, err)

	assert.NotEqual(t, bytes.Repeat([]byte("a"), 1024*1024), out)
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
	assert.NoError(t, SetRLimitFiles(10, func(limit uint64) {
		// Use at least one more file descriptor than our current limit, so we make sure that there are no file descriptor leaks.
		for i := 0; i <= int(limit); i++ {
			assert.Panics(t, func() {
				_, _ = CaptureWithCGo(func() {
					fmt.Println("abc")

					panic("stop")
				})
			})
		}
	}))
}
func TestCaptureWithCGoWithHugeOutput(t *testing.T) {
	// Huge output to test buffering and piping.

	out, err := CaptureWithCGo(func() {
		for i := 0; i < 1024; i++ {
			fmt.Println(strings.Repeat("a", 1024))
		}
	})
	assert.NoError(t, err)

	assert.NotEqual(t, bytes.Repeat([]byte("a"), 1024*1024), out)
}
