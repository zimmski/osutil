package osutil

import (
	"errors"
	"os"
)

// Exist errors
var (
	ErrNotADirectory = errors.New("not a directory")
)

// DirExists check if a directory exists
func DirExists(d string) error {
	f, err := os.Open(d)
	if err != nil {
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	fi, err := f.Stat()
	if err != nil {
		return err
	}

	if !fi.Mode().IsDir() {
		return ErrNotADirectory
	}

	return nil
}
