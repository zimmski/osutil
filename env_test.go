package osutil

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEnvEnabled(t *testing.T) {
	type testCase struct {
		Value         string
		EnabledValues []string

		Expected bool
	}
	envTestName := "_OSUTIL_TEST_ENVIRONMENT_VARIABLE"

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Value, func(t *testing.T) {
			os.Setenv(envTestName, tc.Value)
			defer os.Unsetenv(envTestName)

			actual := IsEnvEnabled(envTestName, tc.EnabledValues...)

			assert.Equal(t, tc.Expected, actual)
		})
	}

	validate(t, &testCase{
		Value:    "1",
		Expected: true,
	})
	validate(t, &testCase{
		Value:    "0",
		Expected: false,
	})
	validate(t, &testCase{
		Value:    "",
		Expected: false,
	})
	validate(t, &testCase{
		Value:    "yes",
		Expected: true,
	})
	validate(t, &testCase{
		Value:    "on",
		Expected: true,
	})
	validate(t, &testCase{
		Value:    "true",
		Expected: true,
	})
	t.Run("Other Cases", func(t *testing.T) {
		validate(t, &testCase{
			Value:    "True",
			Expected: true,
		})
		validate(t, &testCase{
			Value:    "TRUE",
			Expected: true,
		})
	})
	t.Run("Other Values", func(t *testing.T) {
		validate(t, &testCase{
			Value:         "positive",
			EnabledValues: []string{"positive"},
			Expected:      true,
		})
		validate(t, &testCase{
			Value:         "negative",
			EnabledValues: []string{"positive"},
			Expected:      false,
		})
	})
}
