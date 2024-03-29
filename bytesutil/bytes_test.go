package bytesutil

import (
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zimmski/osutil"
)

func TestReplaceBytesInBinary(t *testing.T) {
	type testCase struct {
		Name string

		Binary  string
		Search  string
		Replace string
		N       int

		ExpectedErr    string
		ExpectedOutput string
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			// REMARK Mimic "symflowertesting.TemporaryDirectory" since it cannot be used directly because of a cyclic import.
			tempDir := filepath.Join(os.TempDir(), strconv.Itoa(int(rand.Int31())))
			require.NoError(t, os.Mkdir(tempDir, 0750))
			defer osutil.RemoveTemporaryDirectory(tempDir)

			sourceFilePath := filepath.Join(tempDir, "main.go")
			executableFilePath := filepath.Join(tempDir, "main") + osutil.BinaryExtension()

			require.NoError(t, os.WriteFile(sourceFilePath, []byte(StringTrimIndentations(tc.Binary)), 0640))

			compileCommand := exec.Command("go"+osutil.BinaryExtension(), "build", "-o", executableFilePath, sourceFilePath)
			require.NoError(t, compileCommand.Run(), "cannot compile Go source")

			err := ReplaceBytesInBinary(executableFilePath, "", []byte(tc.Search), []byte(tc.Replace), tc.N)
			if tc.ExpectedErr != "" {
				assert.EqualError(t, err, tc.ExpectedErr)

				return
			}
			assert.NoError(t, err)

			var output strings.Builder
			mainCommand := &exec.Cmd{
				Path:   executableFilePath,
				Stdout: &output,
			}
			err = mainCommand.Run()
			assert.NoError(t, err)
			assert.Equal(t, tc.ExpectedOutput, strings.TrimSpace(output.String()))
		})
	}

	validate(t, &testCase{
		Name: "No Modifications",

		Binary: `
			package main

			import "fmt"

			func main() {
					fmt.Println("Hello world!")
			}
		`,
		Search:  "foo",
		Replace: "bar",

		ExpectedOutput: "Hello world!",
	})
	t.Run("Replace", func(t *testing.T) {
		binaryTwoConstants := `
			package main

			import "fmt"

			const constant1 = "foo"
			const constant2 = "foo"

			func main() {
				fmt.Println("Hello " + constant1 + constant2)
			}
		`
		validate(t, &testCase{
			Name: "One Constant",

			Binary: binaryTwoConstants,

			Search:  "foo",
			Replace: "bar",
			N:       1,

			ExpectedOutput: "Hello barfoo",
		})
		validate(t, &testCase{
			Name: "Two Constants",

			Binary: binaryTwoConstants,

			Search:  "foo",
			Replace: "bar",
			N:       2,

			ExpectedOutput: "Hello barbar",
		})
	})
	validate(t, &testCase{
		Name: "Invalid Length",

		Binary: `
			package main

			func main() {}
		`,
		Search:  "1234",
		Replace: "123",

		ExpectedErr: "can only replace byte sequences of equal length in a binary (was 4 != 3)",
	})
}
