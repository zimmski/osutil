package osutil

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/ulikunitz/xz"
)

// Tar archives a given path into a compressed TAR file.
func Tar(archiveFilePath string, path string) error {
	if DirExists(path) != nil {
		return errors.New("can only archive directories")
	}

	cmd := exec.Command("tar", "-czvf", archiveFilePath, "-C", path, ".")

	return cmd.Run()
}

// CompressionType defines the compression type.
type CompressionType string

const (
	// CompressionTypeNone indicates no compression.
	CompressionTypeNone = CompressionType("")
	// CompressionTypeGNUZipped indicates a GNU zipped compression.
	CompressionTypeGNUZipped = CompressionType("gz")
	// CompressionTypeXZ indicates a XZ compression.
	CompressionTypeXZ = CompressionType("xz")
)

// ExtractFile extracts a compressed file to a given path.
// How the archive is compressed and packed is automatically inferred, e.g. by the file extension.
func ExtractFile(archiveFilePath string, destinationPath string) (err error) {
	if strings.HasSuffix(archiveFilePath, ".tar.gz") || strings.HasSuffix(archiveFilePath, ".tar.xz") {
		return TarExtractFile(archiveFilePath, destinationPath)
	} else if strings.HasSuffix(archiveFilePath, ".zip") {
		return ZipExtractFile(archiveFilePath, destinationPath)
	}

	return fmt.Errorf("unknown compression for %s", archiveFilePath)
}

// TarExtractFile extracts a compressed TAR file to a given path.
func TarExtractFile(archiveFilePath string, destinationPath string) (err error) {
	f, err := os.Open(archiveFilePath)
	if err != nil {
		return err
	}
	defer func() {
		if e := f.Close(); e != nil {
			if err != nil {
				err = errors.Join(err, e)
			} else {
				err = e
			}
		}
	}()

	var compressionType CompressionType
	if strings.HasSuffix(archiveFilePath, ".xz") {
		compressionType = CompressionTypeXZ
	} else {
		compressionType = CompressionTypeGNUZipped
	}

	if err := TarExtract(f, destinationPath, compressionType); err != nil {
		return err
	}

	return nil
}

