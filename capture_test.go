package osutil

import (
	"fmt"
	"sync"
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

func TestCaptureWithCGoMutexBehavior(t *testing.T) {
	numRoutines := 50
	wg := &sync.WaitGroup{}
	wg.Add(numRoutines)
	results := make([]string, numRoutines)
	errors := make([]error, numRoutines)
	for j := 0; j < numRoutines; j++ {
		go func(jc int) {
			defer wg.Done()
			b, e := testCaptureWithCGoWrapper()
			results[jc] = string(b)
			errors[jc] = e
		}(j)
	}
	wg.Wait()
	for gi, err := range errors {
		assert.NoErrorf(t, err, "goroutine %d error", gi)
	}
	for gi, s := range results {
		assert.Equalf(t, s, "Go\nC\n", "goroutine %d", gi)
	}
}

func TestCaptureWithCGoWithPanic(t *testing.T) {
	assert.Panics(t, func() {
		_, _ = CaptureWithCGo(func() {
			fmt.Println("abc")

			panic("stop")
		})
	})
}
