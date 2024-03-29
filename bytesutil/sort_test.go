package bytesutil

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortLines(t *testing.T) {
	type testCase struct {
		In  string
		Out string
	}

	validate := func(name string, tc testCase) {
		assert.Equal(t, tc.Out, SortLines(tc.In))
	}

	validate("One Line", testCase{
		In:  "abc",
		Out: "abc",
	})
	validate("Unsorted Multiple Lines", testCase{
		In: strings.TrimSpace(StringTrimIndentations(`
			c
			b
			a
			d
		`)),
		Out: strings.TrimSpace(StringTrimIndentations(`
			a
			b
			c
			d
		`)),
	})
	validate("Sorted Multiple Lines", testCase{
		In: strings.TrimSpace(StringTrimIndentations(`
			a
			b
			c
			d
		`)),
		Out: strings.TrimSpace(StringTrimIndentations(`
			a
			b
			c
			d
		`)),
	})
}
