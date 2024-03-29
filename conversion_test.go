package osutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnySliceToTypeSlice(t *testing.T) {
	t.Run("String slice", func(t *testing.T) {
		type testCase struct {
			Name string

			AnySlice []any

			ExpectedTypeSlice []string
		}

		validate := func(t *testing.T, tc *testCase) {
			t.Run(tc.Name, func(t *testing.T) {
				actualTypeSlice := AnySliceToTypeSlice[string](tc.AnySlice)

				assert.Equal(t, tc.ExpectedTypeSlice, actualTypeSlice)
			})
		}

		validate(t, &testCase{
			Name: "Nil",

			AnySlice: nil,

			ExpectedTypeSlice: nil,
		})

		validate(t, &testCase{
			Name: "Convert any slice to string slice",

			AnySlice: []any{
				"foo",
				"bar",
			},

			ExpectedTypeSlice: []string{
				"foo",
				"bar",
			},
		})

		validate(t, &testCase{
			Name: "Only values that match the type",

			AnySlice: []any{
				"foo",
				"bar",
				12,
				145.66,
			},

			ExpectedTypeSlice: []string{
				"foo",
				"bar",
			},
		})
	})
	t.Run("Integer slice", func(t *testing.T) {
		type testCase struct {
			Name string

			AnySlice []any

			ExpectedTypeSlice []int
		}

		validate := func(t *testing.T, tc *testCase) {
			t.Run(tc.Name, func(t *testing.T) {
				actualTypeSlice := AnySliceToTypeSlice[int](tc.AnySlice)

				assert.Equal(t, tc.ExpectedTypeSlice, actualTypeSlice)
			})
		}

		validate(t, &testCase{
			Name: "Nil",

			AnySlice: nil,

			ExpectedTypeSlice: nil,
		})

		validate(t, &testCase{
			Name: "Convert any slice to int slice",

			AnySlice: []any{
				15,
				22,
			},

			ExpectedTypeSlice: []int{
				15,
				22,
			},
		})

		validate(t, &testCase{
			Name: "Only values that match the type",

			AnySlice: []any{
				"foo",
				"bar",
				12,
				"13",
				145.66,
			},

			ExpectedTypeSlice: []int{
				12,
			},
		})
	})
}
