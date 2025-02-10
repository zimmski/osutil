package bytesutil

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatToGoObject(t *testing.T) {
	type testCase struct {
		Name string

		Object any

		ExpectedString string
	}

	type structInlineable struct {
		A int
		B string
	}
	type structNotInlineable struct {
		A int
		B string
		C *structNotInlineable
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualString := FormatToGoObject(tc.Object)

			expectedString := StringTrimIndentations(tc.ExpectedString)
			expectedString = strings.TrimSpace(expectedString) // There is no final new-line, because the formatting function prints only what is needed.

			assert.Equal(t, expectedString, actualString)
			assert.Equal(t, []byte(expectedString), []byte(actualString))
		})
	}

	validate(t, &testCase{
		Name: "Anonymous struct",

		Object: struct {
			A int
			B string
		}{
			A: 123,
			B: "abc",
		},

		ExpectedString: `
			struct { A int; B string }{A:123, B:"abc"}
		`,
	})

	validate(t, &testCase{
		Name: "Inline-able struct",

		Object: structInlineable{
			A: 123,
			B: "abc",
		},

		ExpectedString: `
			bytesutil.structInlineable{A:123, B:"abc"}
		`,
	})

	validate(t, &testCase{
		Name: "Pointer to inline-able struct",

		Object: &structInlineable{
			A: 123,
			B: "abc",
		},

		ExpectedString: `
			&bytesutil.structInlineable{A:123, B:"abc"}
		`,
	})

	validate(t, &testCase{
		Name: "Not inline-able struct",

		Object: structNotInlineable{
			A: 123,
			B: "abc",
		},

		ExpectedString: `
			bytesutil.structNotInlineable{
			    A:  123,
			    B:  "abc",
			}
		`,
	})
}
