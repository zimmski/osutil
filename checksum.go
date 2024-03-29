package osutil

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"
)

// WriteChecksumForPath computes a checksum of a file or directory and writes it to the given file.
func WriteChecksumForPath(path string, checksumFile string) error {
	digest, err := ChecksumForPath(path)
	if err != nil {
		panic(err)
	}
	if err := os.WriteFile(checksumFile, []byte(fmt.Sprintf("%x\n", digest)), 0644); err != nil {
		return err
	}

	return nil
}

// ValidateChecksumForPath computes a checksum of a file or directory and returns an error if it does not match the checksum stored in the given file.
func ValidateChecksumForPath(path string, checksumFile string) (valid bool, err error) {
	contents, err := os.ReadFile(checksumFile)
	if err != nil {
		return false, err
	}

	digest, err := ChecksumForPath(path)
	if err != nil {
		return false, err
	}

	return fmt.Sprintf("%x\n", digest) == string(contents), nil
}

// ChecksumForPath computes a checksum of a file or directory.
func ChecksumForPath(path string) (digest []byte, err error) {
	hash := md5.New()

	if err := filepath.Walk(path, func(file string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		contents, err := os.ReadFile(file)
		if err != nil {
			// Even though "filepath.Walk" does not recurse into symlinks, it still invokes us with symlink itself. We can still include symlinked files in the checksum computation. Skip symlinks to directories and broken symlinks instead of failing to compute a checksum.
			if info.Mode()&os.ModeSymlink != 0 {
				if pe, ok := err.(*os.PathError); ok && (pe.Err == syscall.EISDIR || pe.Err == syscall.ENOENT) {
					return nil
				}
			}
			return err
		}

		relativePath, err := filepath.Rel(path, file)
		if err != nil {
			return err
		}

		for _, data := range [][]byte{
			[]byte(relativePath),
			[]byte{'\x00'},
			contents,
			[]byte{'\x00'},
		} {
			if _, err := hash.Write(data); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return hash.Sum(nil), nil
}

// ChecksumsSHA256ForFiles creates checksum-files with SHA-256 recursively for all files in a directory.
func ChecksumsSHA256ForFiles(filePath string) (err error) {
	return filepath.WalkDir(filePath, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Do not generate checksums for directories, but walk into them.
		if info.IsDir() {
			return nil
		}

		// Generate checksum of file.
		file, err := os.Open(path)
		defer func() {
			if e := file.Close(); e != nil {
				err = fmt.Errorf("error during closing of file: %v, %v", e, err)
			}
		}()
		if err != nil {
			return err
		}
		checksum := sha256.New()
		if _, err := io.Copy(checksum, file); err != nil {
			return err
		}

		// Write checksum to checksum-file.
		if err := os.WriteFile(path+".sha256", []byte(hex.EncodeToString(checksum.Sum(nil))), 0644); err != nil {
			return err
		}

		return nil
	})
}
