package bytesutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitter(t *testing.T) {
	type testCase struct {
		Data     string
		Expected []string
	}

	validate := func(name string, tc testCase) {
		var got []string
		for line := range Split([]byte(tc.Data), '\n') {
			got = append(got, string(line))
		}

		assert.Equal(t, tc.Expected, got, "Split %q", tc.Data)
	}

	validate("some items", testCase{
		Data: "a\nb\nc",
		Expected: []string{
			"a",
			"b",
			"c",
		},
	})

	validate("empty item at the beginning", testCase{
		Data: "\nc",
		Expected: []string{
			"",
			"c",
		},
	})

	validate("empty item at the end", testCase{
		Data: "c\n",
		Expected: []string{
			"c",
			"",
		},
	})

	validate("empty string", testCase{
		Data: "",
		Expected: []string{
			"",
		},
	})
}
