package bytesutil

import (
	"bytes"
	"fmt"
	"os"

	pkgerrors "github.com/pkg/errors"
)

// ReplaceBytesInFile replaces a certain amount of occurrences of the given bytes in a file.
// A negative number of occurrences replaces all matches. If no output file is given, the input file is overwritten.
func ReplaceBytesInFile(filePathIn string, filePathOut string, search []byte, replace []byte, n int) (err error) {
	if filePathOut == "" {
		filePathOut = filePathIn
	}

	content, err := os.ReadFile(filePathIn)
	if err != nil {
		return pkgerrors.Wrap(err, filePathIn)
	}
	fileMode, err := os.Stat(filePathIn)
	if err != nil {
		return pkgerrors.Wrap(err, filePathIn)
	}

	contentReplaced := bytes.Replace(content, search, replace, n)

	if err := os.WriteFile(filePathOut, contentReplaced, fileMode.Mode()); err != nil {
		return pkgerrors.Wrap(err, filePathOut)
	}

	return nil
}

// ReplaceBytesInBinary replaces a certain amount of occurrences of the given bytes in a binary.
// A negative number of occurrences replaces all matches. If no output file is given, the input file is overwritten.
func ReplaceBytesInBinary(binaryPathIn string, binaryPathOut string, search []byte, replace []byte, n int) (err error) {
	if len(search) != len(replace) {
		return fmt.Errorf("can only replace byte sequences of equal length in a binary (was %d != %d)", len(search), len(replace))
	}

	return ReplaceBytesInFile(binaryPathIn, binaryPathOut, search, replace, n)
}
