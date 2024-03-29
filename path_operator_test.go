package osutil

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGreatestCommonDirectory(t *testing.T) {
	type testCase struct {
		Name string

		Paths []string

		ExpectedString string
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			for i, p := range tc.Paths {
				tc.Paths[i] = filepath.FromSlash(p)
			}
			tc.ExpectedString = filepath.FromSlash(tc.ExpectedString)

			actualString := GreatestCommonDirectory(tc.Paths)

			assert.Equal(t, tc.ExpectedString, actualString)
		})
	}

	validate(t, &testCase{
		Name: "Equal paths",

		Paths: []string{
			"pkg/a",
			"pkg/a",
		},

		ExpectedString: "pkg/a",
	})
	validate(t, &testCase{
		Name: "Equal paths with trailing slash for one path",

		Paths: []string{
			"pkg/a",
			"pkg/a/",
		},

		ExpectedString: "pkg/a",
	})
	validate(t, &testCase{
		Name: "Equal paths with trailing slash for both paths",

		Paths: []string{
			"pkg/a/",
			"pkg/a/",
		},

		ExpectedString: "pkg/a",
	})
	validate(t, &testCase{
		Name: "Simple with common part",

		Paths: []string{
			"pkg/a",
			"pkg/b",
		},

		ExpectedString: "pkg",
	})
	validate(t, &testCase{
		Name: "Simple without common start",

		Paths: []string{
			"pkg/a",
			"other/a",
		},

		ExpectedString: "",
	})
	validate(t, &testCase{
		Name: "Different number of path components",

		Paths: []string{
			"pkg/a/b",
			"pkg/c",
		},

		ExpectedString: "pkg",
	})
	validate(t, &testCase{
		Name: "Multiple common parts",

		Paths: []string{
			"same/pkg/until/now",
			"same/pkg/until/then",
		},

		ExpectedString: "same/pkg/until",
	})
}
