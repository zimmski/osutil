package bytesutil

import (
	"sort"
	"strings"
)

// SortLines splits the given string into lines, sorts them and then returns the sorted lines as a combined string.
func SortLines(s string) (sorted string) {
	lines := strings.Split(s, "\n")

	sort.Strings(lines)

	return strings.Join(lines, "\n")
}

// SortLinesAndTrimSpace sorts the lines of the given string and removes all leading and trailing whitespaces.
func SortLinesAndTrimSpace(s string) (sorted string) {
	s = SortLines(s)
	s = strings.TrimSpace(s)

	return s
}
