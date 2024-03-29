package bytesutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLineLengthsForFile(t *testing.T) {
	assert.Equal(
		t,
		[]uint{0, 1, 2, 1},
		LineLengths([]byte("a\naa\r\na")),
	)
}

func TestPrefixLines(t *testing.T) {
	type testCase struct {
		Name string

		Data   []byte
		Prefix []byte

		ExpectedResult []byte
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			tc.Data = TrimIndentations(tc.Data)
			tc.ExpectedResult = TrimIndentations(tc.ExpectedResult)

			actualResult := PrefixLines(tc.Data, tc.Prefix)

			assert.Equal(t, tc.ExpectedResult, actualResult)
		})
	}

	validate(t, &testCase{
		Name: "Prefix lines",

		Data: []byte(`
			a
			b
			c
		`),
		Prefix: []byte("- "),

		ExpectedResult: []byte(`
			- a
			- b
			- c
		`),
	})

	validate(t, &testCase{
		Name: "Prefix with empty lines",

		Data: []byte(`
			a

			b

			c
		`),
		Prefix: []byte("- "),

		ExpectedResult: []byte(`
			- a

			- b

			- c
		`),
	})
}
