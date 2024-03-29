package osutil

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
)

// CanonicalizeAndEvaluateSymlinks returns the path after canonicalizing it and the evaluation of any symbolic links.
func CanonicalizeAndEvaluateSymlinks(path string) (resolvedPath string, err error) {
	path, err = filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return filepath.EvalSymlinks(path)
}

// DirectoriesRecursive returns all subdirectories of the given path including the given path.
func DirectoriesRecursive(directoryPath string) (directories []string, err error) {
	if err := filepath.WalkDir(directoryPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			directories = append(directories, path)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return directories, nil
}

// FilesRecursive returns all files in a given path and its subpaths.
func FilesRecursive(path string) (files []string, err error) {
	var fs []string

	err = filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if f.IsDir() {
			return nil
		}

		fs = append(fs, path)

		return err
	})
	if err != nil {
		return nil, err
	}

	return fs, nil
}

// ForEachFile walks through the given path and calls the given callback with every file.
func ForEachFile(path string, handle func(filePath string) error) error {
	return filepath.WalkDir(path, func(filePath string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		return handle(filePath)
	})
}

// FilePathsByHierarchy sorts file paths by their hierarchy.
type FilePathsByHierarchy []string

// Len is the number of elements in the collection.
func (s FilePathsByHierarchy) Len() int {
	return len(s)
}

// Swap swaps the elements with indexes i and j.
func (s FilePathsByHierarchy) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less reports whether the element with index i must sort before the element with index j.
func (s FilePathsByHierarchy) Less(i, j int) bool {
	si := strings.Split(s[i], string(os.PathSeparator))
	sj := strings.Split(s[j], string(os.PathSeparator))

	if len(si) != len(sj) {
		return len(si) < len(sj)
	}

	for i, sie := range si {
		if c := sie < sj[i]; c {
			return c
		}
	}

	return false
}

// AppendToFile opens the named file. If the file does not exist it is created.
func AppendToFile(name string) (*os.File, error) {
	return os.OpenFile(name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
}

// FileChange changes the content of a file.
func FileChange(filePath string, change func(data []byte) (changed []byte, err error)) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	data, err = change(data)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filePath, data, 0000); err != nil { // The permission is invalid and will be ignored, as we know the specified file already exists.
		return err
	}

	return nil
}

// ReplaceVariablesInFile replaces all variables in a file.
// A variable in a file has the syntax `{{$key}}` and which is then replaced by its value.
func ReplaceVariablesInFile(filePath string, variables map[string]string) (err error) {
	d, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	for k, v := range variables {
		d = bytes.ReplaceAll(d, []byte(k), []byte(v))
	}

	return os.WriteFile(filePath, d, 0)
}
