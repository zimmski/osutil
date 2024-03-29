package osutil

import (
	"os"
	"strings"
)

// GreatestCommonDirectory computes the greatest common part of the given paths.
// The resulting string must be the prefix of all the given paths.
func GreatestCommonDirectory(paths []string) string {
	if len(paths) == 0 {
		return ""
	}

	greatestCommonDirectory := paths[0]
	for i := 1; i < len(paths); i++ {
		current := paths[i]

		minLength := len(greatestCommonDirectory)
		if len(current) < minLength {
			minLength = len(current)
		}

		currentCommonIndex := 0
		for i := 0; i < minLength; i++ {
			if current[i] != greatestCommonDirectory[i] {
				break
			}

			if current[i] == os.PathSeparator {
				currentCommonIndex = i
			} else if i == minLength-1 {
				currentCommonIndex = minLength
			}
		}

		greatestCommonDirectory = greatestCommonDirectory[:currentCommonIndex]
	}

	return greatestCommonDirectory
}

// EnvironmentPathList returns the list of file paths contained in the "PATH" environment variable.
func EnvironmentPathList() (filePaths []string) {
	path := os.Getenv(EnvironmentPathIdentifier)

	return strings.Split(path, string(os.PathListSeparator))
}

// RemoveFromEnvironmentPathBySearchTerm returns the content of the "PATH" environment variable where file paths containing the given search terms are removed.
func RemoveFromEnvironmentPathBySearchTerm(searchTerms ...string) (newEnvironmentPath string) {
	filePathsOld := EnvironmentPathList()
	var filePathsNew []string
PATH:
	for _, p := range filePathsOld {
		for _, term := range searchTerms {
			if strings.Contains(p, term) {
				continue PATH
			}
		}

		filePathsNew = append(filePathsNew, p)
	}

	return strings.Join(filePathsNew, string(os.PathListSeparator))
}
