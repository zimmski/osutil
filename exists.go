package osutil

import (
	"errors"
	"io/fs"
	"os"
)

var (
	// ErrNotADirectory indicates that the given directory does not exist.
	ErrNotADirectory = errors.New("not a directory")
	// ErrNotAFile indicates thate the given file does not exist.
	ErrNotAFile = errors.New("not a file")
)

// Stat retuns a FileInfo structure describing the given file.
func Stat(filePath string) (os.FileInfo, error) {
	return os.Stat(filePath)
}

// DirExists checks if a directory exists.
func DirExists(filePath string) error {
	fi, err := Stat(filePath)
	if err != nil {
		return err
	}

	if !fi.Mode().IsDir() {
		return ErrNotADirectory
	}

	return nil
}

// FileExists checks if a file exists while following symlinks.
func FileExists(filePath string) error {
	fi, err := Stat(filePath)
	if err != nil {
		return err
	}

	if fi.Mode().IsDir() {
		return ErrNotAFile
	}

	return nil
}

// RemoveFileIfExists checks if a file exists, and removes the file if it does exist.
// Symlinks are not followed, since they are files and should be removable by this function.
func RemoveFileIfExists(filePath string) error {
	if _, err := os.Lstat(filePath); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}

		return err
	}

	return os.Remove(filePath)
}

// FileOrSymlinkExists checks if a file exists while not following symlinks.
func FileOrSymlinkExists(filepath string) error {
	fi, err := os.Lstat(filepath)
	if err != nil {
		return err
	}

	if fi.Mode().IsDir() {
		return ErrNotAFile
	}

	return nil
}
