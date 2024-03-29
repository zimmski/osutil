//go:build linux

package osutil_test

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnforceProcessTreeLimitsMemory(t *testing.T) {
	type testCase struct {
		Name string

		MemoryLimitInMiB         uint
		MinMemoryToAllocateInMiB uint

		ValidateOutput func(t *testing.T, programErr error, minMemoryAllocatedInMiB uint)
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			assert.NoError(t, exec.Command("go", "build", "-o", "limittest/memory", "limittest/memory.go").Run())

			cmd := exec.Command("limittest/memory", strconv.Itoa(int(tc.MemoryLimitInMiB)), strconv.Itoa(int(tc.MinMemoryToAllocateInMiB)))
			out, programmErr := cmd.CombinedOutput()

			lines := bytes.Split(out, []byte("\n"))
			assert.True(t, len(lines) > 1)
			mem, err := strconv.Atoi(string(lines[len(lines)-2]))
			assert.NoError(t, err)

			tc.ValidateOutput(t, programmErr, uint(mem)) // REMARK The number of loop iterations in the test program is not an accurate metric for how much memory was used but it's an alright sanity check for low memory usage numbers.
		})
	}

	validate(t, &testCase{
		Name: "Limit hit",

		MemoryLimitInMiB:         100,
		MinMemoryToAllocateInMiB: 1000,

		ValidateOutput: func(t *testing.T, programErr error, mem uint) {
			if err, ok := programErr.(*exec.ExitError); ok {
				assert.Equal(t, 5, err.ExitCode())
			} else {
				assert.Fail(t, "Error was not an ExitError: %v", err)
			}

			err := os.Remove("success.txt")
			if err == nil {
				assert.Fail(t, "Expected success.txt to not exist but it existed")
			}
			assert.True(t, errors.Is(err, fs.ErrNotExist), "Unexpected error from removing success.txt: %v", err)

			assert.False(t, mem > 300, "Expected memory usage to stay below 300MiB but was at least %dMiB", mem)
		},
	})

	validate(t, &testCase{
		Name: "No limit hit",

		MemoryLimitInMiB:         1000,
		MinMemoryToAllocateInMiB: 100,

		ValidateOutput: func(t *testing.T, programErr error, mem uint) {
			assert.NoError(t, programErr)

			assert.NoError(t, os.Remove("success.txt"))

			assert.Equal(t, uint(100), mem)
		},
	})
}