// TarExtract reads the gzip-compressed tar file from the reader and writes it into the destination path.
func TarExtract(stream io.Reader, destinationPath string, compressionType CompressionType) (err error) {
	// REMARK This code has been copied from https://cs.opensource.google/go/x/build/+/master:internal/untar/untar.go and then slighlty modified.

	now := time.Now()
	filesCopiedCount := 0
	directoriesCreated := map[string]bool{}

	createDirectoryForFile := func(filePathAbsolute string, fileMode fs.FileMode) (err error) {
		// Make the directory. This is redundant because it should already be made by a directory entry in the tar beforehand. Thus, don't check for errors; the next write will fail with the same error.
		fileDirectoryPath := filepath.Dir(filePathAbsolute)
		if !directoriesCreated[fileDirectoryPath] {
			if err := os.MkdirAll(filepath.Dir(filePathAbsolute), 0755); err != nil {
				return err
			}
			directoriesCreated[fileDirectoryPath] = true
		}
		if IsDarwin() && fileMode&0111 != 0 {
			// The darwin kernel caches binary signatures and SIGKILLs binaries with mismatched signatures. Overwriting a binary with O_TRUNC does not clear the cache, rendering the new copy unusable. Removing the original file first does clear the cache. See #54132.
			err := os.Remove(filePathAbsolute)
			if err != nil && !errors.Is(err, fs.ErrNotExist) {
				return err
			}
		}

		return nil
	}

	// Give a summary of the extraction.
	defer func() {
		td := time.Since(now)
		if err == nil {
			log.Printf("extracted tarball into %s: %d files, %d dirs (%v)", destinationPath, filesCopiedCount, len(directoriesCreated), td)
		} else {
			log.Printf("error extracting tarball into %s after %d files, %d dirs, %v: %v", destinationPath, filesCopiedCount, len(directoriesCreated), td, err)
		}
	}()

	switch compressionType {
	case CompressionTypeGNUZipped:
		stream, err = gzip.NewReader(stream)
		if err != nil {
			return fmt.Errorf("requires gzip-compressed body: %v", err)
		}
	case CompressionTypeXZ:
		stream, err = xz.NewReader(stream)
		if err != nil {
			return fmt.Errorf("requires xz-compressed body: %v", err)
		}
	}
	tarStream := tar.NewReader(stream)
	loggedChangeTimeError := false
	for {
		file, err := tarStream.Next()
		if err != nil {
			if err == io.EOF {
				break
			}

			log.Printf("tar reading error: %v", err)
			return fmt.Errorf("tar error: %v", err)
		}

		if !validRelPath(file.Name) {
			return fmt.Errorf("tar contained invalid name error %q", file.Name)
		}
		filePathRelative := filepath.FromSlash(file.Name)
		filePathAbsolute := filepath.Join(destinationPath, filePathRelative)

		fileInfo := file.FileInfo()
		fileMode := fileInfo.Mode()
		switch {
		case fileMode.IsRegular():
			if err := createDirectoryForFile(filePathAbsolute, fileMode); err != nil {
				return err
			}

			f, err := os.OpenFile(filePathAbsolute, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fileMode.Perm())
			if err != nil {
				return err
			}
			n, err := io.Copy(f, tarStream)
			if closeErr := f.Close(); closeErr != nil && err == nil {
				err = closeErr
			}
			if err != nil {
				return fmt.Errorf("error writing to %s: %v", filePathAbsolute, err)
			}
			if n != file.Size {
				return fmt.Errorf("only wrote %d bytes to %s; expected %d", n, filePathAbsolute, file.Size)
			}

			modTime := file.ModTime
			if modTime.After(now) {
				// Clamp modtimes at system time. See golang.org/issue/19062 when clock on buildlet was behind the gitmirror server doing the git-archive.
				modTime = now
			}
			if !modTime.IsZero() {
				if err := os.Chtimes(filePathAbsolute, modTime, modTime); err != nil && !loggedChangeTimeError {
					// Benign error. Gerrit doesn't even set the modtime in these, and we don't end up relying on it anywhere (the gomote push command relies on digests only), so this is a little pointless for now.
					log.Printf("error changing modtime: %v (further Chtimes errors suppressed)", err)
					loggedChangeTimeError = true // once is enough
				}
			}

			filesCopiedCount++
		case fileMode.IsDir():
			if err := os.MkdirAll(filePathAbsolute, 0755); err != nil {
				return err
			}

			directoriesCreated[filePathAbsolute] = true
		case file.Typeflag == tar.TypeLink:
			if err := createDirectoryForFile(filePathAbsolute, fileMode); err != nil {
				return err
			}

			if err := os.Link(file.Linkname, filePathAbsolute); err != nil {
				return fmt.Errorf("failed writing symbolic link: %s", err)
			}
		case file.Typeflag == tar.TypeSymlink:
			if err := createDirectoryForFile(filePathAbsolute, fileMode); err != nil {
				return err
			}

			if err := os.Symlink(file.Linkname, filePathAbsolute); err != nil {
				return fmt.Errorf("failed writing symbolic link: %s", err)
			}
		default:
			return fmt.Errorf("tar file entry %s contained unsupported file type %v (%v)", file.Name, fileMode, file.Typeflag)
		}
	}

	return nil
}

// ZipExtractFile extracts a zipped file to a given path.
func ZipExtractFile(archiveFilePath string, destinationPath string) (err error) {
	archive, err := zip.OpenReader(archiveFilePath)
	if err != nil {
		return err
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(destinationPath, f.Name)
		if !strings.HasPrefix(filePath, filepath.Clean(destinationPath)+string(os.PathSeparator)) {
			log.Println("invalid file path")

			return
		}
		if f.FileInfo().IsDir() {
			os.MkdirAll(filePath, os.ModePerm)

			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			return err
		}

		destinationFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		fileInArchive, err := f.Open()
		if err != nil {
			return err
		}

		if _, err := io.Copy(destinationFile, fileInArchive); err != nil {
			return err
		}

		destinationFile.Close()
		fileInArchive.Close()
	}

	return nil
}

func validRelPath(p string) bool {
	if p == "" || strings.Contains(p, `\`) || strings.HasPrefix(p, "/") || strings.Contains(p, "../") {
		return false
	}
	return true
}
