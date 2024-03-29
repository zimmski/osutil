package bytesutil

import (
	"bytes"
	"os"
)

// LineLengthsForFile returns a slice of the line lengths of the given data.
// Line endings (Unix+Windows) are ignored.
func LineLengthsForFile(filePath string) ([]uint, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return LineLengths(data), nil
}

// LineLengths returns a slice of the line lengths of the given data.
// Line endings (Unix+Windows) are ignored.
func LineLengths(data []byte) []uint {
	ls := []uint{0} // Line 0 is always filled so that we do not need to subtract one of every line to access it.
	for {
		i := bytes.IndexRune(data, '\n')
		if i == -1 {
			ls = append(ls, uint(len(data)))

			break
		}

		if i > 0 && data[i-1] == '\r' {
			ls = append(ls, uint(i)-1)
		} else {
			ls = append(ls, uint(i))
		}
		data = data[i+1:]
	}

	return ls
}

// PrefixLines prefixes every non-empty line with the given prefix.
func PrefixLines(data []byte, prefix []byte) (result []byte) {
	for l := range Split(data, '\n') {
		if len(result) > 0 {
			result = append(result, '\n')
		}
		if len(l) > 0 {
			result = append(result, prefix...)
			result = append(result, l...)
		}
	}

	return result
}
